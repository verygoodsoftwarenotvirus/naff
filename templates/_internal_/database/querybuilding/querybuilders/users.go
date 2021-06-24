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
			jen.ID("_").Qual(proj.QuerybuildingPackage(), "UserSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID(dbvendor.Singular())).Call(jen.ID("nil")),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildUserHasStatusQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("statuses").Op("...").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Newline(),
			jen.ID("whereStatuses").Op(":=").ID("squirrel").Dot("Or").Values(),
			jen.For(jen.List(jen.ID("_"), jen.ID("status")).Op(":=").Range().ID("statuses")).Body(
				jen.ID("whereStatuses").Op("=").ID("append").Call(
					jen.ID("whereStatuses"),
					jen.ID("squirrel").Dot("Eq").Values(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "UsersTableReputationColumn"),
					).Op(":").ID("status")),
				),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				)).
					Dotln("Prefix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistencePrefix")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))).
					Dotln("Where").Call(jen.ID("whereStatuses")).
					Dotln("Suffix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistenceSuffix")),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("NotEq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetUserWithUnverifiedTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
				).Op(":").ID("nil"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetUserByUsernameQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("username"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableUsernameColumn"),
				).Op(":").ID("username"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("NotEq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildSearchForUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	searchCmd := "%s.%s LIKE ?"
	if dbvendor.SingularPackageName() == "postgres" {
		searchCmd = "%s.%s ILIKE ?"
	}

	lines := []jen.Code{
		jen.Comment("BuildSearchForUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildSearchForUserByUsernameQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("usernameQuery").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachSearchQueryToSpan").Call(
				jen.ID("span"),
				jen.ID("usernameQuery"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Expr").Callln(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit(searchCmd),
						jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
						jen.Qual(proj.QuerybuildingPackage(), "UsersTableUsernameColumn"),
					),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s%%"),
						jen.ID("usernameQuery"),
					),
				)).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).Op(":").ID("nil"))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("NotEq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetAllUsersCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("query").ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
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
		jen.Comment("to a given filter's criteria. It is assumed that this is only accessible to site administrators."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetUsersQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildListQuery").Callln(
				jen.ID("ctx"),
				jen.Qual(proj.QuerybuildingPackage(), "UsersTableName"),
				jen.ID("nil"),
				jen.ID("nil"),
				jen.Lit(""),
				jen.Qual(proj.QuerybuildingPackage(), "UsersTableColumns"),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildTestUserCreationQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("testUserConfig").Op("*").Qual(proj.TypesPackage(), "TestUserCreationConfig")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("testUserConfig").Dot("Username"),
			),
			jen.Newline(),
			jen.ID("serviceRole").Op(":=").Qual(proj.InternalAuthorizationPackage(), "ServiceUserRole"),
			jen.If(jen.ID("testUserConfig").Dot("IsServiceAdmin")).Body(
				jen.ID("serviceRole").Op("=").Qual(proj.InternalAuthorizationPackage(), "ServiceAdminRole"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableUsernameColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableHashedPasswordColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorSekretColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableReputationColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableServiceRolesColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("testUserConfig").Dot("Username"),
					jen.ID("testUserConfig").Dot("HashedPassword"),
					jen.Qual(proj.QuerybuildingPackage(), "DefaultTestUserTwoFactorSecret"),
					jen.Qual(proj.TypesPackage(), "GoodStandingAccountStatus"),
					jen.ID("serviceRole").Dot("String").Call(),
					jen.ID("currentUnixTimeQuery"),
				).Add(utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildingPackage(), "IDColumn"))))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildCreateUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(proj.TypesPackage(), "UserDataStoreCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("Username"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableUsernameColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableHashedPasswordColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorSekretColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableReputationColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableServiceRolesColumn"),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Username"),
					jen.ID("input").Dot("HashedPassword"),
					jen.ID("input").Dot("TwoFactorSecret"),
					jen.Qual(proj.TypesPackage(), "UnverifiedAccountStatus"),
					jen.Qual(proj.InternalAuthorizationPackage(), "ServiceUserRole").Dot("String").Call(),
				).Add(utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildingPackage(), "IDColumn"))))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildUpdateUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(proj.TypesPackage(), "User")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("Username"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableUsernameColumn"),
					jen.ID("input").Dot("Username"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableHashedPasswordColumn"),
					jen.ID("input").Dot("HashedPassword"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableAvatarColumn"),
					jen.ID("input").Dot("AvatarSrc"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorSekretColumn"),
					jen.ID("input").Dot("TwoFactorSecret"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
					jen.ID("input").Dot("TwoFactorSecretVerifiedOn"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("input").Dot("ID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildSetUserStatusQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(proj.TypesPackage(), "UserReputationUpdateInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("TargetUserID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableReputationColumn"),
					jen.ID("input").Dot("NewReputation"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableStatusExplanationColumn"),
					jen.ID("input").Dot("Reason"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("input").Dot("TargetUserID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildUpdateUserPasswordQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newHash").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableHashedPasswordColumn"),
					jen.ID("newHash"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableRequiresPasswordChangeColumn"),
					jen.ID("false"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTablePasswordLastChangedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("userID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildUpdateUserTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newSecret").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
					jen.ID("nil"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorSekretColumn"),
					jen.ID("newSecret"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("userID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildVerifyUserTwoFactorSecretQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableTwoFactorVerifiedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "UsersTableReputationColumn"),
					jen.Qual(proj.TypesPackage(), "GoodStandingAccountStatus"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("userID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildArchiveUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), "UsersTableName")).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").Op(":").ID("userID"), jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").Op(":").ID("nil"))),
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
		jen.Func().Params(jen.ID("b").Op("*").ID(dbvendor.Singular())).ID("BuildGetAuditLogEntriesForUserQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Newline(),
			jen.ID("userIDKey").Op(":=").Qual("fmt", "Sprintf").Callln(
				jen.ID("jsonPluckQuery"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
				utils.ConditionalCode(dbvendor.SingularPackageName() == "mariadb", jen.ID("userID")),
				jen.Qual(proj.InternalAuditPackage(), "UserAssignmentKey"),
			),
			jen.Newline(),
			jen.ID("performedByIDKey").Op(":=").Qual("fmt", "Sprintf").Callln(
				jen.ID("jsonPluckQuery"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
				utils.ConditionalCode(dbvendor.SingularPackageName() == "mariadb", jen.ID("userID")),
				jen.Qual(proj.InternalAuditPackage(), "ActorAssignmentKey"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Or").Valuesln(
					func() jen.Code {
						if dbvendor.SingularPackageName() == "mariadb" {
							return jen.ID("squirrel").Dot("Expr").Call(jen.ID("userIDKey"))
						}
						return jen.ID("squirrel").Dot("Eq").Values(jen.ID("userIDKey").Op(":").ID("userID"))
					}(),
					func() jen.Code {
						if dbvendor.SingularPackageName() == "mariadb" {
							return jen.ID("squirrel").Dot("Expr").Call(jen.ID("performedByIDKey"))
						}
						return jen.ID("squirrel").Dot("Eq").Values(jen.ID("performedByIDKey").Op(":").ID("userID"))
					}(),
				)).
					Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}
