package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func runtimeTestDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestRecordRuntimeStats").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("d").Op(":=").Qual("time", "Second"),
		jen.ID("sf").Op(":=").ID("RecordRuntimeStats").Call(jen.ID("d").Op("/").Lit(5)),
		jen.Qual("time", "Sleep").Call(jen.ID("d")),
		jen.ID("sf").Call(),
	),

		jen.Line(),
	)
	return ret
}
