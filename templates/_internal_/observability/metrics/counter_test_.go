package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_unitCounter_Decrement").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("initiatePrometheusExporter").Call(),
					jen.ID("meterProvider").Op(":=").ID("prometheusExporter").Dot("MeterProvider").Call(),
					jen.ID("mustMeter").Op(":=").ID("metric").Dot("Must").Call(jen.ID("meterProvider").Dot("Meter").Call(
						jen.ID("defaultNamespace"),
						jen.ID("metric").Dot("WithInstrumentationVersion").Call(jen.ID("instrumentationVersion")),
					)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("uc").Op(":=").Op("&").ID("unitCounter").Valuesln(jen.ID("counter").Op(":").ID("mustMeter").Dot("NewInt64Counter").Call(
						jen.ID("t").Dot("Name").Call(),
						jen.ID("metric").Dot("WithUnit").Call(jen.ID("unit").Dot("Dimensionless")),
					)),
					jen.ID("uc").Dot("Decrement").Call(jen.ID("ctx")),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_unitCounter_Increment").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("initiatePrometheusExporter").Call(),
					jen.ID("meterProvider").Op(":=").ID("prometheusExporter").Dot("MeterProvider").Call(),
					jen.ID("mustMeter").Op(":=").ID("metric").Dot("Must").Call(jen.ID("meterProvider").Dot("Meter").Call(
						jen.ID("defaultNamespace"),
						jen.ID("metric").Dot("WithInstrumentationVersion").Call(jen.ID("instrumentationVersion")),
					)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("uc").Op(":=").Op("&").ID("unitCounter").Valuesln(jen.ID("counter").Op(":").ID("mustMeter").Dot("NewInt64Counter").Call(
						jen.ID("t").Dot("Name").Call(),
						jen.ID("metric").Dot("WithUnit").Call(jen.ID("unit").Dot("Dimensionless")),
					)),
					jen.ID("uc").Dot("Increment").Call(jen.ID("ctx")),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_unitCounter_IncrementBy").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("initiatePrometheusExporter").Call(),
			jen.ID("meterProvider").Op(":=").ID("prometheusExporter").Dot("MeterProvider").Call(),
			jen.ID("mustMeter").Op(":=").ID("metric").Dot("Must").Call(jen.ID("meterProvider").Dot("Meter").Call(
				jen.ID("defaultNamespace"),
				jen.ID("metric").Dot("WithInstrumentationVersion").Call(jen.ID("instrumentationVersion")),
			)),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("uc").Op(":=").Op("&").ID("unitCounter").Valuesln(jen.ID("counter").Op(":").ID("mustMeter").Dot("NewInt64Counter").Call(
						jen.ID("t").Dot("Name").Call(),
						jen.ID("metric").Dot("WithUnit").Call(jen.ID("unit").Dot("Dimensionless")),
					)),
					jen.ID("uc").Dot("IncrementBy").Call(
						jen.ID("ctx"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
