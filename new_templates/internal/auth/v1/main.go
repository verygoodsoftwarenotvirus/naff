package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mainDotGo() *jen.File {
	ret := jen.NewFile("auth")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("ErrInvalidTwoFactorCode").Op("=").Qual("errors", "New").Call(jen.Lit("invalid two factor code")).Var().ID("ErrPasswordHashTooWeak").Op("=").Qual("errors", "New").Call(jen.Lit("password's hash is too weak")).Var().ID("Providers").Op("=").ID("wire").Dot(
		"NewSet",
	).Call(jen.ID("ProvideBcryptAuthenticator"), jen.ID("ProvideBcryptHashCost")),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Null().Type().ID("PasswordHasher").Interface(jen.ID("PasswordIsAcceptable").Params(jen.ID("password").ID("string")).Params(jen.ID("bool")), jen.ID("HashPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")), jen.ID("PasswordMatches").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword")).ID("string"), jen.ID("salt").Index().ID("byte")).Params(jen.ID("bool"))).Type().ID("Authenticator").Interface(jen.ID("PasswordHasher"), jen.ID("ValidateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("HashedPassword"), jen.ID("ProvidedPassword"), jen.ID("TwoFactorSecret"), jen.ID("TwoFactorCode")).ID("string"), jen.ID("Salt").Index().ID("byte")).Params(jen.ID("valid").ID("bool"), jen.ID("err").ID("error"))),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
