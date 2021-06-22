package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("APIClientSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("Sqlite")).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(buildBuildGetBatchOfAPIClientsQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAPIClientByClientIDQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAllAPIClientsCountQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAPIClientsQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAPIClientByDatabaseIDQuery(proj, dbvendor)...)
	code.Add(buildBuildCreateAPIClientQuery(proj, dbvendor)...)
	code.Add(buildBuildUpdateAPIClientQuery(proj, dbvendor)...)
	code.Add(buildBuildArchiveAPIClientQuery(proj, dbvendor)...)
	code.Add(buildBuildGetAuditLogEntriesForAPIClientQuery(proj, dbvendor)...)

	return code
}

func buildBuildGetBatchOfAPIClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetBatchOfAPIClientsQuery returns a query that fetches every API client in the database within a bucketed range."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetBatchOfAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("APIClientsTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("beginID"))).Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("endID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAPIClientByClientIDQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAPIClientByClientIDQuery returns a SQL query which requests a given API client by its database ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAPIClientByClientIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("APIClientsTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("APIClientsTableClientIDColumn"),
				).Op(":").ID("clientID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllAPIClientsCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAllAPIClientsCountQuery returns a SQL query for the number of API clients"),
		jen.Newline(),
		jen.Comment("returns the database, regardless of ownership."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAllAPIClientsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQueryOnly").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
				)).Dot("From").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAPIClientsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAPIClientsQuery returns a SQL query (and arguments) that will retrieve a list of API clients that"),
		jen.Newline(),
		jen.Comment("returns the given filter's criteria (if relevant) and belong to a given account."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				jen.ID("querybuilding").Dot("APIClientsTableName"),
				jen.ID("nil"),
				jen.ID("nil"),
				jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn"),
				jen.ID("querybuilding").Dot("APIClientsTableColumns"),
				jen.ID("userID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAPIClientByDatabaseIDQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAPIClientByDatabaseIDQuery returns a SQL query which requests a given API client by its database ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAPIClientByDatabaseIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("APIClientsTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn"),
				).Op(":").ID("userID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("IDColumn"),
				).Op(":").ID("clientID"), jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("APIClientsTableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildCreateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildCreateAPIClientQuery returns a SQL query (and args) that will create the given APIClient in the database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildCreateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("APIClientCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Columns").Call(
					jen.ID("querybuilding").Dot("ExternalIDColumn"),
					jen.ID("querybuilding").Dot("APIClientsTableNameColumn"),
					jen.ID("querybuilding").Dot("APIClientsTableClientIDColumn"),
					jen.ID("querybuilding").Dot("APIClientsTableSecretKeyColumn"),
					jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn"),
				).Dot("Values").Call(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Name"),
					jen.ID("input").Dot("ClientID"),
					jen.ID("input").Dot("ClientSecret"),
					jen.ID("input").Dot("BelongsToUser"),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildUpdateAPIClientQuery returns a SQL query (and args) that will update a given API client in the database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildUpdateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("APIClient")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("APIClientsTableClientIDColumn"),
					jen.ID("input").Dot("ClientID"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("input").Dot("ID"), jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn").Op(":").ID("input").Dot("BelongsToUser"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildArchiveAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildArchiveAPIClientQuery returns a SQL query (and arguments) that will mark an API client as archived."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildArchiveAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("clientID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"), jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn").Op(":").ID("userID"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesForAPIClientQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildGetAuditLogEntriesForAPIClientQuery constructs a SQL query for fetching audit log entries belong to a user with a given ID."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("BuildGetAuditLogEntriesForAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("apiClientIDKey").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("jsonPluckQuery"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
				jen.ID("querybuilding").Dot("AuditLogEntriesTableContextColumn"),
				jen.ID("audit").Dot("APIClientAssignmentKey"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableColumns").Op("...")).Dot("From").Call(jen.ID("querybuilding").Dot("AuditLogEntriesTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("apiClientIDKey").Op(":").ID("clientID"))).Dot("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("querybuilding").Dot("AuditLogEntriesTableName"),
					jen.ID("querybuilding").Dot("CreatedOnColumn"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}
