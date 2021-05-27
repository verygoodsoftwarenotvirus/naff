package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAccount_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("Account").Values(),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Update").Call(jen.Op("&").ID("AccountUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call())),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAccountCreationInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("AccountCreationInput").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
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
		jen.Func().ID("TestAccountUpdateInput_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("AccountUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
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
		jen.Func().ID("TestAccountCreationInputForNewUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("AccountCreationInputForNewUser").Call(jen.Op("&").ID("User").Values()),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
