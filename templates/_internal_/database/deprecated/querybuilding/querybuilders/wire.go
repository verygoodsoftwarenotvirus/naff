package querybuilders

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	suffix := "DB"
	if dbvendor.SingularPackageName() == "mysql" {
		suffix = "Connection"
	}

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injection."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.IDf("Provide%s%s", dbvendor.Singular(), suffix),
				jen.IDf("Provide%s", dbvendor.Singular()),
			),
		),
		jen.Newline(),
	)

	return code
}
