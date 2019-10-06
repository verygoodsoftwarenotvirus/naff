package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func databaseDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("postgresProviderKey").Op("=").Lit("postgres").Var().ID("mariaDBProviderKey").Op("=").Lit("mariadb").Var().ID("sqliteProviderKey").Op("=").Lit("sqlite"))
	ret.Add(jen.Func())
	return ret
}
