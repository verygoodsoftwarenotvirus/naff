package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func argon2DotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("argon2IterationCount").Op("=").Lit(1),
			jen.ID("argon2ThreadCount").Op("=").Lit(2),
			jen.ID("argon2SaltLength").Op("=").Lit(16),
			jen.ID("argon2KeyLength").Op("=").Lit(32),
			jen.ID("sixtyFourMegabytes").Op("=").Lit(64).Op("*").Lit(1024),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("argonParams").Op("=").Op("&").ID("argon2id").Dot("Params").Valuesln(jen.ID("Memory").Op(":").ID("sixtyFourMegabytes"), jen.ID("Iterations").Op(":").ID("argon2IterationCount"), jen.ID("Parallelism").Op(":").ID("argon2ThreadCount"), jen.ID("SaltLength").Op(":").ID("argon2SaltLength"), jen.ID("KeyLength").Op(":").ID("argon2KeyLength")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("Argon2Authenticator").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideArgon2Authenticator returns an argon2 powered Argon2Authenticator."),
		jen.Line(),
		jen.Func().ID("ProvideArgon2Authenticator").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("Authenticator")).Body(
			jen.ID("ba").Op(":=").Op("&").ID("Argon2Authenticator").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("argon2Provider")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("argon2Provider"))),
			jen.Return().ID("ba"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HashPassword takes a password and hashes it using argon2."),
		jen.Line(),
		jen.Func().Params(jen.ID("a").Op("*").ID("Argon2Authenticator")).ID("HashPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("a").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Return().ID("argon2id").Dot("CreateHash").Call(
				jen.ID("password"),
				jen.ID("argonParams"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateLogin validates a login attempt by:"),
		jen.Line(),
		jen.Comment(" - checking that the provided authentication matches the provided hashed passwords."),
		jen.Line(),
		jen.Comment(" - checking that the temporary one-time authentication provided jives with the provided two factor secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("a").Op("*").ID("Argon2Authenticator")).ID("ValidateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hash"), jen.ID("password"), jen.ID("totpSecret"), jen.ID("totpCode")).ID("string")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("a").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("a").Dot("logger"),
			jen.List(jen.ID("passwordMatches"), jen.ID("err")).Op(":=").ID("argon2id").Dot("ComparePasswordAndHash").Call(
				jen.ID("password"),
				jen.ID("hash"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("a").Dot("logger"),
					jen.ID("span"),
					jen.Lit("comparing argon2 hashed password"),
				))).Else().If(jen.Op("!").ID("passwordMatches")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("ErrPasswordDoesNotMatch"))),
			jen.If(jen.Op("!").ID("totp").Dot("Validate").Call(
				jen.ID("totpCode"),
				jen.ID("totpSecret"),
			)).Body(
				jen.ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("password_matches").Op(":").ID("passwordMatches"), jen.Lit("provided_code").Op(":").ID("totpCode"))).Dot("Debug").Call(jen.Lit("invalid code provided")),
				jen.Return().List(jen.ID("passwordMatches"), jen.ID("ErrInvalidTOTPToken")),
			),
			jen.Return().List(jen.ID("passwordMatches"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
