package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func databaseDotGo() *jen.File {
	ret := jen.NewFile("config")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("postgresProviderKey").Op("=").Lit("postgres").Var().ID("mariaDBProviderKey").Op("=").Lit("mariadb").Var().ID("sqliteProviderKey").Op("=").Lit("sqlite"),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
