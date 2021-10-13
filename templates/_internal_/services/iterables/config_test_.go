package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(

		jen.Func().ID("TestConfig_ValidateWithContext").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					constants.CreateCtx(),
					jen.ID("cfg").Assign().AddressOf().ID("Config").Valuesln(
						jen.ID("PreWritesTopicName").MapAssign().Lit("blah"),
						jen.ID("PreUpdatesTopicName").MapAssign().Lit("blah"),
						jen.ID("PreArchivesTopicName").MapAssign().Lit("blah"),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("SearchIndexPath").MapAssign().Lit("blah")),
					),
					jen.Newline(),
					utils.AssertNoError(jen.ID("cfg").Dot("ValidateWithContext").Call(constants.CtxVar()), nil),
				),
			),
		),
	)

	return code
}
