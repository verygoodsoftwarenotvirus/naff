package querier

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildBadFields(varName string, typ models.DataType) []jen.Code {
	fields := []jen.Code{jen.ID(varName).Dot("ID"), jen.ID(varName).Dot("ExternalID")}

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
		jen.Func().IDf("buildMockRowsFrom%s", pn).Params(jen.ID("includeCounts").ID("bool"), jen.ID("filteredCount").ID("uint64"), jen.ID(puvn).Op("...").Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Op("*").Qual("github.com/DATA-DOG/go-sqlmock", "Rows")).Body(
			jen.ID("columns").Assign().Qual(proj.QuerybuildersPackage(), fmt.Sprintf("%sTableColumns", pn)),
			jen.Line(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Op("=").ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.Line(),
			jen.ID("exampleRows").Assign().Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.ID("columns")),
			jen.Line(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().ID(puvn)).Body(
				jen.ID("rowValues").Assign().Index().ID("driver").Dot("Value").Valuesln(
					buildBadFields("x", typ)...,
				),
				jen.Line(),
				jen.If(jen.ID("includeCounts")).Body(
					jen.ID("rowValues").Op("=").ID("append").Call(
						jen.ID("rowValues"),
						jen.ID("filteredCount"),
						jen.ID("len").Call(jen.ID(puvn)),
					)),
				jen.Line(),
				jen.ID("exampleRows").Dot("AddRow").Call(jen.ID("rowValues").Op("...")),
			),
			jen.Line(),
			jen.Return().ID("exampleRows"),
		),
		jen.Line(),
	}
}

