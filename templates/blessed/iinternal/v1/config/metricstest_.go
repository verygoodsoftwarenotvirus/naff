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
		jen.Func().ID("TestServerConfig_ProvideInstrumentationHandler").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("c").Assign().VarPointer().ID("ServerConfig").Valuesln(
					jen.ID("Metrics").MapAssign().ID("MetricsSettings").Valuesln(
						jen.ID("RuntimeMetricsCollectionInterval").MapAssign().Qual("time", "Second"),
						jen.ID("MetricsProvider").MapAssign().ID("DefaultMetricsProvider"),
					),
				),
				jen.Line(),
				jen.List(jen.ID("ih"), jen.Err()).Assign().ID("c").Dot("ProvideInstrumentationHandler").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("ih")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServerConfig_ProvideTracing").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("c").Assign().VarPointer().ID("ServerConfig").Valuesln(
					jen.ID("Metrics").MapAssign().ID("MetricsSettings").Valuesln(
						jen.ID("TracingProvider").MapAssign().ID("DefaultTracingProvider"),
					),
				),
				jen.Line(),
				utils.AssertNoError(jen.ID("c").Dot("ProvideTracing").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()), nil),
			)),
		),
		jen.Line(),
	)
	return ret
}
