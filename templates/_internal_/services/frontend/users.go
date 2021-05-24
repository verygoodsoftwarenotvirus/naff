package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("users").Op("*").ID("types").Dot("UserList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("users").Op("=").ID("fakes").Dot("BuildFakeUserList").Call()).Else().Body(
				jen.ID("filter").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
				jen.List(jen.ID("users"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetUsers").Call(
					jen.ID("ctx"),
					jen.ID("filter"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching user data"),
					))),
			),
			jen.Return().List(jen.ID("users"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("usersTableTemplate").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildUsersTableView").Params(jen.List(jen.ID("includeBaseTemplate"), jen.ID("forSearch")).ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
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
				jen.Var().ID("users").Op("*").ID("types").Dot("UserList"),
				jen.If(jen.ID("forSearch")).Body(
					jen.ID("query").Op(":=").ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.ID("types").Dot("SearchQueryKey")),
					jen.List(jen.ID("searchResults"), jen.ID("searchResultsErr")).Op(":=").ID("s").Dot("dataStore").Dot("SearchForUsersByUsername").Call(
						jen.ID("ctx"),
						jen.ID("query"),
					),
					jen.If(jen.ID("searchResultsErr").Op("!=").ID("nil")).Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("searchResultsErr"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("fetching users from datastore"),
						),
						jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
						jen.Return(),
					),
					jen.ID("users").Op("=").Op("&").ID("types").Dot("UserList").Valuesln(
						jen.ID("Users").Op(":").ID("searchResults")),
				).Else().Body(
					jen.List(jen.ID("users"), jen.ID("err")).Op("=").ID("s").Dot("fetchUsers").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("fetching users from datastore"),
						),
						jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
						jen.Return(),
					),
				),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("usersTableTemplate"),
						jen.ID("nil"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Lit("Users"), jen.ID("ContentData").Op(":").ID("users")),
					jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
						jen.ID("page").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("page"),
						jen.ID("res"),
					),
				).Else().Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit("dashboard"),
						jen.ID("usersTableTemplate"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("users"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	return code
}
