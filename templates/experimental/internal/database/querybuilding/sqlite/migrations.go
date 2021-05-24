package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrationsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("migrations").Op("=").Index().ID("darwin").Dot("Migration").Valuesln(jen.Valuesln(jen.ID("Version").Op(":").Lit(0.00), jen.ID("Description").Op(":").Lit("create sessions table for session manager"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BLOB NOT NULL,
				expiry REAL NOT NULL,
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now'))
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.01), jen.ID("Description").Op(":").Lit("create sessions table for session manager"), jen.ID("Script").Op(":").Lit(`CREATE INDEX sessions_expiry_idx ON sessions(expiry);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.02), jen.ID("Description").Op(":").Lit("create audit log table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS audit_log (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				external_id TEXT NOT NULL,
				event_type TEXT NOT NULL,
				context JSON NOT NULL,
				created_on BIGINT NOT NULL DEFAULT (strftime('%s','now'))
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.03), jen.ID("Description").Op(":").Lit("create users table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				external_id TEXT NOT NULL,
				username TEXT NOT NULL,
				avatar_src TEXT,
				hashed_password TEXT NOT NULL,
				password_last_changed_on INTEGER DEFAULT NULL,
				requires_password_change BOOLEAN NOT NULL DEFAULT 'false',
				two_factor_secret TEXT NOT NULL,
				two_factor_secret_verified_on INTEGER DEFAULT NULL,
				service_roles TEXT NOT NULL DEFAULT 'service_user',
				reputation TEXT NOT NULL DEFAULT 'unverified',
				reputation_explanation TEXT NOT NULL DEFAULT '',
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				last_updated_on INTEGER DEFAULT NULL,
				archived_on INTEGER DEFAULT NULL,
				CONSTRAINT username_unique UNIQUE (username)
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.04), jen.ID("Description").Op(":").Lit("create accounts table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS accounts (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				billing_status TEXT NOT NULL DEFAULT 'unpaid',
				contact_email TEXT NOT NULL DEFAULT '',
				contact_phone TEXT NOT NULL DEFAULT '',
				payment_processor_customer_id TEXT NOT NULL DEFAULT '',
				subscription_plan_id TEXT,
				belongs_to_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				last_updated_on INTEGER DEFAULT NULL,
				archived_on INTEGER DEFAULT NULL,
				CONSTRAINT account_name_unique UNIQUE (name, belongs_to_user)
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.05), jen.ID("Description").Op(":").Lit("create account user memberships table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS account_user_memberships (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				belongs_to_account INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
				belongs_to_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				account_roles TEXT NOT NULL DEFAULT 'account_user',
				default_account BOOLEAN NOT NULL DEFAULT 'false',
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				last_updated_on INTEGER DEFAULT NULL,
				archived_on INTEGER DEFAULT NULL,
				CONSTRAINT plan_name_unique UNIQUE (belongs_to_account, belongs_to_user)
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.06), jen.ID("Description").Op(":").Lit("create API clients table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS api_clients (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				external_id TEXT NOT NULL,
				name TEXT DEFAULT '',
				client_id TEXT NOT NULL,
				secret_key TEXT NOT NULL,
				permissions INTEGER NOT NULL DEFAULT 0,
				admin_permissions INTEGER NOT NULL DEFAULT 0,
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				last_updated_on INTEGER DEFAULT NULL,
				archived_on INTEGER DEFAULT NULL,
				belongs_to_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.07), jen.ID("Description").Op(":").Lit("create webhooks table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS webhooks (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				content_type TEXT NOT NULL,
				url TEXT NOT NULL,
				method TEXT NOT NULL,
				events TEXT NOT NULL,
				data_types TEXT NOT NULL,
				topics TEXT NOT NULL,
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				last_updated_on INTEGER DEFAULT NULL,
				archived_on INTEGER DEFAULT NULL,
				belongs_to_account INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.08), jen.ID("Description").Op(":").Lit("create items table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS items (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				external_id TEXT NOT NULL,
				name TEXT NOT NULL,
				details TEXT NOT NULL DEFAULT '',
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				last_updated_on INTEGER DEFAULT NULL,
				archived_on INTEGER DEFAULT NULL,
				belongs_to_account INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Line(),
		jen.Func().Comment("migrate a sqlite database.").Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()).Body(
			jen.Return().Func().Params().Body(
				jen.ID("d").Op(":=").ID("darwin").Dot("NewGenericDriver").Call(
					jen.ID("db"),
					jen.ID("darwin").Dot("SqliteDialect").Valuesln(),
				),
				jen.If(jen.ID("err").Op(":=").ID("darwin").Dot("Migrate").Call(
					jen.ID("d"),
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
