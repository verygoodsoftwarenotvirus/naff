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

	"github.com/verygoodsoftwarenotvirus/naff/internal/builtins"
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

	// 4. Write .naff.yaml to output directory.
	if err := p.writeNaffConfig(); err != nil {
		return fmt.Errorf("writing %s: %w", config.NaffConfigFileName, err)
	}

	// 5. Print summary.
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
	status, err := diffAwareWrite(outPath, output, pf.Mode)
	if err != nil {
		return FileResult{Path: outPath, Status: FileError, Err: err}
	}

	return FileResult{Path: outPath, Status: status}
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

// findOrphans finds *.generated.* files in the output directory that are not in the plan.
// It also detects orphaned frontend route directories.
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

	// Also detect orphaned frontend route files (which don't use .generated. in names).
	for _, appDir := range []string{"frontend/consumer/src/routes", "frontend/admin/src/routes"} {
		routesDir := filepath.Join(p.outputDir, appDir)
		_ = filepath.Walk(routesDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if !plannedPaths[path] {
				orphans = append(orphans, path)
			}
			return nil
		})
	}

	return orphans
}

// planFiles builds the complete list of files to generate.
// This is the central registry that maps templates to output paths.
func (p *Pipeline) planFiles() []PlannedFile {
	var files []PlannedFile

	generateBackend := p.config.Targets.Backend
	generateIOS := p.config.Targets.IOS

	// Merge built-in domains with user-defined domains.
	allDomains := builtins.BuiltinDomains(p.config.Features)
	allDomains = append(allDomains, p.config.Domains...)

	// Project-level files (generated once per project).
	if generateBackend {
		files = append(files, p.planProjectFiles()...)
		files = append(files, p.planCmdFiles(allDomains)...)
		files = append(files, p.planAuthorizationFiles()...)
		files = append(files, p.planAuthenticationFiles()...)
		files = append(files, p.planConfigFiles()...)
		files = append(files, p.planRepositoriesFiles()...)
		files = append(files, p.planLocalDevFiles()...)
		files = append(files, p.planDomainExtrasFiles()...)
		files = append(files, p.planAuthHandlerFiles()...)
		// Phase 6 ships DDB-specific hand-rolled packages + their codegen output
		// verbatim. These assume all default-on features (webhooks, issuereports,
		// comments, notifications, uploadedmedia, dataprivacy, settings) are
		// enabled AND the project has a `mealplanning` domain. For lightweight
		// projects that disable any of these, phase 6 would emit code that
		// references absent domains. Gate it on the shape.
		if p.isDDBCompat() {
			files = append(files, p.planPhase6Files()...)
			files = append(files, p.planPhase6CodegenFiles()...)
		}
	}
	if generateIOS {
		files = append(files, p.planIOSProjectFiles()...)
	}

	for _, domain := range allDomains {
		d := domain

		// Per-domain files are generated once, using the first entity as context.
		if len(d.Entities) > 0 {
			firstEntity := d.Entities[0]
			domainCtx := TemplateContext{
				Project: p.config,
				Domain:  &d,
				Entity:  &firstEntity,
			}
			if generateBackend {
				files = append(files, p.planPerDomainFiles(domainCtx)...)
			}
			if generateIOS {
				files = append(files, p.planIOSPerDomainFiles(domainCtx)...)
			}
		}

		// Per-entity files are generated for each entity.
		for _, entity := range d.Entities {
			e := entity
			entityCtx := TemplateContext{
				Project: p.config,
				Domain:  &d,
				Entity:  &e,
			}
			if generateBackend {
				files = append(files, p.planPerEntityFiles(entityCtx)...)
			}
			if generateIOS {
				files = append(files, p.planIOSPerEntityFiles(entityCtx)...)
			}
		}
	}

	// Frontend apps (feature-gated).
	generateConsumer := p.config.Features.FeatureEnabled(p.config.Features.ConsumerApp, false)
	generateAdmin := p.config.Features.FeatureEnabled(p.config.Features.AdminApp, false)

	if generateConsumer {
		files = append(files, p.planConsumerAppFiles(allDomains)...)
	}
	if generateAdmin {
		files = append(files, p.planAdminAppFiles(allDomains)...)
	}

	return files
}

// TemplateContext is the data passed to every template.
type TemplateContext struct {
	Project *config.Project
	Domain  *config.Domain
	Entity  *config.Entity
}

// planPerDomainFiles returns files generated once per domain (using the first entity).
func (p *Pipeline) planPerDomainFiles(ctx TemplateContext) []PlannedFile {
	domain := ctx.Domain.Name

	templateMappings := []struct {
		template string
		output   string
		isGo     bool
	}{
		// Domain-level (shared across entities in domain)
		{"domain/repository.go.tmpl", fmt.Sprintf("internal/domain/%s/repository.generated.go", domain), true},
		{"domain_keys/keys.go.tmpl", fmt.Sprintf("internal/domain/%s/keys/keys.generated.go", domain), true},
		{"domain_fakes/fake.go.tmpl", fmt.Sprintf("internal/domain/%s/fakes/fake.generated.go", domain), true},
		{"domain_mock/repository.go.tmpl", fmt.Sprintf("internal/domain/%s/mock/repository.generated.go", domain), true},

		// Manager (one per domain)
		{"manager/interface.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/interface.generated.go", domain), true},
		{"manager/manager.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/manager.generated.go", domain), true},
		{"manager/do.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/do.generated.go", domain), true},
		{"manager/mock_manager.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/mock/manager.generated.go", domain), true},

		// Repository infrastructure (one per domain)
		{"repository/client.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/client.generated.go", domain), true},
		{"repository/do.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/do.generated.go", domain), true},

		// gRPC infrastructure (one per domain)
		{"grpc/service.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/service.generated.go", domain), true},
		{"grpc/converters.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/converters/converters.generated.go", domain), true},
		{"grpc/permissions.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/permissions.generated.go", domain), true},
		{"grpc/do.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/do.generated.go", domain), true},

		// Authorization (one per domain)
		{"authorization/permissions.go.tmpl", fmt.Sprintf("internal/authorization/%s_permissions.generated.go", domain), true},

		// sqlc config (one per domain)
		{"sqlc/block.yaml.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/sqlc.generated.yaml", domain), false},

		// Build extras
		{"build/extras.go.tmpl", fmt.Sprintf("internal/build/services/api/grpc/%s.generated.go", domain), false},
	}

	var files []PlannedFile
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

// planProjectFiles returns files generated once per project (Makefile, scripts, etc.).
func (p *Pipeline) planProjectFiles() []PlannedFile {
	ctx := TemplateContext{
		Project: p.config,
	}

	type mapping struct {
		template string
		output   string
		mode     os.FileMode
	}

	templateMappings := []mapping{
		{"project/go.mod.tmpl", "go.mod", 0},
		{"project/Makefile.tmpl", "Makefile", 0},
		{"project/scripts/configs.sh.tmpl", "scripts/configs.sh", 0o755},
		{"project/scripts/queries.sh.tmpl", "scripts/queries.sh", 0o755},
		{"project/scripts/env_vars.sh.tmpl", "scripts/env_vars.sh", 0o755},
		{"project/scripts/format_golang.sh.tmpl", "scripts/format_golang.sh", 0o755},
		{"project/scripts/goimports.sh.tmpl", "scripts/goimports.sh", 0o755},
		{"project/scripts/format_imports.sh.tmpl", "scripts/format_imports.sh", 0o755},
		{"project/scripts/format_go_fieldalignment.sh.tmpl", "scripts/format_go_fieldalignment.sh", 0o755},
		{"project/scripts/format_go_tag_alignment.sh.tmpl", "scripts/format_go_tag_alignment.sh", 0o755},
		// Top-level shared proto: filtering (QueryFilter, Pagination) — imported
		// by every service's *_service.proto and by testing/integration.
		{"proto/filtering.proto.tmpl", "proto/filtering.proto", 0},
	}

	var files []PlannedFile
	for _, m := range templateMappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			Mode:         m.mode,
		})
	}
	return files
}

// planAuthHandlerFiles returns verbatim templates for
// internal/services/auth/handlers/authentication/, target's 982-LOC hand-rolled
// OAuth2/password/TOTP/PASETO login orchestration. Phase 3c (2026-04-14):
// tried with the original ~2,918 LOC number, actual non-test LOC is 982 across
// 12 files. All dependencies (domain/{auth,identity,oauth}, authentication,
// authorization, testutils) shipped in prior phases.
func (p *Pipeline) planAuthHandlerFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	base := "services/auth/handlers/authentication"
	mappings := []mapping{
		{base + "/authentication_http_routes.go.tmpl", "internal/" + base + "/authentication_http_routes.generated.go"},
		{base + "/config.go.tmpl", "internal/" + base + "/config.generated.go"},
		{base + "/do.go.tmpl", "internal/" + base + "/do.generated.go"},
		{base + "/doc.go.tmpl", "internal/" + base + "/doc.generated.go"},
		{base + "/helpers.go.tmpl", "internal/" + base + "/helpers.generated.go"},
		{base + "/oauth2.go.tmpl", "internal/" + base + "/oauth2.generated.go"},
		{base + "/oauth2_client_info.go.tmpl", "internal/" + base + "/oauth2_client_info.generated.go"},
		{base + "/oauth2_client_store.go.tmpl", "internal/" + base + "/oauth2_client_store.generated.go"},
		{base + "/oauth2_client_token.go.tmpl", "internal/" + base + "/oauth2_client_token.generated.go"},
		{base + "/oauth2_token_store.go.tmpl", "internal/" + base + "/oauth2_token_store.generated.go"},
		{base + "/revoke.go.tmpl", "internal/" + base + "/revoke.generated.go"},
		{base + "/service.go.tmpl", "internal/" + base + "/service.generated.go"},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}

