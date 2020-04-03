package queriers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/stretchr/testify/assert"
)

func Test_makePostgresMigrations(T *testing.T) {
	T.Parallel()

	standardFields := []models.DataField{
		{
			Name: wordsmith.FromSingularPascalCase("FieldOne"),
			Type: "string",
		},
		{
			Name: wordsmith.FromSingularPascalCase("FieldTwo"),
			Type: "string",
		},
		{
			Name: wordsmith.FromSingularPascalCase("FieldThree"),
			Type: "string",
		},
	}

	T.Run("belongs to user", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:          wordsmith.FromSingularPascalCase("ThingOne"),
					BelongsToUser: true,
					Fields:        standardFields,
				},
			},
		}

		expected := []migration{
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
			{
				description: "create thing ones table",
				script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS thing_ones (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"field_one" CHARACTER VARYING NOT NULL,
				"field_two" CHARACTER VARYING NOT NULL,
				"field_three" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_user") REFERENCES "users"("id")
			);`),
			},
		}
		actual := makePostgresMigrations(p)

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to object", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:            wordsmith.FromSingularPascalCase("ThingTwo"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("Owner"),
					Fields:          standardFields,
				},
			},
		}

		expected := []migration{
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
			{
				description: "create thing twos table",
				script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS thing_twos (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"field_one" CHARACTER VARYING NOT NULL,
				"field_two" CHARACTER VARYING NOT NULL,
				"field_three" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_owner" BIGINT NOT NULL,
				FOREIGN KEY ("belongs_to_owner") REFERENCES "owners"("id")
			);`),
			},
		}
		actual := makePostgresMigrations(p)

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to nobody", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:            wordsmith.FromSingularPascalCase("ThingThree"),
					BelongsToNobody: true,
					Fields:          standardFields,
				},
			},
		}

		expected := []migration{
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
			{
				description: "create thing threes table",
				script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS thing_threes (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"field_one" CHARACTER VARYING NOT NULL,
				"field_two" CHARACTER VARYING NOT NULL,
				"field_three" CHARACTER VARYING NOT NULL,
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
			);`),
			},
		}
		actual := makePostgresMigrations(p)

		assert.Equal(t, expected, actual)
	})

}

func Test_makeMariaDBMigrations(T *testing.T) {
	T.Parallel()

	standardFields := []models.DataField{
		{
			Name: wordsmith.FromSingularPascalCase("FieldOne"),
			Type: "string",
		},
		{
			Name: wordsmith.FromSingularPascalCase("FieldTwo"),
			Type: "string",
		},
		{
			Name: wordsmith.FromSingularPascalCase("FieldThree"),
			Type: "string",
		},
	}

	T.Run("belongs to user", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:          wordsmith.FromSingularPascalCase("ThingOne"),
					BelongsToUser: true,
					Fields:        standardFields,
				},
			},
		}

		expected := []migration{
			{
				description: "create users table",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
			{
				description: "create thing ones table",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TABLE IF NOT EXISTS thing_ones ("),
					jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
					jen.Lit("    `field_one` LONGTEXT NOT NULL,"),
					jen.Lit("    `field_two` LONGTEXT NOT NULL,"),
					jen.Lit("    `field_three` LONGTEXT NOT NULL,"),
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
				description: "create thing ones table creation trigger",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS thing_ones_creation_trigger BEFORE INSERT ON thing_ones FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
		}
		actual := makeMariaDBMigrations(p)

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to object", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:            wordsmith.FromSingularPascalCase("ThingTwo"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("Owner"),
					Fields:          standardFields,
				},
			},
		}

		expected := []migration{
			{
				description: "create users table",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
			{
				description: "create thing twos table",

				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TABLE IF NOT EXISTS thing_twos ("),
					jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
					jen.Lit("    `field_one` LONGTEXT NOT NULL,"),
					jen.Lit("    `field_two` LONGTEXT NOT NULL,"),
					jen.Lit("    `field_three` LONGTEXT NOT NULL,"),
					jen.Lit("    `created_on` BIGINT UNSIGNED,"),
					jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `belongs_to_owner` BIGINT UNSIGNED NOT NULL,"),
					jen.Lit("    PRIMARY KEY (`id`),"),
					jen.Lit("    FOREIGN KEY (`belongs_to_owner`) REFERENCES owners(`id`)"),
					jen.Lit(");"),
				), jen.Lit("\n")),
			},
			{
				description: "create thing twos table creation trigger",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS thing_twos_creation_trigger BEFORE INSERT ON thing_twos FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
		}
		actual := makeMariaDBMigrations(p)

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to nobody", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:            wordsmith.FromSingularPascalCase("ThingThree"),
					BelongsToNobody: true,
					Fields:          standardFields,
				},
			},
		}

		expected := []migration{
			{
				description: "create users table",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
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
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
			{
				description: "create thing threes table",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TABLE IF NOT EXISTS thing_threes ("),
					jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
					jen.Lit("    `field_one` LONGTEXT NOT NULL,"),
					jen.Lit("    `field_two` LONGTEXT NOT NULL,"),
					jen.Lit("    `field_three` LONGTEXT NOT NULL,"),
					jen.Lit("    `created_on` BIGINT UNSIGNED,"),
					jen.Lit("    `updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
					jen.Lit("    PRIMARY KEY (`id`),"),
					jen.Lit(");"),
				), jen.Lit("\n")),
			},
			{
				description: "create thing threes table creation trigger",
				script: jen.Qual("strings", "Join").Call(jen.Index().String().Valuesln(
					jen.Lit("CREATE TRIGGER IF NOT EXISTS thing_threes_creation_trigger BEFORE INSERT ON thing_threes FOR EACH ROW"),
					jen.Lit("BEGIN"),
					jen.Lit("  IF (new.created_on is null)"),
					jen.Lit("  THEN"),
					jen.Lit("    SET new.created_on = UNIX_TIMESTAMP();"),
					jen.Lit("  END IF;"),
					jen.Lit("END;"),
				), jen.Lit("\n")),
			},
		}
		actual := makeMariaDBMigrations(p)

		assert.Equal(t, expected, actual)
	})

}

