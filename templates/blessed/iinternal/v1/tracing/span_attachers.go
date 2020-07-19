package encoding

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanAttachersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Const().Defs(
			buildConstants(proj)...,
		),
		jen.Line(),
	)

	code.Add(buildAttachUint64ToSpan()...)
	code.Add(buildAttachStringToSpan()...)
	code.Add(buildAttachFilterToSpan(proj)...)

	searchEnabled := false
	for _, typ := range proj.DataTypes {
		code.Add(buildAttachSomethingIDToSpan(typ)...)
		if typ.SearchEnabled {
			searchEnabled = true
		}
	}

	code.Add(buildAttachUserIDToSpan()...)
	code.Add(buildAttachOAuth2ClientDatabaseIDToSpan()...)
	code.Add(buildAttachOAuth2ClientIDToSpan()...)
	code.Add(buildAttachUsernameToSpan()...)
	code.Add(buildAttachWebhookIDToSpan()...)
	code.Add(buildAttachRequestURIToSpan()...)

	if searchEnabled {
		code.Add(buildAttachSearchQueryToSpan()...)
	}

	return code
}

func buildConstants(proj *models.Project) []jen.Code {
	lines := []jen.Code{}

	searchEnabled := false
	for _, typ := range proj.DataTypes {
		lines = append(lines, jen.IDf("%sIDSpanAttachmentKey", typ.Name.UnexportedVarName()).Equals().Litf("%s_id", typ.Name.RouteName()))
		if typ.SearchEnabled {
			searchEnabled = true
		}
	}

	lines = append(lines,
		jen.ID("userIDSpanAttachmentKey").Equals().Lit("user_id"),
		jen.ID("usernameSpanAttachmentKey").Equals().Lit("username"),
		jen.ID("filterPageSpanAttachmentKey").Equals().Lit("filter_page"),
		jen.ID("filterLimitSpanAttachmentKey").Equals().Lit("filter_limit"),
		jen.ID("oauth2ClientDatabaseIDSpanAttachmentKey").Equals().Lit("oauth2client_id"),
		jen.ID("oauth2ClientIDSpanAttachmentKey").Equals().Lit("client_id"),
		jen.ID("webhookIDSpanAttachmentKey").Equals().Lit("webhook_id"),
		jen.ID("requestURISpanAttachmentKey").Equals().Lit("request_uri"),
	)

	if searchEnabled {
		lines = append(lines,
			jen.ID("searchQuerySpanAttachmentKey").Equals().Lit("search_query"),
		)
	}

	return lines
}

func buildAttachUint64ToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("attachUint64ToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID("attachmentKey").String(),
			jen.ID("id").Uint64(),
		).Block(
			jen.If(jen.ID(utils.SpanVarName).DoesNotEqual().Nil()).Block(
				jen.ID(utils.SpanVarName).Dot("AddAttributes").Call(
					jen.Qual(constants.TracingLibrary, "StringAttribute").Call(
						jen.ID("attachmentKey"),
						jen.Qual("strconv", "FormatUint").Call(jen.ID("id"), jen.Lit(10)),
					),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachStringToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("attachStringToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.List(jen.ID("key"), jen.ID("str").String()),
		).Block(
			jen.If(jen.ID(utils.SpanVarName).Op("!=").Nil()).Block(
				jen.ID(utils.SpanVarName).Dot("AddAttributes").Call(jen.Qual(constants.TracingLibrary, "StringAttribute").Call(jen.ID("key"), jen.ID("str"))),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachFilterToSpan(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachFilterToSpan provides a consistent way to attach a filter's info to a span."),
		jen.Line(),
		jen.Func().ID("AttachFilterToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Block(
			jen.If(jen.ID(constants.FilterVarName).DoesNotEqual().Nil().And().ID(utils.SpanVarName).DoesNotEqual().Nil()).Block(
				jen.ID(utils.SpanVarName).Dot("AddAttributes").Callln(
					jen.Qual(constants.TracingLibrary, "StringAttribute").Call(
						jen.ID("filterPageSpanAttachmentKey"),
						jen.Qual("strconv", "FormatUint").Call(jen.ID(constants.FilterVarName).Dot("QueryPage").Call(), jen.Lit(10)),
					),
					jen.Qual(constants.TracingLibrary, "StringAttribute").Call(
						jen.ID("filterLimitSpanAttachmentKey"),
						jen.Qual("strconv", "FormatUint").Call(jen.Uint64().Call(jen.ID(constants.FilterVarName).Dot("Limit")), jen.Lit(10)),
					),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachSomethingIDToSpan(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	funcName := fmt.Sprintf("Attach%sIDToSpan", sn)
	paramName := fmt.Sprintf("%sID", uvn)

	lines := []jen.Code{
		jen.Commentf("%s attaches %s ID to a given span.", funcName, scnwp),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.IDf("%sIDSpanAttachmentKey", uvn), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachUserIDToSpan() []jen.Code {
	const (
		funcName  = "AttachUserIDToSpan"
		paramName = "userID"
	)

	lines := []jen.Code{
		jen.Commentf("%s provides a consistent way to attach a user's ID to a span.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("userIDSpanAttachmentKey"), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachOAuth2ClientDatabaseIDToSpan() []jen.Code {
	const (
		funcName  = "AttachOAuth2ClientDatabaseIDToSpan"
		paramName = "oauth2ClientID"
	)

	lines := []jen.Code{
		jen.Commentf("%s is a consistent way to attach an oauth2 client's ID to a span.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("oauth2ClientDatabaseIDSpanAttachmentKey"), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachOAuth2ClientIDToSpan() []jen.Code {
	const (
		funcName  = "AttachOAuth2ClientIDToSpan"
		paramName = "clientID"
	)

	lines := []jen.Code{
		jen.Commentf("%s is a consistent way to attach an oauth2 client's Client ID to a span.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).String(),
		).Block(
			jen.ID("attachStringToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("oauth2ClientIDSpanAttachmentKey"), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachUsernameToSpan() []jen.Code {
	const (
		funcName  = "AttachUsernameToSpan"
		paramName = "username"
	)

	lines := []jen.Code{
		jen.Commentf("%s provides a consistent way to attach a user's username to a span.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).String(),
		).Block(
			jen.ID("attachStringToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("usernameSpanAttachmentKey"), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachWebhookIDToSpan() []jen.Code {
	const (
		funcName  = "AttachWebhookIDToSpan"
		paramName = "webhookID"
	)

	lines := []jen.Code{
		jen.Commentf("%s provides a consistent way to attach a webhook's ID to a span.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("webhookIDSpanAttachmentKey"), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachRequestURIToSpan() []jen.Code {
	const (
		funcName  = "AttachRequestURIToSpan"
		paramName = "uri"
	)

	lines := []jen.Code{
		jen.Commentf("%s attaches a given URI to a span.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).String(),
		).Block(
			jen.ID("attachStringToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("requestURISpanAttachmentKey"), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}

func buildAttachSearchQueryToSpan() []jen.Code {
	const (
		funcName  = "AttachSearchQueryToSpan"
		paramName = "query"
	)

	lines := []jen.Code{
		jen.Commentf("%s attaches a given search query to a span.", funcName),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(constants.TracingLibrary, "Span"),
			jen.ID(paramName).String(),
		).Block(
			jen.ID("attachStringToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("searchQuerySpanAttachmentKey"), jen.ID(paramName)),
		),
		jen.Line(),
	}

	return lines
}
