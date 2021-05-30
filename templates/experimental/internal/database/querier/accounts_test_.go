package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildMockRowsFromAccounts").Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").ID("uint64"), jen.ID("accounts").Op("...").Op("*").ID("types").Dot("Account")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("columns").Op(":=").ID("append").Call(
				jen.ID("querybuilding").Dot("AccountsTableColumns"),
				jen.ID("querybuilding").Dot("AccountsUserMembershipTableColumns").Op("..."),
			),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Op("=").ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("columns")),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("accounts")).Body(
				jen.For(jen.List(jen.ID("_"), jen.ID("y")).Op(":=").Range().ID("x").Dot("Members")).Body(
					jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(jen.ID("x").Dot("ID"), jen.ID("x").Dot("ExternalID"), jen.ID("x").Dot("Name"), jen.ID("x").Dot("BillingStatus"), jen.ID("x").Dot("ContactEmail"), jen.ID("x").Dot("ContactPhone"), jen.ID("x").Dot("PaymentProcessorCustomerID"), jen.ID("x").Dot("SubscriptionPlanID"), jen.ID("x").Dot("CreatedOn"), jen.ID("x").Dot("LastUpdatedOn"), jen.ID("x").Dot("ArchivedOn"), jen.ID("x").Dot("BelongsToUser"), jen.ID("y").Dot("ID"), jen.ID("y").Dot("BelongsToUser"), jen.ID("y").Dot("BelongsToAccount"), jen.Qual("strings", "Join").Call(
						jen.ID("y").Dot("AccountRoles"),
						jen.ID("accountMemberRolesSeparator"),
					), jen.ID("y").Dot("DefaultAccount"), jen.ID("y").Dot("CreatedOn"), jen.ID("x").Dot("LastUpdatedOn"), jen.ID("y").Dot("ArchivedOn")),
					jen.If(jen.ID("includeCounts")).Body(
						jen.ID("rowValues").Op("=").ID("append").Call(
							jen.ID("rowValues"),
							jen.ID("filteredCount"),
							jen.ID("len").Call(jen.ID("accounts")),
						)),
					jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
				)),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ScanAccounts").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("surfaces row errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Valuesln(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAccounts").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs row closing errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Valuesln(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAccounts").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
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
		jen.Func().ID("TestQuerier_GetAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAccounts").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleAccount"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAccount"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.ID("exampleUser").Dot("ID"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no returned accounts"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("columns").Op(":=").ID("append").Call(
						jen.ID("querybuilding").Dot("AccountsTableColumns"),
						jen.ID("querybuilding").Dot("AccountsUserMembershipTableColumns").Op("..."),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.ID("columns"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
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

	code.Add(
		jen.Func().ID("TestQuerier_GetAllAccountsCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleCount").Op(":=").ID("uint64").Call(jen.Lit(123)),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllAccountsCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(jen.ID("fakeQuery")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("exampleCount"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAllAccountsCount").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleCount"),
						jen.ID("actual"),
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

	code.Add(
		jen.Func().ID("TestQuerier_GetAllAccounts").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Account")),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("doneChan").Op(":=").ID("make").Call(
						jen.Chan().ID("bool"),
						jen.Lit(1),
					),
					jen.ID("exampleAccountList").Op(":=").ID("fakes").Dot("BuildFakeAccountList").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllAccountsCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Valuesln(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("defaultBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAccounts").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleAccountList").Dot("Accounts").Op("..."),
					)),
					jen.ID("err").Op(":=").ID("c").Dot("GetAllAccounts").Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Var().Defs(
						jen.ID("stillQuerying").Op("=").ID("true"),
					),
					jen.For(jen.ID("stillQuerying")).Body(
						jen.Select().Body(
							jen.Case(jen.ID("batch").Op(":=").Op("<-").ID("results")).Body(
								jen.ID("assert").Dot("NotEmpty").Call(
									jen.ID("t"),
									jen.ID("batch"),
								), jen.ID("doneChan").ReceiveFromChannel().ID("true")),
							jen.Case(jen.Op("<-").Qual("time", "After").Call(jen.Qual("time", "Second"))).Body(
								jen.ID("t").Dot("FailNow").Call()),
							jen.Case(jen.Op("<-").ID("doneChan")).Body(
								jen.ID("stillQuerying").Op("=").ID("false")),
						)),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil results channel"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("GetAllAccounts").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.Lit(0),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with now rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Account")),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllAccountsCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Valuesln(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("err").Op(":=").ID("c").Dot("GetAllAccounts").Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.ID("exampleBatchSize"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching initial count"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Account")),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllAccountsCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Valuesln(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("err").Op(":=").ID("c").Dot("GetAllAccounts").Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.ID("exampleBatchSize"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error querying database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Account")),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllAccountsCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Valuesln(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("err").Op(":=").ID("c").Dot("GetAllAccounts").Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.ID("exampleBatchSize"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Account")),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllAccountsCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Valuesln(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.ID("err").Op(":=").ID("c").Dot("GetAllAccounts").Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.ID("exampleBatchSize"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
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

	code.Add(
		jen.Func().ID("TestQuerier_GetAccounts").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccountList").Op(":=").ID("fakes").Dot("BuildFakeAccountList").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAccounts").Call(
						jen.ID("true"),
						jen.ID("exampleAccountList").Dot("FilteredCount"),
						jen.ID("exampleAccountList").Dot("Accounts").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAccountList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.ID("filter"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil filter"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccountList").Op(":=").ID("fakes").Dot("BuildFakeAccountList").Call(),
					jen.List(jen.ID("exampleAccountList").Dot("Page"), jen.ID("exampleAccountList").Dot("Limit")).Op("=").List(jen.Lit(0), jen.Lit(0)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAccounts").Call(
						jen.ID("true"),
						jen.ID("exampleAccountList").Dot("FilteredCount"),
						jen.ID("exampleAccountList").Dot("Accounts").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAccountList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccounts").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
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

	code.Add(
		jen.Func().ID("TestQuerier_GetAccountsForAdmin").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountList").Op(":=").ID("fakes").Dot("BuildFakeAccountList").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(0)),
						jen.ID("true"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAccounts").Call(
						jen.ID("true"),
						jen.ID("exampleAccountList").Dot("FilteredCount"),
						jen.ID("exampleAccountList").Dot("Accounts").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccountsForAdmin").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAccountList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil filter"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
					jen.ID("exampleAccountList").Op(":=").ID("fakes").Dot("BuildFakeAccountList").Call(),
					jen.List(jen.ID("exampleAccountList").Dot("Page"), jen.ID("exampleAccountList").Dot("Limit")).Op("=").List(jen.Lit(0), jen.Lit(0)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(0)),
						jen.ID("true"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAccounts").Call(
						jen.ID("true"),
						jen.ID("exampleAccountList").Dot("FilteredCount"),
						jen.ID("exampleAccountList").Dot("Accounts").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccountsForAdmin").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAccountList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(0)),
						jen.ID("true"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccountsForAdmin").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(0)),
						jen.ID("true"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAccountsForAdmin").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
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

	code.Add(
		jen.Func().ID("TestQuerier_CreateAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BillingStatus").Op("=").ID("types").Dot("UnpaidAccountBillingStatus"),
					jen.ID("exampleAccount").Dot("PaymentProcessorCustomerID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Call(jen.ID("nil")),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("exampleAccountAdditionInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").Lit("account creation"), jen.ID("UserID").Op(":").ID("exampleUser").Dot("ID"), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAccountCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleCreationInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeCreationAuditLogEntryQuery"), jen.ID("fakeCreationAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeCreationAuditLogEntryQuery"),
						jen.ID("fakeCreationAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationAuditLogEntryArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(123))),
					jen.List(jen.ID("fakeAccountAdditionQuery"), jen.ID("fakeAccountAdditionArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccountAdditionInput"),
					).Dot("Return").Call(
						jen.ID("fakeAccountAdditionQuery"),
						jen.ID("fakeAccountAdditionArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountAdditionQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountAdditionArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeAccountAdditionAuditLogEntryQuery"), jen.ID("fakeAccountAdditionAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeAccountAdditionAuditLogEntryQuery"),
						jen.ID("fakeAccountAdditionAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountAdditionAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountAdditionAuditLogEntryArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(123))),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAccount").Dot("CreatedOn")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleCreationInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAccount"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Call(jen.ID("nil")),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleCreationInput"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.ID("exampleUser").Dot("ID"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Call(jen.ID("nil")),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleCreationInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Valuesln(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAccountCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAccount").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing account creation audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Call(jen.ID("nil")),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAccountCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleCreationInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeCreationAuditLogEntryQuery"), jen.ID("fakeCreationAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeCreationAuditLogEntryQuery"),
						jen.ID("fakeCreationAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationAuditLogEntryArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAccount").Dot("CreatedOn")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleCreationInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing user addition to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Call(jen.ID("nil")),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("exampleAccountAdditionInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").Lit("account creation"), jen.ID("UserID").Op(":").ID("exampleUser").Dot("ID"), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAccountCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleCreationInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeCreationAuditLogEntryQuery"), jen.ID("fakeCreationAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeCreationAuditLogEntryQuery"),
						jen.ID("fakeCreationAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationAuditLogEntryArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(123))),
					jen.List(jen.ID("fakeAccountAdditionQuery"), jen.ID("fakeAccountAdditionArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccountAdditionInput"),
					).Dot("Return").Call(
						jen.ID("fakeAccountAdditionQuery"),
						jen.ID("fakeAccountAdditionArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountAdditionQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountAdditionArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAccount").Dot("CreatedOn")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleCreationInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing user membership addition audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Call(jen.ID("nil")),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("exampleAccountAdditionInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").Lit("account creation"), jen.ID("UserID").Op(":").ID("exampleUser").Dot("ID"), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAccountCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleCreationInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeCreationAuditLogEntryQuery"), jen.ID("fakeCreationAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeCreationAuditLogEntryQuery"),
						jen.ID("fakeCreationAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationAuditLogEntryArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(123))),
					jen.List(jen.ID("fakeAccountAdditionQuery"), jen.ID("fakeAccountAdditionArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccountAdditionInput"),
					).Dot("Return").Call(
						jen.ID("fakeAccountAdditionQuery"),
						jen.ID("fakeAccountAdditionArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountAdditionQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountAdditionArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeAccountAdditionAuditLogEntryQuery"), jen.ID("fakeAccountAdditionAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeAccountAdditionAuditLogEntryQuery"),
						jen.ID("fakeAccountAdditionAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountAdditionAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountAdditionAuditLogEntryArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAccount").Dot("CreatedOn")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleCreationInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("exampleAccount").Dot("Members").Op("=").Index().Op("*").ID("types").Dot("AccountUserMembership").Call(jen.ID("nil")),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("exampleAccount")),
					jen.ID("exampleAccountAdditionInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").Lit("account creation"), jen.ID("UserID").Op(":").ID("exampleUser").Dot("ID"), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAccountCreationQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleCreationInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeCreationAuditLogEntryQuery"), jen.ID("fakeCreationAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeCreationAuditLogEntryQuery"),
						jen.ID("fakeCreationAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationAuditLogEntryArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(123))),
					jen.List(jen.ID("fakeAccountAdditionQuery"), jen.ID("fakeAccountAdditionArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccountAdditionInput"),
					).Dot("Return").Call(
						jen.ID("fakeAccountAdditionQuery"),
						jen.ID("fakeAccountAdditionArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountAdditionQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountAdditionArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeAccountAdditionAuditLogEntryQuery"), jen.ID("fakeAccountAdditionAuditLogEntryArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AuditLogEntrySQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateAuditLogEntryQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("fakeAccountAdditionAuditLogEntryQuery"),
						jen.ID("fakeAccountAdditionAuditLogEntryArgs"),
					).Dot("Once").Call(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountAdditionAuditLogEntryQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountAdditionAuditLogEntryArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.Lit(123))),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleAccount").Dot("CreatedOn")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleCreationInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
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

	code.Add(
		jen.Func().ID("TestQuerier_UpdateAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount"),
							jen.Lit(0),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
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

	code.Add(
		jen.Func().ID("TestQuerier_ArchiveAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("ctx"),
							jen.Lit(0),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.Lit(0),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("exampleUser").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
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

	code.Add(
		jen.Func().ID("TestQuerier_GetAuditLogEntriesForAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAuditLogEntryList").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAuditLogEntriesForAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAuditLogEntries").Call(
						jen.ID("false"),
						jen.ID("exampleAuditLogEntryList").Dot("Entries").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAuditLogEntryList").Dot("Entries"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForAccount").Call(
						jen.ID("ctx"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAuditLogEntriesForAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAuditLogEntriesForAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
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
