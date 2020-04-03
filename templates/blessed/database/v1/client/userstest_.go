package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestClient_GetUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("GetUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleID")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(utils.CtxVar(), jen.ID("exampleID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetUserByUsername").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleUsername").Assign().Lit("username"),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("GetUserByUsername"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleUsername")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("exampleUsername")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllUserCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetAllUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllUserCount").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil filter",
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetAllUserCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllUserCount").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetUsers").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserList").Values(),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUsers"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil filter",
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserList").Values(),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUsers"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_CreateUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "UserInput").Values(),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("CreateUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_UpdateUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("UpdateUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("UpdateUser").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_ArchiveUser").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleInput").Assign().Add(utils.FakeUint64Func()),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(jen.Lit("ArchiveUser"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)
	return ret
}
