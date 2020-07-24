package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bleveTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Type().Defs(
			jen.ID("exampleType").Struct(
				jen.ID("ID").Uint64().Tag(map[string]string{"json": "id"}),
				jen.ID("Name").String().Tag(map[string]string{"json": "name"}),
				jen.ID("BelongsToUser").Uint64().Tag(map[string]string{"json": "belongsToUser"}),
			),
			jen.Line(),
			jen.ID("exampleTypeWithStringID").Struct(
				jen.ID("ID").String().Tag(map[string]string{"json": "id"}),
				jen.ID("Name").String().Tag(map[string]string{"json": "name"}),
				jen.ID("BelongsToUser").Uint64().Tag(map[string]string{"json": "belongsToUser"}),
			),
		),
	)

	code.Add(buildTestNewBleveIndexManager(proj)...)
	code.Add(buildTestBleveIndexManager_Index(proj)...)
	code.Add(buildTestBleveIndexManager_Search(proj)...)
	code.Add(buildTestBleveIndexManager_Delete(proj)...)

	return code
}

func buildTestNewBleveIndexManager(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestNewBleveIndexManager").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("happy path",
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit("constructor_test.bleve"),
				),
				jen.Line(),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "RemoveAll").Call(jen.String().Call(jen.ID("exampleIndexPath"))), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext("invalid path",
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit(""),
				),
				jen.Line(),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext("invalid name",
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit(""),
				),
				jen.Line(),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.Lit("testprojects"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertError(jen.Err(), nil),
			),
		),
	}

	return lines
}

func buildTestBleveIndexManager_Index(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestBleveIndexManager_Index").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("exampleUserID").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call().Dot("ID"),
			jen.Line(),
			utils.BuildSubTest("obligatory",
				jen.Const().ID("exampleQuery").Equals().Lit("index_test"),
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit("index_test.bleve"),
				),
				jen.Line(),
				jen.List(jen.ID("im"), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("im"), nil),
				jen.Line(),
				jen.ID("x").Assign().AddressOf().ID("exampleType").Valuesln(
					jen.ID("ID").MapAssign().Lit(123),
					jen.ID("Name").MapAssign().ID("exampleQuery"),
					jen.ID("BelongsToUser").MapAssign().ID("exampleUserID"),
				),
				utils.AssertNoError(
					jen.ID("im").Dot("Index").Call(
						constants.CtxVar(),
						jen.ID("x").Dot("ID"),
						jen.ID("x"),
					),
					nil,
				),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "RemoveAll").Call(jen.String().Call(jen.ID("exampleIndexPath"))), nil),
			),
		),
	}

	return lines
}

