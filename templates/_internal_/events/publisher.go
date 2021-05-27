package events

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func publisherDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("Publisher").Op("=").Parens(jen.Op("*").ID("publisher")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Publisher").Interface(jen.ID("PublishEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Interface(), jen.ID("extras").Map(jen.ID("string")).ID("string")).Params(jen.ID("error"))),
			jen.ID("publisher").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("topic").Op("*").Qual("gocloud.dev/pubsub", "Topic"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvidePublisher provides a Publisher."),
		jen.Line(),
		jen.Func().ID("ProvidePublisher").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.ID("Publisher"), jen.ID("error")).Body(
			jen.If(jen.ID("cfg").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilConfig"))),
			jen.If(jen.Op("!").ID("cfg").Dot("Enabled")).Body(
				jen.Return().List(jen.Op("&").ID("NoopEventPublisher").Values(), jen.ID("nil"))),
			jen.List(jen.ID("topic"), jen.ID("err")).Op(":=").ID("ProvidePublishTopic").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("error initializing subscription: %w"),
					jen.ID("err"),
				))),
			jen.ID("ep").Op(":=").Op("&").ID("publisher").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("event_publisher_%s"),
				jen.ID("cfg").Dot("SubscriptionIdentifier"),
			)), jen.ID("topic").Op(":").ID("topic")),
			jen.Return().List(jen.ID("ep"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PublishEvent satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("n").Op("*").ID("NoopEventPublisher")).ID("PublishEvent").Params(jen.ID("_").Qual("context", "Context"), jen.ID("_").Interface(), jen.ID("_").Map(jen.ID("string")).ID("string")).Params(jen.ID("error")).Body(
			jen.Return().ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("NoopEventPublisher").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PublishEvent satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("n").Op("*").ID("NoopEventPublisher")).ID("PublishEvent").Params(jen.ID("_").Qual("context", "Context"), jen.ID("_").Interface(), jen.ID("_").Map(jen.ID("string")).ID("string")).Params(jen.ID("error")).Body(
			jen.Return().ID("nil")),
		jen.Line(),
	)

	return code
}
