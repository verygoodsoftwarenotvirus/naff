package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanAttachersDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			buildConstants(proj)...,
		),
		jen.Line(),
	)

	ret.Add(buildAttachUint64ToSpan()...)
	ret.Add(buildAttachStringToSpan()...)
	ret.Add(buildAttachFilterToSpan(proj)...)

	for _, typ := range proj.DataTypes {
		ret.Add(buildAttachSomethingIDToSpan(typ)...)
	}

	ret.Add(buildAttachUserIDToSpan()...)
	ret.Add(buildAttachOAuth2ClientDatabaseIDToSpan()...)
	ret.Add(buildAttachOAuth2ClientIDToSpan()...)
	ret.Add(buildAttachUsernameToSpan()...)
	ret.Add(buildAttachWebhookIDToSpan()...)
	ret.Add(buildAttachRequestURIToSpan()...)

	return ret
}

func buildConstants(proj *models.Project) []jen.Code {
	lines := []jen.Code{}

	for _, typ := range proj.DataTypes {
		lines = append(lines, jen.IDf("%sIDSpanAttachmentKey", typ.Name.UnexportedVarName()).Equals().Litf("%s_id", typ.Name.RouteName()))
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

	return lines
}

func buildAttachUint64ToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("attachUint64ToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID("attachmentKey").String(),
			jen.ID("id").Uint64(),
		).Block(
			jen.If(jen.ID(utils.SpanVarName).DoesNotEqual().Nil()).Block(
				jen.ID(utils.SpanVarName).Dot("AddAttributes").Call(
					jen.Qual(utils.TracingLibrary, "StringAttribute").Call(
						jen.ID("attachmentKey"),
						jen.Qual("strconv", "FormatUint").Call(jen.ID("id"), jen.Lit(10)),
					),
				),
			),
		),
	}

	return lines
}

func buildAttachStringToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("attachStringToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.List(jen.ID("key"), jen.ID("str").String()),
		).Block(
			jen.If(jen.ID(utils.SpanVarName).Op("!=").Nil()).Block(
				jen.ID(utils.SpanVarName).Dot("AddAttributes").Call(jen.Qual(utils.TracingLibrary, "StringAttribute").Call(jen.ID("key"), jen.ID("str"))),
			),
		),
	}

	return lines
}

func buildAttachFilterToSpan(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachFilterToSpan provides a consistent way to attach a filter's info to a span"),
		jen.Line(),
		jen.Func().ID("AttachFilterToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Block(
			jen.If(jen.ID(utils.FilterVarName).DoesNotEqual().Nil().Op("&&").ID(utils.SpanVarName).DoesNotEqual().Nil()).Block(
				jen.ID(utils.SpanVarName).Dot("AddAttributes").Callln(
					jen.Qual(utils.TracingLibrary, "StringAttribute").Call(
						jen.ID("filterPageSpanAttachmentKey"),
						jen.Qual("strconv", "FormatUint").Call(jen.ID(utils.FilterVarName).Dot("QueryPage").Call(), jen.Lit(10)),
					),
					jen.Qual(utils.TracingLibrary, "StringAttribute").Call(
						jen.ID("filterLimitSpanAttachmentKey"),
						jen.Qual("strconv", "FormatUint").Call(jen.ID(utils.FilterVarName).Dot("Limit"), jen.Lit(10)),
					),
				),
			),
		),
	}

	return lines
}

func buildAttachSomethingIDToSpan(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Commentf("Attach%sIDToSpan attaches %s ID to a given span", sn, scnwp),
		jen.Line(),
		jen.Func().IDf("Attach%sIDToSpan", sn).Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.IDf("%sID", uvn).Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.IDf("%sIDSpanAttachmentKey", uvn), jen.IDf("%sID", uvn)),
		),
	}

	return lines
}

func buildAttachUserIDToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachUserIDToSpan provides a consistent way to attach a user's ID to a span"),
		jen.Line(),
		jen.Func().ID("AttachUserIDToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID("userID").Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("userIDSpanAttachmentKey"), jen.ID("userID")),
		),
	}

	return lines
}

func buildAttachOAuth2ClientDatabaseIDToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachOAuth2ClientDatabaseIDToSpan is a consistent way to attach an oauth2 client's ID to a span"),
		jen.Line(),
		jen.Func().ID("AttachOAuth2ClientDatabaseIDToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID("oauth2ClientID").Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("oauth2ClientDatabaseIDSpanAttachmentKey"), jen.ID("oauth2ClientID")),
		),
	}

	return lines
}

func buildAttachOAuth2ClientIDToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachOAuth2ClientIDToSpan is a consistent way to attach an oauth2 client's Client ID to a span"),
		jen.Line(),
		jen.Func().ID("AttachOAuth2ClientIDToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID("clientID").String(),
		).Block(
			jen.ID("attachStringToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("oauth2ClientIDSpanAttachmentKey"), jen.ID("clientID")),
		),
	}

	return lines
}

func buildAttachUsernameToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachUsernameToSpan provides a consistent way to attach a user's username to a span"),
		jen.Line(),
		jen.Func().ID("AttachUsernameToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID("username").String(),
		).Block(
			jen.ID("attachStringToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("usernameSpanAttachmentKey"), jen.ID("username")),
		),
	}

	return lines
}

func buildAttachWebhookIDToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span"),
		jen.Line(),
		jen.Func().ID("AttachWebhookIDToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID("webhookID").Uint64(),
		).Block(
			jen.ID("attachUint64ToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("webhookIDSpanAttachmentKey"), jen.ID("webhookID")),
		),
	}

	return lines
}

func buildAttachRequestURIToSpan() []jen.Code {
	lines := []jen.Code{
		jen.Comment("AttachRequestURIToSpan attaches a given URI to a span"),
		jen.Line(),
		jen.Func().ID("AttachRequestURIToSpan").Params(
			jen.ID(utils.SpanVarName).PointerTo().Qual(utils.TracingLibrary, "Span"),
			jen.ID("uri").String(),
		).Block(
			jen.ID("attachStringToSpan").Call(jen.ID(utils.SpanVarName), jen.ID("requestURISpanAttachmentKey"), jen.ID("uri")),
		),
	}

	return lines
}
