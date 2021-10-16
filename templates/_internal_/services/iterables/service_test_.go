package iterables

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()

	code.Add(
		jen.Func().ID("buildTestService").Params().Params(jen.PointerTo().ID("service")).Body(
			jen.Return().AddressOf().ID("service").Valuesln(
				jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.IDf("%sDataManager", uvn).MapAssign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
				jen.ID("async").MapAssign().True(),
				jen.IDf("%sIDFetcher", uvn).MapAssign().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.String()).SingleLineBody(jen.Return().EmptyString()),
				jen.ID("encoderDecoder").MapAssign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("search").MapAssign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
					}
					return jen.Null()
				}(),
				jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.Lit("test")))),
		jen.Newline(),
	)

	routeParamFetchers := []jen.Code{}
	for _, dep := range proj.FindOwnerTypeChain(typ) {
		tsn := dep.Name.Singular()

		routeParamFetchers = append(routeParamFetchers,
			jen.ID("rpm").Dot("On").Callln(
				jen.Lit("BuildRouteParamStringIDFetcher"),
				jen.Qual(proj.ServicePackage(dep.Name.PackageName()), fmt.Sprintf("%sIDURIParamKey", tsn)),
			).Dot("Return").Call(jen.Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()).SingleLineBody(jen.Return().EmptyString())),
		)
	}

	subtestZeroLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		constants.CreateCtx(),
		jen.ID("rpm").Assign().Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
	}

	subtestZeroLines = append(subtestZeroLines, routeParamFetchers...)
	subtestZeroLines = append(subtestZeroLines,
		jen.ID("rpm").Dot("On").Callln(
			jen.Lit("BuildRouteParamStringIDFetcher"),
			jen.IDf("%sIDURIParamKey", sn),
		).Dot("Return").Call(jen.Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()).SingleLineBody(jen.Return().EmptyString())),
		jen.Newline(),
		jen.ID("cfg").Assign().ID("Config").Valuesln(
			jen.ID("SearchIndexPath").MapAssign().Lit("example/path"),
			jen.ID("PreWritesTopicName").MapAssign().Lit("pre-writes"),
			jen.ID("PreUpdatesTopicName").MapAssign().Lit("pre-updates"),
			jen.ID("PreArchivesTopicName").MapAssign().Lit("pre-archives"),
		),
		jen.Newline(),
		jen.ID("pp").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "ProducerProvider").Values(),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreWritesTopicName")).Dot("Return").Call(jen.AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(), jen.Nil()),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreUpdatesTopicName")).Dot("Return").Call(jen.AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(), jen.Nil()),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreArchivesTopicName")).Dot("Return").Call(jen.AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(), jen.Nil()),
		jen.Newline(),
		jen.List(jen.ID("s"),
			jen.ID("err")).Assign().ID("ProvideService").Callln(
			constants.CtxVar(),
			jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.AddressOf().ID("cfg"),
			jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Func().Params(
						jen.Qual("context", "Context"),
						jen.Qual(proj.InternalLoggingPackage(), "Logger"),
						jen.PointerTo().Qual("net/http", "Client"),
						jen.Qual(proj.InternalSearchPackage(), "IndexPath"),
						jen.Qual(proj.InternalSearchPackage(), "IndexName"),
						jen.Spread().String(),
					).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(
							jen.AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.Nil(),
						),
					)
				}
				return jen.Null()
			}(),
			jen.ID("rpm"),
			jen.ID("pp"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "NotNil").Call(
			jen.ID("t"),
			jen.ID("s"),
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("rpm"),
			jen.ID("pp"),
		),
	)

	subtestZero := jen.ID("T").Dot("Run").Call(
		jen.Lit("standard"),
		jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			subtestZeroLines...,
		),
	)

	subtestOneLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		constants.CreateCtx(),
	}

	subtestOneLines = append(subtestOneLines,
		jen.ID("cfg").Assign().ID("Config").Valuesln(
			jen.ID("SearchIndexPath").MapAssign().Lit("example/path"),
			jen.ID("PreWritesTopicName").MapAssign().Lit("pre-writes"),
			jen.ID("PreUpdatesTopicName").MapAssign().Lit("pre-updates"),
			jen.ID("PreArchivesTopicName").MapAssign().Lit("pre-archives"),
		),
		jen.Newline(),
		jen.ID("pp").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "ProducerProvider").Values(),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreWritesTopicName")).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		jen.Newline(),
		jen.List(jen.ID("s"),
			jen.ID("err")).Assign().ID("ProvideService").Callln(
			constants.CtxVar(),
			jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.AddressOf().ID("cfg"),
			jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Func().Params(
						jen.Qual("context", "Context"),
						jen.Qual(proj.InternalLoggingPackage(), "Logger"),
						jen.PointerTo().Qual("net/http", "Client"),
						jen.Qual(proj.InternalSearchPackage(), "IndexPath"),
						jen.Qual(proj.InternalSearchPackage(), "IndexName"),
						jen.Spread().String(),
					).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(
							jen.AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.Nil(),
						),
					)
				}
				return jen.Null()
			}(),
			jen.Nil(),
			jen.ID("pp"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(
			jen.ID("t"),
			jen.ID("s"),
		),
		jen.Qual(constants.AssertionLibrary, "Error").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("pp"),
		),
	)

	subtestOne := jen.ID("T").Dot("Run").Call(
		jen.Lit("with error providing pre-writes producer"),
		jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			subtestOneLines...,
		),
	)

	subtestTwoLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		constants.CreateCtx(),
	}

	subtestTwoLines = append(subtestTwoLines,
		jen.ID("cfg").Assign().ID("Config").Valuesln(
			jen.ID("SearchIndexPath").MapAssign().Lit("example/path"),
			jen.ID("PreWritesTopicName").MapAssign().Lit("pre-writes"),
			jen.ID("PreUpdatesTopicName").MapAssign().Lit("pre-updates"),
			jen.ID("PreArchivesTopicName").MapAssign().Lit("pre-archives"),
		),
		jen.Newline(),
		jen.ID("pp").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "ProducerProvider").Values(),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreWritesTopicName")).Dot("Return").Call(jen.AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(), jen.Nil()),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreUpdatesTopicName")).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		jen.Newline(),
		jen.List(jen.ID("s"),
			jen.ID("err")).Assign().ID("ProvideService").Callln(
			constants.CtxVar(),
			jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.AddressOf().ID("cfg"),
			jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Func().Params(
						jen.Qual("context", "Context"),
						jen.Qual(proj.InternalLoggingPackage(), "Logger"),
						jen.PointerTo().Qual("net/http", "Client"),
						jen.Qual(proj.InternalSearchPackage(), "IndexPath"),
						jen.Qual(proj.InternalSearchPackage(), "IndexName"),
						jen.Spread().String(),
					).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(
							jen.AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.Nil(),
						),
					)
				}
				return jen.Null()
			}(),
			jen.Nil(),
			jen.ID("pp"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(
			jen.ID("t"),
			jen.ID("s"),
		),
		jen.Qual(constants.AssertionLibrary, "Error").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("pp"),
		),
	)

	subtestTwo := jen.ID("T").Dot("Run").Call(
		jen.Lit("with error providing pre-updates producer"),
		jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			subtestTwoLines...,
		),
	)

	subtestThreeLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		constants.CreateCtx(),
	}

	subtestThreeLines = append(subtestThreeLines,
		jen.ID("cfg").Assign().ID("Config").Valuesln(
			jen.ID("SearchIndexPath").MapAssign().Lit("example/path"),
			jen.ID("PreWritesTopicName").MapAssign().Lit("pre-writes"),
			jen.ID("PreUpdatesTopicName").MapAssign().Lit("pre-updates"),
			jen.ID("PreArchivesTopicName").MapAssign().Lit("pre-archives"),
		),
		jen.Newline(),
		jen.ID("pp").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "ProducerProvider").Values(),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreWritesTopicName")).Dot("Return").Call(jen.AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(), jen.Nil()),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreUpdatesTopicName")).Dot("Return").Call(jen.AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(), jen.Nil()),
		jen.ID("pp").Dot("On").Call(jen.Lit("ProviderPublisher"), jen.ID("cfg").Dot("PreArchivesTopicName")).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		jen.Newline(),
		jen.List(jen.ID("s"),
			jen.ID("err")).Assign().ID("ProvideService").Callln(
			constants.CtxVar(),
			jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.AddressOf().ID("cfg"),
			jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Func().Params(
						jen.Qual("context", "Context"),
						jen.Qual(proj.InternalLoggingPackage(), "Logger"),
						jen.PointerTo().Qual("net/http", "Client"),
						jen.Qual(proj.InternalSearchPackage(), "IndexPath"),
						jen.Qual(proj.InternalSearchPackage(), "IndexName"),
						jen.Spread().String(),
					).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(
							jen.AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.Nil(),
						),
					)
				}
				return jen.Null()
			}(),
			jen.Nil(),
			jen.ID("pp"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(
			jen.ID("t"),
			jen.ID("s"),
		),
		jen.Qual(constants.AssertionLibrary, "Error").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("pp"),
		),
	)

	subtestThree := jen.ID("T").Dot("Run").Call(
		jen.Lit("with error providing pre-archives producer"),
		jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			subtestThreeLines...,
		),
	)

	searchSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		constants.CreateCtx(),
	}

	searchSubtestLines = append(searchSubtestLines,
		jen.ID("cfg").Assign().ID("Config").Valuesln(
			jen.ID("SearchIndexPath").MapAssign().Lit("example/path"),
			jen.ID("PreWritesTopicName").MapAssign().Lit("pre-writes"),
			jen.ID("PreUpdatesTopicName").MapAssign().Lit("pre-updates"),
			jen.ID("PreArchivesTopicName").MapAssign().Lit("pre-archives"),
		),
		jen.Newline(),
		jen.List(jen.ID("s"),
			jen.ID("err")).Assign().ID("ProvideService").Callln(
			constants.CtxVar(),
			jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.AddressOf().ID("cfg"),
			jen.AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
			jen.Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Func().Params(
						jen.Qual("context", "Context"),
						jen.Qual(proj.InternalLoggingPackage(), "Logger"),
						jen.PointerTo().Qual("net/http", "Client"),
						jen.Qual(proj.InternalSearchPackage(), "IndexPath"),
						jen.Qual(proj.InternalSearchPackage(), "IndexName"),
						jen.Spread().String(),
					).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(
							jen.Nil(),
							jen.Qual("errors", "New").Call(jen.Lit("blah")),
						),
					)
				}
				return jen.Null()
			}(),
			jen.Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
			jen.Nil(),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(
			jen.ID("t"),
			jen.ID("s"),
		),
		jen.Qual(constants.AssertionLibrary, "Error").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
	)

	searchSubtest := jen.ID("T").Dot("Run").Call(
		jen.Lit("with error providing index"),
		jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			searchSubtestLines...,
		),
	)

	code.Add(
		jen.Func().IDf("TestProvide%sService", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			subtestZero,
			jen.Newline(),
			subtestOne,
			jen.Newline(),
			subtestTwo,
			jen.Newline(),
			subtestThree,
			jen.Newline(),
			utils.ConditionalCode(typ.SearchEnabled, searchSubtest),
		),
	)

	return code
}
