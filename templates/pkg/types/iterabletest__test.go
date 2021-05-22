package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterableTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterableTestDotGo(proj, typ)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Item{}

		expected := &ItemUpdateInput{
			Name:    v5.Word(),
			Details: v5.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Details, i.Details)
	})
}

func TestItem_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		item := &Item{
			Name:    v5.Word(),
			Details: v5.Word(),
		}

		expected := &ItemUpdateInput{
			Name:    item.Name,
			Details: item.Details,
		}
		actual := item.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestSomething_Update(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestSomething_Update(typ)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Item{}

		expected := &ItemUpdateInput{
			Name:    v5.Word(),
			Details: v5.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Details, i.Details)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestSomething_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestSomething_ToUpdateInput(proj, typ)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		item := &Item{
			Name:    v5.Word(),
			Details: v5.Word(),
		}

		expected := &ItemUpdateInput{
			Name:    item.Name,
			Details: item.Details,
		}
		actual := item.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
