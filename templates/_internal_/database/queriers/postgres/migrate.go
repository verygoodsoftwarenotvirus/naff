package postgres

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func migrateDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(buildConstantsBlock()...)
	code.Add(buildMigrate(proj)...)
	code.Add(buildMigrationScriptDeclarations(proj)...)
	code.Add(buildBuildMigrationFunc()...)

	return code
}

func buildConstantsBlock() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("defaultTestUserTwoFactorSecret is the default TwoFactorSecret we give to test users when we initialize them."),
			jen.Comment("`otpauth://totp/todo:username?secret=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=&issuer=todo`"),
			jen.ID("defaultTestUserTwoFactorSecret").Equals().Lit("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
			jen.Newline(),
			jen.ID("testUserExistenceQuery").Equals().RawString(`
		SELECT users.id, users.username, users.avatar_src, users.hashed_password, users.requires_password_change, users.password_last_changed_on, users.two_factor_secret, users.two_factor_secret_verified_on, users.service_roles, users.reputation, users.reputation_explanation, users.created_on, users.last_updated_on, users.archived_on FROM users WHERE users.archived_on IS NULL AND users.username = $1 AND users.two_factor_secret_verified_on IS NOT NULL
	`),
			jen.Newline(),
			jen.ID("testUserCreationQuery").Equals().RawString(`
		INSERT INTO users (id,username,hashed_password,two_factor_secret,reputation,service_roles,two_factor_secret_verified_on) VALUES ($1,$2,$3,$4,$5,$6,extract(epoch FROM NOW()))
	`),
		),
		jen.Newline(),
	}

	return lines
}

