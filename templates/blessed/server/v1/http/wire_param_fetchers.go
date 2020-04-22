package httpserver

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	buildNewSetArgs := func() []jen.Code {
		args := []jen.Code{}

		for _, typ := range proj.DataTypes {
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

		for _, typ := range proj.DataTypes {
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
			jen.ID("paramFetcherProviders").Equals().Qual("github.com/google/wire", "NewSet").Callln(
				buildNewSetArgs()...,
			),
		),
		jen.Line(),
	)

	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()
		pn := typ.Name.PackageName()

		if typ.BelongsToUser {
			ret.Add(
				jen.Commentf("Provide%sServiceUserIDFetcher provides a UserIDFetcher", sn),
				jen.Line(),
				jen.Func().IDf("Provide%sServiceUserIDFetcher", sn).Params().Params(jen.Qual(proj.ServiceV1Package(pn), "UserIDFetcher")).Block(
					jen.Return().ID("UserIDFetcher"),
				),
				jen.Line(),
			)
		}
		if typ.BelongsToStruct != nil {
			ret.Add(
				jen.Commentf("Provide%sService%sIDFetcher provides a %sIDFetcher", sn, typ.BelongsToStruct.Singular(), typ.BelongsToStruct.Singular()),
				jen.Line(),
				jen.Func().IDf("Provide%sService%sIDFetcher", sn, typ.BelongsToStruct.Singular()).Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1Package(pn), fmt.Sprintf("%sIDFetcher", typ.BelongsToStruct.Singular()))).Block(
					jen.Return().IDf("buildChi%sIDFetcher", typ.BelongsToStruct.Singular()).Call(jen.ID(constants.LoggerVarName)),
				),
				jen.Line(),
			)
		}

		ret.Add(
			jen.Commentf("Provide%sIDFetcher provides an %sIDFetcher", sn, sn),
			jen.Line(),
			jen.Func().IDf("Provide%sIDFetcher", sn).Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1Package(pn), fmt.Sprintf("%sIDFetcher", sn))).Block(
				jen.Return().IDf("buildChi%sIDFetcher", sn).Call(jen.ID(constants.LoggerVarName)),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Comment("ProvideUsernameFetcher provides a UsernameFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideUsernameFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1UsersPackage(), "UserIDFetcher")).Block(
			jen.Return().ID("buildChiUserIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideAuthUserIDFetcher provides a UsernameFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideAuthUserIDFetcher").Params().Params(jen.Qual(proj.ServiceV1AuthPackage(), "UserIDFetcher")).Block(
			jen.Return().ID("UserIDFetcher"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksUserIDFetcher provides a UserIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksUserIDFetcher").Params().Params(jen.Qual(proj.ServiceV1WebhooksPackage(), "UserIDFetcher")).Block(
			jen.Return().ID("UserIDFetcher"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhookIDFetcher provides an WebhookIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideWebhookIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1WebhooksPackage(), "WebhookIDFetcher")).Block(
			jen.Return().ID("buildChiWebhookIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideOAuth2ServiceClientIDFetcher provides a ClientIDFetcher"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ServiceClientIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "ClientIDFetcher")).Block(
			jen.Return().ID("buildChiOAuth2ClientIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserIDFetcher fetches a user ID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("UserIDFetcher").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block( //if userID, ok := req.Context().Value(models.UserIDKey).(uint64); ok {
			jen.If(jen.List(jen.ID("userID"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "UserIDKey")).Assert(jen.Uint64()), jen.ID("ok")).Block(
				jen.Return(jen.ID("userID")),
			),
			jen.Return(jen.Zero()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildChiUserIDFetcher builds a function that fetches a Username from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiUserIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1UsersPackage(), "UserIDFetcher")).Block(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceV1UsersPackage(), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching user ID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)

	for _, typ := range proj.DataTypes {
		n := typ.Name
		sn := n.Singular()
		pn := n.PackageName()

		ret.Add(
			jen.Commentf("buildChi%sIDFetcher builds a function that fetches a %sID from a request routed by chi.", sn, sn),
			jen.Line(),
			jen.Func().IDf("buildChi%sIDFetcher", sn).Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Block(
				jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Comment("we can generally disregard this error only because we should be able to validate"),
					jen.Comment("that the string only contains numbers via chi's regex url param feature."),
					jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceV1Package(pn), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("fetching %sID from request", sn)),
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
		// 		jen.Func().IDf("build%sChi%sIDFetcher",typ.Name.Plural(),, typ.BelongsToStruct.Singular()).Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Block(
		// 			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
		// 				jen.Comment("we can generally disregard this error only because we should be able to validate"),
		// 				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
		// 				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServicesV1Package(typ.BelongsToStruct.PackageName()), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
		// 				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
		// 					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("fetching %sID from request", typ.BelongsToStruct.Singular())),
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
		jen.Func().ID("buildChiWebhookIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Block(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate"),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceV1WebhooksPackage(), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching WebhookID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("chiOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildChiOAuth2ClientIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Block(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate"),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching OAuth2ClientID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	)
	return ret
}
