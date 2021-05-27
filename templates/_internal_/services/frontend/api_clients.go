package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().ID("apiClientIDURLParamKey").Op("=").Lit("api_client"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("apiClient").Op("*").ID("types").Dot("APIClient"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("apiClient").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call()).Else().Body(
				jen.ID("apiClientID").Op(":=").ID("s").Dot("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("apiClientIDURLParamKey"),
					jen.Lit("API client"),
				).Call(jen.ID("req")),
				jen.List(jen.ID("apiClient"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetAPIClientByDatabaseID").Call(
					jen.ID("ctx"),
					jen.ID("apiClientID"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching API client data"),
					))),
			),
			jen.Return().List(jen.ID("apiClient"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("apiClientEditorTemplate").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildAPIClientEditorView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
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
				jen.List(jen.ID("apiClient"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClient").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching item from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("APIClient")).Params(jen.ID("string")).Body(
						jen.Return().Qual("fmt", "Sprintf").Call(
							jen.Lit("Client #%d"),
							jen.ID("x").Dot("ID"),
						))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("apiClientEditorTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("apiClient"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("view").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("apiClientEditorTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Qual("fmt", "Sprintf").Call(
							jen.Lit("APIClient #%d"),
							jen.ID("apiClient").Dot("ID"),
						), jen.ID("ContentData").Op(":").ID("apiClient")),
					jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("view"),
						jen.ID("page"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchAPIClients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("apiClients").Op("*").ID("types").Dot("APIClientList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("apiClients").Op("=").ID("fakes").Dot("BuildFakeAPIClientList").Call()).Else().Body(
				jen.ID("filter").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
				jen.List(jen.ID("apiClients"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetAPIClients").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
					jen.ID("filter"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching API client data"),
					))),
			),
			jen.Return().List(jen.ID("apiClients"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("apiClientsTableTemplate").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildAPIClientsTableView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
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
				jen.List(jen.ID("apiClients"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClients").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching API client from datastore"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("APIClient")).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/api_clients/%d"),
							jen.ID("x").Dot("ID"),
						))), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("APIClient")).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/api_clients/%d"),
							jen.ID("x").Dot("ID"),
						)))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("view").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("apiClientsTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Lit("APIClients"), jen.ID("ContentData").Op(":").ID("apiClients")),
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
						jen.ID("apiClientsTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("apiClients"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	return code
}
