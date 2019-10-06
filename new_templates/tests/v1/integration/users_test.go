package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func usersTestDotGo() *jen.File {
	ret := jen.NewFile("integration")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("init").Params().Block(
		jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
		jen.If(
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")),
			jen.ID("err").Op("!=").ID("nil"),
		).Block(
			jen.ID("panic").Call(jen.ID("err")),
		),
	),
	)
	ret.Add(jen.Func(),
	)
	ret.Add(jen.Func().ID("buildDummyUserInput").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("models").Dot(
		"UserInput",
	)).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("fake").Dot(
			"Seed",
		).Call(jen.Qual("time", "Now").Call().Dot(
			"UnixNano",
		).Call()),
		jen.ID("userInput").Op(":=").Op("&").ID("models").Dot(
			"UserInput",
		).Valuesln(jen.ID("Username").Op(":").ID("fake").Dot(
			"UserName",
		).Call(), jen.ID("Password").Op(":").ID("fake").Dot(
			"Password",
		).Call(jen.Lit(8), jen.Lit(64), jen.ID("true"), jen.ID("true"), jen.ID("true"))),
		jen.Return().ID("userInput"),
	),
	)
	ret.Add(jen.Func().ID("buildDummyUser").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("models").Dot(
		"UserCreationResponse",
	), jen.Op("*").ID("models").Dot(
		"UserInput",
	), jen.Op("*").Qual("net/http", "Cookie")).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.ID("userInput").Op(":=").ID("buildDummyUserInput").Call(jen.ID("t")),
		jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
			"CreateUser",
		).Call(jen.ID("ctx"), jen.ID("userInput")),
		jen.ID("assert").Dot(
			"NotNil",
		).Call(jen.ID("t"), jen.ID("user")),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.If(
			jen.ID("user").Op("==").ID("nil").Op("||").ID("err").Op("!=").ID("nil"),
		).Block(
			jen.ID("t").Dot(
				"FailNow",
			).Call(),
		),
		jen.ID("cookie").Op(":=").ID("loginUser").Call(jen.ID("t"), jen.ID("userInput").Dot(
			"Username",
		), jen.ID("userInput").Dot(
			"Password",
		), jen.ID("user").Dot(
			"TwoFactorSecret",
		)),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.ID("require").Dot(
			"NotNil",
		).Call(jen.ID("t"), jen.ID("cookie")),
		jen.Return().List(jen.ID("user"), jen.ID("userInput"), jen.ID("cookie")),
	),
	)
	ret.Add(jen.Func().ID("checkUserCreationEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("expected").Op("*").ID("models").Dot(
		"UserInput",
	), jen.ID("actual").Op("*").ID("models").Dot(
		"UserCreationResponse",
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
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"Username",
		), jen.ID("actual").Dot(
			"Username",
		)),
		jen.ID("assert").Dot(
			"NotEmpty",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"TwoFactorSecret",
		)),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"CreatedOn",
		)),
		jen.ID("assert").Dot(
			"Nil",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"UpdatedOn",
		)),
		jen.ID("assert").Dot(
			"Nil",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ArchivedOn",
		)),
	),
	)
	ret.Add(jen.Func().ID("checkUserEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("expected").Op("*").ID("models").Dot(
		"UserInput",
	), jen.ID("actual").Op("*").ID("models").Dot(
		"User",
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
			"Equal",
		).Call(jen.ID("t"), jen.ID("expected").Dot(
			"Username",
		), jen.ID("actual").Dot(
			"Username",
		)),
		jen.ID("assert").Dot(
			"NotZero",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"CreatedOn",
		)),
		jen.ID("assert").Dot(
			"Nil",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"UpdatedOn",
		)),
		jen.ID("assert").Dot(
			"Nil",
		).Call(jen.ID("t"), jen.ID("actual").Dot(
			"ArchivedOn",
		)),
	),
	)
	ret.Add(jen.Func().ID("TestUsers").Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
		jen.ID("test").Dot(
			"Parallel",
		).Call(),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be creatable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("expected").Op(":=").ID("buildDummyUserInput").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateUser",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"UserInput",
				).Valuesln(jen.ID("Username").Op(":").ID("expected").Dot(
					"Username",
				), jen.ID("Password").Op(":").ID("expected").Dot(
					"Password",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkUserCreationEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("todoClient").Dot(
					"ArchiveUser",
				).Call(jen.ID("tctx"), jen.ID("actual").Dot(
					"ID",
				))),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should return an error when trying to read something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetUser",
				).Call(jen.ID("tctx"), jen.ID("nonexistentID")),
				jen.ID("assert").Dot(
					"Nil",
				).Call(jen.ID("t"), jen.ID("actual")),
				jen.ID("assert").Dot(
					"Error",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("expected").Op(":=").ID("buildDummyUserInput").Call(jen.ID("t")),
				jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateUser",
				).Call(jen.ID("tctx"), jen.Op("&").ID("models").Dot(
					"UserInput",
				).Valuesln(jen.ID("Username").Op(":").ID("expected").Dot(
					"Username",
				), jen.ID("Password").Op(":").ID("expected").Dot(
					"Password",
				))),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
				jen.ID("assert").Dot(
					"NotEmpty",
				).Call(jen.ID("t"), jen.ID("premade").Dot(
					"TwoFactorSecret",
				)),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetUser",
				).Call(jen.ID("tctx"), jen.ID("premade").Dot(
					"ID",
				)),
				jen.If(
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.ID("t").Dot(
						"Logf",
					).Call(jen.Lit("error encountered trying to fetch user %q: %v\n"), jen.ID("premade").Dot(
						"Username",
					), jen.ID("err")),
				),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("checkUserEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("todoClient").Dot(
					"ArchiveUser",
				).Call(jen.ID("tctx"), jen.ID("actual").Dot(
					"ID",
				))),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("y").Op(":=").ID("buildDummyUserInput").Call(jen.ID("t")),
				jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"CreateUser",
				).Call(jen.ID("tctx"), jen.ID("y")),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot(
					"NotNil",
				).Call(jen.ID("t"), jen.ID("u")),
				jen.If(
					jen.ID("u").Op("==").ID("nil").Op("||").ID("err").Op("!=").ID("nil"),
				).Block(
					jen.ID("t").Dot(
						"Log",
					).Call(jen.Lit("something has gone awry, user returned is nil")),
					jen.ID("t").Dot(
						"FailNow",
					).Call(),
				),
				jen.ID("err").Op("=").ID("todoClient").Dot(
					"ArchiveUser",
				).Call(jen.ID("tctx"), jen.ID("u").Dot(
					"ID",
				)),
				jen.ID("assert").Dot(
					"NoError",
				).Call(jen.ID("t"), jen.ID("err")),
			)),
		)),
		jen.ID("test").Dot(
			"Run",
		).Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot(
				"Run",
			).Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
				jen.Null().Var().ID("expected").Index().Op("*").ID("models").Dot(
					"UserCreationResponse",
				),
				jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
					jen.List(jen.ID("user"), jen.ID("_"), jen.ID("c")).Op(":=").ID("buildDummyUser").Call(jen.ID("t")),
					jen.ID("assert").Dot(
						"NotNil",
					).Call(jen.ID("t"), jen.ID("c")),
					jen.ID("expected").Op("=").ID("append").Call(jen.ID("expected"), jen.ID("user")),
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dot(
					"GetUsers",
				).Call(jen.ID("tctx"), jen.ID("nil")),
				jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
				jen.ID("assert").Dot(
					"True",
				).Call(jen.ID("t"), jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(
					"Users",
				))),
				jen.For(jen.List(jen.ID("_"), jen.ID("user")).Op(":=").Range().ID("actual").Dot(
					"Users",
				)).Block(
					jen.ID("err").Op("=").ID("todoClient").Dot(
						"ArchiveUser",
					).Call(jen.ID("tctx"), jen.ID("user").Dot(
						"ID",
					)),
					jen.ID("assert").Dot(
						"NoError",
					).Call(jen.ID("t"), jen.ID("err")),
				),
			)),
		)),
	),
	)
	return ret
}
