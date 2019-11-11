package oauth2clients

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("Providers are what we provide for dependency injection"),
			jen.ID("Providers").Op("=").Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID("ProvideOAuth2ClientsService"),
				jen.ID("ProvideOAuth2ClientDataServer"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideOAuth2ClientDataServer is an arbitrary function for dependency injection's sake"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)
	return ret
}
