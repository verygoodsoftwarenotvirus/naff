package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsDotGo() *jen.File {
	ret := jen.NewFile("postgres")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("itemsTableName").Op("=").Lit("items"),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("itemsTableColumns").Op("=").Index().ID("string").Valuesln(jen.Lit("id"), jen.Lit("name"), jen.Lit("details"), jen.Lit("created_on"), jen.Lit("updated_on"), jen.Lit("archived_on"), jen.Lit("belongs_to")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// scanItem takes a database Scanner (i.e. *sql.Row) and scans").Comment("// the result into an Item struct").ID("scanItem").Params(jen.ID("scan").ID("database").Dot(
		"Scanner",
	)).Params(jen.Op("*").ID("models").Dot(
		"Item",
	), jen.ID("error")).Block(
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"Item",
		).Valuesln(),
		jen.If(jen.ID("err").Op(":=").ID("scan").Dot(
			"Scan",
		).Call(jen.Op("&").ID("x").Dot(
			"ID",
		), jen.Op("&").ID("x").Dot(
			"Name",
		), jen.Op("&").ID("x").Dot(
			"Details",
		), jen.Op("&").ID("x").Dot(
			"CreatedOn",
		), jen.Op("&").ID("x").Dot(
			"UpdatedOn",
		), jen.Op("&").ID("x").Dot(
			"ArchivedOn",
		), jen.Op("&").ID("x").Dot(
			"BelongsTo",
		)), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.Return().List(jen.ID("x"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// scanItems takes a logger and some database rows and turns them into a slice of items").ID("scanItems").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().ID("models").Dot(
		"Item",
	), jen.ID("error")).Block(
		jen.Null().Var().ID("list").Index().ID("models").Dot(
			"Item",
		),
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
			jen.ID("logger").Dot(
				"Error",
			).Call(jen.ID("closeErr"), jen.Lit("closing database rows")),
		),
		jen.Return().List(jen.ID("list"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetItemQuery").Params(jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot(
			"sqlBuilder",
		).Dot(
			"Select",
		).Call(jen.ID("itemsTableColumns").Op("...")).Dot(
			"From",
		).Call(jen.ID("itemsTableName")).Dot(
			"Where",
		).Call(jen.ID("squirrel").Dot(
			"Eq",
		).Valuesln(jen.Lit("id").Op(":").ID("itemID"), jen.Lit("belongs_to").Op(":").ID("userID"))).Dot(
			"ToSql",
		).Call(),
		jen.ID("p").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetItem fetches an item from the postgres database").Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"Item",
	), jen.ID("error")).Block(
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot(
			"buildGetItemQuery",
		).Call(jen.ID("itemID"), jen.ID("userID")),
		jen.ID("row").Op(":=").ID("p").Dot(
			"db",
		).Dot(
			"QueryRowContext",
		).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
		jen.Return().ID("scanItem").Call(jen.ID("row")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildGetItemCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for").Comment("// fetching the number of items belonging to a given user that meet a given query").Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetItemCountQuery").Params(jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.ID("builder").Op(":=").ID("p").Dot(
			"sqlBuilder",
		).Dot(
			"Select",
		).Call(jen.ID("CountQuery")).Dot(
			"From",
		).Call(jen.ID("itemsTableName")).Dot(
			"Where",
		).Call(jen.ID("squirrel").Dot(
			"Eq",
		).Valuesln(jen.Lit("archived_on").Op(":").ID("nil"), jen.Lit("belongs_to").Op(":").ID("userID"))),
		jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
			jen.ID("builder").Op("=").ID("filter").Dot(
				"ApplyToQueryBuilder",
			).Call(jen.ID("builder")),
		),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
			"ToSql",
		).Call(),
		jen.ID("p").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetItemCount will fetch the count of items from the database that meet a particular filter and belong to a particular user.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetItemCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot(
			"buildGetItemCountQuery",
		).Call(jen.ID("filter"), jen.ID("userID")),
		jen.ID("err").Op("=").ID("p").Dot(
			"db",
		).Dot(
			"QueryRowContext",
		).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
			"Scan",
		).Call(jen.Op("&").ID("count")),
		jen.Return().List(jen.ID("count"), jen.ID("err")),
	),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("allItemsCountQueryBuilder").Qual("sync", "Once").Var().ID("allItemsCountQuery").ID("string"),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database.").Comment("// This query only gets generated once, and is otherwise returned from cache.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetAllItemsCountQuery").Params().Params(jen.ID("string")).Block(
		jen.ID("allItemsCountQueryBuilder").Dot(
			"Do",
		).Call(jen.Func().Params().Block(
			jen.Null().Var().ID("err").ID("error"),
			jen.List(jen.ID("allItemsCountQuery"), jen.ID("_"), jen.ID("err")).Op("=").ID("p").Dot(
				"sqlBuilder",
			).Dot(
				"Select",
			).Call(jen.ID("CountQuery")).Dot(
				"From",
			).Call(jen.ID("itemsTableName")).Dot(
				"Where",
			).Call(jen.ID("squirrel").Dot(
				"Eq",
			).Valuesln(jen.Lit("archived_on").Op(":").ID("nil"))).Dot(
				"ToSql",
			).Call(),
			jen.ID("p").Dot(
				"logQueryBuildingError",
			).Call(jen.ID("err")),
		)),
		jen.Return().ID("allItemsCountQuery"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllItemsCount will fetch the count of items from the database").Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetAllItemsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
		jen.ID("err").Op("=").ID("p").Dot(
			"db",
		).Dot(
			"QueryRowContext",
		).Call(jen.ID("ctx"), jen.ID("p").Dot(
			"buildGetAllItemsCountQuery",
		).Call()).Dot(
			"Scan",
		).Call(jen.Op("&").ID("count")),
		jen.Return().List(jen.ID("count"), jen.ID("err")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,").Comment("// and returns both the query and the relevant args to pass to the query executor.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildGetItemsQuery").Params(jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.ID("builder").Op(":=").ID("p").Dot(
			"sqlBuilder",
		).Dot(
			"Select",
		).Call(jen.ID("itemsTableColumns").Op("...")).Dot(
			"From",
		).Call(jen.ID("itemsTableName")).Dot(
			"Where",
		).Call(jen.ID("squirrel").Dot(
			"Eq",
		).Valuesln(jen.Lit("archived_on").Op(":").ID("nil"), jen.Lit("belongs_to").Op(":").ID("userID"))),
		jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
			jen.ID("builder").Op("=").ID("filter").Dot(
				"ApplyToQueryBuilder",
			).Call(jen.ID("builder")),
		),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot(
			"ToSql",
		).Call(),
		jen.ID("p").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetItems fetches a list of items from the database that meet a particular filter").Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("models").Dot(
		"QueryFilter",
	), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("models").Dot(
		"ItemList",
	), jen.ID("error")).Block(
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot(
			"buildGetItemsQuery",
		).Call(jen.ID("filter"), jen.ID("userID")),
		jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("p").Dot(
			"db",
		).Dot(
			"QueryContext",
		).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("querying database for items"))),
		),
		jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("scanItems").Call(jen.ID("p").Dot(
			"logger",
		), jen.ID("rows")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
		),
		jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("p").Dot(
			"GetItemCount",
		).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching item count: %w"), jen.ID("err"))),
		),
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"ItemList",
		).Valuesln(jen.ID("Pagination").Op(":").ID("models").Dot(
			"Pagination",
		).Valuesln(jen.ID("Page").Op(":").ID("filter").Dot(
			"Page",
		), jen.ID("Limit").Op(":").ID("filter").Dot(
			"Limit",
		), jen.ID("TotalCount").Op(":").ID("count")), jen.ID("Items").Op(":").ID("list")),
		jen.Return().List(jen.ID("x"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// GetAllItemsForUser fetches every item belonging to a user").Params(jen.ID("p").Op("*").ID("Postgres")).ID("GetAllItemsForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID("models").Dot(
		"Item",
	), jen.ID("error")).Block(
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot(
			"buildGetItemsQuery",
		).Call(jen.ID("nil"), jen.ID("userID")),
		jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("p").Dot(
			"db",
		).Dot(
			"QueryContext",
		).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Lit("fetching items for user"))),
		),
		jen.List(jen.ID("list"), jen.ID("err")).Op(":=").ID("scanItems").Call(jen.ID("p").Dot(
			"logger",
		), jen.ID("rows")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("parsing database results: %w"), jen.ID("err"))),
		),
		jen.Return().List(jen.ID("list"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildCreateItemQuery").Params(jen.ID("input").Op("*").ID("models").Dot(
		"Item",
	)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot(
			"sqlBuilder",
		).Dot(
			"Insert",
		).Call(jen.ID("itemsTableName")).Dot(
			"Columns",
		).Call(jen.Lit("name"), jen.Lit("details"), jen.Lit("belongs_to")).Dot(
			"Values",
		).Call(jen.ID("input").Dot(
			"Name",
		), jen.ID("input").Dot(
			"Details",
		), jen.ID("input").Dot(
			"BelongsTo",
		)).Dot(
			"Suffix",
		).Call(jen.Lit("RETURNING id, created_on")).Dot(
			"ToSql",
		).Call(),
		jen.ID("p").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CreateItem creates an item in the database").Params(jen.ID("p").Op("*").ID("Postgres")).ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
		"ItemCreationInput",
	)).Params(jen.Op("*").ID("models").Dot(
		"Item",
	), jen.ID("error")).Block(
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"Item",
		).Valuesln(jen.ID("Name").Op(":").ID("input").Dot(
			"Name",
		), jen.ID("Details").Op(":").ID("input").Dot(
			"Details",
		), jen.ID("BelongsTo").Op(":").ID("input").Dot(
			"BelongsTo",
		)),
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot(
			"buildCreateItemQuery",
		).Call(jen.ID("x")),
		jen.ID("err").Op(":=").ID("p").Dot(
			"db",
		).Dot(
			"QueryRowContext",
		).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
			"Scan",
		).Call(jen.Op("&").ID("x").Dot(
			"ID",
		), jen.Op("&").ID("x").Dot(
			"CreatedOn",
		)),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error executing item creation query: %w"), jen.ID("err"))),
		),
		jen.Return().List(jen.ID("x"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters").Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildUpdateItemQuery").Params(jen.ID("input").Op("*").ID("models").Dot(
		"Item",
	)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot(
			"sqlBuilder",
		).Dot(
			"Update",
		).Call(jen.ID("itemsTableName")).Dot(
			"Set",
		).Call(jen.Lit("name"), jen.ID("input").Dot(
			"Name",
		)).Dot(
			"Set",
		).Call(jen.Lit("details"), jen.ID("input").Dot(
			"Details",
		)).Dot(
			"Set",
		).Call(jen.Lit("updated_on"), jen.ID("squirrel").Dot(
			"Expr",
		).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
			"Where",
		).Call(jen.ID("squirrel").Dot(
			"Eq",
		).Valuesln(jen.Lit("id").Op(":").ID("input").Dot(
			"ID",
		), jen.Lit("belongs_to").Op(":").ID("input").Dot(
			"BelongsTo",
		))).Dot(
			"Suffix",
		).Call(jen.Lit("RETURNING updated_on")).Dot(
			"ToSql",
		).Call(),
		jen.ID("p").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("models").Dot(
		"Item",
	)).Params(jen.ID("error")).Block(
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot(
			"buildUpdateItemQuery",
		).Call(jen.ID("input")),
		jen.Return().ID("p").Dot(
			"db",
		).Dot(
			"QueryRowContext",
		).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot(
			"Scan",
		).Call(jen.Op("&").ID("input").Dot(
			"UpdatedOn",
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived.").Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildArchiveItemQuery").Params(jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot(
			"sqlBuilder",
		).Dot(
			"Update",
		).Call(jen.ID("itemsTableName")).Dot(
			"Set",
		).Call(jen.Lit("updated_on"), jen.ID("squirrel").Dot(
			"Expr",
		).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
			"Set",
		).Call(jen.Lit("archived_on"), jen.ID("squirrel").Dot(
			"Expr",
		).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
			"Where",
		).Call(jen.ID("squirrel").Dot(
			"Eq",
		).Valuesln(jen.Lit("id").Op(":").ID("itemID"), jen.Lit("archived_on").Op(":").ID("nil"), jen.Lit("belongs_to").Op(":").ID("userID"))).Dot(
			"Suffix",
		).Call(jen.Lit("RETURNING archived_on")).Dot(
			"ToSql",
		).Call(),
		jen.ID("p").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ArchiveItem marks an item as archived in the database").Params(jen.ID("p").Op("*").ID("Postgres")).ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("p").Dot(
			"buildArchiveItemQuery",
		).Call(jen.ID("itemID"), jen.ID("userID")),
		jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("p").Dot(
			"db",
		).Dot(
			"ExecContext",
		).Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
		jen.Return().ID("err"),
	),

		jen.Line(),
	)
	return ret
}