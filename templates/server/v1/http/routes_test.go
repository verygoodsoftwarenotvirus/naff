package httpserver

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_routesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := routesDotGo(proj)

		expected := `
package example

import (
	"fmt"
	chi "github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/middleware"
	cors "github.com/go-chi/cors"
	healthcheck "github.com/heptiolabs/healthcheck"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"net/http"
)

const (
	root             = "/"
	searchRoot       = "/search"
	numericIDPattern = "/{%s:[0-9]+}"
	oauth2IDPattern  = "/{%s:[0-9_\\-]+}"
)

func (s *Server) setupRouter(frontendConfig config.FrontendSettings, metricsHandler metrics.Handler) {
	router := chi.NewRouter()

	// Basic CORS, for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	ch := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts,
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Provider",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		// Maximum value not ignored by any of major browsers,
		MaxAge: 300,
	})

	router.Use(
		middleware.RequestID,
		middleware.Timeout(maxTimeout),
		s.loggingMiddleware,
		ch.Handler,
	)

	// all middleware must be defined before routes on a mux.

	router.Route("/_meta_", func(metaRouter chi.Router) {
		health := healthcheck.NewHandler()
		// Expose a liveness check on /live
		metaRouter.Get("/live", health.LiveEndpoint)
		// Expose a readiness check on /ready
		metaRouter.Get("/ready", health.ReadyEndpoint)
	})

	if metricsHandler != nil {
		s.logger.Debug("establishing metrics handler")
		router.Handle("/metrics", metricsHandler)
	}

	// Frontend routes.
	if s.config.Frontend.StaticFilesDirectory != "" {
		s.logger.Debug("setting static file server")
		staticFileServer, err := s.frontendService.StaticDir(frontendConfig.StaticFilesDirectory)
		if err != nil {
			s.logger.Error(err, "establishing static file server")
		}
		router.Get("/*", staticFileServer)
	}

	router.With(
		s.authService.AuthenticationMiddleware(true),
		s.authService.AdminMiddleware,
	).Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Post("/cycle_cookie_secret", s.authService.CycleSecretHandler)
	})

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.With(s.authService.UserLoginInputMiddleware).Post("/login", s.authService.LoginHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Post("/logout", s.authService.LogoutHandler)

		userIDPattern := fmt.Sprintf(oauth2IDPattern, users.URIParamKey)

		userRouter.Get(root, s.usersService.ListHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Get("/status", s.authService.StatusHandler)
		userRouter.With(s.usersService.UserInputMiddleware).Post(root, s.usersService.CreateHandler)
		userRouter.Get(userIDPattern, s.usersService.ReadHandler)
		userRouter.Delete(userIDPattern, s.usersService.ArchiveHandler)

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.TOTPSecretRefreshInputMiddleware,
		).Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler)

		userRouter.With(
			s.usersService.TOTPSecretVerificationInputMiddleware,
		).Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)
		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.PasswordUpdateInputMiddleware,
		).Put("/password/new", s.usersService.UpdatePasswordHandler)
	})

	router.Route("/oauth2", func(oauth2Router chi.Router) {
		oauth2Router.With(
			s.authService.CookieAuthenticationMiddleware,
			s.oauth2ClientsService.CreationInputMiddleware,
		).Post("/client", s.oauth2ClientsService.CreateHandler)

		oauth2Router.With(s.oauth2ClientsService.OAuth2ClientInfoMiddleware).
			Post("/authorize", func(res http.ResponseWriter, req *http.Request) {
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

	router.With(s.authService.AuthenticationMiddleware(true)).
		Route("/api/v1", func(v1Router chi.Router) {
			// Items
			itemPath := "items"
			itemsRouteWithPrefix := fmt.Sprintf("/%s", itemPath)
			itemRouteParam := fmt.Sprintf(numericIDPattern, items.URIParamKey)
			v1Router.Route(itemsRouteWithPrefix, func(itemsRouter chi.Router) {
				itemsRouter.With(s.itemsService.CreationInputMiddleware).Post(root, s.itemsService.CreateHandler)
				itemsRouter.Route(itemRouteParam, func(singleItemRouter chi.Router) {
					singleItemRouter.Get(root, s.itemsService.ReadHandler)
					singleItemRouter.With(s.itemsService.UpdateInputMiddleware).Put(root, s.itemsService.UpdateHandler)
					singleItemRouter.Delete(root, s.itemsService.ArchiveHandler)
					singleItemRouter.Head(root, s.itemsService.ExistenceHandler)
				})
				itemsRouter.Get(root, s.itemsService.ListHandler)
				itemsRouter.Get(searchRoot, s.itemsService.SearchHandler)
			})

			// Webhooks.
			v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, webhooks.URIParamKey)
				webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post(root, s.webhooksService.CreateHandler)
				webhookRouter.Get(sr, s.webhooksService.ReadHandler)
				webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler)
				webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler)
				webhookRouter.Get(root, s.webhooksService.ListHandler)
			})

			// OAuth2 Clients.
			v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, oauth2clients.URIParamKey)
				// CreateHandler is not bound to an OAuth2 authentication token.
				// UpdateHandler not supported for OAuth2 clients.
				clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler)
				clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler)
				clientRouter.Get(root, s.oauth2ClientsService.ListHandler)
			})
		})

	s.router = router
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with search disabled", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		for i := range proj.DataTypes {
			proj.DataTypes[i].SearchEnabled = false
		}
		x := routesDotGo(proj)

		expected := `
