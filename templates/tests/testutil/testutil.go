package testutil

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func testutilDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildInit()...)
	code.Add(buildDetermineServiceURL()...)
	code.Add(buildEnsureServerIsUp()...)
	code.Add(buildIsUp()...)
	code.Add(buildCreateObligatoryUser(proj)...)
	code.Add(buildBuildURL()...)
	code.Add(buildGetLoginCookie(proj)...)
	code.Add(buildCreateObligatoryClient(proj)...)

	return code
}

func buildInit() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.Qual(constants.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
		),
		jen.Line(),
	}

	return lines
}

func buildDetermineServiceURL() []jen.Code {
	lines := []jen.Code{
		jen.Comment("DetermineServiceURL returns the URL, if properly configured."),
		jen.Line(),
		jen.Func().ID("DetermineServiceURL").Params().Params(jen.String()).Body(
			jen.ID("ta").Assign().Qual("os", "Getenv").Call(jen.Lit("TARGET_ADDRESS")),
			jen.If(jen.ID("ta").IsEqualTo().EmptyString()).Body(
				jen.ID("panic").Call(jen.Lit("must provide target address!")),
			),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("ta")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.ID("svcAddr").Assign().ID("u").Dot("String").Call(),
			jen.Line(),
			jen.Qual("log", "Printf").Call(jen.Lit("using target address: %q\n"), jen.ID("svcAddr")),
			jen.Return().ID("svcAddr"),
		),
		jen.Line(),
	}

	return lines
}

