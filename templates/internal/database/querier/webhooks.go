package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Underscore().Qual(proj.TypesPackage(), "WebhookDataManager").Equals().Parens(jen.PointerTo().ID("Client")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildGetWebhook(proj)...)
	code.Add(buildGetWebhooks(proj)...)
	code.Add(buildGetAllWebhooks(proj)...)
	code.Add(buildGetAllWebhooksCount(proj)...)
	code.Add(buildCreateWebhook(proj)...)
	code.Add(buildUpdateWebhook(proj)...)
	code.Add(buildArchiveWebhook(proj)...)

	return code
}

func buildGetWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetWebhook fetches a webhook from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "Webhook"),
			jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetWebhook")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingPackage(), "AttachWebhookIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("webhook_id").MapAssign().ID("webhookID"),
				jen.Lit("user_id").MapAssign().ID(constants.UserIDVarName),
			)).Dot("Debug").Call(jen.Lit("GetWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhook").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID(constants.UserIDVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildGetWebhooks(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetWebhooks").Params(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.TypesPackage(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.TypesPackage(), "WebhookList"), jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetWebhooks")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.FilterVarName)),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetWebhooks").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName), jen.ID(constants.FilterVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllWebhooks(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllWebhooks").Params(constants.CtxParam()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "WebhookList"), jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetAllWebhooks")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetWebhookCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooks").Call(constants.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllWebhooksCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllWebhooksCount").Params(constants.CtxParam()).Params(jen.ID("count").Uint64(), jen.Err().Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetAllWebhooksCount")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetAllWebhooksCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllWebhooksCount").Call(constants.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateWebhook creates a webhook in a database."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("CreateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "WebhookCreationInput")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "Webhook"),
			jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("CreateWebhook")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("input").Dot(constants.UserOwnershipFieldName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("input").Dot(constants.UserOwnershipFieldName)).Dot("Debug").Call(jen.Lit("CreateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("CreateWebhook").Call(constants.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdateWebhook updates a particular webhook."),
		jen.Line(),
		jen.Comment("NOTE: this function expects the provided input to have a non-zero ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("UpdateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "Webhook")).Params(jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("UpdateWebhook")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachWebhookIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("input").Dot("ID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("input").Dot(constants.UserOwnershipFieldName)),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("webhook_id"), jen.ID("input").Dot("ID")).Dot("Debug").Call(jen.Lit("UpdateWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateWebhook").Call(constants.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveWebhook archives a webhook from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("ArchiveWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.Error()).Body(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingPackage(), "StartSpan").Call(constants.CtxVar(), jen.Lit("ArchiveWebhook")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingPackage(), "AttachWebhookIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("webhookID")),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("webhook_id").MapAssign().ID("webhookID"),
				jen.Lit("user_id").MapAssign().ID(constants.UserIDVarName)),
			).Dot("Debug").Call(jen.Lit("ArchiveWebhook called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID(constants.UserIDVarName)),
		),
		jen.Line(),
	}

	return lines
}
