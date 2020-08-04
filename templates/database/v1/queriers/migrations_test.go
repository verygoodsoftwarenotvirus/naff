package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_migrationsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := migrationsDotGo(proj, dbvendor)

		expected := `
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS users (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"salt" BYTEA NOT NULL,
				"password_last_changed_on" integer,
				"requires_password_change" boolean NOT NULL DEFAULT 'false',
				"two_factor_secret" TEXT NOT NULL,
				"two_factor_secret_verified_on" BIGINT DEFAULT NULL,
				"is_admin" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				UNIQUE ("username")
			);` + "`" + `,
		},
		{
			Version:     2,
			Description: "create sessions table for session manager",
			Script: ` + "`" + `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BYTEA NOT NULL,
				expiry TIMESTAMPTZ NOT NULL
			);

			CREATE INDEX sessions_expiry_idx ON sessions (expiry);
		` + "`" + `,
		},
		{
			Version:     3,
			Description: "create oauth2_clients table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     4,
			Description: "create webhooks table",
			Script: ` + "`" + `
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
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     5,
			Description: "create items table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS items (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"details" CHARACTER VARYING NOT NULL DEFAULT '',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
	}
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a postgres database.
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (p *Postgres) Migrate(ctx context.Context) error {
	p.logger.Info("migrating db")
	if !p.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	p.migrateOnce.Do(buildMigrationFunc(p.db))

	return nil
}
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := migrationsDotGo(proj, dbvendor)

		expected := `
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS users (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"salt" TINYBLOB NOT NULL,
				"password_last_changed_on" INTEGER,
				"requires_password_change" BOOLEAN NOT NULL DEFAULT 'false',
				"two_factor_secret" TEXT NOT NULL,
				"two_factor_secret_verified_on" INTEGER DEFAULT NULL,
				"is_admin" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				CONSTRAINT username_unique UNIQUE (username)
			);` + "`" + `,
		},
		{
			Version:     2,
			Description: "create sessions table for session manager",
			Script: ` + "`" + `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BLOB NOT NULL,
				expiry REAL NOT NULL
			);

			CREATE INDEX sessions_expiry_idx ON sessions(expiry);
			` + "`" + `,
		},
		{
			Version:     3,
			Description: "create oauth2_clients table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     4,
			Description: "create webhooks table",
			Script: ` + "`" + `
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
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     5,
			Description: "create items table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS items (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" CHARACTER VARYING NOT NULL,
				"details" CHARACTER VARYING NOT NULL DEFAULT '',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
	}
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a sqlite database.
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.SqliteDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (s *Sqlite) Migrate(ctx context.Context) error {
	s.logger.Info("migrating db")
	if !s.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	s.migrateOnce.Do(buildMigrationFunc(s.db))

	return nil
}
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := migrationsDotGo(proj, dbvendor)

		expected := `
package mariadb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GuiaBolso/darwin"
	"strings"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS users (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `username` + "`" + ` VARCHAR(150) NOT NULL,",
				"    ` + "`" + `hashed_password` + "`" + ` VARCHAR(100) NOT NULL,",
				"    ` + "`" + `salt` + "`" + ` BINARY(16) NOT NULL,",
				"    ` + "`" + `requires_password_change` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `password_last_changed_on` + "`" + ` INTEGER UNSIGNED,",
				"    ` + "`" + `two_factor_secret` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `two_factor_secret_verified_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `is_admin` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    UNIQUE (` + "`" + `username` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     2,
			Description: "create users table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS users_creation_trigger BEFORE INSERT ON users FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     3,
			Description: "create sessions table for session manager",
			Script: strings.Join([]string{
				"CREATE TABLE sessions (",
				"` + "`" + `token` + "`" + ` CHAR(43) PRIMARY KEY,",
				"` + "`" + `data` + "`" + ` BLOB NOT NULL,",
				"` + "`" + `expiry` + "`" + ` TIMESTAMP(6) NOT NULL",
				");",
			}, "\n"),
		},
		{
			Version:     4,
			Description: "create sessions table for session manager",
			Script:      "CREATE INDEX sessions_expiry_idx ON sessions (expiry);",
		},
		{
			Version:     5,
			Description: "create oauth2_clients table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS oauth2_clients (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` VARCHAR(128) DEFAULT '',",
				"    ` + "`" + `client_id` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `client_secret` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `redirect_uri` + "`" + ` VARCHAR(4096) DEFAULT '',",
				"    ` + "`" + `scopes` + "`" + ` VARCHAR(4096) NOT NULL,",
				"    ` + "`" + `implicit_allowed` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY(` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     6,
			Description: "create oauth2_clients table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS oauth2_clients_creation_trigger BEFORE INSERT ON oauth2_clients FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     7,
			Description: "create webhooks table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS webhooks (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` VARCHAR(128) NOT NULL,",
				"    ` + "`" + `content_type` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `url` + "`" + ` VARCHAR(4096) NOT NULL,",
				"    ` + "`" + `method` + "`" + ` VARCHAR(32) NOT NULL,",
				"    ` + "`" + `events` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `data_types` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `topics` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY (` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     8,
			Description: "create webhooks table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     9,
			Description: "create items table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS items (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` LONGTEXT NOT NULL,",
				"    ` + "`" + `details` + "`" + ` LONGTEXT NOT NULL DEFAULT '',",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY (` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     10,
			Description: "create items table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS items_creation_trigger BEFORE INSERT ON items FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
	}
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a maria DB database.
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (m *MariaDB) Migrate(ctx context.Context) error {
	m.logger.Info("migrating db")
	if !m.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	m.migrateOnce.Do(buildMigrationFunc(m.db))

	return nil
}
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMigrationVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildMigrationVarDeclarations(proj, dbvendor)

		expected := `
package example

import (
	darwin "github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS users (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"salt" BYTEA NOT NULL,
				"password_last_changed_on" integer,
				"requires_password_change" boolean NOT NULL DEFAULT 'false',
				"two_factor_secret" TEXT NOT NULL,
				"two_factor_secret_verified_on" BIGINT DEFAULT NULL,
				"is_admin" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				UNIQUE ("username")
			);` + "`" + `,
		},
		{
			Version:     2,
			Description: "create sessions table for session manager",
			Script: ` + "`" + `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BYTEA NOT NULL,
				expiry TIMESTAMPTZ NOT NULL
			);

			CREATE INDEX sessions_expiry_idx ON sessions (expiry);
		` + "`" + `,
		},
		{
			Version:     3,
			Description: "create oauth2_clients table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     4,
			Description: "create webhooks table",
			Script: ` + "`" + `
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
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     5,
			Description: "create items table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS items (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"details" CHARACTER VARYING NOT NULL DEFAULT '',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildMigrationVarDeclarations(proj, dbvendor)

		expected := `
package example

import (
	darwin "github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS users (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"salt" TINYBLOB NOT NULL,
				"password_last_changed_on" INTEGER,
				"requires_password_change" BOOLEAN NOT NULL DEFAULT 'false',
				"two_factor_secret" TEXT NOT NULL,
				"two_factor_secret_verified_on" INTEGER DEFAULT NULL,
				"is_admin" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				CONSTRAINT username_unique UNIQUE (username)
			);` + "`" + `,
		},
		{
			Version:     2,
			Description: "create sessions table for session manager",
			Script: ` + "`" + `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BLOB NOT NULL,
				expiry REAL NOT NULL
			);

			CREATE INDEX sessions_expiry_idx ON sessions(expiry);
			` + "`" + `,
		},
		{
			Version:     3,
			Description: "create oauth2_clients table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     4,
			Description: "create webhooks table",
			Script: ` + "`" + `
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
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     5,
			Description: "create items table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS items (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" CHARACTER VARYING NOT NULL,
				"details" CHARACTER VARYING NOT NULL DEFAULT '',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildMigrationVarDeclarations(proj, dbvendor)

		expected := `
package example

import (
	darwin "github.com/GuiaBolso/darwin"
	"strings"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS users (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `username` + "`" + ` VARCHAR(150) NOT NULL,",
				"    ` + "`" + `hashed_password` + "`" + ` VARCHAR(100) NOT NULL,",
				"    ` + "`" + `salt` + "`" + ` BINARY(16) NOT NULL,",
				"    ` + "`" + `requires_password_change` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `password_last_changed_on` + "`" + ` INTEGER UNSIGNED,",
				"    ` + "`" + `two_factor_secret` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `two_factor_secret_verified_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `is_admin` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    UNIQUE (` + "`" + `username` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     2,
			Description: "create users table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS users_creation_trigger BEFORE INSERT ON users FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     3,
			Description: "create sessions table for session manager",
			Script: strings.Join([]string{
				"CREATE TABLE sessions (",
				"` + "`" + `token` + "`" + ` CHAR(43) PRIMARY KEY,",
				"` + "`" + `data` + "`" + ` BLOB NOT NULL,",
				"` + "`" + `expiry` + "`" + ` TIMESTAMP(6) NOT NULL",
				");",
			}, "\n"),
		},
		{
			Version:     4,
			Description: "create sessions table for session manager",
			Script:      "CREATE INDEX sessions_expiry_idx ON sessions (expiry);",
		},
		{
			Version:     5,
			Description: "create oauth2_clients table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS oauth2_clients (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` VARCHAR(128) DEFAULT '',",
				"    ` + "`" + `client_id` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `client_secret` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `redirect_uri` + "`" + ` VARCHAR(4096) DEFAULT '',",
				"    ` + "`" + `scopes` + "`" + ` VARCHAR(4096) NOT NULL,",
				"    ` + "`" + `implicit_allowed` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY(` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     6,
			Description: "create oauth2_clients table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS oauth2_clients_creation_trigger BEFORE INSERT ON oauth2_clients FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     7,
			Description: "create webhooks table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS webhooks (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` VARCHAR(128) NOT NULL,",
				"    ` + "`" + `content_type` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `url` + "`" + ` VARCHAR(4096) NOT NULL,",
				"    ` + "`" + `method` + "`" + ` VARCHAR(32) NOT NULL,",
				"    ` + "`" + `events` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `data_types` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `topics` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY (` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     8,
			Description: "create webhooks table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     9,
			Description: "create items table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS items (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` LONGTEXT NOT NULL,",
				"    ` + "`" + `details` + "`" + ` LONGTEXT NOT NULL DEFAULT '',",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY (` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     10,
			Description: "create items table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS items_creation_trigger BEFORE INSERT ON items FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildMigrationFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildMigrationFuncDecl(dbvendor)

		expected := `
package example

import (
	"database/sql"
	darwin "github.com/GuiaBolso/darwin"
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a postgres database.
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildMigrationFuncDecl(dbvendor)

		expected := `
package example

import (
	"database/sql"
	darwin "github.com/GuiaBolso/darwin"
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a sqlite database.
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.SqliteDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildMigrationFuncDecl(dbvendor)

		expected := `
package example

import (
	"database/sql"
	darwin "github.com/GuiaBolso/darwin"
)

// buildMigrationFunc returns a sync.Once compatible function closure that will
// migrate a maria DB database.
func buildMigrationFunc(db *sql.DB) func() {
	return func() {
		driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})
		if err := darwin.New(driver, migrations, nil).Migrate(); err != nil {
			panic(err)
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMigrate(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildMigrate(dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
)

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (p *Postgres) Migrate(ctx context.Context) error {
	p.logger.Info("migrating db")
	if !p.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	p.migrateOnce.Do(buildMigrationFunc(p.db))

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildMigrate(dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
)

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (s *Sqlite) Migrate(ctx context.Context) error {
	s.logger.Info("migrating db")
	if !s.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	s.migrateOnce.Do(buildMigrationFunc(s.db))

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildMigrate(dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
)

// Migrate migrates the database. It does so by invoking the migrateOnce function via sync.Once, so it should be
// safe (as in idempotent, though not necessarily recommended) to call this function multiple times.
func (m *MariaDB) Migrate(ctx context.Context) error {
	m.logger.Info("migrating db")
	if !m.IsReady(ctx) {
		return errors.New("db is not ready yet")
	}

	m.migrateOnce.Do(buildMigrationFunc(m.db))

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_typeToPostgresType(T *testing.T) {
	T.Parallel()

	allTypes := []string{
		"string",
		"*string",
		"bool",
		"*bool",
		"int",
		"*int",
		"int8",
		"*int8",
		"int16",
		"*int16",
		"int32",
		"*int32",
		"int64",
		"*int64",
		"uint",
		"*uint",
		"uint8",
		"*uint8",
		"uint16",
		"*uint16",
		"uint32",
		"*uint32",
		"uint64",
		"*uint64",
		"float32",
		"*float32",
		"float64",
		"*float64",
	}

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		for _, typ := range allTypes {
			assert.NotEmpty(t, typeToPostgresType(typ))
		}
	})

	T.Run("panics on unknown type", func(t *testing.T) {
		t.Parallel()

		typ := "fart"

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		typeToPostgresType(typ)
	})
}

func Test_typeToSqliteType(T *testing.T) {
	T.Parallel()

	allTypes := []string{
		"string",
		"*string",
		"bool",
		"*bool",
		"int",
		"*int",
		"int8",
		"*int8",
		"int16",
		"*int16",
		"int32",
		"*int32",
		"int64",
		"*int64",
		"uint",
		"*uint",
		"uint8",
		"*uint8",
		"uint16",
		"*uint16",
		"uint32",
		"*uint32",
		"uint64",
		"*uint64",
		"float32",
		"*float32",
		"float64",
		"*float64",
	}

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		for _, typ := range allTypes {
			assert.NotEmpty(t, typeToSqliteType(typ))
		}
	})

	T.Run("panics on unknown type", func(t *testing.T) {
		t.Parallel()

		typ := "fart"

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		typeToSqliteType(typ)
	})
}

func Test_typeToMariaDBType(T *testing.T) {
	T.Parallel()

	allTypes := []string{
		"string",
		"*string",
		"bool",
		"*bool",
		"int",
		"*int",
		"int8",
		"*int8",
		"int16",
		"*int16",
		"int32",
		"*int32",
		"int64",
		"*int64",
		"uint",
		"*uint",
		"uint8",
		"*uint8",
		"uint16",
		"*uint16",
		"uint32",
		"*uint32",
		"uint64",
		"*uint64",
		"float32",
		"*float32",
		"float64",
		"*float64",
	}

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		for _, typ := range allTypes {
			assert.NotEmpty(t, typeToMariaDBType(typ))
		}
	})

	T.Run("panics on unknown type", func(t *testing.T) {
		t.Parallel()

		typ := "fart"

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		typeToMariaDBType(typ)
	})
}

func Test_makePostgresMigrations(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := makePostgresMigrations(proj)

		assert.NotEmpty(t, x)
	})
}

func Test_makeMariaDBMigrations(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := makeMariaDBMigrations(proj)

		assert.NotEmpty(t, x)
	})
}