func buildEnsureServerIsUp() []jen.Code {
	lines := []jen.Code{
		jen.Comment("EnsureServerIsUp checks that a server is up and doesn't return until it's certain one way or the other."),
		jen.Line(),
		jen.Func().ID("EnsureServerIsUp").Params(jen.ID("address").String()).Body(
			jen.Var().Defs(
				jen.ID("isDown").Equals().True(),
				jen.ID("interval").Equals().Qual("time", "Second"),
				jen.ID("maxAttempts").Equals().Lit(50),
				jen.ID("numberOfAttempts").Equals().Zero(),
			),
			jen.Line(),
			jen.For(jen.ID("isDown")).Body(
				jen.If(jen.Not().ID("IsUp").Call(jen.ID("address"))).Body(
					jen.Qual("log", "Print").Call(jen.Lit("waiting before pinging again")),
					jen.Qual("time", "Sleep").Call(jen.ID("interval")),
					jen.ID("numberOfAttempts").Op("++"),
					jen.If(jen.ID("numberOfAttempts").Op(">=").ID("maxAttempts")).Body(
						jen.Qual("log", "Fatal").Call(jen.Lit("Maximum number of attempts made, something's gone awry")),
					),
				).Else().Body(
					jen.ID("isDown").Equals().False(),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildIsUp() []jen.Code {
	lines := []jen.Code{
		jen.Comment("IsUp can check if an instance of our server is alive."),
		jen.Line(),
		jen.Func().ID("IsUp").Params(jen.ID("address").String()).Params(jen.Bool()).Body(
			jen.ID("uri").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%s/_meta_/ready"), jen.ID("address")),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.ID("uri"), jen.Nil()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID(constants.RequestVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().False(),
			),
			jen.Line(),
			jen.If(jen.Err().Equals().ID(constants.ResponseVarName).Dot("Body").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual("log", "Println").Call(jen.Lit("error closing body")),
			),
			jen.Line(),
			jen.Return().ID(constants.ResponseVarName).Dot("StatusCode").IsEqualTo().Qual("net/http", "StatusOK"),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateObligatoryUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateObligatoryUser creates a user for the sake of having an OAuth2 client."),
		jen.Line(),
		jen.Func().ID("CreateObligatoryUser").Params(jen.ID("address").String(), jen.ID("debug").Bool()).Params(jen.PointerTo().Qual(proj.TypesPackage(),
			"User",
		), jen.Error()).Body(
			constants.CreateCtx(),
			jen.List(jen.ID("tu"), jen.ID("parseErr")).Assign().Qual("net/url", "Parse").Call(jen.ID("address")),
			jen.If(jen.ID("parseErr").DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.ID("parseErr")),
			),
			jen.Line(),
			jen.List(jen.ID("c"), jen.ID("clientInitErr")).Assign().Qual(proj.HTTPClientV1Package(), "NewSimpleClient").Call(constants.CtxVar(), jen.ID("tu"), jen.ID("debug")),
			jen.If(jen.ID("clientInitErr").DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.ID("clientInitErr")),
			),
			jen.Line(),
			jen.Comment("I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,"),
			jen.Comment("could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit"),
			jen.ID("username").Assign().Qual(constants.FakeLibrary, "Username").Call().Plus().Qual(constants.FakeLibrary, "HexColor").Call().Plus().Qual(constants.FakeLibrary, "Country").Call(),
			jen.ID("in").Assign().AddressOf().Qual(proj.TypesPackage(), "UserCreationInput").Valuesln(
				jen.ID("Username").MapAssign().ID("username"),
				jen.ID("Password").MapAssign().Qual(constants.FakeLibrary, "Password").Call(jen.True(), jen.True(), jen.True(), jen.True(), jen.True(), jen.Lit(64)),
			),
			jen.Line(),
			jen.List(jen.ID("ucr"), jen.ID("userCreationErr")).Assign().ID("c").Dot("CreateUser").Call(constants.CtxVar(), jen.ID("in")),
			jen.If(jen.ID("userCreationErr").DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.ID("userCreationErr")),
			).Else().If(jen.ID("ucr").IsEqualTo().ID("nil")).Body(
				jen.Return().List(jen.Nil(), utils.Error("something happened")),
			),
			jen.Line(),
			jen.List(jen.ID("token"), jen.ID("tokenErr")).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
				jen.ID("ucr").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.ID("tokenErr").DoesNotEqual().Nil()).Body(
				jen.Return(
					jen.Nil(),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("generating totp code: %w"),
						jen.ID("tokenErr"),
					),
				),
			),
			jen.Line(),
			jen.If(
				jen.ID("validationErr").Assign().ID("c").Dot("VerifyTOTPSecret").Call(
					constants.CtxVar(),
					jen.ID("ucr").Dot("ID"),
					jen.ID("token"),
				),
				jen.ID("validationErr").DoesNotEqual().Nil(),
			).Body(
				jen.Return(
					jen.Nil(),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("verifying totp code: %w"),
						jen.ID("validationErr"),
					),
				),
			),
			jen.Line(),
			jen.ID("u").Assign().AddressOf().Qual(proj.TypesPackage(), "User").Valuesln(
				jen.ID("ID").MapAssign().ID("ucr").Dot("ID"),
				jen.ID("Username").MapAssign().ID("ucr").Dot("Username"),
				jen.Comment("this is a dirty trick to reuse most of this model"),
				jen.ID("HashedPassword").MapAssign().ID("in").Dot("Password"),
				jen.ID("TwoFactorSecret").MapAssign().ID("ucr").Dot("TwoFactorSecret"),
				jen.ID("PasswordLastChangedOn").MapAssign().ID("ucr").Dot("PasswordLastChangedOn"),
				jen.ID("CreatedOn").MapAssign().ID("ucr").Dot("CreatedOn"),
				jen.ID("LastUpdatedOn").MapAssign().ID("ucr").Dot("LastUpdatedOn"),
				jen.ID("ArchivedOn").MapAssign().ID("ucr").Dot("ArchivedOn"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildURL() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildURL").Params(jen.ID("address").String(), jen.ID("parts").Spread().String()).Params(jen.String()).Body(
			jen.List(jen.ID("tu"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("address")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.Qual("strings", "Join").Call(jen.ID("parts"), jen.Lit("/"))),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.Return().ID("tu").Dot("ResolveReference").Call(jen.ID("u")).Dot("String").Call(),
		),
		jen.Line(),
	}

	return lines
}

func buildGetLoginCookie(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("getLoginCookie").Params(jen.ID("serviceURL").String(), jen.ID("u").PointerTo().Qual(proj.TypesPackage(),
			"User",
		)).Params(jen.PointerTo().Qual("net/http", "Cookie"), jen.Error()).Body(
			jen.ID("uri").Assign().ID("buildURL").Call(jen.ID("serviceURL"), jen.Lit("users"), jen.Lit("login")),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.Qual("strings", "ToUpper").Call(jen.ID("u").Dot("TwoFactorSecret")), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("generating totp token: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodPost"), jen.ID("uri"), jen.Qual("strings", "NewReader").Callln(
					jen.Qual("fmt", "Sprintf").Callln(
						jen.Lit(`
					{
						"username": %q,
						"password": %q,
						"totpToken": %q
					}
				`),
						jen.ID("u").Dot("Username"),
						jen.ID("u").Dot("HashedPassword"),
						jen.ID("code"),
					),
				),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID(constants.RequestVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.Err())),
			),
			jen.Line(),
			jen.If(jen.Err().Equals().ID(constants.ResponseVarName).Dot("Body").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Qual("log", "Println").Call(jen.Lit("error closing body")),
			),
			jen.Line(),
			jen.ID("cookies").Assign().ID(constants.ResponseVarName).Dot("Cookies").Call(),
			jen.If(jen.Len(jen.ID("cookies")).GreaterThan().Zero()).Body(
				jen.Return().List(jen.ID("cookies").Index(jen.Zero()), jen.Nil()),
			),
			jen.Line(),
			jen.Return().List(jen.Nil(), utils.Error("no cookie found :(")),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateObligatoryClient(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateObligatoryClient creates the OAuth2 client we need for tests."),
		jen.Line(),
		jen.Func().ID("CreateObligatoryClient").Params(jen.ID("serviceURL").String(), jen.ID("u").PointerTo().Qual(proj.TypesPackage(), "User")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client"), jen.Error()).Body(
			jen.If(jen.ID("u").IsEqualTo().Nil()).Body(
				jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("user is nil"))),
			),
			jen.Line(),
			jen.ID("firstOAuth2ClientURI").Assign().ID("buildURL").Call(jen.ID("serviceURL"), jen.Lit("oauth2"), jen.Lit("client")),
			jen.Line(),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Callln(
				jen.Qual("strings", "ToUpper").Call(jen.ID("u").Dot("TwoFactorSecret")),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodPost"), jen.ID("firstOAuth2ClientURI"), jen.Qual("strings", "NewReader").Call(
					utils.FormatString(`
	{
		"username": %q,
		"password": %q,
		"totpToken": %q,
		"belongsToUser": %d,
		"scopes": ["*"]
	}
		`,
						jen.ID("u").Dot("Username"),
						jen.ID("u").Dot("HashedPassword"),
						jen.ID("code"),
						jen.ID("u").Dot("ID"),
					),
				),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("getLoginCookie").Call(jen.ID("serviceURL"), jen.ID("u")),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID("cookie").IsEqualTo().ID("nil")).Body(
				jen.Qual("log", "Fatalf").Call(jen.Lit(`
cookie problems!
	cookie == nil: %v
			  err: %v
	`), jen.ID("cookie").IsEqualTo().ID("nil"), jen.Err()),
			),
			jen.ID(constants.RequestVarName).Dot("AddCookie").Call(jen.ID("cookie")),
			jen.Var().ID("o").Qual(proj.TypesPackage(), "OAuth2Client"),
			jen.Line(),
			jen.Var().ID("command").Qual("fmt", "Stringer"),
			jen.If(jen.List(jen.ID("command"), jen.Err()).Equals().Qual("github.com/moul/http2curl", "GetCurlCommand").Call(jen.ID(constants.RequestVarName)), jen.Err().IsEqualTo().ID("nil")).Body(
				jen.Qual("log", "Println").Call(jen.ID("command").Dot("String").Call()),
			),
			jen.Line(),
			jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID(constants.RequestVarName)),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Err()),
			).Else().If(jen.ID(constants.ResponseVarName).Dot("StatusCode").DoesNotEqual().Qual("net/http", "StatusCreated")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("bad status: %d"), jen.ID(constants.ResponseVarName).Dot("StatusCode"))),
			),
			jen.Line(),
			jen.Defer().Func().Params().Body(
				jen.If(jen.Err().Equals().ID(constants.ResponseVarName).Dot("Body").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Qual("log", "Fatal").Call(jen.Err()),
				),
			).Call(),
			jen.Line(),
			jen.List(jen.ID("bdump"), jen.Err()).Assign().ID("httputil").Dot("DumpResponse").Call(jen.ID(constants.ResponseVarName), jen.True()),
			jen.If(jen.Err().IsEqualTo().ID("nil").And().ID(constants.RequestVarName).Dot("Method").DoesNotEqual().Qual("net/http", "MethodGet")).Body(
				jen.Qual("log", "Println").Call(jen.String().Call(jen.ID("bdump"))),
			),
			jen.Line(),
			jen.Return().List(jen.AddressOf().ID("o"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID(constants.ResponseVarName).Dot("Body")).Dot("Decode").Call(jen.AddressOf().ID("o"))),
		),
		jen.Line(),
	}

	return lines
}
