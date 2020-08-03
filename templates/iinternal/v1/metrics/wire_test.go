package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers represents what this library offers to external users in the form of dependencies.
	Providers = wire.NewSet(
		ProvideUnitCounter,
		ProvideUnitCounterProvider,
	)
)

// ProvideUnitCounterProvider provides UnitCounter providers.
func ProvideUnitCounterProvider() UnitCounterProvider {
	return ProvideUnitCounter
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProviders(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildWireProviders()

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers represents what this library offers to external users in the form of dependencies.
	Providers = wire.NewSet(
		ProvideUnitCounter,
		ProvideUnitCounterProvider,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProvideUnitCounterProvider(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildWireProvideUnitCounterProvider()

		expected := `
package example

import ()

// ProvideUnitCounterProvider provides UnitCounter providers.
func ProvideUnitCounterProvider() UnitCounterProvider {
	return ProvideUnitCounter
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
