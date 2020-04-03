package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	loggingImport = "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

func bcryptDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("bcryptCostCompensation").Equals().Lit(2),
			jen.ID("defaultMinimumPasswordSize").Equals().Lit(16),
			jen.Line(),
			jen.Comment("DefaultBcryptHashCost is what it says on the tin"),
			jen.ID("DefaultBcryptHashCost").Equals().ID("BcryptHashCost").Call(jen.Qual("golang.org/x/crypto/bcrypt", "DefaultCost").Op("+").ID("bcryptCostCompensation")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").ID("Authenticator").Equals().Parens(jen.PointerTo().ID("BcryptAuthenticator")).Call(jen.Nil()),
			jen.Line(),
			jen.Comment("ErrCostTooLow indicates that a password has too low a Bcrypt cost"),
			jen.ID("ErrCostTooLow").Equals().Qual("errors", "New").Call(jen.Lit("stored password's cost is too low")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("BcryptAuthenticator is our bcrypt-based authenticator"),
			jen.ID("BcryptAuthenticator").Struct(
				jen.ID("logger").Qual(loggingImport, "Logger"),
				jen.ID("hashCost").ID("uint"),
				jen.ID("minimumPasswordSize").ID("uint"),
			),
			jen.Line(),
			jen.Comment("BcryptHashCost is an arbitrary type alias for dependency injection's sake."),
			jen.ID("BcryptHashCost").ID("uint"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideBcryptAuthenticator returns a bcrypt powered Authenticator"),
		jen.Line(),
		jen.Func().ID("ProvideBcryptAuthenticator").Params(jen.ID("hashCost").ID("BcryptHashCost"), jen.ID("logger").Qual(loggingImport, "Logger")).Params(jen.ID("Authenticator")).Block(
			jen.ID("ba").Assign().VarPointer().ID("BcryptAuthenticator").Valuesln(
				jen.ID("logger").MapAssign().ID("logger").Dot("WithName").Call(jen.Lit("bcrypt")),
				jen.ID("hashCost").MapAssign().ID("uint").Call(jen.Qual("math", "Min").Call(jen.ID("float64").Call(jen.ID("DefaultBcryptHashCost")), jen.ID("float64").Call(jen.ID("hashCost")))),
				jen.ID("minimumPasswordSize").MapAssign().ID("defaultMinimumPasswordSize"),
			),
			jen.Return().ID("ba"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HashPassword takes a password and hashes it using bcrypt"),
		jen.Line(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("BcryptAuthenticator")).ID("HashPassword").Params(jen.ID("c").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.Error()).Block(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("c"), jen.Lit("HashPassword")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.List(jen.ID("hashedPass"), jen.Err()).Assign().Qual("golang.org/x/crypto/bcrypt", "GenerateFromPassword").Call(jen.Index().ID("byte").Call(jen.ID("password")), jen.ID("int").Call(jen.ID("b").Dot("hashCost"))),
			jen.Return().List(jen.ID("string").Call(jen.ID("hashedPass")), jen.Err()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ValidateLogin validates a login attempt by"),
		jen.Line(),
		jen.Comment("1. checking that the provided password matches the stored hashed password"),
		jen.Line(),
		jen.Comment("2. checking that the temporary one-time password provided jives with the stored two factor secret"),
		jen.Line(),
		jen.Comment("3. checking that the provided hashed password isn't too weak, and returning an error otherwise"),
		jen.Line(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("BcryptAuthenticator")).ID("ValidateLogin").Paramsln(
			utils.CtxParam(),
			jen.Listln(jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("twoFactorSecret"),
				jen.ID("twoFactorCode")).ID("string"),
			jen.ID("salt").Index().ID("byte"),
		).Params(jen.ID("passwordMatches").ID("bool"), jen.Err().ID("error")).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("ValidateLogin")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("passwordMatches").Equals().ID("b").Dot("PasswordMatches").Call(utils.CtxVar(), jen.ID("hashedPassword"), jen.ID("providedPassword"), jen.Nil()),
			jen.ID("tooWeak").Assign().ID("b").Dot("hashedPasswordIsTooWeak").Call(jen.ID("hashedPassword")),
			jen.Line(),
			jen.If(jen.Op("!").Qual("github.com/pquerna/otp/totp", "Validate").Call(jen.ID("twoFactorCode"), jen.ID("twoFactorSecret"))).Block(
				jen.ID("b").Dot("logger").Dot("WithValues").Call(
					jen.Map(jen.ID("string")).Interface().Valuesln(
						jen.Lit("password_matches").MapAssign().ID("passwordMatches"),
						jen.Lit("2fa_secret").MapAssign().ID("twoFactorSecret"),
						jen.Lit("provided_code").MapAssign().ID("twoFactorCode"),
					),
				).Dot("Debug").Call(jen.Lit("invalid code provided")),
				jen.Line(),
				jen.Return().List(jen.ID("passwordMatches"), jen.ID("ErrInvalidTwoFactorCode")),
			),
			jen.Line(),
			jen.If(jen.ID("tooWeak")).Block(
				jen.Comment("NOTE: this can end up with a return set where passwordMatches is true and the err is not nil."),
				jen.Comment("This is the valid case in the event the user has logged in with a valid password, but the"),
				jen.Comment("bcrypt cost has been raised since they last logged in."),
				jen.Return().List(jen.ID("passwordMatches"), jen.ID("ErrCostTooLow")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("passwordMatches"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("PasswordMatches validates whether or not a bcrypt-hashed password matches a provided password"),
		jen.Line(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("BcryptAuthenticator")).ID("PasswordMatches").Params(utils.CtxParam(), jen.List(jen.ID("hashedPassword"), jen.ID("providedPassword")).ID("string"), jen.ID("_").Index().ID("byte")).Params(jen.ID("bool")).Block(
			jen.Return().Qual("golang.org/x/crypto/bcrypt", "CompareHashAndPassword").Call(jen.Index().ID("byte").Call(jen.ID("hashedPassword")), jen.Index().ID("byte").Call(jen.ID("providedPassword"))).Op("==").ID("nil"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("hashedPasswordIsTooWeak determines if a given hashed password was hashed with too weak a bcrypt cost"),
		jen.Line(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("BcryptAuthenticator")).ID("hashedPasswordIsTooWeak").Params(jen.ID("hashedPassword").ID("string")).Params(jen.ID("bool")).Block(
			jen.List(jen.ID("cost"), jen.Err()).Assign().Qual("golang.org/x/crypto/bcrypt", "Cost").Call(jen.Index().ID("byte").Call(jen.ID("hashedPassword"))),
			jen.Line(),
			jen.Return().ID("err").DoesNotEqual().ID("nil").Op("||").ID("uint").Call(jen.ID("cost")).Op("<").ID("b").Dot("hashCost"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("PasswordIsAcceptable takes a password and returns whether or not it satisfies the authenticator"),
		jen.Line(),
		jen.Func().Params(jen.ID("b").PointerTo().ID("BcryptAuthenticator")).ID("PasswordIsAcceptable").Params(jen.ID("pass").ID("string")).Params(jen.ID("bool")).Block(
			jen.Return().ID("uint").Call(jen.ID("len").Call(jen.ID("pass"))).Op(">=").ID("b").Dot("minimumPasswordSize"),
		),
		jen.Line(),
	)
	return ret
}
