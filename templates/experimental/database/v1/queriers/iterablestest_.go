package queriers

import (
	"fmt"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildGeneralFields(varName string, typ models.DataType) []jen.Code {
	fields := []jen.Code{jen.ID(varName).Dot("ID")}

	for _, field := range typ.Fields {
		fields = append(fields, jen.ID(varName).Dot(field.Name.Singular()))
	}

	fields = append(fields,
		jen.ID(varName).Dot("CreatedOn"),
		jen.ID(varName).Dot("UpdatedOn"),
		jen.ID(varName).Dot("ArchivedOn"),
		jen.ID(varName).Dot("BelongsTo"),
	)

	return fields
}

func buildBadFields(varName string, typ models.DataType) []jen.Code {
	fields := []jen.Code{jen.ID(varName).Dot("ArchivedOn")}

	for _, field := range typ.Fields {
		fields = append(fields, jen.ID(varName).Dot(field.Name.Singular()))
	}

	fields = append(fields,
		jen.ID(varName).Dot("CreatedOn"),
		jen.ID(varName).Dot("UpdatedOn"),
		jen.ID(varName).Dot("BelongsTo"),
		jen.ID(varName).Dot("ID"),
	)

	return fields
}

func buildStringColumns(typ models.DataType) string {
	out := []string{"id"}

	for _, field := range typ.Fields {
		out = append(out, field.Name.RouteName())
	}

	out = append(out, "created_on", "updated_on", "archived_on", "belongs_to")

	return strings.Join(out, ", ")
}

func buildNonStandardStringColumns(typ models.DataType) string {
	var out []string

	for _, field := range typ.Fields {
		out = append(out, field.Name.RouteName())
	}

	return strings.Join(out, ", ")
}

func buildUpdateQueryParts(typ models.DataType) []string {
	var out []string

	for i, field := range typ.Fields {
		out = append(out, fmt.Sprintf("%s = $%d", field.Name.RouteName(), i+1))
	}

	return out
}

func iterablesTestDotGo(dbvendor *wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(ret)

	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	puvn := n.PluralUnexportedVarName()
	tn := n.PluralRouteName()
	scn := n.SingularCommonName()
	dbvsn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	fieldCols := buildNonStandardStringColumns(typ)
	cols := buildStringColumns(typ)

	var ips []string
	for i := range typ.Fields {
		ips = append(ips, fmt.Sprintf("$%d", i+1))
	}
	ips = append(ips, fmt.Sprintf("$%d", len(ips)+1))
	insertPlaceholders := strings.Join(ips, ",")

	//allColumns := buildTableColumns(typ)
	gFields := buildGeneralFields("x", typ)

	ret.Add(
		jen.Func().IDf("buildMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn)).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Block(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(gFields...),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	badFields := buildBadFields("x", typ)

	ret.Add(
		jen.Func().IDf("buildErroneousMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn)).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Block(
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(badFields...),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGet%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE belongs_to = $1 AND id = $2", cols, tn),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", sn).Call(jen.IDf("example%sID", sn), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.IDf("example%sID", sn), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_Get%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE belongs_to = $1 AND id = $2", cols, tn),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dot("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dot("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE belongs_to = $1 AND id = $2", cols, tn),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGet%sCountQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", tn),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sCountQuery", sn).Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_Get%sCount", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
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
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%sCount", sn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGetAll%sCountQuery", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.Line(),
				jen.ID("actualQuery").Op(":=").ID(dbfl).Dotf("buildGetAll%sCountQuery", pn).Call(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAll%sCount", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sCount", pn).Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildGet%sQuery", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", cols, tn),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", pn).Call(jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.Line(),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
				jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_Get%s", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", cols, tn),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.IDf("expected%s", sn).Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
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
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", cols, tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error executing read query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", cols, tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error scanning %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").ID("uint64").Call(jen.Lit(321)),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", cols, tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying for count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20", cols, tn),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_GetAll%sForUser", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("expected%s", sn).Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1", cols, tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn))),
				jen.Line(),
				jen.ID("expected").Op(":=").Index().Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Values(jen.Op("*").IDf("expected%s", sn)),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1", cols, tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).Dot("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1", cols, tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("example%s", sn).Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = $1", cols, tn),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).Dot("WithArgs").Call(jen.ID("expectedUserID")).
					Dot("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.IDf("example%s", sn))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	creationEqualityExpectations := buildCreationEqualityExpectations("expected", typ)
	createQueryTestBody := []jen.Code{
		jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
			jen.ID("ID").Op(":").Lit(321),
			jen.ID("BelongsTo").Op(":").Lit(123),
		),
		jen.ID("expectedArgCount").Op(":=").Lit(3),
		jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (%s,belongs_to) VALUES (%s) RETURNING id, created_on", tn, strings.ReplaceAll(fieldCols, " ", ""), insertPlaceholders),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildCreate%sQuery", sn).Call(jen.ID("expected")),
		jen.Line(),
		jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	}
	createQueryTestBody = append(createQueryTestBody, creationEqualityExpectations...)

	ret.Add(
		jen.Func().IDf("Test%s_buildCreate%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(createQueryTestBody...,
			)),
		),
		jen.Line(),
	)

	//exampleInputFields := buildFieldMaps("example", typ)
	expectedInputFields := buildFieldMaps("expected", typ)

	nonEssentialFields := func(varName string, typ models.DataType) []jen.Code {
		var out []jen.Code

		for _, field := range typ.Fields {
			out = append(out, jen.ID(varName).Dot(field.Name.Singular()))
		}
		out = append(out, jen.ID(varName).Dot("BelongsTo"))

		return out
	}
	nef := nonEssentialFields("expected", typ)

	ret.Add(
		jen.Func().IDf("Test%s_Create%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", fmt.Sprintf("%sCreationInput", sn)).Valuesln(
					expectedInputFields...,
				),
				jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
					jen.ID("expected").Dot("ID"),
					jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (%s,belongs_to) VALUES (%s) RETURNING id, created_on", tn, strings.ReplaceAll(fieldCols, " ", ""), insertPlaceholders),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					nef...,
				).Dot("WillReturnRows").Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedInput").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", fmt.Sprintf("%sCreationInput", sn)).Valuesln(
					expectedInputFields...,
				),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (%s,belongs_to) VALUES (%s) RETURNING id, created_on", tn, strings.ReplaceAll(fieldCols, " ", ""), insertPlaceholders),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					nef...,
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	updateCols := buildUpdateQueryParts(typ)
	updateColsStr := strings.Join(updateCols, ", ")

	testBuildUpdateQueryBody := []jen.Code{
		jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
			jen.ID("ID").Op(":").Lit(321),
			jen.ID("BelongsTo").Op(":").Lit(123),
		),
		jen.ID("expectedArgCount").Op(":=").Lit(4),
		jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET %s, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $%d AND id = $%d RETURNING updated_on", tn, updateColsStr, len(updateCols)+1, len(updateCols)+2),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildUpdate%sQuery", sn).Call(jen.ID("expected")),
		jen.Line(),
		jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.ID("assert").Dot("Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	}

	testBuildUpdateQueryBody = append(testBuildUpdateQueryBody, creationEqualityExpectations...)
	testBuildUpdateQueryBody = append(testBuildUpdateQueryBody,
		jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(len(creationEqualityExpectations))).Assert(jen.ID("uint64"))),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildUpdate%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				testBuildUpdateQueryBody...,
			)),
		),
		jen.Line(),
	)

	buildExpectQueryArgs := func(varName string, typ models.DataType) []jen.Code {
		var out []jen.Code

		for _, field := range typ.Fields {
			out = append(out, jen.ID(varName).Dot(field.Name.Singular()))
		}

		out = append(out,
			jen.ID(varName).Dot("BelongsTo"),
			jen.ID(varName).Dot("ID"),
		)
		return out
	}
	expectQueryArgs := buildExpectQueryArgs("expected", typ)

	ret.Add(
		jen.Func().IDf("Test%s_Update%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET %s, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on", tn, updateColsStr),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					expectQueryArgs...,
				).Dot("WillReturnRows").Call(jen.ID("exampleRows")),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET %s, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $3 AND id = $4 RETURNING updated_on", tn, updateColsStr),
				jen.Line(),
				jen.List(jen.ID("p"), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).Dot("WithArgs").Callln(
					expectQueryArgs...,
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s_buildArchive%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("p"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),

					jen.ID("BelongsTo").Op(":").Lit(123),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on", tn),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildArchive%sQuery", sn).Call(jen.ID("expected").Dot("ID"), jen.ID("expected").Dot("BelongsTo")),
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
		jen.Func().IDf("Test%s_Archive%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
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
				jen.ID("err").Op(":=").ID(dbfl).Dotf("Archive%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("example").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
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
				jen.ID("err").Op(":=").ID(dbfl).Dotf("Archive%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("example").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)
	return ret
}

func buildCreationEqualityExpectations(varName string, typ models.DataType) []jen.Code {
	var out []jen.Code

	for i, field := range typ.Fields {
		out = append(out,
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID(varName).Dot(field.Name.Singular()), jen.ID("args").Index(jen.Lit(i)).Assert(jen.ID(field.Type))),
		)
	}
	out = append(out,
		jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("BelongsTo"), jen.ID("args").Index(jen.Lit(len(out))).Assert(jen.ID("uint64"))),
	)

	return out
}

func buildFieldMaps(varName string, typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		xn := field.Name.Singular()
		out = append(out, jen.ID(xn).Op(":").ID(varName).Dot(xn))
	}
	out = append(out, jen.ID("BelongsTo").Op(":").ID(varName).Dot("BelongsTo"))

	return out
}
