package postgres

import (
	"fmt"

	"github.com/Masterminds/squirrel"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildPrerequisiteIDsForTest(proj *models.Project, typ models.DataType, includeAccountID, includeSelf, skip bool, skipIndex int) []jen.Code {
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
		lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	}

	return lines
}

func buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p *models.Project, typ models.DataType, skipIndex int, includeSelf, includeAccount bool) []jen.Code {
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
		params = append(params, jen.IDf("example%s", sn).Dot("ID"))
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

func buildMockDBRowFields(varName string, typ models.DataType) []jen.Code {
	fields := []jen.Code{
		jen.ID(varName).Dot("ID"),
	}

	for _, field := range typ.Fields {
		fields = append(fields, jen.ID(varName).Dot(field.Name.Singular()))
	}

	fields = append(fields,
		jen.ID(varName).Dot("CreatedOn"),
		jen.ID(varName).Dot("LastUpdatedOn"),
		jen.ID(varName).Dot("ArchivedOn"),
	)

	if typ.BelongsToStruct != nil {
		fields = append(fields, jen.ID(varName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		fields = append(fields, jen.ID(varName).Dot(constants.AccountOwnershipFieldName))
	}

	return fields
}

func buildBuildMockRowsFromSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()

	return []jen.Code{
		jen.Func().IDf("buildMockRowsFrom%s", pn).Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").Uint64(), jen.ID(puvn).Spread().PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.PointerTo().Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Body(
			jen.ID("columns").Assign().IDf("%sTableColumns", puvn),
			jen.Newline(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Equals().ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.Newline(),
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Newline(),
			jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().ID(puvn)).Body(
				jen.ID("rowValues").Assign().Index().ID("driver").Dot("Value").Valuesln(
					buildMockDBRowFields("x", typ)...,
				),
				jen.Newline(),
				jen.If(jen.ID("includeCounts")).Body(
					jen.ID("rowValues").Equals().ID("append").Call(
						jen.ID("rowValues"),
						jen.ID("filteredCount"),
						jen.ID("len").Call(jen.ID(puvn)),
					)),
				jen.Newline(),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Spread()),
			),
			jen.Newline(),
			jen.Return().ID("exampleRows"),
		),
		jen.Newline(),
	}
}

func buildTestQuerier_ScanSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Scan%s", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("surfaces row errs"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.Underscore(), jen.Underscore(), jen.Err()).Assign().ID("q").Dotf("scan%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.Err(),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs row closing errs"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Assign().AddressOf().Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.Underscore(), jen.Underscore(), jen.Err()).Assign().ID("q").Dotf("scan%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.Err(),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestQuerier_SomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
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
	}

	firstSubtestLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		constants.CreateCtx(),
		jen.Newline(),
		jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, -1), true)...),
		jen.Newline(),
		jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
		jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
		jen.Newline(),
		jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
			Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.ID("true"))),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
			buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)...,
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.Err(),
		),
		jen.Qual(constants.AssertionLibrary, "True").Call(
			jen.ID("t"),
			jen.ID("actual"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("db"),
		),
	}

	bodyLines = append(bodyLines,
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtestLines...,
			),
		),
		jen.Newline(),
	)

	for i, owner := range proj.FindOwnerTypeChain(typ) {
		subtestLines := []jen.Code{
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			constants.CreateCtx(),
			jen.Newline(),
			jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, true, i), true)...),
			jen.Newline(),
			jen.List(jen.ID("c"), jen.Underscore()).Assign().ID("buildTestClient").Call(jen.ID("t")),
			jen.Newline(),
			jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
				buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, i, true, true)...,
			),
			jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.Err()),
			jen.Qual(constants.AssertionLibrary, "False").Call(jen.ID("t"), jen.ID("actual")),
		}

		bodyLines = append(bodyLines,
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
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, false, true)...,
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
		utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid account ID"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, false, true, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, false)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "False").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		)),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with sql.ErrNoRows"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "False").Call(
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
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, -1), true)...),
				jen.Newline(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("%sExists", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "False").Call(
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

	return []jen.Code{
		jen.Func().IDf("TestQuerier_%sExists", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildTestQuerier_GetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	scn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
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

	whereValues := typ.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClauseForTests(proj, true)
	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Where(whereValues)

	_, args, err := qb.ToSql()
	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, -1), true)...),
		jen.Newline(),
		jen.ID("ctx").Assign().Qual("context", "Background").Call(),
		jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
		jen.Newline(),
		jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("get%sQuery", sn))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
			Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
			jen.ID("false"),
			jen.Lit(0),
			jen.IDf("example%s", sn),
		)),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(
			buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)...,
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.Err(),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(
			jen.ID("t"),
			jen.IDf("example%s", sn),
			jen.ID("actual"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("db"),
		),
	}

	bodyLines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
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
					jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, true, i), true)...),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(
						buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, i, true, true)...,
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
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, false, -1), true)...),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, false, true)...,
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

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, false, true, false, -1), true)...),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(
						buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, false)...,
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
			jen.Lit("with error executing query"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, -1), true)...),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("get%sQuery", sn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", sn).Call(
					buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
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

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(bodyLines...),
		jen.Newline(),
	}
}

