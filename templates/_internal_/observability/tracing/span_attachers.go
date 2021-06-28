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
				jen.ID("str")).ID("string")).Body(
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
			jen.ID("sortBy").ID("string")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("FilterPageKey"),
				jen.ID("page"),
			),
			jen.ID("attachUint8ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("FilterLimitKey"),
				jen.ID("limit"),
			),
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("FilterSortByKey"),
				jen.ID("sortBy"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachAuditLogEntryIDToSpan attaches an audit log entry ID to a given span."),
		jen.Newline(),
		jen.Func().ID("AttachAuditLogEntryIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("entryID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("AuditLogEntryIDKey"),
				jen.ID("entryID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachAuditLogEntryEventTypeToSpan attaches an audit log entry ID to a given span."),
		jen.Newline(),
		jen.Func().ID("AttachAuditLogEntryEventTypeToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("eventType").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("AuditLogEntryEventTypeKey"),
				jen.ID("eventType"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachAccountIDToSpan provides a consistent way to attach an account's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachAccountIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("accountID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachActiveAccountIDToSpan provides a consistent way to attach an account's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachActiveAccountIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("accountID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("ActiveAccountIDKey"),
				jen.ID("accountID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachRequestingUserIDToSpan provides a consistent way to attach a user's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachRequestingUserIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("userID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("userID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachChangeSummarySpan provides a consistent way to attach a SessionContextData object to a span."),
		jen.Newline(),
		jen.Func().ID("AttachChangeSummarySpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("typeName").String(),
			jen.ID("changes").Index().PointerTo().Qual(proj.TypesPackage(), "FieldChangeSummary")).Body(
			jen.For(jen.List(jen.ID("i"),
				jen.ID("change")).Assign().Range().ID("changes")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.Qual(constants.TracingAttributionLibrary, "Any").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.field_changes.%d"),
						jen.ID("typeName"),
						jen.ID("i"),
					),
					jen.ID("change"),
				)),
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
						jen.ID("keys").Dot("UserIsServiceAdminKey"),
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
			jen.ID("clientID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("APIClientDatabaseIDKey"),
				jen.ID("clientID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachAPIClientClientIDToSpan is a consistent way to attach an API client's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachAPIClientClientIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("clientID").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("APIClientClientIDKey"),
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
			jen.ID("userID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachUsernameToSpan provides a consistent way to attach a user's username to a span."),
		jen.Newline(),
		jen.Func().ID("AttachUsernameToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("username").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("username"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span."),
		jen.Newline(),
		jen.Func().ID("AttachWebhookIDToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("webhookID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("WebhookIDKey"),
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
				jen.ID("keys").Dot("RequestURIKey"),
				jen.ID("u").Dot("String").Call(),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("AttachRequestURIToSpan attaches a given URI to a span."),
		jen.Newline(),
		jen.Func().ID("AttachRequestURIToSpan").Params(jen.ID("span").Qual(constants.TracingLibrary, "Span"),
			jen.ID("uri").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("RequestURIKey"),
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
					jen.ID("keys").Dot("RequestURIKey"),
					jen.ID("req").Dot("URL").Dot("String").Call(),
				),
				jen.ID("attachStringToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("RequestMethodKey"),
					jen.ID("req").Dot("Method"),
				),
				jen.Newline(),
				jen.ID("htmxHeaderSpanKeys").Assign().Map(jen.ID("string")).String().Valuesln(jen.Lit("HX-Prompt").Op(":").Lit("htmx.prompt"),
					jen.Lit("HX-Target").Op(":").Lit("htmx.target"),
					jen.Lit("HX-Request").Op(":").Lit("htmx.request"),
					jen.Lit("HX-Trigger").Op(":").Lit("htmx.trigger"),
					jen.Lit("HX-Current-URL").Op(":").Lit("htmx.currentURL"),
					jen.Lit("HX-Trigger-LabelName").Op(":").Lit("htmx.triggerName"),
					jen.Lit("HX-History-Restore-Request").Op(":").Lit("htmx.historyRestoreRequest")),
				jen.Newline(),
				jen.For(jen.List(jen.ID("header"),
					jen.ID("spanKey")).Assign().Range().ID("htmxHeaderSpanKeys")).Body(
					jen.If(jen.ID("val").Assign().ID("req").Dot("Header").Dot("Get").Call(jen.ID("header")),
						jen.ID("val").Op("!=").Lit("")).Body(
						jen.ID("attachStringToSpan").Call(
							jen.ID("span"),
							jen.ID("spanKey"),
							jen.ID("val"),
						),
					),
				),
				jen.Newline(),
				jen.For(jen.List(jen.ID("k"),
					jen.ID("v")).Assign().Range().ID("req").Dot("Header")).Body(
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.ID("keys").Dot("RequestHeadersKey"),
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
					jen.ID("keys").Dot("ResponseStatusKey"),
					jen.ID("res").Dot("StatusCode"),
				)),
				jen.Newline(),
				jen.For(jen.List(jen.ID("k"),
					jen.ID("v")).Assign().Range().ID("res").Dot("Header")).Body(
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.ID("keys").Dot("ResponseHeadersKey"),
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
				jen.ID("keys").Dot("DatabaseQueryKey"),
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
					jen.ID("keys").Dot("FilterLimitKey"),
					jen.ID("filter").Dot("Limit"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterPageKey"),
					jen.ID("filter").Dot("Page"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterCreatedAfterKey"),
					jen.ID("filter").Dot("CreatedAfter"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterCreatedBeforeKey"),
					jen.ID("filter").Dot("CreatedBefore"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterUpdatedAfterKey"),
					jen.ID("filter").Dot("UpdatedAfter"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterUpdatedBeforeKey"),
					jen.ID("filter").Dot("UpdatedBefore"),
				),
				jen.ID("attachStringToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterSortByKey"),
					jen.String().Call(jen.ID("filter").Dot("SortBy")),
				),
			).Else().Body(
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterIsNilKey"),
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
			jen.ID("query").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("SearchQueryKey"),
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
					jen.ID("keys").Dot("UserAgentOSKey"),
					jen.ID("ua").Dot("OS").Call(),
				),
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("UserAgentMobileKey"),
					jen.ID("ua").Dot("Mobile").Call(),
				),
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("UserAgentBotKey"),
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
				jen.IDf("%sID", typ.Name.UnexportedVarName()).ID("uint64")).Body(
				jen.ID("attachUint64ToSpan").Call(
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
