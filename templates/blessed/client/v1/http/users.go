package client

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("client")

	utils.AddImports(proj, ret)

	ret.Add(jen.Const().ID("usersBasePath").Op("=").Lit("users"))

	ret.Add(buildBuildGetUserRequest()...)
	ret.Add(buildGetUser(proj)...)
	ret.Add(buildBuildGetUsersRequest(proj)...)
	ret.Add(buildGetUsers(proj)...)
	ret.Add(buildBuildCreateUserRequest(proj)...)
	ret.Add(buildCreateUser(proj)...)
	ret.Add(buildBuildArchiveUserRequest()...)
	ret.Add(buildArchiveUser()...)
	ret.Add(buildBuildLoginRequest(proj)...)
	ret.Add(buildLogin()...)

	return ret
}

func buildBuildGetUserRequest() []jen.Code {
	funcName := "BuildGetUserRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
			jen.ID("nil"),
			jen.ID("usersBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("userID"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildGetUserRequest builds an HTTP request for fetching a user"),
		jen.Line(),
		newClientMethod("BuildGetUserRequest").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetUser(proj *models.Project) []jen.Code {
	funcName := "GetUser"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildGetUserRequest").Call(
			utils.CtxVar(),
			jen.ID("userID"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Op("=").ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Op("&").ID("user"),
		),
		jen.Return().List(
			jen.ID("user"),
			jen.Err(),
		),
	)

	lines := []jen.Code{
		jen.Comment("GetUser retrieves a user"),
		jen.Line(),
		newClientMethod("GetUser").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(
			jen.ID("user").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "User"),
			jen.Err().ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUsersRequest(proj *models.Project) []jen.Code {
	funcName := "BuildGetUsersRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
			jen.ID("filter").Dot("ToValues").Call(),
			jen.ID("usersBasePath"),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildGetUsersRequest builds an HTTP request for fetching a user"),
		jen.Line(),
		newClientMethod("BuildGetUsersRequest").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildGetUsers(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "GetUsers"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.ID("users").Op(":=").Op("&").Qual(filepath.Join(outPath, "models/v1"), "UserList").Values(),
		jen.Line(),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildGetUsersRequest").Call(
			utils.CtxVar(),
			jen.ID("filter"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Op("=").ID("c").Dot("retrieve").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Op("&").ID("users"),
		),
		jen.Return().List(jen.ID("users"), jen.Err()),
	)

	lines := []jen.Code{
		jen.Comment("GetUsers retrieves a list of users"),
		jen.Line(),
		newClientMethod("GetUsers").Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(filepath.Join(outPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual(filepath.Join(outPath, "models/v1"), "UserList"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateUserRequest(proj *models.Project) []jen.Code {
	funcName := "BuildCreateUserRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
			jen.ID("nil"),
			jen.ID("usersBasePath"),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildCreateUserRequest builds an HTTP request for creating a user"),
		jen.Line(),
		newClientMethod("BuildCreateUserRequest").Params(
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "UserInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildCreateUser(proj *models.Project) []jen.Code {
	outPath := proj.OutputPath
	funcName := "CreateUser"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.ID("user").Op(":=").Op("&").Qual(filepath.Join(outPath, "models/v1"), "UserCreationResponse").Values(),
		jen.Line(),
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildCreateUserRequest").Call(
			utils.CtxVar(),
			jen.ID("input"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.Err().Op("=").ID("c").Dot("executeUnathenticatedDataRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.Op("&").ID("user"),
		),
		jen.Return().List(jen.ID("user"), jen.Err()),
	)

	lines := []jen.Code{
		jen.Comment("CreateUser creates a new user"),
		jen.Line(),
		newClientMethod("CreateUser").Params(
			utils.CtxParam(),
			jen.ID("input").Op("*").Qual(filepath.Join(outPath, "models/v1"), "UserInput"),
		).Params(
			jen.Op("*").Qual(filepath.Join(outPath, "models/v1"), "UserCreationResponse"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveUserRequest() []jen.Code {
	funcName := "BuildArchiveUserRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
			jen.ID("nil"),
			jen.ID("usersBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID("userID"),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequest").Call(
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildArchiveUserRequest builds an HTTP request for updating a user"),
		jen.Line(),
		newClientMethod("BuildArchiveUserRequest").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildArchiveUser() []jen.Code {
	funcName := "ArchiveUser"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildArchiveUserRequest").Call(
			utils.CtxVar(),
			jen.ID("userID"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			utils.CtxVar(),
			jen.ID("req"),
			jen.ID("nil"),
		),
	)

	lines := []jen.Code{
		jen.Comment("ArchiveUser archives a user"),
		jen.Line(),
		newClientMethod("ArchiveUser").Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(jen.ID("error")).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildLoginRequest(proj *models.Project) []jen.Code {
	funcName := "BuildLoginRequest"

	block := utils.StartSpan(false, funcName)
	block = append(block,
		jen.List(
			jen.ID("body"),
			jen.Err(),
		).Op(":=").ID("createBodyFromStruct").Call(
			jen.Op("&").Qual(filepath.Join(proj.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
				jen.ID("Username").Op(":").ID("username"),
				jen.ID("Password").Op(":").ID("password"),
				jen.ID("TOTPToken").Op(":").ID("totpToken"),
			),
		),
		jen.Line(),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("creating body from struct: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.ID("uri").Op(":=").ID("c").Dot("buildVersionlessURL").Call(
			jen.ID("nil"),
			jen.ID("usersBasePath"),
			jen.Lit("login"),
		),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			utils.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	)

	lines := []jen.Code{
		jen.Comment("BuildLoginRequest builds an authenticating HTTP request"),
		jen.Line(),
		newClientMethod("BuildLoginRequest").Params(
			utils.CtxParam(),
			jen.List(
				jen.ID("username"),
				jen.ID("password"),
				jen.ID("totpToken"),
			).ID("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines
}

func buildLogin() []jen.Code {
	funcName := "Login"

	block := utils.StartSpan(true, funcName)
	block = append(block,
		jen.List(
			jen.ID("req"),
			jen.Err(),
		).Op(":=").ID("c").Dot("BuildLoginRequest").Call(
			utils.CtxVar(),
			jen.ID("username"),
			jen.ID("password"),
			jen.ID("totpToken"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.List(
			jen.ID("res"),
			jen.Err(),
		).Op(":=").ID("c").Dot("plainClient").Dot("Do").Call(
			jen.ID("req"),
		),
		jen.If(jen.Err().Op("!=").ID("nil")).Block(
			jen.Return().List(
				jen.ID("nil"),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encountered error executing login request: %w"),
					jen.Err(),
				),
			),
		),
		jen.Line(),
		jen.If(jen.ID("c").Dot("Debug")).Block(
			jen.List(
				jen.ID("b"),
				jen.Err(),
			).Op(":=").Qual("net/http/httputil", "DumpResponse").Call(
				jen.ID("res"),
				jen.ID("true"),
			),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("c").Dot("logger").Dot("Error").Call(
					jen.Err(),
					jen.Lit("dumping response"),
				),
			),
			jen.ID("c").Dot("logger").Dot("WithValue").Call(
				jen.Lit("response"),
				jen.ID("string").Call(
					jen.ID("b"),
				),
			).Dot("Debug").Call(
				jen.Lit("login response received"),
			),
		),
		jen.Line(),
		jen.Defer().Func().Params().Block(
			jen.If(jen.Err().Op(":=").ID("res").Dot("Body").Dot("Close").Call(),
				jen.Err().Op("!=").ID("nil"),
			).Block(
				jen.ID("c").Dot("logger").Dot("Error").Call(
					jen.Err(),
					jen.Lit("closing response body"),
				),
			),
		).Call(),
		jen.Line(),
		jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
		jen.If(jen.ID("len").Call(
			jen.ID("cookies"),
		).Op(">").Lit(0),
		).Block(
			jen.Return().List(jen.ID("cookies").Index(
				jen.Lit(0),
			),
				jen.ID("nil"),
			),
		),
		jen.Line(),
		jen.Return().List(
			jen.ID("nil"),
			jen.Qual("errors", "New").Call(
				jen.Lit("no cookies returned from request"),
			),
		),
	)

	lines := []jen.Code{
		jen.Comment("Login will, when provided the correct credentials, fetch a login cookie"),
		jen.Line(),
		newClientMethod("Login").Params(
			utils.CtxParam(),
			jen.List(
				jen.ID("username"),
				jen.ID("password"),
				jen.ID("totpToken"),
			).ID("string"),
		).Params(
			jen.Op("*").Qual("net/http", "Cookie"),
			jen.ID("error"),
		).Block(block...),
		jen.Line(),
	}

	return lines

}
