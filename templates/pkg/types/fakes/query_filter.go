package fakes

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFleshedOutQueryFilter builds a fully fleshed out QueryFilter."),
		jen.Line(),
		jen.Func().ID("BuildFleshedOutQueryFilter").Params().Params(jen.Op("*").ID("types").Dot("QueryFilter")).Body(
			jen.Return().Op("&").ID("types").Dot("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(10), jen.ID("Limit").Op(":").Lit(20), jen.ID("CreatedAfter").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("CreatedBefore").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("UpdatedAfter").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("UpdatedBefore").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("SortBy").Op(":").ID("types").Dot("SortAscending"))),
		jen.Line(),
	)

	return code
}
