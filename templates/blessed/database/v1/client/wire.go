package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkgRoot, types, ret)

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
