package authorization

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authorizationDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(jen.Func().ID("must").Params(jen.ID("err").ID("error")).Body(jen.If(jen.ID("err").Op("!=").ID("nil")).Body(jen.ID("panic").Call(jen.ID("err")))),

		jen.Line(),
	)
	code.Add(jen.Func().ID("hasPermission").Params(jen.ID("p").ID("Permission"), jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(
		jen.For(jen.List(jen.ID("_"), jen.ID("r")).Op(":=").Range().ID("roles")).Body(jen.If(jen.Op("!").ID("globalAuthorizer").Dot(
			"IsGranted",
		).Call(
			jen.ID("r"),
			jen.ID("p"),
			jen.ID("nil"),
		)).Body(jen.Return().ID("false"))),
		jen.Return().ID("true"),
	),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanCreateItems returns whether a user can create items or not.").ID("CanCreateItems").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("CreateItemsPermission"),
		jen.ID("roles").Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeItems returns whether a user can view items or not.").ID("CanSeeItems").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadItemsPermission"),
		jen.ID("roles").Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSearchItems returns whether a user can search items or not.").ID("CanSearchItems").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("SearchItemsPermission"),
		jen.ID("roles").Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanUpdateItems returns whether a user can update items or not.").ID("CanUpdateItems").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("UpdateItemsPermission"),
		jen.ID("roles").Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanDeleteItems returns whether a user can delete items or not.").ID("CanDeleteItems").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ArchiveItemsPermission"),
		jen.ID("roles").Op("..."),
	)),

		jen.Line(),
	)
	return code
}
