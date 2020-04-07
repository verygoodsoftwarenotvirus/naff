package queriers

import (
	"fmt"
	"github.com/Masterminds/squirrel"
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

func buildRequisiteIDDeclarations(proj *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		if varPrefix != "" {
			lines = append(lines, jen.IDf("%s%sID", varPrefix, pt.Name.Singular()).Assign().Add(utils.FakeUint64Func()))
		} else {
			lines = append(lines, jen.IDf("%sID", pt.Name.UnexportedVarName()).Assign().Add(utils.FakeUint64Func()))
		}
	}

	if varPrefix != "" {
		lines = append(lines, jen.IDf("%s%sID", varPrefix, typ.Name.Singular()).Assign().Add(utils.FakeUint64Func()))
	} else {
		lines = append(lines, jen.IDf("%sID", typ.Name.UnexportedVarName()).Assign().Add(utils.FakeUint64Func()))
	}

	if typ.BelongsToUser {
		if varPrefix != "" {
			lines = append(lines, jen.IDf("%sUserID", varPrefix).Assign().Add(utils.FakeUint64Func()))
		} else {
			lines = append(lines, jen.ID("userID").Assign().Add(utils.FakeUint64Func()))
		}
	}

	return lines
}

func buildRequisiteExampleDeclarations(proj *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, pt := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("example%s", pt.Name.Singular()).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", pt.Name.Singular())).Call())
	}
	lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())

	//if typ.BelongsToUser {
	//	lines = append(lines, jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser"))
	//}

	return lines
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

func buildStringColumns(typ models.DataType) []string {
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

func buildStringColumnsAsString(typ models.DataType) string {
	return strings.Join(buildStringColumns(typ), ", ")
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
			jen.ID("includeCount").Bool(),
			jen.ID(puvn).Spread().PointerTo().Qual(proj.ModelsV1Package(), sn),
		).Params(
			jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
		).Block(
			jen.ID("columns").Assign().IDf("%sTableColumns", puvn),
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
			jen.ParamPointer().Qual("github.com/DATA-DOG/go-sqlmock", "Rows"),
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

	qb := queryBuilderForDatabase(dbvendor)
	query, _, _ := qb.Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName).
		Where(eqArgs).
		Suffix(existenceSuffix).
		ToSql()

	expectationArgs := []jen.Code{
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.IDf("example%s", sn).Dot("BelongsToUser")
			}
			return jen.Null()
		}(),
		jen.IDf("example%s", sn).Dot("ID"),
	}
	callArgs := typ.BuildGetSomethingArgsWithExampleVariables(proj)

	return buildQueryTest(
		proj,
		dbvendor,
		typ,
		fmt.Sprintf("%sExists", sn),
		query,
		expectationArgs,
		callArgs[1:],
		false,
		false,
		false,
		false,
	)
}

func buildTestDBSomethingExists(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
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
			utils.CtxVar(), jen.IDf("example%s", sn).Dot("ID"),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", sn).Dot("BelongsToUser"), jen.IDf("example%s", sn).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))
			actualCallArgs = append(actualCallArgs, jen.IDf("example%s", sn).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			btssn := typ.BelongsToStruct.Singular()

			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", btssn).Dot("ID"), jen.IDf("example%s", sn).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True()))
			actualCallArgs = append(actualCallArgs, jen.IDf("example%s", btssn).Dot("ID"))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", sn).Dot("ID")).
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
		jen.Func().IDf("Test%s_%sExists", dbvsn, sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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

	qb := queryBuilderForDatabase(dbvendor)
	query, _, _ := qb.Select(buildStringColumnsAsString(typ)).
		From(tableName).
		Where(eqArgs).
		ToSql()

	expectationArgs := []jen.Code{
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.IDf("example%s", sn).Dot("BelongsToUser")
			}
			return jen.Null()
		}(),
		jen.IDf("example%s", sn).Dot("ID"),
	}
	callArgs := typ.BuildGetSomethingArgsWithExampleVariables(proj)

	return buildQueryTest(
		proj,
		dbvendor,
		typ,
		fmt.Sprintf("Get%s", sn),
		query,
		expectationArgs,
		callArgs[1:],
		false,
		false,
		false,
		false,
	)
}

