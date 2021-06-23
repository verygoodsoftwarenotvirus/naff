package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipsTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestSqlite_BuildGetDefaultAccountIDForUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildUserIsMemberOfAccountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildAddUserToAccountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildRemoveUserFromAccountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildArchiveAccountMembershipsForUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildCreateMembershipForNewUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAccountMembershipsForUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildMarkAccountAsUserDefaultQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildModifyUserPermissionsQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildTransferAccountOwnershipQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildTransferAccountMembershipsQuery(proj, dbvendor)...)

	return code
}

func buildTestSqlite_BuildGetDefaultAccountIDForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetDefaultAccountIDForUserQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT accounts.id FROM accounts JOIN account_user_memberships ON account_user_memberships.belongs_to_account = accounts.id WHERE account_user_memberships.belongs_to_user = ? AND account_user_memberships.default_account = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("true")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetDefaultAccountIDForUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildUserIsMemberOfAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUserIsMemberOfAccountQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT EXISTS ( SELECT account_user_memberships.id FROM account_user_memberships WHERE account_user_memberships.archived_on IS NULL AND account_user_memberships.belongs_to_user = ? )"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUserIsMemberOfAccountQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildAddUserToAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildAddUserToAccountQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").Qual(proj.TypesPackage(), "AddUserToAccountInput").Valuesln(jen.ID("UserID").Op(":").ID("exampleUser").Dot("ID"), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountRoles").Op(":").Index().ID("string").Values(jen.Qual(proj.InternalAuthorizationPackage(), "AccountMemberRole").Dot("String").Call())),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO account_user_memberships (belongs_to_user,belongs_to_account,account_roles) VALUES (?,?,?)"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleInput").Dot("UserID"), jen.ID("exampleAccount").Dot("ID"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleInput").Dot("AccountRoles"),
						jen.ID("accountMemberRolesSeparator"),
					)),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildAddUserToAccountQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildRemoveUserFromAccountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildRemoveUserFromAccountQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("DELETE FROM account_user_memberships WHERE account_user_memberships.archived_on IS NULL AND account_user_memberships.belongs_to_account = ? AND account_user_memberships.belongs_to_user = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildRemoveUserFromAccountQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildArchiveAccountMembershipsForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildArchiveAccountMembershipsForUserQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE account_user_memberships SET archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildArchiveAccountMembershipsForUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildCreateMembershipForNewUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateMembershipForNewUserQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO account_user_memberships (belongs_to_user,belongs_to_account,default_account,account_roles) VALUES (?,?,?,?)"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAccount").Dot("ID"), jen.ID("true"), jen.Qual(proj.InternalAuthorizationPackage(), "AccountAdminRole").Dot("String").Call()),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildCreateMembershipForNewUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAccountMembershipsForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAccountMembershipsForUserQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT account_user_memberships.id, account_user_memberships.belongs_to_user, account_user_memberships.belongs_to_account, account_user_memberships.account_roles, account_user_memberships.default_account, account_user_memberships.created_on, account_user_memberships.last_updated_on, account_user_memberships.archived_on FROM account_user_memberships JOIN accounts ON accounts.id = account_user_memberships.belongs_to_account WHERE account_user_memberships.archived_on IS NULL AND account_user_memberships.belongs_to_user = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAccountMembershipsForUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildMarkAccountAsUserDefaultQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildMarkAccountAsUserDefaultQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE account_user_memberships SET default_account = (belongs_to_user = ? AND belongs_to_account = ?) WHERE archived_on IS NULL AND belongs_to_user = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildMarkAccountAsUserDefaultQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildModifyUserPermissionsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildModifyUserPermissionsQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleRoles").Op(":=").Index().ID("string").Values(jen.Qual(proj.InternalAuthorizationPackage(), "AccountMemberRole").Dot("String").Call()),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE account_user_memberships SET account_roles = ? WHERE belongs_to_account = ? AND belongs_to_user = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.Qual("strings", "Join").Call(
						jen.ID("exampleRoles"),
						jen.ID("accountMemberRolesSeparator"),
					), jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildModifyUserPermissionsQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleRoles"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildTransferAccountOwnershipQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildTransferAccountOwnershipQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleOldOwner").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleNewOwner").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE accounts SET belongs_to_user = ? WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleNewOwner").Dot("ID"), jen.ID("exampleOldOwner").Dot("ID"), jen.ID("exampleAccount").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildTransferAccountOwnershipQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleOldOwner").Dot("ID"),
						jen.ID("exampleNewOwner").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildTransferAccountMembershipsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildTransferAccountMembershipsQuery", dbvendor.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleOldOwner").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleNewOwner").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE account_user_memberships SET belongs_to_user = ? WHERE archived_on IS NULL AND belongs_to_account = ? AND belongs_to_user = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleNewOwner").Dot("ID"), jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleOldOwner").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildTransferAccountMembershipsQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleOldOwner").Dot("ID"),
						jen.ID("exampleNewOwner").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
					jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}
