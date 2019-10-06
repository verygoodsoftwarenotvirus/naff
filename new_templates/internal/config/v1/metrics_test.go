package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func metricsTestDotGo() *jen.File {
	ret := jen.NewFile("config")
	ret.Add(jen.Null(),
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
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with empty config"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("c").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Metrics").Op(":").ID("MetricsSettings").Valuesln(jen.ID("RuntimeMetricsCollectionInterval").Op(":").Qual("time", "Second"))),
			jen.List(jen.ID("ih"), jen.ID("err")).Op(":=").ID("c").Dot(
				"ProvideInstrumentationHandler",
			).Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("err"), jen.ID("ErrInvalidMetricsProvider")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("ih")),
		)),
	),
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
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with empty config"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("c").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Metrics").Op(":").ID("MetricsSettings").Valuesln()),
			jen.ID("err").Op(":=").ID("c").Dot(
				"ProvideTracing",
			).Call(jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("err"), jen.ID("ErrInvalidTracingProvider")),
		)),
	),
	)
	return ret
}
