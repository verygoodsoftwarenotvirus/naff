package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAuthenticationService_issueSessionManagedCookie").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("expectedToken"), jen.ID("err")).Op(":=").ID("random").Dot("GenerateBase64EncodedString").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(32),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("userIDContextKey"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("accountIDContextKey"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Commit"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("expectedToken"),
						jen.Qual("time", "Now").Call().Dot("Add").Call(jen.Lit(24).Op("*").Qual("time", "Hour")),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("issueSessionManagedCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Var().Defs(
						jen.ID("actualToken").ID("string"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("cookieManager").Dot("Decode").Call(
							jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
							jen.ID("cookie").Dot("Value"),
							jen.Op("&").ID("actualToken"),
						),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedToken"),
						jen.ID("actualToken"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error loading from session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("issueSessionManagedCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error renewing token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("issueSessionManagedCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("expectedToken"), jen.ID("err")).Op(":=").ID("random").Dot("GenerateBase64EncodedString").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(32),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("userIDContextKey"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("accountIDContextKey"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Commit"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("expectedToken"),
						jen.Qual("time", "Now").Call(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("issueSessionManagedCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("expectedToken"), jen.ID("err")).Op(":=").ID("random").Dot("GenerateBase64EncodedString").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(32),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("userIDContextKey"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("accountIDContextKey"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Commit"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("expectedToken"),
						jen.Qual("time", "Now").Call().Dot("Add").Call(jen.Lit(24).Op("*").Qual("time", "Hour")),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("cookieManager").Op("=").ID("securecookie").Dot("New").Call(
						jen.ID("securecookie").Dot("GenerateRandomKey").Call(jen.Lit(0)),
						jen.Index().ID("byte").Call(jen.Lit("")),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("helper").Dot("service").Dot("issueSessionManagedCookie").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("require").Dot("Nil").Call(
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

	code.Add(
		jen.Func().ID("TestAuthenticationService_LoginHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("GetDefaultAccountIDForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.ID("auditLog").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(),
					jen.ID("auditLog").Dot("On").Call(
						jen.Lit("LogSuccessfulLoginEvent"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLog"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("membershipDB"),
						jen.ID("auditLog"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing login data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("nil")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no results in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving user from datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with banned user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus").Op("=").ID("types").Dot("BannedUserAccountStatus"),
					jen.ID("helper").Dot("exampleUser").Dot("ReputationExplanation").Op("=").Lit("bad behavior"),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
					jen.ID("auditLog").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(),
					jen.ID("auditLog").Dot("On").Call(
						jen.Lit("LogBannedUserLoginAttemptEvent"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLog"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusForbidden"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("auditLog"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid login"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("auditLog").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(),
					jen.ID("auditLog").Dot("On").Call(
						jen.Lit("LogUnsuccessfulLoginBadPasswordEvent"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLog"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("auditLog"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error validating login"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("authenticator"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid two factor code error returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
						jen.ID("authentication").Dot("ErrInvalidTOTPToken"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("authenticator"),
					jen.ID("auditLog").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(),
					jen.ID("auditLog").Dot("On").Call(
						jen.Lit("LogUnsuccessfulLoginBad2FATokenEvent"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLog"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("auditLog"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with non-matching password error returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
						jen.ID("authentication").Dot("ErrPasswordDoesNotMatch"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("authenticator"),
					jen.ID("auditLog").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(),
					jen.ID("auditLog").Dot("On").Call(
						jen.Lit("LogUnsuccessfulLoginBadPasswordEvent"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLog"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("auditLog"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching default account"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("GetDefaultAccountIDForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("uint64").Call(jen.Lit(0)),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("membershipDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error loading from session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("GetDefaultAccountIDForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("membershipDB"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error renewing token in session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("GetDefaultAccountIDForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("membershipDB"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing to session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("GetDefaultAccountIDForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("userIDContextKey"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("accountIDContextKey"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Commit"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("time", "Now").Call(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("membershipDB"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("cb").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Valuesln(),
					jen.ID("cb").Dot("On").Call(
						jen.Lit("Encode"),
						jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("string")),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("cookieManager").Op("=").ID("cb"),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("GetDefaultAccountIDForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("cb"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("membershipDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building cookie and error encoding cookie response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleLoginInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("cb").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Valuesln(),
					jen.ID("cb").Dot("On").Call(
						jen.Lit("Encode"),
						jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("string")),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("cookieManager").Op("=").ID("cb"),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUserByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("GetDefaultAccountIDForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.ID("helper").Dot("service").Dot("BeginSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("cb"),
						jen.ID("userDataManager"),
						jen.ID("authenticator"),
						jen.ID("membershipDB"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_ChangeActiveAccountHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("nil")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("ChangeActiveAccountInput").Valuesln(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error checking user account membership"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without account authorization"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error loading from session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error renewing token in session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing to session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("userIDContextKey"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Put"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("accountIDContextKey"),
						jen.ID("exampleInput").Dot("AccountID"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Commit"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("time", "Now").Call(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error renewing token in session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("RenewToken"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeChangeActiveAccountInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("accountMembershipManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("accountMembershipManager").Dot("On").Call(
						jen.Lit("UserIsMemberOfAccount"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("AccountID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("accountMembershipManager"),
					jen.ID("cookieManager").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Valuesln(),
					jen.ID("cookieManager").Dot("On").Call(
						jen.Lit("Encode"),
						jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("string")),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("cookieManager").Op("=").ID("cookieManager"),
					jen.ID("helper").Dot("service").Dot("ChangeActiveAccountHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("accountMembershipManager"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_LogoutHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("auditLog").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(),
					jen.ID("auditLog").Dot("On").Call(
						jen.Lit("LogLogoutEvent"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLog"),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("helper").Dot("service").Dot("EndSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusSeeOther"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("actualCookie").Op(":=").ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					jen.ID("assert").Dot("Contains").Call(
						jen.ID("t"),
						jen.ID("actualCookie"),
						jen.Lit("Max-Age=0"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("auditLog"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("EndSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error loading from session manager"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.Qual("context", "Background").Call(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("EndSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("actualCookie").Op(":=").ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actualCookie"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error deleting from session store"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("sm").Op(":=").Op("&").ID("mockSessionManager").Valuesln(),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Load"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Lit(""),
					).Dot("Return").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("sm").Dot("On").Call(
						jen.Lit("Destroy"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("sessionManager").Op("=").ID("sm"),
					jen.ID("helper").Dot("service").Dot("EndSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("actualCookie").Op(":=").ID("helper").Dot("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actualCookie"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("helper").Dot("service").Dot("cookieManager").Op("=").ID("securecookie").Dot("New").Call(
						jen.ID("securecookie").Dot("GenerateRandomKey").Call(jen.Lit(0)),
						jen.Index().ID("byte").Call(jen.Lit("")),
					),
					jen.ID("helper").Dot("service").Dot("EndSessionHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_StatusHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("StatusHandler").Call(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with problem fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("StatusHandler").Call(
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
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_CycleSecretHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("=").Index().ID("string").Valuesln(jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.ID("helper").Dot("setContextFetcher").Call(jen.ID("t")),
					jen.ID("auditLog").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Valuesln(),
					jen.ID("auditLog").Dot("On").Call(
						jen.Lit("LogCycleCookieSecretEvent"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLog"),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("c").Op(":=").ID("helper").Dot("req").Dot("Cookies").Call().Index(jen.Lit(0)),
					jen.Var().Defs(
						jen.ID("token").ID("string"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("cookieManager").Dot("Decode").Call(
							jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
							jen.ID("c").Dot("Value"),
							jen.Op("&").ID("token"),
						),
					),
					jen.ID("helper").Dot("service").Dot("CycleCookieSecretHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected code to be %d, but was %d"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("cookieManager").Dot("Decode").Call(
							jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
							jen.ID("c").Dot("Value"),
							jen.Op("&").ID("token"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("auditLog"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error getting session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("c").Op(":=").ID("helper").Dot("req").Dot("Cookies").Call().Index(jen.Lit(0)),
					jen.Var().Defs(
						jen.ID("token").ID("string"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("cookieManager").Dot("Decode").Call(
							jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
							jen.ID("c").Dot("Value"),
							jen.Op("&").ID("token"),
						),
					),
					jen.ID("helper").Dot("service").Dot("CycleCookieSecretHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected code to be %d, but was %d"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("cookieManager").Dot("Decode").Call(
							jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
							jen.ID("c").Dot("Value"),
							jen.Op("&").ID("token"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid permissions"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.List(jen.ID("helper").Dot("ctx"), jen.ID("helper").Dot("req"), jen.ID("_")).Op("=").ID("attachCookieToRequestForTest").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service"),
						jen.ID("helper").Dot("req"),
						jen.ID("helper").Dot("exampleUser"),
					),
					jen.ID("c").Op(":=").ID("helper").Dot("req").Dot("Cookies").Call().Index(jen.Lit(0)),
					jen.Var().Defs(
						jen.ID("token").ID("string"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("cookieManager").Dot("Decode").Call(
							jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
							jen.ID("c").Dot("Value"),
							jen.Op("&").ID("token"),
						),
					),
					jen.ID("helper").Dot("service").Dot("CycleCookieSecretHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusForbidden"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected code to be %d, but was %d"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("service").Dot("cookieManager").Dot("Decode").Call(
							jen.ID("helper").Dot("service").Dot("config").Dot("Cookies").Dot("Name"),
							jen.ID("c").Dot("Value"),
							jen.Op("&").ID("token"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuthenticationService_PASETOHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("Lifetime").Op("=").Qual("time", "Minute"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("AccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("expected").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("expected"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.Var().Defs(
						jen.ID("bodyBytes").Qual("bytes", "Buffer"),
					),
					jen.ID("marshalErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("bodyBytes")).Dot("Encode").Call(jen.ID("exampleInput")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("marshalErr"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("bodyBytes").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Var().Defs(
						jen.ID("result").Op("*").ID("types").Dot("PASETOResponse"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("helper").Dot("res").Dot("Body")).Dot("Decode").Call(jen.Op("&").ID("result")),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("result").Dot("Token"),
					),
					jen.Var().Defs(
						jen.ID("targetPayload").ID("paseto").Dot("JSONToken"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("paseto").Dot("NewV2").Call().Dot("Decrypt").Call(
							jen.ID("result").Dot("Token"),
							jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey"),
							jen.Op("&").ID("targetPayload"),
							jen.ID("nil"),
						),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("targetPayload").Dot("Expiration").Dot("After").Call(jen.Qual("time", "Now").Call().Dot("UTC").Call()),
					),
					jen.ID("payload").Op(":=").ID("targetPayload").Dot("Get").Call(jen.ID("pasetoDataKey")),
					jen.List(jen.ID("gobEncoding"), jen.ID("err")).Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("DecodeString").Call(jen.ID("payload")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Var().Defs(
						jen.ID("actual").Op("*").ID("types").Dot("SessionContextData"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("encoding/gob", "NewDecoder").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("gobEncoding"))).Dot("Decode").Call(jen.Op("&").ID("actual")),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("userDataManager"),
						jen.ID("membershipDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("does not issue token with longer lifetime than package maximum"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("Lifetime").Op("=").Lit(24).Op("*").Qual("time", "Hour").Op("*").Lit(365),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("expected").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("helper").Dot("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("helper").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("helper").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("helper").Dot("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("ActiveAccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("AccountPermissions").Op(":").ID("helper").Dot("examplePermCheckers")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("expected"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.Var().Defs(
						jen.ID("bodyBytes").Qual("bytes", "Buffer"),
					),
					jen.ID("marshalErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("bodyBytes")).Dot("Encode").Call(jen.ID("exampleInput")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("marshalErr"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("bodyBytes").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Var().Defs(
						jen.ID("result").Op("*").ID("types").Dot("PASETOResponse"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("helper").Dot("res").Dot("Body")).Dot("Decode").Call(jen.Op("&").ID("result")),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("result").Dot("Token"),
					),
					jen.Var().Defs(
						jen.ID("targetPayload").ID("paseto").Dot("JSONToken"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("paseto").Dot("NewV2").Call().Dot("Decrypt").Call(
							jen.ID("result").Dot("Token"),
							jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey"),
							jen.Op("&").ID("targetPayload"),
							jen.ID("nil"),
						),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("targetPayload").Dot("Expiration").Dot("Before").Call(jen.Qual("time", "Now").Call().Dot("UTC").Call().Dot("Add").Call(jen.ID("maxPASETOLifetime"))),
					),
					jen.ID("payload").Op(":=").ID("targetPayload").Dot("Get").Call(jen.ID("pasetoDataKey")),
					jen.List(jen.ID("gobEncoding"), jen.ID("err")).Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("DecodeString").Call(jen.ID("payload")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Var().Defs(
						jen.ID("actual").Op("*").ID("types").Dot("SessionContextData"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("encoding/gob", "NewDecoder").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("gobEncoding"))).Dot("Decode").Call(jen.Op("&").ID("actual")),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("userDataManager"),
						jen.ID("membershipDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("nil")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("Lifetime").Op("=").Qual("time", "Minute"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid request time"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Lit(1)),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error decoding signature header"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("jsonBytes")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base32", "HexEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching API client"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClient")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
					jen.Var().Defs(
						jen.ID("bodyBytes").Qual("bytes", "Buffer"),
					),
					jen.ID("marshalErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("bodyBytes")).Dot("Encode").Call(jen.ID("exampleInput")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("marshalErr"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("bodyBytes").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
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
						jen.ID("apiClientDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
					jen.Var().Defs(
						jen.ID("bodyBytes").Qual("bytes", "Buffer"),
					),
					jen.ID("marshalErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("bodyBytes")).Dot("Encode").Call(jen.ID("exampleInput")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("marshalErr"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("bodyBytes").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
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
						jen.ID("apiClientDataManager"),
						jen.ID("userDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching account memberships"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("SessionContextData")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.Var().Defs(
						jen.ID("bodyBytes").Qual("bytes", "Buffer"),
					),
					jen.ID("marshalErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("bodyBytes")).Dot("Encode").Call(jen.ID("exampleInput")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("marshalErr"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("bodyBytes").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
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
						jen.ID("apiClientDataManager"),
						jen.ID("userDataManager"),
						jen.ID("membershipDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid checksum"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.Index().ID("byte").Call(jen.Lit("lol"))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
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
						jen.ID("apiClientDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with inadequate account permissions"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("fakes").Dot("BuildFakeAPIClient").Call().Dot("ClientSecret"),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("Lifetime").Op("=").Qual("time", "Minute"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("AccountID").Op(":").ID("helper").Dot("exampleAccount").Dot("ID"), jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
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
					jen.ID("delete").Call(
						jen.ID("helper").Dot("sessionCtxData").Dot("AccountPermissions"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					),
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("sessionCtxData"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.Var().Defs(
						jen.ID("bodyBytes").Qual("bytes", "Buffer"),
					),
					jen.ID("marshalErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("bodyBytes")).Dot("Encode").Call(jen.ID("exampleInput")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("marshalErr"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("bodyBytes").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with token encryption error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("config").Dot("PASETO").Dot("LocalModeKey").Op("=").ID("nil"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("helper").Dot("exampleAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByClientID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientManager").Op("=").ID("apiClientDataManager"),
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
					jen.ID("membershipDB").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AccountUserMembershipDataManager").Valuesln(),
					jen.ID("membershipDB").Dot("On").Call(
						jen.Lit("BuildSessionContextDataForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("sessionCtxData"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("accountMembershipManager").Op("=").ID("membershipDB"),
					jen.Var().Defs(
						jen.ID("bodyBytes").Qual("bytes", "Buffer"),
					),
					jen.ID("marshalErr").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("bodyBytes")).Dot("Encode").Call(jen.ID("exampleInput")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("marshalErr"),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
					),
					jen.List(jen.ID("_"), jen.ID("macWriteErr")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("bodyBytes").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("macWriteErr"),
					),
					jen.ID("sigHeader").Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("mac").Dot("Sum").Call(jen.ID("nil"))),
					jen.ID("helper").Dot("req").Dot("Header").Dot("Set").Call(
						jen.ID("signatureHeaderKey"),
						jen.ID("sigHeader"),
					),
					jen.ID("helper").Dot("service").Dot("PASETOHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("userDataManager"),
						jen.ID("membershipDB"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
