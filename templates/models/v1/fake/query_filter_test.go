package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_queryFilterDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := queryFilterDotGo(proj)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *v1.QueryFilter {
	return &v1.QueryFilter{
		Page:          10,
		Limit:         20,
		CreatedAfter:  uint64(uint32(v5.Date().Unix())),
		CreatedBefore: uint64(uint32(v5.Date().Unix())),
		UpdatedAfter:  uint64(uint32(v5.Date().Unix())),
		UpdatedBefore: uint64(uint32(v5.Date().Unix())),
		SortBy:        v1.SortAscending,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildFleshedOutQueryFilter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildFleshedOutQueryFilter(proj)

		expected := `
package example

import (
	v5 "github.com/brianvoe/gofakeit/v5"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter.
func BuildFleshedOutQueryFilter() *v1.QueryFilter {
	return &v1.QueryFilter{
		Page:          10,
		Limit:         20,
		CreatedAfter:  uint64(uint32(v5.Date().Unix())),
		CreatedBefore: uint64(uint32(v5.Date().Unix())),
		UpdatedAfter:  uint64(uint32(v5.Date().Unix())),
		UpdatedBefore: uint64(uint32(v5.Date().Unix())),
		SortBy:        v1.SortAscending,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
