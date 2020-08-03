package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authenticatorTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("auth_test")

	utils.AddImports(proj, code)

	code.Add(buildTestProvideBcryptHashCost(proj)...)

	return code
}

func buildTestProvideBcryptHashCost(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideBcryptHashCost").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Qual(proj.InternalAuthV1Package(), "ProvideBcryptHashCost").Call(),
			),
		),
		jen.Line(),
	}

	return lines
}
