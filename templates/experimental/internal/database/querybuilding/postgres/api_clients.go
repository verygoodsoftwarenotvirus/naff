package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("querybuilding").Dot("APIClientSQLQueryBuilder").Op("=").Parens(jen.Op("*").ID("Postgres")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetBatchOfAPIClientsQuery returns a query that fetches every item in the database within a bucketed range."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetBatchOfAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAPIClientByClientIDQuery returns a SQL query which requests a given API client by its database ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAPIClientByClientIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachAPIClientClientIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
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
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAllAPIClientsCountQuery returns a SQL query for the number of API clients"),
		jen.Line(),
		jen.Func().Comment("in the database regardless of ownership.").Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAllAPIClientsCountQuery").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
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
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAPIClientsQuery returns a SQL query (and arguments) that will retrieve a list of API clients that"),
		jen.Line(),
		jen.Func().Comment("meet the given client's criteria (if relevant) and belong to a given account.").Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAPIClientsQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
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
				jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn"),
				jen.ID("querybuilding").Dot("APIClientsTableColumns"),
				jen.ID("userID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAPIClientByDatabaseIDQuery returns a SQL query which requests a given API client by its database ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAPIClientByDatabaseIDQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
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
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateAPIClientQuery returns a SQL query (and args) that will create the given APIClient in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildCreateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("APIClientCreationInput")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
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
				).Dot("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("RETURNING %s"),
					jen.ID("querybuilding").Dot("IDColumn"),
				)),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUpdateAPIClientQuery returns a SQL query (and args) that will update a given API client in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildUpdateAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("APIClient")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachAPIClientClientIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ClientID"),
			),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("APIClientsTableClientIDColumn"),
					jen.ID("input").Dot("ClientID"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("input").Dot("ID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"), jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn").Op(":").ID("input").Dot("BelongsToUser"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveAPIClientQuery returns a SQL query (and arguments) that will mark an API client as archived."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildArchiveAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.ID("querybuilding").Dot("APIClientsTableName")).Dot("Set").Call(
					jen.ID("querybuilding").Dot("LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Set").Call(
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.ID("querybuilding").Dot("IDColumn").Op(":").ID("clientID"), jen.ID("querybuilding").Dot("APIClientsTableOwnershipColumn").Op(":").ID("userID"), jen.ID("querybuilding").Dot("ArchivedOnColumn").Op(":").ID("nil"))),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesForAPIClientQuery constructs a SQL query for fetching an audit log entry with a given ID belong to a user with a given ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Postgres")).ID("BuildGetAuditLogEntriesForAPIClientQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
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
		jen.Line(),
	)

	return code
}
