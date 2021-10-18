package iterables

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
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
			jen.IDf("%sID", uvn).Assign().ID("s").Dotf("%sIDFetcher", uvn).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
		)
	}

	return idFetches
}

func buildDBClientRetrievalMethodCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildDBClientExistenceMethodCallArgs(proj *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := proj.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildDBClientListRetrievalMethodCallArgs(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToAccountAtSomeLevel(p) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}
	params = append(params, jen.ID("filter"))

	return params
}

func buildUpdateSomethingDBCallArgs(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}
	uvn := typ.Name.UnexportedVarName()

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		params = append(params, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	params = append(params, jen.IDf("%sID", uvn))

	if typ.RestrictedToAccountAtSomeLevel(p) {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	return params
}

func buildDBClientSearchMethodCallArgs(_ *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}

	if typ.BelongsToAccount {
		params = append(params, jen.ID("sessionCtxData").Dot("ActiveAccountID"))
	}

	params = append(params, jen.ID("filter").Dot("Limit"), jen.ID("relevantIDs"))

	return params
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

	return params
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
			jen.IDf("%sIDURIParamKey", sn).Equals().Litf("%sID", uvn),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("parseBool differs from strconv.ParseBool in that it returns false by default."),
		jen.Newline(),
		jen.Func().ID("parseBool").Params(jen.ID("str").String()).Params(jen.ID("bool")).Body(
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("str")))).Body(
				jen.Case(jen.Lit("1"), jen.Lit("t"), jen.Lit("true")).Body(
					jen.Return().True()),
				jen.Default().Body(
					jen.Return().False()),
			)),
		jen.Newline(),
	)

	code.Add(buildCreateHandler(proj, typ)...)
	code.Add(buildReadHandler(proj, typ)...)
	code.Add(buildListHandler(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildSearchHandler(proj, typ)...)
	}

	code.Add(buildUpdateHandler(proj, typ)...)
	code.Add(buildArchiveHandler(proj, typ)...)

	return code
}

func buildCreateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Commentf("CreateHandler is our %s creation route.", scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
			jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
			jen.Newline(),
			jen.Comment("check session context data for parsed input struct."),
			jen.ID("providedInput").Assign().ID("new").Call(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", sn))),
			jen.If(jen.ID("err").Equals().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("ctx"), jen.ID("req"), jen.ID("providedInput")), jen.ID("err").DoesNotEqual().Nil()).Body(
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
			jen.If(jen.ID("err").Equals().ID("providedInput").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Qual(proj.ObservabilityPackage("keys"), "ValidationErrorKey"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("provided input was invalid")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.ID("input").Assign().Qualf(proj.TypesPackage(), "%sDatabaseCreationInputFrom%sCreationInput", sn, sn).Call(jen.ID("providedInput")),
			jen.ID("input").Dot("ID").Equals().Qual(constants.IDGenerationLibrary, "New").Call().Dot("String").Call(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.Commentf("determine %s ID.", typ.BelongsToStruct.SingularCommonName()).Newline().
						IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Assign().ID("s").Dotf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Call(jen.ID("req")).Newline().
						Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).Newline().
						ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.BelongsToStruct.Singular())), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())).Newline().Newline().
						ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()).Equals().IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToAccount {
					return jen.ID("input").Dot("BelongsToAccount").Equals().ID("sessionCtxData").Dot("ActiveAccountID")
				}
				return jen.Null()
			}(),
			jen.Qualf(proj.InternalTracingPackage(), "Attach%sIDToSpan", sn).Call(jen.ID(constants.SpanVarName), jen.ID("input").Dot("ID")),
			jen.Newline(),
			jen.Commentf("create %s in database.", scn),
			jen.If(jen.ID("s").Dot("async")).Body(
				jen.ID("preWrite").Assign().AddressOf().Qual(proj.TypesPackage(), "PreWriteMessage").Valuesln(
					jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn),
					jen.ID(sn).MapAssign().ID("input"),
					jen.ID("AttributableToUserID").MapAssign().ID("sessionCtxData").Dot("Requester").Dot("UserID"),
					jen.ID("AttributableToAccountID").MapAssign().ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.If(jen.Err().Equals().ID("s").Dot("preWritesPublisher").Dot("Publish").Call(constants.CtxVar(), jen.ID("preWrite")), jen.Err().DoesNotEqual().Nil()).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(jen.Err(), constants.LoggerVar(), jen.ID(constants.SpanVarName), jen.Litf("publishing %s write message", scn)),
					jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(constants.CtxVar(), jen.ID(constants.ResponseVarName)),
					jen.Return(),
				),
				jen.Newline(),
				jen.ID("pwr").Assign().Qual(proj.TypesPackage(), "PreWriteResponse").Values(jen.ID("ID").MapAssign().ID("input").Dot("ID")),
				jen.Newline(),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(constants.CtxVar(), jen.ID(constants.ResponseVarName), jen.ID("pwr"), jen.Qual("net/http", "StatusAccepted")),
			).Else().Body(
				jen.List(jen.ID(uvn), jen.IDf("%sCreationErr", uvn)).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("Create%s", sn).Call(
					jen.ID("ctx"),
					jen.ID("input"),
				),
				jen.If(jen.IDf("%sCreationErr", uvn).DoesNotEqual().Nil()).Body(
					jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
						jen.IDf("%sCreationErr", uvn),
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
				utils.ConditionalCode(typ.SearchEnabled,
					jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
						jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
						jen.ID(uvn).Dot("ID"),
					),
				),
				jen.Newline(),
				jen.Comment("notify interested parties."),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.If(jen.ID("searchIndexErr").Assign().ID("s").Dot("search").Dot("Index").Call(jen.ID("ctx"), jen.ID(uvn).Dot("ID"), jen.ID(uvn)), jen.ID("searchIndexErr").DoesNotEqual().Nil()).Body(
							jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(jen.ID("err"), jen.ID("logger"), jen.ID("span"), jen.Litf("adding %s to search index", scn)),
						)
					}
					return jen.Null()
				}(),
				jen.Newline(),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID(uvn),
					jen.Qual("net/http", "StatusCreated"),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildExistenceHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	dbCallArgs := buildDBClientExistenceMethodCallArgs(proj, typ)

	bodyLines := []jen.Code{jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Comment("check the database."),
		jen.List(jen.ID("exists"), jen.ID("err")).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("%sExists", sn).Call(
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
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("ExistenceHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildReadHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	dbCallArgs := buildDBClientRetrievalMethodCallArgs(proj, typ)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, true)...)

	bodyLines = append(bodyLines,
		jen.Commentf("fetch %s from database.", scn),
		jen.List(jen.ID("x"), jen.ID("err")).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("Get%s", sn).Call(
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
		).Else().If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildListHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	dbCallArgs := buildDBClientListRetrievalMethodCallArgs(proj, typ)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("filter").Assign().Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).
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
			jen.String().Call(jen.ID("filter").Dot("SortBy")),
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
			jen.String().Call(jen.ID("filter").Dot("SortBy")),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, false)...)

	bodyLines = append(bodyLines,
		jen.List(jen.ID(puvn), jen.ID("err")).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("Get%s", pn).Call(
			dbCallArgs...,
		),
		jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
			jen.Comment("in the event no rows exist, return an empty list."),
			jen.ID(puvn).Equals().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Values(jen.ID(pn).MapAssign().Index().PointerTo().Qual(proj.TypesPackage(), sn).Values())).Else().If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildUpdateHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
		jen.Comment("check for parsed input attached to session context data."),
		jen.ID("input").Assign().ID("new").Call(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn))),
		jen.If(jen.ID("err").Equals().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
			jen.ID("ctx"),
			jen.ID("req"),
			jen.ID("input"),
		), jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.If(jen.ID("err").Equals().ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").DoesNotEqual().Nil()).Body(
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
				return jen.ID("input").Dot("BelongsToAccount").Equals().ID("sessionCtxData").Dot("ActiveAccountID")
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
		jen.List(jen.ID(uvn), jen.ID("err")).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("Get%s", sn).Call(
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
		).Else().If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.ID(uvn).Dot("Update").Call(jen.ID("input")),
		jen.Newline(),
		jen.If(jen.ID("s").Dot("async")).Body(
			jen.ID("pum").Assign().AddressOf().Qual(proj.TypesPackage(), "PreUpdateMessage").Valuesln(
				jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn),
				jen.ID(sn).MapAssign().ID(uvn),
				jen.ID("AttributableToUserID").MapAssign().ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("AttributableToAccountID").MapAssign().ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.If(jen.Err().Equals().ID("s").Dot("preUpdatesPublisher").Dot("Publish").Call(constants.CtxVar(), jen.ID("pum")), jen.Err().DoesNotEqual().Nil()).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(jen.Err(), constants.LoggerVar(), jen.ID(constants.SpanVarName), jen.Litf("publishing %s update message", scn)),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(constants.CtxVar(), jen.ID(constants.ResponseVarName)),
				jen.Return(),
			),
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID("s").Dot("encoderDecoder").Dot("RespondWithData").Call(constants.CtxVar(), jen.ID(constants.ResponseVarName), jen.ID(uvn)),
		).Else().Body(
			jen.Commentf("update %s in database.", scn),
			jen.If(jen.ID("err").Equals().ID("s").Dotf("%sDataManager", uvn).Dotf("Update%s", sn).Call(
				jen.ID("ctx"),
				jen.ID(uvn),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
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
					return jen.If(jen.ID("searchIndexErr").Assign().ID("s").Dot("search").Dot("Index").Call(
						jen.ID("ctx"),
						jen.ID(uvn).Dot("ID"),
						jen.ID(uvn),
					), jen.ID("searchIndexErr").DoesNotEqual().Nil()).Body(
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
		),
	)

	lines := []jen.Code{
		jen.Commentf("UpdateHandler returns a handler that updates %s.", scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("UpdateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("query").Assign().ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.Qual(proj.TypesPackage(), "SearchQueryKey")),
		jen.ID("filter").Assign().Qual(proj.TypesPackage(), "ExtractQueryFilter").Call(jen.ID("req")),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).
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
			jen.String().Call(jen.ID("filter").Dot("SortBy")),
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
			jen.String().Call(jen.ID("filter").Dot("SortBy")),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
		jen.List(jen.ID("relevantIDs"), jen.ID("err")).Assign().ID("s").Dot("search").Dot("Search").Call(
			jen.ID("ctx"),
			jen.ID("query"),
			jen.ID("sessionCtxData").Dot("ActiveAccountID"),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, false)...)

	dbCallArgs := buildDBClientSearchMethodCallArgs(proj, typ)

	bodyLines = append(bodyLines,
		jen.Commentf("fetch %s from database.", pcn),
		jen.List(jen.ID(puvn), jen.ID("err")).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("Get%sWithIDs", pn).Call(
			dbCallArgs...,
		),
		jen.If(jen.Qual("errors", "Is").Call(
			jen.ID("err"),
			jen.Qual("database/sql", "ErrNoRows"),
		)).Body(
			jen.Comment("in the event no rows exist, return an empty list."),
			jen.ID(puvn).Equals().Index().PointerTo().Qual(proj.TypesPackage(), sn).Values()).Else().If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
	)

	lines := []jen.Code{
		jen.Comment("SearchHandler is our search route."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("SearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildArchiveHandler(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
		jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
			jen.ID("span"),
			jen.ID("req"),
		),
		jen.Newline(),
		jen.Comment("determine user ID."),
		jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDFetchers(proj, typ, true)...)

	dbCallArgs := buildArchiveSomethingArgs(typ)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.If(jen.ID("s").Dot("async")).Body(
			jen.List(jen.ID("exists"), jen.ID("existenceCheckErr")).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("%sExists", sn).Call(
				typ.BuildDBClientExistenceMethodCallArgs(proj)...,
			),
			jen.If(jen.ID("existenceCheckErr").DoesNotEqual().Nil().And().Not().Qual("errors", "Is").Call(jen.ID("existenceCheckErr"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.ID("existenceCheckErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("checking %s existence", scn),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.Not().ID("exists").Or().Qual("errors", "Is").Call(jen.ID("existenceCheckErr"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.ID("pam").Assign().AddressOf().Qual(proj.TypesPackage(), "PreArchiveMessage").Valuesln(
				jen.ID("DataType").MapAssign().Qualf(proj.TypesPackage(), "%sDataType", sn),
				jen.IDf("%sID", sn).MapAssign().IDf("%sID", uvn),
				jen.ID("AttributableToUserID").MapAssign().ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				jen.ID("AttributableToAccountID").MapAssign().ID("sessionCtxData").Dot("ActiveAccountID"),
			),
			jen.If(jen.Err().Equals().ID("s").Dot("preArchivesPublisher").Dot("Publish").Call(constants.CtxVar(), jen.ID("pam")), jen.Err().DoesNotEqual().Nil()).Body(
				jen.Qual(proj.ObservabilityPackage(), "AcknowledgeError").Call(
					jen.Err(),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("publishing %s archive message", scn),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.Newline(),
			jen.Comment("encode our response and peace."),
			jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNoContent")),
		).Else().Body(
			jen.Commentf("archive the %s in the database.", scn),
			jen.ID("err").Equals().ID("s").Dotf("%sDataManager", uvn).Dotf("Archive%s", sn).Call(
				dbCallArgs...,
			),
			jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.If(jen.ID("indexDeleteErr").Assign().ID("s").Dot("search").Dot("Delete").Call(
						jen.ID("ctx"),
						jen.IDf("%sID", uvn),
					), jen.ID("indexDeleteErr").DoesNotEqual().Nil()).Body(
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
		),
	)

	lines := []jen.Code{
		jen.Commentf("ArchiveHandler returns a handler that archives %s.", scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
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
		jen.Func().Params(jen.ID("s").PointerTo().ID("service")).ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.Newline(),
			jen.Comment("determine user ID."),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Assign().ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
			jen.ID("logger").Equals().ID("sessionCtxData").Dot("AttachToLogger").Call(jen.ID("logger")),
			jen.Newline(),
			jen.Commentf("determine %s ID.", scn),
			jen.IDf("%sID", uvn).Assign().ID("s").Dotf("%sIDFetcher", uvn).Call(jen.ID("req")),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
				jen.Qual(proj.ObservabilityPackage("keys"), fmt.Sprintf("%sIDKey", sn)),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("x"), jen.ID("err")).Assign().ID("s").Dotf("%sDataManager", uvn).Dotf("GetAuditLogEntriesFor%s", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual("database/sql", "ErrNoRows"))).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			).Else().If(jen.ID("err").DoesNotEqual().Nil()).Body(
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
