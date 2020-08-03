package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bcryptTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("auth_test")

	utils.AddImports(proj, code)

	code.Add(buildConstDefinitions()...)
	code.Add(buildTestBcrypt_HashPassword(proj)...)
	code.Add(buildTestBcrypt_PasswordMatches(proj)...)
	code.Add(buildTestBcrypt_PasswordIsAcceptable(proj)...)
	code.Add(buildTestBcrypt_ValidateLogin(proj)...)
	code.Add(buildTestProvideBcrypt(proj)...)

	return code
}

func buildConstDefinitions() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID(utils.BuildFakeVarName("Password")).Equals().Lit("Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"),
			jen.ID("weaklyHashedExamplePassword").Equals().Lit("$2a$04$7G7dHZe7MeWjOMsYKO8uCu/CRKnDMMBHOfXaB6YgyQL/cl8nhwf/2"),
			jen.ID("hashedExamplePassword").Equals().Lit("$2a$13$hxMAo/ZRDmyaWcwvIem/vuUJkmeNytg3rwHUj6bRZR1d/cQHXjFvW"),
			jen.ID(utils.BuildFakeVarName("TwoFactorSecret")).Equals().Lit("HEREISASECRETWHICHIVEMADEUPBECAUSEIWANNATESTRELIABLY"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestBcrypt_HashPassword(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestBcrypt_HashPassword").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("t").Dot("Parallel").Call(),
				constants.CtxVar().Assign().Qual("context", "Background").Call(),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("x").Dot("HashPassword").Call(constants.CtxVar(), jen.Lit("password")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotEmpty(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestBcrypt_PasswordMatches(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestBcrypt_PasswordMatches").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"normal usage",
				jen.ID("t").Dot("Parallel").Call(),
				constants.CtxVar().Assign().Qual("context", "Background").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("x").Dot("PasswordMatches").Call(constants.CtxVar(), jen.ID("hashedExamplePassword"), jen.ID(utils.BuildFakeVarName("Password")), jen.Nil()),
				utils.AssertTrue(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"when passwords don't match",
				jen.ID("t").Dot("Parallel").Call(),
				constants.CtxVar().Assign().Qual("context", "Background").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("x").Dot("PasswordMatches").Call(constants.CtxVar(), jen.ID("hashedExamplePassword"), jen.Lit("password"), jen.Nil()),
				utils.AssertFalse(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestBcrypt_PasswordIsAcceptable(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestBcrypt_PasswordIsAcceptable").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("t").Dot("Parallel").Call(),
				jen.Line(),
				utils.AssertTrue(jen.ID("x").Dot("PasswordIsAcceptable").Call(jen.ID(utils.BuildFakeVarName("Password"))), nil),
				utils.AssertFalse(jen.ID("x").Dot("PasswordIsAcceptable").Call(jen.Lit("hi there")), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestBcrypt_ValidateLogin(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestBcrypt_ValidateLogin").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("t").Dot("Parallel").Call(),
				constants.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID(utils.BuildFakeVarName("TwoFactorSecret")), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				utils.AssertNoError(jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					constants.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.ID(utils.BuildFakeVarName("Password")),
					jen.ID(utils.BuildFakeVarName("TwoFactorSecret")),
					jen.ID("code"),
					jen.Nil(),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("unexpected error encountered validating login: %v"),
					jen.Err(),
				),
				utils.AssertTrue(jen.ID("valid"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with weak hash",
				jen.ID("t").Dot("Parallel").Call(),
				constants.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
					jen.ID(utils.BuildFakeVarName("TwoFactorSecret")),
					jen.Qual("time", "Now").Call().Dot("UTC").Call(),
				),
				utils.AssertNoError(jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					constants.CtxVar(),
					jen.ID("weaklyHashedExamplePassword"),
					jen.ID(utils.BuildFakeVarName("Password")),
					jen.ID(utils.BuildFakeVarName("TwoFactorSecret")),
					jen.ID("code"),
					jen.Nil(),
				),
				utils.AssertError(jen.Err(), jen.Lit("unexpected error encountered validating login: %v"), jen.Err()),
				utils.AssertTrue(jen.ID("valid"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with non-matching password",
				jen.ID("t").Dot("Parallel").Call(),
				constants.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID(utils.BuildFakeVarName("TwoFactorSecret")), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				utils.AssertNoError(jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					constants.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.Lit("examplePassword"),
					jen.ID(utils.BuildFakeVarName("TwoFactorSecret")),
					jen.ID("code"),
					jen.Nil(),
				),
				utils.AssertNoError(jen.Err(), jen.Lit("unexpected error encountered validating login: %v"), jen.Err()),
				utils.AssertFalse(jen.ID("valid"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid code",
				jen.ID("t").Dot("Parallel").Call(),
				constants.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					constants.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.ID(utils.BuildFakeVarName("Password")),
					jen.ID(utils.BuildFakeVarName("TwoFactorSecret")),
					jen.Lit("CODE"),
					jen.Nil(),
				),
				utils.AssertError(jen.Err(), jen.Lit("unexpected error encountered validating login: %v"), jen.Err()),
				utils.AssertTrue(jen.ID("valid"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideBcrypt(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideBcrypt").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(
					jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
			),
		),
		jen.Line(),
	}

	return lines
}
