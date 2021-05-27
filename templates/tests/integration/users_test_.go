package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("checkUserCreationEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("expected").Op("*").ID("types").Dot("UserRegistrationInput"), jen.ID("actual").Op("*").ID("types").Dot("UserCreationResponse")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedUserID"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Username"),
				jen.ID("actual").Dot("Username"),
			),
			jen.ID("assert").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("TwoFactorSecret"),
			),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedOn"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("checkUserEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("types").Dot("User")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Username"),
				jen.ID("actual").Dot("Username"),
			),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedOn"),
			),
			jen.ID("assert").Dot("Nil").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("LastUpdatedOn"),
			),
			jen.ID("assert").Dot("Nil").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ArchivedOn"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Creating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be creatable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleUserInput").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
						jen.List(jen.ID("createdUser"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateUser").Call(
							jen.ID("ctx"),
							jen.ID("exampleUserInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdUser"),
							jen.ID("err"),
						),
						jen.ID("checkUserCreationEquality").Call(
							jen.ID("t"),
							jen.ID("exampleUserInput"),
							jen.ID("createdUser"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForUser").Call(
							jen.ID("ctx"),
							jen.ID("createdUser").Dot("CreatedUserID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserAddedToAccountEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdUser").Dot("CreatedUserID"),
							jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("createdUser").Dot("CreatedUserID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Reading_Returns404ForNonexistentUser").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should return an error when trying to read a user that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetUser").Call(
							jen.ID("ctx"),
							jen.ID("nonexistentID"),
						),
						jen.ID("assert").Dot("Nil").Call(
							jen.ID("t"),
							jen.ID("actual"),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Reading").Params().Body(
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
						jen.List(jen.ID("user"), jen.ID("_"), jen.ID("_"), jen.ID("_")).Op(":=").ID("createUserAndClientForTest").Call(
							jen.ID("ctx"),
							jen.ID("t"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetUser").Call(
							jen.ID("ctx"),
							jen.ID("user").Dot("ID"),
						),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
							jen.ID("t").Dot("Logf").Call(
								jen.Lit("error encountered trying to fetch user %q: %v\n"),
								jen.ID("user").Dot("Username"),
								jen.ID("err"),
							)),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkUserEquality").Call(
							jen.ID("t"),
							jen.ID("user"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("actual").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Searching_ReturnsEmptyWhenSearchingForUsernameThatIsNotPresentInTheDatabase").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return empty slice when searching for a username that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("SearchForUsersByUsername").Call(
							jen.ID("ctx"),
							jen.Lit("   this is a really long string that contains characters unlikely to yield any real results   "),
						),
						jen.ID("assert").Dot("Nil").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Searching_OnlyAccessibleToAdmins").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should only be accessible to admins"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("SearchForUsersByUsername").Call(
							jen.ID("ctx"),
							jen.ID("s").Dot("user").Dot("Username"),
						),
						jen.ID("assert").Dot("Nil").Call(
							jen.ID("t"),
							jen.ID("actual"),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Searching").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return be searchable"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleUsername").Op(":=").ID("fakes").Dot("BuildFakeUser").Call().Dot("Username"),
						jen.ID("createdUserIDs").Op(":=").Index().ID("uint64").Values(),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Body(
							jen.List(jen.ID("user"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "CreateServiceUser").Call(
								jen.ID("ctx"),
								jen.ID("urlToUse"),
								jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("%s%d"),
									jen.ID("exampleUsername"),
									jen.ID("i"),
								),
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
							jen.ID("createdUserIDs").Op("=").ID("append").Call(
								jen.ID("createdUserIDs"),
								jen.ID("user").Dot("ID"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("SearchForUsersByUsername").Call(
							jen.ID("ctx"),
							jen.ID("exampleUsername"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("actual"),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("result")).Op(":=").Range().ID("actual")).Body(
							jen.ID("assert").Dot("True").Call(
								jen.ID("t"),
								jen.Qual("strings", "HasPrefix").Call(
									jen.ID("result").Dot("Username"),
									jen.ID("exampleUsername"),
								),
							)),
						jen.For(jen.List(jen.ID("_"), jen.ID("id")).Op(":=").Range().ID("createdUserIDs")).Body(
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
									jen.ID("ctx"),
									jen.ID("id"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Archiving_Returns404ForNonexistentUser").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should fail to archive a non-existent user"),
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
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("nonexistentID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Archiving").Params().Body(
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
						jen.ID("exampleUserInput").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
						jen.List(jen.ID("createdUser"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateUser").Call(
							jen.ID("ctx"),
							jen.ID("exampleUserInput"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("createdUser"),
						),
						jen.If(jen.ID("createdUser").Op("==").ID("nil").Op("||").ID("err").Op("!=").ID("nil")).Body(
							jen.ID("t").Dot("Log").Call(jen.Lit("something has gone awry, user returned is nil")),
							jen.ID("t").Dot("FailNow").Call(),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("createdUser").Dot("CreatedUserID"),
							),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForUser").Call(
							jen.ID("ctx"),
							jen.ID("createdUser").Dot("CreatedUserID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserAddedToAccountEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserArchiveEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdUser").Dot("CreatedUserID"),
							jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserAssignmentKey"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Auditing_Returns404ForNonexistentUser").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should return an error when trying to audit something that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("input").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
						jen.ID("input").Dot("NewReputation").Op("=").ID("types").Dot("BannedUserAccountStatus"),
						jen.ID("input").Dot("TargetUserID").Op("=").ID("nonexistentID"),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("UpdateUserReputation").Call(
								jen.ID("ctx"),
								jen.ID("input"),
							),
						),
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForUser").Call(
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Auditing_InaccessibleToNonAdmins").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("it should not be auditable by a non-admin"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
						jen.ID("exampleUserInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("exampleUser")),
						jen.List(jen.ID("createdUser"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateUser").Call(
							jen.ID("ctx"),
							jen.ID("exampleUserInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdUser"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAuditLogForUser").Call(
							jen.ID("ctx"),
							jen.ID("createdUser").Dot("CreatedUserID"),
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
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("createdUser").Dot("CreatedUserID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_Auditing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to be audited"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
						jen.ID("exampleUserInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("exampleUser")),
						jen.List(jen.ID("createdUser"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateUser").Call(
							jen.ID("ctx"),
							jen.ID("exampleUserInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdUser"),
							jen.ID("err"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForUser").Call(
							jen.ID("ctx"),
							jen.ID("createdUser").Dot("CreatedUserID"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserAddedToAccountEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "AccountCreationEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.Lit(0),
							jen.Lit(""),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("createdUser").Dot("CreatedUserID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestUsers_AvatarManagement").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be able to upload an avatar"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("avatar").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildArbitraryImagePNGBytes").Call(jen.Lit(256)),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("UploadNewAvatar").Call(
								jen.ID("ctx"),
								jen.ID("avatar"),
								jen.Lit("png"),
							),
						),
						jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetUser").Call(
							jen.ID("ctx"),
							jen.ID("s").Dot("user").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("user"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("user").Dot("AvatarSrc"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForUser").Call(
							jen.ID("ctx"),
							jen.ID("s").Dot("user").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "AccountCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserAddedToAccountEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserVerifyTwoFactorSecretEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "SuccessfulLoginEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "APIClientCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "UserUpdateEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("s").Dot("user").Dot("ID"),
							jen.Lit(""),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("s").Dot("user").Dot("ID"),
							),
						),
					)),
				jen.ID("pasetoAuthType"),
			)),
		jen.Line(),
	)

	return code
}
