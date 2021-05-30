package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func actionsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ErrUnavailableYet").Op("=").Qual("errors", "New").Call(jen.Lit("can't do this yet")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("actionFunc").Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")),
			jen.ID("Action").Struct(
				jen.ID("Action").ID("actionFunc"),
				jen.ID("Name").ID("string"),
				jen.ID("Weight").ID("int"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RandomAction takes a client and returns a closure which is an action."),
		jen.Line(),
		jen.Func().ID("RandomAction").Params(jen.ID("c").Op("*").ID("httpclient").Dot("Client"), jen.ID("builder").Op("*").ID("requests").Dot("Builder")).Params(jen.Op("*").ID("Action")).Body(
			jen.ID("allActions").Op(":=").Map(jen.ID("string")).Op("*").ID("Action").Valuesln(jen.Lit("GetHealthCheck").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("GetHealthCheck"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.Return().ID("builder").Dot("BuildHealthCheckRequest").Call(jen.ID("ctx")),
			), jen.ID("Weight").Op(":").Lit(100)), jen.Lit("CreateUser").Op(":").Valuesln(jen.ID("Name").Op(":").Lit("CreateUser"), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
				jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
				jen.ID("ui").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
				jen.Return().ID("builder").Dot("BuildCreateUserRequest").Call(
					jen.ID("ctx"),
					jen.ID("ui"),
				),
			), jen.ID("Weight").Op(":").Lit(100))),
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildItemActions").Call(
				jen.ID("c"),
				jen.ID("builder"),
			)).Body(
				jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v")),
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Op(":=").Range().ID("buildWebhookActions").Call(
				jen.ID("c"),
				jen.ID("builder"),
			)).Body(
				jen.ID("allActions").Index(jen.ID("k")).Op("=").ID("v")),
			jen.ID("totalWeight").Op(":=").Lit(0),
			jen.For(jen.List(jen.ID("_"), jen.ID("rb")).Op(":=").Range().ID("allActions")).Body(
				jen.ID("totalWeight").Op("+=").ID("rb").Dot("Weight")),
			jen.Qual("math/rand", "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
			jen.ID("r").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("totalWeight")),
			jen.For(jen.List(jen.ID("_"), jen.ID("rb")).Op(":=").Range().ID("allActions")).Body(
				jen.ID("r").Op("-=").ID("rb").Dot("Weight"),
				jen.If(jen.ID("r").Op("<=").Lit(0)).Body(
					jen.Return().ID("rb")),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
