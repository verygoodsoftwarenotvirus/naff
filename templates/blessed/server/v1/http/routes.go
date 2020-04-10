package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildCORSHandlerDef() []jen.Code {
	return []jen.Code{
		jen.Comment("Basic CORS, for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing"),
		jen.ID("ch").Assign().Qual("github.com/go-chi/cors", "New").Call(jen.Qual("github.com/go-chi/cors", "Options").Valuesln(
			jen.Comment(`AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts`),
			jen.ID("AllowedOrigins").MapAssign().Index().String().Values(jen.Lit("*")),
			jen.Comment(`AllowOriginFunc:  func(r *http.Request, origin string) bool { return true }`),
			jen.ID("AllowedMethods").MapAssign().Index().String().Valuesln(
				jen.Qual("net/http", "MethodGet"),
				jen.Qual("net/http", "MethodPost"),
				jen.Qual("net/http", "MethodPut"),
				jen.Qual("net/http", "MethodDelete"),
				jen.Qual("net/http", "MethodOptions"),
			),
			jen.ID("AllowedHeaders").MapAssign().Index().String().Valuesln(
				jen.Lit("Accept"),
				jen.Lit("Authorization"),
				jen.Lit("Content-Provider"),
				jen.Lit("X-CSRF-Token"),
			),
			jen.ID("ExposedHeaders").MapAssign().Index().String().Values(jen.Lit("Link")),
			jen.ID("AllowCredentials").MapAssign().True(),
			jen.Comment("Maximum value not ignored by any of major browsers"),
			jen.ID("MaxAge").MapAssign().Lit(300),
		)),
	}
}

func routesDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("numericIDPattern").Equals().Lit(`/{%s:[0-9]+}`),
			jen.ID("oauth2IDPattern").Equals().Lit(`/{%s:[0-9_\-]+}`),
		),
		jen.Line(),
	)

	ret.Add(buildSetupRouterFuncDef(proj)...)

	return ret
}

func buildIterableAPIRoutes(proj *models.Project) jen.Code {
	g := &jen.Group{}

	generatedTypes := map[string]bool{}

	for _, typ := range proj.DataTypes {
		n := typ.Name

		if _, ok := generatedTypes[n.Singular()]; !ok {
			g.Add(
				jen.Comment(n.Plural()),
				jen.Line(),
				buildSubRouterDecl(proj, generatedTypes, "v1", typ),
				jen.Line(),
			)
		}

	}

	return g
}

func buildSubRouterDecl(proj *models.Project, doneMap map[string]bool, routerPrefix string, typ models.DataType) jen.Code {
	x := jen.IDf("%sRouter", routerPrefix).Dot("Route").Call(jen.Litf("/%s", typ.Name.PluralRouteName()), jen.Func().Params(jen.IDf("%sRouter", typ.Name.PluralUnexportedVarName()).Qual("github.com/go-chi/chi", "Router")).Block(
		buildIterableAPIRoutesBlock(proj, doneMap, "", typ)...,
	))

	doneMap[typ.Name.Singular()] = true

	return x
}

