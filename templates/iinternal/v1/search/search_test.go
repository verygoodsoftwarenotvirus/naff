package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_searchDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := searchDotGo(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

type (
	// IndexPath is a type alias for dependency injection's sake
	IndexPath string

	// IndexName is a type alias for dependency injection's sake
	IndexName string

	// IndexManager is our wrapper interface for a text search index
	IndexManager interface {
		Index(ctx context.Context, id uint64, value interface{}) error
		Search(ctx context.Context, query string, userID uint64) (ids []uint64, err error)
		Delete(ctx context.Context, id uint64) (err error)
	}

	// IndexManagerProvider is a function that provides an IndexManager for a given index.
	IndexManagerProvider func(path IndexPath, name IndexName, logger v1.Logger) (IndexManager, error)
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTypeDefs()

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

type (
	// IndexPath is a type alias for dependency injection's sake
	IndexPath string

	// IndexName is a type alias for dependency injection's sake
	IndexName string

	// IndexManager is our wrapper interface for a text search index
	IndexManager interface {
		Index(ctx context.Context, id uint64, value interface{}) error
		Search(ctx context.Context, query string, userID uint64) (ids []uint64, err error)
		Delete(ctx context.Context, id uint64) (err error)
	}

	// IndexManagerProvider is a function that provides an IndexManager for a given index.
	IndexManagerProvider func(path IndexPath, name IndexName, logger v1.Logger) (IndexManager, error)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
