package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestService_fetchUsers").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUserList").Op(":=").ID("fakes").Dot("BuildFakeUserList").Call(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUsers"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleUserList"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/users"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchUsers").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleUserList"),
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
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/users"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchUsers").Call(
						jen.ID("ctx"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUsers"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("UserList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/users"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("fetchUsers").Call(
						jen.ID("ctx"),
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
		jen.Func().ID("TestService_buildUsersTableView").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template but not for search"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleUserList").Op(":=").ID("fakes").Dot("BuildFakeUserList").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUsers"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleUserList"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/users"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildUsersTableView").Call(
						jen.ID("true"),
						jen.ID("false"),
					).Call(
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
				jen.Lit("with base template and for search"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleQuery").Op(":=").Lit("whatever"),
					jen.ID("exampleUserList").Op(":=").ID("fakes").Dot("BuildFakeUserList").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("SearchForUsersByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleQuery"),
					).Dot("Return").Call(
						jen.ID("exampleUserList").Dot("Users"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("uri").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("/users?%s=%s"),
						jen.ID("types").Dot("SearchQueryKey"),
						jen.ID("exampleQuery"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.ID("uri"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildUsersTableView").Call(
						jen.ID("true"),
						jen.ID("true"),
					).Call(
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
				jen.Lit("with base template and for search and error performing search"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleQuery").Op(":=").Lit("whatever"),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("SearchForUsersByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleQuery"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("User").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("uri").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("/users?%s=%s"),
						jen.ID("types").Dot("SearchQueryKey"),
						jen.ID("exampleQuery"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.ID("uri"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildUsersTableView").Call(
						jen.ID("true"),
						jen.ID("true"),
					).Call(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template but for search"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleUserList").Op(":=").ID("fakes").Dot("BuildFakeUserList").Call(),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUsers"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleUserList"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/users"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildUsersTableView").Call(
						jen.ID("false"),
						jen.ID("false"),
					).Call(
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
						jen.Lit("/users"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildUsersTableView").Call(
						jen.ID("true"),
						jen.ID("false"),
					).Call(
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUsers"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("UserList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/users"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildUsersTableView").Call(
						jen.ID("true"),
						jen.ID("false"),
					).Call(
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
