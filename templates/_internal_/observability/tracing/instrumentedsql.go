package tracing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func instrumentedsqlDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("NewInstrumentedSQLTracer wraps a Tracer for instrumentedsql."),
		jen.Line(),
		jen.Func().ID("NewInstrumentedSQLTracer").Params(jen.ID("name").ID("string")).Params(jen.ID("instrumentedsql").Dot("Tracer")).Body(
			jen.Return().Op("&").ID("instrumentedSQLTracerWrapper").Valuesln(jen.ID("tracer").Op(":").ID("NewTracer").Call(jen.ID("name")))),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("instrumentedsql").Dot("Tracer").Op("=").Parens(jen.Op("*").ID("instrumentedSQLTracerWrapper")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("instrumentedSQLTracerWrapper").Struct(jen.ID("tracer").ID("Tracer")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetSpan wraps tracer.GetSpan."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("instrumentedSQLTracerWrapper")).ID("GetSpan").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("instrumentedsql").Dot("Span")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("t").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Return().Op("&").ID("instrumentedSQLSpanWrapper").Valuesln(jen.ID("ctx").Op(":").ID("ctx"), jen.ID("span").Op(":").ID("span")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewInstrumentedSQLLogger wraps a logging.Logger for instrumentedsql."),
		jen.Line(),
		jen.Func().ID("NewInstrumentedSQLLogger").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("instrumentedsql").Dot("Logger")).Body(
			jen.Return().Op("&").ID("instrumentedSQLLoggerWrapper").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.Lit("sql")))),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("instrumentedSQLLoggerWrapper").Struct(jen.ID("logger").ID("logging").Dot("Logger")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("w").Op("*").ID("instrumentedSQLLoggerWrapper")).ID("Log").Params(jen.ID("_").Qual("context", "Context"), jen.ID("msg").ID("string"), jen.ID("keyvals").Op("...").Interface()).Body(),
		jen.Line(),
	)

	return code
}
