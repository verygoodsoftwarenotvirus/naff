package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntriesDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.QuerybuildingPackage(), "AuditLogEntrySQLQueryBuilder").Equals().Parens(jen.PointerTo().ID(dbvendor.Singular())).Call(jen.Nil()),
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
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAuditLogEntryIDToSpan").Call(
				jen.ID("span"),
				jen.ID("entryID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).MapAssign().ID("entryID"))),
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
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAllAuditLogEntriesCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.String()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")),
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
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetBatchOfAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).MapAssign().ID("beginID"))).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).MapAssign().ID("endID"))),
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
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildGetAuditLogEntriesQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("filter").DoesNotEqual().Nil()).Body(
				jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.String().Call(jen.ID("filter").Dot("SortBy")),
				),
			),
			jen.Newline(),
			jen.ID("countQueryBuilder").Assign().ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("allCountQuery")).
				Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")),
			jen.Newline(),
			jen.List(jen.ID("countQuery"), jen.ID("countQueryArgs"), jen.ID("err")).Assign().ID("countQueryBuilder").Dot("ToSql").Call(),
			jen.ID("b").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.Newline(),
			jen.ID("builder").Assign().ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("append").Call(
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s)"),
					jen.ID("countQuery"),
				),
			).Op("...")).
				Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
				Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
			)),
			jen.Newline(),
			jen.If(jen.ID("filter").DoesNotEqual().Nil()).Body(
				jen.ID("builder").Equals().Qual(proj.QuerybuildingPackage(), "ApplyFilterToQueryBuilder").Call(
					jen.ID("filter"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
					jen.ID("builder"),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("selectArgs"), jen.ID("err")).Assign().ID("builder").Dot("ToSql").Call(),
			jen.ID("b").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.Newline(),
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
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).ID("BuildCreateAuditLogEntryQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAuditLogEntryEventTypeToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("EventType"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableEventTypeColumn"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("EventType"),
					jen.ID("input").Dot("Context"),
				).Add(utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildingPackage(), "IDColumn"))))),
			),
		),
		jen.Newline(),
	}

	return lines
}
