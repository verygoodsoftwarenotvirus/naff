package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func argon2TestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("argon2HashedExamplePassword").Op("=").Lit(`$argon2id$v=19$m=65536,t=1,p=2$C+YWiNi21e94acF3ip8UGA$Ru6oL96HZSP7cVcfAbRwOuK9+vwBo/BLhCzOrGrMH0M`),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestArgon2_HashPassword").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("x").Op(":=").ID("authentication").Dot("ProvideArgon2Authenticator").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("x").Dot("HashPassword").Call(
						jen.ID("ctx"),
						jen.ID("examplePassword"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestArgon2_ValidateLogin").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("x").Op(":=").ID("authentication").Dot("ProvideArgon2Authenticator").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
						jen.ID("exampleTwoFactorSecret"),
						jen.Qual("time", "Now").Call().Dot("UTC").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error generating code to validate login"),
					),
					jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
						jen.ID("ctx"),
						jen.ID("argon2HashedExamplePassword"),
						jen.ID("examplePassword"),
						jen.ID("exampleTwoFactorSecret"),
						jen.ID("code"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("unexpected error encountered validating login: %v"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("valid"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error determining if password matches"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
						jen.ID("ctx"),
						jen.Lit("       blah blah blah not a valid hash lol           "),
						jen.ID("examplePassword"),
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("unexpected error encountered validating login: %v"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("valid"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with non-matching password"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot("GenerateCode").Call(
						jen.ID("exampleTwoFactorSecret"),
						jen.Qual("time", "Now").Call().Dot("UTC").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("error generating code to validate login"),
					),
					jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
						jen.ID("ctx"),
						jen.ID("argon2HashedExamplePassword"),
						jen.Lit("examplePassword"),
						jen.ID("exampleTwoFactorSecret"),
						jen.ID("code"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("unexpected error encountered validating login: %v"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("valid"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid code"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("valid"), jen.ID("err")).Op(":=").ID("x").Dot("ValidateLogin").Call(
						jen.ID("ctx"),
						jen.ID("argon2HashedExamplePassword"),
						jen.ID("examplePassword"),
						jen.ID("exampleTwoFactorSecret"),
						jen.Lit("CODE"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
						jen.Lit("unexpected error encountered validating login: %v"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("valid"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideArgon2Authenticator").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("authentication").Dot("ProvideArgon2Authenticator").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
				),
			),
		),
		jen.Line(),
	)

	return code
}
