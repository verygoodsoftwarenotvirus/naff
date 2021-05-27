package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("defaultNamespace").Op("=").Lit("todo_server"),
			jen.ID("instrumentationVersion").Op("=").Lit("1.0.0"),
			jen.ID("minimumRuntimeCollectionInterval").Op("=").Qual("time", "Second"),
			jen.ID("DefaultMetricsCollectionInterval").Op("=").Lit(2).Op("*").Qual("time", "Second"),
			jen.ID("Prometheus").Op("=").Lit("prometheus"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Config").Struct(
				jen.ID("Provider").ID("string"),
				jen.ID("RouteToken").ID("string"),
				jen.ID("RuntimeMetricsCollectionInterval").Qual("time", "Duration"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the config struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("RuntimeMetricsCollectionInterval"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Min").Call(jen.ID("minimumRuntimeCollectionInterval")),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("prometheusExporterInitOnce").Qual("sync", "Once"),
			jen.ID("prometheusExporter").Op("*").Qual("go.opentelemetry.io/otel/exporters/metric/prometheus", "Exporter"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("initiatePrometheusExporter").Params().Body(
			jen.ID("prometheusExporterInitOnce").Dot("Do").Call(jen.Func().Params().Body(
				jen.Var().Defs(
					jen.ID("err").ID("error"),
				),
				jen.If(jen.List(jen.ID("prometheusExporter"), jen.ID("err")).Op("=").Qual("go.opentelemetry.io/otel/exporters/metric/prometheus", "InstallNewPipeline").Call(jen.Qual("go.opentelemetry.io/otel/exporters/metric/prometheus", "Config").Values()), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("panic").Call(jen.ID("err"))),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideInstrumentationHandler provides an instrumentation handler."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideInstrumentationHandler").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("InstrumentationHandler"), jen.ID("error")).Body(
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("metrics_provider"),
				jen.ID("cfg").Dot("Provider"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("setting metrics provider")),
			jen.If(jen.ID("err").Op(":=").ID("runtime").Dot("Start").Call(jen.ID("runtime").Dot("WithMinimumReadMemStatsInterval").Call(jen.ID("cfg").Dot("RuntimeMetricsCollectionInterval"))), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("failed to start runtime metrics collection: %w"),
					jen.ID("err"),
				))),
			jen.Switch(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("cfg").Dot("Provider")))).Body(
				jen.Case(jen.ID("Prometheus")).Body(
					jen.ID("initiatePrometheusExporter").Call(), jen.Return().List(jen.ID("prometheusExporter"), jen.ID("nil"))),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUnitCounterProvider provides an instrumentation handler."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ProvideUnitCounterProvider").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("UnitCounterProvider"), jen.ID("error")).Body(
			jen.ID("p").Op(":=").Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("cfg").Dot("Provider"))),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("metrics_provider"),
				jen.ID("p"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("setting up meter")),
			jen.Switch(jen.ID("p")).Body(
				jen.Case(jen.ID("Prometheus")).Body(
					jen.ID("initiatePrometheusExporter").Call(), jen.ID("meterProvider").Op(":=").ID("prometheusExporter").Dot("MeterProvider").Call(), jen.ID("mustMeter").Op(":=").ID("metric").Dot("Must").Call(jen.ID("meterProvider").Dot("Meter").Call(
						jen.ID("defaultNamespace"),
						jen.ID("metric").Dot("WithInstrumentationVersion").Call(jen.ID("instrumentationVersion")),
					)), jen.ID("logger").Dot("Debug").Call(jen.Lit("meter initialized successfully")), jen.Return().List(jen.Func().Params(jen.List(jen.ID("name"), jen.ID("description")).ID("string")).Params(jen.ID("UnitCounter")).Body(
						jen.ID("l").Op(":=").ID("logger").Dot("WithValue").Call(
							jen.Lit("name"),
							jen.ID("name"),
						),
						jen.ID("counter").Op(":=").ID("mustMeter").Dot("NewInt64Counter").Call(
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("%s_count"),
								jen.ID("name"),
							),
							jen.ID("metric").Dot("WithUnit").Call(jen.ID("unit").Dot("Dimensionless")),
							jen.ID("metric").Dot("WithDescription").Call(jen.ID("description")),
							jen.ID("metric").Dot("WithInstrumentationName").Call(jen.ID("name")),
							jen.ID("metric").Dot("WithInstrumentationVersion").Call(jen.ID("instrumentationVersion")),
						),
						jen.ID("l").Dot("Debug").Call(jen.Lit("returning wrapped unit counter")),
						jen.Return().Op("&").ID("unitCounter").Valuesln(jen.ID("counter").Op(":").ID("counter")),
					), jen.ID("nil"))),
				jen.Default().Body(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("nil unit counter provider")), jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
			),
		),
		jen.Line(),
	)

	return code
}
