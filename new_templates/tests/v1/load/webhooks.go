package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func webhooksDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Func())
	ret.Add(jen.Func().ID("buildWebhookActions").Params(jen.ID("c").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
		jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(jen.Lit("GetWebhooks").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetWebhooks"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.Return().ID("c").Dot(
				"BuildGetWebhooksRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("GetWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.If(
				jen.ID("randomWebhook").Op(":=").ID("fetchRandomWebhook").Call(jen.ID("c")),
				jen.ID("randomWebhook").Op("!=").ID("nil"),
			).Block(
				jen.Return().ID("c").Dot(
					"BuildGetWebhookRequest",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("randomWebhook").Dot(
					"ID",
				)),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("CreateWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("CreateWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.Return().ID("c").Dot(
				"BuildCreateWebhookRequest",
			).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomWebhookInput").Call()),
		), jen.ID("Weight").Op(":").Lit(1)), jen.Lit("UpdateWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("UpdateWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.If(
				jen.ID("randomWebhook").Op(":=").ID("fetchRandomWebhook").Call(jen.ID("c")),
				jen.ID("randomWebhook").Op("!=").ID("nil"),
			).Block(
				jen.ID("randomWebhook").Dot(
					"Name",
				).Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomWebhookInput").Call().Dot(
					"Name",
				),
				jen.Return().ID("c").Dot(
					"BuildUpdateWebhookRequest",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("randomWebhook")),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
		), jen.ID("Weight").Op(":").Lit(50)), jen.Lit("ArchiveWebhook").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("ArchiveWebhook"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.If(
				jen.ID("randomWebhook").Op(":=").ID("fetchRandomWebhook").Call(jen.ID("c")),
				jen.ID("randomWebhook").Op("!=").ID("nil"),
			).Block(
				jen.Return().ID("c").Dot(
					"BuildArchiveWebhookRequest",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("randomWebhook").Dot(
					"ID",
				)),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
		), jen.ID("Weight").Op(":").Lit(50))),
	),
	)
	return ret
}
