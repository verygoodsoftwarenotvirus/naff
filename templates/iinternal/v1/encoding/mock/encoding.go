package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(buildInterfaceImplementationDeclaration(proj)...)
	code.Add(buildEncoderDecoder()...)
	code.Add(buildEncodeResponse()...)
	code.Add(buildDecodeRequest()...)

	return code
}

func buildInterfaceImplementationDeclaration(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual(proj.InternalEncodingV1Package(), "EncoderDecoder").Equals().Parens(jen.PointerTo().ID("EncoderDecoder")).Call(jen.Nil()),
		jen.Line(),
	}

	return lines
}

func buildEncoderDecoder() []jen.Code {
	lines := []jen.Code{
		jen.Comment("EncoderDecoder is a mock EncoderDecoder."),
		jen.Line(),
		jen.Type().ID("EncoderDecoder").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildEncodeResponse() []jen.Code {
	lines := []jen.Code{
		jen.Comment("EncodeResponse satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("EncoderDecoder")).ID("EncodeResponse").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID("v")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildDecodeRequest() []jen.Code {
	lines := []jen.Code{
		jen.Comment("DecodeRequest satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("EncoderDecoder")).ID("DecodeRequest").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID(constants.RequestVarName), jen.ID("v")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
