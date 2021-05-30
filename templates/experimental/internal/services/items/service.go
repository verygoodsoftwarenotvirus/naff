package items

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
			jen.ID("counterName").ID("metrics").Dot("CounterName").Op("=").Lit("items"),
			jen.ID("counterDescription").ID("string").Op("=").Lit("the number of items managed by the items service"),
			jen.ID("serviceName").ID("string").Op("=").Lit("items_service"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("ItemDataService").Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("SearchIndex").ID("search").Dot("IndexManager"),
			jen.ID("Config").Struct(
				jen.ID("Logger").ID("logging").Dot("Config"),
				jen.ID("SearchIndexPath").ID("string"),
			),
			jen.ID("service").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("itemDataManager").ID("types").Dot("ItemDataManager"),
				jen.ID("itemIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
				jen.ID("sessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
				jen.ID("itemCounter").ID("metrics").Dot("UnitCounter"),
				jen.ID("encoderDecoder").ID("encoding").Dot("ServerEncoderDecoder"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("search").ID("SearchIndex"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a Config struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("SearchIndexPath"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new ItemsService."),
		jen.Line(),
		jen.Func().ID("ProvideService").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("cfg").ID("Config"), jen.ID("itemDataManager").ID("types").Dot("ItemDataManager"), jen.ID("encoder").ID("encoding").Dot("ServerEncoderDecoder"), jen.ID("counterProvider").ID("metrics").Dot("UnitCounterProvider"), jen.ID("indexProvider").ID("search").Dot("IndexManagerProvider"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.ID("types").Dot("ItemDataService"), jen.ID("error")).Body(
			jen.ID("logger").Dot("WithValue").Call(
				jen.Lit("index_path"),
				jen.ID("cfg").Dot("SearchIndexPath"),
			).Dot("Debug").Call(jen.Lit("setting up items search index")),
			jen.List(jen.ID("searchIndexManager"), jen.ID("indexInitErr")).Op(":=").ID("indexProvider").Call(
				jen.ID("search").Dot("IndexPath").Call(jen.ID("cfg").Dot("SearchIndexPath")),
				jen.ID("types").Dot("ItemsSearchIndexName"),
				jen.ID("logger"),
			),
			jen.If(jen.ID("indexInitErr").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("indexInitErr"),
					jen.Lit("setting up items search index"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("indexInitErr")),
			),
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("itemIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("ItemIDURIParamKey"),
				jen.Lit("item"),
			), jen.ID("sessionContextDataFetcher").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/authentication", "FetchContextFromRequest"), jen.ID("itemDataManager").Op(":").ID("itemDataManager"), jen.ID("encoderDecoder").Op(":").ID("encoder"), jen.ID("itemCounter").Op(":").ID("metrics").Dot("EnsureUnitCounter").Call(
				jen.ID("counterProvider"),
				jen.ID("logger"),
				jen.ID("counterName"),
				jen.ID("counterDescription"),
			), jen.ID("search").Op(":").ID("searchIndexManager"), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName"))),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
