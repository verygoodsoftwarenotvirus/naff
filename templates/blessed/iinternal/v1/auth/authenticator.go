package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authenticatorDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(
		jen.Var().Defs(
			jen.Comment("ErrInvalidTwoFactorCode indicates that a provided two factor code is invalid"),
			jen.ID("ErrInvalidTwoFactorCode").Equals().Qual("errors", "New").Call(jen.Lit("invalid two factor code")),
			jen.Comment("ErrPasswordHashTooWeak indicates that a provided password hash is too weak"),
			jen.ID("ErrPasswordHashTooWeak").Equals().Qual("errors", "New").Call(jen.Lit("password's hash is too weak")),
			jen.Line(),
			jen.Comment("Providers represents what this package offers to external libraries in the way of constructors"),
			jen.ID("Providers").Equals().Qual("github.com/google/wire", "NewSet").Callln(jen.ID("ProvideBcryptAuthenticator"), jen.ID("ProvideBcryptHashCost")),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Comment("ProvideBcryptHashCost provides a BcryptHashCost"),
		jen.Line(),
		jen.Func().ID("ProvideBcryptHashCost").Params().Params(jen.ID("BcryptHashCost")).Block(
			jen.Return().ID("DefaultBcryptHashCost"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("PasswordHasher hashes passwords"),
			jen.ID("PasswordHasher").Interface(
				jen.ID("PasswordIsAcceptable").Params(jen.ID("password").ID("string")).Params(jen.ID("bool")),
				jen.ID("HashPassword").Params(utils.CtxParam(), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")),
				jen.ID("PasswordMatches").Params(utils.CtxParam(), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword")).ID("string"), jen.ID("salt").Index().ID("byte")).Params(jen.ID("bool")),
			),
			jen.Line(),
			jen.Comment("Authenticator is a poorly named Authenticator interface"),
			jen.ID("Authenticator").Interface(
				jen.ID("PasswordHasher"),
				jen.Line(),
				jen.ID("ValidateLogin").Paramsln(
					utils.CtxParam(),
					jen.Listln(
						jen.ID("HashedPassword"),
						jen.ID("ProvidedPassword"),
						jen.ID("TwoFactorSecret"),
						jen.ID("TwoFactorCode"),
					).ID("string"),
					jen.ID("Salt").Index().ID("byte"),
				).Params(jen.ID("valid").ID("bool"), jen.Err().ID("error")),
			),
		),
		jen.Line(),
	)
	ret.Add(
		jen.Comment("we run this function to ensure that we have no problem reading from crypto/rand"),
		jen.Line(),
		jen.Func().ID("init").Params().Block(
			jen.ID("b").Assign().ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
		),

		jen.Line(),
	)
	return ret
}
