package webhooks

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
Package webhooks provides a series of HTTP handlers for managing webhooks in a compatible database.
*/
package webhooks

import ()
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
