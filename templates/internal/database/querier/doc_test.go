package querier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_docDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := docDotGo()

		expected := `
/*
Package querier provides an abstraction around database queriers. The primary
purpose of this abstraction is to contain all the necessary logging and tracing
steps in a single place, so that the actual package that is responsible for
executing queries and loading their return values into structs isn't burdened
with inconsistent logging
*/
package querier

import ()
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
