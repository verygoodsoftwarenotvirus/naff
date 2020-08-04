package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_databaseTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := databaseTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"regexp"
	"strings"
	"testing"
)

const (
	defaultLimit = uint8(20)
)

func buildTestService(t *testing.T) (*Postgres, gosqlmock.Sqlmock) {
	db, mock, err := gosqlmock.New()
	require.NoError(t, err)
	p := ProvidePostgres(true, db, noop.ProvideNoopLogger())
	return p.(*Postgres), mock
}

var (
	sqlMockReplacer = strings.NewReplacer(
		"$", ` + "`" + `\$` + "`" + `,
		"(", ` + "`" + `\(` + "`" + `,
		")", ` + "`" + `\)` + "`" + `,
		"=", ` + "`" + `\=` + "`" + `,
		"*", ` + "`" + `\*` + "`" + `,
		".", ` + "`" + `\.` + "`" + `,
		"+", ` + "`" + `\+` + "`" + `,
		"?", ` + "`" + `\?` + "`" + `,
		",", ` + "`" + `\,` + "`" + `,
		"-", ` + "`" + `\-` + "`" + `,
		"[", ` + "`" + `\[` + "`" + `,
		"]", ` + "`" + `\]` + "`" + `,
	)
	queryArgRegexp = regexp.MustCompile(` + "`" + `\$\d+` + "`" + `)
)

func formatQueryForSQLMock(query string) string {
	return sqlMockReplacer.Replace(query)
}

func ensureArgCountMatchesQuery(t *testing.T, query string, args []interface{}) {
	t.Helper()

	queryArgCount := len(queryArgRegexp.FindAllString(query, -1))

	if len(args) > 0 {
		assert.Equal(t, queryArgCount, len(args))
	} else {
		assert.Zero(t, queryArgCount)
	}
}

func TestProvidePostgres(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		buildTestService(t)
	})
}

func TestPostgres_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		p, _ := buildTestService(t)
		assert.True(t, p.IsReady(ctx))
	})
}

func TestPostgres_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		p, _ := buildTestService(t)
		p.logQueryBuildingError(errors.New("blah"))
	})
}

func Test_joinUint64s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := []uint64{123, 456, 789}

		expected := "123,456,789"
		actual := joinUint64s(exampleInput)

		assert.Equal(t, expected, actual, "expected %s to equal %s", expected, actual)
	})
}

func TestProvidePostgresDB(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_, err := ProvidePostgresDB(noop.ProvideNoopLogger(), "")
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := databaseTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql/driver"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"regexp"
	"strings"
	"testing"
)

const (
	defaultLimit = uint8(20)
)

func buildTestService(t *testing.T) (*Sqlite, gosqlmock.Sqlmock) {
	db, mock, err := gosqlmock.New()
	require.NoError(t, err)
	s := ProvideSqlite(true, db, noop.ProvideNoopLogger())
	return s.(*Sqlite), mock
}

var (
	sqlMockReplacer = strings.NewReplacer(
		"$", ` + "`" + `\$` + "`" + `,
		"(", ` + "`" + `\(` + "`" + `,
		")", ` + "`" + `\)` + "`" + `,
		"=", ` + "`" + `\=` + "`" + `,
		"*", ` + "`" + `\*` + "`" + `,
		".", ` + "`" + `\.` + "`" + `,
		"+", ` + "`" + `\+` + "`" + `,
		"?", ` + "`" + `\?` + "`" + `,
		",", ` + "`" + `\,` + "`" + `,
		"-", ` + "`" + `\-` + "`" + `,
		"[", ` + "`" + `\[` + "`" + `,
		"]", ` + "`" + `\]` + "`" + `,
	)
	queryArgRegexp = regexp.MustCompile(` + "`" + `\?+` + "`" + `)
)

func formatQueryForSQLMock(query string) string {
	return sqlMockReplacer.Replace(query)
}

