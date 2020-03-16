package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockOauth2ClientDataManagerDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientDataManager").Op("=").Parens(jen.Op("*").ID("OAuth2ClientDataManager")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ClientDataManager is a mocked models.OAuth2ClientDataManager for testing"),
		jen.Line(),
		jen.Type().ID("OAuth2ClientDataManager").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2Client").Params(utils.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("clientID"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientByClientID is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2ClientByClientID").Params(utils.CtxParam(), jen.ID("identifier").ID("string")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("identifier")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientCount is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2ClientCount").Params(utils.CtxParam(), jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientCount is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetAllOAuth2ClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2Clients is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetAllOAuth2Clients").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientsForUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetAllOAuth2ClientsForUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Clients is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("GetOAuth2Clients").Params(utils.CtxParam(), jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientList"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("CreateOAuth2Client").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientCreationInput")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("UpdateOAuth2Client").Params(utils.CtxParam(), jen.ID("updated").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2Client")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("updated")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("OAuth2ClientDataManager")).ID("ArchiveOAuth2Client").Params(utils.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("clientID"), jen.ID("userID")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
