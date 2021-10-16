package mysql

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
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

func buildBuildMockRowsFromSomethings(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
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
			jen.ID("columns").Op(":=").IDf("%sTableColumns", puvn),
			jen.Newline(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("columns").Equals().ID("append").Call(
					jen.ID("columns"),
					jen.Lit("filtered_count"),
					jen.Lit("total_count"),
				)),
			jen.Newline(),
			jen.ID("exampleRows").Op(":=").ID("sqlmock").Dot("NewRows").Call(jen.ID("columns")),
			jen.Newline(),
			jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID(puvn)).Body(
				jen.ID("rowValues").Op(":=").Index().ID("driver").Dot("Value").Valuesln(valuesLines...),
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dotf("scan%s", pn).Call(
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("logs row closing errs"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("mockRows").Op(":=").Op("&").ID("database").Dot("MockResultIterator").Values(),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Next")).Dot("Return").Call(jen.ID("false")),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Err")).Dot("Return").Call(jen.Nil()),
					jen.ID("mockRows").Dot("On").Call(jen.Lit("Close")).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dotf("scan%s", pn).Call(
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
		jen.Newline(),
	}

	return lines
}

func buildTestQuerierSomethingExists(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_%sExists", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sID", sn).Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"),
						jen.IDf("example%sID", sn),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("sqlmock").Dot("NewRows").Call(jen.Index().String().Values(jen.Lit("exists"))).Dot("AddRow").Call(jen.ID("true"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%sID", sn),
						jen.ID("exampleAccountID"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("exampleAccountID"),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.IDf("example%sID", sn).Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%sID", sn),
						jen.Lit(""),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with sql.ErrNoRows"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sID", sn).Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.IDf("example%sID", sn)),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("database/sql", "ErrNoRows")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%sID", sn),
						jen.ID("exampleAccountID"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sID", sn).Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.IDf("example%sID", sn)),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sExistenceQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%sID", sn),
						jen.ID("exampleAccountID"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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

func buildTestQuerierGetSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.IDf("example%s", sn).Dot("ID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("get%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("false"),
						jen.Zero(),
						jen.IDf("example%s", sn),
					)),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccountID"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("exampleAccountID"),
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
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.Lit(""),
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccountID"), jen.IDf("example%s", sn).Dot("ID")),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("get%sQuery", sn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccountID"),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleCount").Op(":=").Uint64().Call(jen.Lit(123)),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("getTotal%sCountQuery", pn))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnRows").Call(jen.ID("newCountDBRowResponse").Call(jen.Uint64().Call(jen.Lit(123)))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetTotal%sCount", sn).Call(jen.ID("ctx")),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("getTotal%sCountQuery", pn))).
						Dotln("WithArgs").Call().
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetTotal%sCount", sn).Call(jen.ID("ctx")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Zero").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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

func buildTestQuerierGetSomethings(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	uvn := typ.Name.UnexportedVarName()

	lines := []jen.Code{
		jen.Func().IDf("TestQuerier_Get%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Op(":=").ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
						jen.ID("ctx"),
						jen.Lit(puvn),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.IDf("%ssTableColumns", uvn),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("true"),
						jen.IDf("example%sList", sn).Dot("FilteredCount"),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.Nil()),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Op(":=").ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.IDf("example%sList", sn).Dot("Page").Equals().Zero(),
					jen.IDf("example%sList", sn).Dot("Limit").Equals().Zero(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
						jen.ID("ctx"),
						jen.Lit(puvn),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.IDf("%ssTableColumns", uvn),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.IDf("buildMockRowsFrom%s", pn).Call(
						jen.ID("true"),
						jen.IDf("example%sList", sn).Dot("FilteredCount"),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
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
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
						jen.ID("ctx"),
						jen.Lit(puvn),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.IDf("%ssTableColumns", uvn),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("c").Dot("buildListQuery").Callln(
						jen.ID("ctx"),
						jen.Lit(puvn),
						jen.Nil(),
						jen.Nil(),
						jen.ID("accountOwnershipColumn"),
						jen.IDf("%ssTableColumns", uvn),
						jen.ID("exampleAccountID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("db").Dot("ExpectQuery").Call(jen.ID("formatQueryForSQLMock").Call(jen.ID("query"))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnRows").Call(jen.ID("buildErroneousMockRow").Call()),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Op(":=").ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
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
						jen.ID("false"),
						jen.Zero(),
						jen.IDf("example%sList", sn).Dot(pn).Op("..."),
					)),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sWithIDs", pn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sWithIDs", pn).Call(
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sWithIDs", pn).Call(
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Op(":=").ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sWithIDs", pn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sList", sn).Op(":=").ID("fakes").Dotf("BuildFake%sList", sn).Call(),
					jen.Newline(),
					jen.Var().ID("exampleIDs").Index().String(),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().IDf("example%sList", sn).Dot(pn)).Body(
						jen.ID("exampleIDs").Equals().ID("append").Call(
							jen.ID("exampleIDs"),
							jen.ID("x").Dot("ID"),
						),
					),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sWithIDs", pn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.IDf("example%s", sn).Dot("ID").Equals().Lit("1"),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dotf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(argsLines...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sCreationQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnResult").Call(jen.ID("newArbitraryDatabaseResult").Call(jen.IDf("example%s", sn).Dot("ID"))),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
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
					jen.ID("expectedErr").Op(":=").Qual("errors", "New").Call(jen.ID("t").Dot("Name").Call()),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dotf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(argsLines...),
					jen.Newline(),
					jen.ID("db").Dot("ExpectExec").Call(jen.ID("formatQueryForSQLMock").Call(jen.IDf("%sCreationQuery", uvn))).
						Dotln("WithArgs").Call(jen.ID("interfaceToDriverValue").Call(jen.ID("args")).Op("...")).
						Dotln("WillReturnError").Call(jen.ID("expectedErr")),
					jen.Newline(),
					jen.ID("c").Dot("timeFunc").Equals().Func().Params().Params(jen.Uint64()).Body(
						jen.Return().IDf("example%s", sn).Dot("CreatedOn")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(argsLines...),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
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
				jen.Lit("with invalid actor ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("c").Dotf("Update%s", sn).Call(
							jen.ID("ctx"),
							jen.IDf("example%s", sn),
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
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(argsLines...),
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sID", sn).Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("db"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid account ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.IDf("example%sID", sn).Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
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
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing to database"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleAccountID").Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.IDf("example%sID", sn).Op(":=").ID("fakes").Dot("BuildFakeID").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("c"), jen.ID("db")).Op(":=").ID("buildTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("args").Op(":=").Index().Interface().Valuesln(
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
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
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
