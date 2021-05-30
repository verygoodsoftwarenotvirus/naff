package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminUserDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AdminUserDataManager").Op("=").Parens(jen.Op("*").ID("AdminUserDataManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AdminUserDataManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUserReputation is a mock function."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AdminUserDataManager")).ID("UpdateUserReputation").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("UserReputationUpdateInput")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("input"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
