package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildTestClient_GetWebhook(proj)...)
	code.Add(buildTestClient_GetAllWebhooksCount(proj)...)
	code.Add(buildTestClient_GetAllWebhooks(proj)...)
	code.Add(buildTestClient_GetWebhooks(proj)...)
	code.Add(buildTestClient_CreateWebhook(proj)...)
	code.Add(buildTestClient_UpdateWebhook(proj)...)
	code.Add(buildTestClient_ArchiveWebhook(proj)...)

	return code
}

func buildTestClient_GetWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetWebhook").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhook"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("Webhook")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhook").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetAllWebhooksCount(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetAllWebhooksCount").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooksCount"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllWebhooksCount").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetAllWebhooks(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetAllWebhooks").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllWebhooks").Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_GetWebhooks(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_GetWebhooks").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "WebhookList"),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhooks"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil filter",
				utils.BuildFakeVar(proj, "WebhookList"),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhooks"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.ID(constants.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("WebhookList")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_CreateWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_CreateWebhook").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook"))),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("CreateWebhook"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("Webhook")), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Input"))),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestClient_UpdateWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_UpdateWebhook").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Var().ID("expected").Error(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("UpdateWebhook"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Webhook")),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("UpdateWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook"))),
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

func buildTestClient_ArchiveWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestClient_ArchiveWebhook").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Var().ID("expected").Error(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("ArchiveWebhook"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"), jen.ID(utils.BuildFakeVarName("Webhook")).Dot(constants.UserOwnershipFieldName)),
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
