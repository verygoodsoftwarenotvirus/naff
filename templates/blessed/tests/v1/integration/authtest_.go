package integration

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("loginUser").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpSecret")).ID("string")).Params(jen.Op("*").Qual("net/http", "Cookie")).Block(
			jen.ID("loginURL").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s://%s:%s/users/login"),
				jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("URL").Dot("Scheme"),
				jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("URL").Dot("Hostname").Call(),
				jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("URL").Dot("Port").Call(),
			),
			jen.Line(),
			jen.List(jen.ID("code"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.Qual("strings", "ToUpper").Call(jen.ID("totpSecret")), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Line(),
			jen.ID("bodyStr").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit(`
	{
		"username": %q,
		"password": %q,
		"totp_token": %q
	}
`), jen.ID("username"), jen.ID("password"), jen.ID("code")),
			jen.Line(),
			jen.ID("body").Op(":=").Qual("strings", "NewReader").Call(jen.ID("bodyStr")),
			jen.List(jen.ID("req"), jen.ID("_")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("loginURL"), jen.ID("body")),
			jen.List(jen.ID("resp"), jen.Err()).Op(":=").Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Qual("log", "Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("resp").Dot("StatusCode"), jen.Lit("login should be successful")),
			jen.Line(),
			jen.ID("cookies").Op(":=").ID("resp").Dot("Cookies").Call(),
			jen.If(jen.ID("len").Call(jen.ID("cookies")).Op("==").Lit(1)).Block(
				jen.Return().ID("cookies").Index(jen.Lit(0)),
			),
			jen.ID("t").Dot("Logf").Call(jen.Lit("wrong number of cookies found: %d"), jen.ID("len").Call(jen.ID("cookies"))),
			jen.ID("t").Dot("FailNow").Call(),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestAuth").Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to login"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Line(),
				jen.Comment("create a user"),
				jen.ID("ui").Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil/rand/model"), "RandomUserInput").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").ID("todoClient").Dot("BuildCreateUserRequest").Call(jen.ID("tctx"), jen.ID("ui")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Line(),
				jen.Comment("load user response"),
				jen.ID("ucr").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserCreationResponse").Values(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("ucr"))),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("token"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("ucr").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("ucr").Dot("Username"),
					jen.ID("Password").Op(":").ID("ui").Dot("Password"),
					jen.ID("TOTPToken").Op(":").ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u").Dot("Path").Op("=").Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to logout"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Line(),
				jen.ID("ui").Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil/rand/model"), "RandomUserInput").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").ID("todoClient").Dot("BuildCreateUserRequest").Call(jen.ID("tctx"), jen.ID("ui")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Line(),
				jen.ID("ucr").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserCreationResponse").Values(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("ucr"))),
				jen.Line(),
				jen.List(jen.ID("token"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("ucr").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("ucr").Dot("Username"),
					jen.ID("Password").Op(":").ID("ui").Dot("Password"),
					jen.ID("TOTPToken").Op(":").ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u").Dot("Path").Op("=").Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.Comment("extract cookie"),
				jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/require", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
				jen.ID("loginCookie").Op(":=").ID("cookies").Index(jen.Lit(0)),
				jen.Line(),
				jen.Comment("build logout request"),
				jen.List(jen.ID("u2"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u2").Dot("Path").Op("=").Lit("/users/logout"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("loginCookie")),
				jen.Line(),
				jen.Comment("execute logout request"),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode")),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("login request without body fails"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u").Dot("Path").Op("=").Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("StatusCode")),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should not be able to log in with the wrong password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Line(),
				jen.Comment("create a user"),
				jen.ID("ui").Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil/rand/model"), "RandomUserInput").Call(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").ID("todoClient").Dot("BuildCreateUserRequest").Call(jen.ID("tctx"), jen.ID("ui")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Line(),
				jen.Comment("load user response"),
				jen.ID("ucr").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserCreationResponse").Values(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("ucr"))),
				jen.Line(),
				jen.Comment("create login request"),
				jen.Var().ID("badPassword").ID("string"),
				jen.For(jen.List(jen.ID("_"), jen.ID("v")).Op(":=").Range().ID("ui").Dot("Password")).Block(
					jen.ID("badPassword").Op("=").ID("string").Call(jen.ID("v")).Op("+").ID("badPassword"),
				),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("token"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("ucr").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("ucr").Dot("Username"),
					jen.ID("Password").Op(":").ID("badPassword"),
					jen.ID("TOTPToken").Op(":").ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u").Dot("Path").Op("=").Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("StatusCode")),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should not be able to login as someone that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("ui").Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil/rand/model"), "RandomUserInput").Call(),
				jen.Line(),
				jen.List(jen.ID("s"), jen.Err()).Op(":=").ID("randString").Call(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("token"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("s"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("ui").Dot("Username"),
					jen.ID("Password").Op(":").ID("ui").Dot("Password"),
					jen.ID("TOTPToken").Op(":").ID("token")),
				jen.List(jen.ID("out"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u").Dot("Path").Op("=").Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(0)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should reject an unauthenticated request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil(), jen.Lit("webhooks")), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("StatusCode")),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to change password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Comment("create user"),
				jen.List(jen.ID("user"), jen.ID("ui"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
				jen.Line(),
				jen.Comment("create login request"),
				jen.Var().ID("backwardsPass").ID("string"),
				jen.For(jen.List(jen.ID("_"), jen.ID("v")).Op(":=").Range().ID("ui").Dot("Password")).Block(
					jen.ID("backwardsPass").Op("=").ID("string").Call(jen.ID("v")).Op("+").ID("backwardsPass"),
				),
				jen.Line(),
				jen.Comment("create password update request"),
				jen.List(jen.ID("token"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("user").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("r").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "PasswordUpdateInput").Valuesln(
					jen.ID("CurrentPassword").Op(":").ID("ui").Dot("Password"),
					jen.ID("TOTPToken").Op(":").ID("token"),
					jen.ID("NewPassword").Op(":").ID("backwardsPass"),
				),
				jen.List(jen.ID("out"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("r")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u").Dot("Path").Op("=").Lit("/users/password/new"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPut"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.Comment("execute password update request"),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.Comment("logout"),
				jen.Line(),
				jen.List(jen.ID("u2"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u2").Dot("Path").Op("=").Lit("/users/logout"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("newToken"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("user").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("newToken"), jen.Err()),
				jen.List(jen.ID("l"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("user").Dot("Username"),
					jen.ID("Password").Op(":").ID("backwardsPass"),
					jen.ID("TOTPToken").Op(":").ID("newToken")),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op("=").Qual("bytes", "NewReader").Call(jen.ID("l")),
				jen.Line(),
				jen.List(jen.ID("u3"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u3").Dot("Path").Op("=").Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u3").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/require", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
				jen.Qual("github.com/stretchr/testify/assert", "NotEqual").Call(jen.ID("t"), jen.ID("cookie"), jen.ID("cookies").Index(jen.Lit(0))),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should be able to change 2FA Token"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Comment("create user"),
				jen.List(jen.ID("user"), jen.ID("ui"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
				jen.Line(),
				jen.Comment("create TOTP secret update request"),
				jen.List(jen.ID("token"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("user").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("token"), jen.Err()),
				jen.ID("ir").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "TOTPSecretRefreshInput").Valuesln(
					jen.ID("CurrentPassword").Op(":").ID("ui").Dot("Password"),
					jen.ID("TOTPToken").Op(":").ID("token"),
				),
				jen.List(jen.ID("out"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("ir")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op(":=").Qual("bytes", "NewReader").Call(jen.ID("out")),
				jen.Line(),
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u").Dot("Path").Op("=").Lit("/users/totp_secret/new"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.Comment("execute TOTP secret update request"),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusAccepted"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.Comment("load user response"),
				jen.ID("r").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "TOTPSecretRefreshResponse").Values(),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.ID("r"))),
				jen.Qual("github.com/stretchr/testify/require", "NotEqual").Call(jen.ID("t"), jen.ID("user").Dot("TwoFactorSecret"), jen.ID("r").Dot("TwoFactorSecret")),
				jen.Line(),
				jen.Comment("logout"),
				jen.Line(),
				jen.List(jen.ID("u2"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u2").Dot("Path").Op("=").Lit("/users/logout"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u2").Dot("String").Call(), jen.Nil()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.Comment("create login request"),
				jen.List(jen.ID("newToken"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("r").Dot("TwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("newToken"), jen.Err()),
				jen.List(jen.ID("l"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("user").Dot("Username"),
					jen.ID("Password").Op(":").ID("ui").Dot("Password"),
					jen.ID("TOTPToken").Op(":").ID("newToken")),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("body").Op("=").Qual("bytes", "NewReader").Call(jen.ID("l")),
				jen.Line(),
				jen.List(jen.ID("u3"), jen.Err()).Op(":=").Qual("net/url", "Parse").Call(jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil())),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("u3").Dot("Path").Op("=").Lit("/users/login"),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op("=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodPost"), jen.ID("u3").Dot("String").Call(), jen.ID("body")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("req"), jen.Err()),
				jen.Line(),
				jen.Comment("execute login request"),
				jen.List(jen.ID("res"), jen.Err()).Op("=").ID("todoClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("res"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusNoContent"), jen.ID("res").Dot("StatusCode")),
				jen.Line(),
				jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
				jen.Qual("github.com/stretchr/testify/require", "Len").Call(jen.ID("t"), jen.ID("cookies"), jen.Lit(1)),
				jen.Qual("github.com/stretchr/testify/assert", "NotEqual").Call(jen.ID("t"), jen.ID("cookie"), jen.ID("cookies").Index(jen.Lit(0))),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should accept a login cookie if a token is missing"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Comment("create user"),
				jen.List(jen.ID("_"), jen.ID("_"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("BuildURL").Call(jen.Nil(), jen.Lit("webhooks")), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
				jen.Line(),
				jen.List(jen.ID("res"), jen.Err()).Op(":=").Parens(jen.Op("&").Qual("net/http", "Client").Values(jen.ID("Timeout").Op(":").Lit(10).Op("*").Qual("time", "Second"))).Dot("Do").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("StatusCode")),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should only allow users to see their own content"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Line(),
				jen.Comment("create user and oauth2 client A"),
				jen.List(jen.ID("userA"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("ca"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("userA")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("clientA"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "NewClient").Callln(
					jen.ID("tctx"),
					jen.ID("ca").Dot("ClientID"),
					jen.ID("ca").Dot("ClientSecret"),
					jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("URL"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("buildHTTPClient").Call(), jen.ID("ca").Dot("Scopes"),
					jen.ID("true"),
				),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientA"), jen.Err()),
				jen.Line(),
				jen.Comment("create user and oauth2 client B"),
				jen.List(jen.ID("userB"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryUser").Call(jen.ID("urlToUse"), jen.ID("debug")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("cb"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil"), "CreateObligatoryClient").Call(jen.ID("urlToUse"), jen.ID("userB")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("clientB"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "NewClient").Callln(
					jen.ID("tctx"),
					jen.ID("cb").Dot("ClientID"),
					jen.ID("cb").Dot("ClientSecret"),
					jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("URL"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("buildHTTPClient").Call(), jen.ID("cb").Dot("Scopes"),
					jen.ID("true"),
				),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("clientA"), jen.Err()),
				jen.Line(),
				jen.Comment("create webhook for user A"),
				jen.List(jen.ID("webhookA"), jen.Err()).Op(":=").ID("clientA").Dot("CreateWebhook").Call(jen.ID("tctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
					jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"),
					jen.ID("Name").Op(":").Add(utils.FakeStringFunc()),
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("webhookA"), jen.Err()),
				jen.Line(),
				jen.Comment("create webhook for user B"),
				jen.List(jen.ID("webhookB"), jen.Err()).Op(":=").ID("clientB").Dot("CreateWebhook").Call(jen.ID("tctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
					jen.ID("Method").Op(":").Qual("net/http", "MethodPatch"),
					jen.ID("Name").Op(":").Add(utils.FakeStringFunc()),
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("webhookB"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("i"), jen.Err()).Op(":=").ID("clientB").Dot("GetWebhook").Call(jen.ID("tctx"), jen.ID("webhookA").Dot("ID")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("i")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err(), jen.Lit("should experience error trying to fetch entry they're not authorized for")),
				jen.Line(),
				jen.Comment("Clean up"),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(jen.ID("tctx"), jen.ID("webhookA").Dot("ID"))),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("ArchiveWebhook").Call(jen.ID("tctx"), jen.ID("webhookB").Dot("ID"))),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("should only allow clients with a given scope to see that scope's content"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.Line(),
				jen.Comment("create user"),
				jen.List(jen.ID("x"), jen.ID("y"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
				jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("test"), jen.ID("cookie")),
				jen.Line(),
				jen.ID("input").Op(":=").ID("buildDummyOAuth2ClientInput").Call(
					jen.ID("test"),
					jen.ID("x").Dot("Username"),
					jen.ID("y").Dot("Password"),
					jen.ID("x").Dot("TwoFactorSecret"),
				),
				jen.ID("input").Dot("Scopes").Op("=").Index().ID("string").Values(jen.Lit("absolutelynevergonnaexistascopelikethis")),
				jen.List(jen.ID("premade"), jen.Err()).Op(":=").ID("todoClient").Dot("CreateOAuth2Client").Call(jen.ID("tctx"), jen.ID("cookie"), jen.ID("input")),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Op(":=").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "NewClient").Callln(
					utils.CtxVar(),
					jen.ID("premade").Dot("ClientID"),
					jen.ID("premade").Dot("ClientSecret"),
					jen.IDf("%sClient", pkg.Name.UnexportedVarName()).Dot("URL"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("buildHTTPClient").Call(), jen.ID("premade").Dot("Scopes"),
					jen.ID("true"),
				),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("c"), jen.Err()),
				jen.Line(),
				jen.List(jen.ID("i"), jen.Err()).Op(":=").ID("c").Dot("GetOAuth2Clients").Call(utils.CtxVar(), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("i")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err(), jen.Lit("should experience error trying to fetch entry they're not authorized for")),
			)),
		),
		jen.Line(),
	)
	return ret
}
