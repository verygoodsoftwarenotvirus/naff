package fakes

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func miscDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeTime builds a fake time."),
		jen.Line(),
		jen.Func().ID("BuildFakeTime").Params().Params(jen.ID("uint64")).Body(
			jen.Return().Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call()),
		jen.Line(),
	)

	return code
}
