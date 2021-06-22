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
		SingularStr:                           "MariaDB",
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
		case string(models.MariaDB):
			dbvendor = buildMariaDBWord()
		case string(models.Postgres), string(models.Sqlite):
			dbvendor = wordsmith.FromSingularPascalCase(rawDBvendor)
		}

		files := map[string]string{
			"generic.go":      genericDotGo(proj),
			"generic_test.go": genericTestDotGo(proj),
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

func genericDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, genericTemplate, nil)
}

//go:embed templates/generic_test.gotpl
var genericTestTemplate string

func genericTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, genericTestTemplate, nil)
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
	case "m":
		return models.RenderCodeFile(proj, mariadbMainTemplate, nil)
	case "p":
		return models.RenderCodeFile(proj, postgresMainTemplate, nil)
	case "s":
		return models.RenderCodeFile(proj, sqliteMainTemplate, nil)
	default:
		panic(fmt.Sprintf("invalid database type! %q", dbvendor.LowercaseAbbreviation()))
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
	case "m":
		return models.RenderCodeFile(proj, mariadbMainTestTemplate, nil)
	case "p":
		return models.RenderCodeFile(proj, postgresMainTestTemplate, nil)
	case "s":
		return models.RenderCodeFile(proj, sqliteMainTestTemplate, nil)
	default:
		panic(fmt.Sprintf("invalid database type! %q", dbvendor.LowercaseAbbreviation()))
	}
}
