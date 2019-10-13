package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func authenticatorDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("ErrInvalidTwoFactorCode").Op("=").Qual("errors", "New").Call(jen.Lit("invalid two factor code")).Var().ID("ErrPasswordHashTooWeak").Op("=").Qual("errors", "New").Call(jen.Lit("password's hash is too weak")).Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideBcryptAuthenticator"), jen.ID("ProvideBcryptHashCost")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideBcryptHashCost provides a BcryptHashCost").ID("ProvideBcryptHashCost").Params().Params(jen.ID("BcryptHashCost")).Block(
		jen.Return().ID("DefaultBcryptHashCost"),
	),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("PasswordHasher").Interface(jen.ID("PasswordIsAcceptable").Params(jen.ID("password").ID("string")).Params(jen.ID("bool")), jen.ID("HashPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")), jen.ID("PasswordMatches").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword")).ID("string"), jen.ID("salt").Index().ID("byte")).Params(jen.ID("bool"))).Type().ID("Authenticator").Interface(jen.ID("PasswordHasher"), jen.ID("ValidateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("HashedPassword"), jen.ID("ProvidedPassword"), jen.ID("TwoFactorSecret"), jen.ID("TwoFactorCode")).ID("string"), jen.ID("Salt").Index().ID("byte")).Params(jen.ID("valid").ID("bool"), jen.ID("err").ID("error"))),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// we run this function to ensure that we have no problem reading from crypto/rand").ID("init").Params().Block(
		jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
		jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
	),

		jen.Line(),
	)
	return ret
}
