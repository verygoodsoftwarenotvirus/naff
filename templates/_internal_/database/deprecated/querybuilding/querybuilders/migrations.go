package querybuilders

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrationsDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabasePackage("queriers", "v1", spn), spn)

	utils.AddImports(proj, code, false)

	if dbvendor.SingularPackageName() == "mysql" {
		code.Add(buildBuildCreationTriggerScript()...)
	}

	code.Add(buildMigrationVarDeclarations(proj, dbvendor))
	code.Add(buildBuildMigrationFuncDecl(dbvendor)...)

	return code
}

func buildBuildCreationTriggerScript() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildCreationTriggerScript").Params(jen.ID("tableName").String()).Params(jen.String()).Body(
			jen.Return().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("CREATE TRIGGER IF NOT EXISTS %s_creation_trigger BEFORE INSERT ON %s FOR EACH ROW"), jen.ID("tableName"), jen.ID("tableName")),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP(NOW());"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				),
				jen.Lit("\n"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildMigrationVarDeclarations(proj *models.Project, dbvendor wordsmith.SuperPalabra) jen.Code {
	lines := makeMigrations(proj, dbvendor)

	return lines
}

func typeToPostgresType(t string) string {
	typeMap := map[string]string{
		//"[]string": "TEXT",
		"string":   "TEXT",
		"*string":  "TEXT",
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
		//"[]string": "TEXT",
		"string":   "TEXT",
		"*string":  "TEXT",
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
		//"[]string": "LONGTEXT",
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
			description: "create sessions table for session manager",
			script: jen.RawString(`
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BYTEA NOT NULL,
				expiry TIMESTAMPTZ NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())
			);`),
		},
		{
			description: "create sessions table for session manager",
			script:      jen.RawString(`CREATE INDEX sessions_expiry_idx ON sessions (expiry);`),
		},
		{
			description: "create audit log table",
			script: jen.RawString(`
			CREATE TABLE IF NOT EXISTS audit_log (
				id BIGSERIAL NOT NULL PRIMARY KEY,
				external_id TEXT NOT NULL,
				event_type TEXT NOT NULL,
				context JSONB NOT NULL,
				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())
			);`),
		},
		{
			description: "create users table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create accounts table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create account user memberships table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create API clients table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create webhooks table",
			script: jen.RawString(`
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
				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
			);`),
		},
	}

	for _, typ := range proj.DataTypes {
		pcn := typ.Name.PluralCommonName()

		scriptParts := []string{
			fmt.Sprintf("\n			CREATE TABLE IF NOT EXISTS %s (", typ.Name.PluralRouteName()),
			`				id BIGSERIAL NOT NULL PRIMARY KEY`,
			`				external_id TEXT NOT NULL`,
		}

		for _, field := range typ.Fields {
			rn := field.Name.RouteName()

			query := fmt.Sprintf("				%s %s", rn, typeToPostgresType(field.Type))
			if !field.IsPointer {
				query += ` NOT NULL`
			}
			if field.DefaultValue != "" {
				query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
			}

			scriptParts = append(scriptParts, fmt.Sprintf("%s", query))
		}

		scriptParts = append(scriptParts,
			`				created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())`,
			`				last_updated_on BIGINT DEFAULT NULL`,
		)

		if !typ.BelongsToAccount && typ.BelongsToStruct == nil {
			scriptParts = append(scriptParts,
				`				archived_on BIGINT DEFAULT NULL`,
			)
		} else {
			scriptParts = append(scriptParts,
				`				archived_on BIGINT DEFAULT NULL`, // note the comma
			)
		}

		if typ.BelongsToAccount {
			scriptParts = append(scriptParts,
				`				belongs_to_account BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE`,
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				fmt.Sprintf(`				"belongs_to_%s" BIGINT NOT NULL REFERENCES %s(id) ON DELETE CASCADE`, typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
			)
		}

		for i, sp := range scriptParts {
			if i != len(scriptParts)-1 {
				if !strings.HasSuffix(sp, "(") {
					scriptParts[i] = fmt.Sprintf("%s,", sp)
				}
			}
		}

		scriptParts = append(scriptParts,
			`			);`,
		)

		migrations = append(migrations,
			migration{
				description: fmt.Sprintf("create %s table", pcn),
				script:      jen.RawString(strings.Join(scriptParts, "\n")),
			},
		)
	}
	return migrations
}

func makeMariaDBMigrations(proj *models.Project) []migration {
	migrations := []migration{
		{
			description: "create sessions table for session manager",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE sessions ("),
				jen.Lit("`token` CHAR(43) PRIMARY KEY,"),
				jen.Lit("`data` BLOB NOT NULL,"),
				jen.Lit("`expiry` TIMESTAMP(6) NOT NULL,"),
				jen.Lit("`created_on` BIGINT UNSIGNED"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create sessions table for session manager",
			script:      jen.Lit("CREATE INDEX sessions_expiry_idx ON sessions (expiry);"),
		},
		{
			description: "create sessions table creation trigger",
			script:      jen.ID("buildCreationTriggerScript").Call(jen.Lit("sessions")),
		},
		{
			description: "create audit log table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS audit_log ("),
				jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
				jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"),
				jen.Lit("    `event_type` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `context` JSON NOT NULL,"),
				jen.Lit("    `created_on` BIGINT UNSIGNED,"),
				jen.Lit("    PRIMARY KEY (`id`)"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create audit log table creation trigger",
			script:      jen.ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")),
		},
		{
			description: "create users table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS users ("),
				jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
				jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"),
				jen.Lit("    `username` VARCHAR(128) NOT NULL,"),
				jen.Lit("    `avatar_src` LONGTEXT NOT NULL DEFAULT '',"),
				jen.Lit("    `hashed_password` VARCHAR(100) NOT NULL,"),
				jen.Lit("    `requires_password_change` BOOLEAN NOT NULL DEFAULT false,"),
				jen.Lit("    `password_last_changed_on` INTEGER UNSIGNED,"),
				jen.Lit("    `two_factor_secret` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `two_factor_secret_verified_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `service_roles` LONGTEXT NOT NULL DEFAULT 'service_user',"),
				jen.Lit("    `reputation` VARCHAR(64) NOT NULL DEFAULT 'unverified',"),
				jen.Lit("    `reputation_explanation` VARCHAR(1024) NOT NULL DEFAULT '',"),
				jen.Lit("    `created_on` BIGINT UNSIGNED,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    UNIQUE (`username`)"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create users table creation trigger",
			script:      jen.ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")),
		},
		{
			description: "create accounts table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS accounts ("),
				jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
				jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"),
				jen.Lit("    `name` LONGTEXT NOT NULL,"),
				jen.Lit("    `billing_status` TEXT NOT NULL DEFAULT 'unpaid',"),
				jen.Lit("    `contact_email` TEXT NOT NULL DEFAULT '',"),
				jen.Lit("    `contact_phone` TEXT NOT NULL DEFAULT '',"),
				jen.Lit("    `payment_processor_customer_id` TEXT NOT NULL DEFAULT '',"),
				jen.Lit("    `subscription_plan_id` VARCHAR(128),"),
				jen.Lit("    `created_on` BIGINT UNSIGNED,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create accounts table creation trigger",
			script:      jen.ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")),
		},
		{
			description: "create account user memberships table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS account_user_memberships ("),
				jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
				jen.Lit("    `belongs_to_account` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    `default_account` BOOLEAN NOT NULL DEFAULT false,"),
				jen.Lit("    `account_roles` LONGTEXT NOT NULL DEFAULT 'account_user',"),
				jen.Lit("    `created_on` BIGINT UNSIGNED,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE,"),
				jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE,"),
				jen.Lit("    UNIQUE (`belongs_to_account`, `belongs_to_user`)"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create accounts membership creation trigger",
			script:      jen.ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")),
		},
		{
			description: "create API clients table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS api_clients ("),
				jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
				jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"),
				jen.Lit("    `name` VARCHAR(128) DEFAULT '',"),
				jen.Lit("    `client_id` VARCHAR(64) NOT NULL,"),
				jen.Lit("    `secret_key` BINARY(128) NOT NULL,"),
				jen.Lit("    `account_roles` LONGTEXT NOT NULL DEFAULT 'account_member',"),
				jen.Lit("    `for_admin` BOOLEAN NOT NULL DEFAULT false,"),
				jen.Lit("    `created_on` BIGINT UNSIGNED,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    UNIQUE (`name`, `belongs_to_user`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create API clients table creation trigger",
			script:      jen.ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildingPackage(), "APIClientsTableName")),
		},
		{
			description: "create webhooks table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS webhooks ("),
				jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
				jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"),
				jen.Lit("    `name` VARCHAR(128) NOT NULL,"),
				jen.Lit("    `content_type` VARCHAR(64) NOT NULL,"),
				jen.Lit("    `url` LONGTEXT NOT NULL,"),
				jen.Lit("    `method` VARCHAR(8) NOT NULL,"),
				jen.Lit("    `events` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `data_types` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `topics` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `created_on` BIGINT UNSIGNED,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `belongs_to_account` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create webhooks table creation trigger",
			script:      jen.ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildingPackage(), "WebhooksTableName")),
		},
	}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		pcn := typ.Name.PluralCommonName()
		tableName := typ.Name.PluralRouteName()

		scriptParts := []string{
			"    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT",
			"    `external_id` VARCHAR(36) NOT NULL",
		}

		for _, field := range typ.Fields {
			rn := field.Name.RouteName()

			query := fmt.Sprintf("    `%s` %s", rn, typeToMariaDBType(field.Type))
			if !field.IsPointer {
				query += ` NOT NULL`
			}
			if field.DefaultValue != "" {
				query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
			}

			scriptParts = append(scriptParts, query)
		}

		scriptParts = append(scriptParts,
			"    `created_on` BIGINT UNSIGNED",
			"    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL",
			"    `archived_on` BIGINT UNSIGNED DEFAULT NULL",
		)

		if typ.BelongsToAccount {
			scriptParts = append(scriptParts,
				"    `belongs_to_account` BIGINT UNSIGNED NOT NULL",
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				fmt.Sprintf("    `belongs_to_%s` BIGINT UNSIGNED NOT NULL", typ.BelongsToStruct.RouteName()),
			)
		}

		scriptParts = append(scriptParts,
			"    PRIMARY KEY (`id`)",
		)

		if typ.BelongsToAccount {
			if typ.BelongsToStruct != nil {
				scriptParts = append(scriptParts,
					"    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE",
				)
			} else {
				scriptParts = append(scriptParts,
					"    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE",
				)
			}
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				fmt.Sprintf("    FOREIGN KEY (`belongs_to_%s`) REFERENCES %s(`id`) ON DELETE CASCADE", typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
			)
		}

		migrations = append(migrations,
			migration{
				description: fmt.Sprintf("create %s table", pcn),
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					codeFromStrings(scriptParts, tableName)...,
				), jen.Lit("\n")),
			},
			migration{
				description: fmt.Sprintf("create %s table creation trigger", pcn),
				script:      jen.ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))),
			},
		)
	}

	return migrations
}

