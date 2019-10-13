package http

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func serverDotGo() *jen.File {
	ret := jen.NewFile("httpserver")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("maxTimeout").Op("=").Lit(120).Op("*").Qual("time", "Second"),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("Server").Struct(jen.ID("DebugMode").ID("bool"), jen.ID("authService").Op("*").ID("auth").Dot(
		"Service",
	), jen.ID("frontendService").Op("*").ID("frontend").Dot(
		"Service",
	), jen.ID("usersService").ID("models").Dot(
		"UserDataServer",
	), jen.ID("oauth2ClientsService").ID("models").Dot(
		"OAuth2ClientDataServer",
	), jen.ID("webhooksService").ID("models").Dot(
		"WebhookDataServer",
	), jen.ID("itemsService").ID("models").Dot(
		"ItemDataServer",
	), jen.ID("db").ID("database").Dot(
		"Database",
	), jen.ID("config").Op("*").ID("config").Dot(
		"ServerConfig",
	), jen.ID("router").Op("*").ID("chi").Dot(
		"Mux",
	), jen.ID("httpServer").Op("*").Qual("net/http", "Server"), jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("encoder").ID("encoding").Dot(
		"EncoderDecoder",
	), jen.ID("newsManager").Op("*").ID("newsman").Dot(
		"Newsman",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideServer builds a new Server instance").ID("ProvideServer").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("config").Dot(
		"ServerConfig",
	), jen.ID("authService").Op("*").ID("auth").Dot(
		"Service",
	), jen.ID("frontendService").Op("*").ID("frontend").Dot(
		"Service",
	), jen.ID("itemsService").ID("models").Dot(
		"ItemDataServer",
	), jen.ID("usersService").ID("models").Dot(
		"UserDataServer",
	), jen.ID("oauth2Service").ID("models").Dot(
		"OAuth2ClientDataServer",
	), jen.ID("webhooksService").ID("models").Dot(
		"WebhookDataServer",
	), jen.ID("db").ID("database").Dot(
		"Database",
	), jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("encoder").ID("encoding").Dot(
		"EncoderDecoder",
	), jen.ID("newsManager").Op("*").ID("newsman").Dot(
		"Newsman",
	)).Params(jen.Op("*").ID("Server"), jen.ID("error")).Block(
		jen.If(jen.ID("len").Call(jen.ID("cfg").Dot(
			"Auth",
		).Dot(
			"CookieSecret",
		)).Op("<").Lit(32)).Block(
			jen.ID("err").Op(":=").ID("errors").Dot(
				"New",
			).Call(jen.Lit("cookie secret is too short, must be at least 32 characters in length")),
			jen.ID("logger").Dot(
				"Error",
			).Call(jen.ID("err"), jen.Lit("cookie secret failure")),
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.ID("srv").Op(":=").Op("&").ID("Server").Valuesln(jen.ID("DebugMode").Op(":").ID("cfg").Dot(
			"Server",
		).Dot(
			"Debug",
		), jen.ID("db").Op(":").ID("db"), jen.ID("config").Op(":").ID("cfg"), jen.ID("encoder").Op(":").ID("encoder"), jen.ID("httpServer").Op(":").ID("provideHTTPServer").Call(), jen.ID("logger").Op(":").ID("logger").Dot(
			"WithName",
		).Call(jen.Lit("api_server")), jen.ID("newsManager").Op(":").ID("newsManager"), jen.ID("webhooksService").Op(":").ID("webhooksService"), jen.ID("frontendService").Op(":").ID("frontendService"), jen.ID("usersService").Op(":").ID("usersService"), jen.ID("authService").Op(":").ID("authService"), jen.ID("itemsService").Op(":").ID("itemsService"), jen.ID("oauth2ClientsService").Op(":").ID("oauth2Service")),
		jen.If(jen.ID("err").Op(":=").ID("cfg").Dot(
			"ProvideTracing",
		).Call(jen.ID("logger")), jen.ID("err").Op("!=").ID("nil").Op("&&").ID("err").Op("!=").ID("config").Dot(
			"ErrInvalidTracingProvider",
		)).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.List(jen.ID("ih"), jen.ID("err")).Op(":=").ID("cfg").Dot(
			"ProvideInstrumentationHandler",
		).Call(jen.ID("logger")),
		jen.If(jen.ID("err").Op("!=").ID("nil").Op("&&").ID("err").Op("!=").ID("config").Dot(
			"ErrInvalidMetricsProvider",
		)).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.If(jen.ID("ih").Op("!=").ID("nil")).Block(
			jen.ID("srv").Dot(
				"setupRouter",
			).Call(jen.ID("cfg").Dot(
				"Frontend",
			), jen.ID("ih")),
		),
		jen.ID("srv").Dot(
			"httpServer",
		).Dot(
			"Handler",
		).Op("=").Op("&").ID("ochttp").Dot(
			"Handler",
		).Valuesln(jen.ID("Handler").Op(":").ID("srv").Dot(
			"router",
		), jen.ID("FormatSpanName").Op(":").ID("formatSpanNameForRequest")),
		jen.List(jen.ID("allWebhooks"), jen.ID("err")).Op(":=").ID("db").Dot(
			"GetAllWebhooks",
		).Call(jen.ID("ctx")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("initializing webhooks: %w"), jen.ID("err"))),
		),
		jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("len").Call(jen.ID("allWebhooks").Dot(
			"Webhooks",
		)), jen.ID("i").Op("++")).Block(
			jen.ID("wh").Op(":=").ID("allWebhooks").Dot(
				"Webhooks",
			).Index(jen.ID("i")),
			jen.ID("l").Op(":=").ID("wh").Dot(
				"ToListener",
			).Call(jen.ID("srv").Dot(
				"logger",
			)),
			jen.ID("srv").Dot(
				"newsManager",
			).Dot(
				"TuneIn",
			).Call(jen.ID("l")),
		),
		jen.Return().List(jen.ID("srv"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// Serve serves HTTP traffic").Params(jen.ID("s").Op("*").ID("Server")).ID("Serve").Params().Block(
		jen.ID("s").Dot(
			"httpServer",
		).Dot(
			"Addr",
		).Op("=").Qual("fmt", "Sprintf").Call(jen.Lit(":%d"), jen.ID("s").Dot(
			"config",
		).Dot(
			"Server",
		).Dot(
			"HTTPPort",
		)),
		jen.ID("s").Dot(
			"logger",
		).Dot(
			"Debug",
		).Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("Listening for HTTP requests on %q"), jen.ID("s").Dot(
			"httpServer",
		).Dot(
			"Addr",
		))),
		jen.If(jen.ID("err").Op(":=").ID("s").Dot(
			"httpServer",
		).Dot(
			"ListenAndServe",
		).Call(), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"Error",
			).Call(jen.ID("err"), jen.Lit("server shutting down")),
			jen.If(jen.ID("err").Op("==").Qual("net/http", "ErrServerClosed")).Block(
				jen.Qual("os", "Exit").Call(jen.Lit(0)),
			),
		),
	),

		jen.Line(),
	)
	return ret
}
