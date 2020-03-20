package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockWebhookDataManagerDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookDataManager").Op("=").Parens(jen.Op("*").ID("WebhookDataManager")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("WebhookDataManager is a mocked models.WebhookDataManager for testing"),
		jen.Line(),
		jen.Type().ID("WebhookDataManager").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("webhookID"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhookCount satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhookCount").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
		).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID"), jen.ID(utils.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksCount satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooksCount").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetWebhooks satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhooks").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
			jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID"), jen.ID(utils.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooks satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooks").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllWebhooksForUser satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooksForUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("CreateWebhook").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("UpdateWebhook").Params(utils.CtxParam(), jen.ID("updated").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveWebhook satisfies our WebhookDataManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("ArchiveWebhook").Params(utils.CtxParam(), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("webhookID"), jen.ID("userID")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
