package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)

	ret.Add(jen.Null().Var().ID("usersTableName").Op("=").Lit("users"))
	ret.Add(jen.Null().Var().ID("usersTableColumns").Op("=").Index().ID("string").Valuesln(jen.Lit("id"), jen.Lit("username"), jen.Lit("hashed_password"), jen.Lit("password_last_changed_on"), jen.Lit("two_factor_secret"), jen.Lit("is_admin"), jen.Lit("created_on"), jen.Lit("updated_on"), jen.Lit("archived_on")))
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func())
	ret.Add(jen.Func().Params(jen.ID("p").Op("*").ID("Postgres")).ID("buildArchiveUserQuery").Params(jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
		jen.Null().Var().ID("err").ID("error"),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("p").Dot(
			"sqlBuilder",
		).Dot(
			"Update",
		).Call(jen.ID("usersTableName")).Dot(
			"Set",
		).Call(jen.Lit("updated_on"), jen.ID("squirrel").Dot(
			"Expr",
		).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
			"Set",
		).Call(jen.Lit("archived_on"), jen.ID("squirrel").Dot(
			"Expr",
		).Call(jen.ID("CurrentUnixTimeQuery"))).Dot(
			"Where",
		).Call(jen.ID("squirrel").Dot(
			"Eq",
		).Valuesln(jen.Lit("id").Op(":").ID("userID"))).Dot(
			"Suffix",
		).Call(jen.Lit("RETURNING archived_on")).Dot(
			"ToSql",
		).Call(),
		jen.ID("p").Dot(
			"logQueryBuildingError",
		).Call(jen.ID("err")),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	),
	)
	ret.Add(jen.Func())
	return ret
}
