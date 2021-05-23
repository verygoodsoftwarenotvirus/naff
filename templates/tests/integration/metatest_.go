package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metaTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestHoldOnForever()...)
	code.Add(buildCheckValueAndError()...)

	return code
}

func buildTestHoldOnForever() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestHoldOnForever").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.If(jen.Qual("os", "Getenv").Call(jen.Lit("WAIT_FOR_COVERAGE")).IsEqualTo().Lit("yes")).Body(
				jen.Comment("snooze for a year."),
				jen.Qual("time", "Sleep").Call(jen.Qual("time", "Hour").Times().Lit(24).Times().Lit(365)),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildCheckValueAndError() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("checkValueAndError").Params(jen.ID("t").PointerTo().Qual("testing", "T"), jen.ID("i").Interface(), jen.Err().Error()).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			utils.RequireNoError(jen.Err(), nil),
			utils.RequireNotNil(jen.ID("i"), nil),
		),
		jen.Line(),
	}

	return lines
}
