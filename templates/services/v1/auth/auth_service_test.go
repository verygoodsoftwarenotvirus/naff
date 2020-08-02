package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_authServiceDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := authServiceDotGo(proj)

		expected := `
package example

import (
	"context"
	v2 "github.com/alexedwards/scs/v2"
	securecookie "github.com/gorilla/securecookie"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

const (
	serviceName = "auth_service"
)

type (
	// OAuth2ClientValidator is a stand-in interface, where we needed to abstract
	// a regular structure with an interface for testing purposes.
	OAuth2ClientValidator interface {
		ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error)
	}

	// cookieEncoderDecoder is a stand-in interface for gorilla/securecookie
	cookieEncoderDecoder interface {
		Encode(name string, value interface{}) (string, error)
		Decode(name, value string, dst interface{}) error
	}

	// Service handles authentication service-wide
	Service struct {
		config               config.AuthSettings
		logger               v11.Logger
		authenticator        auth.Authenticator
		userDB               v1.UserDataManager
		oauth2ClientsService OAuth2ClientValidator
		encoderDecoder       encoding.EncoderDecoder
		cookieManager        cookieEncoderDecoder
		sessionManager       *v2.SessionManager
	}
)

// ProvideAuthService builds a new AuthService.
func ProvideAuthService(
	logger v11.Logger,
	cfg config.AuthSettings,
	authenticator auth.Authenticator,
	database v1.UserDataManager,
	oauth2ClientsService OAuth2ClientValidator,
	sessionManager *v2.SessionManager,
	encoder encoding.EncoderDecoder,
) (*Service, error) {
	svc := &Service{
		logger:               logger.WithName(serviceName),
		encoderDecoder:       encoder,
		config:               cfg,
		userDB:               database,
		oauth2ClientsService: oauth2ClientsService,
		authenticator:        authenticator,
		sessionManager:       sessionManager,
		cookieManager: securecookie.New(
			securecookie.GenerateRandomKey(64),
			[]byte(cfg.CookieSecret),
		),
	}
	svc.sessionManager.Lifetime = cfg.CookieLifetime

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthServiceConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildAuthServiceConstantDefs()

		expected := `
package example

import ()

const (
	serviceName = "auth_service"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthServiceTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildAuthServiceTypeDefs(proj)

		expected := `
package example

import (
	"context"
	v2 "github.com/alexedwards/scs/v2"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

type (
	// OAuth2ClientValidator is a stand-in interface, where we needed to abstract
	// a regular structure with an interface for testing purposes.
	OAuth2ClientValidator interface {
		ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error)
	}

	// cookieEncoderDecoder is a stand-in interface for gorilla/securecookie
	cookieEncoderDecoder interface {
		Encode(name string, value interface{}) (string, error)
		Decode(name, value string, dst interface{}) error
	}

	// Service handles authentication service-wide
	Service struct {
		config               config.AuthSettings
		logger               v11.Logger
		authenticator        auth.Authenticator
		userDB               v1.UserDataManager
		oauth2ClientsService OAuth2ClientValidator
		encoderDecoder       encoding.EncoderDecoder
		cookieManager        cookieEncoderDecoder
		sessionManager       *v2.SessionManager
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideAuthService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildProvideAuthService(proj)

		expected := `
package example

import (
	v2 "github.com/alexedwards/scs/v2"
	securecookie "github.com/gorilla/securecookie"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideAuthService builds a new AuthService.
func ProvideAuthService(
	logger v1.Logger,
	cfg config.AuthSettings,
	authenticator auth.Authenticator,
	database v11.UserDataManager,
	oauth2ClientsService OAuth2ClientValidator,
	sessionManager *v2.SessionManager,
	encoder encoding.EncoderDecoder,
) (*Service, error) {
	svc := &Service{
		logger:               logger.WithName(serviceName),
		encoderDecoder:       encoder,
		config:               cfg,
		userDB:               database,
		oauth2ClientsService: oauth2ClientsService,
		authenticator:        authenticator,
		sessionManager:       sessionManager,
		cookieManager: securecookie.New(
			securecookie.GenerateRandomKey(64),
			[]byte(cfg.CookieSecret),
		),
	}
	svc.sessionManager.Lifetime = cfg.CookieLifetime

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
