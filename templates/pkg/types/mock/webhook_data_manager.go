package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("WebhookDataManager").Op("=").Parens(jen.Op("*").ID("WebhookDataManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("WebhookDataManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("types").Dot("Webhook"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("Webhook")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllWebhooksCount satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetWebhooks satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("WebhookList"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("filter"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("WebhookList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllWebhooks satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").ID("types").Dot("Webhook"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("results"),
				jen.ID("bucketSize"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("CreateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("WebhookCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("types").Dot("Webhook"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("createdByUser"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("Webhook")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("UpdateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("Webhook"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("types").Dot("FieldChangeSummary")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
				jen.ID("changedByUser"),
				jen.ID("changes"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveWebhook satisfies our WebhookDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("ArchiveWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID"), jen.ID("archivedByUserID")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
				jen.ID("accountID"),
				jen.ID("archivedByUserID"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntriesForWebhook is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAuditLogEntriesForWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
			),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Op("*").ID("types").Dot("AuditLogEntry")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
