package converters

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func convertersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestConvertAuditLogEntryCreationInputToEntry").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryCreationInput").Call(),
					jen.ID("actual").Op(":=").ID("ConvertAuditLogEntryCreationInputToEntry").Call(jen.ID("exampleInput")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("EventType"),
						jen.ID("actual").Dot("EventType"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("Context"),
						jen.ID("actual").Dot("Context"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConvertAccountToAccountUpdateInput").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("actual").Op(":=").ID("ConvertAccountToAccountUpdateInput").Call(jen.ID("exampleInput")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("Name"),
						jen.ID("actual").Dot("Name"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleInput").Dot("BelongsToUser"),
						jen.ID("actual").Dot("BelongsToUser"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestConvertItemToItemUpdateInput").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expected").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("actual").Op(":=").ID("ConvertItemToItemUpdateInput").Call(jen.ID("expected")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected").Dot("Name"),
						jen.ID("actual").Dot("Name"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected").Dot("Details"),
						jen.ID("actual").Dot("Details"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
