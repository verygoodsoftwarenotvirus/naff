package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	buildUpdateInputColumns := func() (updateCols []jen.Code, assertCalls []jen.Code) {
		for _, field := range typ.Fields {
			sn := field.Name.Singular()
			updateCols = append(updateCols, jen.ID(sn).MapAssign().Add(utils.FakeFuncForType(field.Type, field.Pointer)()))
			assertCalls = append(assertCalls, utils.AssertEqual(jen.ID("expected").Dot(sn), jen.ID("i").Dot(sn), nil))
		}

		return
	}

	updateCols, assertCalls := buildUpdateInputColumns()

	buildHappyPathBlock := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("i").Assign().AddressOf().ID(sn).Values(),
			jen.Line(),
			jen.ID("expected").Assign().AddressOf().IDf("%sUpdateInput", sn).Valuesln(updateCols...),
			jen.Line(),
			jen.ID("i").Dot("Update").Call(jen.ID("expected")),
		}
		lines = append(lines, assertCalls...)
		return lines
	}

	code.Add(
		jen.Func().IDf("Test%s_Update", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("happy path", buildHappyPathBlock()...)),
		jen.Line(),
	)

	updateInputFields := []jen.Code{}
	expectedFields := []jen.Code{}
	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			fns := field.Name.Singular()
			updateInputFields = append(updateInputFields, jen.ID(fns).MapAssign().Add(utils.FakeCallForField(proj.OutputPath, field)))
			expectedFields = append(expectedFields, jen.ID(fns).MapAssign().ID(uvn).Dot(field.Name.Singular()))
		}
	}

	code.Add(
		jen.Func().IDf("Test%s_ToUpdateInput", sn).Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
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
	)

	return code
}
