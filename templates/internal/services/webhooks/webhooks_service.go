package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildWebhooksServiceConstDefs(proj)...)
	code.Add(buildWebhooksServiceVarDefs(proj)...)
	code.Add(buildWebhooksServiceTypeDefs(proj)...)
	code.Add(buildProvideWebhooksService(proj)...)

	return code
}

func buildWebhooksServiceConstDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("createMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts."),
			jen.ID("createMiddlewareCtxKey").Qual(proj.TypesPackage(), "ContextKey").Equals().Lit("webhook_create_input"),
			jen.Comment("updateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts."),
			jen.ID("updateMiddlewareCtxKey").Qual(proj.TypesPackage(), "ContextKey").Equals().Lit("webhook_update_input"),
			jen.Line(),
			jen.ID("counterName").Qual(proj.InternalMetricsPackage(), "CounterName").Equals().Lit("webhooks"),
			jen.ID("counterDescription").String().Equals().Lit("the number of webhooks managed by the webhooks service"),
			jen.ID("topicName").String().Equals().Lit("webhooks"),
			jen.ID("serviceName").String().Equals().Lit("webhooks_service"),
		),
		jen.Line(),
	}

	return lines
}

func buildWebhooksServiceVarDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Underscore().Qual(proj.TypesPackage(), "WebhookDataServer").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildWebhooksServiceTypeDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.ID("eventManager").Interface(
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
				jen.Line(),
				jen.ID("TuneIn").Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Listener")),
			),
			jen.Line(),
			jen.Comment("Service handles TODO ListHandler webhooks."),
			jen.ID("Service").Struct(
				constants.LoggerParam(),
				jen.ID("webhookCounter").Qual(proj.InternalMetricsPackage(), "UnitCounter"),
				jen.ID("webhookDataManager").Qual(proj.TypesPackage(), "WebhookDataManager"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
				jen.ID("encoderDecoder").Qual(proj.InternalEncodingPackage(), "EncoderDecoder"),
				jen.ID("eventManager").ID("eventManager"),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher is a function that fetches user IDs."),
			jen.ID("UserIDFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
			jen.Line(),
			jen.Comment("WebhookIDFetcher is a function that fetches webhook IDs."),
			jen.ID("WebhookIDFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
		),

		jen.Line(),
	}

	return lines
}

func buildProvideWebhooksService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideWebhooksService builds a new WebhooksService."),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksService").Paramsln(
			constants.LoggerParam(),
			jen.ID("webhookDataManager").Qual(proj.TypesPackage(), "WebhookDataManager"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
			jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
			jen.ID("encoder").Qual(proj.InternalEncodingPackage(), "EncoderDecoder"),
			jen.ID("webhookCounterProvider").Qual(proj.InternalMetricsPackage(), "UnitCounterProvider"),
			jen.ID("em").PointerTo().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		).Params(jen.PointerTo().ID("Service"), jen.Error()).Body(
			jen.List(jen.ID("webhookCounter"), jen.Err()).Assign().ID("webhookCounterProvider").Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("webhookDataManager").MapAssign().ID("webhookDataManager"),
				jen.ID("encoderDecoder").MapAssign().ID("encoder"),
				jen.ID("webhookCounter").MapAssign().ID("webhookCounter"),
				jen.ID("userIDFetcher").MapAssign().ID("userIDFetcher"),
				jen.ID("webhookIDFetcher").MapAssign().ID("webhookIDFetcher"),
				jen.ID("eventManager").MapAssign().ID("em"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
