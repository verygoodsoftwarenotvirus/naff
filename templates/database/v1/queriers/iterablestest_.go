package queriers

import (
	"fmt"
	"log"
	"strings"

	"github.com/Masterminds/squirrel"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	postgresCurrentUnixTimeQuery = "extract(epoch FROM NOW())"
	sqliteCurrentUnixTimeQuery   = "(strftime('%s','now'))"
	mariaDBUnixTimeQuery         = "UNIX_TIMESTAMP()"
)

func applyFleshedOutQueryFilter(qb squirrel.SelectBuilder, tableName string) squirrel.SelectBuilder {
	qb = qb.
		Where(squirrel.Gt{fmt.Sprintf("%s.created_on", tableName): whateverValue}).
		Where(squirrel.Lt{fmt.Sprintf("%s.created_on", tableName): whateverValue}).
		Where(squirrel.Gt{fmt.Sprintf("%s.updated_on", tableName): whateverValue}).
		Where(squirrel.Lt{fmt.Sprintf("%s.updated_on", tableName): whateverValue}).
		GroupBy(fmt.Sprintf("%s.id", tableName)).
		Limit(20).
		Offset(180)

	return qb
}

func applyFleshedOutQueryFilterWithCode(qb squirrel.SelectBuilder, tableName string, eq squirrel.Eq) squirrel.SelectBuilder {
	if eq != nil {
		qb = qb.Where(eq)
	}
	qb = qb.
		Where(squirrel.Gt{fmt.Sprintf("%s.created_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("CreatedAfter"),
		)}).
		Where(squirrel.Lt{fmt.Sprintf("%s.created_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("CreatedBefore"),
		)}).
		Where(squirrel.Gt{fmt.Sprintf("%s.updated_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("UpdatedAfter"),
		)}).
		Where(squirrel.Lt{fmt.Sprintf("%s.updated_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("UpdatedBefore"),
		)}).
		GroupBy(fmt.Sprintf("%s.id", tableName)).
		Limit(20).
		Offset(180)

	return qb
}

func appendFleshedOutQueryFilterArgs(args []jen.Code) []jen.Code {
	args = append(args,
		jen.ID(constants.FilterVarName).Dot("CreatedAfter"),
		jen.ID(constants.FilterVarName).Dot("CreatedBefore"),
		jen.ID(constants.FilterVarName).Dot("UpdatedAfter"),
		jen.ID(constants.FilterVarName).Dot("UpdatedBefore"),
	)

	return args
}

func buildGeneralFields(varName string, typ models.DataType) []jen.Code {
	fields := []jen.Code{jen.ID(varName).Dot("ID")}

	for _, field := range typ.Fields {
		fields = append(fields, jen.ID(varName).Dot(field.Name.Singular()))
	}

	fields = append(fields,
		jen.ID(varName).Dot("CreatedOn"),
		jen.ID(varName).Dot("LastUpdatedOn"),
		jen.ID(varName).Dot("ArchivedOn"),
	)

	if typ.BelongsToUser {
		fields = append(fields, jen.ID(varName).Dot(constants.UserOwnershipFieldName))
	}
	if typ.BelongsToStruct != nil {
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
		jen.ID(varName).Dot("LastUpdatedOn"),
	)

	if typ.BelongsToStruct != nil {
		fields = append(fields, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToUser {
		fields = append(fields, jen.ID(varName).Dot(constants.UserOwnershipFieldName))
	}

	fields = append(fields, jen.ID(varName).Dot("ID"))

	return fields
}

func buildPrefixedStringColumns(typ models.DataType) []string {
	tableName := typ.Name.PluralRouteName()
	out := []string{fmt.Sprintf("%s.id", tableName)}

	for _, field := range typ.Fields {
		out = append(out, fmt.Sprintf("%s.%s", tableName, field.Name.RouteName()))
	}

	out = append(out, fmt.Sprintf("%s.created_on", tableName), fmt.Sprintf("%s.updated_on", tableName), fmt.Sprintf("%s.archived_on", tableName))
	if typ.BelongsToUser {
		out = append(out, fmt.Sprintf("%s.belongs_to_user", tableName))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName()))
	}

	return out
}