func buildTestDBGetSomething(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	rn := typ.Name.RouteName()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	buildFirstSubtestBlock := func() []jen.Code {
		lines := []jen.Code{
			utils.BuildFakeVar(proj, sn),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%s", sn).Dot("BelongsToUser").Equals().IDf("example%s", typ.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToUser {
					return jen.IDf("example%s", sn).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID")
				}
				return jen.Null()
			}(),
		}

		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.IDf("example%s", sn).Dot("ID"),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", sn).Dot("BelongsToUser"), jen.IDf("example%s", sn).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.IDf("example%s", sn)))
			actualCallArgs = append(actualCallArgs, jen.IDf("example%s", sn).Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"), jen.IDf("example%s", sn).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.IDf("example%s", sn)))
			actualCallArgs = append(actualCallArgs, jen.IDf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", sn).Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.IDf("example%s", sn)))
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", sn).Call(actualCallArgs...),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.IDf("example%s", sn), jen.ID("actual"), nil),
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
					return jen.IDf("example%s", sn).Dot("BelongsToUser").Equals().IDf("example%s", typ.BelongsToStruct.Singular()).Dot("ID")
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToUser {
					return jen.IDf("example%s", sn).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID")
				}
				return jen.Null()
			}(),
		}

		actualCallArgs := []jen.Code{
			utils.CtxVar(), jen.IDf("example%s", sn).Dot("ID"),
		}
		var mockDBCall jen.Code
		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.ID("exampleUser").Dot("ID"), jen.IDf("example%s", sn).Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.ID("exampleUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"), jen.IDf("example%s", sn).Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.IDf("example%s", typ.BelongsToStruct.Singular()).Dot("ID"))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Call(jen.IDf("example%s", sn).Dot("ID")).
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
	query, _, _ := qb.Select(buildStringColumnsAsString(typ)).
		From(tableName).
		Where(eqArgs).
		ToSql()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_Get%s", dbvsn, sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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

	return buildQueryTest(
		proj,
		dbvendor,
		typ,
		fmt.Sprintf("GetAll%sCount", pn),
		fmt.Sprintf("SELECT COUNT(%s.id) FROM %s WHERE %s.archived_on IS NULL", tableName, tableName, tableName),
		[]jen.Code{},
		[]jen.Code{},
		true,
		false,
		false,
		false,
	)
}

