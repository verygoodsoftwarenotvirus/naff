package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("sqlite")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideSqliteDB"), jen.ID("ProvideSqlite")),
	jen.Line(),
	)
	return ret
}
