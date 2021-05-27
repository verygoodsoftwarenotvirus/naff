package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("checkWebhookEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("types").Dot("Webhook")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("ContentType"),
				jen.ID("actual").Dot("ContentType"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("URL"),
				jen.ID("actual").Dot("URL"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Method"),
				jen.ID("actual").Dot("Method"),
			),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedOn"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Creating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be createable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
						jen.List(jen.ID("createdWebhook"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdWebhook"),
							jen.ID("err"),
						),
						jen.ID("checkWebhookEquality").Call(
							jen.ID("t"),
							jen.ID("exampleWebhook"),
							jen.ID("createdWebhook"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkWebhookEquality").Call(
							jen.ID("t"),
							jen.ID("exampleWebhook"),
							jen.ID("actual"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookCreationEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdWebhook").Dot("ID"),
							jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Reading_Returns404ForNonexistentWebhook").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should fail to read non-existent webhook"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetWebhook").Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
						),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Reading").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be read"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
						jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("premade"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetWebhook").Call(
							jen.ID("ctx"),
							jen.ID("premade").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkWebhookEquality").Call(
							jen.ID("t"),
							jen.ID("exampleWebhook"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("actual").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Listing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be read in a list"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Var().Defs(
							jen.ID("expected").Index().Op("*").ID("types").Dot("Webhook"),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
							jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
							jen.List(jen.ID("createdWebhook"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
								jen.ID("ctx"),
								jen.ID("exampleWebhookInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdWebhook"),
								jen.ID("err"),
							),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdWebhook"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetWebhooks").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("True").Call(
							jen.ID("t"),
							jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Webhooks")),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("webhook")).Op(":=").Range().ID("actual").Dot("Webhooks")).Body(
							jen.ID("assert").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
									jen.ID("ctx"),
									jen.ID("webhook").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Updating_Returns404ForNonexistentWebhook").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should fail to update a non-existent webhook"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhook").Dot("ID").Op("=").ID("nonexistentID"),
						jen.ID("err").Op(":=").ID("testClients").Dot("main").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook"),
						),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Updating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be updateable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
						jen.List(jen.ID("createdWebhook"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdWebhook"),
							jen.ID("err"),
						),
						jen.ID("createdWebhook").Dot("Name").Op("=").ID("reverseString").Call(jen.ID("createdWebhook").Dot("Name")),
						jen.ID("exampleWebhook").Dot("Name").Op("=").ID("createdWebhook").Dot("Name"),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkWebhookEquality").Call(
							jen.ID("t"),
							jen.ID("exampleWebhook"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("actual").Dot("LastUpdatedOn"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookUpdateEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdWebhook").Dot("ID"),
							jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("actual").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Archiving_Returns404ForNonexistentWebhook").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should fail to archive a non-existent webhook"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("nonexistentID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Archiving").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be archived"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
						jen.List(jen.ID("createdWebhook"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdWebhook"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookArchiveEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdWebhook").Dot("ID"),
							jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "WebhookAssignmentKey"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Auditing_Returns404ForNonexistentWebhook").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should return an error when auditing a non-existent webhook"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhook").Dot("ID").Op("=").ID("nonexistentID"),
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook").Dot("ID"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("Empty").Call(
							jen.ID("t"),
							jen.ID("x"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestWebhooks_Auditing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be auditable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
						jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("premade"),
							jen.ID("err"),
						),
						jen.ID("premade").Dot("Name").Op("=").ID("reverseString").Call(jen.ID("premade").Dot("Name")),
						jen.ID("exampleWebhook").Dot("Name").Op("=").ID("premade").Dot("Name"),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateWebhook").Call(
								jen.ID("ctx"),
								jen.ID("premade"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForWebhook").Call(
							jen.ID("ctx"),
							jen.ID("premade").Dot("ID"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("Len").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.Lit(2),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("premade").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	return code
}
