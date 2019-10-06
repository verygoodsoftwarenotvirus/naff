package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func itemsDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("itemsTableName").Op("=").Lit("items"))
	ret.Add(jen.Null().Var().ID("itemsTableColumns").Op("=").Index().ID("string").Valuesln(jen.Lit("id"), jen.Lit("name"), jen.Lit("details"), jen.Lit("created_on"), jen.Lit("updated_on"), jen.Lit("archived_on"), jen.Lit("belongs_to")))
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Null().Var().ID("allItemsCountQueryBuilder").Qual("sync", "Once").Var().ID("allItemsCountQuery").ID("string"))
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
