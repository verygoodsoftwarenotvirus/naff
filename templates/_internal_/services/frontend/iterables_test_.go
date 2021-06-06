package frontend

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()
	rn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()

	code.Add(
		jen.Func().IDf("TestService_fetch%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.IDf("example%s", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dotf("fetch%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.IDf("example%s", sn),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fake mode"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("useFakeData").Op("=").ID("true"),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dotf("fetch%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.Qual(constants.AssertionLibrary, "NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error fetching %s", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dotf("fetch%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Newline(),
	)

	createInputs := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		createInputs = append(createInputs, jen.IDf("%sCreationInput%sFormKey", uvn, fsn).Op(":").Values(jen.ID("input").Dot(fsn)))
	}

	code.Add(
		jen.Func().IDf("attach%sCreationInputToRequest", sn).Params(jen.ID("input").Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn))).Params(jen.Op("*").Qual("net/http", "Request")).Body(
			jen.ID("form").Op(":=").Qual("net/url", "Values").Valuesln(
				createInputs...,
			),
			jen.Newline(),
			jen.Return().ID("httptest").Dot("NewRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.Litf("/%s", prn),
				jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_build%sCreatorView", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template and error writing to response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").Op("&").Qual(proj.TestUtilsPackage(), "MockHTTPResponseWriter").Values(),
					jen.ID("res").Dot("On").Call(jen.Lit("Write"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template and error writing to response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").Op("&").Qual(proj.TestUtilsPackage(), "MockHTTPResponseWriter").Values(),
					jen.ID("res").Dot("On").Call(jen.Lit("Write"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sCreatorView", sn).Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_parseFormEncoded%sCreationInput", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expected").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInput", sn)).Call(),
					jen.ID("req").Op(":=").IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("expected")),
					jen.ID("sessionCtxData").Op(":=").Op("&").Qual(proj.TypesPackage(), "SessionContextData").Valuesln(
						jen.ID("ActiveAccountID").Op(":").ID("expected").Dot("BelongsToAccount")),
					jen.Newline(),
					jen.ID("actual").Op(":=").ID("s").Dotf("parseFormEncoded%sCreationInput", sn).Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("sessionCtxData"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error extracting form from request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInput", sn)).Call(),
					jen.Newline(),
					jen.ID("badBody").Op(":=").Op("&").Qual(proj.TestUtilsPackage(), "MockReadCloser").Values(),
					jen.ID("badBody").Dot("On").Call(jen.Lit("Read"), jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().ID("byte").Values())).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/test"),
						jen.ID("badBody"),
					),
					jen.ID("sessionCtxData").Op(":=").Op("&").Qual(proj.TypesPackage(), "SessionContextData").Valuesln(
						jen.ID("ActiveAccountID").Op(":").ID("exampleInput").Dot("BelongsToAccount")),
					jen.Newline(),
					jen.ID("actual").Op(":=").ID("s").Dotf("parseFormEncoded%sCreationInput", sn).Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("sessionCtxData"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values(),
					jen.ID("req").Op(":=").IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.ID("sessionCtxData").Op(":=").Op("&").Qual(proj.TypesPackage(), "SessionContextData").Valuesln(
						jen.ID("ActiveAccountID").Op(":").ID("exampleInput").Dot("BelongsToAccount")),
					jen.Newline(),
					jen.ID("actual").Op(":=").ID("s").Dotf("parseFormEncoded%sCreationInput", sn).Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("sessionCtxData"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_handle%sCreationRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.ID("exampleInput").Dot("BelongsToAccount").Op("=").ID("exampleSessionContextData").Dot("ActiveAccountID"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Create%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.IDf("example%s", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sCreationRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusCreated"),
						jen.ID("res").Dot("Code"),
					),
					jen.Qual(constants.AssertionLibrary, "NotEmpty").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("htmxRedirectionHeader")),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sCreationRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.ID("exampleInput").Dot("BelongsToAccount").Op("=").ID("exampleSessionContextData").Dot("ActiveAccountID"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sCreationInputToRequest", sn).Call(jen.Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values()),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Create%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.IDf("example%s", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sCreationRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error creating %s in database", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.ID("exampleInput").Dot("BelongsToAccount").Op("=").ID("exampleSessionContextData").Dot("ActiveAccountID"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Create%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sCreationRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_build%sEditorView", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.IDf("example%s", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.IDf("example%s", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error fetching %s", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sEditorView", sn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_fetch%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.IDf("example%sList", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.IDf("example%sList", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dotf("fetch%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.IDf("example%sList", sn),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fake mode"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("useFakeData").Op("=").ID("true"),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dotf("fetch%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.Qual(constants.AssertionLibrary, "NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dotf("fetch%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleSessionContextData"),
						jen.ID("req"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_build%sTableView", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%sList", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.IDf("example%sList", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sTableView", pn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%sList", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.IDf("example%sList", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sTableView", pn).Call(jen.ID("false")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sTableView", pn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("build%sTableView", pn).Call(jen.ID("true")).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
		),
		jen.Newline(),
	)

	updateInputs := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		updateInputs = append(updateInputs, jen.IDf("%sUpdateInput%sFormKey", uvn, fsn).Op(":").Values(jen.ID("input").Dot(fsn)))
	}

	code.Add(
		jen.Func().IDf("attach%sUpdateInputToRequest", sn).Params(jen.ID("input").Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn))).Params(jen.Op("*").Qual("net/http", "Request")).Body(
			jen.ID("form").Op(":=").Qual("net/url", "Values").Valuesln(
				updateInputs...,
			),
			jen.Newline(),
			jen.Return().ID("httptest").Dot("NewRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.Litf("/%s", prn),
				jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_parseFormEncoded%sUpdateInput", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("expected").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("sessionCtxData").Op(":=").Op("&").Qual(proj.TypesPackage(), "SessionContextData").Valuesln(
						jen.ID("ActiveAccountID").Op(":").ID("expected").Dot("BelongsToAccount")),
					jen.Newline(),
					jen.ID("req").Op(":=").IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("expected")),
					jen.Newline(),
					jen.ID("actual").Op(":=").ID("s").Dotf("parseFormEncoded%sUpdateInput", sn).Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("sessionCtxData"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid form"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("sessionCtxData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("badBody").Op(":=").Op("&").Qual(proj.TestUtilsPackage(), "MockReadCloser").Values(),
					jen.ID("badBody").Dot("On").Call(jen.Lit("Read"), jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().ID("byte").Values())).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/test"),
						jen.ID("badBody"),
					),
					jen.Newline(),
					jen.ID("actual").Op(":=").ID("s").Dotf("parseFormEncoded%sUpdateInput", sn).Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("sessionCtxData"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input attached to valid form"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)).Values(),
					jen.ID("sessionCtxData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("req").Op(":=").IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("actual").Op(":=").ID("s").Dotf("parseFormEncoded%sUpdateInput", sn).Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("sessionCtxData"),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_handle%sUpdateRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.IDf("example%s", sn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
						jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Call(jen.ID("nil")),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sUpdateRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sUpdateRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleInput").Op(":=").Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)).Values(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sUpdateRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sUpdateRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
					).Dot("Return").Call(
						jen.IDf("example%s", sn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
						jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Call(jen.ID("nil")),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sUpdateRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestService_handle%sDeletionRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%sList", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Archive%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.IDf("example%sList", sn),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodDelete"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sDeletionRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodDelete"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sDeletionRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("unauthorizedRedirectResponseCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error archiving %s", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Archive%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodDelete"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sDeletionRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error retrieving new list of %s", pcn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleSessionContextData").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeSessionContextData").Call(),
					jen.Newline(),
					jen.ID("s").Dot("sessionContextDataFetcher").Op("=").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("exampleSessionContextData"), jen.ID("nil")),
					),
					jen.Newline(),
					jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
					jen.ID("rpm").Dot("On").Callln(
						jen.Lit("BuildRouteParamIDFetcher"),
						jen.Qual(constants.MockPkg, "Anything"),
						jen.IDf("%sIDURLParamKey", uvn),
						jen.Lit(rn),
					).Dot("Return").Call(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("ID")),
					),
					jen.ID("s").Dot("routeParamManager").Op("=").ID("rpm"),
					jen.Newline(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Archive%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.ID("exampleSessionContextData").Dot("Requester").Dot("UserID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Get%s", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleSessionContextData").Dot("ActiveAccountID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("dataStore").Op("=").ID("mockDB"),
					jen.Newline(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodDelete"),
						jen.Litf("/%s", prn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("s").Dotf("handle%sDeletionRequest", sn).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
