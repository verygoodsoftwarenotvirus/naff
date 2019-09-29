package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)
	ret.Add(jen.Var().ID("webhooksBasePath").Op("=").Lit("webhooks"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildGetWebhookRequest").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("GetWebhook").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("webhook").Op("*").Qual(modelsPkg, "Webhook"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetWebhookRequest").Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("webhook"),
			),
			jen.Return().List(
				jen.ID("webhook"),
				jen.ID("err"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildGetWebhooksRequest").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("webhooksBasePath"),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("GetWebhooks").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.ID("webhooks").Op("*").Qual(modelsPkg, "WebhookList"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetWebhooksRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("webhooks"),
			),
			jen.Return().List(
				jen.ID("webhooks"),
				jen.ID("err"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildCreateWebhookRequest").Params(
			ctxParam(),
			jen.ID("body").Op("*").Qual(modelsPkg, "WebhookCreationInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
			),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("CreateWebhook").Params(
			ctxParam(),
			jen.ID("input").Op("*").Qual(modelsPkg, "WebhookCreationInput"),
		).Params(
			jen.ID("webhook").Op("*").Qual(modelsPkg, "Webhook"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildCreateWebhookRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("err").Op("=").ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("webhook"),
			),
			jen.Return().List(
				jen.ID("webhook"),
				jen.ID("err"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildUpdateWebhookRequest").Params(
			ctxParam(),
			jen.ID("updated").Op("*").Qual(modelsPkg, "Webhook"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("updated").Dot("ID"),
					jen.Lit(10),
				),
			),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID("updated"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("UpdateWebhook").Params(
			ctxParam(),
			jen.ID("updated").Op("*").Qual(modelsPkg, "Webhook"),
		).Params(
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildUpdateWebhookRequest").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			), jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"), jen.Op("&").ID("updated"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildArchiveWebhookRequest").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("webhooksBasePath"), jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			), jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("ArchiveWebhook").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildArchiveWebhookRequest").Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			), jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
	)
	return ret
}
