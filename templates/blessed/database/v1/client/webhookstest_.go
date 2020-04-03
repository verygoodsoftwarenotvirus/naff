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
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "Webhook").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
					jen.ID("exampleUserID"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhook").Call(utils.CtxVar(), jen.ID("exampleID"), jen.ID("exampleUserID")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetWebhookCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhookCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhookCount").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhookCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhookCount").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllWebhooksCount").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooksCount"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllWebhooksCount").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetAllWebhooks").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "WebhookList").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetAllWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetAllWebhooks").Call(utils.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_GetWebhooks").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "WebhookList").Values(),
				utils.CreateDefaultQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "WebhookList").Values(),
				utils.CreateNilQueryFilter(proj),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("GetWebhooks"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.ID("exampleUserID"),
					jen.ID(utils.FilterVarName),
				),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_CreateWebhook").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "WebhookCreationInput").Values(),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "Webhook").Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("CreateWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleInput"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.Err(), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_UpdateWebhook").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), "Webhook").Values(),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("UpdateWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleInput"),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("UpdateWebhook").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNoError(jen.ID("actual"), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestClient_ArchiveWebhook").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				utils.CreateCtx(),
				jen.ID("exampleID").Assign().Add(utils.FakeUint64Func()),
				jen.ID("exampleUserID").Assign().Add(utils.FakeUint64Func()),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Assign().ID("buildTestClient").Call(),
				jen.ID("mockDB").Dot("WebhookDataManager").Dot("On").Call(
					jen.Lit("ArchiveWebhook"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleID"),
					jen.ID("exampleUserID"),
				).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("actual").Assign().ID("c").Dot("ArchiveWebhook").Call(utils.CtxVar(), jen.ID("exampleID"), jen.ID("exampleUserID")),
				utils.AssertNoError(jen.ID("actual"), nil),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)
	return ret
}
