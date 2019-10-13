package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("dbclient")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"WebhookDataManager",
	).Op("=").Parens(jen.Op("*").ID("Client")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span").ID("attachWebhookIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("webhookID").ID("uint64")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("webhook_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("webhookID"), jen.Lit(10)))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetWebhook fetches a webhook from the database").Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"Webhook",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetWebhook")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValues",
		).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("webhook_id").Op(":").ID("webhookID"), jen.Lit("user_id").Op(":").ID("userID"))).Dot(
			"Debug",
		).Call(jen.Lit("GetWebhook called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetWebhook",
		).Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetWebhookCount fetches the count of webhooks from the database that meet a particular filter").Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhookCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetWebhookCount")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValues",
		).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("filter").Op(":").ID("filter"), jen.Lit("user_id").Op(":").ID("userID"))).Dot(
			"Debug",
		).Call(jen.Lit("GetWebhookCount called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetWebhookCount",
		).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter").Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllWebhooksCount")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"Debug",
		).Call(jen.Lit("GetAllWebhooksCount called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetAllWebhooksCount",
		).Call(jen.ID("ctx")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter").Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").ID("models").Dot(
		"WebhookList",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllWebhooks")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"Debug",
		).Call(jen.Lit("GetWebhookCount called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetAllWebhooks",
		).Call(jen.ID("ctx")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllWebhooksForUser fetches a list of webhooks from the database that meet a particular filter").Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllWebhooksForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("models").Dot(
		"Webhook",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetAllWebhooksForUser")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("user_id"), jen.ID("userID")).Dot(
			"Debug",
		).Call(jen.Lit("GetAllWebhooksForUser called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetAllWebhooksForUser",
		).Call(jen.ID("ctx"), jen.ID("userID")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetWebhooks fetches a list of webhooks from the database that meet a particular filter").Params(jen.ID("c").Op("*").ID("Client")).ID("GetWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"WebhookList",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("GetWebhooks")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.ID("attachFilterToSpan").Call(jen.ID("span"), jen.ID("filter")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("user_id"), jen.ID("userID")).Dot(
			"Debug",
		).Call(jen.Lit("GetWebhookCount called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"GetWebhooks",
		).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateWebhook creates a webhook in a database").Params(jen.ID("c").Op("*").ID("Client")).ID("CreateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
		"WebhookCreationInput",
	)).Params(jen.Op("*").ID("models").Dot(
		"Webhook",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("CreateWebhook")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot(
			"BelongsTo",
		)),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("user_id"), jen.ID("input").Dot(
			"BelongsTo",
		)).Dot(
			"Debug",
		).Call(jen.Lit("CreateWebhook called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"CreateWebhook",
		).Call(jen.ID("ctx"), jen.ID("input")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdateWebhook updates a particular webhook.").Comment("// NOTE: this function expects the provided input to have a non-zero ID.").Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
		"Webhook",
	)).Params(jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("UpdateWebhook")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot(
			"ID",
		)),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot(
			"BelongsTo",
		)),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("webhook_id"), jen.ID("input").Dot(
			"ID",
		)).Dot(
			"Debug",
		).Call(jen.Lit("UpdateWebhook called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"UpdateWebhook",
		).Call(jen.ID("ctx"), jen.ID("input")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveWebhook archives a webhook from the database").Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("ArchiveWebhook")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
		jen.ID("attachWebhookIDToSpan").Call(jen.ID("span"), jen.ID("webhookID")),
		jen.ID("c").Dot(
			"logger",
		).Dot(
			"WithValues",
		).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("webhook_id").Op(":").ID("webhookID"), jen.Lit("user_id").Op(":").ID("userID"))).Dot(
			"Debug",
		).Call(jen.Lit("ArchiveWebhook called")),
		jen.Return().ID("c").Dot(
			"querier",
		).Dot(
			"ArchiveWebhook",
		).Call(jen.ID("ctx"), jen.ID("webhookID"), jen.ID("userID")),
	),

		jen.Line(),
	)
	return ret
}
