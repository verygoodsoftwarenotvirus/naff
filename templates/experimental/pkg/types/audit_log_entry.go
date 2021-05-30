package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntryDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("FieldChangeSummary").Struct(
				jen.ID("OldValue").Interface(),
				jen.ID("NewValue").Interface(),
				jen.ID("FieldName").ID("string"),
			),
			jen.ID("AuditLogContext").Map(jen.ID("string")).Interface(),
			jen.ID("AuditLogEntry").Struct(
				jen.ID("Context").ID("AuditLogContext"),
				jen.ID("ExternalID").ID("string"),
				jen.ID("EventType").ID("string"),
				jen.ID("ID").ID("uint64"),
				jen.ID("CreatedOn").ID("uint64"),
			),
			jen.ID("AuditLogEntryList").Struct(
				jen.ID("Entries").Index().Op("*").ID("AuditLogEntry"),
				jen.ID("Pagination"),
			),
			jen.ID("AuditLogEntryCreationInput").Struct(
				jen.ID("Context").ID("AuditLogContext"),
				jen.ID("EventType").ID("string"),
			),
			jen.ID("AuditLogEntryDataManager").Interface(
				jen.ID("GetAuditLogEntry").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("eventID").ID("uint64")).Params(jen.Op("*").ID("AuditLogEntry"), jen.ID("error")),
				jen.ID("GetAllAuditLogEntriesCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetAllAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("resultChannel").Chan().Index().Op("*").ID("AuditLogEntry"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")),
				jen.ID("GetAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("AuditLogEntryList"), jen.ID("error")),
			),
			jen.ID("AuditLogEntryDataService").Interface(
				jen.ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Value implements the driver.Valuer interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").ID("AuditLogContext")).ID("Value").Params().Params(jen.ID("driver").Dot("Value"), jen.ID("error")).Body(
			jen.Return().Qual("encoding/json", "Marshal").Call(jen.ID("c"))),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errByteAssertionFailed").Op("=").Qual("errors", "New").Call(jen.Lit("type assertion to []byte failed")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Scan implements the sql.Scanner interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("AuditLogContext")).ID("Scan").Params(jen.ID("value").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("b"), jen.ID("ok")).Op(":=").ID("value").Assert(jen.Index().ID("byte")),
			jen.If(jen.Op("!").ID("ok")).Body(
				jen.Return().ID("errByteAssertionFailed")),
			jen.Return().Qual("encoding/json", "Unmarshal").Call(
				jen.ID("b"),
				jen.Op("&").ID("c"),
			),
		),
		jen.Line(),
	)

	return code
}
