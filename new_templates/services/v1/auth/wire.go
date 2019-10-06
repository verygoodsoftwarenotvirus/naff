package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireDotGo() *jen.File {
	ret := jen.NewFile("auth")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideAuthService"), jen.ID("ProvideWebsocketAuthFunc"), jen.ID("ProvideOAuth2ClientValidator")),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
