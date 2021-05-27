package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAPIClientCreationInput_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("APIClientCreationInput").Valuesln(jen.ID("UserLoginInput").Op(":").ID("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("t").Dot("Name").Call(), jen.ID("Password").Op(":").ID("t").Dot("Name").Call(), jen.ID("TOTPToken").Op(":").Lit("123456")), jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid UserLoginInput"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("x").Op(":=").Op("&").ID("APIClientCreationInput").Valuesln(jen.ID("UserLoginInput").Op(":").ID("UserLoginInput").Valuesln(jen.ID("Username").Op(":").Lit(""), jen.ID("Password").Op(":").Lit(""), jen.ID("TOTPToken").Op(":").Lit("")), jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("assert").Dot("Error").Call(
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
