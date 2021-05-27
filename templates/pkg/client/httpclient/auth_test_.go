package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAuth").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.ID("authTestSuite")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("authTestSuite").Struct(
				jen.ID("suite").Dot("Suite"),
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
				jen.ID("exampleCookie").Op("*").Qual("net/http", "Cookie"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").ID("authTestSuite")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("authTestSuite")).ID("SetupTest").Params().Body(
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dot("exampleCookie").Op("=").Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").ID("s").Dot("T").Call().Dot("Name").Call()),
			jen.ID("s").Dot("exampleUser").Op("=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("s").Dot("exampleUser").Dot("HashedPassword").Op("=").Lit(""),
			jen.ID("s").Dot("exampleUser").Dot("TwoFactorSecret").Op("=").Lit(""),
			jen.ID("s").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("authTestSuite")).ID("TestClient_UserStatus").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/auth/status"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("expected").Op(":=").Op("&").ID("types").Dot("UserStatusResponse").Valuesln(jen.ID("UserReputation").Op(":").ID("s").Dot("exampleUser").Dot("ServiceAccountStatus"), jen.ID("UserReputationExplanation").Op(":").ID("s").Dot("exampleUser").Dot("ReputationExplanation"), jen.ID("UserIsAuthenticated").Op(":").ID("true")),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("expected"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("UserStatus").Call(jen.ID("s").Dot("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("UserStatus").Call(jen.ID("s").Dot("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("UserStatus").Call(jen.ID("s").Dot("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("authTestSuite")).ID("TestClient_Login").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/login"),
			),
			jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
				jen.ID("false"),
				jen.Qual("net/http", "MethodPost"),
				jen.Lit(""),
				jen.ID("expectedPath"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
						jen.ID("assertRequestQuality").Call(
							jen.ID("t"),
							jen.ID("req"),
							jen.ID("spec"),
						),
						jen.Qual("net/http", "SetCookie").Call(
							jen.ID("res"),
							jen.Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").ID("s").Dot("exampleUser").Dot("Username")),
						),
					))),
					jen.ID("c").Op(":=").ID("buildTestClient").Call(
						jen.ID("t"),
						jen.ID("ts"),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("c").Dot("BeginSession").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("c").Dot("BeginSession").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
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
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("c").Dot("BeginSession").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
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
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("c").Dot("BeginSession").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with missing cookie"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("c").Dot("BeginSession").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cookie"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("authTestSuite")).ID("TestClient_Logout").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/logout"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusAccepted"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("EndSession").Call(jen.ID("s").Dot("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("EndSession").Call(jen.ID("s").Dot("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("EndSession").Call(jen.ID("s").Dot("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("authTestSuite")).ID("TestClient_ChangePassword").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/password/new"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPut"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("err").Op(":=").ID("c").Dot("ChangePassword").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil cookie"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("err").Op(":=").ID("c").Dot("ChangePassword").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("ChangePassword").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("err").Op(":=").ID("c").Dot("ChangePassword").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("err").Op(":=").ID("c").Dot("ChangePassword").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with unsatisfactory response code"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPut"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusBadRequest"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("err").Op(":=").ID("c").Dot("ChangePassword").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("authTestSuite")).ID("TestClient_CycleTwoFactorSecret").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/totp_secret/new"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("expected").Op(":=").Op("&").ID("types").Dot("TOTPSecretRefreshResponse").Valuesln(jen.ID("TwoFactorQRCode").Op(":").ID("t").Dot("Name").Call(), jen.ID("TwoFactorSecret").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("expected"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CycleTwoFactorSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil cookie"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("expected").Op(":=").Op("&").ID("types").Dot("TOTPSecretRefreshResponse").Valuesln(jen.ID("TwoFactorQRCode").Op(":").ID("t").Dot("Name").Call(), jen.ID("TwoFactorSecret").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("expected"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CycleTwoFactorSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CycleTwoFactorSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("TOTPSecretRefreshInput").Values(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CycleTwoFactorSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CycleTwoFactorSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CycleTwoFactorSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleCookie"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("authTestSuite")).ID("TestClient_VerifyTOTPSecret").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/totp_secret/verify"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusAccepted"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
						jen.ID("exampleInput").Dot("TOTPToken"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid token"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
						jen.Lit(" doesn't parse lol "),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with bad request response"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusBadRequest"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("ErrInvalidTOTPToken"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with otherwise invalid status code response"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusInternalServerError"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with timeout"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("c").Dot("unauthenticatedClient").Dot("Timeout").Op("=").Qual("time", "Millisecond"),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("s").Dot("exampleUser")),
					jen.ID("err").Op(":=").ID("c").Dot("VerifyTOTPSecret").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
