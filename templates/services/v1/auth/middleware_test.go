package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_middlewareDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := middlewareDotGo(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

const (
	// userLoginInputMiddlewareCtxKey is the context key for login input.
	userLoginInputMiddlewareCtxKey v1.ContextKey = "user_login_input"

	// usernameFormKey is the string we look for in request forms for username information.
	usernameFormKey = "username"
	// passwordFormKey is the string we look for in request forms for password information.
	passwordFormKey = "password"
	// totpTokenFormKey is the string we look for in request forms for TOTP token information.
	totpTokenFormKey = "totpToken"
)

// CookieAuthenticationMiddleware checks every request for a user cookie.
func (s *Service) CookieAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CookieAuthenticationMiddleware")
		defer span.End()

		// fetch the user from the request.
		user, err := s.fetchUserFromCookie(ctx, req)
		if err != nil {
			s.logger.Error(err, "error encountered fetching user")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user != nil {
			req = req.WithContext(
				context.WithValue(
					ctx,
					v1.SessionInfoKey,
					user.ToSessionInfo(),
				),
			)
			next.ServeHTTP(res, req)
			return
		}

		// if no error was attached to the request, tell them to login first.
		http.Redirect(res, req, "/login", http.StatusUnauthorized)
	})
}

// AuthenticationMiddleware authenticates based on either an oauth2 token or a cookie.
func (s *Service) AuthenticationMiddleware(allowValidCookieInLieuOfAValidToken bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx, span := tracing.StartSpan(req.Context(), "AuthenticationMiddleware")
			defer span.End()

			// let's figure out who the user is.
			var user *v1.User

			// check for a cookie first if we can.
			if allowValidCookieInLieuOfAValidToken {
				cookieAuth, err := s.DecodeCookieFromRequest(ctx, req)

				if err == nil && cookieAuth != nil {
					user, err = s.userDB.GetUser(ctx, cookieAuth.UserID)
					if err != nil {
						s.logger.Error(err, "error authenticating request")
						http.Error(res, "fetching user", http.StatusInternalServerError)
						// if we get here, then we just don't have a valid cookie, and we need to move on.
						return
					}
				}
			}

			// if the cookie wasn't present, or didn't indicate who the user is.
			if user == nil {
				// check to see if there is an OAuth2 token for a valid client attached to the request.
				// We do this first because it is presumed to be the primary means by which requests are made to the httpServer.
				oauth2Client, err := s.oauth2ClientsService.ExtractOAuth2ClientFromRequest(ctx, req)
				if err != nil || oauth2Client == nil {
					s.logger.Error(err, "fetching oauth2 client")
					http.Redirect(res, req, "/login", http.StatusUnauthorized)
					return
				}

				// attach the oauth2 client and user's info to the request.
				ctx = context.WithValue(ctx, v1.OAuth2ClientKey, oauth2Client)
				user, err = s.userDB.GetUser(ctx, oauth2Client.BelongsToUser)
				if err != nil {
					s.logger.Error(err, "error authenticating request")
					http.Error(res, "fetching user", http.StatusInternalServerError)
					return
				}
			}

			// If your request gets here, you're likely either trying to get here, or desperately trying to get anywhere.
			if user == nil {
				s.logger.Debug("no user attached to request request")
				http.Redirect(res, req, "/login", http.StatusUnauthorized)
				return
			}

			// elsewise, load the request with extra context.
			ctx = context.WithValue(ctx, v1.SessionInfoKey, user.ToSessionInfo())

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}

