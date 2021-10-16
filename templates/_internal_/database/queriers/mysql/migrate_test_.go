package mysql

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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleCreationTime").Op(":=").ID("fakes").Dot("BuildFakeTime").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
					jen.ID("exampleUser").Dot("CreatedOn").Equals().ID("exampleCreationTime"),
					jen.Newline(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccountForUser").Call(jen.ID("exampleUser")),
					jen.Newline(),
					jen.ID("exampleTestUserConfig").Op(":=").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(
						jen.ID("Username").MapAssign().ID("exampleUser").Dot("Username"), jen.ID("Password").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("HashedPassword").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("IsServiceAdmin").MapAssign().ID("true")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleCreationTime"),
					),
					jen.Newline(),
					jen.Comment("called by c.IsReady()"),
					jen.ID("db").Dot("ExpectPing").Call(),
					jen.Newline(),
					jen.ID("c").Dot("migrateOnce").Dot("Do").Call(jen.Func().Params().Body()),
					jen.Newline(),
					jen.Comment("expect TestUser to be queried for"),
					jen.ID("testUserExistenceArgs").Op(":=").Index().Interface().Values(jen.ID("exampleTestUserConfig").Dot("Username")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserExistenceQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserExistenceArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.Comment("expect TestUser to be created"),
					jen.ID("testUserCreationArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(), jen.ID("exampleTestUserConfig").Dot("Username"), jen.ID("exampleTestUserConfig").Dot("HashedPassword"), jen.ID("defaultTestUserTwoFactorSecret"), jen.Lit(""), jen.ID("types").Dot("GoodStandingAccountStatus"), jen.Lit(""), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserCreationArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleUser").Dot("ID"))),
					jen.Newline(),
					jen.Comment("create account for created TestUser"),
					jen.ID("accountCreationInput").Op(":=").ID("types").Dot("AccountCreationInputForNewUser").Call(jen.ID("exampleUser")),
					jen.ID("accountCreationArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(),
						jen.ID("accountCreationInput").Dot("Name"),
						jen.ID("types").Dot("UnpaidAccountBillingStatus"),
						jen.ID("accountCreationInput").Dot("ContactEmail"),
						jen.ID("accountCreationInput").Dot("ContactPhone"),
						jen.Lit(""),
						jen.Op("&").ID("idMatcher").Values(),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("accountCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("accountCreationArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.Newline(),
					jen.Comment("create account user membership for created user"),
					jen.ID("createAccountMembershipForNewUserArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(), jen.Op("&").ID("idMatcher").Values(), jen.Op("&").ID("idMatcher").Values(), jen.ID("true"), jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call()),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("createAccountMembershipForNewUserQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("createAccountMembershipForNewUserArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.Newline(),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dot("Migrate").Call(
						jen.ID("ctx"),
						jen.Lit(1),
						jen.ID("exampleTestUserConfig"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with failure executing creation query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleCreationTime").Op(":=").ID("fakes").Dot("BuildFakeTime").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn").Equals().Nil(),
					jen.ID("exampleUser").Dot("CreatedOn").Equals().ID("exampleCreationTime"),
					jen.Newline(),
					jen.ID("exampleTestUserConfig").Op(":=").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(
						jen.ID("Username").MapAssign().ID("exampleUser").Dot("Username"), jen.ID("Password").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("HashedPassword").MapAssign().ID("exampleUser").Dot("HashedPassword"), jen.ID("IsServiceAdmin").MapAssign().ID("true")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleCreationTime")),
					jen.Newline(),
					jen.Comment("called by c.IsReady()"),
					jen.ID("db").Dot("ExpectPing").Call(),
					jen.Newline(),
					jen.ID("c").Dot("migrateOnce").Dot("Do").Call(jen.Func().Params().Body()),
					jen.Newline(),
					jen.Comment("expect TestUser to be queried for"),
					jen.ID("testUserExistenceArgs").Op(":=").Index().Interface().Values(jen.ID("exampleTestUserConfig").Dot("Username")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserExistenceQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserExistenceArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.Comment("expect TestUser to be created"),
					jen.ID("testUserCreationArgs").Op(":=").Index().Interface().Valuesln(
						jen.Op("&").ID("idMatcher").Values(), jen.ID("exampleTestUserConfig").Dot("Username"), jen.ID("exampleTestUserConfig").Dot("HashedPassword"), jen.ID("defaultTestUserTwoFactorSecret"), jen.Lit(""), jen.ID("types").Dot("GoodStandingAccountStatus"), jen.Lit(""), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("testUserCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("testUserCreationArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dot("Migrate").Call(
						jen.ID("ctx"),
						jen.Lit(1),
						jen.ID("exampleTestUserConfig"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
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
