package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func staticAssetsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().ID("svgFaviconSrc").Index().ID("byte"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("favicon").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.ID("res").Dot("Header").Call().Dot("Set").Call(
				jen.Lit("Content-Type"),
				jen.Lit("image/svg+xml"),
			),
			jen.ID("s").Dot("renderBytesToResponse").Call(
				jen.ID("svgFaviconSrc"),
				jen.ID("res"),
			),
		),
		jen.Line(),
	)

	return code
}
