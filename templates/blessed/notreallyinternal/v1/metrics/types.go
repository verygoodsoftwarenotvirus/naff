package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func typesDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Namespace is a string alias for dependency injection's sake"),
			jen.ID("Namespace").ID("string"),
			jen.Line(),
			jen.Comment("CounterName is a string alias for dependency injection's sake"),
			jen.ID("CounterName").ID("string"),
			jen.Line(),
			jen.Comment("SpanFormatter formats the name of a span given a request"),
			jen.ID("SpanFormatter").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")),
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
				jen.ID("Increment").Params(jen.ID("ctx").Qual("context", "Context")),
				jen.ID("IncrementBy").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("val").ID("uint64")),
				jen.ID("Decrement").Params(jen.ID("ctx").Qual("context", "Context")),
			),
			jen.Comment("UnitCounterProvider is a function that provides a UnitCounter and an error"),
			jen.ID("UnitCounterProvider").Func().Params(
				jen.ID("counterName").ID("CounterName"),
				jen.ID("description").ID("string"),
			).Params(
				jen.ID("UnitCounter"),
				jen.ID("error"),
			),
			jen.Line(),
		),
	)
	return ret
}
