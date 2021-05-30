package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("serviceName").Op("=").Lit("auth_service"),
			jen.ID("userIDContextKey").Op("=").ID("string").Call(jen.ID("types").Dot("UserIDContextKey")),
			jen.ID("accountIDContextKey").Op("=").ID("string").Call(jen.ID("types").Dot("AccountIDContextKey")),
			jen.ID("cookieErrorLogName").Op("=").Lit("_COOKIE_CONSTRUCTION_ERROR_"),
			jen.ID("cookieSecretSize").Op("=").Lit(64),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("cookieEncoderDecoder").Interface(
				jen.ID("Encode").Params(jen.ID("name").ID("string"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")),
				jen.ID("Decode").Params(jen.List(jen.ID("name"), jen.ID("value")).ID("string"), jen.ID("dst").Interface()).Params(jen.ID("error")),
			),
			jen.ID("service").Struct(
				jen.ID("config").Op("*").ID("Config"),
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("authenticator").ID("authentication").Dot("Authenticator"),
				jen.ID("userDataManager").ID("types").Dot("UserDataManager"),
				jen.ID("auditLog").ID("types").Dot("AuthAuditManager"),
				jen.ID("apiClientManager").ID("types").Dot("APIClientDataManager"),
				jen.ID("accountMembershipManager").ID("types").Dot("AccountUserMembershipDataManager"),
				jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
				jen.ID("cookieManager").ID("cookieEncoderDecoder"),
				jen.ID("sessionManager").ID("sessionManager"),
				jen.ID("sessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new AuthService."),
		jen.Line(),
		jen.Func().ID("ProvideService").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("cfg").Op("*").ID("Config"), jen.ID("authenticator").ID("authentication").Dot("Authenticator"), jen.ID("userDataManager").ID("types").Dot("UserDataManager"), jen.ID("auditLog").ID("types").Dot("AuthAuditManager"), jen.ID("apiClientsService").ID("types").Dot("APIClientDataManager"), jen.ID("accountMembershipManager").ID("types").Dot("AccountUserMembershipDataManager"), jen.ID("sessionManager").Op("*").ID("scs").Dot("SessionManager"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.ID("types").Dot("AuthService"), jen.ID("error")).Body(
			jen.ID("hashKey").Op(":=").Index().ID("byte").Call(jen.ID("cfg").Dot("Cookies").Dot("HashKey")),
			jen.If(jen.ID("len").Call(jen.ID("hashKey")).Op("==").Lit(0)).Body(
				jen.ID("hashKey").Op("=").ID("securecookie").Dot("GenerateRandomKey").Call(jen.ID("cookieSecretSize"))),
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("config").Op(":").ID("cfg"), jen.ID("userDataManager").Op(":").ID("userDataManager"), jen.ID("auditLog").Op(":").ID("auditLog"), jen.ID("apiClientManager").Op(":").ID("apiClientsService"), jen.ID("accountMembershipManager").Op(":").ID("accountMembershipManager"), jen.ID("authenticator").Op(":").ID("authenticator"), jen.ID("sessionManager").Op(":").ID("sessionManager"), jen.ID("sessionContextDataFetcher").Op(":").ID("FetchContextFromRequest"), jen.ID("cookieManager").Op(":").ID("securecookie").Dot("New").Call(
				jen.ID("hashKey"),
				jen.Index().ID("byte").Call(jen.ID("cfg").Dot("Cookies").Dot("SigningKey")),
			), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName"))),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("svc").Dot("cookieManager").Dot("Encode").Call(
				jen.ID("cfg").Dot("Cookies").Dot("Name"),
				jen.Lit("blah"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("cookie_signing_key_length"),
					jen.ID("len").Call(jen.ID("cfg").Dot("Cookies").Dot("SigningKey")),
				).Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("building test cookie"),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building test cookie: %w"),
					jen.ID("err"),
				)),
			),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
