package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is our collection of what we provide to other services"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.ID("ProvideFrontendService"),
			),
		),
		jen.Line(),
	)
	return ret
}
