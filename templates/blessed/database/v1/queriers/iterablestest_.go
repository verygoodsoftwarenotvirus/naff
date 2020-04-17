package queriers

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"log"
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
		jen.ID(varName).Dot("UpdatedOn"),
		jen.ID(varName).Dot("ArchivedOn"),
	)

	if typ.BelongsToUser && typ.RestrictedToUser {
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
		jen.ID(varName).Dot("UpdatedOn"),
	)

	if typ.BelongsToUser && typ.RestrictedToUser {
		fields = append(fields, jen.ID(varName).Dot(constants.UserOwnershipFieldName))
	}
	if typ.BelongsToStruct != nil {
		fields = append(fields, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
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
	if typ.BelongsToUser && typ.RestrictedToUser {
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

func buildCreationSringColumnsAndArgs(proj *models.Project, typ models.DataType) (cols []string, args []jen.Code) {
	cols, args = []string{}, []jen.Code{}

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			cols = append(cols, field.Name.RouteName())
			args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToUser && typ.RestrictedToUser {
		cols = append(cols, "belongs_to_user")
		args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot(constants.UserOwnershipFieldName))
	}

	if typ.BelongsToStruct != nil {
		cols = append(cols, fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))
		args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
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

	ret := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, ret)

	n := typ.Name
	sn := n.Singular()
	puvn := n.PluralUnexportedVarName()

	gFields := buildGeneralFields("x", typ)

	ret.Add(
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

	ret.Add(
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

	ret.Add(buildTestDBBuildSomethingExistsQuery(proj, dbvendor, typ)...)
	ret.Add(buildTestDBSomethingExists(proj, dbvendor, typ)...)
	ret.Add(buildTestDBBuildGetSomethingQuery(proj, dbvendor, typ)...)
	ret.Add(buildTestDBGetSomething(proj, dbvendor, typ)...)
	ret.Add(buildTestDBBuildGetAllSomethingCountQuery(proj, dbvendor, typ)...)
	ret.Add(buildTestDBGetAllSomethingCount(proj, dbvendor, typ)...)
	ret.Add(buildTestDBGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	ret.Add(buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)...)
	ret.Add(buildTestDBCreateSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	ret.Add(buildTestDBCreateSomethingFuncDecl(proj, dbvendor, typ)...)
	ret.Add(buildTestBuildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	ret.Add(buildTestDBUpdateSomethingFuncDecl(proj, dbvendor, typ)...)
	ret.Add(buildTestDBArchiveSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	ret.Add(buildTestDBArchiveSomethingFuncDecl(proj, dbvendor, typ)...)

	return ret
}

func buildTestDBBuildSomethingExistsQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
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

	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName).
		Where(eqArgs).
		Suffix(existenceSuffix)

	expectationArgs := []jen.Code{
		func() jen.Code {
			if typ.BelongsToUser && typ.RestrictedToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName)
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
	}
	callArgs := typ.BuildDBBuildSomethingExistsQueryTestCallArgs(proj)
	pql := typ.BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("%sExists", sn), qb, expectationArgs, callArgs, pql)
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

	qb := queryBuilderForDatabase(dbvendor)
	query, _, _ := qb.Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName).
		Where(eqArgs).
		Suffix(existenceSuffix).
		ToSql()
	actualCallArgs := typ.BuildArgsForDBQuerierExistenceMethodTest(proj)

	buildFirstSubtestBlock := func(typ models.DataType) []jen.Code {
		lines := typ.BuildDependentObjectsForDBQueriersExistenceMethodTest(proj)

		var mockDBCall jen.Code

		if typ.BelongsToUser && typ.RestrictedToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))
		}
		if typ.BelongsToStruct != nil {
			btssn := typ.BelongsToStruct.Singular()
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(btssn)).Dot("ID"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))
		}

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

	lines := []jen.Code{
		jen.Func().IDf("Test%s_%sExists", dbvsn, sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(query),
			jen.Line(),
			utils.BuildSubTest("happy path", buildFirstSubtestBlock(typ)...),
		),
		jen.Line(),
	}

	return lines
}

