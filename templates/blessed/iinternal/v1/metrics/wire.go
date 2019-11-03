package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what this library offers to external users in the form of dependencies"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.ID("ProvideUnitCounter"),
				jen.ID("ProvideUnitCounterProvider"),
			),
		),
		jen.Line(),
	)
	return ret
}
