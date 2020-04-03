package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("MaxLimit is the maximum value for list queries"),
			jen.ID("MaxLimit").Equals().Lit(250),
			jen.Comment("DefaultLimit represents how many results we return in a response by default"),
			jen.ID("DefaultLimit").Equals().Lit(20),
			jen.Line(),
			jen.ID("pageKey").Equals().Lit("page"),
			jen.ID("limitKey").Equals().Lit("limit"),
			jen.ID("createdBeforeKey").Equals().Lit("created_before"),
			jen.ID("createdAfterKey").Equals().Lit("created_after"),
			jen.ID("updatedBeforeKey").Equals().Lit("updated_before"),
			jen.ID("updatedAfterKey").Equals().Lit("updated_after"),
			jen.ID("sortByKey").Equals().Lit("sort_by"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("QueryFilter represents all the filters a user could apply to a list query"),
		jen.Line(),
		jen.Type().ID("QueryFilter").Struct(
			jen.ID("Page").Uint64().Tag(jsonTag("page")),
			jen.ID("Limit").Uint64().Tag(jsonTag("limit")),
			jen.ID("CreatedAfter").Uint64().Tag(jsonTag("created_before,omitempty")),
			jen.ID("CreatedBefore").Uint64().Tag(jsonTag("created_after,omitempty")),
			jen.ID("UpdatedAfter").Uint64().Tag(jsonTag("updated_before,omitempty")),
			jen.ID("UpdatedBefore").Uint64().Tag(jsonTag("updated_after,omitempty")),
			jen.ID("SortBy").ID("sortType").Tag(jsonTag("sort_by")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("DefaultQueryFilter builds the default query filter"),
		jen.Line(),
		jen.Func().ID("DefaultQueryFilter").Params().Params(jen.PointerTo().ID("QueryFilter")).Block(
			jen.Return().VarPointer().ID("QueryFilter").Valuesln(
				jen.ID("Page").MapAssign().Lit(1),
				jen.ID("Limit").MapAssign().ID("DefaultLimit"),
				jen.ID("SortBy").MapAssign().ID("SortAscending"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("FromParams overrides the core QueryFilter values with values retrieved from url.Params"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("FromParams").Params(jen.ID("params").Qual("net/url", "Values")).Block(
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("pageKey")), jen.Lit(10), jen.Lit(64)), jen.Err().Op("==").ID("nil")).Block(
				jen.ID("qf").Dot("Page").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(1))),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("limitKey")), jen.Lit(10), jen.Lit(64)), jen.Err().Op("==").ID("nil")).Block(
				jen.ID("qf").Dot("Limit").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0)), jen.ID("MaxLimit"))),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("createdBeforeKey")), jen.Lit(10), jen.Lit(64)), jen.Err().Op("==").ID("nil")).Block(
				jen.ID("qf").Dot("CreatedBefore").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("createdAfterKey")), jen.Lit(10), jen.Lit(64)), jen.Err().Op("==").ID("nil")).Block(
				jen.ID("qf").Dot("CreatedAfter").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("updatedBeforeKey")), jen.Lit(10), jen.Lit(64)), jen.Err().Op("==").ID("nil")).Block(
				jen.ID("qf").Dot("UpdatedBefore").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("i"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.ID("params").Dot("Get").Call(jen.ID("updatedAfterKey")), jen.Lit(10), jen.Lit(64)), jen.Err().Op("==").ID("nil")).Block(
				jen.ID("qf").Dot("UpdatedAfter").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.Line(),
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.ID("params").Dot("Get").Call(jen.ID("sortByKey")))).Block(
				jen.Case(jen.String().Call(jen.ID("SortAscending"))).Block(
					jen.ID("qf").Dot("SortBy").Equals().ID("SortAscending"),
				),
				jen.Case(jen.String().Call(jen.ID("SortDescending"))).Block(
					jen.ID("qf").Dot("SortBy").Equals().ID("SortDescending"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("SetPage sets the current page with certain constraints"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("SetPage").Params(jen.ID("page").Uint64()).Block(
			jen.ID("qf").Dot("Page").Equals().Uint64().Call(jen.Qual("math", "Max").Call(jen.Lit(1), jen.ID("float64").Call(jen.ID("page")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("QueryPage calculates a query page from the current filter values"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("QueryPage").Params().Params(jen.Uint64()).Block(
			jen.Return().ID("qf").Dot("Limit").Times().Parens(jen.ID("qf").Dot("Page").Op("-").Lit(1)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ToValues returns a url.Values from a QueryFilter"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("ToValues").Params().Params(jen.Qual("net/url", "Values")).Block(
			jen.If(jen.ID("qf").Op("==").ID("nil")).Block(
				jen.Return().ID("DefaultQueryFilter").Call().Dot("ToValues").Call(),
			),
			jen.Line(),
			jen.ID("v").Assign().Qual("net/url", "Values").Values(),
			jen.If(jen.ID("qf").Dot("Page").DoesNotEqual().Lit(0)).Block(
				jen.ID("v").Dot("Set").Call(jen.Lit("page"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("Page"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("Limit").DoesNotEqual().Lit(0)).Block(
				jen.ID("v").Dot("Set").Call(jen.Lit("limit"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("Limit"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("SortBy").DoesNotEqual().Lit("")).Block(
				jen.ID("v").Dot("Set").Call(jen.Lit("sort_by"), jen.String().Call(jen.ID("qf").Dot("SortBy"))),
			),
			jen.If(jen.ID("qf").Dot("CreatedBefore").DoesNotEqual().Lit(0)).Block(
				jen.ID("v").Dot("Set").Call(jen.Lit("created_before"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("CreatedBefore"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("CreatedAfter").DoesNotEqual().Lit(0)).Block(
				jen.ID("v").Dot("Set").Call(jen.Lit("created_after"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("CreatedAfter"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").DoesNotEqual().Lit(0)).Block(
				jen.ID("v").Dot("Set").Call(jen.Lit("updated_before"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("UpdatedBefore"), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").DoesNotEqual().Lit(0)).Block(
				jen.ID("v").Dot("Set").Call(jen.Lit("updated_after"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot("UpdatedAfter"), jen.Lit(10))),
			),
			jen.Line(),
			jen.Return().ID("v"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ApplyToQueryBuilder applies the query filter to a query builder"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").PointerTo().ID("QueryFilter")).ID("ApplyToQueryBuilder").Params(jen.ID("queryBuilder").Qual("github.com/Masterminds/squirrel", "SelectBuilder")).Params(jen.Qual("github.com/Masterminds/squirrel", "SelectBuilder")).Block(
			jen.If(jen.ID("qf").Op("==").ID("nil")).Block(
				jen.Return().ID("queryBuilder"),
			),
			jen.Line(),
			jen.ID("qf").Dot("SetPage").Call(jen.ID("qf").Dot("Page")),
			jen.If(jen.ID("qp").Assign().ID("qf").Dot("QueryPage").Call(), jen.ID("qp").Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Offset").Call(jen.ID("qp")),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("Limit").Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Limit").Call(jen.ID("qf").Dot("Limit")),
			).Else().Block(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Limit").Call(jen.ID("MaxLimit")),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("CreatedAfter").Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Gt").Values(jen.Lit("created_on").MapAssign().ID("qf").Dot("CreatedAfter"))),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("CreatedBefore").Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Lt").Values(jen.Lit("created_on").MapAssign().ID("qf").Dot("CreatedBefore"))),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Gt").Values(jen.Lit("updated_on").MapAssign().ID("qf").Dot("UpdatedAfter"))),
			),
			jen.Line(),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Equals().ID("queryBuilder").Dot("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Lt").Values(jen.Lit("updated_on").MapAssign().ID("qf").Dot("UpdatedBefore"))),
			),
			jen.Line(),
			jen.Return().ID("queryBuilder"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ExtractQueryFilter can extract a QueryFilter from a request"),
		jen.Line(),
		jen.Func().ID("ExtractQueryFilter").Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.PointerTo().ID("QueryFilter")).Block(
			jen.ID("qf").Assign().VarPointer().ID("QueryFilter").Values(),
			jen.ID("qf").Dot("FromParams").Call(jen.ID("req").Dot("URL").Dot("Query").Call()),
			jen.Return().ID("qf"),
		),
		jen.Line(),
	)
	return ret
}
