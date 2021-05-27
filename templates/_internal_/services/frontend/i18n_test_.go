package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func i18NTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestService_getSimpleLocalizedString").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Lit(":)"),
						jen.ID("s").Dot("getSimpleLocalizedString").Call(jen.Lit("testing.translation")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_provideLocalizer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("obligatory"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("provideLocalizer").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
