package pipeline

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/renderer"
)

// TemplateFS must be set by the caller to the embedded template filesystem.
var TemplateFS embed.FS

// FileResult describes what happened to a single output file.
type FileResult struct {
	Path   string
	Status FileStatus
	Err    error
}

// FileStatus indicates the outcome for a generated file.
type FileStatus int

const (
	FileCreated   FileStatus = iota
	FileUpdated
	FileUnchanged
	FileOrphaned
	FileError
)

func (s FileStatus) String() string {
	switch s {
	case FileCreated:
		return "created"
	case FileUpdated:
		return "updated"
	case FileUnchanged:
		return "unchanged"
	case FileOrphaned:
		return "orphaned"
	case FileError:
		return "error"
	default:
		return "unknown"
	}
}

// PlannedFile represents a template + context + output path.
type PlannedFile struct {
	TemplatePath string
	OutputPath   string
	Data         any
	IsGo         bool
}

// Pipeline orchestrates the code generation process.
type Pipeline struct {
	config    *config.Project
	outputDir string
	clean     bool
	renderer  *renderer.Renderer
}

// New creates a new Pipeline.
func New(cfg *config.Project, outputDir string, clean bool) *Pipeline {
	return &Pipeline{
		config:    cfg,
		outputDir: outputDir,
		clean:     clean,
		renderer:  renderer.New(TemplateFS),
	}
}

// Summary holds the generation results.
type Summary struct {
	Created   int
	Updated   int
	Unchanged int
	Orphaned  int
	Errors    int
	Results   []FileResult
}

// Run executes the full generation pipeline.
func (p *Pipeline) Run(ctx context.Context) error {
	// 1. Plan all output files.
	planned := p.planFiles()

	fmt.Printf("Planning %d files for project %q\n", len(planned), p.config.ProjectMeta.Name)

	// 2. Render and write all files.
	summary := p.renderAndWrite(ctx, planned)

	// 3. Detect orphaned generated files.
	orphaned := p.findOrphans(planned)
	summary.Orphaned = len(orphaned)

	if p.clean {
		for _, path := range orphaned {
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(os.Stderr, "warning: could not remove orphan %s: %v\n", path, err)
			}
			summary.Results = append(summary.Results, FileResult{Path: path, Status: FileOrphaned})
		}
	} else {
		for _, path := range orphaned {
			fmt.Fprintf(os.Stderr, "warning: orphaned generated file: %s\n", path)
			summary.Results = append(summary.Results, FileResult{Path: path, Status: FileOrphaned})
		}
	}

	// 4. Print summary.
	fmt.Printf("\nGeneration complete: %d created, %d updated, %d unchanged",
		summary.Created, summary.Updated, summary.Unchanged)
	if summary.Orphaned > 0 {
		fmt.Printf(", %d orphaned", summary.Orphaned)
	}
	if summary.Errors > 0 {
		fmt.Printf(", %d errors", summary.Errors)
	}
	fmt.Println()

	if summary.Errors > 0 {
		return fmt.Errorf("%d files had errors during generation", summary.Errors)
	}

	return nil
}

func (p *Pipeline) renderAndWrite(ctx context.Context, planned []PlannedFile) *Summary {
	var (
		summary  Summary
		mu       sync.Mutex
		created  atomic.Int32
		updated  atomic.Int32
		unchanged atomic.Int32
		errors   atomic.Int32
		wg       sync.WaitGroup
	)

	// Use a semaphore to limit concurrency.
	sem := make(chan struct{}, 8)

	for _, pf := range planned {
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		go func(pf PlannedFile) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := p.renderOne(pf)

			switch result.Status {
			case FileCreated:
				created.Add(1)
			case FileUpdated:
				updated.Add(1)
			case FileUnchanged:
				unchanged.Add(1)
			case FileError:
				errors.Add(1)
			}

			mu.Lock()
			summary.Results = append(summary.Results, result)
			mu.Unlock()
		}(pf)
	}

	wg.Wait()

	summary.Created = int(created.Load())
	summary.Updated = int(updated.Load())
	summary.Unchanged = int(unchanged.Load())
	summary.Errors = int(errors.Load())

	return &summary
}

func (p *Pipeline) renderOne(pf PlannedFile) FileResult {
	outPath := filepath.Join(p.outputDir, pf.OutputPath)

	// Render the template.
	content, err := p.renderer.Render(pf.TemplatePath, pf.Data)
	if err != nil {
		return FileResult{Path: outPath, Status: FileError, Err: fmt.Errorf("rendering %s: %w", pf.TemplatePath, err)}
	}

	output := []byte(content)

	// Format Go files.
	if pf.IsGo {
		formatted, fmtErr := renderer.FormatGo(output)
		if fmtErr != nil {
			return FileResult{Path: outPath, Status: FileError, Err: fmt.Errorf("formatting %s: %w", pf.OutputPath, fmtErr)}
		}
		output = formatted
	}

	// Diff-aware write.
	status, err := diffAwareWrite(outPath, output)
	if err != nil {
		return FileResult{Path: outPath, Status: FileError, Err: err}
	}

	return FileResult{Path: outPath, Status: status}
}

