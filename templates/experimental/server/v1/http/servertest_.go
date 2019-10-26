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
			jen.ID("s").Op(":=").Op("&").ID("Server").Valuesln(
				jen.ID("DebugMode").Op(":").ID("true"),
				jen.ID("db").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "BuildMockDatabase").Call(),
				jen.ID("config").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "ServerConfig").Values(),
				jen.ID("encoder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
				jen.ID("httpServer").Op(":").ID("provideHTTPServer").Call(),
				jen.ID("logger").Op(":").ID("noop").Dot("ProvideNoopLogger").Call(),
				jen.ID("frontendService").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/frontend", "ProvideFrontendService").Callln(
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "FrontendSettings").Values(),
				),
				jen.ID("webhooksService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "WebhookDataServer").Values(),
				jen.ID("usersService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataServer").Values(),
				jen.ID("authService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth", "Service").Values(),
				jen.ID("itemsService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "ItemDataServer").Values(),
				jen.ID("oauth2ClientsService").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "OAuth2ClientDataServer").Values(),
			),
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideServer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("mockDB").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(jen.Lit("GetAllWebhooks"), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "WebhookList").Values(), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideServer").Callln(
					jen.Qual("context", "Background").Call(),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "ServerConfig").Valuesln(
						jen.ID("Auth").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "AuthSettings").Valuesln(
							jen.ID("CookieSecret").Op(":").Lit("THISISAVERYLONGSTRINGFORTESTPURPOSES"),
						),
					),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth", "Service").Values(),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/frontend", "Service").Values(),
					jen.Op("&").ID("items").Dot("Service").Values(),
					jen.Op("&").ID("users").Dot("Service").Values(),
					jen.Op("&").ID("oauth2clients").Dot("Service").Values(),
					jen.Op("&").ID("webhooks").Dot("Service").Values(),
					jen.ID("mockDB"), jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
					jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.ID("nil"), jen.ID("nil"))),
				jen.Line(),
				jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)
	return ret
}
