package querybuilders

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	suffix := "DB"
	if dbvendor.SingularPackageName() == "mariadb" {
		suffix = "Connection"
	}

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injection."),
			jen.ID("Providers").Op("=").Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.IDf("Provide%s%s", dbvendor.Singular(), suffix),
				jen.IDf("Provide%s", dbvendor.Singular()),
			),
		),
		jen.Newline(),
	)

	return code
}
