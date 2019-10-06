package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func wireParamFetchersDotGo() *jen.File {
	ret := jen.NewFile("httpserver")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("paramFetcherProviders").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideUserIDFetcher"), jen.ID("ProvideUsernameFetcher"), jen.ID("ProvideOAuth2ServiceClientIDFetcher"), jen.ID("ProvideAuthUserIDFetcher"), jen.ID("ProvideItemIDFetcher"), jen.ID("ProvideWebhooksUserIDFetcher"), jen.ID("ProvideWebhookIDFetcher")),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
