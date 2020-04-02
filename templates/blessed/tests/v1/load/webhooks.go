package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Comment("fetchRandomWebhook retrieves a random webhook from the list of available webhooks"),
		jen.Line(),
		jen.Func().ID("fetchRandomWebhook").Params(jen.ID("c").Op("*").Qual(pkg.HTTPClientV1Package(), "V1Client")).Params(jen.Op("*").Qual(pkg.ModelsV1Package(), "Webhook")).Block(
			jen.List(jen.ID("webhooks"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(utils.InlineCtx(), jen.Nil()),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Op("||").ID("webhooks").Op("==").ID("nil").Op("||").ID("len").Call(jen.ID("webhooks").Dot("Webhooks")).Op("==").Lit(0)).Block(
				jen.Return().ID("nil"),
			),
			jen.Line(),
			jen.ID("randIndex").Assign().Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.ID("webhooks").Dot("Webhooks"))),
			jen.Return().VarPointer().ID("webhooks").Dot("Webhooks").Index(jen.ID("randIndex")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildWebhookActions").Params(jen.ID("c").Op("*").Qual(pkg.HTTPClientV1Package(), "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
			jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(
				jen.Lit("GetWebhooks").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetWebhooks"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dot("BuildGetWebhooksRequest").Call(utils.InlineCtx(), jen.Nil()),
					),
					jen.ID("Weight").MapAssign().Lit(100),
				),
				jen.Lit("GetWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetWebhook"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomWebhook").Assign().ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildGetWebhookRequest").Call(utils.InlineCtx(), jen.ID("randomWebhook").Dot("ID")),
						),
						jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").MapAssign().Lit(100),
				),
				jen.Lit("CreateWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("CreateWebhook"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dot("BuildCreateWebhookRequest").Call(utils.InlineCtx(), jen.Qual(pkg.RandomModelsPackage(), "RandomWebhookInput").Call()),
					),
					jen.ID("Weight").MapAssign().Lit(1),
				),
				jen.Lit("UpdateWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("UpdateWebhook"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomWebhook").Assign().ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").DoesNotEqual().ID("nil")).Block(
							jen.ID("randomWebhook").Dot("Name").Equals().Qual(pkg.RandomModelsPackage(), "RandomWebhookInput").Call().Dot("Name"),
							jen.Return().ID("c").Dot("BuildUpdateWebhookRequest").Call(utils.InlineCtx(), jen.ID("randomWebhook")),
						),
						jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").MapAssign().Lit(50),
				),
				jen.Lit("ArchiveWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("ArchiveWebhook"), jen.ID("Action").MapAssign().Func().Params().Params(jen.ParamPointer().Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.ID("randomWebhook").Assign().ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildArchiveWebhookRequest").Call(utils.InlineCtx(), jen.ID("randomWebhook").Dot("ID")),
						),
						jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").MapAssign().Lit(50),
				),
			),
		),
		jen.Line(),
	)
	return ret
}
