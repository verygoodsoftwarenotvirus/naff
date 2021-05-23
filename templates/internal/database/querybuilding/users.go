package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabasePackage("queriers", "v1", spn), spn)

	utils.AddImports(proj, code, false)

	code.Add(buildUsersFileConstDeclarations()...)
	code.Add(buildUsersFileVarDeclarations()...)
	code.Add(buildScanUser(proj, dbvendor)...)
	code.Add(buildScanUsers(proj, dbvendor)...)
	code.Add(buildBuildGetUserQuery(dbvendor)...)
	code.Add(buildGetUser(proj, dbvendor)...)
	code.Add(buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(dbvendor)...)
	code.Add(buildGetUserWithUnverifiedTwoFactorSecret(proj, dbvendor)...)
	code.Add(buildBuildGetUserByUsernameQuery(dbvendor)...)
	code.Add(buildGetUserByUsername(proj, dbvendor)...)
	code.Add(buildBuildGetAllUsersCountQuery(dbvendor)...)
	code.Add(buildGetAllUsersCount(dbvendor)...)
	code.Add(buildBuildGetUsersQuery(proj, dbvendor)...)
	code.Add(buildGetUsers(proj, dbvendor)...)
	code.Add(buildBuildCreateUserQuery(proj, dbvendor)...)
	code.Add(buildCreateUser(proj, dbvendor)...)
	code.Add(buildBuildUpdateUserQuery(proj, dbvendor)...)
	code.Add(buildUpdateUser(proj, dbvendor)...)
	code.Add(buildBuildUpdateUserPasswordQuery(dbvendor)...)
	code.Add(buildUpdateUserPassword(dbvendor)...)
	code.Add(buildBuildVerifyUserTwoFactorSecretQuery(dbvendor)...)
	code.Add(buildVerifyUserTwoFactorSecret(dbvendor)...)
	code.Add(buildBuildArchiveUserQuery(dbvendor)...)
	code.Add(buildArchiveUser(dbvendor)...)

	return code
}

func buildUsersFileConstDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("usersTableName").Equals().Lit("users"),
			jen.ID("usersTableUsernameColumn").Equals().Lit("username"),
			jen.ID("usersTableHashedPasswordColumn").Equals().Lit("hashed_password"),
			jen.ID("usersTableSaltColumn").Equals().Lit("salt"),
			jen.ID("usersTableRequiresPasswordChangeColumn").Equals().Lit("requires_password_change"),
			jen.ID("usersTablePasswordLastChangedOnColumn").Equals().Lit("password_last_changed_on"),
			jen.ID("usersTableTwoFactorColumn").Equals().Lit("two_factor_secret"),
			jen.ID("usersTableTwoFactorVerifiedOnColumn").Equals().Lit("two_factor_secret_verified_on"),
			jen.ID("usersTableIsAdminColumn").Equals().Lit("is_admin"),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersFileVarDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("usersTableColumns").Equals().Index().String().Valuesln(
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("idColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableUsernameColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableHashedPasswordColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableSaltColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableRequiresPasswordChangeColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTablePasswordLastChangedOnColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableTwoFactorColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableTwoFactorVerifiedOnColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableIsAdminColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("createdOnColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("lastUpdatedOnColumn")),
				utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("archivedOnColumn")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildScanUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("scanUser provides a consistent way to scan something like a *sql.Row into a User struct."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("scanUser").Params(
			jen.ID("scan").Qual(proj.DatabasePackage(), "Scanner"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "User"),
			jen.Error(),
		).Body(
			jen.Var().Defs(
				jen.ID("x").Equals().AddressOf().Qual(proj.TypesPackage(), "User").Values(),
			),
			jen.Line(),
			jen.ID("targetVars").Assign().Index().Interface().Valuesln(
				jen.AddressOf().ID("x").Dot("ID"),
				jen.AddressOf().ID("x").Dot("Username"),
				jen.AddressOf().ID("x").Dot("HashedPassword"),
				jen.AddressOf().ID("x").Dot("Salt"),
				jen.AddressOf().ID("x").Dot("RequiresPasswordChange"),
				jen.AddressOf().ID("x").Dot("PasswordLastChangedOn"),
				jen.AddressOf().ID("x").Dot("TwoFactorSecret"),
				jen.AddressOf().ID("x").Dot("TwoFactorSecretVerifiedOn"),
				jen.AddressOf().ID("x").Dot("IsAdmin"),
				jen.AddressOf().ID("x").Dot("CreatedOn"),
				jen.AddressOf().ID("x").Dot("LastUpdatedOn"),
				jen.AddressOf().ID("x").Dot("ArchivedOn"),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Spread()), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildScanUsers(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("scanUsers takes database rows and loads them into a slice of User structs."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("scanUsers").Params(
			jen.ID("rows").Qual(proj.DatabasePackage(), "ResultIterator"),
		).Params(
			jen.Index().Qual(proj.TypesPackage(), "User"),
			jen.Error(),
		).Body(
			jen.Var().Defs(
				jen.ID("list").Index().Qual(proj.TypesPackage(), "User"),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("user"), jen.Err()).Assign().ID(dbfl).Dot("scanUser").Call(jen.ID("rows")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(
						jen.Nil(),
						jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning user result: %w"), jen.Err()),
					),
				),
				jen.Line(),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.PointerTo().ID("user")),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUserQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetUserQuery").Params(constants.UserIDParam()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Spread()).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("idColumn")).MapAssign().ID(constants.UserIDVarName),
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("archivedOnColumn")).MapAssign().Nil(),
				)).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "NotEq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableTwoFactorVerifiedOnColumn")).MapAssign().Nil(),
				)).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildGetUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetUser fetches a user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "User"), jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserQuery").Call(jen.ID(constants.UserIDVarName)),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().ID(dbfl).Dot("scanUser").Call(jen.ID("row")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Lit("fetching user from database"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user"),
		jen.Line(),
		jen.Comment("by their database ID, who has an unverified two factor secret"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetUserWithUnverifiedTwoFactorSecretQuery").Params(constants.UserIDParam()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Spread()).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("idColumn")).MapAssign().ID(constants.UserIDVarName),
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableTwoFactorVerifiedOnColumn")).MapAssign().Nil(),
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("archivedOnColumn")).MapAssign().Nil(),
				)).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildGetUserWithUnverifiedTwoFactorSecret(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified two factor secret"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetUserWithUnverifiedTwoFactorSecret").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "User"), jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserWithUnverifiedTwoFactorSecretQuery").Call(jen.ID(constants.UserIDVarName)),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().ID(dbfl).Dot("scanUser").Call(jen.ID("row")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Lit("fetching user from database"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUserByUsernameQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetUserByUsernameQuery").Params(jen.ID("username").String()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Spread()).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableUsernameColumn")).MapAssign().ID("username"),
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("archivedOnColumn")).MapAssign().Nil(),
				)).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "NotEq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("usersTableTwoFactorVerifiedOnColumn")).MapAssign().Nil(),
				)).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildGetUserByUsername(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetUserByUsername fetches a user by their username."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetUserByUsername").Params(
			constants.CtxParam(),
			jen.ID("username").String(),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "User"),
			jen.Error(),
		).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserByUsernameQuery").Call(jen.ID("username")),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.List(jen.ID("u"), jen.Err()).Assign().ID(dbfl).Dot("scanUser").Call(jen.ID("row")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Body(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetAllUsersCountQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllUsersCountQuery").Params().Params(
			jen.ID("query").String(),
		).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(utils.FormatStringWithArg(jen.ID("countQuery"), jen.ID("usersTableName"))).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("archivedOnColumn")).MapAssign().ID("nil"),
				),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.Underscore(), jen.Err()).Equals().ID("builder").Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query")),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllUsersCount(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetAllUsersCount fetches a count of users from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllUsersCount").Params(
			constants.CtxParam(),
		).Params(
			jen.ID("count").Uint64(),
			jen.Err().Error(),
		).Body(
			jen.ID("query").Assign().ID(dbfl).Dot("buildGetAllUsersCountQuery").Call(),
			jen.Err().Equals().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(
				constants.CtxVar(),
				jen.ID("query"),
			).Dot("Scan").Call(jen.AddressOf().ID("count")),
			jen.Return(),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUsersQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetUsersQuery").Params(
			utils.QueryFilterParam(proj),
		).Params(
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Spread()).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("archivedOnColumn")).MapAssign().Nil(),
				),
			).
				Dotln("OrderBy").Call(utils.FormatString("%s.%s", jen.ID("usersTableName"), jen.ID("idColumn"))),
			jen.Line(),
			jen.If(jen.ID(constants.FilterVarName).DoesNotEqual().Nil()).Body(
				jen.ID("builder").Equals().ID(constants.FilterVarName).Dot("ApplyToQueryBuilder").Call(jen.ID("builder"), jen.ID("usersTableName")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildGetUsers(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetUsers fetches a list of users from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetUsers").Params(constants.CtxParam(), jen.ID(constants.FilterVarName).PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "UserList"), jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUsersQuery").Call(jen.ID(constants.FilterVarName)),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Lit("querying for user"))),
			),
			jen.Line(),
			jen.List(jen.ID("userList"), jen.Err()).Assign().ID(dbfl).Dot("scanUsers").Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Assign().AddressOf().Qual(proj.TypesPackage(), "UserList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.TypesPackage(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID(constants.FilterVarName).Dot("Page"),
					jen.ID("Limit").MapAssign().ID(constants.FilterVarName).Dot("Limit"),
				),
				jen.ID("Users").MapAssign().ID("userList"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	cols := []jen.Code{
		jen.ID("usersTableUsernameColumn"),
		jen.ID("usersTableHashedPasswordColumn"),
		jen.ID("usersTableSaltColumn"),
		jen.ID("usersTableTwoFactorColumn"),
		jen.ID("usersTableIsAdminColumn"),
	}
	vals := []jen.Code{
		jen.ID("input").Dot("Username"),
		jen.ID("input").Dot("HashedPassword"),
		jen.ID("input").Dot("Salt"),
		jen.ID("input").Dot("TwoFactorSecret"),
		jen.False(),
	}

	buildCreateUserQuery := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Insert").Call(jen.ID("usersTableName")).
		Dotln("Columns").Callln(cols...).
		Dotln("Values").Callln(vals...)

	if isPostgres(dbvendor) {
		buildCreateUserQuery.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s, %s"), jen.ID("idColumn"), jen.ID("createdOnColumn")))
	}

	buildCreateUserQuery.Dotln("ToSql").Call()

	lines := []jen.Code{
		jen.Comment("buildCreateUserQuery returns a SQL query (and arguments) that would create a given User"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildCreateUserQuery").Params(
			jen.ID("input").Qual(proj.TypesPackage(), "UserDatabaseCreationInput"),
		).Params(
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			buildCreateUserQuery,
			jen.Line(),
			jen.Comment("NOTE: we always default is_admin to false, on the assumption that"),
			jen.Comment("admins have DB access and will change that value via SQL query."),
			jen.Comment("There should also be no way to update a user via this structure"),
			jen.Comment("such that they would have admin privileges."),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	uqb := []jen.Code{}
	if isPostgres(dbvendor) {
		uqb = []jen.Code{
			jen.If(jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("ID"), jen.AddressOf().ID("x").Dot("CreatedOn")).Op(";").Err().DoesNotEqual().ID("nil")).Body(
				jen.Switch(jen.ID("e").Assign().Err().Assert(jen.Type())).Body(
					jen.Case(jen.PointerTo().Qual("github.com/lib/pq", "Error")).Body(
						jen.If(jen.ID("e").Dot("Code").IsEqualTo().Qual("github.com/lib/pq", "ErrorCode").Call(jen.ID("postgresRowExistsErrorCode"))).Body(
							jen.Return().List(jen.Nil(), jen.Qual(proj.DatabasePackage("client"), "ErrUserExists")),
						),
					),
					jen.Default().Body(
						jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing user creation query: %w"), jen.Err())),
					),
				),
			),
		}
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		uqb = []jen.Code{
			jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing user creation query: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("fetch the last inserted ID."),
			jen.List(jen.ID("id"), jen.ID("err")).Assign().ID(constants.ResponseVarName).Dot("LastInsertId").Call(),
			jen.ID(dbfl).Dot("logIDRetrievalError").Call(jen.Err()),
			jen.ID("x").Dot("ID").Equals().Uint64().Call(jen.ID("id")),
			jen.Line(),
			jen.Comment("this won't be completely accurate, but it will suffice."),
			jen.ID("x").Dot("CreatedOn").Equals().ID(dbfl).Dot("timeTeller").Dot("Now").Call(),
		}
	}

	createUserBlock := []jen.Code{
		jen.ID("x").Assign().AddressOf().Qual(proj.TypesPackage(), "User").Valuesln(
			jen.ID("Username").MapAssign().ID("input").Dot("Username"),
			jen.ID("HashedPassword").MapAssign().ID("input").Dot("HashedPassword"),
			jen.ID("TwoFactorSecret").MapAssign().ID("input").Dot("TwoFactorSecret")),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildCreateUserQuery").Call(jen.ID("input")),
		jen.Line(),
		jen.Comment("create the user."),
	}
	createUserBlock = append(createUserBlock, uqb...)
	createUserBlock = append(createUserBlock, jen.Line(), jen.Return().List(jen.ID("x"), jen.Nil()))

	lines := []jen.Code{
		jen.Comment("CreateUser creates a user."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("CreateUser").Params(
			constants.CtxParam(),
			jen.ID("input").Qual(proj.TypesPackage(), "UserDatabaseCreationInput"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "User"),
			jen.Error(),
		).Body(
			createUserBlock...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildUpdateUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildUpdateUserQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("usersTableName")).
			Dotln("Set").Call(jen.ID("usersTableUsernameColumn"), jen.ID("input").Dot("Username")).
			Dotln("Set").Call(jen.ID("usersTableHashedPasswordColumn"), jen.ID("input").Dot("HashedPassword")).
			Dotln("Set").Call(jen.ID("usersTableSaltColumn"), jen.ID("input").Dot("Salt")).
			Dotln("Set").Call(jen.ID("usersTableTwoFactorColumn"), jen.ID("input").Dot("TwoFactorSecret")).
			Dotln("Set").Call(jen.ID("usersTableTwoFactorVerifiedOnColumn"), jen.ID("input").Dot("TwoFactorSecretVerifiedOn")).
			Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.ID("idColumn").MapAssign().ID("input").Dot("ID")))

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.ID("lastUpdatedOnColumn")))
		}

		q.Dotln("ToSql").Call()
		return q
	}

	lines := []jen.Code{
		jen.Comment("buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildUpdateUserQuery").Params(jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "User")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			buildUpdateUserQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("UpdateUser receives a complete User struct and updates its place in the db."),
		jen.Line(),
		jen.Comment("NOTE this function uses the ID provided in the input to make its query. Pass in"),
		jen.Line(),
		jen.Comment("incomplete models at your peril."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("UpdateUser").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "User")).Params(jen.Error()).Body(
			func() []jen.Code {
				if isPostgres(dbvendor) {
					return []jen.Code{
						jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("input")),
						jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("input").Dot("LastUpdatedOn")),
					}
				} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
					return []jen.Code{
						jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("input")),
						jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
						jen.Return().Err(),
					}
				}
				panic("invalid database vendor")
			}()...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildUpdateUserPasswordQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's password."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildUpdateUserPasswordQuery").Params(
			jen.ID("userID").Uint64(),
			jen.ID("newHash").String(),
		).Params(
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			func() jen.Code {
				q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Update").Call(jen.ID("usersTableName")).
					Dotln("Set").Call(jen.ID("usersTableHashedPasswordColumn"), jen.ID("newHash")).
					Dotln("Set").Call(jen.ID("usersTableRequiresPasswordChangeColumn"), jen.False()).
					Dotln("Set").Call(jen.ID("usersTablePasswordLastChangedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
					Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.ID("idColumn").MapAssign().ID("userID")))

				if isPostgres(dbvendor) {
					q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.ID("lastUpdatedOnColumn")))
				}

				q.Dotln("ToSql").Call()
				return q
			}(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateUserPassword(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("UpdateUserPassword updates a user's password."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("UpdateUserPassword").Params(
			constants.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID("newHash").String(),
		).Params(jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateUserPasswordQuery").Call(jen.ID("userID"), jen.ID("newHash")),
			jen.Line(),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Line(),
			jen.Return().Err(),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildVerifyUserTwoFactorSecretQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildVerifyUserTwoFactorSecretQuery").Params(constants.UserIDParam()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Update").Call(jen.ID("usersTableName")).
				Dotln("Set").Call(jen.ID("usersTableTwoFactorVerifiedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					jen.ID("idColumn").MapAssign().ID("userID"),
				)).Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return(jen.ID("query"), jen.ID("args")),
		),
	}

	return lines
}

