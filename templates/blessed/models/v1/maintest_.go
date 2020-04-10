package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(utils.FakeSeedFunc(), jen.Line())

	ret.Add(
		jen.Func().ID("TestErrorResponse_Error").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().Parens(jen.AddressOf().ID("ErrorResponse").Values()).Dot("Error").Call(),
			),
		),
		jen.Line(),
	)
	return ret
}
