package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("checkAPIClientEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("types").Dot("APIClient")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
				jen.Lit("expected LabelName for API client #%d to be %q, but it was %q "),
				jen.ID("actual").Dot("ID"),
				jen.ID("expected").Dot("Name"),
				jen.ID("actual").Dot("Name"),
			),
			jen.ID("assert").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ExternalID"),
				jen.Lit("expected ExternalID for API client #%d to not be empty, but it was"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ClientID"),
				jen.Lit("expected ClientID for API client #%d to not be empty, but it was"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("Empty").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("ClientSecret"),
				jen.Lit("expected ClientSecret for API client #%d to not be empty, but it was"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("BelongsToUser"),
				jen.Lit("expected BelongsToUser for API client #%d to not be zero, but it was"),
				jen.ID("actual").Dot("ID"),
			),
			jen.ID("assert").Dot("NotZero").Call(
				jen.ID("t"),
				jen.ID("actual").Dot("CreatedOn"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAPIClients_Creating").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to create API clients"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
						jen.ID("exampleAPIClientInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
						jen.ID("exampleAPIClientInput").Dot("UserLoginInput").Op("=").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("s").Dot("user").Dot("Username"), jen.ID("Password").Op(":").ID("s").Dot("user").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("s").Dot("user"),
						)),
						jen.List(jen.ID("createdAPIClient"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("s").Dot("cookie"),
							jen.ID("exampleAPIClientInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAPIClient"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("createdAPIClient").Dot("ClientID"),
							jen.Lit("expected ClientID for API client #%d to not be empty, but it was"),
							jen.ID("createdAPIClient").Dot("ID"),
						),
						jen.ID("assert").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("createdAPIClient").Dot("ClientSecret"),
							jen.Lit("expected ClientSecret for API client #%d to not be empty, but it was"),
							jen.ID("createdAPIClient").Dot("ID"),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("createdAPIClient").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("APIClientCreationEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdAPIClient").Dot("ID"),
							jen.ID("audit").Dot("APIClientAssignmentKey"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAPIClient").Call(
								jen.ID("ctx"),
								jen.ID("createdAPIClient").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAPIClients_Listing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to read API clients in a list"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.Var().Defs(
							jen.ID("clientsToMake").Op("=").Lit(1),
						),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Var().Defs(
							jen.ID("expected").Index().ID("uint64"),
						),
						jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("clientsToMake"), jen.ID("i").Op("++")).Body(
							jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
							jen.ID("exampleAPIClientInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
							jen.ID("exampleAPIClientInput").Dot("UserLoginInput").Op("=").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("s").Dot("user").Dot("Username"), jen.ID("Password").Op(":").ID("s").Dot("user").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
								jen.ID("t"),
								jen.ID("s").Dot("user"),
							)),
							jen.List(jen.ID("createdAPIClient"), jen.ID("apiClientCreationErr")).Op(":=").ID("testClients").Dot("main").Dot("CreateAPIClient").Call(
								jen.ID("ctx"),
								jen.ID("s").Dot("cookie"),
								jen.ID("exampleAPIClientInput"),
							),
							jen.ID("requireNotNilAndNoProblems").Call(
								jen.ID("t"),
								jen.ID("createdAPIClient"),
								jen.ID("apiClientCreationErr"),
							),
							jen.ID("expected").Op("=").ID("append").Call(
								jen.ID("expected"),
								jen.ID("createdAPIClient").Dot("ID"),
							),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAPIClients").Call(
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
							jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Clients")),
							jen.Lit("expected %d to be <= %d"),
							jen.ID("len").Call(jen.ID("expected")),
							jen.ID("len").Call(jen.ID("actual").Dot("Clients")),
						),
						jen.For(jen.List(jen.ID("_"), jen.ID("createdAPIClientID")).Op(":=").Range().ID("expected")).Body(
							jen.ID("assert").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("testClients").Dot("main").Dot("ArchiveAPIClient").Call(
									jen.ID("ctx"),
									jen.ID("createdAPIClientID"),
								),
							)),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAPIClients_Reading_Returns404ForNonexistentAPIClient").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to read non-existent API clients"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAPIClient").Call(
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
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAPIClients_Reading").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to read API clients"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
						jen.ID("exampleAPIClientInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
						jen.ID("exampleAPIClientInput").Dot("UserLoginInput").Op("=").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("s").Dot("user").Dot("Username"), jen.ID("Password").Op(":").ID("s").Dot("user").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("s").Dot("user"),
						)),
						jen.List(jen.ID("createdAPIClient"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("s").Dot("cookie"),
							jen.ID("exampleAPIClientInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAPIClient"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("GetAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("createdAPIClient").Dot("ID"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("err"),
						),
						jen.ID("checkAPIClientEquality").Call(
							jen.ID("t"),
							jen.ID("exampleAPIClient"),
							jen.ID("actual"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAPIClient").Call(
								jen.ID("ctx"),
								jen.ID("createdAPIClient").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAPIClients_Archiving_Returns404ForNonexistentAPIClient").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to archive non-existent API clients"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.Qual("context", "Background").Call(),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAPIClient").Call(
								jen.ID("ctx"),
								jen.ID("nonexistentID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAPIClients_Archiving").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to archive API clients"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
						jen.ID("exampleAPIClientInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
						jen.ID("exampleAPIClientInput").Dot("UserLoginInput").Op("=").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("s").Dot("user").Dot("Username"), jen.ID("Password").Op(":").ID("s").Dot("user").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("s").Dot("user"),
						)),
						jen.List(jen.ID("createdAPIClient"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("s").Dot("cookie"),
							jen.ID("exampleAPIClientInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAPIClient"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("main").Dot("ArchiveAPIClient").Call(
								jen.ID("ctx"),
								jen.ID("createdAPIClient").Dot("ID"),
							),
						),
						jen.List(jen.ID("auditLogEntries"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("createdAPIClient").Dot("ID"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("expectedAuditLogEntries").Op(":=").Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln(jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("APIClientCreationEvent")), jen.Valuesln(jen.ID("EventType").Op(":").ID("audit").Dot("APIClientArchiveEvent"))),
						jen.ID("validateAuditLogEntries").Call(
							jen.ID("t"),
							jen.ID("expectedAuditLogEntries"),
							jen.ID("auditLogEntries"),
							jen.ID("createdAPIClient").Dot("ID"),
							jen.ID("audit").Dot("APIClientAssignmentKey"),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAPIClients_Auditing").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to audit API clients"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("exampleAPIClient").Op(":=").ID("fakes").Dot("BuildFakeAPIClient").Call(),
						jen.ID("exampleAPIClientInput").Op(":=").ID("fakes").Dot("BuildFakeAPIClientCreationInputFromClient").Call(jen.ID("exampleAPIClient")),
						jen.ID("exampleAPIClientInput").Dot("UserLoginInput").Op("=").ID("types").Dot("UserLoginInput").Valuesln(jen.ID("Username").Op(":").ID("s").Dot("user").Dot("Username"), jen.ID("Password").Op(":").ID("s").Dot("user").Dot("HashedPassword"), jen.ID("TOTPToken").Op(":").ID("generateTOTPTokenForUser").Call(
							jen.ID("t"),
							jen.ID("s").Dot("user"),
						)),
						jen.List(jen.ID("createdAPIClient"), jen.ID("err")).Op(":=").ID("testClients").Dot("main").Dot("CreateAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("s").Dot("cookie"),
							jen.ID("exampleAPIClientInput"),
						),
						jen.ID("requireNotNilAndNoProblems").Call(
							jen.ID("t"),
							jen.ID("createdAPIClient"),
							jen.ID("err"),
						),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClients").Dot("admin").Dot("GetAuditLogForAPIClient").Call(
							jen.ID("ctx"),
							jen.ID("createdAPIClient").Dot("ID"),
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
							jen.ID("testClients").Dot("main").Dot("ArchiveAPIClient").Call(
								jen.ID("ctx"),
								jen.ID("createdAPIClient").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	return code
}
