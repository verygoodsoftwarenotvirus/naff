package postgres

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func migrationsDotGo() *jen.File {
	ret := jen.NewFile("postgres")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("migrations").Op("=").Index().ID("darwin").Dot(
		"Migration",
	).Valuesln(jen.Valuesln(jen.ID("Version").Op(":").Lit(1), jen.ID("Description").Op(":").Lit("create users table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS users (
				"id" bigserial NOT NULL PRIMARY KEY,
				"username" text NOT NULL,
				"hashed_password" text NOT NULL,
				"password_last_changed_on" integer,
				"two_factor_secret" text NOT NULL,
				"is_admin" boolean NOT NULL DEFAULT 'false',
				"created_on" bigint NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" bigint DEFAULT NULL,
				"archived_on" bigint DEFAULT NULL,
				UNIQUE ("username")
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(2), jen.ID("Description").Op(":").Lit("Add oauth2_clients table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" bigserial NOT NULL PRIMARY KEY,
				"name" text DEFAULT '',
				"client_id" text NOT NULL,
				"client_secret" text NOT NULL,
				"redirect_uri" text DEFAULT '',
				"scopes" text NOT NULL,
				"implicit_allowed" boolean NOT NULL DEFAULT 'false',
				"created_on" bigint NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" bigint DEFAULT NULL,
				"archived_on" bigint DEFAULT NULL,
				"belongs_to" bigint NOT NULL,
				FOREIGN KEY(belongs_to) REFERENCES users(id)
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(3), jen.ID("Description").Op(":").Lit("create webhooks table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS webhooks (
				"id" bigserial NOT NULL PRIMARY KEY,
				"name" text NOT NULL,
				"content_type" text NOT NULL,
				"url" text NOT NULL,
				"method" text NOT NULL,
				"events" text NOT NULL,
				"data_types" text NOT NULL,
				"topics" text NOT NULL,
				"created_on" bigint NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" bigint DEFAULT NULL,
				"archived_on" bigint DEFAULT NULL,
				"belongs_to" bigint NOT NULL,
				FOREIGN KEY ("belongs_to") REFERENCES "users"("id")
			);`)), jen.Valuesln(jen.ID("Version").Op(":").Lit(4), jen.ID("Description").Op(":").Lit("create items table"), jen.ID("Script").Op(":").Lit(`
			CREATE TABLE IF NOT EXISTS items (
				"id" bigserial NOT NULL PRIMARY KEY,
				"name" text NOT NULL,
				"details" text NOT NULL DEFAULT '',
				"created_on" bigint NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" bigint DEFAULT NULL,
				"archived_on" bigint DEFAULT NULL,
				"belongs_to" bigint NOT NULL,
				FOREIGN KEY ("belongs_to") REFERENCES "users"("id")
			);`))),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildMigrationFunc returns a sync.Once compatible function closure that will").Comment("// migrate a postgres database").ID("buildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()).Block(
		jen.Return().Func().Params().Block(
			jen.ID("driver").Op(":=").ID("darwin").Dot(
				"NewGenericDriver",
			).Call(jen.ID("db"), jen.ID("darwin").Dot(
				"PostgresDialect",
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
	ret.Add(jen.Func().Comment("// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be").Comment("// safe (as in idempotent, though not recommended) to call this function multiple times.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
		jen.ID("p").Dot(
			"logger",
		).Dot(
			"Info",
		).Call(jen.Lit("migrating db")),
		jen.If(jen.Op("!").ID("p").Dot(
			"IsReady",
		).Call(jen.ID("ctx"))).Block(
			jen.Return().ID("errors").Dot(
				"New",
			).Call(jen.Lit("db is not ready yet")),
		),
		jen.ID("p").Dot(
			"migrateOnce",
		).Dot(
			"Do",
		).Call(jen.ID("buildMigrationFunc").Call(jen.ID("p").Dot(
			"db",
		))),
		jen.Return().ID("nil"),
	),

		jen.Line(),
	)
	return ret
}
