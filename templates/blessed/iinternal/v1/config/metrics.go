package config

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metricsDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("MetricsNamespace is the namespace under which we register metrics"),
			jen.ID("MetricsNamespace").Op("=").Lit("todo_server"),
			jen.Line(),
			jen.Comment("MinimumRuntimeCollectionInterval is the smallest interval we can collect metrics at"),
			jen.Comment("this value is used to guard against zero values"),
			jen.ID("MinimumRuntimeCollectionInterval").Op("=").Qual("time", "Second"),
			jen.Line(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.ID("metricsProvider").ID("string"),
			jen.ID("tracingProvider").ID("string"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("ErrInvalidMetricsProvider is a sentinel error value"),
			jen.ID("ErrInvalidMetricsProvider").Op("=").Qual("errors", "New").Call(jen.Lit("invalid metrics provider")),
			jen.Comment("Prometheus represents the popular time series database"),
			jen.ID("Prometheus").ID("metricsProvider").Op("=").Lit("prometheus"),
			jen.Comment("DefaultMetricsProvider indicates what the preferred metrics provider is"),
			jen.ID("DefaultMetricsProvider").Op("=").ID("Prometheus"),
			jen.Line(),
			jen.Comment("ErrInvalidTracingProvider is a sentinel error value"),
			jen.ID("ErrInvalidTracingProvider").Op("=").Qual("errors", "New").Call(jen.Lit("invalid tracing provider")),
			jen.Comment("Jaeger represents the popular distributed tracing server"),
			jen.ID("Jaeger").ID("tracingProvider").Op("=").Lit("jaeger"),
			jen.Comment("DefaultTracingProvider indicates what the preferred tracing provider is"),
			jen.ID("DefaultTracingProvider").Op("=").ID("Jaeger"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideInstrumentationHandler provides an instrumentation handler"),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("ProvideInstrumentationHandler").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "InstrumentationHandler"), jen.ID("error")).Block(
			jen.If(jen.ID("err").Op(":=").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "RegisterDefaultViews").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("registering default metric views: %w"), jen.ID("err"))),
			),
			jen.ID("_").Op("=").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "RecordRuntimeStats").Call(jen.Qual("time", "Duration").Callln(
				jen.Qual("math", "Max").Callln(
					jen.ID("float64").Call(jen.ID("MinimumRuntimeCollectionInterval")),
					jen.ID("float64").Call(jen.ID("cfg").Dot("Metrics").Dot("RuntimeMetricsCollectionInterval")),
				),
			)),
			jen.Line(),
			jen.ID("log").Op(":=").ID("logger").Dot("WithValue").Call(jen.Lit("metrics_provider"), jen.ID("cfg").Dot("Metrics").Dot("MetricsProvider")),
			jen.ID("log").Dot("Debug").Call(jen.Lit("setting metrics provider")),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Metrics").Dot("MetricsProvider")).Block(
				jen.Case(jen.ID("Prometheus"), jen.ID("DefaultMetricsProvider")).Block(
					jen.List(jen.ID("p"), jen.ID("err")).Op(":=").ID("prometheus").Dot("NewExporter").Call(jen.ID("prometheus").Dot("Options").Valuesln(
						jen.ID("OnError").Op(":").Func().Params(jen.ID("err").ID("error")).Block(
							jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("setting up prometheus export")),
						),
						jen.ID("Namespace").Op(":").ID("string").Call(jen.ID("MetricsNamespace")),
					)),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("failed to create Prometheus exporter: %w"), jen.ID("err"))),
					),
					jen.ID("view").Dot("RegisterExporter").Call(jen.ID("p")), jen.ID("log").Dot("Debug").Call(jen.Lit("metrics provider registered")),
					jen.Return().List(jen.ID("p"), jen.ID("nil"))),
				jen.Default().Block(jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidMetricsProvider"))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideTracing provides an instrumentation handler"),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("ProvideTracing").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
		)).Params(jen.ID("error")).Block(
			jen.Qual("go.opencensus.io/trace", "ApplyConfig").Call(jen.Qual("go.opencensus.io/trace", "Config").Values(jen.ID("DefaultSampler").Op(":").Qual("go.opencensus.io/trace", "ProbabilitySampler").Call(jen.Lit(1)))),
			jen.Line(),
			jen.ID("log").Op(":=").ID("logger").Dot("WithValue").Call(jen.Lit("tracing_provider"), jen.ID("cfg").Dot("Metrics").Dot("TracingProvider")),
			jen.ID("log").Dot("Info").Call(jen.Lit("setting tracing provider")),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Metrics").Dot("TracingProvider")).Block(
				jen.Case(jen.ID("Jaeger"), jen.ID("DefaultTracingProvider")).Block(
					jen.ID("ah").Op(":=").Qual("os", "Getenv").Call(jen.Lit("JAEGER_AGENT_HOST")),
					jen.ID("ap").Op(":=").Qual("os", "Getenv").Call(jen.Lit("JAEGER_AGENT_PORT")),
					jen.ID("sn").Op(":=").Qual("os", "Getenv").Call(jen.Lit("JAEGER_SERVICE_NAME")),
					jen.Line(),
					jen.If(jen.ID("ah").Op("!=").Lit("").Op("&&").ID("ap").Op("!=").Lit("").Op("&&").ID("sn").Op("!=").Lit("")).Block(
						jen.List(jen.ID("je"), jen.ID("err")).Op(":=").ID("jaeger").Dot("NewExporter").Call(jen.ID("jaeger").Dot("Options").Valuesln(
							jen.ID("AgentEndpoint").Op(":").Qual("fmt", "Sprintf").Call(jen.Lit("%s:%s"), jen.ID("ah"), jen.ID("ap")),
							jen.ID("Process").Op(":").ID("jaeger").Dot("Process").Values(jen.ID("ServiceName").Op(":").ID("sn")),
						)),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
							jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("failed to create Jaeger exporter: %w"), jen.ID("err")),
						),
						jen.Line(),
						jen.Qual("go.opencensus.io/trace", "RegisterExporter").Call(jen.ID("je")),
						jen.ID("log").Dot("Debug").Call(jen.Lit("tracing provider registered")),
					)),
			),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)
	return ret
}
