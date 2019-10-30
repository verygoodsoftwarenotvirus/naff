package oauth2clients

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo(pkgRoot string) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(ret)

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
		jen.Func().ID("ProvideOAuth2ClientDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientDataServer")).Block(
			jen.Return().ID("s"),
		),
		jen.Line(),
	)
	return ret
}
