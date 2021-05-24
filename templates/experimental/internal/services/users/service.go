package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("serviceName").Op("=").Lit("users_service").Var().ID("counterDescription").Op("=").Lit("number of users managed by the users service").Var().ID("counterName").Op("=").ID("metrics").Dot("CounterName").Call(jen.Lit("users")),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("types").Dot("UserDataService").Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("RequestValidator").Interface(jen.ID("Validate").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("bool"), jen.ID("error"))).Type().ID("service").Struct(
			jen.ID("userDataManager").ID("types").Dot("UserDataManager"),
			jen.ID("accountDataManager").ID("types").Dot("AccountDataManager"),
			jen.ID("authSettings").Op("*").ID("authentication").Dot("Config"),
			jen.ID("authenticator").ID("authentication").Dot("Authenticator"),
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
			jen.ID("userIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.ID("sessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
			jen.ID("userCounter").ID("metrics").Dot("UnitCounter"),
			jen.ID("secretGenerator").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/random", "Generator"),
			jen.ID("imageUploadProcessor").ID("images").Dot("ImageUploadProcessor"),
			jen.ID("uploadManager").ID("uploads").Dot("UploadManager"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUsersService builds a new UsersService."),
		jen.Line(),
		jen.Func().ID("ProvideUsersService").Params(jen.ID("authSettings").Op("*").ID("authentication").Dot("Config"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("userDataManager").ID("types").Dot("UserDataManager"), jen.ID("accountDataManager").ID("types").Dot("AccountDataManager"), jen.ID("authenticator").ID("authentication").Dot("Authenticator"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("counterProvider").ID("metrics").Dot("UnitCounterProvider"), jen.ID("imageUploadProcessor").ID("images").Dot("ImageUploadProcessor"), jen.ID("uploadManager").ID("uploads").Dot("UploadManager"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.ID("types").Dot("UserDataService")).Body(
			jen.Return().Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("userDataManager").Op(":").ID("userDataManager"), jen.ID("accountDataManager").Op(":").ID("accountDataManager"), jen.ID("authenticator").Op(":").ID("authenticator"), jen.ID("userIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("UserIDURIParamKey"),
				jen.Lit("user"),
			), jen.ID("sessionContextDataFetcher").Op(":").ID("authentication").Dot("FetchContextFromRequest"), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("authSettings").Op(":").ID("authSettings"), jen.ID("userCounter").Op(":").ID("metrics").Dot("EnsureUnitCounter").Call(
				jen.ID("counterProvider"),
				jen.ID("logger"),
				jen.ID("counterName"),
				jen.ID("counterDescription"),
			), jen.ID("secretGenerator").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/random", "NewGenerator").Call(jen.ID("logger")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")), jen.ID("imageUploadProcessor").Op(":").ID("imageUploadProcessor"), jen.ID("uploadManager").Op(":").ID("uploadManager"))),
		jen.Line(),
	)

	return code
}
