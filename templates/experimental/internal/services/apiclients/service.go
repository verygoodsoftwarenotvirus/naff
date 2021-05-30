package apiclients

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
			jen.ID("counterName").ID("metrics").Dot("CounterName").Op("=").Lit("api_clients"),
			jen.ID("counterDescription").ID("string").Op("=").Lit("number of API clients managed by the API client service"),
			jen.ID("serviceName").ID("string").Op("=").Lit("api_clients_service"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("APIClientDataService").Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("config").Struct(jen.List(jen.ID("minimumUsernameLength"), jen.ID("minimumPasswordLength")).ID("uint8")),
			jen.ID("service").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("cfg").Op("*").ID("config"),
				jen.ID("apiClientDataManager").ID("types").Dot("APIClientDataManager"),
				jen.ID("userDataManager").ID("types").Dot("UserDataManager"),
				jen.ID("authenticator").ID("authentication").Dot("Authenticator"),
				jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
				jen.ID("urlClientIDExtractor").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
				jen.ID("sessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
				jen.ID("apiClientCounter").ID("metrics").Dot("UnitCounter"),
				jen.ID("secretGenerator").ID("random").Dot("Generator"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAPIClientsService builds a new APIClientsService."),
		jen.Line(),
		jen.Func().ID("ProvideAPIClientsService").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("clientDataManager").ID("types").Dot("APIClientDataManager"), jen.ID("userDataManager").ID("types").Dot("UserDataManager"), jen.ID("authenticator").ID("authentication").Dot("Authenticator"), jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("counterProvider").ID("metrics").Dot("UnitCounterProvider"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager"), jen.ID("cfg").Op("*").ID("config")).Params(jen.ID("types").Dot("APIClientDataService")).Body(
			jen.Return().Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("cfg").Op(":").ID("cfg"), jen.ID("apiClientDataManager").Op(":").ID("clientDataManager"), jen.ID("userDataManager").Op(":").ID("userDataManager"), jen.ID("authenticator").Op(":").ID("authenticator"), jen.ID("encoderDecoder").Op(":").ID("encoderDecoder"), jen.ID("urlClientIDExtractor").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("APIClientIDURIParamKey"),
				jen.Lit("api client"),
			), jen.ID("sessionContextDataFetcher").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/authentication", "FetchContextFromRequest"), jen.ID("apiClientCounter").Op(":").ID("metrics").Dot("EnsureUnitCounter").Call(
				jen.ID("counterProvider"),
				jen.ID("logger"),
				jen.ID("counterName"),
				jen.ID("counterDescription"),
			), jen.ID("secretGenerator").Op(":").ID("random").Dot("NewGenerator").Call(jen.ID("logger")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")))),
		jen.Line(),
	)

	return code
}
