package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func userTestDotGo() *jen.File {
	ret := jen.NewFile("models")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("TestUser_Update").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("actual").Op(":=").ID("User").Valuesln(jen.ID("Username").Op(":").Lit("username"), jen.ID("HashedPassword").Op(":").Lit("hashed_pass"), jen.ID("TwoFactorSecret").Op(":").Lit("two factor secret")),
			jen.ID("exampleInput").Op(":=").ID("User").Valuesln(jen.ID("Username").Op(":").Lit("newUsername"), jen.ID("HashedPassword").Op(":").Lit("updated_hashed_pass"), jen.ID("TwoFactorSecret").Op(":").Lit("new fancy secret")),
			jen.ID("actual").Dot(
				"Update",
			).Call(jen.Op("&").ID("exampleInput")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleInput"), jen.ID("actual")),
		)),
	),
	)
	return ret
}
