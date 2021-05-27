package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntriesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AuditLogEntryDataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAuditLogEntry takes a database Scanner (i.e. *sql.Row) and scans the result into an AuditLogEntry struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAuditLogEntry").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").ID("database").Dot("Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("entry").Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("totalCount").ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.ID("entry").Op("=").Op("&").ID("types").Dot("AuditLogEntry").Values(),
			jen.ID("targetVars").Op(":=").Index().Interface().Valuesln(jen.Op("&").ID("entry").Dot("ID"), jen.Op("&").ID("entry").Dot("ExternalID"), jen.Op("&").ID("entry").Dot("EventType"), jen.Op("&").ID("entry").Dot("Context"), jen.Op("&").ID("entry").Dot("CreatedOn")),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Op("=").ID("append").Call(
					jen.ID("targetVars"),
					jen.Op("&").ID("totalCount"),
				)),
			jen.If(jen.ID("err").Op("=").ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning API client database result"),
				))),
			jen.Return().List(jen.ID("entry"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAuditLogEntries takes some database rows and turns them into a slice of AuditLogEntry pointers."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("entries").Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("totalCount").ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("x"), jen.ID("tc"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanAuditLogEntry").Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("scanErr"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("scanning audit log entries"),
					))),
				jen.If(jen.ID("includeCounts")).Body(
					jen.If(jen.ID("totalCount").Op("==").Lit(0)).Body(
						jen.ID("totalCount").Op("=").ID("tc"))),
				jen.ID("entries").Op("=").ID("append").Call(
					jen.ID("entries"),
					jen.ID("x"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("checkRowsForErrorAndClose").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("handling rows"),
				))),
			jen.Return().List(jen.ID("entries"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntry fetches an audit log entry from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAuditLogEntry").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("entryID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("tracing").Dot("AttachAuditLogEntryIDToSpan").Call(
				jen.ID("span"),
				jen.ID("entryID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AuditLogEntryIDKey"),
				jen.ID("entryID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAuditLogEntryQuery").Call(
				jen.ID("ctx"),
				jen.ID("entryID"),
			),
			jen.ID("row").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("audit log entry"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.List(jen.ID("entry"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning audit log entry"),
				))),
			jen.Return().List(jen.ID("entry"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAuditLogEntriesCount fetches the count of audit log entries from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllAuditLogEntriesCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAllAuditLogEntriesCountQuery").Call(jen.ID("ctx")),
				jen.Lit("fetching count of audit logs entries"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for count of audit log entries"),
				))),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAuditLogEntries fetches a list of all audit log entries in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("batchSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("results").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("batch_size"),
				jen.ID("batchSize"),
			),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("GetAllAuditLogEntriesCount").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching count of entries"),
				)),
			jen.For(jen.ID("beginID").Op(":=").ID("uint64").Call(jen.Lit(1)), jen.ID("beginID").Op("<=").ID("count"), jen.ID("beginID").Op("+=").ID("uint64").Call(jen.ID("batchSize"))).Body(
				jen.ID("endID").Op(":=").ID("beginID").Op("+").ID("uint64").Call(jen.ID("batchSize")),
				jen.Go().Func().Params(jen.List(jen.ID("begin"), jen.ID("end")).ID("uint64")).Body(
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetBatchOfAuditLogEntriesQuery").Call(
						jen.ID("ctx"),
						jen.ID("begin"),
						jen.ID("end"),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("query").Op(":").ID("query"), jen.Lit("begin").Op(":").ID("begin"), jen.Lit("end").Op(":").ID("end"))),
					jen.List(jen.ID("rows"), jen.ID("queryErr")).Op(":=").ID("q").Dot("db").Dot("Query").Call(
						jen.ID("query"),
						jen.ID("args").Op("..."),
					),
					jen.If(jen.Qual("errors", "Is").Call(
						jen.ID("queryErr"),
						jen.Qual("database/sql", "ErrNoRows"),
					)).Body(
						jen.Return()).Else().If(jen.ID("queryErr").Op("!=").ID("nil")).Body(
						jen.ID("logger").Dot("Error").Call(
							jen.ID("queryErr"),
							jen.Lit("querying for database rows"),
						),
						jen.Return(),
					),
					jen.List(jen.ID("auditLogEntries"), jen.ID("_"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanAuditLogEntries").Call(
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
					jen.ID("results").ReceiveFromChannel().ID("auditLogEntries"),
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
		jen.Comment("GetAuditLogEntries fetches a list of audit log entries from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("x").Op("*").ID("types").Dot("AuditLogEntryList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("q").Dot("logger")),
			jen.ID("x").Op("=").Op("&").ID("types").Dot("AuditLogEntryList").Values(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Op("=").List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAuditLogEntriesQuery").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("audit log entries"),
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
			jen.If(jen.List(jen.ID("x").Dot("Entries"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Op("=").ID("q").Dot("scanAuditLogEntries").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning audit log entry"),
				))),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("createAuditLogEntryInTransaction creates an audit log entry in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("createAuditLogEntryInTransaction").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("transaction").Op("*").Qual("database/sql", "Tx"), jen.ID("input").Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.If(jen.ID("transaction").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilTransactionProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AuditLogEntryEventTypeKey"),
				jen.ID("input").Dot("EventType"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AuditLogEntryContextKey"),
				jen.ID("input").Dot("Context"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildCreateAuditLogEntryQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.ID("tracing").Dot("AttachAuditLogEntryEventTypeToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("EventType"),
			),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("transaction"),
				jen.Lit("audit log entry creation"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("executing audit log entry creation query"),
				),
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("transaction"),
				),
				jen.Return().ID("err"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("audit log entry created")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("createAuditLogEntry creates an audit log entry in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("createAuditLogEntry").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("querier").ID("database").Dot("Querier"), jen.ID("input").Op("*").ID("types").Dot("AuditLogEntryCreationInput")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.ID("observability").Dot("NoteEvent").Call(
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("early return due to nil input"),
				),
				jen.Return(),
			),
			jen.If(jen.ID("querier").Op("==").ID("nil")).Body(
				jen.ID("observability").Dot("NoteEvent").Call(
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("early return due to nil querier"),
				),
				jen.Return(),
			),
			jen.ID("tracing").Dot("AttachAuditLogEntryEventTypeToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("EventType"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AuditLogEntryEventTypeKey"),
				jen.ID("input").Dot("EventType"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildCreateAuditLogEntryQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.List(jen.ID("id"), jen.ID("err")).Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("querier"),
				jen.ID("true"),
				jen.Lit("audit log entry creation"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("executing audit log entry creation query"),
				)),
			jen.ID("tracing").Dot("AttachAuditLogEntryIDToSpan").Call(
				jen.ID("span"),
				jen.ID("id"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("audit log entry created")),
		),
		jen.Line(),
	)

	return code
}
