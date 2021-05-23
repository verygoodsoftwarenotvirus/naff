package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authenticatorDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildAuthenticatorVariableDeclarations()...)
	code.Add(buildProvideBcryptHashCost()...)
	code.Add(buildAuthenticatorTypeDefinitions()...)
	code.Add(buildInit()...)

	return code
}

func buildAuthenticatorVariableDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("ErrInvalidTwoFactorCode indicates that a provided two factor code is invalid."),
			jen.ID("ErrInvalidTwoFactorCode").Equals().Qual("errors", "New").Call(jen.Lit("invalid two factor code")),
			jen.Comment("ErrPasswordHashTooWeak indicates that a provided password hash is too weak."),
			jen.ID("ErrPasswordHashTooWeak").Equals().Qual("errors", "New").Call(jen.Lit("password's hash is too weak")),
			jen.Line(),
			jen.Comment("Providers represents what this package offers to external libraries in the way of constructors."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(jen.ID("ProvideBcryptAuthenticator"), jen.ID("ProvideBcryptHashCost")),
			jen.Line(),
		),
	}

	return lines
}

func buildProvideBcryptHashCost() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideBcryptHashCost provides a BcryptHashCost."),
		jen.Line(),
		jen.Func().ID("ProvideBcryptHashCost").Params().Params(jen.ID("BcryptHashCost")).Body(
			jen.Return().ID("DefaultBcryptHashCost"),
		),
		jen.Line(),
	}

	return lines
}

func buildAuthenticatorTypeDefinitions() []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("PasswordHasher hashes passwords."),
			jen.ID("PasswordHasher").Interface(
				jen.ID("PasswordIsAcceptable").Params(jen.ID("password").String()).Params(jen.Bool()),
				jen.ID("HashPassword").Params(constants.CtxParam(), jen.ID("password").String()).Params(jen.String(), jen.Error()),
				jen.ID("PasswordMatches").Params(constants.CtxParam(), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword")).String(), jen.ID("salt").Index().Byte()).Params(jen.Bool()),
			),
			jen.Line(),
			jen.Comment("Authenticator is a poorly named Authenticator interface."),
			jen.ID("Authenticator").Interface(
				jen.ID("PasswordHasher"),
				jen.Line(),
				jen.ID("ValidateLogin").Paramsln(
					constants.CtxParam(),
					jen.Listln(
						jen.ID("HashedPassword"),
						jen.ID("ProvidedPassword"),
						jen.ID("TwoFactorSecret"),
						jen.ID("TwoFactorCode"),
					).String(),
					jen.ID("Salt").Index().Byte(),
				).Params(jen.ID("valid").Bool(), jen.Err().Error()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildInit() []jen.Code {
	lines := []jen.Code{
		jen.Comment("we run this function to ensure that we have no problem reading from crypto/rand"),
		jen.Line(),
		jen.Func().ID("init").Params().Body(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
		),

		jen.Line(),
	}

	return lines
}
