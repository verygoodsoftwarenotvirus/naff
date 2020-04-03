package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("TestClient_GetWebhook").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleWebhook").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhook").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleWebhook").Dot("ID"),
					jen.ID("exampleWebhook").Dot("BelongsToUser"),
				).Dot("Return").Call(jen.ID("exampleWebhook"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook").Dot("ID"),
					jen.ID("exampleWebhook").Dot("BelongsToUser"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleWebhook"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllWebhooksCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooksCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllWebhooks").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleWebhookList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookList").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("exampleWebhookList"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllWebhooks").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleWebhookList"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetWebhooks").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleWebhookList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookList").Call(),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("exampleWebhookList"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleWebhookList"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with nil filter",
				jen.ID("exampleWebhookList").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookList").Call(),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("exampleWebhookList"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.ID("exampleUser").Dot("ID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleWebhookList"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Func().ID("TestClient_CreateWebhook").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleWebhook").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhook").Call(),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("CreateWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleInput"),
				).Dot("Return").Call(jen.ID("exampleWebhook"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("exampleWebhook"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_UpdateWebhook").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleWebhook").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhook").Call(),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("UpdateWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleWebhook"),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID("exampleWebhook")),
				utils.AssertNoError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_ArchiveWebhook").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("exampleWebhook").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhook").Call(),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("ArchiveWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleWebhook").Dot("ID"),
					jen.ID("exampleWebhook").Dot("BelongsToUser"),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("exampleWebhook").Dot("ID"), jen.ID("exampleWebhook").Dot("BelongsToUser")),
				utils.AssertNoError(jen.ID("actual"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			),
		),
		jen.Line(),
	)
	return ret
}
