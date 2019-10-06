package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func routesDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("numericIDPattern").Op("=").Lit(`/{%s:[0-9]+}`).Var().ID("oauth2IDPattern").Op("=").Lit(`/{%s:[0-9_\-]+}`))
	ret.Add(jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("setupRouter").Params(jen.ID("frontendConfig").ID("config").Dot(
		"FrontendSettings",
	), jen.ID("metricsHandler").ID("metrics").Dot(
		"Handler",
	)).Block(
		jen.ID("router").Op(":=").ID("chi").Dot(
			"NewRouter",
		).Call(),
		jen.ID("ch").Op(":=").ID("cors").Dot(
			"New",
		).Call(jen.ID("cors").Dot(
			"Options",
		).Valuesln(jen.ID("AllowedOrigins").Op(":").Index().ID("string").Valuesln(jen.Lit("*")), jen.ID("AllowedMethods").Op(":").Index().ID("string").Valuesln(jen.Lit("GET"), jen.Lit("POST"), jen.Lit("PUT"), jen.Lit("DELETE"), jen.Lit("OPTIONS")), jen.ID("AllowedHeaders").Op(":").Index().ID("string").Valuesln(jen.Lit("Accept"), jen.Lit("Authorization"), jen.Lit("Content-Provider"), jen.Lit("X-CSRF-Token")), jen.ID("ExposedHeaders").Op(":").Index().ID("string").Valuesln(jen.Lit("Link")), jen.ID("AllowCredentials").Op(":").ID("true"), jen.ID("MaxAge").Op(":").Lit(300))),
		jen.ID("router").Dot(
			"Use",
		).Call(jen.ID("middleware").Dot(
			"RequestID",
		), jen.ID("middleware").Dot(
			"Timeout",
		).Call(jen.ID("maxTimeout")), jen.ID("s").Dot(
			"loggingMiddleware",
		), jen.ID("ch").Dot(
			"Handler",
		)),
		jen.ID("router").Dot(
			"Route",
		).Call(jen.Lit("/_meta_"), jen.Func().Params(jen.ID("metaRouter").ID("chi").Dot(
			"Router",
		)).Block(
			jen.ID("health").Op(":=").ID("healthcheck").Dot(
				"NewHandler",
			).Call(),
			jen.ID("metaRouter").Dot(
				"Get",
			).Call(jen.Lit("/live"), jen.ID("health").Dot(
				"LiveEndpoint",
			)),
			jen.ID("metaRouter").Dot(
				"Get",
			).Call(jen.Lit("/ready"), jen.ID("health").Dot(
				"ReadyEndpoint",
			)),
		)),
		jen.If(
			jen.ID("metricsHandler").Op("!=").ID("nil"),
		).Block(
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"Debug",
			).Call(jen.Lit("establishing metrics handler")),
			jen.ID("router").Dot(
				"Handle",
			).Call(jen.Lit("/metrics"), jen.ID("metricsHandler")),
		),
		jen.If(
			jen.ID("frontendConfig").Dot(
				"StaticFilesDirectory",
			).Op("!=").Lit(""),
		).Block(
			jen.List(jen.ID("staticFileServer"), jen.ID("err")).Op(":=").ID("s").Dot(
				"frontendService",
			).Dot(
				"StaticDir",
			).Call(jen.ID("frontendConfig").Dot(
				"StaticFilesDirectory",
			)),
			jen.If(
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("establishing static file server")),
			),
			jen.ID("router").Dot(
				"Get",
			).Call(jen.Lit("/*"), jen.ID("staticFileServer")),
		),
		jen.For(jen.List(jen.ID("route"), jen.ID("handler")).Op(":=").Range().ID("s").Dot(
			"frontendService",
		).Dot(
			"Routes",
		).Call()).Block(
			jen.ID("router").Dot(
				"Get",
			).Call(jen.ID("route"), jen.ID("handler")),
		),
		jen.ID("router").Dot(
			"With",
		).Call(jen.ID("s").Dot(
			"authService",
		).Dot(
			"AuthenticationMiddleware",
		).Call(jen.ID("true")), jen.ID("s").Dot(
			"authService",
		).Dot(
			"AdminMiddleware",
		)).Dot(
			"Route",
		).Call(jen.Lit("/admin"), jen.Func().Params(jen.ID("adminRouter").ID("chi").Dot(
			"Router",
		)).Block(
			jen.ID("adminRouter").Dot(
				"Post",
			).Call(jen.Lit("/cycle_cookie_secret"), jen.ID("s").Dot(
				"authService",
			).Dot(
				"CycleSecretHandler",
			).Call()),
		)),
		jen.ID("router").Dot(
			"Route",
		).Call(jen.Lit("/users"), jen.Func().Params(jen.ID("userRouter").ID("chi").Dot(
			"Router",
		)).Block(
			jen.ID("userRouter").Dot(
				"With",
			).Call(jen.ID("s").Dot(
				"authService",
			).Dot(
				"UserLoginInputMiddleware",
			)).Dot(
				"Post",
			).Call(jen.Lit("/login"), jen.ID("s").Dot(
				"authService",
			).Dot(
				"LoginHandler",
			).Call()),
			jen.ID("userRouter").Dot(
				"With",
			).Call(jen.ID("s").Dot(
				"authService",
			).Dot(
				"CookieAuthenticationMiddleware",
			)).Dot(
				"Post",
			).Call(jen.Lit("/logout"), jen.ID("s").Dot(
				"authService",
			).Dot(
				"LogoutHandler",
			).Call()),
			jen.ID("userIDPattern").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("oauth2IDPattern"), jen.ID("users").Dot(
				"URIParamKey",
			)),
			jen.ID("userRouter").Dot(
				"Get",
			).Call(jen.Lit("/"), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"ListHandler",
			).Call()),
			jen.ID("userRouter").Dot(
				"With",
			).Call(jen.ID("s").Dot(
				"usersService",
			).Dot(
				"UserInputMiddleware",
			)).Dot(
				"Post",
			).Call(jen.Lit("/"), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"CreateHandler",
			).Call()),
			jen.ID("userRouter").Dot(
				"Get",
			).Call(jen.ID("userIDPattern"), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"ReadHandler",
			).Call()),
			jen.ID("userRouter").Dot(
				"Delete",
			).Call(jen.ID("userIDPattern"), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"ArchiveHandler",
			).Call()),
			jen.ID("userRouter").Dot(
				"With",
			).Call(jen.ID("s").Dot(
				"authService",
			).Dot(
				"CookieAuthenticationMiddleware",
			), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"TOTPSecretRefreshInputMiddleware",
			)).Dot(
				"Post",
			).Call(jen.Lit("/totp_secret/new"), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"NewTOTPSecretHandler",
			).Call()),
			jen.ID("userRouter").Dot(
				"With",
			).Call(jen.ID("s").Dot(
				"authService",
			).Dot(
				"CookieAuthenticationMiddleware",
			), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"PasswordUpdateInputMiddleware",
			)).Dot(
				"Put",
			).Call(jen.Lit("/password/new"), jen.ID("s").Dot(
				"usersService",
			).Dot(
				"UpdatePasswordHandler",
			).Call()),
		)),
		jen.ID("router").Dot(
			"Route",
		).Call(jen.Lit("/oauth2"), jen.Func().Params(jen.ID("oauth2Router").ID("chi").Dot(
			"Router",
		)).Block(
			jen.ID("oauth2Router").Dot(
				"With",
			).Call(jen.ID("s").Dot(
				"authService",
			).Dot(
				"CookieAuthenticationMiddleware",
			), jen.ID("s").Dot(
				"oauth2ClientsService",
			).Dot(
				"CreationInputMiddleware",
			)).Dot(
				"Post",
			).Call(jen.Lit("/client"), jen.ID("s").Dot(
				"oauth2ClientsService",
			).Dot(
				"CreateHandler",
			).Call()),
			jen.ID("oauth2Router").Dot(
				"With",
			).Call(jen.ID("s").Dot(
				"oauth2ClientsService",
			).Dot(
				"OAuth2ClientInfoMiddleware",
			)).Dot(
				"Post",
			).Call(jen.Lit("/authorize"), jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"WithRequest",
				).Call(jen.ID("req")).Dot(
					"Debug",
				).Call(jen.Lit("oauth2 authorize route hit")),
				jen.If(
					jen.ID("err").Op(":=").ID("s").Dot(
						"oauth2ClientsService",
					).Dot(
						"HandleAuthorizeRequest",
					).Call(jen.ID("res"), jen.ID("req")),
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.ID("err").Dot(
						"Error",
					).Call(), jen.Qual("net/http", "StatusBadRequest")),
				),
			)),
			jen.ID("oauth2Router").Dot(
				"Post",
			).Call(jen.Lit("/token"), jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.If(
					jen.ID("err").Op(":=").ID("s").Dot(
						"oauth2ClientsService",
					).Dot(
						"HandleTokenRequest",
					).Call(jen.ID("res"), jen.ID("req")),
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.ID("err").Dot(
						"Error",
					).Call(), jen.Qual("net/http", "StatusBadRequest")),
				),
			)),
		)),
		jen.ID("router").Dot(
			"With",
		).Call(jen.ID("s").Dot(
			"authService",
		).Dot(
			"AuthenticationMiddleware",
		).Call(jen.ID("true"))).Dot(
			"Route",
		).Call(jen.Lit("/api/v1"), jen.Func().Params(jen.ID("v1Router").ID("chi").Dot(
			"Router",
		)).Block(
			jen.ID("v1Router").Dot(
				"Route",
			).Call(jen.Lit("/items"), jen.Func().Params(jen.ID("itemsRouter").ID("chi").Dot(
				"Router",
			)).Block(
				jen.ID("sr").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.ID("items").Dot(
					"URIParamKey",
				)),
				jen.ID("itemsRouter").Dot(
					"With",
				).Call(jen.ID("s").Dot(
					"itemsService",
				).Dot(
					"CreationInputMiddleware",
				)).Dot(
					"Post",
				).Call(jen.Lit("/"), jen.ID("s").Dot(
					"itemsService",
				).Dot(
					"CreateHandler",
				).Call()),
				jen.ID("itemsRouter").Dot(
					"Get",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"itemsService",
				).Dot(
					"ReadHandler",
				).Call()),
				jen.ID("itemsRouter").Dot(
					"With",
				).Call(jen.ID("s").Dot(
					"itemsService",
				).Dot(
					"UpdateInputMiddleware",
				)).Dot(
					"Put",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"itemsService",
				).Dot(
					"UpdateHandler",
				).Call()),
				jen.ID("itemsRouter").Dot(
					"Delete",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"itemsService",
				).Dot(
					"ArchiveHandler",
				).Call()),
				jen.ID("itemsRouter").Dot(
					"Get",
				).Call(jen.Lit("/"), jen.ID("s").Dot(
					"itemsService",
				).Dot(
					"ListHandler",
				).Call()),
			)),
			jen.ID("v1Router").Dot(
				"Route",
			).Call(jen.Lit("/webhooks"), jen.Func().Params(jen.ID("webhookRouter").ID("chi").Dot(
				"Router",
			)).Block(
				jen.ID("sr").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.ID("webhooks").Dot(
					"URIParamKey",
				)),
				jen.ID("webhookRouter").Dot(
					"With",
				).Call(jen.ID("s").Dot(
					"webhooksService",
				).Dot(
					"CreationInputMiddleware",
				)).Dot(
					"Post",
				).Call(jen.Lit("/"), jen.ID("s").Dot(
					"webhooksService",
				).Dot(
					"CreateHandler",
				).Call()),
				jen.ID("webhookRouter").Dot(
					"Get",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"webhooksService",
				).Dot(
					"ReadHandler",
				).Call()),
				jen.ID("webhookRouter").Dot(
					"With",
				).Call(jen.ID("s").Dot(
					"webhooksService",
				).Dot(
					"UpdateInputMiddleware",
				)).Dot(
					"Put",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"webhooksService",
				).Dot(
					"UpdateHandler",
				).Call()),
				jen.ID("webhookRouter").Dot(
					"Delete",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"webhooksService",
				).Dot(
					"ArchiveHandler",
				).Call()),
				jen.ID("webhookRouter").Dot(
					"Get",
				).Call(jen.Lit("/"), jen.ID("s").Dot(
					"webhooksService",
				).Dot(
					"ListHandler",
				).Call()),
			)),
			jen.ID("v1Router").Dot(
				"Route",
			).Call(jen.Lit("/oauth2/clients"), jen.Func().Params(jen.ID("clientRouter").ID("chi").Dot(
				"Router",
			)).Block(
				jen.ID("sr").Op(":=").Qual("fmt", "Sprintf").Call(jen.ID("numericIDPattern"), jen.ID("oauth2clients").Dot(
					"URIParamKey",
				)),
				jen.ID("clientRouter").Dot(
					"Get",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"oauth2ClientsService",
				).Dot(
					"ReadHandler",
				).Call()),
				jen.ID("clientRouter").Dot(
					"Delete",
				).Call(jen.ID("sr"), jen.ID("s").Dot(
					"oauth2ClientsService",
				).Dot(
					"ArchiveHandler",
				).Call()),
				jen.ID("clientRouter").Dot(
					"Get",
				).Call(jen.Lit("/"), jen.ID("s").Dot(
					"oauth2ClientsService",
				).Dot(
					"ListHandler",
				).Call()),
			)),
		)),
		jen.ID("s").Dot(
			"router",
		).Op("=").ID("router"),
	),
	)
	return ret
}
