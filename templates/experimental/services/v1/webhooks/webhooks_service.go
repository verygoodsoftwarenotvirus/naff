package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksServiceDotGo() *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("CreateMiddlewareCtxKey").ID("models").Dot(
			"ContextKey",
		).Op("=").Lit("webhook_create_input").Var().ID("UpdateMiddlewareCtxKey").ID("models").Dot(
			"ContextKey",
		).Op("=").Lit("webhook_update_input").Var().ID("counterName").ID("metrics").Dot(
			"CounterName",
		).Op("=").Lit("webhooks").Var().ID("topicName").ID("string").Op("=").Lit("webhooks").Var().ID("serviceName").ID("string").Op("=").Lit("webhooks_service"),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("models").Dot(
			"WebhookDataServer",
		).Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("eventManager").Interface(jen.ID("newsman").Dot(
			"Reporter",
		), jen.ID("TuneIn").Params(jen.ID("newsman").Dot(
			"Listener",
		))).Type().ID("Service").Struct(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
		), jen.ID("webhookCounter").ID("metrics").Dot(
			"UnitCounter",
		), jen.ID("webhookDatabase").ID("models").Dot(
			"WebhookDataManager",
		), jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"), jen.ID("encoderDecoder").ID("encoding").Dot(
			"EncoderDecoder",
		), jen.ID("eventManager").ID("eventManager")).Type().ID("UserIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Type().ID("WebhookIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksService builds a new WebhooksService"),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksService").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
		), jen.ID("webhookDatabase").ID("models").Dot(
			"WebhookDataManager",
		), jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"), jen.ID("encoder").ID("encoding").Dot(
			"EncoderDecoder",
		), jen.ID("webhookCounterProvider").ID("metrics").Dot(
			"UnitCounterProvider",
		), jen.ID("em").Op("*").ID("newsman").Dot(
			"Newsman",
		)).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID("webhookCounter"), jen.ID("err")).Op(":=").ID("webhookCounterProvider").Call(jen.ID("counterName"), jen.Lit("the number of webhooks managed by the webhooks service")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.ID("err"))),
			),
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("logger").Dot(
				"WithName",
			).Call(jen.ID("serviceName")), jen.ID("webhookDatabase").Op(":").ID("webhookDatabase"), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("webhookCounter").Op(":").ID("webhookCounter"), jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"), jen.ID("webhookIDFetcher").Op(":").ID("webhookIDFetcher"), jen.ID("eventManager").Op(":").ID("em")),
			jen.List(jen.ID("webhookCount"), jen.ID("err")).Op(":=").ID("svc").Dot(
				"webhookDatabase",
			).Dot(
				"GetAllWebhooksCount",
			).Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting current webhook count: %w"), jen.ID("err"))),
			),
			jen.ID("svc").Dot(
				"webhookCounter",
			).Dot(
				"IncrementBy",
			).Call(jen.ID("ctx"), jen.ID("webhookCount")),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Line(),
	)
	return ret
}
