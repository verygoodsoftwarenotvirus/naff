package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(jen.Null())
	ret.Add(jen.Const().Defs(
		jen.ID("oauth2ClientsBasePath").Op("=").Lit("oauth2/clients"),
	))

	ret.Add(
		jen.Comment("BuildGetOAuth2ClientRequest builds an HTTP request for fetching an OAuth2 client"),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientRequest").Params(
			utils.CtxParam(),
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
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Client gets an OAuth2 client"),
		jen.Line(),
		newClientMethod("GetOAuth2Client").Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("oauth2Client").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
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
			jen.Line(),
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
		jen.Comment("BuildGetOAuth2ClientsRequest builds an HTTP request for fetching a list of OAuth2 clients"),
		jen.Line(),
		newClientMethod("BuildGetOAuth2ClientsRequest").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("oauth2ClientsBasePath"),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Clients gets a list of OAuth2 clients"),
		jen.Line(),
		newClientMethod("GetOAuth2Clients").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientList"),
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
			jen.Line(),
			jen.Var().ID("oauth2Clients").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientList"),
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
		jen.Comment("BuildCreateOAuth2ClientRequest builds an HTTP request for creating OAuth2 clients"),
		jen.Line(),
		newClientMethod("BuildCreateOAuth2ClientRequest").Paramsln(
			utils.CtxParam(),
			jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.ID("body").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientCreationInput"),
		).Params(jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error")).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
				jen.ID("nil"),
				jen.Lit("oauth2"),
				jen.Lit("client"),
			),
			jen.Line(),
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
			jen.Line(),
			jen.Return().List(
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(utils.Comments(
		"CreateOAuth2Client creates an OAuth2 client. Note that cookie must not be nil",
		"in order to receive a valid response",
	)...)
	ret.Add(
		newClientMethod("CreateOAuth2Client").Paramsln(
			utils.CtxParam(),
			jen.ID("cookie").Op("*").Qual("net/http", "Cookie"),
			jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientCreationInput"),
		).Params(
			jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
			jen.ID("error"),
		).Block(
			jen.Var().ID("oauth2Client").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
			jen.If(jen.ID("cookie").Op("==").ID("nil")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("errors", "New").Call(
						jen.Lit("cookie required for request"),
					),
				),
			),
			jen.Line(),
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
			jen.Line(),
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
			jen.Line(),
			jen.If(jen.ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusNotFound")).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("ErrNotFound"),
				),
			),
			jen.Line(),
			jen.If(jen.ID("resErr").Op(":=").ID("unmarshalBody").Call(
				jen.ID("res"),
				jen.Op("&").ID("oauth2Client"),
			),
				jen.ID("resErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("loading response from server: %w"),
						jen.ID("resErr"),
					),
				),
			),
			jen.Line(),
			jen.Return().List(
				jen.ID("oauth2Client"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("BuildArchiveOAuth2ClientRequest builds an HTTP request for archiving an oauth2 client"),
		jen.Line(),
		newClientMethod("BuildArchiveOAuth2ClientRequest").Params(
			utils.CtxParam(),
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
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveOAuth2Client archives an OAuth2 client"),
		jen.Line(),
		newClientMethod("ArchiveOAuth2Client").Params(
			utils.CtxParam(),
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
			jen.Line(),
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
