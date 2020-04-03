package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func metaTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestHoldOnForever").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.If(jen.Qual("os", "Getenv").Call(jen.Lit("WAIT_FOR_COVERAGE")).Op("==").Lit("yes")).Block(
				jen.Comment("snooze for a year"),
				jen.Qual("time", "Sleep").Call(jen.Qual("time", "Hour").Times().Lit(24).Times().Lit(365)),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("checkValueAndError").Params(jen.ID("t").ParamPointer().Qual("testing", "T"), jen.ID("i").Interface(), jen.Err().Error()).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("i")),
		),
		jen.Line(),
	)
	return ret
}