func buildMigrate(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("Migrate is a simple wrapper around the core querier Migrate call."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("Migrate").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("maxAttempts").ID("uint8"), jen.ID("testUserConfig").Op("*").Qual(proj.TypesPackage(), "TestUserCreationConfig")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("q").Dot("logger").Dot("Info").Call(jen.Lit("migrating db")),
			jen.Newline(),
			jen.If(jen.Op("!").ID("q").Dot("IsReady").Call(
				jen.ID("ctx"),
				jen.ID("maxAttempts"),
			)).Body(
				jen.Return().Qual(proj.DatabasePackage(), "ErrDatabaseNotReady")),
			jen.Newline(),
			jen.ID("q").Dot("migrateOnce").Dot("Do").Call(jen.ID("q").Dot("migrationFunc")),
			jen.Newline(),
			jen.If(jen.ID("testUserConfig").DoesNotEqual().Nil()).Body(
				jen.ID("q").Dot("logger").Dot("Debug").Call(jen.Lit("creating test user")),
				jen.Newline(),
				jen.ID("testUserExistenceArgs").Assign().Index().Interface().Values(jen.ID("testUserConfig").Dot("Username")),
				jen.Newline(),
				jen.ID("userRow").Assign().ID("q").Dot("getOneRow").Call(
					jen.ID("ctx"),
					jen.ID("q").Dot("db"),
					jen.Lit("user"),
					jen.ID("testUserExistenceQuery"),
					jen.ID("testUserExistenceArgs"),
				),
				jen.List(jen.ID("_"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dot("scanUser").Call(
					jen.ID("ctx"),
					jen.ID("userRow"),
					jen.False(),
				),
				jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
					jen.If(jen.ID("testUserConfig").Dot("ID").Op("==").Lit("")).Body(
						jen.ID("testUserConfig").Dot("ID").Equals().ID("ksuid").Dot("New").Call().Dot("String").Call(),
					),
					jen.Newline(),
					jen.ID("testUserCreationArgs").Assign().Index().Interface().Valuesln(
						jen.ID("testUserConfig").Dot("ID"), jen.ID("testUserConfig").Dot("Username"), jen.ID("testUserConfig").Dot("HashedPassword"), jen.ID("defaultTestUserTwoFactorSecret"), jen.Qual(proj.TypesPackage(), "GoodStandingAccountStatus"), jen.ID("authorization").Dot("ServiceAdminRole").Dot("String").Call()),
					jen.Newline(),
					jen.Comment("these structs will be fleshed out by createUser"),
					jen.ID("user").Assign().Op("&").Qual(proj.TypesPackage(), "User").Valuesln(
						jen.ID("ID").MapAssign().ID("testUserConfig").Dot("ID"), jen.ID("Username").MapAssign().ID("testUserConfig").Dot("Username"),
					),
					jen.ID("account").Assign().Op("&").Qual(proj.TypesPackage(), "Account").Valuesln(
						jen.ID("ID").MapAssign().ID("ksuid").Dot("New").Call().Dot("String").Call(),
					),
					jen.Newline(),
					jen.If(jen.ID("err").Equals().ID("q").Dot("createUser").Call(
						jen.ID("ctx"),
						jen.ID("user"),
						jen.ID("account"),
						jen.ID("testUserCreationQuery"),
						jen.ID("testUserCreationArgs"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
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
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	}

	return lines
}

func buildMigrationScriptDeclarations(proj *models.Project) []jen.Code {
	migrationScriptDecls := []jen.Code{
		jen.Commentf("//go:embed migrations/00001_initial.sql"),
		jen.ID("initMigration").String(),
		jen.Newline(),
	}
	migrationDecls := []jen.Code{
		jen.Valuesln(
			jen.ID("Version").MapAssign().Lit(0.01),
			jen.ID("Description").MapAssign().Lit("basic infrastructural tables"),
			jen.ID("Script").MapAssign().ID("initMigration"),
		),
	}

	for i, typ := range proj.DataTypes {
		migrationScriptDecls = append(migrationScriptDecls,
			jen.Commentf("//go:embed migrations/0000%d_%s.sql", i+2, typ.Name.PluralRouteName()),
			jen.IDf("%sMigration", typ.Name.PluralUnexportedVarName()).String(),
			jen.Newline(),
		)
		migrationDecls = append(migrationDecls,
			jen.Valuesln(
				jen.ID("Version").MapAssign().Lit(0.01*float64(i+2)),
				jen.ID("Description").MapAssign().Litf("create %s table", typ.Name.PluralCommonName()),
				jen.ID("Script").MapAssign().IDf("%sMigration", typ.Name.PluralUnexportedVarName()),
			),
		)
	}

	migrationScriptDecls = append(migrationScriptDecls,
		jen.ID("migrations").Equals().Index().Qual("github.com/GuiaBolso/darwin", "Migration").Valuesln(
			migrationDecls...,
		),
	)

	return []jen.Code{
		jen.Var().Defs(migrationScriptDecls...),
		jen.Newline(),
	}
}

func buildBuildMigrationFunc() []jen.Code {
	lines := []jen.Code{
		jen.Comment("BuildMigrationFunc returns a sync.Once compatible function closure that will"),
		jen.Newline(),
		jen.Comment("migrate a postgres database."),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("migrationFunc").Params().Body(
			jen.ID("driver").Assign().Qual("github.com/GuiaBolso/darwin", "NewGenericDriver").Call(
				jen.ID("q").Dot("db"),
				jen.Qual("github.com/GuiaBolso/darwin", "PostgresDialect").Values(),
			),
			jen.If(jen.ID("err").Assign().Qual("github.com/GuiaBolso/darwin", "New").Call(
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
	}

	return lines
}

func typeToPostgresType(t string) string {
	typeMap := map[string]string{
		"string":   "TEXT",
		"*string":  "TEXT",
		"bool":     "BOOLEAN",
		"*bool":    "BOOLEAN",
		"int":      "INTEGER",
		"*int":     "INTEGER",
		"int8":     "INTEGER",
		"*int8":    "INTEGER",
		"int16":    "INTEGER",
		"*int16":   "INTEGER",
		"int32":    "INTEGER",
		"*int32":   "INTEGER",
		"int64":    "INTEGER",
		"*int64":   "INTEGER",
		"uint":     "INTEGER",
		"*uint":    "INTEGER",
		"uint8":    "INTEGER",
		"*uint8":   "INTEGER",
		"uint16":   "INTEGER",
		"*uint16":  "INTEGER",
		"uint32":   "BIGINT",
		"*uint32":  "BIGINT",
		"uint64":   "BIGINT",
		"*uint64":  "BIGINT",
		"float32":  "DOUBLE PRECISION",
		"*float32": "DOUBLE PRECISION",
		"float64":  "DOUBLE PRECISION",
		"*float64": "DOUBLE PRECISION",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	panic(fmt.Sprintf("unknown type!: %q", t))
}

func buildMigrationScript(typ models.DataType) string {
	scriptParts := []string{
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", typ.Name.PluralRouteName()),
		`	id CHAR(27) NOT NULL PRIMARY KEY`,
	}

	for _, field := range typ.Fields {
		rn := field.Name.RouteName()

		query := fmt.Sprintf("	%s %s", rn, typeToPostgresType(field.Type))
		if !field.IsPointer {
			query += ` NOT NULL`
		}
		if field.DefaultValue != "" {
			query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
		}

		scriptParts = append(scriptParts, fmt.Sprintf("%s", query))
	}

	scriptParts = append(scriptParts,
		`	created_on BIGINT NOT NULL DEFAULT extract(epoch FROM NOW())`,
		`	last_updated_on BIGINT DEFAULT NULL`,
	)

	if !typ.BelongsToAccount && typ.BelongsToStruct == nil {
		scriptParts = append(scriptParts,
			`	archived_on BIGINT DEFAULT NULL`,
		)
	} else {
		scriptParts = append(scriptParts,
			`	archived_on BIGINT DEFAULT NULL`, // note the comma
		)
	}

	if typ.BelongsToAccount {
		scriptParts = append(scriptParts,
			`	belongs_to_account CHAR(27) NOT NULL REFERENCES accounts(id) ON DELETE CASCADE`,
		)
	}
	if typ.BelongsToStruct != nil {
		scriptParts = append(scriptParts,
			fmt.Sprintf(`	"belongs_to_%s" CHAR(27) NOT NULL REFERENCES %s(id) ON DELETE CASCADE`, typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
		)
	}

	for i, sp := range scriptParts {
		if i != len(scriptParts)-1 {
			if !strings.HasSuffix(sp, "(") {
				scriptParts[i] = fmt.Sprintf("%s,", sp)
			}
		}
	}

	scriptParts = append(scriptParts,
		`);`,
	)

	return strings.Join(scriptParts, "\n")
}
