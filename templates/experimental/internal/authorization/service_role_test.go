package authorization

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceRoleTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(jen.Null(),

		jen.Line(),
	)
	code.Add(jen.Func().ID("TestServiceRoles").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
				jen.ID("r").Op(":=").ID("NewServiceRolePermissionChecker").Call(jen.ID("ServiceUserRole").Dot(
					"String",
				).Call()),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"IsServiceAdmin",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanCycleCookieSecrets",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeAccountAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeAPIClientAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeUserAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeWebhookAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanUpdateUserReputations",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeUserData",
					).Call(),
				),
				jen.ID("assert").Dot(
					"False",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSearchUsers",
					).Call(),
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
				jen.ID("r").Op(":=").ID("NewServiceRolePermissionChecker").Call(jen.ID("ServiceAdminRole").Dot(
					"String",
				).Call()),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"IsServiceAdmin",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanCycleCookieSecrets",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeAccountAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeAPIClientAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeUserAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeWebhookAuditLogEntries",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanUpdateUserReputations",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSeeUserData",
					).Call(),
				),
				jen.ID("assert").Dot(
					"True",
				).Call(
					jen.ID("t"),
					jen.ID("r").Dot(
						"CanSearchUsers",
					).Call(),
				),
			),
		),
	),

		jen.Line(),
	)
	return code
}
