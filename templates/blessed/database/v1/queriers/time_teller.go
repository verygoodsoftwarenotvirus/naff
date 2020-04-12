package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func timeTellerDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	ret := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Type().ID("timeTeller").Interface(
			jen.ID("Now").Call().Uint64(),
		),
	)

	ret.Add(
		jen.Type().ID("stdLibTimeTeller").Struct(),
	)

	ret.Add(
		jen.Func().Receiver(jen.ID("t").PointerTo().ID("stdLibTimeTeller")).ID("Now").Params().Uint64().Block(
			jen.Return(jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
		),
	)

	ret.Add(
		jen.Type().ID("mockTimeTeller").Struct(
			jen.Qual(utils.MockPkg, "Mock"),
		),
	)

	ret.Add(
		jen.Func().Receiver(jen.ID("m").PointerTo().ID("mockTimeTeller")).ID("Now").Params().Uint64().Block(
			jen.Return(jen.ID("m").Dot("Called").Call().Dot("Get").Call(jen.Zero()).Assert(jen.Uint64())),
		),
	)

	return ret
}
