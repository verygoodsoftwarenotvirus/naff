package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func counterDotGo(pkgRoot string) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkgRoot, "internal/v1/metrics"), "UnitCounter").Op("=").Parens(jen.Op("*").ID("UnitCounter")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UnitCounter is a mock metrics.UnitCounter"),
		jen.Line(),
		jen.Type().ID("UnitCounter").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Increment implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("Increment").Params(jen.ID("ctx").Qual("context", "Context")).Block(
			jen.ID("m").Dot("Called").Call(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IncrementBy implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("IncrementBy").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("val").ID("uint64")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("val")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Decrement implements our UnitCounter interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("Decrement").Params(jen.ID("ctx").Qual("context", "Context")).Block(
			jen.ID("m").Dot("Called").Call(),
		),
		jen.Line(),
	)
	return ret
}