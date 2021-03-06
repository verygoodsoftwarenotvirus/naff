package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(jen.Const().ID("usersBasePath").Equals().Lit("users"))

	code.Add(buildBuildGetUserRequest(proj)...)
	code.Add(buildGetUser(proj)...)
	code.Add(buildBuildGetUsersRequest(proj)...)
	code.Add(buildGetUsers(proj)...)
	code.Add(buildBuildCreateUserRequest(proj)...)
	code.Add(buildCreateUser(proj)...)
	code.Add(buildBuildArchiveUserRequest(proj)...)
	code.Add(buildArchiveUser(proj)...)
	code.Add(buildBuildLoginRequest(proj)...)
	code.Add(buildLogin(proj)...)
	code.Add(buildBuildVerifyTOTPSecretRequest(proj)...)
	code.Add(buildVerifyTOTPSecret(proj)...)

	return code
}

func buildBuildGetUserRequest(proj *models.Project) []jen.Code {
	const funcName = "BuildGetUserRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID(constants.UserIDVarName),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching a user.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			constants.UserIDParam(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildGetUser(proj *models.Project) []jen.Code {
	const funcName = "GetUser"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetUserRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.UserIDVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
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
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID("user"),
		),
		jen.Return().List(
			jen.ID("user"),
			jen.Err(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves a user.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			constants.UserIDParam(),
		).Params(
			jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"),
			jen.Err().Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildGetUsersRequest(proj *models.Project) []jen.Code {
	const funcName = "BuildGetUsersRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.ID(constants.FilterVarName).Dot("ToValues").Call(),
			jen.ID("usersBasePath"),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodGet"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for fetching a user.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildGetUsers(proj *models.Project) []jen.Code {
	const funcName = "GetUsers"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("users").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserList").Values(),
		jen.Line(),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildGetUsersRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.FilterVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
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
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID("users"),
		),
		jen.Return().List(jen.ID("users"), jen.Err()),
	}

	lines := []jen.Code{
		jen.Commentf("%s retrieves a list of users.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "UserList"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateUserRequest(proj *models.Project) []jen.Code {
	const funcName = "BuildCreateUserRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("buildDataRequest").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for creating a user.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("body").PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildCreateUser(proj *models.Project) []jen.Code {
	const funcName = "CreateUser"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("user").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserCreationResponse").Values(),
		jen.Line(),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildCreateUserRequest").Call(
			constants.CtxVar(),
			jen.ID("input"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
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
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.AddressOf().ID("user"),
		),
		jen.Return().List(jen.ID("user"), jen.Err()),
	}

	lines := []jen.Code{
		jen.Commentf("%s creates a new user.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "UserCreationInput"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), "UserCreationResponse"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveUserRequest(proj *models.Project) []jen.Code {
	const funcName = "BuildArchiveUserRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
			jen.Qual("strconv", "FormatUint").Call(
				jen.ID(constants.UserIDVarName),
				jen.Lit(10),
			),
		),
		jen.Line(),
		jen.Return().Qual("net/http", "NewRequestWithContext").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodDelete"),
			jen.ID("uri"),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an HTTP request for updating a user.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			constants.UserIDParam(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildArchiveUser(proj *models.Project) []jen.Code {
	const funcName = "ArchiveUser"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.List(
			jen.ID(constants.RequestVarName),
			jen.Err(),
		).Assign().ID("c").Dot("BuildArchiveUserRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.UserIDVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("building request: %w"),
				jen.Err(),
			),
		),
		jen.Line(),
		jen.Return().ID("c").Dot("executeRequest").Call(
			constants.CtxVar(),
			jen.ID(constants.RequestVarName),
			jen.Nil(),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s archives a user.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			constants.UserIDParam(),
		).Params(jen.Error()).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildBuildLoginRequest(proj *models.Project) []jen.Code {
	const funcName = "BuildLoginRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.If(jen.ID("input").IsEqualTo().Nil()).Body(
			jen.Return(jen.Nil(), utils.Error("nil input provided")),
		),
		jen.Line(),
		jen.List(jen.ID("body"), jen.Err()).Assign().ID("createBodyFromStruct").Call(jen.AddressOf().ID("input")),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
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
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.ID("body"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds an authenticating HTTP request.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput"),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildLogin(proj *models.Project) []jen.Code {
	const funcName = "Login"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.If(jen.ID("input").IsEqualTo().Nil()).Body(
			jen.Return(jen.Nil(), utils.Error("nil input provided")),
		),
		jen.Line(),
		jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().ID("c").Dot("BuildLoginRequest").Call(constants.CtxVar(), jen.ID("input")),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(jen.Lit("error building login request: %w"), jen.Err()),
			),
		),
		jen.Line(),
		jen.List(
			jen.ID(constants.ResponseVarName),
			jen.Err(),
		).Assign().ID("c").Dot("plainClient").Dot("Do").Call(
			jen.ID(constants.RequestVarName),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encountered error executing login request: %w"),
					jen.Err(),
				),
			),
		),
		jen.ID("c").Dot("closeResponseBody").Call(jen.ID(constants.ResponseVarName)),
		jen.Line(),
		jen.ID("cookies").Assign().ID(constants.ResponseVarName).Dot("Cookies").Call(),
		jen.If(jen.Len(
			jen.ID("cookies"),
		).GreaterThan().Zero(),
		).Body(
			jen.Return().List(jen.ID("cookies").Index(
				jen.Zero(),
			),
				jen.Nil(),
			),
		),
		jen.Line(),
		jen.Return().List(
			jen.Nil(),
			utils.Error("no cookies returned from request"),
		),
	}

	lines := []jen.Code{
		jen.Commentf("%s will, when provided the correct credentials, fetch a login cookie.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			jen.ID("input").PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput"),
		).Params(
			jen.PointerTo().Qual("net/http", "Cookie"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines

}

func buildBuildVerifyTOTPSecretRequest(proj *models.Project) []jen.Code {
	const funcName = "BuildVerifyTOTPSecretRequest"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.ID("uri").Assign().ID("c").Dot("buildVersionlessURL").Call(
			jen.Nil(),
			jen.ID("usersBasePath"),
			jen.Lit("totp_secret"),
			jen.Lit("verify"),
		),
		jen.Line(),
		jen.Return(jen.ID("c").Dot("buildDataRequest").Call(
			constants.CtxVar(),
			jen.Qual("net/http", "MethodPost"),
			jen.ID("uri"),
			jen.AddressOf().Qual(proj.ModelsV1Package(), "TOTPSecretVerificationInput").Valuesln(
				jen.ID("TOTPToken").MapAssign().ID("token"),
				jen.ID(constants.UserIDFieldName).MapAssign().ID(constants.UserIDVarName),
			),
		)),
	}

	lines := []jen.Code{
		jen.Commentf("%s builds a request to validate a TOTP secret.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.ID("token").String(),
		).Params(
			jen.PointerTo().Qual("net/http", "Request"),
			jen.Error(),
		).Body(block...),
		jen.Line(),
	}

	return lines
}

