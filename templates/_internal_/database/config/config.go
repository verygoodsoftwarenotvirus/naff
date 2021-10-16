package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres), jen.Comment("PostgresProvider is the string used to refer to postgres.").Newline().ID("PostgresProvider").Equals().Lit("postgres")),
			utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL), jen.Comment("MySQLProvider is the string used to refer to MySQL.").Newline().ID("MySQLProvider").Equals().Lit("mysql")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Type().Defs(
			jen.Comment("Config represents our database configuration."),
			jen.ID("Config").Struct(
				jen.Underscore().Struct(),
				jen.Newline(),
				jen.ID("CreateTestUser").Op("*").Qual(proj.TypesPackage(), "TestUserCreationConfig").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("CreateTestUser"), false)),
				jen.ID("Provider").String().Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Provider"), false)),
				jen.ID("ConnectionDetails").Qual(proj.DatabasePackage(), "ConnectionDetails").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("ConnectionDetails"), false)),
				jen.ID("Debug").ID("bool").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Debug"), false)),
				jen.ID("RunMigrations").ID("bool").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("RunMigrations"), false)),
				jen.ID("MaxPingAttempts").ID("uint8").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("MaxPingAttempts"), false)),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Equals().Parens(jen.Op("*").ID("Config")).Call(jen.Nil()),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates an DatabaseSettings struct."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Callln(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("ConnectionDetails"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("Provider"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.Postgres), jen.ID("PostgresProvider")),
						utils.ConditionalCode(proj.DatabaseIsEnabled(models.MySQL), jen.ID("MySQLProvider")),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("cfg").Dot("CreateTestUser"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("cfg").Dot("CreateTestUser").DoesNotEqual().Nil(),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Nil")),
				),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideSessionManager provides a session manager based on some settings."),
		jen.Newline(),
		jen.Comment("There's not a great place to put this function. I don't think it belongs in Auth because it accepts a DB connection,"),
		jen.Newline(),
		jen.Comment("but it obviously doesn't belong in the database package, or maybe it does."),
		jen.Newline(),
		jen.Func().ID("ProvideSessionManager").Params(
			jen.ID("cookieConfig").Qual(proj.AuthServicePackage(), "CookieConfig"),
			jen.ID("dm").Qual(proj.DatabasePackage(), "DataManager"),
		).Params(jen.Op("*").Qual("github.com/alexedwards/scs/v2", "SessionManager"), jen.ID("error")).Body(
			jen.ID("sessionManager").Assign().Qual("github.com/alexedwards/scs/v2", "New").Call(),
			jen.Newline(),
			jen.ID("sessionManager").Dot("Lifetime").Equals().ID("cookieConfig").Dot("Lifetime"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Name").Equals().ID("cookieConfig").Dot("Name"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Domain").Equals().ID("cookieConfig").Dot("Domain"),
			jen.ID("sessionManager").Dot("Cookie").Dot("HttpOnly").Equals().True(),
			jen.ID("sessionManager").Dot("Cookie").Dot("Path").Equals().Lit("/"),
			jen.ID("sessionManager").Dot("Cookie").Dot("SameSite").Equals().Qual("net/http", "SameSiteStrictMode"),
			jen.ID("sessionManager").Dot("Cookie").Dot("Secure").Equals().ID("cookieConfig").Dot("SecureOnly"),
			jen.Newline(),
			jen.ID("sessionManager").Dot("Store").Equals().ID("dm").Dot("ProvideSessionStore").Call(),
			jen.Newline(),
			jen.Return().List(jen.ID("sessionManager"), jen.Nil()),
		),
		jen.Newline(),
	)

	return code
}
