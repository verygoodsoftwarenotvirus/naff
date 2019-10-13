package load

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func actionsDotGo() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("ErrUnavailableYet").Op("=").Qual("errors", "New").Call(jen.Lit("can't do this yet")),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("actionFunc").Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Type().ID("Action").Struct(jen.ID("Action").ID("actionFunc"), jen.ID("Weight").ID("int"), jen.ID("Name").ID("string")),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// RandomAction takes a client and returns a closure which is an action").ID("RandomAction").Params(jen.ID("c").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "V1Client")).Params(jen.Op("*").ID("Action")).Block(
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.ID("allActions").Op(":=").Map(jen.ID("string")).Op("*").ID("Action").Valuesln(jen.Lit("GetHealthCheck").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetHealthCheck"), jen.ID("Action").Op(":").ID("c").Dot(
			"BuildHealthCheckRequest",
		), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("CreateUser").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("CreateUser"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
			jen.ID("ui").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "RandomUserInput").Call(),
			jen.Return().ID("c").Dot(
				"BuildCreateUserRequest",
			).Call(jen.ID("ctx"), jen.ID("ui")),
		), jen.ID("Weight").Op(":").Lit(100))),
		jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildItemActions").Call(jen.ID("c"))).Block(
			jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v"),
		),
		jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildWebhookActions").Call(jen.ID("c"))).Block(
			jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v"),
		),
		jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildOAuth2ClientActions").Call(jen.ID("c"))).Block(
			jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v"),
		),
		jen.ID("totalWeight").Op(":=").Lit(0),
		jen.For(jen.List(jen.ID("_"), jen.ID("rb")).Op(":=").Range().ID("allActions")).Block(
			jen.ID("totalWeight").Op("+=").ID("rb").Dot(
				"Weight",
			),
		),
		jen.Qual("math/rand", "Seed").Call(jen.Qual("time", "Now").Call().Dot(
			"UnixNano",
		).Call()),
		jen.ID("r").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("totalWeight")),
		jen.For(jen.List(jen.ID("_"), jen.ID("rb")).Op(":=").Range().ID("allActions")).Block(
			jen.ID("r").Op("-=").ID("rb").Dot(
				"Weight",
			),
			jen.If(jen.ID("r").Op("<=").Lit(0)).Block(
				jen.Return().ID("rb"),
			),
		),
		jen.Return().ID("nil"),
	),

		jen.Line(),
	)
	return ret
}
