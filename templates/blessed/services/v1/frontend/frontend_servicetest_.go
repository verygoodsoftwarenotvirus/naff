package frontend

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func frontendServiceTestDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("TestProvideFrontendService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("ProvideFrontendService").Call(jen.ID("noop").Dot("ProvideNoopLogger").Call(), jen.Qual(filepath.Join(pkgRoot, "internal/v1/config"), "FrontendSettings").Values()),
			)),
		),
		jen.Line(),
	)
	return ret
}
