package oauth2clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_middlewareDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := middlewareDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strings"
)

const (
	scopesSeparator = ","
	apiPathPrefix   = "/api/v1/"
)

// CreationInputMiddleware is a middleware for attaching OAuth2 client info to a request.
func (s *Service) CreationInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreationInputMiddleware")
		defer span.End()
		x := new(v1.OAuth2ClientCreationInput)

		// decode value from request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, creationMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// ExtractOAuth2ClientFromRequest extracts OAuth2 client data from a request.
func (s *Service) ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "ExtractOAuth2ClientFromRequest")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// validate bearer token.
	token, err := s.oauth2Handler.ValidationBearerToken(req)
	if err != nil {
		return nil, fmt.Errorf("validating bearer token: %w", err)
	}

	// fetch client ID.
	clientID := token.GetClientID()
	logger = logger.WithValue("client_id", clientID)

	// fetch client by client ID.
	c, err := s.database.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "error fetching OAuth2 Client")
		return nil, err
	}

	// determine the scope.
	scope := determineScope(req)
	hasScope := c.HasScope(scope)
	logger = logger.WithValue("scope", scope).WithValue("scopes", strings.Join(c.Scopes, scopesSeparator))

	if !hasScope {
		logger.Info("rejecting client for invalid scope")
		return nil, errors.New("client not authorized for scope")
	}

	return c, nil
}

// determineScope determines the scope of a request by its URL.
func determineScope(req *http.Request) string {
	if strings.HasPrefix(req.URL.Path, apiPathPrefix) {
		x := strings.TrimPrefix(req.URL.Path, apiPathPrefix)
		if y := strings.Split(x, "/"); len(y) > 0 {
			x = y[0]
		}
		return x
	}

	return ""
}

// OAuth2TokenAuthenticationMiddleware authenticates Oauth tokens.
func (s *Service) OAuth2TokenAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "OAuth2TokenAuthenticationMiddleware")
		defer span.End()

		c, err := s.ExtractOAuth2ClientFromRequest(ctx, req)
		if err != nil {
			s.logger.Error(err, "error authenticated token-authed request")
			http.Error(res, "invalid token", http.StatusUnauthorized)
			return
		}

		tracing.AttachUserIDToSpan(span, c.BelongsToUser)
		tracing.AttachOAuth2ClientIDToSpan(span, c.ClientID)
		tracing.AttachOAuth2ClientDatabaseIDToSpan(span, c.ID)

		// attach the client object to the request.
		ctx = context.WithValue(ctx, v1.OAuth2ClientKey, c)

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// OAuth2ClientInfoMiddleware fetches clientOAuth2Client info from requests and attaches it explicitly to a request.
func (s *Service) OAuth2ClientInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "OAuth2ClientInfoMiddleware")
		defer span.End()

		if v := req.URL.Query().Get(oauth2ClientIDURIParamKey); v != "" {
			logger := s.logger.WithValue("oauth2_client_id", v)

			client, err := s.database.GetOAuth2ClientByClientID(ctx, v)
			if err != nil {
				logger.Error(err, "error fetching OAuth2 client")
				http.Error(res, "invalid request", http.StatusUnauthorized)
				return
			}

			tracing.AttachUserIDToSpan(span, client.BelongsToUser)
			tracing.AttachOAuth2ClientIDToSpan(span, client.ClientID)
			tracing.AttachOAuth2ClientDatabaseIDToSpan(span, client.ID)

			ctx = context.WithValue(ctx, v1.OAuth2ClientKey, client)

			req = req.WithContext(ctx)
		}

		next.ServeHTTP(res, req)
	})
}

func (s *Service) fetchOAuth2ClientFromRequest(req *http.Request) *v1.OAuth2Client {
	client, ok := req.Context().Value(v1.OAuth2ClientKey).(*v1.OAuth2Client)
	_ = ok // we don't really care, but the linters do
	return client
}

