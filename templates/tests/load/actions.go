package load

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func actionsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildVarDefs()...)
	code.Add(buildTypeDefs()...)
	code.Add(buildRandomAction(proj)...)

	return code
}

func buildVarDefs() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Comment("ErrUnavailableYet is a sentinel error value."),
			jen.ID("ErrUnavailableYet").Equals().Qual("errors", "New").Call(jen.Lit("can't do this yet")),
		),
		jen.Line(),
	}

	return lines
}

func buildTypeDefs() []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("actionFunc represents a thing you can do."),
			jen.ID("actionFunc").Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()),
			jen.Line(),
			jen.Comment("Action is a wrapper struct around some important values."),
			jen.ID("Action").Struct(
				jen.ID("Action").ID("actionFunc"),
				jen.ID("Weight").ID("int"),
				jen.ID("Name").String(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildRandomAction(proj *models.Project) []jen.Code {
	randomActionLines := []jen.Code{
		jen.ID("allActions").Assign().Map(jen.String()).PointerTo().ID("Action").Valuesln(
			jen.Lit("GetHealthCheck").MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Lit("GetHealthCheck"),
				jen.ID("Action").MapAssign().Func().Params().Params(
					jen.PointerTo().Qual("net/http", "Request"),
					jen.Error(),
				).Body(
					constants.CreateCtx(),
					jen.Return(jen.ID("c").Dot("BuildHealthCheckRequest").Call(constants.CtxVar())),
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
			jen.Lit("CreateUser").MapAssign().Valuesln(
				jen.ID("Name").MapAssign().Lit("CreateUser"),
				jen.ID("Action").MapAssign().Func().Params().Params(jen.PointerTo().Qual("net/http", "Request"), jen.Error()).Body(
					constants.CreateCtx(),
					jen.ID("ui").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUserCreationInput").Call(),
					jen.Return().ID("c").Dot("BuildCreateUserRequest").Call(constants.CtxVar(), jen.ID("ui")),
				),
				jen.ID("Weight").MapAssign().Lit(100),
			),
		),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		randomActionLines = append(randomActionLines,
			jen.For(jen.List(jen.ID("k"), jen.ID("v")).Assign().Range().IDf("build%sActions", typ.Name.Singular()).Call(jen.ID("c"))).Body(
				jen.ID("allActions").Index(jen.ID("k")).Equals().ID("v"),
			),
			jen.Line(),
		)
	}

	randomActionLines = append(randomActionLines,
		jen.For(jen.List(jen.ID("k"), jen.ID("v")).Assign().Range().ID("buildWebhookActions").Call(jen.ID("c"))).Body(
			jen.ID("allActions").Index(jen.ID("k")).Equals().ID("v"),
		),
		jen.Line(),
		jen.For(jen.List(jen.ID("k"), jen.ID("v")).Assign().Range().ID("buildOAuth2ClientActions").Call(jen.ID("c"))).Body(
			jen.ID("allActions").Index(jen.ID("k")).Equals().ID("v"),
		),
		jen.Line(),
		jen.ID("totalWeight").Assign().Zero(),
		jen.For(jen.List(jen.Underscore(), jen.ID("rb")).Assign().Range().ID("allActions")).Body(
			jen.ID("totalWeight").Op("+=").ID("rb").Dot("Weight"),
		),
		jen.Line(),
		jen.Qual("math/rand", "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call()),
		jen.ID("r").Assign().Qual("math/rand", "Intn").Call(jen.ID("totalWeight")),
		jen.Line(),
		jen.For(jen.List(jen.Underscore(), jen.ID("rb")).Assign().Range().ID("allActions")).Body(
			jen.ID("r").Op("-=").ID("rb").Dot("Weight"),
			jen.If(jen.ID("r").Op("<=").Zero()).Body(
				jen.Return().ID("rb"),
			),
		),
		jen.Line(),
		jen.Return().ID("nil"),
	)

	lines := []jen.Code{
		jen.Comment("RandomAction takes a client and returns a closure which is an action."),
		jen.Line(),
		jen.Func().ID("RandomAction").Params(jen.ID("c").PointerTo().Qual(proj.HTTPClientV1Package(), "V1Client")).Params(jen.PointerTo().ID("Action")).Body(
			randomActionLines...,
		),
		jen.Line(),
	}

	return lines
}
