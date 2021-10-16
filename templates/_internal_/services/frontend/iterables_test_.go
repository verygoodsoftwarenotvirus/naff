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

	code.Add(buildTestServiceFetchSomething(proj, typ)...)
	code.Add(buildAttachSomethingCreationInputToRequest(proj, typ)...)
	code.Add(buildTestServiceBuildSomethingCreatorView(proj, typ)...)
	code.Add(buildTestServiceParseFormEncodedSomethingCreationInput(proj, typ)...)
	code.Add(buildTestServiceHandleSomethingCreationRequest(proj, typ)...)
	code.Add(buildTestServiceBuildSomethingEditorView(proj, typ)...)
	code.Add(buildTestServiceFetchListOfSomethings(proj, typ)...)
	code.Add(buildTestServiceBuildSomethingTableView(proj, typ)...)
	code.Add(buildAttachSomethingUpdateInputToRequest(proj, typ)...)
	code.Add(buildTestServiceParseFormEncodedSomethingUpdateInput(proj, typ)...)
	code.Add(buildTestServiceHandleSomethingUpdateRequest(proj, typ)...)
	code.Add(buildTestServiceHandleSomethingArchiveRequest(proj, typ)...)

	return code
}

func buildTestIDFetchers(proj *models.Project, typ models.DataType, includeType bool) []jen.Code {
	lines := []jen.Code{}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.IDf("example%sID", dep.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
			jen.ID("s").Dot("service").Dotf("%sIDFetcher", dep.Name.UnexportedVarName()).Equals().Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()).Body(
				jen.Return(jen.IDf("example%sID", dep.Name.Singular())),
			),
			jen.Newline(),
		)
	}

	if includeType {
		lines = append(lines,
			jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%s", typ.Name.Singular()).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToAccount {
					return jen.IDf("example%s", typ.Name.Singular()).Dot("BelongsToAccount").Equals().ID("s").Dot("exampleAccount").Dot("ID")
				}
				return jen.Null()
			}(),
			jen.ID("s").Dot("service").Dotf("%sIDFetcher", typ.Name.UnexportedVarName()).Equals().Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()).Body(
				jen.Return(jen.IDf("example%s", typ.Name.Singular()).Dot("ID")),
			),
			jen.Newline(),
		)
	}

	return lines
}

