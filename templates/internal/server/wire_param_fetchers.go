package server

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireParamFetchersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildWireParamFetchersVarDeclarations(proj)...)

	code.Add(buildProvideUsersServiceUserIDFetcher(proj)...)
	code.Add(buildProvideOAuth2ClientsServiceClientIDFetcher(proj)...)
	code.Add(buildProvideWebhooksServiceWebhookIDFetcher(proj)...)
	code.Add(buildProvideWebhooksServiceUserIDFetcher(proj)...)

	for _, typ := range proj.DataTypes {
		ownerTypes := proj.FindOwnerTypeChain(typ)
		for _, ot := range ownerTypes {
			code.Add(buildProvideSomethingServiceThingIDFetcher(proj, typ, ot)...)
		}

		code.Add(buildProvideSomethingServiceOwnerTypeIDFetcher(proj, typ)...)

		if typ.OwnedByAUserAtSomeLevel(proj) {
			code.Add(buildProvideSomethingServiceUserIDFetcher(proj, typ)...)
		}
	}

	code.Add(buildUserIDFetcherFromRequestContext(proj)...)
	code.Add(buildBuildRouteParamUserIDFetcher(proj)...)

	for _, typ := range proj.DataTypes {
		code.Add(buildBuildRouteParamSomethingIDFetcher(proj, typ)...)
	}

	code.Add(buildBuildRouteParamWebhookIDFetcher(proj)...)
	code.Add(buildBuildRouteParamOAuth2ClientIDFetcher(proj)...)

	return code
}

