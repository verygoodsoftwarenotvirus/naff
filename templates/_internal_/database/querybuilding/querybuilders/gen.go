package querybuilders

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/database/querybuilding"
)

func buildMariaDBWord() wordsmith.SuperPalabra {
	return &wordsmith.ManualWord{
		SingularStr:                           "MySQL",
		PluralStr:                             "MariaDBs",
		RouteNameStr:                          "mariadb",
		KebabNameStr:                          "mariadb",
		AbbreviationStr:                       "M",
		LowercaseAbbreviationStr:              "m",
		PluralRouteNameStr:                    "mariadbs",
		UnexportedVarNameStr:                  "mariaDB",
		PluralUnexportedVarNameStr:            "mariaDBs",
		PackageNameStr:                        "mariadbs",
		SingularPackageNameStr:                "mariadb",
		SingularCommonNameStr:                 "maria DB",
		ProperSingularCommonNameWithPrefixStr: "a Maria DB",
		PluralCommonNameStr:                   "maria DBs",
		SingularCommonNameWithPrefixStr:       "maria DB",
		PluralCommonNameWithPrefixStr:         "maria DBs",
	}
}

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, rawDBvendor := range proj.EnabledDatabases() {
		var dbvendor wordsmith.SuperPalabra

		switch rawDBvendor {
		case string(models.MySQL):
			dbvendor = buildMariaDBWord()
		case string(models.Postgres), string(models.Sqlite):
			dbvendor = wordsmith.FromSingularPascalCase(rawDBvendor)
		}

		files := map[string]string{
			"generic.go":      genericDotGo(proj, dbvendor),
			"generic_test.go": genericTestDotGo(proj, dbvendor),
			fmt.Sprintf("%s.go", dbvendor.RouteName()):      mainDotGo(proj, dbvendor),
			fmt.Sprintf("%s_test.go", dbvendor.RouteName()): mainTestDotGo(proj, dbvendor),
		}

		for path, file := range files {
			if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, dbvendor.RouteName(), path), file); err != nil {
				return err
			}
		}

		jenFiles := map[string]*jen.File{
			"accounts.go":                      accountsDotGo(proj, dbvendor),
			"accounts_test.go":                 accountsTestDotGo(proj, dbvendor),
			"account_user_memberships.go":      accountUserMembershipsDotGo(proj, dbvendor),
			"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj, dbvendor),
			"api_clients.go":                   apiClientsDotGo(proj, dbvendor),
			"api_clients_test.go":              apiClientsTestDotGo(proj, dbvendor),
			"audit_log_entries.go":             auditLogEntriesDotGo(proj, dbvendor),
			"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj, dbvendor),
			"doc.go":                           docDotGo(proj, dbvendor),
			"migrations.go":                    migrationsDotGo(proj, dbvendor),
			"users.go":                         usersDotGo(proj, dbvendor),
			"users_test.go":                    usersTestDotGo(proj, dbvendor),
			"webhooks.go":                      webhooksDotGo(proj, dbvendor),
			"webhooks_test.go":                 webhooksTestDotGo(proj, dbvendor),
			"wire.go":                          wireDotGo(proj, dbvendor),
		}

		for _, typ := range proj.DataTypes {
			jenFiles[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ, dbvendor)
			jenFiles[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ, dbvendor)
		}

		for path, file := range jenFiles {
			if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, dbvendor.RouteName(), path), file); err != nil {
				return err
			}
		}
	}

	return nil
}

//go:embed templates/generic.gotpl
var genericTemplate string

func genericDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) string {
	joinIDsDef := `statement := ""

	for i, id := range ids {
		if i != 0 {
			statement += " "
		}
		statement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}

	statement += " END"

	return statement`

	if dbvendor.SingularPackageName() == "postgres" {
		joinIDsDef = `out := []string{}

	for _, x := range ids {
		out = append(out, strconv.FormatUint(x, 10))
	}

	return strings.Join(out, ",")`
	}

	generated := map[string]string{
		"packageName": dbvendor.SingularPackageName(),
		"structName":  dbvendor.Singular(),
		"joinIDs":     joinIDsDef,
	}
	return models.RenderCodeFile(proj, genericTemplate, generated)
}

