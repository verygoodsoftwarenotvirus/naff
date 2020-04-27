package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockUserDataManagerDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "UserDataManager").Equals().Parens(jen.PointerTo().ID("UserDataManager")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserDataManager is a mocked models.UserDataManager for testing"),
		jen.Line(),
		jen.Type().ID("UserDataManager").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUser").Params(constants.CtxParam(), jen.ID("userID").Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserByUsername is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUserByUsername").Params(constants.CtxParam(), jen.ID("username").String()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("username")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllUserCount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetAllUserCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUsers is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUsers").Params(constants.CtxParam(), jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("CreateUser").Params(constants.CtxParam(), jen.ID("input").Qual(proj.ModelsV1Package(), "UserDatabaseCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("UpdateUser").Params(constants.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("ArchiveUser").Params(constants.CtxParam(), jen.ID("userID").Uint64()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("userID")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	return ret
}
