package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("uploads").Dot("UploadManager").Op("=").Parens(jen.Op("*").ID("UploadManager")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("UploadManager").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SaveFile satisfies the UploadManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UploadManager")).ID("SaveFile").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("path").ID("string"), jen.ID("content").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("path"),
				jen.ID("content"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadFile satisfies the UploadManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UploadManager")).ID("ReadFile").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("path").ID("string")).Params(jen.Index().ID("byte"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("path"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("byte")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ServeFiles satisfies the UploadManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UploadManager")).ID("ServeFiles").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	return code
}
