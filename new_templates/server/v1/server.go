package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func serverDotGo() *jen.File {
	ret := jen.NewFile("server")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Type().ID("Server").Struct(
		jen.ID("config").Op("*").ID("config").Dot(
			"ServerConfig",
		),
		jen.ID("httpServer").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "Server"),
	),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideServer")),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