package example

import (
	"fmt"
	chi "github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/middleware"
	cors "github.com/go-chi/cors"
	healthcheck "github.com/heptiolabs/healthcheck"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"net/http"
)

const (
	root             = "/"
	numericIDPattern = "/{%s:[0-9]+}"
	oauth2IDPattern  = "/{%s:[0-9_\\-]+}"
)

func (s *Server) setupRouter(frontendConfig config.FrontendSettings, metricsHandler metrics.Handler) {
	router := chi.NewRouter()

	// Basic CORS, for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	ch := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts,
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Provider",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		// Maximum value not ignored by any of major browsers,
		MaxAge: 300,
	})

	router.Use(
		middleware.RequestID,
		middleware.Timeout(maxTimeout),
		s.loggingMiddleware,
		ch.Handler,
	)

	// all middleware must be defined before routes on a mux.

	router.Route("/_meta_", func(metaRouter chi.Router) {
		health := healthcheck.NewHandler()
		// Expose a liveness check on /live
		metaRouter.Get("/live", health.LiveEndpoint)
		// Expose a readiness check on /ready
		metaRouter.Get("/ready", health.ReadyEndpoint)
	})

	if metricsHandler != nil {
		s.logger.Debug("establishing metrics handler")
		router.Handle("/metrics", metricsHandler)
	}

	// Frontend routes.
	if s.config.Frontend.StaticFilesDirectory != "" {
		s.logger.Debug("setting static file server")
		staticFileServer, err := s.frontendService.StaticDir(frontendConfig.StaticFilesDirectory)
		if err != nil {
			s.logger.Error(err, "establishing static file server")
		}
		router.Get("/*", staticFileServer)
	}

	router.With(
		s.authService.AuthenticationMiddleware(true),
		s.authService.AdminMiddleware,
	).Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Post("/cycle_cookie_secret", s.authService.CycleSecretHandler)
	})

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.With(s.authService.UserLoginInputMiddleware).Post("/login", s.authService.LoginHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Post("/logout", s.authService.LogoutHandler)

		userIDPattern := fmt.Sprintf(oauth2IDPattern, users.URIParamKey)

		userRouter.Get(root, s.usersService.ListHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Get("/status", s.authService.StatusHandler)
		userRouter.With(s.usersService.UserInputMiddleware).Post(root, s.usersService.CreateHandler)
		userRouter.Get(userIDPattern, s.usersService.ReadHandler)
		userRouter.Delete(userIDPattern, s.usersService.ArchiveHandler)

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.TOTPSecretRefreshInputMiddleware,
		).Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler)

		userRouter.With(
			s.usersService.TOTPSecretVerificationInputMiddleware,
		).Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)
		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.PasswordUpdateInputMiddleware,
		).Put("/password/new", s.usersService.UpdatePasswordHandler)
	})

	router.Route("/oauth2", func(oauth2Router chi.Router) {
		oauth2Router.With(
			s.authService.CookieAuthenticationMiddleware,
			s.oauth2ClientsService.CreationInputMiddleware,
		).Post("/client", s.oauth2ClientsService.CreateHandler)

		oauth2Router.With(s.oauth2ClientsService.OAuth2ClientInfoMiddleware).
			Post("/authorize", func(res http.ResponseWriter, req *http.Request) {
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

	router.With(s.authService.AuthenticationMiddleware(true)).
		Route("/api/v1", func(v1Router chi.Router) {
			// Items
			itemPath := "items"
			itemsRouteWithPrefix := fmt.Sprintf("/%s", itemPath)
			itemRouteParam := fmt.Sprintf(numericIDPattern, items.URIParamKey)
			v1Router.Route(itemsRouteWithPrefix, func(itemsRouter chi.Router) {
				itemsRouter.With(s.itemsService.CreationInputMiddleware).Post(root, s.itemsService.CreateHandler)
				itemsRouter.Route(itemRouteParam, func(singleItemRouter chi.Router) {
					singleItemRouter.Get(root, s.itemsService.ReadHandler)
					singleItemRouter.With(s.itemsService.UpdateInputMiddleware).Put(root, s.itemsService.UpdateHandler)
					singleItemRouter.Delete(root, s.itemsService.ArchiveHandler)
					singleItemRouter.Head(root, s.itemsService.ExistenceHandler)
				})
				itemsRouter.Get(root, s.itemsService.ListHandler)
			})

			// Webhooks.
			v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, webhooks.URIParamKey)
				webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post(root, s.webhooksService.CreateHandler)
				webhookRouter.Get(sr, s.webhooksService.ReadHandler)
				webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler)
				webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler)
				webhookRouter.Get(root, s.webhooksService.ListHandler)
			})

			// OAuth2 Clients.
			v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, oauth2clients.URIParamKey)
				// CreateHandler is not bound to an OAuth2 authentication token.
				// UpdateHandler not supported for OAuth2 clients.
				clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler)
				clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler)
				clientRouter.Get(root, s.oauth2ClientsService.ListHandler)
			})
		})

	s.router = router
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCORSHandlerDef(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildCORSHandlerDef()

		expected := `
package main

import (
	cors "github.com/go-chi/cors"
	"net/http"
)

func main() {
	// Basic CORS, for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	ch := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts,
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Provider",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		// Maximum value not ignored by any of major browsers,
		MaxAge: 300,
	})
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIterableAPIRoutes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildIterableAPIRoutes(proj)

		expected := `
package main

import (
	"fmt"
	chi "github.com/go-chi/chi"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
)

func main() {
	// Items
	itemPath := "items"
	itemsRouteWithPrefix := fmt.Sprintf("/%s", itemPath)
	itemRouteParam := fmt.Sprintf(numericIDPattern, items.URIParamKey)
	v1Router.Route(itemsRouteWithPrefix, func(itemsRouter chi.Router) {
		itemsRouter.With(s.itemsService.CreationInputMiddleware).Post(root, s.itemsService.CreateHandler)
		itemsRouter.Route(itemRouteParam, func(singleItemRouter chi.Router) {
			singleItemRouter.Get(root, s.itemsService.ReadHandler)
			singleItemRouter.With(s.itemsService.UpdateInputMiddleware).Put(root, s.itemsService.UpdateHandler)
			singleItemRouter.Delete(root, s.itemsService.ArchiveHandler)
			singleItemRouter.Head(root, s.itemsService.ExistenceHandler)
		})
		itemsRouter.Get(root, s.itemsService.ListHandler)
		itemsRouter.Get(searchRoot, s.itemsService.SearchHandler)
	})

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildIterableAPIRoutes(proj)

		expected := `
package main

import (
	"fmt"
	chi "github.com/go-chi/chi"
	anotherthings "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/anotherthings"
	things "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/things"
	yetanotherthings "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/yetanotherthings"
	"path/filepath"
)

func main() {
	// Things
	thingPath := "things"
	thingsRouteWithPrefix := fmt.Sprintf("/%s", thingPath)
	thingRouteParam := fmt.Sprintf(numericIDPattern, things.URIParamKey)
	v1Router.Route(thingsRouteWithPrefix, func(thingsRouter chi.Router) {
		thingsRouter.With(s.thingsService.CreationInputMiddleware).Post(root, s.thingsService.CreateHandler)
		thingsRouter.Route(thingRouteParam, func(singleThingRouter chi.Router) {
			singleThingRouter.Get(root, s.thingsService.ReadHandler)
			singleThingRouter.With(s.thingsService.UpdateInputMiddleware).Put(root, s.thingsService.UpdateHandler)
			singleThingRouter.Delete(root, s.thingsService.ArchiveHandler)
			singleThingRouter.Head(root, s.thingsService.ExistenceHandler)
		})
		thingsRouter.Get(root, s.thingsService.ListHandler)
	})

	// AnotherThings
	anotherThingPath := "another_things"
	anotherThingsRoute := filepath.Join(
		thingPath,
		thingRouteParam,
		anotherThingPath,
	)
	anotherThingsRouteWithPrefix := fmt.Sprintf("/%s", anotherThingsRoute)
	anotherThingRouteParam := fmt.Sprintf(numericIDPattern, anotherthings.URIParamKey)
	v1Router.Route(anotherThingsRouteWithPrefix, func(anotherThingsRouter chi.Router) {
		anotherThingsRouter.With(s.anotherThingsService.CreationInputMiddleware).Post(root, s.anotherThingsService.CreateHandler)
		anotherThingsRouter.Route(anotherThingRouteParam, func(singleAnotherThingRouter chi.Router) {
			singleAnotherThingRouter.Get(root, s.anotherThingsService.ReadHandler)
			singleAnotherThingRouter.With(s.anotherThingsService.UpdateInputMiddleware).Put(root, s.anotherThingsService.UpdateHandler)
			singleAnotherThingRouter.Delete(root, s.anotherThingsService.ArchiveHandler)
			singleAnotherThingRouter.Head(root, s.anotherThingsService.ExistenceHandler)
		})
		anotherThingsRouter.Get(root, s.anotherThingsService.ListHandler)
	})

	// YetAnotherThings
	yetAnotherThingPath := "yet_another_things"
	yetAnotherThingsRoute := filepath.Join(
		thingPath,
		thingRouteParam,
		anotherThingPath,
		anotherThingRouteParam,
		yetAnotherThingPath,
	)
	yetAnotherThingsRouteWithPrefix := fmt.Sprintf("/%s", yetAnotherThingsRoute)
	yetAnotherThingRouteParam := fmt.Sprintf(numericIDPattern, yetanotherthings.URIParamKey)
	v1Router.Route(yetAnotherThingsRouteWithPrefix, func(yetAnotherThingsRouter chi.Router) {
		yetAnotherThingsRouter.With(s.yetAnotherThingsService.CreationInputMiddleware).Post(root, s.yetAnotherThingsService.CreateHandler)
		yetAnotherThingsRouter.Route(yetAnotherThingRouteParam, func(singleYetAnotherThingRouter chi.Router) {
			singleYetAnotherThingRouter.Get(root, s.yetAnotherThingsService.ReadHandler)
			singleYetAnotherThingRouter.With(s.yetAnotherThingsService.UpdateInputMiddleware).Put(root, s.yetAnotherThingsService.UpdateHandler)
			singleYetAnotherThingRouter.Delete(root, s.yetAnotherThingsService.ArchiveHandler)
			singleYetAnotherThingRouter.Head(root, s.yetAnotherThingsService.ExistenceHandler)
		})
		yetAnotherThingsRouter.Get(root, s.yetAnotherThingsService.ListHandler)
	})

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSubRouterDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		doneMap := map[string]bool{}
		routerPrefix := "example"
		x := buildSubRouterDecl(proj, doneMap, routerPrefix, typ)

		expected := `
package main

import (
	chi "github.com/go-chi/chi"
)

func main() {
	exampleRouter.Route("/items", func(itemsRouter chi.Router) {
		itemsRouter.With(s.itemsService.CreationInputMiddleware).Post(root, s.itemsService.CreateHandler)
		itemsRouter.Get(itemRouteParam, s.itemsService.ReadHandler)
		itemsRouter.Head(itemRouteParam, s.itemsService.ExistenceHandler)
		itemsRouter.With(s.itemsService.UpdateInputMiddleware).Put(itemRouteParam, s.itemsService.UpdateHandler)
		itemsRouter.Delete(itemRouteParam, s.itemsService.ArchiveHandler)
		itemsRouter.Get(root, s.itemsService.ListHandler)
	})
}
`
		actual := testutils.RenderIndependentStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIterableAPIRoutesBlock(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		doneMap := map[string]bool{}
		routerPrefix := "example"
		x := buildIterableAPIRoutesBlock(proj, doneMap, routerPrefix, typ)

		expected := `
package main

import (
	chi "github.com/go-chi/chi"
)

func main() {

	// Items

	exampleRouter.Route("/items", func(itemsRouter chi.Router) {
		itemsRouter.With(s.itemsService.CreationInputMiddleware).Post(root, s.itemsService.CreateHandler)
		itemsRouter.Get(itemRouteParam, s.itemsService.ReadHandler)
		itemsRouter.Head(itemRouteParam, s.itemsService.ExistenceHandler)
		itemsRouter.With(s.itemsService.UpdateInputMiddleware).Put(itemRouteParam, s.itemsService.UpdateHandler)
		itemsRouter.Delete(itemRouteParam, s.itemsService.ArchiveHandler)
		itemsRouter.Get(root, s.itemsService.ListHandler)
	})

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.DataTypes[0]
		doneMap := map[string]bool{}
		routerPrefix := "example"
		x := buildIterableAPIRoutesBlock(proj, doneMap, routerPrefix, typ)

		expected := `
package main

import (
	chi "github.com/go-chi/chi"
)

func main() {

	// Things

	exampleRouter.Route("/things", func(thingsRouter chi.Router) {
		thingsRouter.With(s.thingsService.CreationInputMiddleware).Post(root, s.thingsService.CreateHandler)
		thingsRouter.Get(thingRouteParam, s.thingsService.ReadHandler)
		thingsRouter.Head(thingRouteParam, s.thingsService.ExistenceHandler)
		thingsRouter.With(s.thingsService.UpdateInputMiddleware).Put(thingRouteParam, s.thingsService.UpdateHandler)
		thingsRouter.Delete(thingRouteParam, s.thingsService.ArchiveHandler)
		thingsRouter.Get(root, s.thingsService.ListHandler)

		// AnotherThings

		thingsRouter.Route("/another_things", func(anotherThingsRouter chi.Router) {
			anotherThingsRouter.With(s.anotherThingsService.CreationInputMiddleware).Post(root, s.anotherThingsService.CreateHandler)
			anotherThingsRouter.Get(anotherThingRouteParam, s.anotherThingsService.ReadHandler)
			anotherThingsRouter.Head(anotherThingRouteParam, s.anotherThingsService.ExistenceHandler)
			anotherThingsRouter.With(s.anotherThingsService.UpdateInputMiddleware).Put(anotherThingRouteParam, s.anotherThingsService.UpdateHandler)
			anotherThingsRouter.Delete(anotherThingRouteParam, s.anotherThingsService.ArchiveHandler)
			anotherThingsRouter.Get(root, s.anotherThingsService.ListHandler)

			// YetAnotherThings

			anotherThingsRouter.Route("/yet_another_things", func(yetAnotherThingsRouter chi.Router) {
				yetAnotherThingsRouter.With(s.yetAnotherThingsService.CreationInputMiddleware).Post(root, s.yetAnotherThingsService.CreateHandler)
				yetAnotherThingsRouter.Get(yetAnotherThingRouteParam, s.yetAnotherThingsService.ReadHandler)
				yetAnotherThingsRouter.Head(yetAnotherThingRouteParam, s.yetAnotherThingsService.ExistenceHandler)
				yetAnotherThingsRouter.With(s.yetAnotherThingsService.UpdateInputMiddleware).Put(yetAnotherThingRouteParam, s.yetAnotherThingsService.UpdateHandler)
				yetAnotherThingsRouter.Delete(yetAnotherThingRouteParam, s.yetAnotherThingsService.ArchiveHandler)
				yetAnotherThingsRouter.Get(root, s.yetAnotherThingsService.ListHandler)
			})

		})

	})

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSetupRouterFuncDef(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildSetupRouterFuncDef(proj)

		expected := `
package example

import (
	"fmt"
	chi "github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/middleware"
	cors "github.com/go-chi/cors"
	healthcheck "github.com/heptiolabs/healthcheck"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	"net/http"
)

func (s *Server) setupRouter(frontendConfig config.FrontendSettings, metricsHandler metrics.Handler) {
	router := chi.NewRouter()

	// Basic CORS, for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	ch := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts,
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Provider",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		// Maximum value not ignored by any of major browsers,
		MaxAge: 300,
	})

	router.Use(
		middleware.RequestID,
		middleware.Timeout(maxTimeout),
		s.loggingMiddleware,
		ch.Handler,
	)

	// all middleware must be defined before routes on a mux.

	router.Route("/_meta_", func(metaRouter chi.Router) {
		health := healthcheck.NewHandler()
		// Expose a liveness check on /live
		metaRouter.Get("/live", health.LiveEndpoint)
		// Expose a readiness check on /ready
		metaRouter.Get("/ready", health.ReadyEndpoint)
	})

	if metricsHandler != nil {
		s.logger.Debug("establishing metrics handler")
		router.Handle("/metrics", metricsHandler)
	}

	// Frontend routes.
	if s.config.Frontend.StaticFilesDirectory != "" {
		s.logger.Debug("setting static file server")
		staticFileServer, err := s.frontendService.StaticDir(frontendConfig.StaticFilesDirectory)
		if err != nil {
			s.logger.Error(err, "establishing static file server")
		}
		router.Get("/*", staticFileServer)
	}

	router.With(
		s.authService.AuthenticationMiddleware(true),
		s.authService.AdminMiddleware,
	).Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Post("/cycle_cookie_secret", s.authService.CycleSecretHandler)
	})

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.With(s.authService.UserLoginInputMiddleware).Post("/login", s.authService.LoginHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Post("/logout", s.authService.LogoutHandler)

		userIDPattern := fmt.Sprintf(oauth2IDPattern, users.URIParamKey)

		userRouter.Get(root, s.usersService.ListHandler)
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Get("/status", s.authService.StatusHandler)
		userRouter.With(s.usersService.UserInputMiddleware).Post(root, s.usersService.CreateHandler)
		userRouter.Get(userIDPattern, s.usersService.ReadHandler)
		userRouter.Delete(userIDPattern, s.usersService.ArchiveHandler)

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.TOTPSecretRefreshInputMiddleware,
		).Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler)

		userRouter.With(
			s.usersService.TOTPSecretVerificationInputMiddleware,
		).Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)
		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.PasswordUpdateInputMiddleware,
		).Put("/password/new", s.usersService.UpdatePasswordHandler)
	})

	router.Route("/oauth2", func(oauth2Router chi.Router) {
		oauth2Router.With(
			s.authService.CookieAuthenticationMiddleware,
			s.oauth2ClientsService.CreationInputMiddleware,
		).Post("/client", s.oauth2ClientsService.CreateHandler)

		oauth2Router.With(s.oauth2ClientsService.OAuth2ClientInfoMiddleware).
			Post("/authorize", func(res http.ResponseWriter, req *http.Request) {
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

	router.With(s.authService.AuthenticationMiddleware(true)).
		Route("/api/v1", func(v1Router chi.Router) {
			// Items
			itemPath := "items"
			itemsRouteWithPrefix := fmt.Sprintf("/%s", itemPath)
			itemRouteParam := fmt.Sprintf(numericIDPattern, items.URIParamKey)
			v1Router.Route(itemsRouteWithPrefix, func(itemsRouter chi.Router) {
				itemsRouter.With(s.itemsService.CreationInputMiddleware).Post(root, s.itemsService.CreateHandler)
				itemsRouter.Route(itemRouteParam, func(singleItemRouter chi.Router) {
					singleItemRouter.Get(root, s.itemsService.ReadHandler)
					singleItemRouter.With(s.itemsService.UpdateInputMiddleware).Put(root, s.itemsService.UpdateHandler)
					singleItemRouter.Delete(root, s.itemsService.ArchiveHandler)
					singleItemRouter.Head(root, s.itemsService.ExistenceHandler)
				})
				itemsRouter.Get(root, s.itemsService.ListHandler)
				itemsRouter.Get(searchRoot, s.itemsService.SearchHandler)
			})

			// Webhooks.
			v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, webhooks.URIParamKey)
				webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post(root, s.webhooksService.CreateHandler)
				webhookRouter.Get(sr, s.webhooksService.ReadHandler)
				webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler)
				webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler)
				webhookRouter.Get(root, s.webhooksService.ListHandler)
			})

			// OAuth2 Clients.
			v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
				sr := fmt.Sprintf(numericIDPattern, oauth2clients.URIParamKey)
				// CreateHandler is not bound to an OAuth2 authentication token.
				// UpdateHandler not supported for OAuth2 clients.
				clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler)
				clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler)
				clientRouter.Get(root, s.oauth2ClientsService.ListHandler)
			})
		})

	s.router = router
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookAPIRoutes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWebhookAPIRoutes(proj)

		expected := `
package main

import (
	"fmt"
	chi "github.com/go-chi/chi"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
)

func main() {
	// Webhooks.
	v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
		sr := fmt.Sprintf(numericIDPattern, webhooks.URIParamKey)
		webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post(root, s.webhooksService.CreateHandler)
		webhookRouter.Get(sr, s.webhooksService.ReadHandler)
		webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(sr, s.webhooksService.UpdateHandler)
		webhookRouter.Delete(sr, s.webhooksService.ArchiveHandler)
		webhookRouter.Get(root, s.webhooksService.ListHandler)
	})

}
`
		actual := testutils.RenderIndependentStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientsAPIRoutes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildOAuth2ClientsAPIRoutes(proj)

		expected := `
package main

import (
	"fmt"
	chi "github.com/go-chi/chi"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
)

func main() {
	// OAuth2 Clients.
	v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
		sr := fmt.Sprintf(numericIDPattern, oauth2clients.URIParamKey)
		// CreateHandler is not bound to an OAuth2 authentication token.
		// UpdateHandler not supported for OAuth2 clients.
		clientRouter.Get(sr, s.oauth2ClientsService.ReadHandler)
		clientRouter.Delete(sr, s.oauth2ClientsService.ArchiveHandler)
		clientRouter.Get(root, s.oauth2ClientsService.ListHandler)
	})
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