func buildIterableAPIRoutesBlock(proj *models.Project, doneMap map[string]bool, routerPrefix string, typ models.DataType) []jen.Code {
	n := typ.Name
	puvn := n.PluralUnexportedVarName()

	lines := []jen.Code{}

	singleRouteVar := jen.IDf("single%sRoute", typ.Name.Singular())

	if routerPrefix != "" {
		lines = append(lines,
			jen.Line(),
			jen.Comment(n.Plural()),
			jen.Line(),
			buildSubRouterDecl(proj, doneMap, routerPrefix, typ),
			jen.Line(),
		)
	} else {
		lines = append(lines,
			jen.IDf("single%sRoute", typ.Name.Singular()).Assign().Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.Qual(proj.ServiceV1Package(n.PackageName()), "URIParamKey")),
			jen.IDf("%sRouter", puvn).Dot("With").Call(jen.ID("s").Dotf("%sService", puvn).Dot("CreationInputMiddleware")).Dot("Post").Call(jen.Lit("/"), jen.ID("s").Dotf("%sService", puvn).Dot("CreateHandler").Call()),
			jen.IDf("%sRouter", puvn).Dot("Get").Call(singleRouteVar, jen.ID("s").Dotf("%sService", puvn).Dot("ReadHandler").Call()),
			jen.IDf("%sRouter", puvn).Dot("Head").Call(singleRouteVar, jen.ID("s").Dotf("%sService", puvn).Dot("ExistenceHandler").Call()),
			jen.IDf("%sRouter", puvn).Dot("With").Call(jen.ID("s").Dotf("%sService", puvn).Dot("UpdateInputMiddleware")).Dot("Put").Call(singleRouteVar, jen.ID("s").Dotf("%sService", puvn).Dot("UpdateHandler").Call()),
			jen.IDf("%sRouter", puvn).Dot("Delete").Call(singleRouteVar, jen.ID("s").Dotf("%sService", puvn).Dot("ArchiveHandler").Call()),
			jen.IDf("%sRouter", puvn).Dot("Get").Call(jen.Lit("/"), jen.ID("s").Dotf("%sService", puvn).Dot("ListHandler").Call()),
		)
	}

	deps := proj.FindDependentsOfType(typ)
	for _, dep := range deps {
		if _, ok := doneMap[n.Singular()]; !ok {
			furtherLines := buildIterableAPIRoutesBlock(proj, doneMap, typ.Name.PluralUnexportedVarName(), dep)
			lines = append(lines, jen.Line())
			lines = append(lines, furtherLines...)
		}
	}

	doneMap[n.Singular()] = true

	return lines
}

