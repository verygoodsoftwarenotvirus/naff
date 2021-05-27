package tracing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func tracerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("errorHandler").Struct(jen.ID("logger").ID("logging").Dot("Logger")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("h").ID("errorHandler")).ID("Handle").Params(jen.ID("err").ID("error")).Body(
			jen.ID("h").Dot("logger").Dot("Error").Call(
				jen.ID("err"),
				jen.Lit("tracer reported issue"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("go.opentelemetry.io/otel", "SetErrorHandler").Call(jen.ID("errorHandler").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("NewNoopLogger").Call().Dot("WithName").Call(jen.Lit("otel_errors"))))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetupJaeger creates a new trace provider instance and registers it as global trace provider."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Config")).ID("SetupJaeger").Params().Params(jen.Func().Params(), jen.ID("error")).Body(
			jen.List(jen.ID("flush"), jen.ID("err")).Op(":=").ID("jaeger").Dot("InstallNewPipeline").Call(
				jen.ID("jaeger").Dot("WithCollectorEndpoint").Call(jen.ID("c").Dot("Jaeger").Dot("CollectorEndpoint")),
				jen.ID("jaeger").Dot("WithProcessFromEnv").Call(),
				jen.ID("jaeger").Dot("WithSDKOptions").Call(
					jen.Qual("go.opentelemetry.io/otel/sdk/trace", "WithSampler").Call(jen.Qual("go.opentelemetry.io/otel/sdk/trace", "TraceIDRatioBased").Call(jen.ID("c").Dot("SpanCollectionProbability"))),
					jen.Qual("go.opentelemetry.io/otel/sdk/trace", "WithResource").Call(jen.ID("resource").Dot("NewWithAttributes").Call(jen.ID("attribute").Dot("String").Call(
						jen.Lit("exporter"),
						jen.Lit("jaeger"),
					))),
				),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("initializing Jaeger: %w"),
					jen.ID("err"),
				))),
			jen.Return().List(jen.ID("flush"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Tracer").Interface(
				jen.ID("StartSpan").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Qual("context", "Context"), jen.ID("Span")),
				jen.ID("StartCustomSpan").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("name").ID("string")).Params(jen.Qual("context", "Context"), jen.ID("Span")),
			),
		),
		jen.Line(),
	)

	return code
}
