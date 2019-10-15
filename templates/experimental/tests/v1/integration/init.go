package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func initDotGo() *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("debug").Op("=").ID("true").Var().ID("nonexistentID").Op("=").Lit(999999999),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("urlToUse").ID("string").Var().ID("todoClient").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "V1Client"),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("init").Params().Block(
		jen.ID("urlToUse").Op("=").ID("testutil").Dot(
			"DetermineServiceURL",
		).Call(),
		jen.ID("logger").Op(":=").ID("zerolog").Dot(
			"NewZeroLogger",
		).Call(),
		jen.ID("logger").Dot(
			"WithValue",
		).Call(jen.Lit("url"), jen.ID("urlToUse")).Dot(
			"Info",
		).Call(jen.Lit("checking server")),
		jen.ID("testutil").Dot(
			"EnsureServerIsUp",
		).Call(jen.ID("urlToUse")),
		jen.ID("fake").Dot(
			"Seed",
		).Call(jen.Qual("time", "Now").Call().Dot(
			"UnixNano",
		).Call()),
		jen.List(jen.ID("ogUser"), jen.ID("err")).Op(":=").ID("testutil").Dot(
			"CreateObligatoryUser",
		).Call(jen.ID("urlToUse"), jen.ID("debug")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot(
				"Fatal",
			).Call(jen.ID("err")),
		),
		jen.List(jen.ID("oa2Client"), jen.ID("err")).Op(":=").ID("testutil").Dot(
			"CreateObligatoryClient",
		).Call(jen.ID("urlToUse"), jen.ID("ogUser")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("logger").Dot(
				"Fatal",
			).Call(jen.ID("err")),
		),
		jen.ID("todoClient").Op("=").ID("initializeClient").Call(jen.ID("oa2Client")),
		jen.ID("fiftySpaces").Op(":=").Qual("strings", "Repeat").Call(jen.Lit("\n"), jen.Lit(50)),
		jen.Qual("fmt", "Printf").Call(jen.Lit("%s\tRunning tests%s"), jen.ID("fiftySpaces"), jen.ID("fiftySpaces")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("buildHTTPClient").Params().Params(jen.Op("*").Qual("net/http", "Client")).Block(
		jen.ID("httpc").Op(":=").Op("&").Qual("net/http", "Client").Valuesln(jen.ID("Transport").Op(":").Qual("net/http", "DefaultTransport"), jen.ID("Timeout").Op(":").Lit(5).Op("*").Qual("time", "Second")),
		jen.Return().ID("httpc"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("initializeClient").Params(jen.ID("oa2Client").Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "V1Client")).Block(
		jen.List(jen.ID("uri"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("urlToUse")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
		jen.List(jen.ID("c"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "NewClient").Call(jen.Qual("context", "Background").Call(), jen.ID("oa2Client").Dot(
			"ClientID",
		), jen.ID("oa2Client").Dot(
			"ClientSecret",
		), jen.ID("uri"), jen.ID("zerolog").Dot(
			"NewZeroLogger",
		).Call(), jen.ID("buildHTTPClient").Call(), jen.ID("oa2Client").Dot(
			"Scopes",
		), jen.ID("debug")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
		jen.Return().ID("c"),
	),

		jen.Line(),
	)
	return ret
}