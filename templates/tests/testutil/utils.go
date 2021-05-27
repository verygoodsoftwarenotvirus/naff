package testutil

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func utilsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("errEmptyAddressUnallowed").Op("=").Qual("errors", "New").Call(jen.Lit("empty address not allowed")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateServiceUser creates a user."),
		jen.Line(),
		jen.Func().ID("CreateServiceUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("address"), jen.ID("username")).ID("string")).Params(jen.Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.If(jen.ID("username").Op("==").Lit("")).Body(
				jen.ID("username").Op("=").ID("gofakeit").Dot("Password").Call(
					jen.ID("true"),
					jen.ID("true"),
					jen.ID("true"),
					jen.ID("false"),
					jen.ID("false"),
					jen.Lit(32),
				)),
			jen.If(jen.ID("address").Op("==").Lit("")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errEmptyAddressUnallowed"))),
			jen.List(jen.ID("parsedAddress"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("address")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("httpclient").Dot("NewClient").Call(jen.ID("parsedAddress")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("initializing client: %w"),
					jen.ID("err"),
				))),
			jen.ID("in").Op(":=").Op("&").ID("types").Dot("UserRegistrationInput").Valuesln(jen.ID("Username").Op(":").ID("username"), jen.ID("Password").Op(":").ID("gofakeit").Dot("Password").Call(
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.ID("true"),
				jen.Lit(64),
			)),
			jen.List(jen.ID("ucr"), jen.ID("err")).Op(":=").ID("c").Dot("CreateUser").Call(
				jen.ID("ctx"),
				jen.ID("in"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.List(jen.ID("token"), jen.ID("tokenErr")).Op(":=").ID("totp").Dot("GenerateCode").Call(
				jen.ID("ucr").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.ID("tokenErr").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("generating totp code: %w"),
					jen.ID("tokenErr"),
				))),
			jen.If(jen.ID("validationErr").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
				jen.ID("ctx"),
				jen.ID("ucr").Dot("CreatedUserID"),
				jen.ID("token"),
			), jen.ID("validationErr").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("verifying totp code: %w"),
					jen.ID("validationErr"),
				))),
			jen.ID("u").Op(":=").Op("&").ID("types").Dot("User").Valuesln(jen.ID("ID").Op(":").ID("ucr").Dot("CreatedUserID"), jen.ID("Username").Op(":").ID("ucr").Dot("Username"), jen.ID("HashedPassword").Op(":").ID("in").Dot("Password"), jen.ID("TwoFactorSecret").Op(":").ID("ucr").Dot("TwoFactorSecret"), jen.ID("CreatedOn").Op(":").ID("ucr").Dot("CreatedOn")),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetLoginCookie fetches a login cookie for a given user."),
		jen.Line(),
		jen.Func().ID("GetLoginCookie").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("serviceURL").ID("string"), jen.ID("u").Op("*").ID("types").Dot("User")).Params(jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")).Body(
			jen.List(jen.ID("tu"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("serviceURL")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.List(jen.ID("lu"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Qual("strings", "Join").Call(
				jen.Index().ID("string").Valuesln(jen.Lit("users"), jen.Lit("login")),
				jen.Lit("/"),
			)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.ID("uri").Op(":=").ID("tu").Dot("ResolveReference").Call(jen.ID("lu")).Dot("String").Call(),
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
				jen.Qual("strings", "ToUpper").Call(jen.ID("u").Dot("TwoFactorSecret")),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("generating totp token: %w"),
					jen.ID("err"),
				))),
			jen.List(jen.ID("body"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("u").Dot("Username"), jen.ID("Password").Op(":").ID("u").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("code"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("generating login request body: %w"),
					jen.ID("err"),
				))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.Qual("bytes", "NewReader").Call(jen.ID("body")),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				))),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("executing request: %w"),
					jen.ID("err"),
				))),
			jen.If(jen.ID("err").Op("=").ID("res").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual("log", "Println").Call(jen.Lit("error closing body"))),
			jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
			jen.If(jen.ID("len").Call(jen.ID("cookies")).Op(">").Lit(0)).Body(
				jen.Return().List(jen.ID("cookies").Index(jen.Lit(0)), jen.ID("nil"))),
			jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "ErrNoCookie")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("DetermineServiceURL returns the url, if properly configured."),
		jen.Line(),
		jen.Func().ID("DetermineServiceURL").Params().Params(jen.Op("*").Qual("net/url", "URL")).Body(
			jen.ID("ta").Op(":=").Qual("os", "Getenv").Call(jen.Lit("TARGET_ADDRESS")),
			jen.If(jen.ID("ta").Op("==").Lit("")).Body(
				jen.ID("panic").Call(jen.Lit("must provide target address!"))),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("ta")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("u"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EnsureServerIsUp checks that a server is up and doesn't return until it's certain one way or the other."),
		jen.Line(),
		jen.Func().ID("EnsureServerIsUp").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("address").ID("string")).Body(
			jen.Var().Defs(
				jen.ID("isDown").Op("=").ID("true"),
				jen.ID("interval").Op("=").Qual("time", "Second"),
				jen.ID("maxAttempts").Op("=").Lit(50),
				jen.ID("numberOfAttempts").Op("=").Lit(0),
			),
			jen.For(jen.ID("isDown")).Body(
				jen.If(jen.Op("!").ID("IsUp").Call(
					jen.ID("ctx"),
					jen.ID("address"),
				)).Body(
					jen.Qual("log", "Printf").Call(
						jen.Lit("waiting %s before pinging %q again"),
						jen.ID("interval"),
						jen.ID("address"),
					),
					jen.Qual("time", "Sleep").Call(jen.ID("interval")),
					jen.ID("numberOfAttempts").Op("++"),
					jen.If(jen.ID("numberOfAttempts").Op(">=").ID("maxAttempts")).Body(
						jen.Qual("log", "Fatal").Call(jen.Lit("Maximum number of attempts made, something's gone awry"))),
				).Else().Body(
					jen.ID("isDown").Op("=").ID("false"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IsUp can check if an instance of our server is alive."),
		jen.Line(),
		jen.Func().ID("IsUp").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("address").ID("string")).Params(jen.ID("bool")).Body(
			jen.ID("uri").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s/_meta_/ready"),
				jen.ID("address"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("false")),
			jen.If(jen.ID("err").Op("=").ID("res").Dot("Body").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual("log", "Println").Call(jen.Lit("error closing body"))),
			jen.Return().ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusOK"),
		),
		jen.Line(),
	)

	return code
}