func buildCallArgsForRetrievalTest(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{jen.Litf("Get%s", typ.Name.Singular()), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("example%sID", pt.Name.Singular()))
	}
	params = append(params, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		params = append(params, jen.ID("s").Dot("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildTestServiceFetchSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	prn := typ.Name.PluralRouteName()

	mockCallArgs := buildCallArgsForRetrievalTest(proj, typ)

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(jen.IDf("example%s", sn), jen.Nil()),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("s").Dot("service").Dotf("fetch%s", sn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.IDf("example%s", sn), jen.ID("actual")),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.ID("s").Dot("service").Dot("useFakeData").Equals().True(),
		jen.Newline(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("s").Dot("service").Dotf("fetch%s", sn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "NotNil").Call(jen.ID("t"), jen.ID("actual")),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(jen.ID("t"), jen.ID("err")),
	}

	thirdSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	thirdSubtestLines = append(thirdSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	thirdSubtestLines = append(thirdSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("s").Dot("service").Dotf("fetch%s", sn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestService_fetch%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(firstSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fake mode"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					secondSubtestLines...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error fetching %s", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					thirdSubtestLines...,
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildAttachSomethingCreationInputToRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	prn := typ.Name.PluralRouteName()

	createInputs := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		if !field.IsPointer {
			createInputs = append(createInputs, jen.IDf("%sCreationInput%sFormKey", uvn, fsn).MapAssign().Values(jen.ID("anyToString").Call(jen.ID("input").Dot(fsn))))
		}
	}

	bodyLines := []jen.Code{
		jen.ID("form").Assign().Qual("net/url", "Values").Valuesln(createInputs...),
		jen.Newline(),
	}

	for _, field := range typ.Fields {
		if field.IsPointer {
			bodyLines = append(bodyLines,
				jen.If(jen.ID("input").Dot(field.Name.Singular()).DoesNotEqual().Nil()).Body(
					jen.ID("form").Dot("Set").Call(jen.IDf("%sCreationInput%sFormKey", uvn, field.Name.Singular()), jen.ID("anyToString").Call(jen.PointerTo().ID("input").Dot(field.Name.Singular()))),
				),
				jen.Newline(),
			)
		}
	}

	bodyLines = append(bodyLines,
		jen.Return().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodPost"),
			jen.Litf("/%s", prn),
			jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
		),
	)

	lines := []jen.Code{
		jen.Func().IDf("attach%sCreationInputToRequest", sn).Params(jen.ID("input").PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sDatabaseCreationInput", sn))).Params(jen.PointerTo().Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildTestServiceBuildSomethingCreatorView(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	prn := typ.Name.PluralRouteName()

	lines := []jen.Code{
		jen.Func().IDf("TestService_build%sCreatorView", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.Newline(),
					jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("s").Dot("service").Dotf("build%sCreatorView", sn).Call(jen.True()).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.Newline(),
					jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("s").Dot("service").Dotf("build%sCreatorView", sn).Call(jen.False()).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("s").Dot("service").Dot("sessionContextDataFetcher").Equals().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("s").Dot("service").Dotf("build%sCreatorView", sn).Call(jen.False()).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("unauthorizedRedirectResponseCode"), jen.ID("res").Dot("Code")),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template and error writing to response"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.Newline(),
					jen.ID("res").Assign().AddressOf().Qual(proj.TestUtilsPackage(), "MockHTTPResponseWriter").Values(),
					jen.ID("res").Dot("On").Call(jen.Lit("Write"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(
						jen.Zero(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("s").Dot("service").Dotf("build%sCreatorView", sn).Call(jen.True()).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template and error writing to response"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.Newline(),
					jen.ID("res").Assign().AddressOf().Qual(proj.TestUtilsPackage(), "MockHTTPResponseWriter").Values(),
					jen.ID("res").Dot("On").Call(jen.Lit("Write"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(
						jen.Zero(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Litf("/%s", prn),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("s").Dot("service").Dotf("build%sCreatorView", sn).Call(jen.False()).Call(
						jen.ID("res"),
						jen.ID("req"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestServiceParseFormEncodedSomethingCreationInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Func().IDf("TestService_parseFormEncoded%sCreationInput", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("expected").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInput", sn)).Call(),
					jen.ID("expected").Dot("ID").Equals().EmptyString(),
					func() jen.Code {
						if typ.BelongsToAccount {
							return jen.ID("expected").Dot("BelongsToAccount").Equals().ID("s").Dot("exampleAccount").Dot("ID")
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().Zero()
						}
						return jen.Null()
					}(),
					jen.ID("req").Assign().IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("expected")),
					jen.Newline(),
					jen.ID("actual").Assign().ID("s").Dot("service").Dotf("parseFormEncoded%sCreationInput", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("req"),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("s").Dot("sessionCtxData")),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error extracting form from request"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("badBody").Assign().AddressOf().Qual(proj.TestUtilsPackage(), "MockReadCloser").Values(),
					jen.ID("badBody").Dot("On").Call(jen.Lit("Read"), jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().ID("byte").Values())).Dot("Return").Call(
						jen.Zero(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/test"),
						jen.ID("badBody"),
					),
					jen.Newline(),
					jen.ID("actual").Assign().ID("s").Dot("service").Dotf("parseFormEncoded%sCreationInput", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("req"),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("s").Dot("sessionCtxData")),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(jen.ID("t"), jen.ID("actual")),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestServiceHandleSomethingCreationRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.ID("exampleInput").Dot("ID").Equals().EmptyString(),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleInput").Dot("BelongsToAccount").Equals().ID("s").Dot("sessionCtxData").Dot("ActiveAccountID")),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			jen.Litf("Create%s", sn),
			jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			jen.ID("exampleInput"),
		).Dot("Return").Call(
			jen.IDf("example%s", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sCreationRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusCreated"), jen.ID("res").Dot("Code")),
		jen.Qual(constants.AssertionLibrary, "NotEmpty").Call(
			jen.ID("t"),
			jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("htmxRedirectionHeader")),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	secondSubtestLines = append(secondSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	secondSubtestLines = append(secondSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.ID("s").Dot("service").Dot("sessionContextDataFetcher").Equals().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
			jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sCreationRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("unauthorizedRedirectResponseCode"), jen.ID("res").Dot("Code")),
	)

	thirdSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	thirdSubtestLines = append(thirdSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	thirdSubtestLines = append(thirdSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleInput").Dot("BelongsToAccount").Equals().ID("s").Dot("sessionCtxData").Dot("ActiveAccountID")),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sCreationInputToRequest", sn).Call(jen.AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values()),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			jen.Litf("Create%s", sn),
			jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			jen.ID("exampleInput"),
			jen.ID("s").Dot("sessionCtxData").Dot("Requester").Dot("UserID"),
		).Dot("Return").Call(
			jen.IDf("example%s", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sCreationRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	fourthSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	fourthSubtestLines = append(fourthSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	fourthSubtestLines = append(fourthSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.ID("exampleInput").Dot("ID").Equals().EmptyString(),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleInput").Dot("BelongsToAccount").Equals().ID("s").Dot("sessionCtxData").Dot("ActiveAccountID")),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sCreationInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			jen.Litf("Create%s", sn),
			jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			jen.ID("exampleInput"),
		).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sCreationRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestService_handle%sCreationRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(firstSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(secondSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error creating %s in database", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(fourthSubtestLines...),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildCallArgsForEditorViewTest(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{jen.Litf("Get%s", typ.Name.Singular()), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("example%sID", pt.Name.Singular()))
	}
	params = append(params, jen.IDf("example%s", typ.Name.Singular()).Dot("ID"))

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		params = append(params, jen.ID("s").Dot("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildTestServiceBuildSomethingEditorView(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	prn := typ.Name.PluralRouteName()

	mockCallArgs := buildCallArgsForEditorViewTest(proj, typ)

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(mockCallArgs...).Dot("Return").Call(
			jen.IDf("example%s", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sEditorView", sn).Call(jen.True()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	secondSubtestLines = append(secondSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	secondSubtestLines = append(secondSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(mockCallArgs...).Dot("Return").Call(
			jen.IDf("example%s", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sEditorView", sn).Call(jen.False()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	thirdSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dot("sessionContextDataFetcher").Equals().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
			jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sEditorView", sn).Call(jen.True()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("unauthorizedRedirectResponseCode"), jen.ID("res").Dot("Code")),
	}

	fourthSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	fourthSubtestLines = append(fourthSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	fourthSubtestLines = append(fourthSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(mockCallArgs...).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sEditorView", sn).Call(jen.True()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestService_build%sEditorView", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(firstSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(secondSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(thirdSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error fetching %s", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					fourthSubtestLines...,
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildCallArgsForListRetrievalTest(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{jen.Litf("Get%s", typ.Name.Plural()), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("example%sID", pt.Name.Singular()))
	}

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		params = append(params, jen.ID("s").Dot("sessionCtxData").Dot("ActiveAccountID"))
	}

	params = append(params, jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), "QueryFilter").Values()))

	return params
}

func buildTestServiceFetchListOfSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	prn := typ.Name.PluralRouteName()

	mockCallArgs := buildCallArgsForListRetrievalTest(proj, typ)

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, false)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(
			jen.IDf("example%sList", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("s").Dot("service").Dotf("fetch%s", pn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.IDf("example%sList", sn), jen.ID("actual")),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.ID("s").Dot("service").Dot("useFakeData").Equals().True(),
		jen.Newline(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("s").Dot("service").Dotf("fetch%s", pn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "NotNil").Call(
			jen.ID("t"),
			jen.ID("actual"),
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("err"),
		),
	}

	withErrorFetchingDataLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	withErrorFetchingDataLines = append(withErrorFetchingDataLines, buildTestIDFetchers(proj, typ, false)...)

	withErrorFetchingDataLines = append(withErrorFetchingDataLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("s").Dot("service").Dotf("fetch%s", pn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(jen.ID("t"), jen.ID("actual")),
		jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.ID("err")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestService_fetch%s", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(firstSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with fake mode"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(secondSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(withErrorFetchingDataLines...),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestServiceBuildSomethingTableView(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	prn := typ.Name.PluralRouteName()

	mockCallArgs := buildCallArgsForListRetrievalTest(proj, typ)

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, false)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		func() jen.Code {
			if typ.BelongsToStruct != nil || typ.BelongsToAccount {
				return jen.For(jen.List(jen.Underscore(), jen.ID(uvn)).Assign().Range().IDf("example%sList", sn)).Dot(pn).Body(
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.ID(uvn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.BelongsToAccount {
							return jen.ID(uvn).Dot("BelongsToAccount").Equals().ID("s").Dot("exampleAccount").Dot("ID")
						}
						return jen.Null()
					}(),
				)
			}
			return jen.Null()
		}(),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(
			jen.IDf("example%sList", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sTableView", pn).Call(jen.True()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	secondSubtestLines = append(secondSubtestLines, buildTestIDFetchers(proj, typ, false)...)

	secondSubtestLines = append(secondSubtestLines,
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(
			jen.IDf("example%sList", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sTableView", pn).Call(jen.False()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	thirdSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.ID("s").Dot("service").Dot("sessionContextDataFetcher").Equals().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
			jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sTableView", pn).Call(jen.True()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("unauthorizedRedirectResponseCode"), jen.ID("res").Dot("Code")),
	}

	fourthSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	fourthSubtestLines = append(fourthSubtestLines, buildTestIDFetchers(proj, typ, false)...)

	fourthSubtestLines = append(fourthSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("build%sTableView", pn).Call(jen.True()).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestService_build%sTableView", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with base template"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(firstSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without base template"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(secondSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(thirdSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(fourthSubtestLines...),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildAttachSomethingUpdateInputToRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	prn := typ.Name.PluralRouteName()

	updateInputs := []jen.Code{}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()
		if !field.IsPointer {
			updateInputs = append(updateInputs, jen.IDf("%sUpdateInput%sFormKey", uvn, fsn).MapAssign().Values(jen.ID("anyToString").Call(jen.ID("input").Dot(fsn))))
		}
	}

	bodyLines := []jen.Code{
		jen.ID("form").Assign().Qual("net/url", "Values").Valuesln(updateInputs...),
		jen.Newline(),
	}

	for _, field := range typ.Fields {
		if field.IsPointer {
			bodyLines = append(bodyLines,
				jen.If(jen.ID("input").Dot(field.Name.Singular()).DoesNotEqual().Nil()).Body(
					jen.ID("form").Dot("Set").Call(jen.IDf("%sUpdateInput%sFormKey", uvn, field.Name.Singular()), jen.ID("anyToString").Call(jen.PointerTo().ID("input").Dot(field.Name.Singular()))),
				),
				jen.Newline(),
			)
		}
	}

	bodyLines = append(bodyLines,
		jen.Return().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodPost"),
			jen.Litf("/%s", prn),
			jen.Qual("strings", "NewReader").Call(jen.ID("form").Dot("Encode").Call()),
		),
	)

	lines := []jen.Code{
		jen.Func().IDf("attach%sUpdateInputToRequest", sn).Params(jen.ID("input").PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn))).Params(jen.PointerTo().Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildTestServiceParseFormEncodedSomethingUpdateInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.ID("expected").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().Zero()
			}
			return jen.Null()
		}(),
		jen.Newline(),
		jen.ID("req").Assign().IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("expected")),
		jen.Newline(),
		jen.ID("actual").Assign().ID("s").Dot("service").Dotf("parseFormEncoded%sUpdateInput", sn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("badBody").Assign().AddressOf().Qual(proj.TestUtilsPackage(), "MockReadCloser").Values(),
		jen.ID("badBody").Dot("On").Call(jen.Lit("Read"), jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().ID("byte").Values())).Dot("Return").Call(
			jen.Zero(),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.Newline(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.Lit("/test"),
			jen.ID("badBody"),
		),
		jen.Newline(),
		jen.ID("actual").Assign().ID("s").Dot("service").Dotf("parseFormEncoded%sUpdateInput", sn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(jen.ID("t"), jen.ID("actual")),
	}

	thirdSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("exampleInput").Assign().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn)).Values(),
		jen.Newline(),
		jen.ID("req").Assign().IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("actual").Assign().ID("s").Dot("service").Dotf("parseFormEncoded%sUpdateInput", sn).Call(
			jen.ID("s").Dot("ctx"),
			jen.ID("req"),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("s").Dot("sessionCtxData")),
		),
		jen.Qual(constants.AssertionLibrary, "Nil").Call(jen.ID("t"), jen.ID("actual")),
	}

	lines := []jen.Code{
		jen.Func().IDf("TestService_parseFormEncoded%sUpdateInput", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(firstSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid form"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(secondSubtestLines...),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input attached to valid form"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(thirdSubtestLines...),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestServiceHandleSomethingUpdateRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	mockGetArgs := buildCallArgsForRetrievalTest(proj, typ)

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockGetArgs...,
		).Dot("Return").Call(
			jen.IDf("example%s", sn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			jen.Litf("Update%s", sn),
			jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			jen.IDf("example%s", sn),
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sUpdateRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	secondSubtestLines = append(secondSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	secondSubtestLines = append(secondSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.Newline(),
		jen.ID("s").Dot("service").Dot("sessionContextDataFetcher").Equals().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
			jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sUpdateRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("unauthorizedRedirectResponseCode"), jen.ID("res").Dot("Code")),
	)

	thirdSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("exampleInput").Assign().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn)).Values(),
		jen.Newline(),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sUpdateRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code")),
	}

	fourthSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	fourthSubtestLines = append(fourthSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	fourthSubtestLines = append(fourthSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockGetArgs...,
		).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sUpdateRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	fifthSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	fifthSubtestLines = append(fifthSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	fifthSubtestLines = append(fifthSubtestLines,
		jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockGetArgs...,
		).Dot("Return").Call(
			jen.IDf("example%s", sn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			jen.Litf("Update%s", sn),
			jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			jen.IDf("example%s", sn),
		).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().IDf("attach%sUpdateInputToRequest", sn).Call(jen.ID("exampleInput")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sUpdateRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestService_handle%sUpdateRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(firstSubtestLines...)),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching session context data"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(secondSubtestLines...)),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid input"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(thirdSubtestLines...)),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching data"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(fourthSubtestLines...)),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error updating data"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(fifthSubtestLines...)),
		),
		jen.Newline(),
	}

	return lines
}

func buildArchiveArgs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	out := []jen.Code{
		jen.Litf("Archive%s", sn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	}

	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
	}

	out = append(out,
		jen.IDf("example%s", sn).Dot("ID"),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID("s").Dot("sessionCtxData").Dot("ActiveAccountID")),
	)

	return out
}

func buildTestServiceHandleSomethingArchiveRequest(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()
	prn := typ.Name.PluralRouteName()

	mockArchiveArgs := buildArchiveArgs(proj, typ)
	mockListArgs := buildCallArgsForListRetrievalTest(proj, typ)

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	firstSubtestLines = append(firstSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	firstSubtestLines = append(firstSubtestLines,
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockArchiveArgs...,
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockListArgs...,
		).Dot("Return").Call(
			jen.IDf("example%sList", sn),
			jen.Nil(),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodDelete"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sArchiveRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	secondSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("s").Dot("service").Dot("sessionContextDataFetcher").Equals().Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")).Body(
			jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodDelete"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sArchiveRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("unauthorizedRedirectResponseCode"), jen.ID("res").Dot("Code")),
	}

	thirdSubtestLines := []jen.Code{jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	thirdSubtestLines = append(thirdSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	thirdSubtestLines = append(thirdSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockArchiveArgs...,
		).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodDelete"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sArchiveRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	fourthSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("s").Assign().ID("buildTestHelper").Call(jen.ID("t")),
		jen.Newline(),
	}

	fourthSubtestLines = append(fourthSubtestLines, buildTestIDFetchers(proj, typ, true)...)

	fourthSubtestLines = append(fourthSubtestLines,
		jen.Newline(),
		jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockArchiveArgs...,
		).Dot("Return").Call(jen.Nil()),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Callln(
			mockListArgs...,
		).Dot("Return").Call(
			jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()),
			jen.Qual("errors", "New").Call(jen.Lit("blah")),
		),
		jen.ID("s").Dot("service").Dot("dataStore").Equals().ID("mockDB"),
		jen.Newline(),
		jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
		jen.ID("req").Assign().ID("httptest").Dot("NewRequest").Call(
			jen.Qual("net/http", "MethodDelete"),
			jen.Litf("/%s", prn),
			jen.Nil(),
		),
		jen.Newline(),
		jen.ID("s").Dot("service").Dotf("handle%sArchiveRequest", sn).Call(
			jen.ID("res"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code")),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(jen.ID("t"), jen.ID("mockDB")),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestService_handle%sArchiveRequest", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					firstSubtestLines...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					secondSubtestLines...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error archiving %s", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					thirdSubtestLines...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error retrieving new list of %s", pcn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(fourthSubtestLines...),
			),
		),
		jen.Newline(),
	}

	return lines
}
