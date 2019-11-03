package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func runtimeTestDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("TestRecordRuntimeStats").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("this is sort of an obligatory test for coverage's sake"),
			jen.Line(),
			jen.ID("d").Op(":=").Qual("time", "Second"),
			jen.ID("sf").Op(":=").ID("RecordRuntimeStats").Call(jen.ID("d").Op("/").Lit(5)),
			jen.Qual("time", "Sleep").Call(jen.ID("d")),
			jen.ID("sf").Call(),
		),
		jen.Line(),
	)
	return ret
}
