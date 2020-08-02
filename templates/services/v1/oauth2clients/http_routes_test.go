package oauth2clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := httpRoutesDotGo(proj)

		expected := `
package example

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strings"
)

const (
	// URIParamKey is used for referring to OAuth2 client IDs in router params.
	URIParamKey = "oauth2ClientID"

	oauth2ClientIDURIParamKey               = "client_id"
	clientIDKey               v1.ContextKey = "client_id"
)

// randString produces a random string.
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func randString() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	// this is so that we don't end up with `+"`"+`=`+"`"+` in IDs
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

// fetchUserID grabs a userID out of the request context.
func (s *Service) fetchUserID(req *http.Request) uint64 {
	if si, ok := req.Context().Value(v1.SessionInfoKey).(*v1.SessionInfo); ok && si != nil {
		return si.UserID
	}
	return 0
}

// ListHandler is a handler that returns a list of OAuth2 clients.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// extract filter.
	filter := v1.ExtractQueryFilter(req)

	// determine user.
	userID := s.fetchUserID(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// fetch oauth2 clients.
	oauth2Clients, err := s.database.GetOAuth2ClientsForUser(ctx, userID, filter)
	if err == sql.ErrNoRows {
		// just return an empty list if there are no results.
		oauth2Clients = &v1.OAuth2ClientList{
			Clients: []v1.OAuth2Client{},
		}
	} else if err != nil {
		logger.Error(err, "encountered error getting list of oauth2 clients from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, oauth2Clients); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CreateHandler is our OAuth2 client creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// fetch creation input from request context.
	input, ok := ctx.Value(creationMiddlewareCtxKey).(*v1.OAuth2ClientCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// set some data.
	input.ClientID, input.ClientSecret = randString(), randString()
	input.BelongsToUser = s.fetchUserID(req)

	// keep relevant data in mind.
	logger = logger.WithValues(map[string]interface{}{
		"username":     input.Username,
		"scopes":       strings.Join(input.Scopes, scopesSeparator),
		"redirect_uri": input.RedirectURI,
	})

	// retrieve user.
	user, err := s.database.GetUserByUsername(ctx, input.Username)
	if err != nil {
		logger.Error(err, "fetching user by username")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// tag span since we have the info.
	tracing.AttachUserIDToSpan(span, user.ID)

	// check credentials.
	valid, err := s.authenticator.ValidateLogin(
		ctx,
		user.HashedPassword,
		input.Password,
		user.TwoFactorSecret,
		input.TOTPToken,
		user.Salt,
	)

	if !valid {
		logger.Debug("invalid credentials provided")
		res.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		logger.Error(err, "validating user credentials")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the client.
	client, err := s.database.CreateOAuth2Client(ctx, input)
	if err != nil {
		logger.Error(err, "creating oauth2Client in the database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify interested parties.
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, client.ID)
	s.oauth2ClientCounter.Increment(ctx)

	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, client); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ReadHandler is a route handler for retrieving an OAuth2 client.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine subject of request.
	userID := s.fetchUserID(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant oauth2 client ID.
	oauth2ClientID := s.urlClientIDExtractor(req)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, oauth2ClientID)
	logger = logger.WithValue("oauth2_client_id", oauth2ClientID)

	// fetch oauth2 client.
	x, err := s.database.GetOAuth2Client(ctx, oauth2ClientID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("ReadHandler called on nonexistent client")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching oauth2Client from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ArchiveHandler is a route handler for archiving an OAuth2 client.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine subject of request.
	userID := s.fetchUserID(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant oauth2 client ID.
	oauth2ClientID := s.urlClientIDExtractor(req)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, oauth2ClientID)
	logger = logger.WithValue("oauth2_client_id", oauth2ClientID)

	// mark client as archived.
	err := s.database.ArchiveOAuth2Client(ctx, oauth2ClientID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "encountered error deleting oauth2 client")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.oauth2ClientCounter.Decrement(ctx)
	res.WriteHeader(http.StatusNoContent)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHTTPRoutesConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientHTTPRoutesConstantDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

const (
	// URIParamKey is used for referring to OAuth2 client IDs in router params.
	URIParamKey = "oauth2ClientID"

	oauth2ClientIDURIParamKey               = "client_id"
	clientIDKey               v1.ContextKey = "client_id"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHTTPRoutesRandString(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2ClientHTTPRoutesRandString()

		expected := `
