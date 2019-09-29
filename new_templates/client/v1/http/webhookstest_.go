package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

const (
	webhookRoute      = "/api/v1/webhooks/%d"
	webhooksListRoute = "/api/v1/webhooks"
)

func webhooksTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildGetWebhookRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildGetWebhookRequest").Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_GetWebhook").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "Webhook").Values(jen.Dict{
					jen.ID("ID"):   jen.Lit(1),
					jen.ID("Name"): jen.Lit("example"),
				}),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.ID("res").Qual("net/http", "ResponseWriter"),
						jen.ID("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertTrue(jen.Qual("strings", "HasSuffix").Call(
							jen.ID("req").Dot("URL").Dot("String").Call(),
							jen.Qual("strconv", "Itoa").Call(
								jen.ID("int").Call(
									jen.ID("expected").Dot("ID"),
								),
							),
						),
							nil,
						),
						assertEqual(
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit(webhookRoute),
								jen.ID("expected").Dot("ID"),
							),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.ID("req").Dot("Method"),
							jen.Qual("net/http", "MethodGet"),
							nil,
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID(t),
							jen.Qual("encoding/json", "NewEncoder").Call(
								jen.ID("res"),
							).Dot("Encode").Call(
								jen.ID("expected"),
							),
						),
					),
				),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetWebhook").Call(
					jen.ID("ctx"),
					jen.ID("expected").Dot("ID"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildGetWebhooksRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildGetWebhooksRequest").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_GetWebhooks").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "WebhookList").Values(jen.Dict{
					jen.ID("Webhooks"): jen.Index().Qual(modelsPkg, "Webhook").Values(jen.Dict{
						jen.ID("ID"):   jen.Lit(1),
						jen.ID("Name"): jen.Lit("example"),
					}),
				}),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit(webhooksListRoute),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodGet"),
								nil,
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID(t),
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(
									jen.ID("expected"),
								),
							),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetWebhooks").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildCreateWebhookRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "WebhookCreationInput").Values(jen.Dict{
					jen.ID("Name"): jen.Lit("expected name"),
				}),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildCreateWebhookRequest").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_CreateWebhook").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "Webhook").Values(jen.Dict{
					jen.ID("ID"):   jen.Lit(1),
					jen.ID("Name"): jen.Lit("example"),
				}),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "WebhookCreationInput").Values(jen.Dict{
					jen.ID("Name"): jen.ID("expected").Dot("Name"),
				}),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Lit(webhooksListRoute),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPost"),
								nil,
							),
							jen.Var().ID("x").Op("*").Qual(modelsPkg, "WebhookCreationInput"),
							jen.ID("require").Dot("NoError").Call(
								jen.ID(t),
								jen.Qual("encoding/json", "NewDecoder").Call(
									jen.ID("req").Dot("Body"),
								).Dot("Decode").Call(
									jen.Op("&").ID("x"),
								),
							),
							assertEqual(
								jen.ID("exampleInput"),
								jen.ID("x"),
								nil,
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID(t),
								jen.Qual("encoding/json", "NewEncoder").Call(
									jen.ID("res"),
								).Dot("Encode").Call(
									jen.ID("expected"),
								),
							),
							writeHeader("StatusOK"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("CreateWebhook").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildUpdateWebhookRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPut"),
				createCtx(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "Webhook").Values(jen.Dict{
					jen.ID("Name"): jen.Lit("changed name"),
				}),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildUpdateWebhookRequest").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_UpdateWebhook").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "Webhook").Values(jen.Dict{
					jen.ID("ID"):   jen.Lit(1),
					jen.ID("Name"): jen.Lit("example"),
				}),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.ID("res").Qual("net/http", "ResponseWriter"),
							jen.ID("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertEqual(
								jen.ID("req").Dot("URL").Dot("Path"),
								jen.Qual("fmt", "Sprintf").Call(
									jen.Lit(webhookRoute),
									jen.ID("expected").Dot("ID"),
								),
								jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.ID("req").Dot("Method"),
								jen.Qual("net/http", "MethodPut"),
								nil,
							),
							writeHeader("StatusOK"),
						),
					),
				),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot("UpdateWebhook").Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildArchiveWebhookRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildArchiveWebhookRequest").Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				requireNotNil(jen.ID("actual"), nil),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID(t),
					jen.ID("actual").Dot("URL"),
				),
				assertTrue(jen.Qual("strings", "HasSuffix").Call(
					jen.ID("actual").Dot("URL").Dot("String").Call(),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%d"),
						jen.ID("expectedID"),
					),
				),
					nil,
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_ArchiveWebhook").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				createCtx(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.ID("res").Qual("net/http", "ResponseWriter"),
						jen.ID("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.ID("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit(webhookRoute),
								jen.ID("expected"),
							),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.ID("req").Dot("Method"),
							jen.Qual("net/http", "MethodDelete"),
							nil,
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot("ArchiveWebhook").Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
	)
	return ret
}
