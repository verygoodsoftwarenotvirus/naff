package client

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(proj.OutputPath, "models/v1"), "WebhookDataManager").Equals().Parens(jen.Op("*").ID("Client")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhook fetches a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("webhook_id").MapAssign().ID("webhookID"),
				jen.Lit("user_id").MapAssign().ID("userID"),
			)).Dot("Debug").Call(jen.Lit("GetWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhooks").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter"),
		).Params(jen.Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "WebhookList"), jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetWebhooks")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID(utils.FilterVarName)),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhooks").Call(utils.CtxVar(), jen.ID("userID"), jen.ID(utils.FilterVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooks").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "WebhookList"), jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetAllWebhooks")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooks").Call(utils.CtxVar()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksForUser fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooksForUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook"), jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetAllWebhooksForUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetAllWebhooksForUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooksForUser").Call(utils.CtxVar(), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhookCount fetches the count of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhookCount").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter"),
		).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetWebhookCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID(utils.FilterVarName)),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("filter").MapAssign().ID("filter"),
				jen.Lit("user_id").MapAssign().ID("userID"),
			)).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhookCount").Call(utils.CtxVar(), jen.ID("userID"), jen.ID(utils.FilterVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooksCount").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetAllWebhooksCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllWebhooksCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateWebhook creates a webhook in a database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CreateWebhook").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "WebhookCreationInput")).Params(jen.Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("CreateWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("BelongsToUser")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("input").Dot("BelongsToUser")).Dot("Debug").Call(jen.Lit("CreateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateWebhook updates a particular webhook."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the provided input to have a non-zero ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateWebhook").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "Webhook")).Params(jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("UpdateWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("ID")),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("BelongsToUser")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("webhook_id"), jen.ID("input").Dot("ID")).Dot("Debug").Call(jen.Lit("UpdateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveWebhook archives a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(utils.CtxVar(), jen.Lit("ArchiveWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(filepath.Join(proj.OutputPath, "internal/v1/tracing"), "AttachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.Lit("webhook_id").MapAssign().ID("webhookID"),
				jen.Lit("user_id").MapAssign().ID("userID")),
			).Dot("Debug").Call(jen.Lit("ArchiveWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
		),
		jen.Line(),
	)
	return ret
}
