package items

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestParseBool").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("expectations").Op(":=").Map(jen.ID("string")).ID("bool").Valuesln(jen.Lit("1").Op(":").ID("true"), jen.ID("t").Dot("Name").Call().Op(":").ID("false"), jen.Lit("true").Op(":").ID("true"), jen.Lit("troo").Op(":").ID("false"), jen.Lit("t").Op(":").ID("true"), jen.Lit("false").Op(":").ID("false")),
			jen.For(jen.List(jen.ID("input"), jen.ID("expected")).Op(":=").Range().ID("expectations")).Body(
				jen.ID("assert").Dot("Equal").Call(
					jen.ID("t"),
					jen.ID("expected"),
					jen.ID("parseBool").Call(jen.ID("input")),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemsService_CreateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("CreateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("ItemCreationInput").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Increment"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusCreated"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("unitCounter"),
						jen.ID("indexManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input attached"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").Op("&").ID("types").Dot("ItemCreationInput").Valuesln(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error creating item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("CreateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("ItemCreationInput").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Item")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error indexing item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("CreateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("ItemCreationInput").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Increment"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusCreated"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("unitCounter"),
						jen.ID("indexManager"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemsService_ReadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Item").Valuesln()),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no such item in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Item")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Item")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemsService_ExistenceHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("ItemExists"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no result in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("ItemExists"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error checking database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("ItemExists"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ExistenceHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemsService_ListHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItems"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.ID("exampleItemList"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("ItemList").Valuesln()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItems"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("ItemList")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("ItemList").Valuesln()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving items from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItems"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("ItemList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemsService_SearchHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("exampleQuery").Op(":=").Lit("whatever"),
			jen.ID("exampleLimit").Op(":=").ID("uint8").Call(jen.Lit(123)),
			jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
			jen.ID("exampleItemIDs").Op(":=").Index().ID("uint64").Valuesln(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
				jen.ID("exampleItemIDs").Op("=").ID("append").Call(
					jen.ID("exampleItemIDs"),
					jen.ID("x").Dot("ID"),
				)),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.ID("types").Dot("SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery")), jen.ID("types").Dot("LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Search"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleItemIDs"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItemsWithIDs"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("exampleLimit"),
						jen.ID("exampleItemIDs"),
					).Dot("Return").Call(
						jen.ID("exampleItemList").Dot("Items"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("Item").Valuesln()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("indexManager"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error conducting search"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.ID("types").Dot("SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery"))).Dot("Encode").Call(),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Search"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().ID("uint64").Valuesln(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("indexManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.ID("types").Dot("SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery")), jen.ID("types").Dot("LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Search"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleItemIDs"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItemsWithIDs"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("exampleLimit"),
						jen.ID("exampleItemIDs"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("Item").Valuesln(),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("Item").Valuesln()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("indexManager"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").Qual("net/url", "Values").Valuesln(jen.ID("types").Dot("SearchQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("exampleQuery")), jen.ID("types").Dot("LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("exampleLimit"))))).Dot("Encode").Call(),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Search"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleQuery"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleItemIDs"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItemsWithIDs"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("exampleLimit"),
						jen.ID("exampleItemIDs"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("Item").Valuesln(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("SearchHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("indexManager"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemsService_UpdateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Item").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("FieldChangeSummary").Valuesln()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").Op("&").ID("types").Dot("ItemUpdateInput").Valuesln(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached to context"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no such item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Item")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving item from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Item")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Item").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("FieldChangeSummary").Valuesln()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating search index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeItemUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleCreationInput"),
					),
					jen.Var().Defs(
						jen.ID("err").ID("error"),
					),
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
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleItem"),
						jen.ID("nil"),
					),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Item").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("FieldChangeSummary").Valuesln()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleItem"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestItemsService_ArchiveHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("ArchiveItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Delete"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Decrement"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
						jen.ID("unitCounter"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no such item in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("ArchiveItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error saving as archived"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("ArchiveItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error removing from search index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("ArchiveItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("indexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Valuesln(),
					jen.ID("indexManager").Dot("On").Call(
						jen.Lit("Delete"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("search").Op("=").ID("indexManager"),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Decrement"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("itemCounter").Op("=").ID("unitCounter"),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("indexManager"),
						jen.ID("unitCounter"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAccountsService_AuditEntryHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("exampleAuditLogEntries").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntries"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Valuesln()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.Lit("unauthenticated"),
						jen.Qual("net/http", "StatusUnauthorized"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with sql.ErrNoRows"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("itemDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "ItemDataManager").Valuesln(),
					jen.ID("itemDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleItem").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("itemDataManager").Op("=").ID("itemDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AuditEntryHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("itemDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
