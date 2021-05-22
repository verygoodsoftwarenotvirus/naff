package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockUserDataServerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Underscore().Qual(proj.ModelsV1Package(), "UserDataServer").Equals().Parens(jen.PointerTo().ID("UserDataServer")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(buildUserDataServer()...)
	code.Add(buildUserLoginInputMiddleware()...)
	code.Add(buildUserInputMiddleware()...)
	code.Add(buildUserPasswordUpdateInputMiddleware()...)
	code.Add(buildUserTOTPSecretVerificationInputMiddleware()...)
	code.Add(buildUserTOTPSecretRefreshInputMiddleware()...)
	code.Add(buildUserListHandler()...)
	code.Add(buildUserCreateHandler()...)
	code.Add(buildUserReadHandler()...)
	code.Add(buildUserTOTPSecretVerificationHandler()...)
	code.Add(buildUserNewTOTPSecretHandler()...)
	code.Add(buildUserUpdatePasswordHandler()...)
	code.Add(buildUserArchiveHandler()...)

	return code
}

func buildUserDataServer() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UserDataServer is a mocked models.UserDataServer for testing"),
		jen.Line(),
		jen.Type().ID("UserDataServer").Struct(jen.Qual(constants.MockPkg, "Mock")),
		jen.Line(),
	}

	return lines
}

func buildUserLoginInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UserLoginInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("UserLoginInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildUserInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UserInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildUserPasswordUpdateInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("PasswordUpdateInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildUserTOTPSecretVerificationInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("TOTPSecretVerificationInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("TOTPSecretVerificationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildUserTOTPSecretRefreshInputMiddleware() []jen.Code {
	lines := []jen.Code{
		jen.Comment("TOTPSecretRefreshInputMiddleware is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(jen.ID("next")),
			jen.Return().ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Qual("net/http", "Handler")),
		),
		jen.Line(),
	}

	return lines
}

func buildUserListHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ListHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("ListHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID(constants.ResponseVarName),
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUserCreateHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("CreateHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("CreateHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID(constants.ResponseVarName),
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUserReadHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ReadHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("ReadHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID(constants.ResponseVarName),
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUserTOTPSecretVerificationHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("TOTPSecretVerificationHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("TOTPSecretVerificationHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID(constants.ResponseVarName),
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUserNewTOTPSecretHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("NewTOTPSecretHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("NewTOTPSecretHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID(constants.ResponseVarName),
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUserUpdatePasswordHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("UpdatePasswordHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("UpdatePasswordHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID(constants.ResponseVarName),
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUserArchiveHandler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ArchiveHandler is a mock method to satisfy our interface requirements."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("UserDataServer")).ID("ArchiveHandler").Params(
			jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params().Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID(constants.ResponseVarName),
				jen.ID(constants.RequestVarName),
			),
		),
		jen.Line(),
	}

	return lines
}
