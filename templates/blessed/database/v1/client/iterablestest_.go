package client

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(buildTestClientSomethingExists(proj, typ)...)
	ret.Add(buildTestClientGetSomething(proj, typ)...)
	ret.Add(buildTestClientGetSomethingCount(proj, typ)...)
	ret.Add(buildTestClientGetAllOfSomethingCount(proj, typ)...)
	ret.Add(buildTestClientGetListOfSomething(proj, typ)...)

	// if typ.BelongsToStruct != nil {
	//		ret.Add(buildTestClientGetAllOfSomethingForSomethingElse(proj, typ)...)
	// }

	if typ.BelongsToUser {
		ret.Add(buildTestClientGetAllOfSomethingForUser(proj, typ)...)
	}

	ret.Add(buildTestClientCreateSomething(proj, typ)...)
	ret.Add(buildTestClientUpdateSomething(proj, typ)...)
	ret.Add(buildTestClientArchiveSomething(proj, typ)...)

	return ret
}

func buildRequisiteIDDeclarations(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		if varPrefix != "" {
			lines = append(lines, jen.IDf("%s%sID", varPrefix, pt.Name.Singular()).Assign().Add(utils.FakeUint64Func()))
		} else {
			lines = append(lines, jen.IDf("%sID", pt.Name.Singular()).Assign().Add(utils.FakeUint64Func()))
		}
	}

	if varPrefix != "" {
		lines = append(lines, jen.IDf("%s%sID", varPrefix, typ.Name.Singular()).Assign().Add(utils.FakeUint64Func()))
	} else {
		lines = append(lines, jen.IDf("%sID", typ.Name.UnexportedVarName()).Assign().Add(utils.FakeUint64Func()))
	}

	if typ.BelongsToUser {
		lines = append(lines, jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()))
	}

	return lines
}

func buildRequisiteIDCallArgs(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		if varPrefix != "" {
			lines = append(lines, jen.IDf("%s%sID", varPrefix, pt.Name.Singular()))
		} else {
			lines = append(lines, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
	}

	if varPrefix != "" {
		lines = append(lines, jen.IDf("%s%sID", varPrefix, typ.Name.Singular()))
	} else {
		lines = append(lines, jen.IDf("%sID", typ.Name.UnexportedVarName()))
	}

	if typ.BelongsToUser {
		lines = append(lines, jen.ID("exampleUserID"))
	}

	return lines
}

func buildTestClientSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	const varPrefix = "example"

	block := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, varPrefix, typ)...)
	mockCallArgs := []jen.Code{
		jen.Litf("%sExists", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, varPrefix, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.ID("expected").Assign().True(),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(callArgs...),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		jen.Line(),
		jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
	)

	return []jen.Code{jen.Func().IDf("TestClient_%sExists", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(block...)),
	),
		jen.Line(),
	}
}

func buildTestClientGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	const varPrefix = "example"

	block := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, varPrefix, typ)...)
	mockCallArgs := []jen.Code{
		jen.Litf("Get%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, varPrefix, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Values(),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(callArgs...),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		jen.Line(),
		jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
	)

	return []jen.Code{jen.Func().IDf("TestClient_Get%s", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(block...)),
	),
		jen.Line(),
	}
}

func buildTestClientGetSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	const (
		varPrefix     = "example"
		filterVarName = "filter"
	)

	buildSubtest := func(typ models.DataType, nilFilter bool) []jen.Code {
		lines := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, varPrefix, typ)[1:]...)

		if !nilFilter {
			lines = append(lines,
				jen.ID(filterVarName).Assign().Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(),
			)
		} else {
			lines = append(lines,
				jen.ID(filterVarName).Assign().Parens(jen.Op("*").Qual(proj.ModelsV1Package(), "QueryFilter")).Call(jen.Nil()),
			)
		}

		mockCallArgs := []jen.Code{
			jen.Litf("Get%sCount", sn),
			jen.Qual("github.com/stretchr/testify/mock", "Anything"),
		}
		idCallArgs := buildRequisiteIDCallArgs(proj, varPrefix, typ)[1:]
		mockCallArgs = append(mockCallArgs, idCallArgs...)
		callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)
		mockCallArgs = append(mockCallArgs, jen.ID(filterVarName))

		if nilFilter {
			callArgs = append(callArgs, jen.ID(filterVarName))
		} else {
			callArgs = append(callArgs, jen.ID(filterVarName))
		}

		lines = append(lines,
			jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
			jen.Line(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
			jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%sCount", sn).Call(callArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%sCount", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(buildSubtest(typ, false)...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(buildSubtest(typ, true)...)),
		),
		jen.Line(),
	}
}

