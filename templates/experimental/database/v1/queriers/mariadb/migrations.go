package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func migrationsDotGo() *jen.File {
	ret := jen.NewFile("mariadb")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("migrations").Op("=").Index().ID("darwin").Dot(
			"Migration",
		).Valuesln(jen.Valuesln(jen.ID("Version").Op(":").Lit(1), jen.ID("Description").Op(":").Lit("create users table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS users ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `username` VARCHAR(150) NOT NULL,"), jen.Lit("    `hashed_password` VARCHAR(100) NOT NULL,"), jen.Lit("    `password_last_changed_on` INTEGER UNSIGNED,"), jen.Lit("    `two_factor_secret` VARCHAR(256) NOT NULL,"), jen.Lit("    `is_admin` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    UNIQUE (`username`)"), jen.Lit(");")), jen.Lit("\n"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(2), jen.ID("Description").Op(":").Lit("create users table creation trigger"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TRIGGER IF NOT EXISTS users_creation_trigger BEFORE INSERT ON users FOR EACH ROW"), jen.Lit("BEGIN"), jen.Lit("  IF (new.created_on is null)"), jen.Lit("  THEN"), jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"), jen.Lit("  END IF;"), jen.Lit("END;")), jen.Lit("\n"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(3), jen.ID("Description").Op(":").Lit("Add oauth2_clients table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS oauth2_clients ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `name` VARCHAR(128) DEFAULT '',"), jen.Lit("    `client_id` VARCHAR(64) NOT NULL,"), jen.Lit("    `client_secret` VARCHAR(64) NOT NULL,"), jen.Lit("    `redirect_uri` VARCHAR(4096) DEFAULT '',"), jen.Lit("    `scopes` VARCHAR(4096) NOT NULL,"), jen.Lit("    `implicit_allowed` BOOLEAN NOT NULL DEFAULT false,"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY(`belongs_to`) REFERENCES users(`id`)"), jen.Lit(");")), jen.Lit("\n"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(4), jen.ID("Description").Op(":").Lit("create oauth2_clients table creation trigger"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TRIGGER IF NOT EXISTS oauth2_clients_creation_trigger BEFORE INSERT ON oauth2_clients FOR EACH ROW"), jen.Lit("BEGIN"), jen.Lit("  IF (new.created_on is null)"), jen.Lit("  THEN"), jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"), jen.Lit("  END IF;"), jen.Lit("END;")), jen.Lit("\n"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(5), jen.ID("Description").Op(":").Lit("create webhooks table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS webhooks ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `name` VARCHAR(128) NOT NULL,"), jen.Lit("    `content_type` VARCHAR(64) NOT NULL,"), jen.Lit("    `url` VARCHAR(4096) NOT NULL,"), jen.Lit("    `method` VARCHAR(32) NOT NULL,"), jen.Lit("    `events` VARCHAR(256) NOT NULL,"), jen.Lit("    `data_types` VARCHAR(256) NOT NULL,"), jen.Lit("    `topics` VARCHAR(256) NOT NULL,"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to`) REFERENCES users(`id`)"), jen.Lit(");")), jen.Lit("\n"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(6), jen.ID("Description").Op(":").Lit("create webhooks table creation trigger"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW"), jen.Lit("BEGIN"), jen.Lit("  IF (new.created_on is null)"), jen.Lit("  THEN"), jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"), jen.Lit("  END IF;"), jen.Lit("END;")), jen.Lit("\n"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(7), jen.ID("Description").Op(":").Lit("create items table"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TABLE IF NOT EXISTS items ("), jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"), jen.Lit("    `name` VARCHAR(256) NOT NULL,"), jen.Lit("    `details` VARCHAR(4096) NOT NULL DEFAULT '',"), jen.Lit("    `created_on` BIGINT UNSIGNED,"), jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("  	 `archived_on` BIGINT UNSIGNED DEFAULT NULL,"), jen.Lit("    `belongs_to` BIGINT UNSIGNED NOT NULL,"), jen.Lit("    PRIMARY KEY (`id`),"), jen.Lit("    FOREIGN KEY (`belongs_to`) REFERENCES users(`id`)"), jen.Lit(");")), jen.Lit("\n"))), jen.Valuesln(jen.ID("Version").Op(":").Lit(8), jen.ID("Description").Op(":").Lit("create items table creation trigger"), jen.ID("Script").Op(":").Qual("strings", "Join").Call(jen.Index().ID("string").Valuesln(jen.Lit("CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW"), jen.Lit("BEGIN"), jen.Lit("  IF (new.created_on is null)"), jen.Lit("  THEN"), jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"), jen.Lit("  END IF;"), jen.Lit("END;")), jen.Lit("\n")))),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Line(),
		jen.Comment("migrate a MariaDB database"),
		jen.Line(),
		jen.Func().ID("buildMigrationFunc").Params(jen.ID("db").Op("*").Qual("database/sql", "DB")).Params(jen.Params()).Block(
			jen.Return().Func().Params().Block(
				jen.ID("driver").Op(":=").ID("darwin").Dot("NewGenericDriver").Call(jen.ID("db"), jen.ID("darwin").Dot("MySQLDialect").Values()),
				jen.If(jen.ID("err").Op(":=").ID("darwin").Dot("New").Call(jen.ID("driver"), jen.ID("migrations"), jen.ID("nil")).Dot("Migrate").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("panic").Call(jen.ID("err")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be"),
		jen.Line(),
		jen.Func().Comment("// safe (as in idempotent, though not recommended) to call this function multiple times.").Params(jen.ID("m").Op("*").ID("MariaDB")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Block(
			jen.ID("m").Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.If(jen.Op("!").ID("m").Dot("IsReady").Call(jen.ID("ctx"))).Block(
				jen.Return().ID("errors").Dot("New").Call(jen.Lit("db is not ready yet")),
			),
			jen.ID("m").Dot("migrateOnce").Dot("Do").Call(jen.ID("buildMigrationFunc").Call(jen.ID("m").Dot(
				"db",
			))),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)
	return ret
}
