package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func timeTellerDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabasePackage("queriers", "v1", spn), spn)

	utils.AddImports(proj, code, false)

	code.Add(jen.Type().ID("timeTeller").Interface(jen.ID("Now").Call().Uint64()))

	code.Add(jen.Type().ID("stdLibTimeTeller").Struct())

	code.Add(
		jen.Func().Receiver(jen.ID("t").PointerTo().ID("stdLibTimeTeller")).ID("Now").Params().Uint64().Body(
			jen.Return(jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
		),
	)

	code.Add(
		jen.Type().ID("mockTimeTeller").Struct(
			jen.Qual(constants.MockPkg, "Mock"),
		),
	)

	code.Add(
		jen.Func().Receiver(jen.ID("m").PointerTo().ID("mockTimeTeller")).ID("Now").Params().Uint64().Body(
			jen.Return(jen.ID("m").Dot("Called").Call().Dot("Get").Call(jen.Zero()).Assert(jen.Uint64())),
		),
	)

	return code
}
