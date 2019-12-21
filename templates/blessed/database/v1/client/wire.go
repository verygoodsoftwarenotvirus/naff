package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what we provide to dependency injectors"),
			jen.ID("Providers").Op("=").Qual("github.com/google/wire", "NewSet".Callln(
				jen.ID("ProvideDatabaseClient"),
			),
		),
		jen.Line(),
	)
	return ret
}
