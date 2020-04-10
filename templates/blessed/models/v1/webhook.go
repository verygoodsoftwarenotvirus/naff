package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("WebhookDataManager describes a structure capable of storing webhooks"),
			jen.ID("WebhookDataManager").Interface(
				jen.ID("GetWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().ID("Webhook"), jen.Error()),
				jen.ID("GetAllWebhooksCount").Params(utils.CtxParam()).Params(jen.Uint64(), jen.Error()),
				jen.ID("GetWebhooks").Params(utils.CtxParam(), jen.ID("userID").Uint64(), utils.QueryFilterParam(nil)).Params(jen.PointerTo().ID("WebhookList"), jen.Error()),
				jen.ID("GetAllWebhooks").Params(utils.CtxParam()).Params(jen.PointerTo().ID("WebhookList"), jen.Error()),
				jen.ID("CreateWebhook").Params(utils.CtxParam(), jen.ID("input").PointerTo().ID("WebhookCreationInput")).Params(jen.PointerTo().ID("Webhook"), jen.Error()),
				jen.ID("UpdateWebhook").Params(utils.CtxParam(), jen.ID("updated").PointerTo().ID("Webhook")).Params(jen.Error()),
				jen.ID("ArchiveWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("WebhookDataServer describes a structure capable of serving traffic related to webhooks"),
			jen.ID("WebhookDataServer").Interface(
				jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Line(),
				jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
			),
			jen.Line(),
			jen.Comment("Webhook represents a webhook listener, an endpoint to send an HTTP request to upon an event"),
			jen.ID("Webhook").Struct(
				jen.ID("ID").Uint64().Tag(jsonTag("id")),
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ContentType").String().Tag(jsonTag("content_type")),
				jen.ID("URL").String().Tag(jsonTag("url")),
				jen.ID("Method").String().Tag(jsonTag("method")),
				jen.ID("Events").Index().String().Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().String().Tag(jsonTag("data_types")),
				jen.ID("Topics").Index().String().Tag(jsonTag("topics")),
				jen.ID("CreatedOn").Uint64().Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").PointerTo().Uint64().Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archived_on")),
				jen.ID("BelongsToUser").Uint64().Tag(jsonTag("belongs_to_user")),
			),
			jen.Line(),
			jen.Comment("WebhookCreationInput represents what a user could set as input for creating a webhook"),
			jen.ID("WebhookCreationInput").Struct(
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ContentType").String().Tag(jsonTag("content_type")),
				jen.ID("URL").String().Tag(jsonTag("url")),
				jen.ID("Method").String().Tag(jsonTag("method")),
				jen.ID("Events").Index().String().Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().String().Tag(jsonTag("data_types")),
				jen.ID("Topics").Index().String().Tag(jsonTag("topics")),
				jen.ID("BelongsToUser").Uint64().Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("WebhookUpdateInput represents what a user could set as input for updating a webhook"),
			jen.ID("WebhookUpdateInput").Struct(
				jen.ID("Name").String().Tag(jsonTag("name")),
				jen.ID("ContentType").String().Tag(jsonTag("content_type")),
				jen.ID("URL").String().Tag(jsonTag("url")),
				jen.ID("Method").String().Tag(jsonTag("method")),
				jen.ID("Events").Index().String().Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().String().Tag(jsonTag("data_types")),
				jen.ID("Topics").Index().String().Tag(jsonTag("topics")),
				jen.ID("BelongsToUser").Uint64().Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("WebhookList represents a list of webhooks"),
			jen.ID("WebhookList").Struct(
				jen.ID("Pagination"),
				jen.ID("Webhooks").Index().ID("Webhook").Tag(jsonTag("webhooks")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Update merges an WebhookCreationInput with an Webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("w").PointerTo().ID("Webhook")).ID("Update").Params(jen.ID("input").PointerTo().ID("WebhookUpdateInput")).Block(
			jen.If(jen.ID("input").Dot("Name").DoesNotEqual().EmptyString()).Block(
				jen.ID("w").Dot("Name").Equals().ID("input").Dot("Name"),
			),
			jen.If(jen.ID("input").Dot("ContentType").DoesNotEqual().EmptyString()).Block(
				jen.ID("w").Dot("ContentType").Equals().ID("input").Dot("ContentType"),
			),
			jen.If(jen.ID("input").Dot("URL").DoesNotEqual().EmptyString()).Block(
				jen.ID("w").Dot("URL").Equals().ID("input").Dot("URL"),
			),
			jen.If(jen.ID("input").Dot("Method").DoesNotEqual().EmptyString()).Block(
				jen.ID("w").Dot("Method").Equals().ID("input").Dot("Method"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("Events").DoesNotEqual().ID("nil").And().ID("len").Call(jen.ID("input").Dot("Events")).Op(">").Zero()).Block(
				jen.ID("w").Dot("Events").Equals().ID("input").Dot("Events"),
			),
			jen.If(jen.ID("input").Dot("DataTypes").DoesNotEqual().ID("nil").And().ID("len").Call(jen.ID("input").Dot("DataTypes")).Op(">").Zero()).Block(
				jen.ID("w").Dot("DataTypes").Equals().ID("input").Dot("DataTypes"),
			),
			jen.If(jen.ID("input").Dot("Topics").DoesNotEqual().ID("nil").And().ID("len").Call(jen.ID("input").Dot("Topics")).Op(">").Zero()).Block(
				jen.ID("w").Dot("Topics").Equals().ID("input").Dot("Topics"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErrorLogFunc").Params(jen.ID("w").PointerTo().ID("Webhook"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.Error())).Block(
			jen.Return().Func().Params(jen.Err().Error()).Block(
				jen.ID("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("url").MapAssign().ID("w").Dot("URL"),
					jen.Lit("method").MapAssign().ID("w").Dot("Method"),
					jen.Lit("content_type").MapAssign().ID("w").Dot("ContentType")),
				).Dot("Error").Call(jen.Err(), jen.Lit("error executing webhook")),
			),
		),
		jen.Line(),
	)

	// if proj.EnableNewsman {
	ret.Add(
		jen.Comment("ToListener creates a newsman Listener from a Webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("w").PointerTo().ID("Webhook")).ID("ToListener").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Listener")).Block(
			jen.Return().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewWebhookListener").Callln(
				jen.ID("buildErrorLogFunc").Call(jen.ID("w"), jen.ID("logger")),
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
	)
	// }

	return ret
}
