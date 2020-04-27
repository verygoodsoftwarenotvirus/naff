package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) *jen.File {
	sn := dbvendor.Singular()
	spn := dbvendor.SingularPackageName()

	ret := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, ret)

	var dbTrail string

	if !isMariaDB(dbvendor) {
		dbTrail = "DB"
	} else {
		dbTrail = "Connection"
	}

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injection."),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				jen.IDf("Provide%s%s", sn, dbTrail),
				jen.IDf("Provide%s", sn),
			),
		),
		jen.Line(),
	)

	return ret
}
