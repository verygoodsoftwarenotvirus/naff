package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareTestDotGo() *jen.File {
	ret := jen.NewFile("webhooks")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("MockHTTPHandler")).Call(jen.ID("nil")),
	jen.Line(),
	)

	ret.Add(
		jen.Type().ID("MockHTTPHandler").Struct(jen.ID("mock").Dot(
		"Mock",
	)),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
		jen.ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_CreationInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"CreationInputMiddleware",
			).Call(jen.ID("mh")),
			jen.ID("actual").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
		jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"CreationInputMiddleware",
			).Call(jen.ID("mh")),
			jen.ID("actual").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusBadRequest")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UpdateInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"UpdateInputMiddleware",
			).Call(jen.ID("mh")),
			jen.ID("actual").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusOK")),
		)),
		jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("mh").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"UpdateInputMiddleware",
			).Call(jen.ID("mh")),
			jen.ID("actual").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusBadRequest")),
		)),
	),
	jen.Line(),
	)
	return ret
}
