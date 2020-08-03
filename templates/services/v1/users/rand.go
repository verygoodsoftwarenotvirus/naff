package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func randDotGo() *jen.File {
	code := jen.NewFile(packageName)

	code.Add(buildRandConstantDefs()...)
	code.Add(buildRandInit()...)
	code.Add(buildRandStandardSecretGeneratorGenerateTwoFactorSecret()...)
	code.Add(buildRandStandardSecretGeneratorGenerateSalt()...)

	return code
}

func buildRandConstantDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("saltSize").Equals().Lit(16),
			jen.ID("randomSecretSize").Equals().Lit(64),
		),
		jen.Line(),
	}

	return lines
}

func buildRandInit() []jen.Code {
	lines := []jen.Code{
		jen.Comment("this function tests that we have appropriate access to crypto/rand"),
		jen.Line(),
		jen.Func().ID("init").Params().Body(
			jen.ID("b").Assign().Make(jen.Index().Byte(), jen.ID("randomSecretSize")),
			jen.If(
				jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")),
				jen.Err().DoesNotEqual().Nil(),
			).Body(
				jen.Panic(jen.Err()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildRandStandardSecretGeneratorGenerateTwoFactorSecret() []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().ID("secretGenerator").Equals().Parens(jen.PointerTo().ID("standardSecretGenerator")).Call(jen.Nil()),
		jen.Line(),
		jen.Type().ID("standardSecretGenerator").Struct(),
		jen.Line(),
		jen.Func().Params(jen.ID("g").PointerTo().ID("standardSecretGenerator")).ID("GenerateTwoFactorSecret").Params().Params(jen.String(), jen.Error()).Body(
			jen.ID("b").Assign().Make(jen.Index().Byte(), jen.ID("randomSecretSize")),
			jen.Line(),
			jen.Comment("Note that err == nil only if we read len(b) bytes."),
			jen.If(
				jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")),
				jen.Err().DoesNotEqual().Nil(),
			).Body(
				jen.Return(jen.EmptyString(), jen.Err()),
			),
			jen.Line(),
			jen.Return(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.Nil()),
		),
	}

	return lines
}

func buildRandStandardSecretGeneratorGenerateSalt() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("g").PointerTo().ID("standardSecretGenerator")).ID("GenerateSalt").Params().Params(
			jen.Index().Byte(),
			jen.Error(),
		).Body(
			jen.ID("b").Assign().Make(jen.Index().Byte(), jen.ID("saltSize")),
			jen.Line(),
			jen.Comment("Note that err == nil only if we read len(b) bytes."),
			jen.If(
				jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")),
				jen.Err().DoesNotEqual().Nil(),
			).Body(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Return(jen.ID("b"), jen.Nil()),
		),
	}

	return lines
}
