package chi

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func chiTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildRouterForTest").Params().Params(jen.ID("routing").Dot("Router")).Body(
			jen.Return().ID("NewRouter").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestNewRouter").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("NewRouter").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildChiMux").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("buildChiMux").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_convertMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("convertMiddleware").Call(jen.Func().Params(jen.Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
							jen.Return().ID("nil"))),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_AddRoute").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("methods").Op(":=").Index().ID("string").Valuesln(
						jen.Qual("net/http", "MethodGet"), jen.Qual("net/http", "MethodHead"), jen.Qual("net/http", "MethodPost"), jen.Qual("net/http", "MethodPut"), jen.Qual("net/http", "MethodPatch"), jen.Qual("net/http", "MethodDelete"), jen.Qual("net/http", "MethodConnect"), jen.Qual("net/http", "MethodOptions"), jen.Qual("net/http", "MethodTrace")),
					jen.For(jen.List(jen.ID("_"), jen.ID("method")).Op(":=").Range().ID("methods")).Body(
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("r").Dot("AddRoute").Call(
								jen.ID("method"),
								jen.Lit("/path"),
								jen.ID("nil"),
							),
						)),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("r").Dot("AddRoute").Call(
							jen.Lit("blah"),
							jen.Lit("/path"),
							jen.ID("nil"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Connect").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Connect").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Delete").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Delete").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Get").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Get").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Handle").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Handle").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_HandleFunc").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("HandleFunc").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Handler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("r").Dot("Handler").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Head").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Head").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_LogRoutes").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("r").Dot("AddRoute").Call(
							jen.Qual("net/http", "MethodGet"),
							jen.Lit("/path"),
							jen.ID("nil"),
						),
					),
					jen.ID("r").Dot("LogRoutes").Call(),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Options").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Options").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Patch").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Patch").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Post").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Post").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Put").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Put").Call(
						jen.Lit("/thing"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Route").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("r").Dot("Route").Call(
							jen.Lit("/test"),
							jen.Func().Params(jen.ID("routing").Dot("Router")).Body(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_Trace").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("r").Dot("Trace").Call(
						jen.Lit("/test"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_WithMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouterForTest").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("r").Dot("WithMiddleware").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_clone").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouter").Call(
						jen.ID("nil"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("r").Dot("clone").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_BuildRouteParamIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouter").Call(
						jen.ID("nil"),
						jen.ID("nil"),
					),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleKey").Op(":=").Lit("blah"),
					jen.ID("rf").Op(":=").ID("r").Dot("BuildRouteParamIDFetcher").Call(
						jen.ID("l"),
						jen.ID("exampleKey"),
						jen.Lit("desc"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("rf"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/blah"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUser").Call().Dot("ID"),
					jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("req").Dot("Context").Call(),
						jen.ID("chi").Dot("RouteCtxKey"),
						jen.Op("&").ID("chi").Dot("Context").Valuesln(
							jen.ID("URLParams").Op(":").ID("chi").Dot("RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Valuesln(
									jen.ID("exampleKey")), jen.ID("Values").Op(":").Index().ID("string").Valuesln(
									jen.Qual("strconv", "FormatUint").Call(
										jen.ID("expected"),
										jen.Lit(10),
									)))),
					)),
					jen.ID("actual").Op(":=").ID("rf").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without appropriate value attached to context"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouter").Call(
						jen.ID("nil"),
						jen.ID("nil"),
					),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleKey").Op(":=").Lit("blah"),
					jen.ID("rf").Op(":=").ID("r").Dot("BuildRouteParamIDFetcher").Call(
						jen.ID("l"),
						jen.ID("exampleKey"),
						jen.Lit("desc"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("rf"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/blah"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("actual").Op(":=").ID("rf").Call(jen.ID("req")),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_router_BuildRouteParamStringIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("buildRouter").Call(
						jen.ID("nil"),
						jen.ID("nil"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleKey").Op(":=").Lit("blah"),
					jen.ID("rf").Op(":=").ID("r").Dot("BuildRouteParamStringIDFetcher").Call(jen.ID("exampleKey")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("rf"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/blah"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("exampleID").Op(":=").ID("fakes").Dot("BuildFakeUser").Call().Dot("ID"),
					jen.ID("expected").Op(":=").Qual("strconv", "FormatUint").Call(
						jen.ID("exampleID"),
						jen.Lit(10),
					),
					jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("req").Dot("Context").Call(),
						jen.ID("chi").Dot("RouteCtxKey"),
						jen.Op("&").ID("chi").Dot("Context").Valuesln(
							jen.ID("URLParams").Op(":").ID("chi").Dot("RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Valuesln(
									jen.ID("exampleKey")), jen.ID("Values").Op(":").Index().ID("string").Valuesln(
									jen.ID("expected")))),
					)),
					jen.ID("actual").Op(":=").ID("rf").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
