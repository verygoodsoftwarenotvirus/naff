package load

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func actionsDotGo(rootPkg string) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().Defs(
			jen.Comment("ErrUnavailableYet is a sentinel error value"),
			jen.ID("ErrUnavailableYet").Op("=").Qual("errors", "New").Call(jen.Lit("can't do this yet")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("actionFunc represents a thing you can do"),
			jen.ID("actionFunc").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")),
			jen.Line(),
			jen.Comment("Action is a wrapper struct around some important values"),
			jen.ID("Action").Struct(
				jen.ID("Action").ID("actionFunc"),
				jen.ID("Weight").ID("int"),
				jen.ID("Name").ID("string"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("RandomAction takes a client and returns a closure which is an action"),
		jen.Line(),
		jen.Func().ID("RandomAction").Params(jen.ID("c").Op("*").Qual(filepath.Join(rootPkg, "client/v1/http"), "V1Client")).Params(jen.Op("*").ID("Action")).Block(
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("allActions").Op(":=").Map(jen.ID("string")).Op("*").ID("Action").Valuesln(
				jen.Lit("GetHealthCheck").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("GetHealthCheck"),
					jen.ID("Action").Op(":").ID("c").Dot("BuildHealthCheckRequest"),
					jen.ID("Weight").Op(":").Lit(100),
				),
				jen.Lit("CreateUser").Op(":").Valuesln(
					jen.ID("Name").Op(":").Lit("CreateUser"),
					jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.ID("ui").Op(":=").Qual(filepath.Join(rootPkg, "tests/v1/testutil/rand/model"), "RandomUserInput").Call(),
						jen.Return().ID("c").Dot("BuildCreateUserRequest").Call(jen.ID("ctx"), jen.ID("ui")),
					),
					jen.ID("Weight").Op(":").Lit(100),
				),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildItemActions").Call(jen.ID("c"))).Block(
				jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v"),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildWebhookActions").Call(jen.ID("c"))).Block(
				jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v"),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildOAuth2ClientActions").Call(jen.ID("c"))).Block(
				jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v"),
			),
			jen.Line(),
			jen.ID("totalWeight").Op(":=").Lit(0),
			jen.For(jen.List(jen.ID("_"), jen.ID("rb")).Op(":=").Range().ID("allActions")).Block(
				jen.ID("totalWeight").Op("+=").ID("rb").Dot("Weight"),
			),
			jen.Line(),
			jen.Qual("math/rand", "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
			jen.ID("r").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("totalWeight")),
			jen.Line(),
			jen.For(jen.List(jen.ID("_"), jen.ID("rb")).Op(":=").Range().ID("allActions")).Block(
				jen.ID("r").Op("-=").ID("rb").Dot("Weight"),
				jen.If(jen.ID("r").Op("<=").Lit(0)).Block(
					jen.Return().ID("rb"),
				),
			),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)
	return ret
}
