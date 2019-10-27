package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersDotGo() *jen.File {
	ret := jen.NewFile("postgres")

	utils.AddImports(ret)

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
		jen.Func().ID("scanUser").Params(jen.ID("scan").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Scanner")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User"), jen.ID("error")).Block(
			jen.Var().ID("x").Op("=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User").Values(),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("scan").Dot("Scan").Callln(
				jen.Op("&").ID("x").Dot("ID"),
				jen.Op("&").ID("x").Dot("Username"),
				jen.Op("&").ID("x").Dot("HashedPassword"),
				jen.Op("&").ID("x").Dot("PasswordLastChangedOn"),
				jen.Op("&").ID("x").Dot("TwoFactorSecret"),
				jen.Op("&").ID("x").Dot("IsAdmin"),
				jen.Op("&").ID("x").Dot("CreatedOn"),
				jen.Op("&").ID("x").Dot("UpdatedOn"),
				jen.Op("&").ID("x").Dot("ArchivedOn"),
			), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanUsers takes database rows and loads them into a slice of User structs"),
		jen.Line(),
		jen.Func().ID("scanUsers").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User"), jen.ID("error")).Block(
			jen.Var().ID("list").Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User"),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("scanUser").Call(jen.ID("rows")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning user result: %w"), jen.ID("err"))),
				),
				jen.ID("list").Op("=").ID("append").Call(jen.ID("list"), jen.Op("*").ID("user")),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Err").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("closing rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetUserQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Values(jen.Lit("id").Op(":").ID("userID"))).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID("p").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUser fetches a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot("buildGetUserQuery").Call(jen.ID("userID")),
			jen.ID("row").Op(":=").ID("p").Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("scanUser").Call(jen.ID("row")),
			jen.Line(),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("fetching user from database"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetUserByUsernameQuery").Params(jen.ID("username").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("username").Op(":").ID("username"))).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID("p").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserByUsername fetches a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot("buildGetUserByUsernameQuery").Call(jen.ID("username")),
			jen.ID("row").Op(":=").ID("p").Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("scanUser").Call(jen.ID("row")),
			jen.Line(),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetUserCountQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("p").Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("CountQuery")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.Line(),
			jen.ID("p").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserCount fetches a count of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetUserCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot("buildGetUserCountQuery").Call(jen.ID("filter")),
			jen.ID("err").Op("=").ID("p").Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserCountQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetUsersQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("p").Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("usersTableColumns").Op("...")).
				Dotln("From").Call(jen.ID("usersTableName")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.ID("p").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUsers fetches a list of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserList"), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot("buildGetUsersQuery").Call(jen.ID("filter")),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("p").Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("querying for user"))),
			),
			jen.Line(),
			jen.List(jen.ID("userList"), jen.ID("err")).Op(":=").ID("scanUsers").Call(jen.ID("p").Dot("logger"), jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading response from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("p").Dot("GetUserCount").Call(jen.ID("ctx"), jen.ID("filter")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user count: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserList").Valuesln(
				jen.ID("Pagination").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Pagination").Valuesln(
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

	ret.Add(
		jen.Comment("buildCreateUserQuery returns a SQL query (and arguments) that would create a given User"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildCreateUserQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot("sqlBuilder").
				Dotln("Insert").Call(jen.ID("usersTableName")).
				Dotln("Columns").Callln(
				jen.Lit("username"),
				jen.Lit("hashed_password"),
				jen.Lit("two_factor_secret"),
				jen.Lit("is_admin"),
			).
				Dotln("Values").Callln(
				jen.ID("input").Dot("Username"),
				jen.ID("input").Dot("Password"),
				jen.ID("input").Dot("TwoFactorSecret"),
				jen.ID("false"),
			).
				Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on")).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.Comment("NOTE: we always default is_admin to false, on the assumption that"),
			jen.Comment("admins have DB access and will change that value via SQL query."),
			jen.Comment("There should also be no way to update a user via this structure"),
			jen.Comment("such that they would have admin privileges"),
			jen.Line(),
			jen.ID("p").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateUser creates a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserInput")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User"), jen.ID("error")).Block(
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User").Valuesln(
				jen.ID("Username").Op(":").ID("input").Dot("Username"),
				jen.ID("TwoFactorSecret").Op(":").ID("input").Dot("TwoFactorSecret")),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot("buildCreateUserQuery").Call(jen.ID("input")),
			jen.Line(),
			jen.Comment("create the user"),
			jen.ID("err").Op(":=").ID("p").Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("CreatedOn")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Switch(jen.ID("e").Op(":=").ID("err").Assert(jen.Type())).Block(
					jen.Case(jen.Op("*").ID("pq").Dot("Error")).Block(
						jen.If(jen.ID("e").Dot("Code").Op("==").ID("pq").Dot("ErrorCode").Call(jen.Lit("23505"))).Block(
							jen.Return().List(jen.ID("nil"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client", "ErrUserExists")),
						),
					),
					jen.Default().Block(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing user creation query: %w"), jen.ID("err"))),
					),
				),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildUpdateUserQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot("sqlBuilder").
				Dotln("Update").Call(jen.ID("usersTableName")).
				Dotln("Set").Call(jen.Lit("username"), jen.ID("input").Dot("Username")).
				Dotln("Set").Call(jen.Lit("hashed_password"), jen.ID("input").Dot("HashedPassword")).
				Dotln("Set").Call(jen.Lit("two_factor_secret"), jen.ID("input").Dot("TwoFactorSecret")).
				Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("input").Dot("ID"))).
				Dotln("Suffix").Call(jen.Lit("RETURNING updated_on")).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID("p").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateUser receives a complete User struct and updates its place in the db."),
		jen.Line(),
		jen.Comment("NOTE this function uses the ID provided in the input to make its query. Pass in"),
		jen.Line(),
		jen.Comment("anonymous structs or incomplete models at your peril."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "User")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot("buildUpdateUserQuery").Call(jen.ID("input")),
			jen.Return().ID("p").Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("input").Dot("UpdatedOn")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildArchiveUserQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot("sqlBuilder").
				Dotln("Update").Call(jen.ID("usersTableName")).
				Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
				Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").ID("userID"))).
				Dotln("Suffix").Call(jen.Lit("RETURNING archived_on")).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID("p").Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveUser archives a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot("buildArchiveUserQuery").Call(jen.ID("userID")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("p").Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