func buildTestQuerier_GetAllSomethingsCount(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_GetTotal%sCount", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("exampleCount").Assign().Uint64().Call(jen.Lit(123)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("getTotal%sCountQuery", pn))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.Uint64().Call(jen.Lit(123)))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetTotal%sCount", sn).Call(jen.ID("ctx")),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.Err(),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("getTotal%sCountQuery", pn))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnError").Call(constants.ObligatoryError()),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetTotal%sCount", sn).Call(jen.ID("ctx")),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.Err(),
					),
					jen.Qual(constants.AssertionLibrary, "Zero").Call(
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
}

func buildArgsForListRetrievalQueryBuilder(p *models.Project, typ models.DataType, includeIncludeArchived, skip bool, skipIndex int) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for i, pt := range owners {
		if skip && i == skipIndex {
			lp = append(lp, jen.EmptyString())
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

func buildTestQuerier_GetListOfSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	uvn := typ.Name.UnexportedVarName()

	callArgs := buildArgsForListRetrievalQueryBuilder(proj, typ, false, false, -1)

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
		jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, false, -1), false)...),
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
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
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
			Dotln("WillReturnRows").Call(
			jen.IDf("buildMockRowsFrom%s", pn).Call(
				jen.ID("true"),
				jen.IDf("example%sList", sn).Dot("FilteredCount"),
				jen.IDf("example%sList", sn).Dot(pn).Spread(),
			)),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
			callArgs...,
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.Err(),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(
			jen.ID("t"),
			jen.IDf("example%sList", sn),
			jen.ID("actual"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("db"),
		),
	}

	bodyLines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
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
					jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
					jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, true, i), false)...),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
						buildArgsForListRetrievalQueryBuilder(proj, typ, false, true, i)...,
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(jen.ID("t"), jen.Err()),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(jen.ID("t"), jen.ID("actual")),
				),
			),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with nil filter"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Parens(jen.PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Call(jen.ID("nil")),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, false, -1), false)...),
				jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
				jen.IDf("example%sList", sn).Dot("Page").Equals().Lit(0),
				jen.IDf("example%sList", sn).Dot("Limit").Equals().Lit(0),
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
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
					jen.ID("true"),
					jen.IDf("example%sList", sn).Dot("FilteredCount"),
					jen.IDf("example%sList", sn).Dot(pn).Spread(),
				)),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
					callArgs...,
				),
				jen.Qual(constants.AssertionLibrary, "NoError").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
		utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid account ID"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
					jen.ID("ctx"),
					jen.EmptyString(),
					jen.ID("filter"),
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
		)),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error executing query"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, false, -1), false)...),
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
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
					callArgs...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
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
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
				jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, false, false, -1), false)...),
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
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%s", pn).Call(
					callArgs...,
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
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

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildTestQuerier_GetSomethingsWithIDs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
			}
			return jen.Null()
		}(),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
		jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
		jen.Newline(),
		jen.Var().ID("exampleIDs").Index().String(),
		jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
			jen.ID("exampleIDs").Equals().ID("append").Call(
				jen.ID("exampleIDs"),
				jen.ID("x").Dot("ID"),
			)),
		jen.Newline(),
		jen.ID("ctx").Assign().Qual("context", "Background").Call(),
		jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")), jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dotf("buildGet%sWithIDsQuery", pn).Call(
			constants.CtxVar(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
			jen.ID("defaultLimit"),
			jen.ID("exampleIDs"),
		),
		jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
			Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
			jen.ID("false"),
			jen.Lit(0),
			jen.IDf("example%sList", sn).Dot(pn).Spread(),
		)),
		jen.Newline(),
		jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
			jen.ID("ctx"),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
			jen.Zero(),
			jen.ID("exampleIDs"),
		),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.Err(),
		),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(
			jen.ID("t"),
			jen.IDf("example%sList", sn).Dot(pn),
			jen.ID("actual"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("db"),
		),
	}

	bodyLines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid IDs"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
					}
					return jen.Null()
				}(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
					constants.CtxVar(),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.Nil(),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
			),
		),
	}

	if typ.BelongsToStruct != nil {
		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", typ.BelongsToStruct.SingularCommonName()),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.EmptyString(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
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
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error executing query"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
					}
					return jen.Null()
				}(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
				jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
				jen.Newline(),
				jen.Var().ID("exampleIDs").Index().String(),
				jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
					jen.ID("exampleIDs").Equals().ID("append").Call(
						jen.ID("exampleIDs"),
						jen.ID("x").Dot("ID"),
					)),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dotf("buildGet%sWithIDsQuery", pn).Call(
					constants.CtxVar(),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.ID("exampleIDs"),
				),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
					jen.ID("ctx"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.ID("exampleIDs"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(
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
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
					}
					return jen.Null()
				}(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
				jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
				jen.Newline(),
				jen.Var().ID("exampleIDs").Index().String(),
				jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
					jen.ID("exampleIDs").Equals().ID("append").Call(
						jen.ID("exampleIDs"),
						jen.ID("x").Dot("ID"),
					)),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("c").Dotf("buildGet%sWithIDsQuery", pn).Call(
					constants.CtxVar(),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.ID("exampleIDs"),
				),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
				jen.Newline(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
					jen.ID("ctx"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.ID("exampleIDs"),
				),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.Err(),
				),
				jen.Qual(constants.AssertionLibrary, "Empty").Call(
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

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Get%sWithIDs", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildTestQuerier_CreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)

	creationColumns := []string{
		"id",
	}
	args := []interface{}{models.NewCodeWrapper(jen.ID("exampleInput").Dot("ID"))}
	for _, field := range typ.Fields {
		creationColumns = append(creationColumns, field.Name.RouteName())
		args = append(args, models.NewCodeWrapper(jen.ID("exampleInput").Dot(field.Name.Singular())))
	}

	if typ.BelongsToStruct != nil {
		creationColumns = append(creationColumns, typ.BelongsToStruct.RouteName())
		args = append(args, models.NewCodeWrapper(jen.ID("exampleInput").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular())))
	}

	if typ.BelongsToAccount {
		creationColumns = append(creationColumns, "belongs_to_account")
		args = append(args, models.NewCodeWrapper(jen.ID("exampleInput").Dot("BelongsToAccount")))
	}

	_, args, err := sqlBuilder.Insert(tableName).
		Columns(creationColumns...).
		Values(args...).
		ToSql()

	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Create%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ID").Equals().Lit("1"),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sCreationQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.Err(),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("nil"),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("expectedErr").Assign().Qual("errors", "New").Call(jen.ID("t").Dot("Name").Call()),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sCreationQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
						Dotln("WillReturnError").Call(jen.ID("expectedErr")),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.Err(),
					),
					jen.Qual(constants.AssertionLibrary, "True").Call(
						jen.ID("t"),
						jen.Qual("errors", "Is").Call(
							jen.Err(),
							jen.ID("expectedErr"),
						),
					),
					jen.Qual(constants.AssertionLibrary, "Nil").Call(
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
}