// planDomainExtrasFiles returns verbatim hand-rolled domain/ content for auth,
// identity, and oauth. Phase 3b (2026-04-14): target doesn't use naff's CRUD
// entity template for these three domains — it hand-rolls each type with
// domain-specific logic (TOTP validation, password-change flows, 2FA state,
// request DTOs, etc.). Phase 3a's authDomain/identityDomain/oauthDomain
// entries have been pulled from builtins.go; this planner now owns all three
// domain packages' type declarations plus their authz permissions sidecar files.
func (p *Pipeline) planDomainExtrasFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	mappings := []mapping{
		// internal/domain/auth/
		{"auth/auth.go.tmpl", "internal/domain/auth/auth.generated.go"},
		{"auth/do.go.tmpl", "internal/domain/auth/do.generated.go"},
		{"auth/password_reset_token.go.tmpl", "internal/domain/auth/password_reset_token.generated.go"},
		{"auth/repository.go.tmpl", "internal/domain/auth/repository.generated.go"},
		{"auth/user.go.tmpl", "internal/domain/auth/user.generated.go"},
		{"auth/user_session.go.tmpl", "internal/domain/auth/user_session.generated.go"},

		// internal/domain/identity/
		{"identity/auth.go.tmpl", "internal/domain/identity/auth.generated.go"},
		{"identity/account.go.tmpl", "internal/domain/identity/account.generated.go"},
		{"identity/account_invitation.go.tmpl", "internal/domain/identity/account_invitation.generated.go"},
		{"identity/account_user_membership.go.tmpl", "internal/domain/identity/account_user_membership.generated.go"},
		{"identity/data_privacy.go.tmpl", "internal/domain/identity/data_privacy.generated.go"},
		{"identity/do.go.tmpl", "internal/domain/identity/do.generated.go"},
		{"identity/password_reset_token.go.tmpl", "internal/domain/identity/password_reset_token.generated.go"},
		{"identity/repository.go.tmpl", "internal/domain/identity/repository.generated.go"},
		{"identity/user.go.tmpl", "internal/domain/identity/user.generated.go"},

		// internal/domain/oauth/
		{"oauth/do.go.tmpl", "internal/domain/oauth/do.generated.go"},
		{"oauth/oauth2_client.go.tmpl", "internal/domain/oauth/oauth2_client.generated.go"},
		{"oauth/oauth2_client_token.go.tmpl", "internal/domain/oauth/oauth2_client_token.generated.go"},
		{"oauth/repository.go.tmpl", "internal/domain/oauth/repository.generated.go"},

		// internal/authorization/ (per-domain permissions sidecar files for
		// auth/identity/oauth, since those domains are no longer in BuiltinDomains()
		// and thus don't trigger planPerDomainFiles' convention-based permissions template).
		{"authorization/identity_permissions.go.tmpl", "internal/authorization/identity_permissions.generated.go"},
		{"authorization/oauth_permissions.go.tmpl", "internal/authorization/oauth_permissions.generated.go"},

		// internal/branding/ stub — only the four constants in-scope generated
		// code actually imports (CompanyName, CompanyNameSlug, CompanySlug,
		// EnvVarPrefix). Users flesh out logos/email-templates/legal-text. Target's
		// real branding.go is 109KB with base64 logo blobs → deliberately NOT shipped.
		{"branding/branding.go.tmpl", "internal/branding/branding.generated.go"},

		// Phase 4a + 4b.5 — integration harness. 4a shipped the stubs (matchers.go,
		// doc.go, constants.go). 4b.5 adds init.go/helpers.go/audit_helpers.go now
		// that 4b.0–4b.4 unblocked the deps: Provide<Domain>Repository rename
		// (4b.0), settings builtin (4b.1), top-level /internal/repositories glue
		// (4b.2), proto/filtering.proto (4b.3), and internal/localdev/server.go
		// (4b.4). init.go and helpers.go inherit the audit import alias from 4b.4
		// (naff emits postgres/audit; target wrote postgres/auditlogentries).
		{"testutils/matchers.go.tmpl", "internal/testutils/matchers.generated.go"},
		{"testing/integration/apiserver/doc.go.tmpl", "testing/integration/apiserver/doc.generated.go"},
		{"testing/integration/apiserver/constants.go.tmpl", "testing/integration/apiserver/constants.generated.go"},
		{"testing/integration/apiserver/init.go.tmpl", "testing/integration/apiserver/init.generated.go"},
		{"testing/integration/apiserver/helpers.go.tmpl", "testing/integration/apiserver/helpers.generated.go"},
		{"testing/integration/apiserver/audit_helpers.go.tmpl", "testing/integration/apiserver/audit_helpers.generated.go"},

		// Phase 3d — pkg/client. Single 248-LOC client plumbing file: gRPC+HTTP
		// transport, OAuth2 token integration, TLS config. Per-entity client methods
		// come from gRPC stubs in internal/grpc/generated/ — target-only protoc
		// output, not naff's concern.
		{"pkg/client/client.go.tmpl", "pkg/client/client.generated.go"},

		// Phase 5 — hand-rolled postgres repos + identity adjuncts. Target ships
		// auth/identity/oauth as hand-rolled domain packages (3b) AND hand-rolled
		// postgres repos. Without these, a consumer can't compile localdev,
		// services/auth/handlers, or the integration harness. Ship verbatim
		// modulo module paths — the exact 3b/4b pattern.

		// postgres repos for hand-rolled domains (auth/identity/oauth).
		// {sqlc}/generated + {sqlc_queries}/ subdirs are consumer's protoc/sqlc
		// responsibility — not emitted here.
		{"repositories/postgres/auth/client.go.tmpl", "internal/repositories/postgres/auth/client.generated.go"},
		{"repositories/postgres/auth/do.go.tmpl", "internal/repositories/postgres/auth/do.generated.go"},
		{"repositories/postgres/auth/password_reset_tokens.go.tmpl", "internal/repositories/postgres/auth/password_reset_tokens.generated.go"},
		{"repositories/postgres/auth/user_sessions.go.tmpl", "internal/repositories/postgres/auth/user_sessions.generated.go"},

		{"repositories/postgres/identity/account_invitations.go.tmpl", "internal/repositories/postgres/identity/account_invitations.generated.go"},
		{"repositories/postgres/identity/account_user_memberships.go.tmpl", "internal/repositories/postgres/identity/account_user_memberships.generated.go"},
		{"repositories/postgres/identity/accounts.go.tmpl", "internal/repositories/postgres/identity/accounts.generated.go"},
		{"repositories/postgres/identity/client.go.tmpl", "internal/repositories/postgres/identity/client.generated.go"},
		{"repositories/postgres/identity/data_privacy.go.tmpl", "internal/repositories/postgres/identity/data_privacy.generated.go"},
		{"repositories/postgres/identity/do.go.tmpl", "internal/repositories/postgres/identity/do.generated.go"},
		{"repositories/postgres/identity/users.go.tmpl", "internal/repositories/postgres/identity/users.generated.go"},
		{"repositories/postgres/identity/webauthn_credentials.go.tmpl", "internal/repositories/postgres/identity/webauthn_credentials.generated.go"},

		{"repositories/postgres/oauth/client.go.tmpl", "internal/repositories/postgres/oauth/client.generated.go"},
		{"repositories/postgres/oauth/do.go.tmpl", "internal/repositories/postgres/oauth/do.generated.go"},
		{"repositories/postgres/oauth/oauth2_client_tokens.go.tmpl", "internal/repositories/postgres/oauth/oauth2_client_tokens.generated.go"},
		{"repositories/postgres/oauth/oauth2_clients.go.tmpl", "internal/repositories/postgres/oauth/oauth2_clients.generated.go"},

		// postgres/testing — `pgtesting.BuildDatabaseContainer` + audit sentinel
		// used by localdev.BuildInProcessServer.
		{"repositories/postgres/testing/audit.go.tmpl", "internal/repositories/postgres/testing/audit.generated.go"},
		{"repositories/postgres/testing/helpers.go.tmpl", "internal/repositories/postgres/testing/helpers.generated.go"},

		// identity adjuncts — converters imported by localdev; fakes imported by
		// postgres/testing; keys (auth/identity/oauth) imported across repos.
		{"identity/converters/account_invitations.go.tmpl", "internal/domain/identity/converters/account_invitations.generated.go"},
		{"identity/converters/account_user_memberships.go.tmpl", "internal/domain/identity/converters/account_user_memberships.generated.go"},
		{"identity/converters/accounts.go.tmpl", "internal/domain/identity/converters/accounts.generated.go"},
		{"identity/converters/admin.go.tmpl", "internal/domain/identity/converters/admin.generated.go"},
		{"identity/converters/users.go.tmpl", "internal/domain/identity/converters/users.generated.go"},

		{"identity/fakes/account_invitation.go.tmpl", "internal/domain/identity/fakes/account_invitation.generated.go"},
		{"identity/fakes/account_user_membership.go.tmpl", "internal/domain/identity/fakes/account_user_membership.generated.go"},
		{"identity/fakes/account.go.tmpl", "internal/domain/identity/fakes/account.generated.go"},
		{"identity/fakes/doc.go.tmpl", "internal/domain/identity/fakes/doc.generated.go"},
		{"identity/fakes/fake.go.tmpl", "internal/domain/identity/fakes/fake.generated.go"},
		{"identity/fakes/user.go.tmpl", "internal/domain/identity/fakes/user.generated.go"},

		{"auth/keys/keys.go.tmpl", "internal/domain/auth/keys/keys.generated.go"},
		{"identity/keys/keys.go.tmpl", "internal/domain/identity/keys/keys.generated.go"},
		{"oauth/keys/keys.go.tmpl", "internal/domain/oauth/keys/keys.generated.go"},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}

// planAuthorizationFiles returns the domain-independent internal/authorization/ core.
// Phase 2d scope (2026-04-14): the role/Permission types, AccountRole/ServiceRole
// checkers, and the service-admin permission constants that service_role.go depends on.
// Per-domain permission files are still emitted by planPerDomainFiles via the existing
// authorization/permissions.go.tmpl (convention-based CRUD). The 520-LOC hand-rolled
// permission lists in target's permissions.go (ServiceAdminPermissions etc.) are NOT
// shipped — they reference per-entity Search/Cancel/Impersonate/etc. permissions that
// naff's convention template doesn't emit. Phase 3 will design a YAML schema that
// lets users declare per-entity verbs, then these lists can be templated properly.
func (p *Pipeline) planAuthorizationFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	mappings := []mapping{
		{"authorization/permissions_core.go.tmpl", "internal/authorization/permissions.generated.go"},
		{"authorization/account_role.go.tmpl", "internal/authorization/account_role.generated.go"},
		{"authorization/service_role.go.tmpl", "internal/authorization/service_role.generated.go"},
		{"authorization/auth_permissions.go.tmpl", "internal/authorization/auth_permissions.generated.go"},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}

// planAuthenticationFiles returns the internal/authentication/ shim layer.
// Phase 2a scope (2026-04-14): files with no dependency on internal/domain/{auth,identity,audit}
// or per-service config packages. Excludes manager.go, authentication/do.go, session_context.go,
// webauthn/service.go, webauthn/user_adapter.go — those land in Phase 2b alongside the
// auth/identity/oauth domains.
func (p *Pipeline) planAuthenticationFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	mappings := []mapping{
		{"authentication/aliases.go.tmpl", "internal/authentication/aliases.generated.go"},
		{"authentication/manager.go.tmpl", "internal/authentication/manager.generated.go"},
		{"authentication/do.go.tmpl", "internal/authentication/do.generated.go"},
		{"authentication/config/config.go.tmpl", "internal/authentication/config/config.generated.go"},
		{"authentication/config/do.go.tmpl", "internal/authentication/config/do.generated.go"},
		{"authentication/mock/mock_authenticator.go.tmpl", "internal/authentication/mock/mock_authenticator.generated.go"},
		{"authentication/mock/authentication_manager.go.tmpl", "internal/authentication/mock/authentication_manager.generated.go"},
		{"authentication/mocks/mock_user.go.tmpl", "internal/authentication/mocks/mock_user.generated.go"},
		{"authentication/sessions/errors.go.tmpl", "internal/authentication/sessions/errors.generated.go"},
		{"authentication/sessions/session_context.go.tmpl", "internal/authentication/sessions/session_context.generated.go"},
		{"authentication/sessions/do.go.tmpl", "internal/authentication/sessions/do.generated.go"},
		{"authentication/webauthn/session_store.go.tmpl", "internal/authentication/webauthn/session_store.generated.go"},
		{"authentication/webauthn/postgres_session_store.go.tmpl", "internal/authentication/webauthn/postgres_session_store.generated.go"},
		{"authentication/webauthn/service.go.tmpl", "internal/authentication/webauthn/service.generated.go"},
		{"authentication/webauthn/user_adapter.go.tmpl", "internal/authentication/webauthn/user_adapter.generated.go"},
		{"authentication/webauthn/config/config.go.tmpl", "internal/authentication/webauthn/config/config.generated.go"},
		{"identity/webauthn_credential.go.tmpl", "internal/domain/identity/webauthn_credential.generated.go"},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}

// planConfigFiles returns the internal/config/ primitives.
// Phase 2a scope (2026-04-14): meta.go, queues.go, doc.go only. The heavier
// configs.go / services_config.go / environment.go / do.go are Phase 2b
// (they import internal/authentication/config, internal/services/auth/handlers/authentication,
// and require per-domain templating).
func (p *Pipeline) planConfigFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	mappings := []mapping{
		{"config/meta.go.tmpl", "internal/config/meta.generated.go"},
		{"config/queues.go.tmpl", "internal/config/queues.generated.go"},
		{"config/doc.go.tmpl", "internal/config/doc.generated.go"},
		{"config/configs.go.tmpl", "internal/config/configs.generated.go"},
		{"config/do.go.tmpl", "internal/config/do.generated.go"},
		{"config/env_vars.go.tmpl", "internal/config/env_vars.generated.go"},
		{"config/environment.go.tmpl", "internal/config/environment.generated.go"},
		{"config/mealplanning_configs.go.tmpl", "internal/config/mealplanning_configs.generated.go"},
		{"config/mealplanning_environment.go.tmpl", "internal/config/mealplanning_environment.generated.go"},
		{"config/services_config.go.tmpl", "internal/config/services_config.generated.go"},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}

