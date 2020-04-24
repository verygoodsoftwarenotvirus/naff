package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Comment("fetchRandomWebhook retrieves a random webhook from the list of available webhooks"),
		jen.Line(),
		jen.Func().ID("fetchRandomWebhook").Params(jen.ID("c").PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "Webhook")).Block(
			jen.List(jen.ID("webhooks"), jen.Err()).Assign().ID("c").Dot("GetWebhooks").Call(constants.InlineCtx(), jen.Nil()),
			jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID("webhooks").IsEqualTo().ID("nil").Or().ID("len").Call(jen.ID("webhooks").Dot("Webhooks")).IsEqualTo().Zero()).Block(
				jen.Return().ID("nil"),
			),
			jen.Line(),
			jen.ID("randIndex").Assign().Qual("math/rand", "Intn").Call(jen.Len(jen.ID("webhooks").Dot("Webhooks"))),
			jen.Return().AddressOf().ID("webhooks").Dot("Webhooks").Index(jen.ID("randIndex")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildWebhookActions").Params(jen.ID("c").PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Params(jen.Map(jen.String()).PointerTo().ID("Action")).Block(
			jen.Return().Map(jen.String()).PointerTo().ID("Action").Valuesln(
				jen.Lit("GetWebhooks").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetWebhooks"),
					jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						constants.CreateCtx(),
						jen.Return().ID("c").Dot("BuildGetWebhooksRequest").Call(constants.CtxVar(), jen.Nil()),
					),
					jen.ID("Weight").MapAssign().Lit(100),
				),
				jen.Lit("GetWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("GetWebhook"),
					jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						constants.CreateCtx(),
						jen.If(jen.ID("randomWebhook").Assign().ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildGetWebhookRequest").Call(constants.CtxVar(), jen.ID("randomWebhook").Dot("ID")),
						),
						jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").MapAssign().Lit(100),
				),
				jen.Lit("CreateWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("CreateWebhook"),
					jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						constants.CreateCtx(),
						utils.BuildFakeVarWithCustomName(
							proj,
							"exampleInput",
							"BuildFakeWebhookCreationInput",
						),
						jen.Return().ID("c").Dot("BuildCreateWebhookRequest").Call(constants.CtxVar(), jen.ID("exampleInput")),
					),
					jen.ID("Weight").MapAssign().One(),
				),
				jen.Lit("UpdateWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("UpdateWebhook"),
					jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						constants.CreateCtx(),
						jen.If(jen.ID("randomWebhook").Assign().ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").DoesNotEqual().ID("nil")).Block(
							jen.ID("randomWebhook").Dot("Name").Equals().Qual(proj.FakeModelsPackage(), "BuildFakeWebhook").Call().Dot("Name"),
							jen.Return().ID("c").Dot("BuildUpdateWebhookRequest").Call(constants.CtxVar(), jen.ID("randomWebhook")),
						),
						jen.Return().List(jen.Nil(), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").MapAssign().Lit(50),
				),
				jen.Lit("ArchiveWebhook").MapAssign().Valuesln(
					jen.ID("Name").MapAssign().Lit("ArchiveWebhook"),
					jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Block(
						constants.CreateCtx(),
						jen.If(jen.ID("randomWebhook").Assign().ID("fetchRandomWebhook").Call(jen.ID("c")), jen.ID("randomWebhook").DoesNotEqual().ID("nil")).Block(
							jen.Return().ID("c").Dot("BuildArchiveWebhookRequest").Call(constants.CtxVar(), jen.ID("randomWebhook").Dot("ID")),
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
