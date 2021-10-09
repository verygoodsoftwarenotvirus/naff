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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Op("&").ID("database").Dot("MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").Op("&").ID("database").Dot("MockDatabase").Values(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah")))),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreUpdateMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("body").Dot("Item").Dot("ID"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
							jen.Return().ID("true"))),
					).Dot("Return").Call(jen.ID("nil")),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreUpdateMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreUpdateMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("body").Dot("Item").Dot("ID"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.ID("nil"))),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreUpdateMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("ItemDataType"), jen.ID("Item").Op(":").ID("fakes").Dot("BuildFakeItem").Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("dbManager").Dot("ItemDataManager").Dot("On").Call(
						jen.Lit("UpdateItem"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/search/mock", "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Call(
						jen.Lit("Index"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("body").Dot("Item").Dot("ID"),
						jen.ID("body").Dot("Item"),
					).Dot("Return").Call(jen.ID("nil")),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
					jen.ID("postArchivesPublisher").Dot("On").Call(
						jen.Lit("Publish"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("mock").Dot("MatchedBy").Call(jen.Func().Params(jen.ID("message").Op("*").ID("types").Dot("DataChangeMessage")).Params(jen.ID("bool")).Body(
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreUpdateMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("UserMembershipDataType")),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
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
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.ID("body").Op(":=").Op("&").ID("types").Dot("PreUpdateMessage").Valuesln(jen.ID("DataType").Op(":").ID("types").Dot("WebhookDataType")),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("dbManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
					jen.ID("searchIndexLocation").Op(":=").ID("search").Dot("IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.ID("logging").Dot("Logger"), jen.Op("*").Qual("net/http", "Client"), jen.ID("search").Dot("IndexPath"), jen.ID("search").Dot("IndexName"), jen.Op("...").ID("string")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("nil"))),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/messagequeue/publishers/mock", "Publisher").Values(),
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
