package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockUserDataManagerDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().ID("_").Qual(proj.ModelsV1Package(), "UserDataManager").Equals().Parens(jen.PointerTo().ID("UserDataManager")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserDataManager is a mocked models.UserDataManager for testing"),
		jen.Line(),
		jen.Type().ID("UserDataManager").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUser").Params(utils.CtxParam(), jen.ID("userID").Uint64()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserByUsername is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUserByUsername").Params(utils.CtxParam(), jen.ID("username").String()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("username")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllUserCount is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetAllUserCount").Params(utils.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUsers is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUsers").Params(utils.CtxParam(), jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID(utils.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("CreateUser").Params(utils.CtxParam(), jen.ID("input").Qual(proj.ModelsV1Package(), "UserDatabaseCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("UpdateUser").Params(utils.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("ArchiveUser").Params(utils.CtxParam(), jen.ID("userID").Uint64()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(utils.CtxVar(), jen.ID("userID")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
