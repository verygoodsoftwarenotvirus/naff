package querier

import (
	"fmt"

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
			params = append(params, jen.Zero())
		} else {
			params = append(params, jen.IDf("example%sID", pt.Name.Singular()))
		}
	}

	if includeSelf {
		params = append(params, jen.IDf("example%s", sn).Dot("ID"))
	} else {
		params = append(params, jen.Zero())
	}

	if typ.RestrictedToAccountAtSomeLevel(p) {
		if includeAccount {
			params = append(params, jen.ID("exampleAccountID"))
		} else {
			params = append(params, jen.Zero())
		}
	}

	return params
}

func buildMockDBRowFields(varName string, typ models.DataType) []jen.Code {
	fields := []jen.Code{jen.ID(varName).Dot("ID"), jen.ID(varName).Dot("ExternalID")}

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
			jen.ID("columns").Assign().Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)),
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
		jen.Newline(),
		jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
		jen.Newline(),
		jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
		jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
			append([]jen.Code{
				jen.Litf("Build%sExistsQuery", sn),
				jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			},
				buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)[1:]...,
			)...,
		).Dot("Return").Call(
			jen.ID("fakeQuery"),
			jen.ID("fakeArgs"),
		),
		jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
		jen.Newline(),
		jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
			jen.ID("mockQueryBuilder"),
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
				jen.Newline(),
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					append([]jen.Code{
						jen.Litf("Build%sExistsQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					},
						buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)[1:]...,
					)...,
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
					jen.ID("mockQueryBuilder"),
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
				jen.Newline(),
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					append([]jen.Code{
						jen.Litf("Build%sExistsQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					},
						buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)[1:]...,
					)...,
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
					jen.ID("mockQueryBuilder"),
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

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.Null().Add(utils.IntersperseWithNewlines(buildPrerequisiteIDsForTest(proj, typ, true, true, false, -1), true)...),
		jen.Newline(),
		jen.ID("ctx").Assign().Qual("context", "Background").Call(),
		jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
		jen.Newline(),
		jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
		jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
			append([]jen.Code{
				jen.Litf("BuildGet%sQuery", sn),
				jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			},
				buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)[1:]...,
			)...,
		).Dot("Return").Call(jen.ID("fakeQuery"), jen.ID("fakeArgs")),
		jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
		jen.Newline(),
		jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
			Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
			jen.ID("false"),
			jen.Zero(),
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
			jen.ID("mockQueryBuilder"),
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
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					append([]jen.Code{
						jen.Litf("BuildGet%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					},
						buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ, -1, true, true)[1:]...,
					)...,
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
					jen.ID("mockQueryBuilder"),
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
	pn := typ.Name.Plural()
	sn := typ.Name.Singular()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_GetAll%sCount", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleCount").Assign().Uint64().Call(jen.Lit(123)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(jen.ID("fakeQuery")),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.Uint64().Call(jen.Lit(123)))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAll%sCount", pn).Call(jen.ID("ctx")),
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
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestQuerier_GetAllSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_GetAll%s", pn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().PointerTo().Qual(proj.TypesPackage(), sn)),
					jen.ID("doneChan").Assign().ID("make").Call(
						jen.Chan().ID("bool"),
						jen.Lit(1),
					),
					jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(20)),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Newline(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Uint64().Call(jen.Lit(1)),
						jen.Uint64().Call(jen.ID("exampleBatchSize").Plus().Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Spread()).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("false"),
						jen.Zero(),
						jen.IDf("example%sList", sn).Dot(pn).Spread(),
					)),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Newline(),
					jen.ID("stillQuerying").Assign().ID("true"),
					jen.For(jen.ID("stillQuerying")).Body(
						jen.Select().Body(
							jen.Case(jen.ID("batch").Assign().Op("<-").ID("results")).Body(
								jen.Qual(constants.AssertionLibrary, "NotEmpty").Call(
									jen.ID("t"),
									jen.ID("batch"),
								), jen.ID("doneChan").ReceiveFromChannel().ID("true")),
							jen.Case(jen.Op("<-").Qual("time", "After").Call(jen.Qual("time", "Second"))).Body(
								jen.ID("t").Dot("FailNow").Call()),
							jen.Case(jen.Op("<-").ID("doneChan")).Body(
								jen.ID("stillQuerying").Equals().ID("false")),
						)),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil results channel"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.ID("exampleBatchSize"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with now rows returned"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().PointerTo().Qual(proj.TypesPackage(), sn)),
					jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Newline(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Uint64().Call(jen.Lit(1)),
						jen.Uint64().Call(jen.ID("exampleBatchSize").Plus().Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Spread()).
						Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Newline(),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching initial count"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().PointerTo().Qual(proj.TypesPackage(), sn)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.Err().Assign().ID("c").Dotf("GetAll%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.ID("exampleBatchSize"),
					),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.Err(),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error querying database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().PointerTo().Qual(proj.TypesPackage(), sn)),
					jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Newline(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Uint64().Call(jen.Lit(1)),
						jen.Uint64().Call(jen.ID("exampleBatchSize").Plus().Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Spread()).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Newline(),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid response from database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().PointerTo().Qual(proj.TypesPackage(), sn)),
					jen.ID("expectedCount").Assign().Uint64().Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Newline(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.Uint64().Call(jen.Lit(1)),
						jen.Uint64().Call(jen.ID("exampleBatchSize").Plus().Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Spread()).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Newline(),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
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

func buildTestQuerier_GetListOfSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	mockCallArgs := append([]jen.Code{
		jen.Litf("BuildGet%sQuery", pn),
		jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
	},
		buildArgsForListRetrievalQueryBuilder(proj, typ, true, false, -1)[1:]...,
	)

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
		jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
		jen.Newline(),
		jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
		jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
			mockCallArgs...,
		).Dot("Return").Call(
			jen.ID("fakeQuery"),
			jen.ID("fakeArgs"),
		),
		jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
		jen.Newline(),
		jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
			jen.ID("mockQueryBuilder"),
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
				jen.IDf("example%sList", sn).Dot("Page").Equals().Zero(),
				jen.IDf("example%sList", sn).Dot("Limit").Equals().Zero(),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					mockCallArgs...,
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
					jen.ID("mockQueryBuilder"),
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
					jen.Zero(),
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
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					mockCallArgs...,
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
					jen.ID("mockQueryBuilder"),
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
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					mockCallArgs...,
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
					jen.ID("mockQueryBuilder"),
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
		jen.Var().ID("exampleIDs").Index().Uint64(),
		jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
			jen.ID("exampleIDs").Equals().ID("append").Call(
				jen.ID("exampleIDs"),
				jen.ID("x").Dot("ID"),
			)),
		jen.Newline(),
		jen.ID("ctx").Assign().Qual("context", "Background").Call(),
		jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
		jen.Newline(),
		jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
		jen.Newline(),
		jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
		jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
			jen.Litf("BuildGet%sWithIDsQuery", pn),
			jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
			jen.ID("defaultLimit"),
			jen.ID("exampleIDs"),
			func() jen.Code {
				if typ.BelongsToAccount {
					if typ.RestrictedToAccountAtSomeLevel(proj) {
						return jen.True()
					} else {
						return jen.False()
					}
				}
				return jen.Null()
			}(),
		).Dot("Return").Call(
			jen.ID("fakeQuery"),
			jen.ID("fakeArgs"),
		),
		jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
		jen.Newline(),
		jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
			Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
			jen.ID("false"),
			jen.Zero(),
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
			jen.ID("defaultLimit"),
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
			jen.ID("mockQueryBuilder"),
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
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().Uint64(),
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
						jen.Zero(),
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
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Var().ID("exampleIDs").Index().Uint64(),
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
						func() jen.Code {
							if typ.BelongsToStruct != nil {
								return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
							}
							return jen.Null()
						}(),
						jen.Zero(),
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
		jen.ID("T").Dot("Run").Call(
			jen.Lit("sets limit if not present"),
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
				jen.Var().ID("exampleIDs").Index().Uint64(),
				jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
					jen.ID("exampleIDs").Equals().ID("append").Call(
						jen.ID("exampleIDs"),
						jen.ID("x").Dot("ID"),
					)),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					jen.Litf("BuildGet%sWithIDsQuery", pn),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.ID("exampleIDs"),
					func() jen.Code {
						if typ.RestrictedToAccountAtSomeLevel(proj) {
							return jen.True()
						} else if typ.BelongsToAccount {
							return jen.False()
						}
						return jen.Null()
					}(),
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
					Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
					jen.ID("false"),
					jen.Zero(),
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
					jen.ID("mockQueryBuilder"),
				),
			),
		),
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
				jen.Var().ID("exampleIDs").Index().Uint64(),
				jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
					jen.ID("exampleIDs").Equals().ID("append").Call(
						jen.ID("exampleIDs"),
						jen.ID("x").Dot("ID"),
					)),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					jen.Litf("BuildGet%sWithIDsQuery", pn),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.ID("exampleIDs"),
					func() jen.Code {
						if typ.RestrictedToAccountAtSomeLevel(proj) {
							return jen.True()
						} else if typ.BelongsToAccount {
							return jen.False()
						}
						return jen.Null()
					}(),
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
					jen.ID("mockQueryBuilder"),
				),
			),
		),
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with erroneous response from database"),
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
				jen.Var().ID("exampleIDs").Index().Uint64(),
				jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
					jen.ID("exampleIDs").Equals().ID("append").Call(
						jen.ID("exampleIDs"),
						jen.ID("x").Dot("ID"),
					)),
				jen.Newline(),
				jen.ID("ctx").Assign().Qual("context", "Background").Call(),
				jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
				jen.Newline(),
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					jen.Litf("BuildGet%sWithIDsQuery", pn),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
					jen.ID("defaultLimit"),
					jen.ID("exampleIDs"),
					func() jen.Code {
						if typ.RestrictedToAccountAtSomeLevel(proj) {
							return jen.True()
						} else if typ.BelongsToAccount {
							return jen.False()
						}
						return jen.Null()
					}(),
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
				jen.Newline(),
				jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
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
				jen.Qual(constants.AssertionLibrary, "Nil").Call(
					jen.ID("t"),
					jen.ID("actual"),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
					jen.ID("mockQueryBuilder"),
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

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Create%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Equals().Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildCreate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
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
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Equals().Lit(""),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.ID("exampleUser").Dot("ID"),
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
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Equals().Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.Zero(),
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
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Equals().Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
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
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("expectedErr").Assign().Qual("errors", "New").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildCreate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
						Dotln("WillReturnError").Call(jen.ID("expectedErr")),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
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
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error creating audit log entry"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Equals().Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildCreate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
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
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Equals().Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildCreate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
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
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestQuerier_UpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Update%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.Zero(),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
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
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Spread()).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Newline(),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry to database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Newline(),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Spread()).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
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

	firstSubtest := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
		jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
		jen.Newline(),
		jen.ID("db").Dot("ExpectBegin").Call(),
		jen.Newline(),
		jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
		jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
			jen.Litf("BuildArchive%sQuery", sn),
			jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			jen.IDf("example%s", sn).Dot("ID"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("exampleAccountID")),
		).Dot("Return").Call(
			jen.ID("fakeQuery"),
			jen.ID("fakeArgs"),
		),
		jen.Newline(),
		jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
			Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
			Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
		jen.Newline(),
		jen.ID("expectAuditLogEntryInTransaction").Call(
			jen.ID("mockQueryBuilder"),
			jen.ID("db"),
			jen.ID("nil"),
		),
		jen.Newline(),
		jen.ID("db").Dot("ExpectCommit").Call(),
		jen.Newline(),
		jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
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
				jen.ID("exampleUserID"),
			),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("db"),
			jen.ID("mockQueryBuilder"),
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
					jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
							jen.Zero(),
							jen.IDf("example%s", sn).Dot("ID"),
							utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
							jen.ID("exampleUserID"),
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
				jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
						jen.Zero(),
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.ID("exampleUserID"),
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
					jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
							jen.Zero(),
							jen.ID("exampleUserID"),
						),
					),
				),
			),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with invalid actor ID"),
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
						utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
						jen.Zero(),
					),
				),
			),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error beginning transaction"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
				jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
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
						jen.ID("exampleUserID"),
					),
				),
			),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error writing to database"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.ID("db").Dot("ExpectBegin").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					jen.Litf("BuildArchive%sQuery", sn),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					jen.IDf("example%s", sn).Dot("ID"),
					utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("exampleAccountID")),
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
					Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.ID("db").Dot("ExpectRollback").Call(),
				jen.Newline(),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
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
						jen.ID("exampleUserID"),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
					jen.ID("mockQueryBuilder"),
				),
			),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error writing audit log entry"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.ID("db").Dot("ExpectBegin").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					jen.Litf("BuildArchive%sQuery", sn),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					jen.IDf("example%s", sn).Dot("ID"),
					utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("exampleAccountID")),
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
					Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
				jen.Newline(),
				jen.ID("expectAuditLogEntryInTransaction").Call(
					jen.ID("mockQueryBuilder"),
					jen.ID("db"),
					jen.Qual("errors", "New").Call(jen.Lit("blah")),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectRollback").Call(),
				jen.Newline(),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
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
						jen.ID("exampleUserID"),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
					jen.ID("mockQueryBuilder"),
				),
			),
		),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("T").Dot("Run").Call(
			jen.Lit("with error committing transaction"),
			jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("exampleUserID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
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
				jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
				jen.Newline(),
				jen.ID("db").Dot("ExpectBegin").Call(),
				jen.Newline(),
				jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
				jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
					jen.Litf("BuildArchive%sQuery", sn),
					jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
						}
						return jen.Null()
					}(),
					jen.IDf("example%s", sn).Dot("ID"),
					utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("exampleAccountID")),
				).Dot("Return").Call(
					jen.ID("fakeQuery"),
					jen.ID("fakeArgs"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
					Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
					Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
				jen.Newline(),
				jen.ID("expectAuditLogEntryInTransaction").Call(
					jen.ID("mockQueryBuilder"),
					jen.ID("db"),
					jen.ID("nil"),
				),
				jen.Newline(),
				jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.Newline(),
				jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
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
						jen.ID("exampleUserID"),
					),
				),
				jen.Newline(),
				jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
					jen.ID("t"),
					jen.ID("db"),
					jen.ID("mockQueryBuilder"),
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

