package v1

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("database")

	modelsImp := fmt.Sprintf("%s/models/v1", proj.OutputPath)

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Defs(
			jen.Underscore().ID("Scanner").Equals().Parens(jen.PointerTo().Qual("database/sql", "Row")).Call(jen.Nil()),
			jen.Underscore().ID("Querier").Equals().Parens(jen.PointerTo().Qual("database/sql", "DB")).Call(jen.Nil()),
			jen.Underscore().ID("Querier").Equals().Parens(jen.PointerTo().Qual("database/sql", "Tx")).Call(jen.Nil()),
		),
		jen.Line(),
	)

	buildInterfaceLines := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("Migrate").Params(constants.CtxParam()).Params(jen.Error()),
			jen.ID("IsReady").Params(constants.CtxParam()).Params(jen.ID("ready").Bool()),
			jen.Line(),
		}

		for _, typ := range proj.DataTypes {
			lines = append(lines,
				jen.Qual(modelsImp, fmt.Sprintf("%sDataManager", typ.Name.Singular())),
			)
		}

		lines = append(lines,
			jen.Qual(modelsImp, "UserDataManager"),
			jen.Qual(modelsImp, "OAuth2ClientDataManager"),
			jen.Qual(modelsImp, "WebhookDataManager"),
		)

		return lines
	}

	ret.Add(
		jen.Type().Defs(
			jen.Comment("Scanner represents any database response (i.e. sql.Row[s])"),
			jen.ID("Scanner").Interface(
				jen.ID("Scan").Params(jen.ID("dest").Spread().Interface()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("Querier is a subset interface for sql.{DB|Tx} objects"),
			jen.ID("Querier").Interface(
				jen.ID("ExecContext").Params(constants.CtxParam(), jen.ID("query").String(), jen.ID("args").Spread().Interface()).Params(jen.Qual("database/sql", "Result"), jen.Error()),
				jen.ID("QueryContext").Params(constants.CtxParam(), jen.ID("query").String(), jen.ID("args").Spread().Interface()).Params(jen.PointerTo().Qual("database/sql", "Rows"), jen.Error()),
				jen.ID("QueryRowContext").Params(constants.CtxParam(), jen.ID("query").String(), jen.ID("args").Spread().Interface()).Params(jen.PointerTo().Qual("database/sql", "Row")),
			),
			jen.Line(),
			jen.Comment("ConnectionDetails is a string alias for dependency injection"),
			jen.ID("ConnectionDetails").String(),
			jen.Line(),
			jen.Comment("Database describes anything that stores data for our services"),
			jen.ID("Database").Interface(
				buildInterfaceLines()...,
			),
		),
		jen.Line(),
	)

	return ret
}
