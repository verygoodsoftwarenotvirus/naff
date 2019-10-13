package frontend

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func httpRoutesTestDotGo() *jen.File {
	ret := jen.NewFile("frontend")
	ret.Add(jen.Null(),

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
	ret.Add(jen.Func().ID("TestService_StaticDir").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.List(jen.ID("cwd"), jen.ID("err")).Op(":=").Qual("os", "Getwd").Call(),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("hf"), jen.ID("err")).Op(":=").ID("s").Dot(
				"StaticDir",
			).Call(jen.ID("cwd")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("hf")),
			jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot(
				"NewRecorder",
			).Call()),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
			).Op("=").Lit("/http_routes_test.go"),
			jen.ID("hf").Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with frontend routing path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("exampleDir").Op(":=").Lit("."),
			jen.List(jen.ID("hf"), jen.ID("err")).Op(":=").ID("s").Dot(
				"StaticDir",
			).Call(jen.ID("exampleDir")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("hf")),
			jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot(
				"NewRecorder",
			).Call()),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
			).Op("=").Lit("/login"),
			jen.ID("hf").Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with frontend items routing path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("logger").Op(":").ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("exampleDir").Op(":=").Lit("."),
			jen.List(jen.ID("hf"), jen.ID("err")).Op(":=").ID("s").Dot(
				"StaticDir",
			).Call(jen.ID("exampleDir")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("hf")),
			jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot(
				"NewRecorder",
			).Call()),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
			).Op("=").Lit("/items/9"),
			jen.ID("hf").Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_Routes").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.Parens(jen.Op("&").ID("Service").Valuesln()).Dot(
				"Routes",
			).Call()),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_buildStaticFileServer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("config").Op(":").ID("config").Dot(
				"FrontendSettings",
			).Valuesln(jen.ID("CacheStaticFiles").Op(":").ID("true"))),
			jen.List(jen.ID("cwd"), jen.ID("err")).Op(":=").Qual("os", "Getwd").Call(),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildStaticFileServer",
			).Call(jen.ID("cwd")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),

		jen.Line(),
	)
	return ret
}
