package httpserver

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func routesDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	// pn := typ.Name.Plural()

	ret.Add(
		jen.Const().Defs(
			jen.ID("numericIDPattern").Op("=").Lit(`/{%s:[0-9]+}`),
			jen.ID("oauth2IDPattern").Op("=").Lit(`/{%s:[0-9_\-]+}`),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("setupRouter").Params(jen.ID("frontendConfig").Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "FrontendSettings"),
			jen.ID("metricsHandler").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "Handler")).Block(
			jen.ID("router").Op(":=").Qual("github.com/go-chi/chi", "NewRouter").Call(),
			jen.Line(),
			jen.Comment("Basic CORS, for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing"),
			jen.ID("ch").Op(":=").ID("cors").Dot("New").Call(jen.ID("cors").Dot("Options").Valuesln(
				jen.Comment(`AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts`),
				jen.ID("AllowedOrigins").Op(":").Index().ID("string").Values(jen.Lit("*")),
				jen.Comment(`AllowOriginFunc:  func(r *http.Request, origin string) bool { return true }`),
				jen.ID("AllowedMethods").Op(":").Index().ID("string").Valuesln(
					jen.Qual("net/http", "MethodGet"),
					jen.Qual("net/http", "MethodPost"),
					jen.Qual("net/http", "MethodPut"),
					jen.Qual("net/http", "MethodDelete"),
					jen.Qual("net/http", "MethodOptions"),
				),
				jen.ID("AllowedHeaders").Op(":").Index().ID("string").Valuesln(
					jen.Lit("Accept"),
					jen.Lit("Authorization"),
					jen.Lit("Content-Provider"),
					jen.Lit("X-CSRF-Token"),
				),
				jen.ID("ExposedHeaders").Op(":").Index().ID("string").Values(jen.Lit("Link")),
				jen.ID("AllowCredentials").Op(":").ID("true"),
				jen.Comment("Maximum value not ignored by any of major browsers"),
				jen.ID("MaxAge").Op(":").Lit(300),
			)),
			jen.Line(),
			jen.ID("router").Dot("Use").Callln(
				jen.ID("middleware").Dot("RequestID"),
				jen.ID("middleware").Dot("Timeout").Call(jen.ID("maxTimeout")),
				jen.ID("s").Dot("loggingMiddleware"),
				jen.ID("ch").Dot("Handler"),
			),
			jen.Line(),
			jen.Comment("all middleware must be defined before routes on a mux"),
			jen.Line(),
			jen.ID("router").Dot("Route").Call(jen.Lit("/_meta_"), jen.Func().Params(jen.ID("metaRouter").Qual("github.com/go-chi/chi", "Router")).Block(
				jen.ID("health").Op(":=").ID("healthcheck").Dot("NewHandler").Call(),
				jen.Comment("Expose a liveness check on /live"),
				jen.ID("metaRouter").Dot("Get").Call(jen.Lit("/live"), jen.ID("health").Dot("LiveEndpoint")),
				jen.Comment("Expose a readiness check on /ready"),
				jen.ID("metaRouter").Dot("Get").Call(jen.Lit("/ready"), jen.ID("health").Dot("ReadyEndpoint")),
			)),
			jen.Line(),
			jen.If(jen.ID("metricsHandler").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("establishing metrics handler")),
				jen.ID("router").Dot("Handle").Call(jen.Lit("/metrics"), jen.ID("metricsHandler")),
			),
			jen.Line(),
			jen.Comment("Frontend routes"),
			jen.If(jen.ID("s").Dot("config").Dot("Frontend").Dot("StaticFilesDirectory").Op("!=").Lit("")).Block(
				jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("setting static file server")),
				jen.List(jen.ID("staticFileServer"), jen.ID("err")).Op(":=").ID("s").Dot("frontendService").Dot("StaticDir").Call(jen.ID("frontendConfig").Dot("StaticFilesDirectory")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("establishing static file server")),
				),
				jen.ID("router").Dot("Get").Call(jen.Lit("/*"), jen.ID("staticFileServer")),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("route"), jen.ID("handler")).Op(":=").Range().ID("s").Dot("frontendService").Dot("Routes").Call()).Block(
				jen.ID("router").Dot("Get").Call(jen.ID("route"), jen.ID("handler")),
			),
			jen.Line(),
			jen.ID("router").Dot("With").Callln(
				jen.ID("s").Dot("authService").Dot("AuthenticationMiddleware").Call(jen.ID("true")),
				jen.ID("s").Dot("authService").Dot("AdminMiddleware"),
			).Dot("Route").Call(jen.Lit("/admin"), jen.Func().Params(jen.ID("adminRouter").Qual("github.com/go-chi/chi", "Router")).Block(
				jen.ID("adminRouter").Dot("Post").Call(jen.Lit("/cycle_cookie_secret"), jen.ID("s").Dot("authService").Dot("CycleSecretHandler").Call()),
			)),
			jen.Line(),
			jen.ID("router").Dot("Route").Call(jen.Lit("/users"), jen.Func().Params(jen.ID("userRouter").Qual("github.com/go-chi/chi", "Router")).Block(
				jen.ID("userRouter").Dot("With").Call(jen.ID("s").Dot("authService").Dot("UserLoginInputMiddleware")).Dot("Post").Call(jen.Lit("/login"), jen.ID("s").Dot("authService").Dot("LoginHandler").Call()),
				jen.ID("userRouter").Dot("With").Call(jen.ID("s").Dot("authService").Dot("CookieAuthenticationMiddleware")).Dot("Post").Call(jen.Lit("/logout"), jen.ID("s").Dot("authService").Dot("LogoutHandler").Call()),
				jen.Line(),
				jen.ID("userIDPattern").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("oauth2IDPattern"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/users"), "URIParamKey")),
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
					Dotln("Post").Call(jen.Lit("/authorize"), jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
					jen.ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("Debug").Call(jen.Lit("oauth2 authorize route hit")),
					jen.If(jen.ID("err").Op(":=").ID("s").Dot("oauth2ClientsService").Dot("HandleAuthorizeRequest").Call(jen.ID("res"), jen.ID("req")), jen.ID("err").Op("!=").ID("nil")).Block(
						jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.ID("err").Dot("Error").Call(), jen.Qual("net/http", "StatusBadRequest")),
					),
				)),
				jen.Line(),
				jen.ID("oauth2Router").Dot("Post").Call(jen.Lit("/token"), jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
					jen.If(jen.ID("err").Op(":=").ID("s").Dot("oauth2ClientsService").Dot("HandleTokenRequest").Call(jen.ID("res"), jen.ID("req")), jen.ID("err").Op("!=").ID("nil")).Block(
						jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.ID("err").Dot("Error").Call(), jen.Qual("net/http", "StatusBadRequest")),
					),
				)),
			)),
			jen.Line(),
			jen.ID("router").Dot("With").Call(jen.ID("s").Dot("authService").Dot("AuthenticationMiddleware").Call(jen.ID("true"))).Dot("Route").Call(jen.Lit("/api/v1"), jen.Func().Params(jen.ID("v1Router").Qual("github.com/go-chi/chi", "Router")).Block(
				buildIterableAPIRoutes(pkg),
				jen.Line(),
				buildWebhookAPIRoutes(pkg),
				jen.Line(),
				buildOAuth2ClientsAPIRoutes(pkg),
			)),
			jen.Line(),
			jen.ID("s").Dot("router").Op("=").ID("router"),
		),
		jen.Line(),
	)
	return ret
}

