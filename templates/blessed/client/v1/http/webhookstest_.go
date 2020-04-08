package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	webhookRoute      = "/api/v1/webhooks/%d"
	webhooksListRoute = "/api/v1/webhooks"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Line(),
		utils.OuterTestFunc("V1Client_BuildGetWebhookRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetWebhookRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("exampleWebhook").Dot("ID"),
						),
					),
					nil,
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_GetWebhook").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertTrue(
						jen.Qual("strings", "HasSuffix").Call(
							jen.ID("req").Dot("URL").Dot("String").Call(),
							jen.Qual("strconv", "Itoa").Call(
								jen.Int().Call(
									jen.ID("exampleWebhook").Dot("ID"),
								),
							),
						),
						nil,
					),
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(webhookRoute),
							jen.ID("exampleWebhook").Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleWebhook")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("GetWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("exampleWebhook"),
					jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("GetWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook").Dot("ID"),
				),
				jen.Line(),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
	)

	ret.Add(
		jen.Line(),
		utils.OuterTestFunc("V1Client_BuildGetWebhooksRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildGetWebhooksRequest").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_GetWebhooks").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "WebhookList"),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit(webhooksListRoute),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleWebhookList")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("GetWebhooks").Call(
					utils.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("exampleWebhookList"),
					jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("GetWebhooks").Call(utils.CtxVar(), jen.Nil()),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
	)

	ret.Add(
		jen.Line(),
		utils.OuterTestFunc("V1Client_BuildCreateWebhookRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.Line(),
				utils.BuildFakeVar(proj, "Webhook"),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildCreateWebhookRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_CreateWebhook").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
				jen.ID("exampleInput").Dot("BelongsToUser").Equals().Zero(),
				jen.Line(),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(jen.ID("req").Dot("URL").Dot("Path"), jen.Lit(webhooksListRoute), jen.Lit("expected and actual paths do not match")),
					utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodPost"), nil),
					jen.Line(),
					jen.Var().ID("x").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput"),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")).Dot("Decode").Call(
							jen.AddressOf().ID("x"),
						),
						nil,
					),
					utils.AssertEqual(jen.ID("exampleInput"), jen.ID("x"), nil),
					jen.Line(),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleWebhook")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleInput"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID("exampleWebhook"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.ID("exampleInput").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("CreateWebhook").Call(utils.CtxVar(), jen.ID("exampleInput")),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
	)

	ret.Add(
		jen.Line(),
		utils.OuterTestFunc("V1Client_BuildUpdateWebhookRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPut"),
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("BuildUpdateWebhookRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_UpdateWebhook").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(webhookRoute),
							jen.ID("exampleWebhook").Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodPut"), nil),
					utils.AssertNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("exampleWebhook")),
						nil),
				),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot("UpdateWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook"),
				),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("UpdateWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook"),
				),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
	)

	ret.Add(
		jen.Line(),
		utils.OuterTestFunc("V1Client_BuildArchiveWebhookRequest").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodDelete"),
				jen.ID("ts").Assign().Qual("net/http/httptest", "NewTLSServer").Call(jen.Nil()),
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildArchiveWebhookRequest").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook").Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.RequireNotNil(
					jen.ID("actual").Dot("URL"),
					nil,
				),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("exampleWebhook").Dot("ID"),
						),
					),
					nil,
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc("V1Client_ArchiveWebhook").Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(webhookRoute),
							jen.ID("exampleWebhook").Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodDelete"),
						nil,
					),
				),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot("ArchiveWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook").Dot("ID"),
				),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("ArchiveWebhook").Call(
					utils.CtxVar(),
					jen.ID("exampleWebhook").Dot("ID"),
				),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
	)
	return ret
}
