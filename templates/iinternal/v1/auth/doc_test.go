package auth

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
Package auth provides functions and structures to facilitate salting and hashing passwords, as well as
validating TOTP tokens
*/
package auth

import ()
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
