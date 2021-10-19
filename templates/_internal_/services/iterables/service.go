package iterables

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

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
	prn := typ.Name.PluralRouteName()

	code.Add(
		jen.Const().Defs(
			jen.ID("serviceName").String().Equals().Litf("%s_service", prn),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Underscore().Qual(proj.TypesPackage(), fmt.Sprintf("%sDataService", sn)).Equals().Parens(jen.PointerTo().ID("service")).Call(jen.Nil()),
		jen.Newline(),
	)

	structLines := []jen.Code{
		jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
		jen.IDf("%sDataManager", uvn).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)),
	}
	for _, dep := range proj.FindOwnerTypeChain(typ) {
		structLines = append(structLines, jen.IDf("%sIDFetcher", dep.Name.UnexportedVarName()).Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()))
	}
	structLines = append(structLines,
		jen.IDf("%sIDFetcher", uvn).Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()),
		jen.ID("sessionContextDataFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")),
		jen.ID("preWritesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
		jen.ID("preUpdatesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
		jen.ID("preArchivesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
		jen.ID("encoderDecoder").Qual(proj.EncodingPackage(), "ServerEncoderDecoder"),
		jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("search").ID("SearchIndex")
			}
			return jen.Null()
		}(),
		jen.ID("async").Bool(),
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

	serviceInitLines := []jen.Code{
		jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")),
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		tuvn := dep.Name.UnexportedVarName()
		tsn := dep.Name.Singular()

		serviceInitLines = append(serviceInitLines,
			jen.IDf("%sIDFetcher", tuvn).MapAssign().ID("routeParamManager").Dot("BuildRouteParamStringIDFetcher").Call(
				jen.Qual(proj.ServicePackage(dep.Name.PackageName()), fmt.Sprintf("%sIDURIParamKey", tsn)),
			),
		)
	}

	serviceInitLines = append(serviceInitLines,
		jen.IDf("%sIDFetcher", uvn).MapAssign().ID("routeParamManager").Dot("BuildRouteParamStringIDFetcher").Call(
			jen.IDf("%sIDURIParamKey", sn),
		),
		jen.ID("sessionContextDataFetcher").MapAssign().Qual(proj.AuthServicePackage(), "FetchContextFromRequest"),
		jen.IDf("%sDataManager", uvn).MapAssign().IDf("%sDataManager", uvn),
		jen.ID("preWritesPublisher").MapAssign().ID("preWritesPublisher"),
		jen.ID("preUpdatesPublisher").MapAssign().ID("preUpdatesPublisher"),
		jen.ID("preArchivesPublisher").MapAssign().ID("preArchivesPublisher"),
		jen.ID("encoderDecoder").MapAssign().ID("encoder"),
		jen.ID("async").MapAssign().ID("cfg").Dot("Async"),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.ID("search").MapAssign().ID("searchIndexManager")
			}
			return jen.Null()
		}(),
		jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("serviceName")),
	)

	serviceBodyLines := []jen.Code{}

	if typ.SearchEnabled {
		searchFields := []jen.Code{}
		for _, field := range typ.Fields {
			if field.Type == "string" {
				searchFields = append(searchFields, jen.Lit(field.Name.UnexportedVarName()))
			}
		}

		serviceBodyLines = append(serviceBodyLines,
			jen.ID("client").Assign().AddressOf().Qual("net/http", "Client").Values(jen.ID("Transport").MapAssign().Qual(proj.InternalTracingPackage(), "BuildTracedHTTPTransport").Call(jen.Qual("time", "Second"))),
			jen.Newline(),
			jen.List(jen.ID("searchIndexManager"), jen.ID("err")).Assign().ID("searchIndexProvider").Call(
				append([]jen.Code{
					constants.CtxVar(),
					jen.ID("logger"),
					jen.ID("client"),
					jen.Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("cfg").Dot("SearchIndexPath")),
					jen.Lit(prn),
				}, searchFields...)...,
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting up search index: %w"), jen.Err())),
			),
			jen.Newline(),
		)
	}

	serviceBodyLines = append(serviceBodyLines,
		jen.List(jen.ID("preWritesPublisher"), jen.Err()).Assign().ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("cfg").Dot("PreWritesTopicName")),
		jen.If(jen.Err().DoesNotEqual().Nil()).Body(
			jen.Return(jen.List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting up queue provider: %w"), jen.Err()))),
		),
		jen.Newline(),
		jen.List(jen.ID("preUpdatesPublisher"), jen.Err()).Assign().ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("cfg").Dot("PreUpdatesTopicName")),
		jen.If(jen.Err().DoesNotEqual().Nil()).Body(
			jen.Return(jen.List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting up queue provider: %w"), jen.Err()))),
		),
		jen.Newline(),
		jen.List(jen.ID("preArchivesPublisher"), jen.Err()).Assign().ID("publisherProvider").Dot("ProviderPublisher").Call(jen.ID("cfg").Dot("PreArchivesTopicName")),
		jen.If(jen.Err().DoesNotEqual().Nil()).Body(
			jen.Return(jen.List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("setting up queue provider: %w"), jen.Err()))),
		),
		jen.Newline(),
	)

	serviceBodyLines = append(serviceBodyLines,
		jen.ID("svc").Assign().AddressOf().ID("service").Valuesln(
			serviceInitLines...,
		),
		jen.Newline(),
		jen.Return().List(jen.ID("svc"), jen.Nil()),
	)

	code.Add(
		jen.Commentf("ProvideService builds a new %sService.", pn),
		jen.Newline(),
		jen.Func().ID("ProvideService").Paramsln(
			constants.CtxParam(),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("cfg").PointerTo().ID("Config"),
			jen.IDf("%sDataManager", uvn).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)),
			jen.ID("encoder").Qual(proj.EncodingPackage(), "ServerEncoderDecoder"),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("searchIndexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider")
				}
				return jen.Null()
			}(),
			jen.ID("routeParamManager").Qual(proj.RoutingPackage(), "RouteParamManager"),
			jen.ID("publisherProvider").Qual(proj.InternalMessageQueuePublishersPackage(), "PublisherProvider"),
		).Params(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sDataService", sn)), jen.ID("error")).Body(
			serviceBodyLines...,
		),
		jen.Newline(),
	)

	return code
}
