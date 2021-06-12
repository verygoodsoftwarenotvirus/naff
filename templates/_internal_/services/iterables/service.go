package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()
	rn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()

	code.Add(
		jen.Const().Defs(
			jen.ID("counterName").Qual(proj.MetricsPackage(), "CounterName").Op("=").Lit(prn),
			jen.ID("counterDescription").ID("string").Op("=").Litf("the number of %s managed by the %s service", pcn, pcn),
			jen.ID("serviceName").ID("string").Op("=").Litf("%s_service", prn),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().ID("_").Qual(proj.TypesPackage(), fmt.Sprintf("%sDataService", sn)).Op("=").Parens(jen.Op("*").ID("service")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	structLines := []jen.Code{
		jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
		jen.IDf("%sDataManager", uvn).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)),
	}
	for _, dep := range proj.FindOwnerTypeChain(typ) {
		structLines = append(structLines, jen.IDf("%sIDFetcher", dep.Name.UnexportedVarName()).Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")))
	}
	structLines = append(structLines,
		jen.IDf("%sIDFetcher", uvn).Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		jen.ID("sessionContextDataFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"),
			jen.ID("error")),
		jen.IDf("%sCounter", uvn).Qual(proj.MetricsPackage(), "UnitCounter"),
		jen.ID("encoderDecoder").Qual(proj.EncodingPackage(), "ServerEncoderDecoder"),
		jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("search").ID("SearchIndex")
			}
			return jen.Null()
		}(),
	)

	code.Add(
		jen.Type().Defs(
			jen.Comment("SearchIndex is a type alias for dependency injection's sake."),
			jen.ID("SearchIndex").Qual(proj.InternalSearchPackage(), "IndexManager"),
			jen.Newline(),
			jen.Commentf("service handles %s.", pcn),
			jen.ID("service").Struct(
				structLines...,
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("ProvideService builds a new %sService.", pn),
		jen.Newline(),
		jen.Func().ID("ProvideService").Paramsln(
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("cfg").ID("Config"),
			jen.IDf("%sDataManager", uvn).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)),
			jen.ID("encoder").Qual(proj.EncodingPackage(), "ServerEncoderDecoder"),
			jen.ID("counterProvider").Qual(proj.MetricsPackage(), "UnitCounterProvider"),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("indexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider")
				}
				return jen.Null()
			}(),
			jen.ID("routeParamManager").Qual(proj.RoutingPackage(), "RouteParamManager"),
		).Params(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sDataService", sn)), jen.ID("error")).Body(

			func() jen.Code {
				if typ.SearchEnabled {
					return jen.List(jen.ID("searchIndexManager"),
						jen.ID("indexInitErr")).Op(":=").ID("indexProvider").Call(
						jen.Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("cfg").Dot("SearchIndexPath")),
						jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sSearchIndexName", pn)),
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
							jen.Litf("setting up %s search index", pcn),
						),
						jen.Return().List(jen.ID("nil"),
							jen.ID("indexInitErr")),
					)
				}
				return jen.Null()
			}(),

			jen.Newline(),
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(jen.ID("logger").Op(":").Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")),
				jen.IDf("%sIDFetcher", uvn).Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
					jen.ID("logger"),
					jen.IDf("%sIDURIParamKey", sn),
					jen.Lit(rn),
				),
				jen.ID("sessionContextDataFetcher").Op(":").Qual(proj.AuthServicePackage(), "FetchContextFromRequest"),
				jen.IDf("%sDataManager", uvn).Op(":").IDf("%sDataManager", uvn),
				jen.ID("encoderDecoder").Op(":").ID("encoder"),
				jen.IDf("%sCounter", uvn).Op(":").Qual(proj.MetricsPackage(), "EnsureUnitCounter").Call(
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
