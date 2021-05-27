package tracing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func tracerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_tracingErrorHandler_Handle").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("errorHandler").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("NewNoopLogger").Call()).Dot("Handle").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConfig_SetupJaeger").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Jaeger").Op(":").Op("&").ID("JaegerConfig").Valuesln(jen.ID("CollectorEndpoint").Op(":").Lit("blah blah blah"), jen.ID("ServiceName").Op(":").ID("t").Dot("Name").Call()), jen.ID("Provider").Op(":").ID("Jaeger"), jen.ID("SpanCollectionProbability").Op(":").Lit(0)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("cfg").Dot("SetupJaeger").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with empty collector endpoint"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Jaeger").Op(":").Op("&").ID("JaegerConfig").Valuesln(jen.ID("CollectorEndpoint").Op(":").Lit(""), jen.ID("ServiceName").Op(":").ID("t").Dot("Name").Call()), jen.ID("Provider").Op(":").ID("Jaeger"), jen.ID("SpanCollectionProbability").Op(":").Lit(0)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("cfg").Dot("SetupJaeger").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
