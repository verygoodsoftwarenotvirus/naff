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

func iterablesTestDotGo(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	n := typ.Name
	sn := n.Singular()
	puvn := n.PluralUnexportedVarName()

	//allColumns := buildTableColumns(typ)
	gFields := buildGeneralFields("x", typ)

	ret.Add(
		jen.Func().IDf("buildMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(gFields...),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	badFields := buildBadFields("x", typ)

	ret.Add(
		jen.Func().IDf("buildErroneousMockRowFrom%s", sn).Params(jen.ID("x").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Block(
			jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(badFields...),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	ret.Add(buildTestDBBuildGetSomethingQuery(dbvendor, typ)...)
	ret.Add(buildTestDBGetSomething(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBBuildGetSomethingCountQuery(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBGetSomethingCount(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBBuildGetAllSomethingCountQuery(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBGetAllSomethingCount(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBGetListOfSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBGetListOfSomethingFuncDecl(pkg, dbvendor, typ)...)

	if typ.BelongsToUser || typ.BelongsToStruct != nil {
		ret.Add(buildTestDBGetAllSomethingForSomethingElseFuncDecl(pkg, dbvendor, typ)...)
	}

	ret.Add(buildTestDBCreateSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBCreateSomethingFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildTestBuildUpdateSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBUpdateSomethingFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBArchiveSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildTestDBArchiveSomethingFuncDecl(pkg, dbvendor, typ)...)

	return ret
}

func buildGeneralFields(varName string, typ models.DataType) []jen.Code {
	fields := []jen.Code{jen.ID(varName).Dot("ID")}

	for _, field := range typ.Fields {
		fields = append(fields, jen.ID(varName).Dot(field.Name.Singular()))
	}

	fields = append(fields,
		jen.ID(varName).Dot("CreatedOn"),
		jen.ID(varName).Dot("UpdatedOn"),
		jen.ID(varName).Dot("ArchivedOn"),
	)

	if typ.BelongsToUser {
		fields = append(fields, jen.ID(varName).Dot("BelongsToUser"))
	} else if typ.BelongsToStruct != nil {
		fields = append(fields, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

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
	)

	if typ.BelongsToUser {
		fields = append(fields, jen.ID(varName).Dot("BelongsToUser"))
	} else if typ.BelongsToStruct != nil {
		fields = append(fields, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	fields = append(fields, jen.ID(varName).Dot("ID"))

	return fields
}

func buildStringColumns(typ models.DataType) string {
	out := []string{"id"}

	for _, field := range typ.Fields {
		out = append(out, field.Name.RouteName())
	}

	out = append(out, "created_on", "updated_on", "archived_on")
	if typ.BelongsToUser {
		out = append(out, "belongs_to_user")
	} else if typ.BelongsToStruct != nil {
		out = append(out, fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))
	}

	return strings.Join(out, ", ")
}

func buildNonStandardStringColumns(typ models.DataType) string {
	var out []string

	for _, field := range typ.Fields {
		out = append(out, field.Name.RouteName())
	}

	return strings.Join(out, ", ")
}

func buildUpdateQueryParts(dbv wordsmith.SuperPalabra, typ models.DataType) []string {
	var out []string

	for i, field := range typ.Fields {
		out = append(out, fmt.Sprintf("%s = %s", field.Name.RouteName(), getIncIndex(dbv, uint(i))))
	}

	return out
}

func getIncIndex(dbv wordsmith.SuperPalabra, index uint) string {
	if isPostgres(dbv) {
		return fmt.Sprintf("$%d", index+1)
	} else if isSqlite(dbv) || isMariaDB(dbv) {
		return "?"
	}
	return ""
}

func getTimeQuery(dbvendor wordsmith.SuperPalabra) string {
	if isPostgres(dbvendor) {
		return postgresCurrentUnixTimeQuery
	} else if isSqlite(dbvendor) {
		return sqliteCurrentUnixTimeQuery
	} else if isMariaDB(dbvendor) {
		return mariaDBUnixTimeQuery
	} else {
		return ""
	}
}

func buildCreationEqualityExpectations(varName string, typ models.DataType) []jen.Code {
	var out []jen.Code

	for i, field := range typ.Fields {
		if field.Pointer {
			out = append(out,
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID(varName).Dot(field.Name.Singular()), jen.ID("args").Index(jen.Lit(i)).Assert(jen.Op("*").ID(field.Type))),
			)
		} else {
			out = append(out,
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID(varName).Dot(field.Name.Singular()), jen.ID("args").Index(jen.Lit(i)).Assert(jen.ID(field.Type))),
			)
		}
	}

	if typ.BelongsToUser {
		out = append(out,
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("BelongsToUser"), jen.ID("args").Index(jen.Lit(len(out))).Assert(jen.ID("uint64"))),
		)
	} else if typ.BelongsToStruct != nil {
		out = append(out,
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()), jen.ID("args").Index(jen.Lit(len(out))).Assert(jen.ID("uint64"))),
		)
	}

	return out
}

func buildFieldMaps(varName string, typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		xn := field.Name.Singular()
		out = append(out, jen.ID(xn).Op(":").ID(varName).Dot(xn))
	}

	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").Op(":").ID(varName).Dot("BelongsToUser"))
	} else if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	return out
}

func buildExpectQueryArgs(varName string, typ models.DataType) []jen.Code {
	var out []jen.Code
	for _, field := range typ.Fields {
		out = append(out, jen.ID(varName).Dot(field.Name.Singular()))
	}

	if typ.BelongsToUser {
		out = append(out, jen.ID(varName).Dot("BelongsToUser"), jen.ID(varName).Dot("ID"))
	} else if typ.BelongsToStruct != nil {
		out = append(out, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()), jen.ID(varName).Dot("ID"))
	}

	return out
}

func buildTestDBUpdateSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	tn := typ.Name.PluralRouteName() // table name

	var (
		expectedQuery string
		queryTail     string
	)

	if isPostgres(dbvendor) {
		queryTail = " RETURNING updated_on"
	}

	updateCols := buildUpdateQueryParts(dbvendor, typ)
	updateColsStr := strings.Join(updateCols, ", ")

	if typ.BelongsToUser {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s", tn, updateColsStr, getTimeQuery(dbvendor), getIncIndex(dbvendor, uint(len(updateCols))), getIncIndex(dbvendor, uint(len(updateCols))+1), queryTail)
	} else if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_%s = %s AND id = %s%s", tn, updateColsStr, getTimeQuery(dbvendor), typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, uint(len(updateCols))), getIncIndex(dbvendor, uint(len(updateCols))+1), queryTail)
	} else {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE id = %s%s", tn, updateColsStr, getTimeQuery(dbvendor), getIncIndex(dbvendor, uint(len(updateCols))), queryTail)
	}

	buildFirstSubTest := func(typ models.DataType) []jen.Code {
		var (
			expectFuncName,
			returnFuncName string
			exRows jen.Code
		)

		dbrn := dbvendor.RouteName()
		sn := typ.Name.Singular()
		dbfl := string(dbrn[0])

		expectQueryArgs := buildExpectQueryArgs("expected", typ)

		if isPostgres(dbvendor) {
			expectFuncName = "ExpectQuery"
			returnFuncName = "WillReturnRows"

			exRows = jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("updated_on"))).Dot("AddRow").Call(jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			expectFuncName = "ExpectExec"
			returnFuncName = "WillReturnResult"

			exRows = jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))
		}

		lines := []jen.Code{}
		expectedValues := []jen.Code{jen.ID("ID").Op(":").Lit(123)}
		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		lines = append(lines,
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
			exRows,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot(returnFuncName).Call(jen.ID("exampleRows")),
			jen.Line(),
			jen.ID("err").Op(":=").ID(dbfl).Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtest := func(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
		dbrn := dbvendor.RouteName()
		sn := typ.Name.Singular()
		dbfl := string(dbrn[0])

		var (
			expectFuncName string
		)
		if isPostgres(dbvendor) {
			expectFuncName = "ExpectQuery"
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			expectFuncName = "ExpectExec"
		}

		out := []jen.Code{}

		expectQueryArgs := buildExpectQueryArgs("expected", typ)
		expectedValues := []jen.Code{
			jen.ID("ID").Op(":").Lit(123),
		}

		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").ID("expectedUserID"))
			out = append(out, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
		} else if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").IDf("expected%sID", typ.BelongsToStruct.Singular()))
			out = append(out, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
		}

		expectedValues = append(expectedValues,
			jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
		)

		out = append(out,
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(expectQueryArgs...).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.ID("err").Op(":=").ID(dbfl).Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expected")),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Update%s", dbvendor.Singular(), typ.Name.Singular()).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildFirstSubTest(typ)...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSecondSubtest(pkg, dbvendor, typ)...)),
		),
		jen.Line(),
	}
}

func buildTestDBArchiveSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	tn := typ.Name.PluralRouteName() // table name

	var (
		expectedQuery string
		queryTail     string
		queryArgCount int
	)

	if isPostgres(dbvendor) {
		queryTail = " RETURNING archived_on"
	}
	expectedValues := []jen.Code{
		jen.ID("ID").Op(":").Lit(321),
	}
	archiveQueryBuildingParams := []jen.Code{
		jen.ID("expected").Dot("ID"),
	}

	if typ.BelongsToUser {
		queryArgCount = 2
		expectedQuery = fmt.Sprintf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s%s", tn, getTimeQuery(dbvendor), getTimeQuery(dbvendor), getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), queryTail)
		expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").Lit(123))
		archiveQueryBuildingParams = append(archiveQueryBuildingParams, jen.ID("expected").Dot("BelongsToUser"))
	} else if typ.BelongsToStruct != nil {
		queryArgCount = 2
		expectedQuery = fmt.Sprintf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_%s = %s AND id = %s%s", tn, getTimeQuery(dbvendor), getTimeQuery(dbvendor), typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), queryTail)
		expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").Lit(123))
		archiveQueryBuildingParams = append(archiveQueryBuildingParams, jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	} else {
		queryArgCount = 1
		expectedQuery = fmt.Sprintf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND id = %s%s", tn, getTimeQuery(dbvendor), getTimeQuery(dbvendor), getIncIndex(dbvendor, 0), queryTail)
	}

	testLines := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
		jen.ID("expectedArgCount").Op(":=").Lit(queryArgCount),
		jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildArchive%sQuery", sn).Call(archiveQueryBuildingParams...),
		jen.Line(),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	}

	if typ.BelongsToUser {
		testLines = append(testLines,
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("BelongsToUser"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)
	} else if typ.BelongsToStruct != nil {
		testLines = append(testLines,
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))),
		)
	}

	testLines = append(testLines,
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
	)

	return []jen.Code{
		jen.Func().IDf("Test%s_buildArchive%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(testLines...)),
		),
		jen.Line(),
	}
}

func buildTestDBArchiveSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	tn := typ.Name.PluralRouteName() // table name

	buildSubtestOne := func() []jen.Code {
		var (
			dbQueryExpectationArgs []jen.Code
			dbQuery                string
			queryTail              string
		)

		if isPostgres(dbvendor) {
			queryTail = " RETURNING archived_on"
		}

		expectedValues := []jen.Code{
			jen.ID("ID").Op(":").Lit(123),
		}
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.ID("expected").Dot("ID"),
		}

		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").ID("expectedUserID"))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").IDf("expected%sID", typ.BelongsToStruct.Singular()))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}

		expectedValues = append(expectedValues, jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		block := []jen.Code{}

		if typ.BelongsToUser {
			dbQuery = fmt.Sprintf(
				"UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s%s",
				tn,
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				queryTail,
			)
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("expected").Dot("BelongsToUser"))
			block = append(block, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
		} else if typ.BelongsToStruct != nil {
			dbQuery = fmt.Sprintf(
				"UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_%s = %s AND id = %s%s",
				tn,
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				typ.BelongsToStruct.RouteName(),
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				queryTail,
			)
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
			block = append(block, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
		} else {
			dbQuery = fmt.Sprintf(
				"UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND id = %s%s",
				tn,
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 0),
				queryTail,
			)
		}

		dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("expected").Dot("ID"))

		block = append(block,
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.Lit(1), jen.Lit(1))),
			jen.Line(),
			jen.ID("err").Op(":=").ID(dbfl).Dotf("Archive%s", sn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return []jen.Code{
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Op(":=").Lit(dbQuery),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(block...)),
		}
	}

	buildSubtestTwo := func() []jen.Code {
		exampleValues := []jen.Code{
			jen.ID("ID").Op(":").Lit(123),
		}

		var dbQueryExpectationArgs []jen.Code
		block := []jen.Code{}
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.ID("example").Dot("ID"),
		}

		if typ.BelongsToUser {
			exampleValues = append(exampleValues, jen.ID("BelongsToUser").Op(":").ID("expectedUserID"))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("example").Dot("BelongsToUser"))
			block = append(block, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			exampleValues = append(exampleValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").IDf("expected%sID", typ.BelongsToStruct.Singular()))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("example").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
			block = append(block, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}

		dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("example").Dot("ID"))
		exampleValues = append(exampleValues, jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		block = append(block,
			jen.ID("example").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(exampleValues...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).
				Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.ID("err").Op(":=").ID(dbfl).Dotf("Archive%s", sn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(
				jen.ID("t"),
				jen.ID("mockDB").Dot("ExpectationsWereMet").Call(),
				jen.Lit("not all database expectations were met"),
			),
		)

		return []jen.Code{
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(block...)),
		}
	}

	var bodyContents []jen.Code

	bodyContents = append(bodyContents, buildSubtestOne()...)
	bodyContents = append(bodyContents, jen.Line())
	bodyContents = append(bodyContents, buildSubtestTwo()...)

	return []jen.Code{
		jen.Func().IDf("Test%s_Archive%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			bodyContents...,
		),
		jen.Line(),
	}
}

func buildTestBuildUpdateSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name

	updateCols := buildUpdateQueryParts(dbvendor, typ)
	updateColsStr := strings.Join(updateCols, ", ")
	creationEqualityExpectations := buildCreationEqualityExpectations("expected", typ)

	var (
		expectedQuery string
		queryTail     string
		varCount      int
	)

	if isPostgres(dbvendor) {
		queryTail = " RETURNING updated_on"
	}

	expectedValues := []jen.Code{jen.ID("ID").Op(":").Lit(321)}

	if typ.BelongsToUser {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s",
			tn,
			updateColsStr,
			getTimeQuery(dbvendor),
			getIncIndex(dbvendor, uint(len(updateCols))),
			getIncIndex(dbvendor, uint(len(updateCols)+1)),
			queryTail,
		)
		expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").Lit(123))
		varCount = len(updateCols) + 2
	} else if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_%s = %s AND id = %s%s",
			tn,
			updateColsStr,
			getTimeQuery(dbvendor),
			typ.BelongsToStruct.RouteName(),
			getIncIndex(dbvendor, uint(len(updateCols))),
			getIncIndex(dbvendor, uint(len(updateCols)+1)),
			queryTail,
		)
		expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").Lit(123))
		varCount = len(updateCols) + 2
	} else {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE id = %s%s",
			tn,
			updateColsStr,
			getTimeQuery(dbvendor),
			getIncIndex(dbvendor, uint(len(updateCols))),
			queryTail,
		)
		varCount = len(updateCols) + 1
	}

	testBuildUpdateQueryBody := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
		jen.ID("expectedArgCount").Op(":=").Lit(varCount), // +2 because of ID and BelongsTo
		jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildUpdate%sQuery", sn).Call(jen.ID("expected")),
		jen.Line(),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	}

	testBuildUpdateQueryBody = append(testBuildUpdateQueryBody, creationEqualityExpectations...)
	testBuildUpdateQueryBody = append(testBuildUpdateQueryBody,
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(len(creationEqualityExpectations))).Assert(jen.ID("uint64"))),
	)

	return []jen.Code{
		jen.Func().IDf("Test%s_buildUpdate%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				testBuildUpdateQueryBody...,
			)),
		),
		jen.Line(),
	}
}

func buildTestDBCreateSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name

	fieldCols := buildNonStandardStringColumns(typ)

	var ips []string
	for i := range typ.Fields {
		ips = append(ips, getIncIndex(dbvendor, uint(i)))
	}
	ips = append(ips, getIncIndex(dbvendor, uint(len(ips))))
	insertPlaceholders := strings.Join(ips, ",")

	var queryTail string
	if isPostgres(dbvendor) {
		queryTail = " RETURNING id, created_on"
	}

	var (
		expectedQuery,
		createdOnAddendum,
		createdOnValueAdd string
	)
	if isMariaDB(dbvendor) {
		createdOnAddendum = ",created_on"
		createdOnValueAdd = ",UNIX_TIMESTAMP()"
	}

	thisFuncExpectedArgCount := len(ips) - 1

	expectedValues := []jen.Code{jen.ID("ID").Op(":").Lit(321)}
	if typ.BelongsToUser {
		expectedQuery = fmt.Sprintf("INSERT INTO %s (%s,belongs_to_user%s) VALUES (%s%s)%s",
			tn,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
		expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").Lit(123))
	} else if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("INSERT INTO %s (%s,belongs_to_%s%s) VALUES (%s%s)%s",
			tn,
			strings.ReplaceAll(fieldCols, " ", ""),
			typ.BelongsToStruct.RouteName(),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
		expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").Lit(123))
	} else {
		expectedQuery = fmt.Sprintf("INSERT INTO %s (%s%s) VALUES (%s%s)%s",
			tn,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	}

	creationEqualityExpectations := buildCreationEqualityExpectations("expected", typ)
	createQueryTestBody := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
		jen.ID("expectedArgCount").Op(":=").Lit(1 + thisFuncExpectedArgCount),
		jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildCreate%sQuery", sn).Call(jen.ID("expected")),
		jen.Line(),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	}
	createQueryTestBody = append(createQueryTestBody, creationEqualityExpectations...)

	return []jen.Code{

		jen.Func().IDf("Test%s_buildCreate%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(createQueryTestBody...)),
		),
		jen.Line(),
	}
}

func buildTestDBCreateSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name

	var (
		ips []string
		createdOnAddendum,
		createdOnValueAdd,
		queryTail,
		expectedCreationQuery string
	)

	for i := range typ.Fields {
		ips = append(ips, getIncIndex(dbvendor, uint(i)))
	}
	ips = append(ips, getIncIndex(dbvendor, uint(len(ips))))
	insertPlaceholders := strings.Join(ips, ",")

	if isPostgres(dbvendor) {
		queryTail = " RETURNING id, created_on"
	} else if isMariaDB(dbvendor) {
		createdOnAddendum = ",created_on"
		createdOnValueAdd = ",UNIX_TIMESTAMP()"
	}

	expectedInputFields := buildFieldMaps("expected", typ)
	fieldCols := buildNonStandardStringColumns(typ)

	if typ.BelongsToUser {
		expectedCreationQuery = fmt.Sprintf(
			"INSERT INTO %s (%s,belongs_to_user%s) VALUES (%s%s)%s",
			tn,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	} else if typ.BelongsToStruct != nil {
		expectedCreationQuery = fmt.Sprintf(
			"INSERT INTO %s (%s,belongs_to_%s%s) VALUES (%s%s)%s",
			tn,
			strings.ReplaceAll(fieldCols, " ", ""),
			typ.BelongsToStruct.RouteName(),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	} else {
		// todo
		expectedCreationQuery = fmt.Sprintf(
			"INSERT INTO %s (%s%s) VALUES (%s%s)%s",
			tn,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	}

	buildNonEssentialFields := func(varName string, typ models.DataType) []jen.Code {
		var out []jen.Code
		for _, field := range typ.Fields {
			out = append(out, jen.ID(varName).Dot(field.Name.Singular()))
		}

		if typ.BelongsToUser {
			out = append(out, jen.ID(varName).Dot("BelongsToUser"))
		} else if typ.BelongsToStruct != nil {
			out = append(out, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		}

		return out
	}

	buildFirstSubtest := func(pkg *models.Project, typ models.DataType) []jen.Code {
		out := []jen.Code{}
		expectedValues := []jen.Code{jen.ID("ID").Op(":").Lit(123)}

		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").ID("expectedUserID"))
			out = append(out, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
		} else if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").IDf("expected%sID", typ.BelongsToStruct.Singular()))
			out = append(out, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out = append(out,
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
			jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(expectedInputFields...),
		)

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("exampleRows").Op(":=").Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
					jen.ID("expected").Dot("ID"),
					jen.ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		)

		nef := buildNonEssentialFields("expected", typ)

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).
					Dotln("WithArgs").Callln(
					nef...,
				).Dot("WillReturnRows").Call(jen.ID("exampleRows")),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			out = append(out,
				jen.Line(),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).
					Dotln("WithArgs").Callln(nef...).Dot("WillReturnResult").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(jen.ID("int64").Call(jen.ID("expected").Dot("ID")), jen.Lit(1))),
				jen.Line(),
				jen.ID("expectedTimeQuery").Op(":=").Litf("SELECT created_on FROM %s WHERE id = %s", tn, getIncIndex(dbvendor, 0)),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("created_on"))).Dot("AddRow").Call(jen.ID("expected").Dot("CreatedOn"))),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	buildSecondSubtest := func() []jen.Code {
		var expectFuncName string
		if isPostgres(dbvendor) {
			expectFuncName = "ExpectQuery"
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			expectFuncName = "ExpectExec"
		}
		nef := buildNonEssentialFields("expected", typ)

		expectedValues := []jen.Code{jen.ID("ID").Op(":").Lit(123)}

		out := []jen.Code{}
		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").Op(":").ID("expectedUserID"))
			out = append(out, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
		} else if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Op(":").IDf("expected%sID", typ.BelongsToStruct.Singular()))
			out = append(out, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out = append(out,
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(expectedValues...),
			jen.ID("expectedInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(expectedInputFields...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).
				Dotln("WithArgs").Callln(nef...).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("expectedInput")),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Create%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedCreationQuery").Op(":=").Lit(expectedCreationQuery),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildFirstSubtest(pkg, typ)...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error writing to database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSecondSubtest()...)),
		),
		jen.Line(),
	}
}

