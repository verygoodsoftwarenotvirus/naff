package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(utils.FakeSeedFunc(), jen.Line())
	code.Add(buildTestErrorResponse_Error()...)

	return code
}

func buildTestErrorResponse_Error() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestErrorResponse_Error").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.Underscore().Equals().Parens(jen.AddressOf().ID("ErrorResponse").Values()).Dot("Error").Call(),
			),
		),
		jen.Line(),
	}

	return lines
}
