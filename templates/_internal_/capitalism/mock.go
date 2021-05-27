package capitalism

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("PaymentManager").Op("=").Parens(jen.Op("*").ID("MockPaymentManager")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("MockPaymentManager").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewMockPaymentManager returns a mockable capitalism.PaymentManager."),
		jen.Line(),
		jen.Func().ID("NewMockPaymentManager").Params().Params(jen.Op("*").ID("MockPaymentManager")).Body(
			jen.Return().Op("&").ID("MockPaymentManager").Values()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HandleSubscriptionEventWebhook satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockPaymentManager")).ID("HandleSubscriptionEventWebhook").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("req")).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateCustomerID satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockPaymentManager")).ID("CreateCustomerID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("account"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("String").Call(jen.Lit(0)), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ListPlans satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockPaymentManager")).ID("ListPlans").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().ID("SubscriptionPlan"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().ID("SubscriptionPlan")), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SubscribeToPlan satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockPaymentManager")).ID("SubscribeToPlan").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("customerID"), jen.ID("paymentMethodToken"), jen.ID("planID")).ID("string")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("customerID"),
				jen.ID("paymentMethodToken"),
				jen.ID("planID"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("String").Call(jen.Lit(0)), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateCheckoutSession satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockPaymentManager")).ID("CreateCheckoutSession").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("subscriptionPlanID").ID("string")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("subscriptionPlanID"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("String").Call(jen.Lit(0)), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UnsubscribeFromPlan satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockPaymentManager")).ID("UnsubscribeFromPlan").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("subscriptionID").ID("string")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("subscriptionID"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
