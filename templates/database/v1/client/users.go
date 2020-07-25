package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("dbclient")

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.ModelsV1Package(), "UserDataManager").Equals().Parens(jen.PointerTo().ID("Client")).Call(jen.Nil()),
			jen.Line(),
			jen.Comment("ErrUserExists is a sentinel error for returning when a username is taken."),
			jen.ID("ErrUserExists").Equals().Qual("errors", "New").Call(jen.Lit("error: username already exists")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUser fetches a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(),
			"User",
		),
			jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetUser")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)).Dot("Debug").Call(jen.Lit("GetUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetUser").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified 2FA secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetUserWithUnverifiedTwoFactorSecret").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(),
			"User",
		),
			jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetUserWithUnverifiedTwoFactorSecret")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)).Dot("Debug").Call(jen.Lit("GetUserWithUnverifiedTwoFactorSecret called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetUserWithUnverifiedTwoFactorSecret").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("VerifyUserTwoFactorSecret marks a user's two factor secret as validated."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("VerifyUserTwoFactorSecret").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("VerifyUserTwoFactorSecret")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)).Dot("Debug").Call(jen.Lit("VerifyUserTwoFactorSecret called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("VerifyUserTwoFactorSecret").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUserByUsername fetches a user by their username."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetUserByUsername").Params(constants.CtxParam(), jen.ID("username").String()).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"), jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetUserByUsername")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("username")),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("username"), jen.ID("username")).Dot("Debug").Call(jen.Lit("GetUserByUsername called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetUserByUsername").Call(constants.CtxVar(), jen.ID("username")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllUsersCount fetches a count of users from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetAllUsersCount").Params(constants.CtxParam()).Params(jen.ID("count").Uint64(), jen.Err().Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetAllUsersCount")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("GetAllUsersCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllUsersCount").Call(constants.CtxVar()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUsers fetches a list of users from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("GetUsers").Params(constants.CtxParam(), jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList"), jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("GetUsers")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachFilterToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.FilterVarName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit(constants.FilterVarName), jen.ID(constants.FilterVarName)).Dot("Debug").Call(jen.Lit("GetUsers called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetUsers").Call(constants.CtxVar(), jen.ID(constants.FilterVarName)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateUser creates a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("CreateUser").Params(
			constants.CtxParam(),
			jen.ID("input").Qual(proj.ModelsV1Package(), "UserDatabaseCreationInput"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Error(),
		).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("CreateUser")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("input").Dot("Username")),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("username"), jen.ID("input").Dot("Username")).Dot("Debug").Call(jen.Lit("CreateUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("CreateUser").Call(constants.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUser receives a complete User struct and updates its record in the database."),
		jen.Line(),
		jen.Comment("NOTE: this function uses the ID provided in the input to make its query."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("UpdateUser").Params(constants.CtxParam(), jen.ID("updated").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("UpdateUser")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("updated").Dot("Username")),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("username"), jen.ID("updated").Dot("Username")).Dot("Debug").Call(jen.Lit("UpdateUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateUser").Call(constants.CtxVar(), jen.ID("updated")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUserPassword updates a user's password hash in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("UpdateUserPassword").Params(
			constants.CtxParam(),
			jen.ID("userID").Uint64(),
			jen.ID("newHash").String(),
		).Params(jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("UpdateUserPassword")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("userID")),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("UpdateUserPassword called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateUserPassword").Call(
				constants.CtxVar(),
				jen.ID("userID"),
				jen.ID("newHash"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveUser archives a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").PointerTo().ID("Client")).ID("ArchiveUser").Params(constants.CtxParam(), constants.UserIDParam()).Params(jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("ArchiveUser")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID(constants.UserIDVarName)),
			jen.ID("c").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID(constants.UserIDVarName)).Dot("Debug").Call(jen.Lit("ArchiveUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("ArchiveUser").Call(constants.CtxVar(), jen.ID(constants.UserIDVarName)),
		),
		jen.Line(),
	)

	return code
}
