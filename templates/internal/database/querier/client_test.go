package querier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_clientDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := clientDotGo(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

var _ v1.DataManager = (*Client)(nil)

/*
	NOTE: the primary purpose of this client is to allow convenient
	wrapping of actual query execution.
*/

// Client is a wrapper around a database querier. Client is where all
// logging and trace propagation should happen, the querier is where
// the actual database querying is performed.
type Client struct {
	db      *sql.DB
	querier v1.DataManager
	debug   bool
	logger  v11.Logger
}

// Migrate is a simple wrapper around the core querier Migrate call.
func (c *Client) Migrate(ctx context.Context) error {
	ctx, span := tracing.StartSpan(ctx, "Migrate")
	defer span.End()

	return c.querier.Migrate(ctx)
}

// IsReady is a simple wrapper around the core querier IsReady call.
func (c *Client) IsReady(ctx context.Context) (ready bool) {
	ctx, span := tracing.StartSpan(ctx, "IsReady")
	defer span.End()

	return c.querier.IsReady(ctx)
}

// ProvideDatabaseClient provides a new DataManager client.
func ProvideDatabaseClient(
	ctx context.Context,
	db *sql.DB,
	querier v1.DataManager,
	debug bool,
	logger v11.Logger,
) (v1.DataManager, error) {
	c := &Client{
		db:      db,
		querier: querier,
		debug:   debug,
		logger:  logger.WithName("db_client"),
	}

	if debug {
		c.logger.SetLevel(v11.DebugLevel)
	}

	c.logger.Debug("migrating querier")
	if err := c.querier.Migrate(ctx); err != nil {
		return nil, err
	}
	c.logger.Debug("querier migrated!")

	return c, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildClientDeclaration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildClientDeclaration(proj)

		expected := `
package example

import (
	"database/sql"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// Client is a wrapper around a database querier. Client is where all
// logging and trace propagation should happen, the querier is where
// the actual database querying is performed.
type Client struct {
	db      *sql.DB
	querier v1.DataManager
	debug   bool
	logger  v11.Logger
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMigrate(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildMigrate(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// Migrate is a simple wrapper around the core querier Migrate call.
func (c *Client) Migrate(ctx context.Context) error {
	ctx, span := tracing.StartSpan(ctx, "Migrate")
	defer span.End()

	return c.querier.Migrate(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildIsReady(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// IsReady is a simple wrapper around the core querier IsReady call.
func (c *Client) IsReady(ctx context.Context) (ready bool) {
	ctx, span := tracing.StartSpan(ctx, "IsReady")
	defer span.End()

	return c.querier.IsReady(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseClient(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// ProvideDatabaseClient provides a new DataManager client.
func ProvideDatabaseClient(
	ctx context.Context,
	db *sql.DB,
	querier v1.DataManager,
	debug bool,
	logger v11.Logger,
) (v1.DataManager, error) {
	c := &Client{
		db:      db,
		querier: querier,
		debug:   debug,
		logger:  logger.WithName("db_client"),
	}

	if debug {
		c.logger.SetLevel(v11.DebugLevel)
	}

	c.logger.Debug("migrating querier")
	if err := c.querier.Migrate(ctx); err != nil {
		return nil, err
	}
	c.logger.Debug("querier migrated!")

	return c, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