func buildTestBleveIndexManager_Search(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestBleveIndexManager_Search").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("exampleUserID").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call().Dot("ID"),
			jen.Line(),
			utils.BuildSubTest("obligatory",
				jen.Const().ID("exampleQuery").Equals().Lit("search_test"),
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit("search_test_1.bleve"),
				),
				jen.List(jen.ID("im"), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("im"), nil),
				jen.Line(),
				jen.ID("x").Assign().AddressOf().ID("exampleType").Valuesln(
					jen.ID("ID").MapAssign().Lit(123),
					jen.ID("Name").MapAssign().ID("exampleQuery"),
					jen.ID("BelongsToUser").MapAssign().ID("exampleUserID"),
				),
				utils.AssertNoError(
					jen.ID("im").Dot("Index").Call(
						constants.CtxVar(),
						jen.ID("x").Dot("ID"),
						jen.ID("x"),
					),
					nil,
				),
				jen.Line(),
				jen.List(jen.ID("results"), jen.Err()).Assign().ID("im").Dot("Search").Call(
					constants.CtxVar(),
					jen.ID("x").Dot("Name"),
					jen.ID("exampleUserID"),
				),
				utils.AssertNotEmpty(jen.ID("results"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "RemoveAll").Call(jen.String().Call(jen.ID("exampleIndexPath"))), nil),
			),
			jen.Line(),
			utils.BuildSubTest("with empty index and search",
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit("search_test_2.bleve"),
				),
				jen.Line(),
				jen.List(jen.ID("im"), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("im"), nil),
				jen.Line(),
				jen.List(jen.ID("results"), jen.Err()).Assign().ID("im").Dot("Search").Call(
					constants.CtxVar(),
					jen.EmptyString(),
					jen.ID("exampleUserID"),
				),
				utils.AssertEmpty(jen.ID("results"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "RemoveAll").Call(jen.String().Call(jen.ID("exampleIndexPath"))), nil),
			),
			jen.Line(),
			utils.BuildSubTest("with closed index",
				jen.Const().ID("exampleQuery").Equals().Lit("search_test"),
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit("search_test_3.bleve"),
				),
				jen.List(jen.ID("im"), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("im"), nil),
				jen.Line(),
				jen.ID("x").Assign().AddressOf().ID("exampleType").Valuesln(
					jen.ID("ID").MapAssign().Lit(123),
					jen.ID("Name").MapAssign().ID("exampleQuery"),
					jen.ID("BelongsToUser").MapAssign().ID("exampleUserID"),
				),
				utils.AssertNoError(
					jen.ID("im").Dot("Index").Call(
						constants.CtxVar(),
						jen.ID("x").Dot("ID"),
						jen.ID("x"),
					),
					nil,
				),
				jen.Line(),
				utils.AssertNoError(
					jen.ID("im").Dot("").Parens(jen.PointerTo().ID("bleveIndexManager")).Dot("index").Dot("Close").Call(),
					nil,
				),
				jen.Line(),
				jen.List(jen.ID("results"), jen.Err()).Assign().ID("im").Dot("Search").Call(
					constants.CtxVar(),
					jen.ID("x").Dot("Name"),
					jen.ID("exampleUserID"),
				),
				utils.AssertEmpty(jen.ID("results"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "RemoveAll").Call(jen.String().Call(jen.ID("exampleIndexPath"))), nil),
			),
			jen.Line(),
			utils.BuildSubTest("with invalid ID",
				jen.Const().ID("exampleQuery").Equals().Lit("search_test"),
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit("search_test_4.bleve"),
				),
				jen.List(jen.ID("im"), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("im"), nil),
				jen.Line(),
				jen.ID("x").Assign().AddressOf().ID("exampleTypeWithStringID").Valuesln(
					jen.ID("ID").MapAssign().Lit("whatever"),
					jen.ID("Name").MapAssign().ID("exampleQuery"),
					jen.ID("BelongsToUser").MapAssign().ID("exampleUserID"),
				),
				utils.AssertNoError(
					jen.ID("im").Dot("").Parens(jen.PointerTo().ID("bleveIndexManager")).Dot("index").Dot("Index").Call(
						jen.ID("x").Dot("ID"),
						jen.ID("x"),
					),
					nil,
				),
				jen.Line(),
				jen.List(jen.ID("results"), jen.Err()).Assign().ID("im").Dot("Search").Call(
					constants.CtxVar(),
					jen.ID("x").Dot("Name"),
					jen.ID("exampleUserID"),
				),
				utils.AssertEmpty(jen.ID("results"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "RemoveAll").Call(jen.String().Call(jen.ID("exampleIndexPath"))), nil),
			),
		),
	}

	return lines
}

func buildTestBleveIndexManager_Delete(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestBleveIndexManager_Delete").Params(jen.ID("T").PointerTo().Qual("testprojects", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("exampleUserID").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call().Dot("ID"),
			jen.Line(),
			utils.BuildSubTest("obligatory",
				jen.Const().ID("exampleQuery").Equals().Lit("delete_test"),
				jen.ID("exampleIndexPath").Assign().Qual(proj.InternalSearchV1Package(), "IndexPath").Call(
					jen.Lit("delete_test.bleve"),
				),
				jen.Line(),
				jen.List(jen.ID("im"), jen.Err()).Assign().ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.ID("testingSearchIndexName"),
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID("im"), nil),
				jen.Line(),
				jen.ID("x").Assign().AddressOf().ID("exampleType").Valuesln(
					jen.ID("ID").MapAssign().Lit(123),
					jen.ID("Name").MapAssign().ID("exampleQuery"),
					jen.ID("BelongsToUser").MapAssign().ID("exampleUserID"),
				),
				utils.AssertNoError(
					jen.ID("im").Dot("Index").Call(
						constants.CtxVar(),
						jen.ID("x").Dot("ID"),
						jen.ID("x"),
					),
					nil,
				),
				jen.Line(),
				utils.AssertNoError(
					jen.ID("im").Dot("Delete").Call(
						constants.CtxVar(),
						jen.ID("x").Dot("ID"),
					),
					nil,
				),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "RemoveAll").Call(jen.String().Call(jen.ID("exampleIndexPath"))), nil),
			),
		),
	}

	return lines
}