// AdminMiddleware restricts requests to admin users only.
func (s *Service) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "AdminMiddleware")
		defer span.End()

		logger := s.logger.WithRequest(req)
		si, ok := ctx.Value(v1.SessionInfoKey).(*v1.SessionInfo)

		if !ok || si == nil {
			logger.Debug("AdminMiddleware called without user attached to context")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !si.UserIsAdmin {
			logger.Debug("AdminMiddleware called by non-admin user")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}

// parseLoginInputFromForm checks a request for a login form, and returns the parsed login data if relevant.
func parseLoginInputFromForm(req *http.Request) *v1.UserLoginInput {
	if err := req.ParseForm(); err == nil {
		uli := &v1.UserLoginInput{
			Username:  req.FormValue(usernameFormKey),
			Password:  req.FormValue(passwordFormKey),
			TOTPToken: req.FormValue(totpTokenFormKey),
		}

		if uli.Username != "" && uli.Password != "" && uli.TOTPToken != "" {
			return uli
		}
	}
	return nil
}

// UserLoginInputMiddleware fetches user login input from requests.
func (s *Service) UserLoginInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UserLoginInputMiddleware")
		defer span.End()

		x := new(v1.UserLoginInput)
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			if x = parseLoginInputFromForm(req); x == nil {
				s.logger.Error(err, "error encountered decoding request body")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		ctx = context.WithValue(ctx, userLoginInputMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMiddlewareConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildMiddlewareConstantDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

const (
	// userLoginInputMiddlewareCtxKey is the context key for login input.
	userLoginInputMiddlewareCtxKey v1.ContextKey = "user_login_input"

	// usernameFormKey is the string we look for in request forms for username information.
	usernameFormKey = "username"
	// passwordFormKey is the string we look for in request forms for password information.
	passwordFormKey = "password"
	// totpTokenFormKey is the string we look for in request forms for TOTP token information.
	totpTokenFormKey = "totpToken"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCookieAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCookieAuthenticationMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// CookieAuthenticationMiddleware checks every request for a user cookie.
func (s *Service) CookieAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "CookieAuthenticationMiddleware")
		defer span.End()

		// fetch the user from the request.
		user, err := s.fetchUserFromCookie(ctx, req)
		if err != nil {
			s.logger.Error(err, "error encountered fetching user")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user != nil {
			req = req.WithContext(
				context.WithValue(
					ctx,
					v1.SessionInfoKey,
					user.ToSessionInfo(),
				),
			)
			next.ServeHTTP(res, req)
			return
		}

		// if no error was attached to the request, tell them to login first.
		http.Redirect(res, req, "/login", http.StatusUnauthorized)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildAuthenticationMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// AuthenticationMiddleware authenticates based on either an oauth2 token or a cookie.
func (s *Service) AuthenticationMiddleware(allowValidCookieInLieuOfAValidToken bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx, span := tracing.StartSpan(req.Context(), "AuthenticationMiddleware")
			defer span.End()

			// let's figure out who the user is.
			var user *v1.User

			// check for a cookie first if we can.
			if allowValidCookieInLieuOfAValidToken {
				cookieAuth, err := s.DecodeCookieFromRequest(ctx, req)

				if err == nil && cookieAuth != nil {
					user, err = s.userDB.GetUser(ctx, cookieAuth.UserID)
					if err != nil {
						s.logger.Error(err, "error authenticating request")
						http.Error(res, "fetching user", http.StatusInternalServerError)
						// if we get here, then we just don't have a valid cookie, and we need to move on.
						return
					}
				}
			}

			// if the cookie wasn't present, or didn't indicate who the user is.
			if user == nil {
				// check to see if there is an OAuth2 token for a valid client attached to the request.
				// We do this first because it is presumed to be the primary means by which requests are made to the httpServer.
				oauth2Client, err := s.oauth2ClientsService.ExtractOAuth2ClientFromRequest(ctx, req)
				if err != nil || oauth2Client == nil {
					s.logger.Error(err, "fetching oauth2 client")
					http.Redirect(res, req, "/login", http.StatusUnauthorized)
					return
				}

				// attach the oauth2 client and user's info to the request.
				ctx = context.WithValue(ctx, v1.OAuth2ClientKey, oauth2Client)
				user, err = s.userDB.GetUser(ctx, oauth2Client.BelongsToUser)
				if err != nil {
					s.logger.Error(err, "error authenticating request")
					http.Error(res, "fetching user", http.StatusInternalServerError)
					return
				}
			}

			// If your request gets here, you're likely either trying to get here, or desperately trying to get anywhere.
			if user == nil {
				s.logger.Debug("no user attached to request request")
				http.Redirect(res, req, "/login", http.StatusUnauthorized)
				return
			}

			// elsewise, load the request with extra context.
			ctx = context.WithValue(ctx, v1.SessionInfoKey, user.ToSessionInfo())

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAdminMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildAdminMiddleware(proj)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// AdminMiddleware restricts requests to admin users only.
func (s *Service) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "AdminMiddleware")
		defer span.End()

		logger := s.logger.WithRequest(req)
		si, ok := ctx.Value(v1.SessionInfoKey).(*v1.SessionInfo)

		if !ok || si == nil {
			logger.Debug("AdminMiddleware called without user attached to context")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !si.UserIsAdmin {
			logger.Debug("AdminMiddleware called by non-admin user")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildparseLoginInputFromForm(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildparseLoginInputFromForm(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// parseLoginInputFromForm checks a request for a login form, and returns the parsed login data if relevant.
func parseLoginInputFromForm(req *http.Request) *v1.UserLoginInput {
	if err := req.ParseForm(); err == nil {
		uli := &v1.UserLoginInput{
			Username:  req.FormValue(usernameFormKey),
			Password:  req.FormValue(passwordFormKey),
			TOTPToken: req.FormValue(totpTokenFormKey),
		}

		if uli.Username != "" && uli.Password != "" && uli.TOTPToken != "" {
			return uli
		}
	}
	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserLoginInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUserLoginInputMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// UserLoginInputMiddleware fetches user login input from requests.
func (s *Service) UserLoginInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx, span := tracing.StartSpan(req.Context(), "UserLoginInputMiddleware")
		defer span.End()

		x := new(v1.UserLoginInput)
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			if x = parseLoginInputFromForm(req); x == nil {
				s.logger.Error(err, "error encountered decoding request body")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		ctx = context.WithValue(ctx, userLoginInputMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
