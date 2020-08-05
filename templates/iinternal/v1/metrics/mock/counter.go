package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(buildUnitCounter(proj)...)
	code.Add(buildIncrement()...)
	code.Add(buildIncrementBy()...)
	code.Add(buildDecrement()...)

	return code
}

func buildUnitCounter(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual(proj.InternalMetricsV1Package(), "UnitCounter").Equals().Parens(jen.PointerTo().ID("UnitCounter")).Call(jen.Nil()),
		jen.Line(),
		jen.Comment("UnitCounter is a mock metrics.UnitCounter"),
		jen.Line(),
		jen.Type().ID("UnitCounter").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildIncrement() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Increment implements our UnitCounter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UnitCounter")).ID("Increment").Params(constants.CtxParam()).Body(
			jen.ID("m").Dot("Called").Call(constants.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}

func buildIncrementBy() []jen.Code {
	lines := []jen.Code{
		jen.Comment("IncrementBy implements our UnitCounter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UnitCounter")).ID("IncrementBy").Params(constants.CtxParam(), jen.ID("val").Uint64()).Body(
			jen.ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("val")),
		),
		jen.Line(),
	}

	return lines
}

func buildDecrement() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Decrement implements our UnitCounter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UnitCounter")).ID("Decrement").Params(constants.CtxParam()).Body(
			jen.ID("m").Dot("Called").Call(constants.CtxVar()),
		),
		jen.Line(),
	}

	return lines
}
