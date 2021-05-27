package stripe

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvideAPIKey").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("obligatory"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("capitalism").Dot("StripeConfig").Valuesln(
						jen.ID("APIKey").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("ProvideAPIKey").Call(jen.ID("cfg")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideWebhookSecret").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("obligatory"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("capitalism").Dot("StripeConfig").Valuesln(
						jen.ID("WebhookSecret").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("ProvideWebhookSecret").Call(jen.ID("cfg")),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
