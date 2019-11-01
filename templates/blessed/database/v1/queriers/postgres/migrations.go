package postgres

import (
	"fmt"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func typeToPostgresType(t string) string {
	typeMap := map[string]string{
		"[]string": "CHARACTER VARYING",
		"string":   "CHARACTER VARYING",
		"*string":  "CHARACTER VARYING",
		"uint64":   "BIGINT",
		"*uint64":  "BIGINT",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"uint":     "INTEGER",
		"*uint":    "INTEGER",
		"float64":  "NUMERIC",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	return t
}

// func typeToSqliteType(t string) string {
// 	typeMap := map[string]string{
// 		"[]string": "CHARACTER VARYING",
// 		"string":   "CHARACTER VARYING",
// 		"*string":  "CHARACTER VARYING",
// 		"uint64":   "INTEGER",
// 		"*uint64":  "INTEGER",
// 		"bool":     "BOOLEAN",
// 		"*bool":    "BOOLEAN",
// 		"int":      "INTEGER",
// 		"*int":     "INTEGER",
// 		"uint":     "INTEGER",
// 		"*uint":    "INTEGER",
// 		"float64":  "REAL",
// 	}

// 	if x, ok := typeMap[t]; ok {
// 		return x
// 	}

// 	return t
// }

// func typeToMariaDBType(t string) string {
// 	typeMap := map[string]string{
// 		"[]string": "LONGTEXT",
// 		"string":   "LONGTEXT",
// 		"*string":  "LONGTEXT",
// 		"uint64":   "INTEGER UNSIGNED",
// 		"*uint64":  "INTEGER UNSIGNED",
// 		"bool":     "BOOLEAN",
// 		"*bool":    "BOOLEAN",
// 		"int":      "INTEGER",
// 		"*int":     "INTEGER",
// 		"uint":     "INTEGER UNSIGNED",
// 		"*uint":    "INTEGER UNSIGNED",
// 		"float64":  "REAL",
// 	}

// 	if x, ok := typeMap[t]; ok {
// 		return x
// 	}

// 	log.Println("typeToMariaDBType called for type: ", t)
// 	return t
// }

// func typeExample(t string, pointer bool) interface{} {
// 	typeMap := map[string]interface{}{
// 		"[]string": `[]string{"example"}`,
// 		"string":   `"example"`,
// 		"*string":  `func(s string) *string { return &s }("example")`,
// 		"uint64":   "uint64(123)",
// 		"*uint64":  "func(u uint64) *uint64 { return &u }(123)",
// 		"bool":     false,
// 		"*bool":    "func(b bool) *bool { return &b }(false)",
// 		"int":      "int(456)",
// 		"*int":     "func(i int) *int { return &i }(123)",
// 		"uint":     "uint(456)",
// 		"*uint":    "func(i uint) *uint { return &i }(123)",
// 		"float64":  "float64(12.34)",
// 	}

// 	tn := t
// 	if pointer {
// 		tn = fmt.Sprintf("*%s", tn)
// 	}

// 	if x, ok := typeMap[tn]; ok {
// 		return x
// 	}

// 	return t
// }

func makeMigrations(types []models.DataType) []jen.Code {
	var out []jen.Code

	type migration struct {
		description, script string
	}

	migrations := []migration{
		{
			description: "create users table",
			script: `
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
			);`,
		},
		{
			description: "create oauth2_clients table",
			script: `
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
			);`,
		},
		{
			description: "create webhooks table",
			script: `
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
			);`,
		},
	}

	for _, typ := range types {
		pcn := typ.Name.PluralCommonName()

		scriptParts := []string{
			"\n			CREATE TABLE IF NOT EXISTS items (",
			`				"id" bigserial NOT NULL PRIMARY KEY,`,
		}

		for _, field := range typ.Fields {
			rn := field.Name.RouteName()

			query := fmt.Sprintf("				%q %s", rn, typeToPostgresType(field.Type))

			if !field.Pointer {
				query += ` NOT NULL`
			}

			if field.DefaultAllowed {
				query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
			}

			scriptParts = append(scriptParts, fmt.Sprintf("%s,", query))
		}

		scriptParts = append(scriptParts,
			`				"created_on" bigint NOT NULL DEFAULT extract(epoch FROM NOW()),`,
			`				"updated_on" bigint DEFAULT NULL,`,
			`				"archived_on" bigint DEFAULT NULL,`,
			`				"belongs_to" bigint NOT NULL,`,
			`				FOREIGN KEY ("belongs_to") REFERENCES "users"("id")`,
			`			);`,
		)

		migrations = append(migrations,
			migration{
				description: fmt.Sprintf("create %s table", pcn),
				script:      strings.Join(scriptParts, "\n"),
			},
		)
	}

	for i, script := range migrations {
		out = append(out, jen.Valuesln(
			jen.ID("Version").Op(":").Lit(i+1),
			jen.ID("Description").Op(":").Lit(script.description),
			jen.ID("Script").Op(":").Lit(script.script),
		))
	}

	return out
}

func migrationsDotGo(types []models.DataType) *jen.File {
	ret := jen.NewFile("postgres")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("migrations").Op("=").Index().Qual("github.com/GuiaBolso/darwin", "Migration").Valuesln(makeMigrations(types)...),
		),
	)

	ret.Add(
		jen.Comment("buildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Line(),
		jen.Comment("migrate a postgres database"),
		jen.Line(),
		jen.Func().ID("buildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Func().Params()).Block(
			jen.Return().Func().Params().Block(
				jen.ID("driver").Op(":=").Qual("github.com/GuiaBolso/darwin", "NewGenericDriver").Call(jen.ID("db"), jen.Qual("github.com/GuiaBolso/darwin", "PostgresDialect").Values()),
				jen.If(jen.ID("err").Op(":=").Qual("github.com/GuiaBolso/darwin", "New").Call(jen.ID("driver"), jen.ID("migrations"), jen.ID("nil")).Dot("Migrate").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("panic").Call(jen.ID("err")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be"),
		jen.Line(),
		jen.Comment("safe (as in idempotent, though not recommended) to call this function multiple times."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
			jen.ID("p").Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.If(jen.Op("!").ID("p").Dot("IsReady").Call(jen.ID("ctx"))).Block(
				jen.Return().ID("errors").Dot("New").Call(jen.Lit("db is not ready yet")),
			),
			jen.Line(),
			jen.ID("p").Dot("migrateOnce").Dot("Do").Call(jen.ID("buildMigrationFunc").Call(jen.ID("p").Dot("db"))),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)
	return ret
}