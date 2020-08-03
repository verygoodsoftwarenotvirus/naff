package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mainDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mainDotGo(proj)

		expected := `
package example

import (
	"fmt"
)

const (
	// SortAscending is the pre-determined Ascending sortType for external use
	SortAscending sortType = "asc"
	// SortDescending is the pre-determined Descending sortType for external use
	SortDescending sortType = "desc"
)

type (
	// ContextKey represents strings to be used in Context objects. From the docs:
	// 		"The provided key must be comparable and should not be of type string or
	// 		any other built-in type to avoid collisions between packages using context."
	ContextKey string
	sortType   string

	// Pagination represents a pagination request.
	Pagination struct {
		Page  uint64 ` + "`" + `json:"page"` + "`" + `
		Limit uint8  ` + "`" + `json:"limit"` + "`" + `
	}

	// CountResponse is what we respond with when a user requests a count of data types.
	CountResponse struct {
		Count uint64 ` + "`" + `json:"count"` + "`" + `
	}

	// ErrorResponse represents a response we might send to the user in the event of an error.
	ErrorResponse struct {
		Message string ` + "`" + `json:"message"` + "`" + `
		Code    uint   ` + "`" + `json:"code"` + "`" + `
	}
)

var _ error = (*ErrorResponse)(nil)

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", er.Code, er.Message)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMainConstantDefs()

		expected := `
package example

import ()

const (
	// SortAscending is the pre-determined Ascending sortType for external use
	SortAscending sortType = "asc"
	// SortDescending is the pre-determined Descending sortType for external use
	SortDescending sortType = "desc"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMainTypeDefs()

		expected := `
package example

import ()

type (
	// ContextKey represents strings to be used in Context objects. From the docs:
	// 		"The provided key must be comparable and should not be of type string or
	// 		any other built-in type to avoid collisions between packages using context."
	ContextKey string
	sortType   string

	// Pagination represents a pagination request.
	Pagination struct {
		Page  uint64 ` + "`" + `json:"page"` + "`" + `
		Limit uint8  ` + "`" + `json:"limit"` + "`" + `
	}

	// CountResponse is what we respond with when a user requests a count of data types.
	CountResponse struct {
		Count uint64 ` + "`" + `json:"count"` + "`" + `
	}

	// ErrorResponse represents a response we might send to the user in the event of an error.
	ErrorResponse struct {
		Message string ` + "`" + `json:"message"` + "`" + `
		Code    uint   ` + "`" + `json:"code"` + "`" + `
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainErrorInterfaceImplementation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMainErrorInterfaceImplementation()

		expected := `
package example

import ()

var _ error = (*ErrorResponse)(nil)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainErrorResponseDotError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMainErrorResponseDotError()

		expected := `
package example

import (
	"fmt"
)

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", er.Code, er.Message)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
