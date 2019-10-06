package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func migrationsDotGo() *jen.File {
	ret := jen.NewFile("sqlite")
	ret.Add(jen.Null(),
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
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
