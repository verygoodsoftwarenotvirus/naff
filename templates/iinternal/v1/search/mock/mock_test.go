package mocksearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mockDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
)

var _ search.IndexManager = (*IndexManager)(nil)

// IndexManager is a mock IndexManager
type IndexManager struct {
	mock.Mock
}

// Index implements our interface
func (m *IndexManager) Index(ctx context.Context, id uint64, value interface{}) error {
	args := m.Called(ctx, id, value)
	return args.Error(0)
}

// Search implements our interface
func (m *IndexManager) Search(ctx context.Context, query string, userID uint64) (ids []uint64, err error) {
	args := m.Called(ctx, query, userID)
	return args.Get(0).([]uint64), args.Error(1)
}

// Delete implements our interface
func (m *IndexManager) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInterfaceImplementationStatement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildInterfaceImplementationStatement(proj)

		expected := `
package example

import (
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
)

var _ search.IndexManager = (*IndexManager)(nil)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIndexManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildIndexManager()

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
)

// IndexManager is a mock IndexManager
type IndexManager struct {
	mock.Mock
}

// Index implements our interface
func (m *IndexManager) Index(ctx context.Context, id uint64, value interface{}) error {
	args := m.Called(ctx, id, value)
	return args.Error(0)
}

// Search implements our interface
func (m *IndexManager) Search(ctx context.Context, query string, userID uint64) (ids []uint64, err error) {
	args := m.Called(ctx, query, userID)
	return args.Get(0).([]uint64), args.Error(1)
}

// Delete implements our interface
func (m *IndexManager) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
