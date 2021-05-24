package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func typesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Null().Type().ID("Namespace").ID("string").Type().ID("CounterName").ID("string").Type().ID("SpanFormatter").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Type().ID("InstrumentationHandler").Qual("net/http", "Handler").Type().ID("Handler").Qual("net/http", "Handler").Type().ID("HandlerInstrumentationFunc").Params(jen.Qual("net/http", "HandlerFunc")).Params(jen.Qual("net/http", "HandlerFunc")).Type().ID("UnitCounter").Interface(
			jen.ID("Increment").Params(jen.ID("ctx").Qual("context", "Context")),
			jen.ID("IncrementBy").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("val").ID("int64")),
			jen.ID("Decrement").Params(jen.ID("ctx").Qual("context", "Context")),
		).Type().ID("UnitCounterProvider").Params(jen.List(jen.ID("name"), jen.ID("description")).ID("string")).Params(jen.ID("UnitCounter")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EnsureUnitCounter always provides a valid UnitCounter."),
		jen.Line(),
		jen.Func().ID("EnsureUnitCounter").Params(jen.ID("ucp").ID("UnitCounterProvider"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("counterName").ID("CounterName"), jen.ID("description").ID("string")).Params(jen.ID("UnitCounter")).Body(
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("counter"),
				jen.ID("counterName"),
			),
			jen.If(jen.ID("ucp").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("building unit counter")),
				jen.Return().ID("ucp").Call(
					jen.ID("string").Call(jen.ID("counterName")),
					jen.ID("description"),
				),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("returning noop counter")),
			jen.Return().Op("&").ID("noopUnitCounter").Valuesln(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("UnitCounter").Op("=").Parens(jen.Op("*").ID("noopUnitCounter")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("noopUnitCounter").Struct(),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("noopUnitCounter")).ID("Increment").Params(jen.ID("_").Qual("context", "Context")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("noopUnitCounter")).ID("IncrementBy").Params(jen.ID("_").Qual("context", "Context"), jen.ID("_").ID("int64")).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("noopUnitCounter")).ID("Decrement").Params(jen.ID("_").Qual("context", "Context")).Body(),
		jen.Line(),
	)

	return code
}