func buildSetupRouterFuncDef(proj *models.Project) []jen.Code {
	block := []jen.Code{
		jen.ID("router").Assign().Qual("github.com/go-chi/chi", "NewRouter").Call(),
		jen.Line(),
	}
	block = append(block, buildCORSHandlerDef()...)

	block = append(block,
		jen.Line(),
		jen.ID("router").Dot("Use").Callln(
			jen.Qual("github.com/go-chi/chi/middleware", "RequestID"),
			jen.Qual("github.com/go-chi/chi/middleware", "Timeout").Call(jen.ID("maxTimeout")),
			jen.ID("s").Dot("loggingMiddleware"),
			jen.ID("ch").Dot("Handler"),
		),
		jen.Line(),
		jen.Comment("all middleware must be defined before routes on a mux"),
		jen.Line(),
		jen.ID("router").Dot("Route").Call(jen.Lit("/_meta_"), jen.Func().Params(jen.ID("metaRouter").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("health").Assign().Qual("github.com/heptiolabs/healthcheck", "NewHandler").Call(),
			jen.Comment("Expose a liveness check on /live"),
			jen.ID("metaRouter").Dot("Get").Call(jen.Lit("/live"), jen.ID("health").Dot("LiveEndpoint")),
			jen.Comment("Expose a readiness check on /ready"),
			jen.ID("metaRouter").Dot("Get").Call(jen.Lit("/ready"), jen.ID("health").Dot("ReadyEndpoint")),
		)),
		jen.Line(),
		jen.If(jen.ID("metricsHandler").DoesNotEqual().ID("nil")).Block(
			jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("establishing metrics handler")),
			jen.ID("router").Dot("Handle").Call(jen.Lit("/metrics"), jen.ID("metricsHandler")),
		),
		jen.Line(),
		jen.Comment("Frontend routes"),
		jen.If(jen.ID("s").Dot("config").Dot("Frontend").Dot("StaticFilesDirectory").DoesNotEqual().EmptyString()).Block(
			jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("setting static file server")),
			jen.List(jen.ID("staticFileServer"), jen.Err()).Assign().ID("s").Dot("frontendService").Dot("StaticDir").Call(jen.ID("frontendConfig").Dot("StaticFilesDirectory")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("establishing static file server")),
			),
			jen.ID("router").Dot("Get").Call(jen.Lit("/*"), jen.ID("staticFileServer")),
		),
		jen.Line(),
		jen.For(jen.List(jen.ID("route"), jen.ID("handler")).Assign().Range().ID("s").Dot("frontendService").Dot("Routes").Call()).Block(
			jen.ID("router").Dot("Get").Call(jen.ID("route"), jen.ID("handler")),
		),
		jen.Line(),
		jen.ID("router").Dot("With").Callln(
			jen.ID("s").Dot("authService").Dot("AuthenticationMiddleware").Call(jen.True()),
			jen.ID("s").Dot("authService").Dot("AdminMiddleware"),
		).Dot("Route").Call(jen.Lit("/admin"), jen.Func().Params(jen.ID("adminRouter").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("adminRouter").Dot("Post").Call(jen.Lit("/cycle_cookie_secret"), jen.ID("s").Dot("authService").Dot("CycleSecretHandler").Call()),
		)),
		jen.Line(),
		jen.ID("router").Dot("Route").Call(jen.Lit("/users"), jen.Func().Params(jen.ID("userRouter").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("userRouter").Dot("With").Call(jen.ID("s").Dot("authService").Dot("UserLoginInputMiddleware")).Dot("Post").Call(jen.Lit("/login"), jen.ID("s").Dot("authService").Dot("LoginHandler").Call()),
			jen.ID("userRouter").Dot("With").Call(jen.ID("s").Dot("authService").Dot("CookieAuthenticationMiddleware")).Dot("Post").Call(jen.Lit("/logout"), jen.ID("s").Dot("authService").Dot("LogoutHandler").Call()),
			jen.Line(),
			jen.ID("userIDPattern").Assign().Qual("fmt", "Sprintf").Call(jen.ID("oauth2IDPattern"), jen.Qual(proj.ServiceV1UsersPackage(), "URIParamKey")),
			jen.Line(),
			jen.ID("userRouter").Dot("Get").Call(jen.Lit("/"), jen.ID("s").Dot("usersService").Dot("ListHandler").Call()),
			jen.ID("userRouter").Dot("With").Call(jen.ID("s").Dot("usersService").Dot("UserInputMiddleware")).Dot("Post").Call(jen.Lit("/"), jen.ID("s").Dot("usersService").Dot("CreateHandler").Call()),
			jen.ID("userRouter").Dot("Get").Call(jen.ID("userIDPattern"), jen.ID("s").Dot("usersService").Dot("ReadHandler").Call()),
			jen.ID("userRouter").Dot("Delete").Call(jen.ID("userIDPattern"), jen.ID("s").Dot("usersService").Dot("ArchiveHandler").Call()),
			jen.Line(),
			jen.ID("userRouter").Dot("With").Callln(
				jen.ID("s").Dot("authService").Dot("CookieAuthenticationMiddleware"),
				jen.ID("s").Dot("usersService").Dot("TOTPSecretRefreshInputMiddleware"),
			).Dot("Post").Call(jen.Lit("/totp_secret/new"), jen.ID("s").Dot("usersService").Dot("NewTOTPSecretHandler").Call()),
			jen.Line(),
			jen.ID("userRouter").Dot("With").Callln(
				jen.ID("s").Dot("authService").Dot("CookieAuthenticationMiddleware"),
				jen.ID("s").Dot("usersService").Dot("PasswordUpdateInputMiddleware"),
			).Dot("Put").Call(jen.Lit("/password/new"), jen.ID("s").Dot("usersService").Dot("UpdatePasswordHandler").Call()),
		)),
		jen.Line(),
		jen.ID("router").Dot("Route").Call(jen.Lit("/oauth2"), jen.Func().Params(jen.ID("oauth2Router").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("oauth2Router").Dot("With").Callln(
				jen.ID("s").Dot("authService").Dot("CookieAuthenticationMiddleware"),
				jen.ID("s").Dot("oauth2ClientsService").Dot("CreationInputMiddleware"),
			).Dot("Post").Call(jen.Lit("/client"), jen.ID("s").Dot("oauth2ClientsService").Dot("CreateHandler").Call()),
			jen.Line(),
			jen.ID("oauth2Router").Dot("With").Call(jen.ID("s").Dot("oauth2ClientsService").Dot("OAuth2ClientInfoMiddleware")).
				Dotln("Post").Call(jen.Lit("/authorize"), jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Block(
				jen.ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("Debug").Call(jen.Lit("oauth2 authorize route hit")),
				jen.If(jen.Err().Assign().ID("s").Dot("oauth2ClientsService").Dot("HandleAuthorizeRequest").Call(jen.ID("res"), jen.ID("req")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Err().Dot("Error").Call(), jen.Qual("net/http", "StatusBadRequest")),
				),
			)),
			jen.Line(),
			jen.ID("oauth2Router").Dot("Post").Call(jen.Lit("/token"), jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Block(
				jen.If(jen.Err().Assign().ID("s").Dot("oauth2ClientsService").Dot("HandleTokenRequest").Call(jen.ID("res"), jen.ID("req")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Err().Dot("Error").Call(), jen.Qual("net/http", "StatusBadRequest")),
				),
			)),
		)),
		jen.Line(),
		jen.ID("router").Dot("With").Call(jen.ID("s").Dot("authService").Dot("AuthenticationMiddleware").Call(jen.True())).Dot("Route").Call(jen.Lit("/api/v1"), jen.Func().Params(jen.ID("v1Router").Qual("github.com/go-chi/chi", "Router")).Block(
			buildIterableAPIRoutes(proj),
			jen.Line(),
			buildWebhookAPIRoutes(proj),
			jen.Line(),
			buildOAuth2ClientsAPIRoutes(proj),
		)),
		jen.Line(),
		jen.ID("s").Dot("router").Equals().ID("router"),
	)

	return []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("Server")).ID("setupRouter").Params(jen.ID("frontendConfig").Qual(proj.InternalConfigV1Package(), "FrontendSettings"),
			jen.ID("metricsHandler").Qual(proj.InternalMetricsV1Package(), "Handler")).Block(block...),
		jen.Line(),
	}
}

