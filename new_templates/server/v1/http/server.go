package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func serverDotGo() *jen.File {
	ret := jen.NewFile("httpserver")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("maxTimeout").Op("=").Lit(120).Op("*").Qual("time", "Second"),
	)
	ret.Add(jen.Null().Type().ID("Server").Struct(
		jen.ID("DebugMode").ID("bool"),
		jen.ID("authService").Op("*").ID("auth").Dot(
			"Service",
		),
		jen.ID("frontendService").Op("*").ID("frontend").Dot(
			"Service",
		),
		jen.ID("usersService").ID("models").Dot(
			"UserDataServer",
		),
		jen.ID("oauth2ClientsService").ID("models").Dot(
			"OAuth2ClientDataServer",
		),
		jen.ID("webhooksService").ID("models").Dot(
			"WebhookDataServer",
		),
		jen.ID("itemsService").ID("models").Dot(
			"ItemDataServer",
		),
		jen.ID("db").ID("database").Dot(
			"Database",
		),
		jen.ID("config").ID("config").Dot(
			"ServerSettings",
		),
		jen.ID("router").Op("*").ID("chi").Dot(
			"Mux",
		),
		jen.ID("httpServer").Op("*").Qual("net/http", "Server"),
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("encoder").ID("encoding").Dot(
			"EncoderDecoder",
		),
		jen.ID("newsManager").Op("*").ID("newsman").Dot(
			"Newsman",
		),
	),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
