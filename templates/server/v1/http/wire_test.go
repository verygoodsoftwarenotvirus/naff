package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

var (
	// Providers is our wire superset of providers this package offers.
	Providers = wire.NewSet(
		paramFetcherProviders,
		ProvideServer,
		ProvideNamespace,
		ProvideNewsmanTypeNameManipulationFunc,
	)
)

// ProvideNamespace provides a namespace.
func ProvideNamespace() metrics.Namespace {
	return serverNamespace
}

// ProvideNewsmanTypeNameManipulationFunc provides an WebhookIDFetcher.
func ProvideNewsmanTypeNameManipulationFunc() newsman.TypeNameManipulationFunc {
	return func(s string) string {
		return s
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWireVarDeclarations(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers is our wire superset of providers this package offers.
	Providers = wire.NewSet(
		paramFetcherProviders,
		ProvideServer,
		ProvideNamespace,
		ProvideNewsmanTypeNameManipulationFunc,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideNamespace(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideNamespace(proj)

		expected := `
package example

import (
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
)

// ProvideNamespace provides a namespace.
func ProvideNamespace() metrics.Namespace {
	return serverNamespace
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideNewsmanTypeNameManipulationFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideNewsmanTypeNameManipulationFunc()

		expected := `
package example

import (
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideNewsmanTypeNameManipulationFunc provides an WebhookIDFetcher.
func ProvideNewsmanTypeNameManipulationFunc() newsman.TypeNameManipulationFunc {
	return func(s string) string {
		return s
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
