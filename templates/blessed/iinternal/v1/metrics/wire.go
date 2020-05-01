package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what this library offers to external users in the form of dependencies."),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID("ProvideUnitCounter"),
				jen.ID("ProvideUnitCounterProvider"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUnitCounterProvider provides UnitCounter providers."),
		jen.Line(),
		jen.Func().ID("ProvideUnitCounterProvider").Params().Params(jen.ID("UnitCounterProvider")).Block(
			jen.Return().ID("ProvideUnitCounter"),
		),
		jen.Line(),
	)

	return ret
}