func ensureArgCountMatchesQuery(t *testing.T, query string, args []interface{}) {
	t.Helper()

	queryArgCount := len(queryArgRegexp.FindAllString(query, -1))

	if len(args) > 0 {
		assert.Equal(t, queryArgCount, len(args))
	} else {
		assert.Zero(t, queryArgCount)
	}
}

func interfacesToDriverValues(in []interface{}) (out []driver.Value) {
	for _, x := range in {
		out = append(out, driver.Value(x))
	}
	return out
}

func TestProvideSqlite(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		buildTestService(t)
	})
}

func TestSqlite_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		s, _ := buildTestService(t)
		assert.True(t, s.IsReady(ctx))
	})
}

func TestSqlite_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		s, _ := buildTestService(t)
		s.logQueryBuildingError(errors.New("blah"))
	})
}

func TestSqlite_logIDRetrievalError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		s, _ := buildTestService(t)
		s.logIDRetrievalError(errors.New("blah"))
	})
}

func TestProvideSqliteDB(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_, err := ProvideSqliteDB(noop.ProvideNoopLogger(), "")
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := databaseTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql/driver"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"regexp"
	"strings"
	"testing"
)

const (
	defaultLimit = uint8(20)
)

func buildTestService(t *testing.T) (*MariaDB, gosqlmock.Sqlmock) {
	db, mock, err := gosqlmock.New()
	require.NoError(t, err)
	m := ProvideMariaDB(true, db, noop.ProvideNoopLogger())
	return m.(*MariaDB), mock
}

var (
	sqlMockReplacer = strings.NewReplacer(
		"$", ` + "`" + `\$` + "`" + `,
		"(", ` + "`" + `\(` + "`" + `,
		")", ` + "`" + `\)` + "`" + `,
		"=", ` + "`" + `\=` + "`" + `,
		"*", ` + "`" + `\*` + "`" + `,
		".", ` + "`" + `\.` + "`" + `,
		"+", ` + "`" + `\+` + "`" + `,
		"?", ` + "`" + `\?` + "`" + `,
		",", ` + "`" + `\,` + "`" + `,
		"-", ` + "`" + `\-` + "`" + `,
		"[", ` + "`" + `\[` + "`" + `,
		"]", ` + "`" + `\]` + "`" + `,
	)
	queryArgRegexp = regexp.MustCompile(` + "`" + `\?+` + "`" + `)
)

func formatQueryForSQLMock(query string) string {
	return sqlMockReplacer.Replace(query)
}

func ensureArgCountMatchesQuery(t *testing.T, query string, args []interface{}) {
	t.Helper()

	queryArgCount := len(queryArgRegexp.FindAllString(query, -1))

	if len(args) > 0 {
		assert.Equal(t, queryArgCount, len(args))
	} else {
		assert.Zero(t, queryArgCount)
	}
}

func interfacesToDriverValues(in []interface{}) (out []driver.Value) {
	for _, x := range in {
		out = append(out, driver.Value(x))
	}
	return out
}

func TestProvideMariaDB(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		buildTestService(t)
	})
}

func TestMariaDB_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		m, _ := buildTestService(t)
		assert.True(t, m.IsReady(ctx))
	})
}

func TestMariaDB_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		m, _ := buildTestService(t)
		m.logQueryBuildingError(errors.New("blah"))
	})
}

func TestMariaDB_logIDRetrievalError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		m, _ := buildTestService(t)
		m.logIDRetrievalError(errors.New("blah"))
	})
}