package example

import (
	"crypto/rand"
	"encoding/base32"
)

// randString produces a random string.
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func randString() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	// this is so that we don't end up with `+"`"+`=`+"`"+` in IDs
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHTTPRoutesFetchUserID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientHTTPRoutesFetchUserID(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// fetchUserID grabs a userID out of the request context.
func (s *Service) fetchUserID(req *http.Request) uint64 {
	if si, ok := req.Context().Value(v1.SessionInfoKey).(*v1.SessionInfo); ok && si != nil {
		return si.UserID
	}
	return 0
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHTTPRoutesListHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientHTTPRoutesListHandler(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// ListHandler is a handler that returns a list of OAuth2 clients.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// extract filter.
	filter := v1.ExtractQueryFilter(req)

	// determine user.
	userID := s.fetchUserID(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// fetch oauth2 clients.
	oauth2Clients, err := s.database.GetOAuth2ClientsForUser(ctx, userID, filter)
	if err == sql.ErrNoRows {
		// just return an empty list if there are no results.
		oauth2Clients = &v1.OAuth2ClientList{
			Clients: []v1.OAuth2Client{},
		}
	} else if err != nil {
		logger.Error(err, "encountered error getting list of oauth2 clients from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, oauth2Clients); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHTTPRoutesCreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientHTTPRoutesCreateHandler(proj)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strings"
)

// CreateHandler is our OAuth2 client creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// fetch creation input from request context.
	input, ok := ctx.Value(creationMiddlewareCtxKey).(*v1.OAuth2ClientCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// set some data.
	input.ClientID, input.ClientSecret = randString(), randString()
	input.BelongsToUser = s.fetchUserID(req)

	// keep relevant data in mind.
	logger = logger.WithValues(map[string]interface{}{
		"username":     input.Username,
		"scopes":       strings.Join(input.Scopes, scopesSeparator),
		"redirect_uri": input.RedirectURI,
	})

	// retrieve user.
	user, err := s.database.GetUserByUsername(ctx, input.Username)
	if err != nil {
		logger.Error(err, "fetching user by username")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// tag span since we have the info.
	tracing.AttachUserIDToSpan(span, user.ID)

	// check credentials.
	valid, err := s.authenticator.ValidateLogin(
		ctx,
		user.HashedPassword,
		input.Password,
		user.TwoFactorSecret,
		input.TOTPToken,
		user.Salt,
	)

	if !valid {
		logger.Debug("invalid credentials provided")
		res.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		logger.Error(err, "validating user credentials")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the client.
	client, err := s.database.CreateOAuth2Client(ctx, input)
	if err != nil {
		logger.Error(err, "creating oauth2Client in the database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify interested parties.
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, client.ID)
	s.oauth2ClientCounter.Increment(ctx)

	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, client); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHTTPRoutesReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientHTTPRoutesReadHandler(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// ReadHandler is a route handler for retrieving an OAuth2 client.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine subject of request.
	userID := s.fetchUserID(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant oauth2 client ID.
	oauth2ClientID := s.urlClientIDExtractor(req)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, oauth2ClientID)
	logger = logger.WithValue("oauth2_client_id", oauth2ClientID)

	// fetch oauth2 client.
	x, err := s.database.GetOAuth2Client(ctx, oauth2ClientID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("ReadHandler called on nonexistent client")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error fetching oauth2Client from database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode response and peace.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHTTPRoutesArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientHTTPRoutesArchiveHandler(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// ArchiveHandler is a route handler for archiving an OAuth2 client.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ArchiveHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine subject of request.
	userID := s.fetchUserID(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant oauth2 client ID.
	oauth2ClientID := s.urlClientIDExtractor(req)
	tracing.AttachOAuth2ClientDatabaseIDToSpan(span, oauth2ClientID)
	logger = logger.WithValue("oauth2_client_id", oauth2ClientID)

	// mark client as archived.
	err := s.database.ArchiveOAuth2Client(ctx, oauth2ClientID, userID)
	if err == sql.ErrNoRows {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "encountered error deleting oauth2 client")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify relevant parties.
	s.oauth2ClientCounter.Decrement(ctx)
	res.WriteHeader(http.StatusNoContent)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
