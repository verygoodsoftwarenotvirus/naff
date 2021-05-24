package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authenticatorDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.Lit(64),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("ErrInvalidTOTPToken").Op("=").Qual("errors", "New").Call(jen.Lit("invalid two factor code")).Var().ID("ErrPasswordDoesNotMatch").Op("=").Qual("errors", "New").Call(jen.Lit("password does not match")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("Hasher").Interface(jen.ID("HashPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error"))).Type().ID("Authenticator").Interface(
			jen.ID("Hasher"),
			jen.ID("ValidateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hash"), jen.ID("password"), jen.ID("totpSecret"), jen.ID("totpCode")).ID("string")).Params(jen.ID("bool"), jen.ID("error")),
		),
		jen.Line(),
	)

	return code
}
