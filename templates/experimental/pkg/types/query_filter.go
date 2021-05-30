package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("MaxLimit").Op("=").Lit(250),
			jen.ID("DefaultLimit").Op("=").Lit(20),
			jen.ID("SearchQueryKey").Op("=").Lit("q"),
			jen.ID("LimitQueryKey").Op("=").Lit("limit"),
			jen.ID("AdminQueryKey").Op("=").Lit("admin"),
			jen.ID("pageQueryKey").Op("=").Lit("page"),
			jen.ID("createdBeforeQueryKey").Op("=").Lit("createdBefore"),
			jen.ID("createdAfterQueryKey").Op("=").Lit("createdAfter"),
			jen.ID("updatedBeforeQueryKey").Op("=").Lit("updatedBefore"),
			jen.ID("updatedAfterQueryKey").Op("=").Lit("updatedAfter"),
			jen.ID("includeArchivedQueryKey").Op("=").Lit("includeArchived"),
			jen.ID("sortByQueryKey").Op("=").Lit("sortBy"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("QueryFilter").Struct(
				jen.ID("SortBy").ID("sortType"),
				jen.ID("Page").ID("uint64"),
				jen.ID("CreatedAfter").ID("uint64"),
				jen.ID("CreatedBefore").ID("uint64"),
				jen.ID("UpdatedAfter").ID("uint64"),
				jen.ID("UpdatedBefore").ID("uint64"),
				jen.ID("Limit").ID("uint8"),
				jen.ID("IncludeArchived").ID("bool"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("DefaultQueryFilter builds the default query filter."),
		jen.Line(),
		jen.Func().ID("DefaultQueryFilter").Params().Params(jen.Op("*").ID("QueryFilter")).Body(
			jen.Return().Op("&").ID("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").ID("DefaultLimit"), jen.ID("SortBy").Op(":").ID("SortAscending"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachToLogger attaches a QueryFilter's values to a logging.Logger."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("AttachToLogger").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("logging").Dot("Logger")).Body(
			jen.ID("l").Op(":=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("Clone").Call(),
			jen.If(jen.ID("qf").Op("==").ID("nil")).Body(
				jen.Return().ID("l").Dot("WithValue").Call(
					jen.ID("keys").Dot("FilterIsNilKey"),
					jen.ID("true"),
				)),
			jen.If(jen.ID("qf").Dot("Page").Op("!=").Lit(0)).Body(
				jen.ID("l").Op("=").ID("l").Dot("WithValue").Call(
					jen.ID("pageQueryKey"),
					jen.ID("qf").Dot("Page"),
				)),
			jen.If(jen.ID("qf").Dot("Limit").Op("!=").Lit(0)).Body(
				jen.ID("l").Op("=").ID("l").Dot("WithValue").Call(
					jen.ID("LimitQueryKey"),
					jen.ID("qf").Dot("Limit"),
				)),
			jen.If(jen.ID("qf").Dot("SortBy").Op("!=").Lit("")).Body(
				jen.ID("l").Op("=").ID("l").Dot("WithValue").Call(
					jen.ID("sortByQueryKey"),
					jen.ID("qf").Dot("SortBy"),
				)),
			jen.If(jen.ID("qf").Dot("CreatedBefore").Op("!=").Lit(0)).Body(
				jen.ID("l").Op("=").ID("l").Dot("WithValue").Call(
					jen.ID("createdBeforeQueryKey"),
					jen.ID("qf").Dot("CreatedBefore"),
				)),
			jen.If(jen.ID("qf").Dot("CreatedAfter").Op("!=").Lit(0)).Body(
				jen.ID("l").Op("=").ID("l").Dot("WithValue").Call(
					jen.ID("createdAfterQueryKey"),
					jen.ID("qf").Dot("CreatedAfter"),
				)),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").Op("!=").Lit(0)).Body(
				jen.ID("l").Op("=").ID("l").Dot("WithValue").Call(
					jen.ID("updatedBeforeQueryKey"),
					jen.ID("qf").Dot("UpdatedBefore"),
				)),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").Op("!=").Lit(0)).Body(
				jen.ID("l").Op("=").ID("l").Dot("WithValue").Call(
					jen.ID("updatedAfterQueryKey"),
					jen.ID("qf").Dot("UpdatedAfter"),
				)),
			jen.Return().ID("l"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("FromParams overrides the core QueryFilter values with values retrieved from url.Params."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("FromParams").Params(jen.ID("params").Qual("net/url", "Values")).Body(
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("params").Dot("Get").Call(jen.ID("pageQueryKey")),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("qf").Dot("Page").Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(
					jen.ID("float64").Call(jen.ID("i")),
					jen.Lit(1),
				))),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("params").Dot("Get").Call(jen.ID("LimitQueryKey")),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("qf").Dot("Limit").Op("=").ID("uint8").Call(jen.Qual("math", "Min").Call(
					jen.Qual("math", "Max").Call(
						jen.ID("float64").Call(jen.ID("i")),
						jen.Lit(0),
					),
					jen.ID("MaxLimit"),
				))),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("params").Dot("Get").Call(jen.ID("createdBeforeQueryKey")),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("qf").Dot("CreatedBefore").Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(
					jen.ID("float64").Call(jen.ID("i")),
					jen.Lit(0),
				))),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("params").Dot("Get").Call(jen.ID("createdAfterQueryKey")),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("qf").Dot("CreatedAfter").Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(
					jen.ID("float64").Call(jen.ID("i")),
					jen.Lit(0),
				))),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("params").Dot("Get").Call(jen.ID("updatedBeforeQueryKey")),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("qf").Dot("UpdatedBefore").Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(
					jen.ID("float64").Call(jen.ID("i")),
					jen.Lit(0),
				))),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("params").Dot("Get").Call(jen.ID("updatedAfterQueryKey")),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("qf").Dot("UpdatedAfter").Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(
					jen.ID("float64").Call(jen.ID("i")),
					jen.Lit(0),
				))),
			jen.If(jen.List(jen.ID("i"), jen.ID("err")).Op(":=").Qual("strconv", "ParseBool").Call(jen.ID("params").Dot("Get").Call(jen.ID("includeArchivedQueryKey"))), jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("qf").Dot("IncludeArchived").Op("=").ID("i")),
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.ID("params").Dot("Get").Call(jen.ID("sortByQueryKey")))).Body(
				jen.Case(jen.ID("string").Call(jen.ID("SortAscending"))).Body(
					jen.ID("qf").Dot("SortBy").Op("=").ID("SortAscending")),
				jen.Case(jen.ID("string").Call(jen.ID("SortDescending"))).Body(
					jen.ID("qf").Dot("SortBy").Op("=").ID("SortDescending")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetPage sets the current page with certain constraints."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("SetPage").Params(jen.ID("page").ID("uint64")).Body(
			jen.ID("qf").Dot("Page").Op("=").ID("uint64").Call(jen.Qual("math", "Max").Call(
				jen.Lit(1),
				jen.ID("float64").Call(jen.ID("page")),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("QueryPage calculates a query page from the current filter values."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("QueryPage").Params().Params(jen.ID("uint64")).Body(
			jen.Return().ID("uint64").Call(jen.ID("qf").Dot("Limit")).Op("*").Parens(jen.ID("qf").Dot("Page").Op("-").Lit(1))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ToValues returns a url.Values from a QueryFilter."),
		jen.Line(),
		jen.Func().Params(jen.ID("qf").Op("*").ID("QueryFilter")).ID("ToValues").Params().Params(jen.Qual("net/url", "Values")).Body(
			jen.If(jen.ID("qf").Op("==").ID("nil")).Body(
				jen.Return().ID("DefaultQueryFilter").Call().Dot("ToValues").Call()),
			jen.ID("v").Op(":=").Qual("net/url", "Values").Valuesln(),
			jen.If(jen.ID("qf").Dot("Page").Op("!=").Lit(0)).Body(
				jen.ID("v").Dot("Set").Call(
					jen.ID("pageQueryKey"),
					jen.Qual("strconv", "FormatUint").Call(
						jen.ID("qf").Dot("Page"),
						jen.Lit(10),
					),
				)),
			jen.If(jen.ID("qf").Dot("Limit").Op("!=").Lit(0)).Body(
				jen.ID("v").Dot("Set").Call(
					jen.ID("LimitQueryKey"),
					jen.Qual("strconv", "FormatUint").Call(
						jen.ID("uint64").Call(jen.ID("qf").Dot("Limit")),
						jen.Lit(10),
					),
				)),
			jen.If(jen.ID("qf").Dot("SortBy").Op("!=").Lit("")).Body(
				jen.ID("v").Dot("Set").Call(
					jen.ID("sortByQueryKey"),
					jen.ID("string").Call(jen.ID("qf").Dot("SortBy")),
				)),
			jen.If(jen.ID("qf").Dot("CreatedBefore").Op("!=").Lit(0)).Body(
				jen.ID("v").Dot("Set").Call(
					jen.ID("createdBeforeQueryKey"),
					jen.Qual("strconv", "FormatUint").Call(
						jen.ID("qf").Dot("CreatedBefore"),
						jen.Lit(10),
					),
				)),
			jen.If(jen.ID("qf").Dot("CreatedAfter").Op("!=").Lit(0)).Body(
				jen.ID("v").Dot("Set").Call(
					jen.ID("createdAfterQueryKey"),
					jen.Qual("strconv", "FormatUint").Call(
						jen.ID("qf").Dot("CreatedAfter"),
						jen.Lit(10),
					),
				)),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").Op("!=").Lit(0)).Body(
				jen.ID("v").Dot("Set").Call(
					jen.ID("updatedBeforeQueryKey"),
					jen.Qual("strconv", "FormatUint").Call(
						jen.ID("qf").Dot("UpdatedBefore"),
						jen.Lit(10),
					),
				)),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").Op("!=").Lit(0)).Body(
				jen.ID("v").Dot("Set").Call(
					jen.ID("updatedAfterQueryKey"),
					jen.Qual("strconv", "FormatUint").Call(
						jen.ID("qf").Dot("UpdatedAfter"),
						jen.Lit(10),
					),
				)),
			jen.ID("v").Dot("Set").Call(
				jen.ID("includeArchivedQueryKey"),
				jen.Qual("strconv", "FormatBool").Call(jen.ID("qf").Dot("IncludeArchived")),
			),
			jen.Return().ID("v"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ExtractQueryFilter can extract a QueryFilter from a request."),
		jen.Line(),
		jen.Func().ID("ExtractQueryFilter").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("QueryFilter")).Body(
			jen.ID("qf").Op(":=").Op("&").ID("QueryFilter").Valuesln(),
			jen.ID("qf").Dot("FromParams").Call(jen.ID("req").Dot("URL").Dot("Query").Call()),
			jen.Return().ID("qf"),
		),
		jen.Line(),
	)

	return code
}
