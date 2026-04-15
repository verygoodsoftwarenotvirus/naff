package pipeline

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
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
	Mode         os.FileMode // 0 means default (0o644); use 0o755 for scripts.
	// Raw bypasses text/template parsing and Go formatting. The file at
	// TemplatePath is read from TemplateFS and written as-is. Used for the
	// verbatim-mirror pipeline where templates/ holds byte-for-byte copies of
	// upstream source.
	Raw bool
	// Symlink, when true, signals that the entry at TemplatePath is a
	// `.symlink` sentinel: its file contents are the symlink target, and
	// renderOne creates a real symlink at OutputPath pointing there. This
	// works around `go:embed` silently skipping real symlinks in the
	// template tree.
	Symlink bool
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
//
// Post-pivot semantics: plan is built by walking the embedded templates FS,
// and every entry is emitted as a raw pass-through. Orphan detection is
// intentionally skipped — `make generate-ddb` wipes the output directory
// before each run. The `--clean` flag is accepted for CLI compatibility but
// has no effect until orphan detection is reintroduced alongside
// parameterization.
func (p *Pipeline) Run(ctx context.Context) error {
	// 1. Plan all output files.
	planned := p.planFiles()

	fmt.Printf("Planning %d files for project %q\n", len(planned), p.config.ProjectMeta.Name)

	// 2. Render and write all files.
	summary := p.renderAndWrite(ctx, planned)

	// 3. Write .naff.yaml to output directory.
	if err := p.writeNaffConfig(); err != nil {
		return fmt.Errorf("writing %s: %w", config.NaffConfigFileName, err)
	}

	// 4. Print summary.
	fmt.Printf("\nGeneration complete: %d created, %d updated, %d unchanged",
		summary.Created, summary.Updated, summary.Unchanged)
	if summary.Errors > 0 {
		fmt.Printf(", %d errors", summary.Errors)
	}
	fmt.Println()

	if summary.Errors > 0 {
		// Print first few errors for debugging.
		errorCount := 0
		for _, r := range summary.Results {
			if r.Status == FileError && errorCount < 5 {
				fmt.Fprintf(os.Stderr, "  error: %s: %v\n", r.Path, r.Err)
				errorCount++
			}
		}
		return fmt.Errorf("%d files had errors during generation", summary.Errors)
	}

	return nil
}

// writeNaffConfig serializes the project config to .naff.yaml in the output directory.
func (p *Pipeline) writeNaffConfig() error {
	data, err := config.Marshal(p.config)
	if err != nil {
		return err
	}

	outPath := filepath.Join(p.outputDir, config.NaffConfigFileName)
	_, err = diffAwareWrite(outPath, data, 0)
	return err
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

	if pf.Symlink {
		target, err := TemplateFS.ReadFile(pf.TemplatePath)
		if err != nil {
			return FileResult{Path: outPath, Status: FileError, Err: fmt.Errorf("reading %s: %w", pf.TemplatePath, err)}
		}
		status, err := writeSymlink(outPath, strings.TrimSpace(string(target)))
		if err != nil {
			return FileResult{Path: outPath, Status: FileError, Err: err}
		}
		return FileResult{Path: outPath, Status: status}
	}

	var output []byte

	if pf.Raw {
		// Raw pass-through: read bytes directly from the embedded FS. No
		// text/template parsing, no Go formatting. Preserves content exactly
		// (including files containing "{{" that would confuse text/template).
		raw, err := TemplateFS.ReadFile(pf.TemplatePath)
		if err != nil {
			return FileResult{Path: outPath, Status: FileError, Err: fmt.Errorf("reading %s: %w", pf.TemplatePath, err)}
		}
		output = raw
	} else {
		// Render the template.
		content, err := p.renderer.Render(pf.TemplatePath, pf.Data)
		if err != nil {
			return FileResult{Path: outPath, Status: FileError, Err: fmt.Errorf("rendering %s: %w", pf.TemplatePath, err)}
		}

		output = []byte(content)

		// Format Go files.
		if pf.IsGo {
			formatted, fmtErr := renderer.FormatGo(output)
			if fmtErr != nil {
				return FileResult{Path: outPath, Status: FileError, Err: fmt.Errorf("formatting %s: %w", pf.OutputPath, fmtErr)}
			}
			output = formatted
		}
	}

	// Diff-aware write.
	status, err := diffAwareWrite(outPath, output, pf.Mode)
	if err != nil {
		return FileResult{Path: outPath, Status: FileError, Err: err}
	}

	return FileResult{Path: outPath, Status: status}
}

// writeSymlink ensures a symlink exists at path pointing to target. Existing
// entries (regular file, dir, or symlink with a different target) are
// replaced. Returns FileUnchanged if the link already points where we want.
func writeSymlink(path, target string) (FileStatus, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return FileError, fmt.Errorf("creating directory %s: %w", dir, err)
	}

	if existing, err := os.Readlink(path); err == nil {
		if existing == target {
			return FileUnchanged, nil
		}
		if err = os.Remove(path); err != nil {
			return FileError, fmt.Errorf("removing stale symlink %s: %w", path, err)
		}
		if err = os.Symlink(target, path); err != nil {
			return FileError, fmt.Errorf("symlinking %s -> %s: %w", path, target, err)
		}
		return FileUpdated, nil
	}

	// Path may exist as a regular file or directory from a prior non-symlink
	// emission; remove it before creating the link.
	if _, err := os.Lstat(path); err == nil {
		if err = os.RemoveAll(path); err != nil {
			return FileError, fmt.Errorf("removing %s: %w", path, err)
		}
	}

	if err := os.Symlink(target, path); err != nil {
		return FileError, fmt.Errorf("symlinking %s -> %s: %w", path, target, err)
	}
	return FileCreated, nil
}

