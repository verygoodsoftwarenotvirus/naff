package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsTestDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)


	ret.Add(jen.Func().ID("buildMockRowFromItem").Params(jen.ID("item").Op("*").ID("models").Dot(
		"Item",
	)).Params(jen.Op("*").ID("sqlmock").Dot(
		"Rows",
	)).Block(
		jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
			"NewRows",
		).Call(jen.ID("itemsTableColumns")).Dot(
			"AddRow",
		).Call(jen.ID("item").Dot(
			"ID",
		), jen.ID("item").Dot(
			"Name",
		), jen.ID("item").Dot(
			"Details",
		), jen.ID("item").Dot(
			"CreatedOn",
		), jen.ID("item").Dot(
			"UpdatedOn",
		), jen.ID("item").Dot(
			"ArchivedOn",
		), jen.ID("item").Dot(
			"BelongsTo",
		)),
		jen.Return().ID("exampleRows"),
	),
	)
	ret.Add(jen.Func().ID("buildErroneousMockRowFromItem").Params(jen.ID("item").Op("*").ID("models").Dot(
		"Item",
	)).Params(jen.Op("*").ID("sqlmock").Dot(
		"Rows",
	)).Block(
		jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
			"NewRows",
		).Call(jen.ID("itemsTableColumns")).Dot(
			"AddRow",
		).Call(jen.ID("item").Dot(
			"ArchivedOn",
		), jen.ID("item").Dot(
			"Name",
		), jen.ID("item").Dot(
			"Details",
		), jen.ID("item").Dot(
			"CreatedOn",
		), jen.ID("item").Dot(
			"UpdatedOn",
		), jen.ID("item").Dot(
			"BelongsTo",
		), jen.ID("item").Dot(
			"ID",
		)),
		jen.Return().ID("exampleRows"),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_buildGetItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleItemID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = $1 AND id = $2"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dot(
				"buildGetItemQuery",
			).Call(jen.ID("exampleItemID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleItemID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_GetItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = $1 AND id = $2"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromItem").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = $1 AND id = $2"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_buildGetItemCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dot(
				"buildGetItemCountQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_GetItemCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItemCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_buildGetAllItemsCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.ID("actualQuery").Op(":=").ID("p").Dot(
				"buildGetAllItemsCountQuery",
			).Call(),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_GetAllItemsCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetAllItemsCount",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_buildGetItemsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dot(
				"buildGetItemsQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_GetItems").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem1").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromItem").Call(jen.ID("expectedItem1"))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"ItemList",
			).Valuesln(jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Items").Op(":").Index().ID("models").Dot(
				"Item",
			).Valuesln(jen.Op("*").ID("expectedItem1"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error executing read query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error scanning item"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem1").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromItem").Call(jen.ID("expectedItem1"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error querying for count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem1").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromItem").Call(jen.ID("expectedItem1"))),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_GetAllItemsForUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromItem").Call(jen.ID("expectedItem"))),
			jen.ID("expected").Op(":=").Index().ID("models").Dot(
				"Item",
			).Valuesln(jen.Op("*").ID("expectedItem")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = $1"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromItem").Call(jen.ID("expectedItem"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_buildCreateItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(3),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO items (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dot(
				"buildCreateItemQuery",
			).Call(jen.ID("expected")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"Name",
			), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"Details",
			), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("uint64"))),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_CreateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"ItemCreationInput",
			).Valuesln(jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
			), jen.ID("BelongsTo").Op(":").ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("id"), jen.Lit("created_on"))).Dot(
				"AddRow",
			).Call(jen.ID("expected").Dot(
				"ID",
			), jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO items (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
			), jen.ID("expected").Dot(
				"Details",
			), jen.ID("expected").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("exampleRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"CreateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("example").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"ItemCreationInput",
			).Valuesln(jen.ID("Name").Op(":").ID("example").Dot(
				"Name",
			), jen.ID("BelongsTo").Op(":").ID("example").Dot(
				"BelongsTo",
			)),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO items (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("example").Dot(
				"Name",
			), jen.ID("example").Dot(
				"Details",
			), jen.ID("example").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot(
				"CreateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_buildUpdateItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(321), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(4),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET name = $1, details = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dot(
				"buildUpdateItemQuery",
			).Call(jen.ID("expected")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"Name",
			), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"Details",
			), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"ID",
			), jen.ID("args").Index(jen.Lit(3)).Assert(jen.ID("uint64"))),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_UpdateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(jen.Lit("updated_on"))).Dot(
				"AddRow",
			).Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET name = $1, details = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
			), jen.ID("expected").Dot(
				"Details",
			), jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("exampleRows")),
			jen.ID("err").Op(":=").ID("p").Dot(
				"UpdateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("example").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET name = $1, details = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("example").Dot(
				"Name",
			), jen.ID("example").Dot(
				"Details",
			), jen.ID("example").Dot(
				"BelongsTo",
			), jen.ID("example").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("err").Op(":=").ID("p").Dot(
				"UpdateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("example")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_buildArchiveItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(321), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dot(
				"buildArchiveItemQuery",
			).Call(jen.ID("expected").Dot(
				"ID",
			), jen.ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected").Dot(
				"ID",
			), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),
	)
	ret.Add(jen.Func().ID("TestPostgres_ArchiveItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"BelongsTo",
			), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.Lit(1), jen.Lit(1))),
			jen.ID("err").Op(":=").ID("p").Dot(
				"ArchiveItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("example").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"),
			jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("example").Dot(
				"BelongsTo",
			), jen.ID("example").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("err").Op(":=").ID("p").Dot(
				"ArchiveItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("example").Dot(
				"ID",
			), jen.ID("expectedUserID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	)
	return ret
}
