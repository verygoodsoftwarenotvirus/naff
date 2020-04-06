package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockWebhookDataManagerDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "WebhookDataManager").Equals().Parens(jen.PointerTo().ID("WebhookDataManager")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("WebhookDataManager is a mocked models.WebhookDataManager for testing"),
		jen.Line(),
		jen.Type().ID("WebhookDataManager").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksCount satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetAllWebhooksCount").Params(utils.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhooks satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetWebhooks").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID"), jen.ID(utils.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooks satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("GetAllWebhooks").Params(utils.CtxParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "WebhookList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("CreateWebhook").Params(utils.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("UpdateWebhook").Params(utils.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("WebhookDataManager")).ID("ArchiveWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).Uint64()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("webhookID"), jen.ID("userID")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)
	return ret
}
