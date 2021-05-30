package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAuthenticationService_getUserIDFromCookie").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.List(jen.ID("_"), jen.ID("userID"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("getUserIDFromCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("userID"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("c").Op(":=").Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"), jen.ID("Value").Op(":").Lit("blah blah blah this is not a real cookie"), jen.ID("Path").Op(":").Lit("/"), jen.ID("HttpOnly").Op(":").ID("true")),
					jen.ID("helper").Dot("req").Dot("AddCookie").Call(jen.ID("c")),
					jen.List(jen.ID("_"), jen.ID("userID"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("getUserIDFromCookie").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("userID"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("_"), jen.ID("userID"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("getUserIDFromCookie").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Qual("net/http", "ErrNoCookie"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("userID"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error loading from session"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("expectedToken").Op(":=").Lit("blahblah"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("expectedToken"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("buildCookie").Call(
						jen.ID("expectedToken"),
						jen.Qual("time", "Now").Call().Dot("Add").Call(jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Lifetime")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("helper").Dot("req").Dot("AddCookie").Call(jen.ID("c")),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("err")).Op("=").ID("helper").Dot("service").Dot("getUserIDFromCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no user ID attached to context"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("ctx"), jen.ID("sessionErr")).Op(":=").ID("helper").Dot("service").Dot("sessionManager").Dot("Load").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("sessionErr"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("sessionManager").Dot("RenewToken").Call(jen.ID("ctx")),
					),
					jen.List(jen.ID("token"), jen.ID("_"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("sessionManager").Dot("Commit").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("token"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("buildCookie").Call(
						jen.ID("token"),
						jen.Qual("time", "Now").Call().Dot("Add").Call(jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Lifetime")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("helper").Dot("req").Dot("AddCookie").Call(jen.ID("c")),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("err")).Op("=").ID("helper").Dot("service").Dot("getUserIDFromCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
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
		jen.Func().ID("TestAuthenticationService_determineUserFromRequestCookie").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
					jen.List(jen.ID("actualUser"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("determineUserFromRequestCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("actualUser"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("actualUser"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("determineUserFromRequestCookie").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actualUser"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving user from datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("expectedError").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.ID("expectedError"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
					jen.List(jen.ID("actualUser"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("determineUserFromRequestCookie").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actualUser"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_validateLogin").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("authenticator").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("authenticator").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("authenticator"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("validateLogin").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("authenticator"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid two factor code"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("authenticator").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("authenticator").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("authentication").Dot("ErrInvalidTOTPToken"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("authenticator"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("validateLogin").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("authenticator"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error returned from validator"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("expectedErr").Op(":=").Qual("errors", "New").Call(jen.Lit("arbitrary")),
					jen.ID("authenticator").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("authenticator").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.ID("expectedErr"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("authenticator"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("validateLogin").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("authenticator"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid login"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("authenticator").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("authenticator").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleLoginInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("authenticator"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("validateLogin").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("authenticator"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_buildCookie").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("buildCookie").Call(
						jen.Lit("example"),
						jen.Qual("time", "Now").Call().Dot("Add").Call(jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Lifetime")),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid cookie builder"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("cookieManager").Op("=").ID("securecookie").Dot("New").Call(
						jen.ID("securecookie").Dot("GenerateRandomKey").Call(jen.Lit(0)),
						jen.Index().ID("byte").Call(jen.Lit("")),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("buildCookie").Call(
						jen.Lit("example"),
						jen.Qual("time", "Now").Call().Dot("Add").Call(jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Lifetime")),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
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

	return code
}
