package v1

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireGenDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("BuildServer").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("config").Dot(
		"ServerConfig",
	), jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("database2").ID("database").Dot(
		"Database",
	)).Params(jen.Op("*").ID("server").Dot(
		"Server",
	), jen.ID("error")).Block(
		jen.ID("bcryptHashCost").Op(":=").ID("auth").Dot(
			"ProvideBcryptHashCost",
		).Call(),
		jen.ID("authenticator").Op(":=").ID("auth").Dot(
			"ProvideBcryptAuthenticator",
		).Call(jen.ID("bcryptHashCost"), jen.ID("logger")),
		jen.ID("userDataManager").Op(":=").ID("users").Dot(
			"ProvideUserDataManager",
		).Call(jen.ID("database2")),
		jen.ID("clientIDFetcher").Op(":=").ID("httpserver").Dot(
			"ProvideOAuth2ServiceClientIDFetcher",
		).Call(jen.ID("logger")),
		jen.ID("encoderDecoder").Op(":=").ID("encoding").Dot(
			"ProvideResponseEncoder",
		).Call(),
		jen.ID("unitCounterProvider").Op(":=").ID("metrics").Dot(
			"ProvideUnitCounterProvider",
		).Call(),
		jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("oauth2clients").Dot(
			"ProvideOAuth2ClientsService",
		).Call(jen.ID("ctx"), jen.ID("logger"), jen.ID("database2"), jen.ID("authenticator"), jen.ID("clientIDFetcher"), jen.ID("encoderDecoder"), jen.ID("unitCounterProvider")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.ID("oAuth2ClientValidator").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth", "ProvideOAuth2ClientValidator").Call(jen.ID("service")),
		jen.ID("userIDFetcher").Op(":=").ID("httpserver").Dot(
			"ProvideAuthUserIDFetcher",
		).Call(),
		jen.ID("authService").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth", "ProvideAuthService").Call(jen.ID("logger"), jen.ID("cfg"), jen.ID("authenticator"), jen.ID("userDataManager"), jen.ID("oAuth2ClientValidator"), jen.ID("userIDFetcher"), jen.ID("encoderDecoder")),
		jen.ID("frontendSettings").Op(":=").ID("config").Dot(
			"ProvideConfigFrontendSettings",
		).Call(jen.ID("cfg")),
		jen.ID("frontendService").Op(":=").ID("frontend").Dot(
			"ProvideFrontendService",
		).Call(jen.ID("logger"), jen.ID("frontendSettings")),
		jen.ID("itemDataManager").Op(":=").ID("items").Dot(
			"ProvideItemDataManager",
		).Call(jen.ID("database2")),
		jen.ID("itemsUserIDFetcher").Op(":=").ID("httpserver").Dot(
			"ProvideUserIDFetcher",
		).Call(),
		jen.ID("itemIDFetcher").Op(":=").ID("httpserver").Dot(
			"ProvideItemIDFetcher",
		).Call(jen.ID("logger")),
		jen.ID("websocketAuthFunc").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth", "ProvideWebsocketAuthFunc").Call(jen.ID("authService")),
		jen.ID("typeNameManipulationFunc").Op(":=").ID("httpserver").Dot(
			"ProvideNewsmanTypeNameManipulationFunc",
		).Call(jen.ID("logger")),
		jen.ID("newsmanNewsman").Op(":=").ID("newsman").Dot(
			"NewNewsman",
		).Call(jen.ID("websocketAuthFunc"), jen.ID("typeNameManipulationFunc")),
		jen.ID("reporter").Op(":=").ID("ProvideReporter").Call(jen.ID("newsmanNewsman")),
		jen.List(jen.ID("itemsService"), jen.ID("err")).Op(":=").ID("items").Dot(
			"ProvideItemsService",
		).Call(jen.ID("ctx"), jen.ID("logger"), jen.ID("itemDataManager"), jen.ID("itemsUserIDFetcher"), jen.ID("itemIDFetcher"), jen.ID("encoderDecoder"), jen.ID("unitCounterProvider"), jen.ID("reporter")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.ID("itemDataServer").Op(":=").ID("items").Dot(
			"ProvideItemDataServer",
		).Call(jen.ID("itemsService")),
		jen.ID("authSettings").Op(":=").ID("config").Dot(
			"ProvideConfigAuthSettings",
		).Call(jen.ID("cfg")),
		jen.ID("usersUserIDFetcher").Op(":=").ID("httpserver").Dot(
			"ProvideUsernameFetcher",
		).Call(jen.ID("logger")),
		jen.List(jen.ID("usersService"), jen.ID("err")).Op(":=").ID("users").Dot(
			"ProvideUsersService",
		).Call(jen.ID("ctx"), jen.ID("authSettings"), jen.ID("logger"), jen.ID("database2"), jen.ID("authenticator"), jen.ID("usersUserIDFetcher"), jen.ID("encoderDecoder"), jen.ID("unitCounterProvider"), jen.ID("reporter")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.ID("userDataServer").Op(":=").ID("users").Dot(
			"ProvideUserDataServer",
		).Call(jen.ID("usersService")),
		jen.ID("oAuth2ClientDataServer").Op(":=").ID("oauth2clients").Dot(
			"ProvideOAuth2ClientDataServer",
		).Call(jen.ID("service")),
		jen.ID("webhookDataManager").Op(":=").ID("webhooks").Dot(
			"ProvideWebhookDataManager",
		).Call(jen.ID("database2")),
		jen.ID("webhooksUserIDFetcher").Op(":=").ID("httpserver").Dot(
			"ProvideWebhooksUserIDFetcher",
		).Call(),
		jen.ID("webhookIDFetcher").Op(":=").ID("httpserver").Dot(
			"ProvideWebhookIDFetcher",
		).Call(jen.ID("logger")),
		jen.List(jen.ID("webhooksService"), jen.ID("err")).Op(":=").ID("webhooks").Dot(
			"ProvideWebhooksService",
		).Call(jen.ID("ctx"), jen.ID("logger"), jen.ID("webhookDataManager"), jen.ID("webhooksUserIDFetcher"), jen.ID("webhookIDFetcher"), jen.ID("encoderDecoder"), jen.ID("unitCounterProvider"), jen.ID("newsmanNewsman")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.ID("webhookDataServer").Op(":=").ID("webhooks").Dot(
			"ProvideWebhookDataServer",
		).Call(jen.ID("webhooksService")),
		jen.List(jen.ID("httpserverServer"), jen.ID("err")).Op(":=").ID("httpserver").Dot(
			"ProvideServer",
		).Call(jen.ID("ctx"), jen.ID("cfg"), jen.ID("authService"), jen.ID("frontendService"), jen.ID("itemDataServer"), jen.ID("userDataServer"), jen.ID("oAuth2ClientDataServer"), jen.ID("webhookDataServer"), jen.ID("database2"), jen.ID("logger"), jen.ID("encoderDecoder"), jen.ID("newsmanNewsman")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.List(jen.ID("serverServer"), jen.ID("err")).Op(":=").ID("server").Dot(
			"ProvideServer",
		).Call(jen.ID("cfg"), jen.ID("httpserverServer")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.Return().List(jen.ID("serverServer"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day").ID("ProvideReporter").Params(jen.ID("n").Op("*").ID("newsman").Dot(
		"Newsman",
	)).Params(jen.ID("newsman").Dot(
		"Reporter",
	)).Block(
		jen.Return().ID("n"),
	),

		jen.Line(),
	)
	return ret
}
