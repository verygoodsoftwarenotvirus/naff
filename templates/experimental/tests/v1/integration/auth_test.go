package integration

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func authTestDotGo() *jen.File {
	ret := jen.NewFile("integration")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("loginUser").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpSecret")).ID("string")).Params(jen.Op("*").Qual("net/http", "Cookie")).Block(
		jen.ID("loginURL").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s://%s:%s/users/login"), jen.ID("todoClient").Dot(
			"URL",
		).Dot(
			"Scheme",
		), jen.ID("todoClient").Dot(
			"URL",
		).Dot(
			"Hostname",
		).Call(), jen.ID("todoClient").Dot(
			"URL",
		).Dot(
			"Port",
		).Call()),
		jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot(
			"GenerateCode",
		).Call(jen.Qual("strings", "ToUpper").Call(jen.ID("totpSecret")), jen.Qual("time", "Now").Call().Dot(
			"UTC",
		).Call()),
		jen.ID("assert").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.ID("bodyStr").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit(`
	{
		"username": %q,
		"password": %q,
		"totp_token": %q
	}
`), jen.ID("username"), jen.ID("password"), jen.ID("code")),
		jen.ID("body").Op(":=").Qual("strings", "NewReader").Call(jen.ID("bodyStr")),
		jen.List(jen.ID("req"), jen.ID("_")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("loginURL"), jen.ID("body")),
		jen.List(jen.ID("resp"), jen.ID("err")).Op(":=").Qual("net/http", "DefaultClient").Dot(
			"Do",
		).Call(jen.ID("req")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Qual("log", "Fatal").Call(jen.ID("err")),
		),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("resp").Dot(
			"StatusCode",
		), jen.Lit("login should be successful")),
		jen.ID("cookies").Op(":=").ID("resp").Dot(
			"Cookies",
		).Call(),
		jen.If(jen.ID("len").Call(jen.ID("cookies")).Op("==").Lit(1)).Block(
			jen.Return().ID("cookies").Index(jen.Lit(0)),
		),
		jen.ID("t").Dot(
			"Logf",
		).Call(jen.Lit("wrong number of cookies found: %d"), jen.ID("len").Call(jen.ID("cookies"))),
		jen.ID("t").Dot(
			"FailNow",
		).Call(),
		jen.Return().ID("nil"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestAuth").Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
		jen.ID("test").Dot(
			"Parallel",
		).Call(),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should be able to login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("ui").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomUserInput").Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"BuildCreateUserRequest",
			).Call(jen.ID("tctx"), jen.ID("ui")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("ucr").Op(":=").Op("&").ID("models").Dot(
				"UserCreationResponse",
			).Valuesln(),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot(
				"Body",
			)).Dot(
				"Decode",
			).Call(jen.ID("ucr"))),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("ucr").Dot(
				"TwoFactorSecret",
			), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.ID("err")),
			jen.ID("r").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("ucr").Dot(
				"Username",
			), jen.ID("Password").Op(":").ID("ui").Dot(
				"Password",
			), jen.ID("TOTPToken").Op(":").ID("token")),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u").Dot(
				"Path",
			).Op("=").Lit("/users/login"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.ID("cookies").Op(":=").ID("res").Dot(
				"Cookies",
			).Call(),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should be able to logout"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("ui").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomUserInput").Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"BuildCreateUserRequest",
			).Call(jen.ID("tctx"), jen.ID("ui")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("ucr").Op(":=").Op("&").ID("models").Dot(
				"UserCreationResponse",
			).Valuesln(),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot(
				"Body",
			)).Dot(
				"Decode",
			).Call(jen.ID("ucr"))),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("ucr").Dot(
				"TwoFactorSecret",
			), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.ID("err")),
			jen.ID("r").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("ucr").Dot(
				"Username",
			), jen.ID("Password").Op(":").ID("ui").Dot(
				"Password",
			), jen.ID("TOTPToken").Op(":").ID("token")),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u").Dot(
				"Path",
			).Op("=").Lit("/users/login"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.ID("cookies").Op(":=").ID("res").Dot(
				"Cookies",
			).Call(),
			jen.ID("require").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
			jen.ID("loginCookie").Op(":=").ID("cookies").Index(jen.Lit(0)),
			jen.List(jen.ID("u2"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u2").Dot(
				"Path",
			).Op("=").Lit("/users/logout"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot(
				"String",
			).Call(), jen.ID("nil")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("loginCookie")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"StatusCode",
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("login request without body fails"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u").Dot(
				"Path",
			).Op("=").Lit("/users/login"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot(
				"String",
			).Call(), jen.ID("nil")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot(
				"StatusCode",
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should not be able to log in with the wrong password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("ui").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomUserInput").Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"BuildCreateUserRequest",
			).Call(jen.ID("tctx"), jen.ID("ui")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("ucr").Op(":=").Op("&").ID("models").Dot(
				"UserCreationResponse",
			).Valuesln(),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot(
				"Body",
			)).Dot(
				"Decode",
			).Call(jen.ID("ucr"))),
			jen.Null().Var().ID("badPassword").ID("string"),
			jen.For(jen.List(jen.ID("_"), jen.ID("v")).Op(":=").Range().ID("ui").Dot(
				"Password",
			)).Block(
				jen.ID("badPassword").Op("=").ID("string").Call(jen.ID("v")).Op("+").ID("badPassword"),
			),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("ucr").Dot(
				"TwoFactorSecret",
			), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.ID("err")),
			jen.ID("r").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("ucr").Dot(
				"Username",
			), jen.ID("Password").Op(":").ID("badPassword"), jen.ID("TOTPToken").Op(":").ID("token")),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u").Dot(
				"Path",
			).Op("=").Lit("/users/login"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"StatusCode",
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should not be able to login as someone that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("ui").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomUserInput").Call(),
			jen.List(jen.ID("s"), jen.ID("err")).Op(":=").ID("randString").Call(),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("s"), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.ID("err")),
			jen.ID("r").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("ui").Dot(
				"Username",
			), jen.ID("Password").Op(":").ID("ui").Dot(
				"Password",
			), jen.ID("TOTPToken").Op(":").ID("token")),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u").Dot(
				"Path",
			).Op("=").Lit("/users/login"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.ID("cookies").Op(":=").ID("res").Dot(
				"Cookies",
			).Call(),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(0)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should reject an unauthenticated request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"), jen.Lit("webhooks")), jen.ID("nil")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"StatusCode",
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should be able to change password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("user"), jen.ID("ui"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("test"), jen.ID("cookie")),
			jen.Null().Var().ID("backwardsPass").ID("string"),
			jen.For(jen.List(jen.ID("_"), jen.ID("v")).Op(":=").Range().ID("ui").Dot(
				"Password",
			)).Block(
				jen.ID("backwardsPass").Op("=").ID("string").Call(jen.ID("v")).Op("+").ID("backwardsPass"),
			),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("user").Dot(
				"TwoFactorSecret",
			), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.ID("err")),
			jen.ID("r").Op(":=").Op("&").ID("models").Dot(
				"PasswordUpdateInput",
			).Valuesln(jen.ID("CurrentPassword").Op(":").ID("ui").Dot(
				"Password",
			), jen.ID("TOTPToken").Op(":").ID("token"), jen.ID("NewPassword").Op(":").ID("backwardsPass")),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u").Dot(
				"Path",
			).Op("=").Lit("/users/password/new"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPut"), jen.ID("u").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("cookie")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.List(jen.ID("u2"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u2").Dot(
				"Path",
			).Op("=").Lit("/users/logout"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot(
				"String",
			).Call(), jen.ID("nil")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("cookie")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.List(jen.ID("newToken"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("user").Dot(
				"TwoFactorSecret",
			), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("newToken"), jen.ID("err")),
			jen.List(jen.ID("l"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("user").Dot(
				"Username",
			), jen.ID("Password").Op(":").ID("backwardsPass"), jen.ID("TOTPToken").Op(":").ID("newToken"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op("=").Qual("bytes", "NewReader").Call(jen.ID("l")),
			jen.List(jen.ID("u3"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u3").Dot(
				"Path",
			).Op("=").Lit("/users/login"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u3").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.ID("cookies").Op(":=").ID("res").Dot(
				"Cookies",
			).Call(),
			jen.ID("require").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
			jen.ID("assert").Dot(
				"NotEqual",
			).Call(jen.ID("t"), jen.ID("cookie"), jen.ID("cookies").Index(jen.Lit(0))),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should be able to change 2FA Token"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("user"), jen.ID("ui"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
			jen.ID("require").Dot(
				"NotNil",
			).Call(jen.ID("test"), jen.ID("cookie")),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("user").Dot(
				"TwoFactorSecret",
			), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.ID("err")),
			jen.ID("ir").Op(":=").Op("&").ID("models").Dot(
				"TOTPSecretRefreshInput",
			).Valuesln(jen.ID("CurrentPassword").Op(":").ID("ui").Dot(
				"Password",
			), jen.ID("TOTPToken").Op(":").ID("token")),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("ir")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u").Dot(
				"Path",
			).Op("=").Lit("/users/totp_secret/new"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("cookie")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.ID("r").Op(":=").Op("&").ID("models").Dot(
				"TOTPSecretRefreshResponse",
			).Valuesln(),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot(
				"Body",
			)).Dot(
				"Decode",
			).Call(jen.ID("r"))),
			jen.ID("require").Dot(
				"NotEqual",
			).Call(jen.ID("t"), jen.ID("user").Dot(
				"TwoFactorSecret",
			), jen.ID("r").Dot(
				"TwoFactorSecret",
			)),
			jen.List(jen.ID("u2"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u2").Dot(
				"Path",
			).Op("=").Lit("/users/logout"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot(
				"String",
			).Call(), jen.ID("nil")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("cookie")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.List(jen.ID("newToken"), jen.ID("err")).Op(":=").ID("totp").Dot(
				"GenerateCode",
			).Call(jen.ID("r").Dot(
				"TwoFactorSecret",
			), jen.Qual("time", "Now").Call().Dot(
				"UTC",
			).Call()),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("newToken"), jen.ID("err")),
			jen.List(jen.ID("l"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("user").Dot(
				"Username",
			), jen.ID("Password").Op(":").ID("ui").Dot(
				"Password",
			), jen.ID("TOTPToken").Op(":").ID("newToken"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("body").Op("=").Qual("bytes", "NewReader").Call(jen.ID("l")),
			jen.List(jen.ID("u3"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("u3").Dot(
				"Path",
			).Op("=").Lit("/users/login"),
			jen.List(jen.ID("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u3").Dot(
				"String",
			).Call(), jen.ID("body")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.ID("err")),
			jen.List(jen.ID("res"), jen.ID("err")).Op("=").ID("todoClient").Dot(
				"PlainClient",
			).Call().Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot(
				"StatusCode",
			)),
			jen.ID("cookies").Op(":=").ID("res").Dot(
				"Cookies",
			).Call(),
			jen.ID("require").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
			jen.ID("assert").Dot(
				"NotEqual",
			).Call(jen.ID("t"), jen.ID("cookie"), jen.ID("cookies").Index(jen.Lit(0))),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should accept a login cookie if a token is missing"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("_"), jen.ID("_"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("test"), jen.ID("cookie")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.ID("todoClient").Dot(
				"BuildURL",
			).Call(jen.ID("nil"), jen.Lit("webhooks")), jen.ID("nil")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"AddCookie",
			).Call(jen.ID("cookie")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").Parens(jen.Op("&").Qual("net/http", "Client").Valuesln(jen.ID("Timeout").Op(":").Lit(10).Op("*").Qual("time", "Second"))).Dot(
				"Do",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"StatusCode",
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should only allow users to see their own content"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
			jen.List(jen.ID("userA"), jen.ID("err")).Op(":=").ID("testutil").Dot(
				"CreateObligatoryUser",
			).Call(jen.ID("urlToUse"), jen.ID("debug")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("ca"), jen.ID("err")).Op(":=").ID("testutil").Dot(
				"CreateObligatoryClient",
			).Call(jen.ID("urlToUse"), jen.ID("userA")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("clientA"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "NewClient").Call(jen.ID("tctx"), jen.ID("ca").Dot(
				"ClientID",
			), jen.ID("ca").Dot(
				"ClientSecret",
			), jen.ID("todoClient").Dot(
				"URL",
			), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("buildHTTPClient").Call(), jen.ID("ca").Dot(
				"Scopes",
			), jen.ID("true")),
			jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientA"), jen.ID("err")),
			jen.List(jen.ID("userB"), jen.ID("err")).Op(":=").ID("testutil").Dot(
				"CreateObligatoryUser",
			).Call(jen.ID("urlToUse"), jen.ID("debug")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("cb"), jen.ID("err")).Op(":=").ID("testutil").Dot(
				"CreateObligatoryClient",
			).Call(jen.ID("urlToUse"), jen.ID("userB")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.List(jen.ID("clientB"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "NewClient").Call(jen.ID("tctx"), jen.ID("cb").Dot(
				"ClientID",
			), jen.ID("cb").Dot(
				"ClientSecret",
			), jen.ID("todoClient").Dot(
				"URL",
			), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("buildHTTPClient").Call(), jen.ID("cb").Dot(
				"Scopes",
			), jen.ID("true")),
			jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientA"), jen.ID("err")),
			jen.List(jen.ID("webhookA"), jen.ID("err")).Op(":=").ID("clientA").Dot(
				"CreateWebhook",
			).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Name").Op(":").Lit("A"))),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("webhookA"), jen.ID("err")),
			jen.List(jen.ID("webhookB"), jen.ID("err")).Op(":=").ID("clientB").Dot(
				"CreateWebhook",
			).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
				"WebhookCreationInput",
			).Valuesln(jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"), jen.ID("Name").Op(":").Lit("B"))),
			jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("webhookB"), jen.ID("err")),
			jen.List(jen.ID("i"), jen.ID("err")).Op(":=").ID("clientB").Dot(
				"GetWebhook",
			).Call(jen.ID("tctx"), jen.ID("webhookA").Dot(
				"ID",
			)),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("i")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err"), jen.Lit("should experience error trying to fetch entry they're not authorized for")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("todoClient").Dot(
				"ArchiveWebhook",
			).Call(jen.ID("tctx"), jen.ID("webhookA").Dot(
				"ID",
			))),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("todoClient").Dot(
				"ArchiveWebhook",
			).Call(jen.ID("tctx"), jen.ID("webhookB").Dot(
				"ID",
			))),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("should only allow clients with a given scope to see that scope's content"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
			jen.List(jen.ID("x"), jen.ID("y"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("test"), jen.ID("cookie")),
			jen.ID("input").Op(":=").ID("buildDummyOAuth2ClientInput").Call(jen.ID("test"), jen.ID("x").Dot(
				"Username",
			), jen.ID("y").Dot(
				"Password",
			), jen.ID("x").Dot(
				"TwoFactorSecret",
			)),
			jen.ID("input").Dot(
				"Scopes",
			).Op("=").Index().ID("string").Valuesln(jen.Lit("absolutelynevergonnaexistascopelikethis")),
			jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
				"CreateOAuth2Client",
			).Call(jen.ID("tctx"), jen.ID("cookie"), jen.ID("input")),
			jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.ID("err")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "NewClient").Call(jen.Qual("context", "Background").Call(), jen.ID("premade").Dot(
				"ClientID",
			), jen.ID("premade").Dot(
				"ClientSecret",
			), jen.ID("todoClient").Dot(
				"URL",
			), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("buildHTTPClient").Call(), jen.ID("premade").Dot(
				"Scopes",
			), jen.ID("true")),
			jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("c"), jen.ID("err")),
			jen.List(jen.ID("i"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetOAuth2Clients",
			).Call(jen.ID("tctx"), jen.ID("nil")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("i")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err"), jen.Lit("should experience error trying to fetch entry they're not authorized for")),
		)),
	),

		jen.Line(),
	)
	return ret
}
