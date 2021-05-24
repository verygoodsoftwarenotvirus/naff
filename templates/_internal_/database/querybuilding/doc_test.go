package querybuilding

import (
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_docDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := "postgres"
		dbDesc := ""

		x := docDotGo(dbvendor, dbDesc)

		expected := `
/*
Package postgres provides a Database implementation that is compatible with
*/
package postgres

import ()
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := "sqlite"
		dbDesc := ""

		x := docDotGo(dbvendor, dbDesc)

		expected := `
/*
Package sqlite provides a Database implementation that is compatible with
*/
package sqlite

import ()
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := "mariadb"
		dbDesc := ""

		x := docDotGo(dbvendor, dbDesc)

		expected := `
/*
Package mariadb provides a Database implementation that is compatible with
*/
package mariadb

import ()
`
		actual := testutils.RenderFileToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
