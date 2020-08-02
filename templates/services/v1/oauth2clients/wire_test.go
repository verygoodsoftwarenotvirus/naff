package oauth2clients

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var (
	// Providers are what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvideOAuth2ClientsService,
		ProvideOAuth2ClientDataServer,
	)
)

// ProvideOAuth2ClientDataServer is an arbitrary function for dependency injection's sake.
func ProvideOAuth2ClientDataServer(s *Service) v1.OAuth2ClientDataServer {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildWireVarDefs()

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers are what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvideOAuth2ClientsService,
		ProvideOAuth2ClientDataServer,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideOAuth2ClientDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildProvideOAuth2ClientDataServer(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideOAuth2ClientDataServer is an arbitrary function for dependency injection's sake.
func ProvideOAuth2ClientDataServer(s *Service) v1.OAuth2ClientDataServer {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