func Test_makeSqliteMigrations(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := makeSqliteMigrations(proj)

		assert.NotEmpty(t, x)
	})
}

func Test_makeMigrations(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := makeMigrations(proj, dbvendor)

		expected := `
package example

import (
	darwin "github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS users (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"salt" BYTEA NOT NULL,
				"password_last_changed_on" integer,
				"requires_password_change" boolean NOT NULL DEFAULT 'false',
				"two_factor_secret" TEXT NOT NULL,
				"two_factor_secret_verified_on" BIGINT DEFAULT NULL,
				"is_admin" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				UNIQUE ("username")
			);` + "`" + `,
		},
		{
			Version:     2,
			Description: "create sessions table for session manager",
			Script: ` + "`" + `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BYTEA NOT NULL,
				expiry TIMESTAMPTZ NOT NULL
			);

			CREATE INDEX sessions_expiry_idx ON sessions (expiry);
		` + "`" + `,
		},
		{
			Version:     3,
			Description: "create oauth2_clients table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" boolean NOT NULL DEFAULT 'false',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     4,
			Description: "create webhooks table",
			Script: ` + "`" + `
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
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     5,
			Description: "create items table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS items (
				"id" BIGSERIAL NOT NULL PRIMARY KEY,
				"name" CHARACTER VARYING NOT NULL,
				"details" CHARACTER VARYING NOT NULL DEFAULT '',
				"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
				"last_updated_on" BIGINT DEFAULT NULL,
				"archived_on" BIGINT DEFAULT NULL,
				"belongs_to_user" BIGINT NOT NULL,
				FOREIGN KEY("belongs_to_user") REFERENCES users(id)
			);` + "`" + `,
		},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := makeMigrations(proj, dbvendor)

		expected := `
package example

import (
	darwin "github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS users (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"username" TEXT NOT NULL,
				"hashed_password" TEXT NOT NULL,
				"salt" TINYBLOB NOT NULL,
				"password_last_changed_on" INTEGER,
				"requires_password_change" BOOLEAN NOT NULL DEFAULT 'false',
				"two_factor_secret" TEXT NOT NULL,
				"two_factor_secret_verified_on" INTEGER DEFAULT NULL,
				"is_admin" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				CONSTRAINT username_unique UNIQUE (username)
			);` + "`" + `,
		},
		{
			Version:     2,
			Description: "create sessions table for session manager",
			Script: ` + "`" + `
			CREATE TABLE sessions (
				token TEXT PRIMARY KEY,
				data BLOB NOT NULL,
				expiry REAL NOT NULL
			);

			CREATE INDEX sessions_expiry_idx ON sessions(expiry);
			` + "`" + `,
		},
		{
			Version:     3,
			Description: "create oauth2_clients table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS oauth2_clients (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" TEXT DEFAULT '',
				"client_id" TEXT NOT NULL,
				"client_secret" TEXT NOT NULL,
				"redirect_uri" TEXT DEFAULT '',
				"scopes" TEXT NOT NULL,
				"implicit_allowed" BOOLEAN NOT NULL DEFAULT 'false',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     4,
			Description: "create webhooks table",
			Script: ` + "`" + `
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
				"last_updated_on" INTEGER,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
		{
			Version:     5,
			Description: "create items table",
			Script: ` + "`" + `
			CREATE TABLE IF NOT EXISTS items (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"name" CHARACTER VARYING NOT NULL,
				"details" CHARACTER VARYING NOT NULL DEFAULT '',
				"created_on" INTEGER NOT NULL DEFAULT (strftime('%s','now')),
				"last_updated_on" INTEGER DEFAULT NULL,
				"archived_on" INTEGER DEFAULT NULL,
				"belongs_to_user" INTEGER NOT NULL,
				FOREIGN KEY(belongs_to_user) REFERENCES users(id)
			);` + "`" + `,
		},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := makeMigrations(proj, dbvendor)

		expected := `
package example

import (
	darwin "github.com/GuiaBolso/darwin"
	"strings"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "create users table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS users (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `username` + "`" + ` VARCHAR(150) NOT NULL,",
				"    ` + "`" + `hashed_password` + "`" + ` VARCHAR(100) NOT NULL,",
				"    ` + "`" + `salt` + "`" + ` BINARY(16) NOT NULL,",
				"    ` + "`" + `requires_password_change` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `password_last_changed_on` + "`" + ` INTEGER UNSIGNED,",
				"    ` + "`" + `two_factor_secret` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `two_factor_secret_verified_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `is_admin` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    UNIQUE (` + "`" + `username` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     2,
			Description: "create users table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS users_creation_trigger BEFORE INSERT ON users FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     3,
			Description: "create sessions table for session manager",
			Script: strings.Join([]string{
				"CREATE TABLE sessions (",
				"` + "`" + `token` + "`" + ` CHAR(43) PRIMARY KEY,",
				"` + "`" + `data` + "`" + ` BLOB NOT NULL,",
				"` + "`" + `expiry` + "`" + ` TIMESTAMP(6) NOT NULL",
				");",
			}, "\n"),
		},
		{
			Version:     4,
			Description: "create sessions table for session manager",
			Script:      "CREATE INDEX sessions_expiry_idx ON sessions (expiry);",
		},
		{
			Version:     5,
			Description: "create oauth2_clients table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS oauth2_clients (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` VARCHAR(128) DEFAULT '',",
				"    ` + "`" + `client_id` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `client_secret` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `redirect_uri` + "`" + ` VARCHAR(4096) DEFAULT '',",
				"    ` + "`" + `scopes` + "`" + ` VARCHAR(4096) NOT NULL,",
				"    ` + "`" + `implicit_allowed` + "`" + ` BOOLEAN NOT NULL DEFAULT false,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY(` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     6,
			Description: "create oauth2_clients table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS oauth2_clients_creation_trigger BEFORE INSERT ON oauth2_clients FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     7,
			Description: "create webhooks table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS webhooks (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` VARCHAR(128) NOT NULL,",
				"    ` + "`" + `content_type` + "`" + ` VARCHAR(64) NOT NULL,",
				"    ` + "`" + `url` + "`" + ` VARCHAR(4096) NOT NULL,",
				"    ` + "`" + `method` + "`" + ` VARCHAR(32) NOT NULL,",
				"    ` + "`" + `events` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `data_types` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `topics` + "`" + ` VARCHAR(256) NOT NULL,",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY (` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     8,
			Description: "create webhooks table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS webhooks_creation_trigger BEFORE INSERT ON webhooks FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
		{
			Version:     9,
			Description: "create items table",
			Script: strings.Join([]string{
				"CREATE TABLE IF NOT EXISTS items (",
				"    ` + "`" + `id` + "`" + ` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,",
				"    ` + "`" + `name` + "`" + ` LONGTEXT NOT NULL,",
				"    ` + "`" + `details` + "`" + ` LONGTEXT NOT NULL DEFAULT '',",
				"    ` + "`" + `created_on` + "`" + ` BIGINT UNSIGNED,",
				"    ` + "`" + `last_updated_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `archived_on` + "`" + ` BIGINT UNSIGNED DEFAULT NULL,",
				"    ` + "`" + `belongs_to_user` + "`" + ` BIGINT UNSIGNED NOT NULL,",
				"    PRIMARY KEY (` + "`" + `id` + "`" + `),",
				"    FOREIGN KEY (` + "`" + `belongs_to_user` + "`" + `) REFERENCES users(` + "`" + `id` + "`" + `)",
				");",
			}, "\n"),
		},
		{
			Version:     10,
			Description: "create items table creation trigger",
			Script: strings.Join([]string{
				"CREATE TRIGGER IF NOT EXISTS items_creation_trigger BEFORE INSERT ON items FOR EACH ROW",
				"BEGIN",
				"  IF (new.created_on is null)",
				"  THEN",
				"    SET new.created_on = UNIX_TIMESTAMP();",
				"  END IF;",
				"END;",
			}, "\n"),
		},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