func buildIterableAPIRoutes(pkg *models.Project) jen.Code {

	g := &jen.Group{}

	for _, typ := range pkg.DataTypes {
		n := typ.Name

		g.Add(
			jen.Comment(n.Plural()),
			jen.Line(),
			jen.ID("v1Router").Dot("Route").Call(jen.Litf("/%s", n.PluralRouteName()), jen.Func().Params(jen.IDf("%sRouter", n.PluralUnexportedVarName()).Qual("github.com/go-chi/chi", "Router")).Block(
				jen.ID("sr").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1", n.PluralRouteName()), "URIParamKey")),
				jen.IDf("%sRouter", n.PluralUnexportedVarName()).Dot("With").Call(jen.ID("s").Dotf("%sService", n.PluralUnexportedVarName()).Dot("CreationInputMiddleware")).Dot("Post").Call(jen.Lit("/"), jen.ID("s").Dotf("%sService", n.PluralUnexportedVarName()).Dot("CreateHandler").Call()),
				jen.IDf("%sRouter", n.PluralUnexportedVarName()).Dot("Get").Call(jen.ID("sr"), jen.ID("s").Dotf("%sService", n.PluralUnexportedVarName()).Dot("ReadHandler").Call()),
				jen.IDf("%sRouter", n.PluralUnexportedVarName()).Dot("With").Call(jen.ID("s").Dotf("%sService", n.PluralUnexportedVarName()).Dot("UpdateInputMiddleware")).Dot("Put").Call(jen.ID("sr"), jen.ID("s").Dotf("%sService", n.PluralUnexportedVarName()).Dot("UpdateHandler").Call()),
				jen.IDf("%sRouter", n.PluralUnexportedVarName()).Dot("Delete").Call(jen.ID("sr"), jen.ID("s").Dotf("%sService", n.PluralUnexportedVarName()).Dot("ArchiveHandler").Call()),
				jen.IDf("%sRouter", n.PluralUnexportedVarName()).Dot("Get").Call(jen.Lit("/"), jen.ID("s").Dotf("%sService", n.PluralUnexportedVarName()).Dot("ListHandler").Call()),
			)),
			jen.Line(),
		)
	}

	return g
}

