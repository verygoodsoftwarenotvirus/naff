package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(proj.ModelsV1Package(), "UserDataManager").Equals().Parens(jen.Op("*").ID("Client")).Call(jen.Nil()),
			jen.Line(),
			jen.Comment("ErrUserExists is a sentinel error for returning when a username is taken"),
			jen.ID("ErrUserExists").Equals().Qual("errors", "New").Call(jen.Lit("error: username already exists")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUser fetches a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(proj.ModelsV1Package(),
			"User",
		),
			jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("GetUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetUser").Call(utils.CtxVar(), jen.ID("userID")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserByUsername fetches a user by their username"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetUserByUsername").Params(utils.CtxParam(), jen.ID("username").ID("string")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "User"), jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetUserByUsername")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID("span"), jen.ID("username")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("username"), jen.ID("username")).Dot("Debug").Call(jen.Lit("GetUserByUsername called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetUserByUsername").Call(utils.CtxVar(), jen.ID("username")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetAllUserCount fetches a count of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAllUserCount").Params(utils.CtxParam()).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetUserCount")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID(utils.FilterVarName)),
			jen.ID("c").Dot("logger").Dot("Debug").Call(jen.Lit("GetAllUserCount called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetAllUserCount").Call(utils.CtxVar(), jen.ID(utils.FilterVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUsers fetches a list of users from the database that meet a particular filter"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetUsers").Params(utils.CtxParam(), jen.ID(utils.FilterVarName).Op("*").Qual(proj.ModelsV1Package(), "QueryFilter")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "UserList"), jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("GetUsers")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID(utils.FilterVarName)),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("filter"), jen.ID(utils.FilterVarName)).Dot("Debug").Call(jen.Lit("GetUsers called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("GetUsers").Call(utils.CtxVar(), jen.ID(utils.FilterVarName)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateUser creates a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CreateUser").Params(utils.CtxParam(), jen.ID("input").Op("*").Qual(proj.ModelsV1Package(), "UserInput")).Params(jen.Op("*").Qual(proj.ModelsV1Package(), "User"), jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("CreateUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID("span"), jen.ID("input").Dot("Username")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("username"), jen.ID("input").Dot("Username")).Dot("Debug").Call(jen.Lit("CreateUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("CreateUser").Call(utils.CtxVar(), jen.ID("input")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateUser receives a complete User struct and updates its record in the database."),
		jen.Line(),
		jen.Comment("NOTE: this function uses the ID provided in the input to make its query."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateUser").Params(utils.CtxParam(), jen.ID("updated").Op("*").Qual(proj.ModelsV1Package(), "User")).Params(jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("UpdateUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID("span"), jen.ID("updated").Dot("Username")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("username"), jen.ID("updated").Dot("Username")).Dot("Debug").Call(jen.Lit("UpdateUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("UpdateUser").Call(utils.CtxVar(), jen.ID("updated")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveUser archives a user"),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveUser").Params(utils.CtxParam(), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("ArchiveUser")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("userID")),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("userID")).Dot("Debug").Call(jen.Lit("ArchiveUser called")),
			jen.Line(),
			jen.Return().ID("c").Dot("querier").Dot("ArchiveUser").Call(utils.CtxVar(), jen.ID("userID")),
		),
		jen.Line(),
	)
	return ret
}
