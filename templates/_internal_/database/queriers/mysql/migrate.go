package mysql

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrateDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("driverName").Equals().Lit("instrumented-mysql"),
			jen.Newline(),
			jen.Comment("defaultTestUserTwoFactorSecret is the default TwoFactorSecret we give to test users when we initialize them."),
			jen.Comment("`otpauth://totp/todo:username?secret=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=&issuer=todo`"),
			jen.ID("defaultTestUserTwoFactorSecret").Equals().Lit("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
			jen.Newline(),
			jen.ID("testUserExistenceQuery").Equals().RawString(`
		SELECT users.id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.username = ? AND users.two_factor_secret_verified_on IS NOT NULL
	`),
			jen.Newline(),
			jen.ID("testUserCreationQuery").Equals().RawString(`
		INSERT INTO users (id,username,hashed_password,two_factor_secret,avatar_src,reputation,reputation_explanation,service_roles,two_factor_secret_verified_on,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
	`),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("Migrate is a simple wrapper around the core querier Migrate call."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8"), jen.ID("testUserConfig").Op("*").Qual(proj.TypesPackage(), "TestUserCreationConfig")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("q").Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.Newline(),
			jen.If(jen.Op("!").ID("q").Dot("IsReady").Call(
				jen.ID("ctx"),
				jen.ID("maxAttempts"),
			)).Body(
				jen.Return().Qual(proj.DatabasePackage(), "ErrDatabaseNotReady"),
			),
			jen.Newline(),
			jen.ID("q").Dot("migrateOnce").Dot("Do").Call(jen.ID("q").Dot("migrationFunc")),
			jen.Newline(),
			jen.If(jen.ID("testUserConfig").DoesNotEqual().Nil()).Body(
				jen.ID("q").Dot("logger").Dot("Debug").Call(jen.Lit("creating test user")),
				jen.Newline(),
				jen.ID("testUserExistenceArgs").Assign().Index().Interface().Values(jen.ID("testUserConfig").Dot("Username")),
				jen.Newline(),
				jen.ID("userRow").Assign().ID("q").Dot("getOneRow").Call(
					jen.ID("ctx"),
					jen.ID("q").Dot("db"),
					jen.Lit("user"),
					jen.ID("testUserExistenceQuery"),
					jen.ID("testUserExistenceArgs"),
				),
				jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dot("scanUser").Call(
					jen.ID("ctx"),
					jen.ID("userRow"),
					jen.False(),
				),
				jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
					jen.If(jen.ID("testUserConfig").Dot("ID").IsEqualTo().Lit("")).Body(
						jen.ID("testUserConfig").Dot("ID").Equals().ID("ksuid").Dot("New").Call().Dot("String").Call(),
					),
					jen.Newline(),
					jen.ID("testUserCreationArgs").Assign().Index().Interface().Valuesln(
						jen.ID("testUserConfig").Dot("ID"),
						jen.ID("testUserConfig").Dot("Username"),
						jen.ID("testUserConfig").Dot("HashedPassword"),
						jen.ID("defaultTestUserTwoFactorSecret"),
						jen.Lit(""),
						jen.Qual(proj.TypesPackage(), "GoodStandingAccountStatus"),
						jen.Lit(""),
						jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call(),
					),
					jen.Newline(),
					jen.Comment("these structs will be fleshed out by createUser"),
					jen.ID("user").Assign().Op("&").Qual(proj.TypesPackage(), "User").Valuesln(
						jen.ID("ID").MapAssign().ID("testUserConfig").Dot("ID"), jen.ID("Username").MapAssign().ID("testUserConfig").Dot("Username")),
					jen.ID("account").Assign().Op("&").Qual(proj.TypesPackage(), "Account").Valuesln(
						jen.ID("ID").MapAssign().ID("ksuid").Dot("New").Call().Dot("String").Call(),
					),
					jen.Newline(),
					jen.If(jen.ID("err").Equals().ID("q").Dot("createUser").Call(
						jen.ID("ctx"),
						jen.ID("user"),
						jen.ID("account"),
						jen.ID("testUserCreationQuery"),
						jen.ID("testUserCreationArgs"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
							jen.ID("err"),
							jen.ID("q").Dot("logger"),
							jen.ID("span"),
							jen.Lit("creating test user"),
						)),
					jen.ID("q").Dot("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("UsernameKey"),
						jen.ID("testUserConfig").Dot("Username"),
					).Dot("Debug").Call(jen.Lit("created test user and account")),
				),
			),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	code.Add(
		makeMigrations(proj),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Newline(),
		jen.Comment("migrate a postgres database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("migrationFunc").Params().Body(
			jen.ID("driver").Assign().ID("darwin").Dot("NewGenericDriver").Call(
				jen.ID("q").Dot("db"),
				jen.ID("darwin").Dot("MySQLDialect").Values(),
			),
			jen.If(jen.ID("err").Assign().ID("darwin").Dot("New").Call(
				jen.ID("driver"),
				jen.ID("migrations"),
				jen.Nil(),
			).Dot("Migrate").Call(), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("panic").Call(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("migrating database: %w"),
					jen.ID("err"),
				))),
		),
		jen.Newline(),
	)

	return code
}

