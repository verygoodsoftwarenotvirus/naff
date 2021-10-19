package iterables

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildReadMockArgs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Litf("Get%s", sn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID("helper").Dotf("example%s", pt.Name.Singular()).Dot("ID"))
	}
	lines = append(lines, jen.ID("helper").Dotf("example%s", sn).Dot("ID"))

	if (typ.BelongsToAccount && typ.RestrictedToAccountMembers) || typ.RestrictedToAccountAtSomeLevel(proj) {
		lines = append(lines, jen.ID("helper").Dot("exampleAccount").Dot("ID"))
	}

	return lines
}

func buildMockCallArgsForExistence(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Litf("%sExists", sn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID("helper").Dotf("example%s", pt.Name.Singular()).Dot("ID"))
	}
	lines = append(lines, jen.ID("helper").Dotf("example%s", sn).Dot("ID"))

	if (typ.BelongsToAccount && typ.RestrictedToAccountMembers) || typ.RestrictedToAccountAtSomeLevel(proj) {
		lines = append(lines, jen.ID("helper").Dot("exampleAccount").Dot("ID"))
	}

	return lines
}

func buildMockCallArgsForListHandler(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	lines := []jen.Code{
		jen.Litf("Get%s", pn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	}

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lines = append(lines, jen.ID("helper").Dotf("example%s", pt.Name.Singular()).Dot("ID"))
	}

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		lines = append(lines, jen.ID("helper").Dot("exampleAccount").Dot("ID"))
	}

	lines = append(lines, jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), "QueryFilter").Values()))

	return lines
}

func buildMockSearchArgs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	lines := []jen.Code{
		jen.Litf("Get%sWithIDs", pn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	}

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID("helper").Dotf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"))
	}

	if typ.BelongsToAccount {
		lines = append(lines, jen.ID("helper").Dot("exampleAccount").Dot("ID"))
	}

	lines = append(lines, jen.ID("exampleLimit"), jen.IDf("example%sIDs", sn))

	return lines
}

func buildMockCallArgsForArchiveExistenceCheck(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Litf("%sExists", sn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	}

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID("helper").Dotf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"))
	}
	lines = append(lines, jen.ID("helper").Dotf("example%s", sn).Dot("ID"))

	if typ.BelongsToAccount {
		lines = append(lines, jen.ID("helper").Dot("exampleAccount").Dot("ID"))
	}

	return lines
}

func buildMockCallArgsForArchive(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Litf("Archive%s", sn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	}

	if typ.BelongsToStruct != nil {
		lines = append(lines, jen.ID("helper").Dotf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"))
	}
	lines = append(lines, jen.ID("helper").Dotf("example%s", sn).Dot("ID"))

	if typ.BelongsToAccount {
		lines = append(lines, jen.ID("helper").Dot("exampleAccount").Dot("ID"))
	}

	return lines
}

func httpRoutesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestParseBool(proj, typ)...)
	code.Add(buildTestSomethingsService_CreateHandler(proj, typ)...)
	code.Add(buildTestSomethingsService_ReadHandler(proj, typ)...)
	code.Add(buildTestSomethingsService_ListHandler(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildTestSomethingsService_SearchHandler(proj, typ)...)
	}

	code.Add(buildTestSomethingsService_UpdateHandler(proj, typ)...)
	code.Add(buildTestSomethingsService_ArchiveHandler(proj, typ)...)

	return code
}

func buildTestParseBool(_ *models.Project, _ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestParseBool").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("expectations").Assign().Map(jen.String()).ID("bool").Valuesln(
				jen.Lit("1").MapAssign().True(), jen.ID("t").Dot("Name").Call().MapAssign().False(), jen.Lit("true").MapAssign().True(), jen.Lit("troo").MapAssign().False(), jen.Lit("t").MapAssign().True(), jen.Lit("false").MapAssign().False()),
			jen.Newline(),
			jen.For(jen.List(jen.ID("input"), jen.ID("expected")).Assign().Range().ID("expectations")).Body(
				jen.Qual(constants.AssertionLibrary, "Equal").Call(
					jen.ID("t"),
					jen.ID("expected"),
					jen.ID("parseBool").Call(jen.ID("input")),
				)),
		),
		jen.Newline(),
	}
}

