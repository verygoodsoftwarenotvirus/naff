package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metaTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestHoldOnForever").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.If(jen.Qual("os", "Getenv").Call(jen.Lit("WAIT_FOR_COVERAGE")).Op("==").Lit("yes")).Body(
				jen.Qual("time", "Sleep").Call(jen.Qual("time", "Hour").Op("*").Lit(24).Op("*").Lit(365))),
		),
		jen.Line(),
	)

	return code
}
