package client

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

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
	ret.Add(buildTestClientCreateSomething(proj, typ)...)
	ret.Add(buildTestClientUpdateSomething(proj, typ)...)
	ret.Add(buildTestClientArchiveSomething(proj, typ)...)

	return ret
}
func buildTestClientSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	block := typ.BuildRequisiteFakeVarsForDBClientExistenceMethodTest(proj)
	mockCallArgs := []jen.Code{
		jen.Litf("%sExists", sn),
		jen.Qual(utils.MockPkg, "Anything"),
	}
	idCallArgs := typ.BuildRequisiteFakeVarCallArgsForDBClientExistenceMethodTest(proj)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{constants.CtxVar()}, idCallArgs...)

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
	idCallArgs := typ.BuildRequisiteFakeVarCallArgsForDBClientRetrievalMethodTest(proj)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{constants.CtxVar()}, idCallArgs...)

	block := append(
		typ.BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(proj),
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
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAll%sCount", pn).Call(constants.CtxVar()),
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

		lines := typ.BuildRequisiteFakeVarsForDBClientListRetrievalMethodTest(proj)

		idCallArgs := typ.BuildCallArgsForDBClientListRetrievalMethodTest(proj)
		callArgs := append([]jen.Code{constants.CtxVar()}, idCallArgs...)
		callArgs = append(callArgs, jen.ID(constants.FilterVarName))

		mockCalls = append(mockCalls, idCallArgs...)
		mockCalls = append(mockCalls, jen.ID(constants.FilterVarName))

		lines = append(lines,
			func() jen.Code {
				if nilFilter {
					return jen.ID(constants.FilterVarName).Assign().Add(utils.NilQueryFilter(proj))
				}
				return jen.ID(constants.FilterVarName).Assign().Add(utils.DefaultQueryFilter(proj))
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

	lines := typ.BuildRequisiteFakeVarsForDBClientCreateMethodTest(proj)
	callArgs := typ.BuildCallArgsForDBClientCreationMethodTest(proj)

	mockCalls = append(mockCalls, callArgs[1:]...)

	lines = append(lines,
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

	mockArgs := []jen.Code{
		jen.Litf("Update%s", sn),
		jen.Qual(utils.MockPkg, "Anything"),
	}

	lines := typ.BuildRequisiteVarsForDBClientUpdateMethodTest(proj)

	callArgs := append(typ.BuildCallArgsForDBClientUpdateMethodTest(proj))
	mockArgs = append(mockArgs, callArgs...)

	lines = append(lines,
		jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
		jen.Var().ID("expected").Error(),
		jen.Line(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(
			mockArgs...,
		).Dot("Return").Call(jen.ID("expected")),
		jen.Line(),
		jen.Err().Assign().ID("c").Dotf("Update%s", sn).Call(
			append([]jen.Code{constants.CtxVar()}, callArgs...)...,
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

	block := typ.BuildRequisiteFakeVarsForDBClientArchiveMethodTest(proj)
	mockCallArgs := []jen.Code{
		jen.Litf("Archive%s", sn),
		jen.Qual(utils.MockPkg, "Anything"),
	}
	idCallArgs := typ.BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest(proj)
	mockCallArgs = append(mockCallArgs, idCallArgs...)
	callArgs := append([]jen.Code{constants.CtxVar()}, idCallArgs...)

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
