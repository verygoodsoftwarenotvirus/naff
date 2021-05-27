package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperBaseTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_newDefaultRoundTripper").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("rt").Op(":=").ID("newDefaultRoundTripper").Call(jen.Lit(0)),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("rt"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_defaultRoundTripper_RoundTrip").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("ts").Op(":=").ID("httptest").Dot("NewServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
						jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK"))))),
					jen.ID("transport").Op(":=").ID("newDefaultRoundTripper").Call(jen.Lit(0)),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("transport"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.ID("ts").Dot("URL"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("transport").Dot("RoundTrip").Call(jen.ID("req")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildWrappedTransport").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("buildWrappedTransport").Call(jen.Qual("time", "Minute")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildRequestLogHook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("actual").Op(":=").ID("buildRequestLogHook").Call(jen.ID("l")),
					jen.ID("actual").Call(
						jen.ID("nil"),
						jen.Op("&").Qual("net/http", "Request").Values(),
						jen.Lit(0),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildResponseLogHook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("actual").Op(":=").ID("buildResponseLogHook").Call(jen.ID("l")),
					jen.ID("actual").Call(
						jen.ID("nil"),
						jen.Op("&").Qual("net/http", "Response").Values(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildCheckRetryFunc").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("f").Op(":=").ID("buildCheckRetryFunc").Call(jen.ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call())),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("f").Call(
						jen.ID("ctx"),
						jen.Op("&").Qual("net/http", "Response").Values(),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildErrorHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("actual").Op(":=").ID("buildErrorHandler").Call(jen.ID("l")),
					jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("actual").Call(
						jen.Op("&").Qual("net/http", "Response").Values(),
						jen.ID("nil"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildRetryingClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("actual").Op(":=").ID("buildRetryingClient").Call(
						jen.Qual("net/http", "DefaultClient"),
						jen.ID("nil"),
						jen.ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call()),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
