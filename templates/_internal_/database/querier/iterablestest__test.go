package querier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterablesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_ItemExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("ItemExists", mock.Anything, exampleItem.ID, exampleItem.BelongsToUser).Return(true, nil)

		actual, err := c.ItemExists(ctx, exampleItem.ID, exampleItem.BelongsToUser)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItem", mock.Anything, exampleItem.ID, exampleItem.BelongsToUser).Return(exampleItem, nil)

		actual, err := c.GetItem(ctx, exampleItem.ID, exampleItem.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetAllItemsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllItemsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllItems(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []v1.Item)

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetAllItems", mock.Anything, results).Return(nil)

		err := c.GetAllItems(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetItems(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()
		exampleItemList := fake.BuildFakeItemList()

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItems", mock.Anything, exampleUser.ID, filter).Return(exampleItemList, nil)

		actual, err := c.GetItems(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)
		exampleItemList := fake.BuildFakeItemList()

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItems", mock.Anything, exampleUser.ID, filter).Return(exampleItemList, nil)

		actual, err := c.GetItems(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleIDs []uint64
		for _, x := range exampleItemList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItemsWithIDs", mock.Anything, exampleUser.ID, defaultLimit, exampleIDs).Return(exampleItemList, nil)

		actual, err := c.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("CreateItem", mock.Anything, exampleInput).Return(exampleItem, nil)

		actual, err := c.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()

		mockDB.ItemDataManager.On("UpdateItem", mock.Anything, exampleItem).Return(expected)

		err := c.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("ArchiveItem", mock.Anything, exampleItem.ID, exampleItem.BelongsToUser).Return(expected)

		err := c.ArchiveItem(ctx, exampleItem.ID, exampleItem.BelongsToUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientSomethingExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientSomethingExists(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_ItemExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("ItemExists", mock.Anything, exampleItem.ID, exampleItem.BelongsToUser).Return(true, nil)

		actual, err := c.ItemExists(ctx, exampleItem.ID, exampleItem.BelongsToUser)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientGetSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientGetSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItem", mock.Anything, exampleItem.ID, exampleItem.BelongsToUser).Return(exampleItem, nil)

		actual, err := c.GetItem(ctx, exampleItem.ID, exampleItem.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientGetAllOfSomethingCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientGetAllOfSomethingCount(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"testing"
)

func TestClient_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetAllItemsCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllItemsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientGetAllOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientGetAllOfSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"testing"
)

func TestClient_GetAllItems(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		results := make(chan []v1.Item)

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetAllItems", mock.Anything, results).Return(nil)

		err := c.GetAllItems(ctx, results)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientGetListOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientGetListOfSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetItems(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()
		exampleItemList := fake.BuildFakeItemList()

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItems", mock.Anything, exampleUser.ID, filter).Return(exampleItemList, nil)

		actual, err := c.GetItems(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)
		exampleItemList := fake.BuildFakeItemList()

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItems", mock.Anything, exampleUser.ID, filter).Return(exampleItemList, nil)

		actual, err := c.GetItems(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with type not belonging to user", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.RestrictedToUser = false
		x := buildTestClientGetListOfSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetItems(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()
		exampleItemList := fake.BuildFakeItemList()

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItems", mock.Anything, filter).Return(exampleItemList, nil)

		actual, err := c.GetItems(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)
		exampleItemList := fake.BuildFakeItemList()

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItems", mock.Anything, filter).Return(exampleItemList, nil)

		actual, err := c.GetItems(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientGetListOfSomethingWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientGetListOfSomethingWithIDs(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItemList := fake.BuildFakeItemList().Items
		var exampleIDs []uint64
		for _, x := range exampleItemList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItemsWithIDs", mock.Anything, exampleUser.ID, defaultLimit, exampleIDs).Return(exampleItemList, nil)

		actual, err := c.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with type not belonging to user", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		x := buildTestClientGetListOfSomethingWithIDs(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleItemList := fake.BuildFakeItemList().Items
		var exampleIDs []uint64
		for _, x := range exampleItemList {
			exampleIDs = append(exampleIDs, x.ID)
		}

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("GetItemsWithIDs", mock.Anything, defaultLimit, exampleIDs).Return(exampleItemList, nil)

		actual, err := c.GetItemsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientCreateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientCreateSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_CreateItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("CreateItem", mock.Anything, exampleInput).Return(exampleItem, nil)

		actual, err := c.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientUpdateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientUpdateSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_UpdateItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()
		var expected error

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()

		mockDB.ItemDataManager.On("UpdateItem", mock.Anything, exampleItem).Return(expected)

		err := c.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClientArchiveSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestClientArchiveSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_ArchiveItem(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		var expected error

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		c, mockDB := buildTestClient()
		mockDB.ItemDataManager.On("ArchiveItem", mock.Anything, exampleItem.ID, exampleItem.BelongsToUser).Return(expected)

		err := c.ArchiveItem(ctx, exampleItem.ID, exampleItem.BelongsToUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
