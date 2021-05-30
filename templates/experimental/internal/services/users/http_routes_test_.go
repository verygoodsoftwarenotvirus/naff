package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestService_validateCredentialChangeRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
					jen.ID("examplePassword").Op(":=").Lit("password"),
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
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("examplePassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleTOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("helper").Dot("service").Dot("validateCredentialChangeRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("examplePassword"),
						jen.ID("exampleTOTPToken"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("sc"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows found in database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
					jen.ID("examplePassword").Op(":=").Lit("password"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("helper").Dot("service").Dot("validateCredentialChangeRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("examplePassword"),
						jen.ID("exampleTOTPToken"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("sc"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
					jen.ID("examplePassword").Op(":=").Lit("password"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("helper").Dot("service").Dot("validateCredentialChangeRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("examplePassword"),
						jen.ID("exampleTOTPToken"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("sc"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error validating login"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
					jen.ID("examplePassword").Op(":=").Lit("password"),
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
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("examplePassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleTOTPToken"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("helper").Dot("service").Dot("validateCredentialChangeRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("examplePassword"),
						jen.ID("exampleTOTPToken"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("sc"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid login"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleTOTPToken").Op(":=").Lit("123456"),
					jen.ID("examplePassword").Op(":=").Lit("password"),
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
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("examplePassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleTOTPToken"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.List(jen.ID("actual"), jen.ID("sc")).Op(":=").ID("helper").Dot("service").Dot("validateCredentialChangeRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("examplePassword"),
						jen.ID("exampleTOTPToken"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("sc"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_UsernameSearchHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleUserList").Op(":=").ID("fakes").Dot("BuildFakeUserList").Call(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("SearchForUsersByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.ID("exampleUserList").Dot("Users"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("v").Op(":=").ID("helper").Dot("req").Dot("URL").Dot("Query").Call(),
					jen.ID("v").Dot("Set").Call(
						jen.ID("types").Dot("SearchQueryKey"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").ID("v").Dot("Encode").Call(),
					jen.ID("helper").Dot("service").Dot("UsernameSearchHandler").Call(
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
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("SearchForUsersByUsername"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("User").Valuesln(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("v").Op(":=").ID("helper").Dot("req").Dot("URL").Dot("Query").Call(),
					jen.ID("v").Dot("Set").Call(
						jen.ID("types").Dot("SearchQueryKey"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
					),
					jen.ID("helper").Dot("req").Dot("URL").Dot("RawQuery").Op("=").ID("v").Dot("Encode").Call(),
					jen.ID("helper").Dot("service").Dot("UsernameSearchHandler").Call(
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
		jen.Func().ID("TestService_ListHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleUserList").Op(":=").ID("fakes").Dot("BuildFakeUserList").Call(),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUsers"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.ID("exampleUserList"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("UserList").Valuesln()),
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
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUsers"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("UserList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
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
		jen.Func().ID("TestService_CreateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("HashPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("Password"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("db").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("db").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("CreateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("UserDataStoreCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("db"),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Increment"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("userCounter").Op("=").ID("unitCounter"),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("UserRegistrationInputContextKey"),
						jen.ID("exampleInput"),
					)),
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("true"),
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
						jen.ID("auth"),
						jen.ID("db"),
						jen.ID("unitCounter"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with user creation disabled"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("false"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusForbidden"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("true"),
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
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("UserRegistrationInput").Valuesln(),
					jen.ID("exampleInput").Dot("Password").Op("=").Lit("a"),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("true"),
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
				jen.Lit("with error validating password"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("exampleInput").Dot("Password").Op("=").Lit("a"),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleAccount").Dot("BelongsToUser").Op("=").ID("helper").Dot("exampleUser").Dot("ID"),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("UserRegistrationInputContextKey"),
						jen.ID("exampleInput"),
					)),
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("true"),
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
				jen.Lit("with error hashing password"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("HashPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("Password"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("UserRegistrationInputContextKey"),
						jen.ID("exampleInput"),
					)),
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("true"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error generating two factor secret"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("HashPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("Password"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("db").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("db").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("CreateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("UserDataStoreCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("db"),
					jen.ID("sg").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/random", "MockGenerator").Valuesln(),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateBase32EncodedString"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("totpSecretSize"),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("secretGenerator").Op("=").ID("sg"),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("UserRegistrationInputContextKey"),
						jen.ID("exampleInput"),
					)),
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("true"),
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
						jen.ID("auth"),
						jen.ID("db"),
						jen.ID("sg"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error creating entry in database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeUserRegistrationInputFromUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("HashPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("Password"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("db").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("db").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("CreateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("UserDataStoreCreationInput").Valuesln()),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("db"),
					jen.ID("helper").Dot("req").Op("=").ID("helper").Dot("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("helper").Dot("req").Dot("Context").Call(),
						jen.ID("types").Dot("UserRegistrationInputContextKey"),
						jen.ID("exampleInput"),
					)),
					jen.ID("helper").Dot("service").Dot("authSettings").Dot("EnableUserSignup").Op("=").ID("true"),
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
						jen.ID("auth"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_buildQRCode").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("actual").Op(":=").ID("helper").Dot("service").Dot("buildQRCode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("helper").Dot("exampleUser").Dot("Username"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("strings", "HasPrefix").Call(
							jen.ID("actual"),
							jen.ID("base64ImagePrefix"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_SelfHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("SelfHandler").Call(
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
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnauthorizedResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("SelfHandler").Call(
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
				jen.Lit("with no rows found"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeNotFoundResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("SelfHandler").Call(
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
						jen.ID("mockDB"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("SelfHandler").Call(
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
		jen.Func().ID("TestService_ReadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
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
						jen.ID("mockDB"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows found"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
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
		jen.Func().ID("TestService_TOTPSecretVerificationHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUserWithUnverifiedTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("MarkUserTwoFactorSecretAsVerified"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached to request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
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
				jen.Lit("without valid input attached"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("TOTPSecretVerificationInput").Valuesln(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUserWithUnverifiedTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUserWithUnverifiedTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
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
				jen.Lit("with invalid TOTP token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("exampleInput").Dot("TOTPToken").Op("=").Lit("000000"),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
						jen.Lit("GetUserWithUnverifiedTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("MarkUserTwoFactorSecretAsVerified"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
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
				jen.Lit("with secret already validated"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("og").Op(":=").ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn"),
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("og"),
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
						jen.Lit("GetUserWithUnverifiedTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAlreadyReported"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid code"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("exampleInput").Dot("TOTPToken").Op("=").Lit("INVALID"),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUserWithUnverifiedTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error verifying two factor secret"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretVerificationInputForUser").Call(jen.ID("helper").Dot("exampleUser")),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUserWithUnverifiedTwoFactorSecret"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("MarkUserTwoFactorSecretAsVerified"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("TOTPSecretVerificationHandler").Call(
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
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_NewTOTPSecretHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("NewTOTPSecretHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusAccepted"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached to request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("helper").Dot("service").Dot("NewTOTPSecretHandler").Call(
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
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("TOTPSecretRefreshInput").Valuesln(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("service").Dot("NewTOTPSecretHandler").Call(
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
				jen.Lit("with input attached but without user information"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("service").Dot("NewTOTPSecretHandler").Call(
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
				jen.Lit("with error validating login"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("NewTOTPSecretHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error generating secret"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("sg").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/random", "MockGenerator").Valuesln(),
					jen.ID("sg").Dot("On").Call(
						jen.Lit("GenerateBase32EncodedString"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("totpSecretSize"),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("secretGenerator").Op("=").ID("sg"),
					jen.ID("helper").Dot("service").Dot("NewTOTPSecretHandler").Call(
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
						jen.ID("auth"),
						jen.ID("sg"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating user in database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeTOTPSecretRefreshInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("NewTOTPSecretHandler").Call(
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
						jen.ID("auth"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_UpdatePasswordHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUserPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("string")),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("HashPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("NewPassword"),
					).Dot("Return").Call(
						jen.Lit("blah"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
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
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without input attached to request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
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
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PasswordUpdateInput").Valuesln(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
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
				jen.Lit("with input but without user info"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
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
				jen.Lit("with error validating login"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUserPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("string")),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("false"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid password"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("exampleInput").Dot("NewPassword").Op("=").Lit("a"),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error hashing password"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUserPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("string")),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("HashPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("NewPassword"),
					).Dot("Return").Call(
						jen.Lit("blah"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
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
						jen.ID("auth"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePasswordUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
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
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUserPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("string")),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("auth").Op(":=").Op("&").ID("authentication").Dot("MockAuthenticator").Valuesln(),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("ValidateLogin"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("HashedPassword"),
						jen.ID("exampleInput").Dot("CurrentPassword"),
						jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret"),
						jen.ID("exampleInput").Dot("TOTPToken"),
					).Dot("Return").Call(
						jen.ID("true"),
						jen.ID("nil"),
					),
					jen.ID("auth").Dot("On").Call(
						jen.Lit("HashPassword"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput").Dot("NewPassword"),
					).Dot("Return").Call(
						jen.Lit("blah"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("authenticator").Op("=").ID("auth"),
					jen.ID("helper").Dot("service").Dot("UpdatePasswordHandler").Call(
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
						jen.ID("auth"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_AvatarUploadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("returnImage").Op(":=").Op("&").ID("images").Dot("Image").Valuesln(),
					jen.ID("ip").Op(":=").Op("&").ID("images").Dot("MockImageUploadProcessor").Valuesln(),
					jen.ID("ip").Dot("On").Call(
						jen.Lit("Process"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
						jen.Lit("avatar"),
					).Dot("Return").Call(
						jen.ID("returnImage"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("imageUploadProcessor").Op("=").ID("ip"),
					jen.ID("um").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads/mock", "UploadManager").Valuesln(),
					jen.ID("um").Dot("On").Call(
						jen.Lit("SaveFile"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("avatar_%d"),
							jen.ID("helper").Dot("exampleUser").Dot("ID"),
						),
						jen.ID("returnImage").Dot("Data"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("uploadManager").Op("=").ID("um"),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("helper").Dot("service").Dot("AvatarUploadHandler").Call(
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
						jen.ID("ip"),
						jen.ID("um"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("without session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnauthorizedResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AvatarUploadHandler").Call(
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
				jen.Lit("with error fetching user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("User")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AvatarUploadHandler").Call(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error processing image"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("ip").Op(":=").Op("&").ID("images").Dot("MockImageUploadProcessor").Valuesln(),
					jen.ID("ip").Dot("On").Call(
						jen.Lit("Process"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
						jen.Lit("avatar"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("images").Dot("Image")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("imageUploadProcessor").Op("=").ID("ip"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeInvalidInputResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AvatarUploadHandler").Call(
						jen.ID("helper").Dot("res"),
						jen.ID("helper").Dot("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mockDB"),
						jen.ID("ip"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error saving file"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
					jen.ID("returnImage").Op(":=").Op("&").ID("images").Dot("Image").Valuesln(),
					jen.ID("ip").Op(":=").Op("&").ID("images").Dot("MockImageUploadProcessor").Valuesln(),
					jen.ID("ip").Dot("On").Call(
						jen.Lit("Process"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
						jen.Lit("avatar"),
					).Dot("Return").Call(
						jen.ID("returnImage"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("imageUploadProcessor").Op("=").ID("ip"),
					jen.ID("um").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads/mock", "UploadManager").Valuesln(),
					jen.ID("um").Dot("On").Call(
						jen.Lit("SaveFile"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("avatar_%d"),
							jen.ID("helper").Dot("exampleUser").Dot("ID"),
						),
						jen.ID("returnImage").Dot("Data"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("uploadManager").Op("=").ID("um"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AvatarUploadHandler").Call(
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
						jen.ID("ip"),
						jen.ID("um"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("GetUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleUser"),
						jen.ID("nil"),
					),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("UpdateUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("User").Valuesln()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("returnImage").Op(":=").Op("&").ID("images").Dot("Image").Valuesln(),
					jen.ID("ip").Op(":=").Op("&").ID("images").Dot("MockImageUploadProcessor").Valuesln(),
					jen.ID("ip").Dot("On").Call(
						jen.Lit("Process"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
						jen.Lit("avatar"),
					).Dot("Return").Call(
						jen.ID("returnImage"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("imageUploadProcessor").Op("=").ID("ip"),
					jen.ID("um").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/uploads/mock", "UploadManager").Valuesln(),
					jen.ID("um").Dot("On").Call(
						jen.Lit("SaveFile"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("avatar_%d"),
							jen.ID("helper").Dot("exampleUser").Dot("ID"),
						),
						jen.ID("returnImage").Dot("Data"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("uploadManager").Op("=").ID("um"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeUnspecifiedInternalServerErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("AvatarUploadHandler").Call(
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
						jen.ID("ip"),
						jen.ID("um"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_ArchiveHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("ArchiveUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Decrement"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("userCounter").Op("=").ID("unitCounter"),
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
						jen.ID("mockDB"),
						jen.ID("unitCounter"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no results in the database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("ArchiveUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("mockDB").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
						jen.Lit("ArchiveUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("mockDB"),
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
						jen.ID("mockDB"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_AuditEntryHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleAuditLogEntries").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntries"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
						jen.ID("userDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with sql.ErrNoRows"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
						jen.ID("userDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("userDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UserDataManager").Valuesln(),
					jen.ID("userDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForUser"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("userDataManager").Op("=").ID("userDataManager"),
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
						jen.ID("userDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
