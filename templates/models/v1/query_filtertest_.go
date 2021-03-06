package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)

	code.Add(buildTestFromParams()...)
	code.Add(buildTestQueryFilter_SetPage()...)
	code.Add(buildTestQueryFilter_QueryPage()...)
	code.Add(buildTestQueryFilter_ToValues()...)
	code.Add(buildTestQueryFilter_ApplyToQueryBuilder()...)
	code.Add(buildTestExtractQueryFilter()...)

	return code
}

func buildTestFromParams() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestFromParams").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("actual").Assign().AddressOf().ID("QueryFilter").Values(),
				jen.ID("expected").Assign().AddressOf().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().ID("MaxLimit"),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("pageQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("Page")))),
					jen.ID("LimitQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("Limit")))),
					jen.ID("createdBeforeQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("CreatedAfter")))),
					jen.ID("createdAfterQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("CreatedBefore")))),
					jen.ID("updatedBeforeQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("UpdatedAfter")))),
					jen.ID("updatedAfterQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("UpdatedBefore")))),
					jen.ID("sortByQueryKey").MapAssign().Index().String().Values(jen.String().Call(jen.ID("expected").Dot("SortBy"))),
				),
				jen.Line(),
				jen.ID("actual").Dot("FromParams").Call(jen.ID(utils.BuildFakeVarName("Input"))),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Input")).Index(jen.ID("sortByQueryKey")).Equals().Index().String().Values(jen.String().Call(jen.ID("SortAscending"))),
				jen.Line(),
				jen.ID("actual").Dot("FromParams").Call(jen.ID(utils.BuildFakeVarName("Input"))),
				utils.AssertEqual(jen.ID("SortAscending"), jen.ID("actual").Dot("SortBy"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestQueryFilter_SetPage() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestQueryFilter_SetPage").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("qf").Assign().AddressOf().ID("QueryFilter").Values(),
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(123)),
				jen.ID("qf").Dot("SetPage").Call(jen.ID("expected")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("qf").Dot("Page"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestQueryFilter_QueryPage() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestQueryFilter_QueryPage").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("qf").Assign().AddressOf().ID("QueryFilter").Values(jen.ID("Limit").MapAssign().Lit(10), jen.ID("Page").MapAssign().Lit(11)),
				jen.ID("expected").Assign().Uint64().Call(jen.Lit(100)),
				jen.ID("actual").Assign().ID("qf").Dot("QueryPage").Call(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestQueryFilter_ToValues() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestQueryFilter_ToValues").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("qf").Assign().AddressOf().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().Lit(50),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.ID("expected").Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("pageQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("qf").Dot("Page")))),
					jen.ID("LimitQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("qf").Dot("Limit")))),
					jen.ID("createdBeforeQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("qf").Dot("CreatedAfter")))),
					jen.ID("createdAfterQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("qf").Dot("CreatedBefore")))),
					jen.ID("updatedBeforeQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("qf").Dot("UpdatedAfter")))),
					jen.ID("updatedAfterQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("qf").Dot("UpdatedBefore")))),
					jen.ID("sortByQueryKey").MapAssign().Index().String().Values(jen.String().Call(jen.ID("qf").Dot("SortBy"))),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("qf").Dot("ToValues").Call(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil",
				jen.ID("qf").Assign().Parens(jen.PointerTo().ID("QueryFilter")).Call(jen.Nil()),
				jen.ID("expected").Assign().ID("DefaultQueryFilter").Call().Dot("ToValues").Call(),
				jen.ID("actual").Assign().ID("qf").Dot("ToValues").Call(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestQueryFilter_ApplyToQueryBuilder() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestQueryFilter_ApplyToQueryBuilder").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID(utils.BuildFakeVarName("TableName")).Assign().Lit("stuff"),
			jen.ID("baseQueryBuilder").Assign().Qual("github.com/Masterminds/squirrel", "StatementBuilder").Dot("PlaceholderFormat").Call(jen.Qual("github.com/Masterminds/squirrel", "Dollar")).
				Dotln("Select").Call(jen.Lit("things")).
				Dotln("From").Call(jen.ID(utils.BuildFakeVarName("TableName"))).
				Dotln("Where").Call(
				jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(
					utils.FormatString("%s.condition", jen.ID(utils.BuildFakeVarName("TableName"))).MapAssign().True(),
				),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("qf").Assign().AddressOf().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().Lit(50),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.Line(),
				jen.ID("sb").Assign().Qual("github.com/Masterminds/squirrel", "StatementBuilder").Dot("Select").Call(jen.Lit("*")).Dot("From").Call(jen.Lit("testing")),
				jen.ID("qf").Dot("ApplyToQueryBuilder").Call(jen.ID("sb"), jen.ID(utils.BuildFakeVarName("TableName"))),
				jen.ID("expected").Assign().Lit("SELECT * FROM testing"),
				jen.List(jen.ID("actual"), jen.Underscore(), jen.Err()).Assign().ID("sb").Dot("ToSql").Call(),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"basic usecase",
				jen.ID(utils.BuildFakeVarName("QF")).Assign().AddressOf().ID("QueryFilter").Values(jen.ID("Limit").MapAssign().Lit(15), jen.ID("Page").MapAssign().Lit(2)),
				jen.Line(),
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE stuff.condition = $1 LIMIT 15 OFFSET 15`),
				jen.ID("x").Assign().ID("exampleQF").Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder"), jen.ID(utils.BuildFakeVarName("TableName"))),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match"), nil),
				utils.AssertNil(jen.Err(), nil),
				utils.AssertNotEmpty(jen.ID("args"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"returns query builder if query filter is nil",
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE stuff.condition = $1`),
				jen.Line(),
				jen.ID("x").Assign().Parens(jen.PointerTo().ID("QueryFilter")).Call(jen.Nil()).Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder"), jen.ID(utils.BuildFakeVarName("TableName"))),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match"), nil),
				utils.AssertNil(jen.Err(), nil),
				utils.AssertNotEmpty(jen.ID("args"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"whole kit and kaboodle",
				jen.ID(utils.BuildFakeVarName("QF")).Assign().AddressOf().ID("QueryFilter").Valuesln(
					jen.ID("Limit").MapAssign().Lit(20), jen.ID("Page").MapAssign().Lit(6),
					jen.ID("CreatedAfter").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("CreatedBefore").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("UpdatedAfter").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("UpdatedBefore").MapAssign().Uint64().Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE stuff.condition = $1 AND stuff.created_on > $2 AND stuff.created_on < $3 AND stuff.last_updated_on > $4 AND stuff.last_updated_on < $5 LIMIT 20 OFFSET 100`),
				jen.ID("x").Assign().ID("exampleQF").Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder"), jen.ID(utils.BuildFakeVarName("TableName"))),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match"), nil),
				utils.AssertNil(jen.Err(), nil),
				utils.AssertNotEmpty(jen.ID("args"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with zero limit",
				jen.ID(utils.BuildFakeVarName("QF")).Assign().AddressOf().ID("QueryFilter").Values(jen.ID("Limit").MapAssign().Zero(), jen.ID("Page").MapAssign().One()),
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE stuff.condition = $1 LIMIT 250`),
				jen.ID("x").Assign().ID("exampleQF").Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder"), jen.ID(utils.BuildFakeVarName("TableName"))),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match"), nil),
				utils.AssertNil(jen.Err(), nil),
				utils.AssertNotEmpty(jen.ID("args"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestExtractQueryFilter() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestExtractQueryFilter").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("expected").Assign().AddressOf().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().ID("MaxLimit"),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.ID(utils.BuildFakeVarName("Input")).Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("pageQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("Page")))),
					jen.ID("LimitQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("Limit")))),
					jen.ID("createdBeforeQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("CreatedAfter")))),
					jen.ID("createdAfterQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("CreatedBefore")))),
					jen.ID("updatedBeforeQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("UpdatedAfter")))),
					jen.ID("updatedAfterQueryKey").MapAssign().Index().String().Values(jen.Qual("strconv", "Itoa").Call(jen.Int().Call(jen.ID("expected").Dot("UpdatedBefore")))),
					jen.ID("sortByQueryKey").MapAssign().Index().String().Values(jen.String().Call(jen.ID("expected").Dot("SortBy"))),
				),
				jen.Line(),
				jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("https://verygoodsoftwarenotvirus.ru"), jen.Nil()),
				utils.AssertNoError(jen.Err(), nil),
				utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
				jen.Line(),
				jen.ID(constants.RequestVarName).Dot("URL").Dot("RawQuery").Equals().ID("exampleInput").Dot("Encode").Call(),
				jen.ID("actual").Assign().ID("ExtractQueryFilter").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
