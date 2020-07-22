package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func randDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.Add(
		jen.Const().Defs(
			jen.ID("randomReadSize").Equals().Lit(64),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("this function tests that we have appropriate access to crypto/rand"),
		jen.Line(),
		jen.Func().ID("init").Params().Block(
			jen.ID("b").Assign().Make(jen.Index().Byte(), jen.ID("randomReadSize")),
			jen.If(
				jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")),
				jen.Err().DoesNotEqual().Nil(),
			).Block(
				jen.Panic(jen.Err()),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("standardSecretGenerator").Struct(),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("g").PointerTo().ID("standardSecretGenerator")).ID("GenerateTwoFactorSecret").Params().Params(jen.String(), jen.Error()).Block(
			jen.ID("b").Assign().Make(jen.Index().Byte(), jen.ID("randomReadSize")),
			jen.Line(),
			jen.Comment("Note that err == nil only if we read len(b) bytes."),
			jen.Line(),
			jen.If(
				jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")),
				jen.Err().DoesNotEqual().Nil(),
			).Block(
				jen.Return(jen.EmptyString(), jen.Err()),
			),
			jen.Line(),
			jen.Return(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.Nil()),
		),
	)

	return code
}
