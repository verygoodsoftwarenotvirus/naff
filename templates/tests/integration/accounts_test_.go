package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("checkAccountEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("types").Dot("Account")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
				jen.Lit("expected BucketName for account #%d to be %v, but it was %v "),
				jen.ID("expected").Dot("ID"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
			),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedOn"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Creating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to create accounts"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
						jen.ID("exampleAccountInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
						jen.List(jen.ID("createdAccount"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccountInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAccount"),
							jen.ID("err"),
						),
						jen.ID("checkAccountEquality").Call(
							jen.ID("t"),
							jen.ID("exampleAccount"),
							jen.ID("createdAccount"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("createdAccount").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("UserAddedToAccountEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdAccount").Dot("ID"),
							jen.ID("audit").Dot("AccountAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("createdAccount").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Listing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to list accounts"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Var().Defs(
							jen.ID("expected").Index().Op("*").ID("types").Dot("Account"),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
							jen.ID("exampleAccountInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
							jen.List(jen.ID("createdAccount"), jen.ID("accountCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
								jen.ID("ctx"),
								jen.ID("exampleAccountInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdAccount"),
								jen.ID("accountCreationErr"),
							),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdAccount"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAccounts").Call(
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
							jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Accounts")),
							jen.Lit("expected %d to be <= %d"),
							jen.ID("len").Call(jen.ID("expected")),
							jen.ID("len").Call(jen.ID("actual").Dot("Accounts")),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("createdAccount")).Op(":=").Range().ID("actual").Dot("Accounts")).Body(
							jen.ID("assert").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
									jen.ID("ctx"),
									jen.ID("createdAccount").Dot("ID"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Reading_Returns404ForNonexistentAccount").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to read a non-existent account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAccount").Call(
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Reading").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to read an account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
						jen.ID("exampleAccountInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
						jen.List(jen.ID("createdAccount"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccountInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAccount"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAccount").Call(
							jen.ID("ctx"),
							jen.ID("createdAccount").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkAccountEquality").Call(
							jen.ID("t"),
							jen.ID("exampleAccount"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("createdAccount").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Updating_Returns404ForNonexistentAccount").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to update a non-existent account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
						jen.ID("exampleAccount").Dot("ID").Op("=").ID("nonexistentID"),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateAccount").Call(
								jen.ID("ctx"),
								jen.ID("exampleAccount"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Updating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to update an account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
						jen.ID("exampleAccountInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
						jen.List(jen.ID("createdAccount"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccountInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAccount"),
							jen.ID("err"),
						),
						jen.ID("createdAccount").Dot("Update").Call(jen.ID("converters").Dot("ConvertAccountToAccountUpdateInput").Call(jen.ID("exampleAccount"))),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UpdateAccount").Call(
								jen.ID("ctx"),
								jen.ID("createdAccount"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAccount").Call(
							jen.ID("ctx"),
							jen.ID("createdAccount").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkAccountEquality").Call(
							jen.ID("t"),
							jen.ID("exampleAccount"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("actual").Dot("LastUpdatedOn"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("createdAccount").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("UserAddedToAccountEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountUpdateEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdAccount").Dot("ID"),
							jen.ID("audit").Dot("AccountAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("createdAccount").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Archiving_Returns404ForNonexistentAccount").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to archiv a non-existent account"),
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
							jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("nonexistentID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Archiving").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to archive an account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
						jen.ID("exampleAccountInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
						jen.List(jen.ID("createdAccount"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccountInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAccount"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("createdAccount").Dot("ID"),
							),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("createdAccount").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("UserAddedToAccountEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountArchiveEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdAccount").Dot("ID"),
							jen.ID("audit").Dot("AccountAssignmentKey"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_ChangingMemberships").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to change members of an account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Var().Defs(
							jen.ID("userCount").Op("=").Lit(1),
						),
						jen.List(jen.ID("currentStatus"), jen.ID("statusErr")).Op(":=").ID("testClients").Dot("main").Dot("UserStatus").Call(jen.ID("s").Dot("ctx")),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("currentStatus"),
							jen.ID("statusErr"),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("initial account is #%d; initial user ID is #%d"),
							jen.ID("currentStatus").Dot("ActiveAccount"),
							jen.ID("s").Dot("user").Dot("ID"),
						),
						jen.ID("accountCreationInput").Op(":=").Op("&").ID("types").Dot("AccountCreationInput").Valuesln(jen.ID("Name").Op(":").ID("fakes").Dot("BuildFakeAccount").Call().Dot("Name")),
						jen.List(jen.ID("account"), jen.ID("accountCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("accountCreationInput"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("accountCreationErr"),
						),
						jen.ID("require").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("account"),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("created account #%d"),
							jen.ID("account").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("SwitchActiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("account").Dot("ID"),
							),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("switched main test client active account to #%d, creating webhook"),
							jen.ID("account").Dot("ID"),
						),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
						jen.List(jen.ID("createdWebhook"), jen.ID("creationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdWebhook"),
							jen.ID("creationErr"),
						),
						jen.ID("require").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("account").Dot("ID"),
							jen.ID("createdWebhook").Dot("BelongsToAccount"),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("created webhook #%d for account #%d"),
							jen.ID("createdWebhook").Dot("ID"),
							jen.ID("createdWebhook").Dot("BelongsToAccount"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("UserAddedToAccountEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("WebhookCreationEvent"))),
						jen.ID("users").Op(":=").Index().Op("*").ID("types").Dot("User").Values(),
						jen.ID("clients").Op(":=").Index().Op("*").ID("httpclient").Dot("Client").Values(),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.List(jen.ID("u"), jen.ID("_"), jen.ID("c"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
								jen.ID("ctx"),
								jen.ID("t"),
							),
							jen.ID("users").Op("=").ID("append").Call(
								jen.ID("users"),
								jen.ID("u"),
							),
							jen.ID("clients").Op("=").ID("append").Call(
								jen.ID("clients"),
								jen.ID("c"),
							),
							jen.List(jen.ID("currentStatus"), jen.ID("statusErr")).Op("=").ID("c").Dot("UserStatus").Call(jen.ID("s").Dot("ctx")),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("currentStatus"),
								jen.ID("statusErr"),
							),
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("created user user #%d with account #%d"),
								jen.ID("u").Dot("ID"),
								jen.ID("currentStatus").Dot("ActiveAccount"),
							),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("checking that user #%d CANNOT see webhook #%d belonging to account #%d"),
								jen.ID("users").Index(jen.ID("i")).Dot("ID"),
								jen.ID("createdWebhook").Dot("ID"),
								jen.ID("createdWebhook").Dot("BelongsToAccount"),
							),
							jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("clients").Index(jen.ID("i")).Dot("GetWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
							jen.ID("require").Dot("Nil").Call(
								jen.ID("t"),
								jen.ID("webhook"),
							),
							jen.ID("require").Dot("Error").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("adding user #%d to account #%d"),
								jen.ID("users").Index(jen.ID("i")).Dot("ID"),
								jen.ID("account").Dot("ID"),
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("AddUserToAccount").Call(
									jen.ID("ctx"),
									jen.Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("UserID").Op(":").ID("users").Index(jen.ID("i")).Dot("ID"), jen.ID("AccountID").Op(":").ID("account").Dot("ID"), jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())),
								),
							),
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("added user #%d to account #%d"),
								jen.ID("users").Index(jen.ID("i")).Dot("ID"),
								jen.ID("account").Dot("ID"),
							),
							jen.ID("expectedAuditLogEntries").Op("=").ID("append").Call(
								jen.ID("expectedAuditLogEntries"),
								jen.Op("&").ID("types").Dot("AuditLogEntry").Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("UserAddedToAccountEvent")),
							),
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("setting user #%d's client to account #%d"),
								jen.ID("users").Index(jen.ID("i")).Dot("ID"),
								jen.ID("account").Dot("ID"),
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("clients").Index(jen.ID("i")).Dot("SwitchActiveAccount").Call(
									jen.ID("ctx"),
									jen.ID("account").Dot("ID"),
								),
							),
							jen.List(jen.ID("currentStatus"), jen.ID("statusErr")).Op("=").ID("clients").Index(jen.ID("i")).Dot("UserStatus").Call(jen.ID("s").Dot("ctx")),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("currentStatus"),
								jen.ID("statusErr"),
							),
							jen.ID("require").Dot("Equal").Call(
								jen.ID("t"),
								jen.ID("currentStatus").Dot("ActiveAccount"),
								jen.ID("account").Dot("ID"),
							),
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("set user #%d's current active account to #%d"),
								jen.ID("users").Index(jen.ID("i")).Dot("ID"),
								jen.ID("account").Dot("ID"),
							),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.ID("input").Op(":=").Op("&").ID("types").Dot("ModifyUserPermissionsInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("NewRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())),
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ModifyMemberPermissions").Call(
									jen.ID("ctx"),
									jen.ID("account").Dot("ID"),
									jen.ID("users").Index(jen.ID("i")).Dot("ID"),
									jen.ID("input"),
								),
							),
							jen.ID("expectedAuditLogEntries").Op("=").ID("append").Call(
								jen.ID("expectedAuditLogEntries"),
								jen.Op("&").ID("types").Dot("AuditLogEntry").Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("UserAccountPermissionsModifiedEvent")),
							),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("checking if user #%d CAN now see webhook #%d belonging to account #%d"),
								jen.ID("users").Index(jen.ID("i")).Dot("ID"),
								jen.ID("createdWebhook").Dot("ID"),
								jen.ID("createdWebhook").Dot("BelongsToAccount"),
							),
							jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("clients").Index(jen.ID("i")).Dot("GetWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("webhook"),
								jen.ID("err"),
							),
						),
						jen.ID("originalWebhookName").Op(":=").ID("createdWebhook").Dot("Name"),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.ID("createdWebhook").Dot("Name").Op("=").Qual("fmt", "Sprintf").Call(
								jen.Lit("%s_%d"),
								jen.ID("originalWebhookName"),
								jen.Qual("time", "Now").Call().Dot("UnixNano").Call(),
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("clients").Index(jen.ID("i")).Dot("UpdateWebhook").Call(
									jen.ID("ctx"),
									jen.ID("createdWebhook"),
								),
							),
							jen.ID("expectedAuditLogEntries").Op("=").ID("append").Call(
								jen.ID("expectedAuditLogEntries"),
								jen.Op("&").ID("types").Dot("AuditLogEntry").Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("WebhookUpdateEvent")),
							),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("RemoveUserFromAccount").Call(
									jen.ID("ctx"),
									jen.ID("account").Dot("ID"),
									jen.ID("users").Index(jen.ID("i")).Dot("ID"),
									jen.ID("t").Dot("Name").Call(),
								),
							)),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("clients").Index(jen.ID("i")).Dot("GetWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
							jen.ID("require").Dot("Nil").Call(
								jen.ID("t"),
								jen.ID("webhook"),
							),
							jen.ID("require").Dot("Error").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("account").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("account").Dot("ID"),
							jen.ID("audit").Dot("AccountAssignmentKey"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("userCount"), jen.ID("i").Op("++")).Body(
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
									jen.ID("ctx"),
									jen.ID("users").Index(jen.ID("i")).Dot("ID"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_OwnershipTransfer").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to transfer ownership of an account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("futureOwner"), jen.ID("_"), jen.ID("_"), jen.ID("futureOwnerClient")).Op(":=").ID("createUserAndClientForTest").Call(
							jen.ID("ctx"),
							jen.ID("t"),
						),
						jen.ID("accountCreationInput").Op(":=").Op("&").ID("types").Dot("AccountCreationInput").Valuesln(jen.ID("Name").Op(":").ID("fakes").Dot("BuildFakeAccount").Call().Dot("Name")),
						jen.List(jen.ID("account"), jen.ID("accountCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("accountCreationInput"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("accountCreationErr"),
						),
						jen.ID("require").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("account"),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("created account #%d"),
							jen.ID("account").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("SwitchActiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("account").Dot("ID"),
							),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("switched to active account: %d"),
							jen.ID("account").Dot("ID"),
						),
						jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
						jen.ID("exampleWebhookInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
						jen.List(jen.ID("createdWebhook"), jen.ID("creationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhookInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdWebhook"),
							jen.ID("creationErr"),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("created webhook #%d belonging to account #%d"),
							jen.ID("createdWebhook").Dot("ID"),
							jen.ID("createdWebhook").Dot("BelongsToAccount"),
						),
						jen.ID("require").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("account").Dot("ID"),
							jen.ID("createdWebhook").Dot("BelongsToAccount"),
						),
						jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("futureOwnerClient").Dot("GetWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("require").Dot("Nil").Call(
							jen.ID("t"),
							jen.ID("webhook"),
						),
						jen.ID("require").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("TransferAccountOwnership").Call(
								jen.ID("ctx"),
								jen.ID("account").Dot("ID"),
								jen.Op("&").ID("types").Dot("AccountOwnershipTransferInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("CurrentOwner").Op(":").ID("account").Dot("BelongsToUser"), jen.ID("NewOwner").Op(":").ID("futureOwner").Dot("ID")),
							),
						),
						jen.ID("t").Dot("Logf").Call(
							jen.Lit("transferred account %d from user %d to user %d"),
							jen.ID("account").Dot("ID"),
							jen.ID("account").Dot("BelongsToUser"),
							jen.ID("futureOwner").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("futureOwnerClient").Dot("SwitchActiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("account").Dot("ID"),
							),
						),
						jen.List(jen.ID("webhook"), jen.ID("err")).Op("=").ID("futureOwnerClient").Dot("GetWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("webhook"),
							jen.ID("err"),
						),
						jen.List(jen.ID("webhook"), jen.ID("err")).Op("=").ID("testClients").Dot("main").Dot("GetWebhook").Call(
							jen.ID("ctx"),
							jen.ID("createdWebhook").Dot("ID"),
						),
						jen.ID("require").Dot("Nil").Call(
							jen.ID("t"),
							jen.ID("webhook"),
						),
						jen.ID("require").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("futureOwnerClient").Dot("UpdateWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook"),
							),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("account").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("UserAddedToAccountEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("WebhookCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("AccountTransferredEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("WebhookUpdateEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("account").Dot("ID"),
							jen.ID("audit").Dot("AccountAssignmentKey"),
						),
						jen.ID("require").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("futureOwnerClient").Dot("ArchiveWebhook").Call(
								jen.ID("ctx"),
								jen.ID("createdWebhook").Dot("ID"),
							),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("futureOwner").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Auditing_Returns404ForNonexistentAccount").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to audit a non-existent account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Auditing_InaccessibleToNonAdmins").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to audit an account as non-admin"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
						jen.ID("exampleAccountInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
						jen.List(jen.ID("createdAccount"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccountInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAccount"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("createdAccount").Dot("ID"),
						),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("Nil").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("createdAccount").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAccounts_Auditing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to audit an account"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
						jen.ID("exampleAccountInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
						jen.List(jen.ID("createdAccount"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccountInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAccount"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAccount").Call(
							jen.ID("ctx"),
							jen.ID("createdAccount").Dot("ID"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAccount").Call(
								jen.ID("ctx"),
								jen.ID("createdAccount").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	return code
}
