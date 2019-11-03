package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockOauth2ClientDataServerDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2ClientDataServer").Op("=").Parens(jen.Op("*").ID("OAuth2ClientDataServer")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ClientDataServer is a mocked models.OAuth2ClientDataServer for testing"),
		jen.Line(),
		jen.Type().ID("OAuth2ClientDataServer").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ListHandler is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateHandler is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ReadHandler is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveHandler is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "HandlerFunc")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreationInputMiddleware is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ClientInfoMiddleware is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ExtractOAuth2ClientFromRequest is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("req")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HandleAuthorizeRequest is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("HandleAuthorizeRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HandleTokenRequest is the obligatory implementation for our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataServer")).ID("HandleTokenRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
