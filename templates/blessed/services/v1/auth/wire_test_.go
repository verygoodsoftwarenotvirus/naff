package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	// if proj.EnableNewsman {
	ret.Add(
		jen.Func().ID("TestProvideWebsocketAuthFunc").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.ID("ProvideWebsocketAuthFunc").Call(jen.ID("buildTestService").Call(jen.ID("t"))),
			)),
		),
		jen.Line(),
	)
	// }

	ret.Add(
		jen.Func().ID("TestProvideOAuth2ClientValidator").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
				jen.ID("ProvideOAuth2ClientValidator").Call(jen.AddressOf().Qual(proj.ServiceV1OAuth2ClientsPackage(), "Service").Values()),
			)),
		),
		jen.Line(),
	)

	return ret
}
