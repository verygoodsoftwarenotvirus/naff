package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestWebhooks").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.ID("webhooksTestSuite")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("webhooksTestSuite").Struct(
				jen.ID("suite").Dot("Suite"),
				jen.ID("ctx").Qual("context", "Context"),
				jen.ID("exampleWebhook").Op("*").ID("types").Dot("Webhook"),
				jen.ID("exampleWebhookList").Op("*").ID("types").Dot("WebhookList"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").ID("webhooksTestSuite")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("webhooksTestSuite")).ID("SetupTest").Params().Body(
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dot("exampleWebhook").Op("=").ID("fakes").Dot("BuildFakeWebhook").Call(),
			jen.ID("s").Dot("exampleWebhookList").Op("=").ID("fakes").Dot("BuildFakeWebhookList").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("webhooksTestSuite")).ID("TestClient_GetWebhook").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/webhooks/%d"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleWebhook"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleWebhook"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid webhook ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("webhooksTestSuite")).ID("TestClient_GetWebhooks").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/webhooks"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleWebhookList"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleWebhookList"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("webhooksTestSuite")).ID("TestClient_CreateWebhook").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/webhooks"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("s").Dot("exampleWebhook")),
					jen.ID("exampleInput").Dot("BelongsToAccount").Op("=").Lit(0),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithRequestBodyValidation").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Op("&").ID("types").Dot("WebhookCreationInput").Valuesln(),
						jen.ID("exampleInput"),
						jen.ID("s").Dot("exampleWebhook"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dot("exampleWebhook"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.Op("&").ID("types").Dot("WebhookCreationInput").Valuesln(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("s").Dot("exampleWebhook")),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("s").Dot("exampleWebhook")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("webhooksTestSuite")).ID("TestClient_UpdateWebhook").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/webhooks/%d"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPut"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleWebhook"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("UpdateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid webhook ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("UpdateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("UpdateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("UpdateWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook"),
					),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("webhooksTestSuite")).ID("TestClient_ArchiveWebhook").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPathFormat").Op("=").Lit("/api/v1/webhooks/%d"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodDelete"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid webhook ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("webhooksTestSuite")).ID("TestClient_GetAuditLogForWebhook").Params().Body(
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/webhooks/%d/audit"),
				jen.ID("expectedMethod").Op("=").Qual("net/http", "MethodGet"),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("exampleAuditLogEntryList").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleAuditLogEntryList"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAuditLogEntryList"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid webhook ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForWebhook").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
