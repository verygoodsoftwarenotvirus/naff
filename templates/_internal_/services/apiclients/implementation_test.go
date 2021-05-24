package apiclients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_implementationDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := implementationDotGo(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	errors1 "errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/observability/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	oauth2v3 "gopkg.in/oauth2.v3"
	errors "gopkg.in/oauth2.v3/errors"
	server "gopkg.in/oauth2.v3/server"
	"net/http"
	"strconv"
)

// gopkg.in/oauth2.v3/server specific implementations

var _ server.InternalErrorHandler = (*Service)(nil).OAuth2InternalErrorHandler

// OAuth2InternalErrorHandler fulfills a role for the OAuth2 server-side provider
func (s *Service) OAuth2InternalErrorHandler(err error) *errors.Response {
	s.logger.Error(err, "OAuth2 Internal Error")

	res := &errors.Response{
		Error:       err,
		Description: "Internal error",
		ErrorCode:   http.StatusInternalServerError,
		StatusCode:  http.StatusInternalServerError,
	}

	return res
}

var _ server.ResponseErrorHandler = (*Service)(nil).OAuth2ResponseErrorHandler

// OAuth2ResponseErrorHandler fulfills a role for the OAuth2 server-side provider
func (s *Service) OAuth2ResponseErrorHandler(re *errors.Response) {
	s.logger.WithValues(map[string]interface{}{
		"error_code":  re.ErrorCode,
		"description": re.Description,
		"uri":         re.URI,
		"status_code": re.StatusCode,
		"header":      re.Header,
	}).Error(re.Error, "OAuth2ResponseErrorHandler")
}

var _ server.AuthorizeScopeHandler = (*Service)(nil).AuthorizeScopeHandler

// AuthorizeScopeHandler satisfies the oauth2server AuthorizeScopeHandler interface.
func (s *Service) AuthorizeScopeHandler(res http.ResponseWriter, req *http.Request) (scope string, err error) {
	ctx, span := tracing.StartSpan(req.Context(), "AuthorizeScopeHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	scope = determineScope(req)
	logger = logger.WithValue("scope", scope)

	// check for client and return if valid.
	var client = s.fetchOAuth2ClientFromRequest(req)
	if client != nil && client.HasScope(scope) {
		res.WriteHeader(http.StatusOK)
		return scope, nil
	}

	// check to see if the client ID is present instead.
	if clientID := s.fetchOAuth2ClientIDFromRequest(req); clientID != "" {
		// fetch oauth2 client from database.
		client, err = s.database.GetOAuth2ClientByClientID(ctx, clientID)

		if err == sql.ErrNoRows {
			logger.Error(err, "error fetching OAuth2 Client")
			res.WriteHeader(http.StatusNotFound)
			return "", err
		} else if err != nil {
			logger.Error(err, "error fetching OAuth2 Client")
			res.WriteHeader(http.StatusInternalServerError)
			return "", err
		}

		// authorization check.
		if !client.HasScope(scope) {
			res.WriteHeader(http.StatusUnauthorized)
			return "", errors1.New("not authorized for scope")
		}

		return scope, nil
	}

	// invalid credentials.
	res.WriteHeader(http.StatusBadRequest)
	return "", errors1.New("no scope information found")
}

var _ server.UserAuthorizationHandler = (*Service)(nil).UserAuthorizationHandler

// UserAuthorizationHandler satisfies the oauth2server UserAuthorizationHandler interface.
func (s *Service) UserAuthorizationHandler(_ http.ResponseWriter, req *http.Request) (userID string, err error) {
	ctx, span := tracing.StartSpan(req.Context(), "UserAuthorizationHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)
	var uid uint64

	// check context for client.
	if client, clientOk := ctx.Value(v1.OAuth2ClientKey).(*v1.OAuth2Client); !clientOk {
		// check for user instead.
		si, userOk := ctx.Value(v1.SessionInfoKey).(*v1.SessionInfo)
		if !userOk || si == nil {
			logger.Debug("no user attached to this request")
			return "", errors1.New("user not found")
		}
		uid = si.UserID
	} else {
		uid = client.BelongsToUser
	}

	return strconv.FormatUint(uid, 10), nil
}

var _ server.ClientAuthorizedHandler = (*Service)(nil).ClientAuthorizedHandler

// ClientAuthorizedHandler satisfies the oauth2server ClientAuthorizedHandler interface.
func (s *Service) ClientAuthorizedHandler(clientID string, grant oauth2v3.GrantType) (allowed bool, err error) {
	// NOTE: it's a shame the interface we're implementing doesn't have this as its first argument
	ctx, span := tracing.StartSpan(context.Background(), "ClientAuthorizedHandler")
	defer span.End()

	logger := s.logger.WithValues(map[string]interface{}{
		"grant":     grant,
		"client_id": clientID,
	})

	// reject invalid grant type.
	if grant == oauth2v3.PasswordCredentials {
		return false, errors1.New("invalid grant type: password")
	}

	// fetch client data.
	client, err := s.database.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "fetching oauth2 client from database")
		return false, fmt.Errorf("fetching oauth2 client from database: %w", err)
	}

	// disallow implicit grants unless authorized.
	if grant == oauth2v3.Implicit && !client.ImplicitAllowed {
		return false, errors1.New("client not authorized for implicit grants")
	}

	return true, nil
}

