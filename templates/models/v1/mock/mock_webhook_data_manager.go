package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockWebhookDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "WebhookDataManager").Equals().Parens(jen.PointerTo().ID("WebhookDataManager")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildMockWebhookDataManager()...)
	code.Add(buildMockGetWebhook(proj)...)
	code.Add(buildMockGetAllWebhooksCount()...)
	code.Add(buildMockGetWebhooks(proj)...)
	code.Add(buildMockGetAllWebhooks(proj)...)
	code.Add(buildMockCreateWebhook(proj)...)
	code.Add(buildMockUpdateWebhook(proj)...)
	code.Add(buildMockArchiveWebhook()...)

	return code
}

func buildMockWebhookDataManager() []jen.Code {
	lines := []jen.Code{
		jen.Comment("WebhookDataManager is a mocked models.WebhookDataManager for testing"),
		jen.Line(),
		jen.Type().ID("WebhookDataManager").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildMockGetWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID(constants.UserIDVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockGetAllWebhooksCount() []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetAllWebhooksCount satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetAllWebhooksCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockGetWebhooks(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetWebhooks satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetWebhooks").Params(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName), jen.ID(constants.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockGetAllWebhooks(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetAllWebhooks satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetAllWebhooks").Params(constants.CtxParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockCreateWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("CreateWebhook").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockUpdateWebhook(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdateWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("UpdateWebhook").Params(constants.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildMockArchiveWebhook() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("ArchiveWebhook").Params(constants.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Uint64()).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("webhookID"), jen.ID(constants.UserIDVarName)).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
