package client

import jen "github.com/dave/jennifer/jen"

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)
	ret.Add(jen.Var().Id("webhooksBasePath").Op("=").Lit("webhooks"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("BuildGetWebhookRequest").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("webhooksBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("id"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("GetWebhook").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Id("webhook").Op("*").Id("models").Dot("Webhook"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetWebhookRequest").Call(
				jen.Id("ctx"),
				jen.Id("id"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("webhook"),
			),
			jen.Return().List(
				jen.Id("webhook"),
				jen.Id("err"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildGetWebhooksRequest").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("filter").Dot("ToValues").Call(),
				jen.Id("webhooksBasePath"),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("GetWebhooks").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Id("webhooks").Op("*").Id("models").Dot("WebhookList"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetWebhooksRequest").Call(
				jen.Id("ctx"),
				jen.Id("filter"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("webhooks"),
			),
			jen.Return().List(
				jen.Id("webhooks"),
				jen.Id("err"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildCreateWebhookRequest").Params(
			ctxParam(),
			jen.Id("body").Op("*").Id("models").Dot("WebhookCreationInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("webhooksBasePath"),
			),
			jen.Return().Id("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.Id("uri"),
				jen.Id("body"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("CreateWebhook").Params(
			ctxParam(),
			jen.Id("input").Op("*").Id("models").Dot("WebhookCreationInput"),
		).Params(
			jen.Id("webhook").Op("*").Id("models").Dot("Webhook"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildCreateWebhookRequest").Call(
				jen.Id("ctx"),
				jen.Id("input"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("webhook"),
			),
			jen.Return().List(
				jen.Id("webhook"),
				jen.Id("err"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildUpdateWebhookRequest").Params(
			ctxParam(),
			jen.Id("updated").Op("*").Id("models").Dot("Webhook"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("webhooksBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("updated").Dot("ID"),
					jen.Lit(10),
				),
			),
			jen.Return().Id("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.Id("uri"),
				jen.Id("updated"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("UpdateWebhook").Params(
			ctxParam(),
			jen.Id("updated").Op("*").Id("models").Dot("Webhook"),
		).Params(
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildUpdateWebhookRequest").Call(
				jen.Id("ctx"),
				jen.Id("updated"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Id("err"),
				),
			), jen.Return().Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"), jen.Op("&").Id("updated"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildArchiveWebhookRequest").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("webhooksBasePath"), jen.Qual("strconv", "FormatUint").Call(
					jen.Id("id"),
					jen.Lit(10),
				),
			), jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("ArchiveWebhook").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildArchiveWebhookRequest").Call(
				jen.Id("ctx"),
				jen.Id("id"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Id("err"),
				),
			), jen.Return().Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Id("nil"),
			),
		),
	)
	return ret
}