// planRepositoriesFiles returns the top-level internal/repositories/ glue —
// ProvideMigrator + its DI registration. Used by localdev/server.go and any
// consumer that wires a database.Migrator.
func (p *Pipeline) planRepositoriesFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	mappings := []mapping{
		{"repositories/do.go.tmpl", "internal/repositories/do.generated.go"},
		{"repositories/migrations.go.tmpl", "internal/repositories/migrations.generated.go"},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}

// planLocalDevFiles ships internal/localdev/server.go — the in-process
// all-in-one bootstrap used by testing/integration/apiserver. Verbatim from
// target modulo module paths, with one deliberate deviation: target imports
// `internal/repositories/postgres/auditlogentries`; naff emits that package
// at `postgres/audit`, so the template aliases it back to `auditlogentries`
// in the import block and leaves all call sites unchanged.
func (p *Pipeline) planLocalDevFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	return []PlannedFile{
		{
			TemplatePath: "localdev/server.go.tmpl",
			OutputPath:   "internal/localdev/server.generated.go",
			Data:         ctx,
			IsGo:         true,
		},
	}
}


// isDDBCompat reports whether the project is shaped like DDB (all default-on
// features enabled + user has a mealplanning domain). Phase 6 packages assume
// this shape. Projects that disable any default-on feature or omit mealplanning
// should skip Phase 6.
func (p *Pipeline) isDDBCompat() bool {
	f := p.config.Features
	if !f.FeatureEnabled(f.Webhooks, true) {
		return false
	}
	if !f.FeatureEnabled(f.IssueReports, true) {
		return false
	}
	if !f.FeatureEnabled(f.Comments, true) {
		return false
	}
	if !f.FeatureEnabled(f.Notifications, true) {
		return false
	}
	if !f.FeatureEnabled(f.UploadedMedia, true) {
		return false
	}
	if !f.FeatureEnabled(f.DataPrivacy, true) {
		return false
	}
	if !f.FeatureEnabled(f.Settings, true) {
		return false
	}
	hasMealPlanning := false
	for _, d := range p.config.Domains {
		if d.Name == "mealplanning" {
			hasMealPlanning = true
			break
		}
	}
	return hasMealPlanning
}

