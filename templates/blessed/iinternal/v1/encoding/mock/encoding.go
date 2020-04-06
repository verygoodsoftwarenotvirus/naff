package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.InternalEncodingV1Package(), "EncoderDecoder").Equals().Parens(jen.PointerTo().ID("EncoderDecoder")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("EncoderDecoder is a mock EncoderDecoder"),
		jen.Line(),
		jen.Type().ID("EncoderDecoder").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("EncodeResponse satisfies our EncoderDecoder interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("EncoderDecoder")).ID("EncodeResponse").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("v")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("DecodeRequest satisfies our EncoderDecoder interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("EncoderDecoder")).ID("DecodeRequest").Params(jen.ID("req").ParamPointer().Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("req"), jen.ID("v")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)
	return ret
}
