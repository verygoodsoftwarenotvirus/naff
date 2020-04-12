package iterables

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("mockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockHTTPHandler").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_CreationInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(utils.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("actual").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
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
		jen.Func().ID("TestService_UpdateInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(utils.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("req"), nil),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("actual").Assign().ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
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
