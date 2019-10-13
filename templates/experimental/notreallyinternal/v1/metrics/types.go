package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func typesDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("Namespace").ID("string").Type().ID("CounterName").ID("string").Type().ID("SpanFormatter").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Type().ID("InstrumentationHandler").Qual("net/http", "Handler").Type().ID("Handler").Qual("net/http", "Handler").Type().ID("HandlerInstrumentationFunc").Params(jen.Qual("net/http", "HandlerFunc")).Params(jen.Qual("net/http", "HandlerFunc")).Type().ID("UnitCounter").Interface(jen.ID("Increment").Params(jen.ID("ctx").Qual("context", "Context")), jen.ID("IncrementBy").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("val").ID("uint64")), jen.ID("Decrement").Params(jen.ID("ctx").Qual("context", "Context"))).Type().ID("UnitCounterProvider").Params(jen.ID("counterName").ID("CounterName"), jen.ID("description").ID("string")).Params(jen.ID("UnitCounter"), jen.ID("error")),

		jen.Line(),
	)
	return ret
}
