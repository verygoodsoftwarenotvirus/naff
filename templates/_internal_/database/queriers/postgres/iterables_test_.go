package postgres

import (
	"fmt"

	"github.com/Masterminds/squirrel"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildBuildMockRowsFromSomethings(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierScanSomethings(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierSomethingExists(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierGetSomething(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierGetTotalSomethingCount(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierGetSomethings(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierGetSomethingsWithIDs(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierCreateSomething(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierUpdateSomething(proj, typ, dbvendor)...)
	code.Add(buildTestQuerierArchiveSomething(proj, typ, dbvendor)...)

	return code
}

func buildPrerequisiteIDsForTest(proj *models.Project, typ models.DataType, includeAccountID, includeSelf, idOnly, skip bool, skipIndex int) []jen.Code {
	lines := []jen.Code{}

	if typ.RestrictedToAccountMembers && includeAccountID {
		lines = append(lines, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
	}

	for i, dep := range proj.FindOwnerTypeChain(typ) {
		if !skip || (skip && i != skipIndex) {
			lines = append(lines, jen.IDf("example%sID", dep.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
		}
	}

	if includeSelf {
		if idOnly {
			lines = append(lines, jen.IDf("example%sID", typ.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
		} else {
			lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
		}
	}

	return lines
}

func buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p *models.Project, typ models.DataType, skipIndex int, includeSelf, idOnly, includeAccount bool) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	owners := p.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for i, pt := range owners {
		if i == skipIndex {
			params = append(params, jen.EmptyString())
		} else {
			params = append(params, jen.IDf("example%sID", pt.Name.Singular()))
		}
	}

	if includeSelf {
		if idOnly {
			params = append(params, jen.IDf("example%sID", sn))
		} else {
			params = append(params, jen.IDf("example%s", sn).Dot("ID"))
		}
	} else {
		params = append(params, jen.EmptyString())
	}

	if typ.RestrictedToAccountAtSomeLevel(p) {
		if includeAccount {
			params = append(params, jen.ID("exampleAccountID"))
		} else {
			params = append(params, jen.EmptyString())
		}
	}

	return params
}

func buildBuildMockRowsFromSomethings(_ *models.Project, typ models.DataType, _ wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()

	valuesLines := []jen.Code{jen.ID("x").Dot("ID")}

	for _, field := range typ.Fields {
		valuesLines = append(valuesLines, jen.ID("x").Dot(field.Name.Singular()))
	}

	valuesLines = append(valuesLines,
		jen.ID("x").Dot("CreatedOn"),
		jen.ID("x").Dot("LastUpdatedOn"),
		jen.ID("x").Dot("ArchivedOn"),
	)

	if typ.BelongsToStruct != nil {
		valuesLines = append(valuesLines, jen.ID("x").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		valuesLines = append(valuesLines, jen.ID("x").Dot("BelongsToAccount"))
	}

	lines := []jen.Code{
		jen.Func().IDf("buildMockRowsFrom%s", pn).Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").Uint64(), jen.ID(puvn).Op("...").Op("*").ID("types").Dot(sn)).Params(jen.Op("*").ID("sqlmock").Dot("Rows")).Body(
			jen.ID("columns").Assign().IDf("%sTableColumns", puvn),
			jen.Newline(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Equals().ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.Newline(),
			jen.ID("exampleRows").Assign().ID("sqlmock").Dot("NewRows").Call(jen.ID("columns")),
			jen.Newline(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().ID(puvn)).Body(
				jen.ID("rowValues").Assign().Index().ID("driver").Dot("Value").Valuesln(valuesLines...),
				jen.Newline(),
				jen.If(jen.ID("includeCounts")).Body(
					jen.ID("rowValues").Equals().ID("append").Call(
						jen.ID("rowValues"),
						jen.ID("filteredCount"),
						jen.ID("len").Call(jen.ID(puvn)),
					),
				),
				jen.Newline(),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Newline(),
			jen.Return().ID("exampleRows"),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierScanSomethings(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Scan%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("surfaces row errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Assign().Op("&").Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.False(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs row closing errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Assign().Op("&").Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.False()),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Nil()),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.False(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierSomethingExists(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	tableName := typ.Name.PluralRouteName()

	eqArgs := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): whatever,
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		eqArgs[fmt.Sprintf("%s.belongs_to_account", tableName)] = whatever
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = whatever
	}

	whereValues := typ.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClauseForTests(proj, false)
	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Suffix(existenceSuffix).
		Where(whereValues)

	_, args, err := qb.ToSql()
	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	bodyLines := []jen.Code{jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, true, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.ID("args").Assign().Index().Interface().Valuesln(
					dbCallArgs...,
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.True())),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true, true)...,
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("True").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		subtestLines := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			constants.CreateCtx(),
			jen.Newline(),
			jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, true, true, i), true)...),
			jen.Newline(),
			jen.List(jen.ID("c"), jen.Underscore()).Assign().ID("buildTestClient").Call(jen.ID("t")),
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
				buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, i, true, true, true)...,
			),
			jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.Err()),
			jen.Qual(constants.AssertionLibrary, "False").Call(jen.ID("t"), jen.ID("actual")),
		}

		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					subtestLines...,
				),
			),
			jen.Newline(),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, true, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, false, true, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "False").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with sql.ErrNoRows"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, true, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.ID("args").Assign().Index().Interface().Valuesln(
					dbCallArgs...,
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true, true)...,
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("False").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error executing query"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, true, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.ID("args").Assign().Index().Interface().Valuesln(
					dbCallArgs...,
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true, true)...,
				),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("False").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_%sExists", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierGetSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	tableName := typ.Name.PluralRouteName()

	eqArgs := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): whatever,
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		eqArgs[fmt.Sprintf("%s.belongs_to_account", tableName)] = whatever
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = whatever
	}

	whereValues := typ.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClauseForTests(proj, true)
	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Suffix(existenceSuffix).
		Where(whereValues)

	_, args, err := qb.ToSql()
	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	bodyLines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, false, -1), true)...),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("get%sQuery", sn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
					jen.False(),
					jen.Zero(),
					jen.IDf("example%s", sn),
				)),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, false, true)...,
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Equal").Call(
					jen.ID("t"),
					jen.IDf("example%s", sn),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
	}

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", owner.Name.SingularCommonName()),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, true, true, i), true)...),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(
						buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, i, true, true, true)...,
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.Err(),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, true, false, -1), true)...),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, false, true, true)...,
				),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid account ID"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, false, -1), true)...),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, false, true)...,
				),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error executing query"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, false, -1), true)...),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("get%sQuery", sn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, false, true)...,
				),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			bodyLines...,
		),
	}

	return lines
}

