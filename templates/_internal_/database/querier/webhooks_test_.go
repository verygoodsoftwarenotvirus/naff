package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildMockRowsFromWebhooks").Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").ID("uint64"), jen.ID("webhooks").Op("...").Op("*").ID("types").Dot("Webhook")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("columns").Op(":=").ID("querybuilding").Dot("WebhooksTableColumns"),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Op("=").ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("columns")),
			jen.For(jen.List(jen.ID("_"), jen.ID("w")).Op(":=").Range().ID("webhooks")).Body(
				jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(jen.ID("w").Dot("ID"), jen.ID("w").Dot("ExternalID"), jen.ID("w").Dot("Name"), jen.ID("w").Dot("ContentType"), jen.ID("w").Dot("URL"), jen.ID("w").Dot("Method"), jen.Qual("strings", "Join").Call(
					jen.ID("w").Dot("Events"),
					jen.ID("querybuilding").Dot("WebhooksTableEventsSeparator"),
				), jen.Qual("strings", "Join").Call(
					jen.ID("w").Dot("DataTypes"),
					jen.ID("querybuilding").Dot("WebhooksTableDataTypesSeparator"),
				), jen.Qual("strings", "Join").Call(
					jen.ID("w").Dot("Topics"),
					jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
				), jen.ID("w").Dot("CreatedOn"), jen.ID("w").Dot("LastUpdatedOn"), jen.ID("w").Dot("ArchivedOn"), jen.ID("w").Dot("BelongsToAccount")),
				jen.If(jen.ID("includeCounts")).Body(
					jen.ID("rowValues").Op("=").ID("append").Call(
						jen.ID("rowValues"),
						jen.ID("filteredCount"),
						jen.ID("len").Call(jen.ID("webhooks")),
					)),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildErroneousMockRowsFromWebhooks").Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").ID("uint64"), jen.ID("webhooks").Op("...").Op("*").ID("types").Dot("Webhook")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("columns").Op(":=").ID("querybuilding").Dot("WebhooksTableColumns"),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Op("=").ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("columns")),
			jen.For(jen.List(jen.ID("_"), jen.ID("w")).Op(":=").Range().ID("webhooks")).Body(
				jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(jen.Qual("strings", "Join").Call(
					jen.ID("w").Dot("Events"),
					jen.ID("querybuilding").Dot("WebhooksTableEventsSeparator"),
				), jen.Qual("strings", "Join").Call(
					jen.ID("w").Dot("DataTypes"),
					jen.ID("querybuilding").Dot("WebhooksTableDataTypesSeparator"),
				), jen.Qual("strings", "Join").Call(
					jen.ID("w").Dot("Topics"),
					jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
				), jen.ID("w").Dot("ID"), jen.ID("w").Dot("ExternalID"), jen.ID("w").Dot("Name"), jen.ID("w").Dot("ContentType"), jen.ID("w").Dot("URL"), jen.ID("w").Dot("Method"), jen.ID("w").Dot("CreatedOn"), jen.ID("w").Dot("LastUpdatedOn"), jen.ID("w").Dot("ArchivedOn"), jen.ID("w").Dot("BelongsToAccount")),
				jen.If(jen.ID("includeCounts")).Body(
					jen.ID("rowValues").Op("=").ID("append").Call(
						jen.ID("rowValues"),
						jen.ID("filteredCount"),
						jen.ID("len").Call(jen.ID("webhooks")),
					)),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ScanWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("surfaces row errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs row closing errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromWebhooks").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleWebhook"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleWebhook"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid webhook ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRowsFromWebhooks").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleWebhook"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetAllWebhooksCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllWebhooksCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expected"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAllWebhooksCount").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleWebhookList").Op(":=").ID("fakes").Dot("BuildFakeWebhookList").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromWebhooks").Call(
						jen.ID("true"),
						jen.ID("exampleWebhookList").Dot("FilteredCount"),
						jen.ID("exampleWebhookList").Dot("Webhooks").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleWebhookList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil filter"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleWebhookList").Op(":=").ID("fakes").Dot("BuildFakeWebhookList").Call(),
					jen.ID("exampleWebhookList").Dot("Page").Op("=").Lit(0),
					jen.ID("exampleWebhookList").Dot("Limit").Op("=").Lit(0),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromWebhooks").Call(
						jen.ID("true"),
						jen.ID("exampleWebhookList").Dot("FilteredCount"),
						jen.ID("exampleWebhookList").Dot("Webhooks").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleWebhookList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous database response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetAllWebhooks").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Webhook")),
					jen.ID("doneChan").Op(":=").ID("make").Call(
						jen.Chan().ID("bool"),
						jen.Lit(1),
					),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleWebhookList").Op(":=").ID("fakes").Dot("BuildFakeWebhookList").Call(),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllWebhooksCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromWebhooks").Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.ID("exampleWebhookList").Dot("Webhooks").Op("..."),
					)),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("GetAllWebhooks").Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.ID("stillQuerying").Op(":=").ID("true"),
					jen.For(jen.ID("stillQuerying")).Body(
						jen.Select().Body(
							jen.Case(jen.ID("batch").Op(":=").Op("<-").ID("results")).Body(
								jen.ID("assert").Dot("NotEmpty").Call(
									jen.ID("t"),
									jen.ID("batch"),
								), jen.ID("doneChan").ReceiveFromChannel().ID("true")),
							jen.Case(jen.Op("<-").Qual("time", "After").Call(jen.Qual("time", "Second"))).Body(
								jen.ID("t").Dot("FailNow").Call()),
							jen.Case(jen.Op("<-").ID("doneChan")).Body(
								jen.ID("stillQuerying").Op("=").ID("false")),
						)),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil results channel"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("GetAllWebhooks").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.Lit(0),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching initial count"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Webhook")),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllWebhooksCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("err").Op(":=").ID("c").Dot("GetAllWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.ID("exampleBatchSize"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with no rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Webhook")),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllWebhooksCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("GetAllWebhooks").Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error querying database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Webhook")),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllWebhooksCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("GetAllWebhooks").Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("results").Op(":=").ID("make").Call(jen.Chan().Index().Op("*").ID("types").Dot("Webhook")),
					jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Op(":=").ID("uint16").Call(jen.Lit(1000)),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAllWebhooksCountQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetBatchOfWebhooksQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("GetAllWebhooks").Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_CreateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("exampleWebhook").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleWebhook").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleWebhook"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("exampleWebhook").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleWebhook").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing creation query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleWebhook").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("exampleWebhook").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleWebhook").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("exampleWebhook").Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeWebhookCreationInputFromWebhook").Call(jen.ID("exampleWebhook")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildCreateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().ID("exampleWebhook").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_UpdateWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook"),
							jen.Lit(0),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildUpdateWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ArchiveWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("actual").Op(":=").ID("c").Dot("ArchiveWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid webhook ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveWebhook").Call(
							jen.ID("ctx"),
							jen.Lit(0),
							jen.ID("exampleWebhook").Dot("BelongsToAccount"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook").Dot("ID"),
							jen.Lit(0),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook").Dot("ID"),
							jen.ID("exampleWebhook").Dot("BelongsToAccount"),
							jen.Lit(0),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildArchiveWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.ID("exampleWebhook").Dot("ID"))),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveWebhook").Call(
							jen.ID("ctx"),
							jen.ID("exampleWebhook").Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetAuditLogEntriesForWebhook").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.ID("auditLogEntries").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAuditLogEntriesForWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromAuditLogEntries").Call(
						jen.ID("false"),
						jen.ID("auditLogEntries").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("auditLogEntries"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid webhook ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForWebhook").Call(
						jen.ID("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAuditLogEntriesForWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleWebhook").Op(":=").ID("fakes").Dot("BuildFakeWebhook").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Op(":=").ID("fakes").Dot("BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Op(":=").ID("database").Dot("BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dot("WebhookSQLQueryBuilder").Dot("On").Call(
						jen.Lit("BuildGetAuditLogEntriesForWebhookQuery"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("exampleWebhook").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogEntriesForWebhook").Call(
						jen.ID("ctx"),
						jen.ID("exampleWebhook").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
