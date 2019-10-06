package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireDotGo() *jen.File {
	ret := jen.NewFile("webhooks")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideWebhooksService"), jen.ID("ProvideWebhookDataManager"), jen.ID("ProvideWebhookDataServer")),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
