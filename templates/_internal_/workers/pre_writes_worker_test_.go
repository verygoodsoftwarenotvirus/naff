package workers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preWritesWorkerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvidePreWritesWorker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Op("&").Qual(proj.DatabasePackage(), "MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil()),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Op("&").Qual(proj.DatabasePackage(), "MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
	)

	code.Add(
		jen.Func().ID("TestPreWritesWorker_HandleMessage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with ItemDataType"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("expectedItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Callln(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.ID("expectedItem"),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("expectedItem").Dot("ID"),
						jen.ID("expectedItem"),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil()),
					),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with ItemDataType and error writing"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Callln(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "Item")).Call(jen.Nil()),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil()),
					),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with ItemDataType and error updating search index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("expectedItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Callln(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.ID("expectedItem"),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("expectedItem").Dot("ID"),
						jen.ID("expectedItem"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with ItemDataType and error publishing data change message"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("expectedItem").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Callln(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.ID("expectedItem"),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("expectedItem").Dot("ID"),
						jen.ID("expectedItem"),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with WebhookDataType"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType"), jen.ID("Webhook").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhookDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("expectedWebhook").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Callln(
						jen.Lit("CreateWebhook"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Webhook"),
					).Dot("Return").Call(
						jen.ID("expectedWebhook"),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with WebhookDataType and error writing"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType"), jen.ID("Webhook").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhookDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Callln(
						jen.Lit("CreateWebhook"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Webhook"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").Qual(proj.TypesPackage(), "Webhook")).Call(jen.Nil()),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil()),
					),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with WebhookDataType and error publishing data change message"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType"), jen.ID("Webhook").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeWebhookDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("expectedWebhook").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeWebhook").Call(),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Callln(
						jen.Lit("CreateWebhook"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Webhook"),
					).Dot("Return").Call(
						jen.ID("expectedWebhook"),
						jen.Nil(),
					),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil()),
					),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "UserMembershipDataType"), jen.ID("UserMembership").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeAddUserToAccountInput").Call(),
					),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("AccountUserMembershipDataManager").Dot("On").Callln(
						jen.Lit("AddUserToAccount"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("UserMembership"),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil()),
					),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with UserMembershipDataType and error writing"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "UserMembershipDataType"), jen.ID("UserMembership").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeAddUserToAccountInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("AccountUserMembershipDataManager").Dot("On").Callln(
						jen.Lit("AddUserToAccount"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("input").Op("*").Qual(proj.TypesPackage(), "AddUserToAccountInput")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil()),
					),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
				jen.Lit("with UserMembershipDataType and error publishing data change message"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "UserMembershipDataType"), jen.ID("UserMembership").MapAssign().Qual(proj.FakeTypesPackage(), "BuildFakeAddUserToAccountInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("AccountUserMembershipDataManager").Dot("On").Callln(
						jen.Lit("AddUserToAccount"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("input").Op("*").Qual(proj.TypesPackage(), "AddUserToAccountInput")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Callln(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Qual(constants.MockPkg, "MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).SingleLineBody(jen.Return().True())),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Callln(
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
		),
		jen.Newline(),
	)

	return code
}
