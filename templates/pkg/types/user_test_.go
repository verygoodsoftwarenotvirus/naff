package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestIsValidAccountStatus").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("IsValidAccountStatus").Call(jen.ID("string").Call(jen.ID("GoodStandingAccountStatus"))),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("IsValidAccountStatus").Call(jen.Lit("blah")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUser_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("actual").Op(":=").ID("User").Valuesln(jen.ID("Username").Op(":").Lit("old_username"), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("TwoFactorSecret").Op(":").Lit("two factor secret")),
					jen.ID("exampleInput").Op(":=").ID("User").Valuesln(jen.ID("Username").Op(":").Lit("new_username"), jen.ID("HashedPassword").Op(":").Lit("updated_hashed_pass"), jen.ID("TwoFactorSecret").Op(":").Lit("new fancy secret")),
					jen.ID("actual").Dot("Update").Call(jen.Op("&").ID("exampleInput")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleInput"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUser_IsBanned").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("User").Valuesln(jen.ID("ServiceAccountStatus").Op(":").ID("BannedUserAccountStatus")),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("x").Dot("IsBanned").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPasswordUpdateInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("PasswordUpdateInput").Valuesln(jen.ID("NewPassword").Op(":").ID("t").Dot("Name").Call(), jen.ID("CurrentPassword").Op(":").ID("t").Dot("Name").Call(), jen.ID("TOTPToken").Op(":").Lit("123456")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("x").Dot("ValidateWithContext").Call(
							jen.ID("ctx"),
							jen.Lit(1),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestTOTPSecretRefreshInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("TOTPSecretRefreshInput").Valuesln(jen.ID("CurrentPassword").Op(":").ID("t").Dot("Name").Call(), jen.ID("TOTPToken").Op(":").Lit("123456")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("x").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestTOTPSecretVerificationInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("TOTPSecretVerificationInput").Valuesln(jen.ID("UserID").Op(":").Lit(123), jen.ID("TOTPToken").Op(":").Lit("123456")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("x").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUserCreationInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("UserRegistrationInput").Valuesln(jen.ID("Username").Op(":").ID("t").Dot("Name").Call(), jen.ID("Password").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("x").Dot("ValidateWithContext").Call(
							jen.ID("ctx"),
							jen.Lit(1),
							jen.Lit(1),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUserLoginInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("t").Dot("Name").Call(), jen.ID("Password").Op(":").ID("t").Dot("Name").Call(), jen.ID("TOTPToken").Op(":").Lit("123456")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("x").Dot("ValidateWithContext").Call(
							jen.ID("ctx"),
							jen.Lit(1),
							jen.Lit(1),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
