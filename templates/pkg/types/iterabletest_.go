package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestSomething_Update(typ)...)
	code.Add(buildTestSomething_ToUpdateInput(proj, typ)...)

	return code
}

func buildTestSomething_Update(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	var (
		updateCols  []jen.Code
		assertCalls []jen.Code
	)

	for _, field := range typ.Fields {
		sn := field.Name.Singular()
		updateCols = append(updateCols, jen.ID(sn).MapAssign().Add(utils.FakeFuncForType(field.Type, field.Pointer)()))
		assertCalls = append(assertCalls, utils.AssertEqual(jen.ID("expected").Dot(sn), jen.ID("i").Dot(sn), nil))
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_Update", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("happy path",
				func() []jen.Code {
					lines := []jen.Code{
						jen.ID("i").Assign().AddressOf().ID(sn).Values(),
						jen.Line(),
						jen.ID("expected").Assign().AddressOf().IDf("%sUpdateInput", sn).Valuesln(updateCols...),
						jen.Line(),
						jen.ID("i").Dot("Update").Call(jen.ID("expected")),
					}
					lines = append(lines, assertCalls...)
					return lines
				}()...,
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestSomething_ToUpdateInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	updateInputFields := []jen.Code{}
	expectedFields := []jen.Code{}
	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			fns := field.Name.Singular()
			updateInputFields = append(updateInputFields, jen.ID(fns).MapAssign().Add(utils.FakeCallForField(proj.OutputPath, field)))
			expectedFields = append(expectedFields, jen.ID(fns).MapAssign().ID(uvn).Dot(field.Name.Singular()))
		}
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_ToUpdateInput", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID(uvn).Assign().AddressOf().ID(sn).Valuesln(updateInputFields...),
				jen.Line(),
				jen.ID("expected").Assign().AddressOf().IDf("%sUpdateInput", sn).Valuesln(expectedFields...),
				jen.ID("actual").Assign().ID(uvn).Dot("ToUpdateInput").Call(),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
