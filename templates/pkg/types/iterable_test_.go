package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()

	var (
		updateCols           []jen.Code
		fieldChangeSummaries []jen.Code
		assertions           []jen.Code
	)

	for _, field := range typ.Fields {
		fsn := field.Name.Singular()

		if field.Type == "bool" {
			if field.IsPointer {
				updateCols = append(updateCols, jen.ID(fsn).MapAssign().Func().Params(jen.ID("x").Bool()).Params(jen.PointerTo().Bool()).SingleLineBody(jen.Return(jen.AddressOf().ID("x"))).Call(jen.True()))
			} else {
				updateCols = append(updateCols, jen.ID(fsn).MapAssign().True())
			}
		} else {
			updateCols = append(updateCols, jen.ID(fsn).MapAssign().Add(utils.FakeFuncForType(field.Type, field.IsPointer)()))
		}

		fieldChangeSummaries = append(fieldChangeSummaries,
			jen.Valuesln(
				jen.ID("FieldName").MapAssign().Lit(fsn),
				jen.ID("OldValue").MapAssign().ID("x").Dot(fsn),
				jen.ID("NewValue").MapAssign().ID("updated").Dot(fsn),
			),
		)

		assertions = append(assertions,
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("updated").Dot(fsn),
				jen.ID("x").Dot(fsn),
			),
		)
	}

	code.Add(
		jen.Func().IDf("Test%s_Update", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					append([]jen.Code{
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("x").Assign().AddressOf().ID(sn).Values(),
						jen.Newline(),
						jen.ID("updated").Assign().AddressOf().IDf("%sUpdateInput", sn).Valuesln(updateCols...),
						jen.Newline(),
						jen.ID("expected").Assign().Index().PointerTo().ID("FieldChangeSummary").Valuesln(
							fieldChangeSummaries...,
						),
						jen.ID("actual").Assign().ID("x").Dot("Update").Call(jen.ID("updated")),
						jen.Newline(),
						jen.List(jen.ID("expectedJSONBytes"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("expected")),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(jen.ID("t"), jen.Err()),
						jen.Newline(),
						jen.List(jen.ID("actualJSONBytes"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("actual")),
						jen.Qual(constants.MustAssertPkg, "NoError").Call(jen.ID("t"), jen.Err()),
						jen.Newline(),
						jen.List(jen.ID("expectedJSON"), jen.ID("actualJSON")).Assign().List(jen.String().Call(jen.ID("expectedJSONBytes")), jen.String().Call(jen.ID("actualJSONBytes"))),
						jen.Newline(),
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("expectedJSON"),
							jen.ID("actualJSON"),
						),
						jen.Newline(),
					},
						assertions...,
					)...,
				),
			),
		),
		jen.Newline(),
	)

	fakeFields := []jen.Code{}
	for _, field := range typ.Fields {
		fakeFields = append(fakeFields, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type, field.IsPointer)()))
	}

	code.Add(
		jen.Func().IDf("Test%sCreationInput_Validate", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("x").Assign().AddressOf().IDf("%sCreationInput", sn).Valuesln(
						fakeFields...,
					),
					jen.Newline(),
					jen.ID("actual").Assign().ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid structure"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("x").Assign().AddressOf().IDf("%sCreationInput", sn).Values(),
					jen.Newline(),
					jen.ID("actual").Assign().ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("Test%sUpdateInput_Validate", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("x").Assign().AddressOf().IDf("%sUpdateInput", sn).Valuesln(
						fakeFields...,
					),
					jen.Newline(),
					jen.ID("actual").Assign().ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with empty strings"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("x").Assign().AddressOf().IDf("%sUpdateInput", sn).Values(),
					jen.Newline(),
					jen.ID("actual").Assign().ID("x").Dot("ValidateWithContext").Call(jen.Qual("context", "Background").Call()),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