//go:embed templates/generic_test.gotpl
var genericTestTemplate string

func genericTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) string {
	generated := map[string]string{
		"packageName":         dbvendor.SingularPackageName(),
		"structName":          dbvendor.Singular(),
		"firstExpectedQuery":  fmt.Sprintf("SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_on IS NULL AND example_table.belongs_to_account = %s AND key = %s) as total_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_on IS NULL AND example_table.belongs_to_account = %s AND key = %s AND example_table.created_on > %s AND example_table.created_on < %s AND example_table.last_updated_on > %s AND example_table.last_updated_on < %s) as filtered_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_on IS NULL AND example_table.belongs_to_account = %s AND key = %s AND example_table.created_on > %s AND example_table.created_on < %s AND example_table.last_updated_on > %s AND example_table.last_updated_on < %s GROUP BY example_table.id LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7), getIncIndex(dbvendor, 8), getIncIndex(dbvendor, 9), getIncIndex(dbvendor, 10), getIncIndex(dbvendor, 11), getIncIndex(dbvendor, 12), getIncIndex(dbvendor, 13)),
		"secondExpectedQuery": fmt.Sprintf("SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.archived_on IS NULL) as total_count, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.archived_on IS NULL AND example_table.created_on > %s AND example_table.created_on < %s AND example_table.last_updated_on > %s AND example_table.last_updated_on < %s) as filtered_count FROM example_table WHERE example_table.created_on > %s AND example_table.created_on < %s AND example_table.last_updated_on > %s AND example_table.last_updated_on < %s GROUP BY example_table.id LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7)),
		"thirdExpectedQuery":  fmt.Sprintf("SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table) as total_count, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.created_on > %s AND example_table.created_on < %s AND example_table.last_updated_on > %s AND example_table.last_updated_on < %s) as filtered_count FROM example_table WHERE example_table.created_on > %s AND example_table.created_on < %s AND example_table.last_updated_on > %s AND example_table.last_updated_on < %s GROUP BY example_table.id LIMIT 20 OFFSET 180", getIncIndex(dbvendor, 0), getIncIndex(dbvendor, 1), getIncIndex(dbvendor, 2), getIncIndex(dbvendor, 3), getIncIndex(dbvendor, 4), getIncIndex(dbvendor, 5), getIncIndex(dbvendor, 6), getIncIndex(dbvendor, 7)),
	}
	return models.RenderCodeFile(proj, genericTestTemplate, generated)
}

var (
	//go:embed templates/main_mariadb.gotpl
	mariadbMainTemplate string
	//go:embed templates/main_postgres.gotpl
	postgresMainTemplate string
	//go:embed templates/main_sqlite.gotpl
	sqliteMainTemplate string
)

func mainDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) string {
	switch dbvendor.SingularPackageName() {
	case "mariadb":
		return models.RenderCodeFile(proj, mariadbMainTemplate, nil)
	case "postgres":
		return models.RenderCodeFile(proj, postgresMainTemplate, nil)
	case "sqlite":
		return models.RenderCodeFile(proj, sqliteMainTemplate, nil)
	default:
		panic(fmt.Sprintf("invalid database type! %q", dbvendor.SingularPackageName()))
	}
}

var (
	//go:embed templates/main_mariadb_test.gotpl
	mariadbMainTestTemplate string
	//go:embed templates/main_postgres_test.gotpl
	postgresMainTestTemplate string
	//go:embed templates/main_sqlite_test.gotpl
	sqliteMainTestTemplate string
)

func mainTestDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra) string {
	switch dbvendor.SingularPackageName() {
	case "mariadb":
		return models.RenderCodeFile(proj, mariadbMainTestTemplate, nil)
	case "postgres":
		return models.RenderCodeFile(proj, postgresMainTestTemplate, nil)
	case "sqlite":
		return models.RenderCodeFile(proj, sqliteMainTestTemplate, nil)
	default:
		panic(fmt.Sprintf("invalid database type! %q", dbvendor.LowercaseAbbreviation()))
	}
}
