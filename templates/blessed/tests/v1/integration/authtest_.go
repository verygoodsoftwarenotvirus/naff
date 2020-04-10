package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("loginUser").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpSecret")).String()).Params(jen.PointerTo().Qual("net/http", "Cookie")).Block(
			jen.ID("loginURL").Assign().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s://%s:%s/users/login"),
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL").Dot("Scheme"),
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL").Dot("Hostname").Call(),
				jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL").Dot("Port").Call(),
			),
			jen.Line(),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.Qual("strings", "ToUpper").Call(jen.ID("totpSecret")), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			utils.AssertNoError(jen.Err(), nil),
			jen.Line(),
			jen.ID("bodyStr").Assign().Qual("fmt", "Sprintf").Call(jen.Lit(`
	{
		"username": %q,
		"password": %q,
		"totp_token": %q
	}
`), jen.ID("username"), jen.ID("password"), jen.ID("code")),
			jen.Line(),
			jen.ID("body").Assign().Qual("strings", "NewReader").Call(jen.ID("bodyStr")),
			jen.List(jen.ID("req"), jen.Underscore()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("loginURL"), jen.ID("body")),
			jen.List(jen.ID("resp"), jen.Err()).Assign().Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Qual("log", "Fatal").Call(jen.Err()),
			),
			jen.Line(),
			utils.AssertEqual(jen.Qual("net/http", "StatusNoContent"), jen.ID("resp").Dot("StatusCode"), jen.Lit("login should be successful"), nil),
			jen.Line(),
			jen.ID("cookies").Assign().ID("resp").Dot("Cookies").Call(),
			jen.If(jen.ID("len").Call(jen.ID("cookies")).Op("==").One()).Block(
				jen.Return().ID("cookies").Index(jen.Zero()),
			),
			jen.ID("t").Dot("Logf").Call(jen.Lit("wrong number of cookies found: %d"), jen.ID("len").Call(jen.ID("cookies"))),
			jen.ID("t").Dot("FailNow").Call(),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestAuth").Params(jen.ID("test").PointerTo().Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to login"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.Line(),
				jen.Comment("create a user"),
				jen.ID("ui").Assign().Qual(proj.FakeModelsPackage(), "RandomUserInput").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().ID("todoClient").Dot("BuildCreateUserRequest").Call(utils.CtxVar(), jen.ID("ui")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Line(),
				jen.Comment("load user response"),
				jen.ID("ucr").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserCreationResponse").Values(),
				utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("ucr")), nil),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("ucr").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("ucr").Dot("Username"),
					jen.ID("Password").MapAssign().ID("ui").Dot("Password"),
					jen.ID("TOTPToken").MapAssign().ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Assign().Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u").Dot("Path").Equals().Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.ID("cookies").Assign().ID("res").Dot("Cookies").Call(),
				utils.AssertLength(jen.ID("cookies"), jen.One(), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to logout"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.Line(),
				jen.ID("ui").Assign().Qual(proj.FakeModelsPackage(), "RandomUserInput").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().ID("todoClient").Dot("BuildCreateUserRequest").Call(utils.CtxVar(), jen.ID("ui")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Line(),
				jen.ID("ucr").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserCreationResponse").Values(),
				utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("ucr")), nil),
				jen.Line(),
				jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("ucr").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("ucr").Dot("Username"),
					jen.ID("Password").MapAssign().ID("ui").Dot("Password"),
					jen.ID("TOTPToken").MapAssign().ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Assign().Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u").Dot("Path").Equals().Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.Comment("extract cookie"),
				jen.ID("cookies").Assign().ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/require", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.One()),
				jen.ID("loginCookie").Assign().ID("cookies").Index(jen.Zero()),
				jen.Line(),
				jen.Comment("build logout request"),
				jen.List(jen.ID("u2"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u2").Dot("Path").Equals().Lit("/users/logout"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("loginCookie")),
				jen.Line(),
				jen.Comment("execute logout request"),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode"), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("login request without body fails"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u").Dot("Path").Equals().Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("StatusCode"), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should not be able to log in with the wrong password"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.Line(),
				jen.Comment("create a user"),
				jen.ID("ui").Assign().Qual(proj.FakeModelsPackage(), "RandomUserInput").Call(),
				jen.List(jen.ID("req"), jen.Err()).Assign().ID("todoClient").Dot("BuildCreateUserRequest").Call(utils.CtxVar(), jen.ID("ui")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Line(),
				jen.Comment("load user response"),
				jen.ID("ucr").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserCreationResponse").Values(),
				utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("ucr")), nil),
				jen.Line(),
				jen.Comment("create login request"),
				jen.Var().ID("badPassword").String(),
				jen.For(jen.List(jen.Underscore(), jen.ID("v")).Assign().Range().ID("ui").Dot("Password")).Block(
					jen.ID("badPassword").Equals().String().Call(jen.ID("v")).Op("+").ID("badPassword"),
				),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("ucr").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("ucr").Dot("Username"),
					jen.ID("Password").MapAssign().ID("badPassword"),
					jen.ID("TOTPToken").MapAssign().ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Assign().Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u").Dot("Path").Equals().Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("StatusCode"), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should not be able to login as someone that doesn't exist"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.ID("ui").Assign().Qual(proj.FakeModelsPackage(), "RandomUserInput").Call(),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Assign().ID("randString").Call(),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("s"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("ui").Dot("Username"),
					jen.ID("Password").MapAssign().ID("ui").Dot("Password"),
					jen.ID("TOTPToken").MapAssign().ID("token")),
				jen.List(jen.ID("out"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Assign().Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u").Dot("Path").Equals().Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.ID("cookies").Assign().ID("res").Dot("Cookies").Call(),
				utils.AssertLength(jen.ID("cookies"), jen.Zero(), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should reject an unauthenticated request"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil(), jen.Lit("webhooks")), jen.Nil()),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("StatusCode"), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to change password"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.Comment("create user"),
				jen.List(jen.ID("user"), jen.ID("ui"), jen.ID("cookie")).Assign().ID("buildDummyUser").Call(jen.ID("test")),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
				jen.Line(),
				jen.Comment("create login request"),
				jen.Var().ID("backwardsPass").String(),
				jen.For(jen.List(jen.Underscore(), jen.ID("v")).Assign().Range().ID("ui").Dot("Password")).Block(
					jen.ID("backwardsPass").Equals().String().Call(jen.ID("v")).Op("+").ID("backwardsPass"),
				),
				jen.Line(),
				jen.Comment("create password update request"),
				jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("user").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Assign().AddressOf().Qual(proj.ModelsV1Package(), "PasswordUpdateInput").Valuesln(
					jen.ID("CurrentPassword").MapAssign().ID("ui").Dot("Password"),
					jen.ID("TOTPToken").MapAssign().ID("token"),
					jen.ID("NewPassword").MapAssign().ID("backwardsPass"),
				),
				jen.List(jen.ID("out"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Assign().Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u").Dot("Path").Equals().Lit("/users/password/new"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPut"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.Comment("execute password update request"),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.Comment("logout"),
				jen.Line(),
				jen.List(jen.ID("u2"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u2").Dot("Path").Equals().Lit("/users/logout"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("newToken"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("user").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("newToken"), jen.Err()),
				jen.List(jen.ID("l"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("Password").MapAssign().ID("backwardsPass"),
					jen.ID("TOTPToken").MapAssign().ID("newToken")),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Equals().Qual("bytes", "NewReader").Call(jen.ID("l")),
				jen.Line(),
				jen.List(jen.ID("u3"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u3").Dot("Path").Equals().Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u3").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.ID("cookies").Assign().ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/require", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.One()),
				utils.AssertNotEqual(jen.ID("cookie"), jen.ID("cookies").Index(jen.Zero()), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to change 2FA Token"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.Comment("create user"),
				jen.List(jen.ID("user"), jen.ID("ui"), jen.ID("cookie")).Assign().ID("buildDummyUser").Call(jen.ID("test")),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
				jen.Line(),
				jen.Comment("create TOTP secret update request"),
				jen.List(jen.ID("token"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("user").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("ir").Assign().AddressOf().Qual(proj.ModelsV1Package(), "TOTPSecretRefreshInput").Valuesln(
					jen.ID("CurrentPassword").MapAssign().ID("ui").Dot("Password"),
					jen.ID("TOTPToken").MapAssign().ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("ir")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Assign().Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u").Dot("Path").Equals().Lit("/users/totp_secret/new"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.Comment("execute TOTP secret update request"),
				jen.List(jen.ID("res"), jen.Err()).Assign().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.Comment("load user response"),
				jen.ID("r").Assign().AddressOf().Qual(proj.ModelsV1Package(), "TOTPSecretRefreshResponse").Values(),
				utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("r")), nil),
				jen.Qual("github.com/stretchr/testify/require", "NotEqual").Call(jen.ID("t"), jen.ID("user").Dot("TwoFactorSecret"), jen.ID("r").Dot("TwoFactorSecret")),
				jen.Line(),
				jen.Comment("logout"),
				jen.Line(),
				jen.List(jen.ID("u2"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u2").Dot("Path").Equals().Lit("/users/logout"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("newToken"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("r").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("newToken"), jen.Err()),
				jen.List(jen.ID("l"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
					jen.ID("Password").MapAssign().ID("ui").Dot("Password"),
					jen.ID("TOTPToken").MapAssign().ID("newToken")),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("body").Equals().Qual("bytes", "NewReader").Call(jen.ID("l")),
				jen.Line(),
				jen.List(jen.ID("u3"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("u3").Dot("Path").Equals().Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Equals().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u3").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Equals().ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				utils.AssertEqual(jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode"), nil),
				jen.Line(),
				jen.ID("cookies").Assign().ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/require", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.One()),
				utils.AssertNotEqual(jen.ID("cookie"), jen.ID("cookies").Index(jen.Zero()), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should accept a login cookie if a token is missing"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.Comment("create user"),
				jen.List(jen.Underscore(), jen.Underscore(), jen.ID("cookie")).Assign().ID("buildDummyUser").Call(jen.ID("test")),
				utils.AssertNotNil(jen.ID("cookie"), nil),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil(), jen.Lit("webhooks")), jen.Nil()),
				utils.AssertNoError(jen.Err(), nil),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Assign().Parens(jen.AddressOf().Qual("net/http", "Client").Values(jen.ID("Timeout").MapAssign().Lit(10).Times().Qual("time", "Second"))).Dot("Do").Call(jen.ID("req")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode"), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should only allow users to see their own content"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.Line(),
				jen.Comment("create user and oauth2 client A"),
				jen.List(jen.ID("userA"), jen.Err()).Assign().Qual(proj.TestutilV1Package(), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("ca"), jen.Err()).Assign().Qual(proj.TestutilV1Package(), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("userA")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("clientA"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
					utils.CtxVar(),
					jen.ID("ca").Dot("ClientID"),
					jen.ID("ca").Dot("ClientSecret"),
					jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("buildHTTPClient").Call(), jen.ID("ca").Dot("Scopes"),
					jen.True(),
				),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientA"), jen.Err()),
				jen.Line(),
				jen.Comment("create user and oauth2 client B"),
				jen.List(jen.ID("userB"), jen.Err()).Assign().Qual(proj.TestutilV1Package(), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("cb"), jen.Err()).Assign().Qual(proj.TestutilV1Package(), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("userB")),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.List(jen.ID("clientB"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
					utils.CtxVar(),
					jen.ID("cb").Dot("ClientID"),
					jen.ID("cb").Dot("ClientSecret"),
					jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("buildHTTPClient").Call(), jen.ID("cb").Dot("Scopes"),
					jen.True(),
				),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientA"), jen.Err()),
				jen.Line(),
				jen.Comment("create webhook for user A"),
				jen.List(jen.ID("webhookA"), jen.Err()).Assign().ID("clientA").Dot("CreateWebhook").Call(utils.CtxVar(), jen.AddressOf().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Valuesln(
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("webhookA"), jen.Err()),
				jen.Line(),
				jen.Comment("create webhook for user B"),
				jen.List(jen.ID("webhookB"), jen.Err()).Assign().ID("clientB").Dot("CreateWebhook").Call(utils.CtxVar(), jen.AddressOf().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Valuesln(
					jen.ID("Method").MapAssign().Qual("net/http", "MethodPatch"),
					jen.ID("Name").MapAssign().Add(utils.FakeStringFunc()),
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("webhookB"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("i"), jen.Err()).Assign().ID("clientB").Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("webhookA").Dot("ID")),
				utils.AssertNil(jen.ID("i"), nil),
				utils.AssertError(jen.Err(), jen.Lit("should experience error trying to fetch entry they're not authorized for")),
				jen.Line(),
				jen.Comment("Clean up"),
				utils.AssertNoError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("webhookA").Dot("ID")), nil),
				utils.AssertNoError(jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("webhookB").Dot("ID")), nil),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should only allow clients with a given scope to see that scope's content"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.Line(),
				jen.Comment("create user"),
				jen.List(jen.ID("x"), jen.ID("y"), jen.ID("cookie")).Assign().ID("buildDummyUser").Call(jen.ID("test")),
				utils.AssertNotNil(jen.ID("cookie"), nil),
				jen.Line(),
				jen.ID("input").Assign().ID("buildDummyOAuth2ClientInput").Call(
					jen.ID("test"),
					jen.ID("x").Dot("Username"),
					jen.ID("y").Dot("Password"),
					jen.ID("x").Dot("TwoFactorSecret"),
				),
				jen.ID("input").Dot("Scopes").Equals().Index().String().Values(jen.Lit("absolutelynevergonnaexistascopelikethis")),
				jen.List(jen.ID("premade"), jen.Err()).Assign().ID("todoClient").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("cookie"), jen.ID("input")),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewClient").Callln(
					utils.CtxVar(),
					jen.ID("premade").Dot("ClientID"),
					jen.ID("premade").Dot("ClientSecret"),
					jen.IDf("%sClient", proj.Name.UnexportedVarName()).Dot("URL"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("buildHTTPClient").Call(), jen.ID("premade").Dot("Scopes"),
					jen.True(),
				),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("c"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("i"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Nil()),
				utils.AssertNil(jen.ID("i"), nil),
				utils.AssertError(jen.Err(), jen.Lit("should experience error trying to fetch entry they're not authorized for")),
			)),
		),
		jen.Line(),
	)
	return ret
}
