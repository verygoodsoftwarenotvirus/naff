package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("DevelopmentRunMode").ID("runMode").Op("=").Lit("development").Var().ID("TestingRunMode").ID("runMode").Op("=").Lit("testing").Var().ID("ProductionRunMode").ID("runMode").Op("=").Lit("production").Var().ID("DefaultRunMode").Op("=").ID("DevelopmentRunMode").Var().ID("DefaultStartupDeadline").Op("=").Qual("time", "Minute"),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("errNilDatabaseConnection").Op("=").Qual("errors", "New").Call(jen.Lit("nil DB connection provided")).Var().ID("errInvalidDatabaseProvider").Op("=").Qual("errors", "New").Call(jen.Lit("invalid database provider")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("runMode").ID("string").Type().ID("ServerConfig").Struct(
			jen.ID("Search").ID("search").Dot("Config"),
			jen.ID("Encoding").ID("encoding").Dot("Config"),
			jen.ID("Capitalism").ID("capitalism").Dot("Config"),
			jen.ID("Uploads").ID("uploads").Dot("Config"),
			jen.ID("Observability").ID("observability").Dot("Config"),
			jen.ID("Routing").ID("routing").Dot("Config"),
			jen.ID("Meta").ID("MetaSettings"),
			jen.ID("Database").ID("config").Dot("Config"),
			jen.ID("Auth").ID("authentication").Dot("Config"),
			jen.ID("Server").ID("server").Dot("Config"),
			jen.ID("AuditLog").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/audit", "Config"),
			jen.ID("Webhooks").ID("webhooks").Dot("Config"),
			jen.ID("Frontend").ID("frontend").Dot("Config"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeToFile renders your config to a file given your favorite encoder."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("EncodeToFile").Params(jen.ID("path").ID("string"), jen.ID("marshaller").Params(jen.ID("v").Interface()).Params(jen.Index().ID("byte"), jen.ID("error"))).Params(jen.ID("error")).Body(
			jen.List(jen.ID("byteSlice"), jen.ID("err")).Op(":=").ID("marshaller").Call(jen.Op("*").ID("cfg")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("err")),
			jen.Return().Qual("os", "WriteFile").Call(
				jen.ID("path"),
				jen.ID("byteSlice"),
				jen.Lit(600),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("ServerConfig")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a ServerConfig struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Search").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Search portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Uploads").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Uploads portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Routing").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Routing portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Meta").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Meta portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Capitalism").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Capitalism portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Encoding").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Encoding portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Encoding").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Encoding portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Observability").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Observability portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Database").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Database portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Auth").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Auth portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Server").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating HTTPServer portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Webhooks").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating Webhooks portion of config: %w"),
					jen.ID("err"),
				)),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("AuditLog").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("error validating AuditLog portion of config: %w"),
					jen.ID("err"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideDatabaseClient provides a database implementation dependent on the configuration."),
		jen.Line(),
		jen.Func().Comment("NOTE: you may be tempted to move this to the database/config package. This is a fool's errand.").Params(jen.ID("cfg").Op("*").ID("ServerConfig")).ID("ProvideDatabaseClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("rawDB").Op("*").Qual("database/sql", "DB")).Params(jen.ID("database").Dot("DataManager"), jen.ID("error")).Body(
			jen.If(jen.ID("rawDB").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilDatabaseConnection"))),
			jen.Var().ID("qb").ID("querybuilding").Dot("SQLQueryBuilder"),
			jen.ID("shouldCreateTestUser").Op(":=").ID("cfg").Dot("Meta").Dot("RunMode").Op("!=").ID("ProductionRunMode"),
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("cfg").Dot("Database").Dot("Provider")))).Body(
				jen.Case(jen.Lit("sqlite")).Body(
					jen.ID("qb").Op("=").ID("sqlite").Dot("ProvideSqlite").Call(jen.ID("logger"))),
				jen.Case(jen.Lit("mariadb")).Body(
					jen.ID("qb").Op("=").ID("mariadb").Dot("ProvideMariaDB").Call(jen.ID("logger"))),
				jen.Case(jen.Lit("postgres")).Body(
					jen.ID("qb").Op("=").ID("postgres").Dot("ProvidePostgres").Call(jen.ID("logger"))),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %q"),
						jen.ID("errInvalidDatabaseProvider"),
						jen.ID("cfg").Dot("Database").Dot("Provider"),
					))),
			),
			jen.Return().ID("querier").Dot("ProvideDatabaseClient").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("rawDB"),
				jen.Op("&").ID("cfg").Dot("Database"),
				jen.ID("qb"),
				jen.ID("shouldCreateTestUser"),
			),
		),
		jen.Line(),
	)

	return code
}
