package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func bcryptDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("bcryptCostCompensation").Op("=").Lit(2).Var().ID("defaultMinimumPasswordSize").Op("=").Lit(16).Var().ID("DefaultBcryptHashCost").Op("=").ID("BcryptHashCost").Call(jen.ID("bcrypt").Dot(
		"DefaultCost",
	).Op("+").ID("bcryptCostCompensation")),
	)
	ret.Add(jen.Null().Var().ID("_").ID("Authenticator").Op("=").Parens(jen.Op("*").ID("BcryptAuthenticator")).Call(jen.ID("nil")).Var().ID("ErrCostTooLow").Op("=").ID("errors").Dot(
		"New",
	).Call(jen.Lit("stored password's cost is too low")),
	)
	ret.Add(jen.Null().Type().ID("BcryptAuthenticator").Struct(
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("hashCost").ID("uint"),
		jen.ID("minimumPasswordSize").ID("uint"),
	).Type().ID("BcryptHashCost").ID("uint"),
	)
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	return ret
}
