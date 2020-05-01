package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestService_CookieAuthenticationMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.ID("md").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
				jen.ID("md").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("User")),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("md"),
				jen.Line(),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID(utils.BuildFakeVarName("User"))),
				utils.RequireNotNil(jen.ID("cookie"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID(constants.RequestVarName).Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("md", "ms"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil user",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.ID("md").Assign().AddressOf().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
				jen.ID("md").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("md"),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID(utils.BuildFakeVarName("User"))),
				utils.RequireNotNil(jen.ID("cookie"), nil),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID(constants.RequestVarName).Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("h").Assign().ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("md", "ms"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without user attached",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("h").Assign().ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("ms"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AuthenticationMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("ocv").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"), jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("User")),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.True()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ocv", "mockDB", "h"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path without allowing cookies",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("ocv").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")), jen.Nil()),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.False()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ocv", "mockDB", "h"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching client but able to use cookie",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")), jen.Nil()),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID(utils.BuildFakeVarName("User"))),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID(constants.RequestVarName).Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.True()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "h"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"able to use cookies but error fetching user info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID(utils.BuildFakeVarName("User"))),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					constants.ObligatoryError(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.ID(constants.RequestVarName).Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.True()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "h"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"no cookies allowed, with error fetching user info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("ocv").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					constants.ObligatoryError(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.False()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ocv", "mockDB", "h"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching client but able to use cookie but unable to decode cookie",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.ID("ocv").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
					constants.ObligatoryError(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("cb").Assign().AddressOf().ID("mockCookieEncoderDecoder").Values(),
				jen.ID("cb").Dot("On").Call(
					jen.Lit("Decode"),
					jen.ID("CookieName"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("cb").Dot("On").Call(
					jen.Lit("Encode"),
					jen.ID("CookieName"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.EmptyString(), jen.Nil()),
				jen.ID("s").Dot("cookieManager").Equals().ID("cb"),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID(utils.BuildFakeVarName("User"))),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID(constants.RequestVarName).Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.True()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("ocv", "cb", "h"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid authentication",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("ocv").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.False()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ocv", "h"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"nightmare path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("ocv").Assign().AddressOf().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("h").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.False()).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ocv", "mockDB", "h"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_parseLoginInputFromForm").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.ID("expected").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.ID(constants.RequestVarName).Dot("Form").Equals().Map(jen.String()).Index().String().Valuesln(
					jen.ID("UsernameFormKey").MapAssign().Values(jen.ID("expected").Dot("Username")),
					jen.ID("PasswordFormKey").MapAssign().Values(jen.ID("expected").Dot("Password")),
					jen.ID("TOTPTokenFormKey").MapAssign().Values(jen.ID("expected").Dot("TOTPToken")),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("parseLoginInputFromForm").Call(jen.ID(constants.RequestVarName)),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"returns nil with error parsing form",
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID(constants.RequestVarName).Dot("URL").Dot("RawQuery").Equals().Lit("%gh&%ij"),
				jen.ID(constants.RequestVarName).Dot("Form").Equals().ID("nil"),
				jen.Line(),
				jen.ID("actual").Assign().ID("parseLoginInputFromForm").Call(jen.ID(constants.RequestVarName)),
				utils.AssertNil(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UserLoginInputMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.Var().ID("b").Qual("bytes", "Buffer"),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(
					jen.ID("t"), jen.Qual("encoding/json", "NewEncoder").Call(
						jen.AddressOf().ID("b"),
					).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Input"))),
				),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.AddressOf().ID("b"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("ms"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.Var().ID("b").Qual("bytes", "Buffer"),
				utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.AddressOf().ID("b")).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Input"))), nil),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.AddressOf().ID("b"),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					constants.ObligatoryError(),
				),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("h").Assign().ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "ms"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request but valid value attached to form",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.ID("form").Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("UsernameFormKey").MapAssign().Values(jen.ID(utils.BuildFakeVarName("Input")).Dot("Username")),
					jen.ID("PasswordFormKey").MapAssign().Values(jen.ID(utils.BuildFakeVarName("Input")).Dot("Password")),
					jen.ID("TOTPTokenFormKey").MapAssign().Values(jen.ID(utils.BuildFakeVarName("Input")).Dot("TOTPToken")),
				),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Dot("Header").Dot("Set").Call(jen.Lit("Content-type"), jen.Lit("application/x-www-form-urlencoded")),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					constants.ObligatoryError(),
				),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertExpectationsFor("ed", "ms"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AdminMiddleware").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("IsAdmin").Equals().True(),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "UserKey"),
						jen.ID(utils.BuildFakeVarName("User")),
					),
				),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ms"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without user attached",
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ms"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with non-admin user",
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				utils.RequireNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("User")).Dot("IsAdmin").Equals().False(),
				jen.Line(),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "UserKey"),
						jen.ID(utils.BuildFakeVarName("User")),
					),
				),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().AddressOf().ID("MockHTTPHandler").Values(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("ms"),
			),
		),
		jen.Line(),
	)

	return ret
}
