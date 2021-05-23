package stripe

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func stripeDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("implementationName").Op("=").Lit("stripe_payment_manager"),
			jen.ID("webhookHeaderName").Op("=").Lit("Stripe-Signature"),
			jen.ID("webhookEventTypeCheckoutCompleted").Op("=").Lit("checkout.session.completed"),
			jen.ID("webhookEventTypeInvoicePaid").Op("=").Lit("invoice.paid"),
			jen.ID("webhookEventTypeInvoicePaymentFailed").Op("=").Lit("invoice.payment_failed"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("WebhookSecret").ID("string").Type().ID("APIKey").ID("string").Type().ID("stripePaymentManager").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
			jen.ID("client").Op("*").ID("client").Dot("API"),
			jen.ID("successURL").ID("string"),
			jen.ID("cancelURL").ID("string"),
			jen.ID("webhookSecret").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewStripePaymentManager builds a Stripe-backed stripePaymentManager."),
		jen.Line(),
		jen.Func().ID("NewStripePaymentManager").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("cfg").Op("*").ID("capitalism").Dot("StripeConfig")).Params(jen.ID("capitalism").Dot("PaymentManager")).Body(
			jen.Return().Op("&").ID("stripePaymentManager").Valuesln(
				jen.ID("client").Op(":").ID("client").Dot("New").Call(
					jen.ID("cfg").Dot("APIKey"),
					jen.ID("nil"),
				), jen.ID("webhookSecret").Op(":").ID("cfg").Dot("WebhookSecret"), jen.ID("successURL").Op(":").ID("cfg").Dot("SuccessURL"), jen.ID("cancelURL").Op(":").ID("cfg").Dot("CancelURL"), jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("implementationName")))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildCustomerName").Params(jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.ID("string")).Body(
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s (%d)"),
				jen.ID("account").Dot("Name"),
				jen.ID("account").Dot("ID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildGetCustomerParams").Params(jen.ID("a").Op("*").ID("types").Dot("Account")).Params(jen.Op("*").ID("stripe").Dot("CustomerParams")).Body(
			jen.ID("p").Op(":=").Op("&").ID("stripe").Dot("CustomerParams").Valuesln(
				jen.ID("Name").Op(":").ID("stripe").Dot("String").Call(jen.ID("buildCustomerName").Call(jen.ID("a"))), jen.ID("Email").Op(":").ID("stripe").Dot("String").Call(jen.ID("a").Dot("ContactEmail")), jen.ID("Phone").Op(":").ID("stripe").Dot("String").Call(jen.ID("a").Dot("ContactPhone")), jen.ID("Address").Op(":").Op("&").ID("stripe").Dot("AddressParams").Values()),
			jen.ID("p").Dot("AddMetadata").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("a").Dot("ExternalID"),
			),
			jen.Return().ID("p"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("stripePaymentManager")).ID("buildCheckoutSessionParams").Params(jen.ID("subscriptionPlanID").ID("string")).Params(jen.Op("*").ID("stripe").Dot("CheckoutSessionParams")).Body(
			jen.Return().Op("&").ID("stripe").Dot("CheckoutSessionParams").Valuesln(
				jen.ID("SuccessURL").Op(":").ID("stripe").Dot("String").Call(jen.ID("s").Dot("successURL")), jen.ID("CancelURL").Op(":").ID("stripe").Dot("String").Call(jen.ID("s").Dot("cancelURL")), jen.ID("Mode").Op(":").ID("stripe").Dot("String").Call(jen.ID("string").Call(jen.ID("stripe").Dot("CheckoutSessionModeSubscription"))), jen.ID("PaymentMethodTypes").Op(":").ID("stripe").Dot("StringSlice").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("card"))), jen.ID("SubscriptionData").Op(":").Op("&").ID("stripe").Dot("CheckoutSessionSubscriptionDataParams").Valuesln(
					jen.ID("Items").Op(":").Index().Op("*").ID("stripe").Dot("CheckoutSessionSubscriptionDataItemsParams").Valuesln(
						jen.Valuesln(
							jen.ID("Plan").Op(":").ID("stripe").Dot("String").Call(jen.ID("subscriptionPlanID")), jen.ID("Quantity").Op(":").ID("stripe").Dot("Int64").Call(jen.Lit(1))))))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("stripePaymentManager")).ID("CreateCheckoutSession").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("subscriptionPlanID").ID("string")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountSubscriptionPlanIDKey"),
				jen.ID("subscriptionPlanID"),
			),
			jen.ID("params").Op(":=").Op("&").ID("stripe").Dot("CheckoutSessionParams").Valuesln(
				jen.ID("SuccessURL").Op(":").ID("stripe").Dot("String").Call(jen.ID("s").Dot("successURL")), jen.ID("CancelURL").Op(":").ID("stripe").Dot("String").Call(jen.ID("s").Dot("cancelURL")), jen.ID("Mode").Op(":").ID("stripe").Dot("String").Call(jen.ID("string").Call(jen.ID("stripe").Dot("CheckoutSessionModeSubscription"))), jen.ID("PaymentMethodTypes").Op(":").ID("stripe").Dot("StringSlice").Call(jen.Index().ID("string").Valuesln(
					jen.Lit("card"))), jen.ID("SubscriptionData").Op(":").Op("&").ID("stripe").Dot("CheckoutSessionSubscriptionDataParams").Valuesln(
					jen.ID("Items").Op(":").Index().Op("*").ID("stripe").Dot("CheckoutSessionSubscriptionDataItemsParams").Valuesln(
						jen.Valuesln(
							jen.ID("Plan").Op(":").ID("stripe").Dot("String").Call(jen.ID("subscriptionPlanID")), jen.ID("Quantity").Op(":").ID("stripe").Dot("Int64").Call(jen.Lit(1)))))),
			jen.List(jen.ID("sess"), jen.ID("err")).Op(":=").ID("s").Dot("client").Dot("CheckoutSessions").Dot("New").Call(jen.ID("params")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating checkout session"),
				))),
			jen.Return().List(jen.ID("sess").Dot("ID"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("stripePaymentManager")).ID("HandleSubscriptionEventWebhook").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.List(jen.ID("b"), jen.ID("err")).Op(":=").Qual("io/ioutil", "ReadAll").Call(jen.ID("req").Dot("Body")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("parsing received webhook content"),
				)),
			jen.List(jen.ID("event"), jen.ID("err")).Op(":=").ID("webhook").Dot("ConstructEvent").Call(
				jen.ID("b"),
				jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("webhookHeaderName")),
				jen.ID("s").Dot("webhookSecret"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("constructing webhook event"),
				)),
			jen.Switch(jen.ID("event").Dot("Type")).Body(
				jen.Case(jen.ID("webhookEventTypeCheckoutCompleted")).Body(),
				jen.Case(jen.ID("webhookEventTypeInvoicePaid")).Body(),
				jen.Case(jen.ID("webhookEventTypeInvoicePaymentFailed")).Body(),
				jen.Default().Body(),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("stripePaymentManager")).ID("CreateCustomerID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("account").Dot("ID"),
			),
			jen.ID("params").Op(":=").ID("buildGetCustomerParams").Call(jen.ID("account")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot("client").Dot("Customers").Dot("New").Call(jen.ID("params")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating customer"),
				))),
			jen.Return().List(jen.ID("c").Dot("ID"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("errSubscriptionNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("subscription not found")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("stripePaymentManager")).ID("findSubscriptionID").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("customerID"), jen.ID("planID")).ID("string")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountSubscriptionPlanIDKey"),
				jen.ID("planID"),
			),
			jen.List(jen.ID("cus"), jen.ID("err")).Op(":=").ID("s").Dot("client").Dot("Customers").Dot("Get").Call(
				jen.ID("customerID"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching customer"),
				))),
			jen.For(jen.List(jen.ID("_"), jen.ID("sub")).Op(":=").Range().ID("cus").Dot("Subscriptions").Dot("Data")).Body(
				jen.If(jen.ID("sub").Dot("Plan").Dot("ID").Op("==").ID("planID")).Body(
					jen.Return().List(jen.ID("sub").Dot("ID"), jen.ID("nil")))),
			jen.Return().List(jen.Lit(""), jen.ID("errSubscriptionNotFound")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("stripePaymentManager")).ID("SubscribeToPlan").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("customerID"), jen.ID("paymentMethodToken"), jen.ID("planID")).ID("string")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountSubscriptionPlanIDKey"),
				jen.ID("planID"),
			),
			jen.List(jen.ID("subscriptionID"), jen.ID("err")).Op(":=").ID("s").Dot("findSubscriptionID").Call(
				jen.ID("ctx"),
				jen.ID("customerID"),
				jen.ID("planID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("&&").Op("!").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.ID("errSubscriptionNotFound"),
			)).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("checking subscription status"),
				))).Else().If(jen.ID("subscriptionID").Op("!=").Lit("")).Body(
				jen.Return().List(jen.ID("subscriptionID"), jen.ID("nil"))),
			jen.ID("params").Op(":=").Op("&").ID("stripe").Dot("SubscriptionParams").Valuesln(
				jen.ID("Customer").Op(":").ID("stripe").Dot("String").Call(jen.ID("customerID")), jen.ID("Plan").Op(":").ID("stripe").Dot("String").Call(jen.ID("planID")), jen.ID("DefaultSource").Op(":").ID("stripe").Dot("String").Call(jen.ID("paymentMethodToken"))),
			jen.List(jen.ID("subscription"), jen.ID("err")).Op(":=").ID("s").Dot("client").Dot("Subscriptions").Dot("New").Call(jen.ID("params")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("subscribing to plan"),
				))),
			jen.Return().List(jen.ID("subscription").Dot("ID"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildCancellationParams").Params().Params(jen.Op("*").ID("stripe").Dot("SubscriptionCancelParams")).Body(
			jen.Return().Op("&").ID("stripe").Dot("SubscriptionCancelParams").Valuesln(
				jen.ID("InvoiceNow").Op(":").ID("stripe").Dot("Bool").Call(jen.ID("true")), jen.ID("Prorate").Op(":").ID("stripe").Dot("Bool").Call(jen.ID("true")))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("stripePaymentManager")).ID("UnsubscribeFromPlan").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("subscriptionID").ID("string")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.Lit("subscription_id"),
				jen.ID("subscriptionID"),
			),
			jen.ID("params").Op(":=").ID("buildCancellationParams").Call(),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("s").Dot("client").Dot("Subscriptions").Dot("Cancel").Call(
				jen.ID("subscriptionID"),
				jen.ID("params"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("unsubscribing from plan"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
