package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
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
	query, _ := buildQuery(queryBuilderForDatabase(dbvendor).
		Select("accounts.id").
		From("accounts").
		Join("account_user_memberships ON account_user_memberships.belongs_to_account = accounts.id").
		Where(squirrel.Eq{
			"account_user_memberships.belongs_to_user": whateverValue,
			"account_user_memberships.default_account": true,
		}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetDefaultAccountIDForUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(query),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("true")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetDefaultAccountIDForUserQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select("account_user_memberships.id").
			Prefix(existencePrefix).
			From("account_user_memberships").
			Where(squirrel.Eq{
				"account_user_memberships.belongs_to_account": whateverValue,
				"account_user_memberships.belongs_to_user":    whateverValue,
				"account_user_memberships.archived_on":        nil,
			}).
			Suffix(")"),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUserIsMemberOfAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUserIsMemberOfAccountQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Insert("account_user_memberships").
			Columns(
				"belongs_to_user",
				"belongs_to_account",
				"account_roles",
			).
			Values(
				whateverValue,
				whateverValue,
				whateverValue,
			),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildAddUserToAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.ID("exampleInput").Assign().AddressOf().Qual(proj.TypesPackage(), "AddUserToAccountInput").Valuesln(jen.ID("UserID").Op(":").ID("exampleUser").Dot("ID"), jen.ID("AccountID").Op(":").ID("exampleAccount").Dot("ID"), jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountRoles").Op(":").Index().String().Values(jen.Qual(proj.InternalAuthorizationPackage(), "AccountMemberRole").Dot("String").Call())),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleInput").Dot("UserID"), jen.ID("exampleAccount").Dot("ID"), jen.Qual("strings", "Join").Call(
						jen.ID("exampleInput").Dot("AccountRoles"),
						jen.ID("accountMemberRolesSeparator"),
					)),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildAddUserToAccountQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Delete("account_user_memberships").
			Where(squirrel.Eq{
				"account_user_memberships.belongs_to_account": whateverValue,
				"account_user_memberships.belongs_to_user":    whateverValue,
				"account_user_memberships.archived_on":        nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildRemoveUserFromAccountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildRemoveUserFromAccountQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("account_user_memberships").
			Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"belongs_to_user": whateverValue,
				"archived_on":     nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildArchiveAccountMembershipsForUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildArchiveAccountMembershipsForUserQuery").Call(
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
	exampleQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Insert("account_user_memberships").
			Columns(
				"belongs_to_user",
				"belongs_to_account",
				"default_account",
				"account_roles",
			).
			Values(
				whateverValue,
				whateverValue,
				true,
				whateverValue,
			),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateMembershipForNewUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(exampleQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAccount").Dot("ID"), jen.ID("true"), jen.Qual(proj.InternalAuthorizationPackage(), "AccountAdminRole").Dot("String").Call()),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildCreateMembershipForNewUserQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Select(
			"account_user_memberships.id",
			"account_user_memberships.belongs_to_user",
			"account_user_memberships.belongs_to_account",
			"account_user_memberships.account_roles",
			"account_user_memberships.default_account",
			"account_user_memberships.created_on",
			"account_user_memberships.last_updated_on",
			"account_user_memberships.archived_on",
		).
			Join("accounts ON accounts.id = account_user_memberships.belongs_to_account").
			From("account_user_memberships").
			Where(squirrel.Eq{
				"account_user_memberships.archived_on":     nil,
				"account_user_memberships.belongs_to_user": whateverValue,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAccountMembershipsForUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAccountMembershipsForUserQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("account_user_memberships").
			Set("default_account", squirrel.And{
				squirrel.Eq{"belongs_to_user": whateverValue},
				squirrel.Eq{"belongs_to_account": whateverValue},
			}).
			Where(squirrel.Eq{
				"belongs_to_user": whateverValue,
				"archived_on":     nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildMarkAccountAsUserDefaultQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildMarkAccountAsUserDefaultQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("account_user_memberships").
			Set("account_roles", whateverValue).
			Where(squirrel.Eq{
				"belongs_to_user":    whateverValue,
				"belongs_to_account": whateverValue,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildModifyUserPermissionsQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleRoles").Assign().Index().String().Values(jen.Qual(proj.InternalAuthorizationPackage(), "AccountMemberRole").Dot("String").Call()),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.Qual("strings", "Join").Call(
						jen.ID("exampleRoles"),
						jen.ID("accountMemberRolesSeparator"),
					), jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildModifyUserPermissionsQuery").Call(
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
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("accounts").
			Set("belongs_to_user", whateverValue).
			Where(squirrel.Eq{
				"id":              whateverValue,
				"belongs_to_user": whateverValue,
				"archived_on":     nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildTransferAccountOwnershipQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleOldOwner").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleNewOwner").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleNewOwner").Dot("ID"), jen.ID("exampleOldOwner").Dot("ID"), jen.ID("exampleAccount").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildTransferAccountOwnershipQuery").Call(
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
	exampleQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("account_user_memberships").
			Set("belongs_to_user", whateverValue).
			Where(squirrel.Eq{
				"belongs_to_account": whateverValue,
				"belongs_to_user":    whateverValue,
				"archived_on":        nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildTransferAccountMembershipsQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleOldOwner").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleNewOwner").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(exampleQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleNewOwner").Dot("ID"), jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleOldOwner").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildTransferAccountMembershipsQuery").Call(
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