func buildTestSomethingsService_CreateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("Test%sService_CreateHandler", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("mockEventProducer").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("mockEventProducer").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Qual(proj.TestUtilsPackage(), "PreWriteMessageMatcher")),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("helper").Dot("service").Dot("preWritesPublisher").Equals().ID("mockEventProducer"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockEventProducer"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard without async"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.ID("helper").Dot("service").Dot("async").Equals().False(),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Create%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sDatabaseCreationInput", sn)).Values()),
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()),
					utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
						jen.ID("helper").Dotf("example%s", sn),
					).Dot("Return").Call(jen.Nil())),
					utils.ConditionalCode(typ.SearchEnabled, jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusCreated"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.Nil()),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input attached"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", sn)).Values(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error publishing event"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("mockEventProducer").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("mockEventProducer").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Qual(proj.TestUtilsPackage(), "PreWriteMessageMatcher")),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("preWritesPublisher").Equals().ID("mockEventProducer"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockEventProducer"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("without async and error creating %s", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.ID("helper").Dot("service").Dot("async").Equals().False(),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Create%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sDatabaseCreationInput", sn)).Values()),
					).Dot("Return").Call(
						jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()),
					utils.ConditionalCode(typ.SearchEnabled, jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
					),
				),
			),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Newline().ID("T").Dot("Run").Call(
						jen.Litf("without async and error indexing %s", scn),
						jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
							jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
								jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
								jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
							),
							jen.ID("helper").Dot("service").Dot("async").Equals().False(),
							jen.Newline(),
							jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationRequestInput", sn)).Call(),
							jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
								jen.ID("helper").Dot("ctx"),
								jen.ID("exampleCreationInput"),
							),
							jen.Newline(),
							jen.Var().ID("err").ID("error"),
							jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
								jen.ID("helper").Dot("ctx"),
								jen.Qual("net/http", "MethodPost"),
								jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
								jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
							),
							jen.Qual(constants.MustAssertPkg, "NoError").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
							jen.Qual(constants.MustAssertPkg, "NotNil").Call(
								jen.ID("t"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								jen.Litf("Create%s", sn),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sDatabaseCreationInput", sn)).Values()),
							).Dot("Return").Call(
								jen.ID("helper").Dotf("example%s", sn),
								jen.Nil(),
							),
							jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
							jen.Newline(),
							utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()),
							utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager").Dot("On").Callln(
								jen.Lit("Index"),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
								jen.ID("helper").Dotf("example%s", sn),
							).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))),
							utils.ConditionalCode(typ.SearchEnabled, jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")),
							jen.Newline(),
							jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
								jen.ID("helper").Dot("res"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.Qual(constants.AssertionLibrary, "Equal").Call(
								jen.ID("t"),
								jen.Qual("net/http", "StatusCreated"),
								jen.ID("helper").Dot("res").Dot("Code"),
							),
							jen.Newline(),
							jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
								jen.ID("t"),
								jen.IDf("%sDataManager", uvn),
								utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
							),
						),
					)
				}
				return jen.Null()
			}(),
		),
		jen.Newline(),
	}
}