// planPhase6Files ships hand-rolled target packages verbatim — 204 files,
// ~70K LOC across 71 packages. Same pattern as 3b/5: byte-identical to
// target modulo module paths. These are a mix of:
//   - infrastructure (build/services/{api,mcp}, grpc/converters,
//     repositories/postgres/migrations, services/errors)
//   - per-service config packages (services/{dataprivacy,identity,
//     mealplanning,oauth,payments,uploadedmedia}/config)
//   - hand-rolled adjuncts to 3b domains (domain/{auth,identity,oauth}/{
//     converters,fakes,managers,manager,mock})
//   - DDB product-specific packages (mealplanning/{bootstrap,emails,
//     grocerylistpreparation,managers,notifications,privacy,recipeanalysis,
//     recipevalidator,registration}, services/mealplanning/{errors,
//     indexing,workers/*}, services/internalops/*, functions/*,
//     build/jobs/*, build/functions/*)
// Phase 6 is honest: to make DDB's generated tree compile, we ship DDB's
// entire non-codegen source tree. A future phase could split these by
// genericness and move product-specific packages out of naff.
func (p *Pipeline) planPhase6Files() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	mappings := []mapping{
		{"internal/build/functions/data_change_message_handler/build.go.tmpl", "internal/build/functions/data_change_message_handler/build.generated.go"},
		{"internal/build/functions/data_change_message_handler/config.go.tmpl", "internal/build/functions/data_change_message_handler/config.generated.go"},
		{"internal/build/functions/data_change_message_handler/misc.go.tmpl", "internal/build/functions/data_change_message_handler/misc.generated.go"},
		{"internal/build/functions/data_change_message_handler/searchers.go.tmpl", "internal/build/functions/data_change_message_handler/searchers.generated.go"},
		{"internal/build/jobs/db_cleaner/build.go.tmpl", "internal/build/jobs/db_cleaner/build.generated.go"},
		{"internal/build/jobs/db_cleaner/config.go.tmpl", "internal/build/jobs/db_cleaner/config.generated.go"},
		{"internal/build/jobs/email_deliverability_test/build.go.tmpl", "internal/build/jobs/email_deliverability_test/build.generated.go"},
		{"internal/build/jobs/email_deliverability_test/config.go.tmpl", "internal/build/jobs/email_deliverability_test/config.generated.go"},
		{"internal/build/jobs/meal_plan_finalizer/build.go.tmpl", "internal/build/jobs/meal_plan_finalizer/build.generated.go"},
		{"internal/build/jobs/meal_plan_finalizer/config.go.tmpl", "internal/build/jobs/meal_plan_finalizer/config.generated.go"},
		{"internal/build/jobs/meal_plan_grocery_list_initializer/build.go.tmpl", "internal/build/jobs/meal_plan_grocery_list_initializer/build.generated.go"},
		{"internal/build/jobs/meal_plan_grocery_list_initializer/config.go.tmpl", "internal/build/jobs/meal_plan_grocery_list_initializer/config.generated.go"},
		{"internal/build/jobs/meal_plan_task_creator/build.go.tmpl", "internal/build/jobs/meal_plan_task_creator/build.generated.go"},
		{"internal/build/jobs/meal_plan_task_creator/config.go.tmpl", "internal/build/jobs/meal_plan_task_creator/config.generated.go"},
		{"internal/build/jobs/mobile_notification_scheduler/build.go.tmpl", "internal/build/jobs/mobile_notification_scheduler/build.generated.go"},
		{"internal/build/jobs/mobile_notification_scheduler/config.go.tmpl", "internal/build/jobs/mobile_notification_scheduler/config.generated.go"},
		{"internal/build/jobs/mobile_notification_scheduler/scheduler.go.tmpl", "internal/build/jobs/mobile_notification_scheduler/scheduler.generated.go"},
		{"internal/build/jobs/queue_test/build.go.tmpl", "internal/build/jobs/queue_test/build.generated.go"},
		{"internal/build/jobs/queue_test/config.go.tmpl", "internal/build/jobs/queue_test/config.generated.go"},
		{"internal/build/jobs/queue_test/result.go.tmpl", "internal/build/jobs/queue_test/result.generated.go"},
		{"internal/build/jobs/search_data_index_scheduler/build.go.tmpl", "internal/build/jobs/search_data_index_scheduler/build.generated.go"},
		{"internal/build/jobs/search_data_index_scheduler/config.go.tmpl", "internal/build/jobs/search_data_index_scheduler/config.generated.go"},
		{"internal/build/jobs/search_data_index_scheduler/indexers.go.tmpl", "internal/build/jobs/search_data_index_scheduler/indexers.generated.go"},
		{"internal/build/services/api/server.go.tmpl", "internal/build/services/api/server.generated.go"},
		{"internal/build/services/api/http/build.go.tmpl", "internal/build/services/api/http/build.generated.go"},
		{"internal/build/services/api/http/config.go.tmpl", "internal/build/services/api/http/config.generated.go"},
		{"internal/build/services/api/http/doc.go.tmpl", "internal/build/services/api/http/doc.generated.go"},
		{"internal/build/services/api/http/http_routes.go.tmpl", "internal/build/services/api/http/http_routes.generated.go"},
		{"internal/build/services/api/http/searchers.go.tmpl", "internal/build/services/api/http/searchers.generated.go"},
		{"internal/build/services/api/http/server.go.tmpl", "internal/build/services/api/http/server.generated.go"},
		{"internal/build/services/mcp/build.go.tmpl", "internal/build/services/mcp/build.generated.go"},
		{"internal/build/services/mcp/config.go.tmpl", "internal/build/services/mcp/config.generated.go"},
		{"internal/domain/auth/converters/password_reset_tokens.go.tmpl", "internal/domain/auth/converters/password_reset_tokens.generated.go"},
		{"internal/domain/auth/fakes/auth.go.tmpl", "internal/domain/auth/fakes/auth.generated.go"},
		{"internal/domain/auth/fakes/doc.go.tmpl", "internal/domain/auth/fakes/doc.generated.go"},
		{"internal/domain/auth/fakes/fake.go.tmpl", "internal/domain/auth/fakes/fake.generated.go"},
		{"internal/domain/auth/fakes/user.go.tmpl", "internal/domain/auth/fakes/user.generated.go"},
		{"internal/domain/auth/managers/auth_manager.go.tmpl", "internal/domain/auth/managers/auth_manager.generated.go"},
		{"internal/domain/auth/managers/do.go.tmpl", "internal/domain/auth/managers/do.generated.go"},
		{"internal/domain/auth/managers/interface.go.tmpl", "internal/domain/auth/managers/interface.generated.go"},
		{"internal/domain/auth/mock/auth_manager.go.tmpl", "internal/domain/auth/mock/auth_manager.generated.go"},
		{"internal/domain/identity/manager/do.go.tmpl", "internal/domain/identity/manager/do.generated.go"},
		{"internal/domain/identity/manager/interface.go.tmpl", "internal/domain/identity/manager/interface.generated.go"},
		{"internal/domain/identity/manager/user_data_manager.go.tmpl", "internal/domain/identity/manager/user_data_manager.generated.go"},
		{"internal/domain/identity/manager/mock/manager.go.tmpl", "internal/domain/identity/manager/mock/manager.generated.go"},
		{"internal/domain/identity/mock/repository.go.tmpl", "internal/domain/identity/mock/repository.generated.go"},
		{"internal/domain/internalops/internalops.go.tmpl", "internal/domain/internalops/internalops.generated.go"},
		{"internal/domain/internalops/queue_test_message.go.tmpl", "internal/domain/internalops/queue_test_message.generated.go"},
		{"internal/domain/internalops/mock/data_manager.go.tmpl", "internal/domain/internalops/mock/data_manager.generated.go"},
		{"internal/domain/mealplanning/bootstrap/enumerations.go.tmpl", "internal/domain/mealplanning/bootstrap/enumerations.generated.go"},
		{"internal/domain/mealplanning/bootstrap/helpers.go.tmpl", "internal/domain/mealplanning/bootstrap/helpers.generated.go"},
		{"internal/domain/mealplanning/bootstrap/meals.go.tmpl", "internal/domain/mealplanning/bootstrap/meals.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0001_pan_seared_butter_basted_steak.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0001_pan_seared_butter_basted_steak.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0002_sous_vide_chicken_breast.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0002_sous_vide_chicken_breast.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0003_whole_roast_chicken.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0003_whole_roast_chicken.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0004_sous_vide_pork_chops.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0004_sous_vide_pork_chops.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0005_cheeseburgers.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0005_cheeseburgers.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0006_simple_white_rice.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0006_simple_white_rice.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0007_ultra_fluffy_mashed_potatoes.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0007_ultra_fluffy_mashed_potatoes.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0008_caesar_roasted_broccoli.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0008_caesar_roasted_broccoli.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0009_haricots_verts_amandine.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0009_haricots_verts_amandine.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0010_mixed_green_salad.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0010_mixed_green_salad.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0011_soy_sauce_braised_chicken_thighs.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0011_soy_sauce_braised_chicken_thighs.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0012_grilled_pork_tenderloin.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0012_grilled_pork_tenderloin.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0013_pan_seared_salmon_fillets.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0013_pan_seared_salmon_fillets.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0014_roasted_brussels_sprouts.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0014_roasted_brussels_sprouts.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0015_refried_beans.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0015_refried_beans.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0016_carne_asada.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0016_carne_asada.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0017_butter_chicken.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0017_butter_chicken.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0017_garlic_parmesan_croutons.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0017_garlic_parmesan_croutons.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0018_caesar_dressing.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0018_caesar_dressing.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0018_stovetop_mac_and_cheese.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0018_stovetop_mac_and_cheese.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0019_caesar_salad.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0019_caesar_salad.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0020_glazed_carrots.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0020_glazed_carrots.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0021_cornbread.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0021_cornbread.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0021_teriyaki_sauce.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0021_teriyaki_sauce.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0022_grilled_whole_cauliflower.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0022_grilled_whole_cauliflower.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0023_stir_fried_green_beans.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0023_stir_fried_green_beans.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0024_tortillas.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0024_tortillas.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0025_caesar_breadcrumbs.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0025_caesar_breadcrumbs.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0026_gochujang_butter_pasta.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0026_gochujang_butter_pasta.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0027_one_pan_pasta.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0027_one_pan_pasta.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0028_peanut_butter_noodles.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0028_peanut_butter_noodles.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0029_chicken_vermicelli_soup.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0029_chicken_vermicelli_soup.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0030_chicken_florentine.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0030_chicken_florentine.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipe_0031_chana_masala.go.tmpl", "internal/domain/mealplanning/bootstrap/recipe_0031_chana_masala.generated.go"},
		{"internal/domain/mealplanning/bootstrap/recipes.go.tmpl", "internal/domain/mealplanning/bootstrap/recipes.generated.go"},
		{"internal/domain/mealplanning/emails/emails.go.tmpl", "internal/domain/mealplanning/emails/emails.generated.go"},
		{"internal/domain/mealplanning/grocerylistpreparation/do.go.tmpl", "internal/domain/mealplanning/grocerylistpreparation/do.generated.go"},
		{"internal/domain/mealplanning/grocerylistpreparation/list_creator.go.tmpl", "internal/domain/mealplanning/grocerylistpreparation/list_creator.generated.go"},
		{"internal/domain/mealplanning/grocerylistpreparation/mocks.go.tmpl", "internal/domain/mealplanning/grocerylistpreparation/mocks.generated.go"},
		{"internal/domain/mealplanning/managers/do.go.tmpl", "internal/domain/mealplanning/managers/do.generated.go"},
		{"internal/domain/mealplanning/managers/meal_planning_manager.go.tmpl", "internal/domain/mealplanning/managers/meal_planning_manager.generated.go"},
		{"internal/domain/mealplanning/managers/recipe_manager.go.tmpl", "internal/domain/mealplanning/managers/recipe_manager.generated.go"},
		{"internal/domain/mealplanning/managers/valid_enumerations_manager.go.tmpl", "internal/domain/mealplanning/managers/valid_enumerations_manager.generated.go"},
		{"internal/domain/mealplanning/managers/mock/meal_planning_manager.go.tmpl", "internal/domain/mealplanning/managers/mock/meal_planning_manager.generated.go"},
		{"internal/domain/mealplanning/managers/mock/recipe_manager.go.tmpl", "internal/domain/mealplanning/managers/mock/recipe_manager.generated.go"},
		{"internal/domain/mealplanning/managers/mock/valid_enumeration_manager.go.tmpl", "internal/domain/mealplanning/managers/mock/valid_enumeration_manager.generated.go"},
		{"internal/domain/mealplanning/mocks/repository.go.tmpl", "internal/domain/mealplanning/mocks/repository.generated.go"},
		{"internal/domain/mealplanning/notifications/constants.go.tmpl", "internal/domain/mealplanning/notifications/constants.generated.go"},
		{"internal/domain/mealplanning/privacy/collector.go.tmpl", "internal/domain/mealplanning/privacy/collector.generated.go"},
		{"internal/domain/mealplanning/recipeanalysis/do.go.tmpl", "internal/domain/mealplanning/recipeanalysis/do.generated.go"},
		{"internal/domain/mealplanning/recipeanalysis/mock_recipe_analyzer.go.tmpl", "internal/domain/mealplanning/recipeanalysis/mock_recipe_analyzer.generated.go"},
		{"internal/domain/mealplanning/recipeanalysis/recipe_analyzer.go.tmpl", "internal/domain/mealplanning/recipeanalysis/recipe_analyzer.generated.go"},
		{"internal/domain/mealplanning/recipevalidator/recipe_validator.go.tmpl", "internal/domain/mealplanning/recipevalidator/recipe_validator.generated.go"},
		{"internal/domain/mealplanning/registration/registration.go.tmpl", "internal/domain/mealplanning/registration/registration.generated.go"},
		{"internal/domain/oauth/converters/oauth2_clients.go.tmpl", "internal/domain/oauth/converters/oauth2_clients.generated.go"},
		{"internal/domain/oauth/fakes/doc.go.tmpl", "internal/domain/oauth/fakes/doc.generated.go"},
		{"internal/domain/oauth/fakes/fake.go.tmpl", "internal/domain/oauth/fakes/fake.generated.go"},
		{"internal/domain/oauth/fakes/oauth2_client.go.tmpl", "internal/domain/oauth/fakes/oauth2_client.generated.go"},
		{"internal/domain/oauth/manager/do.go.tmpl", "internal/domain/oauth/manager/do.generated.go"},
		{"internal/domain/oauth/manager/manager.go.tmpl", "internal/domain/oauth/manager/manager.generated.go"},
		{"internal/domain/oauth/manager/mock/manager.go.tmpl", "internal/domain/oauth/manager/mock/manager.generated.go"},
		{"internal/domain/oauth/mock/doc.go.tmpl", "internal/domain/oauth/mock/doc.generated.go"},
		{"internal/domain/oauth/mock/repository.go.tmpl", "internal/domain/oauth/mock/repository.generated.go"},
		{"internal/functions/datachangemessagehandler/data_change_message_handler.go.tmpl", "internal/functions/datachangemessagehandler/data_change_message_handler.generated.go"},
		{"internal/functions/datachangemessagehandler/data_changes.go.tmpl", "internal/functions/datachangemessagehandler/data_changes.generated.go"},
		{"internal/functions/datachangemessagehandler/do.go.tmpl", "internal/functions/datachangemessagehandler/do.generated.go"},
		{"internal/functions/datachangemessagehandler/identity_handlers.go.tmpl", "internal/functions/datachangemessagehandler/identity_handlers.generated.go"},
		{"internal/functions/datachangemessagehandler/mealplanning_handlers.go.tmpl", "internal/functions/datachangemessagehandler/mealplanning_handlers.generated.go"},
		{"internal/functions/datachangemessagehandler/mobile_notifications.go.tmpl", "internal/functions/datachangemessagehandler/mobile_notifications.generated.go"},
		{"internal/functions/datachangemessagehandler/outbound_emailer.go.tmpl", "internal/functions/datachangemessagehandler/outbound_emailer.generated.go"},
		{"internal/functions/datachangemessagehandler/queue_test_messages.go.tmpl", "internal/functions/datachangemessagehandler/queue_test_messages.generated.go"},
		{"internal/functions/datachangemessagehandler/search_index_requests.go.tmpl", "internal/functions/datachangemessagehandler/search_index_requests.generated.go"},
		{"internal/functions/datachangemessagehandler/user_data_aggregator.go.tmpl", "internal/functions/datachangemessagehandler/user_data_aggregator.generated.go"},
		{"internal/functions/datachangemessagehandler/webhook_executor.go.tmpl", "internal/functions/datachangemessagehandler/webhook_executor.generated.go"},
		{"internal/grpc/converters/helpers.go.tmpl", "internal/grpc/converters/helpers.generated.go"},
		{"internal/grpc/converters/query_filter.go.tmpl", "internal/grpc/converters/query_filter.generated.go"},
		{"internal/repositories/postgres/internalops/client.go.tmpl", "internal/repositories/postgres/internalops/client.generated.go"},
		{"internal/repositories/postgres/internalops/do.go.tmpl", "internal/repositories/postgres/internalops/do.generated.go"},
		{"internal/repositories/postgres/internalops/maintenance.go.tmpl", "internal/repositories/postgres/internalops/maintenance.generated.go"},
		{"internal/repositories/postgres/internalops/queue_test_messages.go.tmpl", "internal/repositories/postgres/internalops/queue_test_messages.generated.go"},
		{"internal/repositories/postgres/migrations/migrate.go.tmpl", "internal/repositories/postgres/migrations/migrate.generated.go"},
		{"internal/services/analytics/grpc/do.go.tmpl", "internal/services/analytics/grpc/do.generated.go"},
		{"internal/services/analytics/grpc/permissions.go.tmpl", "internal/services/analytics/grpc/permissions.generated.go"},
		{"internal/services/analytics/grpc/service.go.tmpl", "internal/services/analytics/grpc/service.generated.go"},
		{"internal/services/auth/grpc/auth.go.tmpl", "internal/services/auth/grpc/auth.generated.go"},
		{"internal/services/auth/grpc/do.go.tmpl", "internal/services/auth/grpc/do.generated.go"},
		{"internal/services/auth/grpc/passkey.go.tmpl", "internal/services/auth/grpc/passkey.generated.go"},
		{"internal/services/auth/grpc/permissions.go.tmpl", "internal/services/auth/grpc/permissions.generated.go"},
		{"internal/services/auth/grpc/service.go.tmpl", "internal/services/auth/grpc/service.generated.go"},
		{"internal/services/auth/grpc/converters/converters.go.tmpl", "internal/services/auth/grpc/converters/converters.generated.go"},
		{"internal/services/auth/grpc/interceptors/authn_interceptor.go.tmpl", "internal/services/auth/grpc/interceptors/authn_interceptor.generated.go"},
		{"internal/services/auth/grpc/interceptors/do.go.tmpl", "internal/services/auth/grpc/interceptors/do.generated.go"},
		{"internal/services/auth/handlers/passkey/cookie.go.tmpl", "internal/services/auth/handlers/passkey/cookie.generated.go"},
		{"internal/services/auth/handlers/passkey/handlers.go.tmpl", "internal/services/auth/handlers/passkey/handlers.generated.go"},
		{"internal/services/dataprivacy/config/config.go.tmpl", "internal/services/dataprivacy/config/config.generated.go"},
		{"internal/services/email/workers/email_deliverability_test/do.go.tmpl", "internal/services/email/workers/email_deliverability_test/do.generated.go"},
		{"internal/services/email/workers/email_deliverability_test/job.go.tmpl", "internal/services/email/workers/email_deliverability_test/job.generated.go"},
		{"internal/services/errors/grpc_mapper.go.tmpl", "internal/services/errors/grpc_mapper.generated.go"},
		{"internal/services/errors/http_mapper.go.tmpl", "internal/services/errors/http_mapper.generated.go"},
		{"internal/services/identity/config/config.go.tmpl", "internal/services/identity/config/config.generated.go"},
		{"internal/services/identity/emails/emails.go.tmpl", "internal/services/identity/emails/emails.generated.go"},
		{"internal/services/identity/grpc/account_invitations.go.tmpl", "internal/services/identity/grpc/account_invitations.generated.go"},
		{"internal/services/identity/grpc/accounts.go.tmpl", "internal/services/identity/grpc/accounts.generated.go"},
		{"internal/services/identity/grpc/admin.go.tmpl", "internal/services/identity/grpc/admin.generated.go"},
		{"internal/services/identity/grpc/do.go.tmpl", "internal/services/identity/grpc/do.generated.go"},
		{"internal/services/identity/grpc/permissions.go.tmpl", "internal/services/identity/grpc/permissions.generated.go"},
		{"internal/services/identity/grpc/service.go.tmpl", "internal/services/identity/grpc/service.generated.go"},
		{"internal/services/identity/grpc/users.go.tmpl", "internal/services/identity/grpc/users.generated.go"},
		{"internal/services/identity/grpc/converters/converters.go.tmpl", "internal/services/identity/grpc/converters/converters.generated.go"},
		{"internal/services/identity/indexing/do.go.tmpl", "internal/services/identity/indexing/do.generated.go"},
		{"internal/services/identity/indexing/doc.go.tmpl", "internal/services/identity/indexing/doc.generated.go"},
		{"internal/services/identity/indexing/indexing.go.tmpl", "internal/services/identity/indexing/indexing.generated.go"},
		{"internal/services/identity/indexing/search_data_index_scheduler.go.tmpl", "internal/services/identity/indexing/search_data_index_scheduler.generated.go"},
		{"internal/services/identity/indexing/search_subsets.go.tmpl", "internal/services/identity/indexing/search_subsets.generated.go"},
		{"internal/services/internalops/grpc/do.go.tmpl", "internal/services/internalops/grpc/do.generated.go"},
		{"internal/services/internalops/grpc/permissions.go.tmpl", "internal/services/internalops/grpc/permissions.generated.go"},
		{"internal/services/internalops/grpc/service.go.tmpl", "internal/services/internalops/grpc/service.generated.go"},
		{"internal/services/internalops/workers/queue_test/do.go.tmpl", "internal/services/internalops/workers/queue_test/do.generated.go"},
		{"internal/services/internalops/workers/queue_test/job.go.tmpl", "internal/services/internalops/workers/queue_test/job.generated.go"},
		{"internal/services/mealplanning/config/config.go.tmpl", "internal/services/mealplanning/config/config.generated.go"},
		{"internal/services/mealplanning/errors/grpc_mapper.go.tmpl", "internal/services/mealplanning/errors/grpc_mapper.generated.go"},
		{"internal/services/mealplanning/errors/http_mapper.go.tmpl", "internal/services/mealplanning/errors/http_mapper.generated.go"},
		{"internal/services/mealplanning/indexing/do.go.tmpl", "internal/services/mealplanning/indexing/do.generated.go"},
		{"internal/services/mealplanning/indexing/doc.go.tmpl", "internal/services/mealplanning/indexing/doc.generated.go"},
		{"internal/services/mealplanning/indexing/indexing.go.tmpl", "internal/services/mealplanning/indexing/indexing.generated.go"},
		{"internal/services/mealplanning/indexing/search_data_index_scheduler.go.tmpl", "internal/services/mealplanning/indexing/search_data_index_scheduler.generated.go"},
		{"internal/services/mealplanning/indexing/search_subsets.go.tmpl", "internal/services/mealplanning/indexing/search_subsets.generated.go"},
		{"internal/services/mealplanning/workers/mocks.go.tmpl", "internal/services/mealplanning/workers/mocks.generated.go"},
		{"internal/services/mealplanning/workers/worker.go.tmpl", "internal/services/mealplanning/workers/worker.generated.go"},
		{"internal/services/mealplanning/workers/meal_plan_finalizer/do.go.tmpl", "internal/services/mealplanning/workers/meal_plan_finalizer/do.generated.go"},
		{"internal/services/mealplanning/workers/meal_plan_finalizer/meal_plan_finalizer.go.tmpl", "internal/services/mealplanning/workers/meal_plan_finalizer/meal_plan_finalizer.generated.go"},
		{"internal/services/mealplanning/workers/meal_plan_grocery_list_initializer/do.go.tmpl", "internal/services/mealplanning/workers/meal_plan_grocery_list_initializer/do.generated.go"},
		{"internal/services/mealplanning/workers/meal_plan_grocery_list_initializer/meal_plan_grocery_list_initializer.go.tmpl", "internal/services/mealplanning/workers/meal_plan_grocery_list_initializer/meal_plan_grocery_list_initializer.generated.go"},
		{"internal/services/mealplanning/workers/meal_plan_task_creator/do.go.tmpl", "internal/services/mealplanning/workers/meal_plan_task_creator/do.generated.go"},
		{"internal/services/mealplanning/workers/meal_plan_task_creator/meal_plan_task_creator.go.tmpl", "internal/services/mealplanning/workers/meal_plan_task_creator/meal_plan_task_creator.generated.go"},
		{"internal/services/oauth/config/config.go.tmpl", "internal/services/oauth/config/config.generated.go"},
		{"internal/services/oauth/grpc/do.go.tmpl", "internal/services/oauth/grpc/do.generated.go"},
		{"internal/services/oauth/grpc/oauth2_clients.go.tmpl", "internal/services/oauth/grpc/oauth2_clients.generated.go"},
		{"internal/services/oauth/grpc/permissions.go.tmpl", "internal/services/oauth/grpc/permissions.generated.go"},
		{"internal/services/oauth/grpc/service.go.tmpl", "internal/services/oauth/grpc/service.generated.go"},
		{"internal/services/oauth/grpc/converters/converters.go.tmpl", "internal/services/oauth/grpc/converters/converters.generated.go"},
		{"internal/services/oauth/workers/db_cleaner/db_cleaner.go.tmpl", "internal/services/oauth/workers/db_cleaner/db_cleaner.generated.go"},
		{"internal/services/oauth/workers/db_cleaner/do.go.tmpl", "internal/services/oauth/workers/db_cleaner/do.generated.go"},
		{"internal/services/payments/adapters/do.go.tmpl", "internal/services/payments/adapters/do.generated.go"},
		{"internal/services/payments/adapters/revenuecat.go.tmpl", "internal/services/payments/adapters/revenuecat.generated.go"},
		{"internal/services/payments/adapters/stripe.go.tmpl", "internal/services/payments/adapters/stripe.generated.go"},
		{"internal/services/payments/adapters/stub.go.tmpl", "internal/services/payments/adapters/stub.generated.go"},
		{"internal/services/payments/config/config.go.tmpl", "internal/services/payments/config/config.generated.go"},
		{"internal/services/payments/http/do.go.tmpl", "internal/services/payments/http/do.generated.go"},
		{"internal/services/payments/http/webhook_handler.go.tmpl", "internal/services/payments/http/webhook_handler.generated.go"},
		{"internal/services/uploadedmedia/config/config.go.tmpl", "internal/services/uploadedmedia/config/config.generated.go"},
		{"internal/services/uploadedmedia/config/do.go.tmpl", "internal/services/uploadedmedia/config/do.generated.go"},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}


