package client

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(jen.Null())
	ret.Add(jen.Var().ID("oauth2ClientsBasePath").Op("=").Lit("oauth2/clients"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID(v1),
		).ID("BuildGetOAuth2ClientRequest").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("oauth2ClientsBasePath"), jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			), jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID(v1),
		).ID("GetOAuth2Client").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("oauth2Client").Op("*").Qual(modelsPkg, "OAuth2Client"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetOAuth2ClientRequest").Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("oauth2Client"),
			),
			jen.Return().List(
				jen.ID("oauth2Client"),
				jen.ID("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientsRequest").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("oauth2ClientsBasePath"),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID(v1),
		).ID("GetOAuth2Clients").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual(modelsPkg, "OAuth2ClientList"),
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetOAuth2ClientsRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Var().ID("oauth2Clients").Op("*").Qual(modelsPkg, "OAuth2ClientList"),
			jen.ID("err").Op("=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("oauth2Clients"),
			),
			jen.Return().List(
				jen.ID("oauth2Clients"),
				jen.ID("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildCreateOAuth2ClientRequest").Params(
			ctxParam(),
			jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.ID("body").Op("*").Qual(modelsPkg, "OAuth2ClientCreationInput"),
		).Params(jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error")).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("nil"),
				jen.Lit("oauth2"),
				jen.Lit("client"),
			),
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("err"),
				),
			),
			jen.ID("req").Dot("AddCookie").Call(
				jen.ID("cookie"),
			),
			jen.Return().List(jen.ID("req"),
				jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID(v1),
		).ID("CreateOAuth2Client").Params(
			ctxParam(),
			jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.ID("input").Op("*").Qual(modelsPkg, "OAuth2ClientCreationInput"),
		).Params(
			jen.Op("*").Qual(modelsPkg, "OAuth2Client"),
			jen.ID("error"),
		).Block(
			jen.Var().ID("oauth2Client").Op("*").Qual(modelsPkg, "OAuth2Client"),
			jen.If(jen.ID("cookie").Op("==").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("errors").Dot("New").Call(
						jen.Lit("cookie required for request"),
					),
				),
			),
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildCreateOAuth2ClientRequest").Call(
				jen.ID("ctx"),
				jen.ID("cookie"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("err"),
				),
			),
			jen.List(
				jen.ID("res"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("executeRawRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("plainClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("executing request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.If(jen.ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("ErrNotFound"),
				),
			),
			jen.If(jen.ID("resErr").Op(":=").ID("unmarshalBody").Call(
				jen.ID("res"),
				jen.Op("&").ID("oauth2Client"),
			),
				jen.ID("resErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("errors").Dot("Wrap").Call(
						jen.ID("resErr"),
						jen.Lit("loading response from server"),
					),
				),
			),
			jen.Return().List(jen.ID("oauth2Client"),
				jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID(v1),
		).ID("BuildArchiveOAuth2ClientRequest").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("oauth2ClientsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.ID("c").Op("*").ID(v1),
		).ID("ArchiveOAuth2Client").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	return ret
}
