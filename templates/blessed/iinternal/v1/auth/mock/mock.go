package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "Authenticator").Op("=").Parens(jen.Op("*").ID("Authenticator")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Authenticator is a mock Authenticator"), jen.Line(),
		jen.Type().ID("Authenticator").Struct(jen.ID("mock").Dot(
			"Mock",
		)),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ValidateLogin satisfies our authenticator interface"), jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("Authenticator")).ID("ValidateLogin").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.Listln(jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("twoFactorSecret"),
				jen.ID("twoFactorCode")).ID("string"),
			jen.ID("salt").Index().ID("byte")).Params(jen.ID("valid").ID("bool"), jen.ID("err").ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Callln(
				jen.ID("ctx"),
				jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("twoFactorSecret"),
				jen.ID("twoFactorCode"),
				jen.ID("salt"),
			),
			jen.Return().List(jen.ID("args").Dot("Bool").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("PasswordIsAcceptable satisfies our authenticator interface"), jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("Authenticator")).ID("PasswordIsAcceptable").Params(jen.ID("password").ID("string")).Params(jen.ID("bool")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("password")).Dot("Bool").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HashPassword satisfies our authenticator interface"), jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("Authenticator")).ID("HashPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("password").ID("string")).Params(jen.ID("string"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("password")),
			jen.Return().List(jen.ID("args").Dot("String").Call(jen.Lit(0)), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("PasswordMatches satisfies our authenticator interface"), jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("Authenticator")).ID("PasswordMatches").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.Listln(jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
			).ID("string"),
			jen.ID("salt").Index().ID("byte")).Params(jen.ID("bool")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Callln(
				jen.ID("ctx"),
				jen.ID("hashedPassword"),
				jen.ID("providedPassword"),
				jen.ID("salt"),
			),
			jen.Return().ID("args").Dot("Bool").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	return ret
}
