package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideOAuth2ClientsService"), jen.ID("ProvideOAuth2ClientDataServer")),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
