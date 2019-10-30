package httpserver

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func wireParamFetchersDotGo(pkgRoot string) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("paramFetcherProviders").Op("=").Qual("github.com/google/wire", "NewSet").Callln(
				jen.ID("ProvideUserIDFetcher"),
				jen.ID("ProvideUsernameFetcher"),
				jen.ID("ProvideOAuth2ServiceClientIDFetcher"),
				jen.ID("ProvideAuthUserIDFetcher"),
				jen.ID("ProvideItemIDFetcher"),
				jen.ID("ProvideWebhooksUserIDFetcher"),
				jen.ID("ProvideWebhookIDFetcher"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUserIDFetcher provides a UserIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideUserIDFetcher").Params().Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/items"), "UserIDFetcher")).Block(
			jen.Return().ID("UserIDFetcher"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideItemIDFetcher provides an ItemIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideItemIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/items"), "ItemIDFetcher")).Block(
			jen.Return().ID("buildChiItemIDFetcher").Call(jen.ID("logger")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUsernameFetcher provides a UsernameFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideUsernameFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/users"), "UserIDFetcher")).Block(
			jen.Return().ID("buildChiUserIDFetcher").Call(jen.ID("logger")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideAuthUserIDFetcher provides a UsernameFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideAuthUserIDFetcher").Params().Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/auth"), "UserIDFetcher")).Block(
			jen.Return().ID("UserIDFetcher"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksUserIDFetcher provides a UserIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksUserIDFetcher").Params().Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/webhooks"), "UserIDFetcher")).Block(
			jen.Return().ID("UserIDFetcher"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhookIDFetcher provides an WebhookIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideWebhookIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/webhooks"), "WebhookIDFetcher")).Block(
			jen.Return().ID("buildChiWebhookIDFetcher").Call(jen.ID("logger")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideOAuth2ServiceClientIDFetcher provides a ClientIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ServiceClientIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/oauth2clients"), "ClientIDFetcher")).Block(
			jen.Return().ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID("logger")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserIDFetcher fetches a user ID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("UserIDFetcher").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().ID("req").Dot("Context").Call().Dot("Value").Call(jen.Qual(filepath.Join(pkgRoot, "models/v1"), "UserIDKey")).Assert(jen.ID("uint64")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildChiUserIDFetcher builds a function that fetches a Username from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiUserIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkgRoot, "services/v1/users"), "UserIDFetcher")).Block(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkgRoot, "services/v1/users"), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("fetching user ID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("chiItemIDFetcher fetches a ItemID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiItemIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate"),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkgRoot, "services/v1/items"), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("fetching ItemID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("chiWebhookIDFetcher fetches a WebhookID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiWebhookIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate"),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkgRoot, "services/v1/webhooks"), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("fetching WebhookID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("chiOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiOAuth2ClientIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate"),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkgRoot, "services/v1/oauth2clients"), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("fetching OAuth2ClientID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)
	return ret
}
