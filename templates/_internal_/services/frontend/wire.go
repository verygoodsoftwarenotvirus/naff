package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(
				jen.ID("ProvideService"),
				jen.ID("ProvideAuthService"),
				jen.ID("ProvideUsersService"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideAuthService does what I hope one day wire figures out how to do."),
		jen.Line(),
		jen.Func().ID("ProvideAuthService").Params(jen.ID("x").ID("types").Dot("AuthService")).Params(jen.ID("AuthService")).Body(
			jen.Return().ID("x")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideUsersService does what I hope one day wire figures out how to do."),
		jen.Line(),
		jen.Func().ID("ProvideUsersService").Params(jen.ID("x").ID("types").Dot("UserDataService")).Params(jen.ID("UsersService")).Body(
			jen.Return().ID("x")),
		jen.Line(),
	)

	return code
}
