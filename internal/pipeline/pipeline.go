package pipeline

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/naming"
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
	FileCreated FileStatus = iota
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
	debug     bool
	renderer  *renderer.Renderer
}

// New creates a new Pipeline. When debug is true, post-generation `make`
// invocations stream their stdout/stderr live to the parent process; otherwise
// their output is captured and only surfaced if the step fails.
func New(cfg *config.Project, outputDir string, clean, debug bool) *Pipeline {
	return &Pipeline{
		config:    cfg,
		outputDir: outputDir,
		clean:     clean,
		debug:     debug,
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
// naff is designed as a one-shot scaffolder: generate into a fresh output
// directory, commit the result, and evolve the tree by hand from there.
// Re-running against a populated output is not a supported workflow, so
// orphan detection is intentionally absent. The `--clean` flag is accepted
// for CLI compatibility but currently has no effect.
func (p *Pipeline) Run(ctx context.Context) error {
	// Total elapsed time for the whole run, surfaced at the end (non-debug
	// only — --debug deliberately leaves output untouched).
	t0 := time.Now()
	timersEnabled := !p.debug

	// 1. Plan all output files.
	planned := p.planFiles()

	fmt.Printf("Planning %d files for project %q\n", len(planned), p.config.ProjectMeta.Name)

	// 2. Render and write all files.
	renderTimer := newPhaseTimer(fmt.Sprintf("rendering %d files", len(planned)), timersEnabled)
	renderTimer.Start()
	summary := p.renderAndWrite(ctx, planned)
	renderTimer.Stop(summary.Errors == 0)

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

	// 5. Run post-generation make targets in the freshly-written tree. Order
	// is load-bearing: `proto` consumes sqlc-generated types from the
	// backend, so the backend `queries` step (rolled into `generated_files`)
	// must run first. `vendor` resolves the module graph; `generated_files`
	// fans out to configs/queries/env_vars/format_golang. Proto runs last.
	backendDir := filepath.Join(p.outputDir, "backend")
	steps := []postGenStep{
		{dir: backendDir, args: []string{"querier", "format", "vendor", "configs", "env_vars"}, debug: p.debug},
		{dir: p.outputDir, args: []string{"proto"}, debug: p.debug},
	}
	for _, s := range steps {
		stepTimer := newPhaseTimer(fmt.Sprintf("make %s", strings.Join(s.args, " ")), timersEnabled)
		stepTimer.Start()
		err := s.run(ctx)
		stepTimer.Stop(err == nil)
		if err != nil {
			return err
		}
	}

	if timersEnabled {
		fmt.Printf("\nDone in %s\n", formatDur(time.Since(t0)))
	}

	return nil
}

// postGenStep is a single `make <args...>` invocation in a specific directory.
type postGenStep struct {
	dir   string
	args  []string
	debug bool
}

// run shells out to `make` in s.dir. When s.debug is true, stdout/stderr
// stream live to the parent process so protoc/go output is visible. Otherwise
// the combined output is captured and only flushed to stderr if the step
// fails — keeping the default `naff generate` output uncluttered while
// preserving diagnosability on error. A non-zero exit is fatal.
func (s postGenStep) run(ctx context.Context) error {
	if s.debug {
		// In debug mode the caller is not wrapping us in a timer, so retain
		// the original start-of-step marker so streamed make output has a
		// header line in front of it.
		fmt.Printf("\n→ make %s  (in %s)\n", strings.Join(s.args, " "), s.dir)
	}

	cmd := exec.CommandContext(ctx, "make", s.args...)
	cmd.Dir = s.dir

	var buf bytes.Buffer
	if s.debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// Combine stdout and stderr into one buffer so interleaving is
		// preserved when we replay the output on failure.
		cmd.Stdout = &buf
		cmd.Stderr = &buf
	}

	if err := cmd.Run(); err != nil {
		if !s.debug && buf.Len() > 0 {
			_, _ = os.Stderr.Write(buf.Bytes())
		}
		return fmt.Errorf("post-gen `make %s` in %s: %w", strings.Join(s.args, " "), s.dir, err)
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
		summary   Summary
		mu        sync.Mutex
		created   atomic.Int32
		updated   atomic.Int32
		unchanged atomic.Int32
		errors    atomic.Int32
		wg        sync.WaitGroup
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

	// Every .tmpl receives the same context. Templates reference the project
	// module (`.Project.ProjectMeta.Module`) and the project-name-derived
	// form methods on `ProjectMeta` (`.Title`, `.Pascal`, `.Kebab`, `.Snake`,
	// `.TitleKebab`, `.Abbrev`, `.LowerAbbrev`), so a single shared map
	// suffices — no per-file data.
	data := map[string]any{"Project": p.config}

	// Verbatim-emission skip set: templates owned by the parameterized layer
	// for user-defined domains. The FS walk below must not emit these as
	// raw pass-throughs — the parameterized planner (below) will emit them
	// instead, rendered from the config's entity schema.
	skipVerbatim := make(map[string]bool, len(p.config.Domains))
	for i := range p.config.Domains {
		d := &p.config.Domains[i]
		skipVerbatim[fmt.Sprintf("_backend/internal/domain/%s/keys/keys.go.tmpl", d.Name)] = true
		// Per-entity hand-written templates are replaced by the parameterized
		// per-entity layer. Path convention mirrors perEntityOutputName: the
		// verbatim template lives at
		// _backend/internal/domain/<domain>/<snake_entity>.go.tmpl.
		for j := range d.Entities {
			e := &d.Entities[j]
			snake := naming.New(e.Name).Snake()
			skipVerbatim[fmt.Sprintf("_backend/internal/domain/%s/%s.go.tmpl", d.Name, snake)] = true
		}
	}

	_ = fs.WalkDir(TemplateFS, ".", func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			// `_backend/_parameterized/` holds templates consumed by the
			// parameterized planner, not the verbatim walk. Skip the whole
			// subtree so nothing here leaks into the output as a raw copy.
			if p == "_backend/_parameterized" {
				return fs.SkipDir
			}
			return nil
		}
		// Ignore the embed.go source file itself, which is embedded as a
		// degenerate side-effect of `//go:embed all:...` in some contexts.
		if p == "embed.go" {
			return nil
		}
		if skipVerbatim[p] {
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

	// Parameterized layer: one generated keys package per user-defined
	// domain. The template at _backend/_parameterized/domain_keys/keys.go.tmpl
	// iterates .Domain.Entities and emits <Name>Key / <Name>IDKey constants.
	// Replaces (not co-resident with) any hand-rolled keys.go.tmpl at the
	// matching path — see skipVerbatim above.
	for i := range p.config.Domains {
		d := &p.config.Domains[i]
		files = append(files, PlannedFile{
			TemplatePath: "_backend/_parameterized/domain_keys/keys.go.tmpl",
			OutputPath:   fmt.Sprintf("backend/internal/domain/%s/keys/keys.go", d.Name),
			Data:         parameterizedCtx{Project: p.config, Domain: d},
			IsGo:         true,
		})
	}

	// Per-entity parameterized layer: one <snake_entity>.go per entity in
	// each user-defined domain. Emits the canonical CRUD shape (domain
	// struct, input variants, DataManager interface, Update, validation).
	// The verbatim hand-written per-entity templates under
	// _backend/internal/domain/<domain>/<entity>.go.tmpl are registered in
	// skipVerbatim above so only the parameterized output reaches the
	// generated tree. Bespoke methods the schema cannot express (custom
	// search indexing, state transitions, etc.) must be added by hand in
	// a sidecar file after generation.
	files = append(files, p.planPerEntityFiles()...)

	// Per-entity converter files. Gated to scalar-only entities —
	// see scalarOnlyConverter for the eligibility rule.
	files = append(files, p.planConverterFiles()...)

	return files
}

// converterEligible reports whether a converter file for the given entity
// can be emitted by the parameterized layer alone, without any hand-
// authored assistance for shapes the schema cannot currently express.
//
// The schema now expresses: scalar fields (via Fields), required and
// optional links (via LinksTo + DomainAccessor), and parent/scope routing
// flags (BelongsTo, BelongsToAccount, BelongsToUser, CreatedByUser). The
// template handles all of these. Every entity is therefore in principle
// eligible — existing hand-written converters, if still present, win via
// the fs-presence check in planConverterFiles.
//
// Shapes not yet expressible in schema (and therefore left to hand-
// written templates): nested-collection fields (e.g. Meal.Components),
// non-denormalized secondary links (e.g. UserIngredientPreference.
// ValidIngredientGroup), range/min/max-pair types that don't reduce to
// scalar pair fields.
func converterEligible(_ *config.Entity) bool {
	return true
}

// planConverterFiles enumerates every entity across every domain and emits
// a parameterized converter file per eligible entity — but only when no
// hand-written converter for that entity exists anywhere in the sibling
// `converters/` directory. The detection considers both filename match
// (`<plural_snake>.go.tmpl`) and function-name collisions (a parent's
// converter template may emit its child's converters alongside its own).
// This makes the layer additive: hand-written templates always win, and
// parameterized output fills only the genuine gaps. To "promote" an
// entity to the parameterized layer, delete every hand-written
// declaration of its Convert<Entity>* functions.
//
// Output lands at backend/internal/domain/<domain>/converters/<plural_snake>.go.
func (p *Pipeline) planConverterFiles() []PlannedFile {
	const tmpl = "_backend/_parameterized/domain_converters/entity_converters.go.tmpl"

	var files []PlannedFile
	for i := range p.config.Domains {
		d := &p.config.Domains[i]
		existing := existingConverterDeclarations(d.Name)
		for j := range d.Entities {
			e := &d.Entities[j]
			if !converterEligible(e) {
				continue
			}
			if existing[e.Name] {
				continue
			}
			pluralSnake := naming.New(e.Name).PluralSnake()
			files = append(files, PlannedFile{
				TemplatePath: tmpl,
				OutputPath:   fmt.Sprintf("backend/internal/domain/%s/converters/%s.go", d.Name, pluralSnake),
				Data:         parameterizedCtx{Project: p.config, Domain: d, Entity: e},
				IsGo:         true,
			})
		}
	}
	return files
}

// existingConverterDeclarations reports which entity names already have
// Convert<Entity>* functions declared by some hand-written converter
// template under _backend/internal/domain/<domain>/converters/. Keys in
// the returned map are entity PascalCase names. Matching is over the
// literal prefix `func Convert<EntityName>` so it catches child converters
// emitted inside a parent's file (e.g. meals.go.tmpl declaring
// ConvertMealComponent*), which the filename alone would miss.
func existingConverterDeclarations(domainName string) map[string]bool {
	out := map[string]bool{}
	root := fmt.Sprintf("_backend/internal/domain/%s/converters", domainName)
	_ = fs.WalkDir(TemplateFS, root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fs.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go.tmpl") {
			return nil
		}
		body, err := TemplateFS.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, line := range strings.Split(string(body), "\n") {
			const marker = "func Convert"
			idx := strings.Index(line, marker)
			if idx < 0 {
				continue
			}
			rest := line[idx+len(marker):]
			// Entity name runs to the next uppercase-delimited word boundary.
			// The next fixed suffix is one of "To", "CreationRequestInput".
			// Extract up to the first of those.
			end := len(rest)
			for _, boundary := range []string{"To", "CreationRequestInput"} {
				if i := strings.Index(rest, boundary); i > 0 && i < end {
					end = i
				}
			}
			if end == 0 || end == len(rest) {
				continue
			}
			name := rest[:end]
			// Sanity: name must start with uppercase and contain only ident chars.
			if name == "" || name[0] < 'A' || name[0] > 'Z' {
				continue
			}
			out[name] = true
		}
		return nil
	})
	return out
}

// planPerEntityFiles enumerates every .tmpl under
// _backend/_parameterized/domain_entity/ and emits one PlannedFile per
// (template, domain, entity) tuple. Output paths derive from the template
// filename by replacing the leading `entity` stem with the snake-cased
// entity name. The emitted file lands directly in the domain directory
// (backend/internal/domain/<domain>/<snake>.go) — the planner's
// skipVerbatim set ensures no hand-written template collides at that
// path.
func (p *Pipeline) planPerEntityFiles() []PlannedFile {
	const root = "_backend/_parameterized/domain_entity"

	var templates []string
	_ = fs.WalkDir(TemplateFS, root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			// Subtree absent — no per-entity templates yet; nothing to plan.
			return fs.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		templates = append(templates, path)
		return nil
	})

	if len(templates) == 0 {
		return nil
	}

	var files []PlannedFile
	for i := range p.config.Domains {
		d := &p.config.Domains[i]
		for j := range d.Entities {
			e := &d.Entities[j]
			for _, tmpl := range templates {
				outName := perEntityOutputName(filepath.Base(tmpl), e.Name)
				files = append(files, PlannedFile{
					TemplatePath: tmpl,
					OutputPath:   fmt.Sprintf("backend/internal/domain/%s/%s", d.Name, outName),
					Data:         parameterizedCtx{Project: p.config, Domain: d, Entity: e},
					IsGo:         strings.HasSuffix(outName, ".go"),
				})
			}
		}
	}
	return files
}

// perEntityOutputName maps a per-entity template filename + entity name to
// the emitted basename. Convention: templates use `entity` as the leading
// stem; we substitute the snake-cased entity name and strip `.tmpl`.
//
//	entity.go.tmpl,      "RecipeRating" -> "recipe_rating.go"
//	entity_test.go.tmpl, "RecipeRating" -> "recipe_rating_test.go"
func perEntityOutputName(templateBase, entityName string) string {
	stem := strings.TrimSuffix(templateBase, ".tmpl") // entity.go | entity_test.go
	snake := naming.New(entityName).Snake()
	return strings.Replace(stem, "entity", snake, 1)
}

// parameterizedCtx is the template data shape for files emitted by the
// parameterized planner layer. Project + Domain are always populated;
// Entity is only populated for per-entity templates under
// _backend/_parameterized/domain_entity/.
type parameterizedCtx struct {
	Project *config.Project
	Domain  *config.Domain
	Entity  *config.Entity
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
