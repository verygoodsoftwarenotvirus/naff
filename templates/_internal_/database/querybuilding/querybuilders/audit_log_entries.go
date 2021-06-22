package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntriesDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("AuditLogEntrySQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("Sqlite")).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(buildBuildGetAuditLogEntryQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAllAuditLogEntriesCountQuery(proj, dbvendor)...)
	code.Add(buildBuildGetBatchOfAuditLogEntriesQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAuditLogEntriesQuery(proj, dbvendor)...)
	code.Add(buildBuildCreateAuditLogEntryQuery(proj, dbvendor)...)

	return code
}

func buildBuildGetAuditLogEntryQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAuditLogEntryQuery constructs a SQL query for fetching an audit log entry with a given ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("entryID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllAuditLogEntriesCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAllAuditLogEntriesCountQuery returns a query that fetches the total number of  in the database."),
		jen.Newline(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAllAuditLogEntriesCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
				)).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetBatchOfAuditLogEntriesQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetBatchOfAuditLogEntriesQuery returns a query that fetches every audit log entry in the database within a bucketed range."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetBatchOfAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("beginID"))).Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAuditLogEntriesQuery builds a SQL query selecting  that adhere to a given QueryFilter and belong to a given account,"),
		jen.Newline(),
		jen.Comment("and returns both the query and the relevant args to pass to the query executor."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				)),
			jen.ID("countQueryBuilder").Op(":=").ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("allCountQuery")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")),
			jen.List(jen.ID("countQuery"), jen.ID("countQueryArgs"), jen.ID("err")).Op(":=").ID("countQueryBuilder").Dot("ToSql").Call(),
			jen.ID("b").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.ID("builder").Op(":=").ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("append").Call(
				jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s)"),
					jen.ID("countQuery"),
				),
			).Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
				jen.ID("querybuilding").Dot("CreatedOnColumn"),
			)),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("builder").Op("=").ID("querybuilding").Dot("ApplyFilterToQueryBuilder").Call(
					jen.ID("filter"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("builder"),
				)),
			jen.List(jen.ID("query"), jen.ID("selectArgs"), jen.ID("err")).Op(":=").ID("builder").Dot("ToSql").Call(),
			jen.ID("b").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.Return().List(jen.ID("query"), jen.ID("append").Call(
				jen.ID("countQueryArgs"),
				jen.ID("selectArgs").Op("..."),
			)),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildCreateAuditLogEntryQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildCreateAuditLogEntryQuery takes an audit log entry and returns a creation query for that audit log entry and the relevant arguments."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildCreateAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachAuditLogEntryEventTypeToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("EventType"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("ExternalIDColumn"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableEventTypeColumn"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableContextColumn"),
				).Dot("Values").Call(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("EventType"),
					jen.ID("input").Dot("Context"),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}
