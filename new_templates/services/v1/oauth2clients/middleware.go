package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("scopesSeparator").Op("=").Lit(",").Var().ID("apiPathPrefix").Op("=").Lit("/api/v1/"))
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("fetchOAuth2ClientFromRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Block(
		jen.List(jen.ID("client"), jen.ID("ok")).Op(":=").ID("req").Dot(
			"Context",
		).Call().Dot(
			"Value",
		).Call(jen.ID("models").Dot(
			"OAuth2ClientKey",
		)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		)),
		jen.ID("_").Op("=").ID("ok"),
		jen.Return().ID("client"),
	),
	)
	ret.Add(jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("fetchOAuth2ClientIDFromRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Block(
		jen.List(jen.ID("clientID"), jen.ID("ok")).Op(":=").ID("req").Dot(
			"Context",
		).Call().Dot(
			"Value",
		).Call(jen.ID("clientIDKey")).Assert(jen.ID("string")),
		jen.ID("_").Op("=").ID("ok"),
		jen.Return().ID("clientID"),
	),
	)
	return ret
}
