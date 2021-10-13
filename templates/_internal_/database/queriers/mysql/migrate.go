package mysql

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrateDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("driverName").Equals().Lit("instrumented-mysql"),
			jen.ID("defaultTestUserTwoFactorSecret").Equals().Lit("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("testUserExistenceQuery").Equals().Lit(`
	SELECT users.id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.username = ? AND users.two_factor_secret_verified_on IS NOT NULL
`),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("testUserCreationQuery").Equals().Lit(`
	INSERT INTO users (id,username,hashed_password,two_factor_secret,avatar_src,reputation,reputation_explanation,service_roles,two_factor_secret_verified_on,created_on) VALUES (?,?,?,?,?,?,?,?,UNIX_TIMESTAMP(),UNIX_TIMESTAMP())
`),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("Migrate is a simple wrapper around the core querier Migrate call."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8"), jen.ID("testUserConfig").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("q").Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.If(jen.Op("!").ID("q").Dot("IsReady").Call(
				jen.ID("ctx"),
				jen.ID("maxAttempts"),
			)).Body(
				jen.Return().ID("database").Dot("ErrDatabaseNotReady")),
			jen.ID("q").Dot("migrateOnce").Dot("Do").Call(jen.ID("q").Dot("migrationFunc")),
			jen.If(jen.ID("testUserConfig").DoesNotEqual().Nil()).Body(
				jen.ID("q").Dot("logger").Dot("Debug").Call(jen.Lit("creating test user")),
				jen.ID("testUserExistenceArgs").Op(":=").Index().Interface().Valuesln(jen.ID("testUserConfig").Dot("Username")),
				jen.ID("userRow").Op(":=").ID("q").Dot("getOneRow").Call(
					jen.ID("ctx"),
					jen.ID("q").Dot("db"),
					jen.Lit("user"),
					jen.ID("testUserExistenceQuery"),
					jen.ID("testUserExistenceArgs"),
				),
				jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanUser").Call(
					jen.ID("ctx"),
					jen.ID("userRow"),
					jen.ID("false"),
				),
				jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
					jen.If(jen.ID("testUserConfig").Dot("ID").Op("==").Lit("")).Body(
						jen.ID("testUserConfig").Dot("ID").Equals().ID("ksuid").Dot("New").Call().Dot("String").Call()),
					jen.ID("testUserCreationArgs").Op(":=").Index().Interface().Valuesln(jen.ID("testUserConfig").Dot("ID"), jen.ID("testUserConfig").Dot("Username"), jen.ID("testUserConfig").Dot("HashedPassword"), jen.ID("defaultTestUserTwoFactorSecret"), jen.Lit(""), jen.ID("types").Dot("GoodStandingAccountStatus"), jen.Lit(""), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.ID("user").Op(":=").Op("&").ID("types").Dot("User").Valuesln(jen.ID("ID").MapAssign().ID("testUserConfig").Dot("ID"), jen.ID("Username").MapAssign().ID("testUserConfig").Dot("Username")),
					jen.ID("account").Op(":=").Op("&").ID("types").Dot("Account").Valuesln(jen.ID("ID").MapAssign().ID("ksuid").Dot("New").Call().Dot("String").Call()),
					jen.If(jen.ID("err").Equals().ID("q").Dot("createUser").Call(
						jen.ID("ctx"),
						jen.ID("user"),
						jen.ID("account"),
						jen.ID("testUserCreationQuery"),
						jen.ID("testUserCreationArgs"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
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
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("migrations").Equals().Index().ID("darwin").Dot("Migration").Valuesln(jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.01), jen.ID("Description").MapAssign().Lit("create sessions table for session manager"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS sessions ("), jen.Lit("`token` CHAR(43) PRIMARY KEY,"), jen.Lit("`data` BLOB NOT NULL,"), jen.Lit("`expiry` TIMESTAMP(6) NOT NULL,"), jen.Lit("`created_on` BIGINT UNSIGNED"), jen.Lit(");")),
				jen.Lit("\n"),
			)), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.02), jen.ID("Description").MapAssign().Lit("create sessions table for session manager"), jen.ID("Script").MapAssign().Lit("CREATE INDEX sessions_expiry_idx ON sessions (expiry);")), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.03), jen.ID("Description").MapAssign().Lit("create users table"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS users ("), jen.Lit("    `id` CHAR(27) NOT NULL,"), jen.Lit("    `username` VARCHAR(128) NOT NULL,"), jen.Lit("    `avatar_src` LONGTEXT NOT NULL,"), jen.Lit("    `hashed_password` VARCHAR(100) NOT NULL,"), jen.Lit("    `requires_password_change` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `password_last_changed_on` INTEGER UNSIGNED,"), jen.Lit("    `two_factor_secret` VARCHAR(256) NOT NULL,"), jen.Lit("    `two_factor_secret_verified_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `service_roles` LONGTEXT NOT NULL,"), jen.Lit("    `reputation` VARCHAR(64) NOT NULL,"), jen.Lit("    `reputation_explanation` VARCHAR(1024) NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    UNIQUE (`username`)"), jen.Lit(");")),
				jen.Lit("\n"),
			)), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.04), jen.ID("Description").MapAssign().Lit("create accounts table"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS accounts ("), jen.Lit("    `id` CHAR(27) NOT NULL,"), jen.Lit("    `name` LONGTEXT NOT NULL,"), jen.Lit("    `billing_status` TEXT NOT NULL,"), jen.Lit("    `contact_email` TEXT NOT NULL,"), jen.Lit("    `contact_phone` TEXT NOT NULL,"), jen.Lit("    `payment_processor_customer_id` TEXT NOT NULL,"), jen.Lit("    `subscription_plan_id` VARCHAR(128),"), jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_user` CHAR(27) NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"), jen.Lit(");")),
				jen.Lit("\n"),
			)), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.05), jen.ID("Description").MapAssign().Lit("create account user memberships table"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS account_user_memberships ("), jen.Lit("    `id` CHAR(27) NOT NULL,"), jen.Lit("    `belongs_to_account` CHAR(27) NOT NULL,"), jen.Lit("    `belongs_to_user` CHAR(27) NOT NULL,"), jen.Lit("    `default_account` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `account_roles` LONGTEXT NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE,"), jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE,"), jen.Lit("    UNIQUE (`belongs_to_account`, `belongs_to_user`)"), jen.Lit(");")),
				jen.Lit("\n"),
			)), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.06), jen.ID("Description").MapAssign().Lit("create API clients table"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS api_clients ("), jen.Lit("    `id` CHAR(27) NOT NULL,"), jen.Lit("    `name` VARCHAR(128),"), jen.Lit("    `client_id` VARCHAR(64) NOT NULL,"), jen.Lit("    `secret_key` BINARY(128) NOT NULL,"), jen.Lit("    `for_admin` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_user` CHAR(27) NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    UNIQUE (`name`, `belongs_to_user`),"), jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"), jen.Lit(");")),
				jen.Lit("\n"),
			)), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.07), jen.ID("Description").MapAssign().Lit("create webhooks table"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS webhooks ("), jen.Lit("    `id` CHAR(27) NOT NULL,"), jen.Lit("    `name` VARCHAR(128) NOT NULL,"), jen.Lit("    `content_type` VARCHAR(64) NOT NULL,"), jen.Lit("    `url` LONGTEXT NOT NULL,"), jen.Lit("    `method` VARCHAR(8) NOT NULL,"), jen.Lit("    `events` VARCHAR(256) NOT NULL,"), jen.Lit("    `data_types` VARCHAR(256) NOT NULL,"), jen.Lit("    `topics` VARCHAR(256) NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_account` CHAR(27) NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE"), jen.Lit(");")),
				jen.Lit("\n"),
			)), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.08), jen.ID("Description").MapAssign().Lit("create notifications table"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS notifications ("), jen.Lit("    `id` CHAR(27) NOT NULL,"), jen.Lit("    `title` LONGTEXT NOT NULL,"), jen.Lit("    `description` LONGTEXT NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `seen_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_account` CHAR(27) NOT NULL,"), jen.Lit("    `belongs_to_user` CHAR(27) NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE,"), jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"), jen.Lit(");")),
				jen.Lit("\n"),
			)), jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.09), jen.ID("Description").MapAssign().Lit("create items table"), jen.ID("Script").MapAssign().Qual("strings", "Join").Call(
				jen.Index().String().Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS items ("), jen.Lit("    `id` CHAR(27) NOT NULL,"), jen.Lit("    `name` LONGTEXT NOT NULL,"), jen.Lit("    `details` LONGTEXT NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_account` CHAR(27) NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE"), jen.Lit(");")),
				jen.Lit("\n"),
			))),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Newline(),
		jen.Comment("migrate a postgres database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("migrationFunc").Params().Body(
			jen.ID("driver").Op(":=").ID("darwin").Dot("NewGenericDriver").Call(
				jen.ID("q").Dot("db"),
				jen.ID("darwin").Dot("MySQLDialect").Values(),
			),
			jen.If(jen.ID("err").Op(":=").ID("darwin").Dot("New").Call(
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
