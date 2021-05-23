package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildWebhookTypeDefinitions()...)
	code.Add(buildWebhookUpdate()...)
	code.Add(buildWebhookbuildErrorLogFunc(proj)...)

	// if proj.EnableNewsman {
	code.Add(buildWebhookToListener(proj)...)
	// }

	return code
}

func buildWebhookTypeDefinitions() []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("Webhook represents a webhook listener, an endpoint to send an HTTP request to upon an event."),
			jen.ID("Webhook").Struct(
				jen.ID("ID").Uint64().Tag(jsonTag("id")),
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ContentType").String().Tag(jsonTag("contentType")),
				jen.ID("URL").String().Tag(jsonTag("url")),
				jen.ID("Method").String().Tag(jsonTag("method")),
				jen.ID("Events").Index().String().Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().String().Tag(jsonTag("dataTypes")),
				jen.ID("Topics").Index().String().Tag(jsonTag("topics")),
				jen.ID("CreatedOn").Uint64().Tag(jsonTag("createdOn")),
				jen.ID("LastUpdatedOn").PointerTo().Uint64().Tag(jsonTag("lastUpdatedOn")),
				jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archivedOn")),
				jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("belongsToUser")),
			),
			jen.Line(),
			jen.Comment("WebhookCreationInput represents what a user could set as input for creating a webhook."),
			jen.ID("WebhookCreationInput").Struct(
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ContentType").String().Tag(jsonTag("contentType")),
				jen.ID("URL").String().Tag(jsonTag("url")),
				jen.ID("Method").String().Tag(jsonTag("method")),
				jen.ID("Events").Index().String().Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().String().Tag(jsonTag("dataTypes")),
				jen.ID("Topics").Index().String().Tag(jsonTag("topics")),
				jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("WebhookUpdateInput represents what a user could set as input for updating a webhook."),
			jen.ID("WebhookUpdateInput").Struct(
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ContentType").String().Tag(jsonTag("contentType")),
				jen.ID("URL").String().Tag(jsonTag("url")),
				jen.ID("Method").String().Tag(jsonTag("method")),
				jen.ID("Events").Index().String().Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().String().Tag(jsonTag("dataTypes")),
				jen.ID("Topics").Index().String().Tag(jsonTag("topics")),
				jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("WebhookList represents a list of webhooks."),
			jen.ID("WebhookList").Struct(
				jen.ID("Pagination"),
				jen.ID("Webhooks").Index().ID("Webhook").Tag(jsonTag("webhooks")),
			),
			jen.Line(),
			jen.Comment("WebhookDataManager describes a structure capable of storing webhooks."),
			jen.ID("WebhookDataManager").Interface(
				jen.ID("GetWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.PointerTo().ID("Webhook"), jen.Error()),
				jen.ID("GetAllWebhooksCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()),
				jen.ID("GetWebhooks").Params(constants.CtxParam(), constants.UserIDParam(), utils.QueryFilterParam(nil)).Params(jen.PointerTo().ID("WebhookList"), jen.Error()),
				jen.ID("GetAllWebhooks").Params(constants.CtxParam()).Params(jen.PointerTo().ID("WebhookList"), jen.Error()),
				jen.ID("CreateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().ID("WebhookCreationInput")).Params(jen.PointerTo().ID("Webhook"), jen.Error()),
				jen.ID("UpdateWebhook").Params(constants.CtxParam(), jen.ID("updated").PointerTo().ID("Webhook")).Params(jen.Error()),
				jen.ID("ArchiveWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("WebhookDataServer describes a structure capable of serving traffic related to webhooks."),
			jen.ID("WebhookDataServer").Interface(
				jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Line(),
				jen.ID("ListHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("CreateHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("ReadHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("UpdateHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
				jen.ID("ArchiveHandler").Params(
					jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
					jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildWebhookUpdate() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Update merges an WebhookCreationInput with an Webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("w").PointerTo().ID("Webhook")).ID("Update").Params(jen.ID("input").PointerTo().ID("WebhookUpdateInput")).Body(
			jen.If(jen.ID("input").Dot("Name").DoesNotEqual().EmptyString()).Body(
				jen.ID("w").Dot("Name").Equals().ID("input").Dot("Name"),
			),
			jen.If(jen.ID("input").Dot("ContentType").DoesNotEqual().EmptyString()).Body(
				jen.ID("w").Dot("ContentType").Equals().ID("input").Dot("ContentType"),
			),
			jen.If(jen.ID("input").Dot("URL").DoesNotEqual().EmptyString()).Body(
				jen.ID("w").Dot("URL").Equals().ID("input").Dot("URL"),
			),
			jen.If(jen.ID("input").Dot("Method").DoesNotEqual().EmptyString()).Body(
				jen.ID("w").Dot("Method").Equals().ID("input").Dot("Method"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("Events").DoesNotEqual().ID("nil").And().ID("len").Call(jen.ID("input").Dot("Events")).GreaterThan().Zero()).Body(
				jen.ID("w").Dot("Events").Equals().ID("input").Dot("Events"),
			),
			jen.If(jen.ID("input").Dot("DataTypes").DoesNotEqual().ID("nil").And().ID("len").Call(jen.ID("input").Dot("DataTypes")).GreaterThan().Zero()).Body(
				jen.ID("w").Dot("DataTypes").Equals().ID("input").Dot("DataTypes"),
			),
			jen.If(jen.ID("input").Dot("Topics").DoesNotEqual().ID("nil").And().ID("len").Call(jen.ID("input").Dot("Topics")).GreaterThan().Zero()).Body(
				jen.ID("w").Dot("Topics").Equals().ID("input").Dot("Topics"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildWebhookbuildErrorLogFunc(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildErrorLogFunc").Params(jen.ID("w").PointerTo().ID("Webhook"), proj.LoggerParam()).Params(jen.Func().Params(jen.Error())).Body(
			jen.Return().Func().Params(jen.Err().Error()).Body(
				jen.ID(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("url").MapAssign().ID("w").Dot("URL"),
					jen.Lit("method").MapAssign().ID("w").Dot("Method"),
					jen.Lit("content_type").MapAssign().ID("w").Dot("ContentType")),
				).Dot("Error").Call(jen.Err(), jen.Lit("error executing webhook")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildWebhookToListener(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ToListener creates a newsman Listener from a Webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("w").PointerTo().ID("Webhook")).ID("ToListener").Params(proj.LoggerParam()).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Listener")).Body(
			jen.Return().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewWebhookListener").Callln(
				jen.ID("buildErrorLogFunc").Call(jen.ID("w"), jen.ID(constants.LoggerVarName)),
				jen.AddressOf().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "WebhookConfig").Valuesln(
					jen.ID("Method").MapAssign().ID("w").Dot("Method"),
					jen.ID("URL").MapAssign().ID("w").Dot("URL"),
					jen.ID("ContentType").MapAssign().ID("w").Dot("ContentType"),
				),
				jen.AddressOf().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "ListenerConfig").Valuesln(
					jen.ID("Events").MapAssign().ID("w").Dot("Events"),
					jen.ID("DataTypes").MapAssign().ID("w").Dot("DataTypes"),
					jen.ID("Topics").MapAssign().ID("w").Dot("Topics"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}
