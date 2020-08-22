package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := httpRoutesDotGo(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	securecookie "github.com/gorilla/securecookie"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"time"
)

const (
	// CookieName is the name of the cookie we attach to requests.
	CookieName         = "t_cookie"
	cookieErrorLogName = "_COOKIE_CONSTRUCTION_ERROR_"

	sessionInfoKey = "session_info"
)

// DecodeCookieFromRequest takes a request object and fetches the cookie data if it is present.
func (s *Service) DecodeCookieFromRequest(ctx context.Context, req *http.Request) (ca *v1.SessionInfo, err error) {
	ctx, span := tracing.StartSpan(ctx, "DecodeCookieFromRequest")
	defer span.End()

	logger := s.logger.WithRequest(req)

	cookie, err := req.Cookie(CookieName)
	if err != http.ErrNoCookie && cookie != nil {
		var token string
		decodeErr := s.cookieManager.Decode(CookieName, cookie.Value, &token)
		if decodeErr != nil {
			logger.Error(err, "decoding request cookie")
			return nil, fmt.Errorf("decoding request cookie: %w", decodeErr)
		}

		var sessionErr error
		ctx, sessionErr = s.sessionManager.Load(ctx, token)
		if sessionErr != nil {
			logger.Error(sessionErr, "error loading token")
			return nil, errors.New("error loading token")
		}

		si, ok := s.sessionManager.Get(ctx, sessionInfoKey).(*v1.SessionInfo)
		if !ok {
			errToReturn := errors.New("no session info attached to context")
			logger.Error(errToReturn, "fetching session data")
			return nil, errToReturn
		}

		return si, nil
	}

	return nil, http.ErrNoCookie
}

// WebsocketAuthFunction is provided to Newsman to determine if a user has access to websockets.
func (s *Service) WebsocketAuthFunction(req *http.Request) bool {
	ctx, span := tracing.StartSpan(req.Context(), "WebsocketAuthFunction")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// First we check to see if there is an OAuth2 token for a valid client attached to the request.
	// We do this first because it is presumed to be the primary means by which requests are made to the httpServer.
	oauth2Client, err := s.oauth2ClientsService.ExtractOAuth2ClientFromRequest(ctx, req)
	if err == nil && oauth2Client != nil {
		return true
	}

	// In the event there's not a valid OAuth2 token attached to the request, or there is some other OAuth2 issue,
	// we next check to see if a valid cookie is attached to the request.
	cookieAuth, cookieErr := s.DecodeCookieFromRequest(ctx, req)
	if cookieErr == nil && cookieAuth != nil {
		return true
	}

	// If your request gets here, you're likely either trying to get here, or desperately trying to get anywhere.
	logger.Error(err, "error authenticated token-authenticated request")
	return false
}

// fetchUserFromCookie takes a request object and fetches the cookie, and then the user for that cookie.
func (s *Service) fetchUserFromCookie(ctx context.Context, req *http.Request) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "fetchUserFromCookie")
	defer span.End()

	ca, decodeErr := s.DecodeCookieFromRequest(ctx, req)
	if decodeErr != nil {
		return nil, fmt.Errorf("fetching cookie data from request: %w", decodeErr)
	}

	user, userFetchErr := s.userDB.GetUser(req.Context(), ca.UserID)
	if userFetchErr != nil {
		return nil, fmt.Errorf("fetching user from request: %w", userFetchErr)
	}
	tracing.AttachUserIDToSpan(span, ca.UserID)

	return user, nil
}