func buildTestQuerier_UpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)

	updateWhere := squirrel.Eq{
		"id":          models.NewCodeWrapper(jen.IDf("example%s", sn).Dot("ID")),
		"archived_on": nil,
	}
	argValues := []jen.Code{}

	if typ.BelongsToStruct != nil {
		updateWhere[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.IDf("example%s", sn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		updateWhere["belongs_to_account"] = models.NewCodeWrapper(jen.IDf("example%s", sn).Dot("BelongsToAccount"))
	}

	updateBuilder := sqlBuilder.Update(tableName)

	for _, field := range typ.Fields {
		argValues = append(argValues, jen.ID("updated").Dot(field.Name.Singular()))
		updateBuilder = updateBuilder.Set(field.Name.RouteName(), models.NewCodeWrapper(jen.IDf("example%s", sn).Dot(field.Name.Singular())))
	}

	if typ.BelongsToStruct != nil {
		argValues = append(argValues, jen.ID("updated").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToAccount {
		argValues = append(argValues, jen.ID("updated").Dot("BelongsToAccount"))
	}

	argValues = append(argValues, jen.ID("updated").Dot("ID"))

	updateBuilder = updateBuilder.Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).Where(updateWhere)

	_, args, err := updateBuilder.ToSql()
	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Update%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("update%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
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
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("update%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
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
}

func buildTestQuerier_ArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)
	archiveWhere := squirrel.Eq{
		"id":          models.NewCodeWrapper(jen.IDf("example%s", sn).Dot("ID")),
		"archived_on": nil,
	}

	if typ.BelongsToStruct != nil {
		archiveWhere[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		archiveWhere["belongs_to_account"] = models.NewCodeWrapper(jen.ID("exampleAccountID"))
	}

	_, args, err := sqlBuilder.Update(tableName).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(archiveWhere).ToSql()

	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
			}
			return jen.Null()
		}(),
		jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
		jen.Newline(),
		jen.ID("ctx").Assign().Qual("context", "Background").Call(),
		jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
		jen.Newline(),
		jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("archive%sQuery", sn))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
			Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
		jen.Newline(),
		jen.Qual(constants.AssertionLibrary, "NoError").Call(
			jen.ID("t"),
			jen.ID("c").Dotf("Archive%s", sn).Call(
				jen.ID("ctx"),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
					}
					return jen.Null()
				}(),
				jen.IDf("example%s", sn).Dot("ID"),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
			),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("db"),
		),
	}

	bodyLines := []jen.Code{
		jen.ID("T").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("standard"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				firstSubtest...,
			),
		),
	}

	if typ.BelongsToStruct != nil {
		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", typ.BelongsToStruct.SingularCommonName()),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.EmptyString(),
							jen.IDf("example%s", sn).Dot("ID"),
							utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						),
					),
				),
			),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Litf("with invalid %s ID", scn),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
					}
					return jen.Null()
				}(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("c").Dotf("Archive%s", sn).Call(
						jen.ID("ctx"),
						func() jen.Code {
							if typ.BelongsToStruct != nil {
								return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
							}
							return jen.Null()
						}(),
						jen.EmptyString(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					),
				),
			),
		),
	)

	if typ.BelongsToAccount {
		bodyLines = append(bodyLines,
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
						}
						return jen.Null()
					}(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							func() jen.Code {
								if typ.BelongsToStruct != nil {
									return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
								}
								return jen.Null()
							}(),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.EmptyString(),
						),
					),
				),
			),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error writing to database"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
					}
					return jen.Null()
				}(),
				jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
				jen.Newline(),
				jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("archive%sQuery", sn))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Spread()).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.Qual(constants.AssertionLibrary, "Error").Call(
					jen.ID("t"),
					jen.ID("c").Dotf("Archive%s", sn).Call(
						jen.ID("ctx"),
						func() jen.Code {
							if typ.BelongsToStruct != nil {
								return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
							}
							return jen.Null()
						}(),
						jen.IDf("example%s", sn).Dot("ID"),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
				),
			),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Archive%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildBuildMockRowsFromSomethings(proj, typ)...)
	code.Add(buildTestQuerier_ScanSomethings(proj, typ)...)
	code.Add(buildTestQuerier_SomethingExists(proj, typ)...)
	code.Add(buildTestQuerier_GetSomething(proj, typ)...)
	code.Add(buildTestQuerier_GetAllSomethingsCount(proj, typ)...)
	code.Add(buildTestQuerier_GetListOfSomethings(proj, typ)...)
	code.Add(buildTestQuerier_GetSomethingsWithIDs(proj, typ)...)
	code.Add(buildTestQuerier_CreateSomething(proj, typ)...)
	code.Add(buildTestQuerier_UpdateSomething(proj, typ)...)
	code.Add(buildTestQuerier_ArchiveSomething(proj, typ)...)

	return code
}
