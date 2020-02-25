package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookDataManager").Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span"),
		jen.Line(),
		jen.Func().ID("attachWebhookIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("webhookID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot("AddAttributes").Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("webhook_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("webhookID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhook fetches a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("webhook_id").Op(":").ID("webhookID"),
				jen.Lit("user_id").Op(":").ID("userID"),
			)).Dot("Debug").Call(jen.Lit("GetWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhook").Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhooks").Params(utils.CtxParam(), jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList"), jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetWebhooks")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhooks").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList"), jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllWebhooks")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooks").Call(jen.ID("ctx")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksForUser fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooksForUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook"), jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllWebhooksForUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetAllWebhooksForUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooksForUser").Call(jen.ID("ctx"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhookCount fetches the count of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhookCount").Params(utils.CtxParam(), jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetWebhookCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("filter").Op(":").ID("filter"),
				jen.Lit("user_id").Op(":").ID("userID"),
			)).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhookCount").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllWebhooksCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllWebhooksCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooksCount").Call(jen.ID("ctx")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateWebhook creates a webhook in a database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CreateWebhook").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("CreateWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("BelongsToUser")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("input").Dot("BelongsToUser")).Dot("Debug").Call(jen.Lit("CreateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("CreateWebhook").Call(jen.ID("ctx"), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateWebhook updates a particular webhook."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the provided input to have a non-zero ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateWebhook").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("UpdateWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("ID")),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("BelongsToUser")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("webhook_id"), jen.ID("input").Dot("ID")).Dot("Debug").Call(jen.Lit("UpdateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateWebhook").Call(jen.ID("ctx"), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveWebhook archives a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("ArchiveWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("webhook_id").Op(":").ID("webhookID"),
				jen.Lit("user_id").Op(":").ID("userID")),
			).Dot("Debug").Call(jen.Lit("ArchiveWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("ArchiveWebhook").Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
		),
		jen.Line(),
	)
	return ret
}
