package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().ID("accountIDURLParamKey").Op("=").Lit("account"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData")).Params(jen.ID("account").Op("*").ID("types").Dot("Account"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("account").Op("=").ID("fakes").Dot("BuildFakeAccount").Call()).Else().Body(
				jen.List(jen.ID("account"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetAccount").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching account data"),
					))),
			),
			jen.Return().List(jen.ID("account"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("accountEditorTemplate").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildAccountEditorView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
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
				jen.List(jen.ID("account"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAccount").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("s").Dot("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("retrieving account information from database"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("templateFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("componentTitle").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Account")).Params(jen.ID("string")).Body(
						jen.Return().Qual("fmt", "Sprintf").Call(
							jen.Lit("Account #%d"),
							jen.ID("x").Dot("ID"),
						))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("accountEditorTemplate"),
						jen.ID("templateFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Qual("fmt", "Sprintf").Call(
							jen.Lit("Account #%d"),
							jen.ID("account").Dot("ID"),
						), jen.ID("ContentData").Op(":").ID("account")),
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
						jen.Lit(""),
						jen.ID("accountEditorTemplate"),
						jen.ID("templateFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("account"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("fetchAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("accounts").Op("*").ID("types").Dot("AccountList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger"),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.If(jen.ID("s").Dot("useFakeData")).Body(
				jen.ID("accounts").Op("=").ID("fakes").Dot("BuildFakeAccountList").Call()).Else().Body(
				jen.ID("qf").Op(":=").ID("types").Dot("ExtractQueryFilter").Call(jen.ID("req")),
				jen.List(jen.ID("accounts"), jen.ID("err")).Op("=").ID("s").Dot("dataStore").Dot("GetAccounts").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
					jen.ID("qf"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("fetching accounts data"),
					))),
			),
			jen.Return().List(jen.ID("accounts"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("accountsTableTemplate").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildAccountsTableView").Params(jen.ID("includeBaseTemplate").ID("bool")).Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body(
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
				jen.List(jen.ID("accounts"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAccounts").Call(
					jen.ID("ctx"),
					jen.ID("sessionCtxData"),
					jen.ID("req"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("s").Dot("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("retrieving account information from database"),
					),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
					jen.Return(),
				),
				jen.ID("tmplFuncMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
					jen.Lit("individualURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Account")).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/accounts/%d"),
							jen.ID("x").Dot("ID"),
						))), jen.Lit("pushURL").Op(":").Func().Params(jen.ID("x").Op("*").ID("types").Dot("Account")).Params(jen.Qual("html/template", "URL")).Body(
						jen.Return().Qual("html/template", "URL").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/accounts/%d"),
							jen.ID("x").Dot("ID"),
						)))),
				jen.If(jen.ID("includeBaseTemplate")).Body(
					jen.ID("tmpl").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
						jen.ID("accountsTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("page").Op(":=").Op("&").ID("pageData").Valuesln(
						jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Lit("Accounts"), jen.ID("ContentData").Op(":").ID("accounts")),
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
						jen.ID("accountsTableTemplate"),
						jen.ID("tmplFuncMap"),
					),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("accounts"),
						jen.ID("res"),
					),
				),
			)),
		jen.Line(),
	)

	return code
}
