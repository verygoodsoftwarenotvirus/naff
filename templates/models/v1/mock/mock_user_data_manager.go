package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockUserDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "UserDataManager").Equals().Parens(jen.PointerTo().ID("UserDataManager")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserDataManager is a mocked models.UserDataManager for testprojects"),
		jen.Line(),
		jen.Type().ID("UserDataManager").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUserWithUnverifiedTwoFactorSecret is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUserWithUnverifiedTwoFactorSecret").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("VerifyUserTwoFactorSecret is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("VerifyUserTwoFactorSecret").Params(
			constants.CtxParam(),
			constants.UserIDParam(),
		).Params(jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.Return(jen.ID("args").Dot("Error").Call(jen.Zero())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUserByUsername is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUserByUsername").Params(constants.CtxParam(), jen.ID("username").String()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("username")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllUsersCount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetAllUsersCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUsers is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUsers").Params(constants.CtxParam(), jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("CreateUser").Params(constants.CtxParam(), jen.ID("input").Qual(proj.ModelsV1Package(), "UserDatabaseCreationInput")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("UpdateUser").Params(constants.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUserPassword is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("UpdateUserPassword").Params(
			constants.CtxParam(),
			jen.ID(constants.UserIDVarName).Uint64(),
			jen.ID("newHash").String(),
		).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(
				constants.CtxVar(),
				jen.ID(constants.UserIDVarName),
				jen.ID("newHash"),
			).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("ArchiveUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	)

	return code
}
