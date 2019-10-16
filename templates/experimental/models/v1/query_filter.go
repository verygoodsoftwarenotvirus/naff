package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func queryFilterDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("MaxLimit").Op("=").Lit(250).Var().ID("DefaultLimit").Op("=").Lit(20).Var().ID("pageKey").Op("=").Lit("page").Var().ID("limitKey").Op("=").Lit("limit").Var().ID("createdBeforeKey").Op("=").Lit("created_before").Var().ID("createdAfterKey").Op("=").Lit("created_after").Var().ID("updatedBeforeKey").Op("=").Lit("updated_before").Var().ID("updatedAfterKey").Op("=").Lit("updated_after").Var().ID("sortByKey").Op("=").Lit("sort_by"),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("QueryFilter").Struct(jen.ID("Page").ID("uint64"), jen.ID("Limit").ID("uint64"), jen.ID("CreatedAfter").ID("uint64"), jen.ID("CreatedBefore").ID("uint64"), jen.ID("UpdatedAfter").ID("uint64"), jen.ID("UpdatedBefore").ID("uint64"), jen.ID("SortBy").ID("sortType")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("DefaultQueryFilter builds the default query filter"),
		jen.Line(),
		jen.Func().ID("DefaultQueryFilter").Params().Params(jen.Op("*").ID("QueryFilter")).Block(
			jen.Return().Op("&").ID("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").ID("DefaultLimit"), jen.ID("SortBy").Op(":").ID("SortAscending")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("FromParams overrides the core QueryFilter values with values retrieved from url.Params"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("FromParams").Params(jen.ID("params").Qual("net/url", "Values")).Block(
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("params").Dot(
				"Get",
			).Call(jen.ID("pageKey")), jen.Lit(10), jen.Lit(64)), jen.ID("err").Op("==").ID("nil")).Block(
				jen.ID("qf").Dot(
					"Page",
				).Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(1))),
			),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("params").Dot(
				"Get",
			).Call(jen.ID("limitKey")), jen.Lit(10), jen.Lit(64)), jen.ID("err").Op("==").ID("nil")).Block(
				jen.ID("qf").Dot(
					"Limit",
				).Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0)), jen.ID("MaxLimit"))),
			),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("params").Dot(
				"Get",
			).Call(jen.ID("createdBeforeKey")), jen.Lit(10), jen.Lit(64)), jen.ID("err").Op("==").ID("nil")).Block(
				jen.ID("qf").Dot(
					"CreatedBefore",
				).Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("params").Dot(
				"Get",
			).Call(jen.ID("createdAfterKey")), jen.Lit(10), jen.Lit(64)), jen.ID("err").Op("==").ID("nil")).Block(
				jen.ID("qf").Dot(
					"CreatedAfter",
				).Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("params").Dot(
				"Get",
			).Call(jen.ID("updatedBeforeKey")), jen.Lit(10), jen.Lit(64)), jen.ID("err").Op("==").ID("nil")).Block(
				jen.ID("qf").Dot(
					"UpdatedBefore",
				).Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.ID("params").Dot(
				"Get",
			).Call(jen.ID("updatedAfterKey")), jen.Lit(10), jen.Lit(64)), jen.ID("err").Op("==").ID("nil")).Block(
				jen.ID("qf").Dot(
					"UpdatedAfter",
				).Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(jen.ID("float64").Call(jen.ID("i")), jen.Lit(0))),
			),
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.ID("params").Dot(
				"Get",
			).Call(jen.ID("sortByKey")))).Block(
				jen.Case(jen.ID("string").Call(jen.ID("SortAscending"))).Block(jen.ID("qf").Dot(
					"SortBy",
				).Op("=").ID("SortAscending")),
				jen.Case(jen.ID("string").Call(jen.ID("SortDescending"))).Block(jen.ID("qf").Dot(
					"SortBy",
				).Op("=").ID("SortDescending")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("SetPage sets the current page with certain constraints"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("SetPage").Params(jen.ID("page").ID("uint64")).Block(
			jen.ID("qf").Dot(
				"Page",
			).Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(jen.Lit(1), jen.ID("float64").Call(jen.ID("page")))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("QueryPage calculates a query page from the current filter values"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("QueryPage").Params().Params(jen.ID("uint64")).Block(
			jen.Return().ID("qf").Dot(
				"Limit",
			).Op("*").Parens(jen.ID("qf").Dot(
				"Page",
			).Op("-").Lit(1)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ToValues returns a url.Values from a QueryFilter"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("ToValues").Params().Params(jen.Qual("net/url", "Values")).Block(
			jen.If(jen.ID("qf").Op("==").ID("nil")).Block(
				jen.Return().ID("DefaultQueryFilter").Call().Dot(
					"ToValues",
				).Call(),
			),
			jen.ID("v").Op(":=").Qual("net/url", "Values").Valuesln(),
			jen.If(jen.ID("qf").Dot(
				"Page",
			).Op("!=").Lit(0)).Block(
				jen.ID("v").Dot(
					"Set",
				).Call(jen.Lit("page"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot(
					"Page",
				), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot(
				"Limit",
			).Op("!=").Lit(0)).Block(
				jen.ID("v").Dot(
					"Set",
				).Call(jen.Lit("limit"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot(
					"Limit",
				), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot(
				"SortBy",
			).Op("!=").Lit("")).Block(
				jen.ID("v").Dot(
					"Set",
				).Call(jen.Lit("sort_by"), jen.ID("string").Call(jen.ID("qf").Dot(
					"SortBy",
				))),
			),
			jen.If(jen.ID("qf").Dot(
				"CreatedBefore",
			).Op("!=").Lit(0)).Block(
				jen.ID("v").Dot(
					"Set",
				).Call(jen.Lit("created_before"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot(
					"CreatedBefore",
				), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot(
				"CreatedAfter",
			).Op("!=").Lit(0)).Block(
				jen.ID("v").Dot(
					"Set",
				).Call(jen.Lit("created_after"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot(
					"CreatedAfter",
				), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot(
				"UpdatedBefore",
			).Op("!=").Lit(0)).Block(
				jen.ID("v").Dot(
					"Set",
				).Call(jen.Lit("updated_before"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot(
					"UpdatedBefore",
				), jen.Lit(10))),
			),
			jen.If(jen.ID("qf").Dot(
				"UpdatedAfter",
			).Op("!=").Lit(0)).Block(
				jen.ID("v").Dot(
					"Set",
				).Call(jen.Lit("updated_after"), jen.Qual("strconv", "FormatUint").Call(jen.ID("qf").Dot(
					"UpdatedAfter",
				), jen.Lit(10))),
			),
			jen.Return().ID("v"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ApplyToQueryBuilder applies the query filter to a query builder"),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("ApplyToQueryBuilder").Params(jen.ID("queryBuilder").ID("squirrel").Dot(
			"SelectBuilder",
		)).Params(jen.ID("squirrel").Dot(
			"SelectBuilder",
		)).Block(
			jen.If(jen.ID("qf").Op("==").ID("nil")).Block(
				jen.Return().ID("queryBuilder"),
			),
			jen.ID("qf").Dot(
				"SetPage",
			).Call(jen.ID("qf").Dot(
				"Page",
			)),
			jen.If(jen.ID("qp").Op(":=").ID("qf").Dot(
				"QueryPage",
			).Call(), jen.ID("qp").Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot(
					"Offset",
				).Call(jen.ID("qp")),
			),
			jen.If(jen.ID("qf").Dot(
				"Limit",
			).Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot(
					"Limit",
				).Call(jen.ID("qf").Dot(
					"Limit",
				)),
			).Else().Block(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot(
					"Limit",
				).Call(jen.ID("MaxLimit")),
			),
			jen.If(jen.ID("qf").Dot(
				"CreatedAfter",
			).Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot(
					"Where",
				).Call(jen.ID("squirrel").Dot(
					"Gt",
				).Valuesln(jen.Lit("created_on").Op(":").ID("qf").Dot(
					"CreatedAfter",
				))),
			),
			jen.If(jen.ID("qf").Dot(
				"CreatedBefore",
			).Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot(
					"Where",
				).Call(jen.ID("squirrel").Dot(
					"Lt",
				).Valuesln(jen.Lit("created_on").Op(":").ID("qf").Dot(
					"CreatedBefore",
				))),
			),
			jen.If(jen.ID("qf").Dot(
				"UpdatedAfter",
			).Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot(
					"Where",
				).Call(jen.ID("squirrel").Dot(
					"Gt",
				).Valuesln(jen.Lit("updated_on").Op(":").ID("qf").Dot(
					"UpdatedAfter",
				))),
			),
			jen.If(jen.ID("qf").Dot(
				"UpdatedBefore",
			).Op(">").Lit(0)).Block(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot(
					"Where",
				).Call(jen.ID("squirrel").Dot(
					"Lt",
				).Valuesln(jen.Lit("updated_on").Op(":").ID("qf").Dot(
					"UpdatedBefore",
				))),
			),
			jen.Return().ID("queryBuilder"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ExtractQueryFilter can extract a QueryFilter from a request"),
		jen.Line(),
		jen.Func().ID("ExtractQueryFilter").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("QueryFilter")).Block(
			jen.ID("qf").Op(":=").Op("&").ID("QueryFilter").Valuesln(),
			jen.ID("qf").Dot(
				"FromParams",
			).Call(jen.ID("req").Dot(
				"URL",
			).Dot(
				"Query",
			).Call()),
			jen.Return().ID("qf"),
		),
		jen.Line(),
	)
	return ret
}
