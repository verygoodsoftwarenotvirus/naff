package queriers

import (
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(pkg *models.Project, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(pkg, ret)
	sn := vendor.Singular()
	dbrn := vendor.RouteName()
	dbfl := strings.ToLower(string([]byte(sn)[0]))

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Const().Defs(
			jen.ID("usersTableName").Op("=").Lit("users"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("usersTableColumns").Op("=").Index().ID("string").Valuesln(
				jen.Lit("id"),
				jen.Lit("username"),
				jen.Lit("hashed_password"),
				jen.Lit("password_last_changed_on"),
				jen.Lit("two_factor_secret"),
				jen.Lit("is_admin"),
				jen.Lit("created_on"),
				jen.Lit("updated_on"),
				jen.Lit("archived_on"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanUser provides a consistent way to scan something like a *sql.Row into a User struct"),
		jen.Line(),
		jen.Func().ID("scanUser").Params(jen.ID("scan").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Scanner")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User"), jen.ID("error")).Block(
			jen.Var().ID("x").Op("=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Values(),
			jen.Line(),
			jen.If(jen.Err().Op(":=").ID("scan").Dot("Scan").Callln(
				jen.Op("&").ID("x").Dot("ID"),
				jen.Op("&").ID("x").Dot("Username"),
				jen.Op("&").ID("x").Dot("HashedPassword"),
				jen.Op("&").ID("x").Dot("PasswordLastChangedOn"),
				jen.Op("&").ID("x").Dot("TwoFactorSecret"),
				jen.Op("&").ID("x").Dot("IsAdmin"),
				jen.Op("&").ID("x").Dot("CreatedOn"),
				jen.Op("&").ID("x").Dot("UpdatedOn"),
				jen.Op("&").ID("x").Dot("ArchivedOn"),
			), jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanUsers takes database rows and loads them into a slice of User structs"),
		jen.Line(),
		jen.Func().ID("scanUsers").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User"), jen.ID("error")).Block(
			jen.Var().ID("list").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User"),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("user"), jen.Err()).Op(":=").ID("scanUser").Call(jen.ID("rows")),
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning user result: %w"), jen.Err())),
				),
				jen.ID("list").Op("=").ID("append").Call(jen.ID("list"), jen.Op("*").ID("user")),
			),
			jen.Line(),
			jen.If(jen.Err().Op(":=").ID("rows").Dot("Err").Call(), jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.Err().Op(":=").ID("rows").Dot("Close").Call(), jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetUserQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("userID"))).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetUser fetches a user"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUserQuery").Call(jen.ID("userID")),
			jen.ID("row").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.List(jen.ID("u"), jen.Err()).Op(":=").ID("scanUser").Call(jen.ID("row")),
			jen.Line(),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.Err(), jen.Lit("fetching user from database"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.Err()),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetUserByUsernameQuery").Params(jen.ID("username").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("username").Op(":").ID("username"))).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetUserByUsername fetches a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetUserByUsername").Params(utils.CtxParam(), jen.ID("username").ID("string")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUserByUsernameQuery").Call(jen.ID("username")),
			jen.ID("row").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.List(jen.ID("u"), jen.Err()).Op(":=").ID("scanUser").Call(jen.ID("row")),
			jen.Line(),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.Err()),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetUserCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetUserCountQuery").Params(jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("CountQuery")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.Line(),
			jen.If(jen.ID(utils.FilterVarName).Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetUserCount fetches a count of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetUserCount").Params(utils.CtxParam(), jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUserCountQuery").Call(jen.ID(utils.FilterVarName)),
			jen.Err().Op("=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return(),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("buildGetUserCountQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildGetUsersQuery").Params(jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.Line(),
			jen.If(jen.ID(utils.FilterVarName).Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("GetUsers fetches a list of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("GetUsers").Params(utils.CtxParam(), jen.ID(utils.FilterVarName).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserList"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildGetUsersQuery").Call(jen.ID(utils.FilterVarName)),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.Err(), jen.Lit("querying for user"))),
			),
			jen.Line(),
			jen.List(jen.ID("userList"), jen.Err()).Op(":=").ID("scanUsers").Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("count"), jen.Err()).Op(":=").ID(dbfl).Dot("GetUserCount").Call(utils.CtxVar(), jen.ID(utils.FilterVarName)),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user count: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserList").Valuesln(
				jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
					jen.ID("Page").Op(":").ID("filter").Dot("Page"),
					jen.ID("Limit").Op(":").ID("filter").Dot("Limit"),
					jen.ID("TotalCount").Op(":").ID("count"),
				),
				jen.ID("Users").Op(":").ID("userList"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	buildBuildCreateUserQueryQuery := func() jen.Code {
		cols := []jen.Code{
			jen.Lit("username"),
			jen.Lit("hashed_password"),
			jen.Lit("two_factor_secret"),
			jen.Lit("is_admin"),
		}
		vals := []jen.Code{
			jen.ID("input").Dot("Username"),
			jen.ID("input").Dot("Password"),
			jen.ID("input").Dot("TwoFactorSecret"),
			jen.ID("false"),
		}

		if isMariaDB {
			cols = append(cols, jen.Lit("created_on"))
			vals = append(vals, jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery")))
		}

		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Insert").Call(jen.ID("usersTableName")).
			Dotln("Columns").Callln(cols...).
			Dotln("Values").Callln(vals...)

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
		}

		q.Dotln("ToSql").Call()

		return q
	}

	ret.Add(
		jen.Comment("buildCreateUserQuery returns a SQL query (and arguments) that would create a given User"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildCreateUserQuery").Params(jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildBuildCreateUserQueryQuery(),
			jen.Line(),
			jen.Comment("NOTE: we always default is_admin to false, on the assumption that"),
			jen.Comment("admins have DB access and will change that value via SQL query."),
			jen.Comment("There should also be no way to update a user via this structure"),
			jen.Comment("such that they would have admin privileges"),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	if isSqlite || isMariaDB {
		ret.Add(
			jen.Comment("buildUserCreationTimeQuery returns a SQL query (and arguments) that would create a given User"),
			jen.Line(),
			jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildUserCreationTimeQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID(dbfl).Dot("sqlBuilder").Dot("Select").Call(jen.Lit("created_on")).
					Dotln("From").Call(jen.ID("usersTableName")).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("userID"))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
				jen.Line(),
				jen.Return(jen.List(jen.ID("query"), jen.ID("args"))),
			),
			jen.Line(),
		)
	}

	////////////

	buildCreateUserQueryBlock := func() []jen.Code {
		if isPostgres {
			out := []jen.Code{
				jen.If(jen.Err().Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("CreatedOn")).Op(";").ID("err").Op("!=").ID("nil")).Block(
					jen.Switch(jen.ID("e").Op(":=").ID("err").Assert(jen.Type())).Block(
						jen.Case(jen.Op("*").Qual("github.com/lib/pq", "Error")).Block(
							jen.If(jen.ID("e").Dot("Code").Op("==").Qual("github.com/lib/pq", "ErrorCode").Call(jen.Lit("23505"))).Block(
								jen.Return().List(jen.ID("nil"), jen.Qual(filepath.Join(pkg.OutputPath, "database/v1/client"), "ErrUserExists")),
							),
						),
						jen.Default().Block(
							jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing user creation query: %w"), jen.Err())),
						),
					),
				),
			}
			return out
		} else if isSqlite || isMariaDB {
			out := []jen.Code{
				jen.List(jen.ID("res"), jen.Err()).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing user creation query: %w"), jen.Err())),
				),
				jen.Line(),
				jen.Comment("fetch the last inserted ID"),
				jen.If(jen.List(jen.ID("id"), jen.ID("idErr")).Op(":=").ID("res").Dot("LastInsertId").Call().Op(";").ID("idErr").Op("==").ID("nil")).Block(
					jen.ID("x").Dot("ID").Op("=").ID("uint64").Call(jen.ID("id")),
					jen.Line(),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUserCreationTimeQuery").Call(jen.ID("x").Dot("ID")),
					jen.ID(dbfl).Dot("logCreationTimeRetrievalError").Call(jen.ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("CreatedOn"))),
				),
			}
			return out
		}

		return nil
	}

	buildCreateUserBlock := func() []jen.Code {
		uqb := buildCreateUserQueryBlock()

		out := []jen.Code{
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User").Valuesln(
				jen.ID("Username").Op(":").ID("input").Dot("Username"),
				jen.ID("TwoFactorSecret").Op(":").ID("input").Dot("TwoFactorSecret")),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildCreateUserQuery").Call(jen.ID("input")),
			jen.Line(),
			jen.Comment("create the user"),
		}

		out = append(out, uqb...)
		out = append(out,
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		)
		return out
	}

	ret.Add(
		jen.Comment("CreateUser creates a user"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("CreateUser").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserInput")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User"), jen.ID("error")).Block(
			buildCreateUserBlock()...,
		),
		jen.Line(),
	)

	////////////

	buildUpdateUserQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("usersTableName")).
			Dotln("Set").Call(jen.Lit("username"), jen.ID("input").Dot("Username")).
			Dotln("Set").Call(jen.Lit("hashed_password"), jen.ID("input").Dot("HashedPassword")).
			Dotln("Set").Call(jen.Lit("two_factor_secret"), jen.ID("input").Dot("TwoFactorSecret")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("input").Dot("ID")))

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING updated_on"))
		}

		q.Dotln("ToSql").Call()
		return q
	}

	ret.Add(
		jen.Comment("buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildUpdateUserQuery").Params(jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildUpdateUserQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	buildUpdateUserBody := func() []jen.Code {
		if isPostgres {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("input")),
				jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("input").Dot("UpdatedOn")),
			}
		} else if isSqlite || isMariaDB {
			return []jen.Code{
				jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildUpdateUserQuery").Call(jen.ID("input")),
				jen.List(jen.ID("_"), jen.Err()).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
				jen.Return().ID("err"),
			}
		}

		return nil
	}

	ret.Add(
		jen.Comment("UpdateUser receives a complete User struct and updates its place in the db."),
		jen.Line(),
		jen.Comment("NOTE this function uses the ID provided in the input to make its query. Pass in"),
		jen.Line(),
		jen.Comment("anonymous structs or incomplete models at your peril."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("UpdateUser").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "User")).Params(jen.ID("error")).Block(
			buildUpdateUserBody()...,
		),
		jen.Line(),
	)

	////////////

	buildArchiveUserQueryQuery := func() jen.Code {
		q := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.ID("usersTableName")).
			Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("userID")))

		if isPostgres {
			q.Dotln("Suffix").Call(jen.Lit("RETURNING archived_on"))
		}

		q.Dotln("ToSql").Call()
		return q
	}

	ret.Add(
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("buildArchiveUserQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildArchiveUserQueryQuery(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Comment("ArchiveUser archives a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(sn)).ID("ArchiveUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dot("buildArchiveUserQuery").Call(jen.ID("userID")),
			jen.List(jen.ID("_"), jen.Err()).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
