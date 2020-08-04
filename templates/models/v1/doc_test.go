package v1

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
Package models provides the core data models for the service
*/
package models

import ()
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
