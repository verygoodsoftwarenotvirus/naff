package postgres

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrateDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Var().Defs(
			jen.ID("defaultTestUserTwoFactorSecret").Equals().Lit("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("testUserExistenceQuery").Equals().Lit(`
	SELECT users.id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.username = $1 AND users.two_factor_secret_verified_on IS NOT NULL
`),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("testUserCreationQuery").Equals().Lit(`
	INSERT INTO users (id,username,hashed_password,two_factor_secret,reputation,service_roles,two_factor_secret_verified_on) VALUES ($1,$2,$3,$4,$5,$6,extract(epoch FROM NOW()))
`),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("Migrate is a simple wrapper around the core querier Migrate call."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8"), jen.ID("testUserConfig").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("q").Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.If(jen.Op("!").ID("q").Dot("IsReady").Call(
				jen.ID("ctx"),
				jen.ID("maxAttempts"),
			)).Body(
				jen.Return().ID("database").Dot("ErrDatabaseNotReady")),
			jen.ID("q").Dot("migrateOnce").Dot("Do").Call(jen.ID("q").Dot("migrationFunc")),
			jen.If(jen.ID("testUserConfig").DoesNotEqual().Nil()).Body(
				jen.ID("q").Dot("logger").Dot("Debug").Call(jen.Lit("creating test user")),
				jen.ID("testUserExistenceArgs").Op(":=").Index().Interface().Valuesln(
					jen.ID("testUserConfig").Dot("Username")),
				jen.ID("userRow").Op(":=").ID("q").Dot("getOneRow").Call(
					jen.ID("ctx"),
					jen.ID("q").Dot("db"),
					jen.Lit("user"),
					jen.ID("testUserExistenceQuery"),
					jen.ID("testUserExistenceArgs"),
				),
				jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanUser").Call(
					jen.ID("ctx"),
					jen.ID("userRow"),
					jen.ID("false"),
				),
				jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
					jen.If(jen.ID("testUserConfig").Dot("ID").Op("==").Lit("")).Body(
						jen.ID("testUserConfig").Dot("ID").Equals().ID("ksuid").Dot("New").Call().Dot("String").Call()),
					jen.ID("testUserCreationArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("testUserConfig").Dot("ID"), jen.ID("testUserConfig").Dot("Username"), jen.ID("testUserConfig").Dot("HashedPassword"), jen.ID("defaultTestUserTwoFactorSecret"), jen.ID("types").Dot("GoodStandingAccountStatus"), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.ID("user").Op(":=").Op("&").ID("types").Dot("User").Valuesln(
						jen.ID("ID").MapAssign().ID("testUserConfig").Dot("ID"), jen.ID("Username").MapAssign().ID("testUserConfig").Dot("Username")),
					jen.ID("account").Op(":=").Op("&").ID("types").Dot("Account").Valuesln(
						jen.ID("ID").MapAssign().ID("ksuid").Dot("New").Call().Dot("String").Call()),
					jen.If(jen.ID("err").Equals().ID("q").Dot("createUser").Call(
						jen.ID("ctx"),
						jen.ID("user"),
						jen.ID("account"),
						jen.ID("testUserCreationQuery"),
						jen.ID("testUserCreationArgs"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("q").Dot("logger"),
							jen.ID("span"),
							jen.Lit("creating test user"),
						)),
					jen.ID("q").Dot("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("UsernameKey"),
						jen.ID("testUserConfig").Dot("Username"),
					).Dot("Debug").Call(jen.Lit("created test user and account")),
				),
			),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("initMigration").String(),
			jen.ID("itemsMigration").String(),
			jen.ID("migrations").Equals().Index().ID("darwin").Dot("Migration").Valuesln(
				jen.Valuesln(jen.ID("Version").MapAssign().Lit(0.01), jen.ID("Description").MapAssign().Lit("basic infrastructural tables"), jen.ID("Script").MapAssign().ID("initMigration")), jen.Valuesln(
					jen.ID("Version").MapAssign().Lit(0.02), jen.ID("Description").MapAssign().Lit("create items table"), jen.ID("Script").MapAssign().ID("itemsMigration"))),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Newline(),
		jen.Comment("migrate a postgres database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("migrationFunc").Params().Body(
			jen.ID("driver").Op(":=").ID("darwin").Dot("NewGenericDriver").Call(
				jen.ID("q").Dot("db"),
				jen.ID("darwin").Dot("PostgresDialect").Values(),
			),
			jen.If(jen.ID("err").Op(":=").ID("darwin").Dot("New").Call(
				jen.ID("driver"),
				jen.ID("migrations"),
				jen.Nil(),
			).Dot("Migrate").Call(), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("panic").Call(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("migrating database: %w"),
					jen.ID("err"),
				))),
		),
		jen.Newline(),
	)

	return code
}
