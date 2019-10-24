package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/heptiolabs/healthcheck"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics"
	"gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items"
)

const (
	numericIDPattern = "/{%s:[0-9]+}"
	oauth2IDPattern  = "/{%s:[0-9_\\-]+}"
)

func (s *Server) setupRouter(frontendConfig config.FrontendSettings, metricsHandler metrics.Handler) {
	router := chi.NewRouter()
	ch := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Provider",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(middleware.RequestID, middleware.Timeout(maxTimeout), s.loggingMiddleware, ch.Handler)
	router.Route("/_meta_", func(metaRouter chi.Router) {
		health := healthcheck.NewHandler()
		metaRouter.Get("/live", health.LiveEndpoint)
		metaRouter.Get("/ready", health.ReadyEndpoint)
	})
	if metricsHandler != nil {
		s.logger.Debug("establishing metrics handler")
		router.Handle("/metrics", metricsHandler)
	}
	if s.config.Frontend.StaticFilesDirectory != "" {
		s.logger.Debug("setting static file server")
		staticFileServer, err := s.frontendService.StaticDir(frontendConfig.StaticFilesDirectory)
		if err != nil {
			s.logger.Error(err, "establishing static file server")
		}
		router.Get("/*", staticFileServer)
	}
	for route, handler := range s.frontendService.Routes() {
		router.Get(route, handler)
	}
	router.With(s.authService.AuthenticationMiddleware(true), s.authService.AdminMiddleware).Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Post("/cycle_cookie_secret", s.authService.CycleSecretHandler())
	})
	router.Route("/users", func(userRouter chi.Router) {
		userRouter.With(s.authService.UserLoginInputMiddleware).Post("/login", s.authService.LoginHandler())
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Post("/logout", s.authService.LogoutHandler())
		userIDPattern := fmt.Sprintf(oauth2IDPattern, users.URIParamKey)
		userRouter.Get("/", s.usersService.ListHandler())
		userRouter.With(s.usersService.UserInputMiddleware).Post("/", s.usersService.CreateHandler())
		userRouter.Get(userIDPattern, s.usersService.ReadHandler())
		userRouter.Delete(userIDPattern, s.usersService.ArchiveHandler())
		userRouter.With(s.authService.CookieAuthenticationMiddleware, s.usersService.TOTPSecretRefreshInputMiddleware).Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler())
		userRouter.With(s.authService.CookieAuthenticationMiddleware, s.usersService.PasswordUpdateInputMiddleware).Put("/password/new", s.usersService.UpdatePasswordHandler())
	})
	router.Route("/oauth2", func(oauth2Router chi.Router) {
		oauth2Router.With(s.authService.CookieAuthenticationMiddleware, s.oauth2ClientsService.CreationInputMiddleware).Post("/client", s.oauth2ClientsService.CreateHandler())
		oauth2Router.With(s.oauth2ClientsService.OAuth2ClientInfoMiddleware).Post("/authorize", func(res http.ResponseWriter, req *http.Request) {
			s.logger.WithRequest(req).Debug("oauth2 authorize route hit")
			if err := s.oauth2ClientsService.HandleAuthorizeRequest(res, req); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
			}
		})
		oauth2Router.Post("/token", func(res http.ResponseWriter, req *http.Request) {
			if err := s.oauth2ClientsService.HandleTokenRequest(res, req); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
			}
		})
	})
	router.With(s.authService.AuthenticationMiddleware(true)).Route("/api/v1", func(v1Router chi.Router) {
		v1Router.Route("/items", func(itemsRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, items.URIParamKey)
			itemsRouter.With(s.itemsService.CreationInputMiddleware).Post("/", s.itemsService.CreateHandler())
			itemsRouter.Get(sr, s.itemsService.ReadHandler())
			itemsRouter.With(s.itemsService.UpdateInputMiddleware).Put(sr, s.itemsService.UpdateHandler())
			itemsRouter.Delete(sr, s.itemsService.ArchiveHandler())
			itemsRouter.Get("/", s.itemsService.ListHandler())
		})
		v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, webhooks.URIParamKey)
			webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post("/", s.webhooksService.CreateHandler())
			webhookRouter.Get(sr, s.webhooksService.ReadHandler())
			webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler())
			webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler())
			webhookRouter.Get("/", s.webhooksService.ListHandler())
		})
		v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
			sr := fmt.Sprintf(numericIDPattern, oauth2clients.URIParamKey)
			clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler())
			clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler())
			clientRouter.Get("/", s.oauth2ClientsService.ListHandler())
		})
	})
	s.router = router
}