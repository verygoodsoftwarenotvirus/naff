package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipsDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("AccountUserMembershipSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("Sqlite")).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("accountMemberRolesSeparator").Op("=").Lit(","),
		),
		jen.Newline(),
	)

	code.Add(buildBuildGetDefaultAccountIDForUserQuery(proj, dbvendor)...)
	code.Add(buildBuildArchiveAccountMembershipsForUserQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAccountMembershipsForUserQuery(proj, dbvendor)...)
	code.Add(buildBuildMarkAccountAsUserDefaultQuery(proj, dbvendor)...)
	code.Add(buildBuildTransferAccountOwnershipQuery(proj, dbvendor)...)
	code.Add(buildBuildTransferAccountMembershipsQuery(proj, dbvendor)...)
	code.Add(buildBuildModifyUserPermissionsQuery(proj, dbvendor)...)
	code.Add(buildBuildCreateMembershipForNewUserQuery(proj, dbvendor)...)
	code.Add(buildBuildUserIsMemberOfAccountQuery(proj, dbvendor)...)
	code.Add(buildBuildAddUserToAccountQuery(proj, dbvendor)...)
	code.Add(buildBuildRemoveUserFromAccountQuery(proj, dbvendor)...)

	return code
}

func buildBuildGetDefaultAccountIDForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetDefaultAccountIDForUserQuery does ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetDefaultAccountIDForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
		jen.Newline(),
	}

	return lines
}

func buildBuildArchiveAccountMembershipsForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildArchiveAccountMembershipsForUserQuery does ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildArchiveAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAccountMembershipsForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAccountMembershipsForUserQuery does ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
		jen.Newline(),
	}

	return lines
}

func buildBuildMarkAccountAsUserDefaultQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildMarkAccountAsUserDefaultQuery does ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildMarkAccountAsUserDefaultQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableDefaultUserAccountColumn"),
					jen.ID("squirrel").Dot("And").Valuesln(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID")), jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID"))),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildTransferAccountOwnershipQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildTransferAccountOwnershipQuery does ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildTransferAccountOwnershipQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn"),
					jen.ID("newOwnerID"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("accountID"), jen.ID("querybuilding").Dot("AccountsTableUserOwnershipColumn").Op(":").ID("currentOwnerID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildTransferAccountMembershipsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildTransferAccountMembershipsQuery does ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildTransferAccountMembershipsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("AccountsUserMembershipTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn"),
					jen.ID("newOwnerID"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID"), jen.ID("querybuilding").Dot("AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("currentOwnerID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildModifyUserPermissionsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildModifyUserPermissionsQuery builds."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildModifyUserPermissionsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64"), jen.ID("newRoles").Index().ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
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
		jen.Newline(),
	}

	return lines
}

func buildBuildCreateMembershipForNewUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildCreateMembershipForNewUserQuery builds a query that ."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildCreateMembershipForNewUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
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
		jen.Newline(),
	}

	return lines
}

func buildBuildUserIsMemberOfAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUserIsMemberOfAccountQuery builds a query that checks to see if the user is the member of a given account."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildUserIsMemberOfAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
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
		jen.Newline(),
	}

	return lines
}

func buildBuildAddUserToAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildAddUserToAccountQuery builds a query that adds a user to an account."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildAddUserToAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
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
		jen.Newline(),
	}

	return lines
}

func buildBuildRemoveUserFromAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildRemoveUserFromAccountQuery builds a query that removes a user from an account."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildRemoveUserFromAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
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
		jen.Newline(),
	}

	return lines
}
