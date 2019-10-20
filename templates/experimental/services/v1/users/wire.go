package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideUsersService"), jen.ID("ProvideUserDataServer"), jen.ID("ProvideUserDataManager")),
	jen.Line(),
	)

	ret.Add(
	jen.Comment("ProvideUserDataManager is an arbitrary function for dependency injection's sake"),
	jen.Line(),
	jen.Func().ID("ProvideUserDataManager").Params(jen.ID("db").ID("database").Dot("Database")).Params(jen.ID("models").Dot(
		"UserDataManager",
	)).Block(
		jen.Return().ID("db"),
	),
	jen.Line(),
	)

	ret.Add(
	jen.Comment("ProvideUserDataServer is an arbitrary function for dependency injection's sake"),
	jen.Line(),
	jen.Func().ID("ProvideUserDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.ID("models").Dot(
		"UserDataServer",
	)).Block(
		jen.Return().ID("s"),
	),
	jen.Line(),
	)
	return ret
}
