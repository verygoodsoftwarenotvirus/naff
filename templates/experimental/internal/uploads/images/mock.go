package images

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("ImageUploadProcessor").Op("=").Parens(jen.Op("*").ID("MockImageUploadProcessor")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockImageUploadProcessor").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Process satisfies the ImageUploadProcessor interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockImageUploadProcessor")).ID("Process").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("filename"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("Image")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAvatarUploadMiddleware satisfies the ImageUploadProcessor interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockImageUploadProcessor")).ID("BuildAvatarUploadMiddleware").Params(jen.ID("next").Qual("net/http", "Handler"), jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("filename").ID("string")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("next"),
				jen.ID("encoderDecoder"),
				jen.ID("filename"),
			),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	return code
}
