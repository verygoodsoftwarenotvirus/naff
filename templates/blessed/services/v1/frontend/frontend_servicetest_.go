package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func frontendServiceTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestProvideFrontendService").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("ProvideFrontendService").Call(jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.Qual(proj.InternalConfigV1Package(), "FrontendSettings").Values()),
			)),
		),
		jen.Line(),
	)
	return ret
}
