package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperCookieTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_newCookieRoundTripper").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("newCookieRoundTripper").Call(
							jen.ID("c"),
							jen.Op("&").Qual("net/http", "Cookie").Valuesln(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("mockRoundTripper").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockRoundTripper")).ID("RoundTrip").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual("net/http", "Response"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(jen.ID("req")),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual("net/http", "Response")), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_cookieRoundtripper_RoundTrip").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("exampleCookie").Op(":=").Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").Lit("testcookie"), jen.ID("Value").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("rt").Op(":=").ID("newCookieRoundTripper").Call(
						jen.ID("c"),
						jen.ID("exampleCookie"),
					),
					jen.ID("exampleResponse").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Header").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.Lit("Set-Cookie").Op(":").Valuesln(jen.ID("exampleCookie").Dot("String").Call())), jen.ID("StatusCode").Op(":").Qual("net/http", "StatusTeapot")),
					jen.ID("mrt").Op(":=").Op("&").ID("mockRoundTripper").Valuesln(),
					jen.ID("mrt").Dot("On").Call(
						jen.Lit("RoundTrip"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").Qual("net/http", "Request").Valuesln()),
					).Dot("Return").Call(
						jen.ID("exampleResponse"),
						jen.ID("nil"),
					),
					jen.ID("rt").Dot("base").Op("=").ID("mrt"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.ID("c").Dot("URL").Call().Dot("String").Call(),
						jen.ID("nil"),
					),
					jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("rt").Dot("RoundTrip").Call(jen.ID("req")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("res"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleResponse"),
						jen.ID("res"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mrt"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing RoundTrip"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("exampleCookie").Op(":=").Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").Lit("testcookie"), jen.ID("Value").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("rt").Op(":=").ID("newCookieRoundTripper").Call(
						jen.ID("c"),
						jen.ID("exampleCookie"),
					),
					jen.ID("mrt").Op(":=").Op("&").ID("mockRoundTripper").Valuesln(),
					jen.ID("mrt").Dot("On").Call(
						jen.Lit("RoundTrip"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").Qual("net/http", "Request").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual("net/http", "Response")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("rt").Dot("base").Op("=").ID("mrt"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.ID("c").Dot("URL").Call().Dot("String").Call(),
						jen.ID("nil"),
					),
					jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("rt").Dot("RoundTrip").Call(jen.ID("req")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
