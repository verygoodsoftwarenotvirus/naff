package queriers

import (
	"fmt"
	"github.com/Masterminds/squirrel"
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

	if typ.BelongsToUser {
		fields = append(fields, jen.ID(varName).Dot("BelongsToUser"))
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

func buildCreationSringColumnsAndArgs(typ models.DataType) (cols []string, args []jen.Code) {
	cols, args = []string{}, []jen.Code{}

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			cols = append(cols, field.Name.RouteName())
			args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToUser {
		cols = append(cols, "belongs_to_user")
		args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot("BelongsToUser"))
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

func buildCreationEqualityExpectations(varName string, typ models.DataType) []jen.Code {
	var out []jen.Code

	for i, field := range typ.Fields {
		if field.Pointer {
			out = append(out,
				utils.AssertEqual(jen.ID(varName).Dot(field.Name.Singular()), jen.ID("args").Index(jen.Lit(i)).Assert(jen.PointerTo().ID(field.Type)), nil),
			)
		} else {
			out = append(out,
				utils.AssertEqual(jen.ID(varName).Dot(field.Name.Singular()), jen.ID("args").Index(jen.Lit(i)).Assert(jen.ID(field.Type)), nil),
			)
		}
	}

	if typ.BelongsToUser {
		out = append(out,
			utils.AssertEqual(jen.ID("expected").Dot("BelongsToUser"), jen.ID("args").Index(jen.Lit(len(out))).Assert(jen.Uint64()), nil),
		)
	}
	if typ.BelongsToStruct != nil {
		out = append(out,
			utils.AssertEqual(jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()), jen.ID("args").Index(jen.Lit(len(out))).Assert(jen.Uint64()), nil),
		)
	}

	return out
}

func buildFieldMaps(varName string, typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		xn := field.Name.Singular()
		out = append(out, jen.ID(xn).MapAssign().ID(varName).Dot(xn))
	}

	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").MapAssign().ID(varName).Dot("BelongsToUser"))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
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
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()), jen.ID(varName).Dot("ID"))
	}

	return out
}

func iterablesTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

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
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For().List(jen.Underscore(), jen.ID("x")).Assign().Range().ID(puvn).Block(
				jen.ID("rowValues").Assign().Index().Qual("database/sql/driver", "Value").Valuesln(gFields...),
				jen.Line(),
				jen.If(jen.ID("includeCount")).Block(
					utils.AppendItemsToList(jen.ID("rowValues"), jen.Len(jen.ID("items"))),
				),
				jen.Line(),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Spread()),
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
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.IDf("%sTableColumns", puvn)).Dot("AddRow").Callln(badFields...),
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

	// if typ.BelongsToUser || typ.BelongsToStruct != nil {
	// 	ret.Add(buildTestDBGetAllSomethingForSomethingElseFuncDecl(proj, dbvendor, typ)...)
	// }

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
		fmt.Sprintf("%s.id", tableName): "fart",
	}
	if typ.BelongsToUser {
		eqArgs[fmt.Sprintf("%s.belongs_to_user", tableName)] = "fart"
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = "fart"
	}

	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName).
		Where(eqArgs).
		Suffix(existenceSuffix)

	expectationArgs := []jen.Code{
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
	}
	callArgs := typ.BuildGetSomethingArgsWithExampleVariables(proj)

	return buildQueryTest(proj, dbvendor, typ, fmt.Sprintf("%sExists", sn), qb, expectationArgs, callArgs[1:], true, false, false, false, true, false, nil)
}

func buildTestDBSomethingExists(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	eqArgs := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): "fart",
	}
	if typ.BelongsToUser {
		eqArgs[fmt.Sprintf("%s.belongs_to_user", tableName)] = "fart"
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = "fart"
	}

	qb := queryBuilderForDatabase(dbvendor)
	query, _, _ := qb.Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName).
		Where(eqArgs).
		Suffix(existenceSuffix).
		ToSql()

	buildFirstSubtestBlock := func(typ models.DataType) []jen.Code {
		lines := []jen.Code{utils.BuildFakeVar(proj, sn)}

		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			utils.CtxVar(), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			btssn := typ.BelongsToStruct.Singular()

			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(btssn)).Dot("ID"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(btssn)).Dot("ID"))
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
	if typ.BelongsToUser {
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
			if typ.BelongsToUser {
				return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser")
			}
			return jen.Null()
		}(),
		jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
	}
	callArgs := typ.BuildGetSomethingArgsWithExampleVariables(proj)

	return buildQueryTest(proj, dbvendor, typ, fmt.Sprintf("Get%s", sn), qb, expectationArgs, callArgs[1:], true, false, false, false, true, false, nil)
}

