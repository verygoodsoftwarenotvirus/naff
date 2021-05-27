package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(
				jen.ID("ProvideUnitCounterProvider"),
				jen.ID("ProvideMetricsInstrumentationHandlerForServer"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideMetricsInstrumentationHandlerForServer provides a metrics.InstrumentationHandler from a config for our server."),
		jen.Line(),
		jen.Func().ID("ProvideMetricsInstrumentationHandlerForServer").Params(jen.ID("cfg").Op("*").ID("Config"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("InstrumentationHandler"), jen.ID("error")).Body(
			jen.Return().ID("cfg").Dot("ProvideInstrumentationHandler").Call(jen.ID("logger"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUnitCounterProvider provides a metrics.InstrumentationHandler from a config for our server."),
		jen.Line(),
		jen.Func().ID("ProvideUnitCounterProvider").Params(jen.ID("cfg").Op("*").ID("Config"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("UnitCounterProvider"), jen.ID("error")).Body(
			jen.Return().ID("cfg").Dot("ProvideUnitCounterProvider").Call(jen.ID("logger"))),
		jen.Line(),
	)

	return code
}
