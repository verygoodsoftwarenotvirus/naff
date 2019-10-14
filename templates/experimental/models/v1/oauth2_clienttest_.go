package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientTestDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestOAuth2Client_GetID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").Lit("uint64(123)"),
			jen.ID("oac").Op(":=").Op("&").ID("OAuth2Client").Valuesln(jen.ID("ClientID").Op(":").ID("expected")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot(
				"GetID",
			).Call()),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestOAuth2Client_GetSecret").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").Lit("uint64(123)"),
			jen.ID("oac").Op(":=").Op("&").ID("OAuth2Client").Valuesln(jen.ID("ClientSecret").Op(":").ID("expected")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot(
				"GetSecret",
			).Call()),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestOAuth2Client_GetDomain").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").Lit("uint64(123)"),
			jen.ID("oac").Op(":=").Op("&").ID("OAuth2Client").Valuesln(jen.ID("RedirectURI").Op(":").ID("expected")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot(
				"GetDomain",
			).Call()),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestOAuth2Client_GetUserID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expectation").Op(":=").ID("uint64").Call(jen.Lit(123)),
			jen.ID("expected").Op(":=").Lit("123"),
			jen.ID("oac").Op(":=").Op("&").ID("OAuth2Client").Valuesln(jen.ID("BelongsTo").Op(":").ID("expectation")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expected"), jen.ID("oac").Dot(
				"GetUserID",
			).Call()),
		)),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().ID("TestOAuth2Client_HasScope").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("oac").Op(":=").Op("&").ID("OAuth2Client").Valuesln(jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(jen.Lit("things"), jen.Lit("and"), jen.Lit("stuff"))),
			jen.ID("assert").Dot(
				"True",
			).Call(jen.ID("t"), jen.ID("oac").Dot(
				"HasScope",
			).Call(jen.ID("oac").Dot(
				"Scopes",
			).Index(jen.Lit(0)))),
			jen.ID("assert").Dot(
				"False",
			).Call(jen.ID("t"), jen.ID("oac").Dot(
				"HasScope",
			).Call(jen.Lit("blah"))),
			jen.ID("assert").Dot(
				"False",
			).Call(jen.ID("t"), jen.ID("oac").Dot(
				"HasScope",
			).Call(jen.Lit(""))),
		)),
	),

		jen.Line(),
	)
	return ret
}
