package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metricsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("config")

	utils.AddImports(proj, code)

	code.Add(
		jen.Const().Defs(
			jen.Comment("MetricsNamespace is the namespace under which we register metrics."),
			jen.ID("MetricsNamespace").Equals().Lit("todo_server"),
			jen.Line(),
			jen.Comment("MinimumRuntimeCollectionInterval is the smallest interval we can collect metrics at"),
			jen.Comment("this value is used to guard against zero values."),
			jen.ID("MinimumRuntimeCollectionInterval").Equals().Qual("time", "Second"),
			jen.Line(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("metricsProvider").String(),
			jen.ID("tracingProvider").String(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.Comment("ErrInvalidMetricsProvider is a sentinel error value."),
			jen.ID("ErrInvalidMetricsProvider").Equals().Qual("errors", "New").Call(jen.Lit("invalid metrics provider")),
			jen.Comment("Prometheus represents the popular time series database."),
			jen.ID("Prometheus").ID("metricsProvider").Equals().Lit("prometheus"),
			jen.Comment("DefaultMetricsProvider indicates what the preferred metrics provider is."),
			jen.ID("DefaultMetricsProvider").Equals().ID("Prometheus"),
			jen.Line(),
			jen.Comment("ErrInvalidTracingProvider is a sentinel error value."),
			jen.ID("ErrInvalidTracingProvider").Equals().Qual("errors", "New").Call(jen.Lit("invalid tracing provider")),
			jen.Comment("Jaeger represents the popular distributed tracing server."),
			jen.ID("Jaeger").ID("tracingProvider").Equals().Lit("jaeger"),
			jen.Comment("DefaultTracingProvider indicates what the preferred tracing provider is."),
			jen.ID("DefaultTracingProvider").Equals().ID("Jaeger"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideInstrumentationHandler provides an instrumentation handler."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("ProvideInstrumentationHandler").Params(
			constants.LoggerParam(),
		).Params(
			jen.Qual(proj.InternalMetricsV1Package(), "InstrumentationHandler"),
		).Block(
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("metrics_provider"), jen.ID("cfg").Dot("Metrics").Dot("MetricsProvider")),
			jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("setting metrics provider")),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Metrics").Dot("MetricsProvider")).Block(
				jen.Case(jen.ID("Prometheus")).Block(
					jen.List(jen.ID("p"), jen.Err()).Assign().Qual("contrib.go.opencensus.io/exporter/prometheus", "NewExporter").Callln(
						jen.Qual("contrib.go.opencensus.io/exporter/prometheus", "Options").Valuesln(
							jen.ID("OnError").MapAssign().Func().Params(jen.Err().Error()).Block(
								jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("setting up prometheus export")),
							),
							jen.ID("Namespace").MapAssign().ID("MetricsNamespace"),
						),
					),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("failed to create Prometheus exporter")),
						jen.Return(jen.Nil()),
					),
					jen.Qual("go.opencensus.io/stats/view", "RegisterExporter").Call(jen.ID("p")), jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("metrics provider registered")),
					jen.Line(),
					jen.If(jen.Err().Assign().Qual(proj.InternalMetricsV1Package(), "RegisterDefaultViews").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("registering default metric views")),
						jen.Return(jen.Nil()),
					),
					jen.Qual(proj.InternalMetricsV1Package(), "RecordRuntimeStats").Call(
						jen.Qual("time", "Duration").Callln(
							jen.Qual("math", "Max").Callln(
								jen.ID("float64").Call(jen.ID("MinimumRuntimeCollectionInterval")),
								jen.ID("float64").Call(jen.ID("cfg").Dot("Metrics").Dot("RuntimeMetricsCollectionInterval")),
							),
						),
					),
					jen.Line(),
					jen.Return(jen.ID("p")),
					jen.Default().Block(jen.Return(jen.Nil())),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideTracing provides an instrumentation handler."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("ServerConfig")).ID("ProvideTracing").Params(
			constants.LoggerParam(),
		).Params(jen.Error()).Block(
			jen.Qual("go.opencensus.io/trace", "ApplyConfig").Call(jen.Qual("go.opencensus.io/trace", "Config").Values(jen.ID("DefaultSampler").MapAssign().Qual("go.opencensus.io/trace", "ProbabilitySampler").Call(jen.One()))),
			jen.Line(),
			jen.ID("log").Assign().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("tracing_provider"), jen.ID("cfg").Dot("Metrics").Dot("TracingProvider")),
			jen.ID("log").Dot("Info").Call(jen.Lit("setting tracing provider")),
			jen.Line(),
			jen.Switch(jen.ID("cfg").Dot("Metrics").Dot("TracingProvider")).Block(
				jen.Case(jen.ID("Jaeger")).Block(
					jen.ID("ah").Assign().Qual("os", "Getenv").Call(jen.Lit("JAEGER_AGENT_HOST")),
					jen.ID("ap").Assign().Qual("os", "Getenv").Call(jen.Lit("JAEGER_AGENT_PORT")),
					jen.ID("sn").Assign().Qual("os", "Getenv").Call(jen.Lit("JAEGER_SERVICE_NAME")),
					jen.Line(),
					jen.If(jen.ID("ah").DoesNotEqual().EmptyString().And().ID("ap").DoesNotEqual().EmptyString().And().ID("sn").DoesNotEqual().EmptyString()).Block(
						jen.List(jen.ID("je"), jen.Err()).Assign().Qual("contrib.go.opencensus.io/exporter/jaeger", "NewExporter").Call(jen.Qual("contrib.go.opencensus.io/exporter/jaeger", "Options").Valuesln(
							jen.ID("AgentEndpoint").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit("%s:%s"), jen.ID("ah"), jen.ID("ap")),
							jen.ID("Process").MapAssign().Qual("contrib.go.opencensus.io/exporter/jaeger", "Process").Values(jen.ID("ServiceName").MapAssign().ID("sn")),
						)),
						jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
							jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("failed to create Jaeger exporter: %w"), jen.Err()),
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

	return code
}
