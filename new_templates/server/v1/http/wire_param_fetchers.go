package httpserver

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireParamFetchersDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("paramFetcherProviders").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideUserIDFetcher"), jen.ID("ProvideUsernameFetcher"), jen.ID("ProvideOAuth2ServiceClientIDFetcher"), jen.ID("ProvideAuthUserIDFetcher"), jen.ID("ProvideItemIDFetcher"), jen.ID("ProvideWebhooksUserIDFetcher"), jen.ID("ProvideWebhookIDFetcher")),
	)

	return ret
}
