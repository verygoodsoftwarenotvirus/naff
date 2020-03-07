package httpserver

import (
	"bytes"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildIterableAPIRoutesBlock(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildIterableAPIRoutes(proj),
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"fmt"
	chi "github.com/go-chi/chi"
	apples "services/v1/apples"
	bananas "services/v1/bananas"
	cherries "services/v1/cherries"
)

func doSomething() {
	// Apples
	v1Router.Route("/apples", func(applesRouter chi.Router) {
		appleRoute := fmt.Sprintf(numericIDPattern, apples.URIParamKey)
		applesRouter.With(s.applesService.CreationInputMiddleware).Post("/", s.applesService.CreateHandler())
		applesRouter.Get(appleRoute, s.applesService.ReadHandler())
		applesRouter.With(s.applesService.UpdateInputMiddleware).Put(appleRoute, s.applesService.UpdateHandler())
		applesRouter.Delete(appleRoute, s.applesService.ArchiveHandler())
		applesRouter.Get("/", s.applesService.ListHandler())

		// Bananas

		appleRouter.Route("/bananas", func(bananasRouter chi.Router) {
			bananaRoute := fmt.Sprintf(numericIDPattern, bananas.URIParamKey)
			bananasRouter.With(s.bananasService.CreationInputMiddleware).Post("/", s.bananasService.CreateHandler())
			bananasRouter.Get(bananaRoute, s.bananasService.ReadHandler())
			bananasRouter.With(s.bananasService.UpdateInputMiddleware).Put(bananaRoute, s.bananasService.UpdateHandler())
			bananasRouter.Delete(bananaRoute, s.bananasService.ArchiveHandler())
			bananasRouter.Get("/", s.bananasService.ListHandler())

			// Cherries

			bananaRouter.Route("/cherries", func(cherriesRouter chi.Router) {
				cherryRoute := fmt.Sprintf(numericIDPattern, cherries.URIParamKey)
				cherriesRouter.With(s.cherriesService.CreationInputMiddleware).Post("/", s.cherriesService.CreateHandler())
				cherriesRouter.Get(cherryRoute, s.cherriesService.ReadHandler())
				cherriesRouter.With(s.cherriesService.UpdateInputMiddleware).Put(cherryRoute, s.cherriesService.UpdateHandler())
				cherriesRouter.Delete(cherryRoute, s.cherriesService.ArchiveHandler())
				cherriesRouter.Get("/", s.cherriesService.ListHandler())
			})

		})

	})

}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

/*

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

	// all middleware must be defined before routes on a mux

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

	// Frontend routes
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

	router.With(
		s.authService.AuthenticationMiddleware(true),
		s.authService.AdminMiddleware,
	).Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Post("/cycle_cookie_secret", s.authService.CycleSecretHandler())
	})

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.With(s.authService.UserLoginInputMiddleware).Post("/login", s.authService.LoginHandler())
		userRouter.With(s.authService.CookieAuthenticationMiddleware).Post("/logout", s.authService.LogoutHandler())

		userIDPattern := "/" + fmt.Sprintf(oauth2IDPattern, users.URIParamKey)

		userRouter.Get("/", s.usersService.ListHandler())
		userRouter.With(s.usersService.UserInputMiddleware).Post("/", s.usersService.CreateHandler())
		userRouter.Get(userIDPattern, s.usersService.ReadHandler())
		userRouter.Delete(userIDPattern, s.usersService.ArchiveHandler())

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.TOTPSecretRefreshInputMiddleware,
		).Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler())

		userRouter.With(
			s.authService.CookieAuthenticationMiddleware,
			s.usersService.PasswordUpdateInputMiddleware,
		).Put("/password/new", s.usersService.UpdatePasswordHandler())
	})

	router.Route("/oauth2", func(oauth2Router chi.Router) {
		oauth2Router.With(
			s.authService.CookieAuthenticationMiddleware,
			s.oauth2ClientsService.CreationInputMiddleware,
		).Post("/client", s.oauth2ClientsService.CreateHandler())

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

	router.With(s.authService.AuthenticationMiddleware(true)).Route("/api/v1", func(v1Router chi.Router) {
		// Items
		itemRoute := fmt.Sprintf(numericIDPattern, items.URIParamKey)
		v1Router.Route("/items", func(itemsRouter chi.Router) {
			itemsRouter.With(s.itemsService.CreationInputMiddleware).Post("/", s.itemsService.CreateHandler())
			itemsRouter.Get(itemRoute, s.itemsService.ReadHandler())
			itemsRouter.With(s.itemsService.UpdateInputMiddleware).Put(itemRoute, s.itemsService.UpdateHandler())
			itemsRouter.Delete(itemRoute, s.itemsService.ArchiveHandler())
			itemsRouter.Get("/", s.itemsService.ListHandler())
		})

		// Forums
		forumsBasePath := "/forums"
		forumRoute := fmt.Sprintf(numericIDPattern, forums.URIParamKey)
		v1Router.Route(forumsBasePath, func(forumsRouter chi.Router) {
			s.logger.WithValue("forumRoute", forumRoute).Debug("forumRoute value")
			forumsRouter.With(s.forumsService.CreationInputMiddleware).Post("/", s.forumsService.CreateHandler())
			forumsRouter.Get(forumRoute, s.forumsService.ReadHandler())
			forumsRouter.With(s.forumsService.UpdateInputMiddleware).Put(forumRoute, s.forumsService.UpdateHandler())
			forumsRouter.Delete(forumRoute, s.forumsService.ArchiveHandler())
			forumsRouter.Get("/", s.forumsService.ListHandler())
		})
		forumRoute = forumsBasePath + forumRoute

		// Threads
		threadsBasePath := "threads"
		threadRoute := filepath.Join(forumRoute, threadsBasePath, fmt.Sprintf(numericIDPatternWithoutPrefix, threads.URIParamKey))
		v1Router.Route(filepath.Join(forumRoute, threadsBasePath), func(threadsRouter chi.Router) {
			s.logger.WithValue("threadRoute", threadRoute).Debug("threadRoute value")
			threadsRouter.With(s.threadsService.CreationInputMiddleware).Post("/", s.threadsService.CreateHandler())
			threadsRouter.Get(threadRoute, s.threadsService.ReadHandler())
			threadsRouter.With(s.threadsService.UpdateInputMiddleware).Put(threadRoute, s.threadsService.UpdateHandler())
			threadsRouter.Delete(threadRoute, s.threadsService.ArchiveHandler())
			threadsRouter.Get("/", s.threadsService.ListHandler())

		})

		// Comments
		commentsBasePath := "comments"
		commentRoute := filepath.Join(threadRoute, commentsBasePath, fmt.Sprintf(numericIDPatternWithoutPrefix, comments.URIParamKey))
		v1Router.Route(filepath.Join(threadRoute, commentsBasePath), func(commentsRouter chi.Router) {
			s.logger.WithValue("commentRoute", commentRoute).Debug("commentRoute value")
			commentsRouter.With(s.commentsService.CreationInputMiddleware).Post("/", s.commentsService.CreateHandler())
			commentsRouter.Get(commentRoute, s.commentsService.ReadHandler())
			commentsRouter.With(s.commentsService.UpdateInputMiddleware).Put(commentRoute, s.commentsService.UpdateHandler())
			commentsRouter.Delete(commentRoute, s.commentsService.ArchiveHandler())
			commentsRouter.Get("/", s.commentsService.ListHandler())
		})

		// Tags
		tagRoute := fmt.Sprintf(numericIDPattern, tags.URIParamKey)
		v1Router.Route("/tags", func(tagsRouter chi.Router) {
			tagsRouter.With(s.tagsService.CreationInputMiddleware).Post("/", s.tagsService.CreateHandler())
			tagsRouter.Get(tagRoute, s.tagsService.ReadHandler())
			tagsRouter.With(s.tagsService.UpdateInputMiddleware).Put(tagRoute, s.tagsService.UpdateHandler())
			tagsRouter.Delete(tagRoute, s.tagsService.ArchiveHandler())
			tagsRouter.Get("/", s.tagsService.ListHandler())
		})

		// Webhooks
		webhooksRoute := fmt.Sprintf(numericIDPattern, webhooks.URIParamKey)
		v1Router.Route("/webhooks", func(webhookRouter chi.Router) {
			webhookRouter.With(s.webhooksService.CreationInputMiddleware).Post("/", s.webhooksService.CreateHandler())
			webhookRouter.Get(webhooksRoute, s.webhooksService.ReadHandler())
			webhookRouter.With(s.webhooksService.UpdateInputMiddleware).Put(webhooksRoute, s.webhooksService.UpdateHandler())
			webhookRouter.Delete(webhooksRoute, s.webhooksService.ArchiveHandler())
			webhookRouter.Get("/", s.webhooksService.ListHandler())
		})

		// OAuth2 Clients
		oauth2ClientsRoute := fmt.Sprintf(numericIDPattern, oauth2clients.URIParamKey)
		v1Router.Route("/oauth2/clients", func(clientRouter chi.Router) {
			// CreateHandler is not bound to an OAuth2 authentication token
			// UpdateHandler not supported for OAuth2 clients.
			clientRouter.Get(oauth2ClientsRoute, s.oauth2ClientsService.ReadHandler())
			clientRouter.Delete(oauth2ClientsRoute, s.oauth2ClientsService.ArchiveHandler())
			clientRouter.Get("/", s.oauth2ClientsService.ListHandler())
		})
	})

	s.router = router
}


*/
