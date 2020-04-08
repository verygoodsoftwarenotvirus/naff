package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("encoding")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("ContentTypeHeader is the HTTP standard header name for content type"),
			jen.ID("ContentTypeHeader").Equals().Lit("Content-type"),
			jen.Comment("XMLContentType represents the XML content type"),
			jen.ID("XMLContentType").Equals().Lit("application/xml"),
			jen.Comment("JSONContentType represents the JSON content type"),
			jen.ID("JSONContentType").Equals().Lit("application/json"),
			jen.Comment("DefaultContentType is what the library defaults to"),
			jen.ID("DefaultContentType").Equals().ID("JSONContentType"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers provides ResponseEncoders for dependency injection"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID("ProvideResponseEncoder"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("EncoderDecoder is an interface that allows for multiple implementations of HTTP response formats"),
			jen.ID("EncoderDecoder").Interface(
				jen.ID("EncodeResponse").Params(jen.Qual("net/http", "ResponseWriter"), jen.Interface()).Params(jen.Error()),
				jen.ID("DecodeRequest").Params(jen.ParamPointer().Qual("net/http", "Request"), jen.Interface()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("ServerEncoderDecoder is our concrete implementation of EncoderDecoder"),
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
	)

	ret.Add(
		jen.Comment("EncodeResponse encodes responses"),
		jen.Line(),
		jen.Func().Params(jen.ID("ed").PointerTo().ID("ServerEncoderDecoder")).ID("EncodeResponse").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Params(jen.Error()).Block(
			jen.Var().ID("ct").Equals().Qual("strings", "ToLower").Call(jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("ContentTypeHeader"))),
			jen.If(jen.ID("ct").Op("==").EmptyString()).Block(
				jen.ID("ct").Equals().ID("DefaultContentType"),
			),
			jen.Line(),
			jen.Var().ID("e").ID("encoder"),
			jen.Switch(jen.ID("ct")).Block(
				jen.Case(jen.ID("XMLContentType")).Block(jen.ID("e").Equals().Qual("encoding/xml", "NewEncoder").Call(jen.ID("res"))),
				jen.Default().Block(jen.ID("e").Equals().Qual("encoding/json", "NewEncoder").Call(jen.ID("res"))),
			),
			jen.Line(),
			jen.ID("res").Dot("Header").Call().Dot("Set").Call(jen.ID("ContentTypeHeader"), jen.ID("ct")),
			jen.Return().ID("e").Dot("Encode").Call(jen.ID("v")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("DecodeRequest decodes responses"),
		jen.Line(),
		jen.Func().Params(jen.ID("ed").PointerTo().ID("ServerEncoderDecoder")).ID("DecodeRequest").Params(jen.ID("req").ParamPointer().Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.Error()).Block(
			jen.Var().ID("ct").Equals().Qual("strings", "ToLower").Call(jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("ContentTypeHeader"))),
			jen.If(jen.ID("ct").Op("==").EmptyString()).Block(
				jen.ID("ct").Equals().ID("DefaultContentType"),
			),
			jen.Line(),
			jen.Var().ID("d").ID("decoder"),
			jen.Switch(jen.ID("ct")).Block(
				jen.Case(jen.ID("XMLContentType")).Block(
					jen.ID("d").Equals().Qual("encoding/xml", "NewDecoder").Call(jen.ID("req").Dot("Body")),
				),
				jen.Default().Block(
					jen.ID("d").Equals().Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")),
				),
			),
			jen.Line(),
			jen.Return().ID("d").Dot("Decode").Call(jen.ID("v")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideResponseEncoder provides a jsonResponseEncoder"),
		jen.Line(),
		jen.Func().ID("ProvideResponseEncoder").Params().Params(jen.ID("EncoderDecoder")).Block(
			jen.Return().AddressOf().ID("ServerEncoderDecoder").Values(),
		),
		jen.Line(),
	)
	return ret
}
