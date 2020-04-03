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
	ret.Add(buildTestClientGetAllOfSomethingCount(proj, typ)...)
	ret.Add(buildTestClientGetListOfSomething(proj, typ)...)

	//if typ.BelongsToUser {
	//	ret.Add(buildTestClientGetAllOfSomethingForUser(proj, typ)...)
	//}

	ret.Add(buildTestClientCreateSomething(proj, typ)...)
	ret.Add(buildTestClientUpdateSomething(proj, typ)...)
	ret.Add(buildTestClientArchiveSomething(proj, typ)...)

	return ret
}

func buildRequisiteIDDeclarations(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("example%s", pt.Name.Singular()).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
	}
	lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	return lines
}

func buildRequisiteIDCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("example%s", pt.Name.Singular()).Dot("ID"))

	}
	lines = append(lines, jen.IDf("example%s", sn).Dot("ID"))

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToUser {
		lines = append(lines, jen.IDf("example%s", sn).Dot("BelongsToUser"))
	}

	return lines
}

func buildRequisiteIDCallArgsWithPreCreatedUser(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("example%s", pt.Name.Singular()).Dot("ID"))

	}
	lines = append(lines, jen.IDf("example%s", sn).Dot("ID"))

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToUser {
		lines = append(lines, jen.ID("exampleUser").Dot("ID"))
	}

	return lines
}

func buildTestClientSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	block := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, typ)...)
	mockCallArgs := []jen.Code{
		jen.Litf("%sExists", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.True(), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(callArgs...),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertTrue(jen.ID("actual"), nil),
		jen.Line(),
		jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_%sExists", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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

	block := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, typ)...)
	mockCallArgs := []jen.Code{
		jen.Litf("Get%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(callArgs...),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertEqual(jen.IDf("example%s", sn), jen.ID("actual"), nil),
		jen.Line(),
		jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%s", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"),
				jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(block...),
			),
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
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertTrue(jen.ID("actual"), nil),
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
				jen.ID("exampleCount").Assign().Uint64().Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sCount", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("exampleCount"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAll%sCount", pn).Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleCount"), jen.ID("actual"), nil, nil),
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

		idDeclarations := buildRequisiteIDDeclarations(proj, typ)
		lines := idDeclarations[1:]

		idCallArgs := buildRequisiteIDCallArgsWithPreCreatedUser(proj, typ)
		idCallArgs = idCallArgs[1:]
		callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)
		callArgs = append(callArgs, jen.ID(utils.FilterVarName))

		mockCalls = append(mockCalls, idCallArgs...)
		mockCalls = append(mockCalls, jen.ID(utils.FilterVarName))

		lines = append(lines,
			func() jen.Code {
				if nilFilter {
					return jen.ID(utils.FilterVarName).Assign().Add(utils.NilQueryFilter(proj))
				}
				return jen.ID(utils.FilterVarName).Assign().Add(utils.DefaultQueryFilter(proj))
			}(),
			jen.IDf("example%sList", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
			jen.Line(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
			jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCalls...).Dot("Return").Call(jen.IDf("example%sList", sn), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(callArgs...),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.IDf("example%sList", sn), jen.ID("actual"), nil),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%s", pn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
			jen.Line(),
			utils.BuildSubTest("obligatory", buildSubtest(false)...),
			jen.Line(),
			utils.BuildSubTest("with nil filter", buildSubtest(true)...),
		),
		jen.Line(),
	}
}

func buildTestClientCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	const (
		inputVarName = "exampleInput"
	)

	mockCalls := []jen.Code{
		jen.Litf("Create%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}

	idDeclarations := buildRequisiteIDDeclarations(proj, typ)
	idDeclarations = idDeclarations[:len(idDeclarations)-1]
	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	idCallArgs = idCallArgs[:len(idCallArgs)-1]

	if typ.BelongsToUser {
		idCallArgs = idCallArgs[:len(idCallArgs)-1]
	}

	lines := append([]jen.Code{utils.CreateCtx()}, idDeclarations...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	callArgs = append(callArgs, jen.ID(inputVarName))

	mockCalls = append(mockCalls, idCallArgs...)
	mockCalls = append(mockCalls, jen.ID(inputVarName))

	lines = append(lines,
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		jen.ID(inputVarName).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(
			mockCalls...,
		).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
			callArgs...,
		),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertEqual(jen.IDf("example%s", sn), jen.ID("actual"), nil),
		jen.Line(),
		jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Create%s", sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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

	inputVarName := fmt.Sprintf("example%s", sn)

	mockArgs := []jen.Code{
		jen.Litf("Update%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}

	idDeclarations := buildRequisiteIDDeclarations(proj, typ)

	if typ.BelongsToStruct != nil {
		idDeclarations = idDeclarations[:len(idDeclarations)-1]
	}

	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	idCallArgs = idCallArgs[:len(idCallArgs)-1]

	if typ.BelongsToStruct != nil {
		idCallArgs = idCallArgs[:len(idCallArgs)-1]
	}
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)
	mockArgs = append(mockArgs, idCallArgs...)

	if typ.BelongsToUser {
		callArgs = callArgs[:len(callArgs)-1]
		mockArgs = mockArgs[:len(mockArgs)-1]
	}

	lines := append([]jen.Code{utils.CreateCtx()}, idDeclarations...)

	mockArgs = append(mockArgs, jen.ID(inputVarName))
	callArgs = append(callArgs, jen.ID(inputVarName))

	lines = append(lines,
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
		utils.AssertNoError(jen.Err(), nil),
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

	block := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, typ)...)
	mockCallArgs := []jen.Code{
		jen.Litf("Archive%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.Var().ID("expected").ID("error"),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected")),
		jen.Line(),
		jen.Err().Assign().ID("c").Dotf("Archive%s", sn).Call(callArgs...),
		utils.AssertNoError(jen.Err(), nil),
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