func buildVerifyUserTwoFactorSecret(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("VerifyUserTwoFactorSecret marks a user's two factor secret as validated."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("VerifyUserTwoFactorSecret").Params(
			constants.CtxParam(),
			jen.ID("userID").Uint64(),
		).Params(
			jen.Error(),
		).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildVerifyUserTwoFactorSecretQuery").Call(jen.ID("userID")),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(
				constants.CtxVar(),
				jen.ID("query"),
				jen.ID("args").Spread(),
			),
			jen.Return().Err(),
		),
	}

	return lines
}

func buildBuildArchiveUserQuery(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildArchiveUserQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("usersTableName")).
			Dotln("Set").Call(jen.ID("archivedOnColumn"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.ID("idColumn").MapAssign().ID(constants.UserIDVarName)))

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.ID("archivedOnColumn")))
		}

		q.Dotln("ToSql").Call()
		return q
	}

	lines := []jen.Code{
		jen.Comment("buildArchiveUserQuery builds a SQL query that marks a user as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildArchiveUserQuery").Params(constants.UserIDParam()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.Var().Err().Error(),
			jen.Line(),
			buildArchiveUserQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveUser(dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("ArchiveUser marks a user as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("ArchiveUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Error()).Body(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveUserQuery").Call(jen.ID(constants.UserIDVarName)),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return().Err(),
		),
		jen.Line(),
	}

	return lines
}