func buildTestQuerier_ScanSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Scan%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("surfaces row errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockRows").Assign().Op("&").Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs row closing errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockRows").Assign().Op("&").Qual(proj.DatabasePackage(), "MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.ID("nil")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("mockRows"),
						jen.ID("false"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_SomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_%sExists", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("Build%sExistsQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.Qual("github.com/DATA-DOG/go-sqlmock", "NewRows").Call(jen.Index().ID("string").Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.ID("true"))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with sql.ErrNoRows"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("Build%sExistsQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("Build%sExistsQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_GetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	scn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf(fmt.Sprintf("BuildFake%s", sn))).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.IDf("example%s", sn),
					)),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.IDf("example%s", sn).Dot("BelongsToAccount"),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.Lit(0),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_GetAllSomethingsCount(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	sn := typ.Name.Singular()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_GetAll%sCount", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.ID("exampleCount").Assign().ID("uint64").Call(jen.Lit(123)),
					jen.Line(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(jen.ID("fakeQuery")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("uint64").Call(jen.Lit(123)))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAll%sCount", pn).Call(jen.ID("ctx")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleCount"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_GetAllSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_GetAll%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().Op("*").Qual(proj.TypesPackage(), sn)),
					jen.ID("doneChan").Assign().ID("make").Call(
						jen.Chan().ID("bool"),
						jen.Lit(1),
					),
					jen.ID("expectedCount").Assign().ID("uint64").Call(jen.Lit(20)),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Line(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Line(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Line(),
					jen.ID("stillQuerying").Assign().ID("true"),
					jen.For(jen.ID("stillQuerying")).Body(
						jen.Select().Body(
							jen.Case(jen.ID("batch").Assign().Op("<-").ID("results")).Body(
								jen.ID("assert").Dot("NotEmpty").Call(
									jen.ID("t"),
									jen.ID("batch"),
								), jen.ID("doneChan").ReceiveFromChannel().ID("true")),
							jen.Case(jen.Op("<-").Qual("time", "After").Call(jen.Qual("time", "Second"))).Body(
								jen.ID("t").Dot("FailNow").Call()),
							jen.Case(jen.Op("<-").ID("doneChan")).Body(
								jen.ID("stillQuerying").Op("=").ID("false")),
						)),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil results channel"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("nil"),
							jen.ID("exampleBatchSize"),
						),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with now rows returned"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().Op("*").Qual(proj.TypesPackage(), sn)),
					jen.ID("expectedCount").Assign().ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Line(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Line(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.Line(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Line(),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error fetching initial count"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().Op("*").Qual(proj.TypesPackage(), sn)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Line(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("err").Assign().ID("c").Dotf("GetAll%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("results"),
						jen.ID("exampleBatchSize"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error querying database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().Op("*").Qual(proj.TypesPackage(), sn)),
					jen.ID("expectedCount").Assign().ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Line(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Line(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Line(),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.Line(),
					jen.ID("results").Assign().ID("make").Call(jen.Chan().Index().Op("*").Qual(proj.TypesPackage(), sn)),
					jen.ID("expectedCount").Assign().ID("uint64").Call(jen.Lit(20)),
					jen.ID("exampleBatchSize").Assign().ID("uint16").Call(jen.Lit(1000)),
					jen.Line(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("_")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAll%sCountQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.ID("expectedCount"))),
					jen.Line(),
					jen.List(jen.ID("secondFakeQuery"), jen.ID("secondFakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetBatchOf%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("uint64").Call(jen.Lit(1)),
						jen.ID("uint64").Call(jen.ID("exampleBatchSize").Op("+").Lit(1)),
					).Dot("Return").Call(
						jen.ID("secondFakeQuery"),
						jen.ID("secondFakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("secondFakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("secondFakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Line(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("GetAll%s", pn).Call(
							jen.ID("ctx"),
							jen.ID("results"),
							jen.ID("exampleBatchSize"),
						),
					),
					jen.Line(),
					jen.Qual("time", "Sleep").Call(jen.Qual("time", "Second")),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_GetSomethings(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(
						jen.IDf("buildMockRowsFrom%s", pn).Call(
							jen.ID("true"),
							jen.IDf("example%sList", sn).Dot("FilteredCount"),
							jen.IDf("example%sList", sn).Dot(pn).Op("..."),
						)),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil filter"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("filter").Assign().Parens(jen.Op("*").Qual(proj.TypesPackage(), "QueryFilter")).Call(jen.ID("nil")),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.IDf("example%sList", sn).Dot("Page").Op("=").Lit(0),
					jen.IDf("example%sList", sn).Dot("Limit").Op("=").Lit(0),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("true"),
						jen.IDf("example%sList", sn).Dot("FilteredCount"),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
						jen.ID("ctx"),
						jen.Lit(0),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("filter").Assign().Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%s", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_GetSomethingsWithIDs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Get%sWithIDs", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Line(),
					jen.Var().ID("exampleIDs").Index().ID("uint64"),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Op("=").ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sWithIDsQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
						jen.ID("false"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("defaultLimit"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Var().ID("exampleIDs").Index().ID("uint64"),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Op("=").ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.Lit(0),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("sets limit if not present"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Var().ID("exampleIDs").Index().ID("uint64"),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Op("=").ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sWithIDsQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
						jen.ID("false"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("false"),
						jen.Lit(0),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.Lit(0),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Var().ID("exampleIDs").Index().ID("uint64"),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Op("=").ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sWithIDsQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
						jen.ID("false"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%sList", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sList", sn)).Call(),
					jen.Var().ID("exampleIDs").Index().ID("uint64"),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Assign().Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Op("=").ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGet%sWithIDsQuery", pn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
						jen.ID("false"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Get%sWithIDs", pn).Call(
						jen.ID("ctx"),
						jen.ID("exampleAccount").Dot("ID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_CreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Create%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildCreate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.Line(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Op("=").Lit(""),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("nil"),
						jen.ID("exampleUser").Dot("ID"),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.Lit(0),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("expectedErr").Assign().Qual("errors", "New").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
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
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.ID("expectedErr")),
					jen.Line(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
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
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error creating audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildCreate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("ExternalID").Op("=").Lit(""),
					jen.ID("exampleInput").Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeCreationQuery"), jen.ID("fakeCreationArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildCreate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.ID("fakeCreationQuery"),
						jen.ID("fakeCreationArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeCreationQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeCreationArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("c").Dot("timeFunc").Op("=").Func().Params().Params(jen.ID("uint64")).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleUser").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_UpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Update%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.Line(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.Lit(0),
							jen.ID("nil"),
						),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeUpdateQuery"), jen.ID("fakeUpdateArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildUpdate%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn),
					).Dot("Return").Call(
						jen.ID("fakeUpdateQuery"),
						jen.ID("fakeUpdateArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeUpdateQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeUpdateArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
							jen.ID("exampleUser").Dot("ID"),
							jen.ID("nil"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_ArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_Archive%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildArchive%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectCommit").Call(),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.Lit(0),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.Lit(0),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.Lit(0),
						),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error beginning transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildArchive%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing audit log entry"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildArchive%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectRollback").Call(),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error committing transaction"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleAccount").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.ID("db").Dot("ExpectBegin").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildArchive%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newSuccessfulDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Line(),
					jen.ID("expectAuditLogEntryInTransaction").Call(
						jen.ID("mockQueryBuilder"),
						jen.ID("db"),
						jen.ID("nil"),
					),
					jen.Line(),
					jen.ID("db").Dot("ExpectCommit").Call().Dot("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Archive%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn).Dot("ID"),
							jen.ID("exampleAccount").Dot("ID"),
							jen.ID("exampleUser").Dot("ID"),
						),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
	}
}

func buildTestQuerier_GetAuditLogEntriesForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Func().IDf("TestQuerier_GetAuditLogEntriesFor%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.ID("exampleAuditLogEntriesList").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryList").Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAuditLogEntriesFor%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildMockRowsFromAuditLogEntries").Call(
						jen.ID("false"),
						jen.ID("exampleAuditLogEntriesList").Dot("Entries").Op("..."),
					)),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAuditLogEntriesList").Dot("Entries"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.Lit(0),
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
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error executing query"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAuditLogEntriesFor%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with erroneous response from database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Line(),
					jen.IDf("example%s", sn).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Line(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Assign().ID("buildTestClient").Call(jen.ID("t")),
					jen.Line(),
					jen.ID("mockQueryBuilder").Assign().Qual(proj.DatabasePackage(), "BuildMockSQLQueryBuilder").Call(),
					jen.Line(),
					jen.List(jen.ID("fakeQuery"), jen.ID("fakeArgs")).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeSQLQuery").Call(),
					jen.ID("mockQueryBuilder").Dotf("%sSQLQueryBuilder", sn).Dot("On").Callln(
						jen.Litf("BuildGetAuditLogEntriesFor%sQuery", sn),
						jen.Qual(proj.TestUtilsPackage(), "ContextMatcher"),
						jen.IDf("example%s", sn).Dot("ID"),
					).Dot("Return").Call(
						jen.ID("fakeQuery"),
						jen.ID("fakeArgs"),
					),
					jen.ID("c").Dot("sqlQueryBuilder").Op("=").ID("mockQueryBuilder"),
					jen.Line(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("fakeQuery"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("fakeArgs")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Assign().ID("c").Dotf("GetAuditLogEntriesFor%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Line(),
					jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
						jen.ID("mockQueryBuilder"),
					),
				),
			),
		),
		jen.Line(),
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
	code.Add(buildTestQuerier_GetSomethings(proj, typ)...)
	code.Add(buildTestQuerier_GetSomethingsWithIDs(proj, typ)...)
	code.Add(buildTestQuerier_CreateSomething(proj, typ)...)
	code.Add(buildTestQuerier_UpdateSomething(proj, typ)...)
	code.Add(buildTestQuerier_ArchiveSomething(proj, typ)...)
	code.Add(buildTestQuerier_GetAuditLogEntriesForSomething(proj, typ)...)

	return code
}
