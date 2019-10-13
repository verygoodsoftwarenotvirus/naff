package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(ret)
	ret.Add(jen.Const().Defs(
		jen.ID("webhooksBasePath").Op("=").Lit("webhooks"),
	))

	ret.Add(
		jen.Comment("BuildGetWebhookRequest builds an HTTP request for fetching a webhook"),
		jen.Line(),
		newClientMethod("BuildGetWebhookRequest").Params(
			utils.CtxParam(),
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
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment("GetWebhook retrieves a webhook"),
		jen.Line(),
		newClientMethod("GetWebhook").Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("webhook").Op("*").Qual(utils.ModelsPkg, "Webhook"),
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
			jen.Line(),
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
		jen.Comment("BuildGetWebhooksRequest builds an HTTP request for fetching webhooks"),
		jen.Line(),
		newClientMethod("BuildGetWebhooksRequest").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("webhooksBasePath"),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment("GetWebhooks gets a list of webhooks"),
		jen.Line(),
		newClientMethod("GetWebhooks").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
		).Params(
			jen.ID("webhooks").Op("*").Qual(utils.ModelsPkg, "WebhookList"),
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
			jen.Line(),
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
		jen.Comment("BuildCreateWebhookRequest builds an HTTP request for creating a webhook"),
		jen.Line(),
		newClientMethod("BuildCreateWebhookRequest").Params(
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(utils.ModelsPkg, "WebhookCreationInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("webhooksBasePath"),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
		),
	)

	ret.Add(
		jen.Comment("CreateWebhook creates a webhook"),
		jen.Line(),
		newClientMethod("CreateWebhook").Params(
			utils.CtxParam(),
			jen.ID("input").Op("*").Qual(utils.ModelsPkg, "WebhookCreationInput"),
		).Params(
			jen.ID("webhook").Op("*").Qual(utils.ModelsPkg, "Webhook"),
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
			jen.Line(),
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
		jen.Comment("BuildUpdateWebhookRequest builds an HTTP request for updating a webhook"),
		jen.Line(),
		newClientMethod("BuildUpdateWebhookRequest").Params(
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(utils.ModelsPkg, "Webhook"),
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
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID("updated"),
			),
		),
	)

	ret.Add(
		jen.Comment("UpdateWebhook updates a webhook"),
		jen.Line(),
		newClientMethod("UpdateWebhook").Params(
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(utils.ModelsPkg, "Webhook"),
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
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"), jen.Op("&").ID("updated"),
			),
		),
	)

	ret.Add(
		jen.Comment("BuildArchiveWebhookRequest builds an HTTP request for updating a webhook"),
		jen.Line(),
		newClientMethod("BuildArchiveWebhookRequest").Params(
			utils.CtxParam(),
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
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment("ArchiveWebhook archives a webhook"),
		jen.Line(),
		newClientMethod("ArchiveWebhook").Params(
			utils.CtxParam(),
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
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
	)
	return ret
}
