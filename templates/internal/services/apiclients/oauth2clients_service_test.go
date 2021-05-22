package apiclients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsServiceDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsServiceDotGo(proj)

		expected := `
package example

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v12 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	oauth2v3 "gopkg.in/oauth2.v3"
	manage "gopkg.in/oauth2.v3/manage"
	server "gopkg.in/oauth2.v3/server"
	store "gopkg.in/oauth2.v3/store"
	"net/http"
)

func init() {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}

const (
	// creationMiddlewareCtxKey is a string alias for referring to OAuth2 client creation data.
	creationMiddlewareCtxKey v1.ContextKey = "create_oauth2_client"

	counterName        metrics.CounterName = "oauth2_clients"
	counterDescription string              = "number of oauth2 clients managed by the oauth2 client service"
	serviceName        string              = "oauth2_clients_service"
)

var (
	_ v1.OAuth2ClientDataServer = (*Service)(nil)
	_ oauth2v3.ClientStore      = (*clientStore)(nil)
)

type (
	oauth2Handler interface {
		SetAllowGetAccessRequest(bool)
		SetClientAuthorizedHandler(handler server.ClientAuthorizedHandler)
		SetClientScopeHandler(handler server.ClientScopeHandler)
		SetClientInfoHandler(handler server.ClientInfoHandler)
		SetUserAuthorizationHandler(handler server.UserAuthorizationHandler)
		SetAuthorizeScopeHandler(handler server.AuthorizeScopeHandler)
		SetResponseErrorHandler(handler server.ResponseErrorHandler)
		SetInternalErrorHandler(handler server.InternalErrorHandler)
		ValidationBearerToken(*http.Request) (oauth2v3.TokenInfo, error)
		HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error
		HandleTokenRequest(res http.ResponseWriter, req *http.Request) error
	}

	// ClientIDFetcher is a function for fetching client IDs out of requests.
	ClientIDFetcher func(req *http.Request) uint64

	// Service manages our OAuth2 clients via HTTP.
	Service struct {
		logger               v11.Logger
		database             v12.DataManager
		authenticator        auth.Authenticator
		encoderDecoder       encoding.EncoderDecoder
		urlClientIDExtractor func(req *http.Request) uint64
		oauth2Handler        oauth2Handler
		oauth2ClientCounter  metrics.UnitCounter
	}

	clientStore struct {
		database v12.DataManager
	}
)

func newClientStore(db v12.DataManager) *clientStore {
	cs := &clientStore{
		database: db,
	}
	return cs
}

// GetByID implements oauth2.ClientStorage
func (s *clientStore) GetByID(id string) (oauth2v3.ClientInfo, error) {
	client, err := s.database.GetOAuth2ClientByClientID(context.Background(), id)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid client")
	} else if err != nil {
		return nil, fmt.Errorf("querying for client: %w", err)
	}

	return client, nil
}

// ProvideOAuth2ClientsService builds a new OAuth2ClientsService.
func ProvideOAuth2ClientsService(
	logger v11.Logger,
	db v12.DataManager,
	authenticator auth.Authenticator,
	clientIDFetcher ClientIDFetcher,
	encoderDecoder encoding.EncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
) (*Service, error) {
	manager := manage.NewDefaultManager()
	clientStore := newClientStore(db)
	manager.MapClientStorage(clientStore)
	tokenStore, tokenStoreErr := store.NewMemoryTokenStore()
	manager.MustTokenStorage(tokenStore, tokenStoreErr)
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	oHandler := server.NewDefaultServer(manager)
	oHandler.SetAllowGetAccessRequest(true)

	svc := &Service{
		database:             db,
		logger:               logger.WithName(serviceName),
		encoderDecoder:       encoderDecoder,
		authenticator:        authenticator,
		urlClientIDExtractor: clientIDFetcher,
		oauth2Handler:        oHandler,
	}
	initializeOAuth2Handler(svc)

	var err error
	if svc.oauth2ClientCounter, err = counterProvider(counterName, counterDescription); err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	return svc, nil
}

// initializeOAuth2Handler.
func initializeOAuth2Handler(svc *Service) {
	svc.oauth2Handler.SetAllowGetAccessRequest(true)
	svc.oauth2Handler.SetClientAuthorizedHandler(svc.ClientAuthorizedHandler)
	svc.oauth2Handler.SetClientScopeHandler(svc.ClientScopeHandler)
	svc.oauth2Handler.SetClientInfoHandler(server.ClientFormHandler)
	svc.oauth2Handler.SetAuthorizeScopeHandler(svc.AuthorizeScopeHandler)
	svc.oauth2Handler.SetResponseErrorHandler(svc.OAuth2ResponseErrorHandler)
	svc.oauth2Handler.SetInternalErrorHandler(svc.OAuth2InternalErrorHandler)
	svc.oauth2Handler.SetUserAuthorizationHandler(svc.UserAuthorizationHandler)

	// this sad type cast is here because I have an arbitrary.
	// test-only interface for OAuth2 interactions.
	if x, ok := svc.oauth2Handler.(*server.Server); ok {
		x.Config.AllowedGrantTypes = []oauth2v3.GrantType{
			oauth2v3.ClientCredentials,
			// oauth2.AuthorizationCode,
			// oauth2.Refreshing,
			// oauth2.Implicit,
		}
	}
}

// HandleAuthorizeRequest is a simple wrapper around the internal server's HandleAuthorizeRequest.
func (s *Service) HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error {
	return s.oauth2Handler.HandleAuthorizeRequest(res, req)
}

// HandleTokenRequest is a simple wrapper around the internal server's HandleTokenRequest.
func (s *Service) HandleTokenRequest(res http.ResponseWriter, req *http.Request) error {
	return s.oauth2Handler.HandleTokenRequest(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServiceInit()

		expected := `
