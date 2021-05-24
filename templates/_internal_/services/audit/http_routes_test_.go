package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAuditLogEntriesService_ListHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("exampleAuditLogEntryList").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call(),
					jen.ID("auditLogEntryManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
					jen.ID("auditLogEntryManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntries"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.ID("exampleAuditLogEntryList"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLogEntryManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryList").Values()),
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
						jen.ID("auditLogEntryManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
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
					jen.ID("auditLogEntryManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
					jen.ID("auditLogEntryManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntries"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("AuditLogEntryList")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLogEntryManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntryList").Values()),
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
						jen.ID("auditLogEntryManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("auditLogEntryManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
					jen.ID("auditLogEntryManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntries"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("QueryFilter").Values()),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("AuditLogEntryList")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLogEntryManager"),
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
						jen.ID("auditLogEntryManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAuditLogEntriesService_ReadHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("auditLogEntryManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
					jen.ID("auditLogEntryManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntry"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAuditLogEntry").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("helper").Dot("exampleAuditLogEntry"),
						jen.ID("nil"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLogEntryManager"),
					jen.ID("encoderDecoder").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "NewMockEncoderDecoder").Call(),
					jen.ID("encoderDecoder").Dot("On").Call(
						jen.Lit("RespondWithData"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "HTTPResponseWriterMatcher"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").ID("types").Dot("AuditLogEntry").Values()),
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
						jen.ID("auditLogEntryManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching session context data"),
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
				jen.Lit("with no audit log entries returned from datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("auditLogEntryManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
					jen.ID("auditLogEntryManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntry"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAuditLogEntry").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("AuditLogEntry")).Call(jen.ID("nil")),
						jen.Qual("database/sql", "ErrNoRows"),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLogEntryManager"),
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
						jen.ID("auditLogEntryManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error reading from datastore"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("auditLogEntryManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuditLogEntryDataManager").Values(),
					jen.ID("auditLogEntryManager").Dot("On").Call(
						jen.Lit("GetAuditLogEntry"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("helper").Dot("exampleAuditLogEntry").Dot("ID"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("AuditLogEntry")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("helper").Dot("service").Dot("auditLog").Op("=").ID("auditLogEntryManager"),
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
						jen.ID("auditLogEntryManager"),
						jen.ID("encoderDecoder"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
