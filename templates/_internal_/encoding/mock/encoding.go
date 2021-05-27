package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func encodingDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("encoding").Dot("ServerEncoderDecoder").Op("=").Parens(jen.Op("*").ID("EncoderDecoder")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewMockEncoderDecoder produces a mock EncoderDecoder."),
		jen.Line(),
		jen.Func().ID("NewMockEncoderDecoder").Params().Params(jen.Op("*").ID("EncoderDecoder")).Body(
			jen.Return().Op("&").ID("EncoderDecoder").Values()),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("EncoderDecoder").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MustEncode satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("MustEncode").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("v").Interface()).Params(jen.Index().ID("byte")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("v"),
			).Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("byte"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MustEncodeJSON satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("MustEncodeJSON").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("v").Interface()).Params(jen.Index().ID("byte")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("v"),
			).Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("byte"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RespondWithData satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("RespondWithData").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("val").Interface()).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("val"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeResponseWithStatus satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeResponseWithStatus").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("val").Interface(), jen.ID("statusCode").ID("int")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("val"),
				jen.ID("statusCode"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.ID("statusCode")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeErrorResponse satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeErrorResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("msg").ID("string"), jen.ID("statusCode").ID("int")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("msg"),
				jen.ID("statusCode"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.ID("statusCode")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeInvalidInputResponse satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeInvalidInputResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeNotFoundResponse satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeNotFoundResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeUnspecifiedInternalServerErrorResponse satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeUnspecifiedInternalServerErrorResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusInternalServerError")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeUnauthorizedResponse satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeUnauthorizedResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeInvalidPermissionsResponse satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("EncodeInvalidPermissionsResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusForbidden")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("DecodeRequest satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("DecodeRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("v"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("DecodeBytes satisfies our EncoderDecoder interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("EncoderDecoder")).ID("DecodeBytes").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Index().ID("byte"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("data"),
				jen.ID("v"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
