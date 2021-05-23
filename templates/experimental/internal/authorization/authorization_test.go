package authorization

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authorizationTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(jen.Null(),

		jen.Line(),
	)
	code.Add(jen.Func().ID("TestAuthorizations").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(
			jen.Lit("service user"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot(
					"Parallel",
				).Call(),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("CanCreateItems").Call(jen.ID("serviceUserRoleName")),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("CanSeeItems").Call(jen.ID("serviceUserRoleName")),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("CanSearchItems").Call(jen.ID("serviceUserRoleName")),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("CanUpdateItems").Call(jen.ID("serviceUserRoleName")),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("CanDeleteItems").Call(jen.ID("serviceUserRoleName")),
				),
			),
		),
		jen.ID("T").Dot(
			"Run",
		).Call(
			jen.Lit("service admin"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot(
					"Parallel",
				).Call(),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanCreateItems").Call(jen.ID("serviceAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanSeeItems").Call(jen.ID("serviceAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanSearchItems").Call(jen.ID("serviceAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanUpdateItems").Call(jen.ID("serviceAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanDeleteItems").Call(jen.ID("serviceAdminRoleName")),
				),
			),
		),
		jen.ID("T").Dot(
			"Run",
		).Call(
			jen.Lit("account admin"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot(
					"Parallel",
				).Call(),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanCreateItems").Call(jen.ID("accountAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanSeeItems").Call(jen.ID("accountAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanSearchItems").Call(jen.ID("accountAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanUpdateItems").Call(jen.ID("accountAdminRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanDeleteItems").Call(jen.ID("accountAdminRoleName")),
				),
			),
		),
		jen.ID("T").Dot(
			"Run",
		).Call(
			jen.Lit("account member"),
			jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
				jen.ID("t").Dot(
					"Parallel",
				).Call(),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanCreateItems").Call(jen.ID("accountMemberRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanSeeItems").Call(jen.ID("accountMemberRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanSearchItems").Call(jen.ID("accountMemberRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanUpdateItems").Call(jen.ID("accountMemberRoleName")),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("CanDeleteItems").Call(jen.ID("accountMemberRoleName")),
				),
			),
		),
	),

		jen.Line(),
	)
	return code
}
