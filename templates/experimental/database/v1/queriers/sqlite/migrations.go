package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func migrationsDotGo() *jen.File {
	ret := jen.NewFile("sqlite")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("migrations").Op("=").Index().ID("darwin").Dot(
		"Migration",
	).Valuesln(jen.Valuesln(jen.ID("Version").Op(":").Lit(1), jen.ID("Description").Op(":").Lit("create users table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS users (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"password_last_changed_on" INTEGER,
				"two_factor_secret" TEXT NOT NULL,
				"is_admin" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				CONSTRAINT username_unique UNIQUE (username)
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(2), jen.ID("Description").Op(":").Lit("Add OAuth2 Clients table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(3), jen.ID("Description").Op(":").Lit("Create items table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS items (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT NOT NULL,
				"details" TEXT,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(4), jen.ID("Description").Op(":").Lit("Create webhooks table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS webhooks (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" text NOT NULL,
				"content_type" text NOT NULL,
				"url" text NOT NULL,
				"method" text NOT NULL,
				"events" text NOT NULL,
				"data_types" text NOT NULL,
				"topics" text NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`))),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildMigrationFunc returns a sync.Once compatible function closure that will").Comment("// migrate a sqlite database").ID("buildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()).Block(
		jen.Return().Func().Params().Block(
			jen.ID("driver").Op(":=").ID("darwin").Dot(
				"NewGenericDriver",
			).Call(jen.ID("db"), jen.ID("darwin").Dot(
				"SqliteDialect",
			).Valuesln()),
			jen.If(jen.ID("err").Op(":=").ID("darwin").Dot(
				"New",
			).Call(jen.ID("driver"), jen.ID("migrations"), jen.ID("nil")).Dot(
				"Migrate",
			).Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be").Comment("// safe (as in idempotent, though not recommended) to call this function multiple times.").Params(jen.ID("s").Op("*").ID("Sqlite")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
		jen.ID("s").Dot(
			"logger",
		).Dot(
			"Info",
		).Call(jen.Lit("migrating db")),
		jen.If(jen.Op("!").ID("s").Dot(
			"IsReady",
		).Call(jen.ID("ctx"))).Block(
			jen.Return().ID("errors").Dot(
				"New",
			).Call(jen.Lit("db is not ready yet")),
		),
		jen.ID("s").Dot(
			"migrateOnce",
		).Dot(
			"Do",
		).Call(jen.ID("buildMigrationFunc").Call(jen.ID("s").Dot(
			"db",
		))),
		jen.ID("s").Dot(
			"logger",
		).Dot(
			"Debug",
		).Call(jen.Lit("database migrated without error!")),
		jen.Return().ID("nil"),
	),

		jen.Line(),
	)
	return ret
}
