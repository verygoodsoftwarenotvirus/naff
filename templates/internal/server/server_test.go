package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_serverDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := serverDotGo(proj)

		expected := `
package example

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	chi "github.com/go-chi/chi"
	v12 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	ochttp "go.opencensus.io/plugin/ochttp"
	"net/http"
	"os"
	"time"
)

const (
	maxTimeout      = 120 * time.Second
	serverNamespace = "todo-service"
)

type (
	// Server is our API httpServer.
	Server struct {
		DebugMode bool

		// Services.
		authService          *authentication.Service
		frontendService      *frontend.Service
		usersService         v1.UserDataServer
		oauth2ClientsService v1.OAuth2ClientDataServer
		webhooksService      v1.WebhookDataServer
		itemsService         v1.ItemDataServer

		// infra things.
		db          v11.DataManager
		config      *config.ServerConfig
		router      *chi.Mux
		httpServer  *http.Server
		logger      v12.Logger
		encoder     encoding.EncoderDecoder
		newsManager *newsman.Newsman
	}
)

// ProvideServer builds a new Server instance.
func ProvideServer(
	ctx context.Context,
	cfg *config.ServerConfig,
	authService *authentication.Service,
	frontendService *frontend.Service,
	itemsService v1.ItemDataServer,
	usersService v1.UserDataServer,
	oauth2Service v1.OAuth2ClientDataServer,
	webhooksService v1.WebhookDataServer,
	db v11.DataManager,
	logger v12.Logger,
	encoder encoding.EncoderDecoder,
	newsManager *newsman.Newsman,
) (*Server, error) {
	if len(cfg.Auth.CookieSecret) < 32 {
		err := errors.New("cookie secret is too short, must be at least 32 characters in length")
		logger.Error(err, "cookie secret failure")
		return nil, err
	}

	srv := &Server{
		DebugMode: cfg.Server.Debug,
		// infra things,
		db:          db,
		config:      cfg,
		encoder:     encoder,
		httpServer:  provideHTTPServer(),
		logger:      logger.WithName("api_server"),
		newsManager: newsManager,
		// services,
		webhooksService:      webhooksService,
		frontendService:      frontendService,
		usersService:         usersService,
		authService:          authService,
		itemsService:         itemsService,
		oauth2ClientsService: oauth2Service,
	}

	if err := cfg.ProvideTracing(logger); err != nil && err != config.ErrInvalidTracingProvider {
		return nil, err
	}

	metricsHandler := cfg.ProvideInstrumentationHandler(logger)
	srv.setupRouter(cfg.Frontend, metricsHandler)

	srv.httpServer.Handler = &ochttp.Handler{
		Handler:        srv.router,
		FormatSpanName: formatSpanNameForRequest,
	}

	allWebhooks, err := db.GetAllWebhooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("initializing webhooks: %w", err)
	}

	for i := 0; i < len(allWebhooks.Webhooks); i++ {
		wh := allWebhooks.Webhooks[i]
		// NOTE: we must guarantee that whatever is stored in the database is valid, otherwise
		// newsman will try (and fail) to execute requests constantly
		l := wh.ToListener(srv.logger)
		srv.newsManager.TuneIn(l)
	}

	return srv, nil
}

/*
func (s *Server) logRoutes() {
	if err := chi.Walk(s.router, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		s.logger.WithValues(map[string]interface{}{
			"method": method,
			"route":  route,
		}).Debug("route found")

		return nil
	}); err != nil {
		s.logger.Error(err, "logging routes")
	}
}
*/

// Serve serves HTTP traffic.
func (s *Server) Serve() {
	s.httpServer.Addr = fmt.Sprintf(":%d", s.config.Server.HTTPPort)
	s.logger.Debug(fmt.Sprintf("Listening for HTTP requests on %q", s.httpServer.Addr))

	// returns ErrServerClosed on graceful close.
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.Error(err, "server shutting down")
		if err == http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop.
			os.Exit(0)
		}
	}
}

// provideHTTPServer provides an HTTP httpServer.
func provideHTTPServer() *http.Server {
	// heavily inspired by https://blog.cloudflare.com/exposing-go-on-the-internet/
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			// "Only use curves which have assembly implementations"
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}
	return srv
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerConstantDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServerConstantDefinitions(proj)

		expected := `
package example

import (
	"time"
)

