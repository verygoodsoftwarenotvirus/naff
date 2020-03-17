package httpserver

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(pkg, ret)

	buildNewSetArgs := func() []jen.Code {
		args := []jen.Code{}

		for _, typ := range pkg.DataTypes {
			sn := typ.Name.Singular()

			if typ.BelongsToUser {
				args = append(args, jen.IDf("Provide%sServiceUserIDFetcher", sn))
			}
		}

		args = append(args,
			jen.ID("ProvideUsernameFetcher"),
			jen.ID("ProvideOAuth2ServiceClientIDFetcher"),
			jen.ID("ProvideAuthUserIDFetcher"),
			jen.ID("ProvideWebhooksUserIDFetcher"),
		)

		for _, typ := range pkg.DataTypes {
			sn := typ.Name.Singular()
			args = append(args, jen.IDf("Provide%sIDFetcher", sn))

			if typ.BelongsToStruct != nil {
				args = append(args, jen.IDf("Provide%sService%sIDFetcher", sn, typ.BelongsToStruct.Singular()))
			}
		}

		args = append(args, jen.ID("ProvideWebhookIDFetcher"))

		return args
	}

	ret.Add(
		jen.Var().Defs(
			jen.ID("paramFetcherProviders").Op("=").Qual("github.com/google/wire", "NewSet").Callln(
				buildNewSetArgs()...,
			),
		),
		jen.Line(),
	)

	for _, typ := range pkg.DataTypes {
		sn := typ.Name.Singular()
		pn := typ.Name.PackageName()

		if typ.BelongsToUser {
			ret.Add(
				jen.Commentf("Provide%sServiceUserIDFetcher provides a UserIDFetcher", sn),
				jen.Line(),
				jen.Func().IDf("Provide%sServiceUserIDFetcher", sn).Params().Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1", pn), "UserIDFetcher")).Block(
					jen.Return().ID("UserIDFetcher"),
				),
				jen.Line(),
			)
		}
		if typ.BelongsToStruct != nil {
			ret.Add(
				jen.Commentf("Provide%sService%sIDFetcher provides a %sIDFetcher", sn, typ.BelongsToStruct.Singular(), typ.BelongsToStruct.Singular()),
				jen.Line(),
				jen.Func().IDf("Provide%sService%sIDFetcher", sn, typ.BelongsToStruct.Singular()).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1", pn), fmt.Sprintf("%sIDFetcher", typ.BelongsToStruct.Singular()))).Block(
					jen.Return().IDf("buildChi%sIDFetcher", typ.BelongsToStruct.Singular()).Call(jen.ID("logger")),
				),
				jen.Line(),
			)
		}

		ret.Add(
			jen.Commentf("Provide%sIDFetcher provides an %sIDFetcher", sn, sn),
			jen.Line(),
			jen.Func().IDf("Provide%sIDFetcher", sn).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1", pn), fmt.Sprintf("%sIDFetcher", sn))).Block(
				jen.Return().IDf("buildChi%sIDFetcher", sn).Call(jen.ID("logger")),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Comment("ProvideUsernameFetcher provides a UsernameFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideUsernameFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/users"), "UserIDFetcher")).Block(
			jen.Return().ID("buildChiUserIDFetcher").Call(jen.ID("logger")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideAuthUserIDFetcher provides a UsernameFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideAuthUserIDFetcher").Params().Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/auth"), "UserIDFetcher")).Block(
			jen.Return().ID("UserIDFetcher"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksUserIDFetcher provides a UserIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksUserIDFetcher").Params().Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/webhooks"), "UserIDFetcher")).Block(
			jen.Return().ID("UserIDFetcher"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhookIDFetcher provides an WebhookIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideWebhookIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/webhooks"), "WebhookIDFetcher")).Block(
			jen.Return().ID("buildChiWebhookIDFetcher").Call(jen.ID("logger")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideOAuth2ServiceClientIDFetcher provides a ClientIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ServiceClientIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/oauth2clients"), "ClientIDFetcher")).Block(
			jen.Return().ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID("logger")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserIDFetcher fetches a user ID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("UserIDFetcher").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().ID("req").Dot("Context").Call().Dot("Value").Call(jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserIDKey")).Assert(jen.ID("uint64")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildChiUserIDFetcher builds a function that fetches a Username from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiUserIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/users"), "UserIDFetcher")).Block(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/users"), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("fetching user ID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)

	for _, typ := range pkg.DataTypes {
		n := typ.Name
		sn := n.Singular()
		pn := n.PackageName()

		ret.Add(
			jen.Commentf("buildChi%sIDFetcher builds a function that fetches a %sID from a request routed by chi.", sn, sn),
			jen.Line(),
			jen.Func().IDf("buildChi%sIDFetcher", sn).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
				jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
					jen.Comment("we can generally disregard this error only because we should be able to validate"),
					jen.Comment("that the string only contains numbers via chi's regex url param feature."),
					jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1", pn), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
					jen.If(jen.Err().Op("!=").ID("nil")).Block(
						jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("fetching %sID from request", sn)),
					),
					jen.Return().ID("u"),
				),
			),
			jen.Line(),
		)

		// if typ.BelongsToStruct != nil {
		// 	ret.Add(
		// 		jen.Commentf(
		// 			"build%sServiceChi%sIDFetcher builds a function that fetches a %sID from a request routed by chi.",
		// 			typ.Name.Plural(),
		// 			typ.BelongsToStruct.Singular(),
		// 			typ.BelongsToStruct.Singular(),
		// 		),
		// 		jen.Line(),
		// 		jen.Func().IDf("build%sChi%sIDFetcher",typ.Name.Plural(),, typ.BelongsToStruct.Singular()).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
		// 			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
		// 				jen.Comment("we can generally disregard this error only because we should be able to validate"),
		// 				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
		// 				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1", typ.BelongsToStruct.PackageName()), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
		// 				jen.If(jen.Err().Op("!=").ID("nil")).Block(
		// 					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Litf("fetching %sID from request", typ.BelongsToStruct.Singular())),
		// 				),
		// 				jen.Return().ID("u"),
		// 			),
		// 		),
		// 		jen.Line(),
		// 	)
		// }

	}

	ret.Add(
		jen.Comment("chiWebhookIDFetcher fetches a WebhookID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiWebhookIDFetcher").Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger")).Params(jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64"))).Block(
			jen.Return().Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate"),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/webhooks"), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("fetching WebhookID from request")),
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
				jen.List(jen.ID("u"), jen.Err()).Op(":=").Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID("req"), jen.Qual(filepath.Join(pkg.OutputPath, "services/v1/oauth2clients"), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().Op("!=").ID("nil")).Block(
					jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("fetching OAuth2ClientID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)
	return ret
}
