package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("mockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockHTTPHandler").Struct(
			jen.Qual("github.com/stretchr/testify/mock", "Mock"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildRequest").Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Params(jen.ParamPointer().Qual("net/http", "Request")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("https://verygoodsoftwarenotvirus.ru"),
				jen.Nil(),
			),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
			utils.AssertNoError(jen.Err(), nil),
			jen.Line(),
			jen.Return().ID("req"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_formatSpanNameForRequest").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Dot("Method").Equals().Qual("net/http", "MethodPatch"),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Lit("/blah"),
				jen.Line(),
				jen.ID("expected").Assign().Lit("PATCH /blah"),
				jen.ID("actual").Assign().ID("formatSpanNameForRequest").Call(jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestServer_loggingMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestServer").Call(),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("res"), jen.ID("req")).Assign().List(jen.ID("httptest").Dot("NewRecorder").Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
				jen.ID("s").Dot("loggingMiddleware").Call(jen.ID("mh")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
