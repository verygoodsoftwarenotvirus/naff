package admin

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("serviceName").Op("=").Lit("auth_service"),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("service").Struct(
			jen.ID("config").Op("*").ID("authentication").Dot("Config"),
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("authenticator").ID("authentication").Dot("Authenticator"),
			jen.ID("userDB").ID("types").Dot("AdminUserDataManager"),
			jen.ID("auditLog").ID("types").Dot("AdminAuditManager"),
			jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
			jen.ID("sessionManager").Op("*").ID("scs").Dot("SessionManager"),
			jen.ID("sessionContextDataFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
			jen.ID("userIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new AuthService."),
		jen.Line(),
		jen.Func().ID("ProvideService").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("cfg").Op("*").ID("authentication").Dot("Config"), jen.ID("authenticator").ID("authentication").Dot("Authenticator"), jen.ID("userDataManager").ID("types").Dot("AdminUserDataManager"), jen.ID("auditLog").ID("types").Dot("AdminAuditManager"), jen.ID("sessionManager").Op("*").ID("scs").Dot("SessionManager"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.ID("types").Dot("AdminService")).Body(
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(
				jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("config").Op(":").ID("cfg"), jen.ID("userDB").Op(":").ID("userDataManager"), jen.ID("auditLog").Op(":").ID("auditLog"), jen.ID("authenticator").Op(":").ID("authenticator"), jen.ID("sessionManager").Op(":").ID("sessionManager"), jen.ID("sessionContextDataFetcher").Op(":").ID("authentication").Dot("FetchContextFromRequest"), jen.ID("userIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("UserIDURIParamKey"),
					jen.Lit("user"),
				), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName"))),
			jen.ID("svc").Dot("sessionManager").Dot("Lifetime").Op("=").ID("cfg").Dot("Cookies").Dot("Lifetime"),
			jen.Return().ID("svc"),
		),
		jen.Line(),
	)

	return code
}
