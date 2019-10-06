package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func metaTestDotGo() *jen.File {
	ret := jen.NewFile("metrics")

	ret.Add(jen.Func().ID("TestRegisterDefaultViews").Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		utils.RequireNoError(jen.ID("t"), jen.ID("RegisterDefaultViews").Call()),
	),
	)
	return ret
}
