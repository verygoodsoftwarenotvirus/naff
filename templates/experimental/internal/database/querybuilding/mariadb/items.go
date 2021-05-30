package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("ItemSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("MariaDB")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildItemExistsQuery constructs a SQL query for checking if an item with a given ID belong to a user with a given ID exists."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildItemExistsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)).Dot("Prefix").Call(jen.ID("querybuilding").Dot("ExistencePrefix")).Dot("From").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Suffix").Call(jen.ID("querybuilding").Dot("ExistenceSuffix")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("itemID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("ItemsTableAccountOwnershipColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("ItemsTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("itemID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("ItemsTableAccountOwnershipColumn"),
				).Op(":").ID("accountID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllItemsCountQuery returns a query that fetches the total number of items in the database."),
		jen.Line(),
		jen.Func().Comment("This query only gets generated once, and is otherwise returned from cache.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetAllItemsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
				)).Dot("From").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfItemsQuery returns a query that fetches every item in the database within a bucketed range."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetBatchOfItemsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("ItemsTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("beginID"))).Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given account,"),
		jen.Line(),
		jen.Func().Comment("and returns both the query and the relevant args to pass to the query executor.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetItemsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("querybuilding").Dot("ItemsTableName"),
				jen.ID("querybuilding").Dot("ItemsTableAccountOwnershipColumn"),
				jen.ID("querybuilding").Dot("ItemsTableColumns"),
				jen.ID("accountID"),
				jen.ID("forAdmin"),
				jen.ID("filter"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetItemsWithIDsQuery builds a SQL query selecting items that belong to a given account,"),
		jen.Line(),
		jen.Func().Comment("and have IDs that exist within a given set of IDs. Returns both the query and the relevant").Comment("args to pass to the query executor. This function is primarily intended for use with a search").Comment("index, which would provide a slice of string IDs to query against. This function accepts a").Comment("slice of uint64s instead of a slice of strings in order to ensure all the provided strings").Comment("are valid database IDs, because there's no way in squirrel to escape them in the unnest join,").Comment("and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.").Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetItemsWithIDsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("uint64"), jen.ID("forAdmin").ID("bool")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("whenThenStatement").Op(":=").ID("buildWhenThenStatement").Call(jen.ID("ids")),
			jen.ID("where").Op(":=").ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.ID("querybuilding").Dot("ItemsTableName"),
				jen.ID("querybuilding").Dot("IDColumn"),
			).Op(":").ID("ids"), jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.ID("querybuilding").Dot("ItemsTableName"),
				jen.ID("querybuilding").Dot("ArchivedOnColumn"),
			).Op(":").ID("nil")),
			jen.If(jen.Op("!").ID("forAdmin")).Body(
				jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("ItemsTableAccountOwnershipColumn"),
				)).Op("=").ID("accountID")),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("ItemsTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Where").Call(jen.ID("where")).Dot("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("CASE %s.%s %s"),
					jen.ID("querybuilding").Dot("ItemsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
					jen.ID("whenThenStatement"),
				)).Dot("Limit").Call(jen.ID("uint64").Call(jen.ID("limit"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildCreateItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("ItemCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("ExternalIDColumn"),
					jen.ID("querybuilding").Dot("ItemsTableNameColumn"),
					jen.ID("querybuilding").Dot("ItemsTableDetailsColumn"),
					jen.ID("querybuilding").Dot("ItemsTableAccountOwnershipColumn"),
				).Dot("Values").Call(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Name"),
					jen.ID("input").Dot("Details"),
					jen.ID("input").Dot("BelongsToAccount"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildUpdateItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("Item")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("BelongsToAccount"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ItemsTableNameColumn"),
					jen.ID("input").Dot("Name"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ItemsTableDetailsColumn"),
					jen.ID("input").Dot("Details"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("input").Dot("ID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"), jen.ID("querybuilding").Dot("ItemsTableAccountOwnershipColumn").Op(":").ID("input").Dot("BelongsToAccount"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildArchiveItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("ItemsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("itemID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"), jen.ID("querybuilding").Dot("ItemsTableAccountOwnershipColumn").Op(":").ID("accountID"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesForItemQuery constructs a SQL query for fetching an audit log entry with a given ID belong to a user with a given ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("MariaDB")).ID("BuildGetAuditLogEntriesForItemQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Expr").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("jsonPluckQuery"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableContextColumn"),
					jen.ID("itemID"),
					jen.ID("audit").Dot("ItemAssignmentKey"),
				))).Dot("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("CreatedOnColumn"),
				)),
			),
		),
		jen.Line(),
	)

	return code
}
