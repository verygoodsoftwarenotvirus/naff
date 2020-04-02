package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)
	sn := typ.Name.Singular()

	buildUpdateInputColumns := func() (updateCols []jen.Code, assertCalls []jen.Code) {
		for _, field := range typ.Fields {
			sn := field.Name.Singular()
			updateCols = append(updateCols, jen.ID(sn).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
			assertCalls = append(assertCalls, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot(sn), jen.ID("i").Dot(sn)))
		}

		return
	}

	updateCols, assertCalls := buildUpdateInputColumns()

	buildHappyPathBlock := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("i").Assign().VarPointer().ID(sn).Values(),
			jen.Line(),
			jen.ID("expected").Assign().VarPointer().IDf("%sUpdateInput", sn).Valuesln(updateCols...),
			jen.Line(),
			jen.ID("i").Dot("Update").Call(jen.ID("expected")),
		}
		lines = append(lines, assertCalls...)
		return lines
	}

	ret.Add(
		jen.Func().IDf("Test%s_Update", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(buildHappyPathBlock()...)),
		),
		jen.Line(),
	)
	return ret
}