// diffAwareWrite writes content to path only if it differs from existing content.
func diffAwareWrite(path string, content []byte) (FileStatus, error) {
	// Check if file exists and has same content.
	existing, err := os.ReadFile(path)
	if err == nil {
		if string(existing) == string(content) {
			return FileUnchanged, nil
		}
		// File exists but content differs.
		if err = os.WriteFile(path, content, 0o644); err != nil {
			return FileError, fmt.Errorf("writing %s: %w", path, err)
		}
		return FileUpdated, nil
	}

	// File doesn't exist, create it.
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return FileError, fmt.Errorf("creating directory %s: %w", dir, err)
	}

	if err = os.WriteFile(path, content, 0o644); err != nil {
		return FileError, fmt.Errorf("writing %s: %w", path, err)
	}

	return FileCreated, nil
}

// findOrphans finds *.generated.* files in the output directory that are not in the plan.
func (p *Pipeline) findOrphans(planned []PlannedFile) []string {
	// Build set of planned output paths.
	plannedPaths := make(map[string]bool)
	for _, pf := range planned {
		absPath := filepath.Join(p.outputDir, pf.OutputPath)
		plannedPaths[absPath] = true
	}

	var orphans []string

	_ = filepath.Walk(p.outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// Only consider files with .generated. in the name.
		base := filepath.Base(path)
		if !strings.Contains(base, ".generated.") {
			return nil
		}

		if !plannedPaths[path] {
			orphans = append(orphans, path)
		}

		return nil
	})

	return orphans
}

// planFiles builds the complete list of files to generate.
// This is the central registry that maps templates to output paths.
func (p *Pipeline) planFiles() []PlannedFile {
	var files []PlannedFile

	for _, domain := range p.config.Domains {
		for _, entity := range domain.Entities {
			ctx := TemplateContext{
				Project: p.config,
				Domain:  &domain,
				Entity:  &entity,
			}

			files = append(files,
				p.planDomainFiles(ctx)...,
			)
		}
	}

	return files
}

// TemplateContext is the data passed to every template.
type TemplateContext struct {
	Project *config.Project
	Domain  *config.Domain
	Entity  *config.Entity
}

func (p *Pipeline) planDomainFiles(ctx TemplateContext) []PlannedFile {
	domain := ctx.Domain.Name
	entitySnake := strings.ToLower(ctx.Entity.Name)

	var files []PlannedFile

	// Domain layer
	templateMappings := []struct {
		template string
		output   string
		isGo     bool
	}{
		{"templates/domain/entity.go.tmpl", fmt.Sprintf("internal/domain/%s/%s.generated.go", domain, entitySnake), true},
		{"templates/domain/repository.go.tmpl", fmt.Sprintf("internal/domain/%s/repository.generated.go", domain), true},
		{"templates/domain_keys/keys.go.tmpl", fmt.Sprintf("internal/domain/%s/keys/keys.generated.go", domain), true},
		{"templates/domain_fakes/fake.go.tmpl", fmt.Sprintf("internal/domain/%s/fakes/fake.generated.go", domain), true},
		{"templates/domain_converters/converters.go.tmpl", fmt.Sprintf("internal/domain/%s/converters/%s.generated.go", domain, entitySnake), true},
		{"templates/domain_mock/repository.go.tmpl", fmt.Sprintf("internal/domain/%s/mock/repository.generated.go", domain), true},

		// Manager layer
		{"templates/manager/interface.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/interface.generated.go", domain), true},
		{"templates/manager/manager.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/manager.generated.go", domain), true},
		{"templates/manager/do.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/do.generated.go", domain), true},
		{"templates/manager/mock_manager.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/mock/manager.generated.go", domain), true},

		// Repository layer
		{"templates/repository/client.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/client.generated.go", domain), true},
		{"templates/repository/entity.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/%s.generated.go", domain, entitySnake), true},
		{"templates/repository/do.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/do.generated.go", domain), true},

		// Query codegen
		{"templates/codegen/queries.go.tmpl", fmt.Sprintf("cmd/tools/codegen/queries/%s_%s.generated.go", domain, entitySnake), true},

		// Migrations
		{"templates/migrations/migration.sql.tmpl", fmt.Sprintf("internal/repositories/postgres/migrations/migration_files/%s.generated.sql", domain), false},

		// Proto
		{"templates/proto/messages.proto.tmpl", fmt.Sprintf("proto/%s/%s_messages.generated.proto", domain, entitySnake), false},
		{"templates/proto/service.proto.tmpl", fmt.Sprintf("proto/%s/%s_service.generated.proto", domain, entitySnake), false},

		// gRPC service
		{"templates/grpc/service.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/service.generated.go", domain), true},
		{"templates/grpc/entity.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/%s.generated.go", domain, entitySnake), true},
		{"templates/grpc/converters.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/converters/converters.generated.go", domain), true},
		{"templates/grpc/permissions.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/permissions.generated.go", domain), true},
		{"templates/grpc/do.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/do.generated.go", domain), true},

		// Authorization
		{"templates/authorization/permissions.go.tmpl", fmt.Sprintf("internal/authorization/%s_permissions.generated.go", domain), true},
	}

	for _, m := range templateMappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         m.isGo,
		})
	}

	return files
}
