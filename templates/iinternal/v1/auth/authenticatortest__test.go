package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_authenticatorTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := authenticatorTestDotGo(proj)

		expected := `
package example

import (
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
)

func TestProvideBcryptHashCost(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		auth.ProvideBcryptHashCost()
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func TestBuildProvideBcryptHashCost(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestProvideBcryptHashCost(proj)

		expected := `
package example

import (
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
)

func TestProvideBcryptHashCost(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		auth.ProvideBcryptHashCost()
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