func codeFromStrings(in []string, tableName string) []jen.Code {
	out := []jen.Code{
		jen.Litf("CREATE TABLE IF NOT EXISTS %s (", tableName),
	}

	for i, x := range in {
		if i == len(in)-1 {
			out = append(out, jen.Lit(x))
		} else {
			out = append(out, jen.Litf("%s,", x))
		}
	}

	out = append(out, jen.Litf(");"))

	return out
}

func makeSqliteMigrations(proj *models.Project) []migration {
	migrations := []migration{
		{
			description: "create sessions table for session manager",
			script: jen.RawString(`
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BLOB NOT NULL,
				expiry REAL NOT NULL,
				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now'))
			);`),
		},
		{
			description: "create sessions table for session manager",
			script:      jen.RawString("CREATE INDEX sessions_expiry_idx ON sessions(expiry);"),
		},
		{
			description: "create audit log table",
			script: jen.RawString(`
			CREATE TABLE IF NOT EXISTS audit_log (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				external_id TEXT NOT NULL,
				event_type TEXT NOT NULL,
				context JSON NOT NULL,
				created_on BIGINT NOT NULL DEFAULT (strftime('%s','now'))
			);`),
		},
		{
			description: "create users table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create accounts table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create account user memberships table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create API clients table",
			script: jen.RawString(`
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
			);`),
		},
		{
			description: "create webhooks table",
			script: jen.RawString(`
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
			);`),
		},
	}

	for _, typ := range proj.DataTypes {
		pcn := typ.Name.PluralCommonName()

		scriptParts := []string{
			"				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT",
			"				external_id TEXT NOT NULL",
		}

		for _, field := range typ.Fields {
			rn := field.Name.RouteName()

			query := fmt.Sprintf("				%s %s", rn, typeToSqliteType(field.Type))
			if !field.IsPointer {
				query += ` NOT NULL`
			}
			if field.DefaultValue != "" {
				query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
			}

			scriptParts = append(scriptParts, fmt.Sprintf("%s", query))
		}

		scriptParts = append(scriptParts,
			`				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now'))`,
			`				last_updated_on INTEGER DEFAULT NULL`,
			`				archived_on INTEGER DEFAULT NULL`,
		)

		if typ.BelongsToAccount {
			scriptParts = append(scriptParts,
				`				belongs_to_account INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE`,
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				fmt.Sprintf(`				belongs_to_%s INTEGER NOT NULL`, typ.BelongsToStruct.RouteName()),
				fmt.Sprintf(`				FOREIGN KEY(belongs_to_%s) REFERENCES %s(id)`, typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
			)
		}

		migrations = append(migrations,
			migration{
				description: fmt.Sprintf("create %s table", pcn),
				script: jen.RawString(fmt.Sprintf("\n			CREATE TABLE IF NOT EXISTS %s (\n%s\n\t\t\t);", typ.Name.PluralRouteName(), strings.Join(scriptParts, ",\n"))),
			},
		)
	}
	return migrations
}

func makeMigrations(proj *models.Project, dbVendor wordsmith.SuperPalabra) jen.Code {
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
	case "mysql", "maria_db":
		migrations = makeMariaDBMigrations(proj)
	}

	for i, script := range migrations {
		out = append(out, jen.Valuesln(
			jen.ID("Version").MapAssign().Lit(float64(i)*.01),
			jen.ID("Description").MapAssign().Lit(script.description),
			jen.ID("Script").MapAssign().Add(script.script),
		))
	}

	return jen.Var().Defs(
		jen.ID("migrations").Equals().Index().Qual("github.com/GuiaBolso/darwin", "Migration").Valuesln(
			out...,
		),
	)
}

