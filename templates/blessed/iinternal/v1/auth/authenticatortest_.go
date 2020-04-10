package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authenticatorTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth_test")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestProvideBcryptHashCost").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Qual(proj.InternalAuthV1Package(), "ProvideBcryptHashCost").Call(),
			),
		),
		jen.Line(),
	)
	return ret
}
