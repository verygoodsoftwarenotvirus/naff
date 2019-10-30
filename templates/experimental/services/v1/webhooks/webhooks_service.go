package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksServiceDotGo() *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("CreateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts"),
			jen.ID("CreateMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("webhook_create_input"),
			jen.Comment("UpdateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts"),
			jen.ID("UpdateMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("webhook_update_input"),
			jen.Line(),
			jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName").Op("=").Lit("webhooks"),
			jen.ID("topicName").ID("string").Op("=").Lit("webhooks"),
			jen.ID("serviceName").ID("string").Op("=").Lit("webhooks_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").ID("models").Dot("WebhookDataServer").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.ID("eventManager").Interface(
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
				jen.Line(),
				jen.ID("TuneIn").Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Listener")),
			),
			jen.Line(),
			jen.Comment("Service handles TODO ListHandler webhooks"),
			jen.ID("Service").Struct(
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("webhookCounter").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
				jen.ID("webhookDatabase").ID("models").Dot("WebhookDataManager"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
				jen.ID("encoderDecoder").ID("encoding").Dot("EncoderDecoder"),
				jen.ID("eventManager").ID("eventManager"),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher is a function that fetches user IDs"),
			jen.ID("UserIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.Line(),
			jen.Comment("WebhookIDFetcher is a function that fetches webhook IDs"),
			jen.ID("WebhookIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		),

		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksService builds a new WebhooksService"),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksService").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("webhookDatabase").ID("models").Dot("WebhookDataManager"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
			jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
			jen.ID("encoder").ID("encoding").Dot("EncoderDecoder"),
			jen.ID("webhookCounterProvider").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider"),
			jen.ID("em").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID("webhookCounter"), jen.ID("err")).Op(":=").ID("webhookCounterProvider").Call(jen.ID("counterName"), jen.Lit("the number of webhooks managed by the webhooks service")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("webhookDatabase").Op(":").ID("webhookDatabase"),
				jen.ID("encoderDecoder").Op(":").ID("encoder"),
				jen.ID("webhookCounter").Op(":").ID("webhookCounter"),
				jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"),
				jen.ID("webhookIDFetcher").Op(":").ID("webhookIDFetcher"),
				jen.ID("eventManager").Op(":").ID("em"),
			),
			jen.Line(),
			jen.List(jen.ID("webhookCount"), jen.ID("err")).Op(":=").ID("svc").Dot("webhookDatabase").Dot("GetAllWebhooksCount").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting current webhook count: %w"), jen.ID("err"))),
			),
			jen.ID("svc").Dot("webhookCounter").Dot("IncrementBy").Call(jen.ID("ctx"), jen.ID("webhookCount")),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Line(),
	)
	return ret
}
