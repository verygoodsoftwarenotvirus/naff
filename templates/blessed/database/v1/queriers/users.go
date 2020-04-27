package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	spn := dbvendor.SingularPackageName()

	ret := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("usersTableName").Equals().Lit("users"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("usersTableColumns").Equals().Index().String().Valuesln(
				utils.FormatString("%s.id", jen.ID("usersTableName")),
				utils.FormatString("%s.username", jen.ID("usersTableName")),
				utils.FormatString("%s.hashed_password", jen.ID("usersTableName")),
				utils.FormatString("%s.password_last_changed_on", jen.ID("usersTableName")),
				utils.FormatString("%s.two_factor_secret", jen.ID("usersTableName")),
				utils.FormatString("%s.is_admin", jen.ID("usersTableName")),
				utils.FormatString("%s.created_on", jen.ID("usersTableName")),
				utils.FormatString("%s.updated_on", jen.ID("usersTableName")),
				utils.FormatString("%s.archived_on", jen.ID("usersTableName")),
			),
		),
		jen.Line(),
	)

	ret.Add(buildScanUser(proj, dbvendor)...)
	ret.Add(buildScanUsers(proj, dbvendor)...)
	ret.Add(buildBuildGetUserQuery(proj, dbvendor)...)
	ret.Add(buildGetUser(proj, dbvendor)...)
	ret.Add(buildBuildGetUserByUsernameQuery(proj, dbvendor)...)
	ret.Add(buildGetUserByUsername(proj, dbvendor)...)
	ret.Add(buildBuildGetAllUsersCountQuery(proj, dbvendor)...)
	ret.Add(buildGetAllUserCount(proj, dbvendor)...)
	ret.Add(buildBuildGetUsersQuery(proj, dbvendor)...)
	ret.Add(buildGetUsers(proj, dbvendor)...)
	ret.Add(buildBuildCreateUserQuery(proj, dbvendor)...)
	ret.Add(buildCreateUser(proj, dbvendor)...)
	ret.Add(buildBuildUpdateUserQuery(proj, dbvendor)...)
	ret.Add(buildUpdateUser(proj, dbvendor)...)
	ret.Add(buildBuildArchiveUserQuery(proj, dbvendor)...)
	ret.Add(buildArchiveUser(proj, dbvendor)...)

	return ret
}

func buildScanUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("scanUser provides a consistent way to scan something like a *sql.Row into a User struct."),
		jen.Line(),
		jen.Func().ID("scanUser").Params(
			jen.ID("scan").Qual(proj.DatabaseV1Package(), "Scanner"),
			jen.ID("includeCount").Bool(),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Uint64(),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("x").Equals().AddressOf().Qual(proj.ModelsV1Package(), "User").Values(),
				jen.ID("count").Uint64(),
			),
			jen.Line(),
			jen.ID("targetVars").Assign().Index().Interface().Valuesln(
				jen.AddressOf().ID("x").Dot("ID"),
				jen.AddressOf().ID("x").Dot("Username"),
				jen.AddressOf().ID("x").Dot("HashedPassword"),
				jen.AddressOf().ID("x").Dot("PasswordLastChangedOn"),
				jen.AddressOf().ID("x").Dot("TwoFactorSecret"),
				jen.AddressOf().ID("x").Dot("IsAdmin"),
				jen.AddressOf().ID("x").Dot("CreatedOn"),
				jen.AddressOf().ID("x").Dot("UpdatedOn"),
				jen.AddressOf().ID("x").Dot("ArchivedOn"),
			),
			jen.Line(),
			jen.If(jen.ID("includeCount")).Block(
				utils.AppendItemsToList(jen.ID("targetVars"), jen.AddressOf().ID("count")),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Spread()), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("count"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildScanUsers(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("scanUsers takes database rows and loads them into a slice of User structs."),
		jen.Line(),
		jen.Func().ID("scanUsers").Params(
			jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
			jen.ID("rows").PointerTo().Qual("database/sql", "Rows"),
		).Params(
			jen.Index().Qual(proj.ModelsV1Package(), "User"),
			jen.Uint64(),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("list").Index().Qual(proj.ModelsV1Package(), "User"),
				jen.ID("count").Uint64(),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("user"), jen.ID("c"), jen.Err()).Assign().ID("scanUser").Call(jen.ID("rows"), jen.True()),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(
						jen.Nil(),
						jen.Zero(),
						jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning user result: %w"), jen.Err()),
					),
				),
				jen.Line(),
				jen.If(jen.ID("count").IsEqualTo().Zero()).Block(
					jen.ID("count").Equals().ID("c"),
				),
				jen.Line(),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.PointerTo().ID("user")),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("count"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetUserQuery").Params(jen.ID("userID").Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Spread()).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				utils.FormatString("%s.id", jen.ID("usersTableName")).MapAssign().ID("userID"))).
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
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetUser").Params(constants.CtxParam(), jen.ID("userID").Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserQuery").Call(jen.ID("userID")),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.List(jen.ID("u"), jen.Underscore(), jen.Err()).Assign().ID("scanUser").Call(jen.ID("row"), jen.False()),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Lit("fetching user from database"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUserByUsernameQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetUserByUsernameQuery").Params(jen.ID("username").String()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Spread()).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				utils.FormatString("%s.username", jen.ID("usersTableName")).MapAssign().ID("username"))).
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
			jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error(),
		).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUserByUsernameQuery").Call(jen.ID("username")),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.List(jen.ID("u"), jen.Underscore(), jen.Err()).Assign().ID("scanUser").Call(jen.ID("row"), jen.False()),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
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

func buildBuildGetAllUsersCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("buildGetAllUserCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildGetAllUserCountQuery").Params().Params(
			jen.ID("query").String(),
		).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(utils.FormatStringWithArg(jen.ID("countQuery"), jen.ID("usersTableName"))).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.archived_on", jen.ID("usersTableName")).MapAssign().ID("nil"),
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

func buildGetAllUserCount(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("GetAllUserCount fetches a count of users from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetAllUserCount").Params(
			constants.CtxParam(),
		).Params(
			jen.ID("count").Uint64(),
			jen.Err().Error(),
		).Block(
			jen.ID("query").Assign().ID(dbfl).Dot("buildGetAllUserCountQuery").Call(),
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
		).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.Append(jen.ID("usersTableColumns"), utils.FormatStringWithArg(jen.ID("countQuery"), jen.ID("usersTableName"))).Spread()).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
					utils.FormatString("%s.archived_on", jen.ID("usersTableName")).MapAssign().Nil(),
				),
			).
				Dotln("GroupBy").Call(utils.FormatString("%s.id", jen.ID("usersTableName"))),
			jen.Line(),
			jen.If(jen.ID(constants.FilterVarName).DoesNotEqual().Nil()).Block(
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
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("GetUsers").Params(constants.CtxParam(), jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList"), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildGetUsersQuery").Call(jen.ID(constants.FilterVarName)),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Lit("querying for user"))),
			),
			jen.Line(),
			jen.List(jen.ID("userList"), jen.ID("count"), jen.Err()).Assign().ID("scanUsers").Call(jen.ID(dbfl).Dot(constants.LoggerVarName), jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserList").Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID(constants.FilterVarName).Dot("Page"),
					jen.ID("Limit").MapAssign().ID(constants.FilterVarName).Dot("Limit"),
					jen.ID("TotalCount").MapAssign().ID("count"),
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
		jen.Lit("username"),
		jen.Lit("hashed_password"),
		jen.Lit("two_factor_secret"),
		jen.Lit("is_admin"),
	}
	vals := []jen.Code{
		jen.ID("input").Dot("Username"),
		jen.ID("input").Dot("HashedPassword"),
		jen.ID("input").Dot("TwoFactorSecret"),
		jen.False(),
	}

	buildCreateUserQuery := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Insert").Call(jen.ID("usersTableName")).
		Dotln("Columns").Callln(cols...).
		Dotln("Values").Callln(vals...)

	if isPostgres(dbvendor) {
		buildCreateUserQuery.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
	}

	buildCreateUserQuery.Dotln("ToSql").Call()

	lines := []jen.Code{
		jen.Comment("buildCreateUserQuery returns a SQL query (and arguments) that would create a given User"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildCreateUserQuery").Params(
			jen.ID("input").Qual(proj.ModelsV1Package(), "UserDatabaseCreationInput"),
		).Params(
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			buildCreateUserQuery,
			jen.Line(),
			jen.Comment("NOTE: we always default is_admin to false, on the assumption that"),
			jen.Comment("admins have DB access and will change that value via SQL query."),
			jen.Comment("There should also be no way to update a user via this structure."),
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
			jen.If(jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("ID"), jen.AddressOf().ID("x").Dot("CreatedOn")).Op(";").Err().DoesNotEqual().ID("nil")).Block(
				jen.Switch(jen.ID("e").Assign().Err().Assert(jen.Type())).Block(
					jen.Case(jen.PointerTo().Qual("github.com/lib/pq", "Error")).Block(
						jen.If(jen.ID("e").Dot("Code").IsEqualTo().Qual("github.com/lib/pq", "ErrorCode").Call(jen.ID("postgresRowExistsErrorCode"))).Block(
							jen.Return().List(jen.Nil(), jen.Qual(proj.DatabaseV1Package("client"), "ErrUserExists")),
						),
					),
					jen.Default().Block(
						jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing user creation query: %w"), jen.Err())),
					),
				),
			),
		}
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		uqb = []jen.Code{
			jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
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
		jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), "User").Valuesln(
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
			jen.ID("input").Qual(proj.ModelsV1Package(), "UserDatabaseCreationInput"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error(),
		).Block(
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
			Dotln("Set").Call(jen.Lit("username"), jen.ID("input").Dot("Username")).
			Dotln("Set").Call(jen.Lit("hashed_password"), jen.ID("input").Dot("HashedPassword")).
			Dotln("Set").Call(jen.Lit("two_factor_secret"), jen.ID("input").Dot("TwoFactorSecret")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").MapAssign().ID("input").Dot("ID")))

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING updated_on"))
		}

		q.Dotln("ToSql").Call()
		return q
	}

	lines := []jen.Code{
		jen.Comment("buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildUpdateUserQuery").Params(jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
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

	buildUpdateUserBody := func() []jen.Code {
		if isPostgres(dbvendor) {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("input")),
				jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("input").Dot("UpdatedOn")),
			}
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("input")),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
				jen.Return().Err(),
			}
		}

		return nil
	}

	lines := []jen.Code{
		jen.Comment("UpdateUser receives a complete User struct and updates its place in the db."),
		jen.Line(),
		jen.Comment("NOTE this function uses the ID provided in the input to make its query. Pass in"),
		jen.Line(),
		jen.Comment("anonymous structs or incomplete models at your peril."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("UpdateUser").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.Error()).Block(
			buildUpdateUserBody()...,
		),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveUserQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	buildArchiveUserQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("usersTableName")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").MapAssign().ID("userID")))

		if isPostgres(dbvendor) {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING archived_on"))
		}

		q.Dotln("ToSql").Call()
		return q
	}

	lines := []jen.Code{
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("buildArchiveUserQuery").Params(jen.ID("userID").Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
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

func buildArchiveUser(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := dbvendor.Singular()
	dbfl := string(dbvendor.LowercaseAbbreviation()[0])

	lines := []jen.Code{
		jen.Comment("ArchiveUser archives a user by their username."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(sn)).ID("ArchiveUser").Params(constants.CtxParam(), jen.ID("userID").Uint64()).Params(jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dot("buildArchiveUserQuery").Call(jen.ID("userID")),
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return().Err(),
		),
		jen.Line(),
	}

	return lines
}
