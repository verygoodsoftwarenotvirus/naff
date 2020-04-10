package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile(packageName)

	utils.AddImports(proj, ret)

	ret.Add(jen.Const().ID("usersBasePath").Equals().Lit("users"))

	ret.Add(buildBuildGetUserRequest(proj)...)
	ret.Add(buildGetUser(proj)...)
	ret.Add(buildBuildGetUsersRequest(proj)...)
	ret.Add(buildGetUsers(proj)...)
	ret.Add(buildBuildCreateUserRequest(proj)...)
	ret.Add(buildCreateUser(proj)...)
	ret.Add(buildBuildArchiveUserRequest(proj)...)
	ret.Add(buildArchiveUser(proj)...)
	ret.Add(buildBuildLoginRequest(proj)...)
	ret.Add(buildLogin(proj)...)

	return ret
}

func buildBuildGetUserRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetUserRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("userID"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildGetUserRequest builds an HTTP request for fetching a user"),
		jen.Line(),
		newClientMethod("BuildGetUserRequest").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetUser(proj *models.Project) []jen.Code {
	funcName := "GetUser"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetUserRequest").Call(
			utils.CtxVar(),
			jen.ID("userID"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Equals().ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.AddressOf().ID("user"),
		),
		jen.Return().List(
			jen.ID("user"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Comment("GetUser retrieves a user"),
		jen.Line(),
		newClientMethod("GetUser").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
		).Params(
			jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Err().Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUsersRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetUsersRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.ID(utils.FilterVarName).Dot("ToValues").Call(),
			jen.ID("usersBasePath"),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildGetUsersRequest builds an HTTP request for fetching a user"),
		jen.Line(),
		newClientMethod("BuildGetUsersRequest").Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetUsers(proj *models.Project) []jen.Code {
	funcName := "GetUsers"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("users").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserList").Values(),
		jen.Line(),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetUsersRequest").Call(
			utils.CtxVar(),
			jen.ID(utils.FilterVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Equals().ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.AddressOf().ID("users"),
		),
		jen.Return().List(jen.ID("users"), jen.Err()),
	}

	lines := []jen.Code{
		jen.Comment("GetUsers retrieves a list of users"),
		jen.Line(),
		newClientMethod("GetUsers").Params(
			utils.CtxParam(),
			jen.ID(utils.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateUserRequest(proj *models.Project) []jen.Code {
	funcName := "BuildCreateUserRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildCreateUserRequest builds an HTTP request for creating a user"),
		jen.Line(),
		newClientMethod("BuildCreateUserRequest").Params(
			utils.CtxParam(),
			jen.ID("body").PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildCreateUser(proj *models.Project) []jen.Code {
	funcName := "CreateUser"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("user").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserCreationResponse").Values(),
		jen.Line(),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildCreateUserRequest").Call(
			utils.CtxVar(),
			jen.ID("input"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Equals().ID("c").Dot("executeUnauthenticatedDataRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.AddressOf().ID("user"),
		),
		jen.Return().List(jen.ID("user"), jen.Err()),
	}

	lines := []jen.Code{
		jen.Comment("CreateUser creates a new user"),
		jen.Line(),
		newClientMethod("CreateUser").Params(
			utils.CtxParam(),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "UserCreationResponse"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveUserRequest(proj *models.Project) []jen.Code {
	funcName := "BuildArchiveUserRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("userID"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildArchiveUserRequest builds an HTTP request for updating a user"),
		jen.Line(),
		newClientMethod("BuildArchiveUserRequest").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildArchiveUser(proj *models.Project) []jen.Code {
	funcName := "ArchiveUser"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Assign().ID("c").Dot("BuildArchiveUserRequest").Call(
			utils.CtxVar(),
			jen.ID("userID"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Comment("ArchiveUser archives a user"),
		jen.Line(),
		newClientMethod("ArchiveUser").Params(
			utils.CtxParam(),
			jen.ID("userID").Uint64(),
		).Params(jen.Error()).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildLoginRequest(proj *models.Project) []jen.Code {
	funcName := "BuildLoginRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.If(jen.ID("input").DoubleEquals().Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("nil input provided"))),
		),
		jen.Line(),
		jen.List(jen.ID("body"), jen.Err()).Assign().ID("createBodyFromStruct").Call(jen.AddressOf().ID("input")),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request body: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
			jen.Lit("login"),
		),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	}

	lines := []jen.Code{
		jen.Comment("BuildLoginRequest builds an authenticating HTTP request"),
		jen.Line(),
		newClientMethod("BuildLoginRequest").Params(
			utils.CtxParam(),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildLogin(proj *models.Project) []jen.Code {
	funcName := "Login"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.If(jen.ID("input").DoubleEquals().Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("nil input provided"))),
		),
		jen.Line(),
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildLoginRequest").Call(utils.CtxVar(), jen.ID("input")),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.List(
			jen.ID("res"),
			jen.Err(),
		).Assign().ID("c").Dot("plainClient").Dot("Do").Call(
			jen.ID("req"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encountered error executing login request: %w"),
					jen.Err(),
				),
			),
		),
		jen.ID("c").Dot("closeResponseBody").Call(jen.ID("res")),
		jen.Line(),
		jen.ID("cookies").Assign().ID("res").Dot("Cookies").Call(),
		jen.If(jen.ID("len").Call(
			jen.ID("cookies"),
		).Op(">").Zero(),
		).Block(
			jen.Return().List(jen.ID("cookies").Index(
				jen.Zero(),
			),
				jen.Nil(),
			),
		),
		jen.Line(),
		jen.Return().List(
			jen.Nil(),
			jen.Qual("errors", "New").Call(
				jen.Lit("no cookies returned from request"),
			),
		),
	}

	lines := []jen.Code{
		jen.Comment("Login will, when provided the correct credentials, fetch a login cookie"),
		jen.Line(),
		newClientMethod("Login").Params(
			utils.CtxParam(),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput"),
		).Params(
			jen.PointerTo().Qual("net/http", "Cookie"),
			jen.Error(),
		).Block(block...),
		jen.Line(),
	}

	return lines

}
