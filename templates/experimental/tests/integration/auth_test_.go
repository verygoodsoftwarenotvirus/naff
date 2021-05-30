package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestLogin").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("logging in and out works"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
						jen.Qual("context", "Background").Call(),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("testUser"), jen.ID("_"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("testUser").Dot("Username"), jen.ID("Password").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("testUser"),
						)),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/authentication", "DefaultCookieName"),
						jen.ID("cookie").Dot("Name"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("cookie").Dot("Value"),
					),
					jen.ID("assert").Dot("NotZero").Call(
						jen.ID("t"),
						jen.ID("cookie").Dot("MaxAge"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("cookie").Dot("HttpOnly"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Lit("/"),
						jen.ID("cookie").Dot("Path"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "SameSiteStrictMode"),
						jen.ID("cookie").Dot("SameSite"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("EndSession").Call(jen.ID("ctx")),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestLogin_WithoutBodyFails").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("login request without body fails"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("testClient").Dot("BuildURL").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
					)),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("u").Dot("Path").Op("=").Lit("/users/login"),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.ID("u").Dot("String").Call(),
						jen.ID("nil"),
					),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("req"),
						jen.ID("err"),
					),
					jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("testClient").Dot("PlainClient").Call().Dot("Do").Call(jen.ID("req")),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("res"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("res").Dot("StatusCode"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestLogin_ShouldNotBeAbleToLoginWithInvalidPassword").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("should not be able to log in with the wrong password"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
						jen.Qual("context", "Background").Call(),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("testUser"), jen.ID("_"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.Var().Defs(
						jen.ID("badPassword").ID("string"),
					),
					jen.For(jen.List(jen.ID("_"), jen.ID("v")).Op(":=").Range().ID("testUser").Dot("HashedPassword")).Body(
						jen.ID("badPassword").Op("=").ID("string").Call(jen.ID("v")).Op("+").ID("badPassword")),
					jen.ID("r").Op(":=").Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("testUser").Dot("Username"), jen.ID("Password").Op(":").ID("badPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
						jen.ID("t"),
						jen.ID("testUser"),
					)),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.ID("r"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestLogin_ShouldNotBeAbleToLoginAsANonexistentUser").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("should not be able to login as someone that does not exist"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
						jen.Qual("context", "Background").Call(),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("testUser"), jen.ID("_"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.ID("exampleUserCreationInput").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
					jen.ID("r").Op(":=").Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("exampleUserCreationInput").Dot("Username"), jen.ID("Password").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").Lit("123456")),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.ID("r"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestLogin_ShouldNotBeAbleToLoginWithoutValidating2FASecret").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("should not be able to login without validating 2FA secret"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.ID("testClient").Op(":=").ID("buildSimpleClient").Call(jen.ID("t")),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleUserCreationInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("exampleUser")),
					jen.List(jen.ID("ucr"), jen.ID("err")).Op(":=").ID("testClient").Dot("CreateUser").Call(
						jen.ID("ctx"),
						jen.ID("exampleUserCreationInput"),
					),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("ucr"),
						jen.ID("err"),
					),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
						jen.ID("ucr").Dot("TwoFactorSecret"),
						jen.Qual("time", "Now").Call().Dot("UTC").Call(),
					),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("token"),
						jen.ID("err"),
					),
					jen.ID("r").Op(":=").Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("exampleUserCreationInput").Dot("Username"), jen.ID("Password").Op(":").ID("exampleUserCreationInput").Dot("Password"), jen.ID("TOTPToken").Op(":").ID("token")),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.ID("r"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestCheckingAuthStatus").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("checking auth status"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
						jen.Qual("context", "Background").Call(),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("testUser"), jen.ID("_"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("testUser").Dot("Username"), jen.ID("Password").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("testUser"),
						)),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClient").Dot("UserStatus").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("true"),
						jen.ID("actual").Dot("UserIsAuthenticated"),
						jen.Lit("expected UserIsAuthenticated to equal %v, but got %v"),
						jen.ID("true"),
						jen.ID("actual").Dot("UserIsAuthenticated"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("types").Dot("GoodStandingAccountStatus"),
						jen.ID("actual").Dot("UserReputation"),
						jen.Lit("expected UserReputation to equal %v, but got %v"),
						jen.ID("types").Dot("GoodStandingAccountStatus"),
						jen.ID("actual").Dot("UserReputation"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Lit(""),
						jen.ID("actual").Dot("UserReputationExplanation"),
						jen.Lit("expected UserReputationExplanation to equal %v, but got %v"),
						jen.Lit(""),
						jen.ID("actual").Dot("UserReputationExplanation"),
					),
					jen.ID("assert").Dot("NotZero").Call(
						jen.ID("t"),
						jen.ID("actual").Dot("ActiveAccount"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("EndSession").Call(jen.ID("ctx")),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestPASETOGeneration").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("checking auth status"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
						jen.Qual("context", "Background").Call(),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("user"), jen.ID("cookie"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
					jen.ID("exampleAPIClientInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
					jen.ID("exampleAPIClientInput").Dot("UserLoginInput").Op("=").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("user").Dot("Username"), jen.ID("Password").Op(":").ID("user").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
						jen.ID("t"),
						jen.ID("user"),
					)),
					jen.List(jen.ID("createdAPIClient"), jen.ID("apiClientCreationErr")).Op(":=").ID("testClient").Dot("CreateAPIClient").Call(
						jen.ID("ctx"),
						jen.ID("cookie"),
						jen.ID("exampleAPIClientInput"),
					),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("createdAPIClient"),
						jen.ID("apiClientCreationErr"),
					),
					jen.List(jen.ID("actualKey"), jen.ID("keyDecodeErr")).Op(":=").Qual("encoding/base64", "RawURLEncoding").Dot("DecodeString").Call(jen.ID("createdAPIClient").Dot("ClientSecret")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("keyDecodeErr"),
					),
					jen.ID("input").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").ID("createdAPIClient").Dot("ClientID"), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("UTC").Call().Dot("UnixNano").Call()),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("testClient").Dot("RequestBuilder").Call().Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("ctx"),
						jen.ID("input"),
						jen.ID("actualKey"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("res"), jen.ID("err")).Op(":=").Qual("net/http", "DefaultClient").Dot("Do").Call(jen.ID("req")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("res"),
					),
					jen.Var().Defs(
						jen.ID("tokenRes").ID("types").Dot("PASETOResponse"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("res").Dot("Body")).Dot("Decode").Call(jen.Op("&").ID("tokenRes")),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("tokenRes").Dot("Token"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("tokenRes").Dot("ExpiresAt"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestPasswordChanging").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("should be possible to change your password"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("testUser"), jen.ID("_"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("testUser").Dot("Username"), jen.ID("Password").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("testUser"),
						)),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Var().Defs(
						jen.ID("backwardsPass").ID("string"),
					),
					jen.For(jen.List(jen.ID("_"), jen.ID("v")).Op(":=").Range().ID("testUser").Dot("HashedPassword")).Body(
						jen.ID("backwardsPass").Op("=").ID("string").Call(jen.ID("v")).Op("+").ID("backwardsPass")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("ChangePassword").Call(
							jen.ID("ctx"),
							jen.ID("cookie"),
							jen.Op("&").ID("types").Dot("PasswordUpdateInput").Valuesln(jen.ID("CurrentPassword").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
								jen.ID("t"),
								jen.ID("testUser"),
							), jen.ID("NewPassword").Op(":").ID("backwardsPass")),
						),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("EndSession").Call(jen.ID("ctx")),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op("=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("testUser").Dot("Username"), jen.ID("Password").Op(":").ID("backwardsPass"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("testUser"),
						)),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestTOTPSecretChanging").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("should be possible to change your TOTP secret"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
						jen.Qual("context", "Background").Call(),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.List(jen.ID("testUser"), jen.ID("_"), jen.ID("testClient"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("testUser").Dot("Username"), jen.ID("Password").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("testUser"),
						)),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("r"), jen.ID("err")).Op(":=").ID("testClient").Dot("CycleTwoFactorSecret").Call(
						jen.ID("ctx"),
						jen.ID("cookie"),
						jen.Op("&").ID("types").Dot("TOTPSecretRefreshInput").Valuesln(jen.ID("CurrentPassword").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("testUser"),
						)),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("secretVerificationToken"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
						jen.ID("r").Dot("TwoFactorSecret"),
						jen.Qual("time", "Now").Call().Dot("UTC").Call(),
					),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("secretVerificationToken"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("VerifyTOTPSecret").Call(
							jen.ID("ctx"),
							jen.ID("testUser").Dot("ID"),
							jen.ID("secretVerificationToken"),
						),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("EndSession").Call(jen.ID("ctx")),
					),
					jen.List(jen.ID("newToken"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
						jen.ID("r").Dot("TwoFactorSecret"),
						jen.Qual("time", "Now").Call().Dot("UTC").Call(),
					),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("newToken"),
						jen.ID("err"),
					),
					jen.List(jen.ID("secondCookie"), jen.ID("err")).Op(":=").ID("testClient").Dot("BeginSession").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("testUser").Dot("Username"), jen.ID("Password").Op(":").ID("testUser").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("newToken")),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("secondCookie"),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestTOTPTokenValidation").Params().Body(
			jen.ID("s").Dot("Run").Call(
				jen.Lit("should be possible to validate TOTP"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.ID("testClient").Op(":=").ID("buildSimpleClient").Call(jen.ID("t")),
					jen.ID("userInput").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
					jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("testClient").Dot("CreateUser").Call(
						jen.ID("ctx"),
						jen.ID("userInput"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("user"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
						jen.ID("user").Dot("TwoFactorSecret"),
						jen.Qual("time", "Now").Call().Dot("UTC").Call(),
					),
					jen.ID("requireNotNilAndNoProblems").Call(
						jen.ID("t"),
						jen.ID("token"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("VerifyTOTPSecret").Call(
							jen.ID("ctx"),
							jen.ID("user").Dot("CreatedUserID"),
							jen.ID("token"),
						),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("should not be possible to validate an invalid TOTP"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
						jen.Qual("context", "Background").Call(),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.ID("testClient").Op(":=").ID("buildSimpleClient").Call(jen.ID("t")),
					jen.ID("userInput").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
					jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("testClient").Dot("CreateUser").Call(
						jen.ID("ctx"),
						jen.ID("userInput"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("user"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("testClient").Dot("VerifyTOTPSecret").Call(
							jen.ID("ctx"),
							jen.ID("user").Dot("CreatedUserID"),
							jen.Lit("NOTREAL"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
