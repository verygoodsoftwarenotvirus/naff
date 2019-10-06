package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareTestDotGo() *jen.File {
	ret := jen.NewFile("users")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("MockHTTPHandler")).Call(jen.ID("nil")))
	ret.Add(jen.Null().Type().ID("MockHTTPHandler").Struct(
		jen.ID("mock").Dot(
			"Mock",
		),
	),
	)
	ret.Add(jen.Func().Params(jen.ID("m").Op("*").ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
		jen.ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
	),
	)
	ret.Add(jen.Func().ID("TestService_UserInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"UserInputMiddleware",
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
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"UserInputMiddleware",
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
	)
	ret.Add(jen.Func().ID("TestService_PasswordUpdateInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"PasswordUpdateInputMiddleware",
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
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("uint64").Call(jen.Lit(123)), jen.ID("nil")),
			jen.ID("s").Dot(
				"database",
			).Op("=").ID("mockDB"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"PasswordUpdateInputMiddleware",
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
	)
	ret.Add(jen.Func().ID("TestService_TOTPSecretRefreshInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"TOTPSecretRefreshInputMiddleware",
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
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
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
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"TOTPSecretRefreshInputMiddleware",
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
	)
	return ret
}
