package authorization

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func rbacDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(jen.Null(),

		jen.Line(),
	)
	code.Add(jen.Null().Type().ID("role").ID("int"),

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("globalAuthorizer").Op("*").ID("gorbac").Dot(
		"RBAC",
	),

		jen.Line(),
	)
	code.Add(jen.Func().ID("init").Params().Body(jen.ID("globalAuthorizer").Op("=").ID("initializeRBAC").Call()),

		jen.Line(),
	)
	code.Add(jen.Func().ID("initializeRBAC").Params().Params(jen.Op("*").ID("gorbac").Dot(
		"RBAC",
	)).Body(
		jen.ID("rbac").Op(":=").ID("gorbac").Dot(
			"New",
		).Call(),
		jen.ID("must").Call(jen.ID("rbac").Dot(
			"Add",
		).Call(jen.ID("serviceUser"))),
		jen.ID("must").Call(jen.ID("rbac").Dot(
			"Add",
		).Call(jen.ID("serviceAdmin"))),
		jen.ID("must").Call(jen.ID("rbac").Dot(
			"Add",
		).Call(jen.ID("accountAdmin"))),
		jen.ID("must").Call(jen.ID("rbac").Dot(
			"Add",
		).Call(jen.ID("accountMember"))),
		jen.ID("must").Call(jen.ID("rbac").Dot(
			"SetParent",
		).Call(
			jen.ID("accountAdminRoleName"),
			jen.ID("accountMemberRoleName"),
		)),
		jen.ID("must").Call(jen.ID("rbac").Dot(
			"SetParent",
		).Call(
			jen.ID("serviceAdminRoleName"),
			jen.ID("accountAdminRoleName"),
		)),
		jen.Return().ID("rbac"),
	),

		jen.Line(),
	)
	return code
}