var _ server.ClientScopeHandler = (*Service)(nil).ClientScopeHandler

// ClientScopeHandler satisfies the oauth2server ClientScopeHandler interface.
func (s *Service) ClientScopeHandler(clientID, scope string) (authed bool, err error) {
	// NOTE: it's a shame the interface we're implementing doesn't have this as its first argument
	ctx, span := tracing.StartSpan(context.Background(), "UserAuthorizationHandler")
	defer span.End()

	logger := s.logger.WithValues(map[string]interface{}{
		"client_id": clientID,
		"scope":     scope,
	})

	// fetch client info.
	c, err := s.database.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "error fetching OAuth2 client for ClientScopeHandler")
		return false, err
	}

	// check for scope.
	if c.HasScope(scope) {
		return true, nil
	}

	return false, errors1.New("unauthorized")
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildImplementationOAuth2InternalErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildImplementationOAuth2InternalErrorHandler()

		expected := `
package example

import (
	errors "gopkg.in/oauth2.v3/errors"
	server "gopkg.in/oauth2.v3/server"
	"net/http"
)

var _ server.InternalErrorHandler = (*Service)(nil).OAuth2InternalErrorHandler

// OAuth2InternalErrorHandler fulfills a role for the OAuth2 server-side provider
func (s *Service) OAuth2InternalErrorHandler(err error) *errors.Response {
	s.logger.Error(err, "OAuth2 Internal Error")

	res := &errors.Response{
		Error:       err,
		Description: "Internal error",
		ErrorCode:   http.StatusInternalServerError,
		StatusCode:  http.StatusInternalServerError,
	}

	return res
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildImplementationOAuth2ResponseErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildImplementationOAuth2ResponseErrorHandler()

		expected := `
package example

import (
	errors "gopkg.in/oauth2.v3/errors"
	server "gopkg.in/oauth2.v3/server"
)

var _ server.ResponseErrorHandler = (*Service)(nil).OAuth2ResponseErrorHandler

