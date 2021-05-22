package fake

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)
	utils.AddImports(proj, code)

	code.Add(buildBuildFleshedOutQueryFilter(proj)...)

	return code
}

func buildBuildFleshedOutQueryFilter(proj *models.Project) []jen.Code {
	funcName := "BuildFleshedOutQueryFilter"

	lines := []jen.Code{
		jen.Commentf("%s builds a fully fleshed out QueryFilter.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Body(
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), "QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(10),
					jen.ID("Limit").MapAssign().Lit(20),
					jen.ID("CreatedAfter").MapAssign().Add(utils.FakeUnixTimeFunc()),
					jen.ID("CreatedBefore").MapAssign().Add(utils.FakeUnixTimeFunc()),
					jen.ID("UpdatedAfter").MapAssign().Add(utils.FakeUnixTimeFunc()),
					jen.ID("UpdatedBefore").MapAssign().Add(utils.FakeUnixTimeFunc()),
					jen.ID("SortBy").MapAssign().Qual(proj.ModelsV1Package(), "SortAscending"),
				),
			),
		),
	}

	return lines
}
