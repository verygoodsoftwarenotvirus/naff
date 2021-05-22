package config_gen

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
Command config_gen generates configuration files in the local repository, configured
via the precise mechanism that parses them to guard against invalid configuration
*/
package main

import ()
`
		actual := testutils.RenderFileToString(t, docDotGo())

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
