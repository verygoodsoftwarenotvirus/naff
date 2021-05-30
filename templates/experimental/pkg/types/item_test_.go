package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestItem_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Item").Valuesln(),
					jen.ID("updated").Op(":=").Op("&").ID("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call(), jen.ID("Details").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call()),
					jen.ID("expected").Op(":=").Index().Op("*").ID("FieldChangeSummary").Valuesln(jen.Valuesln(jen.ID("FieldName").Op(":").Lit("Name"), jen.ID("OldValue").Op(":").ID("i").Dot("Name"), jen.ID("NewValue").Op(":").ID("updated").Dot("Name")), jen.Valuesln(jen.ID("FieldName").Op(":").Lit("Details"), jen.ID("OldValue").Op(":").ID("i").Dot("Details"), jen.ID("NewValue").Op(":").ID("updated").Dot("Details"))),
					jen.ID("actual").Op(":=").ID("i").Dot("Update").Call(jen.ID("updated")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
						jen.Lit("expected and actual diff reports vary"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("updated").Dot("Name"),
						jen.ID("i").Dot("Name"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("updated").Dot("Details"),
						jen.ID("i").Dot("Details"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemCreationInput_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("ItemCreationInput").Valuesln(jen.ID("Name").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call(), jen.ID("Details").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call()),
					jen.ID("actual").Op(":=").ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid structure"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("ItemCreationInput").Valuesln(jen.ID("Name").Op(":").Lit(""), jen.ID("Details").Op(":").Lit("")),
					jen.ID("actual").Op(":=").ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemUpdateInput_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call(), jen.ID("Details").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call()),
					jen.ID("actual").Op(":=").ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with empty strings"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("ItemUpdateInput").Valuesln(jen.ID("Name").Op(":").Lit(""), jen.ID("Details").Op(":").Lit("")),
					jen.ID("actual").Op(":=").ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
