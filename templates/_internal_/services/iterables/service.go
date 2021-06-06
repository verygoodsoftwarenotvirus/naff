package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	pcn := typ.Name.PluralCommonName()

	code.Add(
		jen.Const().Defs(
			jen.ID("counterName").Qual(proj.MetricsPackage(), "CounterName").Op("=").Lit("items"),
			jen.ID("counterDescription").ID("string").Op("=").Lit("the number of items managed by the items service"),
			jen.ID("serviceName").ID("string").Op("=").Lit("items_service"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().ID("_").Qual(proj.TypesPackage(), "ItemDataService").Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(
		jen.Type().Defs(
			jen.Comment("SearchIndex is a type alias for dependency injection's sake."),
			jen.ID("SearchIndex").Qual(proj.InternalSearchPackage(), "IndexManager"),
			jen.Newline(),
			jen.Commentf("service handles %s.", pcn),
			jen.ID("service").Struct(
				jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
				jen.ID("itemDataManager").Qual(proj.TypesPackage(), "ItemDataManager"),
				jen.ID("itemIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
				jen.ID("sessionContextDataFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"),
					jen.ID("error")),
				jen.ID("itemCounter").Qual(proj.MetricsPackage(), "UnitCounter"),
				jen.ID("encoderDecoder").Qual(proj.EncodingPackage(), "ServerEncoderDecoder"),
				jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("search").ID("SearchIndex")
					}
					return jen.Null()
				}(),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new ItemsService."),
		jen.Newline(),
		jen.Func().ID("ProvideService").Paramsln(
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("cfg").ID("Config"),
			jen.ID("itemDataManager").Qual(proj.TypesPackage(), "ItemDataManager"),
			jen.ID("encoder").Qual(proj.EncodingPackage(), "ServerEncoderDecoder"),
			jen.ID("counterProvider").Qual(proj.MetricsPackage(), "UnitCounterProvider"),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("indexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider")
				}
				return jen.Null()
			}(),
			jen.ID("routeParamManager").Qual(proj.RoutingPackage(), "RouteParamManager"),
		).Params(jen.Qual(proj.TypesPackage(), "ItemDataService"), jen.ID("error")).Body(

			func() jen.Code {
				if typ.SearchEnabled {
					return jen.List(jen.ID("searchIndexManager"),
						jen.ID("indexInitErr")).Op(":=").ID("indexProvider").Call(
						jen.Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("cfg").Dot("SearchIndexPath")),
						jen.Qual(proj.TypesPackage(), "ItemsSearchIndexName"),
						jen.ID("logger"),
					)
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.If(jen.ID("indexInitErr").Op("!=").ID("nil")).Body(
						jen.ID("logger").Dot("Error").Call(
							jen.ID("indexInitErr"),
							jen.Lit("setting up items search index"),
						),
						jen.Return().List(jen.ID("nil"),
							jen.ID("indexInitErr")),
					)
				}
				return jen.Null()
			}(),

			jen.Newline(),
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("itemIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.ID("ItemIDURIParamKey"),
					jen.Lit("item"),
				),
				jen.ID("sessionContextDataFetcher").Op(":").Qual(proj.AuthServicePackage(), "FetchContextFromRequest"),
				jen.ID("itemDataManager").Op(":").ID("itemDataManager"),
				jen.ID("encoderDecoder").Op(":").ID("encoder"),
				jen.ID("itemCounter").Op(":").Qual(proj.MetricsPackage(), "EnsureUnitCounter").Call(
					jen.ID("counterProvider"),
					jen.ID("logger"),
					jen.ID("counterName"),
					jen.ID("counterDescription"),
				),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("search").Op(":").ID("searchIndexManager")
					}
					return jen.Null()
				}(),
				jen.ID("tracer").Op(":").Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("serviceName"))),
			jen.Newline(),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	return code
}