type migration struct {
	description string
	script      jen.Code
}

func makeMigrations(proj *models.Project) jen.Code {
	out := []jen.Code{}

	migrations := makeMySQLMigrations(proj)

	for i, script := range migrations {
		out = append(out, jen.Valuesln(
			jen.ID("Version").MapAssign().Lit(float64(i+1)*.01),
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

func typeToMySQLType(t string) string {
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

func makeMySQLMigrations(proj *models.Project) []migration {
	migrations := []migration{
		{
			description: "create sessions table for session manager",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS sessions ("),
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
			description: "create users table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS users ("),
				jen.Lit("    `id` CHAR(27) NOT NULL,"),
				jen.Lit("    `username` VARCHAR(128) NOT NULL,"),
				jen.Lit("    `avatar_src` LONGTEXT NOT NULL,"),
				jen.Lit("    `hashed_password` VARCHAR(100) NOT NULL,"),
				jen.Lit("    `requires_password_change` BOOLEAN NOT NULL DEFAULT false,"),
				jen.Lit("    `password_last_changed_on` INTEGER UNSIGNED,"),
				jen.Lit("    `two_factor_secret` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `two_factor_secret_verified_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `service_roles` LONGTEXT NOT NULL,"),
				jen.Lit("    `reputation` VARCHAR(64) NOT NULL,"),
				jen.Lit("    `reputation_explanation` VARCHAR(1024) NOT NULL,"),
				jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    UNIQUE (`username`)"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create accounts table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS accounts ("),
				jen.Lit("    `id` CHAR(27) NOT NULL,"),
				jen.Lit("    `name` LONGTEXT NOT NULL,"),
				jen.Lit("    `billing_status` TEXT NOT NULL,"),
				jen.Lit("    `contact_email` TEXT NOT NULL,"),
				jen.Lit("    `contact_phone` TEXT NOT NULL,"),
				jen.Lit("    `payment_processor_customer_id` TEXT NOT NULL,"),
				jen.Lit("    `subscription_plan_id` VARCHAR(128),"),
				jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `belongs_to_user` CHAR(27) NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create account user memberships table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS account_user_memberships ("),
				jen.Lit("    `id` CHAR(27) NOT NULL,"),
				jen.Lit("    `belongs_to_account` CHAR(27) NOT NULL,"),
				jen.Lit("    `belongs_to_user` CHAR(27) NOT NULL,"),
				jen.Lit("    `default_account` BOOLEAN NOT NULL DEFAULT false,"),
				jen.Lit("    `account_roles` LONGTEXT NOT NULL,"),
				jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"),
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
			description: "create API clients table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS api_clients ("),
				jen.Lit("    `id` CHAR(27) NOT NULL,"),
				jen.Lit("    `name` VARCHAR(128),"),
				jen.Lit("    `client_id` VARCHAR(64) NOT NULL,"),
				jen.Lit("    `secret_key` BINARY(128) NOT NULL,"),
				jen.Lit("    `for_admin` BOOLEAN NOT NULL DEFAULT false,"),
				jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `belongs_to_user` CHAR(27) NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    UNIQUE (`name`, `belongs_to_user`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
		{
			description: "create webhooks table",
			script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
				jen.Lit("CREATE TABLE IF NOT EXISTS webhooks ("),
				jen.Lit("    `id` CHAR(27) NOT NULL,"),
				jen.Lit("    `name` VARCHAR(128) NOT NULL,"),
				jen.Lit("    `content_type` VARCHAR(64) NOT NULL,"),
				jen.Lit("    `url` LONGTEXT NOT NULL,"),
				jen.Lit("    `method` VARCHAR(8) NOT NULL,"),
				jen.Lit("    `events` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `data_types` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `topics` VARCHAR(256) NOT NULL,"),
				jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"),
				jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
				jen.Lit("    `belongs_to_account` CHAR(27) NOT NULL,"),
				jen.Lit("    PRIMARY KEY (`id`),"),
				jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE"),
				jen.Lit(");"),
			), jen.Lit("\n")),
		},
	}

	for _, typ := range proj.DataTypes {
		pcn := typ.Name.PluralCommonName()
		tableName := typ.Name.PluralRouteName()

		scriptParts := []string{
			"    `id` CHAR(27) NOT NULL",
		}

		for _, field := range typ.Fields {
			rn := field.Name.RouteName()

			query := fmt.Sprintf("    `%s` %s", rn, typeToMySQLType(field.Type))
			if !field.IsPointer {
				query += ` NOT NULL`
			}

			scriptParts = append(scriptParts, query)
		}

		scriptParts = append(scriptParts,
			"    `created_on` BIGINT UNSIGNED NOT NULL",
			"    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL",
			"    `archived_on` BIGINT UNSIGNED DEFAULT NULL",
		)

		if typ.BelongsToAccount {
			scriptParts = append(scriptParts,
				"    `belongs_to_account` CHAR(27) NOT NULL",
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				fmt.Sprintf("    `belongs_to_%s` CHAR(27) NOT NULL", typ.BelongsToStruct.RouteName()),
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
