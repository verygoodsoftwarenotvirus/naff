package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("AccountUserMembershipSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("Postgres")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("accountMemberRolesSeparator").Op("=").Lit(","),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetDefaultAccountIDForUserQuery does ."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetDefaultAccountIDForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)).Dot("From").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Join").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s ON %s.%s = %s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableDefaultUserAccountColumn"),
				).Op(":").ID("true"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveAccountMembershipsForUserQuery does ."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildArchiveAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAccountMembershipsForUserQuery does ."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableColumns").Op("...")).Dot("Join").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s ON %s.%s = %s.%s"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("AccountsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn"),
				)).Dot("From").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
				).Op(":").ID("userID"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateMembershipForNewUserQuery builds a query that ."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildCreateMembershipForNewUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableDefaultUserAccountColumn"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountRolesColumn"),
				).Dot("Values").Call(
					jen.ID("userID"),
					jen.ID("accountID"),
					jen.ID("true"),
					jen.Qual("strings", "Join").Call(
						jen.Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call()),
						jen.ID("accountMemberRolesSeparator"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildMarkAccountAsUserDefaultQuery builds a query that marks a user's account as their primary."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildMarkAccountAsUserDefaultQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableDefaultUserAccountColumn"),
					jen.ID("squirrel").Dot("And").Valuesln(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID")), jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID"))),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildModifyUserPermissionsQuery builds."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildModifyUserPermissionsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64"), jen.ID("newRoles").Index().ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountRolesColumn"),
					jen.Qual("strings", "Join").Call(
						jen.ID("newRoles"),
						jen.ID("accountMemberRolesSeparator"),
					),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildTransferAccountOwnershipQuery builds."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildTransferAccountOwnershipQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("newOwnerID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn"),
					jen.ID("newOwnerID"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("accountID"), jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn").Op(":").ID("currentOwnerID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildTransferAccountMembershipsQuery does ."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildTransferAccountMembershipsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("newOwnerID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
					jen.ID("newOwnerID"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID"), jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("currentOwnerID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserIsMemberOfAccountQuery builds a query that checks to see if the user is the member of a given account."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildUserIsMemberOfAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)).Dot("Prefix").Call(jen.ID("querybuilding").Dot("ExistencePrefix")).Dot("From").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))).Dot("Suffix").Call(jen.ID("querybuilding").Dot("ExistenceSuffix")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAddUserToAccountQuery builds a query that adds a user to an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildAddUserToAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("UserID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("AccountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountRolesColumn"),
				).Dot("Values").Call(
					jen.ID("input").Dot("UserID"),
					jen.ID("input").Dot("AccountID"),
					jen.Qual("strings", "Join").Call(
						jen.ID("input").Dot("AccountRoles"),
						jen.ID("accountMemberRolesSeparator"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildRemoveUserFromAccountQuery builds a query that removes a user from an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildRemoveUserFromAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Delete").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	return code
}
