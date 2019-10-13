package mock

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func mockDotGo() *jen.File {
	ret := jen.NewFile("mock")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("auth").Dot(
		"Authenticator",
	).Op("=").Parens(jen.Op("*").ID("Authenticator")).Call(jen.ID("nil")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("Authenticator").Struct(jen.ID("mock").Dot(
		"Mock",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ValidateLogin satisfies our authenticator interface").Params(jen.ID("m").Op("*").ID("Authenticator")).ID("ValidateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword"), jen.ID("twoFactorSecret"), jen.ID("twoFactorCode")).ID("string"), jen.ID("salt").Index().ID("byte")).Params(jen.ID("valid").ID("bool"), jen.ID("err").ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("hashedPassword"), jen.ID("providedPassword"), jen.ID("twoFactorSecret"), jen.ID("twoFactorCode"), jen.ID("salt")),
		jen.Return().List(jen.ID("args").Dot(
			"Bool",
		).Call(jen.Lit(0)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// PasswordIsAcceptable satisfies our authenticator interface").Params(jen.ID("m").Op("*").ID("Authenticator")).ID("PasswordIsAcceptable").Params(jen.ID("password").ID("string")).Params(jen.ID("bool")).Block(
		jen.Return().ID("m").Dot(
			"Called",
		).Call(jen.ID("password")).Dot(
			"Bool",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// HashPassword satisfies our authenticator interface").Params(jen.ID("m").Op("*").ID("Authenticator")).ID("HashPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("password")),
		jen.Return().List(jen.ID("args").Dot(
			"String",
		).Call(jen.Lit(0)), jen.ID("args").Dot(
			"Error",
		).Call(jen.Lit(1))),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// PasswordMatches satisfies our authenticator interface").Params(jen.ID("m").Op("*").ID("Authenticator")).ID("PasswordMatches").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword")).ID("string"), jen.ID("salt").Index().ID("byte")).Params(jen.ID("bool")).Block(
		jen.ID("args").Op(":=").ID("m").Dot(
			"Called",
		).Call(jen.ID("ctx"), jen.ID("hashedPassword"), jen.ID("providedPassword"), jen.ID("salt")),
		jen.Return().ID("args").Dot(
			"Bool",
		).Call(jen.Lit(0)),
	),

		jen.Line(),
	)
	return ret
}
