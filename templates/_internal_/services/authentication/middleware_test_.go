package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildArbitraryPASETO").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("helper").Op("*").ID("authServiceHTTPRoutesTestHelper"), jen.ID("issueTime").Qual("time", "Time"), jen.ID("lifetime").Qual("time", "Duration"), jen.ID("pasetoData").ID("string")).Params(jen.Op("*").ID("types").Dot("PASETOResponse")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("jsonToken").Op(":=").ID("paseto").Dot("JSONToken").Valuesln(jen.ID("Audience").Op(":").Qual("strconv", "FormatUint").Call(
				jen.ID("helper").Dot("exampleAPIClient").Dot("BelongsToUser"),
				jen.Lit(10),
			), jen.ID("Subject").Op(":").Qual("strconv", "FormatUint").Call(
				jen.ID("helper").Dot("exampleAPIClient").Dot("BelongsToUser"),
				jen.Lit(10),
			), jen.ID("Jti").Op(":").ID("uuid").Dot("NewString").Call(), jen.ID("Issuer").Op(":").ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("Issuer"), jen.ID("IssuedAt").Op(":").ID("issueTime"), jen.ID("NotBefore").Op(":").ID("issueTime"), jen.ID("Expiration").Op(":").ID("issueTime").Dot("Add").Call(jen.ID("lifetime"))),
			jen.ID("jsonToken").Dot("Set").Call(
				jen.ID("pasetoDataKey"),
				jen.ID("pasetoData"),
			),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("paseto").Dot("NewV2").Call().Dot("Encrypt").Call(
				jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey"),
				jen.ID("jsonToken"),
				jen.Lit(""),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().Op("&").ID("types").Dot("PASETOResponse").Valuesln(jen.ID("Token").Op(":").ID("token"), jen.ID("ExpiresAt").Op(":").ID("jsonToken").Dot("Expiration").Dot("String").Call()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_fetchSessionContextDataFromPASETO").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("tokenRes"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("buildPASETOResponse").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("sessionCtxData"),
						jen.ID("helper").Dot("exampleAPIClient"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("pasetoAuthorizationHeaderKey"),
						jen.ID("tokenRes").Dot("Token"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("fetchSessionContextDataFromPASETO").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
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
				jen.Lit("with invalid PASETO"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("pasetoAuthorizationHeaderKey"),
						jen.Lit("blah"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("fetchSessionContextDataFromPASETO").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with expired token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("tokenRes").Op(":=").ID("buildArbitraryPASETO").Call(
						jen.ID("t"),
						jen.ID("helper"),
						jen.Qual("time", "Now").Call().Dot("Add").Call(jen.Op("-").Lit(24).Op("*").Qual("time", "Hour")),
						jen.Qual("time", "Minute"),
						jen.Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("helper").Dot("sessionCtxData").Dot("ToBytes").Call()),
					),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("pasetoAuthorizationHeaderKey"),
						jen.ID("tokenRes").Dot("Token"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("fetchSessionContextDataFromPASETO").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid base64 encoding"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("tokenRes").Op(":=").ID("buildArbitraryPASETO").Call(
						jen.ID("t"),
						jen.ID("helper"),
						jen.Qual("time", "Now").Call(),
						jen.Qual("time", "Hour"),
						jen.Lit(`       \\\\\\\\\\\\               lololo`),
					),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("pasetoAuthorizationHeaderKey"),
						jen.ID("tokenRes").Dot("Token"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("fetchSessionContextDataFromPASETO").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid GOB string"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("tokenRes").Op(":=").ID("buildArbitraryPASETO").Call(
						jen.ID("t"),
						jen.ID("helper"),
						jen.Qual("time", "Now").Call(),
						jen.Qual("time", "Hour"),
						jen.Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.Index().ID("byte").Call(jen.Lit("blah"))),
					),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("pasetoAuthorizationHeaderKey"),
						jen.ID("tokenRes").Dot("Token"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("fetchSessionContextDataFromPASETO").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("fetchSessionContextDataFromPASETO").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_CookieAuthenticationMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("accountUserMembershipDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Values(),
					jen.ID("accountUserMembershipDataManager").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("sessionCtxData"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountUserMembershipDataManager"),
					jen.ID("mockHandler").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("mockHandler").Dot("On").Call(
						jen.Lit("ServeHTTP"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
					).Dot("Return").Call(),
					jen.List(jen.ID("_"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("helper").Dot("service").Dot("CookieRequirementMiddleware").Call(jen.ID("mockHandler")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockHandler"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_UserAttributionMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("mockAccountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Values(),
					jen.ID("mockAccountMembershipManager").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("sessionCtxData"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("mockAccountMembershipManager"),
					jen.List(jen.ID("_"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("h").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("h").Dot("On").Call(
						jen.Lit("ServeHTTP"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("UserAttributionMiddleware").Call(jen.ID("h")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAccountMembershipManager"),
						jen.ID("h"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building session context data for user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("mockAccountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Values(),
					jen.ID("mockAccountMembershipManager").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("SessionContextData")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("mockAccountMembershipManager"),
					jen.List(jen.ID("_"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("mh").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("helper").Dot("service").Dot("UserAttributionMiddleware").Call(jen.ID("mh")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAccountMembershipManager"),
						jen.ID("mh"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with PASETO"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("tokenRes"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("buildPASETOResponse").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("sessionCtxData"),
						jen.ID("helper").Dot("exampleAPIClient"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("pasetoAuthorizationHeaderKey"),
						jen.ID("tokenRes").Dot("Token"),
					),
					jen.ID("h").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("h").Dot("On").Call(
						jen.Lit("ServeHTTP"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("UserAttributionMiddleware").Call(jen.ID("h")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("h"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with PASETO and issue parsing token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("pasetoAuthorizationHeaderKey"),
						jen.Lit("blah"),
					),
					jen.ID("h").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("h").Dot("On").Call(
						jen.Lit("ServeHTTP"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("UserAttributionMiddleware").Call(jen.ID("h")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_AuthorizationMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("mockUserDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Values(),
					jen.ID("mockUserDataManager").Dot("On").Call(
						jen.Lit("GetSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("sessionCtxData"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockUserDataManager"),
					jen.ID("h").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("h").Dot("On").Call(
						jen.Lit("ServeHTTP"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("types").Dot("SessionContextDataKey"),
						jen.ID("sessionCtxData"),
					)),
					jen.ID("helper").Dot("service").Dot("AuthorizationMiddleware").Call(jen.ID("h")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("h"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with banned user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus").Op("=").ID("types").Dot("BannedUserAccountStatus"),
					jen.ID("helper").Dot("setContextFetcher").Call(jen.ID("t")),
					jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("mockUserDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Values(),
					jen.ID("mockUserDataManager").Dot("On").Call(
						jen.Lit("GetSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("sessionCtxData"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockUserDataManager"),
					jen.ID("h").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("h").Dot("On").Call(
						jen.Lit("ServeHTTP"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("types").Dot("SessionContextDataKey"),
						jen.ID("sessionCtxData"),
					)),
					jen.ID("mh").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("helper").Dot("service").Dot("AuthorizationMiddleware").Call(jen.ID("mh")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusForbidden"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mh"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("mh").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("helper").Dot("service").Dot("AuthorizationMiddleware").Call(jen.ID("mh")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mh"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without authorization for account"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("sessionCtxData").Dot("AccountPermissions").Op("=").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Values(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("sessionCtxData"), jen.ID("nil"))),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("types").Dot("SessionContextDataKey"),
						jen.ID("sessionCtxData"),
					)),
					jen.ID("helper").Dot("service").Dot("AuthorizationMiddleware").Call(jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values()).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_AdminMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("=").Index().ID("string").Valuesln(jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.ID("helper").Dot("setContextFetcher").Call(jen.ID("t")),
					jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("SessionContextDataKey"),
						jen.ID("sessionCtxData"),
					)),
					jen.ID("mockHandler").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("mockHandler").Dot("On").Call(
						jen.Lit("ServeHTTP"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("ServiceAdminMiddleware").Call(jen.ID("mockHandler")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockHandler"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("=").Index().ID("string").Valuesln(jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("SessionContextDataKey"),
						jen.ID("sessionCtxData"),
					)),
					jen.ID("mockHandler").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("helper").Dot("service").Dot("ServiceAdminMiddleware").Call(jen.ID("mockHandler")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockHandler"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with non-admin user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("SessionContextDataKey"),
						jen.ID("sessionCtxData"),
					)),
					jen.ID("mockHandler").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPHandler").Values(),
					jen.ID("helper").Dot("service").Dot("ServiceAdminMiddleware").Call(jen.ID("mockHandler")).Dot("ServeHTTP").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockHandler"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