// planPhase6CodegenFiles ships the codegen output that consumer normally runs
// protoc/sqlc to produce. 193 files, ~146K LOC across internal/grpc/generated/*
// (14 proto-derived packages) and internal/repositories/postgres/*/generated
// (19 sqlc-derived packages). Verbatim from target modulo module paths. This
// turns `make build` into a one-step operation: consumers don't need protoc or
// sqlc installed just to compile the generated tree. Regenerating these is
// still possible via the embedded sqlc.yaml and per-domain proto templates.
func (p *Pipeline) planPhase6CodegenFiles() []PlannedFile {
	ctx := TemplateContext{Project: p.config}

	type mapping struct {
		template string
		output   string
	}

	mappings := []mapping{
		{"internal/grpc/generated/filtering/filtering.pb.go.tmpl", "internal/grpc/generated/filtering/filtering.pb.generated.go"},
		{"internal/grpc/generated/services/analytics/analytics_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/analytics/analytics_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/analytics/analytics_service_types.pb.go.tmpl", "internal/grpc/generated/services/analytics/analytics_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/analytics/analytics_service.pb.go.tmpl", "internal/grpc/generated/services/analytics/analytics_service.pb.generated.go"},
		{"internal/grpc/generated/services/audit/audit_messages.pb.go.tmpl", "internal/grpc/generated/services/audit/audit_messages.pb.generated.go"},
		{"internal/grpc/generated/services/audit/audit_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/audit/audit_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/audit/audit_service_types.pb.go.tmpl", "internal/grpc/generated/services/audit/audit_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/audit/audit_service.pb.go.tmpl", "internal/grpc/generated/services/audit/audit_service.pb.generated.go"},
		{"internal/grpc/generated/services/auth/auth_messages.pb.go.tmpl", "internal/grpc/generated/services/auth/auth_messages.pb.generated.go"},
		{"internal/grpc/generated/services/auth/auth_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/auth/auth_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/auth/auth_service_types.pb.go.tmpl", "internal/grpc/generated/services/auth/auth_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/auth/auth_service.pb.go.tmpl", "internal/grpc/generated/services/auth/auth_service.pb.generated.go"},
		{"internal/grpc/generated/services/comments/comments_messages.pb.go.tmpl", "internal/grpc/generated/services/comments/comments_messages.pb.generated.go"},
		{"internal/grpc/generated/services/comments/comments_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/comments/comments_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/comments/comments_service_types.pb.go.tmpl", "internal/grpc/generated/services/comments/comments_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/comments/comments_service.pb.go.tmpl", "internal/grpc/generated/services/comments/comments_service.pb.generated.go"},
		{"internal/grpc/generated/services/dataprivacy/dataprivacy_messages.pb.go.tmpl", "internal/grpc/generated/services/dataprivacy/dataprivacy_messages.pb.generated.go"},
		{"internal/grpc/generated/services/dataprivacy/dataprivacy_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/dataprivacy/dataprivacy_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/dataprivacy/dataprivacy_service_types.pb.go.tmpl", "internal/grpc/generated/services/dataprivacy/dataprivacy_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/dataprivacy/dataprivacy_service.pb.go.tmpl", "internal/grpc/generated/services/dataprivacy/dataprivacy_service.pb.generated.go"},
		{"internal/grpc/generated/services/identity/identity_messages.pb.go.tmpl", "internal/grpc/generated/services/identity/identity_messages.pb.generated.go"},
		{"internal/grpc/generated/services/identity/identity_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/identity/identity_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/identity/identity_service_types.pb.go.tmpl", "internal/grpc/generated/services/identity/identity_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/identity/identity_service.pb.go.tmpl", "internal/grpc/generated/services/identity/identity_service.pb.generated.go"},
		{"internal/grpc/generated/services/internalops/internal_ops_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/internalops/internal_ops_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/internalops/internal_ops_service_types.pb.go.tmpl", "internal/grpc/generated/services/internalops/internal_ops_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/internalops/internal_ops_service.pb.go.tmpl", "internal/grpc/generated/services/internalops/internal_ops_service.pb.generated.go"},
		{"internal/grpc/generated/services/issuereports/issue_reports_messages.pb.go.tmpl", "internal/grpc/generated/services/issuereports/issue_reports_messages.pb.generated.go"},
		{"internal/grpc/generated/services/issuereports/issue_reports_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/issuereports/issue_reports_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/issuereports/issue_reports_service_types.pb.go.tmpl", "internal/grpc/generated/services/issuereports/issue_reports_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/issuereports/issue_reports_service.pb.go.tmpl", "internal/grpc/generated/services/issuereports/issue_reports_service.pb.generated.go"},
		{"internal/grpc/generated/services/mealplanning/mealplanning_messages.pb.go.tmpl", "internal/grpc/generated/services/mealplanning/mealplanning_messages.pb.generated.go"},
		{"internal/grpc/generated/services/mealplanning/mealplanning_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/mealplanning/mealplanning_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/mealplanning/mealplanning_service_types.pb.go.tmpl", "internal/grpc/generated/services/mealplanning/mealplanning_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/mealplanning/mealplanning_service.pb.go.tmpl", "internal/grpc/generated/services/mealplanning/mealplanning_service.pb.generated.go"},
		{"internal/grpc/generated/services/notifications/notifications_messages.pb.go.tmpl", "internal/grpc/generated/services/notifications/notifications_messages.pb.generated.go"},
		{"internal/grpc/generated/services/notifications/notifications_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/notifications/notifications_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/notifications/notifications_service_types.pb.go.tmpl", "internal/grpc/generated/services/notifications/notifications_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/notifications/notifications_service.pb.go.tmpl", "internal/grpc/generated/services/notifications/notifications_service.pb.generated.go"},
		{"internal/grpc/generated/services/oauth/oauth_messages.pb.go.tmpl", "internal/grpc/generated/services/oauth/oauth_messages.pb.generated.go"},
		{"internal/grpc/generated/services/oauth/oauth_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/oauth/oauth_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/oauth/oauth_service_types.pb.go.tmpl", "internal/grpc/generated/services/oauth/oauth_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/oauth/oauth_service.pb.go.tmpl", "internal/grpc/generated/services/oauth/oauth_service.pb.generated.go"},
		{"internal/grpc/generated/services/payments/payments_messages.pb.go.tmpl", "internal/grpc/generated/services/payments/payments_messages.pb.generated.go"},
		{"internal/grpc/generated/services/payments/payments_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/payments/payments_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/payments/payments_service_types.pb.go.tmpl", "internal/grpc/generated/services/payments/payments_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/payments/payments_service.pb.go.tmpl", "internal/grpc/generated/services/payments/payments_service.pb.generated.go"},
		{"internal/grpc/generated/services/settings/settings_messages.pb.go.tmpl", "internal/grpc/generated/services/settings/settings_messages.pb.generated.go"},
		{"internal/grpc/generated/services/settings/settings_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/settings/settings_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/settings/settings_service_types.pb.go.tmpl", "internal/grpc/generated/services/settings/settings_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/settings/settings_service.pb.go.tmpl", "internal/grpc/generated/services/settings/settings_service.pb.generated.go"},
		{"internal/grpc/generated/services/uploadedmedia/uploaded_media_messages.pb.go.tmpl", "internal/grpc/generated/services/uploadedmedia/uploaded_media_messages.pb.generated.go"},
		{"internal/grpc/generated/services/uploadedmedia/uploaded_media_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/uploadedmedia/uploaded_media_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/uploadedmedia/uploaded_media_service_types.pb.go.tmpl", "internal/grpc/generated/services/uploadedmedia/uploaded_media_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/uploadedmedia/uploaded_media_service.pb.go.tmpl", "internal/grpc/generated/services/uploadedmedia/uploaded_media_service.pb.generated.go"},
		{"internal/grpc/generated/services/waitlists/waitlists_messages.pb.go.tmpl", "internal/grpc/generated/services/waitlists/waitlists_messages.pb.generated.go"},
		{"internal/grpc/generated/services/waitlists/waitlists_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/waitlists/waitlists_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/waitlists/waitlists_service_types.pb.go.tmpl", "internal/grpc/generated/services/waitlists/waitlists_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/waitlists/waitlists_service.pb.go.tmpl", "internal/grpc/generated/services/waitlists/waitlists_service.pb.generated.go"},
		{"internal/grpc/generated/services/webhooks/webhooks_messages.pb.go.tmpl", "internal/grpc/generated/services/webhooks/webhooks_messages.pb.generated.go"},
		{"internal/grpc/generated/services/webhooks/webhooks_service_grpc.pb.go.tmpl", "internal/grpc/generated/services/webhooks/webhooks_service_grpc.pb.generated.go"},
		{"internal/grpc/generated/services/webhooks/webhooks_service_types.pb.go.tmpl", "internal/grpc/generated/services/webhooks/webhooks_service_types.pb.generated.go"},
		{"internal/grpc/generated/services/webhooks/webhooks_service.pb.go.tmpl", "internal/grpc/generated/services/webhooks/webhooks_service.pb.generated.go"},
		{"internal/grpc/generated/types/common.pb.go.tmpl", "internal/grpc/generated/types/common.pb.generated.go"},
		{"internal/repositories/postgres/audit/generated/audit_logs.generated.sql_generated.go.tmpl", "internal/repositories/postgres/audit/generated/audit_logs.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/audit/generated/db_generated.go.tmpl", "internal/repositories/postgres/audit/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/audit/generated/models_generated.go.tmpl", "internal/repositories/postgres/audit/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/audit/generated/querier_generated.go.tmpl", "internal/repositories/postgres/audit/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/auth/generated/db_generated.go.tmpl", "internal/repositories/postgres/auth/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/auth/generated/models_generated.go.tmpl", "internal/repositories/postgres/auth/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/auth/generated/password_reset_tokens.generated.sql_generated.go.tmpl", "internal/repositories/postgres/auth/generated/password_reset_tokens.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/auth/generated/querier_generated.go.tmpl", "internal/repositories/postgres/auth/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/auth/generated/user_sessions.generated.sql_generated.go.tmpl", "internal/repositories/postgres/auth/generated/user_sessions.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/comments/generated/comments.generated.sql_generated.go.tmpl", "internal/repositories/postgres/comments/generated/comments.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/comments/generated/db_generated.go.tmpl", "internal/repositories/postgres/comments/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/comments/generated/models_generated.go.tmpl", "internal/repositories/postgres/comments/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/comments/generated/querier_generated.go.tmpl", "internal/repositories/postgres/comments/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/dataprivacy/generated/db_generated.go.tmpl", "internal/repositories/postgres/dataprivacy/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/dataprivacy/generated/models_generated.go.tmpl", "internal/repositories/postgres/dataprivacy/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/dataprivacy/generated/querier_generated.go.tmpl", "internal/repositories/postgres/dataprivacy/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/dataprivacy/generated/user_data_disclosures.generated.sql_generated.go.tmpl", "internal/repositories/postgres/dataprivacy/generated/user_data_disclosures.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/account_invitations.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/account_invitations.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/account_user_memberships.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/account_user_memberships.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/accounts.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/accounts.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/admin.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/admin.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/db_generated.go.tmpl", "internal/repositories/postgres/identity/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/models_generated.go.tmpl", "internal/repositories/postgres/identity/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/permissions.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/permissions.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/querier_generated.go.tmpl", "internal/repositories/postgres/identity/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/user_avatars.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/user_avatars.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/user_role_assignments.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/user_role_assignments.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/user_role_hierarchy.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/user_role_hierarchy.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/user_role_permissions.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/user_role_permissions.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/user_roles.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/user_roles.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/users.generated.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/users.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/identity/generated/webauthn_credentials.sql_generated.go.tmpl", "internal/repositories/postgres/identity/generated/webauthn_credentials.sql_generated.generated.go"},
		{"internal/repositories/postgres/internalops/generated/db_generated.go.tmpl", "internal/repositories/postgres/internalops/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/internalops/generated/internalops.generated.sql_generated.go.tmpl", "internal/repositories/postgres/internalops/generated/internalops.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/internalops/generated/models_generated.go.tmpl", "internal/repositories/postgres/internalops/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/internalops/generated/querier_generated.go.tmpl", "internal/repositories/postgres/internalops/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/issuereports/generated/db_generated.go.tmpl", "internal/repositories/postgres/issuereports/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/issuereports/generated/issue_reports.generated.sql_generated.go.tmpl", "internal/repositories/postgres/issuereports/generated/issue_reports.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/issuereports/generated/models_generated.go.tmpl", "internal/repositories/postgres/issuereports/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/issuereports/generated/querier_generated.go.tmpl", "internal/repositories/postgres/issuereports/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/account_instrument_ownerships.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/account_instrument_ownerships.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/db_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/ingredient_media.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/ingredient_media.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_components.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_components.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_images.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_images.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_list_items.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_list_items.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_lists.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_lists.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plan_events.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plan_events.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plan_grocery_list_items.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plan_grocery_list_items.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plan_option_votes.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plan_option_votes.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plan_options.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plan_options.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plan_recipe_option_selections.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plan_recipe_option_selections.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plan_task_notification_context.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plan_task_notification_context.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plan_tasks.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plan_tasks.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meal_plans.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meal_plans.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/meals.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/meals.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/models_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/preparation_media.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/preparation_media.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/querier_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_images.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_images.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_list_items.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_list_items.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_lists.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_lists.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_media.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_media.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_prep_task_steps.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_prep_task_steps.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_prep_tasks.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_prep_tasks.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_ratings.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_ratings.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_step_completion_condition_ingredients.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_step_completion_condition_ingredients.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_step_completion_conditions.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_step_completion_conditions.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_step_images.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_step_images.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_step_ingredients.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_step_ingredients.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_step_instruments.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_step_instruments.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_step_products.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_step_products.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_step_vessels.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_step_vessels.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipe_steps.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipe_steps.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/recipes.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/recipes.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/uploaded_media_fetch.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/uploaded_media_fetch.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/user_ingredient_preferences.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/user_ingredient_preferences.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_ingredient_groups.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_ingredient_groups.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_ingredient_measurement_units.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_ingredient_measurement_units.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_ingredient_preparations.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_ingredient_preparations.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_ingredient_state_ingredients.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_ingredient_state_ingredients.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_ingredient_states.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_ingredient_states.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_ingredients.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_ingredients.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_instruments.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_instruments.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_measurement_unit_conversions.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_measurement_unit_conversions.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_measurement_units.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_measurement_units.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_prep_task_configs.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_prep_task_configs.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_preparation_instruments.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_preparation_instruments.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_preparation_vessels.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_preparation_vessels.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_preparations.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_preparations.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/mealplanning/generated/valid_vessels.generated.sql_generated.go.tmpl", "internal/repositories/postgres/mealplanning/generated/valid_vessels.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/notifications/generated/db_generated.go.tmpl", "internal/repositories/postgres/notifications/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/notifications/generated/models_generated.go.tmpl", "internal/repositories/postgres/notifications/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/notifications/generated/querier_generated.go.tmpl", "internal/repositories/postgres/notifications/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/notifications/generated/user_device_tokens.sql_generated.go.tmpl", "internal/repositories/postgres/notifications/generated/user_device_tokens.sql_generated.generated.go"},
		{"internal/repositories/postgres/notifications/generated/user_notifications.generated.sql_generated.go.tmpl", "internal/repositories/postgres/notifications/generated/user_notifications.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/oauth/generated/db_generated.go.tmpl", "internal/repositories/postgres/oauth/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/oauth/generated/models_generated.go.tmpl", "internal/repositories/postgres/oauth/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/oauth/generated/oauth2_client_tokens.generated.sql_generated.go.tmpl", "internal/repositories/postgres/oauth/generated/oauth2_client_tokens.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/oauth/generated/oauth2_clients.generated.sql_generated.go.tmpl", "internal/repositories/postgres/oauth/generated/oauth2_clients.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/oauth/generated/querier_generated.go.tmpl", "internal/repositories/postgres/oauth/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/payments/generated/db_generated.go.tmpl", "internal/repositories/postgres/payments/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/payments/generated/models_generated.go.tmpl", "internal/repositories/postgres/payments/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/payments/generated/payment_transactions.generated.sql_generated.go.tmpl", "internal/repositories/postgres/payments/generated/payment_transactions.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/payments/generated/products.generated.sql_generated.go.tmpl", "internal/repositories/postgres/payments/generated/products.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/payments/generated/purchases.generated.sql_generated.go.tmpl", "internal/repositories/postgres/payments/generated/purchases.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/payments/generated/querier_generated.go.tmpl", "internal/repositories/postgres/payments/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/payments/generated/subscriptions.generated.sql_generated.go.tmpl", "internal/repositories/postgres/payments/generated/subscriptions.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/settings/generated/db_generated.go.tmpl", "internal/repositories/postgres/settings/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/settings/generated/models_generated.go.tmpl", "internal/repositories/postgres/settings/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/settings/generated/querier_generated.go.tmpl", "internal/repositories/postgres/settings/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/settings/generated/service_setting_configurations.generated.sql_generated.go.tmpl", "internal/repositories/postgres/settings/generated/service_setting_configurations.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/settings/generated/service_settings.generated.sql_generated.go.tmpl", "internal/repositories/postgres/settings/generated/service_settings.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/uploadedmedia/generated/db_generated.go.tmpl", "internal/repositories/postgres/uploadedmedia/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/uploadedmedia/generated/models_generated.go.tmpl", "internal/repositories/postgres/uploadedmedia/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/uploadedmedia/generated/querier_generated.go.tmpl", "internal/repositories/postgres/uploadedmedia/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/uploadedmedia/generated/uploaded_media.generated.sql_generated.go.tmpl", "internal/repositories/postgres/uploadedmedia/generated/uploaded_media.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/waitlists/generated/db_generated.go.tmpl", "internal/repositories/postgres/waitlists/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/waitlists/generated/issue_reports.sql_generated.go.tmpl", "internal/repositories/postgres/waitlists/generated/issue_reports.sql_generated.generated.go"},
		{"internal/repositories/postgres/waitlists/generated/models_generated.go.tmpl", "internal/repositories/postgres/waitlists/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/waitlists/generated/querier_generated.go.tmpl", "internal/repositories/postgres/waitlists/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/waitlists/generated/waitlist_signups.generated.sql_generated.go.tmpl", "internal/repositories/postgres/waitlists/generated/waitlist_signups.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/waitlists/generated/waitlists.generated.sql_generated.go.tmpl", "internal/repositories/postgres/waitlists/generated/waitlists.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/webhooks/generated/db_generated.go.tmpl", "internal/repositories/postgres/webhooks/generated/db_generated.generated.go"},
		{"internal/repositories/postgres/webhooks/generated/models_generated.go.tmpl", "internal/repositories/postgres/webhooks/generated/models_generated.generated.go"},
		{"internal/repositories/postgres/webhooks/generated/querier_generated.go.tmpl", "internal/repositories/postgres/webhooks/generated/querier_generated.generated.go"},
		{"internal/repositories/postgres/webhooks/generated/webhook_trigger_configs.generated.sql_generated.go.tmpl", "internal/repositories/postgres/webhooks/generated/webhook_trigger_configs.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/webhooks/generated/webhook_trigger_events.generated.sql_generated.go.tmpl", "internal/repositories/postgres/webhooks/generated/webhook_trigger_events.generated.sql_generated.generated.go"},
		{"internal/repositories/postgres/webhooks/generated/webhooks.generated.sql_generated.go.tmpl", "internal/repositories/postgres/webhooks/generated/webhooks.generated.sql_generated.generated.go"},

	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, m := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
			IsGo:         true,
		})
	}
	return files
}

