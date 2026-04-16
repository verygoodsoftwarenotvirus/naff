// Package pipeline — parked planner functions from the pre-pivot,
// registry-driven generator. Go's toolchain skips directories that start
// with `_`, so this file is NOT compiled as part of the module. Kept here
// for reference when re-introducing parameterization file-by-file.
//
// The file intentionally keeps `package pipeline` plus the original imports
// so it can be pasted back with minimal edits.
package pipeline

import (
	"fmt"
	"os"
	"strings"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/naming"
)

// planFiles builds the complete list of files to generate.
// This is the central registry that maps templates to output paths.
func (p *Pipeline) planFiles() []PlannedFile {
	var files []PlannedFile

	generateBackend := p.config.Targets.BackendEnabled()
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
		// planLocalDevFiles emits internal/localdev/server.go, shipped verbatim
		// from target in Phase 4b.4. It imports internal/build/services/api,
		// internal/domain/mealplanning, internal/repositories/postgres/auth+identity+oauth,
		// and internal/grpc/generated/services/auth — all DDB-specific or
		// codegen-dependent. Skipped in the generic path.
		// files = append(files, p.planLocalDevFiles()...)
		files = append(files, p.planDomainExtrasFiles()...)
		// planAuthHandlerFiles ships internal/services/auth/handlers/
		// authentication/*.go verbatim from target. Its do.go imports internal/
		// domain/identity/manager — a hand-rolled target package naff doesn't
		// emit. Omitted in the generic path.
		// files = append(files, p.planAuthHandlerFiles()...)
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
		{"domain/repository.go.tmpl", fmt.Sprintf("internal/domain/%s/repository.go", domain), true},
		{"domain_keys/keys.go.tmpl", fmt.Sprintf("internal/domain/%s/keys/keys.go", domain), true},
		{"domain_fakes/fake.go.tmpl", fmt.Sprintf("internal/domain/%s/fakes/fake.go", domain), true},
		{"domain_mock/repository.go.tmpl", fmt.Sprintf("internal/domain/%s/mock/repository.go", domain), true},

		// Manager (one per domain)
		{"manager/interface.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/interface.go", domain), true},
		{"manager/manager.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/manager.go", domain), true},
		{"manager/do.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/do.go", domain), true},
		{"manager/mock_manager.go.tmpl", fmt.Sprintf("internal/domain/%s/manager/mock/manager.go", domain), true},

		// Repository infrastructure (one per domain)
		{"repository/client.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/client.go", domain), true},
		{"repository/do.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/do.go", domain), true},

		// gRPC infrastructure (one per domain)
		{"grpc/service.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/service.go", domain), true},
		{"grpc/converters.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/converters/converters.go", domain), true},
		{"grpc/permissions.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/permissions.go", domain), true},
		{"grpc/do.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/do.go", domain), true},

		// Authorization (one per domain)
		{"authorization/permissions.go.tmpl", fmt.Sprintf("internal/authorization/%s_permissions.go", domain), true},

		// sqlc config (one per domain)
		{"sqlc/block.yaml.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/sqlc.yaml", domain), false},

		// Build extras
		{"build/extras.go.tmpl", fmt.Sprintf("internal/build/services/api/grpc/%s.go", domain), false},
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
		{"project/go.sum.tmpl", "go.sum", 0},
		{"project/Makefile.tmpl", "Makefile", 0},
		{"project/scripts/configs.sh.tmpl", "scripts/configs.sh", 0o755},
		{"project/scripts/queries.sh.tmpl", "scripts/queries.sh", 0o755},
		{"project/scripts/env_vars.sh.tmpl", "scripts/env_vars.sh", 0o755},
		{"project/scripts/format_golang.sh.tmpl", "scripts/format_golang.sh", 0o755},
		{"project/scripts/goimports.sh.tmpl", "scripts/goimports.sh", 0o755},
		{"project/scripts/format_imports.sh.tmpl", "scripts/format_imports.sh", 0o755},
		{"project/scripts/format_go_fieldalignment.sh.tmpl", "scripts/format_go_fieldalignment.sh", 0o755},
		{"project/scripts/format_go_tag_alignment.sh.tmpl", "scripts/format_go_tag_alignment.sh", 0o755},
		// Top-level shared protos:
		// - filtering.proto: QueryFilter + Pagination (imported by per-domain services)
		// - types/common.proto: ResponseDetails + NamedID (imported by per-domain services)
		{"proto/filtering.proto.tmpl", "proto/filtering.proto", 0},
		{"proto/common.proto.tmpl", "proto/types/common.proto", 0},
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
		{base + "/authentication_http_routes.go.tmpl", "internal/" + base + "/authentication_http_routes.go"},
		{base + "/config.go.tmpl", "internal/" + base + "/config.go"},
		{base + "/do.go.tmpl", "internal/" + base + "/do.go"},
		{base + "/doc.go.tmpl", "internal/" + base + "/doc.go"},
		{base + "/helpers.go.tmpl", "internal/" + base + "/helpers.go"},
		{base + "/oauth2.go.tmpl", "internal/" + base + "/oauth2.go"},
		{base + "/oauth2_client_info.go.tmpl", "internal/" + base + "/oauth2_client_info.go"},
		{base + "/oauth2_client_store.go.tmpl", "internal/" + base + "/oauth2_client_store.go"},
		{base + "/oauth2_client_token.go.tmpl", "internal/" + base + "/oauth2_client_token.go"},
		{base + "/oauth2_token_store.go.tmpl", "internal/" + base + "/oauth2_token_store.go"},
		{base + "/revoke.go.tmpl", "internal/" + base + "/revoke.go"},
		{base + "/service.go.tmpl", "internal/" + base + "/service.go"},
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
		{"auth/auth.go.tmpl", "internal/domain/auth/auth.go"},
		{"auth/do.go.tmpl", "internal/domain/auth/do.go"},
		{"auth/password_reset_token.go.tmpl", "internal/domain/auth/password_reset_token.go"},
		{"auth/repository.go.tmpl", "internal/domain/auth/repository.go"},
		{"auth/user.go.tmpl", "internal/domain/auth/user.go"},
		{"auth/user_session.go.tmpl", "internal/domain/auth/user_session.go"},

		// internal/domain/identity/
		{"identity/auth.go.tmpl", "internal/domain/identity/auth.go"},
		{"identity/account.go.tmpl", "internal/domain/identity/account.go"},
		{"identity/account_invitation.go.tmpl", "internal/domain/identity/account_invitation.go"},
		{"identity/account_user_membership.go.tmpl", "internal/domain/identity/account_user_membership.go"},
		{"identity/data_privacy.go.tmpl", "internal/domain/identity/data_privacy.go"},
		{"identity/do.go.tmpl", "internal/domain/identity/do.go"},
		{"identity/password_reset_token.go.tmpl", "internal/domain/identity/password_reset_token.go"},
		{"identity/repository.go.tmpl", "internal/domain/identity/repository.go"},
		{"identity/user.go.tmpl", "internal/domain/identity/user.go"},

		// internal/domain/oauth/
		{"oauth/do.go.tmpl", "internal/domain/oauth/do.go"},
		{"oauth/oauth2_client.go.tmpl", "internal/domain/oauth/oauth2_client.go"},
		{"oauth/oauth2_client_token.go.tmpl", "internal/domain/oauth/oauth2_client_token.go"},
		{"oauth/repository.go.tmpl", "internal/domain/oauth/repository.go"},

		// internal/authorization/ (per-domain permissions sidecar files for
		// auth/identity/oauth, since those domains are no longer in BuiltinDomains()
		// and thus don't trigger planPerDomainFiles' convention-based permissions template).
		{"authorization/identity_permissions.go.tmpl", "internal/authorization/identity_permissions.go"},
		{"authorization/oauth_permissions.go.tmpl", "internal/authorization/oauth_permissions.go"},

		// internal/branding/ stub — only the four constants in-scope generated
		// code actually imports (CompanyName, CompanyNameSlug, CompanySlug,
		// EnvVarPrefix). Users flesh out logos/email-templates/legal-text. Target's
		// real branding.go is 109KB with base64 logo blobs → deliberately NOT shipped.
		{"branding/branding.go.tmpl", "internal/branding/branding.go"},

		// Phase 4a + 4b.5 — integration harness. 4a shipped the stubs (matchers.go,
		// doc.go, constants.go). 4b.5 adds init.go/helpers.go/audit_helpers.go now
		// that 4b.0–4b.4 unblocked the deps: Provide<Domain>Repository rename
		// (4b.0), settings builtin (4b.1), top-level /internal/repositories glue
		// (4b.2), proto/filtering.proto (4b.3), and internal/localdev/server.go
		// (4b.4). init.go and helpers.go inherit the audit import alias from 4b.4
		// (naff emits postgres/audit; target wrote postgres/auditlogentries).
		{"testutils/matchers.go.tmpl", "internal/testutils/matchers.go"},
		{"testing/integration/apiserver/doc.go.tmpl", "testing/integration/apiserver/doc.go"},
		{"testing/integration/apiserver/constants.go.tmpl", "testing/integration/apiserver/constants.go"},
		// init.go/helpers.go/audit_helpers.go omitted in the generic path — they
		// reference internal/localdev (DDB-specific, dropped), internal/services/
		// identity/grpc/converters (the per-entity converters, not emitted
		// generically), and internal/grpc/generated/filtering (protoc output).
		// Re-enable per project once those are either shipped or stubbed.

		// pkg/client/client.go omitted in the generic path — target's version
		// hardcodes imports for every DDB service's protoc output. A future
		// templated version should iterate .AllDomains.
		// {"pkg/client/client.go.tmpl", "pkg/client/client.go"},

		// Phase 5 — hand-rolled postgres repos + identity adjuncts. Target ships
		// auth/identity/oauth as hand-rolled domain packages (3b) AND hand-rolled
		// postgres repos. Without these, a consumer can't compile localdev,
		// services/auth/handlers, or the integration harness. Ship verbatim
		// modulo module paths — the exact 3b/4b pattern.

		// Phase 5 postgres repos for auth/identity/oauth are omitted in the
		// generic path: their client.go and entity files reference `generated.
		// Querier` from internal/repositories/postgres/<X>/generated (sqlc
		// output from queries that target hand-writes and naff doesn't ship).
		// Without those SQL queries, sqlc emits an empty package and the
		// Phase 5 repos fail to compile. Re-enable after naff ships the auth/
		// identity/oauth SQL query sources.
		// postgres/testing is similarly dropped — its helpers.go uses
		// identity/fakes which is present, but its audit.go imports the full
		// DDB auditlogentries + transitive dependencies, not generic.

		// Ship grpc/converters — generic helpers for gRPC<->domain type
		// conversions, used by naff's per-domain grpc/*.go output.
		{"internal/grpc/converters/helpers.go.tmpl", "internal/grpc/converters/helpers.go"},
		{"internal/grpc/converters/query_filter.go.tmpl", "internal/grpc/converters/query_filter.go"},

		// audit domain extras — DataChangeMessage type + constructor used by
		// per-domain manager templates and internal/authentication/manager.go
		// for async worker dispatch. Target hand-rolls these in audit/; naff's
		// CRUD only emits AuditLogEntry.
		{"domain/audit_extras.go.tmpl", "internal/domain/audit/data_change_message.go"},

		// identity adjuncts — converters imported by localdev; fakes imported by
		// postgres/testing; keys (auth/identity/oauth) imported across repos.
		{"identity/converters/account_invitations.go.tmpl", "internal/domain/identity/converters/account_invitations.go"},
		{"identity/converters/account_user_memberships.go.tmpl", "internal/domain/identity/converters/account_user_memberships.go"},
		{"identity/converters/accounts.go.tmpl", "internal/domain/identity/converters/accounts.go"},
		// identity/converters/admin.go omitted — its only contents is a single
		// GRPC → domain converter that imports grpc/generated/services/identity
		// (protoc output, absent in a fresh project before `make generate`).
		{"identity/converters/users.go.tmpl", "internal/domain/identity/converters/users.go"},

		{"identity/fakes/account_invitation.go.tmpl", "internal/domain/identity/fakes/account_invitation.go"},
		{"identity/fakes/account_user_membership.go.tmpl", "internal/domain/identity/fakes/account_user_membership.go"},
		{"identity/fakes/account.go.tmpl", "internal/domain/identity/fakes/account.go"},
		{"identity/fakes/doc.go.tmpl", "internal/domain/identity/fakes/doc.go"},
		{"identity/fakes/fake.go.tmpl", "internal/domain/identity/fakes/fake.go"},
		{"identity/fakes/user.go.tmpl", "internal/domain/identity/fakes/user.go"},

		{"auth/keys/keys.go.tmpl", "internal/domain/auth/keys/keys.go"},
		{"identity/keys/keys.go.tmpl", "internal/domain/identity/keys/keys.go"},
		{"oauth/keys/keys.go.tmpl", "internal/domain/oauth/keys/keys.go"},
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
		{"authorization/permissions_core.go.tmpl", "internal/authorization/permissions.go"},
		{"authorization/account_role.go.tmpl", "internal/authorization/account_role.go"},
		{"authorization/service_role.go.tmpl", "internal/authorization/service_role.go"},
		{"authorization/auth_permissions.go.tmpl", "internal/authorization/auth_permissions.go"},
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
		{"authentication/aliases.go.tmpl", "internal/authentication/aliases.go"},
		{"authentication/manager.go.tmpl", "internal/authentication/manager.go"},
		{"authentication/do.go.tmpl", "internal/authentication/do.go"},
		{"authentication/config/config.go.tmpl", "internal/authentication/config/config.go"},
		{"authentication/config/do.go.tmpl", "internal/authentication/config/do.go"},
		{"authentication/mock/mock_authenticator.go.tmpl", "internal/authentication/mock/mock_authenticator.go"},
		{"authentication/mock/authentication_manager.go.tmpl", "internal/authentication/mock/authentication_manager.go"},
		{"authentication/mocks/mock_user.go.tmpl", "internal/authentication/mocks/mock_user.go"},
		{"authentication/sessions/errors.go.tmpl", "internal/authentication/sessions/errors.go"},
		{"authentication/sessions/session_context.go.tmpl", "internal/authentication/sessions/session_context.go"},
		{"authentication/sessions/do.go.tmpl", "internal/authentication/sessions/do.go"},
		{"authentication/webauthn/session_store.go.tmpl", "internal/authentication/webauthn/session_store.go"},
		{"authentication/webauthn/postgres_session_store.go.tmpl", "internal/authentication/webauthn/postgres_session_store.go"},
		{"authentication/webauthn/service.go.tmpl", "internal/authentication/webauthn/service.go"},
		{"authentication/webauthn/user_adapter.go.tmpl", "internal/authentication/webauthn/user_adapter.go"},
		{"authentication/webauthn/config/config.go.tmpl", "internal/authentication/webauthn/config/config.go"},
		{"identity/webauthn_credential.go.tmpl", "internal/domain/identity/webauthn_credential.go"},
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
		{"config/meta.go.tmpl", "internal/config/meta.go"},
		{"config/queues.go.tmpl", "internal/config/queues.go"},
		{"config/doc.go.tmpl", "internal/config/doc.go"},
		{"config/configs.go.tmpl", "internal/config/configs.go"},
		{"config/do.go.tmpl", "internal/config/do.go"},
		{"config/env_vars.go.tmpl", "internal/config/env_vars.go"},
		// environment.go + mealplanning_* omitted: they hardcode references
		// to target's ServicesConfig fields (DataPrivacy, MealPlanning, etc.)
		// and per-entity config paths. Users write their own env→config
		// loader. Re-enable after templating on .AllDomains.
		// config/services_config.go is DDB-specific (hardcoded service list:
		// dataprivacy/identity/mealplanning/oauth/payments/uploadedmedia). It
		// imports services/<X>/config packages that naff doesn't emit in the
		// generic path. The generic template emits an empty ServicesConfig
		// stub so config/{configs,do}.go can still reference the type.
		{"config/services_config.go.tmpl", "internal/config/services_config.go"},
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
		{"repositories/do.go.tmpl", "internal/repositories/do.go"},
		{"repositories/migrations.go.tmpl", "internal/repositories/migrations.go"},
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
			OutputPath:   "internal/localdev/server.go",
			Data:         ctx,
			IsGo:         true,
		},
	}
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
		// services/api root — generic API server entry point. Simplified to
		// inline DI; does not depend on internal/build/services/api/*.
		{"cmd/services/api/main.go.tmpl", "cmd/services/api/main.go", true, projectCtx},
		{"cmd/services/api/doc.go.tmpl", "cmd/services/api/doc.go", true, projectCtx},

		// tools/codegen/queries — emits the per-entity sqlc query SQL files that
		// feed `sqlc generate`. Required by `make generate` in the fresh-project
		// flow. main.go uses a per-entity dispatch map (needs codegenCtx).
		{"cmd/tools/codegen/queries/main.go.tmpl", "cmd/tools/codegen/queries/main.go", true, codegenCtx},
		{"cmd/tools/codegen/queries/helpers.go.tmpl", "cmd/tools/codegen/queries/helpers.go", true, projectCtx},
		{"cmd/tools/codegen/queries/sqlc.go.tmpl", "cmd/tools/codegen/queries/sqlc.go", true, projectCtx},

		// playground — trivial id-gen dev utility, always generic
		{"cmd/playground/main.go.tmpl", "cmd/playground/main.go", true, projectCtx},

		// DDB-specific binaries (services/api/grpc & http split, services/mcp,
		// tools/*, workers/*, localdev/server, tools/codegen/configs,
		// tools/codegen/valid_env_vars) are not emitted in the generic path:
		// their templates were lifted from target verbatim and reference
		// packages (internal/build/services/api, internal/build/jobs/*,
		// internal/build/functions/*, internal/domain/mealplanning/*, etc.)
		// that don't exist in a fresh project. Re-enabling them requires either
		// templating on .AllDomains or reintroducing those packages as
		// project-specific extensions. See plan pivot (2026-04-14).
		//
		// functions/async_message_handler is opted into via
		// Targets.AsyncMessageHandler (see below) — the first re-enabled
		// DDB-shaped binary.
	}

	// Opt-in DDB-shaped binaries. Each guard gates one binary whose template
	// imports a hand-rolled target package naff does not generate; a project
	// that authors those packages can flip the corresponding Targets flag.
	if p.config.Targets.AsyncMessageHandlerEnabled() {
		mappings = append(mappings,
			m{"cmd/functions/async_message_handler/main.go.tmpl", "cmd/functions/async_message_handler/main.go", true, projectCtx},
			m{"cmd/functions/async_message_handler/README.md.tmpl", "cmd/functions/async_message_handler/README.md", false, projectCtx},
		)
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
		{"domain/entity.go.tmpl", fmt.Sprintf("internal/domain/%s/%s.go", domain, entitySnake), true},
		{"domain_converters/converters.go.tmpl", fmt.Sprintf("internal/domain/%s/converters/%s.go", domain, entitySnake), true},

		// Repository entity CRUD
		{"repository/entity.go.tmpl", fmt.Sprintf("internal/repositories/postgres/%s/%s.go", domain, entitySnake), true},

		// Query codegen — filename matches dinnerdonebetter pattern: <domain>_<plural_snake>.go
		{"codegen/queries.go.tmpl", fmt.Sprintf("cmd/tools/codegen/queries/%s_%s.go", domain, entityPluralSnake), true},

		// Migration (one per entity)
		{"migrations/migration.sql.tmpl", fmt.Sprintf("internal/repositories/postgres/migrations/migration_files/%s_%s.sql", domain, entitySnake), false},

		// Proto
		{"proto/messages.proto.tmpl", fmt.Sprintf("proto/%s/%s_messages.proto", domain, entitySnake), false},
		{"proto/service.proto.tmpl", fmt.Sprintf("proto/%s/%s_service.proto", domain, entitySnake), false},

		// gRPC entity handlers
		{"grpc/entity.go.tmpl", fmt.Sprintf("internal/services/%s/grpc/%s.go", domain, entitySnake), true},
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
		{"ios/project/App.swift.tmpl", fmt.Sprintf("%s/App.swift", src)},
		{"ios/project/ContentView.swift.tmpl", fmt.Sprintf("%s/ContentView.swift", src)},
		{"ios/project/Configuration.swift.tmpl", fmt.Sprintf("%s/Configuration.swift", src)},
		{"ios/api/GRPCClient.swift.tmpl", fmt.Sprintf("%s/API/GRPCClient.swift", src)},
		{"ios/auth/AuthManager.swift.tmpl", fmt.Sprintf("%s/Auth/AuthManager.swift", src)},
		{"ios/auth/LoginView.swift.tmpl", fmt.Sprintf("%s/Auth/LoginView.swift", src)},
		{"ios/auth/RegistrationView.swift.tmpl", fmt.Sprintf("%s/Auth/RegistrationView.swift", src)},
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
			OutputPath:   fmt.Sprintf("%s/API/%sService.swift", src, strings.Title(domain)), //nolint:staticcheck
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
			OutputPath:   fmt.Sprintf("%s/Models/%s.swift", src, entityName),
			Data:         ctx,
		},
		{
			TemplatePath: "ios/models/entity_input.swift.tmpl",
			OutputPath:   fmt.Sprintf("%s/Models/%sInput.swift", src, entityName),
			Data:         ctx,
		},
		{
			TemplatePath: "proto/messages.proto.tmpl",
			OutputPath:   fmt.Sprintf("%s/Proto/%s/%s_messages.proto", p.iosBasePath(), ctx.Domain.Name, strings.ToLower(entityName)),
			Data:         ctx,
		},
		{
			TemplatePath: "proto/service.proto.tmpl",
			OutputPath:   fmt.Sprintf("%s/Proto/%s/%s_service.proto", p.iosBasePath(), ctx.Domain.Name, strings.ToLower(entityName)),
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
		{"frontend/shared/api_client.ts.tmpl", "frontend/consumer/src/lib/api/client.ts"},
		{"frontend/shared/types.ts.tmpl", "frontend/consumer/src/lib/api/types.ts"},
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
		OutputPath:   "frontend/consumer/src/lib/components/Nav.svelte",
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
			OutputPath:   fmt.Sprintf("frontend/consumer/src/lib/api/%s/types.ts", d.Name),
			Data:         domainCtx,
		})
		files = append(files, PlannedFile{
			TemplatePath: "frontend/consumer/domain/api.ts.tmpl",
			OutputPath:   fmt.Sprintf("frontend/consumer/src/lib/api/%s/api.ts", d.Name),
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
		{"frontend/shared/api_client.ts.tmpl", "frontend/admin/src/lib/api/client.ts"},
		{"frontend/shared/types.ts.tmpl", "frontend/admin/src/lib/api/types.ts"},
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
		OutputPath:   "frontend/admin/src/lib/components/Nav.svelte",
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
			OutputPath:   fmt.Sprintf("frontend/admin/src/lib/api/%s/types.ts", d.Name),
			Data:         domainCtx,
		})
		files = append(files, PlannedFile{
			TemplatePath: "frontend/admin/domain/api.ts.tmpl",
			OutputPath:   fmt.Sprintf("frontend/admin/src/lib/api/%s/api.ts", d.Name),
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
