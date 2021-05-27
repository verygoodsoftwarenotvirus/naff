package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestService_fetchAPIClient").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("Anything"),
						jen.ID("apiClientIDURLParamKey"),
						jen.Lit("API client"),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAPIClient").Dot("ID"))),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.ID("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAPIClient"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fake mode"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("useFakeData").Op("=").ID("true"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching apiClient"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("Anything"),
						jen.ID("apiClientIDURLParamKey"),
						jen.Lit("API client"),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAPIClient").Dot("ID"))),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClient")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_buildAPIClientEditorView").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("Anything"),
						jen.ID("apiClientIDURLParamKey"),
						jen.Lit("API client"),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAPIClient").Dot("ID"))),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.ID("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientEditorView").Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("Anything"),
						jen.ID("apiClientIDURLParamKey"),
						jen.Lit("API client"),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAPIClient").Dot("ID"))),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.ID("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientEditorView").Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientEditorView").Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching apiClient"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.ID("mock").Dot("Anything"),
						jen.ID("apiClientIDURLParamKey"),
						jen.Lit("API client"),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAPIClient").Dot("ID"))),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAPIClient").Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClient")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientEditorView").Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_fetchAPIClients").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("exampleAPIClientList").Op(":=").ID("fakes").Dot("BuildFakeAPIClientList").Call(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleAPIClientList"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClients").Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAPIClientList"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fake mode"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("useFakeData").Op("=").ID("true"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClients").Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClientList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchAPIClients").Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_buildAPIClientsTableView").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleAPIClientList").Op(":=").ID("fakes").Dot("BuildFakeAPIClientList").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleAPIClientList"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientsTableView").Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleAPIClientList").Op(":=").ID("fakes").Dot("BuildFakeAPIClientList").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleAPIClientList"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientsTableView").Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientsTableView").Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClientList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/api_clients"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildAPIClientsTableView").Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
