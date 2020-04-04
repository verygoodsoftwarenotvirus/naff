package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestService_CookieAuthenticationMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Line(),
				jen.ID("md").Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
				jen.ID("md").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleUser"),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("md"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("cookie")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil user",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Line(),
				jen.ID("md").Assign().VarPointer().Qual(proj.ModelsV1Package("mock"), "UserDataManager").Values(),
				jen.ID("md").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("md"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("cookie")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without user attached",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("CookieAuthenticationMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AuthenticationMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleOAuth2Client"),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(
					jen.ID("exampleUser"),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path without allowing cookies",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleOAuth2Client"),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(jen.ID("exampleUser"), jen.Nil()),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching client but able to use cookie",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUser").Dot("ID"),
				).Dot("Return").Call(jen.ID("exampleUser"), jen.Nil()),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"able to use cookies but error fetching user info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleOAuth2Client"),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUser").Dot("ID"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"no cookies allowed, with error fetching user info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleOAuth2Client"),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching client but able to use cookie but unable to decode cookie",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.Line(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID("exampleUser")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("c")),
				jen.Line(),
				jen.ID("cb").Assign().VarPointer().ID("mockCookieEncoderDecoder").Values(),
				jen.ID("cb").Dot("On").Call(
					jen.Lit("Decode"),
					jen.ID("CookieName"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.ID("s").Dot("cookieManager").Equals().ID("cb"),
				jen.Line(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("true")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid authentication",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"nightmare path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.ID("ocv").Assign().VarPointer().ID("mockOAuth2ClientValidator").Values(),
				jen.ID("ocv").Dot("On").Call(
					jen.Lit("ExtractOAuth2ClientFromRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleOAuth2Client"),
					jen.Nil(),
				),
				jen.ID("s").Dot("oauth2ClientsService").Equals().ID("ocv"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call().Dot("UserDataManager"),
				jen.ID("mockDB").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")).Call(jen.Nil()),
					jen.Nil(),
				),
				jen.ID("s").Dot("userDB").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("h").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("s").Dot("AuthenticationMiddleware").Call(jen.ID("false")).Call(jen.ID("h")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_parseLoginInputFromForm").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("expected").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.ID("req").Dot("Form").Equals().Map(jen.String()).Index().String().Valuesln(
					jen.ID("UsernameFormKey").MapAssign().Values(jen.ID("expected").Dot("Username")),
					jen.ID("PasswordFormKey").MapAssign().Values(jen.ID("expected").Dot("Password")),
					jen.ID("TOTPTokenFormKey").MapAssign().Values(jen.ID("expected").Dot("TOTPToken")),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("parseLoginInputFromForm").Call(jen.ID("req")),
				utils.AssertNotNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"returns nil with error parsing form",
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodGet"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("req").Dot("URL").Dot("RawQuery").Equals().Lit("%gh&%ij"),
				jen.ID("req").Dot("Form").Equals().ID("nil"),
				jen.Line(),
				jen.ID("actual").Assign().ID("parseLoginInputFromForm").Call(jen.ID("req")),
				utils.AssertNil(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UserLoginInputMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.Var().ID("b").Qual("bytes", "Buffer"),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(
					jen.ID("t"), jen.Qual("encoding/json", "NewEncoder").Call(
						jen.VarPointer().ID("b"),
					).Dot("Encode").Call(jen.ID("exampleInput")),
				),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.VarPointer().ID("b"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.Var().ID("b").Qual("bytes", "Buffer"),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.VarPointer().ID("b")).Dot("Encode").Call(jen.ID("exampleInput"))),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.VarPointer().ID("b"),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("h").Assign().ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request but valid value attached to form",
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserLoginInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.ID("form").Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("UsernameFormKey").MapAssign().Values(jen.ID("exampleInput").Dot("Username")),
					jen.ID("PasswordFormKey").MapAssign().Values(jen.ID("exampleInput").Dot("Password")),
					jen.ID("TOTPTokenFormKey").MapAssign().Values(jen.ID("exampleInput").Dot("TOTPToken")),
				),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Dot("Header").Dot("Set").Call(jen.Lit("Content-type"), jen.Lit("application/x-www-form-urlencoded")),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("DecodeRequest"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("UserLoginInputMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AdminMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleUser").Dot("IsAdmin").Equals().True(),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "UserKey"),
						jen.ID("exampleUser"),
					),
				),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.ID("ms").Dot("On").Call(
					jen.Lit("ServeHTTP"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without user attached",
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with non-admin user",
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"),
					jen.Nil(),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
				jen.ID("exampleUser").Dot("IsAdmin").Equals().False(),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						jen.ID("req").Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "UserKey"),
						jen.ID("exampleUser"),
					),
				),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("ms").Assign().VarPointer().ID("MockHTTPHandler").Values(),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("AdminMiddleware").Call(jen.ID("ms")),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				jen.ID("ms").Dot("AssertExpectations").Call(jen.ID("t")),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