func buildTestDBGetAllSomethingForSomethingElseFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name
	cols := buildStringColumns(typ)

	var (
		baseFuncName        string
		testFuncName        string
		expectedQuery       string
		expectedSomethingID string
	)

	if typ.BelongsToUser {
		expectedSomethingID = "expectedUserID"
		baseFuncName = fmt.Sprintf("GetAll%sForUser", pn)
		testFuncName = fmt.Sprintf("Test%s_%s", dbvsn, baseFuncName)
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_user = %s", cols, tn, getIncIndex(dbvendor, 0))
	} else if typ.BelongsToStruct != nil {
		expectedSomethingID = fmt.Sprintf("expected%sID", typ.BelongsToStruct.Singular())
		baseFuncName = fmt.Sprintf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular())
		testFuncName = fmt.Sprintf("Test%s_%s", dbvsn, baseFuncName)
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_%s = %s", cols, tn, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0))
	}
	// we don't need to consider the case where this object belongs to nothing

	return []jen.Code{
		jen.Func().ID(testFuncName).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Op(":=").Lit(expectedQuery),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID(expectedSomethingID).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("expected%s", sn).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn))),
				jen.Line(),
				jen.ID("expected").Op(":=").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(jen.Op("*").IDf("expected%s", sn)),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot(baseFuncName).Call(jen.Qual("context", "Background").Call(), jen.ID(expectedSomethingID)),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID(expectedSomethingID).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot(baseFuncName).Call(jen.Qual("context", "Background").Call(), jen.ID(expectedSomethingID)),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID(expectedSomethingID).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot(baseFuncName).Call(jen.Qual("context", "Background").Call(), jen.ID(expectedSomethingID)),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with unscannable response"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID(expectedSomethingID).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("example%s", sn).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.ID("ID").Op(":").Lit(321),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.IDf("example%s", sn))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dot(baseFuncName).Call(jen.Qual("context", "Background").Call(), jen.ID(expectedSomethingID)),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	}
}

func buildTestDBGetListOfSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name
	cols := buildStringColumns(typ)

	var (
		expectedQuery string
	)

	if typ.BelongsToUser {
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", cols, tn, getIncIndex(dbvendor, 0))
	} else if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_%s = %s LIMIT 20", cols, tn, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0))
	} else {
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL LIMIT 20", cols, tn)
	}

	buildFirstSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var expectQueryMock jen.Code
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUserID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn)))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn)))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.IDf("expected%s", sn)))
		}

		lines = append(lines,
			jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
			jen.IDf("expected%s", sn).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(321),
			),
			jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
					jen.ID("Page").Op(":").Lit(1),
					jen.ID("Limit").Op(":").Lit(20),
					jen.ID("TotalCount").Op(":").ID("expectedCount"),
				),
				jen.ID(pn).Op(":").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
					jen.Op("*").IDf("expected%s", sn),
				),
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			expectQueryMock,
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).
				Dot("AddRow").Call(jen.ID("expectedCount"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUserID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildThirdSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUserID")).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildFourthSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUserID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected")))
		}

		lines = append(lines,
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(321),
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildFifthSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
			jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUserID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected")))
		}

		lines = append(lines,
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(321),
			),
			jen.ID("expectedCountQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Get%s", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Op(":=").Lit(expectedQuery),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildFirstSubtest()...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSecondSubtest()...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error executing read query"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildThirdSubtest()...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Litf("with error scanning %s", scn), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildFourthSubtest()...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error querying for count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildFifthSubtest()...)),
		),
		jen.Line(),
	}
}

func buildTestDBGetListOfSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name
	cols := buildStringColumns(typ)

	var (
		expectedQuery   string
		expectedOwnerID string
	)

	bodyBlock := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
	}

	if typ.BelongsToUser {
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", cols, tn, getIncIndex(dbvendor, 0))
		expectedOwnerID = "exampleUserID"
		bodyBlock = append(bodyBlock, jen.ID(expectedOwnerID).Op(":=").ID("uint64").Call(jen.Lit(321)))
	} else if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_%s = %s LIMIT 20", cols, tn, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0))
		expectedOwnerID = fmt.Sprintf("example%sID", typ.BelongsToStruct.Singular())
		bodyBlock = append(bodyBlock, jen.IDf(expectedOwnerID).Op(":=").ID("uint64").Call(jen.Lit(321)))
	} else if typ.BelongsToNobody {
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL LIMIT 20", cols, tn)
	}

	bodyBlock = append(bodyBlock,
		jen.Line(),
		jen.ID("expectedArgCount").Op(":=").Lit(1),
		jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", pn).Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID(expectedOwnerID)),
		jen.Line(),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	)

	if typ.BelongsToUser || typ.BelongsToStruct != nil {
		bodyBlock = append(bodyBlock, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID(expectedOwnerID), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))))
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_buildGet%sQuery", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(bodyBlock...)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBBuildGetSomethingQuery(dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name
	cols := buildStringColumns(typ)

	var query string
	if typ.BelongsToUser {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE belongs_to_user = %s AND id = %s", cols, tn, getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1))
	} else if typ.BelongsToStruct != nil {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE belongs_to_%s = %s AND id = %s", cols, tn, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1))
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s id = %s", cols, tn, getIncIndex(dbvendor, 0))
	}

	block := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
	}

	if typ.BelongsToUser {
		block = append(block, jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
	}

	block = append(block,
		jen.Line(),
		jen.ID("expectedArgCount").Op(":=").Lit(2),
		jen.ID("expectedQuery").Op(":=").Lit(query),
	)

	if typ.BelongsToUser {
		block = append(block, jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", sn).Call(jen.IDf("example%sID", sn), jen.ID("exampleUserID")))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", sn).Call(jen.IDf("example%sID", sn), jen.IDf("example%sID", typ.BelongsToStruct.Singular())))
	} else if typ.BelongsToNobody {
		block = append(block, jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", sn).Call(jen.IDf("example%sID", sn)))
	}

	block = append(block,
		jen.Line(),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	)

	if typ.BelongsToUser {
		block = append(block, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.IDf("example%sID", typ.BelongsToStruct.Singular()), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))))
	}

	block = append(block,
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.IDf("example%sID", sn), jen.ID("args").Index(jen.Lit(1)).Assert(jen.ID("uint64"))),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_buildGet%sQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(block...)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBGetSomething(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name
	cols := buildStringColumns(typ)

	var query string
	if typ.BelongsToUser {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE belongs_to_user = %s AND id = %s", cols, tn, getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1))
	} else if typ.BelongsToStruct != nil {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE belongs_to_%s = %s AND id = %s", cols, tn, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1))
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s id = %s", cols, tn, getIncIndex(dbvendor, 0))
	}

	buildFirstSubtestBlock := func(typ models.DataType) []jen.Code {
		lines := []jen.Code{
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(123),
			),
		}

		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"),
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular()), jen.ID("expected").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowFrom%s", sn).Call(jen.ID("expected")))
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", sn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtestBlock := func(typ models.DataType) []jen.Code {
		lines := []jen.Code{
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
				jen.ID("ID").Op(":").Lit(123),
			),
		}

		actualCallArgs := []jen.Code{
			jen.Qual("context", "Background").Call(), jen.ID("expected").Dot("ID"),
		}
		var mockDBCall jen.Code
		if typ.BelongsToUser {
			lines = append(lines, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUserID"), jen.ID("expected").Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular()), jen.ID("expected").Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%s", sn).Call(actualCallArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("actual")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.Qual("database/sql", "ErrNoRows"), jen.ID("err")),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_Get%s", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Op(":=").Lit(query),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildFirstSubtestBlock(typ)...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("surfaces sql.ErrNoRows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSecondSubtestBlock(typ)...)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBBuildGetSomethingCountQuery(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name

	var (
		query            string
		expectedArgCount int
	)

	block := []jen.Code{
		jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
	}
	actualCallArgs := []jen.Code{
		jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
	}

	if typ.BelongsToUser {
		expectedArgCount = 1
		query = fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", tn, getIncIndex(dbvendor, 0))
		block = append(block, jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
		actualCallArgs = append(actualCallArgs, jen.ID("exampleUserID"))
	} else if typ.BelongsToStruct != nil {
		expectedArgCount = 1
		query = fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to_%s = %s LIMIT 20", tn, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0))
		block = append(block, jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
		actualCallArgs = append(actualCallArgs, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
	} else {
		query = fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL LIMIT 20", tn)
	}

	block = append(block,
		jen.Line(),
		jen.ID("expectedArgCount").Op(":=").Lit(expectedArgCount),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.Line(),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sCountQuery", sn).Call(actualCallArgs...),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("args"), jen.ID("expectedArgCount")),
	)

	if typ.BelongsToUser {
		block = append(block, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("exampleUserID"), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.IDf("example%sID", typ.BelongsToStruct.Singular()), jen.ID("args").Index(jen.Lit(0)).Assert(jen.ID("uint64"))))
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_buildGet%sCountQuery", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(block...)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBGetSomethingCount(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name

	var (
		query      string
		mockDBCall jen.Code
	)
	block := []jen.Code{}

	callArgs := []jen.Code{
		jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(),
	}

	if typ.BelongsToUser {
		query = fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to_user = %s LIMIT 20", tn, getIncIndex(dbvendor, 0))
		block = append(block, jen.ID("expectedUserID").Op(":=").ID("uint64").Call(jen.Lit(321)))
		mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WithArgs").Call(jen.ID("expectedUserID")).
			Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount")))
		callArgs = append(callArgs, jen.ID("expectedUserID"))
	} else if typ.BelongsToStruct != nil {
		query = fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL AND belongs_to_%s = %s LIMIT 20", tn, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0))
		block = append(block, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(321)))
		mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
			Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount")))
		callArgs = append(callArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
	} else if typ.BelongsToNobody {
		mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount")))
		query = fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL LIMIT 20", tn)
	}

	block = append(block,
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.ID("expectedCount").Op(":=").ID("uint64").Call(jen.Lit(666)),
		jen.Line(),
		jen.List(jen.ID(dbfl), jen.ID("mockDB")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		mockDBCall,
		jen.Line(),
		jen.List(jen.ID("actualCount"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%sCount", sn).Call(callArgs...),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
		jen.Line(),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_Get%sCount", dbvsn, sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(block...)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBBuildGetAllSomethingCountQuery(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name

	lines := []jen.Code{
		jen.Func().IDf("Test%s_buildGetAll%sCountQuery", dbvsn, pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID(dbfl), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expectedQuery").Op(":=").Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tn),
				jen.Line(),
				jen.ID("actualQuery").Op(":=").ID(dbfl).Dotf("buildGetAll%sCountQuery", pn).Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBGetAllSomethingCount(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tn := typ.Name.PluralRouteName() // table name

	lines := []jen.Code{
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
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedCount"), jen.ID("actualCount")),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			)),
		),
		jen.Line(),
	}

	return lines
}
