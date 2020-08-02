package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_timeTellerTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := timeTellerTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func Test_stdLibTimeTeller_Now(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		tt := &stdLibTimeTeller{}

		assert.NotZero(t, tt.Now())
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
		x := timeTellerTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func Test_stdLibTimeTeller_Now(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		tt := &stdLibTimeTeller{}

		assert.NotZero(t, tt.Now())
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
		x := timeTellerTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func Test_stdLibTimeTeller_Now(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		tt := &stdLibTimeTeller{}

		assert.NotZero(t, tt.Now())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
