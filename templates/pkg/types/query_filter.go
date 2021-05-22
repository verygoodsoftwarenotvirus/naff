package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildQueryFilterConstantDeclarations0()...)
	code.Add(buildQueryFilter()...)
	code.Add(buildDefaultQueryFilter()...)
	code.Add(buildFromParams()...)
	code.Add(buildSetPage()...)
	code.Add(buildQueryPage()...)
	code.Add(buildToValues()...)
	code.Add(buildApplyToQueryBuilder()...)
	code.Add(buildExtractQueryFilter()...)

	return code
}

func buildQueryFilterConstantDeclarations0() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("MaxLimit is the maximum value for list queries."),
			jen.ID("MaxLimit").Equals().Lit(250),
			jen.Comment("DefaultLimit represents how many results we return in a response by default."),
			jen.ID("DefaultLimit").Equals().Lit(20),
			jen.Line(),
			jen.Comment("SearchQueryKey is the query param key we use to find search queries in requests"),
			jen.ID("SearchQueryKey").Equals().Lit("q"),
			jen.Comment("LimitQueryKey is the query param key we use to specify a limit in a query"),
			jen.ID("LimitQueryKey").Equals().Lit("limit"),
			jen.Line(),
			jen.ID("pageQueryKey").Equals().Lit("page"),
			jen.ID("createdBeforeQueryKey").Equals().Lit("createdBefore"),
			jen.ID("createdAfterQueryKey").Equals().Lit("createdAfter"),
			jen.ID("updatedBeforeQueryKey").Equals().Lit("updatedBefore"),
			jen.ID("updatedAfterQueryKey").Equals().Lit("updatedAfter"),
			jen.ID("sortByQueryKey").Equals().Lit("sortBy"),
		),
		jen.Line(),
	}

	return lines
}

func buildQueryFilter() []jen.Code {
	lines := []jen.Code{
		jen.Comment("QueryFilter represents all the filters a user could apply to a list query."),
		jen.Line(),
		jen.Type().ID("QueryFilter").Struct(
			jen.ID("Page").Uint64().Tag(jsonTag("page")),
			jen.ID("Limit").Uint8().Tag(jsonTag("limit")),
			jen.ID("CreatedAfter").Uint64().Tag(jsonTag("createdBefore,omitempty")),
			jen.ID("CreatedBefore").Uint64().Tag(jsonTag("createdAfter,omitempty")),
			jen.ID("UpdatedAfter").Uint64().Tag(jsonTag("updatedBefore,omitempty")),
			jen.ID("UpdatedBefore").Uint64().Tag(jsonTag("updatedAfter,omitempty")),
			jen.ID("SortBy").ID("sortType").Tag(jsonTag("sortBy")),
		),
		jen.Line(),
	}

	return lines
}

