package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func migrationsDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

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
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
