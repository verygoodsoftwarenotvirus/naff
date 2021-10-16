package workers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preUpdatesWorkerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestProvidePreUpdatesWorker(proj)...)
	code.Add(buildTestPreUpdatesWorker_HandleMessage(proj)...)

	return code
}

func buildTestProvidePreUpdatesWorker(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvidePreUpdatesWorker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
						jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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

func buildTestPreUpdatesWorker_HandleMessage(proj *models.Project) []jen.Code {
	testCases := []jen.Code{
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
					jen.Return().List(jen.Nil(), jen.Nil()),
				),
				jen.Newline(),
				jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
	}

	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()
		scn := typ.Name.SingularCommonName()

		testCases = append(testCases,
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with %sDataType", sn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn), jen.ID(sn).MapAssign().Qualf(proj.FakeTypesPackage(), "BuildFake%s", sn).Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot(sn),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot(sn).Dot("ID"),
						jen.ID("body").Dot(sn),
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
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
				jen.Litf("with %sDataType with error updating %s", sn, scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn), jen.ID(sn).MapAssign().Qualf(proj.FakeTypesPackage(), "BuildFake%s", sn).Call(),
					),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot(sn),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.Nil(), jen.Nil()),
					),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
				jen.Litf("with %sDataType and error updating search index", sn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn), jen.ID(sn).MapAssign().Qualf(proj.FakeTypesPackage(), "BuildFake%s", sn).Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot(sn),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot(sn).Dot("ID"),
						jen.ID("body").Dot(sn),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
					jen.Newline(),
					jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
						jen.Return().List(jen.ID("searchIndexManager"), jen.Nil()),
					),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
				jen.Litf("with %sDataType and error publishing data change event", sn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
					jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
					jen.Newline(),
					jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
						jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn), jen.ID(sn).MapAssign().Qualf(proj.FakeTypesPackage(), "BuildFake%s", sn).Call()),
					jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.ID("dbManager").Dotf("%sDataManager", sn).Dot("On").Callln(
						jen.Litf("Update%s", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot(sn),
					).Dot("Return").Call(jen.Nil()),
					jen.Newline(),
					jen.ID("searchIndexManager").Op(":=").Op("&").Qual(proj.InternalSearchPackage("mock"), "IndexManager").Values(),
					jen.ID("searchIndexManager").Dot("On").Callln(
						jen.Lit("Index"),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("body").Dot(sn).Dot("ID"),
						jen.ID("body").Dot(sn),
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
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
			jen.Lit("with UserMembershipDataType"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
				jen.Newline(),
				jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
					jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "UserMembershipDataType")),
				jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
				jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
					jen.Return().List(jen.Nil(), jen.Nil())),
				jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
				jen.Newline(),
				jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
			jen.Lit("with WebhookDataType"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				jen.ID("client").Op(":=").Op("&").Qual("net/http", "Client").Values(),
				jen.Newline(),
				jen.ID("body").Op(":=").Op("&").Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
					jen.ID("DataType").MapAssign().Qual(proj.TypesPackage(), "WebhookDataType")),
				jen.List(jen.ID("examplePayload"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("body")),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.Newline(),
				jen.ID("dbManager").Op(":=").Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("searchIndexLocation").Op(":=").Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("t").Dot("Name").Call()),
				jen.ID("searchIndexProvider").Op(":=").Func().Params(jen.Qual("context", "Context"), jen.Qual(proj.InternalLoggingPackage(), "Logger"), jen.Op("*").Qual("net/http", "Client"), jen.Qual(proj.InternalSearchPackage(), "IndexPath"), jen.Qual(proj.InternalSearchPackage(), "IndexName"), jen.Op("...").String()).Params(jen.Qual(proj.InternalSearchPackage(), "IndexManager"), jen.ID("error")).Body(
					jen.Return().List(jen.Nil(), jen.Nil())),
				jen.ID("postArchivesPublisher").Op(":=").Op("&").Qual(proj.InternalMessageQueuePublishersPackage("mock"), "Publisher").Values(),
				jen.Newline(),
				jen.List(jen.ID("worker"), jen.ID("err")).Op(":=").ID("ProvidePreUpdatesWorker").Callln(
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
	)

	lines := []jen.Code{
		jen.Func().ID("TestPreUpdatesWorker_HandleMessage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			append([]jen.Code{
				jen.ID("T").Dot("Parallel").Call(),
				jen.Newline(),
			},
				testCases...,
			)...,
		),
		jen.Newline(),
	}

	return lines
}
