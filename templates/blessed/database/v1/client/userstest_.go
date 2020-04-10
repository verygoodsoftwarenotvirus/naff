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
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUser").Dot("ID"),
				).Dot("Return").Call(
					jen.ID("exampleUser"),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUser").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleUser"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	ret.Add(
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
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUser").Dot("Username"),
				).Dot("Return").Call(
					jen.ID("exampleUser"), jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("Username")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleUser"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	ret.Add(
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
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("exampleUserList"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleUserList"), jen.ID("actual"), nil),
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
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("exampleUserList"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetUsers").Call(
					utils.CtxVar(),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleUserList"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_CreateUser").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "User"),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserDatabaseCreationInputFromUser").Call(jen.ID("exampleUser")),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
					jen.Lit("CreateUser"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleInput"),
				).Dot("Return").Call(
					jen.ID("exampleUser"),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleUser"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	ret.Add(
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
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUser"),
				).Dot("Return").Call(
					jen.ID("expected"),
				),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("UpdateUser").Call(utils.CtxVar(), jen.ID("exampleUser")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)

	ret.Add(
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
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUser").Dot("ID"),
				).Dot("Return").Call(jen.Nil()),
				jen.Line(),
				jen.Err().Assign().ID("c").Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("exampleUser").Dot("ID")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	)
	return ret
}