func buildBuildMigrationFuncDecl(dbvendor wordsmith.SuperPalabra) []jen.Code {
	dbcn := dbvendor.SingularCommonName()
	dbvsn := dbvendor.Singular()
	isMariaDB := dbvendor.RouteName() == "mysql" || dbvendor.RouteName() == "maria_db"

	var dialectName string
	if !isMariaDB {
		dialectName = fmt.Sprintf("%sDialect", dbvsn)
	} else {
		dialectName = "MySQLDialect"
	}

	return []jen.Code{
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Newline(),
		jen.Commentf("migrate a %s database.", dbcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvsn)).ID("BuildMigrationFunc").Params(jen.ID("db").PointerTo().Qual("database/sql", "DB")).Params(jen.Func().Params()).Body(
			jen.Return().Func().Params().Body(
				jen.ID("driver").Assign().Qual("github.com/GuiaBolso/darwin", "NewGenericDriver").Call(jen.ID("db"), jen.Qual("github.com/GuiaBolso/darwin", dialectName).Values()),
				jen.If(jen.Err().Assign().Qual("github.com/GuiaBolso/darwin", "New").Call(jen.ID("driver"), jen.ID("migrations"), jen.Nil()).Dot("Migrate").Call(), jen.Err().DoesNotEqual().Nil()).Body(
					jen.ID("panic").Call(jen.Qual("fmt", "Errorf").Call(jen.Lit("migrating database: %w"), jen.Err())),
				),
			),
		),
		jen.Newline(),
	}
}
