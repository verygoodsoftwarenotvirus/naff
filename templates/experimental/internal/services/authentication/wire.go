package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(
				jen.ID("ProvideService"),
				jen.ID("wire").Dot("FieldsOf").Call(
					jen.ID("new").Call(jen.Op("*").ID("Config")),
					jen.Lit("Cookies"),
					jen.Lit("PASETO"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
