package users

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
	// userCreationMiddlewareCtxKey is the context key for creation input.
	userCreationMiddlewareCtxKey v1.ContextKey = "user_creation_input"

	// passwordChangeMiddlewareCtxKey is the context key for password changes.
	passwordChangeMiddlewareCtxKey v1.ContextKey = "user_password_change"

	// totpSecretVerificationMiddlewareCtxKey is the context key for TOTP token refreshes.
	totpSecretVerificationMiddlewareCtxKey v1.ContextKey = "totp_verify"

	// totpSecretRefreshMiddlewareCtxKey is the context key for TOTP token refreshes.
	totpSecretRefreshMiddlewareCtxKey v1.ContextKey = "totp_refresh"
)

// UserInputMiddleware fetches user input from requests.
func (s *Service) UserInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.UserCreationInput)
		ctx, span := tracing.StartSpan(req.Context(), "UserInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, userCreationMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// PasswordUpdateInputMiddleware fetches password update input from requests.
func (s *Service) PasswordUpdateInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.PasswordUpdateInput)
		ctx, span := tracing.StartSpan(req.Context(), "PasswordUpdateInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, passwordChangeMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// TOTPSecretVerificationInputMiddleware fetches 2FA update input from requests.
func (s *Service) TOTPSecretVerificationInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.TOTPSecretVerificationInput)
		ctx, span := tracing.StartSpan(req.Context(), "TOTPSecretVerificationInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, totpSecretVerificationMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// TOTPSecretRefreshInputMiddleware fetches 2FA update input from requests.
func (s *Service) TOTPSecretRefreshInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.TOTPSecretRefreshInput)
		ctx, span := tracing.StartSpan(req.Context(), "TOTPSecretRefreshInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, totpSecretRefreshMiddlewareCtxKey, x)
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
	// userCreationMiddlewareCtxKey is the context key for creation input.
	userCreationMiddlewareCtxKey v1.ContextKey = "user_creation_input"

	// passwordChangeMiddlewareCtxKey is the context key for password changes.
	passwordChangeMiddlewareCtxKey v1.ContextKey = "user_password_change"

	// totpSecretVerificationMiddlewareCtxKey is the context key for TOTP token refreshes.
	totpSecretVerificationMiddlewareCtxKey v1.ContextKey = "totp_verify"

	// totpSecretRefreshMiddlewareCtxKey is the context key for TOTP token refreshes.
	totpSecretRefreshMiddlewareCtxKey v1.ContextKey = "totp_refresh"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUserInputMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// UserInputMiddleware fetches user input from requests.
func (s *Service) UserInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.UserCreationInput)
		ctx, span := tracing.StartSpan(req.Context(), "UserInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, userCreationMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildPasswordUpdateInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildPasswordUpdateInputMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// PasswordUpdateInputMiddleware fetches password update input from requests.
func (s *Service) PasswordUpdateInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.PasswordUpdateInput)
		ctx, span := tracing.StartSpan(req.Context(), "PasswordUpdateInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, passwordChangeMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTOTPSecretVerificationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTOTPSecretVerificationInputMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// TOTPSecretVerificationInputMiddleware fetches 2FA update input from requests.
func (s *Service) TOTPSecretVerificationInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.TOTPSecretVerificationInput)
		ctx, span := tracing.StartSpan(req.Context(), "TOTPSecretVerificationInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, totpSecretVerificationMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTOTPSecretRefreshInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTOTPSecretRefreshInputMiddleware(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// TOTPSecretRefreshInputMiddleware fetches 2FA update input from requests.
func (s *Service) TOTPSecretRefreshInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(v1.TOTPSecretRefreshInput)
		ctx, span := tracing.StartSpan(req.Context(), "TOTPSecretRefreshInputMiddleware")
		defer span.End()

		// decode the request.
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// attach parsed value to request context.
		ctx = context.WithValue(ctx, totpSecretRefreshMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
