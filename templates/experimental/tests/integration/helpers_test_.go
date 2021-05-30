package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("reverseString").Params(jen.ID("s").ID("string")).Params(jen.ID("string")).Body(
			jen.ID("runes").Op(":=").Index().ID("rune").Call(jen.ID("s")),
			jen.For(jen.List(jen.ID("i"), jen.ID("j")).Op(":=").List(jen.Lit(0), jen.ID("len").Call(jen.ID("runes")).Op("-").Lit(1)), jen.ID("i").Op("<").ID("j"), jen.List(jen.ID("i"), jen.ID("j")).Op("=").List(jen.ID("i").Op("+").Lit(1), jen.ID("j").Op("-").Lit(1))).Body(
				jen.List(jen.ID("runes").Index(jen.ID("i")), jen.ID("runes").Index(jen.ID("j"))).Op("=").List(jen.ID("runes").Index(jen.ID("j")), jen.ID("runes").Index(jen.ID("i")))),
			jen.Return().ID("string").Call(jen.ID("runes")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("requireNotNilAndNoProblems").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("i").Interface(), jen.ID("err").ID("error")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("i"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("createUserAndClientForTest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.ID("user").Op("*").ID("types").Dot("User"), jen.ID("cookie").Op("*").Qual("net/http", "Cookie"), jen.List(jen.ID("cookieClient"), jen.ID("pasetoClient")).Op("*").ID("httpclient").Dot("Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "CreateServiceUser").Call(
				jen.ID("ctx"),
				jen.ID("urlToUse"),
				jen.ID("fakes").Dot("BuildFakeUser").Call().Dot("Username"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("t").Dot("Logf").Call(
				jen.Lit("created user #%d: %q"),
				jen.ID("user").Dot("ID"),
				jen.ID("user").Dot("Username"),
			),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "GetLoginCookie").Call(
				jen.ID("ctx"),
				jen.ID("urlToUse"),
				jen.ID("user"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("cookieClient"), jen.ID("err")).Op("=").ID("initializeCookiePoweredClient").Call(jen.ID("cookie")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("apiClient"), jen.ID("err")).Op(":=").ID("cookieClient").Dot("CreateAPIClient").Call(
				jen.ID("ctx"),
				jen.ID("cookie"),
				jen.Op("&").ID("types").Dot("APIClientCreationInput").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call(), jen.ID("UserLoginInput").Op(":").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("user").Dot("Username"), jen.ID("Password").Op(":").ID("user").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
					jen.ID("t"),
					jen.ID("user"),
				))),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("secretKey"), jen.ID("err")).Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("DecodeString").Call(jen.ID("apiClient").Dot("ClientSecret")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("pasetoClient"), jen.ID("err")).Op("=").ID("initializePASETOPoweredClient").Call(
				jen.ID("apiClient").Dot("ClientID"),
				jen.ID("secretKey"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().List(jen.ID("user"), jen.ID("cookie"), jen.ID("cookieClient"), jen.ID("pasetoClient")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("initializeCookiePoweredClient").Params(jen.ID("cookie").Op("*").Qual("net/http", "Cookie")).Params(jen.Op("*").ID("httpclient").Dot("Client"), jen.ID("error")).Body(
			jen.If(jen.ID("parsedURLToUse").Op("==").ID("nil")).Body(
				jen.ID("panic").Call(jen.Lit("url not set!"))),
			jen.ID("logger").Op(":=").ID("logging").Dot("ProvideLogger").Call(jen.ID("logging").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("logging").Dot("ProviderZerolog"))),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("httpclient").Dot("NewClient").Call(
				jen.ID("parsedURLToUse"),
				jen.ID("httpclient").Dot("UsingLogger").Call(jen.ID("logger")),
				jen.ID("httpclient").Dot("UsingCookie").Call(jen.ID("cookie")),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.If(jen.ID("debug")).Body(
				jen.If(jen.ID("setOptionErr").Op(":=").ID("c").Dot("SetOptions").Call(jen.ID("httpclient").Dot("UsingDebug").Call(jen.ID("true"))), jen.ID("setOptionErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("setOptionErr")))),
			jen.Return().List(jen.ID("c"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("initializePASETOPoweredClient").Params(jen.ID("clientID").ID("string"), jen.ID("secretKey").Index().ID("byte")).Params(jen.Op("*").ID("httpclient").Dot("Client"), jen.ID("error")).Body(
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("httpclient").Dot("NewClient").Call(
				jen.ID("parsedURLToUse"),
				jen.ID("httpclient").Dot("UsingLogger").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
				jen.ID("httpclient").Dot("UsingPASETO").Call(
					jen.ID("clientID"),
					jen.ID("secretKey"),
				),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.If(jen.ID("debug")).Body(
				jen.If(jen.ID("setOptionErr").Op(":=").ID("c").Dot("SetOptions").Call(jen.ID("httpclient").Dot("UsingDebug").Call(jen.ID("true"))), jen.ID("setOptionErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("setOptionErr")))),
			jen.Return().List(jen.ID("c"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildSimpleClient").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("httpclient").Dot("Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("httpclient").Dot("NewClient").Call(jen.ID("parsedURLToUse")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("c"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("generateTOTPTokenForUser").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("u").Op("*").ID("types").Dot("User")).Params(jen.ID("string")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
				jen.ID("u").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.ID("require").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("code"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("code"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildAdminCookieAndPASETOClients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.List(jen.ID("cookieClient"), jen.ID("pasetoClient")).Op("*").ID("httpclient").Dot("Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("u").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "DetermineServiceURL").Call(),
			jen.ID("urlToUse").Op("=").ID("u").Dot("String").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("ProvideLogger").Call(jen.ID("logging").Dot("Config").Valuesln(jen.ID("Provider").Op(":").ID("logging").Dot("ProviderZerolog"))),
			jen.ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("URLKey"),
				jen.ID("urlToUse"),
			).Dot("Info").Call(jen.Lit("checking server")),
			jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "EnsureServerIsUp").Call(
				jen.ID("ctx"),
				jen.ID("urlToUse"),
			),
			jen.List(jen.ID("adminCookie"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "GetLoginCookie").Call(
				jen.ID("ctx"),
				jen.ID("urlToUse"),
				jen.ID("premadeAdminUser"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("cClient"), jen.ID("err")).Op(":=").ID("initializeCookiePoweredClient").Call(jen.ID("adminCookie")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
				jen.ID("premadeAdminUser").Dot("TwoFactorSecret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("apiClient"), jen.ID("err")).Op(":=").ID("cClient").Dot("CreateAPIClient").Call(
				jen.ID("ctx"),
				jen.ID("adminCookie"),
				jen.Op("&").ID("types").Dot("APIClientCreationInput").Valuesln(jen.ID("Name").Op(":").Qual("fmt", "Sprintf").Call(
					jen.Lit("admin_paseto_client_%d"),
					jen.Qual("time", "Now").Call().Dot("UnixNano").Call(),
				), jen.ID("UserLoginInput").Op(":").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("premadeAdminUser").Dot("Username"), jen.ID("Password").Op(":").ID("premadeAdminUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("code"))),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("secretKey"), jen.ID("err")).Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("DecodeString").Call(jen.ID("apiClient").Dot("ClientSecret")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("PASETOClient"), jen.ID("err")).Op(":=").ID("initializePASETOPoweredClient").Call(
				jen.ID("apiClient").Dot("ClientID"),
				jen.ID("secretKey"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().List(jen.ID("cClient"), jen.ID("PASETOClient")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("validateAuditLogEntries").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expectedEntries"), jen.ID("actualEntries")).Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("relevantID").ID("uint64"), jen.ID("key").ID("string")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("expectedEventTypes").Op(":=").Index().ID("string").Valuesln(),
			jen.ID("actualEventTypes").Op(":=").Index().ID("string").Valuesln(),
			jen.For(jen.List(jen.ID("_"), jen.ID("e")).Op(":=").Range().ID("expectedEntries")).Body(
				jen.ID("expectedEventTypes").Op("=").ID("append").Call(
					jen.ID("expectedEventTypes"),
					jen.ID("e").Dot("EventType"),
				)),
			jen.For(jen.List(jen.ID("_"), jen.ID("e")).Op(":=").Range().ID("actualEntries")).Body(
				jen.ID("actualEventTypes").Op("=").ID("append").Call(
					jen.ID("actualEventTypes"),
					jen.ID("e").Dot("EventType"),
				),
				jen.If(jen.ID("relevantID").Op("!=").Lit(0).Op("&&").ID("key").Op("!=").Lit("")).Body(
					jen.If(jen.ID("assert").Dot("Contains").Call(
						jen.ID("t"),
						jen.ID("e").Dot("Context"),
						jen.ID("key"),
					)).Body(
						jen.ID("assert").Dot("EqualValues").Call(
							jen.ID("t"),
							jen.ID("relevantID"),
							jen.ID("e").Dot("Context").Index(jen.ID("key")),
						))),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("len").Call(jen.ID("expectedEntries")),
				jen.ID("len").Call(jen.ID("actualEntries")),
				jen.Lit("expected %q, got %q"),
				jen.Qual("strings", "Join").Call(
					jen.ID("expectedEventTypes"),
					jen.Lit(","),
				),
				jen.Qual("strings", "Join").Call(
					jen.ID("actualEventTypes"),
					jen.Lit(","),
				),
			),
			jen.ID("assert").Dot("Subset").Call(
				jen.ID("t"),
				jen.ID("expectedEventTypes"),
				jen.ID("actualEventTypes"),
			),
		),
		jen.Line(),
	)

	return code
}
