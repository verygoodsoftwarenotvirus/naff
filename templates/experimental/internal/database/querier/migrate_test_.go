package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrateTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestQuerier_Migrate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleCreationTime").Op(":=").ID("fakes").Dot("BuildFakeTime").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleUser").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("exampleUser").Dot("CreatedOn").Op("=").ID("exampleCreationTime"),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccountForUser").Call(jen.ID("exampleUser")),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccountCreationInput").Op(":=").Op("&").ID("types").Dot("AccountCreationInput").Valuesln(jen.ID("Name").Op(":").Qual("fmt", "Sprintf").Call(
						jen.Lit("%s_default"),
						jen.ID("exampleUser").Dot("Username"),
					), jen.ID("BelongsToUser").Op(":").ID("exampleUser").Dot("ID")),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(jen.ID("Username").Op(":").ID("exampleUser").Dot("Username"), jen.ID("Password").Op(":").ID("exampleUser").Dot("HashedPassword"), jen.ID("HashedPassword").Op(":").ID("exampleUser").Dot("HashedPassword"), jen.ID("IsServiceAdmin").Op(":").ID("true")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleCreationTime")),
					jen.ID("db").Dot("ExpectPing").Call(),
					jen.ID("migrationFuncCalled").Op(":=").ID("false"),
					jen.ID("mockQueryBuilder").Dot("On").Call(
						jen.Lit("BuildMigrationFunc"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").Qual("database/sql", "DB").Valuesln()),
					).Dot("Return").Call(jen.Func().Params().Body(
						jen.ID("migrationFuncCalled").Op("=").ID("true"))),
					jen.List(jen.ID("fakeTestUserExistenceQuery"), jen.ID("fakeTestUserExistenceArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("UserSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetUserByUsernameQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("fakeTestUserExistenceQuery"),
						jen.ID("fakeTestUserExistenceArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeTestUserExistenceQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeTestUserExistenceArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeTestUserCreationQuery"), jen.ID("fakeTestUserCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTestUserCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeTestUserCreationQuery"),
						jen.ID("fakeTestUserCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeTestUserCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeTestUserCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleUser").Dot("ID"))),
					jen.List(jen.ID("firstFakeAuditLogEntryEventQuery"), jen.ID("firstFakeAuditLogEntryEventArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "AuditLogEntryCreationInputMatcher").Call(jen.ID("audit").Dot("UserCreationEvent"))),
					).Dot("Return").Call(
						jen.ID("firstFakeAuditLogEntryEventQuery"),
						jen.ID("firstFakeAuditLogEntryEventArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("firstFakeAuditLogEntryEventQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("firstFakeAuditLogEntryEventArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("sqlmock").Dot("NewResult").Call(
						jen.Lit(1),
						jen.Lit(1),
					)),
					jen.List(jen.ID("fakeAccountCreationQuery"), jen.ID("fakeAccountCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAccountCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccountCreationInput"),
					).Dot("Return").Call(
						jen.ID("fakeAccountCreationQuery"),
						jen.ID("fakeAccountCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("secondFakeAuditLogEntryEventQuery"), jen.ID("secondFakeAuditLogEntryEventArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "AuditLogEntryCreationInputMatcher").Call(jen.ID("audit").Dot("AccountCreationEvent"))),
					).Dot("Return").Call(
						jen.ID("secondFakeAuditLogEntryEventQuery"),
						jen.ID("secondFakeAuditLogEntryEventArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeAuditLogEntryEventQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeAuditLogEntryEventArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("sqlmock").Dot("NewResult").Call(
						jen.Lit(1),
						jen.Lit(1),
					)),
					jen.List(jen.ID("fakeMembershipCreationQuery"), jen.ID("fakeMembershipCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateMembershipForNewUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeMembershipCreationQuery"),
						jen.ID("fakeMembershipCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeMembershipCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeMembershipCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("thirdFakeAuditLogEntryEventQuery"), jen.ID("thirdFakeAuditLogEntryEventArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "AuditLogEntryCreationInputMatcher").Call(jen.ID("audit").Dot("UserAddedToAccountEvent"))),
					).Dot("Return").Call(
						jen.ID("thirdFakeAuditLogEntryEventQuery"),
						jen.ID("thirdFakeAuditLogEntryEventArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("thirdFakeAuditLogEntryEventQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("thirdFakeAuditLogEntryEventArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("sqlmock").Dot("NewResult").Call(
						jen.Lit(1),
						jen.Lit(1),
					)),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("err").Op(":=").ID("c").Dot("Migrate").Call(
						jen.ID("ctx"),
						jen.Lit(1),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("migrationFuncCalled"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with failure executing creation query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleCreationTime").Op(":=").ID("fakes").Dot("BuildFakeTime").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleUser").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("exampleUser").Dot("CreatedOn").Op("=").ID("exampleCreationTime"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(jen.ID("Username").Op(":").ID("exampleUser").Dot("Username"), jen.ID("Password").Op(":").ID("exampleUser").Dot("HashedPassword"), jen.ID("HashedPassword").Op(":").ID("exampleUser").Dot("HashedPassword"), jen.ID("IsServiceAdmin").Op(":").ID("true")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleCreationTime")),
					jen.ID("db").Dot("ExpectPing").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("On").Call(
						jen.Lit("BuildMigrationFunc"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").Qual("database/sql", "DB").Valuesln()),
					).Dot("Return").Call(jen.Func().Params().Body()),
					jen.List(jen.ID("fakeTestUserExistenceQuery"), jen.ID("fakeTestUserExistenceArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("UserSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetUserByUsernameQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("fakeTestUserExistenceQuery"),
						jen.ID("fakeTestUserExistenceArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeTestUserExistenceQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeTestUserExistenceArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.List(jen.ID("fakeTestUserCreationQuery"), jen.ID("fakeTestUserCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTestUserCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeTestUserCreationQuery"),
						jen.ID("fakeTestUserCreationArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("Migrate").Call(
							jen.ID("ctx"),
							jen.Lit(1),
							jen.ID("exampleInput"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
