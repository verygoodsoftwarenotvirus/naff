package users

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo(pkgRoot string) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injectors"),
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Callln(
				jen.ID("ProvideUsersService"),
				jen.ID("ProvideUserDataServer"),
				jen.ID("ProvideUserDataManager"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUserDataManager is an arbitrary function for dependency injection's sake"),
		jen.Line(),
		jen.Func().ID("ProvideUserDataManager").Params(jen.ID("db").Qual(filepath.Join(pkgRoot, "database/v1"), "Database")).Params(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "UserDataManager")).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUserDataServer is an arbitrary function for dependency injection's sake"),
		jen.Line(),
		jen.Func().ID("ProvideUserDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "UserDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)
	return ret
}
