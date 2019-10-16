package http

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func serverTestDotGo() *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildTestServer").Params().Params(jen.Op("*").ID("Server")).Block(
		jen.ID("s").Op(":=").Op("&").ID("Server").Valuesln(jen.ID("DebugMode").Op(":").ID("true"), jen.ID("db").Op(":").ID("database").Dot(
			"BuildMockDatabase",
		).Call(), jen.ID("config").Op(":").Op("&").ID("config").Dot(
			"ServerConfig",
		).Valuesln(), jen.ID("encoder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(), jen.ID("httpServer").Op(":").ID("provideHTTPServer").Call(), jen.ID("logger").Op(":").ID("noop").Dot(
			"ProvideNoopLogger",
		).Call(), jen.ID("frontendService").Op(":").ID("frontend").Dot(
			"ProvideFrontendService",
		).Call(jen.ID("noop").Dot(
			"ProvideNoopLogger",
		).Call(), jen.ID("config").Dot(
			"FrontendSettings",
		).Valuesln()), jen.ID("webhooksService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataServer").Valuesln(), jen.ID("usersService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataServer").Valuesln(), jen.ID("authService").Op(":").Op("&").ID("auth").Dot(
			"Service",
		).Valuesln(), jen.ID("itemsService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataServer").Valuesln(), jen.ID("oauth2ClientsService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "OAuth2ClientDataServer").Valuesln()),
		jen.Return().ID("s"),
	),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideServer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"WebhookDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllWebhooks"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"WebhookList",
			).Valuesln(), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideServer").Call(jen.Qual("context", "Background").Call(), jen.Op("&").ID("config").Dot(
				"ServerConfig",
			).Valuesln(jen.ID("Auth").Op(":").ID("config").Dot(
				"AuthSettings",
			).Valuesln(jen.ID("CookieSecret").Op(":").Lit("THISISAVERYLONGSTRINGFORTESTPURPOSES"))), jen.Op("&").ID("auth").Dot(
				"Service",
			).Valuesln(), jen.Op("&").ID("frontend").Dot(
				"Service",
			).Valuesln(), jen.Op("&").ID("items").Dot(
				"Service",
			).Valuesln(), jen.Op("&").ID("users").Dot(
				"Service",
			).Valuesln(), jen.Op("&").ID("oauth2clients").Dot(
				"Service",
			).Valuesln(), jen.Op("&").ID("webhooks").Dot(
				"Service",
			).Valuesln(), jen.ID("mockDB"), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(), jen.ID("newsman").Dot(
				"NewNewsman",
			).Call(jen.ID("nil"), jen.ID("nil"))),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),
		jen.Line(),
	)
	return ret
}
