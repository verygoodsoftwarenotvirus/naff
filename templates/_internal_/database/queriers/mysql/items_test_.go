package mysql

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func itemsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildMockRowsFromItems").Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").Uint64(), jen.ID("items").Op("...").Op("*").ID("types").Dot("Item")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("columns").Op(":=").ID("itemsTableColumns"),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Equals().ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("columns")),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("items")).Body(
				jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(jen.ID("x").Dot("ID"), jen.ID("x").Dot("Name"), jen.ID("x").Dot("Details"), jen.ID("x").Dot("CreatedOn"), jen.ID("x").Dot("LastUpdatedOn"), jen.ID("x").Dot("ArchivedOn"), jen.ID("x").Dot("BelongsToAccount")),
				jen.If(jen.ID("includeCounts")).Body(
					jen.ID("rowValues").Equals().ID("append").Call(
						jen.ID("rowValues"),
						jen.ID("filteredCount"),
						jen.ID("len").Call(jen.ID("items")),
					)),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Return().ID("exampleRows"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ScanItems").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanItems").Call(
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
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Nil()),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanItems").Call(
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
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ItemExists").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemExistenceQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().String().Valuesln(jen.Lit("exists"))).Dot("AddRow").Call(jen.ID("true"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("ctx"),
						jen.ID("exampleItemID"),
						jen.ID("exampleAccountID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("exampleAccountID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("ctx"),
						jen.ID("exampleItemID"),
						jen.Lit(""),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with sql.ErrNoRows"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemExistenceQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("ctx"),
						jen.ID("exampleItemID"),
						jen.ID("exampleAccountID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemExistenceQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("ctx"),
						jen.ID("exampleItemID"),
						jen.ID("exampleAccountID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID"), jen.ID("exampleItem").Dot("ID")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getItemQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("false"),
						jen.Zero(),
						jen.ID("exampleItem"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
						jen.ID("ctx"),
						jen.ID("exampleItem").Dot("ID"),
						jen.ID("exampleAccountID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleItem"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("exampleAccountID"),
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
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
						jen.ID("ctx"),
						jen.ID("exampleItem").Dot("ID"),
						jen.Lit(""),
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID"), jen.ID("exampleItem").Dot("ID")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getItemQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
						jen.ID("ctx"),
						jen.ID("exampleItem").Dot("ID"),
						jen.ID("exampleAccountID"),
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
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetTotalItemCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleCount").Op(":=").Uint64().Call(jen.Lit(123)),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getAllItemsCountQuery"))).Dot("WithArgs").Call().Dot("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.Uint64().Call(jen.Lit(123)))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetTotalItemCount").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleCount"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getAllItemsCountQuery"))).Dot("WithArgs").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetTotalItemCount").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetItems").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Call(
						jen.ID("ctx"),
						jen.Lit("items"),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.ID("itemsTableColumns"),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("true"),
						jen.ID("exampleItemList").Dot("FilteredCount"),
						jen.ID("exampleItemList").Dot("Items").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleItemList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil filter"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.Nil()),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("exampleItemList").Dot("Page").Equals().Zero(),
					jen.ID("exampleItemList").Dot("Limit").Equals().Zero(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Call(
						jen.ID("ctx"),
						jen.Lit("items"),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.ID("itemsTableColumns"),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("true"),
						jen.ID("exampleItemList").Dot("FilteredCount"),
						jen.ID("exampleItemList").Dot("Items").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleItemList"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
						jen.ID("ctx"),
						jen.Lit(""),
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Call(
						jen.ID("ctx"),
						jen.Lit("items"),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.ID("itemsTableColumns"),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
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
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Call(
						jen.ID("ctx"),
						jen.Lit("items"),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.ID("itemsTableColumns"),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
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
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_GetItemsWithIDs").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("exampleArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID")),
					jen.Var().Defs(
						jen.ID("exampleIDs").Index().String(),
					),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
						jen.ID("exampleArgs").Equals().ID("append").Call(
							jen.ID("exampleArgs"),
							jen.ID("x").Dot("ID"),
						),
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("query").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.ID("getItemsWithIDsQuery"),
						jen.ID("joinIDs").Call(jen.ID("exampleIDs")),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("exampleArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("false"),
						jen.Zero(),
						jen.ID("exampleItemList").Dot("Items").Op("..."),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItemsWithIDs").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.Zero(),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleItemList").Dot("Items"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItemsWithIDs").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("defaultLimit"),
						jen.Nil(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid IDs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItemsWithIDs").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("defaultLimit"),
						jen.Nil(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("exampleArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID")),
					jen.Var().Defs(
						jen.ID("exampleIDs").Index().String(),
					),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
						jen.ID("exampleArgs").Equals().ID("append").Call(
							jen.ID("exampleArgs"),
							jen.ID("x").Dot("ID"),
						),
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("query").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.ID("getItemsWithIDsQuery"),
						jen.ID("joinIDs").Call(jen.ID("exampleIDs")),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("exampleArgs")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItemsWithIDs").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error scanning query results"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("exampleArgs").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID")),
					jen.Var().Defs(
						jen.ID("exampleIDs").Index().String(),
					),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
						jen.ID("exampleArgs").Equals().ID("append").Call(
							jen.ID("exampleArgs"),
							jen.ID("x").Dot("ID"),
						),
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("query").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.ID("getItemsWithIDsQuery"),
						jen.ID("joinIDs").Call(jen.ID("exampleIDs")),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("exampleArgs")).Op("...")).Dot("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItemsWithIDs").Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_CreateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleItem").Dot("ID").Equals().Lit("1"),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemDatabaseCreationInputFromItem").Call(jen.ID("exampleItem")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleInput").Dot("ID"), jen.ID("exampleInput").Dot("Name"), jen.ID("exampleInput").Dot("Details"), jen.ID("exampleInput").Dot("BelongsToAccount")),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleItem").Dot("ID"))),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleItem").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleItem"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
						jen.ID("ctx"),
						jen.Nil(),
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
					jen.ID("expectedErr").Op(":=").Qual("errors", "New").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemDatabaseCreationInputFromItem").Call(jen.ID("exampleItem")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleInput").Dot("ID"), jen.ID("exampleInput").Dot("Name"), jen.ID("exampleInput").Dot("Details"), jen.ID("exampleInput").Dot("BelongsToAccount")),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemCreationQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnError").Call(jen.ID("expectedErr")),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleItem").Dot("CreatedOn")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.ID("expectedErr"),
						),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_UpdateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleItem").Dot("Name"), jen.ID("exampleItem").Dot("Details"), jen.ID("exampleItem").Dot("BelongsToAccount"), jen.ID("exampleItem").Dot("ID")),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("updateItemQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleItem").Dot("ID"))),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.Nil(),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleItem").Dot("Name"), jen.ID("exampleItem").Dot("Details"), jen.ID("exampleItem").Dot("BelongsToAccount"), jen.ID("exampleItem").Dot("ID")),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("updateItemQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ArchiveItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("archiveItemQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleItemID"))),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemID"),
							jen.ID("exampleAccountID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveItem").Call(
							jen.ID("ctx"),
							jen.Lit(""),
							jen.ID("exampleAccountID"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemID"),
							jen.Lit(""),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("archiveItemQuery"))).Dot("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemID"),
							jen.ID("exampleAccountID"),
						),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
