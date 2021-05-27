package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestService_buildLoginView").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildLoginView").Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("buildLoginView").Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildFormFromLoginRequest").Params(jen.ID("input").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.Qual("net/url", "Values")).Body(
			jen.ID("form").Op(":=").Qual("net/url", "Values").Values(),
			jen.ID("form").Dot("Set").Call(
				jen.ID("usernameFormKey"),
				jen.ID("input").Dot("Username"),
			),
			jen.ID("form").Dot("Set").Call(
				jen.ID("passwordFormKey"),
				jen.ID("input").Dot("Password"),
			),
			jen.ID("form").Dot("Set").Call(
				jen.ID("totpTokenFormKey"),
				jen.ID("input").Dot("TOTPToken"),
			),
			jen.Return().ID("form"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_parseFormEncodedLoginRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
					jen.ID("expectedRedirectTo").Op(":=").Lit("/somewheres"),
					jen.ID("form").Op(":=").ID("buildFormFromLoginRequest").Call(jen.ID("expected")),
					jen.ID("form").Dot("Set").Call(
						jen.ID("redirectToQueryKey"),
						jen.ID("expectedRedirectTo"),
					),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.List(jen.ID("actual"), jen.ID("actualRedirectTo")).Op(":=").ID("s").Dot("parseFormEncodedLoginRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedRedirectTo"),
						jen.ID("actualRedirectTo"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid request body"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("badBody").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockReadCloser").Values(),
					jen.ID("badBody").Dot("On").Call(
						jen.Lit("Read"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("byte").Values()),
					).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("badBody"),
					),
					jen.List(jen.ID("actual"), jen.ID("actualRedirectTo")).Op(":=").ID("s").Dot("parseFormEncodedLoginRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actualRedirectTo"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid form"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("nil"),
					),
					jen.List(jen.ID("actual"), jen.ID("actualRedirectTo")).Op(":=").ID("s").Dot("parseFormEncodedLoginRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actualRedirectTo"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_handleLoginSubmission").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedCookie").Op(":=").Op("&").Qual("net/http", "Cookie").Valuesln(
						jen.ID("Name").Op(":").Lit("testing"), jen.ID("Value").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
					jen.ID("mockAuthService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuthService").Values(),
					jen.ID("mockAuthService").Dot("On").Call(
						jen.Lit("AuthenticateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("expected"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.ID("expectedCookie"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("authService").Op("=").ID("mockAuthService"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromLoginRequest").Call(jen.ID("expected")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("s").Dot("handleLoginSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAuthService"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("htmxRedirectionHeader")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid request content"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("handleLoginSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error authenticating user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
					jen.ID("mockAuthService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuthService").Values(),
					jen.ID("mockAuthService").Dot("On").Call(
						jen.Lit("AuthenticateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("expected"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Parens(jen.Op("*").Qual("net/http", "Cookie")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("authService").Op("=").ID("mockAuthService"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromLoginRequest").Call(jen.ID("expected")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("s").Dot("handleLoginSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAuthService"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.Lit("Set-Cookie")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_handleLogoutSubmission").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("mockAuthService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuthService").Values(),
					jen.ID("mockAuthService").Dot("On").Call(
						jen.Lit("LogoutUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
						jen.ID("res"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("s").Dot("authService").Op("=").ID("mockAuthService"),
					jen.ID("s").Dot("handleLogoutSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAuthService"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("htmxRedirectionHeader")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("handleLogoutSubmission").Call(
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
				jen.Lit("with error logging user out"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleSessionContextData").Op(":=").ID("fakes").Dot("BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil"))),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("mockAuthService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuthService").Values(),
					jen.ID("mockAuthService").Dot("On").Call(
						jen.Lit("LogoutUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
						jen.ID("res"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("s").Dot("authService").Op("=").ID("mockAuthService"),
					jen.ID("s").Dot("handleLogoutSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAuthService"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_registrationComponent").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("registrationComponent").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_registrationView").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("registrationView").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildFormFromRegistrationRequest").Params(jen.ID("input").Op("*").ID("types").Dot("UserRegistrationInput")).Params(jen.Qual("net/url", "Values")).Body(
			jen.ID("form").Op(":=").Qual("net/url", "Values").Values(),
			jen.ID("form").Dot("Set").Call(
				jen.ID("usernameFormKey"),
				jen.ID("input").Dot("Username"),
			),
			jen.ID("form").Dot("Set").Call(
				jen.ID("passwordFormKey"),
				jen.ID("input").Dot("Password"),
			),
			jen.Return().ID("form"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_parseFormEncodedRegistrationRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInput").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromRegistrationRequest").Call(jen.ID("expected")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("actual").Op(":=").ID("s").Dot("parseFormEncodedRegistrationRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid request body"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("badBody").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockReadCloser").Values(),
					jen.ID("badBody").Dot("On").Call(
						jen.Lit("Read"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("byte").Values()),
					).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("badBody"),
					),
					jen.ID("actual").Op(":=").ID("s").Dot("parseFormEncodedRegistrationRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid form"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/verify"),
						jen.ID("nil"),
					),
					jen.ID("actual").Op(":=").ID("s").Dot("parseFormEncodedRegistrationRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_handleRegistrationSubmission").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInput").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromRegistrationRequest").Call(jen.ID("expected")),
					jen.ID("mockUsersService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UsersService").Values(),
					jen.ID("mockUsersService").Dot("On").Call(
						jen.Lit("RegisterUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("expected"),
					).Dot("Return").Call(
						jen.Op("&").ID("types").Dot("UserCreationResponse").Values(),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("usersService").Op("=").ID("mockUsersService"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Dot("handleRegistrationSubmission").Call(
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
						jen.ID("mockUsersService"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("nil"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Dot("handleRegistrationSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error registering user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInput").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromRegistrationRequest").Call(jen.ID("expected")),
					jen.ID("mockUsersService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UsersService").Values(),
					jen.ID("mockUsersService").Dot("On").Call(
						jen.Lit("RegisterUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("expected"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("UserCreationResponse")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("usersService").Op("=").ID("mockUsersService"),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Dot("handleRegistrationSubmission").Call(
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
						jen.ID("mockUsersService"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fake data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("useFakeData").Op("=").ID("true"),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInput").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromRegistrationRequest").Call(jen.ID("expected")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Dot("handleRegistrationSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildFormFromTOTPSecretVerificationRequest").Params(jen.ID("input").Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Params(jen.Qual("net/url", "Values")).Body(
			jen.ID("form").Op(":=").Qual("net/url", "Values").Values(),
			jen.ID("form").Dot("Set").Call(
				jen.ID("totpTokenFormKey"),
				jen.ID("input").Dot("TOTPToken"),
			),
			jen.ID("form").Dot("Set").Call(
				jen.ID("userIDFormKey"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("input").Dot("UserID"),
					jen.Lit(10),
				),
			),
			jen.Return().ID("form"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_parseFormEncodedTOTPSecretVerificationRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInput").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromTOTPSecretVerificationRequest").Call(jen.ID("expected")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/verify"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("actual").Op(":=").ID("s").Dot("parseFormEncodedTOTPSecretVerificationRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid request body"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("badBody").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockReadCloser").Values(),
					jen.ID("badBody").Dot("On").Call(
						jen.Lit("Read"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("byte").Values()),
					).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("badBody"),
					),
					jen.ID("actual").Op(":=").ID("s").Dot("parseFormEncodedTOTPSecretVerificationRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID format"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("form").Op(":=").Qual("net/url", "Values").Valuesln(
						jen.ID("userIDFormKey").Op(":").Valuesln(
							jen.Lit("not a number lol"))),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/verify"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("actual").Op(":=").ID("s").Dot("parseFormEncodedTOTPSecretVerificationRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid form"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("form").Op(":=").Qual("net/url", "Values").Valuesln(
						jen.ID("userIDFormKey").Op(":").Valuesln(
							jen.Lit("0"))),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/verify"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("actual").Op(":=").ID("s").Dot("parseFormEncodedTOTPSecretVerificationRequest").Call(
						jen.ID("ctx"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_handleTOTPVerificationSubmission").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInput").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromTOTPSecretVerificationRequest").Call(jen.ID("expected")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/verify"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("mockUsersService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UsersService").Values(),
					jen.ID("mockUsersService").Dot("On").Call(
						jen.Lit("VerifyUserTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("expected"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("s").Dot("usersService").Op("=").ID("mockUsersService"),
					jen.ID("s").Dot("handleTOTPVerificationSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/verify"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("handleTOTPVerificationSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInput").Call(),
					jen.ID("form").Op(":=").ID("buildFormFromTOTPSecretVerificationRequest").Call(jen.ID("expected")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/verify"),
						jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
					),
					jen.ID("mockUsersService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UsersService").Values(),
					jen.ID("mockUsersService").Dot("On").Call(
						jen.Lit("VerifyUserTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("expected"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("s").Dot("usersService").Op("=").ID("mockUsersService"),
					jen.ID("s").Dot("handleTOTPVerificationSubmission").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
