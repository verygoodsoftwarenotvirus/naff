package testutil

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func testutilDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("testutil")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("DetermineServiceURL returns the URL, if properly configured"),
		jen.Line(),
		jen.Func().ID("DetermineServiceURL").Params().Params(jen.String()).Block(
			jen.ID("ta").Assign().Qual("os", "Getenv").Call(jen.Lit("TARGET_ADDRESS")),
			jen.If(jen.ID("ta").Op("==").EmptyString()).Block(
				jen.ID("panic").Call(jen.Lit("must provide target address!")),
			),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("ta")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("EnsureServerIsUp checks that a server is up and doesn't return until it's certain one way or the other"),
		jen.Line(),
		jen.Func().ID("EnsureServerIsUp").Params(jen.ID("address").String()).Block(
			jen.Var().Defs(
				jen.ID("isDown").Equals().True(),
				jen.ID("interval").Equals().Qual("time", "Second"),
				jen.ID("maxAttempts").Equals().Lit(50),
				jen.ID("numberOfAttempts").Equals().Zero(),
			),
			jen.Line(),
			jen.For(jen.ID("isDown")).Block(
				jen.If(jen.Op("!").ID("IsUp").Call(jen.ID("address"))).Block(
					jen.Qual("log", "Print").Call(jen.Lit("waiting before pinging again")),
					jen.Qual("time", "Sleep").Call(jen.ID("interval")),
					jen.ID("numberOfAttempts").Op("++"),
					jen.If(jen.ID("numberOfAttempts").Op(">=").ID("maxAttempts")).Block(
						jen.Qual("log", "Fatal").Call(jen.Lit("Maximum number of attempts made, something's gone awry")),
					),
				).Else().Block(
					jen.ID("isDown").Equals().False(),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsUp can check if an instance of our server is alive"),
		jen.Line(),
		jen.Func().ID("IsUp").Params(jen.ID("address").String()).Params(jen.Bool()).Block(
			jen.ID("uri").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%s/_meta_/ready"), jen.ID("address")),
			jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.ID("uri"), jen.Nil()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.Err()).Assign().Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().False(),
			),
			jen.Line(),
			jen.If(jen.Err().Equals().ID("res").Dot("Body").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Qual("log", "Println").Call(jen.Lit("error closing body")),
			),
			jen.Line(),
			jen.Return().ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusOK"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateObligatoryUser creates a user for the sake of having an OAuth2 client"),
		jen.Line(),
		jen.Func().ID("CreateObligatoryUser").Params(jen.ID("address").String(), jen.ID("debug").Bool()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(),
			"User",
		), jen.Error()).Block(
			jen.List(jen.ID("tu"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("address")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("c"), jen.Err()).Assign().Qual(proj.HTTPClientV1Package(), "NewSimpleClient").Call(utils.InlineCtx(), jen.ID("tu"), jen.ID("debug")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Comment("I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,"),
			jen.Comment("could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit"),
			jen.ID("username").Assign().Qual(utils.FakeLibrary, "Username").Call().Op("+").Qual(utils.FakeLibrary, "HexColor").Call().Op("+").Qual(utils.FakeLibrary, "Country").Call(),
			jen.ID("in").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Valuesln(
				jen.ID("Username").MapAssign().ID("username"),
				jen.ID("Password").MapAssign().Qual(utils.FakeLibrary, "Password").Call(jen.True(), jen.True(), jen.True(), jen.True(), jen.True(), jen.Lit(64)),
			),
			jen.Line(),
			jen.List(jen.ID("ucr"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(utils.InlineCtx(), jen.ID("in")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			).Else().If(jen.ID("ucr").Op("==").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("something happened"))),
			),
			jen.Line(),
			jen.ID("u").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Valuesln(
				jen.ID("ID").MapAssign().ID("ucr").Dot("ID"),
				jen.ID("Username").MapAssign().ID("ucr").Dot("Username"),
				jen.Comment("this is a dirty trick to reuse most of this model"),
				jen.ID("HashedPassword").MapAssign().ID("in").Dot("Password"),
				jen.ID("TwoFactorSecret").MapAssign().ID("ucr").Dot("TwoFactorSecret"),
				jen.ID("PasswordLastChangedOn").MapAssign().ID("ucr").Dot("PasswordLastChangedOn"),
				jen.ID("CreatedOn").MapAssign().ID("ucr").Dot("CreatedOn"),
				jen.ID("UpdatedOn").MapAssign().ID("ucr").Dot("UpdatedOn"),
				jen.ID("ArchivedOn").MapAssign().ID("ucr").Dot("ArchivedOn"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildURL").Params(jen.ID("address").String(), jen.ID("parts").Spread().String()).Params(jen.String()).Block(
			jen.List(jen.ID("tu"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.ID("address")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().Qual("net/url", "Parse").Call(jen.Qual("strings", "Join").Call(jen.ID("parts"), jen.Lit("/"))),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
			jen.Line(),
			jen.Return().ID("tu").Dot("ResolveReference").Call(jen.ID("u")).Dot("String").Call(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("getLoginCookie").Params(jen.ID("serviceURL").String(), jen.ID("u").PointerTo().Qual(proj.ModelsV1Package(),
			"User",
		)).Params(jen.ParamPointer().Qual("net/http", "Cookie"), jen.Error()).Block(
			jen.ID("uri").Assign().ID("buildURL").Call(jen.ID("serviceURL"), jen.Lit("users"), jen.Lit("login")),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.Qual("strings", "ToUpper").Call(jen.ID("u").Dot("TwoFactorSecret")), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("generating totp token: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodPost"), jen.ID("uri"), jen.Qual("strings", "NewReader").Callln(
					jen.Qual("fmt", "Sprintf").Callln(
						jen.Lit(`
					{
						"username": %q,
						"password": %q,
						"totp_token": %q
					}
				`),
						jen.ID("u").Dot("Username"),
						jen.ID("u").Dot("HashedPassword"),
						jen.ID("code"),
					),
				),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.Err()).Assign().Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.Err())),
			),
			jen.Line(),
			jen.If(jen.Err().Equals().ID("res").Dot("Body").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Qual("log", "Println").Call(jen.Lit("error closing body")),
			),
			jen.Line(),
			jen.ID("cookies").Assign().ID("res").Dot("Cookies").Call(),
			jen.If(jen.ID("len").Call(jen.ID("cookies")).Op(">").Zero()).Block(
				jen.Return().List(jen.ID("cookies").Index(jen.Zero()), jen.Nil()),
			),
			jen.Line(),
			jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("no cookie found :("))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateObligatoryClient creates the OAuth2 client we need for tests"),
		jen.Line(),
		jen.Func().ID("CreateObligatoryClient").Params(jen.ID("serviceURL").String(), jen.ID("u").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.ID("firstOAuth2ClientURI").Assign().ID("buildURL").Call(jen.ID("serviceURL"), jen.Lit("oauth2"), jen.Lit("client")),
			jen.Line(),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Callln(
				jen.Qual("strings", "ToUpper").Call(jen.ID("u").Dot("TwoFactorSecret")),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodPost"), jen.ID("firstOAuth2ClientURI"), jen.Qual("strings", "NewReader").Call(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit(`
	{
		"username": %q,
		"password": %q,
		"totp_token": %q,

		"belongs_to_user": %d,
		"scopes": ["*"]
	}
		`),
						jen.ID("u").Dot("Username"),
						jen.ID("u").Dot("HashedPassword"),
						jen.ID("code"),
						jen.ID("u").Dot("ID"),
					),
				),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("getLoginCookie").Call(jen.ID("serviceURL"), jen.ID("u")),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID("cookie").Op("==").ID("nil")).Block(
				jen.Qual("log", "Fatalf").Call(jen.Lit(`
cookie problems!
	cookie == nil: %v
			  err: %v
	`), jen.ID("cookie").Op("==").ID("nil"), jen.Err()),
			),
			jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
			jen.Var().ID("o").Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Line(),
			jen.Var().ID("command").Qual("fmt", "Stringer"),
			jen.If(jen.List(jen.ID("command"), jen.Err()).Equals().Qual("github.com/moul/http2curl", "GetCurlCommand").Call(jen.ID("req")), jen.Err().Op("==").ID("nil")).Block(
				jen.Qual("log", "Println").Call(jen.ID("command").Dot("String").Call()),
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.Err()).Assign().Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			).Else().If(jen.ID("res").Dot("StatusCode").DoesNotEqual().Qual("net/http", "StatusCreated")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("bad status: %d"), jen.ID("res").Dot("StatusCode"))),
			),
			jen.Line(),
			jen.Defer().Func().Params().Block(
				jen.If(jen.Err().Equals().ID("res").Dot("Body").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Qual("log", "Fatal").Call(jen.Err()),
				),
			).Call(),
			jen.Line(),
			jen.List(jen.ID("bdump"), jen.Err()).Assign().ID("httputil").Dot("DumpResponse").Call(jen.ID("res"), jen.True()),
			jen.If(jen.Err().Op("==").ID("nil").And().ID("req").Dot("Method").DoesNotEqual().Qual("net/http", "MethodGet")).Block(
				jen.Qual("log", "Println").Call(jen.String().Call(jen.ID("bdump"))),
			),
			jen.Line(),
			jen.Return().List(jen.AddressOf().ID("o"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.AddressOf().ID("o"))),
		),
		jen.Line(),
	)

	return ret
}