const (
	maxTimeout      = 120 * time.Second
	serverNamespace = "todo-service"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServerTypeDefinitions(proj)

		expected := `
package example

import (
	chi "github.com/go-chi/chi"
	v12 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

type (
	// Server is our API httpServer.
	Server struct {
		DebugMode bool

		// Services.
		authService          *authentication.Service
		frontendService      *frontend.Service
		usersService         v1.UserDataServer
		oauth2ClientsService v1.OAuth2ClientDataServer
		webhooksService      v1.WebhookDataServer
		itemsService         v1.ItemDataServer

		// infra things.
		db          v11.DataManager
		config      *config.ServerConfig
		router      *chi.Mux
		httpServer  *http.Server
		logger      v12.Logger
		encoder     encoding.EncoderDecoder
		newsManager *newsman.Newsman
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerProvideServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildServerProvideServer(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	v12 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	ochttp "go.opencensus.io/plugin/ochttp"
)

// ProvideServer builds a new Server instance.
func ProvideServer(
	ctx context.Context,
	cfg *config.ServerConfig,
	authService *authentication.Service,
	frontendService *frontend.Service,
	itemsService v1.ItemDataServer,
	usersService v1.UserDataServer,
	oauth2Service v1.OAuth2ClientDataServer,
	webhooksService v1.WebhookDataServer,
	db v11.DataManager,
	logger v12.Logger,
	encoder encoding.EncoderDecoder,
	newsManager *newsman.Newsman,
) (*Server, error) {
	if len(cfg.Auth.CookieSecret) < 32 {
		err := errors.New("cookie secret is too short, must be at least 32 characters in length")
		logger.Error(err, "cookie secret failure")
		return nil, err
	}

	srv := &Server{
		DebugMode: cfg.Server.Debug,
		// infra things,
		db:          db,
		config:      cfg,
		encoder:     encoder,
		httpServer:  provideHTTPServer(),
		logger:      logger.WithName("api_server"),
		newsManager: newsManager,
		// services,
		webhooksService:      webhooksService,
		frontendService:      frontendService,
		usersService:         usersService,
		authService:          authService,
		itemsService:         itemsService,
		oauth2ClientsService: oauth2Service,
	}

	if err := cfg.ProvideTracing(logger); err != nil && err != config.ErrInvalidTracingProvider {
		return nil, err
	}

	metricsHandler := cfg.ProvideInstrumentationHandler(logger)
	srv.setupRouter(cfg.Frontend, metricsHandler)

	srv.httpServer.Handler = &ochttp.Handler{
		Handler:        srv.router,
		FormatSpanName: formatSpanNameForRequest,
	}

	allWebhooks, err := db.GetAllWebhooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("initializing webhooks: %w", err)
	}

	for i := 0; i < len(allWebhooks.Webhooks); i++ {
		wh := allWebhooks.Webhooks[i]
		// NOTE: we must guarantee that whatever is stored in the database is valid, otherwise
		// newsman will try (and fail) to execute requests constantly
		l := wh.ToListener(srv.logger)
		srv.newsManager.TuneIn(l)
	}

	return srv, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCommentedOutLogRoutesMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildCommentedOutLogRoutesMethod()

		expected := `
package example

import ()

/*
func (s *Server) logRoutes() {
	if err := chi.Walk(s.router, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		s.logger.WithValues(map[string]interface{}{
			"method": method,
			"route":  route,
		}).Debug("route found")

		return nil
	}); err != nil {
		s.logger.Error(err, "logging routes")
	}
}
*/
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerServe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServerServe()

		expected := `
package example

import (
	"fmt"
	"net/http"
	"os"
)

// Serve serves HTTP traffic.
func (s *Server) Serve() {
	s.httpServer.Addr = fmt.Sprintf(":%d", s.config.Server.HTTPPort)
	s.logger.Debug(fmt.Sprintf("Listening for HTTP requests on %q", s.httpServer.Addr))

	// returns ErrServerClosed on graceful close.
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.Error(err, "server shutting down")
		if err == http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop.
			os.Exit(0)
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerProvideHTTPServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServerProvideHTTPServer()

		expected := `
package example

import (
	"crypto/tls"
	"net/http"
	"time"
)

// provideHTTPServer provides an HTTP httpServer.
func provideHTTPServer() *http.Server {
	// heavily inspired by https://blog.cloudflare.com/exposing-go-on-the-internet/
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			// "Only use curves which have assembly implementations"
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}
	return srv
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
