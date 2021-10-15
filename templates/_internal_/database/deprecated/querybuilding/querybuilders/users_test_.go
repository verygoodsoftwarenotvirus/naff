package querybuilders

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

func usersTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestSqlite_BuildUserIsBannedQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetUserWithUnverifiedTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetUsersQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildTestUserCreationQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetUserByUsernameQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildSearchForUserByUsernameQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAllUsersCountQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildCreateUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildUpdateUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildUpdateUserPasswordQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildSetUserStatusQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildUpdateUserTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildVerifyUserTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildArchiveUserQuery(proj, dbvendor)...)
	code.Add(buildTestSqlite_BuildGetAuditLogEntriesForUserQuery(proj, dbvendor)...)

	return code
}

func buildTestSqlite_BuildUserIsBannedQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUserIsBannedQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleStatuses").Assign().Index().String().Valuesln(
						jen.String().Call(jen.Qual(proj.TypesPackage(), "BannedUserAccountStatus")),
						jen.String().Call(jen.Qual(proj.TypesPackage(), "TerminatedUserReputation")),
					),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("SELECT EXISTS ( SELECT users.id FROM users WHERE users.archived_on IS NULL AND users.id = %s AND (users.reputation = %s OR users.reputation = %s) )", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.String().Call(jen.Qual(proj.TypesPackage(), "BannedUserAccountStatus")), jen.String().Call(jen.Qual(proj.TypesPackage(), "TerminatedUserReputation"))),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUserHasStatusQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleStatuses").Spread(),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Litf("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.id = %s AND users.two_factor_secret_verified_on IS NOT NULL", getIncIndex(dbvendor, 0)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetUserWithUnverifiedTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetUserWithUnverifiedTwoFactorSecretQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Litf("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.id = %s AND users.two_factor_secret_verified_on IS NULL", getIncIndex(dbvendor, 0)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetUserWithUnverifiedTwoFactorSecretQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetUsersQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetUsersQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("filter").Assign().Qual(proj.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on, (SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL) as total_count, (SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL AND users.created_on > %s AND users.created_on < %s AND users.last_updated_on > %s AND users.last_updated_on < %s) as filtered_count FROM users WHERE users.archived_on IS NULL AND users.created_on > %s AND users.created_on < %s AND users.last_updated_on > %s AND users.last_updated_on < %s GROUP BY users.id LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetUsersQuery").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildTestUserCreationQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	qb := queryBuilderForDatabase(dbvendor).Insert("users").
		Columns(
			"external_id",
			"username",
			"hashed_password",
			"two_factor_secret",
			"reputation",
			"service_roles",
			"two_factor_secret_verified_on",
		).
		Values(
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			whateverValue,
			squirrel.Expr(unixTimeForDatabase(dbvendor)),
		)

	if dbvendor.SingularPackageName() == "postgres" {
		qb = qb.Suffix("RETURNING id")
	}

	expectedQuery, _ := buildQuery(qb)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildTestUserCreationQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("fakeUUID").Assign().Lit("blahblah"),
					jen.ID("mockExternalIDGenerator").Assign().AddressOf().Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("mockExternalIDGenerator").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("fakeUUID")),
					jen.ID("q").Dot("externalIDGenerator").Equals().ID("mockExternalIDGenerator"),
					jen.Newline(),
					jen.ID("exampleInput").Assign().AddressOf().Qual(proj.TypesPackage(), "TestUserCreationConfig").Valuesln(
						jen.ID("Username").MapAssign().Lit("username"),
						jen.ID("Password").MapAssign().Lit("password"),
						jen.ID("HashedPassword").MapAssign().Lit("hashashashash"),
						jen.ID("IsServiceAdmin").MapAssign().ID("true"),
					),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Lit(expectedQuery),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("fakeUUID"), jen.ID("exampleInput").Dot("Username"), jen.ID("exampleInput").Dot("HashedPassword"), jen.Qual(proj.QuerybuildingPackage(), "DefaultTestUserTwoFactorSecret"), jen.Qual(proj.TypesPackage(), "GoodStandingAccountStatus"), jen.Qual(proj.InternalAuthorizationPackage(), "ServiceAdminRole").Dot("String").Call()),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildTestUserCreationQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockExternalIDGenerator")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetUserByUsernameQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Litf("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.username = %s AND users.two_factor_secret_verified_on IS NOT NULL", getIncIndex(dbvendor, 0)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("Username")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetUserByUsernameQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("Username"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildSearchForUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var postgresPrefix string
	if dbvendor.SingularPackageName() == "postgres" {
		postgresPrefix = "I"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildSearchForUserByUsernameQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Litf("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.username %sLIKE %s AND users.archived_on IS NULL AND users.two_factor_secret_verified_on IS NOT NULL", postgresPrefix, getIncIndex(dbvendor, 0)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s%%"),
						jen.ID("exampleUser").Dot("Username"),
					)),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildSearchForUserByUsernameQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("Username"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAllUsersCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAllUsersCountQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Lit("SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL"),
					jen.ID("actualQuery").Assign().ID("q").Dot("BuildGetAllUsersCountQuery").Call(jen.ID("ctx")),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildCreateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var querySuffix string
	if dbvendor.SingularPackageName() == "postgres" {
		querySuffix = " RETURNING id"
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildCreateUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUserDataStoreCreationInputFromUser").Call(jen.ID("exampleUser")),
					jen.Newline(),
					jen.ID("exIDGen").Assign().AddressOf().Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleUser").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Equals().ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("INSERT INTO users (external_id,username,hashed_password,two_factor_secret,reputation,service_roles) VALUES (%s,%s,%s,%s,%s,%s)%s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), querySuffix),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ExternalID"), jen.ID("exampleUser").Dot("Username"), jen.ID("exampleUser").Dot("HashedPassword"), jen.ID("exampleUser").Dot("TwoFactorSecret"), jen.Qual(proj.TypesPackage(), "UnverifiedAccountStatus"), jen.Qual(proj.InternalAuthorizationPackage(), "ServiceUserRole").Dot("String").Call()),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildCreateUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("exIDGen"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildUpdateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("users").
			Set("username", whateverValue).
			Set("hashed_password", whateverValue).
			Set("avatar_src", whateverValue).
			Set("two_factor_secret", whateverValue).
			Set("two_factor_secret_verified_on", whateverValue).
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":          whateverValue,
				"archived_on": nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUpdateUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(
						jen.ID("exampleUser").Dot("Username"),
						jen.ID("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleUser").Dot("AvatarSrc"),
						jen.ID("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUpdateUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildUpdateUserPasswordQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("users").
			Set("hashed_password", whateverValue).
			Set("requires_password_change", false).
			Set("password_last_changed_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":          whateverValue,
				"archived_on": nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUpdateUserPasswordQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("HashedPassword"), jen.ID("false"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUpdateUserPasswordQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleUser").Dot("HashedPassword"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildSetUserStatusQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildSetUserStatusQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUserReputationUpdateInputFromUser").Call(jen.ID("exampleUser")),
					jen.Newline(),
					jen.ID("expectedQuery").Assign().Litf("UPDATE users SET reputation = %s, reputation_explanation = %s WHERE archived_on IS NULL AND id = %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleInput").Dot("NewReputation"), jen.ID("exampleInput").Dot("Reason"), jen.ID("exampleInput").Dot("TargetUserID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildSetUserStatusQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildUpdateUserTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildUpdateUserTwoFactorSecretQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedQuery").Assign().Litf("UPDATE users SET two_factor_secret_verified_on = %s, two_factor_secret = %s WHERE archived_on IS NULL AND id = %s", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2)),
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.Nil(), jen.ID("exampleUser").Dot("TwoFactorSecret"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildUpdateUserTwoFactorSecretQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleUser").Dot("TwoFactorSecret"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildVerifyUserTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("users").
			Set("two_factor_secret_verified_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Set("reputation", whateverValue).
			Where(squirrel.Eq{
				"id":          whateverValue,
				"archived_on": nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildVerifyUserTwoFactorSecretQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.Qual(proj.TypesPackage(), "GoodStandingAccountStatus"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildVerifyUserTwoFactorSecretQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildArchiveUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	expectedQuery, _ := buildQuery(
		queryBuilderForDatabase(dbvendor).Update("users").
			Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
			Where(squirrel.Eq{
				"id":          whateverValue,
				"archived_on": nil,
			}),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildArchiveUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildArchiveUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestSqlite_BuildGetAuditLogEntriesForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	var (
		userIDKey        string
		performedByIDKey string
	)

	switch dbvendor.LowercaseAbbreviation() {
	case "m":
		userIDKey = fmt.Sprintf(`JSON_CONTAINS(%s.%s, '%s', '$.%s')`, "audit_log", "context", "%d", "user_id")
		performedByIDKey = fmt.Sprintf(`JSON_CONTAINS(%s.%s, '%s', '$.%s')`, "audit_log", "context", "%d", "performed_by")
	case "p":
		userIDKey = fmt.Sprintf(`%s.%s->'%s'`, "audit_log", "context", "user_id")
		performedByIDKey = fmt.Sprintf(`%s.%s->'%s'`, "audit_log", "context", "performed_by")
	case "s":
		userIDKey = fmt.Sprintf(`json_extract(%s.%s, '$.%s')`, "audit_log", "context", "user_id")
		performedByIDKey = fmt.Sprintf(`json_extract(%s.%s, '$.%s')`, "audit_log", "context", "performed_by")
	}

	queryBuilder := queryBuilderForDatabase(dbvendor).Select(
		"audit_log.id",
		"audit_log.external_id",
		"audit_log.event_type",
		"audit_log.context",
		"audit_log.created_on",
	).
		From("audit_log")

	if dbvendor.SingularPackageName() == "mysql" {
		queryBuilder = queryBuilder.Where(squirrel.Or{squirrel.Expr(userIDKey), squirrel.Expr(performedByIDKey)})
	} else {
		queryBuilder = queryBuilder.Where(squirrel.Or{squirrel.Eq{userIDKey: whateverValue}, squirrel.Eq{performedByIDKey: whateverValue}})
	}

	queryBuilder = queryBuilder.OrderBy("audit_log.created_on")

	expectedQuery, _ := buildQuery(queryBuilder)

	expectedQueryDecl := jen.ID("expectedQuery").Assign().Lit(expectedQuery)
	expectedArgs := jen.ID("expectedArgs").Assign().Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"))

	if dbvendor.SingularPackageName() == "mysql" {
		expectedQueryDecl = jen.ID("expectedQuery").Assign().Qual("fmt", "Sprintf").Call(jen.Lit(expectedQuery), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"))
		expectedArgs = jen.ID("expectedArgs").Assign().Index().Interface().Call(jen.Nil())
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_BuildGetAuditLogEntriesForUserQuery", dbvendor.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
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
					expectedQueryDecl,
					expectedArgs,
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Assign().ID("q").Dot("BuildGetAuditLogEntriesForUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}
