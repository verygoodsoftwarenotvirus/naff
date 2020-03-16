package webhooks

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("URIParamKey is a standard string that we'll use to refer to webhook IDs with"),
			jen.ID("URIParamKey").Op("=").Lit("webhookID"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachWebhookIDToSpan provides a consistent way to attach a webhook ID to a given span"),
		jen.Line(),
		jen.Func().ID("attachWebhookIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("webhookID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot(
					"AddAttributes",
				).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("webhook_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("webhookID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("attachUserIDToSpan provides a consistent way to attach a user ID to a given span"),
		jen.Line(),
		jen.Func().ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
				jen.ID("span").Dot(
					"AddAttributes",
				).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is our list route"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ListHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("figure out how specific we need to be"),
				jen.ID("qf").Op(":=").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ExtractQueryFilter").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("figure out who this is all for"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("find the webhooks"),
				jen.List(jen.ID("webhooks"), jen.ID("err")).Op(":=").ID("s").Dot("webhookDatabase").Dot("GetWebhooks").Call(jen.ID("ctx"), jen.ID("qf"), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("webhooks").Op("=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList").Valuesln(
						jen.ID("Webhooks").Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Values()),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered fetching webhooks")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode the response"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("webhooks")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("validateWebhook does some validation on a WebhookCreationInput and returns an error if anything runs foul"),
		jen.Line(),
		jen.Func().ID("validateWebhook").Params(jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"),
			"WebhookCreationInput",
		)).Params(jen.ID("error")).Block(
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("input").Dot("URL")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("invalid URL provided: %w"), jen.ID("err")),
			),
			jen.Line(),
			jen.ID("input").Dot("Method").Op("=").Qual("strings", "ToUpper").Call(jen.ID("input").Dot("Method")),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreateHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("figure out who this is all for"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("user"), jen.ID("userID")),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.Line(),
				jen.Comment("try to pluck the parsed input from the request context"),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot("Value").Call(jen.ID("CreateMiddlewareCtxKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("logger").Dot("Info").Call(jen.Lit("valid input not attached to request")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.ID("input").Dot("BelongsToUser").Op("=").ID("userID"),
				jen.Line(),
				jen.Comment("ensure everythings on the up-and-up"),
				jen.If(jen.ID("err").Op(":=").ID("validateWebhook").Call(jen.ID("input")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Info").Call(jen.Lit("invalid method provided")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("create the webhook"),
				jen.List(jen.ID("wh"), jen.ID("err")).Op(":=").ID("s").Dot("webhookDatabase").Dot("CreateWebhook").Call(jen.ID("ctx"), jen.ID("input")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error creating webhook")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify the relevant parties"),
				jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("wh").Dot("ID")),
				jen.ID("s").Dot("webhookCounter").Dot("Increment").Call(jen.ID("ctx")),
				jen.ID("s").Dot("eventManager").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").Op(":").ID("string").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Create")),
					jen.ID("Data").Op(":").ID("wh"),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName"))),
				),
				jen.Line(),
				jen.ID("l").Op(":=").ID("wh").Dot("ToListener").Call(jen.ID("s").Dot("logger")),
				jen.ID("s").Dot("eventManager").Dot("TuneIn").Call(jen.ID("l")),
				jen.Line(),
				jen.Comment("let everybody know we're good"),
				utils.WriteXHeader("res", "StatusCreated"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("wh")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler returns a GET handler that returns an webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("ReadHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("figure out what this is for and who it belongs to"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("webhookID").Op(":=").ID("s").Dot("webhookIDFetcher").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("document it for posterity"),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("user").Op(":").ID("userID"),
					jen.Lit("webhook").Op(":").ID("webhookID"),
				)),
				jen.Line(),
				jen.Comment("fetch the webhook from the database"),
				jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dot("webhookDatabase").Dot("GetWebhook").Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("No rows found in webhookDatabase")),
					utils.WriteXHeader("res", "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("Error fetching webhook from webhookDatabase")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("encode the response"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateHandler returns a handler that updates an webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UpdateHandler")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("figure out what this is for and who it belongs to"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("webhookID").Op(":=").ID("s").Dot("webhookIDFetcher").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("document it for posterity"),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("user_id").Op(":").ID("userID"),
					jen.Lit("webhook_id").Op(":").ID("webhookID")),
				),
				jen.Line(),
				jen.Comment("fetch parsed creation input from request context"),
				jen.List(jen.ID("input"), jen.ID("ok")).Op(":=").ID("ctx").Dot("Value").Call(jen.ID("UpdateMiddlewareCtxKey")).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookUpdateInput")),
				jen.If(jen.Op("!").ID("ok")).Block(
					jen.ID("s").Dot("logger").Dot("Info").Call(jen.Lit("no input attached to request")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("fetch the webhook in question"),
				jen.List(jen.ID("wh"), jen.ID("err")).Op(":=").ID("s").Dot("webhookDatabase").Dot("GetWebhook").Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("no rows found for webhook")),
					utils.WriteXHeader("res", "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered getting webhook")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("update it"),
				jen.ID("wh").Dot("Update").Call(jen.ID("input")),
				jen.Line(),
				jen.Comment("save the update in the database"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("webhookDatabase").Dot("UpdateWebhook").Call(jen.ID("ctx"), jen.ID("wh")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered updating webhook")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("notify the relevant parties"),
				jen.ID("s").Dot("eventManager").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").Op(":").ID("string").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Update")),
					jen.ID("Data").Op(":").ID("wh"),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName"))),
				),
				jen.Line(),
				jen.Comment("let everybody know we're good"),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID("res"), jen.ID("wh")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("encoding response")),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler returns a handler that archives an webhook"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("delete_route")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("figure out what this is for and who it belongs to"),
				jen.ID("userID").Op(":=").ID("s").Dot("userIDFetcher").Call(jen.ID("req")),
				jen.ID("webhookID").Op(":=").ID("s").Dot("webhookIDFetcher").Call(jen.ID("req")),
				jen.Line(),
				jen.Comment("document it for posterity"),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
				jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("webhook_id").Op(":").ID("webhookID"),
					jen.Lit("user_id").Op(":").ID("userID"),
				),
				),
				jen.Line(),
				jen.Comment("do the deed"),
				jen.ID("err").Op(":=").ID("s").Dot("webhookDatabase").Dot("ArchiveWebhook").Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("no rows found for webhook")),
					utils.WriteXHeader("res", "StatusNotFound"),
					jen.Return(),
				).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered deleting webhook")),
					utils.WriteXHeader("res", "StatusInternalServerError"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("let the interested parties know"),
				jen.ID("s").Dot("webhookCounter").Dot("Decrement").Call(jen.ID("ctx")),
				jen.ID("s").Dot("eventManager").Dot("Report").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Event").Valuesln(
					jen.ID("EventType").Op(":").ID("string").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Archive")),
					jen.ID("Data").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Values(jen.ID("ID").Op(":").ID("webhookID")),
					jen.ID("Topics").Op(":").Index().ID("string").Values(jen.ID("topicName"))),
				),
				jen.Line(),
				jen.Comment("let everybody go home"),
				utils.WriteXHeader("res", "StatusNoContent"),
			),
		),
		jen.Line(),
	)
	return ret
}
