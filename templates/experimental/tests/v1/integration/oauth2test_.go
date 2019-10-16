package integration

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2TestDotGo() *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("mustBuildCode").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("totpSecret").ID("string")).Params(jen.ID("string")).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.List(jen.ID("code"), jen.ID("err")).Op(":=").ID("totp").Dot(
			"GenerateCode",
		).Call(jen.ID("totpSecret"), jen.Qual("time", "Now").Call().Dot(
			"UTC",
		).Call()),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.Return().ID("code"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyOAuth2ClientInput").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpToken")).ID("string")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2ClientCreationInput",
	)).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("x").Op(":=").Op("&").ID("models").Dot(
			"OAuth2ClientCreationInput",
		).Valuesln(jen.ID("UserLoginInput").Op(":").ID("models").Dot(
			"UserLoginInput",
		).Valuesln(jen.ID("Username").Op(":").ID("username"), jen.ID("Password").Op(":").ID("password"), jen.ID("TOTPToken").Op(":").ID("mustBuildCode").Call(jen.ID("t"), jen.ID("totpToken"))), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(jen.Lit("*")), jen.ID("RedirectURI").Op(":").Lit("http://localhost")),
		jen.Return().ID("x"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("convertInputToClient").Params(jen.ID("input").Op("*").ID("models").Dot(
		"OAuth2ClientCreationInput",
	)).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Block(
		jen.Return().Op("&").ID("models").Dot(
			"OAuth2Client",
		).Valuesln(jen.ID("ClientID").Op(":").ID("input").Dot(
			"ClientID",
		), jen.ID("ClientSecret").Op(":").ID("input").Dot(
			"ClientSecret",
		), jen.ID("RedirectURI").Op(":").ID("input").Dot(
			"RedirectURI",
		), jen.ID("Scopes").Op(":").ID("input").Dot(
			"Scopes",
		), jen.ID("BelongsTo").Op(":").ID("input").Dot(
			"BelongsTo",
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("checkOAuth2ClientEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ID",
		)),
		jen.ID("assert").Dot(
			"NotEmpty",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ClientID",
		)),
		jen.ID("assert").Dot(
			"NotEmpty",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ClientSecret",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"RedirectURI",
		), jen.ID("actual").Dot(
			"RedirectURI",
		)),
		jen.ID("assert").Dot(
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"Scopes",
		), jen.ID("actual").Dot(
			"Scopes",
		)),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"CreatedOn",
		)),
		jen.ID("assert").Dot(
			"Nil",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ArchivedOn",
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestOAuth2Clients").Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
		jen.ID("test").Dot(
			"Parallel",
		).Call(),
		jen.List(jen.ID("x"), jen.ID("y"), jen.ID("cookie")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
		jen.ID("assert").Dot(
			"NotNil",
		).Call(jen.ID("test"), jen.ID("cookie")),
		jen.ID("input").Op(":=").ID("buildDummyOAuth2ClientInput").Call(jen.ID("test"), jen.ID("x").Dot(
			"Username",
		), jen.ID("y").Dot(
			"Password",
		), jen.ID("x").Dot(
			"TwoFactorSecret",
		)),
		jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
			"CreateOAuth2Client",
		).Call(jen.Qual("context", "Background").Call(), jen.ID("cookie"), jen.ID("input")),
		jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.ID("err")),
		jen.List(jen.ID("testClient"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "NewClient").Call(jen.Qual("context", "Background").Call(), jen.ID("premade").Dot(
			"ClientID",
		), jen.ID("premade").Dot(
			"ClientSecret",
		), jen.ID("todoClient").Dot(
			"URL",
		), jen.ID("noop").Dot(
			"ProvideNoopLogger",
		).Call(), jen.ID("todoClient").Dot(
			"PlainClient",
		).Call(), jen.ID("premade").Dot(
			"Scopes",
		), jen.ID("debug")),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("test"), jen.ID("err"), jen.Lit("error setting up auxiliary client")),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be creatable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClient").Dot(
					"CreateOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("cookie"), jen.ID("input")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkOAuth2ClientEquality").Call(jen.ID("t"), jen.ID("convertInputToClient").Call(jen.ID("input")), jen.ID("actual")),
				jen.ID("err").Op("=").ID("testClient").Dot(
					"ArchiveOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("actual").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should return an error when trying to read one that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("testClient").Dot(
					"GetOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("nonexistentID")),
				jen.ID("assert").Dot(
					"Error",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("input").Op(":=").ID("buildDummyOAuth2ClientInput").Call(jen.ID("t"), jen.ID("x").Dot(
					"Username",
				), jen.ID("y").Dot(
					"Password",
				), jen.ID("x").Dot(
					"TwoFactorSecret",
				)),
				jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("testClient").Dot(
					"CreateOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("cookie"), jen.ID("input")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("c"), jen.ID("err")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClient").Dot(
					"GetOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("c").Dot(
					"ID",
				)),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkOAuth2ClientEquality").Call(jen.ID("t"), jen.ID("convertInputToClient").Call(jen.ID("input")), jen.ID("actual")),
				jen.ID("err").Op("=").ID("testClient").Dot(
					"ArchiveOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("actual").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("input").Op(":=").ID("buildDummyOAuth2ClientInput").Call(jen.ID("t"), jen.ID("x").Dot(
					"Username",
				), jen.ID("y").Dot(
					"Password",
				), jen.ID("x").Dot(
					"TwoFactorSecret",
				)),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("testClient").Dot(
					"CreateOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("cookie"), jen.ID("input")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.ID("err").Op("=").ID("testClient").Dot(
					"ArchiveOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be unable to authorize after being deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("createdUser"), jen.ID("createdUserInput"), jen.ID("_")).Op(":=").ID("buildDummyUser").Call(jen.ID("test")),
				jen.ID("assert").Dot(
					"NotNil",
				).Call(jen.ID("test"), jen.ID("cookie")),
				jen.ID("input").Op(":=").ID("buildDummyOAuth2ClientInput").Call(jen.ID("test"), jen.ID("createdUserInput").Dot(
					"Username",
				), jen.ID("createdUserInput").Dot(
					"Password",
				), jen.ID("createdUser").Dot(
					"TwoFactorSecret",
				)),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateOAuth2Client",
				).Call(jen.Qual("context", "Background").Call(), jen.ID("cookie"), jen.ID("input")),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("premade"), jen.ID("err")),
				jen.ID("err").Op("=").ID("testClient").Dot(
					"ArchiveOAuth2Client",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
				jen.List(jen.ID("c2"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "NewClient").Call(jen.Qual("context", "Background").Call(), jen.ID("premade").Dot(
					"ClientID",
				), jen.ID("premade").Dot(
					"ClientSecret",
				), jen.ID("todoClient").Dot(
					"URL",
				), jen.ID("noop").Dot(
					"ProvideNoopLogger",
				).Call(), jen.ID("buildHTTPClient").Call(), jen.ID("premade").Dot(
					"Scopes",
				), jen.ID("true")),
				jen.ID("checkValueAndError").Call(jen.ID("test"), jen.ID("c2"), jen.ID("err")),
				jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("c2").Dot(
					"GetOAuth2Clients",
				).Call(jen.ID("tctx"), jen.ID("nil")),
				jen.ID("assert").Dot(
					"Error",
				).Call(jen.ID("t"), jen.ID("err"), jen.Lit("expected error from what should be an unauthorized client")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),

		jen.Var().ID("expected").Index().Op("*").ID("models").Dot(
					"OAuth2Client",
				),
				jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
					jen.ID("input").Op(":=").ID("buildDummyOAuth2ClientInput").Call(jen.ID("t"), jen.ID("x").Dot(
						"Username",
					), jen.ID("y").Dot(
						"Password",
					), jen.ID("x").Dot(
						"TwoFactorSecret",
					)),
					jen.List(jen.ID("oac"), jen.ID("err")).Op(":=").ID("testClient").Dot(
						"CreateOAuth2Client",
					).Call(jen.ID("tctx"), jen.ID("cookie"), jen.ID("input")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("oac"), jen.ID("err")),
					jen.ID("expected").Op("=").ID("append").Call(jen.ID("expected"), jen.ID("oac")),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("testClient").Dot(
					"GetOAuth2Clients",
				).Call(jen.ID("tctx"), jen.ID("nil")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("assert").Dot(
					"True",
				).Call(jen.ID("t"), jen.ID("len").Call(jen.ID("actual").Dot(
					"Clients",
				)).Op("-").ID("len").Call(jen.ID("expected")).Op(">").Lit(0), jen.Lit("expected %d - %d to be > 0"), jen.ID("len").Call(jen.ID("actual").Dot(
					"Clients",
				)), jen.ID("len").Call(jen.ID("expected"))),
				jen.For(jen.List(jen.ID("_"), jen.ID("oAuth2Client")).Op(":=").Range().ID("expected")).Block(
					jen.ID("clientFound").Op(":=").ID("false"),
					jen.For(jen.List(jen.ID("_"), jen.ID("c")).Op(":=").Range().ID("actual").Dot(
						"Clients",
					)).Block(
						jen.If(jen.ID("c").Dot(
							"ID",
						).Op("==").ID("oAuth2Client").Dot(
							"ID",
						)).Block(
							jen.ID("clientFound").Op("=").ID("true"),
							jen.Break(),
						),
					),
					jen.ID("assert").Dot(
						"True",
					).Call(jen.ID("t"), jen.ID("clientFound"), jen.Lit("expected oAuth2Client ID %d to be present in results"), jen.ID("oAuth2Client").Dot(
						"ID",
					)),
				),
				jen.For(jen.List(jen.ID("_"), jen.ID("oa2c")).Op(":=").Range().ID("expected")).Block(
					jen.ID("err").Op("=").ID("testClient").Dot(
						"ArchiveOAuth2Client",
					).Call(jen.ID("tctx"), jen.ID("oa2c").Dot(
						"ID",
					)),
					jen.ID("assert").Dot(
						"NoError",
					).Call(jen.ID("t"), jen.ID("err"), jen.Lit("error deleting client %d: %v"), jen.ID("oa2c").Dot(
						"ID",
					), jen.ID("err")),
				),
			)),
		)),
	),
	jen.Line(),
	)
	return ret
}
