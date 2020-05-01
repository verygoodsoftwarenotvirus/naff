package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
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
						utils.FormatString("%d",
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
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
							jen.ID(constants.RequestVarName).Dot("URL").Dot("String").Call(),
							jen.Qual("strconv", "Itoa").Call(
								jen.Int().Call(
									jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
								),
							),
						),
						nil,
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						utils.FormatString(
							webhookRoute,
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Webhook"))),
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Webhook")),
					jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("GetWebhook").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
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
					constants.CtxVar(),
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
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						jen.Lit(webhooksListRoute),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("WebhookList"))),
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
					constants.CtxVar(),
					jen.Nil(),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.Err(),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("WebhookList")),
					jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("GetWebhooks").Call(constants.CtxVar(), jen.Nil()),
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
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook"))),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(
					jen.ID("t"),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.Err(),
				).Assign().ID("c").Dot("BuildCreateWebhookRequest").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
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
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook"))),
				jen.ID(utils.BuildFakeVarName("Input")).Dot(constants.UserOwnershipFieldName).Equals().Zero(),
				jen.Line(),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"), jen.Lit(webhooksListRoute), jen.Lit("expected and actual paths do not match")),
					utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodPost"), nil),
					jen.Line(),
					jen.Var().ID("x").PointerTo().Qual(proj.ModelsV1Package(), "WebhookCreationInput"),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewDecoder").Call(jen.ID(constants.RequestVarName).Dot("Body")).Dot("Decode").Call(
							jen.AddressOf().ID("x"),
						),
						nil,
					),
					utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Input")), jen.ID("x"), nil),
					jen.Line(),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Webhook"))),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("CreateWebhook").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Input")),
				),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("Webhook")), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID(utils.BuildFakeVarName("Webhook"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("CreateWebhook").Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName("Input"))),
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")),
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
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						utils.FormatString(
							webhookRoute,
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(jen.ID(constants.RequestVarName).Dot("Method"), jen.Qual("net/http", "MethodPut"), nil),
					utils.AssertNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID(constants.ResponseVarName)).Dot("Encode").Call(jen.ID(utils.BuildFakeVarName("Webhook"))),
						nil),
				),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot("UpdateWebhook").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")),
				),
				utils.AssertNoError(jen.Err(), jen.Lit("no error should be returned")),
			),
			jen.Line(),
			utils.BuildSubTest(
				"with invalid client URL",
				utils.BuildFakeVar(proj, "Webhook"),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClientWithInvalidURL").Call(jen.ID("t")).Dot("UpdateWebhook").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")),
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
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
						utils.FormatString("%d",
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
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
						jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"),
						utils.FormatString(
							webhookRoute,
							jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
						),
						jen.Lit("expected and actual paths do not match"),
					),
					utils.AssertEqual(
						jen.ID(constants.RequestVarName).Dot("Method"),
						jen.Qual("net/http", "MethodDelete"),
						nil,
					),
				),
				jen.Line(),
				jen.Err().Assign().ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot("ArchiveWebhook").Call(
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
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
					constants.CtxVar(),
					jen.ID(utils.BuildFakeVarName("Webhook")).Dot("ID"),
				),
				utils.AssertError(jen.Err(), jen.Lit("error should be returned")),
			),
		),
	)

	return ret
}
