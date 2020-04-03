package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bcryptTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth_test")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("examplePassword").Equals().Lit("Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"),
			jen.ID("weaklyHashedExamplePassword").Equals().Lit("$2a$04$7G7dHZe7MeWjOMsYKO8uCu/CRKnDMMBHOfXaB6YgyQL/cl8nhwf/2"),
			jen.ID("hashedExamplePassword").Equals().Lit("$2a$13$hxMAo/ZRDmyaWcwvIem/vuUJkmeNytg3rwHUj6bRZR1d/cQHXjFvW"),
			jen.ID("exampleTwoFactorSecret").Equals().Lit("HEREISASECRETWHICHIVEMADEUPBECAUSEIWANNATESTRELIABLY"),
		),
		jen.Line(),
	)
	ret.Add(
		jen.Func().ID("TestBcrypt_HashPassword").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("t").Dot("Parallel").Call(),
				jen.ID("tctx").Assign().Qual("context", "Background").Call(),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("x").Dot("HashPassword").Call(jen.ID("tctx"), jen.Lit("password")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotEmpty(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_PasswordMatches").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"normal usage",
				jen.ID("t").Dot("Parallel").Call(),
				jen.ID("tctx").Assign().Qual("context", "Background").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("x").Dot("PasswordMatches").Call(jen.ID("tctx"), jen.ID("hashedExamplePassword"), jen.ID("examplePassword"), jen.Nil()),
				utils.AssertTrue(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"when passwords don't match",
				jen.ID("t").Dot("Parallel").Call(),
				jen.ID("tctx").Assign().Qual("context", "Background").Call(),
				jen.Line(),
				jen.ID("actual").Assign().ID("x").Dot("PasswordMatches").Call(jen.ID("tctx"), jen.ID("hashedExamplePassword"), jen.Lit("password"), jen.Nil()),
				utils.AssertFalse(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_PasswordIsAcceptable").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("t").Dot("Parallel").Call(),
				jen.Line(),
				utils.AssertTrue(jen.ID("x").Dot("PasswordIsAcceptable").Call(jen.ID("examplePassword")), nil),
				utils.AssertFalse(jen.ID("x").Dot("PasswordIsAcceptable").Call(jen.Lit("hi there")), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_ValidateLogin").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Assign().Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("t").Dot("Parallel").Call(),
				utils.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("exampleTwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				utils.AssertNoError(jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
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
				utils.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
					jen.ID("exampleTwoFactorSecret"),
					jen.Qual("time", "Now").Call().Dot("UTC").Call(),
				),
				utils.AssertNoError(jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("weaklyHashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
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
				utils.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("exampleTwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				utils.AssertNoError(jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.Lit("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
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
				utils.CreateCtx(),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Assign().ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.Lit("CODE"),
					jen.Nil(),
				),
				utils.AssertError(jen.Err(), jen.Lit("unexpected error encountered validating login: %v"), jen.Err()),
				utils.AssertTrue(jen.ID("valid"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideBcrypt").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Qual(proj.InternalAuthV1Package(), "ProvideBcryptAuthenticator").Call(
					jen.Qual(proj.InternalAuthV1Package(), "DefaultBcryptHashCost"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
			),
		),

		jen.Line(),
	)
	return ret
}
