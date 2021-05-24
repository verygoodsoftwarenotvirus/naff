package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("serviceName").Op("=").Lit("audit_log_entries_service"),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("types").Dot("AuditLogEntryDataService").Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("service").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("auditLog").ID("types").Dot("AuditLogEntryDataManager"),
			jen.ID("auditLogEntryIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.ID("sessionContextDataFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
			jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new service."),
		jen.Line(),
		jen.Func().ID("ProvideService").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("auditLog").ID("types").Dot("AuditLogEntryDataManager"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.ID("types").Dot("AuditLogEntryDataService")).Body(
			jen.Return().Op("&").ID("service").Valuesln(
				jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("auditLog").Op(":").ID("auditLog"), jen.ID("auditLogEntryIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("LogEntryURIParamKey"),
					jen.Lit("audit log entry"),
				), jen.ID("sessionContextDataFetcher").Op(":").ID("authentication").Dot("FetchContextFromRequest"), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")))),
		jen.Line(),
	)

	return code
}
