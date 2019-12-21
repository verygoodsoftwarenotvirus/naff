package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("encoding")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("ContentTypeHeader is the HTTP standard header name for content type"),
			jen.ID("ContentTypeHeader").Op("=").Lit("Content-type"),
			jen.Comment("XMLContentType represents the XML content type"),
			jen.ID("XMLContentType").Op("=").Lit("application/xml"),
			jen.Comment("JSONContentType represents the JSON content type"),
			jen.ID("JSONContentType").Op("=").Lit("application/json"),
			jen.Comment("DefaultContentType is what the library defaults to"),
			jen.ID("DefaultContentType").Op("=").ID("JSONContentType"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers provides ResponseEncoders for dependency injection"),
			jen.ID("Providers").Op("=").Qual("github.com/google/wire", "NewSet".Callln(
				jen.ID("ProvideResponseEncoder"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("EncoderDecoder is an interface that allows for multiple implementations of HTTP response formats"),
			jen.ID("EncoderDecoder").Interface(
				jen.ID("EncodeResponse").Params(jen.Qual("net/http", "ResponseWriter"), jen.Interface()).Params(jen.ID("error")),
				jen.ID("DecodeRequest").Params(jen.Op("*").Qual("net/http", "Request"), jen.Interface()).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("ServerEncoderDecoder is our concrete implementation of EncoderDecoder"),
			jen.ID("ServerEncoderDecoder").Struct(),
			jen.Line(),
			jen.ID("encoder").Interface(
				jen.ID("Encode").Params(jen.ID("v").Interface()).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.ID("decoder").Interface(
				jen.ID("Decode").Params(jen.ID("v").Interface()).Params(jen.ID("error")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("EncodeResponse encodes responses"),
		jen.Line(),
		jen.Func().Params(jen.ID("ed").Op("*").ID("ServerEncoderDecoder")).ID("EncodeResponse").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Params(jen.ID("error")).Block(
			jen.Var().ID("ct").Op("=").Qual("strings", "ToLower").Call(jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("ContentTypeHeader"))),
			jen.If(jen.ID("ct").Op("==").Lit("")).Block(
				jen.ID("ct").Op("=").ID("DefaultContentType"),
			),
			jen.Line(),
			jen.Var().ID("e").ID("encoder"),
			jen.Switch(jen.ID("ct")).Block(
				jen.Case(jen.ID("XMLContentType")).Block(jen.ID("e").Op("=").Qual("encoding/xml", "NewEncoder").Call(jen.ID("res"))),
				jen.Default().Block(jen.ID("e").Op("=").Qual("encoding/json", "NewEncoder").Call(jen.ID("res"))),
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
		jen.Func().Params(jen.ID("ed").Op("*").ID("ServerEncoderDecoder")).ID("DecodeRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.ID("error")).Block(
			jen.Var().ID("ct").Op("=").Qual("strings", "ToLower").Call(jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("ContentTypeHeader"))),
			jen.If(jen.ID("ct").Op("==").Lit("")).Block(
				jen.ID("ct").Op("=").ID("DefaultContentType"),
			),
			jen.Line(),
			jen.Var().ID("d").ID("decoder"),
			jen.Switch(jen.ID("ct")).Block(
				jen.Case(jen.ID("XMLContentType")).Block(
					jen.ID("d").Op("=").Qual("encoding/xml", "NewDecoder").Call(jen.ID("req").Dot("Body")),
				),
				jen.Default().Block(
					jen.ID("d").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")),
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
			jen.Return().Op("&").ID("ServerEncoderDecoder").Values(),
		),
		jen.Line(),
	)
	return ret
}
