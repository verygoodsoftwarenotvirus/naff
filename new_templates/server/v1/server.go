package server

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func serverDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

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
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
