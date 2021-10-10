package config

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.Comment("DevelopmentRunMode is the run mode for a development environment."),
			jen.ID("DevelopmentRunMode").ID("runMode").Equals().Lit("development"),
			jen.Comment("TestingRunMode is the run mode for a testing environment."),
			jen.ID("TestingRunMode").ID("runMode").Equals().Lit("testing"),
			jen.Comment("ProductionRunMode is the run mode for a production environment."),
			jen.ID("ProductionRunMode").ID("runMode").Equals().Lit("production"),
			jen.Comment("DefaultRunMode is the default run mode."),
			jen.ID("DefaultRunMode").Equals().ID("DevelopmentRunMode"),
			jen.Comment("DefaultStartupDeadline is the default amount of time we allow for server startup."),
			jen.ID("DefaultStartupDeadline").Equals().Qual("time", "Minute"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errNilDatabaseConnection").Equals().Qual("errors", "New").Call(jen.Lit("nil DB connection provided")),
			jen.ID("errNilConfig").Equals().Qual("errors", "New").Call(jen.Lit("nil config provided")),
			jen.ID("errInvalidDatabaseProvider").Equals().Qual("errors", "New").Call(jen.Lit("invalid database provider")),
		),
		jen.Newline(),
	)

	serviceConfigurations := []jen.Code{}
	for _, typ := range proj.DataTypes {
		serviceConfigurations = append(serviceConfigurations, jen.ID(typ.Name.Plural()).Qual(proj.ServicePackage(typ.Name.PackageName()), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase(typ.Name.Singular()), true)))
	}

	serviceConfigurations = append(serviceConfigurations,
		jen.ID("Auth").Qual(proj.AuthServicePackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Auth"), false)),
		jen.ID("Webhooks").Qual(proj.WebhooksServicePackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Webhook"), true)),
		jen.ID("Frontend").Qual(proj.FrontendServicePackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Frontend"), false)),
	)

	code.Add(
		jen.Type().Defs(
			jen.Comment("runMode describes what method of operation the server is under."),
			jen.ID("runMode").String(),
			jen.Newline(),
			jen.Comment("ServicesConfigurations collects the various service configurations."),
			jen.ID("ServicesConfigurations").Struct(
				serviceConfigurations...,
			),
			jen.Newline(),
			jen.Comment("InstanceConfig configures an instance of the service. It is composed of all the other setting structs."),
			jen.ID("InstanceConfig").Struct(
				jen.ID("Search").Qual(proj.InternalSearchPackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Search"), false)),
				jen.ID("Encoding").Qual(proj.EncodingPackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Encoding"), false)),
				jen.ID("Uploads").Qual(proj.UploadsPackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Upload"), true)),
				jen.ID("Observability").Qual(proj.ObservabilityPackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Observability"), false)),
				jen.ID("Routing").Qual(proj.RoutingPackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Routing"), false)),
				jen.ID("Meta").ID("MetaSettings").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Meta"), false)),
				jen.ID("Database").Qual(proj.DatabasePackage("config"), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Database"), false)),
				jen.ID("Services").ID("ServicesConfigurations").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Service"), true)),
				jen.ID("Server").Qual(proj.HTTPServerPackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Server"), false)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("EncodeToFile renders your config to a file given your favorite encoder."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("InstanceConfig")).ID("EncodeToFile").Params(jen.ID("path").String(),
			jen.ID("marshaller").Func().Params(jen.ID("v").Interface()).Params(jen.Index().ID("byte"),
				jen.ID("error"))).Params(jen.ID("error")).Body(
			jen.List(jen.ID("byteSlice"),
				jen.ID("err")).Assign().ID("marshaller").Call(jen.PointerTo().ID("cfg")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("err"),
			),
			jen.Newline(),
			jen.Return().Qual("os", "WriteFile").Call(
				jen.ID("path"),
				jen.ID("byteSlice"),
				jen.Octal(600),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Underscore().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Equals().Parens(jen.PointerTo().ID("InstanceConfig")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	validateLines := []jen.Code{
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Search").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Search portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Uploads").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Uploads portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Routing").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Routing portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Meta").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Meta portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Capitalism").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Capitalism portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Encoding").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Encoding portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Encoding").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Encoding portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Observability").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Observability portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Database").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Database portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Server").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating HTTPServer portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Services").Dot("AuditLog").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating AuditLog portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Services").Dot("Auth").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Auth service portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Services").Dot("Frontend").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Frontend service portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("cfg").Dot("Services").Dot("Webhooks").Dot("ValidateWithContext").Call(jen.ID("ctx")),
			jen.ID("err").Op("!=").ID("nil")).Body(
			jen.Return().Qual("fmt", "Errorf").Call(
				jen.Lit("error validating Webhooks service portion of config: %w"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()

		validateLines = append(validateLines,
			jen.If(jen.Err().Assign().ID("cfg").Dot("Services").Dot(pn).Dot("ValidateWithContext").Call(jen.ID("ctx")),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit(fmt.Sprintf("error validating %s service portion of config: ", pn)+"%w"),
					jen.ID("err"),
				),
			),
			jen.Newline(),
		)
	}

	validateLines = append(validateLines, jen.Return().Nil())

	code.Add(
		jen.Comment("ValidateWithContext validates a InstanceConfig struct."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("InstanceConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			validateLines...,
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideDatabaseClient provides a database implementation dependent on the configuration."),
		jen.Newline(),
		jen.Comment("NOTE: you may be tempted to move this to the database/config package. This is a fool's errand."),
		jen.Newline(),
		jen.Func().ID("ProvideDatabaseClient").Params(
			jen.ID("ctx").Qual("context", "Context"),
			constants.LoggerVar().Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("rawDB").PointerTo().Qual("database/sql", "DB"),
			jen.ID("cfg").PointerTo().ID("InstanceConfig"),
		).Params(
			jen.Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("error"),
		).Body(
			jen.If(jen.ID("rawDB").IsEqualTo().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("errNilDatabaseConnection"),
				),
			),
			jen.Newline(),
			jen.If(jen.ID("cfg").IsEqualTo().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("errNilConfig"),
				),
			),
			jen.Newline(),
			jen.Var().ID("qb").Qual(proj.DatabasePackage("querybuilding"), "SQLQueryBuilder"),
			jen.ID("shouldCreateTestUser").Assign().ID("cfg").Dot("Meta").Dot("RunMode").Op("!=").ID("ProductionRunMode"),
			jen.Newline(),
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("cfg").Dot("Database").Dot("Provider")))).Body(
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL), jen.Case(jen.Lit("mariadb")).Body(
					jen.ID("qb").Equals().Qual(proj.DatabasePackage("querybuilding", "mariadb"), "ProvideMariaDB").Call(constants.LoggerVar()),
				)),
				utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres), jen.Case(jen.Lit("postgres")).Body(
					jen.ID("qb").Equals().Qual(proj.DatabasePackage("querybuilding", "postgres"), "ProvidePostgres").Call(constants.LoggerVar()),
				)),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"),
						jen.Qual("fmt", "Errorf").Call(
							jen.Lit("%w: %q"),
							jen.ID("errInvalidDatabaseProvider"),
							jen.ID("cfg").Dot("Database").Dot("Provider"),
						),
					),
				),
			),
			jen.Newline(),
			jen.Return().Qual(proj.DatabasePackage("querier"), "ProvideDatabaseClient").Call(
				jen.ID("ctx"),
				constants.LoggerVar(),
				jen.ID("rawDB"),
				jen.AddressOf().ID("cfg").Dot("Database"),
				jen.ID("qb"),
				jen.ID("shouldCreateTestUser"),
			),
		),
		jen.Newline(),
	)

	return code
}
