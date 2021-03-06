package v1

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("database")

	utils.AddImports(proj, code)

	code.Add(buildVarDeclarations()...)
	code.Add(buildTypeDeclarations(proj)...)

	return code
}

func buildVarDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Underscore().ID("Scanner").Equals().Parens(jen.PointerTo().Qual("database/sql", "Row")).Call(jen.Nil()),
			jen.Underscore().ID("Querier").Equals().Parens(jen.PointerTo().Qual("database/sql", "DB")).Call(jen.Nil()),
			jen.Underscore().ID("Querier").Equals().Parens(jen.PointerTo().Qual("database/sql", "Tx")).Call(jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildTypeDeclarations(proj *models.Project) []jen.Code {
	modelsImp := fmt.Sprintf("%s/models/v1", proj.OutputPath)

	interfaceLines := []jen.Code{
		jen.ID("Migrate").Params(constants.CtxParam()).Params(jen.Error()),
		jen.ID("IsReady").Params(constants.CtxParam()).Params(jen.ID("ready").Bool()),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		interfaceLines = append(interfaceLines,
			jen.Qual(modelsImp, fmt.Sprintf("%sDataManager", typ.Name.Singular())),
		)
	}

	interfaceLines = append(interfaceLines,
		jen.Qual(modelsImp, "UserDataManager"),
		jen.Qual(modelsImp, "OAuth2ClientDataManager"),
		jen.Qual(modelsImp, "WebhookDataManager"),
	)

	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("Scanner represents any database response (i.e. sql.Row[s])"),
			jen.ID("Scanner").Interface(
				jen.ID("Scan").Params(jen.ID("dest").Spread().Interface()).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("ResultIterator represents any iterable database response (i.e. sql.Rows)"),
			jen.ID("ResultIterator").Interface(
				jen.ID("Next").Params().Bool(),
				jen.ID("Err").Params().Error(),
				jen.ID("Scanner"),
				jen.Qual("io", "Closer"),
			),
			jen.Line(),
			jen.Comment("Querier is a subset interface for sql.{DB|Tx} objects"),
			jen.ID("Querier").Interface(
				jen.ID("ExecContext").Params(constants.CtxParam(), jen.ID("query").String(), jen.ID("args").Spread().Interface()).Params(jen.Qual("database/sql", "Result"), jen.Error()),
				jen.ID("QueryContext").Params(constants.CtxParam(), jen.ID("query").String(), jen.ID("args").Spread().Interface()).Params(jen.PointerTo().Qual("database/sql", "Rows"), jen.Error()),
				jen.ID("QueryRowContext").Params(constants.CtxParam(), jen.ID("query").String(), jen.ID("args").Spread().Interface()).Params(jen.PointerTo().Qual("database/sql", "Row")),
			),
			jen.Line(),
			jen.Comment("ConnectionDetails is a string alias for dependency injection."),
			jen.ID("ConnectionDetails").String(),
			jen.Line(),
			jen.Comment("DataManager describes anything that stores data for our services."),
			jen.ID("DataManager").Interface(
				interfaceLines...,
			),
		),
		jen.Line(),
	}

	return lines
}
