package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestUser_Update").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("actual").Assign().ID("User").Valuesln(
					jen.ID("Username").MapAssign().Lit("old_username"),
					jen.ID("HashedPassword").MapAssign().Lit("hashed_pass"),
					jen.ID("TwoFactorSecret").MapAssign().Lit("two factor secret"),
				),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().ID("User").Valuesln(
					jen.ID("Username").MapAssign().Lit("new_username"),
					jen.ID("HashedPassword").MapAssign().Lit("updated_hashed_pass"),
					jen.ID("TwoFactorSecret").MapAssign().Lit("new fancy secret"),
				),
				jen.Line(),
				jen.ID("actual").Dot("Update").Call(jen.AddressOf().ID("exampleInput")),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Input")), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	return ret
}
