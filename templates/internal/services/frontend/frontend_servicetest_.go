package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func frontendServiceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestProvideFrontendService(proj)...)

	return code
}

func buildTestProvideFrontendService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideFrontendService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("ProvideFrontendService").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.Qual(proj.InternalConfigPackage(), "FrontendSettings").Values()),
			),
		),
		jen.Line(),
	}

	return lines
}
