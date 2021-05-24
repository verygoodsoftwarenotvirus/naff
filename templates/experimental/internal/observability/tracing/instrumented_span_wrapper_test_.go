package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func instrumentedSpanWrapperTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_instrumentedSQLSpanWrapper_NewChild").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("w").Op(":=").Op("&").ID("instrumentedSQLSpanWrapper").Valuesln(jen.ID("ctx").Op(":").ID("ctx"), jen.ID("span").Op(":").ID("span")),
					jen.ID("w").Dot("NewChild").Call(jen.Lit("test")),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_instrumentedSQLSpanWrapper_SetLabel").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("w").Op(":=").Op("&").ID("instrumentedSQLSpanWrapper").Valuesln(jen.ID("ctx").Op(":").ID("ctx"), jen.ID("span").Op(":").ID("span")),
					jen.ID("w").Dot("SetLabel").Call(
						jen.Lit("things"),
						jen.Lit("stuff"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_instrumentedSQLSpanWrapper_SetError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("w").Op(":=").Op("&").ID("instrumentedSQLSpanWrapper").Valuesln(jen.ID("ctx").Op(":").ID("ctx"), jen.ID("span").Op(":").ID("span")),
					jen.ID("w").Dot("SetError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_instrumentedSQLSpanWrapper_Finish").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("w").Op(":=").Op("&").ID("instrumentedSQLSpanWrapper").Valuesln(jen.ID("ctx").Op(":").ID("ctx"), jen.ID("span").Op(":").ID("span")),
					jen.ID("w").Dot("Finish").Call(),
				),
			),
		),
		jen.Line(),
	)

	return code
}
