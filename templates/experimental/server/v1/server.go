package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func serverDotGo() *jen.File {
	ret := jen.NewFile("server")

	utils.AddImports(ret)

	ret.Add(
		jen.Type().ID("Server").Struct(
			jen.ID("config").Op("*").ID("config").Dot("ServerConfig"),
			jen.ID("httpServer").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "Server"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(jen.ID("ProvideServer")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideServer builds a new Server instance"),
		jen.Line(),
		jen.Func().ID("ProvideServer").Params(jen.ID("cfg").Op("*").ID("config").Dot("ServerConfig"), jen.ID("httpServer").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "Server")).Params(jen.Op("*").ID("Server"), jen.ID("error")).Block(
			jen.ID("srv").Op(":=").Op("&").ID("Server").Valuesln(
				jen.ID("config").Op(":").ID("cfg"),
				jen.ID("httpServer").Op(":").ID("httpServer"),
			),
			jen.Return().List(jen.ID("srv"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Serve serves HTTP traffic"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Server")).ID("Serve").Params().Block(
			jen.ID("s").Dot("httpServer").Dot("Serve").Call(),
		),
		jen.Line(),
	)
	return ret
}
