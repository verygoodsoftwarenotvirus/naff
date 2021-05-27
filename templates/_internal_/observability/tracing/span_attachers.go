package tracing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanAttachersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("attachUint8ToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("attachmentKey").ID("string"), jen.ID("id").ID("uint8")).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Int64").Call(
					jen.ID("attachmentKey"),
					jen.ID("int64").Call(jen.ID("id")),
				)))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("attachUint64ToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("attachmentKey").ID("string"), jen.ID("id").ID("uint64")).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Int64").Call(
					jen.ID("attachmentKey"),
					jen.ID("int64").Call(jen.ID("id")),
				)))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("attachStringToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.List(jen.ID("key"), jen.ID("str")).ID("string")).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("String").Call(
					jen.ID("key"),
					jen.ID("str"),
				)))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("attachBooleanToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("key").ID("string"), jen.ID("b").ID("bool")).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Bool").Call(
					jen.ID("key"),
					jen.ID("b"),
				)))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("attachSliceToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("key").ID("string"), jen.ID("slice").Interface()).Body(
			jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Array").Call(
				jen.ID("key"),
				jen.ID("slice"),
			))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachToSpan allows a user to attach any value to a span."),
		jen.Line(),
		jen.Func().ID("AttachToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("key").ID("string"), jen.ID("val").Interface()).Body(
			jen.If(jen.ID("span").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Any").Call(
					jen.ID("key"),
					jen.ID("val"),
				)))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachFilterToSpan provides a consistent way to attach a filter's info to a span."),
		jen.Line(),
		jen.Func().ID("AttachFilterToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("page").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("sortBy").ID("string")).Body(
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
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachAuditLogEntryIDToSpan attaches an audit log entry ID to a given span."),
		jen.Line(),
		jen.Func().ID("AttachAuditLogEntryIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("entryID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("AuditLogEntryIDKey"),
				jen.ID("entryID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachAuditLogEntryEventTypeToSpan attaches an audit log entry ID to a given span."),
		jen.Line(),
		jen.Func().ID("AttachAuditLogEntryEventTypeToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("eventType").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("AuditLogEntryEventTypeKey"),
				jen.ID("eventType"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachAccountIDToSpan provides a consistent way to attach an account's ID to a span."),
		jen.Line(),
		jen.Func().ID("AttachAccountIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("accountID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachRequestingUserIDToSpan provides a consistent way to attach a user's ID to a span."),
		jen.Line(),
		jen.Func().ID("AttachRequestingUserIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachChangeSummarySpan provides a consistent way to attach a SessionContextData object to a span."),
		jen.Line(),
		jen.Func().ID("AttachChangeSummarySpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("typeName").ID("string"), jen.ID("changes").Index().Op("*").ID("types").Dot("FieldChangeSummary")).Body(
			jen.For(jen.List(jen.ID("i"), jen.ID("change")).Op(":=").Range().ID("changes")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Any").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.field_changes.%d"),
						jen.ID("typeName"),
						jen.ID("i"),
					),
					jen.ID("change"),
				)))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachSessionContextDataToSpan provides a consistent way to attach a SessionContextData object to a span."),
		jen.Line(),
		jen.Func().ID("AttachSessionContextDataToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData")).Body(
			jen.If(jen.ID("sessionCtxData").Op("!=").ID("nil")).Body(
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("RequesterIDKey"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("UserID"),
				),
				jen.ID("attachUint64ToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("ActiveAccountIDKey"),
					jen.ID("sessionCtxData").Dot("ActiveAccountID"),
				),
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("ServiceRoleKey"),
					jen.ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("IsServiceAdmin").Call(),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachAPIClientDatabaseIDToSpan is a consistent way to attach an API client's database row ID to a span."),
		jen.Line(),
		jen.Func().ID("AttachAPIClientDatabaseIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("clientID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("APIClientDatabaseIDKey"),
				jen.ID("clientID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachAPIClientClientIDToSpan is a consistent way to attach an API client's ID to a span."),
		jen.Line(),
		jen.Func().ID("AttachAPIClientClientIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("clientID").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("APIClientClientIDKey"),
				jen.ID("clientID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachUserToSpan provides a consistent way to attach a user to a span."),
		jen.Line(),
		jen.Func().ID("AttachUserToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("user").Op("*").ID("types").Dot("User")).Body(
			jen.If(jen.ID("user").Op("!=").ID("nil")).Body(
				jen.ID("AttachUserIDToSpan").Call(
					jen.ID("span"),
					jen.ID("user").Dot("ID"),
				),
				jen.ID("AttachUsernameToSpan").Call(
					jen.ID("span"),
					jen.ID("user").Dot("Username"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachUserIDToSpan provides a consistent way to attach a user's ID to a span."),
		jen.Line(),
		jen.Func().ID("AttachUserIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachUsernameToSpan provides a consistent way to attach a user's username to a span."),
		jen.Line(),
		jen.Func().ID("AttachUsernameToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("username").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("username"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span."),
		jen.Line(),
		jen.Func().ID("AttachWebhookIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("webhookID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("WebhookIDKey"),
				jen.ID("webhookID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachURLToSpan attaches a given URI to a span."),
		jen.Line(),
		jen.Func().ID("AttachURLToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("u").Op("*").Qual("net/url", "URL")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("RequestURIKey"),
				jen.ID("u").Dot("String").Call(),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachRequestURIToSpan attaches a given URI to a span."),
		jen.Line(),
		jen.Func().ID("AttachRequestURIToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("uri").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("RequestURIKey"),
				jen.ID("uri"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachRequestToSpan attaches a given *http.Request to a span."),
		jen.Line(),
		jen.Func().ID("AttachRequestToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
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
				jen.ID("htmxHeaderSpanKeys").Op(":=").Map(jen.ID("string")).ID("string").Valuesln(jen.Lit("HX-Prompt").Op(":").Lit("htmx.prompt"), jen.Lit("HX-Target").Op(":").Lit("htmx.target"), jen.Lit("HX-Request").Op(":").Lit("htmx.request"), jen.Lit("HX-Trigger").Op(":").Lit("htmx.trigger"), jen.Lit("HX-Current-URL").Op(":").Lit("htmx.currentURL"), jen.Lit("HX-Trigger-LabelName").Op(":").Lit("htmx.triggerName"), jen.Lit("HX-History-Restore-Request").Op(":").Lit("htmx.historyRestoreRequest")),
				jen.For(jen.List(jen.ID("header"), jen.ID("spanKey")).Op(":=").Range().ID("htmxHeaderSpanKeys")).Body(
					jen.If(jen.ID("val").Op(":=").ID("req").Dot("Header").Dot("Get").Call(jen.ID("header")), jen.ID("val").Op("!=").Lit("")).Body(
						jen.ID("attachStringToSpan").Call(
							jen.ID("span"),
							jen.ID("spanKey"),
							jen.ID("val"),
						))),
				jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("req").Dot("Header")).Body(
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.ID("keys").Dot("RequestHeadersKey"),
							jen.ID("k"),
						),
						jen.ID("v"),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachResponseToSpan attaches a given *http.Response to a span."),
		jen.Line(),
		jen.Func().ID("AttachResponseToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("res").Op("*").Qual("net/http", "Response")).Body(
			jen.If(jen.ID("res").Op("!=").ID("nil")).Body(
				jen.ID("AttachRequestToSpan").Call(
					jen.ID("span"),
					jen.ID("res").Dot("Request"),
				),
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Int").Call(
					jen.ID("keys").Dot("ResponseStatusKey"),
					jen.ID("res").Dot("StatusCode"),
				)),
				jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("res").Dot("Header")).Body(
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s.%s"),
							jen.ID("keys").Dot("ResponseHeadersKey"),
							jen.ID("k"),
						),
						jen.ID("v"),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachErrorToSpan attaches a given error to a span."),
		jen.Line(),
		jen.Func().ID("AttachErrorToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("description").ID("string"), jen.ID("err").ID("error")).Body(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("span").Dot("RecordError").Call(
					jen.ID("err"),
					jen.ID("trace").Dot("WithTimestamp").Call(jen.Qual("time", "Now").Call()),
					jen.ID("trace").Dot("WithAttributes").Call(jen.ID("attribute").Dot("String").Call(
						jen.Lit("error.description"),
						jen.ID("description"),
					)),
				))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachDatabaseQueryToSpan attaches a given search query to a span."),
		jen.Line(),
		jen.Func().ID("AttachDatabaseQueryToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.List(jen.ID("queryDescription"), jen.ID("query")).ID("string"), jen.ID("args").Index().Interface()).Body(
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
			jen.For(jen.List(jen.ID("i"), jen.ID("arg")).Op(":=").Range().ID("args")).Body(
				jen.ID("span").Dot("SetAttributes").Call(jen.ID("attribute").Dot("Any").Call(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("query_args_%d"),
						jen.ID("i"),
					),
					jen.ID("arg"),
				))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachQueryFilterToSpan attaches a given query filter to a span."),
		jen.Line(),
		jen.Func().ID("AttachQueryFilterToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Body(
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
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				),
			).Else().Body(
				jen.ID("attachBooleanToSpan").Call(
					jen.ID("span"),
					jen.ID("keys").Dot("FilterIsNilKey"),
					jen.ID("true"),
				))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachSearchQueryToSpan attaches a given search query to a span."),
		jen.Line(),
		jen.Func().ID("AttachSearchQueryToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("query").ID("string")).Body(
			jen.ID("attachStringToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("SearchQueryKey"),
				jen.ID("query"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachUserAgentDataToSpan attaches a given search query to a span."),
		jen.Line(),
		jen.Func().ID("AttachUserAgentDataToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("ua").Op("*").Qual("github.com/mssola/user_agent", "UserAgent")).Body(
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
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AttachItemIDToSpan attaches an item ID to a given span."),
		jen.Line(),
		jen.Func().ID("AttachItemIDToSpan").Params(jen.ID("span").ID("trace").Dot("Span"), jen.ID("itemID").ID("uint64")).Body(
			jen.ID("attachUint64ToSpan").Call(
				jen.ID("span"),
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			)),
		jen.Line(),
	)

	return code
}
