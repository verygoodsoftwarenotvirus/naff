package auth

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
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideAuthService,
		ProvideWebsocketAuthFunc,
		ProvideOAuth2ClientValidator,
	)
)

// ProvideWebsocketAuthFunc provides a WebsocketAuthFunc.
func ProvideWebsocketAuthFunc(svc *Service) newsman.WebsocketAuthFunc {
	return svc.WebsocketAuthFunction
}

// ProvideOAuth2ClientValidator converts an oauth2clients.Service to an OAuth2ClientValidator
func ProvideOAuth2ClientValidator(s *oauth2clients.Service) OAuth2ClientValidator {
	return s
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
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideAuthService,
		ProvideWebsocketAuthFunc,
		ProvideOAuth2ClientValidator,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProvideWebsocketAuthFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildWireProvideWebsocketAuthFunc()

		expected := `
package example

import (
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideWebsocketAuthFunc provides a WebsocketAuthFunc.
func ProvideWebsocketAuthFunc(svc *Service) newsman.WebsocketAuthFunc {
	return svc.WebsocketAuthFunction
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProvideOAuth2ClientValidator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWireProvideOAuth2ClientValidator(proj)

		expected := `
package example

import (
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
)

// ProvideOAuth2ClientValidator converts an oauth2clients.Service to an OAuth2ClientValidator
func ProvideOAuth2ClientValidator(s *oauth2clients.Service) OAuth2ClientValidator {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
