package webhooks

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestWebhooksService_CreateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInput").Call(),
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
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Increment"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("webhookCounter").Op("=").ID("unitCounter"),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("CreateWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("WebhookCreationInput").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleWebhook"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("unitCounter"),
						jen.ID("wd"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("sessionContextDataFetcher").Op("=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BrokenSessionContextDataFetcher"),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInput").Call(),
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
						jen.Qual("net/http", "StatusUnauthorized"),
						jen.ID("helper").Dot("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error decoding request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("DecodeRequest"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("WebhookCreationInput").Valuesln()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("")),
						jen.Qual("net/http", "StatusBadRequest"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
					jen.ID("helper").Dot("service").Dot("CreateHandler").Call(
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
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid content attached to request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").Op("&").ID("types").Dot("WebhookCreationInput").Valuesln(),
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
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Increment"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("webhookCounter").Op("=").ID("unitCounter"),
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
				jen.Lit("with error creating webhook"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleCreationInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInput").Call(),
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
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("CreateWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("WebhookCreationInput").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Webhook")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhooksService_ListHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleWebhookList").Op(":=").ID("fakes").Dot("BuildFakeWebhookList").Call(),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhooks"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.ID("exampleWebhookList"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("WebhookList").Valuesln()),
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
						jen.ID("wd"),
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
					jen.ID("helper").Dot("service").Dot("ListHandler").Call(
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
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhooks"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("WebhookList")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("WebhookList").Valuesln()),
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
						jen.ID("wd"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching webhooks from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhooks"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Valuesln()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("WebhookList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhooksService_ReadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleWebhook"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Webhook").Valuesln()),
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
						jen.ID("wd"),
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
					jen.ID("helper").Dot("service").Dot("ReadHandler").Call(
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
				jen.Lit("with no such webhook in database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Webhook")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching webhook from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Webhook")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhooksService_UpdateHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleUpdateInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleUpdateInput"),
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
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleWebhook"),
						jen.ID("nil"),
					),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("UpdateWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Webhook").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("FieldChangeSummary").Valuesln()),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
				jen.Lit("with error decoding request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("DecodeRequest"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPRequestMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("WebhookUpdateInput").Valuesln()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("EncodeErrorResponse"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("")),
						jen.Qual("net/http", "StatusBadRequest"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoderDecoder"),
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
				jen.Lit("with invalid content attached to request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleUpdateInput").Op(":=").Op("&").ID("types").Dot("WebhookUpdateInput").Valuesln(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleUpdateInput"),
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
				jen.Lit("with no rows fetching webhook"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleUpdateInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleUpdateInput"),
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
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Webhook")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching webhook"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleUpdateInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleUpdateInput"),
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
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Webhook")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows updating webhook"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleUpdateInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleUpdateInput"),
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
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleWebhook"),
						jen.ID("nil"),
					),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("UpdateWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Webhook").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("FieldChangeSummary").Valuesln()),
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
					jen.ID("helper").Dot("service").Dot("UpdateHandler").Call(
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
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error updating webhook"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("helper").Dot("service").Dot("encoderDecoder").Op("=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNoopLogger").Call(),
						jen.ID("encoding").Dot("ContentTypeJSON"),
					),
					jen.ID("exampleUpdateInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookUpdateInput").Call(),
					jen.ID("jsonBytes").Op(":=").ID("helper").Dot("service").Dot("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleUpdateInput"),
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
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("GetWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleWebhook"),
						jen.ID("nil"),
					),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("UpdateWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("Webhook").Valuesln()),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().Op("*").ID("types").Dot("FieldChangeSummary").Valuesln()),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhooksService_ArchiveHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("unitCounter").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/observability/metrics/mock", "UnitCounter").Valuesln(),
					jen.ID("unitCounter").Dot("On").Call(
						jen.Lit("Decrement"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(),
					jen.ID("helper").Dot("service").Dot("webhookCounter").Op("=").ID("unitCounter"),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("ArchiveWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("unitCounter"),
						jen.ID("wd"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error retrieving session context data"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
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
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no webhook in database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("ArchiveWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("wd").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("wd").Dot("On").Call(
						jen.Lit("ArchiveWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
						jen.ID("helper").Dot("exampleAccount").Dot("ID"),
						jen.ID("helper").Dot("exampleUser").Dot("ID"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("wd"),
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
						jen.ID("wd"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestWebhooksService_AuditEntryHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("exampleAuditLogEntries").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.ID("webhookDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("webhookDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntries"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("webhookDataManager"),
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
						jen.ID("webhookDataManager"),
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
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("webhookDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("webhookDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("webhookDataManager"),
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
						jen.ID("webhookDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("newTestHelper").Call(jen.ID("t")),
					jen.ID("webhookDataManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "WebhookDataManager").Valuesln(),
					jen.ID("webhookDataManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntriesForWebhook"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleWebhook").Dot("ID"),
					).Dot("Return").Call(
						jen.Index().Op("*").ID("types").Dot("AuditLogEntry").Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("webhookDataManager").Op("=").ID("webhookDataManager"),
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
						jen.ID("webhookDataManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
