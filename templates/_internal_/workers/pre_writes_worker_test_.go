package workers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preWritesWorkerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvidePreWritesWorker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Op("&").ID("database").Dot("MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error providing search index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Op("&").ID("database").Dot("MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.Index().ID("byte").Call(jen.Lit("} bad JSON lol")),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with ItemDataType"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("expectedItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.ID("expectedItem"),
						jen.ID("nil"),
					),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("expectedItem").Dot("ID"),
						jen.ID("expectedItem"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.ID("nil")),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with ItemDataType and error writing"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Item")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with ItemDataType and error updating search index"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("expectedItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.ID("expectedItem"),
						jen.ID("nil"),
					),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("expectedItem").Dot("ID"),
						jen.ID("expectedItem"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with ItemDataType and error publishing data change message"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItemDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("expectedItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("CreateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(
						jen.ID("expectedItem"),
						jen.ID("nil"),
					),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("expectedItem").Dot("ID"),
						jen.ID("expectedItem"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
						jen.ID("searchIndexManager"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with WebhookDataType"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("WebhookDataType"), jen.ID("Webhook").Op(":").ID("fakes").Dot("BuildFakeWebhookDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("expectedWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Call(
						jen.Lit("CreateWebhook"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Webhook"),
					).Dot("Return").Call(
						jen.ID("expectedWebhook"),
						jen.ID("nil"),
					),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.ID("nil")),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with WebhookDataType and error writing"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("WebhookDataType"), jen.ID("Webhook").Op(":").ID("fakes").Dot("BuildFakeWebhookDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Call(
						jen.Lit("CreateWebhook"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Webhook"),
					).Dot("Return").Call(
						jen.Parens(jen.Op("*").ID("types").Dot("Webhook")).Call(jen.ID("nil")),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with WebhookDataType and error publishing data change message"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("WebhookDataType"), jen.ID("Webhook").Op(":").ID("fakes").Dot("BuildFakeWebhookDatabaseCreationInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("expectedWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("WebhookDataManager").Dot("On").Call(
						jen.Lit("CreateWebhook"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Webhook"),
					).Dot("Return").Call(
						jen.ID("expectedWebhook"),
						jen.ID("nil"),
					),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with UserMembershipDataType"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("UserMembershipDataType"), jen.ID("UserMembership").Op(":").ID("fakes").Dot("BuildFakeAddUserToAccountInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("AccountUserMembershipDataManager").Dot("On").Call(
						jen.Lit("AddUserToAccount"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("UserMembership"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.ID("nil")),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with UserMembershipDataType and error writing"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("UserMembershipDataType"), jen.ID("UserMembership").Op(":").ID("fakes").Dot("BuildFakeAddUserToAccountInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("AccountUserMembershipDataManager").Dot("On").Call(
						jen.Lit("AddUserToAccount"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("dbManager"),
						jen.ID("postArchivesPublisher"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with UserMembershipDataType and error publishing data change message"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreWriteMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("UserMembershipDataType"), jen.ID("UserMembership").Op(":").ID("fakes").Dot("BuildFakeAddUserToAccountInput").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("AccountUserMembershipDataManager").Dot("On").Call(
						jen.Lit("AddUserToAccount"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreWritesWorker").Call(
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
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("worker").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.ID("examplePayload"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
