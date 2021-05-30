package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrationsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("migrations").Op("=").Index().ID("darwin").Dot("Migration").Valuesln(jen.Valuesln(jen.ID("Version").Op(":").Lit(0.00), jen.ID("Description").Op(":").Lit("create sessions table for session manager"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BYTEA NOT NULL,
				expiry TIMESTAMPTZ NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.01), jen.ID("Description").Op(":").Lit("create sessions table for session manager"), jen.ID("Script").Op(":").Lit(`CREATE INDEX sessions_expiry_idx ON sessions (expiry);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.02), jen.ID("Description").Op(":").Lit("create audit log table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS audit_log (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				event_type TEXT NOT NULL,
				context JSONB NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.03), jen.ID("Description").Op(":").Lit("create users table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				username TEXT NOT NULL,
				avatar_src TEXT,
				hashed_password TEXT NOT NULL,
				password_last_changed_on INTEGER,
				requires_password_change BOOLEAN NOT NULL DEFAULT 'false',
				two_factor_secret TEXT NOT NULL,
				two_factor_secret_verified_on BIGINT DEFAULT NULL,
				service_roles TEXT NOT NULL DEFAULT 'service_user',
				reputation TEXT NOT NULL DEFAULT 'unverified',
				reputation_explanation TEXT NOT NULL DEFAULT '',
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				UNIQUE("username")
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.04), jen.ID("Description").Op(":").Lit("create accounts table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS accounts (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				billing_status TEXT NOT NULL DEFAULT 'unpaid',
				contact_email TEXT NOT NULL DEFAULT '',
				contact_phone TEXT NOT NULL DEFAULT '',
				payment_processor_customer_id TEXT NOT NULL DEFAULT '',
				subscription_plan_id TEXT,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_user BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				UNIQUE("belongs_to_user", "name")
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.05), jen.ID("Description").Op(":").Lit("create account user memberships table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS account_user_memberships (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
				belongs_to_user BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				default_account BOOLEAN NOT NULL DEFAULT 'false',
				account_roles TEXT NOT NULL DEFAULT 'account_user',
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				UNIQUE("belongs_to_account", "belongs_to_user")
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.06), jen.ID("Description").Op(":").Lit("create API clients table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS api_clients (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT DEFAULT '',
				client_id TEXT NOT NULL,
				secret_key BYTEA NOT NULL,
				permissions BIGINT NOT NULL DEFAULT 0,
				admin_permissions BIGINT NOT NULL DEFAULT 0,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_user BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.07), jen.ID("Description").Op(":").Lit("create webhooks table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS webhooks (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				content_type TEXT NOT NULL,
				url TEXT NOT NULL,
				method TEXT NOT NULL,
				events TEXT NOT NULL,
				data_types TEXT NOT NULL,
				topics TEXT NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_account INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.8), jen.ID("Description").Op(":").Lit("create items table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS items (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				details TEXT NOT NULL DEFAULT '',
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				last_updated_on BIGINT DEFAULT NULL,
				archived_on BIGINT DEFAULT NULL,
				belongs_to_account INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Line(),
		jen.Func().Comment("migrate a postgres database.").Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()).Body(
			jen.Return().Func().Params().Body(
				jen.ID("driver").Op(":=").ID("darwin").Dot("NewGenericDriver").Call(
					jen.ID("db"),
					jen.ID("darwin").Dot("PostgresDialect").Valuesln(),
				),
				jen.If(jen.ID("err").Op(":=").ID("darwin").Dot("Migrate").Call(
					jen.ID("driver"),
					jen.ID("migrations"),
					jen.ID("nil"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("panic").Call(jen.Qual("fmt", "Errorf").Call(
						jen.Lit("migrating database: %w"),
						jen.ID("err"),
					))),
			)),
		jen.Line(),
	)

	return code
}
