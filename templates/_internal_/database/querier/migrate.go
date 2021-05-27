package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrateDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("Migrate is a simple wrapper around the core querier Migrate call."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8"), jen.ID("testUserConfig").Op("*").ID("types").Dot("TestUserCreationConfig")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("q").Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.If(jen.Op("!").ID("q").Dot("IsReady").Call(
				jen.ID("ctx"),
				jen.ID("maxAttempts"),
			)).Body(
				jen.Return().ID("database").Dot("ErrDatabaseNotReady")),
			jen.ID("q").Dot("migrateOnce").Dot("Do").Call(jen.ID("q").Dot("sqlQueryBuilder").Dot("BuildMigrationFunc").Call(jen.ID("q").Dot("db"))),
			jen.If(jen.ID("testUserConfig").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("logger").Dot("Debug").Call(jen.Lit("creating test user")),
				jen.List(jen.ID("testUserExistenceQuery"), jen.ID("testUserExistenceArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetUserByUsernameQuery").Call(
					jen.ID("ctx"),
					jen.ID("testUserConfig").Dot("Username"),
				),
				jen.ID("userRow").Op(":=").ID("q").Dot("getOneRow").Call(
					jen.ID("ctx"),
					jen.ID("q").Dot("db"),
					jen.Lit("user"),
					jen.ID("testUserExistenceQuery"),
					jen.ID("testUserExistenceArgs").Op("..."),
				),
				jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanUser").Call(
					jen.ID("ctx"),
					jen.ID("userRow"),
					jen.ID("false"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.List(jen.ID("testUserCreationQuery"), jen.ID("testUserCreationArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildTestUserCreationQuery").Call(
						jen.ID("ctx"),
						jen.ID("testUserConfig"),
					),
					jen.ID("user").Op(":=").Op("&").ID("types").Dot("User").Valuesln(jen.ID("Username").Op(":").ID("testUserConfig").Dot("Username")),
					jen.ID("account").Op(":=").Op("&").ID("types").Dot("Account").Values(),
					jen.If(jen.ID("err").Op("=").ID("q").Dot("createUser").Call(
						jen.ID("ctx"),
						jen.ID("user"),
						jen.ID("account"),
						jen.ID("testUserCreationQuery"),
						jen.ID("testUserCreationArgs"),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
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
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