func TestProvideMariaDBConnection(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_, err := ProvideMariaDBConnection(noop.ProvideNoopLogger(), "")
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConstDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildConstDecls()

		expected := `
package example

import ()

const (
	defaultLimit = uint8(20)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildTestService(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildTestService(dbvendor)

		expected := `
package example

import (
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func buildTestService(t *testing.T) (*Postgres, gosqlmock.Sqlmock) {
	db, mock, err := gosqlmock.New()
	require.NoError(t, err)
	p := ProvidePostgres(true, db, noop.ProvideNoopLogger())
	return p.(*Postgres), mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildTestService(dbvendor)

		expected := `
package example

import (
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func buildTestService(t *testing.T) (*Sqlite, gosqlmock.Sqlmock) {
	db, mock, err := gosqlmock.New()
	require.NoError(t, err)
	s := ProvideSqlite(true, db, noop.ProvideNoopLogger())
	return s.(*Sqlite), mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildTestService(dbvendor)

		expected := `
package example

import (
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func buildTestService(t *testing.T) (*MariaDB, gosqlmock.Sqlmock) {
	db, mock, err := gosqlmock.New()
	require.NoError(t, err)
	m := ProvideMariaDB(true, db, noop.ProvideNoopLogger())
	return m.(*MariaDB), mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDBVendorTestVarDecls(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildDBVendorTestVarDecls(dbvendor)

		expected := `
package example

import (
	"regexp"
	"strings"
)

var (
	sqlMockReplacer = strings.NewReplacer(
		"$", ` + "`" + `\$` + "`" + `,
		"(", ` + "`" + `\(` + "`" + `,
		")", ` + "`" + `\)` + "`" + `,
		"=", ` + "`" + `\=` + "`" + `,
		"*", ` + "`" + `\*` + "`" + `,
		".", ` + "`" + `\.` + "`" + `,
		"+", ` + "`" + `\+` + "`" + `,
		"?", ` + "`" + `\?` + "`" + `,
		",", ` + "`" + `\,` + "`" + `,
		"-", ` + "`" + `\-` + "`" + `,
		"[", ` + "`" + `\[` + "`" + `,
		"]", ` + "`" + `\]` + "`" + `,
	)
	queryArgRegexp = regexp.MustCompile(` + "`" + `\$\d+` + "`" + `)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildDBVendorTestVarDecls(dbvendor)

		expected := `
package example

import (
	"regexp"
	"strings"
)

var (
	sqlMockReplacer = strings.NewReplacer(
		"$", ` + "`" + `\$` + "`" + `,
		"(", ` + "`" + `\(` + "`" + `,
		")", ` + "`" + `\)` + "`" + `,
		"=", ` + "`" + `\=` + "`" + `,
		"*", ` + "`" + `\*` + "`" + `,
		".", ` + "`" + `\.` + "`" + `,
		"+", ` + "`" + `\+` + "`" + `,
		"?", ` + "`" + `\?` + "`" + `,
		",", ` + "`" + `\,` + "`" + `,
		"-", ` + "`" + `\-` + "`" + `,
		"[", ` + "`" + `\[` + "`" + `,
		"]", ` + "`" + `\]` + "`" + `,
	)
	queryArgRegexp = regexp.MustCompile(` + "`" + `\?+` + "`" + `)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildDBVendorTestVarDecls(dbvendor)

		expected := `
package example

import (
	"regexp"
	"strings"
)

var (
	sqlMockReplacer = strings.NewReplacer(
		"$", ` + "`" + `\$` + "`" + `,
		"(", ` + "`" + `\(` + "`" + `,
		")", ` + "`" + `\)` + "`" + `,
		"=", ` + "`" + `\=` + "`" + `,
		"*", ` + "`" + `\*` + "`" + `,
		".", ` + "`" + `\.` + "`" + `,
		"+", ` + "`" + `\+` + "`" + `,
		"?", ` + "`" + `\?` + "`" + `,
		",", ` + "`" + `\,` + "`" + `,
		"-", ` + "`" + `\-` + "`" + `,
		"[", ` + "`" + `\[` + "`" + `,
		"]", ` + "`" + `\]` + "`" + `,
	)
	queryArgRegexp = regexp.MustCompile(` + "`" + `\?+` + "`" + `)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFormatQueryForSQLMock(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildFormatQueryForSQLMock(dbvendor)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildFormatQueryForSQLMock(dbvendor)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildFormatQueryForSQLMock(dbvendor)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEnsureArgCountMatchesQuery(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildEnsureArgCountMatchesQuery()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func ensureArgCountMatchesQuery(t *testing.T, query string, args []interface{}) {
	t.Helper()

	queryArgCount := len(queryArgRegexp.FindAllString(query, -1))

	if len(args) > 0 {
		assert.Equal(t, queryArgCount, len(args))
	} else {
		assert.Zero(t, queryArgCount)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInterfacesToDriverValues(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildInterfacesToDriverValues()

		expected := `
package example

import (
	"database/sql/driver"
)

func interfacesToDriverValues(in []interface{}) (out []driver.Value) {
	for _, x := range in {
		out = append(out, driver.Value(x))
	}
	return out
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDBVendorProviderTest(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildDBVendorProviderTest(dbvendor)

		expected := `
package example

import (
	"testing"
)

func TestProvidePostgres(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		buildTestService(t)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildDBVendorProviderTest(dbvendor)

		expected := `
package example

import (
	"testing"
)

func TestProvideSqlite(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		buildTestService(t)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildDBVendorProviderTest(dbvendor)

		expected := `
package example

import (
	"testing"
)

func TestProvideMariaDB(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		buildTestService(t)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBVendor_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildTestDBVendor_IsReady(dbvendor)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		p, _ := buildTestService(t)
		assert.True(t, p.IsReady(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildTestDBVendor_IsReady(dbvendor)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		s, _ := buildTestService(t)
		assert.True(t, s.IsReady(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildTestDBVendor_IsReady(dbvendor)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestMariaDB_IsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		m, _ := buildTestService(t)
		assert.True(t, m.IsReady(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBVendor_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildTestDBVendor_logQueryBuildingError(dbvendor)

		expected := `
package example

import (
	"errors"
	"testing"
)

func TestPostgres_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		p, _ := buildTestService(t)
		p.logQueryBuildingError(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildTestDBVendor_logQueryBuildingError(dbvendor)

		expected := `
package example

import (
	"errors"
	"testing"
)

func TestSqlite_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		s, _ := buildTestService(t)
		s.logQueryBuildingError(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildTestDBVendor_logQueryBuildingError(dbvendor)

		expected := `
package example

import (
	"errors"
	"testing"
)

func TestMariaDB_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		m, _ := buildTestService(t)
		m.logQueryBuildingError(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBVendor_logIDRetrievalError(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildTestDBVendor_logIDRetrievalError(dbvendor)

		expected := `
package example

import (
	"errors"
	"testing"
)

func TestPostgres_logIDRetrievalError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		p, _ := buildTestService(t)
		p.logIDRetrievalError(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildTestDBVendor_logIDRetrievalError(dbvendor)

		expected := `
package example

import (
	"errors"
	"testing"
)

func TestSqlite_logIDRetrievalError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		s, _ := buildTestService(t)
		s.logIDRetrievalError(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildTestDBVendor_logIDRetrievalError(dbvendor)

		expected := `
package example

import (
	"errors"
	"testing"
)

func TestMariaDB_logIDRetrievalError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		m, _ := buildTestService(t)
		m.logIDRetrievalError(errors.New("blah"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_joinUint64s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTest_joinUint64s()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func Test_joinUint64s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		exampleInput := []uint64{123, 456, 789}

		expected := "123,456,789"
		actual := joinUint64s(exampleInput)

		assert.Equal(t, expected, actual, "expected %s to equal %s", expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

}

func Test_buildTestProviderFunc(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildTestProviderFunc(dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvidePostgresDB(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_, err := ProvidePostgresDB(noop.ProvideNoopLogger(), "")
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildTestProviderFunc(dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvideSqliteDB(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_, err := ProvideSqliteDB(noop.ProvideNoopLogger(), "")
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildTestProviderFunc(dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"testing"
)

func TestProvideMariaDBConnection(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		_, err := ProvideMariaDBConnection(noop.ProvideNoopLogger(), "")
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("invalid", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("invalid")

		assert.Panics(t, func() { buildTestProviderFunc(dbvendor) })
	})
}