func buildTestDBGetSomething(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	buildFirstSubtestBlock := func() []jen.Code {
		lines := []jen.Code{
			utils.BuildFakeVar(proj, sn),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToUser {
					return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID")
				}
				return jen.Null()
			}(),
		}

		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
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
		lines := []jen.Code{
			utils.BuildFakeVar(proj, sn),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID")
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToUser {
					return jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID")
				}
				return jen.Null()
			}(),
		}

		actualCallArgs := []jen.Code{
			utils.CtxVar(), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		}
		var mockDBCall jen.Code
		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"), jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
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
	if typ.BelongsToUser {
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
			utils.BuildFakeVar(proj, "User"),
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

	return buildQueryTest(proj, dbvendor, typ, fmt.Sprintf("GetAll%sCount", pn), qb, []jen.Code{}, []jen.Code{}, false, false, false, false, false, false, nil)
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
				jen.List(jen.ID("actualCount"), jen.Err()).Assign().ID(dbfl).Dotf("GetAll%sCount", pn).Call(utils.CtxVar()),
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
		jen.ID(utils.FilterVarName).Dot("CreatedAfter"),
		jen.ID(utils.FilterVarName).Dot("CreatedBefore"),
		jen.ID(utils.FilterVarName).Dot("UpdatedAfter"),
		jen.ID(utils.FilterVarName).Dot("UpdatedBefore"),
	)

	return args
}

func buildTestDBGetListOfSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	tableName := typ.Name.PluralRouteName()
	cols := buildPrefixedStringColumns(typ)
	expectedArgs := []jen.Code{}
	equals := squirrel.Eq{fmt.Sprintf("%s.archived_on", tableName): nil}

	if typ.BelongsToUser {
		equals[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID("exampleUser").Dot("ID"))
	}

	qb := queryBuilderForDatabase(dbvendor).Select(append(cols, fmt.Sprintf(countQuery, tableName))...).
		From(tableName).
		Where(equals)
	qb = applyFleshedOutQueryFilter(qb, tableName)
	expectedArgs = appendFleshedOutQueryFilterArgs(expectedArgs)

	return buildQueryTest(proj, dbvendor, typ, fmt.Sprintf("Get%s", pn), qb, expectedArgs, []jen.Code{}, true, true, true, true, false, false, nil)
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
	if typ.BelongsToUser {
		equals[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
	}
	expectedQuery, _, _ := queryBuilderForDatabase(dbvendor).Select(append(cols, fmt.Sprintf(countQuery, tableName))...).
		From(tableName).
		Where(equals).
		GroupBy(fmt.Sprintf("%s.id", tableName)).
		Limit(20).
		ToSql()

	buildFirstSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var expectQueryMock jen.Code

		actualCallArgs := []jen.Code{utils.CtxVar()}
		if typ.BelongsToUser {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("exampleUser").Dot("ID")).
				Dotln("WillReturnRows").Callln(
				jen.IDf("buildMockRowsFrom%s", sn).Callln(
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Zero()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.One()),
					jen.AddressOf().ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))).Dot(pn).Index(jen.Lit(2)),
				),
			)
			actualCallArgs = append(actualCallArgs, jen.ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Callln(
				jen.IDf("buildMockRowsFrom%s", sn).Callln(
					jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))),
				),
			)
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Callln(
				jen.IDf("buildMockRowsFrom%s", sn).Callln(
					jen.ID(utils.BuildFakeVarName(fmt.Sprintf("%sList", sn))),
				),
			)
		}
		actualCallArgs = append(actualCallArgs, jen.ID(utils.FilterVarName))

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
		lines := []jen.Code{}
		var mockDBCall jen.Code

		actualCallArgs := []jen.Code{utils.CtxVar()}
		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("exampleUser").Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
		}
		actualCallArgs = append(actualCallArgs, jen.ID(utils.FilterVarName))

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
		lines := []jen.Code{}
		var mockDBCall jen.Code

		actualCallArgs := []jen.Code{utils.CtxVar()}
		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("exampleUser").Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
			actualCallArgs = append(actualCallArgs, jen.ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
		}
		actualCallArgs = append(actualCallArgs, jen.ID(utils.FilterVarName))

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
		lines := []jen.Code{}
		var mockDBCall jen.Code

		actualCallArgs := []jen.Code{utils.CtxVar()}
		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser")).
				Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(
				jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn))),
			)
		}
		actualCallArgs = append(actualCallArgs, jen.ID(utils.FilterVarName))

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
			utils.BuildFakeVar(proj, "User"),
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

	fieldCols, expectedArgs := buildCreationSringColumnsAndArgs(typ)
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

	return buildQueryTest(proj, dbvendor, typ, fmt.Sprintf("Create%s", sn), qb, expectedArgs, []jen.Code{jen.ID(utils.BuildFakeVarName(sn))}, true, false, false, false, true, false, nil)
}

func buildTestDBCreateSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	const (
		expectedQueryVarName = "expectedCreationQuery"
	)

	fieldCols, expectedArgs := buildCreationSringColumnsAndArgs(typ)
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
		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out := []jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
			utils.BuildFakeVar(proj, sn),
			utils.BuildFakeVarWithCustomName(proj, "exampleInput", fmt.Sprintf("%sCreationInputFrom%s", sn, sn), jen.ID(utils.BuildFakeVarName(sn))),
			jen.Line(),
		}

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
					jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
					jen.ID(utils.BuildFakeVarName(sn)).Dot("CreatedOn"),
				),
			)
		}

		var nef []jen.Code
		for _, field := range typ.Fields {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(field.Name.Singular()))
		}

		if typ.BelongsToUser {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		}

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).
					Dot("WillReturnRows").Call(jen.ID("exampleRows")),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			out = append(out,
				jen.Line(),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).Dot("WillReturnResult").Call(
					jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
						jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")),
						jen.Lit(123),
					),
				),
				jen.Line(),
				jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM %s WHERE id = %s", tableName, getIncIndex(dbvendor, 0)),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(
					jen.Index().String().Values(jen.Lit("created_on")),
				).Dot("AddRow").Call(
					jen.ID(utils.BuildFakeVarName(sn)).Dot("CreatedOn")),
				),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Create%s", sn).Call(utils.CtxVar(), jen.ID("exampleInput")),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	buildSecondSubtest := func(proj *models.Project, typ models.DataType) []jen.Code {
		expectedValues := []jen.Code{}
		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out := []jen.Code{
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
			utils.BuildFakeVar(proj, sn),
			utils.BuildFakeVarWithCustomName(proj, "exampleInput", fmt.Sprintf("%sCreationInputFrom%s", sn, sn), jen.ID(utils.BuildFakeVarName(sn))),
			jen.Line(),
		}

		var nef []jen.Code
		for _, field := range typ.Fields {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot(field.Name.Singular()))
		}

		if typ.BelongsToUser {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			nef = append(nef, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		}

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).
					Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			out = append(out,
				jen.Line(),
				jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID(expectedQueryVarName))).
					Dotln("WithArgs").Callln(nef...).Dot("WillReturnResult").Call(
					jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
						jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")),
						jen.Lit(123),
					),
				),
				jen.Line(),
				jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM %s WHERE id = %s", tableName, getIncIndex(dbvendor, 0)),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
					Dotln("WithArgs").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")).
					Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Create%s", sn).Call(utils.CtxVar(), jen.ID("exampleInput")),
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

	if typ.BelongsToUser {
		eq["belongs_to_user"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
	}
	if typ.BelongsToStruct != nil {
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", sn))
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

	return buildQueryTest(proj, dbvendor, typ, fmt.Sprintf("Update%s", sn), qb, expectedArgs, []jen.Code{jen.ID(utils.BuildFakeVarName(sn))}, true, false, false, false, true, false, nil)
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

			exRows = jen.ID("exampleRows").Assign().
				Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").
				Call(jen.Index().String().Values(jen.Lit("updated_on"))).
				Dot("AddRow").
				Call(jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			expectFuncName = "ExpectExec"
			returnFuncName = "WillReturnResult"
			exRows = jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
				jen.ID("int64").Call(jen.ID(utils.BuildFakeVarName(sn)).Dot("ID")),
				jen.Lit(123),
			)
		}

		lines := []jen.Code{}
		expectedValues := []jen.Code{}
		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		lines = append(lines,
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.BuildFakeVar(proj, sn),
			jen.Line(),
			exRows,
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot(returnFuncName).Call(jen.ID("exampleRows")),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName(sn))),
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
		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("exampleUser").Dot("ID"))
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
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(utils.CtxVar(), jen.ID(utils.BuildFakeVarName(sn))),
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
	if typ.BelongsToUser {
		eq["belongs_to_user"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
	}
	callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))
	if typ.BelongsToStruct != nil {
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToUser {
		callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
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

	return buildQueryTest(proj, dbvendor, typ, fmt.Sprintf("Archive%s", sn), qb, expectedArgs, callArgs, true, false, false, false, true, false, nil)
}

func buildTestDBArchiveSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()

	qb, dbQueryExpectationArgs, _ := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

	dbQuery, _, _ := qb.ToSql()

	buildSubtestOne := func() []jen.Code {
		expectedValues := []jen.Code{}
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		}

		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("exampleUser").Dot("ID"))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(typ.BelongsToStruct.Singular())).Dot("ID"))
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
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"),
		}

		if typ.BelongsToUser {
			exampleValues = append(exampleValues, jen.ID("BelongsToUser").MapAssign().ID("exampleUser").Dot("ID"))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			btss := typ.BelongsToStruct.Singular()
			exampleValues = append(exampleValues, jen.IDf("BelongsTo%s", btss).MapAssign().ID(utils.BuildFakeVarName(btss)).Dot("ID"))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", btss))
			actualCallArgs = append(actualCallArgs, jen.ID(utils.BuildFakeVarName(btss)).Dot("ID"))
		}

		dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))
		exampleValues = append(exampleValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		block = append(block,
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			utils.BuildFakeVar(proj, sn),
			jen.Line(),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).
				Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
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

func buildTestDBGetAllSomethingForSomethingElseFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := dbvendor.LowercaseAbbreviation()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()
	cols := buildPrefixedStringColumnsAsString(typ)

	var (
		baseFuncName        string
		testFuncName        string
		expectedQuery       string
		expectedSomethingID string
	)

	if typ.BelongsToUser {
		expectedSomethingID = "exampleUserID"
		baseFuncName = fmt.Sprintf("GetAll%sForUser", pn)
		testFuncName = fmt.Sprintf("Test%s_%s", dbvsn, baseFuncName)
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_user = %s", cols, tableName, getIncIndex(dbvendor, 0))
	}
	if typ.BelongsToStruct != nil {
		expectedSomethingID = fmt.Sprintf("expected%sID", typ.BelongsToStruct.Singular())
		baseFuncName = fmt.Sprintf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular())
		testFuncName = fmt.Sprintf("Test%s_%s", dbvsn, baseFuncName)
		expectedQuery = fmt.Sprintf("SELECT %s FROM %s WHERE archived_on IS NULL AND belongs_to_%s = %s", cols, tableName, typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0))
	}
	// we don't need to consider the case where this object belongs to nothing

	return []jen.Code{
		jen.Func().ID(testFuncName).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID(expectedSomethingID).Assign().Add(utils.FakeUint64Func()),
				jen.IDf("expected%s", sn).Assign().AddressOf().Qual(proj.ModelsV1Package(), sn).Valuesln(
					jen.ID("ID").MapAssign().Lit(321),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.IDf("expected%s", sn))),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName(sn)).Assign().Index().Qual(proj.ModelsV1Package(), sn).Values(jen.PointerTo().IDf("expected%s", sn)),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot(baseFuncName).Call(utils.CtxVar(), jen.ID(expectedSomethingID)),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName(sn)), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"surfaces sql.ErrNoRows",
				jen.ID(expectedSomethingID).Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot(baseFuncName).Call(utils.CtxVar(), jen.ID(expectedSomethingID)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				utils.AssertEqual(jen.Qual("database/sql", "ErrNoRows"), jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error querying database",
				jen.ID(expectedSomethingID).Assign().Add(utils.FakeUint64Func()),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot(baseFuncName).Call(utils.CtxVar(), jen.ID(expectedSomethingID)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with unscannable response",
				jen.ID(expectedSomethingID).Assign().Add(utils.FakeUint64Func()),
				jen.ID(utils.BuildFakeVarName(sn)).Assign().AddressOf().Qual(proj.ModelsV1Package(), sn).Valuesln(
					jen.ID("ID").MapAssign().Lit(321),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID(utils.BuildFakeVarName(sn)))),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot(baseFuncName).Call(utils.CtxVar(), jen.ID(expectedSomethingID)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
			),
		),
		jen.Line(),
	}
}
