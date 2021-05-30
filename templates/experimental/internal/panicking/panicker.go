package panicking

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func panickerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("Panicker").Interface(
				jen.ID("Panic").Params(jen.Interface()),
				jen.ID("Panicf").Params(jen.ID("format").ID("string"), jen.ID("args").Op("...").Interface()),
			),
		),
		jen.Line(),
	)

	return code
}
