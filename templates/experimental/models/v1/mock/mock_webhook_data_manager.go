package mock

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockWebhookDataManagerDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"WebhookDataManager",
	).Op("=").Parens(jen.Op("*").ID("WebhookDataManager")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("WebhookDataManager").Struct(jen.ID("mock").Dot(
		"Mock",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetWebhook satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"Webhook",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"Webhook",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetWebhookCount satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhookCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllWebhooksCount satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetWebhooks satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"WebhookList",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"WebhookList",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllWebhooks satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").ID("models").Dot(
		"WebhookList",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"WebhookList",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllWebhooksForUser satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("GetAllWebhooksForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("models").Dot(
		"Webhook",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("userID")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Index().ID("models").Dot(
			"Webhook",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateWebhook satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("CreateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
		"WebhookCreationInput",
	)).Params(jen.Op("*").ID("models").Dot(
		"Webhook",
	), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("input")),
		jen.Return().List(jen.ID("args").Dot(
			"Get",
		).Call(jen.Lit(0)).Assert(jen.Op("*").ID("models").Dot(
			"Webhook",
		)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdateWebhook satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("UpdateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("models").Dot(
		"Webhook",
	)).Params(jen.ID("error")).Block(
		jen.Return().ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("updated")).Dot(
			"Error",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveWebhook satisfies our WebhookDataManager interface").Params(jen.ID("m").Op("*").ID("WebhookDataManager")).ID("ArchiveWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
		jen.Return().ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")).Dot(
			"Error",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	return ret
}
