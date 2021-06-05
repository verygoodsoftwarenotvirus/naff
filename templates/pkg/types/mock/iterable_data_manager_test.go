package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockIterableDataManagerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := mockIterableDataManagerDotGo(proj, typ)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

var _ v1.ItemDataManager = (*ItemDataManager)(nil)

// ItemDataManager is a mocked models.ItemDataManager for testing.
type ItemDataManager struct {
	mock.Mock
}

// ItemExists is a mock function.
func (m *ItemDataManager) ItemExists(ctx context.Context, itemID, userID uint64) (bool, error) {
	args := m.Called(ctx, itemID, userID)
	return args.Bool(0), args.Error(1)
}

// GetItem is a mock function.
func (m *ItemDataManager) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	args := m.Called(ctx, itemID, userID)
	return args.Get(0).(*v1.Item), args.Error(1)
}

// GetAllItemsCount is a mock function.
func (m *ItemDataManager) GetAllItemsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllItems is a mock function.
func (m *ItemDataManager) GetAllItems(ctx context.Context, results chan []v1.Item) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetItems is a mock function.
func (m *ItemDataManager) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*v1.ItemList), args.Error(1)
}

// GetItemsWithIDs is a mock function.
func (m *ItemDataManager) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v1.Item, error) {
	args := m.Called(ctx, userID, limit, ids)
	return args.Get(0).([]v1.Item), args.Error(1)
}

// CreateItem is a mock function.
func (m *ItemDataManager) CreateItem(ctx context.Context, input *v1.ItemCreationInput) (*v1.Item, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.Item), args.Error(1)
}

// UpdateItem is a mock function.
func (m *ItemDataManager) UpdateItem(ctx context.Context, updated *v1.Item) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveItem is a mock function.
func (m *ItemDataManager) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	return m.Called(ctx, itemID, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingExists(proj, typ)

		expected := `
package example

import (
	"context"
)

// ItemExists is a mock function.
func (m *ItemDataManager) ItemExists(ctx context.Context, itemID, userID uint64) (bool, error) {
	args := m.Called(ctx, itemID, userID)
	return args.Bool(0), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomething(proj, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItem is a mock function.
func (m *ItemDataManager) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	args := m.Called(ctx, itemID, userID)
	return args.Get(0).(*v1.Item), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllSomethingsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllSomethingsCount(proj, typ)

		expected := `
package example

import (
	"context"
)

// GetAllItemsCount is a mock function.
func (m *ItemDataManager) GetAllItemsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllSomethings(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllSomethings(proj, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetAllItems is a mock function.
func (m *ItemDataManager) GetAllItems(ctx context.Context, results chan []v1.Item) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetListOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomething(proj, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItems is a mock function.
func (m *ItemDataManager) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*v1.ItemList), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetSomethingsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingsWithIDs(proj, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItemsWithIDs is a mock function.
func (m *ItemDataManager) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v1.Item, error) {
	args := m.Called(ctx, userID, limit, ids)
	return args.Get(0).([]v1.Item), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomething(proj, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateItem is a mock function.
func (m *ItemDataManager) CreateItem(ctx context.Context, input *v1.ItemCreationInput) (*v1.Item, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.Item), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomething(proj, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// UpdateItem is a mock function.
func (m *ItemDataManager) UpdateItem(ctx context.Context, updated *v1.Item) error {
	return m.Called(ctx, updated).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomething(typ)

		expected := `
package example

import (
	"context"
)

// ArchiveItem is a mock function.
func (m *ItemDataManager) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	return m.Called(ctx, itemID, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
