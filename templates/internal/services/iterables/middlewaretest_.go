package iterables

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("mockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildMockHTTPHandler()...)
	code.Add(buildMockHTTPHandlerServeHTTP()...)
	code.Add(buildTestService_CreationInputMiddleware(proj)...)
	code.Add(buildTestService_UpdateInputMiddleware(proj)...)

	return code
}

func buildMockHTTPHandler() []jen.Code {
	lines := []jen.Code{
		jen.Type().ID("mockHTTPHandler").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildMockHTTPHandlerServeHTTP() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockHTTPHandler")).ID("ServeHTTP").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_CreationInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_CreationInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("actual").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_UpdateInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_UpdateInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(jen.Lit("ServeHTTP"), jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().ID("buildTestService").Call(),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(jen.Lit("DecodeRequest"), jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("actual").Assign().ID("s").Dot("UpdateInputMiddleware").Call(jen.ID("mh")),
				jen.ID("actual").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "mh"),
			),
		),
		jen.Line(),
	}

	return lines
}
