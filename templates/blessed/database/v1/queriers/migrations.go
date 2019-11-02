package queriers

import (
	"fmt"
	"log"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
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

func typeToSqliteType(t string) string {
	typeMap := map[string]string{
		"[]string": "CHARACTER VARYING",
		"string":   "CHARACTER VARYING",
		"*string":  "CHARACTER VARYING",
		"uint64":   "INTEGER",
		"*uint64":  "INTEGER",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"uint":     "INTEGER",
		"*uint":    "INTEGER",
		"float64":  "REAL",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	log.Println("typeToSqliteType called for type: ", t)
	return t
}

func typeToMariaDBType(t string) string {
	typeMap := map[string]string{
		"[]string": "LONGTEXT",
		"string":   "LONGTEXT",
		"*string":  "LONGTEXT",
		"uint64":   "INTEGER UNSIGNED",
		"*uint64":  "INTEGER UNSIGNED",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"uint":     "INTEGER UNSIGNED",
		"*uint":    "INTEGER UNSIGNED",
		"float64":  "REAL",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	return t
}

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

type migration struct {
	description string
	script      jen.Code
}

func makeMigrations(dbVendor wordsmith.SuperPalabra, types []models.DataType) []jen.Code {
	var (
		out        []jen.Code
		migrations []migration
	)

	dbrn := strings.ToLower(dbVendor.RouteName())

	switch dbrn {
	case "postgres":
		migrations = []migration{
			{
				description: "create users table",
				script: jen.Lit(`
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
			);`),
			},
			{
				description: "create oauth2_clients table",
				script: jen.Lit(`
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
			);`),
			},
			{
				description: "create webhooks table",
				script: jen.Lit(`
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
			);`),
			},
		}

		for _, typ := range types {
			pcn := typ.Name.PluralCommonName()

			scriptParts := []string{
				fmt.Sprintf("\n			CREATE TABLE IF NOT EXISTS %s (", typ.Name.PluralRouteName()),
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
					script:      jen.Lit(strings.Join(scriptParts, "\n")),
				},
			)
		}
	case "sqlite":
		migrations = []migration{
			{
				description: "create users table",
				script: jen.Lit(`
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
			);`),
			},
			{
				description: "create oauth2_clients table",
				script: jen.Lit(`
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
			);`),
			},
			{
				description: "create webhooks table",
				script: jen.Lit(`
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
			);`),
			},
		}

		idType := "INTEGER"
		idAddendum := " AUTOINCREMENT"

		for _, typ := range types {
			pcn := typ.Name.PluralCommonName()

			scriptParts := []string{
				fmt.Sprintf("\n			CREATE TABLE IF NOT EXISTS %s (", typ.Name.PluralRouteName()),
				fmt.Sprintf(`				"id" %s NOT NULL PRIMARY KEY%s,`, idType, idAddendum),
			}

			for _, field := range typ.Fields {
				rn := field.Name.RouteName()

				query := fmt.Sprintf("				%q %s", rn, typeToSqliteType(field.Type))

				if !field.Pointer {
					query += ` NOT NULL`
				}

				if field.DefaultAllowed {
					query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
				}

				scriptParts = append(scriptParts, fmt.Sprintf("%s,", query))
			}

			scriptParts = append(scriptParts,
				`				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),`,
				`				"updated_on" INTEGER DEFAULT NULL,`,
				`				"archived_on" INTEGER DEFAULT NULL,`,
				`				"belongs_to" INTEGER NOT NULL,`,
				`				FOREIGN KEY(belongs_to) REFERENCES users(id)`,
				`			);`,
			)

			migrations = append(migrations,
				migration{
					description: fmt.Sprintf("create %s table", pcn),
					script:      jen.Lit(strings.Join(scriptParts, "\n")),
				},
			)
		}
	case "mariadb", "maria_db":
		migrations = []migration{
			{
				description: "create users table",
				script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("CREATE TABLE IF NOT EXISTS users ("),
					jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
					jen.Lit("    `username` VARCHAR(150) NOT NULL,"),
					jen.Lit("    `hashed_password` VARCHAR(100) NOT NULL,"),
					jen.Lit("    `password_last_changed_on` INTEGER UNSIGNED,"),
					jen.Lit("    `two_factor_secret` VARCHAR(256) NOT NULL,"),
					jen.Lit("    `is_admin` BOOLEAN NOT NULL DEFAULT false,"),
					jen.Lit("    `created_on` BIGINT UNSIGNED,"),
					jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    PRIMARY KEY (`id`),"),
					jen.Lit("    UNIQUE (`username`)"),
					jen.Lit(");"),
				), jen.Lit("\n")),
			},
			{
				description: "create users table creation trigger",
				script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS users_creation_trigger BEFORE INSERT ON users FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
			{
				description: "create oauth2_clients table",
				script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("CREATE TABLE IF NOT EXISTS oauth2_clients ("),
					jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
					jen.Lit("    `name` VARCHAR(128) DEFAULT '',"),
					jen.Lit("    `client_id` VARCHAR(64) NOT NULL,"),
					jen.Lit("    `client_secret` VARCHAR(64) NOT NULL,"),
					jen.Lit("    `redirect_uri` VARCHAR(4096) DEFAULT '',"),
					jen.Lit("    `scopes` VARCHAR(4096) NOT NULL,"),
					jen.Lit("    `implicit_allowed` BOOLEAN NOT NULL DEFAULT false,"),
					jen.Lit("    `created_on` BIGINT UNSIGNED,"),
					jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `belongs_to` BIGINT UNSIGNED NOT NULL,"),
					jen.Lit("    PRIMARY KEY (`id`),"),
					jen.Lit("    FOREIGN KEY(`belongs_to`) REFERENCES users(`id`)"),
					jen.Lit(");"),
				), jen.Lit("\n")),
			},
			{
				description: "create oauth2_clients table creation trigger",
				script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS oauth2_clients_creation_trigger BEFORE INSERT ON oauth2_clients FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
			{
				description: "create webhooks table",
				script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("CREATE TABLE IF NOT EXISTS webhooks ("),
					jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
					jen.Lit("    `name` VARCHAR(128) NOT NULL,"),
					jen.Lit("    `content_type` VARCHAR(64) NOT NULL,"),
					jen.Lit("    `url` VARCHAR(4096) NOT NULL,"),
					jen.Lit("    `method` VARCHAR(32) NOT NULL,"),
					jen.Lit("    `events` VARCHAR(256) NOT NULL,"),
					jen.Lit("    `data_types` VARCHAR(256) NOT NULL,"),
					jen.Lit("    `topics` VARCHAR(256) NOT NULL,"),
					jen.Lit("    `created_on` BIGINT UNSIGNED,"),
					jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `belongs_to` BIGINT UNSIGNED NOT NULL,"),
					jen.Lit("    PRIMARY KEY (`id`),"),
					jen.Lit("    FOREIGN KEY (`belongs_to`) REFERENCES users(`id`)"),
					jen.Lit(");"),
				), jen.Lit("\n")),
			},
			{
				description: "create webhooks table creation trigger",
				script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
		}

		for _, typ := range types {
			pcn := typ.Name.PluralCommonName()

			scriptParts := []jen.Code{
				jen.Litf("CREATE TABLE IF NOT EXISTS %s (", typ.Name.PluralRouteName()),
				jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
			}

			for _, field := range typ.Fields {
				rn := field.Name.RouteName()
				query := fmt.Sprintf("    `%s` %s", rn, typeToMariaDBType(field.Type))
				if !field.Pointer {
					query += ` NOT NULL`
				}

				if field.DefaultAllowed {
					query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
				}
				scriptParts = append(scriptParts, jen.Lit(fmt.Sprintf("%s,", query)))
			}

			scriptParts = append(scriptParts,
				jen.Lit("    `created_on` BIGINT UNSIGNED,"),
				jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `belongs_to` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to`) REFERENCES users(`id`)"),
				jen.Lit(");"),
			)

			migrations = append(migrations,
				migration{
					description: fmt.Sprintf("create %s table", pcn),
					script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
						scriptParts...,
					), jen.Lit("\n")),
				},
				migration{
					description: fmt.Sprintf("create %s table creation trigger", pcn),
					script: jen.Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(
						jen.Litf("CREATE TRIGGER IF NOT EXISTS %s_creation_trigger BEFORE INSERT ON %s FOR EACH ROW", pcn, pcn),
						jen.Lit("BEGIN"),
						jen.Lit("  IF (new.created_on is null)"),
						jen.Lit("  THEN"),
						jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
						jen.Lit("  END IF;"),
						jen.Lit("END;"),
					), jen.Lit("\n")),
				},
			)
		}
	}

	for i, script := range migrations {
		out = append(out, jen.Valuesln(
			jen.ID("Version").Op(":").Lit(i+1),
			jen.ID("Description").Op(":").Lit(script.description),
			jen.ID("Script").Op(":").Add(script.script),
		))
	}

	return out
}

func migrationsDotGo(vendor wordsmith.SuperPalabra, types []models.DataType) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(ret)
	dbvsn := vendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	dbcn := vendor.SingularCommonName()

	dbrn := vendor.RouteName()

	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Var().Defs(
			jen.ID("migrations").Op("=").Index().Qual("github.com/GuiaBolso/darwin", "Migration").Valuesln(makeMigrations(vendor, types)...),
		),
	)

	var dialectName string
	if !isMariaDB {
		dialectName = fmt.Sprintf("%sDialect", dbvsn)
	} else {
		dialectName = "MySQLDialect"
	}

	ret.Add(
		jen.Comment("buildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Line(),
		jen.Commentf("migrate a %s database", dbcn),
		jen.Line(),
		jen.Func().ID("buildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Func().Params()).Block(
			jen.Return().Func().Params().Block(
				jen.ID("driver").Op(":=").Qual("github.com/GuiaBolso/darwin", "NewGenericDriver").Call(jen.ID("db"), jen.Qual("github.com/GuiaBolso/darwin", dialectName).Values()),
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
		jen.Comment("safe (as in idempotent, though not necessarily recommended) to call this function multiple times."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
			jen.ID(dbfl).Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.If(jen.Op("!").ID(dbfl).Dot("IsReady").Call(jen.ID("ctx"))).Block(
				jen.Return().ID("errors").Dot("New").Call(jen.Lit("db is not ready yet")),
			),
			jen.Line(),
			jen.ID(dbfl).Dot("migrateOnce").Dot("Do").Call(jen.ID("buildMigrationFunc").Call(jen.ID(dbfl).Dot("db"))),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return ret
}