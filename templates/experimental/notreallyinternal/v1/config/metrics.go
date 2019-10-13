package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func metricsDotGo() *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("MetricsNamespace").Op("=").Lit("todo_server").Var().ID("MinimumRuntimeCollectionInterval").Op("=").Qual("time", "Second"),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("metricsProvider").ID("string").Type().ID("tracingProvider").ID("string"),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("ErrInvalidMetricsProvider").Op("=").Qual("errors", "New").Call(jen.Lit("invalid metrics provider")).Var().ID("Prometheus").ID("metricsProvider").Op("=").Lit("prometheus").Var().ID("DefaultMetricsProvider").Op("=").ID("Prometheus").Var().ID("ErrInvalidTracingProvider").Op("=").Qual("errors", "New").Call(jen.Lit("invalid tracing provider")).Var().ID("Jaeger").ID("tracingProvider").Op("=").Lit("jaeger").Var().ID("DefaultTracingProvider").Op("=").ID("Jaeger"),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideInstrumentationHandler provides an instrumentation handler").Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("ProvideInstrumentationHandler").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("metrics").Dot(
		"InstrumentationHandler",
	), jen.ID("error")).Block(
		jen.If(jen.ID("err").Op(":=").ID("metrics").Dot(
			"RegisterDefaultViews",
		).Call(), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("registering default metric views: %w"), jen.ID("err"))),
		),
		jen.ID("_").Op("=").ID("metrics").Dot(
			"RecordRuntimeStats",
		).Call(jen.Qual("time", "Duration").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("MinimumRuntimeCollectionInterval")), jen.ID("float64").Call(jen.ID("cfg").Dot(
			"Metrics",
		).Dot(
			"RuntimeMetricsCollectionInterval",
		))))),
		jen.ID("log").Op(":=").ID("logger").Dot(
			"WithValue",
		).Call(jen.Lit("metrics_provider"), jen.ID("cfg").Dot(
			"Metrics",
		).Dot(
			"MetricsProvider",
		)),
		jen.ID("log").Dot(
			"Debug",
		).Call(jen.Lit("setting metrics provider")),
		jen.Switch(jen.ID("cfg").Dot(
			"Metrics",
		).Dot(
			"MetricsProvider",
		)).Block(
			jen.Case(jen.ID("Prometheus"), jen.ID("DefaultMetricsProvider")).Block(jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("prometheus").Dot(
				"NewExporter",
			).Call(jen.ID("prometheus").Dot(
				"Options",
			).Valuesln(jen.ID("OnError").Op(":").Func().Params(jen.ID("err").ID("error")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("setting up prometheus export")),
			), jen.ID("Namespace").Op(":").ID("string").Call(jen.ID("MetricsNamespace")))), jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("failed to create Prometheus exporter: %w"), jen.ID("err"))),
			), jen.ID("view").Dot(
				"RegisterExporter",
			).Call(jen.ID("p")), jen.ID("log").Dot(
				"Debug",
			).Call(jen.Lit("metrics provider registered")), jen.Return().List(jen.ID("p"), jen.ID("nil"))),
			jen.Default().Block(jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidMetricsProvider"))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideTracing provides an instrumentation handler").Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("ProvideTracing").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("error")).Block(
		jen.Qual("go.opencensus.io/trace", "ApplyConfig").Call(jen.Qual("go.opencensus.io/trace", "Config").Valuesln(jen.ID("DefaultSampler").Op(":").Qual("go.opencensus.io/trace", "ProbabilitySampler").Call(jen.Lit(1)))),
		jen.ID("log").Op(":=").ID("logger").Dot(
			"WithValue",
		).Call(jen.Lit("tracing_provider"), jen.ID("cfg").Dot(
			"Metrics",
		).Dot(
			"TracingProvider",
		)),
		jen.ID("log").Dot(
			"Info",
		).Call(jen.Lit("setting tracing provider")),
		jen.Switch(jen.ID("cfg").Dot(
			"Metrics",
		).Dot(
			"TracingProvider",
		)).Block(
			jen.Case(jen.ID("Jaeger"), jen.ID("DefaultTracingProvider")).Block(jen.ID("ah").Op(":=").Qual("os", "Getenv").Call(jen.Lit("JAEGER_AGENT_HOST")), jen.ID("ap").Op(":=").Qual("os", "Getenv").Call(jen.Lit("JAEGER_AGENT_PORT")), jen.ID("sn").Op(":=").Qual("os", "Getenv").Call(jen.Lit("JAEGER_SERVICE_NAME")), jen.If(jen.ID("ah").Op("!=").Lit("").Op("&&").ID("ap").Op("!=").Lit("").Op("&&").ID("sn").Op("!=").Lit("")).Block(
				jen.List(jen.ID("je"), jen.ID("err")).Op(":=").ID("jaeger").Dot(
					"NewExporter",
				).Call(jen.ID("jaeger").Dot(
					"Options",
				).Valuesln(jen.ID("AgentEndpoint").Op(":").Qual("fmt", "Sprintf").Call(jen.Lit("%s:%s"), jen.ID("ah"), jen.ID("ap")), jen.ID("Process").Op(":").ID("jaeger").Dot(
					"Process",
				).Valuesln(jen.ID("ServiceName").Op(":").ID("sn")))),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("failed to create Jaeger exporter: %w"), jen.ID("err")),
				),
				jen.Qual("go.opencensus.io/trace", "RegisterExporter").Call(jen.ID("je")),
				jen.ID("log").Dot(
					"Debug",
				).Call(jen.Lit("tracing provider registered")),
			)),
		),
		jen.Return().ID("nil"),
	),

		jen.Line(),
	)
	return ret
}
