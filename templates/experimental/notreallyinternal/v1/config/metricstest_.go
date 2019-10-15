package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func metricsTestDotGo() *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestServerConfig_ProvideInstrumentationHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("c").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Metrics").Op(":").ID("MetricsSettings").Valuesln(jen.ID("RuntimeMetricsCollectionInterval").Op(":").Qual("time", "Second"), jen.ID("MetricsProvider").Op(":").ID("DefaultMetricsProvider"))),
			jen.List(jen.ID("ih"), jen.ID("err")).Op(":=").ID("c").Dot(
				"ProvideInstrumentationHandler",
			).Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("ih")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestServerConfig_ProvideTracing").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("c").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Metrics").Op(":").ID("MetricsSettings").Valuesln(jen.ID("TracingProvider").Op(":").ID("DefaultTracingProvider"))),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("c").Dot(
				"ProvideTracing",
			).Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call())),
		)),
	),

		jen.Line(),
	)
	return ret
}