package logging

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func loggingTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestEnsureLogger").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("EnsureLogger").Call(jen.ID("NewNoopLogger").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("EnsureLogger").Call(jen.ID("nil")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildLoggingMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("middleware").Op(":=").ID("BuildLoggingMiddleware").Call(jen.ID("NewNoopLogger").Call()),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("middleware"),
					),
					jen.ID("hf").Op(":=").Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body()),
					jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/nil"),
						jen.ID("nil"),
					), jen.ID("httptest").Dot("NewRecorder").Call()),
					jen.ID("middleware").Call(jen.ID("hf")).Dot("ServeHTTP").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with non-logged route"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("middleware").Op(":=").ID("BuildLoggingMiddleware").Call(jen.ID("NewNoopLogger").Call()),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("middleware"),
					),
					jen.ID("hf").Op(":=").Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body()),
					jen.If(jen.ID("len").Call(jen.ID("doNotLog")).Op("==").Lit(0)).Body(
						jen.ID("t").Dot("SkipNow").Call()),
					jen.Var().Defs(
						jen.ID("route").ID("string"),
					),
					jen.For(jen.ID("k").Op(":=").Range().ID("doNotLog")).Body(
						jen.ID("route").Op("=").ID("k"),
						jen.Break(),
					),
					jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.ID("route"),
						jen.ID("nil"),
					), jen.ID("httptest").Dot("NewRecorder").Call()),
					jen.ID("middleware").Call(jen.ID("hf")).Dot("ServeHTTP").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
