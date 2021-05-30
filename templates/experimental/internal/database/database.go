package database

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func databaseDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("Scanner").Op("=").Parens(jen.Op("*").Qual("database/sql", "Row")).Call(jen.ID("nil")),
			jen.ID("_").ID("Querier").Op("=").Parens(jen.Op("*").Qual("database/sql", "DB")).Call(jen.ID("nil")),
			jen.ID("_").ID("Querier").Op("=").Parens(jen.Op("*").Qual("database/sql", "Tx")).Call(jen.ID("nil")),
			jen.ID("ErrDatabaseNotReady").Op("=").Qual("errors", "New").Call(jen.Lit("database is not ready yet")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Scanner").Interface(jen.ID("Scan").Params(jen.ID("dest").Op("...").Interface()).Params(jen.ID("error"))),
			jen.ID("ResultIterator").Interface(
				jen.ID("Next").Params().Params(jen.ID("bool")),
				jen.ID("Err").Params().Params(jen.ID("error")),
				jen.ID("Scanner"),
				jen.Qual("io", "Closer"),
			),
			jen.ID("Querier").Interface(
				jen.ID("ExecContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Qual("database/sql", "Result"), jen.ID("error")),
				jen.ID("QueryContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Rows"), jen.ID("error")),
				jen.ID("QueryRowContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("args").Op("...").Interface()).Params(jen.Op("*").Qual("database/sql", "Row")),
			),
			jen.ID("MetricsCollectionInterval").Qual("time", "Duration"),
			jen.ID("ConnectionDetails").ID("string"),
			jen.ID("DataManager").Interface(
				jen.ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8"), jen.ID("testUserConfig").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("error")),
				jen.ID("IsReady").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8")).Params(jen.ID("ready").ID("bool")),
				jen.ID("types").Dot("AdminUserDataManager"),
				jen.ID("types").Dot("AccountDataManager"),
				jen.ID("types").Dot("AccountUserMembershipDataManager"),
				jen.ID("types").Dot("UserDataManager"),
				jen.ID("types").Dot("AuditLogEntryDataManager"),
				jen.ID("types").Dot("APIClientDataManager"),
				jen.ID("types").Dot("WebhookDataManager"),
				jen.ID("types").Dot("ItemDataManager"),
				jen.ID("types").Dot("AdminAuditManager"),
				jen.ID("types").Dot("AuthAuditManager"),
			),
		),
		jen.Line(),
	)

	return code
}
