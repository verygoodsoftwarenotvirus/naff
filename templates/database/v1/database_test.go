package v1

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
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"io"
)

var (
	_ Scanner = (*sql.Row)(nil)
	_ Querier = (*sql.DB)(nil)
	_ Querier = (*sql.Tx)(nil)
)

type (
	// Scanner represents any database response (i.e. sql.Row[s])
	Scanner interface {
		Scan(dest ...interface{}) error
	}

	// ResultIterator represents any iterable database response (i.e. sql.Rows)
	ResultIterator interface {
		Next() bool
		Err() error
		Scanner
		io.Closer
	}

	// Querier is a subset interface for sql.{DB|Tx} objects
	Querier interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	// ConnectionDetails is a string alias for dependency injection.
	ConnectionDetails string

	// DataManager describes anything that stores data for our services.
	DataManager interface {
		Migrate(ctx context.Context) error
		IsReady(ctx context.Context) (ready bool)

		v1.ItemDataManager
		v1.UserDataManager
		v1.OAuth2ClientDataManager
		v1.WebhookDataManager
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildVarDeclarations()

		expected := `
package example

import (
	"database/sql"
)

var (
	_ Scanner = (*sql.Row)(nil)
	_ Querier = (*sql.DB)(nil)
	_ Querier = (*sql.Tx)(nil)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTypeDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTypeDeclarations(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"io"
)

type (
	// Scanner represents any database response (i.e. sql.Row[s])
	Scanner interface {
		Scan(dest ...interface{}) error
	}

	// ResultIterator represents any iterable database response (i.e. sql.Rows)
	ResultIterator interface {
		Next() bool
		Err() error
		Scanner
		io.Closer
	}

	// Querier is a subset interface for sql.{DB|Tx} objects
	Querier interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	// ConnectionDetails is a string alias for dependency injection.
	ConnectionDetails string

	// DataManager describes anything that stores data for our services.
	DataManager interface {
		Migrate(ctx context.Context) error
		IsReady(ctx context.Context) (ready bool)

		v1.ItemDataManager
		v1.UserDataManager
		v1.OAuth2ClientDataManager
		v1.WebhookDataManager
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
