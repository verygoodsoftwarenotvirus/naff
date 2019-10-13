package http

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func middlewareTestDotGo() *jen.File {
	ret := jen.NewFile("httpserver")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("mockHTTPHandler")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("mockHTTPHandler").Struct(jen.ID("mock").Dot(
		"Mock",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Params(jen.ID("m").Op("*").ID("mockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
		jen.ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("buildRequest").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual("net/http", "Request")).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
		jen.ID("require").Dot(
			"NotNil",
		).Call(jen.ID("t"), jen.ID("req")),
		jen.ID("assert").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.Return().ID("req"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("Test_formatSpanNameForRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Dot(
				"Method",
			).Op("=").Qual("net/http", "MethodPatch"),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
			).Op("=").Lit("/blah"),
			jen.ID("expected").Op(":=").Lit("PATCH /blah"),
			jen.ID("actual").Op(":=").ID("formatSpanNameForRequest").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestServer_loggingMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestServer").Call(),
			jen.ID("mh").Op(":=").Op("&").ID("mockHTTPHandler").Valuesln(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("s").Dot(
				"loggingMiddleware",
			).Call(jen.ID("mh")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),

		jen.Line(),
	)
	return ret
}
