package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("UserSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("Sqlite")).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(buildBuildUserHasStatusQuery(proj, dbvendor)...)
	code.Add(buildBuildGetUserQuery(proj, dbvendor)...)
	code.Add(buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildBuildGetUserByUsernameQuery(proj, dbvendor)...)
	code.Add(buildBuildSearchForUserByUsernameQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAllUsersCountQuery(proj, dbvendor)...)
	code.Add(buildBuildGetUsersQuery(proj, dbvendor)...)
	code.Add(buildBuildTestUserCreationQuery(proj, dbvendor)...)
	code.Add(buildBuildCreateUserQuery(proj, dbvendor)...)
	code.Add(buildBuildUpdateUserQuery(proj, dbvendor)...)
	code.Add(buildBuildSetUserStatusQuery(proj, dbvendor)...)
	code.Add(buildBuildUpdateUserPasswordQuery(proj, dbvendor)...)
	code.Add(buildBuildUpdateUserTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildBuildVerifyUserTwoFactorSecretQuery(proj, dbvendor)...)
	code.Add(buildBuildArchiveUserQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAuditLogEntriesForUserQuery(proj, dbvendor)...)

	return code
}

func buildBuildUserHasStatusQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUserHasStatusQuery returns a SQL query (and argument) for retrieving a user by their database ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildUserHasStatusQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("statuses").Op("...").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("whereStatuses").Op(":=").ID("squirrel").Dot("Or").Values(),
			jen.For(jen.List(jen.ID("_"), jen.ID("status")).Op(":=").Range().ID("statuses")).Body(
				jen.ID("whereStatuses").Op("=").ID("append").Call(
					jen.ID("whereStatuses"),
					jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("querybuilding").Dot("UsersTableName"),
						jen.ID("querybuilding").Dot("UsersTableReputationColumn"),
					).Op(":").ID("status")),
				)),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)).Dot("Prefix").Call(jen.ID("querybuilding").Dot("ExistencePrefix")).Dot("From").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))).Dot("Where").Call(jen.ID("whereStatuses")).Dot("Suffix").Call(jen.ID("querybuilding").Dot("ExistenceSuffix")),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("UsersTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))).Dot("Where").Call(jen.ID("squirrel").Dot("NotEq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user"),
		jen.Newline(),
		jen.Comment("by their database ID, who has an unverified two factor secret."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetUserWithUnverifiedTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("UsersTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
				).Op(":").ID("nil"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetUserByUsernameQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("UsersTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("UsersTableUsernameColumn"),
				).Op(":").ID("username"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))).Dot("Where").Call(jen.ID("squirrel").Dot("NotEq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildSearchForUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildSearchForUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildSearchForUserByUsernameQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("usernameQuery").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("UsersTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Expr").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s LIKE ?"),
						jen.ID("querybuilding").Dot("UsersTableName"),
						jen.ID("querybuilding").Dot("UsersTableUsernameColumn"),
					),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s%%"),
						jen.ID("usernameQuery"),
					),
				)).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))).Dot("Where").Call(jen.ID("squirrel").Dot("NotEq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllUsersCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere"),
		jen.Newline(),
		jen.Comment("to a given filter's criteria."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAllUsersCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("query").ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.ID("querybuilding").Dot("UsersTableName"),
				)).Dot("From").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("UsersTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetUsersQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere"),
		jen.Newline(),
		jen.Comment("to a given filter's criteria."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetUsersQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				)),
			jen.Return().ID("b").Dot("buildListQuery").Call(
				jen.ID("ctx"),
				jen.ID("querybuilding").Dot("UsersTableName"),
				jen.ID("nil"),
				jen.ID("nil"),
				jen.Lit(""),
				jen.ID("querybuilding").Dot("UsersTableColumns"),
				jen.Lit(0),
				jen.ID("false"),
				jen.ID("filter"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildTestUserCreationQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildTestUserCreationQuery builds a query and arguments that creates a test user."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildTestUserCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("testUserConfig").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("serviceRole").Op(":=").ID("authorization").Dot("ServiceUserRole"),
			jen.If(jen.ID("testUserConfig").Dot("IsServiceAdmin")).Body(
				jen.ID("serviceRole").Op("=").ID("authorization").Dot("ServiceAdminRole")),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("ExternalIDColumn"),
					jen.ID("querybuilding").Dot("UsersTableUsernameColumn"),
					jen.ID("querybuilding").Dot("UsersTableHashedPasswordColumn"),
					jen.ID("querybuilding").Dot("UsersTableTwoFactorSekretColumn"),
					jen.ID("querybuilding").Dot("UsersTableReputationColumn"),
					jen.ID("querybuilding").Dot("UsersTableServiceRolesColumn"),
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
				).Dot("Values").Call(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("testUserConfig").Dot("Username"),
					jen.ID("testUserConfig").Dot("HashedPassword"),
					jen.ID("querybuilding").Dot("DefaultTestUserTwoFactorSecret"),
					jen.ID("types").Dot("GoodStandingAccountStatus"),
					jen.ID("serviceRole").Dot("String").Call(),
					jen.ID("currentUnixTimeQuery"),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildCreateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildCreateUserQuery returns a SQL query (and arguments) that would create a given Requester."),
		jen.Newline(),
		jen.Comment("NOTE: we always default is_admin to false, on the assumption that"),
		jen.Newline(),
		jen.Comment("admins have DB access and will change that value via SQL query."),
		jen.Newline(),
		jen.Comment("There should be no way to update a user via this structure"),
		jen.Newline(),
		jen.Comment("such that they would have admin privileges."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildCreateUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserDataStoreCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("ExternalIDColumn"),
					jen.ID("querybuilding").Dot("UsersTableUsernameColumn"),
					jen.ID("querybuilding").Dot("UsersTableHashedPasswordColumn"),
					jen.ID("querybuilding").Dot("UsersTableTwoFactorSekretColumn"),
					jen.ID("querybuilding").Dot("UsersTableReputationColumn"),
					jen.ID("querybuilding").Dot("UsersTableServiceRolesColumn"),
				).Dot("Values").Call(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Username"),
					jen.ID("input").Dot("HashedPassword"),
					jen.ID("input").Dot("TwoFactorSecret"),
					jen.ID("types").Dot("UnverifiedAccountStatus"),
					jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call(),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildUpdateUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("User")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableUsernameColumn"),
					jen.ID("input").Dot("Username"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableHashedPasswordColumn"),
					jen.ID("input").Dot("HashedPassword"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableAvatarColumn"),
					jen.ID("input").Dot("AvatarSrc"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableTwoFactorSekretColumn"),
					jen.ID("input").Dot("TwoFactorSecret"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
					jen.ID("input").Dot("TwoFactorSecretVerifiedOn"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("input").Dot("ID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildSetUserStatusQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildSetUserStatusQuery returns a SQL query (and arguments) that would change a user's account status."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildSetUserStatusQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserReputationUpdateInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableReputationColumn"),
					jen.ID("input").Dot("NewReputation"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableStatusExplanationColumn"),
					jen.ID("input").Dot("Reason"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("input").Dot("TargetUserID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateUserPasswordQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's passwords."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildUpdateUserPasswordQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newHash").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableHashedPasswordColumn"),
					jen.ID("newHash"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableRequiresPasswordChangeColumn"),
					jen.ID("false"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTablePasswordLastChangedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateUserTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUpdateUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildUpdateUserTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newSecret").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
					jen.ID("nil"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableTwoFactorSekretColumn"),
					jen.ID("newSecret"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildVerifyUserTwoFactorSecretQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildVerifyUserTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableTwoFactorVerifiedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("UsersTableReputationColumn"),
					jen.ID("types").Dot("GoodStandingAccountStatus"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildArchiveUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildArchiveUserQuery builds a SQL query that marks a user as archived."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildArchiveUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("UsersTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesForUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAuditLogEntriesForUserQuery constructs a SQL query for fetching audit log entries belong to a user with a given ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAuditLogEntriesForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("userIDKey").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("jsonPluckQuery"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableContextColumn"),
				jen.ID("audit").Dot("UserAssignmentKey"),
			),
			jen.ID("performedByIDKey").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("jsonPluckQuery"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableContextColumn"),
				jen.ID("audit").Dot("ActorAssignmentKey"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Or").Valuesln(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("userIDKey").Op(":").ID("userID")), jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("performedByIDKey").Op(":").ID("userID")))).Dot("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("CreatedOnColumn"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}
