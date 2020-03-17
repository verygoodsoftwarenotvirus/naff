package v1

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("database")

	modelsImp := fmt.Sprintf("%s/models/v1", pkg.OutputPath)

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").ID("Scanner").Op("=").Parens(jen.Op("*").Qual("database/sql", "Row")).Call(jen.ID("nil")),
			jen.ID("_").ID("Querier").Op("=").Parens(jen.Op("*").Qual("database/sql", "DB")).Call(jen.ID("nil")),
			jen.ID("_").ID("Querier").Op("=").Parens(jen.Op("*").Qual("database/sql", "Tx")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	buildInterfaceLines := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("Migrate").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("error")),
			jen.ID("IsReady").Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("ready").ID("bool")),
			jen.Line(),
		}

		for _, typ := range pkg.DataTypes {
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
				jen.ID("Scan").Params(jen.ID("dest").Op("...").Interface()).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("Querier is a subset interface for sql.{DB|Tx} objects"),
			jen.ID("Querier").Interface(
				jen.ID("ExecContext").Params(utils.CtxParam(), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")),
				jen.ID("QueryContext").Params(utils.CtxParam(), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")),
				jen.ID("QueryRowContext").Params(utils.CtxParam(), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row")),
			),
			jen.Line(),
			jen.Comment("ConnectionDetails is a string alias for dependency injection"),
			jen.ID("ConnectionDetails").ID("string"),
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
