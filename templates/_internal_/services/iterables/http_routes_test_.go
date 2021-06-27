package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestParseBool(proj, typ)...)
	code.Add(buildTestSomethingsService_CreateHandler(proj, typ)...)
	code.Add(buildTestSomethingsService_ReadHandler(proj, typ)...)
	code.Add(buildTestSomethingsService_ExistenceHandler(proj, typ)...)
	code.Add(buildTestSomethingsService_ListHandler(proj, typ)...)
	if typ.SearchEnabled {
		code.Add(buildTestSomethingsService_SearchHandler(proj, typ)...)
	}
	code.Add(buildTestSomethingsService_UpdateHandler(proj, typ)...)
	code.Add(buildTestSomethingsService_ArchiveHandler(proj, typ)...)
	code.Add(buildTestAccountsService_AuditEntryHandler(proj, typ)...)

	return code
}

func buildTestParseBool(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestParseBool").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("expectations").Op(":=").Map(jen.ID("string")).ID("bool").Valuesln(jen.Lit("1").Op(":").ID("true"), jen.ID("t").Dot("Name").Call().Op(":").ID("false"), jen.Lit("true").Op(":").ID("true"), jen.Lit("troo").Op(":").ID("false"), jen.Lit("t").Op(":").ID("true"), jen.Lit("false").Op(":").ID("false")),
			jen.Newline(),
			jen.For(jen.List(jen.ID("input"), jen.ID("expected")).Op(":=").Range().ID("expectations")).Body(
				jen.ID("assert").Dot("Equal").Call(
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
		jen.Func().IDf("Test%sService_CreateHandler", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInput", sn)).Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Create%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					jen.ID("unitCounter").Dot("On").Call(jen.Lit("Increment"), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dotf("%sCounter", uvn).Op("=").ID("unitCounter"),
					jen.Newline(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
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
							).Dot("Return").Call(jen.ID("nil"))
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager")
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusCreated"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("unitCounter"),
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
				jen.Lit("without input attached"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("nil")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input attached"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInput", sn)).Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with error creating %s", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInput", sn)).Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Create%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Newline().ID("T").Dot("Run").Call(
						jen.Litf("with error indexing %s", scn),
						jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
							jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
								jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
								jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
							),
							jen.Newline(),
							jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInput", sn)).Call(),
							jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
								jen.ID("helper").Dot("ctx"),
								jen.ID("exampleCreationInput"),
							),
							jen.Newline(),
							jen.Var().ID("err").ID("error"),
							jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
								jen.ID("helper").Dot("ctx"),
								jen.Qual("net/http", "MethodPost"),
								jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
								jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
							jen.ID("require").Dot("NotNil").Call(
								jen.ID("t"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								jen.Litf("Create%s", sn),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Values()),
								jen.ID("helper").Dot("exampleUser").Dot("ID"),
							).Dot("Return").Call(
								jen.ID("helper").Dotf("example%s", sn),
								jen.ID("nil"),
							),
							jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
							jen.Newline(),
							jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
							jen.ID("unitCounter").Dot("On").Call(jen.Lit("Increment"), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")).Dot("Return").Call(),
							jen.ID("helper").Dot("service").Dotf("%sCounter", uvn).Op("=").ID("unitCounter"),
							jen.Newline(),
							jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.ID("indexManager").Dot("On").Callln(
								jen.Lit("Index"),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
								jen.ID("helper").Dotf("example%s", sn),
							).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
							jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
							jen.Newline(),
							jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
								jen.ID("helper").Dot("res"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.ID("assert").Dot("Equal").Call(
								jen.ID("t"),
								jen.Qual("net/http", "StatusCreated"),
								jen.ID("helper").Dot("res").Dot("Code"),
							),
							jen.Newline(),
							jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
								jen.ID("t"),
								jen.IDf("%sDataManager", uvn),
								jen.ID("unitCounter"),
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

func buildTestSomethingsService_ReadHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	getSomethingExpectationArgs := buildReadMockArgs(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ReadHandler", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(getSomethingExpectationArgs...).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), sn).Values()),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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

func buildTestSomethingsService_ExistenceHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	dbMockCallArgs := buildMockCallArgsForExistence(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ExistenceHandler", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						dbMockCallArgs...,
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						dbMockCallArgs...,
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						dbMockCallArgs...,
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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

func buildTestSomethingsService_ListHandler(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	getSomethingExpectedArgs := buildMockCallArgsForListHandler(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ListHandler", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("example%sList", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectedArgs...,
					).Dot("Return").Call(
						jen.IDf("example%sList", sn),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectedArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectedArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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

func buildTestSomethingsService_SearchHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	mockDBCallArgs := buildMockSearchArgs(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_SearchHandler", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("exampleQuery").Op(":=").Lit("whatever"),
			jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(123)),
			jen.IDf("example%sList", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
			jen.IDf("example%sIDs", sn).Op(":=").Index().ID("uint64").Values(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().IDf("example%sList", sn).Dot(pn)).Body(
				jen.IDf("example%sIDs", sn).Op("=").ID("append").Call(
					jen.IDf("example%sIDs", sn),
					jen.ID("x").Dot("ID"),
				)),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(
						jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Values(jen.ID("exampleQuery")),
						jen.Qual(proj.TypesPackage(), "LimitQueryKey").Op(":").Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.IDf("example%sIDs", sn),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						mockDBCallArgs...,
					).Dot("Return").Call(
						jen.IDf("example%sList", sn).Dot(pn),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), sn).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
						jen.ID("indexManager"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Values(jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Values(jen.ID("exampleQuery"))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().ID("uint64").Values(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("indexManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(
						jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Values(jen.ID("exampleQuery")),
						jen.Qual(proj.TypesPackage(), "LimitQueryKey").Op(":").Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.IDf("example%sIDs", sn),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						mockDBCallArgs...,
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), sn).Values(),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), sn).Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
						jen.ID("indexManager"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(
						jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Values(jen.ID("exampleQuery")),
						jen.Qual(proj.TypesPackage(), "LimitQueryKey").Op(":").Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.IDf("example%sIDs", sn),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						mockDBCallArgs...,
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), sn).Values(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("indexManager"),
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
		jen.Func().IDf("Test%sService_UpdateHandler", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInput", sn)).Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), sn).Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
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
							).Dot("Return").Call(jen.ID("nil"))
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager")
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)).Values(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached to context"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("nil")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with no such %s", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInput", sn)).Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInput", sn)).Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), sn)).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Litf("with error updating %s", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInput", sn)).Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Newline(),
					jen.Var().ID("err").ID("error"),
					jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("helper").Dot("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						getSomethingExpectationArgs...,
					).Dot("Return").Call(
						jen.ID("helper").Dotf("example%s", sn),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), sn).Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Newline().ID("T").Dot("Run").Call(
						jen.Lit("with error updating search index"),
						jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
							jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
								jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
								jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
							),
							jen.Newline(),
							jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sUpdateInput", sn)).Call(),
							jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
								jen.ID("helper").Dot("ctx"),
								jen.ID("exampleCreationInput"),
							),
							jen.Newline(),
							jen.Var().ID("err").ID("error"),
							jen.List(jen.ID("helper").Dot("req"), jen.ID("err")).Op("=").Qual("net/http", "NewRequestWithContext").Call(
								jen.ID("helper").Dot("ctx"),
								jen.Qual("net/http", "MethodPost"),
								jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
								jen.Qual("bytes", "NewReader").Call(jen.ID("jsonBytes")),
							),
							jen.ID("require").Dot("NoError").Call(
								jen.ID("t"),
								jen.ID("err"),
							),
							jen.ID("require").Dot("NotNil").Call(
								jen.ID("t"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								getSomethingExpectationArgs...,
							).Dot("Return").Call(
								jen.ID("helper").Dotf("example%s", sn),
								jen.ID("nil"),
							),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								jen.Litf("Update%s", sn),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), sn).Values()),
								jen.ID("helper").Dot("exampleUser").Dot("ID"),
								jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Values()),
							).Dot("Return").Call(jen.ID("nil")),
							jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
							jen.Newline(),
							jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
							jen.ID("indexManager").Dot("On").Callln(
								jen.Lit("Index"),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
								jen.ID("helper").Dotf("example%s", sn),
							).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
							jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
							jen.Newline(),
							jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
								jen.ID("helper").Dot("res"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.ID("assert").Dot("Equal").Call(
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

	lines = append(lines, jen.ID("helper").Dot("exampleUser").Dot("ID"))

	return lines
}

func buildTestSomethingsService_ArchiveHandler(proj *models.Project, typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()

	expectedCallArgs := buildMockCallArgsForArchive(proj, typ)

	return []jen.Code{
		jen.Func().IDf("Test%sService_ArchiveHandler", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedCallArgs...,
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("indexManager").Dot("On").Callln(
								jen.Lit("Delete"),
								jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
								jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
							).Dot("Return").Call(jen.ID("nil"))
						}
						return jen.Null()
					}(),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager")
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					jen.ID("unitCounter").Dot("On").Call(jen.Lit("Decrement"), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dotf("%sCounter", uvn).Op("=").ID("unitCounter"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.IDf("%sDataManager", uvn),
						jen.ID("unitCounter"),
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedCallArgs...,
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Lit("with error saving as archived"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						expectedCallArgs...,
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.Newline().ID("T").Dot("Run").Call(
						jen.Lit("with error removing from search index"),
						jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
							jen.Newline(),
							jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
							jen.IDf("%sDataManager", uvn).Dot("On").Callln(
								expectedCallArgs...,
							).Dot("Return").Call(jen.ID("nil")),
							jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
							jen.Newline(),
							func() jen.Code {
								if typ.SearchEnabled {
									return jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values()
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
									return jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager")
								}
								return jen.Null()
							}(),
							jen.Newline(),
							jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
							jen.ID("unitCounter").Dot("On").Call(jen.Lit("Decrement"), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")).Dot("Return").Call(),
							jen.ID("helper").Dot("service").Dotf("%sCounter", uvn).Op("=").ID("unitCounter"),
							jen.Newline(),
							jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
								jen.ID("helper").Dot("res"),
								jen.ID("helper").Dot("req"),
							),
							jen.Newline(),
							jen.ID("assert").Dot("Equal").Call(
								jen.ID("t"),
								jen.Qual("net/http", "StatusNoContent"),
								jen.ID("helper").Dot("res").Dot("Code"),
							),
							jen.Newline(),
							jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
								jen.ID("t"),
								jen.IDf("%sDataManager", uvn),
								jen.ID("unitCounter"),
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
		jen.Func().ID("TestAccountsService_AuditEntryHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("exampleAuditLogEntries").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("GetAuditLogEntriesFor%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntries"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("GetAuditLogEntriesFor%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.IDf("%sDataManager", uvn).Op(":=").Op("&").Qual(proj.TypesPackage("mock"), fmt.Sprintf("%sDataManager", sn)).Values(),
					jen.IDf("%sDataManager", uvn).Dot("On").Callln(
						jen.Litf("GetAuditLogEntriesFor%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dotf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dotf("%sDataManager", uvn).Op("=").IDf("%sDataManager", uvn),
					jen.Newline(),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Equal").Call(
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
