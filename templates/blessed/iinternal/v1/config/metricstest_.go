package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metricsTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Func().ID("TestServerConfig_ProvideInstrumentationHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("c").Op(":=").Op("&").ID("ServerConfig").Valuesln(
					jen.ID("Metrics").Op(":").ID("MetricsSettings").Valuesln(
						jen.ID("RuntimeMetricsCollectionInterval").Op(":").Qual("time", "Second"),
						jen.ID("MetricsProvider").Op(":").ID("DefaultMetricsProvider"),
					),
				),
				jen.Line(),
				jen.List(jen.ID("ih"), jen.ID("err")).Op(":=").ID("c").Dot("ProvideInstrumentationHandler").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("ih")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServerConfig_ProvideTracing").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("c").Op(":=").Op("&").ID("ServerConfig").Valuesln(
					jen.ID("Metrics").Op(":").ID("MetricsSettings").Valuesln(
						jen.ID("TracingProvider").Op(":").ID("DefaultTracingProvider"),
					),
				),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("c").Dot("ProvideTracing").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call())),
			)),
		),
		jen.Line(),
	)
	return ret
}
