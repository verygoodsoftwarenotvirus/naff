package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksServiceDotGo() *jen.File {
	ret := jen.NewFile("webhooks")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("CreateMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("webhook_create_input").Var().ID("UpdateMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("webhook_update_input").Var().ID("counterName").ID("metrics").Dot(
		"CounterName",
	).Op("=").Lit("webhooks").Var().ID("topicName").ID("string").Op("=").Lit("webhooks").Var().ID("serviceName").ID("string").Op("=").Lit("webhooks_service"),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"WebhookDataServer",
	).Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("eventManager").Interface(jen.ID("newsman").Dot(
		"Reporter",
	), jen.ID("TuneIn").Params(jen.ID("newsman").Dot(
		"Listener",
	))).Type().ID("Service").Struct(
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("webhookCounter").ID("metrics").Dot(
			"UnitCounter",
		),
		jen.ID("webhookDatabase").ID("models").Dot(
			"WebhookDataManager",
		),
		jen.ID("userIDFetcher").ID("UserIDFetcher"),
		jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
		jen.ID("encoderDecoder").ID("encoding").Dot(
			"EncoderDecoder",
		),
		jen.ID("eventManager").ID("eventManager"),
	).Type().ID("UserIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Type().ID("WebhookIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
	)
	ret.Add(jen.Func())
	return ret
}
