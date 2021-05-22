package querybuilding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	sn := dbvendor.Singular()
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", spn), spn)

	utils.AddImports(proj, code)

	var dbTrail string

	if !isMariaDB(dbvendor) {
		dbTrail = "DB"
	} else {
		dbTrail = "Connection"
	}

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injection."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.IDf("Provide%s%s", sn, dbTrail),
				jen.IDf("Provide%s", sn),
			),
		),
		jen.Line(),
	)

	return code
}