// planCmdFiles returns the cmd/ binaries that NAFF emits. Modeled on the
// dinnerdonebetter/backend/cmd/ tree:
//   - services/api: cobra root + grpc/http sub-binaries
//   - services/mcp: MCP server scaffold (auth, oauth2, schema)
//   - tools/codegen: queries, configs, valid_env_vars
//   - tools: bootstrap, migrate, encryptor, data_exporter, data_importer,
//     search_index_initializer, push_tester, aiagent
//   - workers: 8 cron-style job binaries
//   - functions/async_message_handler: pubsub consumer
//   - localdev/server: in-process all-in-one bootstrap
//   - playground: id-generator dev utility
//
// The codegen/queries main.go template needs every domain (built-in + user-defined)
// for its dispatch map, so it gets a project context with allDomains baked in.
func (p *Pipeline) planCmdFiles(allDomains []config.Domain) []PlannedFile {
	projectCtx := TemplateContext{Project: p.config}

	codegenCtx := TemplateContext{
		Project: &config.Project{
			ProjectMeta: p.config.ProjectMeta,
			Features:    p.config.Features,
			Targets:     p.config.Targets,
			Domains:     allDomains,
		},
	}

	type m struct {
		template string
		output   string
		isGo     bool
		data     TemplateContext
	}

	mappings := []m{
		// services/api root + grpc + http
		{"cmd/services/api/main.go.tmpl", "cmd/services/api/main.go", true, projectCtx},
		{"cmd/services/api/doc.go.tmpl", "cmd/services/api/doc.go", true, projectCtx},
		{"cmd/services/api/grpc/main.go.tmpl", "cmd/services/api/grpc/main.go", true, projectCtx},
		{"cmd/services/api/grpc/doc.go.tmpl", "cmd/services/api/grpc/doc.go", true, projectCtx},
		{"cmd/services/api/http/main.go.tmpl", "cmd/services/api/http/main.go", true, projectCtx},
		{"cmd/services/api/http/doc.go.tmpl", "cmd/services/api/http/doc.go", true, projectCtx},

		// services/mcp
		{"cmd/services/mcp/main.go.tmpl", "cmd/services/mcp/main.go", true, projectCtx},
		{"cmd/services/mcp/auth.go.tmpl", "cmd/services/mcp/auth.go", true, projectCtx},
		{"cmd/services/mcp/oauth2_handler.go.tmpl", "cmd/services/mcp/oauth2_handler.go", true, projectCtx},
		{"cmd/services/mcp/schema.go.tmpl", "cmd/services/mcp/schema.go", true, projectCtx},

		// tools/codegen/queries — main.go drives a per-entity dispatch map
		{"cmd/tools/codegen/queries/main.go.tmpl", "cmd/tools/codegen/queries/main.go", true, codegenCtx},
		{"cmd/tools/codegen/queries/helpers.go.tmpl", "cmd/tools/codegen/queries/helpers.go", true, projectCtx},
		{"cmd/tools/codegen/queries/sqlc.go.tmpl", "cmd/tools/codegen/queries/sqlc.go", true, projectCtx},

		// tools/codegen/configs
		{"cmd/tools/codegen/configs/main.go.tmpl", "cmd/tools/codegen/configs/main.go", true, projectCtx},
		{"cmd/tools/codegen/configs/doc.go.tmpl", "cmd/tools/codegen/configs/doc.go", true, projectCtx},
		{"cmd/tools/codegen/configs/utils.go.tmpl", "cmd/tools/codegen/configs/utils.go", true, projectCtx},
		{"cmd/tools/codegen/configs/localdev.go.tmpl", "cmd/tools/codegen/configs/localdev.go", true, projectCtx},
		{"cmd/tools/codegen/configs/integrationtests.go.tmpl", "cmd/tools/codegen/configs/integrationtests.go", true, projectCtx},
		{"cmd/tools/codegen/configs/prod.go.tmpl", "cmd/tools/codegen/configs/prod.go", true, projectCtx},

		// tools/codegen/valid_env_vars
		{"cmd/tools/codegen/valid_env_vars/main.go.tmpl", "cmd/tools/codegen/valid_env_vars/main.go", true, projectCtx},

		// tools/* (one main.go per tool)
		{"cmd/tools/bootstrap/main.go.tmpl", "cmd/tools/bootstrap/main.go", true, projectCtx},
		{"cmd/tools/migrate/main.go.tmpl", "cmd/tools/migrate/main.go", true, projectCtx},
		{"cmd/tools/encryptor/main.go.tmpl", "cmd/tools/encryptor/main.go", true, projectCtx},
		{"cmd/tools/data_exporter/main.go.tmpl", "cmd/tools/data_exporter/main.go", true, projectCtx},
		{"cmd/tools/data_importer/main.go.tmpl", "cmd/tools/data_importer/main.go", true, projectCtx},
		{"cmd/tools/search_index_initializer/main.go.tmpl", "cmd/tools/search_index_initializer/main.go", true, projectCtx},
		{"cmd/tools/push_tester/main.go.tmpl", "cmd/tools/push_tester/main.go", true, projectCtx},
		{"cmd/tools/aiagent/main.go.tmpl", "cmd/tools/aiagent/main.go", true, projectCtx},

		// workers/*
		{"cmd/workers/db_cleaner/main.go.tmpl", "cmd/workers/db_cleaner/main.go", true, projectCtx},
		{"cmd/workers/email_deliverability_test/main.go.tmpl", "cmd/workers/email_deliverability_test/main.go", true, projectCtx},
		{"cmd/workers/meal_plan_finalizer/main.go.tmpl", "cmd/workers/meal_plan_finalizer/main.go", true, projectCtx},
		{"cmd/workers/meal_plan_grocery_list_initializer/main.go.tmpl", "cmd/workers/meal_plan_grocery_list_initializer/main.go", true, projectCtx},
		{"cmd/workers/meal_plan_task_creator/main.go.tmpl", "cmd/workers/meal_plan_task_creator/main.go", true, projectCtx},
		{"cmd/workers/mobile_notification_scheduler/main.go.tmpl", "cmd/workers/mobile_notification_scheduler/main.go", true, projectCtx},
		{"cmd/workers/search_data_index_scheduler/main.go.tmpl", "cmd/workers/search_data_index_scheduler/main.go", true, projectCtx},
		{"cmd/workers/queue_test/main.go.tmpl", "cmd/workers/queue_test/main.go", true, projectCtx},

		// functions
		{"cmd/functions/async_message_handler/main.go.tmpl", "cmd/functions/async_message_handler/main.go", true, projectCtx},

		// localdev
		{"cmd/localdev/server/main.go.tmpl", "cmd/localdev/server/main.go", true, projectCtx},

		// playground
		{"cmd/playground/main.go.tmpl", "cmd/playground/main.go", true, projectCtx},
	}

	files := make([]PlannedFile, 0, len(mappings))
	for _, mp := range mappings {
		files = append(files, PlannedFile{
			TemplatePath: mp.template,
			OutputPath:   mp.output,
			Data:         mp.data,
			IsGo:         mp.isGo,
		})
	}
	return files
}