func buildTestDBGetAllSomethingCount(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	lines := []jen.Code{
		jen.Func().IDf("Test%s_GetAll%sCount", dbvsn, pn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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

func applyFleshedOutQueryFilter(qb *squirrel.SelectBuilder, tableName string) string {
	expectedQuery, _, _ := qb.
		Where(squirrel.Gt{fmt.Sprintf("%s.created_on", tableName): whateverValue}).
		Where(squirrel.Lt{fmt.Sprintf("%s.created_on", tableName): whateverValue}).
		Where(squirrel.Gt{fmt.Sprintf("%s.updated_on", tableName): whateverValue}).
		Where(squirrel.Lt{fmt.Sprintf("%s.updated_on", tableName): whateverValue}).
		GroupBy(fmt.Sprintf("%s.id", tableName)).
		Limit(20).
		Offset(180).
		ToSql()

	return expectedQuery
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
	cols := buildStringColumns(typ)
	expectedArgs := []jen.Code{}
	equals := squirrel.Eq{fmt.Sprintf("%s.archived_on", tableName): nil}

	if typ.BelongsToUser {
		equals[fmt.Sprintf("%s.belongs_to_user", tableName)] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID("exampleUser").Dot("ID"))
	}

	qb := queryBuilderForDatabase(dbvendor).Select(append(cols, fmt.Sprintf(countQuery, tableName))...).
		From(tableName).
		Where(equals)
	expectedQuery := applyFleshedOutQueryFilter(&qb, tableName)
	expectedArgs = appendFleshedOutQueryFilterArgs(expectedArgs)

	return buildQueryTest(
		proj,
		dbvendor,
		typ,
		fmt.Sprintf("Get%s", pn),
		expectedQuery,
		expectedArgs,
		[]jen.Code{},
		false,
		true,
		true,
		true,
	)
}

func buildTestDBGetListOfSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbfl := string(dbvendor.RouteName()[0])
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()
	cols := buildStringColumns(typ)

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
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUser").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.IDf("expected%s", sn)))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.IDf("expected%s", sn)))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			expectQueryMock = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.IDf("expected%s", sn)))
		}

		lines = append(lines,
			jen.ID("expectedCountQuery").Assign().Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tableName),
			jen.IDf("expected%s", sn).Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(
				jen.ID("ID").MapAssign().Lit(321),
			),
			jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(666)),
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().Lit(1),
					jen.ID("Limit").MapAssign().Lit(20),
					jen.ID("TotalCount").MapAssign().ID("expectedCount"),
				),
				jen.ID(pn).MapAssign().Index().Qual(proj.ModelsV1Package(), sn).Valuesln(
					jen.PointerTo().IDf("expected%s", sn),
				),
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			expectQueryMock,
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
				Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("count"))).
				Dot("AddRow").Call(jen.ID("expectedCount"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Get%s", pn).Call(actualCallArgs...),
			jen.Line(),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUser").Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows"))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUser").Dot("ID"))
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

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
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
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUser").Dot("ID")).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah")))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUser").Dot("ID"))
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

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
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
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUser").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.ID("expected")))
		}

		lines = append(lines,
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(
				jen.ID("ID").MapAssign().Lit(321),
			),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
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

	buildFifthSubtest := func() []jen.Code {
		lines := []jen.Code{}
		var mockDBCall jen.Code
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.Qual(proj.ModelsV1Package(), "DefaultQueryFilter").Call(),
		}

		if typ.BelongsToUser {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.ID("expectedUser").Dot("ID")).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WithArgs").Call(jen.IDf("expected%sID", typ.BelongsToStruct.Singular())).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.ID("expected")))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			mockDBCall = jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
				Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.ID("expected")))
		}

		lines = append(lines,
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(
				jen.ID("ID").MapAssign().Lit(321),
			),
			jen.ID("expectedCountQuery").Assign().Litf("SELECT COUNT(id) FROM %s WHERE archived_on IS NULL", tableName),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			mockDBCall,
			jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCountQuery"))).
				Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
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
		jen.Func().IDf("Test%s_Get%s", dbvsn, pn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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
			jen.Line(),
			utils.BuildSubTest("with error querying for count", buildFifthSubtest()...),
		),
		jen.Line(),
	}
}

func buildTestDBCreateSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

	fieldCols := buildNonStandardStringColumns(typ)

	var ips []string
	for i := range typ.Fields {
		ips = append(ips, getIncIndex(dbvendor, uint(i)))
	}

	if typ.BelongsToStruct != nil || typ.BelongsToUser {
		ips = append(ips, getIncIndex(dbvendor, uint(len(ips))))
	}

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

	expectedValues := []jen.Code{jen.ID("ID").MapAssign().Lit(321)}
	if typ.BelongsToUser {
		expectedQuery = fmt.Sprintf("INSERT INTO %s (%s,belongs_to_user%s) VALUES (%s%s)%s",
			tableName,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	}
	if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("INSERT INTO %s (%s,belongs_to_%s%s) VALUES (%s%s)%s",
			tableName,
			strings.ReplaceAll(fieldCols, " ", ""),
			typ.BelongsToStruct.RouteName(),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	} else {
		expectedQuery = fmt.Sprintf("INSERT INTO %s (%s%s) VALUES (%s%s)%s",
			tableName,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	}

	creationEqualityExpectations := buildCreationEqualityExpectations("expected", typ)
	createQueryTestBody := []jen.Code{
		jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
		jen.ID("expectedArgCount").Assign().Lit(1 + thisFuncExpectedArgCount),
		jen.ID("expectedQuery").Assign().Lit(expectedQuery),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildCreate%sQuery", sn).Call(jen.ID("expected")),
		jen.Line(),
		utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
		utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
	}
	createQueryTestBody = append(createQueryTestBody, creationEqualityExpectations...)

	return []jen.Code{
		jen.Func().IDf("Test%s_buildCreate%sQuery", dbvsn, sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path", createQueryTestBody...),
		),
		jen.Line(),
	}
}

func buildTestDBCreateSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

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
	if typ.BelongsToUser || typ.BelongsToStruct != nil {
		ips = append(ips, getIncIndex(dbvendor, uint(len(ips))))
	}
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
			tableName,
			strings.ReplaceAll(fieldCols, " ", ""),
			createdOnAddendum,
			insertPlaceholders,
			createdOnValueAdd,
			queryTail,
		)
	}
	if typ.BelongsToStruct != nil {
		expectedCreationQuery = fmt.Sprintf(
			"INSERT INTO %s (%s,belongs_to_%s%s) VALUES (%s%s)%s",
			tableName,
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
			tableName,
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
		}
		if typ.BelongsToStruct != nil {
			out = append(out, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		}

		return out
	}

	buildFirstSubtest := func(proj *models.Project, typ models.DataType) []jen.Code {
		out := []jen.Code{}
		expectedValues := []jen.Code{}

		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out = append(out,
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
			jen.ID("expectedInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn)).Valuesln(expectedInputFields...),
		)

		if isPostgres(dbvendor) {
			out = append(out,
				jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("id"), jen.Lit("created_on"))).Dot("AddRow").Call(
					jen.ID("expected").Dot("ID"),
					jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
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
					Dotln("WithArgs").Callln(nef...).Dot("WillReturnResult").Call(
					jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
						jen.ID("int64").Call(jen.ID("expected").Dot("ID")),
						jen.Lit(123),
					),
				),
				jen.Line(),
				jen.ID("expectedTimeQuery").Assign().Litf("SELECT created_on FROM %s WHERE id = %s", tableName, getIncIndex(dbvendor, 0)),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedTimeQuery"))).
					Dotln("WithArgs").Call(jen.ID("expected").Dot("ID")).
					Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(
					jen.Index().String().Values(jen.Lit("created_on")),
				).Dot("AddRow").Call(
					jen.ID("expected").Dot("CreatedOn")),
				),
				jen.Line(),
			)
		}

		out = append(out,
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Create%s", sn).Call(utils.CtxVar(), jen.ID("expectedInput")),
			utils.AssertNoError(jen.Err(), nil),
			utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
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

		expectedValues := []jen.Code{}

		out := []jen.Code{}
		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		out = append(out,
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
			jen.ID("expectedInput").Assign().VarPointer().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn)).Valuesln(expectedInputFields...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedCreationQuery"))).
				Dotln("WithArgs").Callln(nef...).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dotf("Create%s", sn).Call(utils.CtxVar(), jen.ID("expectedInput")),
			utils.AssertError(jen.Err(), nil),
			utils.AssertNil(jen.ID("actual"), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Create%s", dbvsn, sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedCreationQuery").Assign().Lit(expectedCreationQuery),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path", buildFirstSubtest(proj, typ)...),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error writing to database", buildSecondSubtest()...),
		),
		jen.Line(),
	}
}

func buildTestBuildUpdateSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()

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

	expectedValues := []jen.Code{jen.ID("ID").MapAssign().Lit(321)}

	if typ.BelongsToUser {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s",
			tableName,
			updateColsStr,
			getTimeQuery(dbvendor),
			getIncIndex(dbvendor, uint(len(updateCols))),
			getIncIndex(dbvendor, uint(len(updateCols)+1)),
			queryTail,
		)
		varCount = len(updateCols) + 2
	}
	if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_%s = %s AND id = %s%s",
			tableName,
			updateColsStr,
			getTimeQuery(dbvendor),
			typ.BelongsToStruct.RouteName(),
			getIncIndex(dbvendor, uint(len(updateCols))),
			getIncIndex(dbvendor, uint(len(updateCols)+1)),
			queryTail,
		)
		varCount = len(updateCols) + 2
	} else {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE id = %s%s",
			tableName,
			updateColsStr,
			getTimeQuery(dbvendor),
			getIncIndex(dbvendor, uint(len(updateCols))),
			queryTail,
		)
		varCount = len(updateCols) + 1
	}

	testBuildUpdateQueryBody := []jen.Code{
		jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
		jen.ID("expectedArgCount").Assign().Lit(varCount), // +2 because of ID and BelongsTo
		jen.ID("expectedQuery").Assign().Lit(expectedQuery),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildUpdate%sQuery", sn).Call(jen.ID("expected")),
		jen.Line(),
		utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
		utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
	}

	testBuildUpdateQueryBody = append(testBuildUpdateQueryBody, creationEqualityExpectations...)
	testBuildUpdateQueryBody = append(testBuildUpdateQueryBody,
		utils.AssertEqual(jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(len(creationEqualityExpectations))).Assert(jen.Uint64()), nil),
	)

	return []jen.Code{
		jen.Func().IDf("Test%s_buildUpdate%sQuery", dbvsn, sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				testBuildUpdateQueryBody...,
			),
		),
		jen.Line(),
	}
}

func buildTestDBUpdateSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()

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
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_user = %s AND id = %s%s", tableName, updateColsStr, getTimeQuery(dbvendor), getIncIndex(dbvendor, uint(len(updateCols))), getIncIndex(dbvendor, uint(len(updateCols))+1), queryTail)
	}
	if typ.BelongsToStruct != nil {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE belongs_to_%s = %s AND id = %s%s", tableName, updateColsStr, getTimeQuery(dbvendor), typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, uint(len(updateCols))), getIncIndex(dbvendor, uint(len(updateCols))+1), queryTail)
	} else {
		expectedQuery = fmt.Sprintf("UPDATE %s SET %s, updated_on = %s WHERE id = %s%s", tableName, updateColsStr, getTimeQuery(dbvendor), getIncIndex(dbvendor, uint(len(updateCols))), queryTail)
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
		if typ.BelongsToNobody {
			expectQueryArgs = append(expectQueryArgs, jen.ID("expected").Dot("ID"))
		}

		if isPostgres(dbvendor) {
			expectFuncName = "ExpectQuery"
			returnFuncName = "WillReturnRows"

			exRows = jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(
				jen.Index().String().Values(jen.Lit("updated_on")),
			).Dot("AddRow").Call(
				jen.Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
			)
		} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
			expectFuncName = "ExpectExec"
			returnFuncName = "WillReturnResult"

			exRows = jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
				jen.ID("int64").Call(jen.ID("expected").Dot("ID")),
				jen.Lit(123),
			)
		}

		lines := []jen.Code{}
		expectedValues := []jen.Code{}
		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}
		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		lines = append(lines,
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
			exRows,
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(
				expectQueryArgs...,
			).Dot(returnFuncName).Call(jen.ID("exampleRows")),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(utils.CtxVar(), jen.ID("expected")),
			utils.AssertNoError(jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return lines
	}

	buildSecondSubtest := func(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
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
		expectedValues := []jen.Code{}

		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
		} else if typ.BelongsToNobody {
			expectQueryArgs = append(expectQueryArgs, jen.ID("expected").Dot("ID"))
		}

		expectedValues = append(expectedValues,
			jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
		)

		out = append(out,
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot(expectFuncName).Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(expectQueryArgs...).Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Update%s", sn).Call(utils.CtxVar(), jen.ID("expected")),
			utils.AssertError(jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return out
	}

	return []jen.Code{
		jen.Func().IDf("Test%s_Update%s", dbvendor.Singular(), typ.Name.Singular()).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTestWithoutContext("happy path", buildFirstSubTest(typ)...),
			jen.Line(),
			utils.BuildSubTestWithoutContext("with error writing to database", buildSecondSubtest(proj, dbvendor, typ)...),
			jen.Line(),
		),
	}
}

func buildTestDBArchiveSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	var (
		expectedQuery string
		queryTail     string
		queryArgCount int
	)

	if isPostgres(dbvendor) {
		queryTail = " RETURNING archived_on"
	}
	expectedValues := []jen.Code{
		jen.ID("ID").MapAssign().Lit(321),
	}
	archiveQueryBuildingParams := []jen.Code{
		jen.ID("expected").Dot("ID"),
	}

	if typ.BelongsToUser {
		queryArgCount = 2
		expectedQuery = fmt.Sprintf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s%s", tableName, getTimeQuery(dbvendor), getTimeQuery(dbvendor), getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), queryTail)
		archiveQueryBuildingParams = append(archiveQueryBuildingParams, jen.ID("expected").Dot("BelongsToUser"))
	}
	if typ.BelongsToStruct != nil {
		queryArgCount = 2
		expectedQuery = fmt.Sprintf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_%s = %s AND id = %s%s", tableName, getTimeQuery(dbvendor), getTimeQuery(dbvendor), typ.BelongsToStruct.RouteName(), getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), queryTail)
		archiveQueryBuildingParams = append(archiveQueryBuildingParams, jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	} else {
		queryArgCount = 1
		expectedQuery = fmt.Sprintf("UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND id = %s%s", tableName, getTimeQuery(dbvendor), getTimeQuery(dbvendor), getIncIndex(dbvendor, 0), queryTail)
	}

	testLines := []jen.Code{
		jen.List(jen.ID(dbfl), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
		jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
		jen.ID("expectedArgCount").Assign().Lit(queryArgCount),
		jen.ID("expectedQuery").Assign().Lit(expectedQuery),
		jen.List(jen.ID("actualQuery"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildArchive%sQuery", sn).Call(archiveQueryBuildingParams...),
		jen.Line(),
		utils.AssertEqual(jen.ID("expectedQuery"), jen.ID("actualQuery"), nil),
		utils.AssertLength(jen.ID("args"), jen.ID("expectedArgCount"), nil),
	}

	var assertIndesx int
	if typ.BelongsToUser {
		assertIndesx = 1
		testLines = append(testLines,
			utils.AssertEqual(jen.ID("expected").Dot("BelongsToUser"), jen.ID("args").Index(jen.Zero()).Assert(jen.Uint64()), nil),
		)
	}
	if typ.BelongsToStruct != nil {
		assertIndesx = 1
		testLines = append(testLines,
			utils.AssertEqual(jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()), jen.ID("args").Index(jen.Zero()).Assert(jen.Uint64()), nil),
		)
	} else if typ.BelongsToNobody {
		assertIndesx = 0
	}

	testLines = append(testLines,
		utils.AssertEqual(jen.ID("expected").Dot("ID"), jen.ID("args").Index(jen.Lit(assertIndesx)).Assert(jen.Uint64()), nil),
	)

	return []jen.Code{
		jen.Func().IDf("Test%s_buildArchive%sQuery", dbvsn, sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path", testLines...),
		),
		jen.Line(),
	}
}

func buildTestDBArchiveSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	buildSubtestOne := func() []jen.Code {
		var (
			dbQueryExpectationArgs []jen.Code
			dbQuery                string
			queryTail              string
		)

		if isPostgres(dbvendor) {
			queryTail = " RETURNING archived_on"
		}

		expectedValues := []jen.Code{}
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.ID("expected").Dot("ID"),
		}

		if typ.BelongsToUser {
			expectedValues = append(expectedValues, jen.ID("BelongsToUser").MapAssign().ID("expectedUser").Dot("ID"))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			expectedValues = append(expectedValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}

		expectedValues = append(expectedValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))
		block := []jen.Code{}

		if typ.BelongsToUser {
			dbQuery = fmt.Sprintf(
				"UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_user = %s AND id = %s%s",
				tableName,
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				queryTail,
			)
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("expected").Dot("BelongsToUser"))
		}
		if typ.BelongsToStruct != nil {
			dbQuery = fmt.Sprintf(
				"UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND belongs_to_%s = %s AND id = %s%s",
				tableName,
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				typ.BelongsToStruct.RouteName(),
				getIncIndex(dbvendor, 0),
				getIncIndex(dbvendor, 1),
				queryTail,
			)
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("expected").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
		} else {
			dbQuery = fmt.Sprintf(
				"UPDATE %s SET updated_on = %s, archived_on = %s WHERE archived_on IS NULL AND id = %s%s",
				tableName,
				getTimeQuery(dbvendor),
				getTimeQuery(dbvendor),
				getIncIndex(dbvendor, 0),
				queryTail,
			)
		}

		dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("expected").Dot("ID"))

		block = append(block,
			jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(expectedValues...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mockDB").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedQuery"))).
				Dotln("WithArgs").Callln(dbQueryExpectationArgs...).Dot("WillReturnResult").Call(
				jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewResult").Call(
					jen.Lit(123),
					jen.Lit(123),
				),
			),
			jen.Line(),
			jen.Err().Assign().ID(dbfl).Dotf("Archive%s", sn).Call(actualCallArgs...),
			utils.AssertNoError(jen.Err(), nil),
			jen.Line(),
			utils.AssertNoError(jen.ID("mockDB").Dot("ExpectationsWereMet").Call(), jen.Lit("not all database expectations were met")),
		)

		return []jen.Code{
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedQuery").Assign().Lit(dbQuery),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path", block...),
		}
	}

	buildSubtestTwo := func() []jen.Code {
		exampleValues := []jen.Code{}

		var dbQueryExpectationArgs []jen.Code
		block := []jen.Code{}
		actualCallArgs := []jen.Code{
			utils.CtxVar(),
			jen.ID("example").Dot("ID"),
		}

		if typ.BelongsToUser {
			exampleValues = append(exampleValues, jen.ID("BelongsToUser").MapAssign().ID("expectedUser").Dot("ID"))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("example").Dot("BelongsToUser"))
			actualCallArgs = append(actualCallArgs, jen.ID("expectedUser").Dot("ID"))
		}
		if typ.BelongsToStruct != nil {
			exampleValues = append(exampleValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().IDf("expected%sID", typ.BelongsToStruct.Singular()))
			dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("example").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
			actualCallArgs = append(actualCallArgs, jen.IDf("expected%sID", typ.BelongsToStruct.Singular()))
		}

		dbQueryExpectationArgs = append(dbQueryExpectationArgs, jen.ID("example").Dot("ID"))
		exampleValues = append(exampleValues, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()))

		block = append(block,
			jen.ID("example").Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(exampleValues...),
			jen.Line(),
			jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
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

		return []jen.Code{
			utils.BuildSubTestWithoutContext(
				"with error writing to database", block...),
		}
	}

	var bodyContents []jen.Code

	bodyContents = append(bodyContents, buildSubtestOne()...)
	bodyContents = append(bodyContents, jen.Line())
	bodyContents = append(bodyContents, buildSubtestTwo()...)

	return []jen.Code{
		jen.Func().IDf("Test%s_Archive%s", dbvsn, sn).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			bodyContents...,
		),
		jen.Line(),
	}
}

