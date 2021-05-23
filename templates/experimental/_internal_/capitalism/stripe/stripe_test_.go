package stripe

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func stripeTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("fakeAPIKey").Op("=").Lit("fake_api_key"),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestPaymentManager").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("stripePaymentManager")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
			jen.ID("pm").Op(":=").ID("NewStripePaymentManager").Call(
				jen.ID("logger"),
				jen.Op("&").ID("capitalism").Dot("StripeConfig").Values(),
			),
			jen.Return().ID("pm").Assert(jen.Op("*").ID("stripePaymentManager")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestNewStripePaymentManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("pm").Op(":=").ID("NewStripePaymentManager").Call(
						jen.ID("logger"),
						jen.Op("&").ID("capitalism").Dot("StripeConfig").Values(),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("pm"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stripePaymentManager_CreateCheckoutSession").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("exampleSubscriptionPlanID").Op(":=").Lit("example_subscription_plan_id"),
					jen.ID("expected").Op(":=").Lit("example_session_id"),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.Op("&").ID("stripe").Dot("CheckoutSession").Valuesln(
							jen.ID("ID").Op(":").ID("expected")),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/v1/checkout/sessions"),
						jen.ID("exampleAPIKey"),
						jen.ID("pm").Dot("buildCheckoutSessionParams").Call(jen.ID("exampleSubscriptionPlanID")),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("CheckoutSession").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("exampleAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("CreateCheckoutSession").Call(
						jen.ID("ctx"),
						jen.ID("exampleSubscriptionPlanID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("exampleSubscriptionPlanID").Op(":=").Lit("example_subscription_plan_id"),
					jen.ID("expected").Op(":=").Lit("example_session_id"),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.Op("&").ID("stripe").Dot("CheckoutSession").Valuesln(
							jen.ID("ID").Op(":").ID("expected")),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/v1/checkout/sessions"),
						jen.ID("exampleAPIKey"),
						jen.ID("pm").Dot("buildCheckoutSessionParams").Call(jen.ID("exampleSubscriptionPlanID")),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("CheckoutSession").Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("exampleAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("CreateCheckoutSession").Call(
						jen.ID("ctx"),
						jen.ID("exampleSubscriptionPlanID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stripePaymentManager_HandleSubscriptionEventWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard with webhookEventTypeCheckoutCompleted"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("pm").Dot("webhookSecret").Op("=").Lit("example_webhook_secret"),
					jen.ID("testEventTypes").Op(":=").Index().ID("string").Valuesln(
						jen.ID("webhookEventTypeCheckoutCompleted"), jen.ID("webhookEventTypeInvoicePaid"), jen.ID("webhookEventTypeInvoicePaymentFailed"), jen.ID("t").Dot("Name").Call()),
					jen.For(jen.List(jen.ID("_"), jen.ID("et")).Op(":=").Range().ID("testEventTypes")).Body(
						jen.ID("exampleEvent").Op(":=").Op("&").ID("stripe").Dot("Event").Valuesln(
							jen.ID("Account").Op(":").Lit("whatever"), jen.ID("Created").Op(":").Qual("time", "Now").Call().Dot("Unix").Call(), jen.ID("Data").Op(":").Op("&").ID("stripe").Dot("EventData").Valuesln(
								jen.ID("Object").Op(":").Map(jen.ID("string")).Interface().Valuesln(
									jen.Lit("things").Op(":").Lit("stuff")), jen.ID("PreviousAttributes").Op(":").Map(jen.ID("string")).Interface().Valuesln(
									jen.Lit("things").Op(":").Lit("stuff"))), jen.ID("ID").Op(":").Lit("example"), jen.ID("Type").Op(":").ID("et")),
						jen.Var().ID("b").Qual("bytes", "Buffer"),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b")).Dot("Encode").Call(jen.ID("exampleEvent")),
						),
						jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
							jen.Qual("crypto/sha256", "New"),
							jen.Index().ID("byte").Call(jen.ID("pm").Dot("webhookSecret")),
						),
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("b").Dot("Bytes").Call()),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("now").Op(":=").Qual("time", "Now").Call(),
						jen.ID("exampleSig").Op(":=").ID("webhook").Dot("ComputeSignature").Call(
							jen.ID("now"),
							jen.ID("b").Dot("Bytes").Call(),
							jen.ID("pm").Dot("webhookSecret"),
						),
						jen.ID("exampleSignature").Op(":=").Qual("fmt", "Sprintf").Call(
							jen.Lit("t=%d,v1=%s"),
							jen.ID("now").Dot("Unix").Call(),
							jen.Qual("encoding/hex", "EncodeToString").Call(jen.ID("exampleSig")),
						),
						jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
							jen.Qual("net/http", "MethodPost"),
							jen.Lit("/webhook_update"),
							jen.Qual("bytes", "NewReader").Call(jen.ID("b").Dot("Bytes").Call()),
						),
						jen.ID("req").Dot("Header").Dot("Set").Call(
							jen.ID("webhookHeaderName"),
							jen.ID("exampleSignature"),
						),
						jen.ID("err").Op("=").ID("pm").Dot("HandleSubscriptionEventWebhook").Call(jen.ID("req")),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid body"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("pm").Dot("webhookSecret").Op("=").Lit("example_webhook_secret"),
					jen.ID("exampleEvent").Op(":=").Op("&").ID("stripe").Dot("Event").Valuesln(
						jen.ID("Account").Op(":").Lit("whatever"), jen.ID("Created").Op(":").Qual("time", "Now").Call().Dot("Unix").Call(), jen.ID("Data").Op(":").Op("&").ID("stripe").Dot("EventData").Valuesln(
							jen.ID("Object").Op(":").Map(jen.ID("string")).Interface().Valuesln(
								jen.Lit("things").Op(":").Lit("stuff")), jen.ID("PreviousAttributes").Op(":").Map(jen.ID("string")).Interface().Valuesln(
								jen.Lit("things").Op(":").Lit("stuff"))), jen.ID("ID").Op(":").Lit("example"), jen.ID("Type").Op(":").ID("webhookEventTypeCheckoutCompleted")),
					jen.Var().ID("b").Qual("bytes", "Buffer"),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b")).Dot("Encode").Call(jen.ID("exampleEvent")),
					),
					jen.ID("mac").Op(":=").Qual("crypto/hmac", "New").Call(
						jen.Qual("crypto/sha256", "New"),
						jen.Index().ID("byte").Call(jen.ID("pm").Dot("webhookSecret")),
					),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("mac").Dot("Write").Call(jen.ID("b").Dot("Bytes").Call()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mrc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockReadCloser").Values(),
					jen.ID("mrc").Dot("On").Call(
						jen.Lit("Read"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("byte").Call(jen.Lit(""))),
					).Dot("Return").Call(
						jen.ID("int").Call(jen.Lit(0)),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/webhook_update"),
						jen.ID("mrc"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.ID("webhookHeaderName"),
						jen.Lit("bad-sig"),
					),
					jen.ID("err").Op("=").ID("pm").Dot("HandleSubscriptionEventWebhook").Call(jen.ID("req")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid signature"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/webhook_update"),
						jen.ID("nil"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.ID("webhookHeaderName"),
						jen.Lit("bad-sig"),
					),
					jen.ID("err").Op(":=").ID("pm").Dot("HandleSubscriptionEventWebhook").Call(jen.ID("req")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stripePaymentManager_GetCustomerID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").Lit("fake_customer_id"),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.Op("&").ID("stripe").Dot("Customer").Valuesln(
							jen.ID("ID").Op(":").ID("expected")),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/v1/customers"),
						jen.ID("exampleAPIKey"),
						jen.ID("buildGetCustomerParams").Call(jen.ID("exampleAccount")),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Customer").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("exampleAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("CreateCustomerID").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").Lit("fake_customer_id"),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.Op("&").ID("stripe").Dot("Customer").Valuesln(
							jen.ID("ID").Op(":").ID("expected")),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/v1/customers"),
						jen.ID("exampleAPIKey"),
						jen.ID("buildGetCustomerParams").Call(jen.ID("exampleAccount")),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Customer").Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("exampleAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("CreateCustomerID").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stripePaymentManager_findSubscriptionID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("exampleCustomerID").Op(":=").Lit("fake_customer"),
					jen.ID("examplePlanID").Op(":=").Lit("fake_plan"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("expected").Op(":=").Lit("fake_subscription_id"),
					jen.ID("expectedCustomer").Op(":=").Op("&").ID("stripe").Dot("Customer").Valuesln(
						jen.ID("ID").Op(":").ID("exampleCustomerID"), jen.ID("Subscriptions").Op(":").Op("&").ID("stripe").Dot("SubscriptionList").Valuesln(
							jen.ID("Data").Op(":").Index().Op("*").ID("stripe").Dot("Subscription").Valuesln(
								jen.Valuesln(
									jen.ID("ID").Op(":").ID("expected"), jen.ID("Plan").Op(":").Op("&").ID("stripe").Dot("Plan").Valuesln(
										jen.ID("ID").Op(":").ID("examplePlanID")))))),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedCustomer"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/v1/customers/%s"),
							jen.ID("exampleCustomerID"),
						),
						jen.ID("exampleAPIKey"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("CustomerParams").Values()),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Customer").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("fakeAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("findSubscriptionID").Call(
						jen.ID("ctx"),
						jen.ID("exampleCustomerID"),
						jen.ID("examplePlanID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stripePaymentManager_SubscribeToPlan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with pre-existing subscription"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("exampleCustomerID").Op(":=").Lit("fake_customer"),
					jen.ID("examplePlanID").Op(":=").Lit("fake_plan"),
					jen.ID("examplePaymentMethodToken").Op(":=").Lit("fake_payment_token"),
					jen.ID("expected").Op(":=").Lit("fake_subscription"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("expectedCustomer").Op(":=").Op("&").ID("stripe").Dot("Customer").Valuesln(
						jen.ID("ID").Op(":").ID("exampleCustomerID"), jen.ID("Subscriptions").Op(":").Op("&").ID("stripe").Dot("SubscriptionList").Valuesln(
							jen.ID("Data").Op(":").Index().Op("*").ID("stripe").Dot("Subscription").Valuesln(
								jen.Valuesln(
									jen.ID("ID").Op(":").ID("expected"), jen.ID("Plan").Op(":").Op("&").ID("stripe").Dot("Plan").Valuesln(
										jen.ID("ID").Op(":").ID("examplePlanID")))))),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedCustomer"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/v1/customers/%s"),
							jen.ID("exampleCustomerID"),
						),
						jen.ID("exampleAPIKey"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("CustomerParams").Values()),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Customer").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("fakeAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("SubscribeToPlan").Call(
						jen.ID("ctx"),
						jen.ID("exampleCustomerID"),
						jen.ID("examplePaymentMethodToken"),
						jen.ID("examplePlanID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error checking pre-existing subscription"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("exampleCustomerID").Op(":=").Lit("fake_customer"),
					jen.ID("examplePlanID").Op(":=").Lit("fake_plan"),
					jen.ID("examplePaymentMethodToken").Op(":=").Lit("fake_payment_token"),
					jen.ID("expected").Op(":=").Lit("fake_subscription"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("expectedCustomer").Op(":=").Op("&").ID("stripe").Dot("Customer").Valuesln(
						jen.ID("ID").Op(":").ID("exampleCustomerID"), jen.ID("Subscriptions").Op(":").Op("&").ID("stripe").Dot("SubscriptionList").Valuesln(
							jen.ID("Data").Op(":").Index().Op("*").ID("stripe").Dot("Subscription").Valuesln(
								jen.Valuesln(
									jen.ID("ID").Op(":").ID("expected"), jen.ID("Plan").Op(":").Op("&").ID("stripe").Dot("Plan").Valuesln(
										jen.ID("ID").Op(":").ID("examplePlanID")))))),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedCustomer"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/v1/customers/%s"),
							jen.ID("exampleCustomerID"),
						),
						jen.ID("exampleAPIKey"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("CustomerParams").Values()),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Customer").Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("fakeAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("SubscribeToPlan").Call(
						jen.ID("ctx"),
						jen.ID("exampleCustomerID"),
						jen.ID("examplePaymentMethodToken"),
						jen.ID("examplePlanID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without pre-existing subscription"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("exampleCustomerID").Op(":=").Lit("fake_customer"),
					jen.ID("exampleSubscriptionID").Op(":=").Lit("fake_subscription"),
					jen.ID("examplePlanID").Op(":=").Lit("fake_plan"),
					jen.ID("examplePaymentMethodToken").Op(":=").Lit("fake_payment_token"),
					jen.ID("expected").Op(":=").Lit("fake_subscription"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("expectedCustomer").Op(":=").Op("&").ID("stripe").Dot("Customer").Valuesln(
						jen.ID("ID").Op(":").ID("exampleCustomerID"), jen.ID("Subscriptions").Op(":").Op("&").ID("stripe").Dot("SubscriptionList").Valuesln(
							jen.ID("Data").Op(":").Index().Op("*").ID("stripe").Dot("Subscription").Values())),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedCustomer"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/v1/customers/%s"),
							jen.ID("exampleCustomerID"),
						),
						jen.ID("exampleAPIKey"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("CustomerParams").Values()),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Customer").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("expectedSubscription").Op(":=").Op("&").ID("stripe").Dot("Subscription").Valuesln(
						jen.ID("ID").Op(":").ID("exampleSubscriptionID")),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedSubscription"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/v1/subscriptions"),
						jen.ID("exampleAPIKey"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("SubscriptionParams").Values()),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Subscription").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("fakeAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("SubscribeToPlan").Call(
						jen.ID("ctx"),
						jen.ID("exampleCustomerID"),
						jen.ID("examplePaymentMethodToken"),
						jen.ID("examplePlanID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without pre-existing subscription and with error creating subscription"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("exampleCustomerID").Op(":=").Lit("fake_customer"),
					jen.ID("exampleSubscriptionID").Op(":=").Lit("fake_subscription"),
					jen.ID("examplePlanID").Op(":=").Lit("fake_plan"),
					jen.ID("examplePaymentMethodToken").Op(":=").Lit("fake_payment_token"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("expectedCustomer").Op(":=").Op("&").ID("stripe").Dot("Customer").Valuesln(
						jen.ID("ID").Op(":").ID("exampleCustomerID"), jen.ID("Subscriptions").Op(":").Op("&").ID("stripe").Dot("SubscriptionList").Valuesln(
							jen.ID("Data").Op(":").Index().Op("*").ID("stripe").Dot("Subscription").Values())),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedCustomer"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodGet"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/v1/customers/%s"),
							jen.ID("exampleCustomerID"),
						),
						jen.ID("exampleAPIKey"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("CustomerParams").Values()),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Customer").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("expectedSubscription").Op(":=").Op("&").ID("stripe").Dot("Subscription").Valuesln(
						jen.ID("ID").Op(":").ID("exampleSubscriptionID")),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedSubscription"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/v1/subscriptions"),
						jen.ID("exampleAPIKey"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("SubscriptionParams").Values()),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Subscription").Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("fakeAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("pm").Dot("SubscribeToPlan").Call(
						jen.ID("ctx"),
						jen.ID("exampleCustomerID"),
						jen.ID("examplePaymentMethodToken"),
						jen.ID("examplePlanID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stripePaymentManager_UnsubscribeFromPlan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("exampleSubscriptionID").Op(":=").Lit("fake_subscription_id"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("expectedCustomer").Op(":=").Op("&").ID("stripe").Dot("Customer").Valuesln(
						jen.ID("Subscriptions").Op(":").Op("&").ID("stripe").Dot("SubscriptionList").Valuesln(
							jen.ID("Data").Op(":").Index().Op("*").ID("stripe").Dot("Subscription").Valuesln(
								jen.Valuesln(
									jen.ID("ID").Op(":").ID("exampleSubscriptionID"))))),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedCustomer"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodDelete"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/v1/subscriptions/%s"),
							jen.ID("exampleSubscriptionID"),
						),
						jen.ID("exampleAPIKey"),
						jen.ID("buildCancellationParams").Call(),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Subscription").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("fakeAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.ID("err").Op(":=").ID("pm").Dot("UnsubscribeFromPlan").Call(
						jen.ID("ctx"),
						jen.ID("exampleSubscriptionID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAPIKey").Op(":=").ID("fakeAPIKey"),
					jen.ID("exampleSubscriptionID").Op(":=").Lit("fake_subscription_id"),
					jen.ID("pm").Op(":=").ID("buildTestPaymentManager").Call(jen.ID("t")),
					jen.ID("mockAPIBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockConnectBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("mockUploadsBackend").Op(":=").Op("&").ID("mockBackend").Values(),
					jen.ID("expectedCustomer").Op(":=").Op("&").ID("stripe").Dot("Customer").Valuesln(
						jen.ID("Subscriptions").Op(":").Op("&").ID("stripe").Dot("SubscriptionList").Valuesln(
							jen.ID("Data").Op(":").Index().Op("*").ID("stripe").Dot("Subscription").Valuesln(
								jen.Valuesln(
									jen.ID("ID").Op(":").ID("exampleSubscriptionID"))))),
					jen.ID("mockAPIBackend").Dot("AnticipateCall").Call(
						jen.ID("t"),
						jen.ID("expectedCustomer"),
					),
					jen.ID("mockAPIBackend").Dot("On").Call(
						jen.Lit("Call"),
						jen.Qual("net/http", "MethodDelete"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/v1/subscriptions/%s"),
							jen.ID("exampleSubscriptionID"),
						),
						jen.ID("exampleAPIKey"),
						jen.ID("buildCancellationParams").Call(),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("stripe").Dot("Subscription").Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("mockedBackends").Op(":=").Op("&").ID("stripe").Dot("Backends").Valuesln(
						jen.ID("API").Op(":").ID("mockAPIBackend"), jen.ID("Connect").Op(":").ID("mockConnectBackend"), jen.ID("Uploads").Op(":").ID("mockUploadsBackend")),
					jen.ID("pm").Dot("client").Op("=").ID("client").Dot("New").Call(
						jen.ID("fakeAPIKey"),
						jen.ID("mockedBackends"),
					),
					jen.ID("err").Op(":=").ID("pm").Dot("UnsubscribeFromPlan").Call(
						jen.ID("ctx"),
						jen.ID("exampleSubscriptionID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockAPIBackend"),
						jen.ID("mockConnectBackend"),
						jen.ID("mockUploadsBackend"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
