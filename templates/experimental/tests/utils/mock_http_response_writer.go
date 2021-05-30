package utils

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockHTTPResponseWriterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("net/http", "ResponseWriter").Op("=").Parens(jen.Op("*").ID("MockHTTPResponseWriter")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockHTTPResponseWriter").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Header satisfies our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPResponseWriter")).ID("Header").Params().Params(jen.Qual("net/http", "Header")).Body(
			jen.Return().ID("m").Dot("Called").Call().Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Header"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Write satisfies our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPResponseWriter")).ID("Write").Params(jen.ID("in").Index().ID("byte")).Params(jen.ID("int"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(jen.ID("in")),
			jen.Return().List(jen.ID("returnValues").Dot("Int").Call(jen.Lit(0)), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("WriteHeader satisfies our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPResponseWriter")).ID("WriteHeader").Params(jen.ID("statusCode").ID("int")).Body(
			jen.ID("m").Dot("Called").Call(jen.ID("statusCode"))),
		jen.Line(),
	)

	return code
}
