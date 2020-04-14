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
		lines = append(lines, jen.ID(utils.BuildFakeVarName(pt.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
	}
	lines = append(lines, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	return lines
}

func buildRequisiteIDCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.ID(utils.BuildFakeVarName(pt.Name.Singular())).Dot("ID"))

	}
	lines = append(lines, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	//if typ.BelongsToStruct != nil {
	//	lines = append(lines, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	//}

	if typ.BelongsToUser {
		lines = append(lines, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
	}

	return lines
}

func buildRequisiteIDCallArgsWithPreCreatedUser(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}
	sn := typ.Name.Singular()

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.ID(utils.BuildFakeVarName(pt.Name.Singular())).Dot("ID"))

	}
	lines = append(lines, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToUser {
		lines = append(lines, jen.ID(utils.BuildFakeVarName("User")).Dot("ID"))
	}

	return lines
}

func buildTestClientSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	block := append([]jen.Code{utils.CreateCtx()}, buildRequisiteIDDeclarations(proj, typ)...)
	mockCallArgs := []jen.Code{
		jen.Litf("%sExists", sn),
		jen.Qual(utils.MockPkg, "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.True(), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(callArgs...),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertTrue(jen.ID("actual"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("mockDB"),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_%sExists", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("obligatory", block...)),
		jen.Line(),
	}
}

func buildTestClientGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	mockCallArgs := []jen.Code{
		jen.Litf("Get%s", sn),
		jen.Qual(utils.MockPkg, "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block := append(
		buildRequisiteIDDeclarations(proj, typ),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(callArgs...),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertEqual(jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("mockDB"),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest("obligatory", block...),
		),
		jen.Line(),
	}
}

func buildTestClientGetAllOfSomethingForUser(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	return []jen.Code{
		jen.Func().IDf("TestClient_GetAll%sForUser", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.ID(utils.BuildFakeVarName("UserID")).Assign().Add(utils.FakeUint64Func()),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("expected").Assign().Slice().Qual(proj.ModelsV1Package(), sn).Values(),
				jen.Line(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sForUser", pn), jen.Qual(utils.MockPkg, "Anything"), jen.ID(utils.BuildFakeVarName("UserID"))).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAll%sForUser", pn).Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName("UserID"))),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertTrue(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}
}

func buildTestClientGetAllOfSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	return []jen.Code{
		jen.Func().IDf("TestClient_GetAll%sCount", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"obligatory",
				jen.ID(utils.BuildFakeVarName("Count")).Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sCount", pn), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("Count")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAll%sCount", pn).Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Count")), jen.ID("actual"), nil, nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
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
			jen.Qual(utils.MockPkg, "Anything"),
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
			utils.AssertExpectationsFor("mockDB"),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%s", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			func() jen.Code {
				if typ.BelongsToUser && typ.RestrictedToUser {
					return utils.BuildFakeVar(proj, "User")
				}
				return jen.Null()
			}(),
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
		jen.Qual(utils.MockPkg, "Anything"),
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
		jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		jen.ID(inputVarName).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.ID(utils.BuildFakeVarName(sn))),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(
			mockCalls...,
		).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)), jen.Nil()),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
			callArgs...,
		),
		utils.AssertNoError(jen.Err(), nil),
		utils.AssertEqual(jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual"), nil),
		jen.Line(),
		utils.AssertExpectationsFor("mockDB"),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Create%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("obligatory", lines...),
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
		jen.Qual(utils.MockPkg, "Anything"),
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
		jen.Var().ID("expected").Error(),
		jen.Line(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(
			mockArgs...,
		).Dot("Return").Call(jen.ID("expected")),
		jen.Line(),
		jen.Err().Assign().ID("c").Dotf("Update%s", sn).Call(
			callArgs...,
		),
		utils.AssertNoError(jen.Err(), nil),
		jen.Line(),
		utils.AssertExpectationsFor("mockDB"),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Update%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("obligatory", lines...),
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
		jen.Qual(utils.MockPkg, "Anything"),
	}
	idCallArgs := buildRequisiteIDCallArgs(proj, typ)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{utils.CtxVar()}, idCallArgs...)

	block = append(block,
		jen.Var().ID("expected").Error(),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected")),
		jen.Line(),
		jen.Err().Assign().ID("c").Dotf("Archive%s", sn).Call(callArgs...),
		utils.AssertNoError(jen.Err(), nil),
		jen.Line(),
		utils.AssertExpectationsFor("mockDB"),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Archive%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("obligatory", block...)),
		jen.Line(),
	}
}