func buildTestSomethingsService_ReadHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	getSomethingExpectationArgs := buildReadMockArgs(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ReadHandler", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(getSomethingExpectationArgs...).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), sn).Values()),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with no such %s in the database", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching from database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSomethingsService_ExistenceHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	dbMockCallArgs := buildMockCallArgsForExistence(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ExistenceHandler", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						dbMockCallArgs...,
					).Dot("Return").Call(
						jen.True(),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no result in the database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						dbMockCallArgs...,
					).Dot("Return").Call(
						jen.False(),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error checking database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						dbMockCallArgs...,
					).Dot("Return").Call(
						jen.False(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSomethingsService_ListHandler(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	getSomethingExpectedArgs := buildMockCallArgsForListHandler(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ListHandler", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectedArgs...,
					).Dot("Return").Call(
						jen.IDf("example%sList", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectedArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error retrieving %s from database", pcn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectedArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.Nil()),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSomethingsService_SearchHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	mockDBCallArgs := buildMockSearchArgs(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_SearchHandler", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("exampleQuery").Assign().Lit("whatever"),
			jen.ID("exampleLimit").Assign().ID("uint8").Call(jen.Lit(123)),
			jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
			jen.IDf("example%sIDs", sn).Assign().Index().String().Values(),
			jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
				jen.IDf("example%sIDs", sn).Equals().ID("append").Call(
					jen.IDf("example%sIDs", sn),
					jen.ID("x").Dot("ID"),
				)),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Equals().Qual("net/url", "Values").Valuesln(
						jen.Qual(proj.TypesPackage(), "SearchQueryKey").MapAssign().Index().String().Values(jen.ID("exampleQuery")),
						jen.Qual(proj.TypesPackage(), "LimitQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.IDf("example%sIDs", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						mockDBCallArgs...,
					).Dot("Return").Call(
						jen.IDf("example%sList", sn).Dot(pn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().PointerTo().Qual(proj.TypesPackage(), sn).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error conducting search"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Equals().Qual("net/url", "Values").Values(jen.Qual(proj.TypesPackage(), "SearchQueryKey").MapAssign().Index().String().Values(jen.ID("exampleQuery"))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().String().Values(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager"),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Equals().Qual("net/url", "Values").Valuesln(
						jen.Qual(proj.TypesPackage(), "SearchQueryKey").MapAssign().Index().String().Values(jen.ID("exampleQuery")),
						jen.Qual(proj.TypesPackage(), "LimitQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.IDf("example%sIDs", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						mockDBCallArgs...,
					).Dot("Return").Call(
						jen.Index().PointerTo().Qual(proj.TypesPackage(), sn).Values(),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().PointerTo().Qual(proj.TypesPackage(), sn).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving from database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Equals().Qual("net/url", "Values").Valuesln(
						jen.Qual(proj.TypesPackage(), "SearchQueryKey").MapAssign().Index().String().Values(jen.ID("exampleQuery")),
						jen.Qual(proj.TypesPackage(), "LimitQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.IDf("example%sIDs", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						mockDBCallArgs...,
					).Dot("Return").Call(
						jen.Index().PointerTo().Qual(proj.TypesPackage(), sn).Values(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSomethingsService_UpdateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	getSomethingExpectationArgs := buildReadMockArgs(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_UpdateHandler", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("mockEventProducer").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("mockEventProducer").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Qual(proj.TestUtilsPackage(), "PreUpdateMessageMatcher")),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("helper").Dot("service").Dot("preUpdatesPublisher").Equals().ID("mockEventProducer"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("mockEventProducer"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard without async"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.ID("helper").Dot("service").Dot("async").Equals().False(),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.Nil(),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), sn).Values()),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Dot("On").Callln(
								jen.Lit("Index"),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
								jen.ID("helper").Dotf("example%s", sn),
							).Dot("Return").Call(jen.Nil())
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn)).Values(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached to context"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.Nil()),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with no such %s", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error retrieving %s from database", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Call(jen.Nil()),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("without async and with issue fetching %s", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.ID("helper").Dot("service").Dot("async").Equals().False(),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.Nil(),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), sn).Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()),
					utils.ConditionalCode(typ.SearchEnabled, jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")),
					utils.ConditionalCode(typ.SearchEnabled, jen.Newline()),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						utils.ConditionalCode(typ.SearchEnabled, jen.ID("indexManager")),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error publishing to message queue"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInput", sn)).Call(),
					jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Qual(constants.MustAssertPkg, "NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("mockEventProducer").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("mockEventProducer").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Qual(proj.TestUtilsPackage(), "PreUpdateMessageMatcher")),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("preUpdatesPublisher").Equals().ID("mockEventProducer"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("mockEventProducer"),
					),
				),
			),
			jen.Newline(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Newline().ID("T").Dot("Run").Call(
						jen.Lit("without async and error indexing for search"),
						jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
							jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
								jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
								jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
							),
							jen.ID("helper").Dot("service").Dot("async").Equals().False(),
							jen.Newline(),
							jen.ID("exampleCreationInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateRequestInput", sn)).Call(),
							jen.ID("jsonBytes").Assign().ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
								jen.ID("helper").Dot("ctx"),
								jen.ID("exampleCreationInput"),
							),
							jen.Newline(),
							jen.Var().ID("err").ID("error"),
							jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Equals().Qual("net/http", "NewRequestWithContext").Call(
								jen.ID("helper").Dot("ctx"),
								jen.Qual("net/http", "MethodPost"),
								jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
								jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
							),
							jen.Qual(constants.MustAssertPkg, "NoError").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
							jen.Qual(constants.MustAssertPkg, "NotNil").Call(
								jen.ID("t"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								getSomethingExpectationArgs...,
							).Dot("Return").Call(
								jen.ID("helper").Dotf("example%s", sn),
								jen.Nil(),
							),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								jen.Litf("Update%s", sn),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.Qual(constants.MockPkg, "IsType").Call(jen.AddressOf().Qual(proj.TypesPackage(), sn).Values()),
							).Dot("Return").Call(jen.Nil()),
							jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
							jen.Newline(),
							jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.ID("indexManager").Dot("On").Callln(
								jen.Lit("Index"),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
								jen.ID("helper").Dotf("example%s", sn),
							).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
							jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager"),
							jen.Newline(),
							jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
								jen.ID("helper").Dot("res"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.Qual(constants.AssertionLibrary, "Equal").Call(
								jen.ID("t"),
								jen.Qual("net/http", "StatusOK"),
								jen.ID("helper").Dot("res").Dot("Code"),
								jen.Lit("expected %d in status response, got %d"),
								jen.Qual("net/http", "StatusOK"),
								jen.ID("helper").Dot("res").Dot("Code"),
							),
							jen.Newline(),
							jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
								jen.ID("t"),
								jen.IDf("%sDataManager", uvn),
								jen.ID("indexManager"),
							),
						),
					)
				}
				return jen.Null()
			}(),
		),
		jen.Newline(),
	}
}

