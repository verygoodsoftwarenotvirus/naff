package tracing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanAttachersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("attachUint8ToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("attachmentKey").String(),
			jen.ID("id").ID("uint8")).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Int64").Call(
					jen.ID("attachmentKey"),
					jen.ID("int64").Call(jen.ID("id")),
				)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("attachUint64ToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("attachmentKey").String(),
			jen.ID("id").ID("uint64")).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Int64").Call(
					jen.ID("attachmentKey"),
					jen.ID("int64").Call(jen.ID("id")),
				)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("attachStringToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.List(jen.ID("key"),
				jen.ID("str")).String()).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "String").Call(
					jen.ID("key"),
					jen.ID("str"),
				)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("attachBooleanToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("key").String(),
			jen.ID("b").ID("bool")).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Bool").Call(
					jen.ID("key"),
					jen.ID("b"),
				)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("attachSliceToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("key").String(),
			jen.ID("slice").Interface()).Body(
			jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Array").Call(
				jen.ID("key"),
				jen.ID("slice"),
			)),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachToSpan allows a user to attach any value to a span."),
		jen.Newline(),
		jen.Func().ID("AttachToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("key").String(),
			jen.ID("val").Interface()).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Any").Call(
					jen.ID("key"),
					jen.ID("val"),
				)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachFilterToSpan provides a consistent way to attach a filter's info to a span."),
		jen.Newline(),
		jen.Func().ID("AttachFilterToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("page").Uint64(),
			jen.ID("limit").ID("uint8"),
			jen.ID("sortBy").String()).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "FilterPageKey"),
				jen.ID("page"),
			),
			jen.ID("attachUint8ToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "FilterLimitKey"),
				jen.ID("limit"),
			),
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "FilterSortByKey"),
				jen.ID("sortBy"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachAccountIDToSpan provides a consistent way to attach an account's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachAccountIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("accountID").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "AccountIDKey"),
				jen.ID("accountID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachActiveAccountIDToSpan provides a consistent way to attach an account's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachActiveAccountIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("accountID").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "ActiveAccountIDKey"),
				jen.ID("accountID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachRequestingUserIDToSpan provides a consistent way to attach a user's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachRequestingUserIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("userID").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "RequesterIDKey"),
				jen.ID("userID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachSessionContextDataToSpan provides a consistent way to attach a SessionContextData object to a span."),
		jen.Newline(),
		jen.Func().ID("AttachSessionContextDataToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("sessionCtxData").PointerTo().Qual(proj.TypesPackage(), "SessionContextData")).Body(
			jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
				jen.ID("AttachRequestingUserIDToSpan").Call(
					jen.ID("span"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				),
				jen.ID("AttachActiveAccountIDToSpan").Call(
					jen.ID("span"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.If(jen.ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").DoesNotEqual().Nil()).Body(
					jen.ID("attachBooleanToSpan").Call(
						jen.ID("span"),
						jen.Qual(proj.ConstantKeysPackage(), "UserIsServiceAdminKey"),
						jen.ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachAPIClientDatabaseIDToSpan is a consistent way to attach an API client's database row ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachAPIClientDatabaseIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("clientID").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "APIClientDatabaseIDKey"),
				jen.ID("clientID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachAPIClientClientIDToSpan is a consistent way to attach an API client's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachAPIClientClientIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("clientID").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "APIClientClientIDKey"),
				jen.ID("clientID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachUserToSpan provides a consistent way to attach a user to a span."),
		jen.Newline(),
		jen.Func().ID("AttachUserToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("user").PointerTo().Qual(proj.TypesPackage(), "User")).Body(
			jen.If(jen.ID("user").Op("!=").ID("nil")).Body(
				jen.ID("AttachUserIDToSpan").Call(
					jen.ID("span"),
					jen.ID("user").Dot("ID"),
				),
				jen.ID("AttachUsernameToSpan").Call(
					jen.ID("span"),
					jen.ID("user").Dot("Username"),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachUserIDToSpan provides a consistent way to attach a user's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachUserIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("userID").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "UserIDKey"),
				jen.ID("userID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachUsernameToSpan provides a consistent way to attach a user's username to a span."),
		jen.Newline(),
		jen.Func().ID("AttachUsernameToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("username").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "UsernameKey"),
				jen.ID("username"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachWebhookIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("webhookID").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "WebhookIDKey"),
				jen.ID("webhookID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachURLToSpan attaches a given URI to a span."),
		jen.Newline(),
		jen.Func().ID("AttachURLToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("u").PointerTo().Qual("net/url", "URL")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "RequestURIKey"),
				jen.ID("u").Dot("String").Call(),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachRequestURIToSpan attaches a given URI to a span."),
		jen.Newline(),
		jen.Func().ID("AttachRequestURIToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("uri").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "RequestURIKey"),
				jen.ID("uri"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachRequestToSpan attaches a given *http.Request to a span."),
		jen.Newline(),
		jen.Func().ID("AttachRequestToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("req").PointerTo().Qual("net/http", "Request")).Body(
			jen.If(jen.ID("req").Op("!=").ID("nil")).Body(
				jen.ID("attachStringToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "RequestURIKey"),
					jen.ID("req").Dot("URL").Dot("String").Call(),
				),
				jen.ID("attachStringToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "RequestMethodKey"),
					jen.ID("req").Dot("Method"),
				),
				jen.Newline(),
				jen.For(jen.List(jen.ID("k"),
					jen.ID("v")).Assign().Range().ID("req").Dot("Header")).Body(
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.Qual(proj.ConstantKeysPackage(), "RequestHeadersKey"),
							jen.ID("k"),
						),
						jen.ID("v"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachResponseToSpan attaches a given *http.Response to a span."),
		jen.Newline(),
		jen.Func().ID("AttachResponseToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("res").PointerTo().Qual("net/http", "Response")).Body(
			jen.If(jen.ID("res").Op("!=").ID("nil")).Body(
				jen.ID("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("res").Dot("Request"),
				),
				jen.Newline(),
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Int").Call(
					jen.Qual(proj.ConstantKeysPackage(), "ResponseStatusKey"),
					jen.ID("res").Dot("StatusCode"),
				)),
				jen.Newline(),
				jen.For(jen.List(jen.ID("k"),
					jen.ID("v")).Assign().Range().ID("res").Dot("Header")).Body(
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.Qual(proj.ConstantKeysPackage(), "ResponseHeadersKey"),
							jen.ID("k"),
						),
						jen.ID("v"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachErrorToSpan attaches a given error to a span."),
		jen.Newline(),
		jen.Func().ID("AttachErrorToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("description").String(),
			jen.ID("err").ID("error")).Body(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("RecordError").Callln(
					jen.ID("err"),
					jen.Qual(constants.TracingLibrary, "WithTimestamp").Call(jen.Qual("time", "Now").Call()),
					jen.Qual(constants.TracingLibrary, "WithAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "String").Call(
						jen.Lit("error.description"),
						jen.ID("description"),
					),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachDatabaseQueryToSpan attaches a given search query to a span."),
		jen.Newline(),
		jen.Func().ID("AttachDatabaseQueryToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.List(jen.ID("queryDescription"),
				jen.ID("query")).String(),
			jen.ID("args").Index().Interface()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "DatabaseQueryKey"),
				jen.ID("query"),
			),
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Lit("query_description"),
				jen.ID("queryDescription"),
			),
			jen.Newline(),
			jen.For(jen.List(jen.ID("i"),
				jen.ID("arg")).Assign().Range().ID("args")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Any").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("query_args_%d"),
						jen.ID("i"),
					),
					jen.ID("arg"),
				),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachQueryFilterToSpan attaches a given query filter to a span."),
		jen.Newline(),
		jen.Func().ID("AttachQueryFilterToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("filter").PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Body(
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("attachUint8ToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterLimitKey"),
					jen.ID("filter").Dot("Limit"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterPageKey"),
					jen.ID("filter").Dot("Page"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterCreatedAfterKey"),
					jen.ID("filter").Dot("CreatedAfter"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterCreatedBeforeKey"),
					jen.ID("filter").Dot("CreatedBefore"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterUpdatedAfterKey"),
					jen.ID("filter").Dot("UpdatedAfter"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterUpdatedBeforeKey"),
					jen.ID("filter").Dot("UpdatedBefore"),
				),
				jen.ID("attachStringToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterSortByKey"),
					jen.String().Call(jen.ID("filter").Dot("SortBy")),
				),
			).Else().Body(
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "FilterIsNilKey"),
					jen.ID("true"),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachSearchQueryToSpan attaches a given search query to a span."),
		jen.Newline(),
		jen.Func().ID("AttachSearchQueryToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("query").String()).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.Qual(proj.ConstantKeysPackage(), "SearchQueryKey"),
				jen.ID("query"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachUserAgentDataToSpan attaches a given search query to a span."),
		jen.Newline(),
		jen.Func().ID("AttachUserAgentDataToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("ua").PointerTo().Qual("github.com/mssola/user_agent", "UserAgent")).Body(
			jen.If(jen.ID("ua").Op("!=").ID("nil")).Body(
				jen.ID("attachStringToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "UserAgentOSKey"),
					jen.ID("ua").Dot("OS").Call(),
				),
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "UserAgentMobileKey"),
					jen.ID("ua").Dot("Mobile").Call(),
				),
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.Qual(proj.ConstantKeysPackage(), "UserAgentBotKey"),
					jen.ID("ua").Dot("Bot").Call(),
				),
			),
		),
		jen.Newline(),
	)

	for _, typ := range proj.DataTypes {
		code.Add(
			jen.Commentf("Attach%sIDToSpan attaches %s ID to a given span.", typ.Name.Singular(), typ.Name.SingularCommonNameWithPrefix()),
			jen.Newline(),
			jen.Func().IDf("Attach%sIDToSpan", typ.Name.Singular()).Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
				jen.IDf("%sID", typ.Name.UnexportedVarName()).String()).Body(
				jen.ID("attachStringToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dotf("%sIDKey", typ.Name.Singular()),
					jen.IDf("%sID", typ.Name.UnexportedVarName()),
				),
			),
			jen.Newline(),
		)
	}

	return code
}
