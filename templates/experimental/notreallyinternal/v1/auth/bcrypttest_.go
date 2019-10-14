package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func bcryptTestDotGo() *jen.File {
	ret := jen.NewFile("auth_test")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("examplePassword").Op("=").Lit("Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"),
			jen.ID("weaklyHashedExamplePassword").Op("=").Lit("$2a$04$7G7dHZe7MeWjOMsYKO8uCu/CRKnDMMBHOfXaB6YgyQL/cl8nhwf/2"),
			jen.ID("hashedExamplePassword").Op("=").Lit("$2a$13$hxMAo/ZRDmyaWcwvIem/vuUJkmeNytg3rwHUj6bRZR1d/cQHXjFvW"),
			jen.ID("exampleTwoFactorSecret").Op("=").Lit("HEREISASECRETWHICHIVEMADEUPBECAUSEIWANNATESTRELIABLY"),
		),
		jen.Line(),
	)
	ret.Add(
		jen.Func().ID("TestBcrypt_HashPassword").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("x").Op(":=").ID("auth").Dot("ProvideBcryptAuthenticator").Call(jen.ID("auth").Dot("DefaultBcryptHashCost"),
				jen.ID("noop").Dot("ProvideNoopLogger").Call()),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("x").Dot("HashPassword").Call(jen.ID("tctx"), jen.Lit("password")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NotEmpty").Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_PasswordMatches").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Parallel",
			).Call(),
			jen.ID("x").Op(":=").ID("auth").Dot(
				"ProvideBcryptAuthenticator",
			).Call(jen.ID("auth").Dot(
				"DefaultBcryptHashCost",
			), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("normal usage"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot(
					"Parallel",
				).Call(),
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("actual").Op(":=").ID("x").Dot(
					"PasswordMatches",
				).Call(jen.ID("tctx"), jen.ID("hashedExamplePassword"), jen.ID("examplePassword"), jen.ID("nil")),
				jen.ID("assert").Dot(
					"True",
				).Call(jen.ID("t"), jen.ID("actual")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("when passwords don't match"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot(
					"Parallel",
				).Call(),
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("actual").Op(":=").ID("x").Dot(
					"PasswordMatches",
				).Call(jen.ID("tctx"), jen.ID("hashedExamplePassword"), jen.Lit("password"), jen.ID("nil")),
				jen.ID("assert").Dot(
					"False",
				).Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_PasswordIsAcceptable").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("x").Op(":=").ID("auth").Dot("ProvideBcryptAuthenticator").Call(jen.ID("auth").Dot(
				"DefaultBcryptHashCost",
			), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call()),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot(
					"Parallel",
				).Call(),
				jen.ID("assert").Dot(
					"True",
				).Call(jen.ID("t"), jen.ID("x").Dot(
					"PasswordIsAcceptable",
				).Call(jen.ID("examplePassword"))),
				jen.ID("assert").Dot(
					"False",
				).Call(jen.ID("t"), jen.ID("x").Dot(
					"PasswordIsAcceptable",
				).Call(jen.Lit("hi there"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_ValidateLogin").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("x").Op(":=").ID("auth").Dot("ProvideBcryptAuthenticator").Call(jen.ID("auth").Dot("DefaultBcryptHashCost"), jen.ID("noop").Dot("ProvideNoopLogger").Call()),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(jen.ID("exampleTwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err"), jen.Lit("error generating code to validate login")),
				jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
					jen.Qual("context", "Background").Call(),
					jen.ID("hashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.ID("code"),
					jen.ID("nil"),
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
					jen.Lit("unexpected error encountered validating login: %v"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("True").Call(jen.ID("t"), jen.ID("valid")),
			)),

			jen.ID("T").Dot("Run").Call(jen.Lit("with weak hash"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
					jen.ID("exampleTwoFactorSecret"),
					jen.Qual("time", "Now").Call().Dot("UTC").Call(),
				),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err"), jen.Lit("error generating code to validate login")),
				jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
					jen.Qual("context", "Background").Call(),
					jen.ID("weaklyHashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.ID("code"),
					jen.ID("nil"),
				),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err"), jen.Lit("unexpected error encountered validating login: %v"), jen.ID("err")),
				jen.ID("assert").Dot("True").Call(jen.ID("t"), jen.ID("valid")),
			)),

			jen.ID("T").Dot("Run").Call(jen.Lit("with non-matching password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(jen.ID("exampleTwoFactorSecret"),
					jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err"), jen.Lit("error generating code to validate login")),
				jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
					jen.Qual("context", "Background").Call(),
					jen.ID("hashedExamplePassword"),
					jen.Lit("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.ID("code"),
					jen.ID("nil"),
				),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err"), jen.Lit("unexpected error encountered validating login: %v"), jen.ID("err")),
				jen.ID("assert").Dot("False").Call(jen.ID("t"), jen.ID("valid")),
			)),

			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid code"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
					jen.Qual("context", "Background").Call(),
					jen.ID("hashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.Lit("CODE"),
					jen.ID("nil"),
				),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err"), jen.Lit("unexpected error encountered validating login: %v"), jen.ID("err")),
				jen.ID("assert").Dot("True").Call(jen.ID("t"), jen.ID("valid")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideBcrypt").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("auth").Dot("ProvideBcryptAuthenticator").Call(
					jen.ID("auth").Dot("DefaultBcryptHashCost"),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
				),
			)),
		),

		jen.Line(),
	)
	return ret
}
