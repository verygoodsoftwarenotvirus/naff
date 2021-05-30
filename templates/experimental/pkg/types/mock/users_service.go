package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("UsersService").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserRegistrationInputMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("UserRegistrationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PasswordUpdateInputMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TOTPSecretRefreshInputMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TOTPSecretVerificationInputMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("TOTPSecretVerificationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ListHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuditEntryHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SelfHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("SelfHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UsernameSearchHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("UsernameSearchHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewTOTPSecretHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("NewTOTPSecretHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TOTPSecretVerificationHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("TOTPSecretVerificationHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdatePasswordHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("UpdatePasswordHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AvatarUploadHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("AvatarUploadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("res"),
				jen.ID("req"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RegisterUser satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("RegisterUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("registrationInput").Op("*").ID("types").Dot("UserRegistrationInput")).Params(jen.Op("*").ID("types").Dot("UserCreationResponse"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("registrationInput"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("UserCreationResponse")), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("VerifyUserTwoFactorSecret satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UsersService")).ID("VerifyUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
