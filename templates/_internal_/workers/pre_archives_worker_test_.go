package workers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preArchivesWorkerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestProvidePreArchivesWorker(proj)...)
	code.Add(buildTestPreArchivesWorker_HandleMessage(proj)...)

	return code
}

func buildTestProvidePreArchivesWorker(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvidePreArchivesWorker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Assign().Op("&").Qual(proj.DatabasePackage(), "MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("client"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexLocation"),
						jen.ID("searchIndexProvider"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error providing search index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Assign().Op("&").Qual(proj.DatabasePackage(), "MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("client"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexLocation"),
						jen.ID("searchIndexProvider"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestPreArchivesWorker_HandleMessage(proj *models.Project) []jen.Code {
	testCases := []jen.Code{}

	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()

		testCases = append(testCases,
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with %sDataType", sn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn)),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Archive%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						func() jen.Code {
							if typ.BelongsToStruct != nil {
								return jen.ID("body").Dotf("%sID", typ.BelongsToStruct.Singular())
							}
							return jen.Null()
						}(),
						jen.ID("body").Dotf("%sID", sn),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("body").Dot("AttributableToAccountID")),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexManager").Assign().Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Delete"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dotf("%sID", sn),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("client"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexLocation"),
						jen.ID("searchIndexProvider"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("worker"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexManager"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with %sDataType and error archiving", sn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn)),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Archive%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						func() jen.Code {
							if typ.BelongsToStruct != nil {
								return jen.ID("body").Dotf("%sID", typ.BelongsToStruct.Singular())
							}
							return jen.Null()
						}(),
						jen.ID("body").Dotf("%sID", sn),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("body").Dot("AttributableToAccountID")),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("client"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexLocation"),
						jen.ID("searchIndexProvider"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("worker"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with %sDataType and error removing from search index", sn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn)),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Archive%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						func() jen.Code {
							if typ.BelongsToStruct != nil {
								return jen.ID("body").Dotf("%sID", typ.BelongsToStruct.Singular())
							}
							return jen.Null()
						}(),
						jen.ID("body").Dotf("%sID", sn),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("body").Dot("AttributableToAccountID")),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexManager").Assign().Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Delete"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dotf("%sID", sn),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.Newline(),
					jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("client"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexLocation"),
						jen.ID("searchIndexProvider"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("worker"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexManager"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with %sDataType and error publishing post-archive message", sn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn)),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Archive%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						func() jen.Code {
							if typ.BelongsToStruct != nil {
								return jen.ID("body").Dotf("%sID", typ.BelongsToStruct.Singular())
							}
							return jen.Null()
						}(),
						jen.ID("body").Dotf("%sID", sn),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("body").Dot("AttributableToAccountID")),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("searchIndexManager").Assign().Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Delete"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dotf("%sID", sn),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("client"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexLocation"),
						jen.ID("searchIndexProvider"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("worker"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexManager"),
					),
				),
			),
			jen.Newline(),
		)
	}

	testCases = append(testCases,
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with WebhookDataType"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
				jen.Newline(),
				jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
					jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType")),
				jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Callln(
					jen.Lit("ArchiveWebhook"),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					jen.ID("body").Dot("WebhookID"),
					jen.ID("body").Dot("AttributableToAccountID"),
				).Dot("Return").Call(jen.Nil()),
				jen.Newline(),
				jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
				jen.ID("postArchivesPublisher").Dot("On").Callln(
					jen.Lit("Publish"),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
				).Dot("Return").Call(jen.Nil()),
				jen.Newline(),
				jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
				jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
					jen.Return().List(jen.Nil(), jen.Nil())),
				jen.Newline(),
				jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
					jen.ID("ctx"),
					jen.ID("logger"),
					jen.ID("client"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
					jen.ID("searchIndexLocation"),
					jen.ID("searchIndexProvider"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID("t"),
					jen.ID("worker"),
				),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("worker").Dot("HandleMessage").Call(
						jen.ID("ctx"),
						jen.ID("examplePayload"),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with WebhookDataType and error archiving"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
				jen.Newline(),
				jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
					jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType")),
				jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Callln(
					jen.Lit("ArchiveWebhook"),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					jen.ID("body").Dot("WebhookID"),
					jen.ID("body").Dot("AttributableToAccountID"),
				).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
				jen.ID("searchIndexManager").Assign().Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
				jen.Newline(),
				jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
				jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
					jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
				jen.Newline(),
				jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
					jen.ID("ctx"),
					jen.ID("logger"),
					jen.ID("client"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
					jen.ID("searchIndexLocation"),
					jen.ID("searchIndexProvider"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID("t"),
					jen.ID("worker"),
				),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("worker").Dot("HandleMessage").Call(
						jen.ID("ctx"),
						jen.ID("examplePayload"),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
					jen.ID("searchIndexManager"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with WebhookDataType and error publishing post-archive message"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
				jen.Newline(),
				jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
					jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType")),
				jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Callln(
					jen.Lit("ArchiveWebhook"),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					jen.ID("body").Dot("WebhookID"),
					jen.ID("body").Dot("AttributableToAccountID"),
				).Dot("Return").Call(jen.Nil()),
				jen.Newline(),
				jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
				jen.ID("postArchivesPublisher").Dot("On").Callln(
					jen.Lit("Publish"),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
				).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
				jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
					jen.Return().List(jen.Nil(), jen.Nil())),
				jen.Newline(),
				jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
					jen.ID("ctx"),
					jen.ID("logger"),
					jen.ID("client"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
					jen.ID("searchIndexLocation"),
					jen.ID("searchIndexProvider"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID("t"),
					jen.ID("worker"),
				),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("worker").Dot("HandleMessage").Call(
						jen.ID("ctx"),
						jen.ID("examplePayload"),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with UserMembershipDataType"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
				jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
				jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
				jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
					jen.Return().List(jen.Nil(), jen.Nil())),
				jen.Newline(),
				jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
					jen.ID("ctx"),
					jen.ID("logger"),
					jen.ID("client"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
					jen.ID("searchIndexLocation"),
					jen.ID("searchIndexProvider"),
				),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID("t"),
					jen.ID("worker"),
				),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("body").Assign().Op("&").Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
					jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "UserMembershipDataType")),
				jen.List(jen.ID("examplePayload"), jen.ID("err")).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("body")),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("worker").Dot("HandleMessage").Call(
						jen.ID("ctx"),
						jen.ID("examplePayload"),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("dbManager"),
					jen.ID("postArchivesPublisher"),
				),
			),
		),
	)

	lines := []jen.Code{
		jen.Func().ID("TestPreArchivesWorker_HandleMessage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			append([]jen.Code{
				jen.ID("T").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("with invalid input"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("ctx").Assign().Qual("context", "Background").Call(),
						jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
						jen.ID("client").Assign().Op("&").Qual("net/http", "Client").Values(),
						jen.ID("dbManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
						jen.ID("postArchivesPublisher").Assign().Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
						jen.ID("searchIndexLocation").Assign().Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
						jen.ID("searchIndexProvider").Assign().Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
							jen.Return().List(jen.Nil(), jen.Nil())),
						jen.Newline(),
						jen.List(jen.ID("worker"), jen.ID("err")).Assign().ID("ProvidePreArchivesWorker").Callln(
							jen.ID("ctx"),
							jen.ID("logger"),
							jen.ID("client"),
							jen.ID("dbManager"),
							jen.ID("postArchivesPublisher"),
							jen.ID("searchIndexLocation"),
							jen.ID("searchIndexProvider"),
						),
						jen.ID("require").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("worker"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("worker").Dot("HandleMessage").Call(
								jen.ID("ctx"),
								jen.Index().ID("byte").Call(jen.Lit("} bad JSON lol")),
							),
						),
						jen.Newline(),
						jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
							jen.ID("t"),
							jen.ID("dbManager"),
							jen.ID("postArchivesPublisher"),
						),
					),
				),
				jen.Newline(),
			},
				testCases...,
			)...,
		),
		jen.Newline(),
	}

	return lines
}
