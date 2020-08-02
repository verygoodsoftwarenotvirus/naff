package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_authServiceTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := authServiceTestDotGo(proj)

		expected := `
package example

import (
	v2 "github.com/alexedwards/scs/v2"
	memstore "github.com/alexedwards/scs/v2/memstore"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"testing"
)

func buildTestService(t *testing.T) *Service {
	t.Helper()

	logger := noop.ProvideNoopLogger()
	cfg := config.AuthSettings{
		CookieSecret: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
	}
	auth := &mock.Authenticator{}
	userDB := &mock1.UserDataManager{}
	oauth := &mockOAuth2ClientValidator{}
	ed := encoding.ProvideResponseEncoder()

	sm := v2.New()
	// this is currently the default, but in case that changes
	sm.Store = memstore.New()

	service, err := ProvideAuthService(
		logger,
		cfg,
		auth,
		userDB,
		oauth,
		sm,
		ed,
	)
	require.NoError(t, err)

	return service
}

func TestProvideAuthService(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		cfg := config.AuthSettings{
			CookieSecret: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
		}
		auth := &mock.Authenticator{}
		userDB := &mock1.UserDataManager{}
		oauth := &mockOAuth2ClientValidator{}
		ed := encoding.ProvideResponseEncoder()
		sm := v2.New()

		service, err := ProvideAuthService(
			noop.ProvideNoopLogger(),
			cfg,
			auth,
			userDB,
			oauth,
			sm,
			ed,
		)
		assert.NotNil(t, service)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildTestService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildTestService(proj)

		expected := `
package example

import (
	v2 "github.com/alexedwards/scs/v2"
	memstore "github.com/alexedwards/scs/v2/memstore"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"testing"
)

func buildTestService(t *testing.T) *Service {
	t.Helper()

	logger := noop.ProvideNoopLogger()
	cfg := config.AuthSettings{
		CookieSecret: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
	}
	auth := &mock.Authenticator{}
	userDB := &mock1.UserDataManager{}
	oauth := &mockOAuth2ClientValidator{}
	ed := encoding.ProvideResponseEncoder()

	sm := v2.New()
	// this is currently the default, but in case that changes
	sm.Store = memstore.New()

	service, err := ProvideAuthService(
		logger,
		cfg,
		auth,
		userDB,
		oauth,
		sm,
		ed,
	)
	require.NoError(t, err)

	return service
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideAuthService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestProvideAuthService(proj)

		expected := `
package example

import (
	v2 "github.com/alexedwards/scs/v2"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"testing"
)

func TestProvideAuthService(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		cfg := config.AuthSettings{
			CookieSecret: "BLAHBLAHBLAHPRETENDTHISISSECRET!",
		}
		auth := &mock.Authenticator{}
		userDB := &mock1.UserDataManager{}
		oauth := &mockOAuth2ClientValidator{}
		ed := encoding.ProvideResponseEncoder()
		sm := v2.New()

		service, err := ProvideAuthService(
			noop.ProvideNoopLogger(),
			cfg,
			auth,
			userDB,
			oauth,
			sm,
			ed,
		)
		assert.NotNil(t, service)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