func buildWebhookAPIRoutes(proj *models.Project) jen.Code {
	g := &jen.Group{}

	g.Add(
		jen.Comment("Webhooks"),
		jen.Line(),
		jen.ID("v1Router").Dot("Route").Call(jen.Lit("/webhooks"), jen.Func().Params(jen.ID("webhookRouter").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("sr").Assign().Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.Qual(proj.ServiceV1WebhooksPackage(), "URIParamKey")),
			jen.ID("webhookRouter").Dot("With").Call(jen.ID("s").Dot("webhooksService").Dot("CreationInputMiddleware")).Dot("Post").Call(jen.Lit("/"), jen.ID("s").Dot("webhooksService").Dot("CreateHandler").Call()),
			jen.ID("webhookRouter").Dot("Get").Call(jen.ID("sr"), jen.ID("s").Dot("webhooksService").Dot("ReadHandler").Call()),
			jen.ID("webhookRouter").Dot("With").Call(jen.ID("s").Dot("webhooksService").Dot("UpdateInputMiddleware")).Dot("Put").Call(jen.ID("sr"), jen.ID("s").Dot("webhooksService").Dot("UpdateHandler").Call()),
			jen.ID("webhookRouter").Dot("Delete").Call(jen.ID("sr"), jen.ID("s").Dot("webhooksService").Dot("ArchiveHandler").Call()),
			jen.ID("webhookRouter").Dot("Get").Call(jen.Lit("/"), jen.ID("s").Dot("webhooksService").Dot("ListHandler").Call()),
		)),
	)

	return g
}

func buildOAuth2ClientsAPIRoutes(proj *models.Project) jen.Code {
	g := &jen.Group{}

	g.Add(
		jen.Comment("OAuth2 Clients"),
		jen.Line(),
		jen.ID("v1Router").Dot("Route").Call(jen.Lit("/oauth2/clients"), jen.Func().Params(jen.ID("clientRouter").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("sr").Assign().Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "URIParamKey")),
			jen.Comment("CreateHandler is not bound to an OAuth2 authentication token"),
			jen.Comment("UpdateHandler not supported for OAuth2 clients."),
			jen.ID("clientRouter").Dot("Get").Call(jen.ID("sr"), jen.ID("s").Dot("oauth2ClientsService").Dot("ReadHandler").Call()),
			jen.ID("clientRouter").Dot("Delete").Call(jen.ID("sr"), jen.ID("s").Dot("oauth2ClientsService").Dot("ArchiveHandler").Call()),
			jen.ID("clientRouter").Dot("Get").Call(jen.Lit("/"), jen.ID("s").Dot("oauth2ClientsService").Dot("ListHandler").Call()),
		)),
	)

	return g
}
