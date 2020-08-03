package queriers

import (
	"github.com/Masterminds/squirrel"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_isPostgres(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		actual := isPostgres(dbvendor)

		assert.True(t, actual)
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		actual := isPostgres(dbvendor)

		assert.False(t, actual)
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		actual := isPostgres(dbvendor)

		assert.False(t, actual)
	})
}

func Test_isSqlite(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		actual := isSqlite(dbvendor)

		assert.False(t, actual)
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		actual := isSqlite(dbvendor)

		assert.True(t, actual)
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		actual := isSqlite(dbvendor)

		assert.False(t, actual)
	})
}

func Test_isMariaDB(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		actual := isMariaDB(dbvendor)

		assert.False(t, actual)
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		actual := isMariaDB(dbvendor)

		assert.False(t, actual)
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		actual := isMariaDB(dbvendor)

		assert.True(t, actual)
	})
}

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()
		assert.NoError(t, RenderPackage(proj))
	})

	T.Run("with invalid output directory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = `/\0/\0/\0`

		assert.Error(t, RenderPackage(proj))
	})
}

func TestGetOAuth2ClientPalabra(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := GetOAuth2ClientPalabra()

		expected := `OAuth2Client`
		actual := x.Singular()

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func TestGetUserPalabra(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := GetUserPalabra()

		expected := `User`
		actual := x.Singular()

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func TestGetWebhookPalabra(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := GetWebhookPalabra()

		expected := `Webhook`
		actual := x.Singular()

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_renderDatabasePackage(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := postgres
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()

		assert.NoError(t, renderDatabasePackage(proj, dbvendor))
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := sqlite
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()
		assert.NoError(t, renderDatabasePackage(proj, dbvendor))
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := mariadb
		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()
		assert.NoError(t, renderDatabasePackage(proj, dbvendor))
	})
}

func Test_buildMariaDBWord(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMariaDBWord()

		expected := `MariaDB`
		actual := x.Singular()

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_convertArgsToCode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := convertArgsToCode([]interface{}{"things", "and", "stuff"})

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		actual := testutils.RenderCallArgsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildQueryTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		qb := squirrel.Select("*")
		x := buildQueryTest(dbvendor, "Example", qb, nil, nil, nil)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_buildExampleQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT *"
		actualQuery := p.buildExampleQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_unixTimeForDatabase(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		expected := `extract(epoch FROM NOW())`
		actual := unixTimeForDatabase(dbvendor)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		expected := `(strftime('%s','now'))`
		actual := unixTimeForDatabase(dbvendor)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Mariadb")

		expected := `UNIX_TIMESTAMP()`
		actual := unixTimeForDatabase(dbvendor)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_queryBuilderForDatabase(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := queryBuilderForDatabase(dbvendor).Select("*").From("table").Where(squirrel.Eq{"id": 123})

		expected := `SELECT * FROM table WHERE id = $1`
		actual, _, err := x.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := queryBuilderForDatabase(dbvendor).Select("*").From("table").Where(squirrel.Eq{"id": 123})

		expected := `SELECT * FROM table WHERE id = ?`
		actual, _, err := x.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Mariadb")

		x := queryBuilderForDatabase(dbvendor).Select("*").From("table").Where(squirrel.Eq{"id": 123})

		expected := `SELECT * FROM table WHERE id = ?`
		actual, _, err := x.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