func buildTestQuerier_GetAuditLogEntriesForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_GetAuditLogEntriesFor%s", sn).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleAuditLogEntriesList").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryList").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAuditLogEntriesFor%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
						Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromAuditLogEntries").Call(
						jen.ID("false"),
						jen.ID("exampleAuditLogEntriesList").Dot("Entries").Spread(),
					)),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.Err(),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAuditLogEntriesList").Dot("Entries"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.Zero(),
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
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAuditLogEntriesFor%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
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
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Newline(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAuditLogEntriesFor%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Equals().ID("mockQueryBuilder"),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Spread()).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
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
						jen.ID("mockQueryBuilder"),
					),
				),
			),
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
	code.Add(buildTestQuerier_GetAllSomethings(proj, typ)...)
	code.Add(buildTestQuerier_GetListOfSomethings(proj, typ)...)
	code.Add(buildTestQuerier_GetSomethingsWithIDs(proj, typ)...)
	code.Add(buildTestQuerier_CreateSomething(proj, typ)...)
	code.Add(buildTestQuerier_UpdateSomething(proj, typ)...)
	code.Add(buildTestQuerier_ArchiveSomething(proj, typ)...)
	code.Add(buildTestQuerier_GetAuditLogEntriesForSomething(proj, typ)...)

	return code
}