func buildTestDBBuildGetSomethingQuery(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	eqArgs := squirrel.Eq{fmt.Sprintf("%s.id", tableName): whateverValue}
	if typ.BelongsToUser && typ.RestrictedToUser {
		eqArgs[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = whateverValue
	}

	qb := queryBuilderForDatabase(dbvendor).Select(buildPrefixedStringColumnsAsString(typ)).
		From(tableName).
		Where(eqArgs)

	expectationArgs := []jen.Code{
		func() jen.Code {
			if typ.BelongsToUser && typ.RestrictedToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName)
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
	}
	callArgs := typ.BuildDBQuerierRetrievalQueryTestCallArgs(proj)
	pql := typ.BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Get%s", sn), qb, expectationArgs, callArgs, pql)
}

func buildTestDBGetSomething(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	actualCallArgs := typ.BuildArgsForDBQuerierRetrievalMethodTest(proj)

	buildFirstSubtestBlock := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(proj)

		var mockDBCall jen.Code

		if typ.BelongsToUser && typ.RestrictedToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))))
		}

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

		var mockDBCall jen.Code
		if typ.BelongsToUser && typ.RestrictedToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		}

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

	eqArgs := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): whateverValue,
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		eqArgs[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, rn)] = whateverValue
	}

	qb := queryBuilderForDatabase(dbvendor)
	query, _, _ := qb.Select(buildPrefixedStringColumnsAsString(typ)).
		From(tableName).
		Where(eqArgs).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_Get%s", dbvsn, sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			func() jen.Code {
				if typ.BelongsToUser {
					return utils.BuildFakeVar(proj, "User")
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
	expectedArgs := []jen.Code{}
	equals := squirrel.Eq{fmt.Sprintf("%s.archived_on", tableName): nil}

	if typ.BelongsToUser && typ.RestrictedToUser {
		equals[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName("User")).Dot("ID"))
	}

	qb := queryBuilderForDatabase(dbvendor).Select(append(cols, fmt.Sprintf(countQuery, tableName))...).
		From(tableName).
		Where(equals)
	qb = applyFleshedOutQueryFilter(qb, tableName)
	expectedArgs = appendFleshedOutQueryFilterArgs(expectedArgs)

	callArgs := typ.BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(proj)
	pql := typ.BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Get%s", pn), qb, expectedArgs, callArgs, pql)
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
	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).Select(append(cols, fmt.Sprintf(countQuery, tableName))...).
		From(tableName).
		Where(equals).
		GroupBy(fmt.Sprintf("%s.id", tableName)).
		Limit(20).
		ToSql()

	expectedQueryArgs := typ.BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(proj)
	actualCallArgs := typ.BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(proj)

	buildFirstSubtest := func() []jen.Code {
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj)

		var expectQueryMock jen.Code

		if typ.BelongsToUser && typ.RestrictedToUser {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(
				expectedQueryArgs...,
			).
				Dotln("WillReturnRows").Callln(
				jen.IDf("buildMockRowsFrom%s", sn).Callln(
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Zero()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.One()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Lit(2)),
				),
			)
		}
		if typ.BelongsToStruct != nil {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(
				expectedQueryArgs...,
			).
				Dotln("WillReturnRows").Callln(
				jen.IDf("buildMockRowsFrom%s", sn).Callln(
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Zero()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.One()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Lit(2)),
				),
			)
		} else if typ.BelongsToNobody {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Callln(
				jen.IDf("buildMockRowsFrom%s", sn).Callln(
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Zero()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.One()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Lit(2)),
				),
			)
		}

		lines = append(lines,
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
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj)
		var mockDBCall jen.Code

		if typ.BelongsToUser && typ.RestrictedToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(expectedQueryArgs...).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(expectedQueryArgs...).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))

		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		}

		lines = append(lines,
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
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj)
		var mockDBCall jen.Code

		if typ.BelongsToUser && typ.RestrictedToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(expectedQueryArgs...).
				Dotln("WillReturnError").Call(constants.ObligatoryError())
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(expectedQueryArgs...).
				Dotln("WillReturnError").Call(constants.ObligatoryError())
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(constants.ObligatoryError())
		}

		lines = append(lines,
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
		lines := typ.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(proj)
		var mockDBCall jen.Code

		if typ.BelongsToUser && typ.RestrictedToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(expectedQueryArgs...).
				Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(expectedQueryArgs...).
				Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
		}

		lines = append(lines,
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.CreateDefaultQueryFilter(proj),
			utils.BuildFakeVar(proj, sn),
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
				if typ.BelongsToUser && typ.RestrictedToUser {
					return utils.BuildFakeVar(proj, "User")
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

	fieldCols, expectedArgs := buildCreationSringColumnsAndArgs(proj, typ)
	callArgs := typ.BuildArgsToUseForDBQuerierCreationQueryBuildingTest(proj)
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

	fieldCols, expectedArgs := buildCreationSringColumnsAndArgs(proj, typ)
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

		if typ.BelongsToUser && typ.RestrictedToUser {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
		}
		if typ.BelongsToStruct != nil {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
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

		if typ.BelongsToUser && typ.RestrictedToUser {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
		}
		if typ.BelongsToStruct != nil {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
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

	if typ.BelongsToUser && typ.RestrictedToUser {
		eq["belongs_to_user"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}
	if typ.BelongsToStruct != nil {
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
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

		lines := []jen.Code{}
		expectedValues := []jen.Code{}
		if typ.BelongsToUser && typ.RestrictedToUser {
			expectedValues = append(expectedValues, jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		lines = append(lines,
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
			utils.BuildFakeVar(proj, sn),
			jen.Line(),
			exRows,
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot(returnFuncName).Call(jen.ID(utils.BuildFakeVarName("Rows"))),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName(sn))),
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

		lines := []jen.Code{}
		expectedValues := []jen.Code{}
		if typ.BelongsToUser && typ.RestrictedToUser {
			expectedValues = append(expectedValues, jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		lines = append(lines,
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.BuildFakeVar(proj, sn),
			jen.Line(),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot("WillReturnError").Call(utils.FakeError()),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(constants.CtxVar(), jen.ID(utils.BuildFakeVarName(sn))),
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

func buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor wordsmith.SuperPalabra, typ models.DataType) (qb squirrel.UpdateBuilder, expectedArgs []jen.Code, callArgs []jen.Code) {
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
	if typ.BelongsToUser && typ.RestrictedToUser {
		eq["belongs_to_user"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}
	callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))
	if typ.BelongsToStruct != nil {
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToUser && typ.RestrictedToUser {
		callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
	}
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

	qb, expectedArgs, callArgs := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)
	pql := typ.BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(proj)

	return buildQueryTest(proj, dbvendor, fmt.Sprintf("Archive%s", sn), qb, expectedArgs, callArgs, pql)
}

func buildTestDBArchiveSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()

	qb, dbQueryExpectationArgs, _ := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)
	actualCallArgs := typ.BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest(proj)

	dbQuery, _, _ := qb.ToSql()

	buildSubtestOne := func() []jen.Code {
		expectedValues := []jen.Code{}

		if typ.BelongsToUser && typ.RestrictedToUser {
			expectedValues = append(expectedValues, jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}

		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		block := []jen.Code{}

		block = append(block,
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.BuildFakeVar(proj, sn),
			jen.Line(),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).Dot("WillReturnResult").Call(
				jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
					jen.Lit(1),
					jen.Lit(1),
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
		exampleValues := []jen.Code{}

		var dbQueryExpectationArgs []jen.Code
		block := []jen.Code{}

		if typ.BelongsToUser && typ.RestrictedToUser {
			exampleValues = append(exampleValues, jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("exampleUser").Dot("ID"))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.UserOwnershipFieldName))
		}
		if typ.BelongsToStruct != nil {
			btss := typ.BelongsToStruct.Singular()
			exampleValues = append(exampleValues, jen.IDf("BelongsTo%s", btss).MapAssign().ID(utils.BuildFakeVarName(btss)).Dot("ID"))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", btss))
		}

		dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))
		exampleValues = append(exampleValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		block = append(block,
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.BuildFakeVar(proj, sn),
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
			utils.BuildSubTest("with error writing to database", buildSubtestTwo()...),
		),
		jen.Line(),
	}
}
