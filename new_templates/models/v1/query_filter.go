package models

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func queryFilterDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("MaxLimit").Op("=").Lit(250).Var().ID("DefaultLimit").Op("=").Lit(20).Var().ID("pageKey").Op("=").Lit("page").Var().ID("limitKey").Op("=").Lit("limit").Var().ID("createdBeforeKey").Op("=").Lit("created_before").Var().ID("createdAfterKey").Op("=").Lit("created_after").Var().ID("updatedBeforeKey").Op("=").Lit("updated_before").Var().ID("updatedAfterKey").Op("=").Lit("updated_after").Var().ID("sortByKey").Op("=").Lit("sort_by"))
	ret.Add(jen.Null().Type().ID("QueryFilter").Struct(
		jen.ID("Page").ID("uint64"),
		jen.ID("Limit").ID("uint64"),
		jen.ID("CreatedAfter").ID("uint64"),
		jen.ID("CreatedBefore").ID("uint64"),
		jen.ID("UpdatedAfter").ID("uint64"),
		jen.ID("UpdatedBefore").ID("uint64"),
		jen.ID("SortBy").ID("sortType"),
	),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
