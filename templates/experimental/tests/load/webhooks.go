package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("fetchRandomWebhook retrieves a random webhook from the list of available webhooks."),
		jen.Line(),
		jen.Func().ID("fetchRandomWebhook").Params(jen.ID("c").Op("*").ID("httpclient").Dot("Client")).Params(jen.Op("*").ID("types").Dot("Webhook")).Body(
			jen.List(jen.ID("webhooks"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
				jen.Qual("context", "Background").Call(),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("webhooks").Op("==").ID("nil").Op("||").ID("len").Call(jen.ID("webhooks").Dot("Webhooks")).Op("==").Lit(0)).Body(
				jen.Return().ID("nil")),
			jen.ID("randIndex").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.ID("webhooks").Dot("Webhooks"))),
			jen.Return().ID("webhooks").Dot("Webhooks").Index(jen.ID("randIndex")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildWebhookActions").Params(jen.ID("c").Op("*").ID("httpclient").Dot("Client"), jen.ID("builder").Op("*").ID("requests").Dot("Builder")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Body(
			jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(jen.Lit("GetWebhooks").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetWebhooks"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Return().ID("builder").Dot("BuildGetWebhooksRequest").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
			), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.If(jen.ID("randomWebhook").Op(":=").ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").Op("!=").ID("nil")).Body(
					jen.Return().ID("builder").Dot("BuildGetWebhookRequest").Call(
						jen.ID("ctx"),
						jen.ID("randomWebhook").Dot("ID"),
					)),
				jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
			), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("CreateWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("CreateWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInput").Call(),
				jen.Return().ID("builder").Dot("BuildCreateWebhookRequest").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
			), jen.ID("Weight").Op(":").Lit(1)), jen.Lit("UpdateWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("UpdateWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.If(jen.ID("randomWebhook").Op(":=").ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").Op("!=").ID("nil")).Body(
					jen.ID("randomWebhook").Dot("Name").Op("=").ID("fakes").Dot("BuildFakeWebhook").Call().Dot("Name"),
					jen.Return().ID("builder").Dot("BuildUpdateWebhookRequest").Call(
						jen.ID("ctx"),
						jen.ID("randomWebhook"),
					),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
			), jen.ID("Weight").Op(":").Lit(50)), jen.Lit("ArchiveWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("ArchiveWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.If(jen.ID("randomWebhook").Op(":=").ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").Op("!=").ID("nil")).Body(
					jen.Return().ID("builder").Dot("BuildArchiveWebhookRequest").Call(
						jen.ID("ctx"),
						jen.ID("randomWebhook").Dot("ID"),
					)),
				jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
			), jen.ID("Weight").Op(":").Lit(50)))),
		jen.Line(),
	)

	return code
}