func Test_makeSqliteMigrations(T *testing.T) {
	T.Parallel()

	standardFields := []models.DataField{
		{
			Name: wordsmith.FromSingularPascalCase("FieldOne"),
			Type: "string",
		},
		{
			Name: wordsmith.FromSingularPascalCase("FieldTwo"),
			Type: "string",
		},
		{
			Name: wordsmith.FromSingularPascalCase("FieldThree"),
			Type: "string",
		},
	}

	T.Run("belongs to user", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:          wordsmith.FromSingularPascalCase("ThingOne"),
					BelongsToUser: true,
					Fields:        standardFields,
				},
			},
		}

		expected := []migration{
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
			{
				description: "create thing ones table",
				script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS thing_ones (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"field_one" CHARACTER VARYING NOT NULL,
				"field_two" CHARACTER VARYING NOT NULL,
				"field_three" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);`),
			},
		}
		actual := makeSqliteMigrations(p)

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to object", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:            wordsmith.FromSingularPascalCase("ThingTwo"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("Owner"),
					Fields:          standardFields,
				},
			},
		}

		expected := []migration{
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
			{
				description: "create thing twos table",
				script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS thing_twos (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"field_one" CHARACTER VARYING NOT NULL,
				"field_two" CHARACTER VARYING NOT NULL,
				"field_three" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_owner" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_owner) REFERENCES owners(id)
			);`),
			},
		}
		actual := makeSqliteMigrations(p)

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to nobody", func(t *testing.T) {
		p := &models.Project{
			OutputPath:    "",
			EnableNewsman: false,
			Name:          nil,
			DataTypes: []models.DataType{
				{
					Name:            wordsmith.FromSingularPascalCase("ThingThree"),
					BelongsToNobody: true,
					Fields:          standardFields,
				},
			},
		}

		expected := []migration{
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
			{
				description: "create thing threes table",
				script: jen.Lit(`
			CREATE TABLE IF NOT EXISTS thing_threes (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"field_one" CHARACTER VARYING NOT NULL,
				"field_two" CHARACTER VARYING NOT NULL,
				"field_three" CHARACTER VARYING NOT NULL,
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
			);`),
			},
		}
		actual := makeSqliteMigrations(p)

		assert.Equal(t, expected, actual)
	})

}
