package testutil

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func testutilDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("testutil")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.Qual(utils.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("DetermineServiceURL returns the URL, if properly configured"),
		jen.Line(),
		jen.Func().ID("DetermineServiceURL").Params().Params(jen.ID("string")).Block(
			jen.ID("ta").Op(":=").Qual("os", "Getenv").Call(jen.Lit("TARGET_ADDRESS")),
			jen.If(jen.ID("ta").Op("==").Lit("")).Block(
				jen.ID("panic").Call(jen.Lit("must provide target address!")),
			),
			jen.Line(),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("ta")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("EnsureServerIsUp checks that a server is up and doesn't return until it's certain one way or the other"),
		jen.Line(),
		jen.Func().ID("EnsureServerIsUp").Params(jen.ID("address").ID("string")).Block(
			jen.Var().Defs(
				jen.ID("isDown").Op("=").ID("true"),
				jen.ID("interval").Op("=").Qual("time", "Second"),
				jen.ID("maxAttempts").Op("=").Lit(50),
				jen.ID("numberOfAttempts").Op("=").Lit(0),
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
					jen.ID("isDown").Op("=").ID("false"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsUp can check if an instance of our server is alive"),
		jen.Line(),
		jen.Func().ID("IsUp").Params(jen.ID("address").ID("string")).Params(jen.ID("bool")).Block(
			jen.ID("uri").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s/_meta_/ready"), jen.ID("address")),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.ID("uri"), jen.ID("nil")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().ID("false"),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("res").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
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
		jen.Func().ID("CreateObligatoryUser").Params(jen.ID("address").ID("string"), jen.ID("debug").ID("bool")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"),
			"User",
		), jen.ID("error")).Block(
			jen.List(jen.ID("tu"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("address")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").Qual(filepath.Join(pkgRoot, "client/v1/http"), "NewSimpleClient").Call(jen.Qual("context", "Background").Call(), jen.ID("tu"), jen.ID("debug")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.Comment("I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,"),
			jen.Comment("could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit"),
			jen.ID("username").Op(":=").Qual(utils.FakeLibrary, "Username").Call().Op("+").Qual(utils.FakeLibrary, "HexColor").Call().Op("+").Qual(utils.FakeLibrary, "Country").Call(),
			jen.ID("in").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "UserInput").Valuesln(
				jen.ID("Username").Op(":").ID("username"),
				jen.ID("Password").Op(":").Qual(utils.FakeLibrary, "Password").Call(jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.ID("true"), jen.Lit(64)),
			),
			jen.Line(),
			jen.List(jen.ID("ucr"), jen.ID("err")).Op(":=").ID("c").Dot("CreateUser").Call(jen.Qual("context", "Background").Call(), jen.ID("in")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			).Else().If(jen.ID("ucr").Op("==").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("errors").Dot("New").Call(jen.Lit("something happened"))),
			),
			jen.Line(),
			jen.ID("u").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), "User").Valuesln(
				jen.ID("ID").Op(":").ID("ucr").Dot("ID"),
				jen.ID("Username").Op(":").ID("ucr").Dot("Username"),
				jen.Comment("this is a dirty trick to reuse most of this model"),
				jen.ID("HashedPassword").Op(":").ID("in").Dot("Password"),
				jen.ID("TwoFactorSecret").Op(":").ID("ucr").Dot("TwoFactorSecret"),
				jen.ID("PasswordLastChangedOn").Op(":").ID("ucr").Dot("PasswordLastChangedOn"),
				jen.ID("CreatedOn").Op(":").ID("ucr").Dot("CreatedOn"),
				jen.ID("UpdatedOn").Op(":").ID("ucr").Dot("UpdatedOn"),
				jen.ID("ArchivedOn").Op(":").ID("ucr").Dot("ArchivedOn"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildURL").Params(jen.ID("address").ID("string"), jen.ID("parts").Op("...").ID("string")).Params(jen.ID("string")).Block(
			jen.List(jen.ID("tu"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("address")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Qual("strings", "Join").Call(jen.ID("parts"), jen.Lit("/"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Return().ID("tu").Dot("ResolveReference").Call(jen.ID("u")).Dot("String").Call(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("getLoginCookie").Params(jen.ID("serviceURL").ID("string"), jen.ID("u").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"),
			"User",
		)).Params(jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")).Block(
			jen.ID("uri").Op(":=").ID("buildURL").Call(jen.ID("serviceURL"), jen.Lit("users"), jen.Lit("login")),
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(jen.Qual("strings", "ToUpper").Call(jen.ID("u").Dot("TwoFactorSecret")), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("generating totp token: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Callln(
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
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("res").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Qual("log", "Println").Call(jen.Lit("error closing body")),
			),
			jen.Line(),
			jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
			jen.If(jen.ID("len").Call(jen.ID("cookies")).Op(">").Lit(0)).Block(
				jen.Return().List(jen.ID("cookies").Index(jen.Lit(0)), jen.ID("nil")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("nil"), jen.ID("errors").Dot("New").Call(jen.Lit("no cookie found :("))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateObligatoryClient creates the OAuth2 client we need for tests"),
		jen.Line(),
		jen.Func().ID("CreateObligatoryClient").Params(jen.ID("serviceURL").ID("string"), jen.ID("u").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"), jen.ID("error")).Block(
			jen.ID("firstOAuth2ClientURI").Op(":=").ID("buildURL").Call(jen.ID("serviceURL"), jen.Lit("oauth2"), jen.Lit("client")),
			jen.Line(),
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Callln(
				jen.Qual("strings", "ToUpper").Call(jen.ID("u").Dot("TwoFactorSecret")),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodPost"), jen.ID("firstOAuth2ClientURI"), jen.Qual("strings", "NewReader").Call(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit(`
	{
		"username": %q,
		"password": %q,
		"totp_token": %q,

		"belongs_to": %d,
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
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("getLoginCookie").Call(jen.ID("serviceURL"), jen.ID("u")),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("cookie").Op("==").ID("nil")).Block(
				jen.Qual("log", "Fatalf").Call(jen.Lit(`
cookie problems!
	cookie == nil: %v
			  err: %v
	`), jen.ID("cookie").Op("==").ID("nil"), jen.ID("err")),
			),
			jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
			jen.Var().ID("o").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
			jen.Line(),
			jen.Var().ID("command").Qual("fmt", "Stringer"),
			jen.If(jen.List(jen.ID("command"), jen.ID("err")).Op("=").ID("http2curl").Dot("GetCurlCommand").Call(jen.ID("req")), jen.ID("err").Op("==").ID("nil")).Block(
				jen.Qual("log", "Println").Call(jen.ID("command").Dot("String").Call()),
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			).Else().If(jen.ID("res").Dot("StatusCode").Op("!=").Qual("net/http", "StatusCreated")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("bad status: %d"), jen.ID("res").Dot("StatusCode"))),
			),
			jen.Line(),
			jen.Defer().Func().Params().Block(
				jen.If(jen.ID("err").Op("=").ID("res").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Qual("log", "Fatal").Call(jen.ID("err")),
				),
			).Call(),
			jen.Line(),
			jen.List(jen.ID("bdump"), jen.ID("err")).Op(":=").ID("httputil").Dot("DumpResponse").Call(jen.ID("res"), jen.ID("true")),
			jen.If(jen.ID("err").Op("==").ID("nil").Op("&&").ID("req").Dot("Method").Op("!=").Qual("net/http", "MethodGet")).Block(
				jen.Qual("log", "Println").Call(jen.ID("string").Call(jen.ID("bdump"))),
			),
			jen.Line(),
			jen.Return().List(jen.Op("&").ID("o"), jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.Op("&").ID("o"))),
		),
		jen.Line(),
	)
	return ret
}
