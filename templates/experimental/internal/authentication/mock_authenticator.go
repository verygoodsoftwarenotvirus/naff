package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockAuthenticatorDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("Authenticator").Op("=").Parens(jen.Op("*").ID("MockAuthenticator")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("MockAuthenticator").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateLogin satisfies our authenticator interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockAuthenticator")).ID("ValidateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("hash"), jen.ID("password"), jen.ID("totpSecret"), jen.ID("totpCode")).ID("string")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("hash"),
				jen.ID("password"),
				jen.ID("totpSecret"),
				jen.ID("totpCode"),
			),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HashPassword satisfies our authenticator interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockAuthenticator")).ID("HashPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("password"),
			),
			jen.Return().List(jen.ID("args").Dot("String").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	return code
}
