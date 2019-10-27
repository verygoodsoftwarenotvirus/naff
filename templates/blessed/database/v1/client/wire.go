package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what we provide to dependency injectors"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.ID("ProvideDatabaseClient"),
			),
		),
		jen.Line(),
	)
	return ret
}
