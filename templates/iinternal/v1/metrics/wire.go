package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("metrics")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents what this library offers to external users in the form of dependencies."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideUnitCounter"),
				jen.ID("ProvideUnitCounterProvider"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUnitCounterProvider provides UnitCounter providers."),
		jen.Line(),
		jen.Func().ID("ProvideUnitCounterProvider").Params().Params(jen.ID("UnitCounterProvider")).Body(
			jen.Return().ID("ProvideUnitCounter"),
		),
		jen.Line(),
	)

	return code
}
