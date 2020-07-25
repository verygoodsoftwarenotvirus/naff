package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
)

func randTestDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.Add(
		jen.Var().Underscore().ID("secretGenerator").Equals().Parens(jen.PointerTo().ID("mockSecretGenerator")).Call(jen.Nil()),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("mockSecretGenerator").Struct(
			jen.Qual(constants.MockPkg, "Mock"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockSecretGenerator")).ID("GenerateTwoFactorSecret").Params().Params(jen.String(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Line(),
			jen.Return(jen.ID("args").Dot("String").Call(jen.Zero()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockSecretGenerator")).ID("GenerateSalt").Params().Params(jen.Index().Byte(), jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(),
			jen.Line(),
			jen.Return(jen.ID("args").Dot("Get").Call(jen.Zero()).Assert(jen.Index().Byte()), jen.ID("args").Dot("Error").Call(jen.One())),
		),
		jen.Line(),
	)

	return code
}