// diffAwareWrite writes content to path only if it differs from existing content.
func diffAwareWrite(path string, content []byte, mode os.FileMode) (FileStatus, error) {
	if mode == 0 {
		mode = 0o644
	}

	// Check if file exists and has same content.
	existing, err := os.ReadFile(path)
	if err == nil {
		if string(existing) == string(content) {
			return FileUnchanged, nil
		}
		// File exists but content differs.
		if err = os.WriteFile(path, content, mode); err != nil {
			return FileError, fmt.Errorf("writing %s: %w", path, err)
		}
		return FileUpdated, nil
	}

	// File doesn't exist, create it.
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return FileError, fmt.Errorf("creating directory %s: %w", dir, err)
	}

	if err = os.WriteFile(path, content, mode); err != nil {
		return FileError, fmt.Errorf("writing %s: %w", path, err)
	}

	return FileCreated, nil
}


// planFiles builds the list of files to generate by walking the embedded
// templates FS. Every regular file becomes a PlannedFile. Entries whose name
// ends in `.tmpl` are rendered through text/template (strip `.tmpl` from the
// output path); everything else is pass-through copied.
//
// The top-level `_backend` directory is emitted as `backend/` — the leading
// `_` is a trick to keep Go tooling from trying to compile the embedded Go
// sources as part of the naff module. The sentinel filenames `_go.mod` /
// `_go.sum` inside `_backend/` are emitted as `go.mod` / `go.sum` for the
// same reason (Go's module resolver flags nested go.mod files as a module
// boundary, which `go:embed` refuses to cross).
func (p *Pipeline) planFiles() []PlannedFile {
	var files []PlannedFile

	// Every .tmpl receives the same context. Templates currently only
	// reference `.Project.ProjectMeta.Module` (verified with a grep across
	// templates/), so a single shared map suffices — no per-file data.
	data := map[string]any{"Project": p.config}

	_ = fs.WalkDir(TemplateFS, ".", func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		// Ignore the embed.go source file itself, which is embedded as a
		// degenerate side-effect of `//go:embed all:...` in some contexts.
		if p == "embed.go" {
			return nil
		}

		out := rewriteOutputPath(p)
		if out == "" {
			// Nothing to emit for this file (filtered).
			return nil
		}

		// `.symlink` sentinel: a regular file whose contents are the link
		// target. Emitted as a real symlink at the de-suffixed output path.
		// Used because `go:embed` silently drops actual symlinks.
		if strings.HasSuffix(p, ".symlink") {
			files = append(files, PlannedFile{
				TemplatePath: p,
				OutputPath:   strings.TrimSuffix(out, ".symlink"),
				Symlink:      true,
			})
			return nil
		}

		raw := !strings.HasSuffix(p, ".tmpl")
		if !raw {
			out = strings.TrimSuffix(out, ".tmpl")
		}

		pf := PlannedFile{
			TemplatePath: p,
			OutputPath:   out,
			Raw:          raw,
			Mode:         modeFor(out),
		}
		if !raw {
			pf.Data = data
		}
		files = append(files, pf)
		return nil
	})

	return files
}

// rewriteOutputPath maps a template-FS path to the output-tree path.
// Handles the `_`-prefix hiding conventions described on planFiles.
// Returns "" to skip the file.
func rewriteOutputPath(p string) string {
	// `_root/` is a special collector for upstream repo-root files (e.g.
	// .claude/, .github/, Makefile, README.md). Strip the whole segment so
	// these emit at the output root.
	if strings.HasPrefix(p, "_root/") {
		p = p[len("_root/"):]
	} else if idx := strings.IndexByte(p, '/'); idx > 0 && p[0] == '_' {
		// Strip leading `_` on the top-level directory (e.g. `_backend/...`
		// → `backend/...`). Only the top-level segment gets this treatment
		// — a genuine `_terraform_.tf` file further down retains its
		// underscore.
		p = p[1:idx] + p[idx:]
	}

	// Sentinel filename renames. Only applies to path basenames that exactly
	// match these sentinels — a file literally named `_go.mod.something`
	// would not match.
	base := filepath.Base(p)
	switch base {
	case "_go.mod":
		p = filepath.Join(filepath.Dir(p), "go.mod")
	case "_go.sum":
		p = filepath.Join(filepath.Dir(p), "go.sum")
	}

	return p
}

// modeFor returns the file mode for an output path. Executable bit is set
// for shell scripts and anything living under a `scripts/` directory —
// embed.FS flattens modes to 0o444, so we restore +x here for things that
// must be runnable in the generated project.
func modeFor(outPath string) os.FileMode {
	if strings.HasSuffix(outPath, ".sh") {
		return 0o755
	}
	for _, seg := range strings.Split(filepath.ToSlash(outPath), "/") {
		if seg == "scripts" {
			return 0o755
		}
	}
	return 0
}
