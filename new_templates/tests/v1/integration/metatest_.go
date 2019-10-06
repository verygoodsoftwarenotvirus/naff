package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func metaTestDotGo() *jen.File {
	ret := jen.NewFile("integration")
	utils.AddImports(ret)

	ret.Add(jen.Func().ID("TestHoldOnForever").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.If(
			jen.Qual("os", "Getenv").Call(jen.Lit("WAIT_FOR_COVERAGE")).Op("==").Lit("yes"),
		).Block(
			jen.Qual("time", "Sleep").Call(jen.Qual("time", "Hour").Op("*").Lit(24).Op("*").Lit(365)),
		),
	),
	)
	ret.Add(jen.Func().ID("checkValueAndError").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("i").Interface(), jen.ID("err").ID("error")).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		utils.RequireNoError(jen.ID("t"), jen.ID("err")),
		jen.ID("require").Dot(
			"NotNil",
		).Call(jen.ID("t"), jen.ID("i")),
	),
	)
	return ret
}
