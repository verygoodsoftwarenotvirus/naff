package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)

	code.Add(utils.FakeSeedFunc(), jen.Line())

	code.Add(
		jen.Func().ID("TestErrorResponse_Error").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().Parens(jen.AddressOf().ID("ErrorResponse").Values()).Dot("Error").Call(),
			),
		),
		jen.Line(),
	)

	return code
}
