package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildTestClient_GetOAuth2Client(proj)...)
	code.Add(buildTestClient_GetOAuth2ClientByClientID(proj)...)
	code.Add(buildTestClient_GetAllOAuth2ClientCount(proj)...)
	code.Add(buildTestClient_GetOAuth2ClientsForUser(proj)...)
	code.Add(buildTestClient_CreateOAuth2Client(proj)...)
	code.Add(buildTestClient_UpdateOAuth2Client(proj)...)
	code.Add(buildTestClient_ArchiveOAuth2Client(proj)...)

	return code
}

func buildTestClient_GetOAuth2Client(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2Client").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID("expected"), constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2Client").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
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
		jen.Func().ID("TestClient_GetOAuth2ClientByClientID").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientByClientID").Call(
					constants.CtxVar(), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					constants.ObligatoryError(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientByClientID").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
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
		jen.Func().ID("TestClient_GetOAuth2ClientCount").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				utils.CreateDefaultQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientCount"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("UserID")),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("UserID")),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil filter",
				jen.ID(utils.BuildFakeVarName("UserID")).Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				utils.CreateNilQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientCount"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("UserID")),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("UserID")),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.ID(utils.BuildFakeVarName("UserID")).Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				utils.CreateDefaultQueryFilter(proj),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientCount"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("UserID")),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientCount").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("UserID")),
					jen.ID(constants.FilterVarName),
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

func buildTestClient_GetAllOAuth2ClientCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetAllOAuth2ClientCount").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID(utils.BuildFakeVarName("Count")).Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetAllOAuth2ClientCount"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(
					jen.ID(utils.BuildFakeVarName("Count")),
					jen.Nil(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllOAuth2ClientCount").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Count")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetOAuth2ClientsForUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetOAuth2ClientsForUser").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				utils.BuildFakeVar(proj, "OAuth2ClientList"),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientsForUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientsForUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil filter",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				utils.BuildFakeVar(proj, "OAuth2ClientList"),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientsForUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientsForUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID(utils.BuildFakeVarName("OAuth2ClientList")).Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList")).Call(jen.Nil()),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetOAuth2ClientsForUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), constants.ObligatoryError()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetOAuth2ClientsForUser").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), jen.ID("actual"), nil),
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
		jen.Func().ID("TestClient_CreateOAuth2Client").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Input"))),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.Line(),
				jen.ID("expected").Assign().Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(
					jen.ID("expected"),
					constants.ObligatoryError(),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateOAuth2Client").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Input"))),
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
		jen.Func().ID("TestClient_UpdateOAuth2Client").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.Var().ID("expected").Error(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("UpdateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("UpdateOAuth2Client").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
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
		jen.Func().ID("TestClient_ArchiveOAuth2Client").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.Var().ID("expected").Error(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("ArchiveOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveOAuth2Client").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				),
				utils.AssertNoError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with error returned from querier",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("expected").Assign().Qual("fmt", "Errorf").Call(jen.Lit("blah")),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("ArchiveOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveOAuth2Client").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
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
