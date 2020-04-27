package httpserver

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("httpserver")

	utils.AddImports(proj, ret)

	buildNewSetArgs := func() []jen.Code {
		args := []jen.Code{}

		args = append(args,
			jen.ID("ProvideUsersServiceUserIDFetcher"),
			jen.ID("ProvideOAuth2ClientsServiceClientIDFetcher"),
			jen.ID("ProvideAuthServiceUserIDFetcher"),
		)

		for _, typ := range proj.DataTypes {
			sn := typ.Name.Singular()
			pn := typ.Name.Plural()

			if typ.OwnedByAUserAtSomeLevel(proj) {
				args = append(args, jen.IDf("Provide%sServiceUserIDFetcher", pn))
			}

			for _, ot := range proj.FindOwnerTypeChain(typ) {
				args = append(args, jen.IDf("Provide%sService%sIDFetcher", pn, ot.Name.Singular()))
			}

			args = append(args, jen.IDf("Provide%sService%sIDFetcher", pn, sn))
		}

		args = append(args,
			jen.ID("ProvideWebhooksServiceUserIDFetcher"),
			jen.ID("ProvideWebhooksServiceWebhookIDFetcher"),
		)

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
		pn := typ.Name.Plural()
		pkgN := typ.Name.PackageName()

		if typ.OwnedByAUserAtSomeLevel(proj) {
			ret.Add(
				jen.Commentf("Provide%sServiceUserIDFetcher provides a UserIDFetcher.", pn),
				jen.Line(),
				jen.Func().IDf("Provide%sServiceUserIDFetcher", pn).Params().Params(jen.Qual(proj.ServiceV1Package(pkgN), "UserIDFetcher")).Block(
					jen.Return().ID("userIDFetcherFromRequestContext"),
				),
				jen.Line(),
			)
		}

		for _, ot := range proj.FindOwnerTypeChain(typ) {
			ots := ot.Name.Singular()
			ret.Add(
				jen.Commentf("Provide%sService%sIDFetcher provides %s %sIDFetcher.", pn, ots, wordsmith.AOrAn(ots), ots),
				jen.Line(),
				jen.Func().IDf("Provide%sService%sIDFetcher", pn, ots).Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1Package(pkgN), fmt.Sprintf("%sIDFetcher", ots))).Block(
					jen.Return().IDf("buildRouteParam%sIDFetcher", ots).Call(jen.ID(constants.LoggerVarName)),
				),
				jen.Line(),
			)
		}

		ret.Add(
			jen.Commentf("Provide%sService%sIDFetcher provides %s %sIDFetcher.", pn, sn, wordsmith.AOrAn(sn), sn),
			jen.Line(),
			jen.Func().IDf("Provide%sService%sIDFetcher", pn, sn).Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1Package(pkgN), fmt.Sprintf("%sIDFetcher", sn))).Block(
				jen.Return().IDf("buildRouteParam%sIDFetcher", sn).Call(jen.ID(constants.LoggerVarName)),
			),
			jen.Line(),
		)
	}

	ret.Add(
		jen.Comment("ProvideUsersServiceUserIDFetcher provides a UsernameFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideUsersServiceUserIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1UsersPackage(), "UserIDFetcher")).Block(
			jen.Return().ID("buildRouteParamUserIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideAuthServiceUserIDFetcher provides a UsernameFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideAuthServiceUserIDFetcher").Params().Params(jen.Qual(proj.ServiceV1AuthPackage(), "UserIDFetcher")).Block(
			jen.Return().ID("userIDFetcherFromRequestContext"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksServiceUserIDFetcher").Params().Params(jen.Qual(proj.ServiceV1WebhooksPackage(), "UserIDFetcher")).Block(
			jen.Return().ID("userIDFetcherFromRequestContext"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksServiceWebhookIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1WebhooksPackage(), "WebhookIDFetcher")).Block(
			jen.Return().ID("buildRouteParamWebhookIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientsServiceClientIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1OAuth2ClientsPackage(), "ClientIDFetcher")).Block(
			jen.Return().ID("buildRouteParamOAuth2ClientIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("userIDFetcherFromRequestContext fetches a user ID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("userIDFetcherFromRequestContext").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block( //if userID, ok := req.Context().Value(models.UserIDKey).(uint64); ok {
			jen.If(jen.List(jen.ID("userID"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "UserIDKey")).Assert(jen.Uint64()), jen.ID("ok")).Block(
				jen.Return(jen.ID("userID")),
			),
			jen.Return(jen.Zero()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("buildRouteParamUserIDFetcher builds a function that fetches a Username from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildRouteParamUserIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Qual(proj.ServiceV1UsersPackage(), "UserIDFetcher")).Block(
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
			jen.Commentf("buildRouteParam%sIDFetcher builds a function that fetches a %sID from a request routed by chi..", sn, sn),
			jen.Line(),
			jen.Func().IDf("buildRouteParam%sIDFetcher", sn).Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Block(
				jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
					jen.Comment("we can generally disregard this error only because we should be able to validate."),
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

	}

	ret.Add(
		jen.Comment("buildRouteParamWebhookIDFetcher fetches a WebhookID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildRouteParamWebhookIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Block(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate."),
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
		jen.Comment("buildRouteParamOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildRouteParamOAuth2ClientIDFetcher").Params(jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger")).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Block(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Block(
				jen.Comment("we can generally disregard this error only because we should be able to validate."),
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
