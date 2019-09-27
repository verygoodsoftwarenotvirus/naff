package client

import jen "github.com/dave/jennifer/jen"

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(jen.Null())
	ret.Add(jen.Var().Id("oauth2ClientsBasePath").Op("=").Lit("oauth2/clients"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildGetOAuth2ClientRequest").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("oauth2ClientsBasePath"), jen.Qual("strconv", "FormatUint").Call(
					jen.Id("id"),
					jen.Lit(10),
				),
			), jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("GetOAuth2Client").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Id("oauth2Client").Op("*").Id("models").Dot("OAuth2Client"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetOAuth2ClientRequest").Call(
				jen.Id("ctx"),
				jen.Id("id"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("oauth2Client"),
			),
			jen.Return().List(
				jen.Id("oauth2Client"),
				jen.Id("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("BuildGetOAuth2ClientsRequest").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("filter").Dot("ToValues").Call(),
				jen.Id("oauth2ClientsBasePath"),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("GetOAuth2Clients").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Op("*").Id("models").Dot("OAuth2ClientList"),
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetOAuth2ClientsRequest").Call(
				jen.Id("ctx"),
				jen.Id("filter"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Var().Id("oauth2Clients").Op("*").Id("models").Dot("OAuth2ClientList"),
			jen.Id("err").Op("=").Id("c").Dot("retrieve").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("oauth2Clients"),
			),
			jen.Return().List(
				jen.Id("oauth2Clients"),
				jen.Id("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("BuildCreateOAuth2ClientRequest").Params(
			ctxParam(),
			jen.Id("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.Id("body").Op("*").Id("models").Dot("OAuth2ClientCreationInput"),
		).Params(jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error")).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("buildVersionlessURL").Call(
				jen.Id("nil"),
				jen.Lit("oauth2"),
				jen.Lit("client"),
			),
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.Id("uri"),
				jen.Id("body"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("err"),
				),
			),
			jen.Id("req").Dot("AddCookie").Call(
				jen.Id("cookie"),
			),
			jen.Return().List(jen.Id("req"),
				jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("CreateOAuth2Client").Params(
			ctxParam(),
			jen.Id("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.Id("input").Op("*").Id("models").Dot("OAuth2ClientCreationInput"),
		).Params(
			jen.Op("*").Id("models").Dot("OAuth2Client"),
			jen.Id("error"),
		).Block(
			jen.Var().Id("oauth2Client").Op("*").Id("models").Dot("OAuth2Client"),
			jen.If(jen.Id("cookie").Op("==").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("errors").Dot("New").Call(
						jen.Lit("cookie required for request"),
					),
				),
			),
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildCreateOAuth2ClientRequest").Call(
				jen.Id("ctx"),
				jen.Id("cookie"),
				jen.Id("input"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("err"),
				),
			),
			jen.List(
				jen.Id("res"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("executeRawRequest").Call(
				jen.Id("ctx"),
				jen.Id("c").Dot("plainClient"),
				jen.Id("req"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("executing request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.If(jen.Id("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("ErrNotFound"),
				),
			),
			jen.If(jen.Id("resErr").Op(":=").Id("unmarshalBody").Call(
				jen.Id("res"),
				jen.Op("&").Id("oauth2Client"),
			),
				jen.Id("resErr").Op("!=").Id("nil"),
			).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("errors").Dot("Wrap").Call(
						jen.Id("resErr"),
						jen.Lit("loading response from server"),
					),
				),
			),
			jen.Return().List(jen.Id("oauth2Client"),
				jen.Id("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildArchiveOAuth2ClientRequest").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("oauth2ClientsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("id"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("ArchiveOAuth2Client").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildArchiveOAuth2ClientRequest").Call(
				jen.Id("ctx"),
				jen.Id("id"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Id("err"),
				),
			),
			jen.Return().Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	return ret
}
