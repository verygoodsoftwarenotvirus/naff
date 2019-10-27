package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersDotGo() *jen.File {
	ret := jen.NewFile("sqlite")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("usersTableName").Op("=").Lit("users"),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("usersTableColumns").Op("=").Index().ID("string").Valuesln(
			jen.Lit("id"), jen.Lit("username"), jen.Lit("hashed_password"), jen.Lit("password_last_changed_on"), jen.Lit("two_factor_secret"), jen.Lit("is_admin"), jen.Lit("created_on"), jen.Lit("updated_on"), jen.Lit("archived_on")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanUser provides a consistent way to scan something like a *sql.Row into a User struct"),
		jen.Line(),
		jen.Func().ID("scanUser").Params(jen.ID("scan").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1",
			"Scanner",
		)).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"User",
		),
			jen.ID("error")).Block(

			jen.Var().ID("x").Op("=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"User",
			).Values(),
			jen.If(jen.ID("err").Op(":=").ID("scan").Dot(
				"Scan",
			).Call(jen.Op("&").ID("x").Dot("ID"),
				jen.Op("&").ID("x").Dot(
					"Username",
				),
				jen.Op("&").ID("x").Dot(
					"HashedPassword",
				),
				jen.Op("&").ID("x").Dot(
					"PasswordLastChangedOn",
				),
				jen.Op("&").ID("x").Dot(
					"TwoFactorSecret",
				),
				jen.Op("&").ID("x").Dot(
					"IsAdmin",
				),
				jen.Op("&").ID("x").Dot("CreatedOn"),
				jen.Op("&").ID("x").Dot("UpdatedOn"),
				jen.Op("&").ID("x").Dot("ArchivedOn")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanUsers takes database rows and loads them into a slice of User structs"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("scanUsers").Params(jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"User",
		),
			jen.ID("error")).Block(

			jen.Var().ID("list").Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"User",
			),
			jen.For(jen.ID("rows").Dot(
				"Next",
			).Call()).Block(
				jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("scanUser").Call(jen.ID("rows")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning user result: %w"), jen.ID("err"))),
				),
				jen.ID("list").Op("=").ID("append").Call(jen.ID("list"), jen.Op("*").ID("user")),
			),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot(
				"Err",
			).Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("rows").Dot(
				"Close",
			).Call()),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetUserQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("usersTableColumns").Op("...")).Dot(
				"From",
			).Call(jen.ID("usersTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("userID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUser fetches a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"User",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetUserQuery",
			).Call(jen.ID("userID")),
			jen.ID("row").Op(":=").ID("s").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("scanUser").Call(jen.ID("row")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("fetching user from database"))),
			),
			jen.Return().List(jen.ID("u"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetUserByUsernameQuery").Params(jen.ID("username").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("usersTableColumns").Op("...")).Dot(
				"From",
			).Call(jen.ID("usersTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("username").Op(":").ID("username"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserByUsername fetches a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"User",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetUserByUsernameQuery",
			).Call(jen.ID("username")),
			jen.ID("row").Op(":=").ID("s").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("scanUser").Call(jen.ID("row")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user from database: %w"), jen.ID("err"))),
			),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetUserCountQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("CountQuery")).Dot(
				"From",
			).Call(jen.ID("usersTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot(
					"ApplyToQueryBuilder",
				).Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserCount fetches a count of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetUserCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetUserCountQuery",
			).Call(jen.ID("filter")),
			jen.ID("err").Op("=").ID("s").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
				"Scan",
			).Call(jen.Op("&").ID("count")),
			jen.Return(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetUserCountQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere"),
		jen.Line(),
		jen.Comment("to a given filter's criteria."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildGetUsersQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("usersTableColumns").Op("...")).Dot(
				"From",
			).Call(jen.ID("usersTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("archived_on").Op(":").ID("nil"))),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot(
					"ApplyToQueryBuilder",
				).Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUsers fetches a list of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"UserList",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildGetUsersQuery",
			).Call(jen.ID("filter")),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"QueryContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("querying for user"))),
			),
			jen.List(jen.ID("userList"), jen.ID("err")).Op(":=").ID("s").Dot(
				"scanUsers",
			).Call(jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading response from database: %w"), jen.ID("err"))),
			),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("s").Dot(
				"GetUserCount",
			).Call(jen.ID("ctx"), jen.ID("filter")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user count: %w"), jen.ID("err"))),
			),
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"UserList",
			).Valuesln(
				jen.ID("Pagination").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
					"Pagination",
				).Valuesln(
					jen.ID("Page").Op(":").ID("filter").Dot(
						"Page",
					),
					jen.ID("Limit").Op(":").ID("filter").Dot(
						"Limit",
					),
					jen.ID("TotalCount").Op(":").ID("count")), jen.ID("Users").Op(":").ID("userList")),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildCreateUserQuery returns a SQL query (and arguments) that would create a given User"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildCreateUserQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"UserInput",
		)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Insert",
			).Call(jen.ID("usersTableName")).Dot(
				"Columns",
			).Call(jen.Lit("username"), jen.Lit("hashed_password"), jen.Lit("two_factor_secret"), jen.Lit("is_admin")).Dot(
				"Values",
			).Call(jen.ID("input").Dot(
				"Username",
			),
				jen.ID("input").Dot(
					"Password",
				),
				jen.ID("input").Dot(
					"TwoFactorSecret",
				),
				jen.ID("false")).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildCreateUserQuery returns a SQL query (and arguments) that would create a given User"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildUserCreationTimeQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.Lit("created_on")).Dot(
				"From",
			).Call(jen.ID("usersTableName")).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("userID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateUser creates a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"UserInput",
		)).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"User",
		),
			jen.ID("error")).Block(
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"User",
			).Valuesln(
				jen.ID("Username").Op(":").ID("input").Dot(
					"Username",
				),
				jen.ID("TwoFactorSecret").Op(":").ID("input").Dot(
					"TwoFactorSecret",
				)),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildCreateUserQuery",
			).Call(jen.ID("input")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing user creation query: %w"), jen.ID("err"))),
			),
			jen.If(jen.List(jen.ID("id"), jen.ID("idErr")).Op(":=").ID("res").Dot(
				"LastInsertId",
			).Call(), jen.ID("idErr").Op("==").ID("nil")).Block(
				jen.ID("x").Dot("ID").Op("=").ID("uint64").Call(jen.ID("id")),
				jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
					"buildUserCreationTimeQuery",
				).Call(jen.ID("x").Dot("ID")),
				jen.ID("s").Dot(
					"logCreationTimeRetrievalError",
				).Call(jen.ID("s").Dot("db").Dot(
					"QueryRowContext",
				).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
					"Scan",
				).Call(jen.Op("&").ID("x").Dot("CreatedOn"))),
			),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildUpdateUserQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"User",
		)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Update",
			).Call(jen.ID("usersTableName")).Dot("Set").Call(jen.Lit("username"), jen.ID("input").Dot(
				"Username",
			)).Dot("Set").Call(jen.Lit("hashed_password"), jen.ID("input").Dot(
				"HashedPassword",
			)).Dot("Set").Call(jen.Lit("two_factor_secret"), jen.ID("input").Dot(
				"TwoFactorSecret",
			)).Dot("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("input").Dot("ID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
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
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"User",
		)).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildUpdateUserQuery",
			).Call(jen.ID("input")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("buildArchiveUserQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("s").Dot(
				"sqlBuilder",
			).Dot(
				"Update",
			).Call(jen.ID("usersTableName")).Dot("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).Dot("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
				"Where",
			).Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").ID("userID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("s").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveUser archives a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Sqlite")).ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("s").Dot(
				"buildArchiveUserQuery",
			).Call(jen.ID("userID")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("s").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
