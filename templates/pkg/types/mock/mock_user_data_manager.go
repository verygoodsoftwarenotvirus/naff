package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockUserDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Underscore().Qual(proj.TypesPackage(), "UserDataManager").Equals().Parens(jen.PointerTo().ID("UserDataManager")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildUserDataManager()...)
	code.Add(buildGetUser(proj)...)
	code.Add(buildGetUserWithUnverifiedTwoFactorSecret(proj)...)
	code.Add(buildVerifyUserTwoFactorSecret()...)
	code.Add(buildGetUserByUsername(proj)...)
	code.Add(buildGetAllUsersCount()...)
	code.Add(buildGetUsers(proj)...)
	code.Add(buildCreateUser(proj)...)
	code.Add(buildUpdateUser(proj)...)
	code.Add(buildUpdateUserPassword()...)
	code.Add(buildArchiveUser()...)

	return code
}

func buildUserDataManager() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UserDataManager is a mocked models.UserDataManager for testing"),
		jen.Line(),
		jen.Type().ID("UserDataManager").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildGetUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "User"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildGetUserWithUnverifiedTwoFactorSecret(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetUserWithUnverifiedTwoFactorSecret is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUserWithUnverifiedTwoFactorSecret").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "User"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildVerifyUserTwoFactorSecret() []jen.Code {
	lines := []jen.Code{
		jen.Comment("VerifyUserTwoFactorSecret is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("VerifyUserTwoFactorSecret").Params(
			constants.CtxParam(),
			constants.UserIDParam(),
		).Params(jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
			jen.Return(jen.ID("args").Dot("Error").Call(jen.Zero())),
		),
		jen.Line(),
	}

	return lines
}

func buildGetUserByUsername(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetUserByUsername is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUserByUsername").Params(constants.CtxParam(), jen.ID("username").String()).Params(jen.PointerTo().Qual(proj.TypesPackage(), "User"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("username")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildGetAllUsersCount() []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetAllUsersCount is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetAllUsersCount").Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar()),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Uint64()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildGetUsers(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetUsers is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("GetUsers").Params(constants.CtxParam(), jen.ID(constants.FilterVarName).PointerTo().Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "UserList"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "UserList")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("CreateUser").Params(constants.CtxParam(), jen.ID("input").Qual(proj.TypesPackage(), "UserDatabaseCreationInput")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "User"),
			jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.PointerTo().Qual(proj.TypesPackage(), "User")), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateUser(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdateUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("UpdateUser").Params(constants.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.TypesPackage(), "User")).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("updated")).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateUserPassword() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdateUserPassword is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("UpdateUserPassword").Params(
			constants.CtxParam(),
			jen.ID(constants.UserIDVarName).Uint64(),
			jen.ID("newHash").String(),
		).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(
				constants.CtxVar(),
				jen.ID(constants.UserIDVarName),
				jen.ID("newHash"),
			).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveUser() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveUser is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataManager")).ID("ArchiveUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Error()).Body(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
