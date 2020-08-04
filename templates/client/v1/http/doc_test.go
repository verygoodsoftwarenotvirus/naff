package client

import (
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_docDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		expected := `
/*
Package client provides an HTTP client that can communicate with and interpret the responses
of an instance of the todo service.
*/
package client

import ()
`

		actual := testutils.RenderFileToString(t, docDotGo())

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
