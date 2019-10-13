package oauth2clients

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideOAuth2ClientsService"), jen.ID("ProvideOAuth2ClientDataServer")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideOAuth2ClientDataServer is an arbitrary function for dependency injection's sake").ID("ProvideOAuth2ClientDataServer").Params(jen.ID("s").Op("*").ID("Service")).Params(jen.ID("models").Dot(
		"OAuth2ClientDataServer",
	)).Block(
		jen.Return().ID("s"),
	),

		jen.Line(),
	)
	return ret
}