func buildPrefixedStringColumnsAsString(typ models.DataType) string {
	return strings.Join(buildPrefixedStringColumns(typ), ", ")
}

func buildCreationStringColumnsAndArgs(proj *models.Project, typ models.DataType) (cols []string, args []jen.Code) {
	cols, args = []string{}, []jen.Code{}

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			cols = append(cols, field.Name.RouteName())
			args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToStruct != nil {
		cols = append(cols, fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))
		args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToUser {
		cols = append(cols, "belongs_to_user")
		args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot(constants.UserOwnershipFieldName))
	}

	return
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

func iterablesTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	n := typ.Name
	sn := n.Singular()
	puvn := n.PluralUnexportedVarName()

	gFields := buildGeneralFields("x", typ)

	code.Add(
		jen.Func().IDf("buildMockRowsFrom%s", sn).Params(
			jen.ID(puvn).Spread().PointerTo().Qual(proj.ModelsV1Package(), sn),
		).Params(
			jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Block(
			jen.ID("includeCount").Assign().Len(jen.ID(puvn)).GreaterThan().One(),
			jen.ID("columns").Assign().IDf("%sTableColumns", puvn),
			jen.Line(),
			jen.If(jen.ID("includeCount")).Block(
				jen.ID("columns").Equals().Append(jen.ID("columns"), jen.Lit("count")),
			),
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For().List(jen.Underscore(), jen.ID("x")).Assign().Range().ID(puvn).Block(
				jen.ID("rowValues").Assign().Index().Qual("database/sql/driver", "Value").Valuesln(gFields...),
				jen.Line(),
				jen.If(jen.ID("includeCount")).Block(
					utils.AppendItemsToList(jen.ID("rowValues"), jen.Len(jen.ID(puvn))),
				),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Rows")).Dot("AddRow").Call(jen.ID("rowValues").Spread()),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	badFields := buildBadFields("x", typ)

	code.Add(
		jen.Func().IDf("buildErroneousMockRowFrom%s", sn).Params(
			jen.ID("x").PointerTo().Qual(proj.ModelsV1Package(), sn),
		).Params(
			jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Block(
			jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(badFields...),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	)

	code.Add(buildTestScanListOfThings(proj, dbvendor, typ)...)
	code.Add(buildTestDBBuildSomethingExistsQuery(proj, dbvendor, typ)...)
	code.Add(buildTestDBSomethingExists(proj, dbvendor, typ)...)
	code.Add(buildTestDBBuildGetSomethingQuery(proj, dbvendor, typ)...)
	code.Add(buildTestDBGetSomething(proj, dbvendor, typ)...)
	code.Add(buildTestDBBuildGetAllSomethingCountQuery(proj, dbvendor, typ)...)
	code.Add(buildTestDBGetAllSomethingCount(proj, dbvendor, typ)...)
	code.Add(buildTestDBGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildTestDBCreateSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildTestDBCreateSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildTestBuildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildTestDBUpdateSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildTestDBArchiveSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildTestDBArchiveSomethingFuncDecl(proj, dbvendor, typ)...)

	return code
}

func buildTestScanListOfThings(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbv := dbvendor.Singular()
	pn := typ.Name.Plural()
	dbfl := strings.ToLower(string([]byte(dbv)[0]))

	lines := []jen.Code{
		jen.Func().IDf("Test%s_Scan%s", dbv, pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces row errors",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabaseV1Package(), "MockResultIterator").Values(),
				jen.Line(),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(
					jen.Underscore(),
					jen.Underscore(),
					jen.Err(),
				).Assign().ID(dbfl).Dotf("scan%s", typ.Name.Plural()).Call(jen.ID("mockRows")),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"logs row closing errors",
				jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabaseV1Package(), "MockResultIterator").Values(),
				jen.Line(),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Nil()),
				jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(constants.ObligatoryError()),
				jen.Line(),
				jen.List(
					jen.Underscore(),
					jen.Underscore(),
					jen.Err(),
				).Assign().ID(dbfl).Dotf("scan%s", typ.Name.Plural()).Call(jen.ID("mockRows")),
				utils.AssertNoError(jen.Err(), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBBuildSomethingExistsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	whereValues := typ.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(proj)
	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Suffix(existenceSuffix).
		Where(whereValues)

	callArgs := typ.BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(proj)
	pql := typ.BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("%sExists", sn), qb, nil, callArgs, pql)
}

func buildTestDBSomethingExists(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	eqArgs := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): whateverValue,
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		eqArgs[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = whateverValue
	}

	whereValues := typ.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(proj)
	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Suffix(existenceSuffix).
		Where(whereValues)

	query, args, _ := qb.ToSql()
	actualCallArgs := typ.BuildArgsForDBQuerierExistenceMethodTest(proj)
	mockDBCallArgs := convertArgsToCode(args)

	buildFirstSubtestBlock := func(typ models.DataType) []jen.Code {
		lines := typ.BuildDependentObjectsForDBQueriersExistenceMethodTest(proj)

		mockDBCall := jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WithArgs").Callln(mockDBCallArgs...).
			Dotln("WillReturnRows").
			Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("%sExists", sn).Call(actualCallArgs...),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertTrue(jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtestBlock := func(typ models.DataType) []jen.Code {
		lines := typ.BuildDependentObjectsForDBQueriersExistenceMethodTest(proj)

		mockDBCall := jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WithArgs").Callln(mockDBCallArgs...).
			Dotln("WillReturnError").
			Call(jen.Qual("database/sql", "ErrNoRows"))

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("%sExists", sn).Call(actualCallArgs...),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertFalse(jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_%sExists", dbvsn, sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(query),
			jen.Line(),
			utils.BuildSubTest("happy path", buildFirstSubtestBlock(typ)...),
			jen.Line(),
			utils.BuildSubTest("with no rows", buildSecondSubtestBlock(typ)...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBBuildGetSomethingQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	whereValues := typ.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(proj)
	qb := queryBuilderForDatabase(dbvendor).Select(buildPrefixedStringColumnsAsString(typ)).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Where(whereValues)

	callArgs := typ.BuildDBQuerierRetrievalQueryTestCallArgs(proj)
	pql := typ.BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Get%s", sn), qb, []jen.Code{}, callArgs, pql)
}

func buildTestDBGetSomething(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	whereValues := typ.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(proj)
	qb := queryBuilderForDatabase(dbvendor).Select(buildPrefixedStringColumnsAsString(typ)).
		From(tableName)
	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)
	qb = qb.Where(whereValues)
	query, args, _ := qb.ToSql()
	mockDBCallArgs := convertArgsToCode(args)
	actualCallArgs := typ.BuildArgsForDBQuerierRetrievalMethodTest(proj)

	buildFirstSubtestBlock := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(proj)

		mockDBCall := jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WithArgs").Callln(mockDBCallArgs...).
			Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))))

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", sn).Call(actualCallArgs...),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtestBlock := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(proj)

		mockDBCall := jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
			Dotln("WithArgs").Callln(mockDBCallArgs...).
			Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", sn).Call(actualCallArgs...),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_Get%s", dbvsn, sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			func() jen.Code {
				if typ.OwnedByAUserAtSomeLevel(proj) {
					return jen.ID(utils.BuildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call()
				}
				return jen.Null()
			}(),
			jen.ID("expectedQuery").Assign().Lit(query),
			jen.Line(),
			utils.BuildSubTest("happy path", buildFirstSubtestBlock()...),
			jen.Line(),
			utils.BuildSubTest("surfaces sql.ErrNoRows", buildSecondSubtestBlock()...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBBuildGetAllSomethingCountQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	pn := typ.Name.Plural()

	qb := queryBuilderForDatabase(dbvendor).
		Select(fmt.Sprintf(countQuery, tableName)).
		From(tableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", tableName): nil,
		})

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("GetAll%sCount", pn), qb, []jen.Code{}, []jen.Code{}, nil)
}

func buildTestDBGetAllSomethingCount(_ *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAll%sCount", dbvsn, pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.Line(),
				jen.ID("expectedQuery").Assign().Litf("SELECT COUNT(%s.id) FROM %s WHERE %s.archived_on IS NULL", tableName, tableName, tableName),
				jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).Dot("AddRow").Call(jen.ID("expectedCount"))),
				jen.Line(),
				jen.List(jen.ID("actualCount"), jen.Err()).Assign().ID(dbfl).Dotf("GetAll%sCount", pn).Call(constants.CtxVar()),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expectedCount"), jen.ID("actualCount"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBGetListOfSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	tableName := typ.Name.PluralRouteName()
	cols := buildPrefixedStringColumns(typ)

	whereValues := typ.BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(proj)

	qb := queryBuilderForDatabase(dbvendor).Select(append(cols, fmt.Sprintf(countQuery, tableName))...).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)
	qb = applyFleshedOutQueryFilterWithCode(qb, tableName, whereValues)

	callArgs := typ.BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(proj)
	pql := typ.BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Get%s", pn), qb, nil, callArgs, pql)
}

func buildTestDBGetListOfSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := string(dbvendor.RouteName()[0])
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()
	cols := buildPrefixedStringColumns(typ)

	equals := squirrel.Eq{fmt.Sprintf("%s.archived_on", tableName): nil}
	if typ.BelongsToUser && typ.RestrictedToUser {
		equals[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
	}

	whereValues := typ.BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(proj)

	qb := queryBuilderForDatabase(dbvendor).Select(append(cols, fmt.Sprintf(countQuery, tableName))...).
		From(tableName)
	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	expectedQuery, args, _ := qb.
		Where(whereValues).
		GroupBy(fmt.Sprintf("%s.id", tableName)).
		Limit(20).
		ToSql()

	withArgs := convertArgsToCode(args)
	actualCallArgs := typ.BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(proj)

	buildFirstSubtest := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj, false)

		expectQueryMock := jen.ID("mockDB").Dot("ExpectQuery").
			Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery")))

		if len(withArgs) > 0 {
			expectQueryMock = expectQueryMock.Dotln("WithArgs").Callln(withArgs...)
		}
		expectQueryMock = expectQueryMock.Dotln("WillReturnRows").Callln(
			jen.IDf("buildMockRowsFrom%s", sn).Callln(
				jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Zero()),
				jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.One()),
				jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Lit(2)),
			),
		)

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.CreateDefaultQueryFilter(proj),
			jen.Line(),
			utils.BuildFakeVar(proj, fmt.Sprintf("%sList", sn)),
			jen.Line(),
			expectQueryMock,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			jen.Line(),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.IDf("example%sList", sn), jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtest := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj, false)
		var mockDBCall *jen.Statement

		if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		} else {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery")))
			if len(withArgs) > 0 {
				mockDBCall = mockDBCall.Dotln("WithArgs").Callln(withArgs...)
			}
			mockDBCall = mockDBCall.Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))

		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.CreateDefaultQueryFilter(proj),
			jen.Line(),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildThirdSubtest := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj, false)
		var mockDBCall *jen.Statement

		if (typ.BelongsToUser && typ.RestrictedToUser) || typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery")))

			if len(withArgs) > 0 {
				mockDBCall = mockDBCall.Dotln("WithArgs").Callln(withArgs...)
			}
			mockDBCall = mockDBCall.Dotln("WillReturnError").Call(constants.ObligatoryError())
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(constants.ObligatoryError())
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.CreateDefaultQueryFilter(proj),
			jen.Line(),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildFourthSubtest := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj, false)

		var mockDBCall *jen.Statement
		if (typ.BelongsToUser && typ.RestrictedToUser) || typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery")))

			if len(withArgs) > 0 {
				mockDBCall = mockDBCall.Dotln("WithArgs").Callln(withArgs...)
			}
			mockDBCall = mockDBCall.Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
		}
		if mockDBCall == nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").
				Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Callln(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.IDf("example%s", sn)),
			)
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.CreateDefaultQueryFilter(proj),
			jen.Line(),
			func() jen.Code {
				if typ.BelongsToUser && typ.RestrictedToUser {
					return utils.BuildFakeVar(proj, "User")
				}
				return jen.Null()
			}(),
			utils.BuildFakeVar(proj, sn),
			func() jen.Code {
				if typ.BelongsToUser && typ.RestrictedToUser {
					return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName("User")).Dot("ID")
				}
				return jen.Null()
			}(),
			jen.Line(),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Get%s", dbvsn, pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			func() jen.Code {
				if typ.RestrictedToUserAtSomeLevel(proj) {
					return jen.ID(utils.BuildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call()
				}
				return jen.Null()
			}(),
			jen.ID("expectedListQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest("happy path", buildFirstSubtest()...),
			jen.Line(),
			utils.BuildSubTest("surfaces sql.ErrNoRows", buildSecondSubtest()...),
			jen.Line(),
			utils.BuildSubTest("with error executing read query", buildThirdSubtest()...),
			jen.Line(),
			utils.BuildSubTest(fmt.Sprintf("with error scanning %s", scn), buildFourthSubtest()...),
		),
		jen.Line(),
	}
}

func buildTestDBCreateSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	fieldCols, expectedArgs := buildCreationStringColumnsAndArgs(proj, typ)
	valueArgs := []interface{}{}
	for range expectedArgs {
		valueArgs = append(valueArgs, whateverValue)
	}

	qb := queryBuilderForDatabase(dbvendor).
		Insert(tableName).
		Columns(fieldCols...).
		Values(valueArgs...)

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING id, created_on")
	}
	pql := typ.BuildDBQuerierCreateSomethingQueryBuilderTestPreQueryLines(proj)
	callArgs := typ.BuildArgsToUseForDBQuerierCreationQueryBuildingTest(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Create%s", sn), qb, expectedArgs, callArgs, pql)
}

func buildTestDBCreateSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	const (
		expectedQueryVarName = "expectedCreationQuery"
	)

	fieldCols, expectedArgs := buildCreationStringColumnsAndArgs(proj, typ)
	callArgs := typ.BuildDBQuerierCreationMethodArgsToUseFromMethodTest(proj)

	valueArgs := []interface{}{}
	for range expectedArgs {
		valueArgs = append(valueArgs, whateverValue)
	}

	qb := queryBuilderForDatabase(dbvendor).
		Insert(tableName).
		Columns(fieldCols...).
		Values(valueArgs...)

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING id, created_on")
	}
	expectedQuery, _, _ := qb.ToSql()

	buildFirstSubtest := func(proj *models.Project, typ models.DataType) []jen.Code {
		expectedValues := []jen.Code{}
		if typ.BelongsToUser && typ.RestrictedToUser {
			expectedValues = append(expectedValues, jen.ID(constants.UserOwnershipFieldName).MapAssign().ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out := []jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		}

		out = append(out, typ.BuildDependentObjectsForDBQueriersCreationMethodTest(proj)...)
		out = append(out, jen.Line())

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID(utils.BuildFakeVarName("Rows")).Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
					jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
					jen.ID(utils.BuildFakeVarName(sn)).Dot("CreatedOn"),
				),
			)
		}

		var nef []jen.Code
		for _, field := range typ.Fields {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(field.Name.Singular()))
		}

		if typ.BelongsToStruct != nil {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		}
		if typ.BelongsToUser {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
		}

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).
					Dot("WillReturnRows").Call(jen.ID(utils.BuildFakeVarName("Rows"))),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			out = append(out,
				jen.Line(),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).Dot("WillReturnResult").Call(
					jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
						jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")),
						jen.Lit(1),
					),
				),
				jen.Line(),
				jen.ID("mtt").Assign().AddressOf().ID("mockTimeTeller").Values(),
				jen.ID("mtt").Dot("On").Call(jen.Lit("Now")).Dot("Return").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("CreatedOn")),
				jen.ID(dbfl).Dot("timeTeller").Equals().ID("mtt"),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Create%s", sn).Call(
				callArgs...,
			),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual"), nil),
			jen.Line(),
			func() jen.Code {
				if isSqlite(dbvendor) || isMariaDB(dbvendor) {
					return utils.AssertExpectationsFor("mtt")
				}
				return jen.Null()
			}(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	buildSecondSubtest := func(proj *models.Project, typ models.DataType) []jen.Code {
		expectedValues := []jen.Code{}
		if typ.BelongsToUser && typ.RestrictedToUser {
			expectedValues = append(expectedValues, jen.ID(constants.UserOwnershipFieldName).MapAssign().ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out := []jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		}

		out = append(out, typ.BuildDependentObjectsForDBQueriersCreationMethodTest(proj)...)
		out = append(out, jen.Line())

		var nef []jen.Code
		for _, field := range typ.Fields {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(field.Name.Singular()))
		}

		if typ.BelongsToStruct != nil {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		}
		if typ.BelongsToUser {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
		}

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).
					Dot("WillReturnError").Call(constants.ObligatoryError()),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			out = append(out,
				jen.Line(),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).Dot("WillReturnError").Call(
					constants.ObligatoryError(),
				),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Create%s", sn).Call(
				callArgs...,
			),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Create%s", dbvsn, sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID(expectedQueryVarName).Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest("happy path", buildFirstSubtest(proj, typ)...),
			jen.Line(),
			utils.BuildSubTest("with error writing to database", buildSecondSubtest(proj, typ)...),
		),
		jen.Line(),
	}
}

func buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(dbvendor wordsmith.SuperPalabra, typ models.DataType) (squirrel.UpdateBuilder, []jen.Code) {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	updateCols := buildUpdateQueryParts(dbvendor, typ)
	expectedArgs := []jen.Code{}

	valueArgs := []interface{}{}
	for range updateCols {
		valueArgs = append(valueArgs, whateverValue)
	}

	eq := squirrel.Eq{"id": whateverValue}
	qb := queryBuilderForDatabase(dbvendor).Update(tableName)
	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			qb = qb.Set(field.Name.RouteName(), jen.ID("input").Dot(field.Name.Singular()))
			expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToStruct != nil {
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToUser {
		eq["belongs_to_user"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}
	expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	qb = qb.Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor)))
	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING updated_on")
	}

	qb = qb.Where(eq)

	return qb, expectedArgs
}

func buildTestBuildUpdateSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	qb, expectedArgs := buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)
	callArgs := typ.BuildArgsForDBQuerierTestOfUpdateQueryBuilder(proj)
	pql := typ.BuildDBQuerierUpdateSomethingQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Update%s", sn), qb, expectedArgs, callArgs, pql)
}

func buildTestDBUpdateSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	qb, expectQueryArgs := buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

	expectedQuery, _, err := qb.ToSql()
	if err != nil {
		log.Fatalf("error running buildTestDBUpdateSomethingFuncDecl: %v", err)
	}

	callArgs := typ.BuildArgsForDBQuerierTestOfUpdateMethod(proj)

	buildFirstSubTest := func(typ models.DataType) []jen.Code {
		var (
			expectFuncName,
			returnFuncName string
			exRows jen.Code
		)

		dbrn := dbvendor.RouteName()
		sn := typ.Name.Singular()
		dbfl := string(dbrn[0])

		if isPostgres(dbvendor) {
			expectFuncName = "ExpectQuery"
			returnFuncName = "WillReturnRows"

			exRows = jen.ID(utils.BuildFakeVarName("Rows")).Assign().
				Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").
				Call(jen.Index().String().Values(jen.Lit("updated_on"))).Dot("AddRow").
				Call(jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			expectFuncName = "ExpectExec"
			returnFuncName = "WillReturnResult"
			exRows = jen.ID(utils.BuildFakeVarName("Rows")).Assign().
				Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").
				Call(jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")), jen.Lit(1))
		}

		lines := append([]jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		},
			typ.BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(proj)...)

		lines = append(lines,
			jen.Line(),
			exRows,
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot(returnFuncName).Call(jen.ID(utils.BuildFakeVarName("Rows"))),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(callArgs...),
			utils.AssertNoError(jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtest := func(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
		var (
			expectFuncName string
		)

		dbrn := dbvendor.RouteName()
		sn := typ.Name.Singular()
		dbfl := string(dbrn[0])

		if isPostgres(dbvendor) {
			expectFuncName = "ExpectQuery"
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			expectFuncName = "ExpectExec"
		}

		lines := append([]jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		},
			typ.BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(proj)...)

		lines = append(lines,
			jen.Line(),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot("WillReturnError").Call(utils.FakeError()),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(callArgs...),
			utils.AssertError(jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}
	return []jen.Code{
		jen.Func().IDf("Test%s_Update%s", dbvendor.Singular(), typ.Name.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTest("happy path", buildFirstSubTest(typ)...),
			jen.Line(),
			utils.BuildSubTest("with error writing to database", buildSecondSubtest(proj, dbvendor, typ)...),
		),
		jen.Line(),
	}
}

func buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) (qb squirrel.UpdateBuilder, expectedArgs []jen.Code, callArgs []jen.Code) {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	updateCols := buildUpdateQueryParts(dbvendor, typ)
	valueArgs := []interface{}{}
	for range updateCols {
		valueArgs = append(valueArgs, whateverValue)
	}

	eq := squirrel.Eq{
		"id":          whateverValue,
		"archived_on": nil,
	}
	if typ.BelongsToStruct != nil {
		btssn := typ.BelongsToStruct.Singular()
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(btssn)).Dot("ID"))
		callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(btssn)).Dot("ID"))
	}
	if typ.BelongsToUser {
		eq["belongs_to_user"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName("User")).Dot("ID"))
	}
	callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	qb = queryBuilderForDatabase(dbvendor).
		Update(tableName).
		Set("updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(eq)

	if isPostgres(dbvendor) {
		qb = qb.Suffix("RETURNING archived_on")
	}

	return
}

func buildTestDBArchiveSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	qb, expectedArgs, _ := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(proj, dbvendor, typ)
	pql := typ.BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(proj)
	callArgs := typ.BuildArgsForDBQuerierTestOfArchiveQueryBuilder(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Archive%s", sn), qb, expectedArgs, callArgs, pql)
}

func buildTestDBArchiveSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()

	qb, dbQueryExpectationArgs, _ := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(proj, dbvendor, typ)
	actualCallArgs := typ.BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest(proj)

	dbQuery, _, _ := qb.ToSql()

	buildSubtestOne := func() []jen.Code {
		block := append([]jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		},
			typ.BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(proj)...)

		block = append(block,
			jen.Line(),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).Dot("WillReturnResult").Call(
				jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
					jen.One(),
					jen.One(),
				),
			),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Archive%s", sn).Call(actualCallArgs...),
			utils.AssertNoError(jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return block
	}

	buildSubtestTwo := func() []jen.Code {
		block := append([]jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		},
			typ.BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(proj)...)

		block = append(block,
			jen.Line(),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).Dot("WillReturnResult").Call(
				jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
					jen.Zero(),
					jen.Zero(),
				),
			),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Archive%s", sn).Call(actualCallArgs...),
			utils.AssertError(jen.Err(), nil),
			utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return block
	}

	buildSubtestThree := func() []jen.Code {
		block := append([]jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		},
			typ.BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(proj)...)

		block = append(block,
			jen.Line(),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).
				Dot("WillReturnError").Call(constants.ObligatoryError()),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Archive%s", sn).Call(actualCallArgs...),
			utils.AssertError(jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(
				jen.ID("mockDB").Dot("ExpectationsWereMet").Call(),
				jen.Lit("not all database expectations were met"),
			),
		)

		return block
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Archive%s", dbvsn, sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(dbQuery),
			jen.Line(),
			utils.BuildSubTest("happy path", buildSubtestOne()...),
			jen.Line(),
			utils.BuildSubTest("returns sql.ErrNoRows with no rows affected", buildSubtestTwo()...),
			jen.Line(),
			utils.BuildSubTest("with error writing to database", buildSubtestThree()...),
		),
		jen.Line(),
	}
}
