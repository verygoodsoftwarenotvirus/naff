package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsDotGo() *jen.File {
	ret := jen.NewFile("mariadb")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("itemsTableName").Op("=").Lit("items"),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("itemsTableColumns").Op("=").Index().ID("string").Valuesln(
			jen.Lit("id"), jen.Lit("name"), jen.Lit("details"), jen.Lit("created_on"), jen.Lit("updated_on"), jen.Lit("archived_on"), jen.Lit("belongs_to")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanItem takes a database Scanner (i.e. *sql.Row) and scans"),
		jen.Line(),
		jen.Comment("the result into an Item struct"),
		jen.Line(),
		jen.Func().ID("scanItem").Params(jen.ID("scan").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "Scanner")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item"),
			jen.ID("error")).Block(
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item").Values(),
			jen.If(jen.ID("err").Op(":=").ID("scan").Dot(
				"Scan",
			).Call(jen.Op("&").ID("x").Dot("ID"),
				jen.Op("&").ID("x").Dot("Name"),
				jen.Op("&").ID("x").Dot(
					"Details",
				),
				jen.Op("&").ID("x").Dot("CreatedOn"),
				jen.Op("&").ID("x").Dot("UpdatedOn"),
				jen.Op("&").ID("x").Dot("ArchivedOn"),
				jen.Op("&").ID("x").Dot("BelongsTo")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("scanItems takes a logger and some database rows and turns them into a slice of items"),
		jen.Line(),
		jen.Func().ID("scanItems").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1",
			"Logger",
		),
			jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item"),
			jen.ID("error")).Block(

			jen.Var().ID("list").Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item"),
			jen.For(jen.ID("rows").Dot(
				"Next",
			).Call()).Block(
				jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("scanItem").Call(jen.ID("rows")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.ID("list").Op("=").ID("append").Call(jen.ID("list"), jen.Op("*").ID("x")),
			),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot(
				"Err",
			).Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.If(jen.ID("closeErr").Op(":=").ID("rows").Dot(
				"Close",
			).Call(), jen.ID("closeErr").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("closeErr"), jen.Lit("closing database rows")),
			),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildGetItemQuery").Params(jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("m").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("itemsTableColumns").Op("...")).Dot(
				"From",
			).Call(jen.ID("itemsTableName")).Dot(
				"Where",
			).Call(jen.ID("squirrel").Dot(
				"Eq",
			).Valuesln(
				jen.Lit("id").Op(":").ID("itemID"), jen.Lit("belongs_to").Op(":").ID("userID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("m").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetItem fetches an item from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item"),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetItemQuery",
			).Call(jen.ID("itemID"), jen.ID("userID")),
			jen.ID("row").Op(":=").ID("m").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("scanItem").Call(jen.ID("row")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetItemCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for"),
		jen.Line(),
		jen.Comment("fetching the number of items belonging to a given user that meet a given query"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildGetItemCountQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("m").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("CountQuery")).Dot(
				"From",
			).Call(jen.ID("itemsTableName")).Dot(
				"Where",
			).Call(jen.ID("squirrel").Dot(
				"Eq",
			).Valuesln(
				jen.Lit("archived_on").Op(":").ID("nil"), jen.Lit("belongs_to").Op(":").ID("userID"))),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot(
					"ApplyToQueryBuilder",
				).Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
				"ToSql",
			).Call(),
			jen.ID("m").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetItemCount will fetch the count of items from the database that meet a particular filter and belong to a particular user."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("GetItemCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetItemCountQuery",
			).Call(jen.ID("filter"), jen.ID("userID")),
			jen.ID("err").Op("=").ID("m").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
				"Scan",
			).Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("allItemsCountQueryBuilder").Qual("sync", "Once"),
			jen.ID("allItemsCountQuery").ID("string"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database."),
		jen.Line(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildGetAllItemsCountQuery").Params().Params(jen.ID("string")).Block(
			jen.ID("allItemsCountQueryBuilder").Dot(
				"Do",
			).Call(jen.Func().Params().Block(

				jen.Var().ID("err").ID("error"),
				jen.List(jen.ID("allItemsCountQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID("m").Dot(
					"sqlBuilder",
				).Dot(
					"Select",
				).Call(jen.ID("CountQuery")).Dot(
					"From",
				).Call(jen.ID("itemsTableName")).Dot(
					"Where",
				).Call(jen.ID("squirrel").Dot(
					"Eq",
				).Valuesln(
					jen.Lit("archived_on").Op(":").ID("nil"))).Dot(
					"ToSql",
				).Call(),
				jen.ID("m").Dot(
					"logQueryBuildingError",
				).Call(jen.ID("err")),
			)),
			jen.Return().ID("allItemsCountQuery"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllItemsCount will fetch the count of items from the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("GetAllItemsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.ID("err").Op("=").ID("m").Dot("db").Dot(
				"QueryRowContext",
			).Call(jen.ID("ctx"), jen.ID("m").Dot(
				"buildGetAllItemsCountQuery",
			).Call()).Dot(
				"Scan",
			).Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,"),
		jen.Line(),
		jen.Comment("and returns both the query and the relevant args to pass to the query executor."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildGetItemsQuery").Params(jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID("m").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("itemsTableColumns").Op("...")).Dot(
				"From",
			).Call(jen.ID("itemsTableName")).Dot(
				"Where",
			).Call(jen.ID("squirrel").Dot(
				"Eq",
			).Valuesln(
				jen.Lit("archived_on").Op(":").ID("nil"), jen.Lit("belongs_to").Op(":").ID("userID"))),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot(
					"ApplyToQueryBuilder",
				).Call(jen.ID("builder")),
			),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
				"ToSql",
			).Call(),
			jen.ID("m").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetItems fetches a list of items from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"ItemList",
		),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetItemsQuery",
			).Call(jen.ID("filter"), jen.ID("userID")),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("m").Dot("db").Dot(
				"QueryContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("querying database for items"))),
			),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("scanItems").Call(jen.ID("m").Dot("logger"),
				jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
			),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItemCount",
			).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching item count: %w"), jen.ID("err"))),
			),
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
				"ItemList",
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
					jen.ID("TotalCount").Op(":").ID("count")), jen.ID("Items").Op(":").ID("list")),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllItemsForUser fetches every item belonging to a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("GetAllItemsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item"),
			jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetItemsQuery",
			).Call(jen.ID("nil"), jen.ID("userID")),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("m").Dot("db").Dot(
				"QueryContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("fetching items for user"))),
			),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("scanItems").Call(jen.ID("m").Dot("logger"),
				jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("parsing database results: %w"), jen.ID("err"))),
			),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildCreateItemQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("m").Dot(
				"sqlBuilder",
			).Dot(
				"Insert",
			).Call(jen.ID("itemsTableName")).Dot(
				"Columns",
			).Call(jen.Lit("name"), jen.Lit("details"), jen.Lit("belongs_to"), jen.Lit("created_on")).Dot(
				"Values",
			).Call(jen.ID("input").Dot("Name"),
				jen.ID("input").Dot(
					"Details",
				),
				jen.ID("input").Dot("BelongsTo"),
				jen.ID("squirrel").Dot(
					"Expr",
				).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("m").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildItemCreationTimeQuery").Params(jen.ID("itemID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("m").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.Lit("created_on")).Dot(
				"From",
			).Call(jen.ID("itemsTableName")).Dot(
				"Where",
			).Call(jen.ID("squirrel").Dot(
				"Eq",
			).Valuesln(
				jen.Lit("id").Op(":").ID("itemID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("m").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateItem creates an item in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1",
			"ItemCreationInput",
		)).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item"),
			jen.ID("error")).Block(
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item").Valuesln(
				jen.ID("Name").Op(":").ID("input").Dot("Name"),
				jen.ID("Details").Op(":").ID("input").Dot(
					"Details",
				),
				jen.ID("BelongsTo").Op(":").ID("input").Dot("BelongsTo")),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildCreateItemQuery",
			).Call(jen.ID("x")),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("m").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing item creation query: %w"), jen.ID("err"))),
			),
			jen.List(jen.ID("id"), jen.ID("idErr")).Op(":=").ID("res").Dot(
				"LastInsertId",
			).Call(),
			jen.If(jen.ID("idErr").Op("==").ID("nil")).Block(
				jen.ID("x").Dot("ID").Op("=").ID("uint64").Call(jen.ID("id")),
				jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
					"buildItemCreationTimeQuery",
				).Call(jen.ID("x").Dot("ID")),
				jen.ID("m").Dot(
					"logCreationTimeRetrievalError",
				).Call(jen.ID("m").Dot("db").Dot(
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
		jen.Comment("buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildUpdateItemQuery").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("m").Dot(
				"sqlBuilder",
			).Dot(
				"Update",
			).Call(jen.ID("itemsTableName")).Dot("Set").Call(jen.Lit("name"), jen.ID("input").Dot("Name")).Dot("Set").Call(jen.Lit("details"), jen.ID("input").Dot(
				"Details",
			)).Dot("Set").Call(jen.Lit("updated_on"), jen.ID("squirrel").Dot(
				"Expr",
			).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
				"Where",
			).Call(jen.ID("squirrel").Dot(
				"Eq",
			).Valuesln(
				jen.Lit("id").Op(":").ID("input").Dot("ID"),
				jen.Lit("belongs_to").Op(":").ID("input").Dot("BelongsTo"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("m").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Item")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildUpdateItemQuery",
			).Call(jen.ID("input")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("m").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("buildArchiveItemQuery").Params(jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("m").Dot(
				"sqlBuilder",
			).Dot(
				"Update",
			).Call(jen.ID("itemsTableName")).Dot("Set").Call(jen.Lit("updated_on"), jen.ID("squirrel").Dot(
				"Expr",
			).Call(jen.ID("CurrentUnixTimeQuery"))).Dot("Set").Call(jen.Lit("archived_on"), jen.ID("squirrel").Dot(
				"Expr",
			).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
				"Where",
			).Call(jen.ID("squirrel").Dot(
				"Eq",
			).Valuesln(
				jen.Lit("id").Op(":").ID("itemID"), jen.Lit("archived_on").Op(":").ID("nil"), jen.Lit("belongs_to").Op(":").ID("userID"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("m").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveItem marks an item as archived in the database"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MariaDB")).ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildArchiveItemQuery",
			).Call(jen.ID("itemID"), jen.ID("userID")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("m").Dot("db").Dot(
				"ExecContext",
			).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
