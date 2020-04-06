package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(buildTestClient_GetOAuth2Client(proj)...)
	ret.Add(buildTestClient_GetOAuth2ClientByClientID(proj)...)
	ret.Add(buildTestClient_GetAllOAuth2ClientCount(proj)...)
	//ret.Add(buildTestClient_GetAllOAuth2Clients(proj)...)
	ret.Add(buildTestClient_GetOAuth2Clients(proj)...)
	ret.Add(buildTestClient_CreateOAuth2Client(proj)...)
	ret.Add(buildTestClient_UpdateOAuth2Client(proj)...)
	ret.Add(buildTestClient_ArchiveOAuth2Client(proj)...)

	return ret
}

func buildTestClient_GetOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(jen.ID("exampleOAuth2Client"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(
					utils.CtxVar(),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleOAuth2Client"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(jen.ID("expected"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(
					utils.CtxVar(),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetOAuth2ClientByClientID(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2ClientByClientID").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleOAuth2Client").Dot("ClientID"),
				).Dot("Return").Call(jen.ID("exampleOAuth2Client"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientByClientID").Call(
					utils.CtxVar(), jen.ID("exampleOAuth2Client").Dot("ClientID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleOAuth2Client"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleOAuth2Client").Dot("ClientID"),
				).Dot("Return").Call(
					jen.ID("exampleOAuth2Client"),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientByClientID").Call(
					utils.CtxVar(),
					jen.ID("exampleOAuth2Client").Dot("ClientID"),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetOAuth2ClientCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2ClientCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				utils.CreateDefaultQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientCount"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			)),
		jen.Line(),
		utils.BuildSubTest(
			"with nil filter",
			jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
			jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
			utils.CreateNilQueryFilter(proj),
			jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
				jen.Lit("GetOAuth2ClientCount"),
				jen.Qual(utils.MockPkg, "Anything"),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
				utils.CtxVar(),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertExpectationsFor("mockDB"),
		),
		jen.Line(),
		utils.BuildSubTest(
			"with error returned from querier",
			jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
			jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
			utils.CreateDefaultQueryFilter(proj),
			jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
				jen.Lit("GetOAuth2ClientCount"),
				jen.Qual(utils.MockPkg, "Anything"),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			).Dot("Return").Call(jen.ID("expected"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
				utils.CtxVar(),
				jen.ID("exampleUserID"),
				jen.ID(utils.FilterVarName),
			),
			utils.AssertError(jen.Err(), nil),
			utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertExpectationsFor("mockDB"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetAllOAuth2ClientCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetAllOAuth2ClientCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleCount").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID("exampleCount"),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllOAuth2ClientCount").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleCount"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

//func buildTestClient_GetAllOAuth2Clients(proj *models.Project) []jen.Code {
//	lines := []jen.Code{
//		jen.Func().ID("TestClient_GetAllOAuth2Clients").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
//			jen.ID("T").Dot("Parallel").Call(),
//			jen.Line(),
//			utils.BuildSubTest(
//				"happy path",
//				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
//				jen.Var().ID("expected").Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
//				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(jen.Lit("GetAllOAuth2Clients"), jen.Qual(utils.MockPkg, "Anything")).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
//				jen.Line(),
//				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllOAuth2Clients").Call(utils.CtxVar()),
//				utils.AssertNoError(jen.Err(), nil),
//				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
//				jen.Line(),
//				utils.AssertExpectationsFor("mockDB"),
//			),
//		),
//		jen.Line(),
//	}
//
//	return lines
//}

func buildTestClient_GetOAuth2Clients(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2Clients").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("exampleOAuth2ClientList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2ClientList").Call(),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Clients"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("exampleOAuth2ClientList"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleOAuth2ClientList"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil filter",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("exampleOAuth2ClientList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2ClientList").Call(),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Clients"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("exampleOAuth2ClientList"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleOAuth2ClientList"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("exampleOAuth2ClientList").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList")).Call(jen.Nil()),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Clients"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("exampleOAuth2ClientList"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Clients").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleOAuth2ClientList"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_CreateOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_CreateOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.Line(),
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2ClientCreationInputFromClient").Call(jen.ID("exampleOAuth2Client")),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleInput"),
				).Dot("Return").Call(jen.ID("exampleOAuth2Client"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleOAuth2Client"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.Line(),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2ClientCreationInputFromClient").Call(jen.ID("exampleOAuth2Client")),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleInput"),
				).Dot("Return").Call(
					jen.ID("expected"),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_UpdateOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_UpdateOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.Var().ID("expected").Error(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("UpdateOAuth2Client"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleOAuth2Client"),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("UpdateOAuth2Client").Call(utils.CtxVar(), jen.ID("exampleOAuth2Client")),
				utils.AssertNoError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_ArchiveOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_ArchiveOAuth2Client").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.Var().ID("expected").Error(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("ArchiveOAuth2Client"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveOAuth2Client").Call(
					utils.CtxVar(),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				),
				utils.AssertNoError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.ID("exampleOAuth2Client").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeOAuth2Client").Call(),
				jen.Line(),
				jen.ID("expected").Assign().Qual("fmt", "Errorf").Call(jen.Lit("blah")),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("ArchiveOAuth2Client"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveOAuth2Client").Call(
					utils.CtxVar(),
					jen.ID("exampleOAuth2Client").Dot("ID"),
					jen.ID("exampleOAuth2Client").Dot("BelongsToUser"),
				),
				utils.AssertError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}
