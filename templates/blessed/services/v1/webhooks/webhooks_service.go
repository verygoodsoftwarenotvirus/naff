package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksServiceDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("CreateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts"),
			jen.ID("CreateMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("webhook_create_input"),
			jen.Comment("UpdateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts"),
			jen.ID("UpdateMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("webhook_update_input"),
			jen.Line(),
			jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName").Equals().Lit("webhooks"),
			jen.ID("topicName").String().Equals().Lit("webhooks"),
			jen.ID("serviceName").String().Equals().Lit("webhooks_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.ModelsV1Package(), "WebhookDataServer").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()),
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
				jen.ID("webhookCounter").Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
				jen.ID("webhookDatabase").Qual(proj.ModelsV1Package(), "WebhookDataManager"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
				jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
				jen.ID("eventManager").ID("eventManager"),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher is a function that fetches user IDs"),
			jen.ID("UserIDFetcher").Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()),
			jen.Line(),
			jen.Comment("WebhookIDFetcher is a function that fetches webhook IDs"),
			jen.ID("WebhookIDFetcher").Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()),
		),

		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksService builds a new WebhooksService"),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksService").Paramsln(
			utils.CtxParam(),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("webhookDatabase").Qual(proj.ModelsV1Package(), "WebhookDataManager"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
			jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
			jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
			jen.ID("webhookCounterProvider").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider"),
			jen.ID("em").ParamPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		).Params(jen.PointerTo().ID("Service"), jen.Error()).Block(
			jen.List(jen.ID("webhookCounter"), jen.Err()).Assign().ID("webhookCounterProvider").Call(jen.ID("counterName"), jen.Lit("the number of webhooks managed by the webhooks service")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID("logger").MapAssign().ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("webhookDatabase").MapAssign().ID("webhookDatabase"),
				jen.ID("encoderDecoder").MapAssign().ID("encoder"),
				jen.ID("webhookCounter").MapAssign().ID("webhookCounter"),
				jen.ID("userIDFetcher").MapAssign().ID("userIDFetcher"),
				jen.ID("webhookIDFetcher").MapAssign().ID("webhookIDFetcher"),
				jen.ID("eventManager").MapAssign().ID("em"),
			),
			jen.Line(),
			jen.List(jen.ID("webhookCount"), jen.Err()).Assign().ID("svc").Dot("webhookDatabase").Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting current webhook count: %w"), jen.Err())),
			),
			jen.ID("svc").Dot("webhookCounter").Dot("IncrementBy").Call(utils.CtxVar(), jen.ID("webhookCount")),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.Nil()),
		),
		jen.Line(),
	)
	return ret
}