func buildTestClientGetAllOfSomethingForUser(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	return []jen.Code{
		jen.Func().IDf("TestClient_GetAll%sForUser", pn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("expected").Assign().Slice().Qual(proj.ModelsV1Package(), sn).Values(),
				jen.Line(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sForUser", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAll%sForUser", pn).Call(utils.CtxVar(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	}
}

func buildTestClientGetAllOfSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	return []jen.Code{
		jen.Func().IDf("TestClient_GetAll%sCount", pn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sCount", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAll%sCount", pn).Call(utils.CtxVar()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	}
}

func buildTestClientGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	buildSubtest := func(nilFilter bool) []jen.Code {
		mockCalls := []jen.Code{
			jen.Litf("Get%s", pn),
			jen.Qual("github.com/stretchr/testify/mock", "Anything"),
		}
		const varPrefix = "example"

		idDeclarations := buildRequisiteIDDeclarations(proj, varPrefix, typ)
		idDeclarations = idDeclarations[1:]
		lines := append([]jen.Code{utils.CreateCtx()}, idDeclarations...)

		idCallArgs := buildRequisiteIDCallArgs(proj, varPrefix, typ)
		idCallArgs = idCallArgs[1:]
		callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)
		callArgs = append(callArgs, jen.ID(utils.FilterVarName))

		mockCalls = append(mockCalls, idCallArgs...)
		mockCalls = append(mockCalls, jen.ID(utils.FilterVarName))

		lines = append(lines,
			jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)).Values(),
			func() jen.Code {
				if nilFilter {
					return jen.ID(utils.FilterVarName).Assign().Add(utils.NilQueryFilter(proj))
				}
				return jen.ID(utils.FilterVarName).Assign().Add(utils.DefaultQueryFilter(proj))
			}(),
			jen.Line(),
			jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCalls...).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(callArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%s", pn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(buildSubtest(false)...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(buildSubtest(true)...)),
		),
		jen.Line(),
	}
}

func buildTestClientCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	const (
		varPrefix    = "example"
		inputVarName = varPrefix + "Input"
	)

	mockCalls := []jen.Code{
		jen.Litf("Create%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}

	idDeclarations := buildRequisiteIDDeclarations(proj, varPrefix, typ)
	idDeclarations = idDeclarations[:len(idDeclarations)-1]
	idCallArgs := buildRequisiteIDCallArgs(proj, varPrefix, typ)
	idCallArgs = idCallArgs[:len(idCallArgs)-1]

	if typ.BelongsToUser {
		idDeclarations = idDeclarations[:len(idDeclarations)-1]
		idCallArgs = idCallArgs[:len(idCallArgs)-1]
	}

	lines := append([]jen.Code{utils.CreateCtx()}, idDeclarations...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	callArgs = append(callArgs, jen.ID(inputVarName))

	mockCalls = append(mockCalls, idCallArgs...)
	mockCalls = append(mockCalls, jen.ID(inputVarName))

	lines = append(lines,
		jen.ID(inputVarName).Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn)).Values(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Values(),
		jen.Line(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(
			mockCalls...,
		).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
			callArgs...,
		),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		jen.Line(),
		jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
	)

	return []jen.Code{jen.Func().IDf("TestClient_Create%s", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
			lines...,
		)),
	),
		jen.Line(),
	}
}

func buildTestClientUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	const (
		varPrefix    = "example"
		inputVarName = varPrefix + "Input"
	)

	mockArgs := []jen.Code{
		jen.Litf("Update%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}

	idDeclarations := buildRequisiteIDDeclarations(proj, varPrefix, typ)
	idDeclarations = idDeclarations[:len(idDeclarations)-1]

	if typ.BelongsToStruct != nil {
		idDeclarations = idDeclarations[:len(idDeclarations)-1]
	}

	idCallArgs := buildRequisiteIDCallArgs(proj, varPrefix, typ)
	idCallArgs = idCallArgs[:len(idCallArgs)-1]

	if typ.BelongsToStruct != nil {
		idCallArgs = idCallArgs[:len(idCallArgs)-1]
	}
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)
	mockArgs = append(mockArgs, idCallArgs...)

	if typ.BelongsToUser {
		idDeclarations = idDeclarations[:len(idDeclarations)-1]
		callArgs = callArgs[:len(callArgs)-1]
		mockArgs = mockArgs[:len(mockArgs)-1]
	}

	lines := append([]jen.Code{utils.CreateCtx()}, idDeclarations...)

	mockArgs = append(mockArgs, jen.ID(inputVarName))
	callArgs = append(callArgs, jen.ID(inputVarName))

	lines = append(lines,
		jen.ID(inputVarName).Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Values(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.Var().ID("expected").ID("error"),
		jen.Line(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(
			mockArgs...,
		).Dot("Return").Call(jen.ID("expected")),
		jen.Line(),
		jen.Err().Assign().ID("c").Dotf("Update%s", sn).Call(
			callArgs...,
		),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Update%s", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				lines...,
			)),
		),
		jen.Line(),
	}
}

func buildTestClientArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	const varPrefix = "example"

	block := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, varPrefix, typ)...)
	mockCallArgs := []jen.Code{
		jen.Litf("Archive%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, varPrefix, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.Var().ID("expected").ID("error"),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected")),
		jen.Line(),
		jen.Err().Assign().ID("c").Dotf("Archive%s", sn).Call(callArgs...),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Archive%s", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(block...)),
		),
		jen.Line(),
	}
}
