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
			jen.ID("_").Qual(proj.QuerybuildingPackage(), "AccountUserMembershipSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID(dbvendor.Singular())).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Const().Defs(
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetDefaultAccountIDForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Callln(
					jen.Lit("%s ON %s.%s = %s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				)).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"),
				).Op(":").ID("userID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableDefaultUserAccountColumn"),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildArchiveAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Set").Call(jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"), jen.ID("currentUnixTimeQuery")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetAccountMembershipsForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableColumns").Op("...")).
					Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Callln(
					jen.Lit("%s ON %s.%s = %s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn"),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"),
					).Op(":").ID("userID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildMarkAccountAsUserDefaultQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildMarkAccountAsUserDefaultQuery builds a query that marks a user's account as their primary."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildMarkAccountAsUserDefaultQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"),
				jen.ID("accountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Set").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableDefaultUserAccountColumn"),
					jen.ID("squirrel").Dot("And").Valuesln(
						jen.ID("squirrel").Dot("Eq").Values(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID")),
						jen.ID("squirrel").Dot("Eq").Values(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID")),
					)).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"),
				)),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildTransferAccountOwnershipQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("newOwnerID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableName")).
					Dotln("Set").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn"), jen.ID("newOwnerID")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("accountID"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsTableUserOwnershipColumn").Op(":").ID("currentOwnerID"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildTransferAccountMembershipsQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("currentOwnerID"), jen.ID("newOwnerID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("newOwnerID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Set").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"), jen.ID("newOwnerID")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("currentOwnerID"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildModifyUserPermissionsQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64"),
			jen.ID("newRoles").Index().ID("string")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Set").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountRolesColumn"), jen.Qual("strings", "Join").Call(jen.ID("newRoles"),
					jen.ID("accountMemberRolesSeparator"),
				),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn").Op(":").ID("userID"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn").Op(":").ID("accountID"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildCreateMembershipForNewUserQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableDefaultUserAccountColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountRolesColumn"),
				).Dotln("Values").Callln(
					jen.ID("userID"),
					jen.ID("accountID"),
					jen.ID("true"),
					jen.Qual("strings", "Join").Call(jen.Index().ID("string").Values(jen.Qual(proj.InternalAuthorizationPackage(), "AccountAdminRole").Dot("String").Call()), jen.ID("accountMemberRolesSeparator")),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildUserIsMemberOfAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				)).
					Dotln("Prefix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistencePrefix")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"),
				).Op(":").ID("accountID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"),
					).Op(":").ID("userID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					).Op(":").ID("nil")),
				).Dotln("Suffix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistenceSuffix")),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildAddUserToAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").Op("*").Qual(proj.TypesPackage(), "AddUserToAccountInput")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("UserID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("AccountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountRolesColumn"),
				).
					Dotln("Values").Callln(
					jen.ID("input").Dot("UserID"),
					jen.ID("input").Dot("AccountID"),
					jen.Qual("strings", "Join").Call(jen.ID("input").Dot("AccountRoles"), jen.ID("accountMemberRolesSeparator")),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildRemoveUserFromAccountQuery").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Delete").Call(jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableAccountOwnershipColumn"),
					).Op(":").ID("accountID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableUserOwnershipColumn"),
					).Op(":").ID("userID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "AccountsUserMembershipTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					).Op(":").ID("nil")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}
