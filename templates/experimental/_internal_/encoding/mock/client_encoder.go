package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientEncoderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("ClientEncoder").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ContentType satisfies the ClientEncoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ClientEncoder")).ID("ContentType").Params().Params(jen.ID("string")).Body(
			jen.Return().ID("m").Dot("Called").Call().Dot("String").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Unmarshal satisfies the ClientEncoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ClientEncoder")).ID("Unmarshal").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Index().ID("byte"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("data"),
				jen.ID("v"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Encode satisfies the ClientEncoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ClientEncoder")).ID("Encode").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("dest").Qual("io", "Writer"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("dest"),
				jen.ID("v"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeReader satisfies the ClientEncoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("ClientEncoder")).ID("EncodeReader").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Interface()).Params(jen.Qual("io", "Reader"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("data"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("io", "Reader")), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
