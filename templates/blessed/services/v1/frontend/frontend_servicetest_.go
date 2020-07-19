package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func frontendServiceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("frontend")

	utils.AddImports(proj, code)

	code.Add(
		jen.Func().ID("TestProvideFrontendService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("ProvideFrontendService").Call(jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.Qual(proj.InternalConfigV1Package(), "FrontendSettings").Values()),
			),
		),
		jen.Line(),
	)

	return code
}