package example

import (
	"crypto/rand"
)

func init() {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceConstDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServiceConstDefs(proj)

		expected := `
package example

import (
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

const (
	// creationMiddlewareCtxKey is a string alias for referring to OAuth2 client creation data.
	creationMiddlewareCtxKey v1.ContextKey = "create_oauth2_client"

	counterName        metrics.CounterName = "oauth2_clients"
	counterDescription string              = "number of oauth2 clients managed by the oauth2 client service"
	serviceName        string              = "oauth2_clients_service"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServiceVarDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	oauth2v3 "gopkg.in/oauth2.v3"
)

var (
	_ v1.OAuth2ClientDataServer = (*Service)(nil)
	_ oauth2v3.ClientStore      = (*clientStore)(nil)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServiceTypeDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	oauth2v3 "gopkg.in/oauth2.v3"
	server "gopkg.in/oauth2.v3/server"
	"net/http"
)

type (
	oauth2Handler interface {
		SetAllowGetAccessRequest(bool)
		SetClientAuthorizedHandler(handler server.ClientAuthorizedHandler)
		SetClientScopeHandler(handler server.ClientScopeHandler)
		SetClientInfoHandler(handler server.ClientInfoHandler)
		SetUserAuthorizationHandler(handler server.UserAuthorizationHandler)
		SetAuthorizeScopeHandler(handler server.AuthorizeScopeHandler)
		SetResponseErrorHandler(handler server.ResponseErrorHandler)
		SetInternalErrorHandler(handler server.InternalErrorHandler)
		ValidationBearerToken(*http.Request) (oauth2v3.TokenInfo, error)
		HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error
		HandleTokenRequest(res http.ResponseWriter, req *http.Request) error
	}

	// ClientIDFetcher is a function for fetching client IDs out of requests.
	ClientIDFetcher func(req *http.Request) uint64

	// Service manages our OAuth2 clients via HTTP.
	Service struct {
		logger               v1.Logger
		database             v11.DataManager
		authenticator        auth.Authenticator
		encoderDecoder       encoding.EncoderDecoder
		urlClientIDExtractor func(req *http.Request) uint64
		oauth2Handler        oauth2Handler
		oauth2ClientCounter  metrics.UnitCounter
	}

	clientStore struct {
		database v11.DataManager
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceNewClientStore(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServiceNewClientStore(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

func newClientStore(db v1.DataManager) *clientStore {
	cs := &clientStore{
		database: db,
	}
	return cs
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceGetByID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServiceGetByID()

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	oauth2v3 "gopkg.in/oauth2.v3"
)

// GetByID implements oauth2.ClientStorage
func (s *clientStore) GetByID(id string) (oauth2v3.ClientInfo, error) {
	client, err := s.database.GetOAuth2ClientByClientID(context.Background(), id)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid client")
	} else if err != nil {
		return nil, fmt.Errorf("querying for client: %w", err)
	}

	return client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceProvideOAuth2ClientsService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServiceProvideOAuth2ClientsService(proj)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	manage "gopkg.in/oauth2.v3/manage"
	server "gopkg.in/oauth2.v3/server"
	store "gopkg.in/oauth2.v3/store"
)

// ProvideOAuth2ClientsService builds a new OAuth2ClientsService.
func ProvideOAuth2ClientsService(
	logger v1.Logger,
	db v11.DataManager,
	authenticator auth.Authenticator,
	clientIDFetcher ClientIDFetcher,
	encoderDecoder encoding.EncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
) (*Service, error) {
	manager := manage.NewDefaultManager()
	clientStore := newClientStore(db)
	manager.MapClientStorage(clientStore)
	tokenStore, tokenStoreErr := store.NewMemoryTokenStore()
	manager.MustTokenStorage(tokenStore, tokenStoreErr)
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	oHandler := server.NewDefaultServer(manager)
	oHandler.SetAllowGetAccessRequest(true)

	svc := &Service{
		database:             db,
		logger:               logger.WithName(serviceName),
		encoderDecoder:       encoderDecoder,
		authenticator:        authenticator,
		urlClientIDExtractor: clientIDFetcher,
		oauth2Handler:        oHandler,
	}
	initializeOAuth2Handler(svc)

	var err error
	if svc.oauth2ClientCounter, err = counterProvider(counterName, counterDescription); err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceInitializeOAuth2Handler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServiceInitializeOAuth2Handler()

		expected := `
package example

import (
	oauth2v3 "gopkg.in/oauth2.v3"
	server "gopkg.in/oauth2.v3/server"
)

// initializeOAuth2Handler.
func initializeOAuth2Handler(svc *Service) {
	svc.oauth2Handler.SetAllowGetAccessRequest(true)
	svc.oauth2Handler.SetClientAuthorizedHandler(svc.ClientAuthorizedHandler)
	svc.oauth2Handler.SetClientScopeHandler(svc.ClientScopeHandler)
	svc.oauth2Handler.SetClientInfoHandler(server.ClientFormHandler)
	svc.oauth2Handler.SetAuthorizeScopeHandler(svc.AuthorizeScopeHandler)
	svc.oauth2Handler.SetResponseErrorHandler(svc.OAuth2ResponseErrorHandler)
	svc.oauth2Handler.SetInternalErrorHandler(svc.OAuth2InternalErrorHandler)
	svc.oauth2Handler.SetUserAuthorizationHandler(svc.UserAuthorizationHandler)

	// this sad type cast is here because I have an arbitrary.
	// test-only interface for OAuth2 interactions.
	if x, ok := svc.oauth2Handler.(*server.Server); ok {
		x.Config.AllowedGrantTypes = []oauth2v3.GrantType{
			oauth2v3.ClientCredentials,
			// oauth2.AuthorizationCode,
			// oauth2.Refreshing,
			// oauth2.Implicit,
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceHandleAuthorizeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServiceHandleAuthorizeRequest()

		expected := `
package example

import (
	"net/http"
)

// HandleAuthorizeRequest is a simple wrapper around the internal server's HandleAuthorizeRequest.
func (s *Service) HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error {
	return s.oauth2Handler.HandleAuthorizeRequest(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceHandleTokenRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServiceHandleTokenRequest()

		expected := `
package example

import (
	"net/http"
)

// HandleTokenRequest is a simple wrapper around the internal server's HandleTokenRequest.
func (s *Service) HandleTokenRequest(res http.ResponseWriter, req *http.Request) error {
	return s.oauth2Handler.HandleTokenRequest(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
