package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildMockRowsFromAccountUserMemberships").Params(jen.ID("memberships").Op("...").Op("*").ID("types").Dot("AccountUserMembership")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableColumns")),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("memberships")).Body(
				jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("BelongsToUser"), jen.Op("&").ID("x").Dot("BelongsToAccount"), jen.Qual("strings", "Join").Call(
					jen.ID("x").Dot("AccountRoles"),
					jen.ID("accountMemberRolesSeparator"),
				), jen.Op("&").ID("x").Dot("DefaultAccount"), jen.Op("&").ID("x").Dot("CreatedOn"), jen.Op("&").ID("x").Dot("LastUpdatedOn"), jen.Op("&").ID("x").Dot("ArchivedOn")),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildInvalidMockRowsFromAccountUserMemberships").Params(jen.ID("memberships").Op("...").Op("*").ID("types").Dot("AccountUserMembership")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableColumns")),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("memberships")).Body(
				jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(jen.Op("&").ID("x").Dot("DefaultAccount"), jen.Op("&").ID("x").Dot("BelongsToUser"), jen.Op("&").ID("x").Dot("BelongsToAccount"), jen.Qual("strings", "Join").Call(
					jen.ID("x").Dot("AccountRoles"),
					jen.ID("accountMemberRolesSeparator"),
				), jen.Op("&").ID("x").Dot("CreatedOn"), jen.Op("&").ID("x").Dot("LastUpdatedOn"), jen.Op("&").ID("x").Dot("ArchivedOn"), jen.Op("&").ID("x").Dot("ID")),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ScanAccountUserMemberships").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("surfaces row errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAccountUserMemberships").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
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
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAccountUserMemberships").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
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
		jen.Func().ID("TestQuerier_BuildSessionContextDataForUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("Members").Index(jen.Lit(0)).Dot("DefaultAccount").Op("=").ID("true"),
					jen.ID("examplePermsMap").Op(":=").Map(jen.ID("uint64")).Op("*").ID("types").Dot("UserAccountMembershipInfo").Values(),
					jen.For(jen.List(jen.ID("_"), jen.ID("membership")).Op(":=").Range().ID("exampleAccount").Dot("Members")).Body(
						jen.ID("examplePermsMap").Index(jen.ID("membership").Dot("BelongsToAccount")).Op("=").Op("&").ID("types").Dot("UserAccountMembershipInfo").Valuesln(jen.ID("AccountName").Op(":").ID("exampleAccount").Dot("Name"), jen.ID("AccountID").Op(":").ID("membership").Dot("BelongsToAccount"), jen.ID("AccountRoles").Op(":").ID("membership").Dot("AccountRoles"))),
					jen.ID("exampleAccountPermissionsMap").Op(":=").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Values(),
					jen.For(jen.List(jen.ID("_"), jen.ID("membership")).Op(":=").Range().ID("exampleAccount").Dot("Members")).Body(
						jen.ID("exampleAccountPermissionsMap").Index(jen.ID("membership").Dot("BelongsToAccount")).Op("=").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("membership").Dot("AccountRoles").Op("..."))),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeUserRetrievalQuery"), jen.ID("fakeUserRetrievalArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("UserSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeUserRetrievalQuery"),
						jen.ID("fakeUserRetrievalArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUserRetrievalQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUserRetrievalArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromUsers").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleUser"),
					)),
					jen.List(jen.ID("fakeAccountMembershipsQuery"), jen.ID("fakeAccountMembershipsArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountMembershipsForUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountMembershipsQuery"),
						jen.ID("fakeAccountMembershipsArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountMembershipsQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountMembershipsArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAccountUserMemberships").Call(jen.ID("exampleAccount").Dot("Members").Op("..."))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("expectedActiveAccountID").Op(":=").ID("exampleAccount").Dot("Members").Index(jen.Lit(0)).Dot("BelongsToAccount"),
					jen.ID("expected").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("exampleUser").Dot("ID"), jen.ID("Reputation").Op(":").ID("exampleUser").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("exampleUser").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("exampleUser").Dot("ServiceRoles").Op("..."))), jen.ID("AccountPermissions").Op(":").ID("exampleAccountPermissionsMap"), jen.ID("ActiveAccountID").Op(":").ID("expectedActiveAccountID")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("BuildSessionContextDataForUser").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
						jen.Lit("expected and actual RequestContextData do not match"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("BuildSessionContextDataForUser").Call(
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
				jen.Lit("with error retrieving user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("examplePermsMap").Op(":=").Map(jen.ID("uint64")).Op("*").ID("types").Dot("UserAccountMembershipInfo").Values(),
					jen.For(jen.List(jen.ID("_"), jen.ID("membership")).Op(":=").Range().ID("exampleAccount").Dot("Members")).Body(
						jen.ID("examplePermsMap").Index(jen.ID("membership").Dot("BelongsToAccount")).Op("=").Op("&").ID("types").Dot("UserAccountMembershipInfo").Valuesln(jen.ID("AccountName").Op(":").ID("exampleAccount").Dot("Name"), jen.ID("AccountID").Op(":").ID("membership").Dot("BelongsToAccount"), jen.ID("AccountRoles").Op(":").ID("membership").Dot("AccountRoles"))),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeUserRetrievalQuery"), jen.ID("fakeUserRetrievalArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("UserSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeUserRetrievalQuery"),
						jen.ID("fakeUserRetrievalArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUserRetrievalQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUserRetrievalArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("BuildSessionContextDataForUser").Call(
						jen.ID("ctx"),
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
				jen.Lit("with error retrieving account memberships"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("examplePermsMap").Op(":=").Map(jen.ID("uint64")).Op("*").ID("types").Dot("UserAccountMembershipInfo").Values(),
					jen.For(jen.List(jen.ID("_"), jen.ID("membership")).Op(":=").Range().ID("exampleAccount").Dot("Members")).Body(
						jen.ID("examplePermsMap").Index(jen.ID("membership").Dot("BelongsToAccount")).Op("=").Op("&").ID("types").Dot("UserAccountMembershipInfo").Valuesln(jen.ID("AccountName").Op(":").ID("exampleAccount").Dot("Name"), jen.ID("AccountID").Op(":").ID("membership").Dot("BelongsToAccount"), jen.ID("AccountRoles").Op(":").ID("membership").Dot("AccountRoles"))),
					jen.ID("exampleAccountPermissionsMap").Op(":=").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Values(),
					jen.For(jen.List(jen.ID("_"), jen.ID("membership")).Op(":=").Range().ID("exampleAccount").Dot("Members")).Body(
						jen.ID("exampleAccountPermissionsMap").Index(jen.ID("membership").Dot("BelongsToAccount")).Op("=").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("membership").Dot("AccountRoles").Op("..."))),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeUserRetrievalQuery"), jen.ID("fakeUserRetrievalArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("UserSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeUserRetrievalQuery"),
						jen.ID("fakeUserRetrievalArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUserRetrievalQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUserRetrievalArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromUsers").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleUser"),
					)),
					jen.List(jen.ID("fakeAccountMembershipsQuery"), jen.ID("fakeAccountMembershipsArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountMembershipsForUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountMembershipsQuery"),
						jen.ID("fakeAccountMembershipsArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountMembershipsQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountMembershipsArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("BuildSessionContextDataForUser").Call(
						jen.ID("ctx"),
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
				jen.Lit("with error scanning account user memberships"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("examplePermsMap").Op(":=").Map(jen.ID("uint64")).Op("*").ID("types").Dot("UserAccountMembershipInfo").Values(),
					jen.For(jen.List(jen.ID("_"), jen.ID("membership")).Op(":=").Range().ID("exampleAccount").Dot("Members")).Body(
						jen.ID("examplePermsMap").Index(jen.ID("membership").Dot("BelongsToAccount")).Op("=").Op("&").ID("types").Dot("UserAccountMembershipInfo").Valuesln(jen.ID("AccountName").Op(":").ID("exampleAccount").Dot("Name"), jen.ID("AccountID").Op(":").ID("membership").Dot("BelongsToAccount"), jen.ID("AccountRoles").Op(":").ID("membership").Dot("AccountRoles"))),
					jen.ID("exampleAccountPermissionsMap").Op(":=").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Values(),
					jen.For(jen.List(jen.ID("_"), jen.ID("membership")).Op(":=").Range().ID("exampleAccount").Dot("Members")).Body(
						jen.ID("exampleAccountPermissionsMap").Index(jen.ID("membership").Dot("BelongsToAccount")).Op("=").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("membership").Dot("AccountRoles").Op("..."))),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeUserRetrievalQuery"), jen.ID("fakeUserRetrievalArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("UserSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeUserRetrievalQuery"),
						jen.ID("fakeUserRetrievalArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUserRetrievalQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUserRetrievalArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromUsers").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleUser"),
					)),
					jen.List(jen.ID("fakeAccountMembershipsQuery"), jen.ID("fakeAccountMembershipsArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAccountMembershipsForUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountMembershipsQuery"),
						jen.ID("fakeAccountMembershipsArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountMembershipsQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountMembershipsArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildInvalidMockRowsFromAccountUserMemberships").Call(jen.ID("exampleAccount").Dot("Members").Op("..."))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("BuildSessionContextDataForUser").Call(
						jen.ID("ctx"),
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
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetDefaultAccountIDForUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("expected").Op(":=").ID("exampleAccount").Dot("ID"),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetDefaultAccountIDForUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Valuesln(jen.Lit("id"))).Dot("AddRow").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetDefaultAccountIDForUser").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("db").Dot("ExpectationsWereMet").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetDefaultAccountIDForUser").Call(
						jen.ID("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Zero").Call(
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
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetDefaultAccountIDForUserQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetDefaultAccountIDForUser").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("db").Dot("ExpectationsWereMet").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_MarkAccountAsUserDefault").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildMarkAccountAsUserDefaultQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
						jen.ID("c").Dot("MarkAccountAsUserDefault").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
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
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAccountAsUserDefault").Call(
							jen.ID("ctx"),
							jen.Lit(0),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAccountAsUserDefault").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAccountAsUserDefault").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error marking account as default"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildMarkAccountAsUserDefaultQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("MarkAccountAsUserDefault").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildMarkAccountAsUserDefaultQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
						jen.ID("c").Dot("MarkAccountAsUserDefault").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildMarkAccountAsUserDefaultQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
						jen.ID("c").Dot("MarkAccountAsUserDefault").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_UserIsMemberOfAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUserIsMemberOfAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Valuesln(jen.Lit("result"))).Dot("AddRow").Call(jen.ID("true"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("UserIsMemberOfAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid user ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("UserIsMemberOfAccount").Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("UserIsMemberOfAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error performing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUserIsMemberOfAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("UserIsMemberOfAccount").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("False").Call(
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
		jen.Func().ID("TestQuerier_ModifyUserPermissions").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildModifyUserPermissionsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleInput").Dot("NewRoles"),
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
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
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
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.Lit(0),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account id"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.Lit(0),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildModifyUserPermissionsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleInput").Dot("NewRoles"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildModifyUserPermissionsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleInput").Dot("NewRoles"),
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
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserPermissionModificationInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildModifyUserPermissionsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleInput").Dot("NewRoles"),
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
						jen.ID("c").Dot("ModifyUserPermissions").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_TransferAccountOwnership").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeAccountTransferQuery"), jen.ID("fakeAccountTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountOwnershipQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountTransferQuery"),
						jen.ID("fakeAccountTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountTransferArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeAccountMembershipsTransferQuery"), jen.ID("fakeAccountMembershipsTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountMembershipsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountMembershipsTransferQuery"),
						jen.ID("fakeAccountMembershipsTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountMembershipsTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountMembershipsTransferArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.Lit(0),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error starting transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing account transfer"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeAccountTransferQuery"), jen.ID("fakeAccountTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountOwnershipQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountTransferQuery"),
						jen.ID("fakeAccountTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountTransferArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing membership transfers"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeAccountTransferQuery"), jen.ID("fakeAccountTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountOwnershipQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountTransferQuery"),
						jen.ID("fakeAccountTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountTransferArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeAccountMembershipsTransferQuery"), jen.ID("fakeAccountMembershipsTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountMembershipsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountMembershipsTransferQuery"),
						jen.ID("fakeAccountMembershipsTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountMembershipsTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountMembershipsTransferArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing membership transfers audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeAccountTransferQuery"), jen.ID("fakeAccountTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountOwnershipQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountTransferQuery"),
						jen.ID("fakeAccountTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountTransferArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeAccountMembershipsTransferQuery"), jen.ID("fakeAccountMembershipsTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountMembershipsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountMembershipsTransferQuery"),
						jen.ID("fakeAccountMembershipsTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountMembershipsTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountMembershipsTransferArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeAccountTransferQuery"), jen.ID("fakeAccountTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountOwnershipQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountTransferQuery"),
						jen.ID("fakeAccountTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountTransferArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.List(jen.ID("fakeAccountMembershipsTransferQuery"), jen.ID("fakeAccountMembershipsTransferArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildTransferAccountMembershipsQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("CurrentOwner"),
						jen.ID("exampleInput").Dot("NewOwner"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeAccountMembershipsTransferQuery"),
						jen.ID("fakeAccountMembershipsTransferArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeAccountMembershipsTransferQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeAccountMembershipsTransferArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccount").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("TransferAccountOwnership").Call(
							jen.ID("ctx"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleInput"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_AddUserToAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccountUserMembership").Op(":=").ID("fakes").Dot("BuildFakeAccountUserMembership").Call(),
					jen.ID("exampleAccountUserMembership").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("UserID").Op(":").ID("exampleAccount").Dot("BelongsToUser"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("accountMemberRolesSeparator"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccountUserMembership").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleInput"),
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
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccountUserMembership").Op(":=").ID("fakes").Dot("BuildFakeAccountUserMembership").Call(),
					jen.ID("exampleAccountUserMembership").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("UserID").Op(":").ID("exampleAccount").Dot("BelongsToUser"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("accountMemberRolesSeparator"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleInput"),
							jen.Lit(0),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccountUserMembership").Op(":=").ID("fakes").Dot("BuildFakeAccountUserMembership").Call(),
					jen.ID("exampleAccountUserMembership").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
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
					jen.ID("exampleAccountUserMembership").Op(":=").ID("fakes").Dot("BuildFakeAccountUserMembership").Call(),
					jen.ID("exampleAccountUserMembership").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("UserID").Op(":").ID("exampleAccount").Dot("BelongsToUser"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("accountMemberRolesSeparator"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleInput"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing add query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccountUserMembership").Op(":=").ID("fakes").Dot("BuildFakeAccountUserMembership").Call(),
					jen.ID("exampleAccountUserMembership").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("UserID").Op(":").ID("exampleAccount").Dot("BelongsToUser"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("accountMemberRolesSeparator"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleInput"),
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
					jen.ID("exampleAccountUserMembership").Op(":=").ID("fakes").Dot("BuildFakeAccountUserMembership").Call(),
					jen.ID("exampleAccountUserMembership").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("UserID").Op(":").ID("exampleAccount").Dot("BelongsToUser"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("accountMemberRolesSeparator"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccountUserMembership").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleInput"),
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
					jen.ID("exampleAccountUserMembership").Op(":=").ID("fakes").Dot("BuildFakeAccountUserMembership").Call(),
					jen.ID("exampleAccountUserMembership").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("UserID").Op(":").ID("exampleAccount").Dot("BelongsToUser"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("accountMemberRolesSeparator"))),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildAddUserToAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleAccountUserMembership").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("AddUserToAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleInput"),
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
		jen.Func().ID("TestQuerier_RemoveUserFromAccount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildRemoveUserFromAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
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
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.Lit(0),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.Lit(0),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.Lit(0),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with empty reason"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildRemoveUserFromAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.Lit(""),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing removal to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildRemoveUserFromAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildRemoveUserFromAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("AccountUserMembershipSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildRemoveUserFromAccountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
						jen.ID("c").Dot("RemoveUserFromAccount").Call(
							jen.ID("ctx"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("t").Dot("Name").Call(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
