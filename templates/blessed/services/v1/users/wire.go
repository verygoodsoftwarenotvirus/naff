package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injectors"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(
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
		jen.Func().ID("ProvideUserDataManager").Params(jen.ID("db").Qual(proj.DatabaseV1Package(), "Database")).Params(jen.Qual(proj.ModelsV1Package(), "UserDataManager")).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUserDataServer is an arbitrary function for dependency injection's sake"),
		jen.Line(),
		jen.Func().ID("ProvideUserDataServer").Params(jen.ID("s").PointerTo().ID("Service")).Params(jen.Qual(proj.ModelsV1Package(), "UserDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	return ret
}
