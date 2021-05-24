package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().ID("webhookIDURLParamKey").Op("=").Lit("webhook"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("webhook").Op("*").ID("types").Dot("Webhook"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("webhook").Op("=").ID("fakes").Dot("BuildFakeWebhook").Call()).Else().Body(
				jen.ID("webhookID").Op(":=").ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("webhookIDURLParamKey"),
					jen.Lit("webhook"),
				).Call(jen.ID("req")),
				jen.List(jen.ID("webhook"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetWebhook").Call(
					jen.ID("ctx"),
					jen.ID("webhookID"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching webhook data"),
					))),
			),
			jen.Return().List(jen.ID("webhook"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("webhookEditorTemplate").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildWebhookEditorView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("no session context data attached to request"),
					),
					jen.Qual("net/http", "Redirect").Call(
						jen.ID("res"),
						jen.ID("req"),
						jen.Lit("/login"),
						jen.ID("unauthorizedRedirectResponseCode"),
					),
					jen.Return(),
				),
				jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("s").Dot("fetchWebhook").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching webhook from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Webhook")).Params(jen.ID("string")).Body(
						jen.Return().Qual("fmt", "Sprintf").Call(
							jen.Lit("Webhook #%d"),
							jen.ID("x").Dot("ID"),
						))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("webhookEditorTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Qual("fmt", "Sprintf").Call(
							jen.Lit("Webhook #%d"),
							jen.ID("webhook").Dot("ID"),
						), jen.ID("ContentData").Op(":").ID("webhook")),
					jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("view"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("webhookEditorTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("webhook"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("webhooks").Op("*").ID("types").Dot("WebhookList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("webhooks").Op("=").ID("fakes").Dot("BuildFakeWebhookList").Call()).Else().Body(
				jen.ID("filter").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
				jen.List(jen.ID("webhooks"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetWebhooks").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
					jen.ID("filter"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching webhook data"),
					))),
			),
			jen.Return().List(jen.ID("webhooks"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("webhooksTableTemplate").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildWebhooksTableView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
			jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("tracing").Dot("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("req"),
				),
				jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("no session context data attached to request"),
					),
					jen.Qual("net/http", "Redirect").Call(
						jen.ID("res"),
						jen.ID("req"),
						jen.Lit("/login"),
						jen.ID("unauthorizedRedirectResponseCode"),
					),
					jen.Return(),
				),
				jen.List(jen.ID("webhooks"), jen.ID("err")).Op(":=").ID("s").Dot("fetchWebhooks").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching webhooks from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Webhook")).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/webhooks/%d"),
							jen.ID("x").Dot("ID"),
						))), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Webhook")).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/webhooks/%d"),
							jen.ID("x").Dot("ID"),
						)))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("webhooksTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Lit("Webhooks"), jen.ID("ContentData").Op(":").ID("webhooks")),
					jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("view"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit("dashboard"),
						jen.ID("webhooksTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("webhooks"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	return code
}
