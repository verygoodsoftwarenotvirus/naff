package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("users")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers is what we provide for dependency injectors."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideUsersService"),
				jen.ID("ProvideUserDataServer"),
				jen.ID("ProvideUserDataManager"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUserDataManager is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideUserDataManager").Params(jen.ID("db").Qual(proj.DatabaseV1Package(), "DataManager")).Params(jen.Qual(proj.ModelsV1Package(), "UserDataManager")).Block(
			jen.Return().ID("db"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUserDataServer is an arbitrary function for dependency injection's sake."),
		jen.Line(),
		jen.Func().ID("ProvideUserDataServer").Params(jen.ID("s").PointerTo().ID("Service")).Params(jen.Qual(proj.ModelsV1Package(), "UserDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)

	return code
}
