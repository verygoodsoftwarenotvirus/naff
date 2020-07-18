package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metricsTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestServerConfig_ProvideInstrumentationHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
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
				utils.AssertNotNil(jen.ID("c").Dot("ProvideInstrumentationHandler").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServerConfig_ProvideTracing").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
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
				utils.AssertNoError(jen.ID("c").Dot("ProvideTracing").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()), nil),
			),
		),
		jen.Line(),
	)

	return ret
}
