package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metricsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestServerConfig_ProvideInstrumentationHandler()...)
	code.Add(buildTestServerConfig_ProvideTracing()...)

	return code
}

func buildTestServerConfig_ProvideInstrumentationHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestServerConfig_ProvideInstrumentationHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("c").Assign().AddressOf().ID("ServerConfig").Valuesln(
					jen.ID("Metrics").MapAssign().ID("MetricsSettings").Valuesln(
						jen.ID("RuntimeMetricsCollectionInterval").MapAssign().Qual("time", "Second"),
						jen.ID("MetricsProvider").MapAssign().ID("DefaultMetricsProvider"),
					),
				),
				jen.Line(),
				utils.AssertNotNil(jen.ID("c").Dot("ProvideInstrumentationHandler").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestServerConfig_ProvideTracing() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestServerConfig_ProvideTracing").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("c").Assign().AddressOf().ID("ServerConfig").Valuesln(
					jen.ID("Metrics").MapAssign().ID("MetricsSettings").Valuesln(
						jen.ID("TracingProvider").MapAssign().ID("DefaultTracingProvider"),
					),
				),
				jen.Line(),
				utils.AssertNoError(jen.ID("c").Dot("ProvideTracing").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