func buildWireParamFetchersVarDeclarations(proj *models.Project) []jen.Code {
	newSetArgs := []jen.Code{
		jen.ID("ProvideUsersServiceUserIDFetcher"),
		jen.ID("ProvideOAuth2ClientsServiceClientIDFetcher"),
		jen.ID("ProvideWebhooksServiceWebhookIDFetcher"),
		jen.ID("ProvideWebhooksServiceUserIDFetcher"),
	}

	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()
		pn := typ.Name.Plural()

		for _, ot := range proj.FindOwnerTypeChain(typ) {
			newSetArgs = append(newSetArgs, jen.IDf("Provide%sService%sIDFetcher", pn, ot.Name.Singular()))
		}

		newSetArgs = append(newSetArgs, jen.IDf("Provide%sService%sIDFetcher", pn, sn))

		if typ.OwnedByAUserAtSomeLevel(proj) {
			newSetArgs = append(newSetArgs, jen.IDf("Provide%sServiceUserIDFetcher", pn))
		}
	}

	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("paramFetcherProviders").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				newSetArgs...,
			),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideSomethingServiceUserIDFetcher(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pkgN := typ.Name.PackageName()

	lines := []jen.Code{
		jen.Commentf("Provide%sServiceUserIDFetcher provides a UserIDFetcher.", pn),
		jen.Line(),
		jen.Func().IDf("Provide%sServiceUserIDFetcher", pn).Params().Params(jen.Qual(proj.ServiceV1Package(pkgN), "UserIDFetcher")).Body(
			jen.Return().ID("userIDFetcherFromRequestContext"),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideSomethingServiceThingIDFetcher(proj *models.Project, typ models.DataType, ownerType models.DataType) []jen.Code {
	ots := ownerType.Name.Singular()
	pn := typ.Name.Plural()
	pkgN := typ.Name.PackageName()

	lines := []jen.Code{
		jen.Commentf("Provide%sService%sIDFetcher provides %s %sIDFetcher.", pn, ots, wordsmith.AOrAn(ots), ots),
		jen.Line(),
		jen.Func().IDf("Provide%sService%sIDFetcher", pn, ots).Params(constants.LoggerParam()).Params(jen.Qual(proj.ServiceV1Package(pkgN), fmt.Sprintf("%sIDFetcher", ots))).Body(
			jen.Return().IDf("buildRouteParam%sIDFetcher", ots).Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideSomethingServiceOwnerTypeIDFetcher(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pkgN := typ.Name.PackageName()

	lines := []jen.Code{
		jen.Commentf("Provide%sService%sIDFetcher provides %s %sIDFetcher.", pn, sn, wordsmith.AOrAn(sn), sn),
		jen.Line(),
		jen.Func().IDf("Provide%sService%sIDFetcher", pn, sn).Params(constants.LoggerParam()).Params(jen.Qual(proj.ServiceV1Package(pkgN), fmt.Sprintf("%sIDFetcher", sn))).Body(
			jen.Return().IDf("buildRouteParam%sIDFetcher", sn).Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideUsersServiceUserIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideUsersServiceUserIDFetcher provides a UsernameFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideUsersServiceUserIDFetcher").Params(constants.LoggerParam()).Params(jen.Qual(proj.ServiceUsersPackage(), "UserIDFetcher")).Body(
			jen.Return().ID("buildRouteParamUserIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideWebhooksServiceUserIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideWebhooksServiceUserIDFetcher provides a UserIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksServiceUserIDFetcher").Params().Params(jen.Qual(proj.ServiceWebhooksPackage(), "UserIDFetcher")).Body(
			jen.Return().ID("userIDFetcherFromRequestContext"),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideWebhooksServiceWebhookIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideWebhooksServiceWebhookIDFetcher provides an WebhookIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideWebhooksServiceWebhookIDFetcher").Params(constants.LoggerParam()).Params(jen.Qual(proj.ServiceWebhooksPackage(), "WebhookIDFetcher")).Body(
			jen.Return().ID("buildRouteParamWebhookIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideOAuth2ClientsServiceClientIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideOAuth2ClientsServiceClientIDFetcher provides a ClientIDFetcher."),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientsServiceClientIDFetcher").Params(constants.LoggerParam()).Params(jen.Qual(proj.ServiceOAuth2ClientsPackage(), "ClientIDFetcher")).Body(
			jen.Return().ID("buildRouteParamOAuth2ClientIDFetcher").Call(jen.ID(constants.LoggerVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildUserIDFetcherFromRequestContext(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("userIDFetcherFromRequestContext fetches a user ID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("userIDFetcherFromRequestContext").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body( //if userID, ok := req.Context().Value(models.UserIDKey).(uint64); ok {
			jen.If(jen.List(jen.ID("si"), jen.ID("ok")).Assign().ID(constants.RequestVarName).Dot("Context").Call().Dot("Value").Call(jen.Qual(proj.TypesPackage(), "SessionInfoKey")).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "SessionInfo")), jen.ID("ok").And().ID("si").DoesNotEqual().Nil()).Body(
				jen.Return(jen.ID("si").Dot(constants.UserIDFieldName)),
			),
			jen.Return(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildRouteParamUserIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildRouteParamUserIDFetcher builds a function that fetches a Username from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildRouteParamUserIDFetcher").Params(constants.LoggerParam()).Params(jen.Qual(proj.ServiceUsersPackage(), "UserIDFetcher")).Body(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceUsersPackage(), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching user ID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildRouteParamSomethingIDFetcher(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.PackageName()

	lines := []jen.Code{
		jen.Commentf("buildRouteParam%sIDFetcher builds a function that fetches a %sID from a request routed by chi.", sn, sn),
		jen.Line(),
		jen.Func().IDf("buildRouteParam%sIDFetcher", sn).Params(constants.LoggerParam()).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Body(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
				jen.Comment("we can generally disregard this error only because we should be able to validate."),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceV1Package(pn), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Litf("fetching %sID from request", sn)),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildRouteParamWebhookIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildRouteParamWebhookIDFetcher fetches a WebhookID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildRouteParamWebhookIDFetcher").Params(constants.LoggerParam()).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Body(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
				jen.Comment("we can generally disregard this error only because we should be able to validate."),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceWebhooksPackage(), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching WebhookID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildRouteParamOAuth2ClientIDFetcher(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildRouteParamOAuth2ClientIDFetcher fetches a OAuth2ClientID from a request routed by chi."),
		jen.Line(),
		jen.Func().ID("buildRouteParamOAuth2ClientIDFetcher").Params(constants.LoggerParam()).Params(jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64())).Body(
			jen.Return().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
				jen.Comment("we can generally disregard this error only because we should be able to validate."),
				jen.Comment("that the string only contains numbers via chi's regex url param feature."),
				jen.List(jen.ID("u"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(jen.Qual("github.com/go-chi/chi", "URLParam").Call(jen.ID(constants.RequestVarName), jen.Qual(proj.ServiceOAuth2ClientsPackage(), "URIParamKey")), jen.Lit(10), jen.Lit(64)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching OAuth2ClientID from request")),
				),
				jen.Return().ID("u"),
			),
		),
		jen.Line(),
	}

	return lines
}
