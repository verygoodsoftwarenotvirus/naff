package postgres

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
					jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
					jen.ID("exampleUser").Dot("CreatedOn").Equals().ID("exampleCreationTime"),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccountForUser").Call(jen.ID("exampleUser")),
					jen.ID("exampleTestUserConfig").Op(":=").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(
						jen.ID("Username").MapAssign().ID("exampleUser").Dot("Username"), jen.ID("Password").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("HashedPassword").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("IsServiceAdmin").MapAssign().ID("true")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleCreationTime")),
					jen.ID("db").Dot("ExpectPing").Call(),
					jen.ID("c").Dot("migrateOnce").Dot("Do").Call(jen.Func().Params().Body()),
					jen.ID("testUserExistenceArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleTestUserConfig").Dot("Username")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserExistenceQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserExistenceArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("testUserCreationArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(), jen.ID("exampleTestUserConfig").Dot("Username"), jen.ID("exampleTestUserConfig").Dot("HashedPassword"), jen.ID("defaultTestUserTwoFactorSecret"), jen.ID("types").Dot("GoodStandingAccountStatus"), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleUser").Dot("ID"))),
					jen.ID("accountCreationInput").Op(":=").ID("types").Dot("AccountCreationInputForNewUser").Call(jen.ID("exampleUser")),
					jen.ID("accountCreationArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(), jen.ID("accountCreationInput").Dot("Name"), jen.ID("types").Dot("UnpaidAccountBillingStatus"), jen.ID("accountCreationInput").Dot("ContactEmail"), jen.ID("accountCreationInput").Dot("ContactPhone"), jen.Op("&").ID("idMatcher").Values()),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("accountCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("accountCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("createAccountMembershipForNewUserArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(), jen.Op("&").ID("idMatcher").Values(), jen.Op("&").ID("idMatcher").Values(), jen.ID("true"), jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call()),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("createAccountMembershipForNewUserQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("createAccountMembershipForNewUserArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("err").Op(":=").ID("c").Dot("Migrate").Call(
						jen.ID("ctx"),
						jen.Lit(1),
						jen.ID("exampleTestUserConfig"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with failure executing creation query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleCreationTime").Op(":=").ID("fakes").Dot("BuildFakeTime").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
					jen.ID("exampleUser").Dot("CreatedOn").Equals().ID("exampleCreationTime"),
					jen.ID("exampleTestUserConfig").Op(":=").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(
						jen.ID("Username").MapAssign().ID("exampleUser").Dot("Username"), jen.ID("Password").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("HashedPassword").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("IsServiceAdmin").MapAssign().ID("true")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleCreationTime")),
					jen.ID("db").Dot("ExpectPing").Call(),
					jen.ID("c").Dot("migrateOnce").Dot("Do").Call(jen.Func().Params().Body()),
					jen.ID("testUserExistenceArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleTestUserConfig").Dot("Username")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserExistenceQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserExistenceArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("testUserCreationArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(), jen.ID("exampleTestUserConfig").Dot("Username"), jen.ID("exampleTestUserConfig").Dot("HashedPassword"), jen.ID("defaultTestUserTwoFactorSecret"), jen.ID("types").Dot("GoodStandingAccountStatus"), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserCreationArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("err").Op(":=").ID("c").Dot("Migrate").Call(
						jen.ID("ctx"),
						jen.Lit(1),
						jen.ID("exampleTestUserConfig"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
