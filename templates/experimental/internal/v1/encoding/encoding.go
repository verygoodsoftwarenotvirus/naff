package encoding

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func encodingDotGo() *jen.File {
	ret := jen.NewFile("encoding")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("ContentTypeHeader").Op("=").Lit("Content-type").Var().ID("XMLContentType").Op("=").Lit("application/xml").Var().ID("JSONContentType").Op("=").Lit("application/json").Var().ID("DefaultContentType").Op("=").ID("JSONContentType"),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideResponseEncoder")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("EncoderDecoder").Interface(jen.ID("EncodeResponse").Params(jen.Qual("net/http", "ResponseWriter"), jen.Interface()).Params(jen.ID("error")), jen.ID("DecodeRequest").Params(jen.Op("*").Qual("net/http", "Request"), jen.Interface()).Params(jen.ID("error"))).Type().ID("ServerEncoderDecoder").Struct().Type().ID("encoder").Interface(jen.ID("Encode").Params(jen.ID("v").Interface()).Params(jen.ID("error"))).Type().ID("decoder").Interface(jen.ID("Decode").Params(jen.ID("v").Interface()).Params(jen.ID("error"))),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// EncodeResponse encodes responses").Params(jen.ID("ed").Op("*").ID("ServerEncoderDecoder")).ID("EncodeResponse").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Params(jen.ID("error")).Block(
		jen.Null().Var().ID("ct").Op("=").Qual("strings", "ToLower").Call(jen.ID("res").Dot(
			"Header",
		).Call().Dot(
			"Get",
		).Call(jen.ID("ContentTypeHeader"))),
		jen.If(jen.ID("ct").Op("==").Lit("")).Block(
			jen.ID("ct").Op("=").ID("DefaultContentType"),
		),
		jen.Null().Var().ID("e").ID("encoder"),
		jen.Switch(jen.ID("ct")).Block(
			jen.Case(jen.ID("XMLContentType")).Block(jen.ID("e").Op("=").Qual("encoding/xml", "NewEncoder").Call(jen.ID("res"))),
			jen.Default().Block(jen.ID("e").Op("=").Qual("encoding/json", "NewEncoder").Call(jen.ID("res"))),
		),
		jen.ID("res").Dot(
			"Header",
		).Call().Dot(
			"Set",
		).Call(jen.ID("ContentTypeHeader"), jen.ID("ct")),
		jen.Return().ID("e").Dot(
			"Encode",
		).Call(jen.ID("v")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// DecodeRequest decodes responses").Params(jen.ID("ed").Op("*").ID("ServerEncoderDecoder")).ID("DecodeRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.ID("error")).Block(
		jen.Null().Var().ID("ct").Op("=").Qual("strings", "ToLower").Call(jen.ID("req").Dot(
			"Header",
		).Dot(
			"Get",
		).Call(jen.ID("ContentTypeHeader"))),
		jen.If(jen.ID("ct").Op("==").Lit("")).Block(
			jen.ID("ct").Op("=").ID("DefaultContentType"),
		),
		jen.Null().Var().ID("d").ID("decoder"),
		jen.Switch(jen.ID("ct")).Block(
			jen.Case(jen.ID("XMLContentType")).Block(jen.ID("d").Op("=").Qual("encoding/xml", "NewDecoder").Call(jen.ID("req").Dot(
				"Body",
			))),
			jen.Default().Block(jen.ID("d").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot(
				"Body",
			))),
		),
		jen.Return().ID("d").Dot(
			"Decode",
		).Call(jen.ID("v")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideResponseEncoder provides a jsonResponseEncoder").ID("ProvideResponseEncoder").Params().Params(jen.ID("EncoderDecoder")).Block(
		jen.Return().Op("&").ID("ServerEncoderDecoder").Valuesln(),
	),

		jen.Line(),
	)
	return ret
}
