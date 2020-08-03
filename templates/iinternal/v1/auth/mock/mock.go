package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("mock")

	utils.AddImports(proj, code)

	code.Add(buildInterfaceImplementationDeclaration(proj)...)
	code.Add(buildMockAuthenticator()...)
	code.Add(buildMockValidateLogin()...)
	code.Add(buildMockPasswordIsAcceptable()...)
	code.Add(buildMockHashPassword()...)
	code.Add(buildMockPasswordMatches()...)

	return code
}

func buildInterfaceImplementationDeclaration(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual(proj.InternalAuthV1Package(), "Authenticator").Equals().Parens(jen.PointerTo().ID("Authenticator")).Call(jen.Nil()),
		jen.Line(),
	}

	return lines
}

func buildMockAuthenticator() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Authenticator is a mock Authenticator."), jen.Line(),
		jen.Type().ID("Authenticator").Struct(jen.Qual(constants.MockPkg,
			"Mock",
		)),
		jen.Line(),
	}

	return lines
}

func buildMockValidateLogin() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ValidateLogin satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("ValidateLogin").Paramsln(
			constants.CtxParam(),
			jen.Listln(jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("twoFactorSecret"),
				jen.ID("twoFactorCode")).String(),
			jen.ID("salt").Index().Byte()).Params(jen.ID("valid").Bool(), jen.Err().Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Callln(
				constants.CtxVar(),
				jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("twoFactorSecret"),
				jen.ID("twoFactorCode"),
				jen.ID("salt"),
			),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Zero()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockPasswordIsAcceptable() []jen.Code {
	lines := []jen.Code{
		jen.Comment("PasswordIsAcceptable satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("PasswordIsAcceptable").Params(jen.ID("password").String()).Params(jen.Bool()).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("password")).Dot("Bool").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildMockHashPassword() []jen.Code {
	lines := []jen.Code{
		jen.Comment("HashPassword satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("HashPassword").Params(constants.CtxParam(), jen.ID("password").String()).Params(jen.String(), jen.Error()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("password")),
			jen.Return().List(jen.ID("args").Dot("String").Call(jen.Zero()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	}

	return lines
}

func buildMockPasswordMatches() []jen.Code {
	lines := []jen.Code{
		jen.Comment("PasswordMatches satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("PasswordMatches").Paramsln(
			constants.CtxParam(),
			jen.Listln(jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
			).String(),
			jen.ID("salt").Index().Byte()).Params(jen.Bool()).Body(
			jen.ID("args").Assign().ID("m").Dot("Called").Callln(
				constants.CtxVar(),
				jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("salt"),
			),
			jen.Return().ID("args").Dot("Bool").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
