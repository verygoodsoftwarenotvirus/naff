package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func baseTemplateDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Type().ID("pageData").Struct(
			jen.ID("ContentData").Interface(),
			jen.ID("Title").ID("string"),
			jen.ID("PageDescription").ID("string"),
			jen.ID("PageTitle").ID("string"),
			jen.ID("PageImagePreview").ID("string"),
			jen.ID("PageImagePreviewDescription").ID("string"),
			jen.ID("InheritedQuery").ID("string"),
			jen.ID("IsLoggedIn").ID("bool"),
			jen.ID("IsServiceAdmin").ID("bool"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("baseTemplateSrc").ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("homepage").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("_").Op("=").ID("err")),
			jen.ID("tmpl").Op(":=").ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
				jen.Lit(""),
				jen.ID("nil"),
			),
			jen.ID("x").Op(":=").Op("&").ID("pageData").Valuesln(
				jen.ID("IsLoggedIn").Op(":").ID("sessionCtxData").Op("!=").ID("nil"), jen.ID("Title").Op(":").Lit("Home"), jen.ID("ContentData").Op(":").Lit("")),
			jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
				jen.ID("x").Dot("IsServiceAdmin").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call()),
			jen.ID("s").Dot("renderTemplateToResponse").Call(
				jen.ID("ctx"),
				jen.ID("tmpl"),
				jen.ID("x"),
				jen.ID("res"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("wrapTemplateInContentDefinition").Params(jen.ID("tmpl").ID("string")).Params(jen.ID("string")).Body(
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.Lit(`{{ define "content" }}
	%s
{{ end }}
`),
				jen.ID("tmpl"),
			)),
		jen.Line(),
	)

	return code
}
