package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
	"go.opencensus.io/plugin/ochttp"
)

var maxTimeout = 120 * time.Second

type Server struct {
	DebugMode            bool
	authService          *auth.Service
	frontendService      *frontend.Service
	usersService         models.UserDataServer
	oauth2ClientsService models.OAuth2ClientDataServer
	webhooksService      models.WebhookDataServer
	itemsService         models.ItemDataServer
	db                   database.Database
	config               *config.ServerConfig
	router               *chi.Mux
	httpServer           *http.Server
	logger               logging.Logger
	encoder              encoding.EncoderDecoder
	newsManager          *newsman.Newsman
}

// ProvideServer builds a new Server instance
func ProvideServer(ctx context.Context, cfg *config.ServerConfig, authService *auth.Service, frontendService *frontend.Service, itemsService models.ItemDataServer, usersService models.UserDataServer, oauth2Service models.OAuth2ClientDataServer, webhooksService models.WebhookDataServer, db database.Database, logger logging.Logger, encoder encoding.EncoderDecoder, newsManager *newsman.Newsman) (*Server, error) {
	if len(cfg.Auth.CookieSecret) < 32 {
		err := errors.New("cookie secret is too short, must be at least 32 characters in length")
		logger.Error(err, "cookie secret failure")
		return nil, err
	}
	srv := &Server{
		DebugMode:            cfg.Server.Debug,
		db:                   db,
		config:               cfg,
		encoder:              encoder,
		httpServer:           provideHTTPServer(),
		logger:               logger.WithName("api_server"),
		newsManager:          newsManager,
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
	ih, err := cfg.ProvideInstrumentationHandler(logger)
	if err != nil && err != config.ErrInvalidMetricsProvider {
		return nil, err
	}
	if ih != nil {
		srv.setupRouter(cfg.Frontend, ih)
	}
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
		l := wh.ToListener(srv.logger)
		srv.newsManager.TuneIn(l)
	}
	return srv, nil
}

// Serve serves HTTP traffic
func (s *Server) Serve() {
	s.httpServer.Addr = fmt.Sprintf(":%d", s.config.Server.HTTPPort)
	s.logger.Debug(fmt.Sprintf("Listening for HTTP requests on %q", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.Error(err, "server shutting down")
		if err == http.ErrServerClosed {
			os.Exit(0)
		}
	}
}
