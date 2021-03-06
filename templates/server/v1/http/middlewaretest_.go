package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("mockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildMiddlewareTestTypeDefinitions()...)
	code.Add(buildMockHTTPHandlerServeHTTP()...)
	code.Add(buildBuildRequest()...)
	code.Add(buildTest_formatSpanNameForRequest()...)
	code.Add(buildTestServer_loggingMiddleware()...)

	return code
}

func buildMiddlewareTestTypeDefinitions() []jen.Code {
	lines := []jen.Code{
		jen.Type().ID("mockHTTPHandler").Struct(
			jen.Qual(constants.MockPkg, "Mock"),
		),
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

func buildBuildRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildRequest").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().Qual("net/http", "Request")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("https://verygoodsoftwarenotvirus.ru"),
				jen.Nil(),
			),
			jen.Line(),
			utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
			utils.AssertNoError(jen.Err(), nil),
			jen.Line(),
			jen.Return().ID(constants.RequestVarName),
		),
		jen.Line(),
	}

	return lines
}

func buildTest_formatSpanNameForRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_formatSpanNameForRequest").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Dot("Method").Equals().Qual("net/http", "MethodPatch"),
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Lit("/blah"),
				jen.Line(),
				jen.ID("expected").Assign().Lit("PATCH /blah"),
				jen.ID("actual").Assign().ID("formatSpanNameForRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestServer_loggingMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestServer_loggingMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestServer").Call(),
				jen.Line(),
				jen.ID("mh").Assign().AddressOf().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)).Assign().List(jen.ID("httptest").Dot("NewRecorder").Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
				jen.ID("s").Dot("loggingMiddleware").Call(jen.ID("mh")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				utils.AssertExpectationsFor("mh"),
			),
		),
		jen.Line(),
	}

	return lines
}
