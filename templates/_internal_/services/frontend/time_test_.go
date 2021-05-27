package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func timeTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_mustGoment").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("ts").Op(":=").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("mustGoment").Call(jen.ID("ts")),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_relativeTime").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("ts").Op(":=").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("relativeTime").Call(jen.ID("ts")),
					),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_relativeTimeFromPtr").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("ts").Op(":=").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("relativeTimeFromPtr").Call(jen.Op("&").ID("ts")),
					),
				),
			)),
		jen.Line(),
	)

	return code
}
