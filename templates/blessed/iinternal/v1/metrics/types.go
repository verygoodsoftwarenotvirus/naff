package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func typesDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Namespace is a string alias for dependency injection's sake"),
			jen.ID("Namespace").String(),
			jen.Line(),
			jen.Comment("CounterName is a string alias for dependency injection's sake"),
			jen.ID("CounterName").String(),
			jen.Line(),
			jen.Comment("SpanFormatter formats the name of a span given a request"),
			jen.ID("SpanFormatter").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()),
			jen.Line(),
			jen.Comment("InstrumentationHandler is an obligatory alias"),
			jen.ID("InstrumentationHandler").Qual("net/http", "Handler"),
			jen.Line(),
			jen.Comment("Handler is the Handler that provides metrics data to scraping services"),
			jen.ID("Handler").Qual("net/http", "Handler"),
			jen.Line(),
			jen.Comment("HandlerInstrumentationFunc blah"), // LOL
			jen.ID("HandlerInstrumentationFunc").Func().Params(jen.Qual("net/http", "HandlerFunc")).Params(jen.Qual("net/http", "HandlerFunc")),
			jen.Line(),
			jen.Comment("UnitCounter describes a counting interface for things like total user counts"),
			jen.Comment("Meant to handle integers exclusively"),
			jen.ID("UnitCounter").Interface(
				jen.ID("Increment").Params(constants.CtxParam()),
				jen.ID("IncrementBy").Params(constants.CtxParam(), jen.ID("val").Uint64()),
				jen.ID("Decrement").Params(constants.CtxParam()),
			),
			jen.Line(),
			jen.Comment("UnitCounterProvider is a function that provides a UnitCounter and an error"),
			jen.ID("UnitCounterProvider").Func().Params(jen.ID("counterName").ID("CounterName"), jen.ID("description").String()).Params(jen.ID("UnitCounter"), jen.Error()),
			jen.Line(),
		),
	)
	return ret
}
