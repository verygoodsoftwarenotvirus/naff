package querybuilding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_timeTellerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := timeTellerDotGo(proj, dbvendor)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	"time"
)

type timeTeller interface {
	Now() uint64
}

type stdLibTimeTeller struct{}

func (t *stdLibTimeTeller) Now() uint64 {
	return uint64(time.Now().Unix())
}

type mockTimeTeller struct {
	mock.Mock
}

func (m *mockTimeTeller) Now() uint64 {
	return m.Called().Get(0).(uint64)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := timeTellerDotGo(proj, dbvendor)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	"time"
)

type timeTeller interface {
	Now() uint64
}

type stdLibTimeTeller struct{}

func (t *stdLibTimeTeller) Now() uint64 {
	return uint64(time.Now().Unix())
}

type mockTimeTeller struct {
	mock.Mock
}

func (m *mockTimeTeller) Now() uint64 {
	return m.Called().Get(0).(uint64)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := timeTellerDotGo(proj, dbvendor)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	"time"
)

type timeTeller interface {
	Now() uint64
}

type stdLibTimeTeller struct{}

func (t *stdLibTimeTeller) Now() uint64 {
	return uint64(time.Now().Unix())
}

type mockTimeTeller struct {
	mock.Mock
}

func (m *mockTimeTeller) Now() uint64 {
	return m.Called().Get(0).(uint64)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
