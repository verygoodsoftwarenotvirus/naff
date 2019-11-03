package queriers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkgRoot string, types []models.DataType, vendor wordsmith.SuperPalabra) *jen.File {
	ret := jen.NewFile(vendor.SingularPackageName())

	utils.AddImports(pkgRoot, types, ret)
	sn := vendor.Singular()

	isMariaDB := vendor.RouteName() == "mariadb" || vendor.RouteName() == "maria_db"
	var (
		dbTrail string
	)
	if !isMariaDB {
		dbTrail = "DB"
	} else {
		dbTrail = "Connection"
	}

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injection"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.IDf("Provide%s%s", sn, dbTrail),
				jen.IDf("Provide%s", sn),
			),
		),
		jen.Line(),
	)
	return ret
}
