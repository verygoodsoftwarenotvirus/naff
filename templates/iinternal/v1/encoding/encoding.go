package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("encoding")

	utils.AddImports(proj, code)

	code.Add(buildEncodingConstDeclarations()...)
	code.Add(buildEncodingVarDeclarations()...)
	code.Add(buildEncodingTypeDeclarations()...)
	code.Add(buildEncodeResponse()...)
	code.Add(buildDecodeRequest()...)
	code.Add(buildProvideResponseEncoder()...)

	return code
}

func buildEncodingConstDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("ContentTypeHeader is the HTTP standard header name for content type."),
			jen.ID("ContentTypeHeader").Equals().Lit("Content-type"),
			jen.Comment("XMLContentType represents the XML content type."),
			jen.ID("XMLContentType").Equals().Lit("application/xml"),
			jen.Comment("JSONContentType represents the JSON content type."),
			jen.ID("JSONContentType").Equals().Lit("application/json"),
			jen.Comment("DefaultContentType is what the library defaults to."),
			jen.ID("DefaultContentType").Equals().ID("JSONContentType"),
		),
		jen.Line(),
	}

	return lines
}

func buildEncodingVarDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("Providers provides ResponseEncoders for dependency injection."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideResponseEncoder"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildEncodingTypeDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("EncoderDecoder is an interface that allows for multiple implementations of HTTP response formats."),
			jen.ID("EncoderDecoder").Interface(
				jen.ID("EncodeResponse").Params(jen.Qual("net/http", "ResponseWriter"), jen.Interface()).Params(jen.Error()),
				jen.ID("DecodeRequest").Params(jen.PointerTo().Qual("net/http", "Request"), jen.Interface()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("ServerEncoderDecoder is our concrete implementation of EncoderDecoder."),
			jen.ID("ServerEncoderDecoder").Struct(),
			jen.Line(),
			jen.ID("encoder").Interface(
				jen.ID("Encode").Params(jen.ID("v").Interface()).Params(jen.Error()),
			),
			jen.Line(),
			jen.ID("decoder").Interface(
				jen.ID("Decode").Params(jen.ID("v").Interface()).Params(jen.Error()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildEncodeResponse() []jen.Code {
	lines := []jen.Code{
		jen.Comment("EncodeResponse encodes responses."),
		jen.Line(),
		jen.Func().Params(jen.ID("ed").PointerTo().ID("ServerEncoderDecoder")).ID("EncodeResponse").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Params(jen.Error()).Block(
			jen.Var().ID("ct").Equals().Qual("strings", "ToLower").Call(jen.ID(constants.ResponseVarName).Dot("Header").Call().Dot("Get").Call(jen.ID("ContentTypeHeader"))),
			jen.If(jen.ID("ct").IsEqualTo().EmptyString()).Block(
				jen.ID("ct").Equals().ID("DefaultContentType"),
			),
			jen.Line(),
			jen.Var().ID("e").ID("encoder"),
			jen.Switch(jen.ID("ct")).Block(
				jen.Case(jen.ID("XMLContentType")).Block(jen.ID("e").Equals().Qual("encoding/xml", "NewEncoder").Call(jen.ID(constants.ResponseVarName))),
				jen.Default().Block(jen.ID("e").Equals().Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName))),
			),
			jen.Line(),
			jen.ID(constants.ResponseVarName).Dot("Header").Call().Dot("Set").Call(jen.ID("ContentTypeHeader"), jen.ID("ct")),
			jen.Return().ID("e").Dot("Encode").Call(jen.ID("v")),
		),
		jen.Line(),
	}

	return lines
}

func buildDecodeRequest() []jen.Code {
	lines := []jen.Code{
		jen.Comment("DecodeRequest decodes responses."),
		jen.Line(),
		jen.Func().Params(jen.ID("ed").PointerTo().ID("ServerEncoderDecoder")).ID("DecodeRequest").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.Error()).Block(
			jen.Var().ID("ct").Equals().Qual("strings", "ToLower").Call(jen.ID(constants.RequestVarName).Dot("Header").Dot("Get").Call(jen.ID("ContentTypeHeader"))),
			jen.If(jen.ID("ct").IsEqualTo().EmptyString()).Block(
				jen.ID("ct").Equals().ID("DefaultContentType"),
			),
			jen.Line(),
			jen.Var().ID("d").ID("decoder"),
			jen.Switch(jen.ID("ct")).Block(
				jen.Case(jen.ID("XMLContentType")).Block(
					jen.ID("d").Equals().Qual("encoding/xml", "NewDecoder").Call(jen.ID(constants.RequestVarName).Dot("Body")),
				),
				jen.Default().Block(
					jen.ID("d").Equals().Qual("encoding/json", "NewDecoder").Call(jen.ID(constants.RequestVarName).Dot("Body")),
				),
			),
			jen.Line(),
			jen.Return().ID("d").Dot("Decode").Call(jen.ID("v")),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideResponseEncoder() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideResponseEncoder provides a jsonResponseEncoder."),
		jen.Line(),
		jen.Func().ID("ProvideResponseEncoder").Params().Params(jen.ID("EncoderDecoder")).Block(
			jen.Return().AddressOf().ID("ServerEncoderDecoder").Values(),
		),
		jen.Line(),
	}

	return lines
}
