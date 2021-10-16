package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildMockRowsFromItems").Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").Uint64(), jen.ID("items").Op("...").Op("*").ID("types").Dot("Item")).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("columns").Op(":=").ID("itemsTableColumns"),
			jen.Newline(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Equals().ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.Newline(),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("columns")),
			jen.Newline(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("items")).Body(
				jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(
					jen.ID("x").Dot("ID"),
					jen.ID("x").Dot("Name"),
					jen.ID("x").Dot("Details"),
					jen.ID("x").Dot("CreatedOn"),
					jen.ID("x").Dot("LastUpdatedOn"),
					jen.ID("x").Dot("ArchivedOn"),
					jen.ID("x").Dot("BelongsToAccount"),
				),
				jen.Newline(),
				jen.If(jen.ID("includeCounts")).Body(
					jen.ID("rowValues").Equals().ID("append").Call(
						jen.ID("rowValues"),
						jen.ID("filteredCount"),
						jen.ID("len").Call(jen.ID("items")),
					),
				),
				jen.Newline(),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Newline(),
			jen.Return().ID("exampleRows"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestQuerier_ScanItems").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("surfaces row errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs row closing errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Nil()),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"),
						jen.ID("exampleItemID"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemExistenceQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.ID("true"))),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with sql.ErrNoRows"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemExistenceQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemExistenceQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
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
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.ID("exampleItem").Dot("ID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getItemQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("false"),
						jen.Zero(),
						jen.ID("exampleItem"),
					)),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.ID("exampleItem").Dot("ID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getItemQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
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
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleCount").Op(":=").Uint64().Call(jen.Lit(123)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getAllItemsCountQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.Uint64().Call(jen.Lit(123)))),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("getAllItemsCountQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetTotalItemCount").Call(jen.ID("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
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
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("true"),
						jen.ID("exampleItemList").Dot("FilteredCount"),
						jen.ID("exampleItemList").Dot("Items").Op("..."),
					)),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil filter"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.Nil()),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.ID("exampleItemList").Dot("Page").Equals().Zero(),
					jen.ID("exampleItemList").Dot("Limit").Equals().Zero(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
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
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("true"),
						jen.ID("exampleItemList").Dot("FilteredCount"),
						jen.ID("exampleItemList").Dot("Items").Op("..."),
					)),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
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
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
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
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Newline(),
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
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dot("buildGetItemsWithIDsQuery").Call(
						constants.CtxVar(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromItems").Call(
						jen.ID("false"),
						jen.Zero(),
						jen.ID("exampleItemList").Dot("Items").Op("..."),
					)),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid IDs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dot("buildGetItemsWithIDsQuery").Call(
						constants.CtxVar(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(constants.ObligatoryError()),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error scanning query results"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemList").Op(":=").ID("fakes").Dot("BuildFakeItemList").Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("exampleItemList").Dot("Items")).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dot("buildGetItemsWithIDsQuery").Call(
						constants.CtxVar(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Newline(),
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
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleItem").Dot("ID").Equals().Lit("1"),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemDatabaseCreationInputFromItem").Call(jen.ID("exampleItem")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleInput").Dot("ID"),
						jen.ID("exampleInput").Dot("Name"),
						jen.ID("exampleInput").Dot("Details"),
						jen.ID("exampleInput").Dot("BelongsToAccount"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleItem").Dot("ID"))),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleItem").Dot("CreatedOn")),
					jen.Newline(),
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
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("expectedErr").Op(":=").Qual("errors", "New").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemDatabaseCreationInputFromItem").Call(jen.ID("exampleItem")),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleInput").Dot("ID"),
						jen.ID("exampleInput").Dot("Name"),
						jen.ID("exampleInput").Dot("Details"),
						jen.ID("exampleInput").Dot("BelongsToAccount"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("itemCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.ID("expectedErr")),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().ID("exampleItem").Dot("CreatedOn")),
					jen.Newline(),
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
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleItem").Dot("Name"), jen.ID("exampleItem").Dot("Details"), jen.ID("exampleItem").Dot("BelongsToAccount"), jen.ID("exampleItem").Dot("ID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("updateItemQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleItem").Dot("ID"))),
					jen.Newline(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem"),
						),
					),
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.Nil(),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleItem").Dot("Name"),
						jen.ID("exampleItem").Dot("Details"),
						jen.ID("exampleItem").Dot("BelongsToAccount"),
						jen.ID("exampleItem").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("updateItemQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("UpdateItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItem"),
						),
					),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("archiveItemQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.ID("exampleItemID"))),
					jen.Newline(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemID"),
							jen.ID("exampleAccountID"),
						),
					),
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("exampleItemID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.ID("exampleItemID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("archiveItemQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dot("ArchiveItem").Call(
							jen.ID("ctx"),
							jen.ID("exampleItemID"),
							jen.ID("exampleAccountID"),
						),
					),
					jen.Newline(),
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