func buildWebhookAPIRoutes(pkg *models.Project) jen.Code {
	g := &jen.Group{}

	g.Add(
		jen.Comment("Webhooks"),
		jen.Line(),
		jen.ID("v1Router").Dot("Route").Call(jen.Lit("/webhooks"), jen.Func().Params(jen.ID("webhookRouter").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("sr").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/webhooks"), "URIParamKey")),
			jen.ID("webhookRouter").Dot("With").Call(jen.ID("s").Dot("webhooksService").Dot("CreationInputMiddleware")).Dot("Post").Call(jen.Lit("/"), jen.ID("s").Dot("webhooksService").Dot("CreateHandler").Call()),
			jen.ID("webhookRouter").Dot("Get").Call(jen.ID("sr"), jen.ID("s").Dot("webhooksService").Dot("ReadHandler").Call()),
			jen.ID("webhookRouter").Dot("With").Call(jen.ID("s").Dot("webhooksService").Dot("UpdateInputMiddleware")).Dot("Put").Call(jen.ID("sr"), jen.ID("s").Dot("webhooksService").Dot("UpdateHandler").Call()),
			jen.ID("webhookRouter").Dot("Delete").Call(jen.ID("sr"), jen.ID("s").Dot("webhooksService").Dot("ArchiveHandler").Call()),
			jen.ID("webhookRouter").Dot("Get").Call(jen.Lit("/"), jen.ID("s").Dot("webhooksService").Dot("ListHandler").Call()),
		)),
	)

	return g
}

func buildOAuth2ClientsAPIRoutes(pkg *models.Project) jen.Code {
	g := &jen.Group{}

	g.Add(
		jen.Comment("OAuth2 Clients"),
		jen.Line(),
		jen.ID("v1Router").Dot("Route").Call(jen.Lit("/oauth2/clients"), jen.Func().Params(jen.ID("clientRouter").Qual("github.com/go-chi/chi", "Router")).Block(
			jen.ID("sr").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/oauth2clients"), "URIParamKey")),
			jen.Comment("CreateHandler is not bound to an OAuth2 authentication token"),
			jen.Comment("UpdateHandler not supported for OAuth2 clients."),
			jen.ID("clientRouter").Dot("Get").Call(jen.ID("sr"), jen.ID("s").Dot("oauth2ClientsService").Dot("ReadHandler").Call()),
			jen.ID("clientRouter").Dot("Delete").Call(jen.ID("sr"), jen.ID("s").Dot("oauth2ClientsService").Dot("ArchiveHandler").Call()),
			jen.ID("clientRouter").Dot("Get").Call(jen.Lit("/"), jen.ID("s").Dot("oauth2ClientsService").Dot("ListHandler").Call()),
		)),
	)

	return g
}
