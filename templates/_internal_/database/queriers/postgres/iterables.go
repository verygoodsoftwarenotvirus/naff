package postgres

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"github.com/Masterminds/squirrel"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	pn := typ.Name.Plural()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("ItemDataManager").Equals().Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.Nil()),
			jen.Newline(),
			jen.Comment("itemsTableColumns are the columns for the items table."),
			jen.ID("itemsTableColumns").Equals().Index().String().Valuesln(jen.Lit("items.id"), jen.Lit("items.name"), jen.Lit("items.details"), jen.Lit("items.created_on"), jen.Lit("items.last_updated_on"), jen.Lit("items.archived_on"), jen.Lit("items.belongs_to_account")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an item struct."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").ID("database").Dot("Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("x").Op("*").ID("types").Dot("Item"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).Uint64(), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.Newline(),
			jen.ID("x").Equals().Op("&").ID("types").Dot("Item").Values(),
			jen.Newline(),
			jen.ID("targetVars").Op(":=").Index().Interface().Valuesln(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("Name"), jen.Op("&").ID("x").Dot("Details"), jen.Op("&").ID("x").Dot("CreatedOn"), jen.Op("&").ID("x").Dot("LastUpdatedOn"), jen.Op("&").ID("x").Dot("ArchivedOn"), jen.Op("&").ID("x").Dot("BelongsToAccount")),
			jen.Newline(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Equals().ID("append").Call(
					jen.ID("targetVars"),
					jen.Op("&").ID("filteredCount"),
					jen.Op("&").ID("totalCount"),
				)),
			jen.Newline(),
			jen.If(jen.ID("err").Equals().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Zero(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit(""),
				))),
			jen.Newline(),
			jen.Return().List(jen.ID("x"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.Nil()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("scanItems takes some database rows and turns them into a slice of items."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("items").Index().Op("*").ID("types").Dot("Item"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).Uint64(), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.Newline(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("x"), jen.ID("fc"), jen.ID("tc"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanItem").Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").DoesNotEqual().Nil()).Body(
					jen.Return().List(jen.Nil(), jen.Zero(), jen.Zero(), jen.ID("scanErr"))),
				jen.Newline(),
				jen.If(jen.ID("includeCounts")).Body(
					jen.If(jen.ID("filteredCount").Op("==").Zero()).Body(
						jen.ID("filteredCount").Equals().ID("fc")),
					jen.Newline(),
					jen.If(jen.ID("totalCount").Op("==").Zero()).Body(
						jen.ID("totalCount").Equals().ID("tc")),
				),
				jen.Newline(),
				jen.ID("items").Equals().ID("append").Call(
					jen.ID("items"),
					jen.ID("x"),
				),
			),
			jen.Newline(),
			jen.If(jen.ID("err").Equals().ID("q").Dot("checkRowsForErrorAndClose").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Zero(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("handling rows"),
				))),
			jen.Newline(),
			jen.Return().List(jen.ID("items"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.Nil()),
		),
		jen.Newline(),
	)

	query, _, err := sqlBuilder.Select(fmt.Sprintf("%s.id", tableName)).
		Prefix("SELECT EXISTS (").
		From(tableName).
		Suffix(")").
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", tableName):                 whatever,
			fmt.Sprintf("%s.archived_on", tableName):        nil,
			fmt.Sprintf("%s.belongs_to_account", tableName): whatever,
		}).ToSql()
	if err != nil {
		panic(err)
	}

	code.Add(
		jen.Const().ID("itemExistenceQuery").Equals().Lit(query),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ItemExists fetches whether an item exists from the database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("ItemExists").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).String()).Params(jen.ID("exists").ID("bool"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("itemID").Op("==").Lit("")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.If(jen.ID("accountID").Op("==").Lit("")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("accountID"), jen.ID("itemID")),
			jen.Newline(),
			jen.List(jen.ID("result"), jen.ID("err")).Op(":=").ID("q").Dot("performBooleanQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("itemExistenceQuery"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.ID("false"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("performing item existence check"),
				))),
			jen.Newline(),
			jen.Return().List(jen.ID("result"), jen.Nil()),
		),
		jen.Newline(),
	)

	selectColumns := []string{
		fmt.Sprintf("%s.id", tableName),
	}
	for _, field := range typ.Fields {
		selectColumns = append(selectColumns, fmt.Sprintf("%s.%s", tableName, field.Name.RouteName()))
	}
	selectColumns = append(selectColumns,
		fmt.Sprintf("%s.created_on", tableName),
		fmt.Sprintf("%s.last_updated_on", tableName),
		fmt.Sprintf("%s.archived_on", tableName),
	)

	singleSelectWhereClause := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName):          whatever,
		fmt.Sprintf("%s.archived_on", tableName): nil,
	}

	if typ.BelongsToStruct != nil {
		selectColumns = append(selectColumns, fmt.Sprintf("%s.belongs_to_%s", typ.BelongsToStruct.RouteName()))
		singleSelectWhereClause[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = whatever
	}
	if typ.BelongsToAccount {
		selectColumns = append(selectColumns, fmt.Sprintf("%s.belongs_to_account", tableName))
		singleSelectWhereClause[fmt.Sprintf("%s.belongs_to_account", tableName)] = whatever
	}

	query, _, err = sqlBuilder.Select(selectColumns...).From(tableName).Where(singleSelectWhereClause).ToSql()
	if err != nil {
		panic(err)
	}

	code.Add(
		jen.Const().ID("getItemQuery").Equals().Lit(query),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("GetItem fetches an item from the database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).String()).Params(jen.Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("itemID").Op("==").Lit("")).Body(
				jen.Return().List(jen.Nil(), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.If(jen.ID("accountID").Op("==").Lit("")).Body(
				jen.Return().List(jen.Nil(), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("accountID"), jen.ID("itemID")),
			jen.Newline(),
			jen.ID("row").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("item"),
				jen.ID("getItemQuery"),
				jen.ID("args"),
			),
			jen.Newline(),
			jen.List(jen.ID("item"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanItem").Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning item"),
				))),
			jen.Newline(),
			jen.Return().List(jen.ID("item"), jen.Nil()),
		),
		jen.Newline(),
	)

	query, _, err = sqlBuilder.Select(fmt.Sprintf("COUNT(%s.id)", tableName)).
		From(tableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", tableName): nil,
		}).ToSql()
	if err != nil {
		panic(err)
	}

	code.Add(
		jen.Const().ID("getAllItemsCountQuery").Equals().Lit(query),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("GetTotalItemCount fetches the count of items from the database that meet a particular filter."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetTotalItemCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Uint64(), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.Newline(),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("getAllItemsCountQuery"),
				jen.Lit("fetching count of items"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Zero(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for count of items"),
				))),
			jen.Newline(),
			jen.Return().List(jen.ID("count"), jen.Nil()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("GetItems fetches a list of items from the database that meet a particular filter."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").String(), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("x").Op("*").ID("types").Dot("ItemList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("accountID").Op("==").Lit("")).Body(
				jen.Return().List(jen.Nil(), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.ID("x").Equals().Op("&").ID("types").Dot("ItemList").Values(),
			jen.ID("logger").Equals().ID("filter").Dot("AttachToLogger").Call(jen.ID("logger")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.Newline(),
			jen.If(jen.ID("filter").DoesNotEqual().Nil()).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Equals().List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit")),
			),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("buildListQuery").Callln(
				jen.ID("ctx"),
				jen.Lit("items"),
				jen.Nil(),
				jen.Nil(),
				jen.ID("accountOwnershipColumn"),
				jen.ID("itemsTableColumns"),
				jen.ID("accountID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
			jen.Newline(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("items"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing items list retrieval query"),
				))),
			jen.Newline(),
			jen.If(jen.List(jen.ID("x").Dot("Items"), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Equals().ID("q").Dot("scanItems").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning items"),
				))),
			jen.Newline(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("buildGetItemsWithIDsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("string"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("string")).Params(jen.ID("string"), jen.Index().Interface()).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("withIDsWhere").Op(":=").ID("squirrel").Dot("Eq").Valuesln(
				jen.Lit("items.id").Op(":").ID("ids"),
				jen.Lit("items.archived_on").Op(":").ID("nil"),
				jen.Lit("items.belongs_to_account").Op(":").ID("accountID"),
			),
			jen.Newline(),
			jen.ID("subqueryBuilder").Op(":=").ID("q").Dot("sqlBuilder").Dot("Select").Call(jen.ID("itemsTableColumns").Op("...")).
				Dotln("From").Call(jen.Lit("items")).
				Dotln("Join").Call(jen.Lit("unnest('{%s}'::text[])")).
				Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d"), jen.ID("limit"))),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op(":=").ID("q").Dot("sqlBuilder").Dot("Select").Call(jen.ID("itemsTableColumns").Op("...")).
				Dotln("FromSelect").Call(jen.ID("subqueryBuilder"), jen.Lit("items")).
				Dotln("Where").Call(jen.ID("withIDsWhere")).Dot("ToSql").Call(),
			jen.ID("query").Op("=").Qual("fmt", "Sprintf").Call(
				jen.ID("query"),
				jen.ID("joinIDs").Call(jen.ID("ids")),
			),
			jen.Newline(),
			jen.ID("q").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("GetItemsWithIDs fetches items from the database within a given set of IDs."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetItemsWithIDs").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").String(), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().String()).Params(jen.Index().Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("accountID").Op("==").Lit("")).Body(
				jen.Return().List(jen.Nil(), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.If(jen.ID("limit").Op("==").Zero()).Body(
				jen.ID("limit").Equals().ID("uint8").Call(jen.ID("types").Dot("DefaultLimit")),
			),
			jen.Newline(),
			jen.ID("logger").Equals().ID("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(jen.Lit("limit").MapAssign().ID("limit"), jen.Lit("id_count").MapAssign().ID("len").Call(jen.ID("ids")))),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dotf("buildGet%sWithIDsQuery", pn).Call(
				constants.CtxVar(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("accountID")),
				jen.ID("limit"),
				jen.ID("ids"),
			),
			jen.Newline(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("items with IDs"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching items from database"),
				))),
			jen.Newline(),
			jen.List(jen.ID("items"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanItems").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning items"),
				))),
			jen.Newline(),
			jen.Return().List(jen.ID("items"), jen.Nil()),
		),
		jen.Newline(),
	)

	creationColumns := []string{
		"id",
	}
	args := []interface{}{whatever}
	for _, field := range typ.Fields {
		creationColumns = append(creationColumns, field.Name.RouteName())
		args = append(args, whatever)
	}

	if typ.BelongsToStruct != nil {
		creationColumns = append(creationColumns, typ.BelongsToStruct.RouteName())
		args = append(args, whatever)
	}

	if typ.BelongsToAccount {
		creationColumns = append(creationColumns, "belongs_to_account")
		args = append(args, whatever)
	}

	query, _, err = sqlBuilder.Insert(tableName).
		Columns(creationColumns...).
		Values(args...).
		ToSql()

	code.Add(
		jen.Const().ID("itemCreationQuery").Equals().Lit(query),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CreateItem creates an item in the database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("ItemDatabaseCreationInput")).Params(jen.Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("input").Op("==").Nil()).Body(
				jen.Return().List(jen.Nil(), jen.ID("ErrNilInputProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("input").Dot("ID"),
			),
			jen.Newline(),
			jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("input").Dot("ID"), jen.ID("input").Dot("Name"), jen.ID("input").Dot("Details"), jen.ID("input").Dot("BelongsToAccount")),
			jen.Newline(),
			jen.Comment("create the item."),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("item creation"),
				jen.ID("itemCreationQuery"),
				jen.ID("args"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating item"),
				))),
			jen.Newline(),
			jen.ID("x").Op(":=").Op("&").ID("types").Dot("Item").Valuesln(
				jen.ID("ID").MapAssign().ID("input").Dot("ID"),
				jen.ID("Name").MapAssign().ID("input").Dot("Name"),
				jen.ID("Details").MapAssign().ID("input").Dot("Details"),
				jen.ID("BelongsToAccount").MapAssign().ID("input").Dot("BelongsToAccount"),
				jen.ID("CreatedOn").MapAssign().ID("q").Dot("currentTime").Call(),
			),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("x").Dot("ID"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("item created")),
			jen.Newline(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		),
		jen.Newline(),
	)

	updateWhere := squirrel.Eq{
		"id":          whatever,
		"archived_on": nil,
	}

	if typ.BelongsToStruct != nil {
		updateWhere[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whatever
	}
	if typ.BelongsToAccount {
		updateWhere["belongs_to_account"] = whatever
	}

	updateBuilder := sqlBuilder.Update(tableName)

	for _, field := range typ.Fields {
		updateBuilder = updateBuilder.Set(field.Name.RouteName(), whatever)
	}

	updateBuilder = updateBuilder.Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).Where(updateWhere)

	query, _, err = updateBuilder.ToSql()
	if err != nil {
		panic(err)
	}

	code.Add(
		jen.Const().ID("updateItemQuery").Equals().Lit(query),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("Item")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("updated").Op("==").Nil()).Body(
				jen.Return().ID("ErrNilInputProvided"),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("updated").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("BelongsToAccount"),
			),
			jen.Newline(),
			jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("updated").Dot("Name"), jen.ID("updated").Dot("Details"), jen.ID("updated").Dot("BelongsToAccount"), jen.ID("updated").Dot("ID")),
			jen.Newline(),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("item update"),
				jen.ID("updateItemQuery"),
				jen.ID("args"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating item"),
				)),
			jen.Newline(),
			jen.ID("logger").Dot("Info").Call(jen.Lit("item updated")),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	archiveWhere := squirrel.Eq{
		"id":          whatever,
		"archived_on": nil,
	}

	if typ.BelongsToStruct != nil {
		archiveWhere[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whatever
	}
	if typ.BelongsToAccount {
		archiveWhere["belongs_to_account"] = whatever
	}

	query, _, err = sqlBuilder.Update(tableName).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(archiveWhere).ToSql()

	code.Add(
		jen.Const().ID("archiveItemQuery").Equals().Lit(query),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ArchiveItem archives an item from the database by its ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("itemID"), jen.ID("accountID")).String()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("itemID").Op("==").Lit("")).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.ID("tracing").Dot("AttachItemIDToSpan").Call(
				jen.ID("span"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.If(jen.ID("accountID").Op("==").Lit("")).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("accountID"), jen.ID("itemID")),
			jen.Newline(),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("item archive"),
				jen.ID("archiveItemQuery"),
				jen.ID("args"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating item"),
				)),
			jen.Newline(),
			jen.ID("logger").Dot("Info").Call(jen.Lit("item archived")),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	return code
}