// planPerEntityFiles returns files generated once per entity.
func (p *Pipeline) planPerEntityFiles(ctx TemplateContext) []PlannedFile {
	domain := ctx.Domain.Name
	entitySnake := strings.ToLower(ctx.Entity.Name)
	entityPluralSnake := naming.New(ctx.Entity.Name).PluralSnake()

	templateMappings := []struct {
		template string
		output   string
		isGo     bool
	}{
		// Domain entity types
		{"domain/entity.go.tmpl", fmt.Sprintf("internal/domain/%s/%s.generated.go", domain, entitySnake), true},
		{"domain_converters/converters.go.tmpl", fmt.Sprintf("internal/domain/%s/converters/%s.generated.go", domain, entitySnake), true},

		// Repository entity CRUD
		{"repository/entity.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/%s.generated.go", domain, entitySnake), true},

		// Query codegen — filename matches dinnerdonebetter pattern: <domain>_<plural_snake>.go
		{"codegen/queries.go.tmpl", fmt.Sprintf("cmd/tools/codegen/queries/%s_%s.generated.go", domain, entityPluralSnake), true},

		// Migration (one per entity)
		{"migrations/migration.sql.tmpl", fmt.Sprintf("internal/repositories/postgres/migrations/migration_files/%s_%s.generated.sql", domain, entitySnake), false},

		// Proto
		{"proto/messages.proto.tmpl", fmt.Sprintf("proto/%s/%s_messages.generated.proto", domain, entitySnake), false},
		{"proto/service.proto.tmpl", fmt.Sprintf("proto/%s/%s_service.generated.proto", domain, entitySnake), false},

		// gRPC entity handlers
		{"grpc/entity.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/%s.generated.go", domain, entitySnake), true},
	}

	var files []PlannedFile
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

// iosBasePath returns the base output path for iOS files.
func (p *Pipeline) iosBasePath() string {
	moduleName := p.config.ProjectMeta.IOSModuleName
	return fmt.Sprintf("ios/%s", moduleName)
}

// planIOSProjectFiles returns iOS project-level files generated once per project.
func (p *Pipeline) planIOSProjectFiles() []PlannedFile {
	ctx := TemplateContext{
		Project: p.config,
	}

	base := p.iosBasePath()
	src := fmt.Sprintf("%s/Sources/%s", base, p.config.ProjectMeta.IOSModuleName)

	type mapping struct {
		template string
		output   string
	}

	templateMappings := []mapping{
		{"ios/project/Package.swift.tmpl", fmt.Sprintf("%s/Package.swift", base)},
		{"ios/project/App.swift.tmpl", fmt.Sprintf("%s/App.generated.swift", src)},
		{"ios/project/ContentView.swift.tmpl", fmt.Sprintf("%s/ContentView.generated.swift", src)},
		{"ios/project/Configuration.swift.tmpl", fmt.Sprintf("%s/Configuration.generated.swift", src)},
		{"ios/api/GRPCClient.swift.tmpl", fmt.Sprintf("%s/API/GRPCClient.generated.swift", src)},
		{"ios/auth/AuthManager.swift.tmpl", fmt.Sprintf("%s/Auth/AuthManager.generated.swift", src)},
		{"ios/auth/LoginView.swift.tmpl", fmt.Sprintf("%s/Auth/LoginView.generated.swift", src)},
		{"ios/auth/RegistrationView.swift.tmpl", fmt.Sprintf("%s/Auth/RegistrationView.generated.swift", src)},
	}

	var files []PlannedFile
	for _, m := range templateMappings {
		files = append(files, PlannedFile{
			TemplatePath: m.template,
			OutputPath:   m.output,
			Data:         ctx,
		})
	}
	return files
}

