package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func bcryptDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("bcryptCostCompensation").Op("=").Lit(2).Var().ID("defaultMinimumPasswordSize").Op("=").Lit(16).Var().ID("DefaultBcryptHashCost").Op("=").ID("BcryptHashCost").Call(jen.ID("bcrypt").Dot(
		"DefaultCost",
	).Op("+").ID("bcryptCostCompensation")),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("_").ID("Authenticator").Op("=").Parens(jen.Op("*").ID("BcryptAuthenticator")).Call(jen.ID("nil")).Var().ID("ErrCostTooLow").Op("=").ID("errors").Dot(
		"New",
	).Call(jen.Lit("stored password's cost is too low")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("BcryptAuthenticator").Struct(jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("hashCost").ID("uint"), jen.ID("minimumPasswordSize").ID("uint")).Type().ID("BcryptHashCost").ID("uint"),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideBcryptAuthenticator returns a bcrypt powered Authenticator").ID("ProvideBcryptAuthenticator").Params(jen.ID("hashCost").ID("BcryptHashCost"), jen.ID("logger").ID("logging").Dot(
		"Logger",
	)).Params(jen.ID("Authenticator")).Block(
		jen.ID("ba").Op(":=").Op("&").ID("BcryptAuthenticator").Valuesln(jen.ID("logger").Op(":").ID("logger").Dot(
			"WithName",
		).Call(jen.Lit("bcrypt")), jen.ID("hashCost").Op(":").ID("uint").Call(jen.Qual("math", "Min").Call(jen.ID("float64").Call(jen.ID("DefaultBcryptHashCost")), jen.ID("float64").Call(jen.ID("hashCost")))), jen.ID("minimumPasswordSize").Op(":").ID("defaultMinimumPasswordSize")),
		jen.Return().ID("ba"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// HashPassword takes a password and hashes it using bcrypt").Params(jen.ID("b").Op("*").ID("BcryptAuthenticator")).ID("HashPassword").Params(jen.ID("c").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")).Block(
		jen.List(jen.ID("_"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("c"), jen.Lit("HashPassword")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.List(jen.ID("hashedPass"), jen.ID("err")).Op(":=").ID("bcrypt").Dot(
			"GenerateFromPassword",
		).Call(jen.Index().ID("byte").Call(jen.ID("password")), jen.ID("int").Call(jen.ID("b").Dot(
			"hashCost",
		))),
		jen.Return().List(jen.ID("string").Call(jen.ID("hashedPass")), jen.ID("err")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ValidateLogin validates a login attempt by").Comment("// 1. checking that the provided password matches the stored hashed password").Comment("// 2. checking that the temporary one-time password provided jives with the stored two factor secret").Comment("// 3. checking that the provided hashed password isn't too weak, and returning an error otherwise").Params(jen.ID("b").Op("*").ID("BcryptAuthenticator")).ID("ValidateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword"), jen.ID("twoFactorSecret"), jen.ID("twoFactorCode")).ID("string"), jen.ID("salt").Index().ID("byte")).Params(jen.ID("passwordMatches").ID("bool"), jen.ID("err").ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("ValidateLogin")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("passwordMatches").Op("=").ID("b").Dot(
			"PasswordMatches",
		).Call(jen.ID("ctx"), jen.ID("hashedPassword"), jen.ID("providedPassword"), jen.ID("nil")),
		jen.ID("tooWeak").Op(":=").ID("b").Dot(
			"hashedPasswordIsTooWeak",
		).Call(jen.ID("hashedPassword")),
		jen.If(jen.Op("!").ID("totp").Dot(
			"Validate",
		).Call(jen.ID("twoFactorCode"), jen.ID("twoFactorSecret"))).Block(
			jen.ID("b").Dot(
				"logger",
			).Dot(
				"WithValues",
			).Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("password_matches").Op(":").ID("passwordMatches"), jen.Lit("2fa_secret").Op(":").ID("twoFactorSecret"), jen.Lit("provided_code").Op(":").ID("twoFactorCode"))).Dot(
				"Debug",
			).Call(jen.Lit("invalid code provided")),
			jen.Return().List(jen.ID("passwordMatches"), jen.ID("ErrInvalidTwoFactorCode")),
		),
		jen.If(jen.ID("tooWeak")).Block(
			jen.Return().List(jen.ID("passwordMatches"), jen.ID("ErrCostTooLow")),
		),
		jen.Return().List(jen.ID("passwordMatches"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// PasswordMatches validates whether or not a bcrypt-hashed password matches a provided password").Params(jen.ID("b").Op("*").ID("BcryptAuthenticator")).ID("PasswordMatches").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword")).ID("string"), jen.ID("_").Index().ID("byte")).Params(jen.ID("bool")).Block(
		jen.Return().ID("bcrypt").Dot(
			"CompareHashAndPassword",
		).Call(jen.Index().ID("byte").Call(jen.ID("hashedPassword")), jen.Index().ID("byte").Call(jen.ID("providedPassword"))).Op("==").ID("nil"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// hashedPasswordIsTooWeak determines if a given hashed password was hashed with too weak a bcrypt cost").Params(jen.ID("b").Op("*").ID("BcryptAuthenticator")).ID("hashedPasswordIsTooWeak").Params(jen.ID("hashedPassword").ID("string")).Params(jen.ID("bool")).Block(
		jen.List(jen.ID("cost"), jen.ID("err")).Op(":=").ID("bcrypt").Dot(
			"Cost",
		).Call(jen.Index().ID("byte").Call(jen.ID("hashedPassword"))),
		jen.Return().ID("err").Op("!=").ID("nil").Op("||").ID("uint").Call(jen.ID("cost")).Op("<").ID("b").Dot(
			"hashCost",
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// PasswordIsAcceptable takes a password and returns whether or not it satisfies the authenticator").Params(jen.ID("b").Op("*").ID("BcryptAuthenticator")).ID("PasswordIsAcceptable").Params(jen.ID("pass").ID("string")).Params(jen.ID("bool")).Block(
		jen.Return().ID("uint").Call(jen.ID("len").Call(jen.ID("pass"))).Op(">=").ID("b").Dot(
			"minimumPasswordSize",
		),
	),

		jen.Line(),
	)
	return ret
}