// OAuth2ResponseErrorHandler fulfills a role for the OAuth2 server-side provider
func (s *Service) OAuth2ResponseErrorHandler(re *errors.Response) {
	s.logger.WithValues(map[string]interface{}{
		"error_code":  re.ErrorCode,
		"description": re.Description,
		"uri":         re.URI,
		"status_code": re.StatusCode,
		"header":      re.Header,
	}).Error(re.Error, "OAuth2ResponseErrorHandler")
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildImplementationAuthorizeScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildImplementationAuthorizeScopeHandler(proj)

		expected := `
package example

import (
	"database/sql"
	"errors"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/observability/tracing"
	server "gopkg.in/oauth2.v3/server"
	"net/http"
)

var _ server.AuthorizeScopeHandler = (*Service)(nil).AuthorizeScopeHandler

// AuthorizeScopeHandler satisfies the oauth2server AuthorizeScopeHandler interface.
func (s *Service) AuthorizeScopeHandler(res http.ResponseWriter, req *http.Request) (scope string, err error) {
	ctx, span := tracing.StartSpan(req.Context(), "AuthorizeScopeHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	scope = determineScope(req)
	logger = logger.WithValue("scope", scope)

	// check for client and return if valid.
	var client = s.fetchOAuth2ClientFromRequest(req)
	if client != nil && client.HasScope(scope) {
		res.WriteHeader(http.StatusOK)
		return scope, nil
	}

	// check to see if the client ID is present instead.
	if clientID := s.fetchOAuth2ClientIDFromRequest(req); clientID != "" {
		// fetch oauth2 client from database.
		client, err = s.database.GetOAuth2ClientByClientID(ctx, clientID)

		if err == sql.ErrNoRows {
			logger.Error(err, "error fetching OAuth2 Client")
			res.WriteHeader(http.StatusNotFound)
			return "", err
		} else if err != nil {
			logger.Error(err, "error fetching OAuth2 Client")
			res.WriteHeader(http.StatusInternalServerError)
			return "", err
		}

		// authorization check.
		if !client.HasScope(scope) {
			res.WriteHeader(http.StatusUnauthorized)
			return "", errors.New("not authorized for scope")
		}

		return scope, nil
	}

	// invalid credentials.
	res.WriteHeader(http.StatusBadRequest)
	return "", errors.New("no scope information found")
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildImplementationUserAuthorizationHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildImplementationUserAuthorizationHandler(proj)

		expected := `
package example

import (
	"errors"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/observability/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	server "gopkg.in/oauth2.v3/server"
	"net/http"
	"strconv"
)

var _ server.UserAuthorizationHandler = (*Service)(nil).UserAuthorizationHandler

// UserAuthorizationHandler satisfies the oauth2server UserAuthorizationHandler interface.
func (s *Service) UserAuthorizationHandler(_ http.ResponseWriter, req *http.Request) (userID string, err error) {
	ctx, span := tracing.StartSpan(req.Context(), "UserAuthorizationHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)
	var uid uint64

	// check context for client.
	if client, clientOk := ctx.Value(v1.OAuth2ClientKey).(*v1.OAuth2Client); !clientOk {
		// check for user instead.
		si, userOk := ctx.Value(v1.SessionInfoKey).(*v1.SessionInfo)
		if !userOk || si == nil {
			logger.Debug("no user attached to this request")
			return "", errors.New("user not found")
		}
		uid = si.UserID
	} else {
		uid = client.BelongsToUser
	}

	return strconv.FormatUint(uid, 10), nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildImplementationClientAuthorizedHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildImplementationClientAuthorizedHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/observability/tracing"
	oauth2v3 "gopkg.in/oauth2.v3"
	server "gopkg.in/oauth2.v3/server"
)

var _ server.ClientAuthorizedHandler = (*Service)(nil).ClientAuthorizedHandler

// ClientAuthorizedHandler satisfies the oauth2server ClientAuthorizedHandler interface.
func (s *Service) ClientAuthorizedHandler(clientID string, grant oauth2v3.GrantType) (allowed bool, err error) {
	// NOTE: it's a shame the interface we're implementing doesn't have this as its first argument
	ctx, span := tracing.StartSpan(context.Background(), "ClientAuthorizedHandler")
	defer span.End()

	logger := s.logger.WithValues(map[string]interface{}{
		"grant":     grant,
		"client_id": clientID,
	})

	// reject invalid grant type.
	if grant == oauth2v3.PasswordCredentials {
		return false, errors.New("invalid grant type: password")
	}

	// fetch client data.
	client, err := s.database.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "fetching oauth2 client from database")
		return false, fmt.Errorf("fetching oauth2 client from database: %w", err)
	}

	// disallow implicit grants unless authorized.
	if grant == oauth2v3.Implicit && !client.ImplicitAllowed {
		return false, errors.New("client not authorized for implicit grants")
	}

	return true, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildImplementationClientScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildImplementationClientScopeHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/observability/tracing"
	server "gopkg.in/oauth2.v3/server"
)

var _ server.ClientScopeHandler = (*Service)(nil).ClientScopeHandler

// ClientScopeHandler satisfies the oauth2server ClientScopeHandler interface.
func (s *Service) ClientScopeHandler(clientID, scope string) (authed bool, err error) {
	// NOTE: it's a shame the interface we're implementing doesn't have this as its first argument
	ctx, span := tracing.StartSpan(context.Background(), "UserAuthorizationHandler")
	defer span.End()

	logger := s.logger.WithValues(map[string]interface{}{
		"client_id": clientID,
		"scope":     scope,
	})

	// fetch client info.
	c, err := s.database.GetOAuth2ClientByClientID(ctx, clientID)
	if err != nil {
		logger.Error(err, "error fetching OAuth2 client for ClientScopeHandler")
		return false, err
	}

	// check for scope.
	if c.HasScope(scope) {
		return true, nil
	}

	return false, errors.New("unauthorized")
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
