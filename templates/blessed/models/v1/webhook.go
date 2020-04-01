package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("WebhookDataManager describes a structure capable of storing webhooks"),
			jen.ID("WebhookDataManager").Interface(
				jen.ID("GetWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("Webhook"), jen.ID("error")),
				jen.ID("GetWebhookCount").Params(utils.CtxParam(), jen.ID("userID").ID("uint64"), jen.ID(utils.FilterVarName).Op("*").ID("QueryFilter")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetAllWebhooksCount").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetWebhooks").Params(utils.CtxParam(), jen.ID("userID").ID("uint64"), jen.ID(utils.FilterVarName).Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("WebhookList"), jen.ID("error")),
				jen.ID("GetAllWebhooks").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.Op("*").ID("WebhookList"), jen.ID("error")),
				jen.ID("GetAllWebhooksForUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("Webhook"), jen.ID("error")),
				jen.ID("CreateWebhook").Params(utils.CtxParam(), jen.ID("input").Op("*").ID("WebhookCreationInput")).Params(jen.Op("*").ID("Webhook"), jen.ID("error")),
				jen.ID("UpdateWebhook").Params(utils.CtxParam(), jen.ID("updated").Op("*").ID("Webhook")).Params(jen.ID("error")),
				jen.ID("ArchiveWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")),
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
				jen.ID("ID").ID("uint64").Tag(jsonTag("id")),
				jen.ID("Name").ID("string").Tag(jsonTag("name")),
				jen.ID("ContentType").ID("string").Tag(jsonTag("content_type")),
				jen.ID("URL").ID("string").Tag(jsonTag("url")),
				jen.ID("Method").ID("string").Tag(jsonTag("method")),
				jen.ID("Events").Index().ID("string").Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().ID("string").Tag(jsonTag("data_types")),
				jen.ID("Topics").Index().ID("string").Tag(jsonTag("topics")),
				jen.ID("CreatedOn").ID("uint64").Tag(jsonTag("created_on")),
				jen.ID("UpdatedOn").Op("*").ID("uint64").Tag(jsonTag("updated_on")),
				jen.ID("ArchivedOn").Op("*").ID("uint64").Tag(jsonTag("archived_on")),
				jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("belongs_to_user")),
			),
			jen.Line(),
			jen.Comment("WebhookCreationInput represents what a user could set as input for creating a webhook"),
			jen.ID("WebhookCreationInput").Struct(
				jen.ID("Name").ID("string").Tag(jsonTag("name")),
				jen.ID("ContentType").ID("string").Tag(jsonTag("content_type")),
				jen.ID("URL").ID("string").Tag(jsonTag("url")),
				jen.ID("Method").ID("string").Tag(jsonTag("method")),
				jen.ID("Events").Index().ID("string").Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().ID("string").Tag(jsonTag("data_types")),
				jen.ID("Topics").Index().ID("string").Tag(jsonTag("topics")),
				jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("WebhookUpdateInput represents what a user could set as input for updating a webhook"),
			jen.ID("WebhookUpdateInput").Struct(
				jen.ID("Name").ID("string").Tag(jsonTag("name")),
				jen.ID("ContentType").ID("string").Tag(jsonTag("content_type")),
				jen.ID("URL").ID("string").Tag(jsonTag("url")),
				jen.ID("Method").ID("string").Tag(jsonTag("method")),
				jen.ID("Events").Index().ID("string").Tag(jsonTag("events")),
				jen.ID("DataTypes").Index().ID("string").Tag(jsonTag("data_types")),
				jen.ID("Topics").Index().ID("string").Tag(jsonTag("topics")),
				jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("-")),
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
		jen.Func().Params(jen.ID("w").Op("*").ID("Webhook")).ID("Update").Params(jen.ID("input").Op("*").ID("WebhookUpdateInput")).Block(
			jen.If(jen.ID("input").Dot("Name").DoesNotEqual().Lit("")).Block(
				jen.ID("w").Dot("Name").Equals().ID("input").Dot("Name"),
			),
			jen.If(jen.ID("input").Dot("ContentType").DoesNotEqual().Lit("")).Block(
				jen.ID("w").Dot("ContentType").Equals().ID("input").Dot("ContentType"),
			),
			jen.If(jen.ID("input").Dot("URL").DoesNotEqual().Lit("")).Block(
				jen.ID("w").Dot("URL").Equals().ID("input").Dot("URL"),
			),
			jen.If(jen.ID("input").Dot("Method").DoesNotEqual().Lit("")).Block(
				jen.ID("w").Dot("Method").Equals().ID("input").Dot("Method"),
			),
			jen.Line(),
			jen.If(jen.ID("input").Dot("Events").DoesNotEqual().ID("nil").Op("&&").ID("len").Call(jen.ID("input").Dot("Events")).Op(">").Lit(0)).Block(
				jen.ID("w").Dot("Events").Equals().ID("input").Dot("Events"),
			),
			jen.If(jen.ID("input").Dot("DataTypes").DoesNotEqual().ID("nil").Op("&&").ID("len").Call(jen.ID("input").Dot("DataTypes")).Op(">").Lit(0)).Block(
				jen.ID("w").Dot("DataTypes").Equals().ID("input").Dot("DataTypes"),
			),
			jen.If(jen.ID("input").Dot("Topics").DoesNotEqual().ID("nil").Op("&&").ID("len").Call(jen.ID("input").Dot("Topics")).Op(">").Lit(0)).Block(
				jen.ID("w").Dot("Topics").Equals().ID("input").Dot("Topics"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErrorLogFunc").Params(jen.ID("w").Op("*").ID("Webhook"), jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.ID("error"))).Block(
			jen.Return().Func().Params(jen.Err().ID("error")).Block(
				jen.ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("url").MapAssign().ID("w").Dot("URL"),
					jen.Lit("method").MapAssign().ID("w").Dot("Method"),
					jen.Lit("content_type").MapAssign().ID("w").Dot("ContentType")),
				).Dot("Error").Call(jen.Err(), jen.Lit("error executing webhook")),
			),
		),
		jen.Line(),
	)

	// if pkg.EnableNewsman {
	ret.Add(
		jen.Comment("ToListener creates a newsman Listener from a Webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("w").Op("*").ID("Webhook")).ID("ToListener").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Listener")).Block(
			jen.Return().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewWebhookListener").Callln(
				jen.ID("buildErrorLogFunc").Call(jen.ID("w"), jen.ID("logger")),
				jen.VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "WebhookConfig").Valuesln(
					jen.ID("Method").MapAssign().ID("w").Dot("Method"),
					jen.ID("URL").MapAssign().ID("w").Dot("URL"),
					jen.ID("ContentType").MapAssign().ID("w").Dot("ContentType"),
				),
				jen.VarPointer().Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "ListenerConfig").Valuesln(
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
