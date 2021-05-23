package capitalism

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func paymentManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("SubscriptionPlan").Struct(
			jen.ID("ID").ID("string"),
			jen.ID("Name").ID("string"),
			jen.ID("Price").ID("uint32"),
		).Type().ID("PaymentManager").Interface(
			jen.ID("CreateCustomerID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.ID("string"), jen.ID("error")),
			jen.ID("HandleSubscriptionEventWebhook").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")),
			jen.ID("SubscribeToPlan").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("customerID"), jen.ID("paymentMethodToken"), jen.ID("planID")).ID("string")).Params(jen.ID("string"), jen.ID("error")),
			jen.ID("CreateCheckoutSession").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("subscriptionPlanID").ID("string")).Params(jen.ID("string"), jen.ID("error")),
			jen.ID("UnsubscribeFromPlan").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("subscriptionID").ID("string")).Params(jen.ID("error")),
		),
		jen.Line(),
	)

	return code
}
