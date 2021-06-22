package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
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
		jen.Func().ID("TestSqlite_BuildUserIsBannedQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT EXISTS ( SELECT users.id FROM users WHERE users.archived_on IS NULL AND users.id = ? AND (users.reputation = ? OR users.reputation = ?) )"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("string").Call(jen.ID("types").Dot("BannedUserAccountStatus")), jen.ID("string").Call(jen.ID("types").Dot("TerminatedUserReputation"))),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUserHasStatusQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("string").Call(jen.ID("types").Dot("BannedUserAccountStatus")),
						jen.ID("string").Call(jen.ID("types").Dot("TerminatedUserReputation")),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
		jen.Func().ID("TestSqlite_BuildGetUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.id = ? AND users.two_factor_secret_verified_on IS NOT NULL"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
		jen.Func().ID("TestSqlite_BuildGetUserWithUnverifiedTwoFactorSecretQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.id = ? AND users.two_factor_secret_verified_on IS NULL"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetUserWithUnverifiedTwoFactorSecretQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
		jen.Func().ID("TestSqlite_BuildGetUsersQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("filter").Op(":=").ID("fakes").Dot("BuildFleshedOutQueryFilter").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on, (SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL) as total_count, (SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL AND users.created_on > ? AND users.created_on < ? AND users.last_updated_on > ? AND users.last_updated_on < ?) as filtered_count FROM users WHERE users.archived_on IS NULL AND users.created_on > ? AND users.created_on < ? AND users.last_updated_on > ? AND users.last_updated_on < ? GROUP BY users.id LIMIT 20 OFFSET 180"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetUsersQuery").Call(
						jen.ID("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildTestUserCreationQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("fakeUUID").Op(":=").Lit("blahblah"),
					jen.ID("mockExternalIDGenerator").Op(":=").Op("&").ID("querybuilding").Dot("MockExternalIDGenerator").Values(),
					jen.ID("mockExternalIDGenerator").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("fakeUUID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("mockExternalIDGenerator"),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("HashedPassword").Op(":").Lit("hashashashash"), jen.ID("IsServiceAdmin").Op(":").ID("true")),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO users (external_id,username,hashed_password,two_factor_secret,reputation,service_roles,two_factor_secret_verified_on) VALUES (?,?,?,?,?,?,(strftime('%s','now')))"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("fakeUUID"), jen.ID("exampleInput").Dot("Username"), jen.ID("exampleInput").Dot("HashedPassword"), jen.ID("querybuilding").Dot("DefaultTestUserTwoFactorSecret"), jen.ID("types").Dot("GoodStandingAccountStatus"), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildTestUserCreationQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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

func buildTestSqlite_BuildGetUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetUserByUsernameQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.username = ? AND users.two_factor_secret_verified_on IS NOT NULL"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("Username")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetUserByUsernameQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("Username"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildSearchForUserByUsernameQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT users.id, users.external_id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.username LIKE ? AND users.archived_on IS NULL AND users.two_factor_secret_verified_on IS NOT NULL"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s%%"),
						jen.ID("exampleUser").Dot("Username"),
					)),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildSearchForUserByUsernameQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("Username"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
		jen.Func().ID("TestSqlite_BuildGetAllUsersCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(users.id) FROM users WHERE users.archived_on IS NULL"),
					jen.ID("actualQuery").Op(":=").ID("q").Dot("BuildGetAllUsersCountQuery").Call(jen.ID("ctx")),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("assert").Dot("Equal").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildCreateUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserDataStoreCreationInputFromUser").Call(jen.ID("exampleUser")),
					jen.ID("exIDGen").Op(":=").Op("&").ID("querybuilding").Dot("MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleUser").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO users (external_id,username,hashed_password,two_factor_secret,reputation,service_roles) VALUES (?,?,?,?,?,?)"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ExternalID"), jen.ID("exampleUser").Dot("Username"), jen.ID("exampleUser").Dot("HashedPassword"), jen.ID("exampleUser").Dot("TwoFactorSecret"), jen.ID("types").Dot("UnverifiedAccountStatus"), jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call()),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildCreateUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildUpdateUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET username = ?, hashed_password = ?, avatar_src = ?, two_factor_secret = ?, two_factor_secret_verified_on = ?, last_updated_on = (strftime('%s','now')) WHERE archived_on IS NULL AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("Username"), jen.ID("exampleUser").Dot("HashedPassword"), jen.ID("exampleUser").Dot("AvatarSrc"), jen.ID("exampleUser").Dot("TwoFactorSecret"), jen.ID("exampleUser").Dot("TwoFactorSecretVerifiedOn"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUpdateUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildUpdateUserPasswordQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET hashed_password = ?, requires_password_change = ?, password_last_changed_on = (strftime('%s','now')), last_updated_on = (strftime('%s','now')) WHERE archived_on IS NULL AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("HashedPassword"), jen.ID("false"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUpdateUserPasswordQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleUser").Dot("HashedPassword"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
		jen.Func().ID("TestSqlite_BuildSetUserStatusQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInputFromUser").Call(jen.ID("exampleUser")),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET reputation = ?, reputation_explanation = ? WHERE archived_on IS NULL AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleInput").Dot("NewReputation"), jen.ID("exampleInput").Dot("Reason"), jen.ID("exampleInput").Dot("TargetUserID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildSetUserStatusQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
		jen.Func().ID("TestSqlite_BuildUpdateUserTwoFactorSecretQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET two_factor_secret_verified_on = ?, two_factor_secret = ? WHERE archived_on IS NULL AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("nil"), jen.ID("exampleUser").Dot("TwoFactorSecret"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUpdateUserTwoFactorSecretQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("exampleUser").Dot("TwoFactorSecret"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildVerifyUserTwoFactorSecretQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET two_factor_secret_verified_on = (strftime('%s','now')), reputation = ? WHERE archived_on IS NULL AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("types").Dot("GoodStandingAccountStatus"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildVerifyUserTwoFactorSecretQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildArchiveUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE users SET archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildArchiveUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
	lines := []jen.Code{
		jen.Func().ID("TestSqlite_BuildGetAuditLogEntriesForUserQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("expectedQuery").Op(":=").Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE (json_extract(audit_log.context, '$.user_id') = ? OR json_extract(audit_log.context, '$.performed_by') = ?) ORDER BY audit_log.created_on"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAuditLogEntriesForUserQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
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
