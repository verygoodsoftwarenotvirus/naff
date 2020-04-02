package queriers

import (
	"fmt"
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
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"int8":     "INTEGER",
		"*int8":    "INTEGER",
		"int16":    "INTEGER",
		"*int16":   "INTEGER",
		"int32":    "INTEGER",
		"*int32":   "INTEGER",
		"int64":    "INTEGER",
		"*int64":   "INTEGER",
		"uint":     "INTEGER",
		"*uint":    "INTEGER",
		"uint8":    "INTEGER",
		"*uint8":   "INTEGER",
		"uint16":   "INTEGER",
		"*uint16":  "INTEGER",
		"uint32":   "BIGINT",
		"*uint32":  "BIGINT",
		"uint64":   "BIGINT",
		"*uint64":  "BIGINT",
		"float32":  "DOUBLE PRECISION",
		"*float32": "DOUBLE PRECISION",
		"float64":  "DOUBLE PRECISION",
		"*float64": "DOUBLE PRECISION",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	panic(fmt.Sprintf("unknown type!: %q", t))
}

func typeToSqliteType(t string) string {
	typeMap := map[string]string{
		"[]string": "CHARACTER VARYING",
		"string":   "CHARACTER VARYING",
		"*string":  "CHARACTER VARYING",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"int8":     "INTEGER",
		"*int8":    "INTEGER",
		"int16":    "INTEGER",
		"*int16":   "INTEGER",
		"int32":    "INTEGER",
		"*int32":   "INTEGER",
		"int64":    "INTEGER",
		"*int64":   "INTEGER",
		"uint":     "INTEGER",
		"*uint":    "INTEGER",
		"uint8":    "INTEGER",
		"*uint8":   "INTEGER",
		"uint16":   "INTEGER",
		"*uint16":  "INTEGER",
		"uint32":   "INTEGER",
		"*uint32":  "INTEGER",
		"uint64":   "INTEGER",
		"*uint64":  "INTEGER",
		"float32":  "REAL",
		"*float32": "REAL",
		"float64":  "REAL",
		"*float64": "REAL",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	panic(fmt.Sprintf("unknown type!: %q", t))
}

func typeToMariaDBType(t string) string {
	typeMap := map[string]string{
		"[]string": "LONGTEXT",
		"string":   "LONGTEXT",
		"*string":  "LONGTEXT",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"int8":     "INTEGER",
		"*int8":    "INTEGER",
		"int16":    "INTEGER",
		"*int16":   "INTEGER",
		"int32":    "INTEGER",
		"*int32":   "INTEGER",
		"int64":    "INTEGER",
		"*int64":   "INTEGER",
		"uint":     "INTEGER UNSIGNED",
		"*uint":    "INTEGER UNSIGNED",
		"uint8":    "INTEGER UNSIGNED",
		"*uint8":   "INTEGER UNSIGNED",
		"uint16":   "INTEGER UNSIGNED",
		"*uint16":  "INTEGER UNSIGNED",
		"uint32":   "INTEGER UNSIGNED",
		"*uint32":  "INTEGER UNSIGNED",
		"uint64":   "INTEGER UNSIGNED",
		"*uint64":  "INTEGER UNSIGNED",
		"float32":  "DOUBLE PRECISION",
		"*float32": "DOUBLE PRECISION",
		"float64":  "DOUBLE PRECISION",
		"*float64": "DOUBLE PRECISION",
	}

	if x, ok := typeMap[strings.TrimSpace(t)]; ok {
		return x
	}

	panic(fmt.Sprintf("unknown type!: %q", t))
}

type migration struct {
	description string
	script      jen.Code
}

func makePostgresMigrations(proj *models.Project) []migration {
	migrations := []migration{
		{
			description: "create users table",
			script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS users (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"password_last_changed_on" integer,
				"two_factor_secret" TEXT NOT NULL,
				"is_admin" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				UNIQUE ("username")
			);`),
		},
		{
			description: "create oauth2_clients table",
			script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_user") REFERENCES "users"("id")
			);`),
		},
		{
			description: "create webhooks table",
			script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS webhooks (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" TEXT NOT NULL,
				"content_type" TEXT NOT NULL,
				"url" TEXT NOT NULL,
				"method" TEXT NOT NULL,
				"events" TEXT NOT NULL,
				"data_types" TEXT NOT NULL,
				"topics" TEXT NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_user") REFERENCES "users"("id")
			);`),
		},
	}

	for _, typ := range proj.DataTypes {
		pcn := typ.Name.PluralCommonName()

		scriptParts := []string{
			fmt.Sprintf("\n			CREATE TABLE IF NOT EXISTS %s (", typ.Name.PluralRouteName()),
			`				"id" BIGSERIAL NOT NULL PRIMARY KEY,`,
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
			`				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),`,
			`				"updated_on" BIGINT DEFAULT NULL,`,
		)

		if typ.BelongsToUser {
			scriptParts = append(scriptParts,
				`				"archived_on" BIGINT DEFAULT NULL,`,
				`				"belongs_to_user" BIGINT NOT NULL,`,
				`				FOREIGN KEY ("belongs_to_user") REFERENCES "users"("id")`,
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				`				"archived_on" BIGINT DEFAULT NULL,`,
				fmt.Sprintf(`				"belongs_to_%s" BIGINT NOT NULL,`, typ.BelongsToStruct.RouteName()),
				fmt.Sprintf(`				FOREIGN KEY ("belongs_to_%s") REFERENCES "%s"("id")`, typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
			)
		} else if typ.BelongsToNobody {
			scriptParts = append(scriptParts, `				"archived_on" BIGINT DEFAULT NULL`)
		}

		scriptParts = append(scriptParts,
			`			);`,
		)

		migrations = append(migrations,
			migration{
				description: fmt.Sprintf("create %s table", pcn),
				script:      jen.Lit(strings.Join(scriptParts, "\n")),
			},
		)
	}
	return migrations
}

func makeMariaDBMigrations(proj *models.Project) []migration {
	migrations := []migration{
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
				jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY(`belongs_to_user`) REFERENCES users(`id`)"),
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
				jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`)"),
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

	for _, typ := range proj.DataTypes {
		pcn := typ.Name.PluralCommonName()
		tableName := typ.Name.PluralRouteName()

		scriptParts := []jen.Code{
			jen.Litf("CREATE TABLE IF NOT EXISTS %s (", tableName),
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
		)

		if typ.BelongsToUser {
			scriptParts = append(scriptParts,
				jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"),
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				jen.Litf("    `belongs_to_%s` BIGINT UNSIGNED NOT NULL,", typ.BelongsToStruct.RouteName()),
			)
		}

		scriptParts = append(scriptParts,
			jen.Lit("    PRIMARY KEY (`id`),"),
		)

		if typ.BelongsToUser {
			scriptParts = append(scriptParts,
				jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`)"),
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				jen.Litf("    FOREIGN KEY (`belongs_to_%s`) REFERENCES %s(`id`)", typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
			)
		}

		scriptParts = append(scriptParts,
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
					jen.Litf("CREATE TRIGGER IF NOT EXISTS %s_creation_trigger BEFORE INSERT ON %s FOR EACH ROW", tableName, tableName),
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

	return migrations
}

func makeSqliteMigrations(proj *models.Project) []migration {
	migrations := []migration{
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
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);`),
		},
		{
			description: "create webhooks table",
			script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS webhooks (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT NOT NULL,
				"content_type" TEXT NOT NULL,
				"url" TEXT NOT NULL,
				"method" TEXT NOT NULL,
				"events" TEXT NOT NULL,
				"data_types" TEXT NOT NULL,
				"topics" TEXT NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);`),
		},
	}

	idType := "INTEGER"
	idAddendum := " AUTOINCREMENT"

	for _, typ := range proj.DataTypes {
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
		)

		if typ.BelongsToUser {
			scriptParts = append(scriptParts,
				`				"belongs_to_user" INTEGER NOT NULL,`,
				`				FOREIGN KEY(belongs_to_user) REFERENCES users(id)`,
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				fmt.Sprintf(`				"belongs_to_%s" INTEGER NOT NULL,`, typ.BelongsToStruct.RouteName()),
				fmt.Sprintf(`				FOREIGN KEY(belongs_to_%s) REFERENCES %s(id)`, typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
			)
		}

		scriptParts = append(scriptParts,
			`			);`,
		)

		migrations = append(migrations,
			migration{
				description: fmt.Sprintf("create %s table", pcn),
				script:      jen.Lit(strings.Join(scriptParts, "\n")),
			},
		)
	}
	return migrations
}

func makeMigrations(proj *models.Project, dbVendor wordsmith.SuperPalabra) []jen.Code {
	var (
		out        []jen.Code
		migrations []migration
	)

	dbrn := strings.ToLower(dbVendor.RouteName())

	switch dbrn {
	case "postgres":
		migrations = makePostgresMigrations(proj)
	case "sqlite":
		migrations = makeSqliteMigrations(proj)
	case "mariadb", "maria_db":
		migrations = makeMariaDBMigrations(proj)
	}

	for i, script := range migrations {
		out = append(out, jen.Valuesln(
			jen.ID("Version").MapAssign().Lit(i+1),
			jen.ID("Description").MapAssign().Lit(script.description),
			jen.ID("Script").MapAssign().Add(script.script),
		))
	}

	return out
}

func migrationsDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("migrations").Equals().Index().Qual("github.com/GuiaBolso/darwin", "Migration").Valuesln(makeMigrations(proj, dbvendor)...),
		),
	)

	ret.Add(
		buildBuildMigrationFuncDecl(dbvendor)...,
	)

	ret.Add(
		buildMigrate(dbvendor)...,
	)

	return ret
}

func buildBuildMigrationFuncDecl(dbvendor wordsmith.SuperPalabra) []jen.Code {
	dbcn := dbvendor.SingularCommonName()
	dbvsn := dbvendor.Singular()
	isMariaDB := dbvendor.RouteName() == "mariadb" || dbvendor.RouteName() == "maria_db"

	var dialectName string
	if !isMariaDB {
		dialectName = fmt.Sprintf("%sDialect", dbvsn)
	} else {
		dialectName = "MySQLDialect"
	}

	return []jen.Code{
		jen.Comment("buildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Line(),
		jen.Commentf("migrate a %s database", dbcn),
		jen.Line(),
		jen.Func().ID("buildMigrationFunc").Params(jen.ID("db").ParamPointer().Qual("database/sql", "DB")).Params(jen.Func().Params()).Block(
			jen.Return().Func().Params().Block(
				jen.ID("driver").Assign().Qual("github.com/GuiaBolso/darwin", "NewGenericDriver").Call(jen.ID("db"), jen.Qual("github.com/GuiaBolso/darwin", dialectName).Values()),
				jen.If(jen.Err().Assign().Qual("github.com/GuiaBolso/darwin", "New").Call(jen.ID("driver"), jen.ID("migrations"), jen.Nil()).Dot("Migrate").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("panic").Call(jen.Err()),
				),
			),
		),
		jen.Line(),
	}
}

func buildMigrate(dbvendor wordsmith.SuperPalabra) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	return []jen.Code{
		jen.Comment("Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be"),
		jen.Line(),
		jen.Comment("safe (as in idempotent, though not necessarily recommended) to call this function multiple times."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).ID("Migrate").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("error")).Block(
			jen.ID(dbfl).Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.If(jen.Op("!").ID(dbfl).Dot("IsReady").Call(utils.CtxVar())).Block(
				jen.Return().ID("errors").Dot("New").Call(jen.Lit("db is not ready yet")),
			),
			jen.Line(),
			jen.ID(dbfl).Dot("migrateOnce").Dot("Do").Call(jen.ID("buildMigrationFunc").Call(jen.ID(dbfl).Dot("db"))),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}
}
