package workers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preUpdatesWorkerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvidePreUpdatesWorker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Op("&").Qual(proj.DatabasePackage(), "MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
					jen.ID("dbManager").Op(":=").Op("&").Qual(proj.DatabasePackage(), "MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
		jen.Func().ID("TestPreUpdatesWorker_HandleMessage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item").Dot("ID"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.Nil()),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
				jen.Lit("with ItemDataType with error updating item"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item").Dot("ID"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
				jen.Lit("with ItemDataType and error publishing data change event"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "ItemDataType"), jen.ID("Item").MapAssign().ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot("Item").Dot("ID"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Nil()),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil())),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
				jen.Lit("with UserMembershipDataType"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "UserMembershipDataType")),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
				jen.Lit("with WebhookDataType"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType")),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil())),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Call(
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
		),
		jen.Newline(),
	)

	return code
}
