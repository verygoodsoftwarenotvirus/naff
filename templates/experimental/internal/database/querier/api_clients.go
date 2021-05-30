package querier

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
			jen.ID("_").ID("types").Dot("APIClientDataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAPIClient takes a Scanner (i.e. *sql.Row) and scans its results into an APIClient struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").ID("database").Dot("Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("client").Op("*").ID("types").Dot("APIClient"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.ID("client").Op("=").Op("&").ID("types").Dot("APIClient").Valuesln(),
			jen.ID("targetVars").Op(":=").Index().Interface().Valuesln(jen.Op("&").ID("client").Dot("ID"), jen.Op("&").ID("client").Dot("ExternalID"), jen.Op("&").ID("client").Dot("Name"), jen.Op("&").ID("client").Dot("ClientID"), jen.Op("&").ID("client").Dot("ClientSecret"), jen.Op("&").ID("client").Dot("CreatedOn"), jen.Op("&").ID("client").Dot("LastUpdatedOn"), jen.Op("&").ID("client").Dot("ArchivedOn"), jen.Op("&").ID("client").Dot("BelongsToUser")),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Op("=").ID("append").Call(
					jen.ID("targetVars"),
					jen.Op("&").ID("filteredCount"),
					jen.Op("&").ID("totalCount"),
				)),
			jen.If(jen.ID("err").Op("=").ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning API client database result"),
				))),
			jen.Return().List(jen.ID("client"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAPIClients takes sql rows and turns them into a slice of API Clients."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAPIClients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("clients").Index().Op("*").ID("types").Dot("APIClient"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("client"), jen.ID("fc"), jen.ID("tc"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanAPIClient").Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("scanErr"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("scanning API client"),
					))),
				jen.If(jen.ID("includeCounts")).Body(
					jen.If(jen.ID("filteredCount").Op("==").Lit(0)).Body(
						jen.ID("filteredCount").Op("=").ID("fc")),
					jen.If(jen.ID("totalCount").Op("==").Lit(0)).Body(
						jen.ID("totalCount").Op("=").ID("tc")),
				),
				jen.ID("clients").Op("=").ID("append").Call(
					jen.ID("clients"),
					jen.ID("client"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("checkRowsForErrorAndClose").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("handling rows"),
				))),
			jen.Return().List(jen.ID("clients"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAPIClientByClientID gets an API client from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAPIClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.Op("*").ID("types").Dot("APIClient"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit("")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrEmptyInputProvided"))),
			jen.ID("tracing").Dot("AttachAPIClientClientIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientClientIDKey"),
				jen.ID("clientID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAPIClientByClientIDQuery").Call(
				jen.ID("ctx"),
				jen.ID("clientID"),
			),
			jen.ID("row").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("API client"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.List(jen.ID("client"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAPIClient").Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("err"))),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for API client"),
				)),
			),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAPIClientByDatabaseID gets an API client from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAPIClientByDatabaseID").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("types").Dot("APIClient"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit(0).Op("||").ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("APIClientDatabaseIDKey").Op(":").ID("clientID"), jen.ID("keys").Dot("UserIDKey").Op(":").ID("userID"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAPIClientByDatabaseIDQuery").Call(
				jen.ID("ctx"),
				jen.ID("clientID"),
				jen.ID("userID"),
			),
			jen.ID("row").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("API client"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.List(jen.ID("client"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAPIClient").Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("err"))),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for API client"),
				)),
			),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetTotalAPIClientCount gets the count of API clients that match the current filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetTotalAPIClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAllAPIClientsCountQuery").Call(jen.ID("ctx")),
				jen.Lit("fetching count of API clients"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for count of API clients"),
				))),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAPIClients loads all API clients into a channel."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllAPIClients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").ID("types").Dot("APIClient"), jen.ID("batchSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("results").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("batch_size"),
				jen.ID("batchSize"),
			),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("GetTotalAPIClientCount").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching count of API clients"),
				)),
			jen.For(jen.ID("beginID").Op(":=").ID("uint64").Call(jen.Lit(1)), jen.ID("beginID").Op("<=").ID("count"), jen.ID("beginID").Op("+=").ID("uint64").Call(jen.ID("batchSize"))).Body(
				jen.ID("endID").Op(":=").ID("beginID").Op("+").ID("uint64").Call(jen.ID("batchSize")),
				jen.Go().Func().Params(jen.List(jen.ID("begin"), jen.ID("end")).ID("uint64")).Body(
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetBatchOfAPIClientsQuery").Call(
						jen.ID("ctx"),
						jen.ID("begin"),
						jen.ID("end"),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("query").Op(":").ID("query"), jen.Lit("begin").Op(":").ID("begin"), jen.Lit("end").Op(":").ID("end"))),
					jen.List(jen.ID("rows"), jen.ID("queryErr")).Op(":=").ID("q").Dot("db").Dot("Query").Call(
						jen.ID("query"),
						jen.ID("args").Op("..."),
					),
					jen.If(jen.ID("queryErr").Op("!=").ID("nil")).Body(
						jen.If(jen.Op("!").Qual("errors", "Is").Call(
							jen.ID("queryErr"),
							jen.Qual("database/sql", "ErrNoRows"),
						)).Body(
							jen.ID("logger").Dot("Error").Call(
								jen.ID("queryErr"),
								jen.Lit("querying for database rows"),
							)),
						jen.Return(),
					),
					jen.List(jen.ID("clients"), jen.ID("_"), jen.ID("_"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanAPIClients").Call(
						jen.ID("ctx"),
						jen.ID("rows"),
						jen.ID("false"),
					),
					jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
						jen.ID("logger").Dot("Error").Call(
							jen.ID("scanErr"),
							jen.Lit("scanning database rows"),
						),
						jen.Return(),
					),
					jen.ID("results").ReceiveFromChannel().ID("clients"),
				).Call(
					jen.ID("beginID"),
					jen.ID("endID"),
				),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAPIClients gets a list of API clients."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAPIClients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("x").Op("*").ID("types").Dot("APIClientList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("q").Dot("logger")).Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("x").Op("=").Op("&").ID("types").Dot("APIClientList").Valuesln(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Op("=").List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAPIClientsQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("API clients"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("err"))),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for API clients"),
				)),
			),
			jen.If(jen.List(jen.ID("x").Dot("Clients"), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Op("=").ID("q").Dot("scanAPIClients").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning response from database"),
				))),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateAPIClient creates an API client."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("CreateAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("APIClientCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("types").Dot("APIClient"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("createdByUser").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("createdByUser"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("APIClientClientIDKey").Op(":").ID("input").Dot("ClientID"), jen.ID("keys").Dot("UserIDKey").Op(":").ID("input").Dot("BelongsToUser"))),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildCreateAPIClientQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.List(jen.ID("id"), jen.ID("err")).Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("false"),
				jen.Lit("API client creation"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating API client"),
				)),
			),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("id"),
			),
			jen.ID("client").Op(":=").Op("&").ID("types").Dot("APIClient").Valuesln(jen.ID("ID").Op(":").ID("id"), jen.ID("Name").Op(":").ID("input").Dot("Name"), jen.ID("ClientID").Op(":").ID("input").Dot("ClientID"), jen.ID("ClientSecret").Op(":").ID("input").Dot("ClientSecret"), jen.ID("BelongsToUser").Op(":").ID("input").Dot("BelongsToUser"), jen.ID("CreatedOn").Op(":").ID("q").Dot("currentTime").Call()),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildAPIClientCreationEventEntry").Call(
					jen.ID("client"),
					jen.ID("createdByUser"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing API client creation audit log entry"),
				)),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				))),
			jen.ID("logger").Dot("Info").Call(jen.Lit("API client created")),
			jen.Return().List(jen.ID("client"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveAPIClient archives an API client."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("ArchiveAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("accountID"), jen.ID("archivedByUser")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit(0).Op("||").ID("accountID").Op("==").Lit(0).Op("||").ID("archivedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("archivedByUser"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("APIClientDatabaseIDKey").Op(":").ID("clientID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("UserIDKey").Op(":").ID("archivedByUser"))),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				)),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildArchiveAPIClientQuery").Call(
				jen.ID("ctx"),
				jen.ID("clientID"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("API client archive"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating API client"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildAPIClientArchiveEventEntry").Call(
					jen.ID("accountID"),
					jen.ID("clientID"),
					jen.ID("archivedByUser"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing API client archive audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("API client archived")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntriesForAPIClient fetches a list of audit log entries from the database that relate to a given client."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAuditLogEntriesForAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("clientID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("APIClientDatabaseIDKey"),
				jen.ID("clientID"),
			),
			jen.ID("tracing").Dot("AttachAPIClientDatabaseIDToSpan").Call(
				jen.ID("span"),
				jen.ID("clientID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAuditLogEntriesForAPIClientQuery").Call(
				jen.ID("ctx"),
				jen.ID("clientID"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("audit log entries for API client"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying database for audit log entries"),
				))),
			jen.List(jen.ID("auditLogEntries"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAuditLogEntries").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning response from database"),
				))),
			jen.Return().List(jen.ID("auditLogEntries"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
