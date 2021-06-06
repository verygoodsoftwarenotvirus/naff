package iterables

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildTestParseBool(proj, typ)...)
	code.Add(buildTestItemsService_CreateHandler(proj, typ)...)
	code.Add(buildTestItemsService_ReadHandler(proj, typ)...)
	code.Add(buildTestItemsService_ExistenceHandler(proj, typ)...)
	code.Add(buildTestItemsService_ListHandler(proj, typ)...)
	if typ.SearchEnabled {
		code.Add(buildTestItemsService_SearchHandler(proj, typ)...)
	}
	code.Add(buildTestItemsService_UpdateHandler(proj, typ)...)
	code.Add(buildTestItemsService_ArchiveHandler(proj, typ)...)
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

func buildTestItemsService_CreateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestItemsService_CreateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "ItemCreationInput").Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					jen.ID("unitCounter").Dot("On").Call(jen.Lit("Increment"), jen.Qual(proj.TestUtilsPackage(), "ContextMatcher")).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusCreated"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("unitCounter"),
						jen.ID("indexManager"),
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
					jen.ID("exampleCreationInput").Op(":=").Op("&").Qual(proj.TypesPackage(), "ItemCreationInput").Values(),
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
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInput").Call(),
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
				jen.Lit("with error creating item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "ItemCreationInput").Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "Item")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
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
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error indexing item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemCreationInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "ItemCreationInput").Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					jen.Newline(),
					jen.ID("unitCounter").Dot("On").Callln(
						jen.Lit("Increment"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusCreated"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("unitCounter"),
						jen.ID("indexManager"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestItemsService_ReadHandler(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestItemsService_ReadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "Item").Values()),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
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
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
				jen.Lit("with no such item in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "Item")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "Item")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestItemsService_ExistenceHandler(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestItemsService_ExistenceHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("ItemExists"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
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
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("ItemExists"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("ItemExists"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestItemsService_ListHandler(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestItemsService_ListHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("exampleItemList").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemList").Call(),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItems"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleItemList"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "ItemList").Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
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
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItems"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "ItemList")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "ItemList").Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving items from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItems"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "ItemList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestItemsService_SearchHandler(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestItemsService_SearchHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("exampleQuery").Op(":=").Lit("whatever"),
			jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(123)),
			jen.ID("exampleItemList").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemList").Call(),
			jen.ID("exampleItemIDs").Op(":=").Index().ID("uint64").Values(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
				jen.ID("exampleItemIDs").Op("=").ID("append").Call(
					jen.ID("exampleItemIDs"),
					jen.ID("x").Dot("ID"),
				)),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery")), jen.Qual(proj.TypesPackage(), "LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleItemIDs"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItemsWithIDs"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("exampleLimit"),
						jen.ID("exampleItemIDs"),
					).Dot("Return").Call(
						jen.ID("exampleItemList").Dot("Items"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "Item").Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
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
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery"))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
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
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery")), jen.Qual(proj.TypesPackage(), "LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleItemIDs"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItemsWithIDs"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("exampleLimit"),
						jen.ID("exampleItemIDs"),
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), "Item").Values(),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("RespondWithData"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "Item").Values()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
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
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.Qual(proj.TypesPackage(), "SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery")), jen.Qual(proj.TypesPackage(), "LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Search"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleItemIDs"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItemsWithIDs"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("exampleLimit"),
						jen.ID("exampleItemIDs"),
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), "Item").Values(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("indexManager"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestItemsService_UpdateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestItemsService_UpdateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemUpdateInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("UpdateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "Item").Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
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
					jen.ID("exampleCreationInput").Op(":=").Op("&").Qual(proj.TypesPackage(), "ItemUpdateInput").Values(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no such item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemUpdateInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "Item")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving item from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemUpdateInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "Item")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").Qual(proj.EncodingPackage(), "ProvideServerEncoderDecoder").Call(
						jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
					),
					jen.Newline(),
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemUpdateInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("UpdateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "Item").Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Values()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
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
					jen.ID("exampleCreationInput").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItemUpdateInput").Call(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("UpdateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Op("&").Qual(proj.TypesPackage(), "Item").Values()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.Qual(constants.MockPkg, "IsType").Call(jen.Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary").Values()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestItemsService_ArchiveHandler(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestItemsService_ArchiveHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("ArchiveItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Delete"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					jen.Newline(),
					jen.ID("unitCounter").Dot("On").Callln(
						jen.Lit("Decrement"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
						jen.ID("unitCounter"),
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
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
					jen.ID("encoderDecoder").Dot("On").Callln(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(proj.TestUtilsPackage(), "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual(proj.TestUtilsPackage(), "BrokenSessionContextDataFetcher"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
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
				jen.Lit("with no such item in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("ArchiveItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("ArchiveItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error removing from search index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("ArchiveItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.Newline(),
					jen.ID("indexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.Newline(),
					jen.ID("indexManager").Dot("On").Callln(
						jen.Lit("Delete"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.Newline(),
					jen.ID("unitCounter").Op(":=").Op("&").Qual(proj.MetricsPackage("mock"), "UnitCounter").Values(),
					jen.Newline(),
					jen.ID("unitCounter").Dot("On").Callln(
						jen.Lit("Decrement"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.Newline(),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
						jen.ID("unitCounter"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestAccountsService_AuditEntryHandler(proj *models.Project, typ models.DataType) []jen.Code {
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
					jen.ID("exampleAuditLogEntries").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.Newline(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetAuditLogEntriesForItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntries"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
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
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetAuditLogEntriesForItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "ItemDataManager").Values(),
					jen.Newline(),
					jen.ID("itemDataManager").Dot("On").Callln(
						jen.Lit("GetAuditLogEntriesForItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual(proj.EncodingPackage("mock"), "NewMockEncoderDecoder").Call(),
					jen.Newline(),
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
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}
