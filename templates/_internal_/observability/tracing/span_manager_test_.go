package tracing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanManagerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestNewTracer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("NewTracer").Call(jen.ID("t").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_otelSpanManager_StartSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("NewTracer").Call(jen.ID("t").Dot("Name").Call()).Dot("StartSpan").Call(jen.Qual("context", "Background").Call()),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_otelSpanManager_StartCustomSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("NewTracer").Call(jen.ID("t").Dot("Name").Call()).Dot("StartCustomSpan").Call(
						jen.ID("ctx"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
