package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func timeTellerTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	code.Add(
		utils.OuterTestFunc("_stdLibTimeTeller_Now").Body(
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

	return code
}
