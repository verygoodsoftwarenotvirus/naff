package data_scaffolder

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("uri").ID("string").Var().List(jen.ID("userCount"), jen.ID("dataCount")).ID("uint16"),
			jen.ID("debug").ID("bool"),
			jen.ID("singleUserMode").ID("bool"),
			jen.ID("singleUser").Op("*").ID("types").Dot("User"),
			jen.ID("quitter").Op("=").ID("fatalQuitter").Values(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("github.com/spf13/pflag", "StringVarP").Call(
				jen.Op("&").ID("uri"),
				jen.Lit("url"),
				jen.Lit("u"),
				jen.Lit(""),
				jen.Lit("where the target instance is hosted"),
			),
			jen.Qual("github.com/spf13/pflag", "Uint16VarP").Call(
				jen.Op("&").ID("userCount"),
				jen.Lit("user-count"),
				jen.Lit("c"),
				jen.Lit(0),
				jen.Lit("how many users to create"),
			),
			jen.Qual("github.com/spf13/pflag", "Uint16VarP").Call(
				jen.Op("&").ID("dataCount"),
				jen.Lit("data-count"),
				jen.Lit("d"),
				jen.Lit(0),
				jen.Lit("how many accounts/api clients/etc per user to create"),
			),
			jen.Qual("github.com/spf13/pflag", "BoolVarP").Call(
				jen.Op("&").ID("debug"),
				jen.Lit("debug"),
				jen.Lit("z"),
				jen.ID("false"),
				jen.Lit("whether debug mode is enabled"),
			),
			jen.Qual("github.com/spf13/pflag", "BoolVarP").Call(
				jen.Op("&").ID("singleUserMode"),
				jen.Lit("single-user-mode"),
				jen.Lit("s"),
				jen.ID("false"),
				jen.Lit("whether single user mode is enabled"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("clearTheScreen").Params().Body(
			jen.Qual("fmt", "Println").Call(jen.Lit("\x1b[2J")),
			jen.Qual("fmt", "Printf").Call(jen.Lit("\x1b[0;0H")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTOTPTokenForSecret").Params(jen.ID("secret").ID("string")).Params(jen.ID("string")).Body(
			jen.ID("secret").Op("=").Qual("strings", "ToUpper").Call(jen.ID("secret")),
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
				jen.ID("secret"),
				jen.Qual("time", "Now").Call().Dot("UTC").Call(),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.If(jen.Op("!").ID("totp").Dot("Validate").Call(
				jen.ID("code"),
				jen.ID("secret"),
			)).Body(
				jen.ID("panic").Call(jen.Lit("this shouldn't happen"))),
			jen.Return().ID("code"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("main").Params().Body(
			jen.Qual("github.com/spf13/pflag", "Parse").Call(),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("logger").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/logging/zerolog", "NewLogger").Call(),
			jen.If(jen.ID("debug")).Body(
				jen.ID("logger").Dot("SetLevel").Call(jen.ID("logging").Dot("DebugLevel"))),
			jen.If(jen.ID("dataCount").Op("<=").Lit(0)).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("exiting early because the requested amount is already satisfied")),
				jen.ID("quitter").Dot("Quit").Call(jen.Lit(0)),
			),
			jen.If(jen.ID("dataCount").Op("==").Lit(1).Op("&&").Op("!").ID("singleUserMode")).Body(
				jen.ID("singleUserMode").Op("=").ID("true")),
			jen.If(jen.ID("uri").Op("==").Lit("")).Body(
				jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Lit("uri must be valid"))),
			jen.List(jen.ID("parsedURI"), jen.ID("uriParseErr")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("uri")),
			jen.If(jen.ID("uriParseErr").Op("!=").ID("nil")).Body(
				jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("parsing provided url: %w"),
					jen.ID("uriParseErr"),
				))),
			jen.If(jen.ID("parsedURI").Dot("Scheme").Op("==").Lit("")).Body(
				jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Lit("provided URI missing scheme"))),
			jen.ID("wg").Op(":=").Op("&").Qual("sync", "WaitGroup").Values(),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("int").Call(jen.ID("userCount")), jen.ID("i").Op("++")).Body(
				jen.ID("wg").Dot("Add").Call(jen.Lit(1)),
				jen.Go().Func().Params(jen.ID("x").ID("int"), jen.ID("wg").Op("*").Qual("sync", "WaitGroup")).Body(
					jen.List(jen.ID("createdUser"), jen.ID("userCreationErr")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "CreateServiceUser").Call(
						jen.ID("ctx"),
						jen.ID("uri"),
						jen.Lit(""),
					),
					jen.If(jen.ID("userCreationErr").Op("!=").ID("nil")).Body(
						jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
							jen.Lit("creating user #%d: %w"),
							jen.ID("x"),
							jen.ID("userCreationErr"),
						))),
					jen.If(jen.ID("x").Op("==").Lit(0).Op("&&").ID("singleUserMode")).Body(
						jen.ID("singleUser").Op("=").ID("createdUser")),
					jen.ID("userLogger").Op(":=").ID("logger").Dot("WithValue").Call(
						jen.Lit("username"),
						jen.ID("createdUser").Dot("Username"),
					).Dot("WithValue").Call(
						jen.Lit("password"),
						jen.ID("createdUser").Dot("HashedPassword"),
					).Dot("WithValue").Call(
						jen.Lit("totp_secret"),
						jen.ID("createdUser").Dot("TwoFactorSecret"),
					).Dot("WithValue").Call(
						jen.Lit("user_id"),
						jen.ID("createdUser").Dot("ID"),
					).Dot("WithValue").Call(
						jen.Lit("user_number"),
						jen.ID("x"),
					),
					jen.ID("userLogger").Dot("Debug").Call(jen.Lit("created user")),
					jen.List(jen.ID("cookie"), jen.ID("cookieErr")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "GetLoginCookie").Call(
						jen.ID("ctx"),
						jen.ID("uri"),
						jen.ID("createdUser"),
					),
					jen.If(jen.ID("cookieErr").Op("!=").ID("nil")).Body(
						jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
							jen.Lit("getting cookie: %v"),
							jen.ID("cookieErr"),
						))),
					jen.List(jen.ID("userClient"), jen.ID("err")).Op(":=").ID("httpclient").Dot("NewClient").Call(
						jen.ID("parsedURI"),
						jen.ID("httpclient").Dot("UsingLogger").Call(jen.ID("userLogger")),
						jen.ID("httpclient").Dot("UsingCookie").Call(jen.ID("cookie")),
					),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
							jen.Lit("initializing client: %w"),
							jen.ID("err"),
						))),
					jen.ID("userLogger").Dot("Debug").Call(jen.Lit("assigned user API client")),
					jen.ID("wg").Dot("Add").Call(jen.Lit(1)),
					jen.Go().Func().Params(jen.ID("wg").Op("*").Qual("sync", "WaitGroup")).Body(
						jen.For(jen.ID("j").Op(":=").Lit(0), jen.ID("j").Op("<").ID("int").Call(jen.ID("dataCount")), jen.ID("j").Op("++")).Body(
							jen.ID("iterationLogger").Op(":=").ID("userLogger").Dot("WithValue").Call(
								jen.Lit("creating"),
								jen.Lit("accounts"),
							).Dot("WithValue").Call(
								jen.Lit("iteration"),
								jen.ID("j"),
							),
							jen.List(jen.ID("createdAccount"), jen.ID("accountCreationError")).Op(":=").ID("userClient").Dot("CreateAccount").Call(
								jen.ID("ctx"),
								jen.ID("fakes").Dot("BuildFakeAccountCreationInput").Call(),
							),
							jen.If(jen.ID("accountCreationError").Op("!=").ID("nil")).Body(
								jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
									jen.Lit("creating account #%d: %w"),
									jen.ID("j"),
									jen.ID("accountCreationError"),
								))),
							jen.ID("iterationLogger").Dot("WithValue").Call(
								jen.ID("keys").Dot("AccountIDKey"),
								jen.ID("createdAccount").Dot("ID"),
							).Dot("Debug").Call(jen.Lit("created account")),
						),
						jen.ID("wg").Dot("Done").Call(),
					).Call(jen.ID("wg")),
					jen.ID("wg").Dot("Add").Call(jen.Lit(1)),
					jen.Go().Func().Params(jen.ID("wg").Op("*").Qual("sync", "WaitGroup")).Body(
						jen.For(jen.ID("j").Op(":=").Lit(0), jen.ID("j").Op("<").ID("int").Call(jen.ID("dataCount")), jen.ID("j").Op("++")).Body(
							jen.ID("iterationLogger").Op(":=").ID("userLogger").Dot("WithValue").Call(
								jen.Lit("creating"),
								jen.Lit("api_clients"),
							).Dot("WithValue").Call(
								jen.Lit("iteration"),
								jen.ID("j"),
							),
							jen.List(jen.ID("code"), jen.ID("codeErr")).Op(":=").ID("totp").Dot("GenerateCode").Call(
								jen.Qual("strings", "ToUpper").Call(jen.ID("createdUser").Dot("TwoFactorSecret")),
								jen.Qual("time", "Now").Call().Dot("UTC").Call(),
							),
							jen.If(jen.ID("codeErr").Op("!=").ID("nil")).Body(
								jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
									jen.Lit("creating API Client #%d: %w"),
									jen.ID("j"),
									jen.ID("codeErr"),
								))),
							jen.ID("fakeInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInput").Call(),
							jen.List(jen.ID("createdAPIClient"), jen.ID("apiClientCreationErr")).Op(":=").ID("userClient").Dot("CreateAPIClient").Call(
								jen.ID("ctx"),
								jen.ID("cookie"),
								jen.Op("&").ID("types").Dot("APIClientCreationInput").Valuesln(
									jen.ID("UserLoginInput").Op(":").ID("types").Dot("UserLoginInput").Valuesln(
										jen.ID("Username").Op(":").ID("createdUser").Dot("Username"), jen.ID("Password").Op(":").ID("createdUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("code")), jen.ID("Name").Op(":").ID("fakeInput").Dot("Name")),
							),
							jen.If(jen.ID("apiClientCreationErr").Op("!=").ID("nil")).Body(
								jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
									jen.Lit("API Client webhook #%d: %w"),
									jen.ID("j"),
									jen.ID("apiClientCreationErr"),
								))),
							jen.ID("iterationLogger").Dot("WithValue").Call(
								jen.ID("keys").Dot("APIClientDatabaseIDKey"),
								jen.ID("createdAPIClient").Dot("ID"),
							).Dot("Debug").Call(jen.Lit("created API Client")),
						),
						jen.ID("wg").Dot("Done").Call(),
					).Call(jen.ID("wg")),
					jen.ID("wg").Dot("Add").Call(jen.Lit(1)),
					jen.Go().Func().Params(jen.ID("wg").Op("*").Qual("sync", "WaitGroup")).Body(
						jen.For(jen.ID("j").Op(":=").Lit(0), jen.ID("j").Op("<").ID("int").Call(jen.ID("dataCount")), jen.ID("j").Op("++")).Body(
							jen.ID("iterationLogger").Op(":=").ID("userLogger").Dot("WithValue").Call(
								jen.Lit("creating"),
								jen.Lit("webhooks"),
							).Dot("WithValue").Call(
								jen.Lit("iteration"),
								jen.ID("j"),
							),
							jen.List(jen.ID("createdWebhook"), jen.ID("webhookCreationErr")).Op(":=").ID("userClient").Dot("CreateWebhook").Call(
								jen.ID("ctx"),
								jen.ID("fakes").Dot("BuildFakeWebhookCreationInput").Call(),
							),
							jen.If(jen.ID("webhookCreationErr").Op("!=").ID("nil")).Body(
								jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
									jen.Lit("creating webhook #%d: %w"),
									jen.ID("j"),
									jen.ID("webhookCreationErr"),
								))),
							jen.ID("iterationLogger").Dot("WithValue").Call(
								jen.ID("keys").Dot("WebhookIDKey"),
								jen.ID("createdWebhook").Dot("ID"),
							).Dot("Debug").Call(jen.Lit("created webhook")),
						),
						jen.ID("wg").Dot("Done").Call(),
					).Call(jen.ID("wg")),
					jen.ID("wg").Dot("Add").Call(jen.Lit(1)),
					jen.Go().Func().Params(jen.ID("wg").Op("*").Qual("sync", "WaitGroup")).Body(
						jen.For(jen.ID("j").Op(":=").Lit(0), jen.ID("j").Op("<").ID("int").Call(jen.ID("dataCount")), jen.ID("j").Op("++")).Body(
							jen.ID("iterationLogger").Op(":=").ID("userLogger").Dot("WithValue").Call(
								jen.Lit("creating"),
								jen.Lit("items"),
							).Dot("WithValue").Call(
								jen.Lit("iteration"),
								jen.ID("j"),
							),
							jen.List(jen.ID("createdItem"), jen.ID("itemCreationErr")).Op(":=").ID("userClient").Dot("CreateItem").Call(
								jen.ID("ctx"),
								jen.ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
							),
							jen.If(jen.ID("itemCreationErr").Op("!=").ID("nil")).Body(
								jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(
									jen.Lit("creating item #%d: %w"),
									jen.ID("j"),
									jen.ID("itemCreationErr"),
								))),
							jen.ID("iterationLogger").Dot("WithValue").Call(
								jen.ID("keys").Dot("WebhookIDKey"),
								jen.ID("createdItem").Dot("ID"),
							).Dot("Debug").Call(jen.Lit("created item")),
						),
						jen.ID("wg").Dot("Done").Call(),
					).Call(jen.ID("wg")),
					jen.ID("wg").Dot("Done").Call(),
				).Call(
					jen.ID("i"),
					jen.ID("wg"),
				),
			),
			jen.ID("wg").Dot("Wait").Call(),
			jen.If(jen.ID("singleUserMode").Op("&&").ID("singleUser").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Debug").Call(jen.Lit("engage single user mode!")),
				jen.For(jen.Range().Qual("time", "Tick").Call(jen.Lit(1).Op("*").Qual("time", "Second"))).Body(
					jen.ID("clearTheScreen").Call(),
					jen.Qual("fmt", "Printf").Call(
						jen.Lit(`

username:  %s
passwords:  %s
2FA token: %s

`),
						jen.ID("singleUser").Dot("Username"),
						jen.ID("singleUser").Dot("HashedPassword"),
						jen.ID("buildTOTPTokenForSecret").Call(jen.ID("singleUser").Dot("TwoFactorSecret")),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
