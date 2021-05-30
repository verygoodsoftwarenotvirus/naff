package httpclient

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("GetAuditLogEntry retrieves an entry."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAuditLogEntry").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("entryID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AuditLogEntryIDKey"),
				jen.ID("entryID"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetAuditLogEntryRequest").Call(
				jen.ID("ctx"),
				jen.ID("entryID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building get audit log entry request"),
				))),
			jen.ID("c").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("Debug").Call(jen.Lit("Fetching audit log entry")),
			jen.Var().Defs(
				jen.ID("entry").Op("*").ID("types").Dot("AuditLogEntry"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("entry"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving audit log entry"),
				))),
			jen.Return().List(jen.ID("entry"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntries retrieves a list of audit log entries."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAuditLogEntries").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("AuditLogEntryList"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("loggerWithFilter").Call(jen.ID("filter")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetAuditLogEntriesRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building fetch audit log entries request"),
				))),
			jen.ID("logger").Op("=").ID("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Var().Defs(
				jen.ID("entries").Op("*").ID("types").Dot("AuditLogEntryList"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("entries"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching audit log entries"),
				))),
			jen.Return().List(jen.ID("entries"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
