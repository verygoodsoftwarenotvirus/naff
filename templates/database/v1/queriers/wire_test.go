package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj, dbvendor)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers is what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvidePostgresDB,
		ProvidePostgres,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj, dbvendor)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers is what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvideSqliteDB,
		ProvideSqlite,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj, dbvendor)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers is what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvideMariaDBConnection,
		ProvideMariaDB,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
