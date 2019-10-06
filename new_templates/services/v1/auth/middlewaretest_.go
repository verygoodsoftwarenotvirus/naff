package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareTestDotGo() *jen.File {
	ret := jen.NewFile("auth")
	utils.AddImports(ret)

	ret.Add(jen.Func().ID("TestService_CookieAuthenticationMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("md").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager").Valuesln(),
			jen.ID("md").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("md"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("exampleUser")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("cookie")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("cookie")),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("ms").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"CookieAuthenticationMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with nil user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("md").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager").Valuesln(),
			jen.ID("md").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("md"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("exampleUser")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("cookie")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("cookie")),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("ms").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"CookieAuthenticationMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusUnauthorized")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without user attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("ms").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"CookieAuthenticationMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestService_AuthenticationMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"), jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"), jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleClient"), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call().Dot(
				"UserDataManager",
			),
			jen.ID("mockDB").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleClient").Dot(
				"BelongsTo",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("mockDB"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("true")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path without allowing cookies"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"), jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"), jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleClient"), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call().Dot(
				"UserDataManager",
			),
			jen.ID("mockDB").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleClient").Dot(
				"BelongsTo",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("mockDB"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("false")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching client but able to use cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username")),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("exampleUser")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call().Dot(
				"UserDataManager",
			),
			jen.ID("mockDB").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleUser").Dot(
				"ID",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("mockDB"),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("true")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("able to use cookies but error fetching user info"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username")),
			jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"), jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"), jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleClient"), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call().Dot(
				"UserDataManager",
			),
			jen.ID("mockDB").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleClient").Dot(
				"BelongsTo",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("mockDB"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("exampleUser")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("true")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("no cookies allowed, with error fetching user info"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"), jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"), jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleClient"), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call().Dot(
				"UserDataManager",
			),
			jen.ID("mockDB").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleClient").Dot(
				"BelongsTo",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("mockDB"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("false")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching client but able to use cookie but unable to decode cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username"))),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.ID("cb").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Valuesln(),
			jen.ID("cb").Dot(
				"On",
			).Call(jen.Lit("Decode"), jen.ID("CookieName"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"cookieManager",
			).Op("=").ID("cb"),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("true")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid authentication"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("false")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusUnauthorized")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("nightmare path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123)),
			jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(jen.ID("ClientID").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_ID"), jen.ID("ClientSecret").Op(":").Lit("PRETEND_THIS_IS_A_REAL_CLIENT_SECRET"), jen.ID("BelongsTo").Op(":").ID("exampleUser").Dot(
				"ID",
			)),
			jen.ID("ocv").Op(":=").Op("&").ID("mockOAuth2ClientValidator").Valuesln(),
			jen.ID("ocv").Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("exampleClient"), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Op("=").ID("ocv"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call().Dot(
				"UserDataManager",
			),
			jen.ID("mockDB").Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("exampleClient").Dot(
				"BelongsTo",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Op("=").ID("mockDB"),
			jen.ID("h").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"AuthenticationMiddleware",
			).Call(jen.ID("false")).Call(jen.ID("h")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	)
	ret.Add(jen.Func().ID("Test_parseLoginInputFromForm").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Dot(
				"Form",
			).Op("=").Map(jen.ID("string")).Index().ID("string").Valuesln(jen.ID("UsernameFormKey").Op(":").Valuesln(jen.ID("expected").Dot(
				"Username",
			)), jen.ID("PasswordFormKey").Op(":").Valuesln(jen.ID("expected").Dot(
				"Password",
			)), jen.ID("TOTPTokenFormKey").Op(":").Valuesln(jen.ID("expected").Dot(
				"TOTPToken",
			))),
			jen.ID("actual").Op(":=").ID("parseLoginInputFromForm").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("returns nil with error parsing form"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"RawQuery",
			).Op("=").Lit("%gh&%ij"),
			jen.ID("req").Dot(
				"Form",
			).Op("=").ID("nil"),
			jen.ID("actual").Op(":=").ID("parseLoginInputFromForm").Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestService_UserLoginInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("1233456")),
			jen.Null().Var().ID("b").Qual("bytes", "Buffer"),
			utils.RequireNoError(jen.ID("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b")).Dot(
				"Encode",
			).Call(jen.ID("exampleInput"))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Op("&").ID("b")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("ms").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"UserLoginInputMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("ms").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("1233456")),
			jen.Null().Var().ID("b").Qual("bytes", "Buffer"),
			utils.RequireNoError(jen.ID("t"), jen.Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b")).Dot(
				"Encode",
			).Call(jen.ID("exampleInput"))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Op("&").ID("b")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"UserLoginInputMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("ms").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error decoding request but valid value attached to form"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("1233456")),
			jen.ID("form").Op(":=").Qual("net/url", "Values").Valuesln(jen.ID("UsernameFormKey").Op(":").Valuesln(jen.ID("exampleInput").Dot(
				"Username",
			)), jen.ID("PasswordFormKey").Op(":").Valuesln(jen.ID("exampleInput").Dot(
				"Password",
			)), jen.ID("TOTPTokenFormKey").Op(":").Valuesln(jen.ID("exampleInput").Dot(
				"TOTPToken",
			))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot(
				"Encode",
			).Call())),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Dot(
				"Header",
			).Dot(
				"Set",
			).Call(jen.Lit("Content-type"), jen.Lit("application/x-www-form-urlencoded")),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("ms").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"UserLoginInputMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("ms").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestService_AdminMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserKey",
			), jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("IsAdmin").Op(":").ID("true")))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("ms").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"AdminMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("ms").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without user attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"AdminMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("ms").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with non-admin user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru"), jen.ID("nil")),
			utils.RequireNoError(jen.ID("t"), jen.ID("err")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("models").Dot(
				"UserKey",
			), jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("IsAdmin").Op(":").ID("false")))),
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ms").Op(":=").Op("&").ID("MockHTTPHandler").Valuesln(),
			jen.ID("h").Op(":=").ID("s").Dot(
				"AdminMiddleware",
			).Call(jen.ID("ms")),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("ms").Dot(
				"AssertExpectations",
			).Call(jen.ID("t")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	)
	return ret
}
