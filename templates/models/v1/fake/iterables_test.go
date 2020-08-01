package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterablesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesDotGo(proj, typ)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeItem builds a faked item.
func BuildFakeItem() *v1.Item {
	return &v1.Item{
		ID:            v5.Uint64(),
		Name:          v5.Word(),
		Details:       v5.Word(),
		CreatedOn:     uint64(uint32(v5.Date().Unix())),
		BelongsToUser: v5.Uint64(),
	}
}

// BuildFakeItemList builds a faked ItemList.
func BuildFakeItemList() *v1.ItemList {
	exampleItem1 := BuildFakeItem()
	exampleItem2 := BuildFakeItem()
	exampleItem3 := BuildFakeItem()

	return &v1.ItemList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Items: []v1.Item{
			*exampleItem1,
			*exampleItem2,
			*exampleItem3,
		},
	}
}

// BuildFakeItemUpdateInputFromItem builds a faked ItemUpdateInput from an item.
func BuildFakeItemUpdateInputFromItem(item *v1.Item) *v1.ItemUpdateInput {
	return &v1.ItemUpdateInput{
		Name:          item.Name,
		Details:       item.Details,
		BelongsToUser: item.BelongsToUser,
	}
}

// BuildFakeItemCreationInput builds a faked ItemCreationInput.
func BuildFakeItemCreationInput() *v1.ItemCreationInput {
	item := BuildFakeItem()
	return BuildFakeItemCreationInputFromItem(item)
}

// BuildFakeItemCreationInputFromItem builds a faked ItemCreationInput from an item.
func BuildFakeItemCreationInputFromItem(item *v1.Item) *v1.ItemCreationInput {
	return &v1.ItemCreationInput{
		Name:          item.Name,
		Details:       item.Details,
		BelongsToUser: item.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildFakeSomething(proj, typ)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeItem builds a faked item.
func BuildFakeItem() *v1.Item {
	return &v1.Item{
		ID:            v5.Uint64(),
		Name:          v5.Word(),
		Details:       v5.Word(),
		CreatedOn:     uint64(uint32(v5.Date().Unix())),
		BelongsToUser: v5.Uint64(),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeSomethingList(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildFakeSomethingList(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeItemList builds a faked ItemList.
func BuildFakeItemList() *v1.ItemList {
	exampleItem1 := BuildFakeItem()
	exampleItem2 := BuildFakeItem()
	exampleItem3 := BuildFakeItem()

	return &v1.ItemList{
		Pagination: v1.Pagination{
			Page:  1,
			Limit: 20,
		},
		Items: []v1.Item{
			*exampleItem1,
			*exampleItem2,
			*exampleItem3,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeSomethingUpdateInputFromSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildFakeSomethingUpdateInputFromSomething(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeItemUpdateInputFromItem builds a faked ItemUpdateInput from an item.
func BuildFakeItemUpdateInputFromItem(item *v1.Item) *v1.ItemUpdateInput {
	return &v1.ItemUpdateInput{
		Name:          item.Name,
		Details:       item.Details,
		BelongsToUser: item.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeSomethingCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildFakeSomethingCreationInput(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeItemCreationInput builds a faked ItemCreationInput.
func BuildFakeItemCreationInput() *v1.ItemCreationInput {
	item := BuildFakeItem()
	return BuildFakeItemCreationInputFromItem(item)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFakeSomethingCreationInputFromSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildFakeSomethingCreationInputFromSomething(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFakeItemCreationInputFromItem builds a faked ItemCreationInput from an item.
func BuildFakeItemCreationInputFromItem(item *v1.Item) *v1.ItemCreationInput {
	return &v1.ItemCreationInput{
		Name:          item.Name,
		Details:       item.Details,
		BelongsToUser: item.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
