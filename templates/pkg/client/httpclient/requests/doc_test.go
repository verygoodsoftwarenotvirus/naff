package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_docDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()

		expected := `
/*
Package requests provides an HTTP client that can communicate with and interpret the responses
of an instance of the todo service.
*/
package requests

import ()
`

		actual := testutils.RenderFileToString(t, docDotGo(proj))

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
