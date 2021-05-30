package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAPIClientsService_ListHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("exampleAPIClientList").Op(":=").ID("fakes").Dot("BuildFakeAPIClientList").Call(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.ID("exampleAPIClientList"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("APIClientList").Valuesln()),
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
						jen.ID("mockDB"),
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
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no results returned from datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClientList")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("APIClientList").Valuesln()),
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
						jen.ID("mockDB"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving clients from the datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClients"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClientList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAPIClientsService_CreateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("a").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("a").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("a"),
					jen.ID("sg").Op(":=").Op("&").ID("random").Dot("MockGenerator").Valuesln(),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateBase64EncodedString"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("clientIDSize"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
						jen.ID("nil"),
					),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateRawBytes"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("clientSecretSize"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("secretGenerator").Op("=").ID("sg"),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("CreateAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleInput"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("uc").Dot("On").Call(
						jen.Lit("Increment"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("apiClientCounter").Op("=").ID("uc"),
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
						jen.ID("mockDB"),
						jen.ID("a"),
						jen.ID("sg"),
						jen.ID("uc"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached to request"),
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
				jen.Lit("with invalid input attached to request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("helper").Dot("service").Dot("cfg").Dot("minimumPasswordLength").Op("=").Qual("math", "MaxUint8"),
					jen.ID("helper").Dot("exampleInput").Dot("Password").Op("=").Lit(""),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
				jen.Lit("with error retrieving user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid credentials"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("CreateAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleInput"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("a").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("a").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("a"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
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
						jen.ID("mockDB"),
						jen.ID("a"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid password"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("CreateAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleInput"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("a").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("a").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("a"),
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
						jen.ID("mockDB"),
						jen.ID("a"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error generating client ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("a").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("a").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("a"),
					jen.ID("sg").Op(":=").Op("&").ID("random").Dot("MockGenerator").Valuesln(),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateBase64EncodedString"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("clientIDSize"),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("secretGenerator").Op("=").ID("sg"),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("CreateAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleInput"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
						jen.ID("a"),
						jen.ID("sg"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error generating client secret"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("a").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("a").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("a"),
					jen.ID("sg").Op(":=").Op("&").ID("random").Dot("MockGenerator").Valuesln(),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateBase64EncodedString"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("clientIDSize"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
						jen.ID("nil"),
					),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateRawBytes"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("clientSecretSize"),
					).Dot("Return").Call(
						jen.Index().ID("byte").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("secretGenerator").Op("=").ID("sg"),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("CreateAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleInput"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
						jen.ID("a"),
						jen.ID("sg"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error creating API Client in data store"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleInput"),
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
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("a").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("a").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("helper").Dot("exampleInput").Dot("Password"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("helper").Dot("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("a"),
					jen.ID("sg").Op(":=").Op("&").ID("random").Dot("MockGenerator").Valuesln(),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateBase64EncodedString"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("clientIDSize"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientID"),
						jen.ID("nil"),
					),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateRawBytes"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("clientSecretSize"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient").Dot("ClientSecret"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("secretGenerator").Op("=").ID("sg"),
					jen.ID("mockDB").Dot("APIClientDataManager").Dot("On").Call(
						jen.Lit("CreateAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleInput"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClient")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
						jen.ID("a"),
						jen.ID("sg"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAPIClientsService_ReadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAPIClient"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("APIClient").Valuesln()),
					).Dot("Return").Call(),
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
						jen.ID("apiClientDataManager"),
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
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no such API client in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClient")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
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
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAPIClientByDatabaseID"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("APIClient")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
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
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAPIClientsService_ArchiveHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("ArchiveAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Decrement"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("apiClientCounter").Op("=").ID("unitCounter"),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNoContent"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("unitCounter"),
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
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no such API client in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("ArchiveAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
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
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("ArchiveAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("ArchiveHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
						jen.Lit("expected %d in status response, got %d"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("apiClientDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAPIClientsService_AuditEntryHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("exampleAuditLogEntries").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntries"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
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
						jen.ID("apiClientDataManager"),
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
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
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
						jen.ID("apiClientDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("apiClientDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "APIClientDataManager").Valuesln(),
					jen.ID("apiClientDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForAPIClient"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAPIClient").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("apiClientDataManager").Op("=").ID("apiClientDataManager"),
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
						jen.ID("apiClientDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
