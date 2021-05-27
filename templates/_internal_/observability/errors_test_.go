package observability

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func errorsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestPrepareError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("descriptionFmt"), jen.ID("descriptionArgs")).Op(":=").List(jen.Lit("things and %s"), jen.Lit("stuff")),
					jen.ID("err").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.ID("descriptionFmt"),
							jen.ID("descriptionArgs"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAcknowledgeError").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("descriptionFmt"), jen.ID("descriptionArgs")).Op(":=").List(jen.Lit("things and %s"), jen.Lit("stuff")),
					jen.ID("err").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
					jen.ID("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.ID("descriptionFmt"),
						jen.ID("descriptionArgs"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestNoteEvent").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("descriptionFmt"), jen.ID("descriptionArgs")).Op(":=").List(jen.Lit("things and %s"), jen.Lit("stuff")),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
					jen.ID("NoteEvent").Call(
						jen.ID("logger"),
						jen.ID("span"),
						jen.ID("descriptionFmt"),
						jen.ID("descriptionArgs"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
