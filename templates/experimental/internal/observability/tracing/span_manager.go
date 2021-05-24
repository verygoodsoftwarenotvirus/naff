package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("Tracer").Op("=").Parens(jen.Op("*").ID("otelSpanManager")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("otelSpanManager").Struct(jen.ID("tracer").ID("trace").Dot("Tracer")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewTracer creates a Tracer."),
		jen.Line(),
		jen.Func().ID("NewTracer").Params(jen.ID("name").ID("string")).Params(jen.ID("Tracer")).Body(
			jen.Return().Op("&").ID("otelSpanManager").Valuesln(jen.ID("tracer").Op(":").Qual("go.opentelemetry.io/otel", "Tracer").Call(jen.ID("name")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("StartSpan wraps tracer.Start."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("otelSpanManager")).ID("StartSpan").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Qual("context", "Context"), jen.ID("Span")).Body(
			jen.Return().ID("t").Dot("tracer").Dot("Start").Call(
				jen.ID("ctx"),
				jen.ID("GetCallerName").Call(),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("StartCustomSpan wraps tracer.Start."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("otelSpanManager")).ID("StartCustomSpan").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("name").ID("string")).Params(jen.Qual("context", "Context"), jen.ID("Span")).Body(
			jen.Return().ID("t").Dot("tracer").Dot("Start").Call(
				jen.ID("ctx"),
				jen.ID("name"),
			)),
		jen.Line(),
	)

	return code
}
