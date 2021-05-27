package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuilder_BuildUserStatusRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/auth/status"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildUserStatusRequest").Call(jen.ID("helper").Dot("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildLoginRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/login"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserLoginInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildLoginRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildLoginRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("req"),
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
		jen.Func().ID("TestBuilder_BuildLogoutRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/logout"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildLogoutRequest").Call(jen.ID("helper").Dot("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildChangePasswordRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/password/new"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPut"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildChangePasswordRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").Qual("net/http", "Cookie").Values(),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildChangePasswordRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").Qual("net/http", "Cookie").Values(),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("helper").Dot("builder").Op("=").ID("buildTestRequestBuilderWithInvalidURL").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildChangePasswordRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").Qual("net/http", "Cookie").Values(),
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
		jen.Func().ID("TestBuilder_BuildCycleTwoFactorSecretRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/totp_secret/new"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildCycleTwoFactorSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").Qual("net/http", "Cookie").Values(),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil cookie"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildCycleTwoFactorSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildCycleTwoFactorSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").Qual("net/http", "Cookie").Values(),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildCycleTwoFactorSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").Qual("net/http", "Cookie").Values(),
						jen.Op("&").ID("types").Dot("TOTPSecretRefreshInput").Values(),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("helper").Dot("builder").Op("=").ID("buildTestRequestBuilderWithInvalidURL").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildCycleTwoFactorSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").Qual("net/http", "Cookie").Values(),
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
		jen.Func().ID("TestBuilder_BuildVerifyTOTPSecretRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/users/totp_secret/verify"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildVerifyTOTPSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildVerifyTOTPSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Lit(0),
						jen.ID("exampleInput").Dot("TOTPToken"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildVerifyTOTPSecretRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.Lit(" nope lol "),
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

	return code
}