// planIOSPerDomainFiles returns iOS files generated once per domain.
func (p *Pipeline) planIOSPerDomainFiles(ctx TemplateContext) []PlannedFile {
	src := fmt.Sprintf("%s/Sources/%s", p.iosBasePath(), p.config.ProjectMeta.IOSModuleName)
	domain := ctx.Domain.Name

	return []PlannedFile{
		{
			TemplatePath: "ios/api/entity_service.swift.tmpl",
			OutputPath:   fmt.Sprintf("%s/API/%sService.generated.swift", src, strings.Title(domain)), //nolint:staticcheck
			Data:         ctx,
		},
	}
}

// planIOSPerEntityFiles returns iOS files generated once per entity.
func (p *Pipeline) planIOSPerEntityFiles(ctx TemplateContext) []PlannedFile {
	src := fmt.Sprintf("%s/Sources/%s", p.iosBasePath(), p.config.ProjectMeta.IOSModuleName)
	entityName := ctx.Entity.Name

	return []PlannedFile{
		{
			TemplatePath: "ios/models/entity.swift.tmpl",
			OutputPath:   fmt.Sprintf("%s/Models/%s.generated.swift", src, entityName),
			Data:         ctx,
		},
		{
			TemplatePath: "ios/models/entity_input.swift.tmpl",
			OutputPath:   fmt.Sprintf("%s/Models/%sInput.generated.swift", src, entityName),
			Data:         ctx,
		},
		{
			TemplatePath: "proto/messages.proto.tmpl",
			OutputPath:   fmt.Sprintf("%s/Proto/%s/%s_messages.generated.proto", p.iosBasePath(), ctx.Domain.Name, strings.ToLower(entityName)),
			Data:         ctx,
		},
		{
			TemplatePath: "proto/service.proto.tmpl",
			OutputPath:   fmt.Sprintf("%s/Proto/%s/%s_service.generated.proto", p.iosBasePath(), ctx.Domain.Name, strings.ToLower(entityName)),
			Data:         ctx,
		},
	}
}

// planConsumerAppFiles returns all files for the consumer SvelteKit app.
func (p *Pipeline) planConsumerAppFiles(allDomains []config.Domain) []PlannedFile {
	var files []PlannedFile

	projectCtx := TemplateContext{Project: p.config}

	// Project scaffolding.
	for _, m := range []struct{ template, output string }{
		{"frontend/consumer/project/package.json.tmpl", "frontend/consumer/package.json"},
		{"frontend/consumer/project/svelte.config.js.tmpl", "frontend/consumer/svelte.config.js"},
		{"frontend/consumer/project/tsconfig.json.tmpl", "frontend/consumer/tsconfig.json"},
		{"frontend/consumer/project/vite.config.ts.tmpl", "frontend/consumer/vite.config.ts"},
		{"frontend/consumer/project/app.html.tmpl", "frontend/consumer/src/app.html"},
		{"frontend/consumer/project/app.d.ts.tmpl", "frontend/consumer/src/app.d.ts"},
	} {
		files = append(files, PlannedFile{TemplatePath: m.template, OutputPath: m.output, Data: projectCtx})
	}

	// Shared library files.
	for _, m := range []struct{ template, output string }{
		{"frontend/shared/api_client.ts.tmpl", "frontend/consumer/src/lib/api/client.generated.ts"},
		{"frontend/shared/types.ts.tmpl", "frontend/consumer/src/lib/api/types.generated.ts"},
	} {
		files = append(files, PlannedFile{TemplatePath: m.template, OutputPath: m.output, Data: projectCtx})
	}

	// Layout files (need full project context for nav).
	// Build a project context with all domains for navigation.
	navCtx := TemplateContext{
		Project: &config.Project{
			ProjectMeta: p.config.ProjectMeta,
			Features:    p.config.Features,
			Domains:     allDomains,
		},
	}
	files = append(files, PlannedFile{
		TemplatePath: "frontend/consumer/layout/layout.svelte.tmpl",
		OutputPath:   "frontend/consumer/src/routes/+layout.svelte",
		Data:         projectCtx,
	})
	files = append(files, PlannedFile{
		TemplatePath: "frontend/consumer/layout/layout.ts.tmpl",
		OutputPath:   "frontend/consumer/src/routes/+layout.ts",
		Data:         projectCtx,
	})
	files = append(files, PlannedFile{
		TemplatePath: "frontend/consumer/layout/nav.svelte.tmpl",
		OutputPath:   "frontend/consumer/src/lib/components/Nav.generated.svelte",
		Data:         navCtx,
	})

	// Per-domain files.
	for _, domain := range allDomains {
		d := domain
		domainCtx := TemplateContext{
			Project: p.config,
			Domain:  &d,
		}
		files = append(files, PlannedFile{
			TemplatePath: "frontend/consumer/domain/types.ts.tmpl",
			OutputPath:   fmt.Sprintf("frontend/consumer/src/lib/api/%s/types.generated.ts", d.Name),
			Data:         domainCtx,
		})
		files = append(files, PlannedFile{
			TemplatePath: "frontend/consumer/domain/api.ts.tmpl",
			OutputPath:   fmt.Sprintf("frontend/consumer/src/lib/api/%s/api.generated.ts", d.Name),
			Data:         domainCtx,
		})

		// Per-entity route files.
		for _, entity := range d.Entities {
			e := entity
			entityCtx := TemplateContext{
				Project: p.config,
				Domain:  &d,
				Entity:  &e,
			}
			entityKebab := naming.New(e.Name).PluralKebab()

			// List and detail pages (always present).
			files = append(files, PlannedFile{
				TemplatePath: "frontend/consumer/entity/list_page.svelte.tmpl",
				OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/+page.svelte", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/consumer/entity/list_page_ts.tmpl",
				OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/+page.ts", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/consumer/entity/detail_page.svelte.tmpl",
				OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/[id]/+page.svelte", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/consumer/entity/detail_page_ts.tmpl",
				OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/[id]/+page.ts", entityKebab),
				Data:         entityCtx,
			})

			// Create and edit pages (only for consumer_editable entities).
			if e.ConsumerEditable {
				files = append(files, PlannedFile{
					TemplatePath: "frontend/consumer/entity/create_page.svelte.tmpl",
					OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/new/+page.svelte", entityKebab),
					Data:         entityCtx,
				})
				files = append(files, PlannedFile{
					TemplatePath: "frontend/consumer/entity/create_page_ts.tmpl",
					OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/new/+page.ts", entityKebab),
					Data:         entityCtx,
				})
				files = append(files, PlannedFile{
					TemplatePath: "frontend/consumer/entity/edit_page.svelte.tmpl",
					OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/[id]/edit/+page.svelte", entityKebab),
					Data:         entityCtx,
				})
				files = append(files, PlannedFile{
					TemplatePath: "frontend/consumer/entity/edit_page_ts.tmpl",
					OutputPath:   fmt.Sprintf("frontend/consumer/src/routes/%s/[id]/edit/+page.ts", entityKebab),
					Data:         entityCtx,
				})
			}
		}
	}

	return files
}

// planAdminAppFiles returns all files for the admin SvelteKit app.
func (p *Pipeline) planAdminAppFiles(allDomains []config.Domain) []PlannedFile {
	var files []PlannedFile

	projectCtx := TemplateContext{Project: p.config}

	// Project scaffolding.
	for _, m := range []struct{ template, output string }{
		{"frontend/admin/project/package.json.tmpl", "frontend/admin/package.json"},
		{"frontend/admin/project/svelte.config.js.tmpl", "frontend/admin/svelte.config.js"},
		{"frontend/admin/project/tsconfig.json.tmpl", "frontend/admin/tsconfig.json"},
		{"frontend/admin/project/vite.config.ts.tmpl", "frontend/admin/vite.config.ts"},
		{"frontend/admin/project/app.html.tmpl", "frontend/admin/src/app.html"},
		{"frontend/admin/project/app.d.ts.tmpl", "frontend/admin/src/app.d.ts"},
	} {
		files = append(files, PlannedFile{TemplatePath: m.template, OutputPath: m.output, Data: projectCtx})
	}

	// Shared library files.
	for _, m := range []struct{ template, output string }{
		{"frontend/shared/api_client.ts.tmpl", "frontend/admin/src/lib/api/client.generated.ts"},
		{"frontend/shared/types.ts.tmpl", "frontend/admin/src/lib/api/types.generated.ts"},
	} {
		files = append(files, PlannedFile{TemplatePath: m.template, OutputPath: m.output, Data: projectCtx})
	}

	// Layout files.
	navCtx := TemplateContext{
		Project: &config.Project{
			ProjectMeta: p.config.ProjectMeta,
			Features:    p.config.Features,
			Domains:     allDomains,
		},
	}
	files = append(files, PlannedFile{
		TemplatePath: "frontend/admin/layout/layout.svelte.tmpl",
		OutputPath:   "frontend/admin/src/routes/+layout.svelte",
		Data:         projectCtx,
	})
	files = append(files, PlannedFile{
		TemplatePath: "frontend/admin/layout/layout.ts.tmpl",
		OutputPath:   "frontend/admin/src/routes/+layout.ts",
		Data:         projectCtx,
	})
	files = append(files, PlannedFile{
		TemplatePath: "frontend/admin/layout/nav.svelte.tmpl",
		OutputPath:   "frontend/admin/src/lib/components/Nav.generated.svelte",
		Data:         navCtx,
	})

	// Per-domain files.
	for _, domain := range allDomains {
		d := domain
		domainCtx := TemplateContext{
			Project: p.config,
			Domain:  &d,
		}
		files = append(files, PlannedFile{
			TemplatePath: "frontend/admin/domain/types.ts.tmpl",
			OutputPath:   fmt.Sprintf("frontend/admin/src/lib/api/%s/types.generated.ts", d.Name),
			Data:         domainCtx,
		})
		files = append(files, PlannedFile{
			TemplatePath: "frontend/admin/domain/api.ts.tmpl",
			OutputPath:   fmt.Sprintf("frontend/admin/src/lib/api/%s/api.generated.ts", d.Name),
			Data:         domainCtx,
		})

		// Per-entity route files (admin always has full CRUD).
		for _, entity := range d.Entities {
			e := entity
			entityCtx := TemplateContext{
				Project: p.config,
				Domain:  &d,
				Entity:  &e,
			}
			entityKebab := naming.New(e.Name).PluralKebab()

			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/list_page.svelte.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/+page.svelte", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/list_page_ts.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/+page.ts", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/detail_page.svelte.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/[id]/+page.svelte", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/detail_page_ts.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/[id]/+page.ts", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/create_page.svelte.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/new/+page.svelte", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/create_page_ts.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/new/+page.ts", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/edit_page.svelte.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/[id]/edit/+page.svelte", entityKebab),
				Data:         entityCtx,
			})
			files = append(files, PlannedFile{
				TemplatePath: "frontend/admin/entity/edit_page_ts.tmpl",
				OutputPath:   fmt.Sprintf("frontend/admin/src/routes/%s/[id]/edit/+page.ts", entityKebab),
				Data:         entityCtx,
			})
		}
	}

	return files
}
