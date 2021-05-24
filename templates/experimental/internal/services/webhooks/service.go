package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("counterName").ID("metrics").Dot("CounterName").Op("=").Lit("webhooks").Var().ID("counterDescription").ID("string").Op("=").Lit("the number of webhooks managed by the webhooks service").Var().ID("serviceName").ID("string").Op("=").Lit("webhooks_service"),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("types").Dot("WebhookDataService").Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("service").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("webhookCounter").ID("metrics").Dot("UnitCounter"),
			jen.ID("webhookDataManager").ID("types").Dot("WebhookDataManager"),
			jen.ID("sessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
			jen.ID("webhookIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideWebhooksService builds a new WebhooksService."),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksService").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("webhookDataManager").ID("types").Dot("WebhookDataManager"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("counterProvider").ID("metrics").Dot("UnitCounterProvider"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.ID("types").Dot("WebhookDataService")).Body(
			jen.Return().Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("webhookDataManager").Op(":").ID("webhookDataManager"), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("webhookCounter").Op(":").ID("metrics").Dot("EnsureUnitCounter").Call(
				jen.ID("counterProvider"),
				jen.ID("logger"),
				jen.ID("counterName"),
				jen.ID("counterDescription"),
			), jen.ID("sessionContextDataFetcher").Op(":").ID("authentication").Dot("FetchContextFromRequest"), jen.ID("webhookIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("WebhookIDURIParamKey"),
				jen.Lit("webhook"),
			), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")))),
		jen.Line(),
	)

	return code
}
