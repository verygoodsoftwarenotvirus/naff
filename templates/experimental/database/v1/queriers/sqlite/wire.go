package sqlite

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireDotGo() *jen.File {
	ret := jen.NewFile("sqlite")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideSqliteDB"), jen.ID("ProvideSqlite")),

		jen.Line(),
	)
	return ret
}
