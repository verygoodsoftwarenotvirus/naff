package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	// if proj.EnableNewsman {
	code.Add(
		jen.Func().ID("TestProvideWebsocketAuthFunc").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				utils.AssertNotNil(jen.ID("ProvideWebsocketAuthFunc").Call(jen.ID("buildTestService").Call(jen.ID("t"))), nil),
			)),
		),
		jen.Line(),
	)
	// }

	code.Add(
		jen.Func().ID("TestProvideOAuth2ClientValidator").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testprojects", "T")).Block(
				utils.AssertNotNil(jen.ID("ProvideOAuth2ClientValidator").Call(jen.AddressOf().Qual(proj.ServiceV1OAuth2ClientsPackage(), "Service").Values()), nil),
			)),
		),
		jen.Line(),
	)

	return code
}
