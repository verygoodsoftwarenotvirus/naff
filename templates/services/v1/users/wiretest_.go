package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	// if proj.EnableNewsman {
	code.Add(buildTestProvideUserDataManager(proj)...)
	// }

	code.Add(buildTestProvideUserDataServer()...)

	return code
}

func buildTestProvideUserDataManager(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideUserDataManager").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				utils.AssertNotNil(jen.ID("ProvideUserDataManager").Call(
					jen.Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				), nil),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideUserDataServer() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideUserDataServer").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				utils.AssertNotNil(jen.ID("ProvideUserDataServer").Call(jen.ID("buildTestService").Call(jen.ID("t"))), nil),
			)),
		),
		jen.Line(),
	}

	return lines
}
