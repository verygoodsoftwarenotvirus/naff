package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAddUserToAccountInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("AddUserToAccountInput").Valuesln(jen.ID("UserID").Op(":").Lit(123)),
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
		jen.Func().ID("TestTransferAccountOwnershipInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("AccountOwnershipTransferInput").Valuesln(jen.ID("CurrentOwner").Op(":").Lit(123), jen.ID("NewOwner").Op(":").Lit(321), jen.ID("Reason").Op(":").ID("t").Dot("Name").Call()),
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
		jen.Func().ID("TestModifyUserPermissionsInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("ModifyUserPermissionsInput").Valuesln(jen.ID("NewRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()), jen.ID("Reason").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("x").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
