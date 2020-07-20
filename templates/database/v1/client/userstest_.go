package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("dbclient")

	utils.AddImports(proj, code)

	code.Add(
		jen.Func().ID("TestClient_GetUser").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("User")),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_GetUserByUsername").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUserByUsername"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Username"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("User")), jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUserByUsername").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("Username")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_GetAllUserCount").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID(utils.BuildFakeVarName("Count")).Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetAllUserCount"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("Count")), jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllUserCount").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Count")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_GetUsers").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "UserList"),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUsers"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("UserList")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					constants.CtxVar(),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("UserList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil filter",
				utils.BuildFakeVar(proj, "UserList"),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("GetUsers"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("UserList")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					constants.CtxVar(),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("UserList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_CreateUser").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserDatabaseCreationInputFromUser").Call(jen.ID(utils.BuildFakeVarName("User"))),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("CreateUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("User")),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Input"))),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_UpdateUser").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Var().ID("expected").Error(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("UpdateUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")),
				).Dot("Return").Call(
					jen.ID("expected"),
				),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("UpdateUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User"))),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestClient_ArchiveUser").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("ArchiveUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
				).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("ArchiveUser").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	return code
}
