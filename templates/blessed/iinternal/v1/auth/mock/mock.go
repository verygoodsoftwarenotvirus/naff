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

	code.Add(
		jen.Var().Underscore().Qual(proj.InternalAuthV1Package(), "Authenticator").Equals().Parens(jen.PointerTo().ID("Authenticator")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Authenticator is a mock Authenticator."), jen.Line(),
		jen.Type().ID("Authenticator").Struct(jen.Qual(constants.MockPkg,
			"Mock",
		)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateLogin satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("ValidateLogin").Paramsln(
			constants.CtxParam(),
			jen.Listln(jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("twoFactorSecret"),
				jen.ID("twoFactorCode")).String(),
			jen.ID("salt").Index().Byte()).Params(jen.ID("valid").Bool(), jen.Err().Error()).Block(
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
	)

	code.Add(
		jen.Comment("PasswordIsAcceptable satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("PasswordIsAcceptable").Params(jen.ID("password").String()).Params(jen.Bool()).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("password")).Dot("Bool").Call(jen.Zero()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HashPassword satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("HashPassword").Params(constants.CtxParam(), jen.ID("password").String()).Params(jen.String(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(constants.CtxVar(), jen.ID("password")),
			jen.Return().List(jen.ID("args").Dot("String").Call(jen.Zero()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PasswordMatches satisfies our authenticator interface."), jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("Authenticator")).ID("PasswordMatches").Paramsln(
			constants.CtxParam(),
			jen.Listln(jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
			).String(),
			jen.ID("salt").Index().Byte()).Params(jen.Bool()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Callln(
				constants.CtxVar(),
				jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("salt"),
			),
			jen.Return().ID("args").Dot("Bool").Call(jen.Zero()),
		),
		jen.Line(),
	)

	return code
}
