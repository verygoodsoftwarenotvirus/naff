package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrationsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildCreationTriggerScript").Params(jen.ID("tableName").ID("string")).Params(jen.ID("string")).Body(
			jen.Return().Qual("strings", "Join").Call(
				jen.Index().ID("string").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS %s_creation_trigger BEFORE INSERT ON %s FOR EACH ROW"),
					jen.ID("tableName"),
					jen.ID("tableName"),
				), jen.Lit("BEGIN"), jen.Lit("  IF (new.created_on is null)"), jen.Lit("  THEN"), jen.Lit("    SET new.created_on = UNIX_TIMESTAMP(NOW());"), jen.Lit("  END IF;"), jen.Lit("END;")),
				jen.Lit("\n"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("migrations").Op("=").Index().ID("darwin").Dot("Migration").Valuesln(jen.Valuesln(jen.ID("Version").Op(":").Lit(0.00), jen.ID("Description").Op(":").Lit("create sessions table for session manager"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE sessions ("), jen.Lit("`token` CHAR(43) PRIMARY KEY,"), jen.Lit("`data` BLOB NOT NULL,"), jen.Lit("`expiry` TIMESTAMP(6) NOT NULL,"), jen.Lit("`created_on` BIGINT UNSIGNED"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.01), jen.ID("Description").Op(":").Lit("create sessions table for session manager"), jen.ID("Script").Op(":").Lit("CREATE INDEX sessions_expiry_idx ON sessions (expiry);")), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.02), jen.ID("Description").Op(":").Lit("create sessions table creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.Lit("sessions"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.03), jen.ID("Description").Op(":").Lit("create audit log table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS audit_log ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"), jen.Lit("    `event_type` VARCHAR(256) NOT NULL,"), jen.Lit("    `context` JSON NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    PRIMARY KEY (`id`)"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.04), jen.ID("Description").Op(":").Lit("create audit_log table creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.05), jen.ID("Description").Op(":").Lit("create users table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS users ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"), jen.Lit("    `username` VARCHAR(128) NOT NULL,"), jen.Lit("    `avatar_src` LONGTEXT NOT NULL DEFAULT '',"), jen.Lit("    `hashed_password` VARCHAR(100) NOT NULL,"), jen.Lit("    `requires_password_change` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `password_last_changed_on` INTEGER UNSIGNED,"), jen.Lit("    `two_factor_secret` VARCHAR(256) NOT NULL,"), jen.Lit("    `two_factor_secret_verified_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `service_roles` LONGTEXT NOT NULL DEFAULT 'service_user',"), jen.Lit("    `reputation` VARCHAR(64) NOT NULL DEFAULT 'unverified',"), jen.Lit("    `reputation_explanation` VARCHAR(1024) NOT NULL DEFAULT '',"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    UNIQUE (`username`)"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.06), jen.ID("Description").Op(":").Lit("create users table creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.ID("querybuilding").Dot("UsersTableName"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.07), jen.ID("Description").Op(":").Lit("create accounts table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS accounts ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"), jen.Lit("    `name` LONGTEXT NOT NULL,"), jen.Lit("	 `billing_status` TEXT NOT NULL DEFAULT 'unpaid',"), jen.Lit("	 `contact_email` TEXT NOT NULL DEFAULT '',"), jen.Lit("	 `contact_phone` TEXT NOT NULL DEFAULT '',"), jen.Lit("	 `payment_processor_customer_id` TEXT NOT NULL DEFAULT '',"), jen.Lit("    `subscription_plan_id` VARCHAR(128),"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.08), jen.ID("Description").Op(":").Lit("create accounts table creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.ID("querybuilding").Dot("AccountsTableName"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.09), jen.ID("Description").Op(":").Lit("create account user memberships table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS account_user_memberships ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `belongs_to_account` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    `default_account` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `account_roles` LONGTEXT NOT NULL DEFAULT 'account_user',"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE,"), jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE,"), jen.Lit("    UNIQUE (`belongs_to_account`, `belongs_to_user`)"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.10), jen.ID("Description").Op(":").Lit("create accounts membership creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.Lit("account_user_memberships"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.11), jen.ID("Description").Op(":").Lit("create API clients table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS api_clients ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"), jen.Lit("    `name` VARCHAR(128) DEFAULT '',"), jen.Lit("    `client_id` VARCHAR(64) NOT NULL,"), jen.Lit("    `secret_key` BINARY(128) NOT NULL,"), jen.Lit("    `account_roles` LONGTEXT NOT NULL DEFAULT 'account_member',"), jen.Lit("    `for_admin` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_user` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    UNIQUE (`name`, `belongs_to_user`),"), jen.Lit("    FOREIGN KEY (`belongs_to_user`) REFERENCES users(`id`) ON DELETE CASCADE"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.12), jen.ID("Description").Op(":").Lit("create api_clients table creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.ID("querybuilding").Dot("APIClientsTableName"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.13), jen.ID("Description").Op(":").Lit("create webhooks table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS webhooks ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"), jen.Lit("    `name` VARCHAR(128) NOT NULL,"), jen.Lit("    `content_type` VARCHAR(64) NOT NULL,"), jen.Lit("    `url` LONGTEXT NOT NULL,"), jen.Lit("    `method` VARCHAR(8) NOT NULL,"), jen.Lit("    `events` VARCHAR(256) NOT NULL,"), jen.Lit("    `data_types` VARCHAR(256) NOT NULL,"), jen.Lit("    `topics` VARCHAR(256) NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_account` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.14), jen.ID("Description").Op(":").Lit("create webhooks table creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.ID("querybuilding").Dot("WebhooksTableName"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.15), jen.ID("Description").Op(":").Lit("create items table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(
			jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS items ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"), jen.Lit("    `name` LONGTEXT NOT NULL,"), jen.Lit("    `details` LONGTEXT NOT NULL DEFAULT '',"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to_account` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE"), jen.Lit(");")),
			jen.Lit("\n"),
		)), jen.Valuesln(jen.ID("Version").Op(":").Lit(0.16), jen.ID("Description").Op(":").Lit("create items table creation trigger"), jen.ID("Script").Op(":").ID("buildCreationTriggerScript").Call(jen.ID("querybuilding").Dot("ItemsTableName")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Line(),
		jen.Func().Comment("migrate a maria DB database.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()).Body(
			jen.Return().Func().Params().Body(
				jen.ID("driver").Op(":=").ID("darwin").Dot("NewGenericDriver").Call(
					jen.ID("db"),
					jen.ID("darwin").Dot("MySQLDialect").Valuesln(),
				),
				jen.If(jen.ID("err").Op(":=").ID("darwin").Dot("New").Call(
					jen.ID("driver"),
					jen.ID("migrations"),
					jen.ID("nil"),
				).Dot("Migrate").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("panic").Call(jen.Qual("fmt", "Errorf").Call(
						jen.Lit("migrating database: %w"),
						jen.ID("err"),
					))),
			)),
		jen.Line(),
	)

	return code
}
