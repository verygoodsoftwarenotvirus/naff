package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestConfig_ProvideInstrumentationHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("Prometheus"), jen.ID("RuntimeMetricsCollectionInterval").Op(":").ID("minimumRuntimeCollectionInterval")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideInstrumentationHandler").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
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
				jen.Lit("without provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("RuntimeMetricsCollectionInterval").Op(":").ID("minimumRuntimeCollectionInterval")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("cfg").Dot("ProvideInstrumentationHandler").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
					jen.ID("assert").Dot("NoError").Call(
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

	code.Add(
		jen.Func().ID("TestConfig_ProvideUnitCounterProvider").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("Prometheus"), jen.ID("RuntimeMetricsCollectionInterval").Op(":").ID("minimumRuntimeCollectionInterval")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideUnitCounterProvider").Call(
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNoopLogger").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("actual").Call(
						jen.Lit("things"),
						jen.Lit("stuff"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("RuntimeMetricsCollectionInterval").Op(":").ID("minimumRuntimeCollectionInterval")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideUnitCounterProvider").Call(
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNoopLogger").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
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

	code.Add(
		jen.Func().ID("TestConfig_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("RuntimeMetricsCollectionInterval").Op(":").ID("minimumRuntimeCollectionInterval")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_initiatePrometheusExporter").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("initiatePrometheusExporter").Call(),
				),
			),
		),
		jen.Line(),
	)

	return code
}
