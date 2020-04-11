package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func timeTellerTestDotGo(vendor, dbDesc string) *jen.File {
	ret := jen.NewFile(vendor)

	ret.Add(
		utils.OuterTestFunc("_stdLibTimeTeller_Now").Block(
			utils.ParallelTest(jen.ID("T")),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("tt").Assign().AddressOf().ID("stdLibTimeTeller").Values(),
				jen.Line(),
				utils.AssertNotZero(jen.ID("tt").Dot("Now").Call(), nil),
			),
		),
	)

	return ret
}
