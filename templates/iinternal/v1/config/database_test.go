package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_databaseDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := databaseDotGo(proj)

		expected := `
package example

import (
	"context"
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	"errors"
	"fmt"
	v2 "github.com/alexedwards/scs/v2"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1/client"
)

const ()

// ProvideDatabaseConnection provides a database implementation dependent on the configuration.
func (cfg *ServerConfig) ProvideDatabaseConnection(logger v1.Logger) (*sql.DB, error) {
	switch cfg.Database.Provider {
	default:
		return nil, fmt.Errorf("invalid database type selected: %q", cfg.Database.Provider)
	}
}

// ProvideDatabaseClient provides a database implementation dependent on the configuration.
func (cfg *ServerConfig) ProvideDatabaseClient(ctx context.Context, logger v1.Logger, rawDB *sql.DB) (v11.DataManager, error) {
	if rawDB == nil {
		return nil, errors.New("nil DB connection provided")
	}

	debug := cfg.Database.Debug || cfg.Meta.Debug

	ocsql.RegisterAllViews()
	ocsql.RecordStats(rawDB, cfg.Metrics.DBMetricsCollectionInterval)

	var dbc v11.DataManager
	switch cfg.Database.Provider {
	default:
		return nil, fmt.Errorf("invalid database type selected: %q", cfg.Database.Provider)
	}

	return client.ProvideDatabaseClient(ctx, rawDB, dbc, debug, logger)
}

// ProvideSessionManager provides a session manager based on some settings.
// There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,
// but it obviously doesn't belong in the database package, or maybe it does
func ProvideSessionManager(authConf AuthSettings, dbConf DatabaseSettings, db *sql.DB) *v2.SessionManager {
	sessionManager := v2.New()

	switch dbConf.Provider {
	}

	sessionManager.Lifetime = authConf.CookieLifetime
	// elaborate further here later if you so choose

	return sessionManager
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDatabaseConstantDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildDatabaseConstantDeclarations(proj)

		expected := `
package example

import ()

const ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideDatabaseConnection(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseConnection(proj)

		expected := `
package example

import (
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

// ProvideDatabaseConnection provides a database implementation dependent on the configuration.
func (cfg *ServerConfig) ProvideDatabaseConnection(logger v1.Logger) (*sql.DB, error) {
	switch cfg.Database.Provider {
	default:
		return nil, fmt.Errorf("invalid database type selected: %q", cfg.Database.Provider)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseClient(proj)

		expected := `
package example

import (
	"context"
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	"errors"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1/client"
)

// ProvideDatabaseClient provides a database implementation dependent on the configuration.
func (cfg *ServerConfig) ProvideDatabaseClient(ctx context.Context, logger v1.Logger, rawDB *sql.DB) (v11.DataManager, error) {
	if rawDB == nil {
		return nil, errors.New("nil DB connection provided")
	}

	debug := cfg.Database.Debug || cfg.Meta.Debug

	ocsql.RegisterAllViews()
	ocsql.RecordStats(rawDB, cfg.Metrics.DBMetricsCollectionInterval)

	var dbc v11.DataManager
	switch cfg.Database.Provider {
	default:
		return nil, fmt.Errorf("invalid database type selected: %q", cfg.Database.Provider)
	}

	return client.ProvideDatabaseClient(ctx, rawDB, dbc, debug, logger)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideSessionManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildProvideSessionManager(proj)

		expected := `
package example

import (
	"database/sql"
	v2 "github.com/alexedwards/scs/v2"
)

// ProvideSessionManager provides a session manager based on some settings.
// There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,
// but it obviously doesn't belong in the database package, or maybe it does
func ProvideSessionManager(authConf AuthSettings, dbConf DatabaseSettings, db *sql.DB) *v2.SessionManager {
	sessionManager := v2.New()

	switch dbConf.Provider {
	}

	sessionManager.Lifetime = authConf.CookieLifetime
	// elaborate further here later if you so choose

	return sessionManager
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
