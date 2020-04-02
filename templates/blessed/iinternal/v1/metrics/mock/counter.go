package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("_").Qual(pkg.InternalMetricsV1Package(), "UnitCounter").Equals().Parens(jen.Op("*").ID("UnitCounter")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UnitCounter is a mock metrics.UnitCounter"),
		jen.Line(),
		jen.Type().ID("UnitCounter").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Increment implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("Increment").Params(utils.CtxVar().Qual("context", "Context")).Block(
			jen.ID("m").Dot("Called").Call(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IncrementBy implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("IncrementBy").Params(utils.CtxParam(), jen.ID("val").ID("uint64")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("val")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Decrement implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("Decrement").Params(utils.CtxVar().Qual("context", "Context")).Block(
			jen.ID("m").Dot("Called").Call(),
		),
		jen.Line(),
	)
	return ret
}