func buildTestDBGetAllSomethingForSomethingElseFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbrn := dbvendor.RouteName()
	dbfl := string(dbrn[0])
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	tableName := typ.Name.PluralRouteName()
	cols := buildStringColumnsAsString(typ)

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
		jen.Func().ID(testFuncName).Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("expectedListQuery").Assign().Lit(expectedQuery),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID(expectedSomethingID).Assign().Add(utils.FakeUint64Func()),
				jen.IDf("expected%s", sn).Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(
					jen.ID("ID").MapAssign().Lit(321),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", sn).Call(jen.False(), jen.IDf("expected%s", sn))),
				jen.Line(),
				jen.ID("expected").Assign().Index().Qual(proj.ModelsV1Package(), sn).Values(jen.PointerTo().IDf("expected%s", sn)),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID(dbfl).Dot(baseFuncName).Call(utils.CtxVar(), jen.ID(expectedSomethingID)),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
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
				jen.IDf("example%s", sn).Assign().VarPointer().Qual(proj.ModelsV1Package(), sn).Valuesln(
					jen.ID("ID").MapAssign().Lit(321),
				),
				jen.Line(),
				jen.List(jen.ID(dbfl), jen.ID("mockDB")).Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("mockDB").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("expectedListQuery"))).
					Dotln("WithArgs").Call(jen.ID(expectedSomethingID)).
					Dotln("WillReturnRows").Call(jen.IDf("buildErroneousMockRowFrom%s", sn).Call(jen.IDf("example%s", sn))),
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
