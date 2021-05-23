package accounts

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
			jen.ID("counterName").ID("metrics").Dot("CounterName").Op("=").Lit("accounts"),
			jen.ID("counterDescription").ID("string").Op("=").Lit("the number of accounts managed by the accounts service"),
			jen.ID("serviceName").ID("string").Op("=").Lit("accounts_service"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("types").Dot("AccountDataService").Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("SearchIndex").ID("search").Dot("IndexManager").Type().ID("service").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("accountDataManager").ID("types").Dot("AccountDataManager"),
			jen.ID("accountMembershipDataManager").ID("types").Dot("AccountUserMembershipDataManager"),
			jen.ID("accountIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.ID("userIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.ID("sessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
			jen.ID("accountCounter").ID("metrics").Dot("UnitCounter"),
			jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new AccountsService."),
		jen.Line(),
		jen.Func().ID("ProvideService").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("accountDataManager").ID("types").Dot("AccountDataManager"), jen.ID("accountMembershipDataManager").ID("types").Dot("AccountUserMembershipDataManager"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("counterProvider").ID("metrics").Dot("UnitCounterProvider"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.ID("types").Dot("AccountDataService")).Body(
			jen.Return().Op("&").ID("service").Valuesln(
				jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("accountIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("AccountIDURIParamKey"),
					jen.Lit("account"),
				), jen.ID("userIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("UserIDURIParamKey"),
					jen.Lit("user"),
				), jen.ID("sessionContextDataFetcher").Op(":").ID("authentication").Dot("FetchContextFromRequest"), jen.ID("accountDataManager").Op(":").ID("accountDataManager"), jen.ID("accountMembershipDataManager").Op(":").ID("accountMembershipDataManager"), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("accountCounter").Op(":").ID("metrics").Dot("EnsureUnitCounter").Call(
					jen.ID("counterProvider"),
					jen.ID("logger"),
					jen.ID("counterName"),
					jen.ID("counterDescription"),
				), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")))),
		jen.Line(),
	)

	return code
}
