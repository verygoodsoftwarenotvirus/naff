package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("redirectToQueryKey").Op("=").Lit("redirectTo"),
			jen.ID("htmxRedirectionHeader").Op("=").Lit("HX-Redirect"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("buildRedirectURL").Params(jen.List(jen.ID("basePath"), jen.ID("redirectTo")).ID("string")).Params(jen.ID("string")).Body(
			jen.ID("u").Op(":=").Op("&").Qual("net/url", "URL").Valuesln(jen.ID("Path").Op(":").ID("basePath"), jen.ID("RawQuery").Op(":").Qual("net/url", "Values").Valuesln(jen.ID("redirectToQueryKey").Op(":").Valuesln(jen.ID("redirectTo"))).Dot("Encode").Call()),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("pluckRedirectURL").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
			jen.Return().ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.ID("redirectToQueryKey"))),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("htmxRedirectTo").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("path").ID("string")).Body(
			jen.ID("res").Dot("Header").Call().Dot("Set").Call(
				jen.ID("htmxRedirectionHeader"),
				jen.ID("path"),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("parseListOfTemplates").Params(jen.ID("funcMap").Qual("html/template", "FuncMap"), jen.ID("name").ID("string"), jen.ID("templates").Op("...").ID("string")).Params(jen.Op("*").Qual("html/template", "Template")).Body(
			jen.ID("tmpl").Op(":=").Qual("html/template", "New").Call(jen.ID("name")).Dot("Funcs").Call(jen.ID("funcMap")),
			jen.For(jen.List(jen.ID("_"), jen.ID("t")).Op(":=").Range().ID("templates")).Body(
				jen.ID("tmpl").Op("=").Qual("html/template", "Must").Call(jen.ID("tmpl").Dot("Parse").Call(jen.ID("t")))),
			jen.Return().ID("tmpl"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("renderStringToResponse").Params(jen.ID("thing").ID("string"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.ID("s").Dot("renderBytesToResponse").Call(
				jen.Index().ID("byte").Call(jen.ID("thing")),
				jen.ID("res"),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("renderBytesToResponse").Params(jen.ID("thing").Index().ID("byte"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("res").Dot("Write").Call(jen.ID("thing")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("writing response"),
				))),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("mergeFuncMaps").Params(jen.List(jen.ID("a"), jen.ID("b")).Qual("html/template", "FuncMap")).Params(jen.Qual("html/template", "FuncMap")).Body(
			jen.ID("out").Op(":=").Map(jen.ID("string")).Interface().Valuesln(),
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("a")).Body(
				jen.ID("out").Index(jen.ID("k")).Op("=").ID("v")),
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("b")).Body(
				jen.ID("out").Index(jen.ID("k")).Op("=").ID("v")),
			jen.Return().ID("out"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("extractFormFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Qual("net/url", "Values"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.List(jen.ID("bodyBytes"), jen.ID("err")).Op(":=").Qual("io", "ReadAll").Call(jen.ID("req").Dot("Body")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading form from request"),
				))),
			jen.List(jen.ID("form"), jen.ID("err")).Op(":=").Qual("net/url", "ParseQuery").Call(jen.ID("string").Call(jen.ID("bodyBytes"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("parsing request form"),
				))),
			jen.Return().List(jen.ID("form"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	return code
}
