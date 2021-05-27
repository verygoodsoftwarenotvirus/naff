package stripe

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(
				jen.ID("ProvideAPIKey"),
				jen.ID("ProvideWebhookSecret"),
				jen.ID("NewStripePaymentManager"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAPIKey is an arbitrary wrapper for wire."),
		jen.Line(),
		jen.Func().ID("ProvideAPIKey").Params(jen.ID("cfg").Op("*").ID("capitalism").Dot("StripeConfig")).Params(jen.ID("APIKey")).Body(
			jen.Return().ID("APIKey").Call(jen.ID("cfg").Dot("APIKey"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideWebhookSecret is an arbitrary wrapper for wire."),
		jen.Line(),
		jen.Func().ID("ProvideWebhookSecret").Params(jen.ID("cfg").Op("*").ID("capitalism").Dot("StripeConfig")).Params(jen.ID("WebhookSecret")).Body(
			jen.Return().ID("WebhookSecret").Call(jen.ID("cfg").Dot("WebhookSecret"))),
		jen.Line(),
	)

	return code
}
