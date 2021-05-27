package events

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func subscriberDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("Subscriber").Op("=").Parens(jen.Op("*").ID("subscriber")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Subscriber").Interface(jen.ID("HandleEvents").Params(jen.ID("waitPeriod").Qual("time", "Duration"), jen.ID("stopCh").Chan().ID("bool"), jen.ID("handler").Func().Params(jen.ID("body").Index().ID("byte")))),
			jen.ID("subscriber").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("subscription").Op("*").Qual("gocloud.dev/pubsub", "Subscription"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideSubscriber provides an Subscriber."),
		jen.Line(),
		jen.Func().ID("ProvideSubscriber").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("sub").Op("*").Qual("gocloud.dev/pubsub", "Subscription"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.ID("Subscriber"), jen.ID("error")).Body(
			jen.If(jen.ID("sub").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilSubscription"))),
			jen.If(jen.ID("cfg").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilConfig"))),
			jen.If(jen.Op("!").ID("cfg").Dot("Enabled")).Body(
				jen.Return().List(jen.Op("&").ID("NoopEventSubscriber").Values(), jen.ID("nil"))),
			jen.ID("ep").Op(":=").Op("&").ID("subscriber").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("event_publisher_%s"),
				jen.ID("cfg").Dot("Topic"),
			)), jen.ID("subscription").Op(":").ID("sub")),
			jen.Return().List(jen.ID("ep"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HandleEvents satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("n").Op("*").ID("NoopEventSubscriber")).ID("HandleEvents").Params(jen.ID("_").Qual("time", "Duration"), jen.ID("_").Chan().ID("bool"), jen.ID("_").Func().Params(jen.ID("body").Index().ID("byte"))).Body(),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("Subscriber").Op("=").Parens(jen.Op("*").ID("NoopEventSubscriber")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("NoopEventSubscriber").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HandleEvents satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("n").Op("*").ID("NoopEventSubscriber")).ID("HandleEvents").Params(jen.ID("_").Qual("time", "Duration"), jen.ID("_").Chan().ID("bool"), jen.ID("_").Func().Params(jen.ID("body").Index().ID("byte"))).Body(),
		jen.Line(),
	)

	return code
}
