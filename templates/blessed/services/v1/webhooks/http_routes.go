package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("URIParamKey is a standard string that we'll use to refer to webhook IDs with"),
			jen.ID("URIParamKey").Equals().Lit("webhookID"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is our list route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ListHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out how specific we need to be"),
				jen.ID(constants.FilterVarName).Assign().Qual(proj.ModelsV1Package(), "ExtractQueryFilter").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out who this is all for"),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Line(),
				jen.Comment("find the webhooks"),
				jen.List(jen.ID("webhooks"), jen.Err()).Assign().ID("s").Dot("webhookDatabase").Dot("GetWebhooks").Call(constants.CtxVar(), jen.ID("userID"), jen.ID(constants.FilterVarName)),
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("webhooks").Equals().AddressOf().Qual(proj.ModelsV1Package(), "WebhookList").Valuesln(
						jen.ID("Webhooks").MapAssign().Index().Qual(proj.ModelsV1Package(), "Webhook").Values()),
				).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered fetching webhooks")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode the response"),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("webhooks")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("validateWebhook does some validation on a WebhookCreationInput and returns an error if anything runs foul"),
		jen.Line(),
		jen.Func().ID("validateWebhook").Params(jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(),
			"WebhookCreationInput",
		)).Params(jen.Error()).Block(
			jen.List(jen.Underscore(), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("input").Dot("URL")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("invalid URL provided: %w"), jen.Err()),
			),
			jen.Line(),
			jen.ID("input").Dot("Method").Equals().Qual("strings", "ToUpper").Call(jen.ID("input").Dot("Method")),
			jen.Switch(jen.ID("input").Dot("Method")).Block(
				jen.Comment("allowed methods"),
				jen.Caseln(
					jen.Qual("net/http", "MethodGet"),
					jen.Qual("net/http", "MethodPost"),
					jen.Qual("net/http", "MethodPut"),
					jen.Qual("net/http", "MethodPatch"),
					jen.Qual("net/http", "MethodDelete"),
					jen.Qual("net/http", "MethodHead"),
				).Block(
					jen.Break(),
				),
				jen.Default().Block(
					jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("invalid method provided: %q"), jen.ID("input").Dot(
						"Method",
					)),
				),
			),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler is our webhook creation route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CreateHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out who this is all for"),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user"), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Line(),
				jen.Comment("try to pluck the parsed input from the request context"),
				jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("CreateMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("valid input not attached to request")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.ID("input").Dot(constants.UserOwnershipFieldName).Equals().ID("userID"),
				jen.Line(),
				jen.Comment("ensure everythings on the up-and-up"),
				jen.If(jen.Err().Assign().ID("validateWebhook").Call(jen.ID("input")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("invalid method provided")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("create the webhook"),
				jen.List(jen.ID("wh"), jen.Err()).Assign().ID("s").Dot("webhookDatabase").Dot("CreateWebhook").Call(constants.CtxVar(), jen.ID("input")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error creating webhook")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify the relevant parties"),
				jen.Qual(proj.InternalTracingV1Package(), "AttachWebhookIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("wh").Dot("ID")),
				jen.ID("s").Dot("webhookCounter").Dot("Increment").Call(constants.CtxVar()),
				jen.ID("s").Dot("eventManager").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Create")),
					jen.ID("Data").MapAssign().ID("wh"),
					jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName"))),
				),
				jen.Line(),
				jen.ID("l").Assign().ID("wh").Dot("ToListener").Call(jen.ID("s").Dot(constants.LoggerVarName)),
				jen.ID("s").Dot("eventManager").Dot("TuneIn").Call(jen.ID("l")),
				jen.Line(),
				jen.Comment("let everybody know we're good"),
				utils.WriteXHeader(constants.ResponseVarName, "StatusCreated"),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("wh")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler returns a GET handler that returns an webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("ReadHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out what this is for and who it belongs to"),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.ID("webhookID").Assign().ID("s").Dot("webhookIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("document it for posterity"),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachWebhookIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("webhookID")),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("user").MapAssign().ID("userID"),
					jen.Lit("webhook").MapAssign().ID("webhookID"),
				)),
				jen.Line(),
				jen.Comment("fetch the webhook from the database"),
				jen.List(jen.ID("x"), jen.Err()).Assign().ID("s").Dot("webhookDatabase").Dot("GetWebhook").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("No rows found in webhookDatabase")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("Error fetching webhook from webhookDatabase")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode the response"),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateHandler returns a handler that updates an webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UpdateHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out what this is for and who it belongs to"),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.ID("webhookID").Assign().ID("s").Dot("webhookIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("document it for posterity"),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachWebhookIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("webhookID")),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("user_id").MapAssign().ID("userID"),
					jen.Lit("webhook_id").MapAssign().ID("webhookID")),
				),
				jen.Line(),
				jen.Comment("fetch parsed creation input from request context"),
				jen.List(jen.ID("input"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("UpdateMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookUpdateInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Info").Call(jen.Lit("no input attached to request")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("fetch the webhook in question"),
				jen.List(jen.ID("wh"), jen.Err()).Assign().ID("s").Dot("webhookDatabase").Dot("GetWebhook").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no rows found for webhook")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered getting webhook")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("update it"),
				jen.ID("wh").Dot("Update").Call(jen.ID("input")),
				jen.Line(),
				jen.Comment("save the update in the database"),
				jen.If(jen.Err().Equals().ID("s").Dot("webhookDatabase").Dot("UpdateWebhook").Call(constants.CtxVar(), jen.ID("wh")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered updating webhook")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify the relevant parties"),
				jen.ID("s").Dot("eventManager").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Update")),
					jen.ID("Data").MapAssign().ID("wh"),
					jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName"))),
				),
				jen.Line(),
				jen.Comment("let everybody know we're good"),
				jen.If(jen.Err().Equals().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("wh")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler returns a handler that archives an webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("delete_route")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("figure out what this is for and who it belongs to"),
				jen.ID("userID").Assign().ID("s").Dot("userIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.ID("webhookID").Assign().ID("s").Dot("webhookIDFetcher").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.Comment("document it for posterity"),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachWebhookIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("webhookID")),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
					jen.Lit("webhook_id").MapAssign().ID("webhookID"),
					jen.Lit("user_id").MapAssign().ID("userID"),
				),
				),
				jen.Line(),
				jen.Comment("do the deed"),
				jen.Err().Assign().ID("s").Dot("webhookDatabase").Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
					jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no rows found for webhook")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered deleting webhook")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("let the interested parties know"),
				jen.ID("s").Dot("webhookCounter").Dot("Decrement").Call(constants.CtxVar()),
				jen.ID("s").Dot("eventManager").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").MapAssign().String().Call(jen.Qual(proj.ModelsV1Package(), "Archive")),
					jen.ID("Data").MapAssign().Qual(proj.ModelsV1Package(), "Webhook").Values(jen.ID("ID").MapAssign().ID("webhookID")),
					jen.ID("Topics").MapAssign().Index().String().Values(jen.ID("topicName"))),
				),
				jen.Line(),
				jen.Comment("let everybody go home"),
				utils.WriteXHeader(constants.ResponseVarName, "StatusNoContent"),
			),
		),
		jen.Line(),
	)

	return ret
}
