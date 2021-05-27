package requests

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
			jen.ID("auditLogBasePath").Op("=").Lit("audit_log"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntryRequest builds an HTTP request for fetching a given audit log entry."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetAuditLogEntryRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("entryID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AuditLogEntryIDKey"),
				jen.ID("entryID"),
			),
			jen.ID("tracing").Dot("AttachAuditLogEntryIDToSpan").Call(
				jen.ID("span"),
				jen.ID("entryID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("adminBasePath"),
				jen.ID("auditLogBasePath"),
				jen.ID("id").Call(jen.ID("entryID")),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogEntriesRequest builds an HTTP request for fetching audit log entries."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetAuditLogEntriesRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("b").Dot("logger")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("adminBasePath"),
				jen.ID("auditLogBasePath"),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