func buildVerifyTOTPSecret(proj *models.Project) []jen.Code {
	const funcName = "VerifyTOTPSecret"

	block := []jen.Code{
		utils.StartSpan(proj, true, funcName),
		jen.Line(),
		jen.List(jen.ID("req"), jen.Err()).Assign().ID("c").Dot("BuildVerifyTOTPSecretRequest").Call(
			constants.CtxVar(),
			constants.UserIDVar(),
			jen.ID("token"),
		),
		jen.If(jen.Err().DoesNotEqual().Nil()).Body(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("error building TOTP validation request: %w"), jen.Err())),
		),
		jen.Line(),
		jen.List(jen.ID("res"), jen.Err()).Assign().ID("c").Dot("executeRawRequest").Call(
			constants.CtxVar(),
			jen.ID("c").Dot("plainClient"),
			jen.ID("req"),
		),
		jen.If(jen.Err().DoesNotEqual().Nil()).Body(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("executing request: %w"), jen.Err())),
		),
		jen.ID("c").Dot("closeResponseBody").Call(jen.ID("res")),
		jen.Line(),
		jen.If(jen.ID("res").Dot("StatusCode").IsEqualTo().Qual("net/http", "StatusBadRequest")).Body(
			jen.Return(jen.ID("ErrInvalidTOTPToken")),
		).Else().If(jen.ID("res").Dot("StatusCode").DoesNotEqual().Qual("net/http", "StatusAccepted")).Body(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("erroneous response code when validating TOTP secret: %d"), jen.ID("res").Dot("StatusCode"))),
		),
		jen.Line(),
		jen.Return(jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("%s executes a request to verify a TOTP secret.", funcName),
		jen.Line(),
		newClientMethod(funcName).Params(
			constants.CtxParam(),
			constants.UserIDParam(),
			jen.ID("token").String(),
		).Params(jen.Error()).Body(block...),
		jen.Line(),
	}

	return lines
}
