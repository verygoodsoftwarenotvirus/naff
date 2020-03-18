package auth

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bcryptTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("auth_test")

	utils.AddImports(pkg, ret)

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
			jen.Line(),
			jen.ID("x").Op(":=").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "ProvideBcryptAuthenticator").Call(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("x").Dot("HashPassword").Call(jen.ID("tctx"), jen.Lit("password")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_PasswordMatches").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Op(":=").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "ProvideBcryptAuthenticator").Call(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("normal usage"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("x").Dot("PasswordMatches").Call(jen.ID("tctx"), jen.ID("hashedExamplePassword"), jen.ID("examplePassword"), jen.ID("nil")),
				jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("when passwords don't match"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Line(),
				jen.ID("actual").Op(":=").ID("x").Dot("PasswordMatches").Call(jen.ID("tctx"), jen.ID("hashedExamplePassword"), jen.Lit("password"), jen.ID("nil")),
				jen.Qual("github.com/stretchr/testify/assert", "False").Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_PasswordIsAcceptable").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Op(":=").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "ProvideBcryptAuthenticator").Call(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("x").Dot("PasswordIsAcceptable").Call(jen.ID("examplePassword"))),
				jen.Qual("github.com/stretchr/testify/assert", "False").Call(jen.ID("t"), jen.ID("x").Dot("PasswordIsAcceptable").Call(jen.Lit("hi there"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestBcrypt_ValidateLogin").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("x").Op(":=").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "ProvideBcryptAuthenticator").Call(jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "DefaultBcryptHashCost"), jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call()),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("exampleTwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Op(":=").ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.ID("code"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(
					jen.ID("t"),
					jen.Err(),
					jen.Lit("unexpected error encountered validating login: %v"),
					jen.Err(),
				),
				jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("valid")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with weak hash"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
					jen.ID("exampleTwoFactorSecret"),
					jen.Qual("time", "Now").Call().Dot("UTC").Call(),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Op(":=").ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("weaklyHashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.ID("code"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err(), jen.Lit("unexpected error encountered validating login: %v"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("valid")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with non-matching password"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Line(),
				jen.List(jen.ID("code"), jen.Err()).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("exampleTwoFactorSecret"), jen.Qual("time", "Now").Call().Dot("UTC").Call()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err(), jen.Lit("error generating code to validate login")),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Op(":=").ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.Lit("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.ID("code"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err(), jen.Lit("unexpected error encountered validating login: %v"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "False").Call(jen.ID("t"), jen.ID("valid")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid code"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Line(),
				jen.List(jen.ID("valid"), jen.Err()).Op(":=").ID("x").Dot("ValidateLogin").Callln(
					utils.CtxVar(),
					jen.ID("hashedExamplePassword"),
					jen.ID("examplePassword"),
					jen.ID("exampleTwoFactorSecret"),
					jen.Lit("CODE"),
					jen.ID("nil"),
				),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err(), jen.Lit("unexpected error encountered validating login: %v"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("valid")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideBcrypt").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "ProvideBcryptAuthenticator").Call(
					jen.Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "DefaultBcryptHashCost"),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
			)),
		),

		jen.Line(),
	)
	return ret
}
