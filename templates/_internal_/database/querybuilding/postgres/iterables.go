package postgres

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	pcn := typ.Name.PluralCommonName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	code.Add(
		jen.Var().ID("_").ID("querybuilding").Dotf("%sSQLQueryBuilder", sn).Op("=").Parens(jen.Op("*").ID("Postgres")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Build%sExistsQuery constructs a SQL query for checking if %s with a given ID belong to a user with a given ID exists.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("Build%sExistsQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.IDf("%sID", uvn),
				jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
					jen.Qual(proj.QuerybuildersPackage(), "IDColumn"),
				)).
					Dotln("Prefix").Call(jen.Qual(proj.QuerybuildersPackage(), "ExistencePrefix")).
					Dotln("From").Call(jen.ID("querybuilding").Dotf("%sTableName", pn)).
					Dotln("Suffix").Call(jen.Qual(proj.QuerybuildersPackage(), "ExistenceSuffix")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
					jen.Qual(proj.QuerybuildersPackage(), "IDColumn"),
				).MapAssign().IDf("%sID", uvn),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("querybuilding").Dotf("%sTableName", pn),
						jen.ID("querybuilding").Dotf("%sTableAccountOwnershipColumn", pn),
					).MapAssign().ID("accountID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("querybuilding").Dotf("%sTableName", pn),
						jen.Qual(proj.QuerybuildersPackage(), "ArchivedOnColumn"),
					).MapAssign().ID("nil"),
				),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to a user with a given ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildGet%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.IDf("%sID", uvn),
				jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.ID("querybuilding").Dotf("%sTableColumns", pn).Op("...")).
					Dotln("From").Call(jen.ID("querybuilding").Dotf("%sTableName", pn)).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
					jen.Qual(proj.QuerybuildersPackage(), "IDColumn"),
				).MapAssign().IDf("%sID", uvn),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("querybuilding").Dotf("%sTableName", pn),
						jen.ID("querybuilding").Dotf("%sTableAccountOwnershipColumn", pn),
					).MapAssign().ID("accountID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("querybuilding").Dotf("%sTableName", pn),
						jen.Qual(proj.QuerybuildersPackage(), "ArchivedOnColumn"),
					).MapAssign().ID("nil"),
				),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGetAll%sCountQuery returns a query that fetches the total number of %s in the database.", pn, pcn),
		jen.Newline(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildGetAll%sCountQuery", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),

			jen.Return().ID("b").Dot("buildQueryOnly").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
				)).
					Dotln("From").Call(jen.ID("querybuilding").Dotf("%sTableName", pn)).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
					jen.Qual(proj.QuerybuildersPackage(), "ArchivedOnColumn"),
				).MapAssign().ID("nil"),
				),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGetBatchOf%sQuery returns a query that fetches every %s in the database within a bucketed range.", pn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildGetBatchOf%sQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("beginID"),
				jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.ID("querybuilding").Dotf("%sTableColumns", pn).Op("...")).
					Dotln("From").Call(jen.ID("querybuilding").Dotf("%sTableName", pn)).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
					jen.Qual(proj.QuerybuildersPackage(), "IDColumn"),
				).MapAssign().ID("beginID"))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
					jen.Qual(proj.QuerybuildersPackage(), "IDColumn"),
				).MapAssign().ID("endID"),
				),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given account,", pn, pcn),
		jen.Newline(),
		jen.Comment("and returns both the query and the relevant args to pass to the query executor."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildGet%sQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("accountID").ID("uint64"),
			jen.ID("forAdmin").ID("bool"),
			jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
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
			jen.Return().ID("b").Dot("buildListQuery").Call(
				jen.ID("ctx"),
				jen.ID("querybuilding").Dotf("%sTableName", pn),
				jen.ID("querybuilding").Dotf("%sTableAccountOwnershipColumn", pn),
				jen.ID("querybuilding").Dotf("%sTableColumns", pn),
				jen.ID("accountID"),
				jen.ID("forAdmin"),
				jen.ID("filter"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGet%sWithIDsQuery builds a SQL query selecting %s that belong to a given account,", pn, pcn),
		jen.Newline(),
		jen.Comment("and have IDs that exist within a given set of IDs. Returns both the query and the relevant"),
		jen.Newline(),
		jen.Comment("args to pass to the query executor. This function is primarily intended for use with a search"),
		jen.Newline(),
		jen.Comment("index, which would provide a slice of string IDs to query against. This function accepts a"),
		jen.Newline(),
		jen.Comment("slice of uint64s instead of a slice of strings in order to ensure all the provided strings"),
		jen.Newline(),
		jen.Comment("are valid database IDs, because there's no way in squirrel to escape them in the unnest join,"),
		jen.Newline(),
		jen.Comment("and if we accept strings we could leave ourselves vulnerable to SQL injection attacks."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildGet%sWithIDsQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("accountID").ID("uint64"),
			jen.ID("limit").ID("uint8"),
			jen.ID("ids").Index().ID("uint64"),
			jen.ID("forAdmin").ID("bool")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.ID("where").Op(":=").ID("squirrel").Dot("Eq").Valuesln(
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.ID("querybuilding").Dotf("%sTableName", pn), jen.Qual(proj.QuerybuildersPackage(), "ArchivedOnColumn")).MapAssign().ID("nil"),
			),
			jen.Newline(),
			jen.If(jen.Op("!").ID("forAdmin")).Body(
				jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dotf("%sTableName", pn),
					jen.ID("querybuilding").Dotf("%sTableAccountOwnershipColumn", pn),
				)).Op("=").ID("accountID"),
			),
			jen.Newline(),
			jen.ID("subqueryBuilder").Assign().ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildersPackage(), fmt.Sprintf("%sTableColumns", pn)).Spread()).
				Dotln("From").Call(jen.Qual(proj.QuerybuildersPackage(), fmt.Sprintf("%sTableName", pn))).
				Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("unnest('{%s}'::int[])"), jen.ID("joinUint64s").Call(jen.ID("ids")))).
				Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d"), jen.ID("limit"))),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.ID("querybuilding").Dotf("%sTableColumns", pn).Op("...")).
					Dotln("FromSelect").Call(jen.ID("subqueryBuilder"), jen.Qual(proj.QuerybuildersPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Where").Call(jen.ID("where")),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildCreate%sQuery takes %s and returns a creation query for that %s and the relevant arguments.", sn, scnwp, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildCreate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").Op("*").ID("types").Dotf("%sCreationInput", sn)).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),

			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dotf("%sTableName", pn)).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildersPackage(), "ExternalIDColumn"),
					jen.ID("querybuilding").Dotf("%sTableNameColumn", pn),
					jen.ID("querybuilding").Dotf("%sTableDetailsColumn", pn),
					jen.ID("querybuilding").Dotf("%sTableAccountOwnershipColumn", pn),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Name"),
					jen.ID("input").Dot("Details"),
					jen.ID("input").Dot("BelongsToAccount"),
				).
					Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildersPackage(), "IDColumn"))),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildUpdate%sQuery takes %s and returns an update SQL query, with the relevant query parameters.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildUpdate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").Op("*").ID("types").Dot(sn)).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("BelongsToAccount"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dotf("%sTableName", pn)).
					Dotln("Set").Call(
					jen.ID("querybuilding").Dotf("%sTableNameColumn", pn),
					jen.ID("input").Dot("Name"),
				).
					Dotln("Set").Call(
					jen.ID("querybuilding").Dotf("%sTableDetailsColumn", pn),
					jen.ID("input").Dot("Details"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildersPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildersPackage(), "IDColumn").MapAssign().ID("input").Dot("ID"),
					jen.Qual(proj.QuerybuildersPackage(), "ArchivedOnColumn").MapAssign().ID("nil"),
					jen.ID("querybuilding").Dotf("%sTableAccountOwnershipColumn", pn).MapAssign().ID("input").Dot("BelongsToAccount"),
				),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildArchive%sQuery returns a SQL query which marks a given %s belonging to a given account as archived.", sn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildArchive%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.IDf("%sID", uvn),
				jen.ID("accountID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dotf("%sTableName", pn)).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildersPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildersPackage(), "ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildersPackage(), "IDColumn").MapAssign().IDf("%sID", uvn),
					jen.Qual(proj.QuerybuildersPackage(), "ArchivedOnColumn").MapAssign().ID("nil"),
					jen.ID("querybuilding").Dotf("%sTableAccountOwnershipColumn", pn).MapAssign().ID("accountID"),
				),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildGetAuditLogEntriesFor%sQuery constructs a SQL query for fetching audit log entries relating to %s with a given ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).IDf("BuildGetAuditLogEntriesFor%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("%sID", uvn).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"),
				jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.IDf("%sIDKey", typ.Name.UnexportedVarName()).Assign().Qual("fmt", "Sprintf").Call(
				jen.ID("jsonPluckQuery"),
				jen.Qual(proj.QuerybuildersPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildersPackage(), "AuditLogEntriesTableContextColumn"),
				jen.ID("audit").Dotf("%sAssignmentKey", sn),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual(proj.QuerybuildersPackage(), "AuditLogEntriesTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildersPackage(), "AuditLogEntriesTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Values(jen.IDf("%sIDKey", typ.Name.UnexportedVarName()).MapAssign().IDf("%sID", uvn))).
					Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildersPackage(), "AuditLogEntriesTableName"),
					jen.Qual(proj.QuerybuildersPackage(), "CreatedOnColumn"),
				),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
