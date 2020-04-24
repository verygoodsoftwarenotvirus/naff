package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockOauth2ClientDataManagerDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "OAuth2ClientDataManager").Equals().Parens(jen.PointerTo().ID("OAuth2ClientDataManager")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ClientDataManager is a mocked models.OAuth2ClientDataManager for testing"),
		jen.Line(),
		jen.Type().ID("OAuth2ClientDataManager").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("GetOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("clientID"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetOAuth2ClientByClientID is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("GetOAuth2ClientByClientID").Params(constants.CtxParam(), jen.ID("identifier").String()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("identifier")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	//ret.Add(
	//	jen.Comment("GetOAuth2ClientCount is a mock function"),
	//	jen.Line(),
	//	jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("GetOAuth2ClientCount").Params(
	//		utils.CtxParam(),
	//		jen.ID("userID").Uint64(),
	//		jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
	//	).Params(jen.Uint64(), jen.Error()).Block(
	//		jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID"), jen.ID(utils.FilterVarName)),
	//		jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
	//	),
	//	jen.Line(),
	//)

	ret.Add(
		jen.Comment("GetAllOAuth2ClientCount is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("GetAllOAuth2ClientCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllOAuth2Clients is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("GetAllOAuth2Clients").Params(constants.CtxParam()).Params(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	//ret.Add(
	//	jen.Comment("GetAllOAuth2ClientsForUser is a mock function"),
	//	jen.Line(),
	//	jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("GetAllOAuth2ClientsForUser").Params(utils.CtxParam(), jen.ID("userID").Uint64()).Params(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
	//		jen.Error()).Block(
	//		jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID")),
	//		jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.One())),
	//	),
	//	jen.Line(),
	//)

	ret.Add(
		jen.Comment("GetOAuth2Clients is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("GetOAuth2Clients").Params(
			constants.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("userID"), jen.ID(constants.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("CreateOAuth2Client").Params(constants.CtxParam(), jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("UpdateOAuth2Client").Params(constants.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveOAuth2Client is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("OAuth2ClientDataManager")).ID("ArchiveOAuth2Client").Params(constants.CtxParam(), jen.List(jen.ID("clientID"), jen.ID("userID")).Uint64()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("clientID"), jen.ID("userID")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	return ret
}