func buildTestQuerierGetTotalSomethingCount(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_GetTotal%sCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("exampleCount").Assign().Uint64().Call(jen.Lit(123)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("getTotal%sCountQuery", pn))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.Uint64().Call(jen.Lit(123)))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetTotal%sCount", sn).Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleCount"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("getTotal%sCountQuery", pn))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetTotal%sCount", sn).Call(jen.ID("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildArgsForListRetrievalQueryBuilder(p *models.Project, typ models.DataType, includeIncludeArchived, skip bool, skipIndex int) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if skip && i == skipIndex {
			lp = append(lp, jen.Zero())
		} else {
			lp = append(lp, jen.IDf("example%sID", pt.Name.Singular()))
		}
	}

	if typ.RestrictedToAccountAtSomeLevel(p) {
		lp = append(lp, jen.ID("exampleAccountID"))
	}
	lp = append(lp, utils.ConditionalCode(includeIncludeArchived, jen.False()), jen.ID(constants.FilterVarName))

	params = append(params, lp...)

	return params
}

func buildTestQuerierGetSomethings(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	uvn := typ.Name.UnexportedVarName()

	callArgs := buildArgsForListRetrievalQueryBuilder(proj, typ, false, false, -1)

	bodyLines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, false, false, -1), false)...),
				jen.IDf("example%sList", sn).Assign().ID("fakes").Dotf("BuildFake%sList", sn).Call(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dot("buildListQuery").Callln(
					jen.ID("ctx"),
					jen.Lit(puvn),
					func() jen.Code {
						if len(proj.FindOwnerTypeChain(typ)) > 0 {
							return jen.IDf("get%sJoins", pn)
						}
						return jen.Nil()
					}(),
					jen.Nil(),
					jen.ID("accountOwnershipColumn"),
					jen.IDf("%ssTableColumns", uvn),
					func() jen.Code {
						if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
							return jen.ID("exampleAccountID")
						}
						return jen.EmptyString()
					}(),
					jen.False(),
					jen.ID("filter"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
					jen.True(),
					jen.IDf("example%sList", sn).Dot("FilteredCount"),
					jen.IDf("example%sList", sn).Dot(pn).Op("..."),
				)),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					callArgs...,
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Equal").Call(
					jen.ID("t"),
					jen.IDf("example%sList", sn),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with nil filter"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Parens(jen.Op("*").Qual(proj.TypesPackage(), "QueryFilter")).Call(jen.Nil()),
				jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
				jen.IDf("example%sList", sn).Assign().ID("fakes").Dotf("BuildFake%sList", sn).Call(),
				jen.IDf("example%sList", sn).Dot("Page").Equals().Zero(),
				jen.IDf("example%sList", sn).Dot("Limit").Equals().Zero(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dot("buildListQuery").Callln(
					jen.ID("ctx"),
					jen.Lit(puvn),
					func() jen.Code {
						if len(proj.FindOwnerTypeChain(typ)) > 0 {
							return jen.IDf("get%sJoins", pn)
						}
						return jen.Nil()
					}(),
					jen.Nil(),
					jen.ID("accountOwnershipColumn"),
					jen.IDf("%ssTableColumns", uvn),
					func() jen.Code {
						if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
							return jen.ID("exampleAccountID")
						}
						return jen.EmptyString()
					}(),
					jen.False(),
					jen.ID("filter"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
					jen.True(),
					jen.IDf("example%sList", sn).Dot("FilteredCount"),
					jen.IDf("example%sList", sn).Dot(pn).Op("..."),
				)),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					jen.ID("ctx"),
					jen.ID("exampleAccountID"),
					jen.ID("filter"),
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Equal").Call(
					jen.ID("t"),
					jen.IDf("example%sList", sn),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid account ID"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					jen.ID("ctx"),
					jen.Lit(""),
					jen.ID("filter"),
				),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error executing query"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
				jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dot("buildListQuery").Callln(
					jen.ID("ctx"),
					jen.Lit(puvn),
					func() jen.Code {
						if len(proj.FindOwnerTypeChain(typ)) > 0 {
							return jen.IDf("get%sJoins", pn)
						}
						return jen.Nil()
					}(),
					jen.Nil(),
					jen.ID("accountOwnershipColumn"),
					jen.IDf("%ssTableColumns", uvn),
					func() jen.Code {
						if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
							return jen.ID("exampleAccountID")
						}
						return jen.EmptyString()
					}(),
					jen.False(),
					jen.ID("filter"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					jen.ID("ctx"),
					jen.ID("exampleAccountID"),
					jen.ID("filter"),
				),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with erroneous response from database"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
				jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dot("buildListQuery").Callln(
					jen.ID("ctx"),
					jen.Lit(puvn),
					func() jen.Code {
						if len(proj.FindOwnerTypeChain(typ)) > 0 {
							return jen.IDf("get%sJoins", pn)
						}
						return jen.Nil()
					}(),
					jen.Nil(),
					jen.ID("accountOwnershipColumn"),
					jen.IDf("%ssTableColumns", uvn),
					func() jen.Code {
						if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
							return jen.ID("exampleAccountID")
						}
						return jen.EmptyString()
					}(),
					jen.False(),
					jen.ID("filter"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
					jen.ID("ctx"),
					jen.ID("exampleAccountID"),
					jen.ID("filter"),
				),
				jen.ID("assert").Dot("Error").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
				jen.ID("assert").Dot("Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
	)

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierGetSomethingsWithIDs(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Get%sWithIDs", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Assign().ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dotf("buildGet%sWithIDsQuery", pn).Call(
						constants.CtxVar(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.False(),
						jen.Zero(),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.Zero(),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.IDf("example%sList", sn).Dot(pn),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("defaultLimit"),
						jen.Nil(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid IDs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("defaultLimit"),
						jen.Nil(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Assign().ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dotf("buildGet%sWithIDsQuery", pn).Call(
						constants.CtxVar(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(constants.ObligatoryError()),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error scanning query results"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Assign().ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dotf("buildGet%sWithIDsQuery", pn).Call(
						constants.CtxVar(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccountID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierCreateSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	argsLines := []jen.Code{jen.ID("exampleInput").Dot("ID")}
	for _, field := range typ.Fields {
		argsLines = append(argsLines, jen.ID("exampleInput").Dot(field.Name.Singular()))
	}

	if typ.BelongsToStruct != nil {
		argsLines = append(argsLines, jen.ID("exampleInput").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		argsLines = append(argsLines, jen.ID("exampleInput").Dot("BelongsToAccount"))
	}

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Create%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.IDf("example%s", sn).Dot("ID").Equals().Lit("1"),
					jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(argsLines...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sCreationQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.IDf("example%s", sn),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.Nil(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("expectedErr").Assign().Qual("errors", "New").Call(jen.ID("t").Dot("Name").Call()),
					jen.IDf("example%s", sn).Assign().ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.ID("exampleInput").Assign().ID("fakes").Dotf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(argsLines...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sCreationQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.ID("expectedErr")),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.ID("err"),
							jen.ID("expectedErr"),
						),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierUpdateSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()

	argsLines := []jen.Code{}

	for _, field := range typ.Fields {
		argsLines = append(argsLines, jen.IDf("example%s", sn).Dot(field.Name.Singular()))
	}

	if typ.BelongsToStruct != nil {
		argsLines = append(argsLines, jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		argsLines = append(argsLines, jen.IDf("example%s", sn).Dot("BelongsToAccount"))
	}

	argsLines = append(argsLines, jen.IDf("example%s", sn).Dot("ID"))

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Update%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(argsLines...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("update%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.Nil(),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(argsLines...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("update%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierArchiveSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	scn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Archive%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.IDf("example%sID", sn).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.IDf("example%sID", sn)),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("archive%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.IDf("example%sID", sn))),
					jen.Newline(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%sID", sn),
							jen.ID("exampleAccountID"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID("T").Dot("Run").Call(
						jen.Litf("with invalid %s ID", typ.BelongsToStruct.SingularCommonName()),
						jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
							jen.ID("t").Dot("Parallel").Call(),
							jen.Newline(),
							jen.IDf("example%sID", sn).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
							jen.Newline(),
							jen.ID("ctx").Assign().Qual("context", "Background").Call(),
							jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
							jen.Newline(),
							jen.ID("assert").Dot("Error").Call(
								jen.ID("t"),
								jen.ID("c").Dotf("Archive%s", sn).Call(
									jen.ID("ctx"),
									jen.IDf("example%sID", sn),
									jen.Lit(""),
								),
							),
						),
					)
				}
				return jen.Null()
			}(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.Lit(""),
							jen.ID("exampleAccountID"),
						),
					),
				),
			),
			jen.Newline(),
			utils.ConditionalCode(typ.BelongsToAccount,
				jen.ID("T").Dot("Run").Call(
					jen.Lit("with invalid account ID"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.IDf("example%sID", sn).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
						jen.Newline(),
						jen.ID("ctx").Assign().Qual("context", "Background").Call(),
						jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
						jen.Newline(),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("c").Dotf("Archive%s", sn).Call(
								jen.ID("ctx"),
								jen.IDf("example%sID", sn),
								jen.Lit(""),
							),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.IDf("example%sID", sn).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.IDf("example%sID", sn)),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("archive%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%sID", sn),
							jen.ID("exampleAccountID"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}
