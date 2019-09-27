package client

import jen "github.com/dave/jennifer/jen"

func webhooksTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildGetWebhookRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("expectedID").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildGetWebhookRequest").Call(
					jen.Id("ctx"),
					jen.Id("expectedID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertTrue(
					jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(
						jen.Id("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.Id("expectedID"),
						),
					),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_GetWebhook").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Webhook").Values(jen.Id("ID").Op(":").Lit(1),
					jen.Id("Name").Op(":").Lit("example")),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertTrue(
							jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(
								jen.Id("req").Dot("URL").Dot("String").Call(), jen.Qual("strconv", "Itoa").Call(
									jen.Id("int").Call(
										jen.Id("expected").Dot("ID"),
									),
								),
							),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/api/v1/webhooks/%d"),
								jen.Id("expected").Dot("ID"),
							), jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("GetWebhook").Call(
					jen.Id("ctx"),
					jen.Id("expected").Dot("ID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildGetWebhooksRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildGetWebhooksRequest").Call(
					jen.Id("ctx"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_GetWebhooks").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("WebhookList").Values(jen.Id("Webhooks").Op(":").Index().Id("models").Dot("Webhook").Values(jen.Values(jen.Id("ID").Op(":").Lit(1),
					jen.Id("Name").Op(":").Lit("example"),
				),
				),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/api/v1/webhooks"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("GetWebhooks").Call(
					jen.Id("ctx"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildCreateWebhookRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("WebhookCreationInput").Values(
					jen.Id("Name").Op(":").Lit("expected name"),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildCreateWebhookRequest").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_CreateWebhook").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Webhook").Values(jen.Id("ID").Op(":").Lit(1),
					jen.Id("Name").Op(":").Lit("example")),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("WebhookCreationInput").Values(
					jen.Id("Name").Op(":").Id("expected").Dot("Name"),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"), jen.Lit("/api/v1/webhooks"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost"),
						), jen.Var().Id("x").Op("*").Id("models").Dot("WebhookCreationInput"),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewDecoder").Call(
								jen.Id("req").Dot("Body"),
							).Dot("Decode").Call(
								jen.Op("&").Id("x"),
							),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("exampleInput"),
							jen.Id("x"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("CreateWebhook").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildUpdateWebhookRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPut"),
				createCtx(),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("Webhook").Values(
					jen.Id("Name").Op(":").Lit("changed name"),
				),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildUpdateWebhookRequest").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_UpdateWebhook").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Webhook").Values(
					jen.Id("ID").Op(":").Lit(1),
					jen.Id("Name").Op(":").Lit("example"),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/api/v1/webhooks/%d"),
								jen.Id("expected").Dot("ID"),
							),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodPut"),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("err").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				).Dot("UpdateWebhook").Call(
					jen.Id("ctx"),
					jen.Id("expected"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_BuildArchiveWebhookRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("expectedID").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildArchiveWebhookRequest").Call(
					jen.Id("ctx"),
					jen.Id("expectedID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				jen.Id("require").Dot("NotNil").Call(
					jen.Id("t"),
					jen.Id("actual").Dot("URL"),
				),
				assertTrue(
					jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(
						jen.Id("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.Id("expectedID"),
						),
					),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
	)

	ret.Add(
		jen.Line(),
		testFunc("V1Client_ArchiveWebhook").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/api/v1/webhooks/%d"),
								jen.Id("expected"),
							),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"),
							jen.Qual("net/http", "MethodDelete"),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("err").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				).Dot("ArchiveWebhook").Call(
					jen.Id("ctx"),
					jen.Id("expected"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
	)
	return ret
}