func buildTestSomethingsService_ArchiveHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	expectedExistenceCallArgs := buildMockCallArgsForExistence(proj, typ)
	expectedCallArgs := buildMockCallArgsForArchive(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ArchiveHandler", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedExistenceCallArgs...,
					).Dot("Return").Call(jen.True(), jen.Nil()),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("mockEventProducer").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("mockEventProducer").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Qual(proj.TestUtilsPackage(), "PreArchiveMessageMatcher")),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("helper").Dot("service").Dot("preArchivesPublisher").Equals().ID("mockEventProducer"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("mockEventProducer"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard without async"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("async").Equals().False(),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedCallArgs...,
					).Dot("Return").Call(jen.Nil()),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Dot("On").Callln(
								jen.Lit("Delete"),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
							).Dot("Return").Call(jen.Nil())
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						func() jen.Code {
							if typ.SearchEnabled {
								return jen.ID("indexManager")
							}
							return jen.Null()
						}(),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with no such %s in the database", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedExistenceCallArgs...,
					).Dot("Return").Call(jen.False(), jen.Nil()),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error checking for item in database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedExistenceCallArgs...,
					).Dot("Return").Call(jen.False(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error publishing to message queue"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedExistenceCallArgs...,
					).Dot("Return").Call(jen.True(), jen.Nil()),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("mockEventProducer").Assign().AddressOf().Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("mockEventProducer").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Qual(proj.TestUtilsPackage(), "PreArchiveMessageMatcher")),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("preArchivesPublisher").Equals().ID("mockEventProducer"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("mockEventProducer"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without async and no rows returned"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("async").Equals().False(),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedCallArgs...,
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						func() jen.Code {
							if typ.SearchEnabled {
								return jen.ID("indexManager")
							}
							return jen.Null()
						}(),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without async and with error archiving item"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("async").Equals().False(),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedCallArgs...,
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						func() jen.Code {
							if typ.SearchEnabled {
								return jen.ID("indexManager")
							}
							return jen.Null()
						}(),
					),
				),
			),
			jen.Newline(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Newline().ID("T").Dot("Run").Call(
						jen.Lit("without async and issue deleting from search"),
						jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
							jen.ID("helper").Dot("service").Dot("async").Equals().False(),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								expectedCallArgs...,
							).Dot("Return").Call(jen.Nil()),
							jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
							jen.Newline(),
							func() jen.Code {
								if typ.SearchEnabled {
									return jen.ID("indexManager").Assign().AddressOf().Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
								}
								return jen.Null()
							}(),
							func() jen.Code {
								if typ.SearchEnabled {
									return jen.ID("indexManager").Dot("On").Callln(
										jen.Lit("Delete"),
										jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
										jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
									).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
								}
								return jen.Null()
							}(),
							func() jen.Code {
								if typ.SearchEnabled {
									return jen.ID("helper").Dot("service").Dot("search").Equals().ID("indexManager")
								}
								return jen.Null()
							}(),
							jen.Newline(),
							jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
								jen.ID("helper").Dot("res"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.Qual(constants.AssertionLibrary, "Equal").Call(
								jen.ID("t"),
								jen.Qual("net/http", "StatusNoContent"),
								jen.ID("helper").Dot("res").Dot("Code"),
							),
							jen.Newline(),
							jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
								jen.ID("t"),
								jen.IDf("%sDataManager", uvn),
								func() jen.Code {
									if typ.SearchEnabled {
										return jen.ID("indexManager")
									}
									return jen.Null()
								}(),
							),
						),
					)
				}
				return jen.Null()
			}(),
		),
		jen.Newline(),
	}
}

func buildTestAccountsService_AuditEntryHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()

	return []jen.Code{
		jen.Func().ID("TestAccountsService_AuditEntryHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleAuditLogEntries").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("GetAuditLogEntriesFor%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntries"),
						jen.Nil(),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry").Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Equals().Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with sql.ErrNoRows"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("GetAuditLogEntriesFor%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry").Call(jen.Nil()),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Assign().ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Assign().AddressOf().Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("GetAuditLogEntriesFor%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry").Call(jen.Nil()),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Equals().IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Assign().Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Equals().ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}
