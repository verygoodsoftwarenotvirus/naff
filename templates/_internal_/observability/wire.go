package observability

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(jen.ID("wire").Dot("FieldsOf").Call(
			jen.ID("new").Call(jen.Op("*").ID("Config")),
			jen.Lit("Metrics"),
			jen.Lit("Tracing"),
		)),
		jen.Line(),
	)

	return code
}
