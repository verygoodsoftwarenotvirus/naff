package webhooks

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksServiceDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("CreateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts"),
			jen.ID("CreateMiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Op("=").Lit("webhook_create_input"),
			jen.Comment("UpdateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts"),
			jen.ID("UpdateMiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Op("=").Lit("webhook_update_input"),
			jen.Line(),
			jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName").Op("=").Lit("webhooks"),
			jen.ID("topicName").ID("string").Op("=").Lit("webhooks"),
			jen.ID("serviceName").ID("string").Op("=").Lit("webhooks_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookDataServer").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.Nil()),
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
				jen.ID("webhookCounter").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
				jen.ID("webhookDatabase").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookDataManager"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
				jen.ID("encoderDecoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
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
			utils.CtxParam(),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("webhookDatabase").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookDataManager"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
			jen.ID("webhookIDFetcher").ID("WebhookIDFetcher"),
			jen.ID("encoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
			jen.ID("webhookCounterProvider").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider"),
			jen.ID("em").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Newsman"),
		).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID("webhookCounter"), jen.Err()).Op(":=").ID("webhookCounterProvider").Call(jen.ID("counterName"), jen.Lit("the number of webhooks managed by the webhooks service")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
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
			jen.List(jen.ID("webhookCount"), jen.Err()).Op(":=").ID("svc").Dot("webhookDatabase").Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
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