func (s *Service) fetchOAuth2ClientIDFromRequest(req *http.Request) string {
	clientID, ok := req.Context().Value(clientIDKey).(string)
	_ = ok // we don't really care, but the linters do
	return clientID
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMiddlewareConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMiddlewareConstantDefs()

		expected := `
package example

import ()

const (
	scopesSeparator = ","
	apiPathPrefix   = "/api/v1/"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildCreationInputMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// CreationInputMiddleware is a middleware for attaching OAuth2 client info to a request.
func (s *Service) CreationInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CreationInputMiddleware")
		defer span.End()
		x := new(v1.OAuth2ClientCreationInput)

		// decode value from request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, creationMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExtractOAuth2ClientFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildExtractOAuth2ClientFromRequest(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strings"
)

// ExtractOAuth2ClientFromRequest extracts OAuth2 client data from a request.
func (s *Service) ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "ExtractOAuth2ClientFromRequest")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// validate bearer token.
	token, err := s.oauth2Handler.ValidationBearerToken(req)
	if err != nil {
		return nil, fmt.Errorf("validating bearer token: %w", err)
	}

	// fetch client ID.
	clientID := token.GetClientID()
	logger = logger.WithValue("client_id", clientID)

	// fetch client by client ID.
	c, err := s.database.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "error fetching OAuth2 Client")
		return nil, err
	}

	// determine the scope.
	scope := determineScope(req)
	hasScope := c.HasScope(scope)
	logger = logger.WithValue("scope", scope).WithValue("scopes", strings.Join(c.Scopes, scopesSeparator))

	if !hasScope {
		logger.Info("rejecting client for invalid scope")
		return nil, errors.New("client not authorized for scope")
	}

	return c, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDetermineScope(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildDetermineScope()

		expected := `
package example

import (
	"net/http"
	"strings"
)

// determineScope determines the scope of a request by its URL.
func determineScope(req *http.Request) string {
	if strings.HasPrefix(req.URL.Path, apiPathPrefix) {
		x := strings.TrimPrefix(req.URL.Path, apiPathPrefix)
		if y := strings.Split(x, "/"); len(y) > 0 {
			x = y[0]
		}
		return x
	}

	return ""
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2TokenAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2TokenAuthenticationMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// OAuth2TokenAuthenticationMiddleware authenticates Oauth tokens.
func (s *Service) OAuth2TokenAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "OAuth2TokenAuthenticationMiddleware")
		defer span.End()

		c, err := s.ExtractOAuth2ClientFromRequest(ctx, req)
		if err != nil {
			s.logger.Error(err, "error authenticated token-authed request")
			http.Error(res, "invalid token", http.StatusUnauthorized)
			return
		}

		tracing.AttachUserIDToSpan(span, c.BelongsToUser)
		tracing.AttachOAuth2ClientIDToSpan(span, c.ClientID)
		tracing.AttachOAuth2ClientDatabaseIDToSpan(span, c.ID)

		// attach the client object to the request.
		ctx = context.WithValue(ctx, v1.OAuth2ClientKey, c)

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientInfoMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientInfoMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// OAuth2ClientInfoMiddleware fetches clientOAuth2Client info from requests and attaches it explicitly to a request.
func (s *Service) OAuth2ClientInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "OAuth2ClientInfoMiddleware")
		defer span.End()

		if v := req.URL.Query().Get(oauth2ClientIDURIParamKey); v != "" {
			logger := s.logger.WithValue("oauth2_client_id", v)

			client, err := s.database.GetOAuth2ClientByClientID(ctx, v)
			if err != nil {
				logger.Error(err, "error fetching OAuth2 client")
				http.Error(res, "invalid request", http.StatusUnauthorized)
				return
			}

			tracing.AttachUserIDToSpan(span, client.BelongsToUser)
			tracing.AttachOAuth2ClientIDToSpan(span, client.ClientID)
			tracing.AttachOAuth2ClientDatabaseIDToSpan(span, client.ID)

			ctx = context.WithValue(ctx, v1.OAuth2ClientKey, client)

			req = req.WithContext(ctx)
		}

		next.ServeHTTP(res, req)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceFetchOAuth2ClientFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildServiceFetchOAuth2ClientFromRequest(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

func (s *Service) fetchOAuth2ClientFromRequest(req *http.Request) *v1.OAuth2Client {
	client, ok := req.Context().Value(v1.OAuth2ClientKey).(*v1.OAuth2Client)
	_ = ok // we don't really care, but the linters do
	return client
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceFetchOAuth2ClientIDFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildServiceFetchOAuth2ClientIDFromRequest()

		expected := `
package example

import (
	"net/http"
)

func (s *Service) fetchOAuth2ClientIDFromRequest(req *http.Request) string {
	clientID, ok := req.Context().Value(clientIDKey).(string)
	_ = ok // we don't really care, but the linters do
	return clientID
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
