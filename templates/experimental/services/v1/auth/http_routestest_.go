package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func httpRoutesTestDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_DecodeCookieFromRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/api/v1/something"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"DecodeCookieFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("cookie")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/api/v1/something"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("c").Op(":=").Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").ID("CookieName"), jen.ID("Value").Op(":").Lit("blah blah blah this is not a real cookie"), jen.ID("Path").Op(":").Lit("/"), jen.ID("HttpOnly").Op(":").ID("true")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"DecodeCookieFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("cookie")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/api/v1/something"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"DecodeCookieFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("err"), jen.Qual("net/http", "ErrNoCookie")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("cookie")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_WebsocketAuthFunction").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with valid oauth2 client"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Assert(jen.Op("*").ID("mockOAuth2ClientValidator")).Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"WebsocketAuthFunction",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"True",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with valid cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("oac").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Assert(jen.Op("*").ID("mockOAuth2ClientValidator")).Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("oac"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"WebsocketAuthFunction",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"True",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with nothing"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("oac").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(),
			jen.ID("s").Dot(
				"oauth2ClientsService",
			).Assert(jen.Op("*").ID("mockOAuth2ClientValidator")).Dot(
				"On",
			).Call(jen.Lit("ExtractOAuth2ClientFromRequest"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("oac"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"WebsocketAuthFunction",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"False",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_FetchUserFromRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("userID"), jen.ID("Username").Op(":").Lit("username")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("expectedUser")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.List(jen.ID("actualUser"), jen.ID("err")).Op(":=").ID("s").Dot(
				"FetchUserFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedUser"), jen.ID("actualUser")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("userID"), jen.ID("Username").Op(":").Lit("username")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.List(jen.ID("actualUser"), jen.ID("err")).Op(":=").ID("s").Dot(
				"FetchUserFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actualUser")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("userID").Op(":=").ID("uint64").Call(jen.Lit(1)),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").ID("userID"), jen.ID("Username").Op(":").Lit("username")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("expectedUser")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.ID("expectedError").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("userID")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.ID("expectedError")),
			jen.List(jen.ID("actualUser"), jen.ID("err")).Op(":=").ID("s").Dot(
				"FetchUserFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actualUser")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_Login").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.ID("s").Dot(
				"LoginHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusNoContent")),
			jen.ID("assert").Dot(
				"NotEmpty",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching login data from request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.Qual("errors", "New").Call(jen.Lit("arbitrary"))),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.ID("s").Dot(
				"LoginHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusUnauthorized")),
			jen.ID("assert").Dot(
				"Empty",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding error fetching login dataa"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.Qual("errors", "New").Call(jen.Lit("arbitrary"))),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.ID("s").Dot(
				"LoginHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusUnauthorized")),
			jen.ID("assert").Dot(
				"Empty",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("false"), jen.ID("nil")),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.ID("s").Dot(
				"LoginHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusUnauthorized")),
			jen.ID("assert").Dot(
				"Empty",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error validating login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.ID("s").Dot(
				"LoginHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusUnauthorized")),
			jen.ID("assert").Dot(
				"Empty",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error building cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("cb").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Valuesln(),
			jen.ID("cb").Dot(
				"On",
			).Call(jen.Lit("Encode"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"cookieManager",
			).Op("=").ID("cb"),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.ID("s").Dot(
				"LoginHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
			jen.ID("assert").Dot(
				"Empty",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie"))),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error building cookie and error encoding cookie response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("cb").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Valuesln(),
			jen.ID("cb").Dot(
				"On",
			).Call(jen.Lit("Encode"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"cookieManager",
			).Op("=").ID("cb"),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Valuesln(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("EncodeResponse"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("nil")),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.ID("s").Dot(
				"LoginHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Code",
			), jen.Qual("net/http", "StatusInternalServerError")),
			jen.ID("assert").Dot(
				"Empty",
			).Call(jen.ID("t"), jen.ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie"))),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_Logout").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(1), jen.ID("Username").Op(":").Lit("username"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("c")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"LogoutHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("actualCookie").Op(":=").ID("res").Dot(
				"Header",
			).Call().Dot(
				"Get",
			).Call(jen.Lit("Set-Cookie")),
			jen.ID("assert").Dot(
				"Contains",
			).Call(jen.ID("t"), jen.ID("actualCookie"), jen.Lit("Max-Age=0")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without cookie"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"LogoutHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_fetchLoginDataFromRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedUser"), jen.ID("nil")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.List(jen.ID("loginData"), jen.ID("err")).Op(":=").ID("s").Dot(
				"fetchLoginDataFromRequest",
			).Call(jen.ID("req")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("loginData")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("loginData").Dot(
				"user",
			), jen.ID("expectedUser")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("without login data attached to request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("s").Dot(
				"fetchLoginDataFromRequest",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with DB error fetching user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("s").Dot(
				"fetchLoginDataFromRequest",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error fetching user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("GetUserByUsername"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("expectedUser").Dot(
				"Username",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"User",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("http://todo.verygoodsoftwarenotvirus.ru/testing"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("exampleLoginData").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("expectedUser").Dot(
				"Username",
			), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("123456")),
			jen.ID("req").Op("=").ID("req").Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("exampleLoginData"))),
			jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("s").Dot(
				"fetchLoginDataFromRequest",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_validateLogin").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("true"),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("exampleInput").Op(":=").ID("loginData").Valuesln(jen.ID("loginInput").Op(":").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("012345")), jen.ID("user").Op(":").Op("&").ID("models").Dot(
				"User",
			).Valuesln()),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"validateLogin",
			).Call(jen.ID("ctx"), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with too weak a password hash"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("true"),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("exampleInput").Op(":=").ID("loginData").Valuesln(jen.ID("loginInput").Op(":").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("012345")), jen.ID("user").Op(":").Op("&").ID("models").Dot(
				"User",
			).Valuesln()),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("auth").Dot(
				"ErrPasswordHashTooWeak",
			)),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Lit("blah"), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"validateLogin",
			).Call(jen.ID("ctx"), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with too weak a password hash and error hashing the password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("false"),
			jen.ID("expectedErr").Op(":=").Qual("errors", "New").Call(jen.Lit("arbitrary")),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("exampleInput").Op(":=").ID("loginData").Valuesln(jen.ID("loginInput").Op(":").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("012345")), jen.ID("user").Op(":").Op("&").ID("models").Dot(
				"User",
			).Valuesln()),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("auth").Dot(
				"ErrPasswordHashTooWeak",
			)),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Lit(""), jen.ID("expectedErr")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"validateLogin",
			).Call(jen.ID("ctx"), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with too weak a password hash and error updating user"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("false"),
			jen.ID("expectedErr").Op(":=").Qual("errors", "New").Call(jen.Lit("arbitrary")),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("exampleInput").Op(":=").ID("loginData").Valuesln(jen.ID("loginInput").Op(":").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("012345")), jen.ID("user").Op(":").Op("&").ID("models").Dot(
				"User",
			).Valuesln()),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("true"), jen.ID("auth").Dot(
				"ErrPasswordHashTooWeak",
			)),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("HashPassword"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Lit("blah"), jen.ID("nil")),
			jen.ID("s").Dot(
				"userDB",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "UserDataManager")).Dot(
				"On",
			).Call(jen.Lit("UpdateUser"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expectedErr")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"validateLogin",
			).Call(jen.ID("ctx"), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error validating login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("false"),
			jen.ID("expectedErr").Op(":=").Qual("errors", "New").Call(jen.Lit("arbitrary")),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("exampleInput").Op(":=").ID("loginData").Valuesln(jen.ID("loginInput").Op(":").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("012345")), jen.ID("user").Op(":").Op("&").ID("models").Dot(
				"User",
			).Valuesln()),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("expectedErr")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"validateLogin",
			).Call(jen.ID("ctx"), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with invalid login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").ID("false"),
			jen.ID("s").Dot(
				"authenticator",
			).Assert(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator")).Dot(
				"On",
			).Call(jen.Lit("ValidateLogin"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("exampleInput").Op(":=").ID("loginData").Valuesln(jen.ID("loginInput").Op(":").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("TOTPToken").Op(":").Lit("012345")), jen.ID("user").Op(":").Op("&").ID("models").Dot(
				"User",
			).Valuesln()),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"validateLogin",
			).Call(jen.ID("ctx"), jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_buildCookie").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("cookie")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error encoding"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("cb").Op(":=").Op("&").ID("mockCookieEncoderDecoder").Valuesln(),
			jen.ID("cb").Dot(
				"On",
			).Call(jen.Lit("Encode"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"cookieManager",
			).Op("=").ID("cb"),
			jen.ID("exampleInput").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("Username").Op(":").Lit("username")),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("exampleInput")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("cookie")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestService_CycleSecret").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("normal operation"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
				"User",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Username").Op(":").Lit("username")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("exampleUser")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("c")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.Null().Var().ID("ca").ID("models").Dot(
				"CookieAuth",
			),
			jen.ID("decodeErr").Op(":=").ID("s").Dot(
				"cookieManager",
			).Dot(
				"Decode",
			).Call(jen.ID("CookieName"), jen.ID("c").Dot(
				"Value",
			), jen.Op("&").ID("ca")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("decodeErr")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.Lit("https://blah.com"), jen.ID("nil")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("req")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("s").Dot(
				"CycleSecretHandler",
			).Call().Call(jen.ID("res"), jen.ID("req")),
			jen.ID("decodeErr2").Op(":=").ID("s").Dot(
				"cookieManager",
			).Dot(
				"Decode",
			).Call(jen.ID("CookieName"), jen.ID("c").Dot(
				"Value",
			), jen.Op("&").ID("ca")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("decodeErr2")),
		)),
	),

		jen.Line(),
	)
	return ret
}