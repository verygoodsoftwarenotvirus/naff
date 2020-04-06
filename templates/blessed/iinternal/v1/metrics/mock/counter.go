package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.InternalMetricsV1Package(), "UnitCounter").Equals().Parens(jen.PointerTo().ID("UnitCounter")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UnitCounter is a mock metrics.UnitCounter"),
		jen.Line(),
		jen.Type().ID("UnitCounter").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Increment implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UnitCounter")).ID("Increment").Params(utils.CtxParam()).Block(
			jen.ID("m").Dot("Called").Call(utils.CtxVar()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IncrementBy implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UnitCounter")).ID("IncrementBy").Params(utils.CtxParam(), jen.ID("val").Uint64()).Block(
			jen.ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("val")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Decrement implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UnitCounter")).ID("Decrement").Params(utils.CtxParam()).Block(
			jen.ID("m").Dot("Called").Call(utils.CtxVar()),
		),
		jen.Line(),
	)
	return ret
}
