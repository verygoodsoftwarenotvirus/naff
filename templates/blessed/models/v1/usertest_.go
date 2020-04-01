package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("TestUser_Update").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("actual").Assign().ID("User").Valuesln(
					jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
					jen.ID("HashedPassword").MapAssign().Lit("hashed_pass"),
					jen.ID("TwoFactorSecret").MapAssign().Lit("two factor secret"),
				),
				jen.ID("exampleInput").Assign().ID("User").Valuesln(
					jen.ID("Username").MapAssign().Add(utils.FakeUsernameFunc()),
					jen.ID("HashedPassword").MapAssign().Lit("updated_hashed_pass"),
					jen.ID("TwoFactorSecret").MapAssign().Lit("new fancy secret"),
				),
				jen.Line(),
				jen.ID("actual").Dot("Update").Call(jen.VarPointer().ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleInput"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)
	return ret
}
