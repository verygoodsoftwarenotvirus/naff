package iterables

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildIDFetchers(proj *models.Project, typ models.DataType, includePrimaryType bool) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	idFetches := []jen.Code{}
	for _, dep := range proj.FindOwnerTypeChain(typ) {
		tsn := dep.Name.Singular()
		tuvn := dep.Name.UnexportedVarName()

		idFetches = append(idFetches,
			jen.Commentf("determine %s ID.", dep.Name.SingularCommonName()),
			jen.IDf("%sID", tuvn).Assign().ID("s").Dotf("%sIDFetcher", tuvn).Call(jen.ID(constants.RequestVarName)),
			jen.Qualf(proj.InternalTracingPackage(), "Attach%sIDToSpan", tsn).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", tuvn)),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qualf(proj.ConstantKeysPackage(), "%sIDKey", tsn), jen.IDf("%sID", tuvn)),
			jen.Newline(),
		)
	}

	if includePrimaryType {
		idFetches = append(idFetches,
			jen.Commentf("determine %s ID.", scn),
			jen.IDf("%sID", uvn).Op(":=").ID("s").Dotf("%sIDFetcher", uvn).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
		)
	}

	return idFetches
}

func httpRoutesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	code.Add(
		jen.Const().Defs(
			jen.Commentf("%sIDURIParamKey is a standard string that we'll use to refer to %s IDs with.", sn, scn),
			jen.IDf("%sIDURIParamKey", sn).Op("=").Litf("%sID", uvn),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("parseBool differs from strconv.ParseBool in that it returns false by default."),
		jen.Newline(),
		jen.Func().ID("parseBool").Params(jen.ID("str").ID("string")).Params(jen.ID("bool")).Body(
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("str")))).Body(
				jen.Case(jen.Lit("1"), jen.Lit("t"), jen.Lit("true")).Body(
					jen.Return().ID("true")),
				jen.Default().Body(
					jen.Return().ID("false")),
			)),
		jen.Newline(),
	)

	code.Add(buildCreateHandler(proj, typ)...)
	code.Add(buildReadHandler(proj, typ)...)
	code.Add(buildExistenceHandler(proj, typ)...)
	code.Add(buildListHandler(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildSearchHandler(proj, typ)...)
	}

	code.Add(buildUpdateHandler(proj, typ)...)
	code.Add(buildArchiveHandler(proj, typ)...)
	code.Add(buildAuditEntryHandler(proj, typ)...)

	return code
}

func buildDBClientRetrievalMethodCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildReadHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	dbCallArgs := buildDBClientRetrievalMethodCallArgs(proj, typ)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("retrieving session context data"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("unauthenticated"),
				jen.Qual("net/http", "StatusUnauthorized"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
			jen.ID("span"),
			jen.ID("sessionCtxData"),
		),
		jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Commentf("fetch %s from database.", scn),
		jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dotf("%sDataManager", uvn).Dotf("Get%s", sn).Call(
			dbCallArgs...,
		),
		jen.If(jen.Qual("errors", "Is").Call(
			jen.ID("err"),
			jen.Qual("database/sql", "ErrNoRows"),
		)).Body(
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("retrieving %s", scn),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Comment("encode our response and peace."),
		jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
			jen.ID("ctx"),
			jen.ID("res"),
			jen.ID("x"),
		),
	)

	lines := []jen.Code{
		jen.Commentf("ReadHandler returns a GET handler that returns %s.", scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildDBClientExistenceMethodCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToUserAtSomeLevel(proj) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildExistenceHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	dbCallArgs := buildDBClientExistenceMethodCallArgs(proj, typ)

	bodyLines := []jen.Code{jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.ID("s").Dot("logger").Dot("Error").Call(
				jen.ID("err"),
				jen.Lit("retrieving session context data"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("unauthenticated"),
				jen.Qual("net/http", "StatusUnauthorized"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
			jen.ID("span"),
			jen.ID("sessionCtxData"),
		),
		jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Comment("check the database."),
		jen.List(jen.ID("exists"), jen.ID("err")).Op(":=").ID("s").Dotf("%sDataManager", uvn).Dotf("%sExists", sn).Call(
			dbCallArgs...,
		),
		jen.If(jen.Op("!").Qual("errors", "Is").Call(
			jen.ID("err"),
			jen.Qual("database/sql", "ErrNoRows"),
		)).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("checking %s existence", scn),
			)),
		jen.Newline(),
		jen.If(jen.Op("!").ID("exists").Op("||").Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
		),
	)

	lines := []jen.Code{
		jen.Commentf("ExistenceHandler returns a HEAD handler that returns 200 if %s exists, 404 otherwise.", scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ExistenceHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildDBClientListRetrievalMethodCallArgs(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func buildListHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	dbCallArgs := buildDBClientListRetrievalMethodCallArgs(proj, typ)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("filter").Op(":=").Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).
			Dotln("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "FilterLimitKey"),
			jen.ID("filter").Dot("Limit"),
		).
			Dotln("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "FilterPageKey"),
			jen.ID("filter").Dot("Page"),
		).
			Dotln("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "FilterSortByKey"),
			jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(
			jen.ID("span"),
			jen.ID("filter").Dot("Page"),
			jen.ID("filter").Dot("Limit"),
			jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("retrieving session context data"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("unauthenticated"),
				jen.Qual("net/http", "StatusUnauthorized"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
			jen.ID("span"),
			jen.ID("sessionCtxData"),
		),
		jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.List(jen.ID(puvn), jen.ID("err")).Op(":=").ID("s").Dotf("%sDataManager", uvn).Dotf("Get%s", pn).Call(
			dbCallArgs...,
		),
		jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
			jen.Comment("in the event no rows exist, return an empty list."),
			jen.ID(puvn).Op("=").Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Values(jen.ID(pn).Op(":").Index().Op("*").Qual(proj.TypesPackage(), sn).Values())).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("retrieving %s", pcn),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Comment("encode our response and peace."),
		jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
			jen.ID("ctx"),
			jen.ID("res"),
			jen.ID(puvn),
		),
	)

	lines := []jen.Code{
		jen.Comment("ListHandler is our list route."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildUpdateSomethingDBCallArgs(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToUserAtSomeLevel(p) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildUpdateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	rn := typ.Name.RouteName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("retrieving session context data"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("unauthenticated"),
				jen.Qual("net/http", "StatusUnauthorized"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
			jen.ID("span"),
			jen.ID("sessionCtxData"),
		),
		jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
		jen.Comment("check for parsed input attached to session context data."),
		jen.ID("input").Op(":=").ID("new").Call(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn))),
		jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.ID("input"),
		), jen.ID("err").Op("!=").ID("nil")).Body(
			jen.ID("logger").Dot("Error").Call(
				jen.ID("err"),
				jen.Lit("error encountered decoding request body"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("invalid request content"),
				jen.Qual("net/http", "StatusBadRequest"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Op("=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
			jen.ID("logger").Dot("Error").Call(
				jen.ID("err"),
				jen.Lit("provided input was invalid"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("err").Dot("Error").Call(),
				jen.Qual("net/http", "StatusBadRequest"),
			),
			jen.Return(),
		),
		func() jen.Code {
			if typ.BelongsToAccount {
				return jen.ID("input").Dot("BelongsToAccount").Op("=").ID("sessionCtxData").Dot("ActiveAccountID")
			}
			return jen.Null()
		}(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, true)...)

	dbCallArgs := buildUpdateSomethingDBCallArgs(proj, typ)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("fetch %s from database.", scn),
		jen.List(jen.ID(uvn), jen.ID("err")).Op(":=").ID("s").Dotf("%sDataManager", uvn).Dotf("Get%s", sn).Call(
			dbCallArgs...,
		),
		jen.If(jen.Qual("errors", "Is").Call(
			jen.ID("err"),
			jen.Qual("database/sql", "ErrNoRows"),
		)).Body(
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("retrieving %s for update", scn),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Commentf("update the %s.", scn),
		jen.ID("changeReport").Op(":=").ID(uvn).Dot("Update").Call(jen.ID("input")),
		jen.Qual(proj.InternalTracingPackage(), "AttachChangeSummarySpan").Call(
			jen.ID("span"),
			jen.Lit(rn),
			jen.ID("changeReport"),
		),
		jen.Newline(),
		jen.Commentf("update %s in database.", scn),
		jen.If(jen.ID("err").Op("=").ID("s").Dotf("%sDataManager", uvn).Dotf("Update%s", sn).Call(
			jen.ID("ctx"),
			jen.ID(uvn),
			jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("changeReport"),
		), jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("updating %s", scn),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Comment("notify interested parties."),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.If(jen.ID("searchIndexErr").Op(":=").ID("s").Dot("search").Dot("Index").Call(
					jen.ID("ctx"),
					jen.ID(uvn).Dot("ID"),
					jen.ID(uvn),
				), jen.ID("searchIndexErr").Op("!=").ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("updating %s in search index", scn),
					),
				)
			}
			return jen.Null()
		}(),
		jen.Newline(),
		jen.Comment("encode our response and peace."),
		jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
			jen.ID("ctx"),
			jen.ID("res"),
			jen.ID(uvn),
		),
	)

	lines := []jen.Code{
		jen.Commentf("UpdateHandler returns a handler that updates %s.", scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("UpdateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildSearchHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	scn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("query").Op(":=").ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.Qual(proj.TypesPackage(), "SearchQueryKey")),
		jen.ID("filter").Op(":=").Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).
			Dotln("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "FilterLimitKey"),
			jen.ID("filter").Dot("Limit"),
		).
			Dotln("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "FilterPageKey"),
			jen.ID("filter").Dot("Page"),
		).
			Dotln("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "FilterSortByKey"),
			jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
		).
			Dotln("WithValue").Call(
			jen.Qual(proj.ObservabilityPackage("keys"), "SearchQueryKey"),
			jen.ID("query"),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(
			jen.ID("span"),
			jen.ID("filter").Dot("Page"),
			jen.ID("filter").Dot("Limit"),
			jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.ID("s").Dot("logger").Dot("Error").Call(
				jen.ID("err"),
				jen.Lit("retrieving session context data"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("unauthenticated"),
				jen.Qual("net/http", "StatusUnauthorized"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
			jen.ID("span"),
			jen.ID("sessionCtxData"),
		),
		jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
		jen.List(jen.ID("relevantIDs"), jen.ID("err")).Op(":=").ID("s").Dot("search").Dot("Search").Call(
			jen.ID("ctx"),
			jen.ID("query"),
			jen.ID("sessionCtxData").Dot("ActiveAccountID"),
		),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("executing %s search query", scn),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Commentf("fetch %s from database.", pcn),
		jen.List(jen.ID(puvn), jen.ID("err")).Op(":=").ID("s").Dotf("%sDataManager", uvn).Dotf("Get%sWithIDs", pn).Call(
			jen.ID("ctx"),
			jen.ID("sessionCtxData").Dot("ActiveAccountID"),
			jen.ID("filter").Dot("Limit"),
			jen.ID("relevantIDs"),
		),
		jen.If(jen.Qual("errors", "Is").Call(
			jen.ID("err"),
			jen.Qual("database/sql", "ErrNoRows"),
		)).Body(
			jen.Comment("in the event no rows exist, return an empty list."),
			jen.ID(puvn).Op("=").Index().Op("*").Qual(proj.TypesPackage(), sn).Values()).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("searching %s", pcn),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Comment("encode our response and peace."),
		jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
			jen.ID("ctx"),
			jen.ID("res"),
			jen.ID(puvn),
		),
	}

	lines := []jen.Code{
		jen.Comment("SearchHandler is our search route."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("SearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildCreateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Commentf("CreateHandler is our %s creation route.", scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
			jen.Newline(),
			jen.Comment("check session context data for parsed input struct."),
			jen.ID("input").Op(":=").ID("new").Call(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn))),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("ctx"), jen.ID("req"), jen.ID("input")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("invalid request content"),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.If(jen.ID("err").Op("=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "ValidationErrorKey"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("invalid input attached to request")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.Newline(),
			func() jen.Code {
				if typ.BelongsToAccount {
					return jen.ID("input").Dot("BelongsToAccount").Op("=").ID("sessionCtxData").Dot("ActiveAccountID")
				}
				return jen.Null()
			}(),
			jen.Newline(),
			jen.Commentf("create %s in database.", scn),
			jen.List(jen.ID(uvn), jen.ID("err")).Op(":=").ID("s").Dotf("%sDataManager", uvn).Dotf("Create%s", sn).Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("creating %s", scn),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.ID(uvn).Dot("ID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.ID(uvn).Dot("ID"),
			),
			jen.Newline(),
			jen.Comment("notify interested parties."),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.If(jen.ID("searchIndexErr").Op(":=").ID("s").Dot("search").Dot("Index").Call(jen.ID("ctx"), jen.ID(uvn).Dot("ID"), jen.ID(uvn)), jen.ID("searchIndexErr").Op("!=").ID("nil")).Body(
						jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(jen.ID("err"), jen.ID("logger"), jen.ID("span"), jen.Litf("adding %s to search index", scn)),
					)
				}
				return jen.Null()
			}(),
			jen.ID("s").Dotf("%sCounter", uvn).Dot("Increment").Call(jen.ID("ctx")),
			jen.Newline(),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID(uvn),
				jen.Qual("net/http", "StatusCreated"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildArchiveSomethingArgs(typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.BelongsToAccount {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	params = append(params, jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"))

	return params
}

func buildArchiveHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("retrieving session context data"),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("unauthenticated"),
				jen.Qual("net/http", "StatusUnauthorized"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
			jen.ID("span"),
			jen.ID("sessionCtxData"),
		),
		jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, true)...)

	dbCallArgs := buildArchiveSomethingArgs(typ)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Commentf("archive the %s in the database.", scn),
		jen.ID("err").Op("=").ID("s").Dotf("%sDataManager", uvn).Dotf("Archive%s", sn).Call(
			dbCallArgs...,
		),
		jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("archiving %s", scn),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return(),
		),
		jen.Newline(),
		jen.Comment("notify interested parties."),
		jen.ID("s").Dotf("%sCounter", uvn).Dot("Decrement").Call(jen.ID("ctx")),
		jen.Newline(),
		func() jen.Code {
			if typ.SearchEnabled {
				return jen.If(jen.ID("indexDeleteErr").Op(":=").ID("s").Dot("search").Dot("Delete").Call(
					jen.ID("ctx"),
					jen.IDf("%sID", uvn),
				), jen.ID("indexDeleteErr").Op("!=").ID("nil")).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("removing from search index"),
					),
				)
			}
			return jen.Null()
		}(),
		jen.Newline(),
		jen.Comment("encode our response and peace."),
		jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNoContent")),
	)

	lines := []jen.Code{
		jen.Commentf("ArchiveHandler returns a handler that archives %s.", scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildAuditEntryHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Commentf("AuditEntryHandler returns a GET handler that returns all audit log entries related to %s.", scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("unauthenticated"),
					jen.Qual("net/http", "StatusUnauthorized"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.ID("logger").Op("=").ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
			jen.Newline(),
			jen.Commentf("determine %s ID.", scn),
			jen.IDf("%sID", uvn).Op(":=").ID("s").Dotf("%sIDFetcher", uvn).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("s").Dotf("%sDataManager", uvn).Dotf("GetAuditLogEntriesFor%s", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("retrieving audit log entries for %s", scn),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("x"),
			),
		),
		jen.Newline(),
	}

	return lines
}
