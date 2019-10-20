package mariadb

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsTestDotGo() *jen.File {
	ret := jen.NewFile("mariadb")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildMockRowFromItem").Params(jen.ID("item").Op("*").ID("models").Dot(
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
	),
	jen.ID("item").Dot(
			"Name",
	),
	jen.ID("item").Dot(
			"Details",
	),
	jen.ID("item").Dot(
			"CreatedOn",
	),
	jen.ID("item").Dot(
			"UpdatedOn",
	),
	jen.ID("item").Dot(
			"ArchivedOn",
	),
	jen.ID("item").Dot(
			"BelongsTo",
		)),
		jen.Return().ID("exampleRows"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildErroneousMockRowFromItem").Params(jen.ID("item").Op("*").ID("models").Dot(
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
	),
	jen.ID("item").Dot(
			"Name",
	),
	jen.ID("item").Dot(
			"Details",
	),
	jen.ID("item").Dot(
			"CreatedOn",
	),
	jen.ID("item").Dot(
			"UpdatedOn",
	),
	jen.ID("item").Dot(
			"BelongsTo",
	),
	jen.ID("item").Dot(
			"ID",
		)),
		jen.Return().ID("exampleRows"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleItemID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetItemQuery",
			).Call(jen.ID("exampleItemID"), jen.ID("exampleUserID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleItemID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = ? AND id = ?"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromItem").Call(jen.ID("expected"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("expectedUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = ? AND id = ?"),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details")),
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID"), jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetItemCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetItemCountQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetItemCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItemCount",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetAllItemsCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.ID("actualQuery").Op(":=").ID("m").Dot(
				"buildGetAllItemsCountQuery",
			).Call(),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetAllItemsCount").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllItemsCount",
			).Call(jen.Qual("context", "Background").Call()),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildGetItemsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expectedArgCount").Op(":=").Lit(1),
			jen.ID("expectedQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildGetItemsQuery",
			).Call(jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("exampleUserID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetItems").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem1").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
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
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("count"))).Dot(
				"AddRow",
			).Call(jen.ID("expectedCount"))),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"ItemList",
			).Valuesln(
	jen.ID("Pagination").Op(":").ID("models").Dot(
				"Pagination",
			).Valuesln(
	jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("TotalCount").Op(":").ID("expectedCount")), jen.ID("Items").Op(":").Index().ID("models").Dot(
				"Item",
			).Valuesln(
	jen.Op("*").ID("expectedItem1"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error executing read query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error scanning item"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem1").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromItem").Call(jen.ID("expectedItem1"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying for count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem1").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"),
			jen.ID("expectedCountQuery").Op(":=").Lit("SELECT COUNT(id) FROM items WHERE archived_on IS NULL"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
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
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetItems",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"DefaultQueryFilter",
			).Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_GetAllItemsForUser").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildMockRowFromItem").Call(jen.ID("expectedItem"))),
			jen.ID("expected").Op(":=").Index().ID("models").Dot(
				"Item",
			).Valuesln(
	jen.Op("*").ID("expectedItem")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("database/sql", "ErrNoRows")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expectedItem").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("Name").Op(":").Lit("name")),
			jen.ID("expectedListQuery").Op(":=").Lit("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expectedUserID")).Dot(
				"WillReturnRows",
			).Call(jen.ID("buildErroneousMockRowFromItem").Call(jen.ID("expectedItem"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"GetAllItemsForUser",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildCreateItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(3),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO items (name,details,belongs_to,created_on) VALUES (?,?,?,UNIX_TIMESTAMP())"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildCreateItemQuery",
			).Call(jen.ID("expected")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"Name",
	),
	jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"Details",
	),
	jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_CreateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"ItemCreationInput",
			).Valuesln(
	jen.ID("Name").Op(":").ID("expected").Dot(
				"Name",
	),
	jen.ID("BelongsTo").Op(":").ID("expected").Dot(
				"BelongsTo",
			)),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedCreationQuery").Op(":=").Lit("INSERT INTO items (name,details,belongs_to,created_on) VALUES (?,?,?,UNIX_TIMESTAMP())"),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
	),
	jen.ID("expected").Dot(
				"Details",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.ID("int64").Call(jen.ID("expected").Dot(
				"ID",
			)), jen.Lit(1))),
			jen.ID("expectedTimeQuery").Op(":=").Lit("SELECT created_on FROM items WHERE id = ?"),
			jen.ID("mockDB").Dot(
				"ExpectQuery",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnRows",
			).Call(jen.ID("sqlmock").Dot(
				"NewRows",
			).Call(jen.Index().ID("string").Valuesln(
	jen.Lit("created_on"))).Dot(
				"AddRow",
			).Call(jen.ID("expected").Dot(
				"CreatedOn",
			))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"CreateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("example").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedInput").Op(":=").Op("&").ID("models").Dot(
				"ItemCreationInput",
			).Valuesln(
	jen.ID("Name").Op(":").ID("example").Dot(
				"Name",
	),
	jen.ID("BelongsTo").Op(":").ID("example").Dot(
				"BelongsTo",
			)),
			jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO items (name,details,belongs_to,created_on) VALUES (?,?,?,UNIX_TIMESTAMP())"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("example").Dot(
				"Name",
	),
	jen.ID("example").Dot(
				"Details",
	),
	jen.ID("example").Dot(
				"BelongsTo",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("m").Dot(
				"CreateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildUpdateItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(321), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(4),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET name = ?, details = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildUpdateItemQuery",
			).Call(jen.ID("expected")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"Name",
	),
	jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"Details",
	),
	jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("args").Index(jen.Lit(3)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_UpdateItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET name = ?, details = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"Name",
	),
	jen.ID("expected").Dot(
				"Details",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.Lit(1), jen.Lit(1))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"UpdateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("example").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET name = ?, details = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("example").Dot(
				"Name",
	),
	jen.ID("example").Dot(
				"Details",
	),
	jen.ID("example").Dot(
				"BelongsTo",
	),
	jen.ID("example").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"UpdateItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("example")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_buildArchiveItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("m"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(321), jen.ID("Name").Op(":").Lit("name"), jen.ID("Details").Op(":").Lit("details"), jen.ID("BelongsTo").Op(":").Lit(123)),
			jen.ID("expectedArgCount").Op(":=").Lit(2),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("m").Dot(
				"buildArchiveItemQuery",
			).Call(jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("expected").Dot(
				"BelongsTo",
			)),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			jen.ID("assert").Dot(
				"Len",
			).Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestMariaDB_ArchiveItem").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("expected").Dot(
				"BelongsTo",
	),
	jen.ID("expected").Dot(
				"ID",
			)).Dot(
				"WillReturnResult",
			).Call(jen.ID("sqlmock").Dot(
				"NewResult",
			).Call(jen.Lit(1), jen.Lit(1))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"ArchiveItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot(
				"ID",
	),
	jen.ID("expectedUserID")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("example").Op(":=").Op("&").ID("models").Dot(
				"Item",
			).Valuesln(
	jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").Lit("name"), jen.ID("BelongsTo").Op(":").ID("expectedUserID"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot(
				"Unix",
			).Call())),
			jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"),
			jen.List(jen.ID("m"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(
				"ExpectExec",
			).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot(
				"WithArgs",
			).Call(jen.ID("example").Dot(
				"BelongsTo",
	),
	jen.ID("example").Dot(
				"ID",
			)).Dot(
				"WillReturnError",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("err").Op(":=").ID("m").Dot(
				"ArchiveItem",
			).Call(jen.Qual("context", "Background").Call(), jen.ID("example").Dot(
				"ID",
	),
	jen.ID("expectedUserID")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot(
				"ExpectationsWereMet",
			).Call(), jen.Lit("not all database expectations were met")),
		)),
	),
	jen.Line(),
	)
	return ret
}
