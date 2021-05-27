package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("webhooksBasePath").Op("=").Lit("webhooks"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetWebhookRequest builds an HTTP request for fetching a webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetWebhookRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("webhookID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("WebhookIDKey"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
				jen.ID("id").Call(jen.ID("webhookID")),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetWebhooksRequest builds an HTTP request for fetching a list of webhooks."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetWebhooksRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("b").Dot("logger")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("webhooksBasePath"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateWebhookRequest builds an HTTP request for creating a webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildCreateWebhookRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("WebhookCreationInput")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("NameKey"),
				jen.ID("input").Dot("Name"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				))),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
			),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateWebhookRequest builds an HTTP request for updating a webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildUpdateWebhookRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("Webhook")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("updated").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("ID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("updated").Dot("ID"),
					jen.Lit(10),
				),
			),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID("updated"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveWebhookRequest builds an HTTP request for archiving a webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildArchiveWebhookRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("webhookID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("WebhookIDKey"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
				jen.ID("id").Call(jen.ID("webhookID")),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogForWebhookRequest builds an HTTP request for fetching a list of audit log entries pertaining to a webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetAuditLogForWebhookRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("webhookID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("WebhookIDKey"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
				jen.ID("id").Call(jen.ID("webhookID")),
				jen.Lit("audit"),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
