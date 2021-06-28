package authorization

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authorizationDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("must").Params(jen.ID("err").ID("error")).Body(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")))),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("hasPermission").Params(jen.ID("p").ID("Permission"), jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(
			jen.For(jen.List(jen.ID("_"), jen.ID("r")).Assign().Range().ID("roles")).Body(
				jen.If(jen.Op("!").ID("globalAuthorizer").Dot("IsGranted").Call(
					jen.ID("r"),
					jen.ID("p"),
					jen.ID("nil"),
				)).Body(
					jen.Return().ID("false"))),
			jen.Newline(),
			jen.Return().ID("true"),
		),
		jen.Newline(),
	)

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		pcn := typ.Name.PluralCommonName()

		code.Add(
			jen.Commentf("CanCreate%s returns whether a user can create %s or not.", pn, pcn),
			jen.Newline(),
			jen.Func().IDf("CanCreate%s", pn).Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(
				jen.Return().ID("hasPermission").Call(
					jen.IDf("Create%sPermission", pn),
					jen.ID("roles").Op("..."),
				)),
			jen.Newline(),
		)

		code.Add(
			jen.Commentf("CanSee%s returns whether a user can view %s or not.", pn, pcn),
			jen.Newline(),
			jen.Func().IDf("CanSee%s", pn).Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(
				jen.Return().ID("hasPermission").Call(
					jen.IDf("Read%sPermission", pn),
					jen.ID("roles").Op("..."),
				)),
			jen.Newline(),
		)

		if typ.SearchEnabled {
			code.Add(
				jen.Commentf("CanSearch%s returns whether a user can search %s or not.", pn, pcn),
				jen.Newline(),
				jen.Func().IDf("CanSearch%s", pn).Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(
					jen.Return().ID("hasPermission").Call(
						jen.IDf("Search%sPermission", pn),
						jen.ID("roles").Op("..."),
					)),
				jen.Newline(),
			)
		}

		code.Add(
			jen.Commentf("CanUpdate%s returns whether a user can update %s or not.", pn, pcn),
			jen.Newline(),
			jen.Func().IDf("CanUpdate%s", pn).Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(
				jen.Return().ID("hasPermission").Call(
					jen.IDf("Update%sPermission", pn),
					jen.ID("roles").Op("..."),
				)),
			jen.Newline(),
		)

		code.Add(
			jen.Commentf("CanDelete%s returns whether a user can delete %s or not.", pn, pcn),
			jen.Newline(),
			jen.Func().IDf("CanDelete%s", pn).Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("bool")).Body(
				jen.Return().ID("hasPermission").Call(
					jen.IDf("Archive%sPermission", pn),
					jen.ID("roles").Op("..."),
				)),
			jen.Newline(),
		)
	}

	return code
}
