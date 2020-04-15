package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("MockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("MockHTTPHandler").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UserInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().AddressOf().ID("Service").Valuesln(
					jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("UserInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().AddressOf().ID("Service").Valuesln(
					jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("actual").Assign().ID("s").Dot("UserInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_PasswordUpdateInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().AddressOf().ID("Service").Valuesln(
					jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("PasswordUpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().AddressOf().ID("Service").Valuesln(
					jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserCount"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Uint64().Call(jen.Lit(123)), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("actual").Assign().ID("s").Dot("PasswordUpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_TOTPSecretRefreshInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().AddressOf().ID("Service").Valuesln(
					jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("TOTPSecretRefreshInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().AddressOf().ID("Service").Valuesln(
					jen.ID("logger").MapAssign().Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("actual").Assign().ID("s").Dot("TOTPSecretRefreshInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
		),
		jen.Line(),
	)
	return ret
}