// LoginHandler is our login route.
func (s *Service) LoginHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "LoginHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	loginData, errRes := s.fetchLoginDataFromRequest(req)
	if errRes != nil || loginData == nil {
		logger.Error(errRes, "error encountered fetching login data from request")
		res.WriteHeader(http.StatusUnauthorized)
		if err := s.encoderDecoder.EncodeResponse(res, errRes); err != nil {
			logger.Error(err, "encoding response")
		}
		return
	}

	tracing.AttachUserIDToSpan(span, loginData.user.ID)
	tracing.AttachUsernameToSpan(span, loginData.user.Username)

	logger = logger.WithValue("user", loginData.user.ID)
	loginValid, err := s.validateLogin(ctx, *loginData)
	if err != nil {
		logger.Error(err, "error encountered validating login")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	logger = logger.WithValue("valid", loginValid)

	if !loginValid {
		logger.Debug("login was invalid")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var sessionErr error
	ctx, sessionErr = s.sessionManager.Load(ctx, "")
	if sessionErr != nil {
		logger.Error(sessionErr, "error loading token")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if renewTokenErr := s.sessionManager.RenewToken(ctx); renewTokenErr != nil {
		logger.Error(err, "error encountered renewing token")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.sessionManager.Put(ctx, sessionInfoKey, loginData.user.ToSessionInfo())

	token, expiry, err := s.sessionManager.Commit(ctx)
	if err != nil {
		logger.Error(err, "error encountered writing to session store")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := s.buildCookie(token, expiry)
	if err != nil {
		logger.Error(err, "error encountered building cookie")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, cookie)
	res.WriteHeader(http.StatusNoContent)
}

// LogoutHandler is our logout route.
func (s *Service) LogoutHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "LogoutHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	ctx, sessionErr := s.sessionManager.Load(ctx, "")
	if sessionErr != nil {
		logger.Error(sessionErr, "error loading token")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.sessionManager.Clear(ctx); err != nil {
		logger.Error(err, "clearing user session")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cookie, cookieRetrievalErr := req.Cookie(CookieName); cookieRetrievalErr == nil && cookie != nil {
		if c, cookieBuildingErr := s.buildCookie("deleted", time.Time{}); cookieBuildingErr == nil && c != nil {
			c.MaxAge = -1
			http.SetCookie(res, c)
		} else {
			logger.Error(cookieBuildingErr, "error encountered building cookie")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		logger.WithError(cookieRetrievalErr).Debug("logout was called, no cookie was found")
	}

	res.WriteHeader(http.StatusOK)
}

// StatusHandler returns the user info for the user making the request.
func (s *Service) StatusHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "StatusHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	var sr *v1.StatusResponse
	userInfo, err := s.fetchUserFromCookie(ctx, req)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		sr = &v1.StatusResponse{
			Authenticated: false,
			IsAdmin:       false,
		}
	} else {
		sr = &v1.StatusResponse{
			Authenticated: true,
			IsAdmin:       userInfo.IsAdmin,
		}
	}

	if err := s.encoderDecoder.EncodeResponse(res, sr); err != nil {
		logger.Error(err, "encoding response")
	}
}

// CycleSecretHandler rotates the cookie building secret with a new random secret.
func (s *Service) CycleSecretHandler(res http.ResponseWriter, req *http.Request) {
	_, span := tracing.StartSpan(req.Context(), "CycleSecretHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)
	logger.Info("cycling cookie secret!")

	s.cookieManager = securecookie.New(
		securecookie.GenerateRandomKey(64),
		[]byte(s.config.CookieSecret),
	)

	res.WriteHeader(http.StatusCreated)
}

type loginData struct {
	loginInput *v1.UserLoginInput
	user       *v1.User
}

// fetchLoginDataFromRequest searches a given HTTP request for parsed login input data, and
// returns a helper struct with the relevant login information.
func (s *Service) fetchLoginDataFromRequest(req *http.Request) (*loginData, *v1.ErrorResponse) {
	ctx, span := tracing.StartSpan(req.Context(), "fetchLoginDataFromRequest")
	defer span.End()

	logger := s.logger.WithRequest(req)

	loginInput, ok := ctx.Value(userLoginInputMiddlewareCtxKey).(*v1.UserLoginInput)
	if !ok {
		logger.Debug("no UserLoginInput found for /login request")
		return nil, &v1.ErrorResponse{
			Code: http.StatusUnauthorized,
		}
	}

	username := loginInput.Username
	tracing.AttachUsernameToSpan(span, username)

	// you could ensure there isn't an unsatisfied password reset token
	// requested before allowing login here.

	user, err := s.userDB.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		logger.Error(err, "no matching user")
		return nil, &v1.ErrorResponse{Code: http.StatusBadRequest}
	} else if err != nil {
		logger.Error(err, "error fetching user")
		return nil, &v1.ErrorResponse{Code: http.StatusInternalServerError}
	}
	tracing.AttachUserIDToSpan(span, user.ID)

	ld := &loginData{
		loginInput: loginInput,
		user:       user,
	}

	return ld, nil
}

// validateLogin takes login information and returns whether or not the login is valid.
// In the event that there's an error, this function will return false and the error.
func (s *Service) validateLogin(ctx context.Context, loginInfo loginData) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "validateLogin")
	defer span.End()

	// alias the relevant data.
	user, loginInput := loginInfo.user, loginInfo.loginInput
	logger := s.logger.WithValue("username", user.Username)

	// check for login validity.
	loginValid, err := s.authenticator.ValidateLogin(
		ctx,
		user.HashedPassword,
		loginInput.Password,
		user.TwoFactorSecret,
		loginInput.TOTPToken,
		user.Salt,
	)

	// if the login is otherwise valid, but the password is too weak, try to rehash it.
	if err == auth.ErrCostTooLow && loginValid {
		logger.Debug("hashed password was deemed to weak, updating its hash")

		// re-hash the password
		updated, hashErr := s.authenticator.HashPassword(ctx, loginInput.Password)
		if hashErr != nil {
			return false, fmt.Errorf("updating password hash: %w", hashErr)
		}

		// update stored hashed password in the database.
		user.HashedPassword = updated
		if updateErr := s.userDB.UpdateUser(ctx, user); updateErr != nil {
			return false, fmt.Errorf("saving updated password hash: %w", updateErr)
		}

		return loginValid, nil
	} else if err != nil && err != auth.ErrCostTooLow {
		logger.Error(err, "issue validating login")
		return false, fmt.Errorf("validating login: %w", err)
	}

	return loginValid, err
}

// buildCookie provides a consistent way of constructing an HTTP cookie.
func (s *Service) buildCookie(value string, expiry time.Time) (*http.Cookie, error) {
	encoded, err := s.cookieManager.Encode(CookieName, value)
	if err != nil {
		// NOTE: these errors should be infrequent, and should cause alarm when they do occur
		s.logger.WithName(cookieErrorLogName).Error(err, "error encoding cookie")
		return nil, err
	}

	// https://www.calhoun.io/securing-cookies-in-go/
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.config.SecureCookiesOnly,
		Domain:   s.config.CookieDomain,
		Expires:  expiry,
	}

	return cookie, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildHTTPRoutesConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildHTTPRoutesConstantDefs(proj)

		expected := `
package example

import ()

const (
	// CookieName is the name of the cookie we attach to requests.
	CookieName         = "t_cookie"
	cookieErrorLogName = "_COOKIE_CONSTRUCTION_ERROR_"

	sessionInfoKey = "session_info"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDecodeCookieFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildDecodeCookieFromRequest(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// DecodeCookieFromRequest takes a request object and fetches the cookie data if it is present.
func (s *Service) DecodeCookieFromRequest(ctx context.Context, req *http.Request) (ca *v1.SessionInfo, err error) {
	ctx, span := tracing.StartSpan(ctx, "DecodeCookieFromRequest")
	defer span.End()

	logger := s.logger.WithRequest(req)

	cookie, err := req.Cookie(CookieName)
	if err != http.ErrNoCookie && cookie != nil {
		var token string
		decodeErr := s.cookieManager.Decode(CookieName, cookie.Value, &token)
		if decodeErr != nil {
			logger.Error(err, "decoding request cookie")
			return nil, fmt.Errorf("decoding request cookie: %w", decodeErr)
		}

		var sessionErr error
		ctx, sessionErr = s.sessionManager.Load(ctx, token)
		if sessionErr != nil {
			logger.Error(sessionErr, "error loading token")
			return nil, errors.New("error loading token")
		}

		si, ok := s.sessionManager.Get(ctx, sessionInfoKey).(*v1.SessionInfo)
		if !ok {
			errToReturn := errors.New("no session info attached to context")
			logger.Error(errToReturn, "fetching session data")
			return nil, errToReturn
		}

		return si, nil
	}

	return nil, http.ErrNoCookie
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebsocketAuthFunction(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWebsocketAuthFunction(proj)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// WebsocketAuthFunction is provided to Newsman to determine if a user has access to websockets.
func (s *Service) WebsocketAuthFunction(req *http.Request) bool {
	ctx, span := tracing.StartSpan(req.Context(), "WebsocketAuthFunction")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// First we check to see if there is an OAuth2 token for a valid client attached to the request.
	// We do this first because it is presumed to be the primary means by which requests are made to the httpServer.
	oauth2Client, err := s.oauth2ClientsService.ExtractOAuth2ClientFromRequest(ctx, req)
	if err == nil && oauth2Client != nil {
		return true
	}

	// In the event there's not a valid OAuth2 token attached to the request, or there is some other OAuth2 issue,
	// we next check to see if a valid cookie is attached to the request.
	cookieAuth, cookieErr := s.DecodeCookieFromRequest(ctx, req)
	if cookieErr == nil && cookieAuth != nil {
		return true
	}

	// If your request gets here, you're likely either trying to get here, or desperately trying to get anywhere.
	logger.Error(err, "error authenticated token-authenticated request")
	return false
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFetchUserFromCookie(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildFetchUserFromCookie(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// fetchUserFromCookie takes a request object and fetches the cookie, and then the user for that cookie.
func (s *Service) fetchUserFromCookie(ctx context.Context, req *http.Request) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "fetchUserFromCookie")
	defer span.End()

	ca, decodeErr := s.DecodeCookieFromRequest(ctx, req)
	if decodeErr != nil {
		return nil, fmt.Errorf("fetching cookie data from request: %w", decodeErr)
	}

	user, userFetchErr := s.userDB.GetUser(req.Context(), ca.UserID)
	if userFetchErr != nil {
		return nil, fmt.Errorf("fetching user from request: %w", userFetchErr)
	}
	tracing.AttachUserIDToSpan(span, ca.UserID)

	return user, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildLoginHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildLoginHandler(proj)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// LoginHandler is our login route.
func (s *Service) LoginHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "LoginHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	loginData, errRes := s.fetchLoginDataFromRequest(req)
	if errRes != nil || loginData == nil {
		logger.Error(errRes, "error encountered fetching login data from request")
		res.WriteHeader(http.StatusUnauthorized)
		if err := s.encoderDecoder.EncodeResponse(res, errRes); err != nil {
			logger.Error(err, "encoding response")
		}
		return
	}

	tracing.AttachUserIDToSpan(span, loginData.user.ID)
	tracing.AttachUsernameToSpan(span, loginData.user.Username)

	logger = logger.WithValue("user", loginData.user.ID)
	loginValid, err := s.validateLogin(ctx, *loginData)
	if err != nil {
		logger.Error(err, "error encountered validating login")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	logger = logger.WithValue("valid", loginValid)

	if !loginValid {
		logger.Debug("login was invalid")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var sessionErr error
	ctx, sessionErr = s.sessionManager.Load(ctx, "")
	if sessionErr != nil {
		logger.Error(sessionErr, "error loading token")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if renewTokenErr := s.sessionManager.RenewToken(ctx); renewTokenErr != nil {
		logger.Error(err, "error encountered renewing token")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.sessionManager.Put(ctx, sessionInfoKey, loginData.user.ToSessionInfo())

	token, expiry, err := s.sessionManager.Commit(ctx)
	if err != nil {
		logger.Error(err, "error encountered writing to session store")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := s.buildCookie(token, expiry)
	if err != nil {
		logger.Error(err, "error encountered building cookie")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, cookie)
	res.WriteHeader(http.StatusNoContent)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildLogoutHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildLogoutHandler(proj)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"time"
)

// LogoutHandler is our logout route.
func (s *Service) LogoutHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "LogoutHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	ctx, sessionErr := s.sessionManager.Load(ctx, "")
	if sessionErr != nil {
		logger.Error(sessionErr, "error loading token")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.sessionManager.Clear(ctx); err != nil {
		logger.Error(err, "clearing user session")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cookie, cookieRetrievalErr := req.Cookie(CookieName); cookieRetrievalErr == nil && cookie != nil {
		if c, cookieBuildingErr := s.buildCookie("deleted", time.Time{}); cookieBuildingErr == nil && c != nil {
			c.MaxAge = -1
			http.SetCookie(res, c)
		} else {
			logger.Error(cookieBuildingErr, "error encountered building cookie")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		logger.WithError(cookieRetrievalErr).Debug("logout was called, no cookie was found")
	}

	res.WriteHeader(http.StatusOK)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildStatusHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildStatusHandler(proj)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// StatusHandler returns the user info for the user making the request.
func (s *Service) StatusHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "StatusHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	var sr *v1.StatusResponse
	userInfo, err := s.fetchUserFromCookie(ctx, req)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		sr = &v1.StatusResponse{
			Authenticated: false,
			IsAdmin:       false,
		}
	} else {
		sr = &v1.StatusResponse{
			Authenticated: true,
			IsAdmin:       userInfo.IsAdmin,
		}
	}

	if err := s.encoderDecoder.EncodeResponse(res, sr); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCycleSecretHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCycleSecretHandler(proj)

		expected := `
package example

import (
	securecookie "github.com/gorilla/securecookie"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// CycleSecretHandler rotates the cookie building secret with a new random secret.
func (s *Service) CycleSecretHandler(res http.ResponseWriter, req *http.Request) {
	_, span := tracing.StartSpan(req.Context(), "CycleSecretHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)
	logger.Info("cycling cookie secret!")

	s.cookieManager = securecookie.New(
		securecookie.GenerateRandomKey(64),
		[]byte(s.config.CookieSecret),
	)

	res.WriteHeader(http.StatusCreated)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFetchLoginDataFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildFetchLoginDataFromRequest(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// fetchLoginDataFromRequest searches a given HTTP request for parsed login input data, and
// returns a helper struct with the relevant login information.
func (s *Service) fetchLoginDataFromRequest(req *http.Request) (*loginData, *v1.ErrorResponse) {
	ctx, span := tracing.StartSpan(req.Context(), "fetchLoginDataFromRequest")
	defer span.End()

	logger := s.logger.WithRequest(req)

	loginInput, ok := ctx.Value(userLoginInputMiddlewareCtxKey).(*v1.UserLoginInput)
	if !ok {
		logger.Debug("no UserLoginInput found for /login request")
		return nil, &v1.ErrorResponse{
			Code: http.StatusUnauthorized,
		}
	}

	username := loginInput.Username
	tracing.AttachUsernameToSpan(span, username)

	// you could ensure there isn't an unsatisfied password reset token
	// requested before allowing login here.

	user, err := s.userDB.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		logger.Error(err, "no matching user")
		return nil, &v1.ErrorResponse{Code: http.StatusBadRequest}
	} else if err != nil {
		logger.Error(err, "error fetching user")
		return nil, &v1.ErrorResponse{Code: http.StatusInternalServerError}
	}
	tracing.AttachUserIDToSpan(span, user.ID)

	ld := &loginData{
		loginInput: loginInput,
		user:       user,
	}

	return ld, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildValidateLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildValidateLogin(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// validateLogin takes login information and returns whether or not the login is valid.
// In the event that there's an error, this function will return false and the error.
func (s *Service) validateLogin(ctx context.Context, loginInfo loginData) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "validateLogin")
	defer span.End()

	// alias the relevant data.
	user, loginInput := loginInfo.user, loginInfo.loginInput
	logger := s.logger.WithValue("username", user.Username)

	// check for login validity.
	loginValid, err := s.authenticator.ValidateLogin(
		ctx,
		user.HashedPassword,
		loginInput.Password,
		user.TwoFactorSecret,
		loginInput.TOTPToken,
		user.Salt,
	)

	// if the login is otherwise valid, but the password is too weak, try to rehash it.
	if err == auth.ErrCostTooLow && loginValid {
		logger.Debug("hashed password was deemed to weak, updating its hash")

		// re-hash the password
		updated, hashErr := s.authenticator.HashPassword(ctx, loginInput.Password)
		if hashErr != nil {
			return false, fmt.Errorf("updating password hash: %w", hashErr)
		}

		// update stored hashed password in the database.
		user.HashedPassword = updated
		if updateErr := s.userDB.UpdateUser(ctx, user); updateErr != nil {
			return false, fmt.Errorf("saving updated password hash: %w", updateErr)
		}

		return loginValid, nil
	} else if err != nil && err != auth.ErrCostTooLow {
		logger.Error(err, "issue validating login")
		return false, fmt.Errorf("validating login: %w", err)
	}

	return loginValid, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildCookie(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildCookie()

		expected := `
package example

import (
	"net/http"
	"time"
)

// buildCookie provides a consistent way of constructing an HTTP cookie.
func (s *Service) buildCookie(value string, expiry time.Time) (*http.Cookie, error) {
	encoded, err := s.cookieManager.Encode(CookieName, value)
	if err != nil {
		// NOTE: these errors should be infrequent, and should cause alarm when they do occur
		s.logger.WithName(cookieErrorLogName).Error(err, "error encoding cookie")
		return nil, err
	}

	// https://www.calhoun.io/securing-cookies-in-go/
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.config.SecureCookiesOnly,
		Domain:   s.config.CookieDomain,
		Expires:  expiry,
	}

	return cookie, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
