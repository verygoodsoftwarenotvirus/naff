package models

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhookDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Type().ID("WebhookDataManager").Interface(jen.ID("GetWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("Webhook"), jen.ID("error")), jen.ID("GetWebhookCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")), jen.ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")), jen.ID("GetWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("WebhookList"), jen.ID("error")), jen.ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").ID("WebhookList"), jen.ID("error")), jen.ID("GetAllWebhooksForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("Webhook"), jen.ID("error")), jen.ID("CreateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("WebhookCreationInput")).Params(jen.Op("*").ID("Webhook"), jen.ID("error")), jen.ID("UpdateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("Webhook")).Params(jen.ID("error")), jen.ID("ArchiveWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error"))).Type().ID("WebhookDataServer").Interface(jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")), jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")), jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc"))).Type().ID("Webhook").Struct(
		jen.ID("ID").ID("uint64"),
		jen.ID("Name").ID("string"),
		jen.ID("ContentType").ID("string"),
		jen.ID("URL").ID("string"),
		jen.ID("Method").ID("string"),
		jen.ID("Events").Index().ID("string"),
		jen.ID("DataTypes").Index().ID("string"),
		jen.ID("Topics").Index().ID("string"),
		jen.ID("CreatedOn").ID("uint64"),
		jen.ID("UpdatedOn").Op("*").ID("uint64"),
		jen.ID("ArchivedOn").Op("*").ID("uint64"),
		jen.ID("BelongsTo").ID("uint64"),
	).Type().ID("WebhookCreationInput").Struct(
		jen.ID("Name").ID("string"),
		jen.ID("ContentType").ID("string"),
		jen.ID("URL").ID("string"),
		jen.ID("Method").ID("string"),
		jen.ID("Events").Index().ID("string"),
		jen.ID("DataTypes").Index().ID("string"),
		jen.ID("Topics").Index().ID("string"),
		jen.ID("BelongsTo").ID("uint64"),
	).Type().ID("WebhookUpdateInput").Struct(
		jen.ID("Name").ID("string"),
		jen.ID("ContentType").ID("string"),
		jen.ID("URL").ID("string"),
		jen.ID("Method").ID("string"),
		jen.ID("Events").Index().ID("string"),
		jen.ID("DataTypes").Index().ID("string"),
		jen.ID("Topics").Index().ID("string"),
		jen.ID("BelongsTo").ID("uint64"),
	).Type().ID("WebhookList").Struct(
		jen.ID("Pagination"),
		jen.ID("Webhooks").Index().ID("Webhook"),
	),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func().ID("buildErrorLogFunc").Params(jen.ID("w").Op("*").ID("Webhook"), jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.Params(jen.ID("error"))).Block(
		jen.Return().Func().Params(jen.ID("err").ID("error")).Block(
			jen.ID("logger").Dot(
				"WithValues",
			).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("url").Op(":").ID("w").Dot(
				"URL",
			), jen.Lit("method").Op(":").ID("w").Dot(
				"Method",
			), jen.Lit("content_type").Op(":").ID("w").Dot(
				"ContentType",
			))).Dot(
				"Error",
			).Call(jen.ID("err"), jen.Lit("error executing webhook")),
		),
	),
	)
	ret.Add(jen.Func())
	return ret
}
