package users

import (
	"context"
	"net/http"

	"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
	"go.opencensus.io/trace"
)

var (
	UserCreationMiddlewareCtxKey      models.ContextKey = "user_creation_input"
	PasswordChangeMiddlewareCtxKey    models.ContextKey = "user_password_change"
	TOTPSecretRefreshMiddlewareCtxKey models.ContextKey = "totp_refresh"
)

// UserInputMiddleware fetches user input from requests
func (s *Service) UserInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(models.UserInput)
		ctx, span := trace.StartSpan(req.Context(), "UserInputMiddleware")
		defer span.End()
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, UserCreationMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// PasswordUpdateInputMiddleware fetches password update input from requests
func (s *Service) PasswordUpdateInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(models.PasswordUpdateInput)
		ctx, span := trace.StartSpan(req.Context(), "PasswordUpdateInputMiddleware")
		defer span.End()
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, PasswordChangeMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// TOTPSecretRefreshInputMiddleware fetches 2FA update input from requests
func (s *Service) TOTPSecretRefreshInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		x := new(models.TOTPSecretRefreshInput)
		ctx, span := trace.StartSpan(req.Context(), "TOTPSecretRefreshInputMiddleware")
		defer span.End()
		if err := s.encoderDecoder.DecodeRequest(req, x); err != nil {
			s.logger.Error(err, "error encountered decoding request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, TOTPSecretRefreshMiddlewareCtxKey, x)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
