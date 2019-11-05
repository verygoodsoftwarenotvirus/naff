package queriers

import (
	"fmt"
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	postgresCurrentUnixTimeQuery = "extract(epoch FROM NOW())"
	sqliteCurrentUnixTimeQuery   = "(strftime('%s','now'))"
	mariaDBUnixTimeQuery         = "UNIX_TIMESTAMP()"
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

func buildUpdateQueryParts(dbrn string, typ models.DataType) []string {
	var out []string

	for i, field := range typ.Fields {
		out = append(out, fmt.Sprintf("%s = %s", field.Name.RouteName(), getIncIndex(dbrn, uint(i))))
	}

	return out
}

func getIncIndex(dbrn string, index uint) string {
	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	if isPostgres {
		return fmt.Sprintf("$%d", index+1)
	} else if isSqlite || isMariaDB {
		return "?"
	}
	return ""
}

func getTimeQuery(dbrn string) string {
	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	if isPostgres {
		return postgresCurrentUnixTimeQuery
	}
	if isSqlite {
		return sqliteCurrentUnixTimeQuery
	}
	if isMariaDB {
		return mariaDBUnixTimeQuery
	}
	return ""
}

func iterablesTestDotGo(pkgRoot string, dbvendor wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(pkgRoot, []models.DataType{typ}, ret)

	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	puvn := n.PluralUnexportedVarName()
	tn := n.PluralRouteName()
	scn := n.SingularCommonName()
	dbvsn := dbvendor.Singular()
	dbrn := dbvendor.RouteName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	fieldCols := buildNonStandardStringColumns(typ)
	cols := buildStringColumns(typ)

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	var ips []string
	for i := range typ.Fields {
		ips = append(ips, getIncIndex(dbrn, uint(i)))
	}
	ips = append(ips, getIncIndex(dbrn, uint(len(ips))))
	insertPlaceholders := strings.Join(ips, ",")

	//allColumns := buildTableColumns(typ)
	gFields := buildGeneralFields("x", typ)

	ret.Add(
		jen.Func().IDf("buildMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), sn)).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(gFields...),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	badFields := buildBadFields("x", typ)

	ret.Add(
		jen.Func().IDf("buildErroneousMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), sn)).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(badFields...),
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
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE belongs_to = %s AND id = %s", cols, tn, getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", sn).Call(jen.IDf("example%sID", sn), jen.ID("exampleUserID")),
				jen.Line(),
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
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE belongs_to = %s AND id = %s", cols, tn, getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE belongs_to = %s AND id = %s", cols, tn, getIncIndex(dbrn, 0), getIncIndex(dbrn, 1)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
				),
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
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
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sCountQuery", sn).Call(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
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
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", tn, getIncIndex(dbrn, 0)),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%sCount", sn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.Line(),
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
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
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
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sCount", pn).Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.Line(),
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
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.Line(),
				jen.ID("expectedArgCount").Op(":=").Lit(1),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", cols, tn, getIncIndex(dbrn, 0)),
				jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", pn).Call(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
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
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", cols, tn, getIncIndex(dbrn, 0)),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.IDf("expected%s", sn).Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sList", sn)).Valuesln(
					jen.ID("Pagination").Op(":").Qual(filepath.Join(pkgRoot, "models/v1"), "Pagination").Valuesln(
						jen.ID("Page").Op(":").Lit(1),
						jen.ID("Limit").Op(":").Lit(20),
						jen.ID("TotalCount").Op(":").ID("expectedCount"),
					),
					jen.ID(pn).Op(":").Index().Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
						jen.Op("*").IDf("expected%s", sn),
					),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn))),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).
					Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", cols, tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error executing read query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", cols, tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error scanning %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", cols, tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying for count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s LIMIT 20", cols, tn, getIncIndex(dbrn, 0)),
				jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected"))),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkgRoot, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
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
				jen.IDf("expected%s", sn).Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s", cols, tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn))),
				jen.Line(),
				jen.ID("expected").Op(":=").Index().Qual(filepath.Join(pkgRoot, "models/v1"), sn).Values(jen.Op("*").IDf("expected%s", sn)),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s", cols, tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s", cols, tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("example%s", sn).Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.ID("expectedListQuery").Op(":=").Litf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to = %s", cols, tn, getIncIndex(dbrn, 0)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID("expectedUserID")).
					Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.IDf("example%s", sn))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	)

	////////////

	var queryTail string
	if isPostgres {
		queryTail = " RETURNING id, created_on"
	}

	var (
		createdOnAddendum,
		createdOnValueAdd string
	)
	if isMariaDB {
		createdOnAddendum = ",created_on"
		createdOnValueAdd = ",UNIX_TIMESTAMP()"
	}

	thisFuncExpectedArgCount := len(ips) - 1

	creationEqualityExpectations := buildCreationEqualityExpectations("expected", typ)
	createQueryTestBody := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
			jen.ID("ID").Op(":").Lit(321), jen.ID("BelongsTo").Op(":").Lit(123),
		),
		jen.ID("expectedArgCount").Op(":=").Lit(1 + thisFuncExpectedArgCount),
		jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (%s,belongs_to%s) VALUES (%s%s)%s",
			tn,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		),
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
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(createQueryTestBody...)),
		),
		jen.Line(),
	)

	////////////

	//exampleInputFields := buildFieldMaps("example", typ)
	expectedInputFields := buildFieldMaps("expected", typ)

	nonEssentialFields := func(varName string, typ models.DataType) []jen.Code {
		var out []jen.Code
		for _, field := range typ.Fields {
			out = append(out, jen.ID(varName).Dot(field.Name.Singular()))
		}
		return append(out, jen.ID(varName).Dot("BelongsTo"))
	}
	nef := nonEssentialFields("expected", typ)

	buildTestCreateHappyPathBody := func() []jen.Code {
		out := []jen.Code{
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(123),
				jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
				jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
				expectedInputFields...,
			),
		}

		if isPostgres {
			out = append(out,
				jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
					jen.ID("expected").Dot("ID"),
					jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (%s,belongs_to) VALUES (%s)%s", tn, strings.ReplaceAll(fieldCols, " ", ""), insertPlaceholders, queryTail),
			)
		}

		var (
			createdOnAddendum,
			createdOnValueAdd string
		)
		if isMariaDB {
			createdOnAddendum = ",created_on"
			createdOnValueAdd = ",UNIX_TIMESTAMP()"
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		)

		if isPostgres {
			out = append(out,
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Callln(
					nef...,
				).Dot("WillReturnRows").Call(jen.ID("exampleRows")),
			)
		} else if isSqlite || isMariaDB {
			out = append(out,
				jen.Line(),
				jen.ID("expectedCreationQuery").Op(":=").Litf("INSERT INTO %s (%s,belongs_to%s) VALUES (%s%s)",
					tn,
					strings.ReplaceAll(fieldCols, " ", ""),
					createdOnAddendum,
					insertPlaceholders,
					createdOnValueAdd,
				),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).
					Dotln("WithArgs").Callln(nef...).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))),
				jen.Line(),
				jen.ID("expectedTimeQuery").Op(":=").Litf("SELECT created_on FROM %s WHERE id = %s", tn, getIncIndex(dbrn, 0)),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	buildTestCreateDbErrorBody := func() []jen.Code {
		var expectFuncName string
		if isPostgres {
			expectFuncName = "ExpectQuery"
		} else if isSqlite || isMariaDB {
			expectFuncName = "ExpectExec"
		}

		var (
			createdOnAddendum,
			createdOnValueAdd string
		)
		if isMariaDB {
			createdOnAddendum = ",created_on"
			createdOnValueAdd = ",UNIX_TIMESTAMP()"
		}

		out := []jen.Code{
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(123),
				jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
				jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
				expectedInputFields...,
			),
			jen.ID("expectedQuery").Op(":=").Litf("INSERT INTO %s (%s,belongs_to%s) VALUES (%s%s)%s",
				tn,
				strings.ReplaceAll(fieldCols, " ", ""),
				createdOnAddendum,
				insertPlaceholders,
				createdOnValueAdd,
				queryTail,
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(nef...).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.Line(),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		}

		return out
	}

	ret.Add(
		jen.Func().IDf("Test%s_Create%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				buildTestCreateHappyPathBody()...,
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildTestCreateDbErrorBody()...)),
		),
		jen.Line(),
	)

	////////////

	updateCols := buildUpdateQueryParts(dbrn, typ)
	updateColsStr := strings.Join(updateCols, ", ")

	queryTail = ""
	if isPostgres {
		queryTail = " RETURNING updated_on"
	}

	testBuildUpdateQueryBody := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
			jen.ID("ID").Op(":").Lit(321),
			jen.ID("BelongsTo").Op(":").Lit(123),
		),
		jen.ID("expectedArgCount").Op(":=").Lit(len(updateCols) + 2), // +2 because of ID and BelongsTo
		jen.ID("expectedQuery").Op(":=").Litf(
			"UPDATE %s SET %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
			tn,
			updateColsStr,
			getTimeQuery(dbrn),
			getIncIndex(dbrn, uint(len(updateCols))),
			getIncIndex(dbrn, uint(len(updateCols)+1)),
			queryTail,
		),
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

	////////////

	buildExpectQueryArgs := func(varName string, typ models.DataType) []jen.Code {
		var out []jen.Code
		for _, field := range typ.Fields {
			out = append(out, jen.ID(varName).Dot(field.Name.Singular()))
		}
		return append(out, jen.ID(varName).Dot("BelongsTo"), jen.ID(varName).Dot("ID"))
	}
	expectQueryArgs := buildExpectQueryArgs("expected", typ)

	var exRows jen.Code
	queryTail = ""
	if isPostgres {
		exRows = jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		queryTail = " RETURNING updated_on"
	} else if isSqlite || isMariaDB {
		exRows = jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))
	}

	buildHappyPathUpdateBody := func() []jen.Code {
		var (
			expectFuncName,
			returnFuncName string
		)
		if isPostgres {
			expectFuncName = "ExpectQuery"
			returnFuncName = "WillReturnRows"
		} else if isSqlite || isMariaDB {
			expectFuncName = "ExpectExec"
			returnFuncName = "WillReturnResult"
		}

		out := []jen.Code{
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(123),
				jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
				jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			exRows,
			jen.ID("expectedQuery").Op(":=").Litf(
				"UPDATE %s SET %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s",
				tn,
				updateColsStr,
				getTimeQuery(dbrn),
				getIncIndex(dbrn, uint(len(updateCols))),
				getIncIndex(dbrn, uint(len(updateCols)+1)),
				queryTail,
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot(returnFuncName).Call(jen.ID("exampleRows")),
			jen.Line(),
			jen.ID("err").Op(":=").ID(dbfl).Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Line(),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		}

		return out
	}

	buildDBErrorUpdateBody := func() []jen.Code {
		var (
			expectFuncName string
		)
		if isPostgres {
			expectFuncName = "ExpectQuery"
		} else if isSqlite || isMariaDB {
			expectFuncName = "ExpectExec"
		}

		out := []jen.Code{
			jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(123),
				jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
				jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			),
			jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to = %s AND id = %s%s", tn, updateColsStr, getTimeQuery(dbrn), getIncIndex(dbrn, uint(len(updateCols))), getIncIndex(dbrn, uint(len(updateCols))+1), queryTail),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(expectQueryArgs...).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.ID("err").Op(":=").ID(dbfl).Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Line(),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		}

		return out
	}

	ret.Add(
		jen.Func().IDf("Test%s_Update%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildHappyPathUpdateBody()...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildDBErrorUpdateBody()...)),
		),
		jen.Line(),
	)

	////////////

	queryTail = ""
	if isPostgres {
		queryTail = " RETURNING archived_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_buildArchive%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
					jen.ID("BelongsTo").Op(":").Lit(123),
				),
				jen.ID("expectedArgCount").Op(":=").Lit(2),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to = %s AND id = %s%s", tn, getTimeQuery(dbrn), getTimeQuery(dbrn), getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), queryTail),
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

	////////////

	queryTail = ""
	if isPostgres {
		queryTail = " RETURNING archived_on"
	}

	ret.Add(
		jen.Func().IDf("Test%s_Archive%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to = %s AND id = %s%s", tn, getTimeQuery(dbrn), getTimeQuery(dbrn), getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), queryTail),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Callln(
					jen.ID("expected").Dot("BelongsTo"),
					jen.ID("expected").Dot("ID"),
				).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Lit(1), jen.Lit(1))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dotf("Archive%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("example").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(123),
					jen.ID("BelongsTo").Op(":").ID("expectedUserID"),
					jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.ID("expectedQuery").Op(":=").Litf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to = %s AND id = %s%s", tn, getTimeQuery(dbrn), getTimeQuery(dbrn), getIncIndex(dbrn, 0), getIncIndex(dbrn, 1), queryTail),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WithArgs").Callln(
					jen.ID("example").Dot("BelongsTo"),
					jen.ID("example").Dot("ID"),
				).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.ID("err").Op(":=").ID(dbfl).Dotf("Archive%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("example").Dot("ID"), jen.ID("expectedUserID")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.Line(),
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
		if field.Pointer {
			out = append(out,
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID(varName).Dot(field.Name.Singular()), jen.ID("args").Index(jen.Lit(i)).Assert(jen.Op("*").ID(field.Type))),
			)
		} else {
			out = append(out,
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID(varName).Dot(field.Name.Singular()), jen.ID("args").Index(jen.Lit(i)).Assert(jen.ID(field.Type))),
			)
		}
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
