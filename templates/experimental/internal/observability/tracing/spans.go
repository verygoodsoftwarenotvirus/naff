package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spansDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("StartCustomSpan starts an anonymous custom span."),
		jen.Line(),
		jen.Func().ID("StartCustomSpan").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("name").ID("string")).Params(jen.Qual("context", "Context"), jen.ID("trace").Dot("Span")).Body(
			jen.If(jen.ID("ctx").Op("==").ID("nil")).Body(
				jen.ID("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return().Qual("go.opentelemetry.io/otel", "Tracer").Call(jen.Lit("_anon_")).Dot("Start").Call(
				jen.ID("ctx"),
				jen.ID("name"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("StartSpan starts an anonymous span."),
		jen.Line(),
		jen.Func().ID("StartSpan").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Qual("context", "Context"), jen.ID("trace").Dot("Span")).Body(
			jen.If(jen.ID("ctx").Op("==").ID("nil")).Body(
				jen.ID("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return().Qual("go.opentelemetry.io/otel", "Tracer").Call(jen.Lit("_anon_")).Dot("Start").Call(
				jen.ID("ctx"),
				jen.ID("GetCallerName").Call(),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("uriIDReplacementRegex").Op("=").Qual("regexp", "MustCompile").Call(jen.Lit(`/\d+`)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("FormatSpan formats a span."),
		jen.Line(),
		jen.Func().ID("FormatSpan").Params(jen.ID("operation").ID("string"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s %s: %s"),
				jen.ID("req").Dot("Method"),
				jen.ID("uriIDReplacementRegex").Dot("ReplaceAllString").Call(
					jen.ID("req").Dot("URL").Dot("Path"),
					jen.Lit("/<id>"),
				),
				jen.ID("operation"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("Span").ID("trace").Dot("Span"),
		jen.Line(),
	)

	return code
}
