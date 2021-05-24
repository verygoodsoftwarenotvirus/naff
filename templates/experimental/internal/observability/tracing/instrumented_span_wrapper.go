package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func instrumentedSpanWrapperDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("instrumentedsql").Dot("Span").Op("=").Parens(jen.Op("*").ID("instrumentedSQLSpanWrapper")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("instrumentedSQLSpanWrapper").Struct(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("span").ID("trace").Dot("Span"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("w").Op("*").ID("instrumentedSQLSpanWrapper")).ID("NewChild").Params(jen.ID("s").ID("string")).Params(jen.ID("instrumentedsql").Dot("Span")).Body(
			jen.List(jen.ID("w").Dot("ctx"), jen.ID("w").Dot("span")).Op("=").ID("w").Dot("span").Dot("Tracer").Call().Dot("Start").Call(
				jen.ID("w").Dot("ctx"),
				jen.ID("s"),
			),
			jen.Return().ID("w"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("w").Op("*").ID("instrumentedSQLSpanWrapper")).ID("SetLabel").Params(jen.List(jen.ID("k"), jen.ID("v")).ID("string")).Body(
			jen.ID("w").Dot("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("String").Call(
				jen.ID("k"),
				jen.ID("v"),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("w").Op("*").ID("instrumentedSQLSpanWrapper")).ID("SetError").Params(jen.ID("err").ID("error")).Body(
			jen.ID("w").Dot("span").Dot("RecordError").Call(jen.ID("err"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("w").Op("*").ID("instrumentedSQLSpanWrapper")).ID("Finish").Params().Body(
			jen.ID("w").Dot("span").Dot("End").Call()),
		jen.Line(),
	)

	return code
}