func buildDefaultQueryFilter() []jen.Code {
	lines := []jen.Code{
		jen.Comment("DefaultQueryFilter builds the default query filter."),
		jen.Line(),
		jen.Func().ID("DefaultQueryFilter").Params().Params(jen.PointerTo().ID("QueryFilter")).Body(
			jen.Return().AddressOf().ID("QueryFilter").Valuesln(
				jen.ID("Page").MapAssign().One(),
				jen.ID("Limit").MapAssign().ID("DefaultLimit"),
				jen.ID("SortBy").MapAssign().ID("SortAscending"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildFromParams() []jen.Code {
	lines := []jen.Code{
		jen.Comment("FromParams overrides the core QueryFilter values with values retrieved from url.Params"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("FromParams").Params(jen.ID("params").Qual("net/url", "Values")).Body(
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("pageQueryKey")), jen.Lit(10), jen.Lit(64)), jen.Err().IsEqualTo().ID("nil")).Body(
				jen.ID("qf").Dot("Page").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.One())),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("LimitQueryKey")), jen.Lit(10), jen.Lit(64)), jen.Err().IsEqualTo().ID("nil")).Body(
				jen.ID("qf").Dot("Limit").Equals().Uint8().Call(jen.Qual("math", "Min").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Zero()), jen.ID("MaxLimit"))),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("createdBeforeQueryKey")), jen.Lit(10), jen.Lit(64)), jen.Err().IsEqualTo().ID("nil")).Body(
				jen.ID("qf").Dot("CreatedBefore").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Zero())),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("createdAfterQueryKey")), jen.Lit(10), jen.Lit(64)), jen.Err().IsEqualTo().ID("nil")).Body(
				jen.ID("qf").Dot("CreatedAfter").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Zero())),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("updatedBeforeQueryKey")), jen.Lit(10), jen.Lit(64)), jen.Err().IsEqualTo().ID("nil")).Body(
				jen.ID("qf").Dot("UpdatedBefore").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Zero())),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("updatedAfterQueryKey")), jen.Lit(10), jen.Lit(64)), jen.Err().IsEqualTo().ID("nil")).Body(
				jen.ID("qf").Dot("UpdatedAfter").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Zero())),
			),
			jen.Line(),
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.ID("params").Dot("Get").Call(jen.ID("sortByQueryKey")))).Body(
				jen.Case(jen.String().Call(jen.ID("SortAscending"))).Body(
					jen.ID("qf").Dot("SortBy").Equals().ID("SortAscending"),
				),
				jen.Case(jen.String().Call(jen.ID("SortDescending"))).Body(
					jen.ID("qf").Dot("SortBy").Equals().ID("SortDescending"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildSetPage() []jen.Code {
	lines := []jen.Code{
		jen.Comment("SetPage sets the current page with certain constraints."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("SetPage").Params(jen.ID("page").Uint64()).Body(
			jen.ID("qf").Dot("Page").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.One(), jen.ID("float64").Call(jen.ID("page")))),
		),
		jen.Line(),
	}

	return lines
}

func buildQueryPage() []jen.Code {
	lines := []jen.Code{
		jen.Comment("QueryPage calculates a query page from the current filter values."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("QueryPage").Params().Params(jen.Uint64()).Body(
			jen.Return().Uint64().Call(jen.ID("qf").Dot("Limit")).Times().Parens(jen.ID("qf").Dot("Page").Minus().One()),
		),
		jen.Line(),
	}

	return lines
}

func buildToValues() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ToValues returns a url.Values from a QueryFilter"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("ToValues").Params().Params(jen.Qual("net/url", "Values")).Body(
			jen.If(jen.ID("qf").IsEqualTo().ID("nil")).Body(
				jen.Return().ID("DefaultQueryFilter").Call().Dot("ToValues").Call(),
			),
			jen.Line(),
			jen.ID("v").Assign().Qual("net/url", "Values").Values(),
			jen.If(jen.ID("qf").Dot("Page").DoesNotEqual().Zero()).Body(
				jen.ID("v").Dot("Set").Call(jen.ID("pageQueryKey"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("Page"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("Limit").DoesNotEqual().Zero()).Body(
				jen.ID("v").Dot("Set").Call(jen.ID("LimitQueryKey"), jen.Qual("strconv", "FormatUint").Call(jen.Uint64().Call(jen.ID("qf").Dot("Limit")), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("SortBy").DoesNotEqual().EmptyString()).Body(
				jen.ID("v").Dot("Set").Call(jen.ID("sortByQueryKey"), jen.String().Call(jen.ID("qf").Dot("SortBy"))),
			),
			jen.If(jen.ID("qf").Dot("CreatedBefore").DoesNotEqual().Zero()).Body(
				jen.ID("v").Dot("Set").Call(jen.ID("createdBeforeQueryKey"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("CreatedBefore"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("CreatedAfter").DoesNotEqual().Zero()).Body(
				jen.ID("v").Dot("Set").Call(jen.ID("createdAfterQueryKey"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("CreatedAfter"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").DoesNotEqual().Zero()).Body(
				jen.ID("v").Dot("Set").Call(jen.ID("updatedBeforeQueryKey"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("UpdatedBefore"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").DoesNotEqual().Zero()).Body(
				jen.ID("v").Dot("Set").Call(jen.ID("updatedAfterQueryKey"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("UpdatedAfter"), jen.Lit(10))),
			),
			jen.Line(),
			jen.Return().ID("v"),
		),
		jen.Line(),
	}

	return lines
}

func buildApplyToQueryBuilder() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ApplyToQueryBuilder applies the query filter to a query builder."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("ApplyToQueryBuilder").Params(
			jen.ID("queryBuilder").Qual("github.com/Masterminds/squirrel", "SelectBuilder"),
			jen.ID("tableName").String(),
		).Params(
			jen.Qual("github.com/Masterminds/squirrel", "SelectBuilder"),
		).Body(
			jen.If(jen.ID("qf").IsEqualTo().ID("nil")).Body(
				jen.Return().ID("queryBuilder"),
			),
			jen.Line(),
			jen.Const().Defs(
				jen.ID("createdOnKey").Equals().Lit("created_on"),
				jen.ID("updatedOnKey").Equals().Lit("last_updated_on"),
			),
			jen.Line(),
			jen.ID("qf").Dot("SetPage").Call(jen.ID("qf").Dot("Page")),
			jen.If(jen.ID("qp").Assign().ID("qf").Dot("QueryPage").Call(), jen.ID("qp").GreaterThan().Zero()).Body(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Offset").Call(jen.ID("qp")),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("Limit").GreaterThan().Zero()).Body(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Limit").Call(jen.Uint64().Call(jen.ID("qf").Dot("Limit"))),
			).Else().Body(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Limit").Call(jen.ID("MaxLimit")),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("CreatedAfter").GreaterThan().Zero()).Body(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Gt").Values(
						utils.FormatString("%s.%s", jen.ID("tableName"), jen.ID("createdOnKey")).MapAssign().ID("qf").Dot("CreatedAfter"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("CreatedBefore").GreaterThan().Zero()).Body(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Lt").Values(
						utils.FormatString("%s.%s", jen.ID("tableName"), jen.ID("createdOnKey")).MapAssign().ID("qf").Dot("CreatedBefore"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").GreaterThan().Zero()).Body(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Gt").Values(
						utils.FormatString("%s.%s", jen.ID("tableName"), jen.ID("updatedOnKey")).MapAssign().ID("qf").Dot("UpdatedAfter"),
					),
				),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").GreaterThan().Zero()).Body(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Lt").Values(
						utils.FormatString("%s.%s", jen.ID("tableName"), jen.ID("updatedOnKey")).MapAssign().ID("qf").Dot("UpdatedBefore"),
					),
				),
			),
			jen.Line(),
			jen.Return().ID("queryBuilder"),
		),
		jen.Line(),
	}

	return lines
}

func buildExtractQueryFilter() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ExtractQueryFilter can extract a QueryFilter from a request."),
		jen.Line(),
		jen.Func().ID("ExtractQueryFilter").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().ID("QueryFilter")).Body(
			jen.ID("qf").Assign().AddressOf().ID("QueryFilter").Values(),
			jen.ID("qf").Dot("FromParams").Call(jen.ID(constants.RequestVarName).Dot("URL").Dot("Query").Call()),
			jen.Return().ID("qf"),
		),
		jen.Line(),
	}

	return lines
}
