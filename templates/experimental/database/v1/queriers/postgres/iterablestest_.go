package postgres

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile("postgres")

	utils.AddImports(ret)

	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	puvn := n.PluralUnexportedVarName()
	tn := n.PluralRouteName()
	scn := n.SingularCommonName()

	//allColumns := buildTableColumns(typ)

	ret.Add(
		jen.Func().IDf("buildMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn)).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Block(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(
				jen.ID("x").Dot("ID"),
				jen.ID("x").Dot("Name"),
				jen.ID("x").Dot("Details"),
				jen.ID("x").Dot("CreatedOn"),
				jen.ID("x").Dot("UpdatedOn"),
				jen.ID("x").Dot("ArchivedOn"),
				jen.ID("x").Dot("BelongsTo"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("buildErroneousMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn)).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Block(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(
				jen.ID("x").Dot("ArchivedOn"),
				jen.ID("x").Dot("Name"),
				jen.ID("x").Dot("Details"),
				jen.ID("x").Dot("CreatedOn"),
				jen.ID("x").Dot("UpdatedOn"),
				jen.ID("x").Dot("BelongsTo"),
				jen.ID("x").Dot("ID"),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_buildGet%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE belongs_to = $1 AND id = $2", tn),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dotf("buildGet%sQuery", sn).Call(jen.IDf("example%sID", sn), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.IDf("example%sID", sn), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_Get%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE belongs_to = $1 AND id = $2", tn),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Details").Op(":").Lit("details"),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dot("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dot("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE belongs_to = $1 AND id = $2", tn),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Details").Op(":").Lit("details"),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_buildGet%sCountQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dotf("buildGet%sCountQuery", sn).Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_Get%sCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%sCount", sn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_buildGetAll%sCountQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.Line(),
				jen.ID("actualQuery").Op(":=").ID("p").Dotf("buildGetAll%sCountQuery", pn).Call(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_GetAll%sCount", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID("p").Dotf("GetAll%sCount", pn).Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_buildGet%sQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dotf("buildGet%sQuery", pn).Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Line(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_Get%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.IDf("expected%s", sn).Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", fmt.Sprintf("%sList", sn)).Valuesln(
					jen.ID("Pagination").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").ID("expectedCount"),
					),
					jen.ID(pn).Op(":").Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
						jen.Op("*").IDf("expected%s", sn),
					),
				),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn))),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dot("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).
					Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error executing read query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error scanning %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.Line(),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying for count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.Line(),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_GetAll%sForUser", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("expected%s", sn).Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn))),
				jen.Line(),
				jen.ID("expected").Op(":=").Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Values(jen.Op("*").IDf("expected%s", sn)),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("example%s", sn).Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM %s WHERE archived_on IS NULL AND belongs_to = $1", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.IDf("example%s", sn))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_buildCreate%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Details").Op(":").Lit("details"),
					jen.ID("BelongsTo").Op(":").Lit(123),
				),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(3),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on", tn),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dotf("buildCreate%sQuery", sn).Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Name"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Details"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("BelongsTo"), jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_Create%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", fmt.Sprintf("%sCreationInput", sn)).Valuesln(
					jen.ID("Name").Op(":").ID("expected").Dot("Name"),
					jen.ID("BelongsTo").Op(":").ID("expected").Dot("BelongsTo"),
				),
				jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
					jen.ID("expected").Dot("ID"),
					jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("Details"),
					jen.ID("expected").Dot("BelongsTo"),
				).Dot("WillReturnRows").Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("example").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", fmt.Sprintf("%sCreationInput", sn)).Valuesln(
					jen.ID("Name").Op(":").ID("example").Dot("Name"),
					jen.ID("BelongsTo").Op(":").ID("example").Dot("BelongsTo"),
				),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (name,details,belongs_to) VALUES ($1,$2,$3) RETURNING id, created_on", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("example").Dot("Name"),
					jen.ID("example").Dot("Details"),
					jen.ID("example").Dot("BelongsTo"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_buildUpdate%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Details").Op(":").Lit("details"),
					jen.ID("BelongsTo").Op(":").Lit(123),
				),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(4),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET name = $1, details = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on", tn),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dotf("buildUpdate%sQuery", sn).Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Name"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("string"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Details"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("string"))), jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("BelongsTo"), jen.ID("args").Index(jen.Lit(2)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(3)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_Update%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET name = $1, details = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("Name"),
					jen.ID("expected").Dot("Details"),
					jen.ID("expected").Dot("BelongsTo"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnRows").Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.ID("err").Op(":=").ID("p").Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("example").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET name = $1, details = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("example").Dot("Name"),
					jen.ID("example").Dot("Details"),
					jen.ID("example").Dot("BelongsTo"),
					jen.ID("example").Dot("ID"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("err").Op(":=").ID("p").Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("example")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_buildArchive%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("Details").Op(":").Lit("details"),
					jen.ID("BelongsTo").Op(":").Lit(123),
				),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on", tn),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID("p").Dotf("buildArchive%sQuery", sn).Call(jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsTo")),
				jen.Line(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("BelongsTo"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestPostgres_Archive%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("expected").Dot("BelongsTo"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Lit(1), jen.Lit(1))),
				jen.Line(),
				jen.ID("err").Op(":=").ID("p").Dotf("Archive%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("example").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("Name").Op(":").Lit("name"),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					jen.ID("example").Dot("BelongsTo"),
					jen.ID("example").Dot("ID"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("err").Op(":=").ID("p").Dotf("Archive%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("example").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	return ret
}
