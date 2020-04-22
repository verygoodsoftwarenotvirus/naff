package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "WebhookDataManager").Equals().Parens(jen.PointerTo().ID("Client")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhook fetches a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("webhook_id").MapAssign().ID("webhookID"),
				jen.Lit("user_id").MapAssign().ID("userID"),
			)).Dot("Debug").Call(jen.Lit("GetWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhook").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetWebhooks").Params(
			constants.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetWebhooks")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID(constants.FilterVarName)),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhooks").Call(constants.CtxVar(), jen.ID("userID"), jen.ID(constants.FilterVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllWebhooks").Params(constants.CtxParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetAllWebhooks")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooks").Call(constants.CtxVar()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllWebhooksCount").Params(constants.CtxParam()).Params(jen.ID("count").Uint64(), jen.Err().Error()).Block(
			jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetAllWebhooksCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetAllWebhooksCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooksCount").Call(constants.CtxVar()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateWebhook creates a webhook in a database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("CreateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("CreateWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot(constants.UserOwnershipFieldName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("input").Dot(constants.UserOwnershipFieldName)).Dot("Debug").Call(jen.Lit("CreateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("CreateWebhook").Call(constants.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateWebhook updates a particular webhook."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the provided input to have a non-zero ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("UpdateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("UpdateWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("ID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot(constants.UserOwnershipFieldName)),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("webhook_id"), jen.ID("input").Dot("ID")).Dot("Debug").Call(jen.Lit("UpdateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateWebhook").Call(constants.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveWebhook archives a webhook from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("ArchiveWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("ArchiveWebhook")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("webhook_id").MapAssign().ID("webhookID"),
				jen.Lit("user_id").MapAssign().ID("userID")),
			).Dot("Debug").Call(jen.Lit("ArchiveWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
		),
		jen.Line(),
	)
	return ret
}
