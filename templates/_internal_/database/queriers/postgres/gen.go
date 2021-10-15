package postgres

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/Masterminds/squirrel"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	whatever        = "fart"
	packageName     = "postgres"
	basePackagePath = "internal/database/queriers/postgres"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"accounts.go":                      accountsDotGo(proj),
		"accounts_test.go":                 accountsTestDotGo(proj),
		"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
		"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
		"admin.go":                         adminDotGo(proj),
		"admin_test.go":                    adminTestDotGo(proj),
		"api_clients.go":                   apiClientsDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
		"doc.go":                           docDotGo(proj),
		"errors.go":                        errorsDotGo(proj),
		"postgres.go":                      postgresDotGo(proj),
		"postgres_test.go":                 postgresTestDotGo(proj),
		"queries.go":                       queriesDotGo(proj),
		"queries_test.go":                  queriesTestDotGo(proj),
		"query_filters.go":                 queryFiltersDotGo(proj),
		"query_filter_test.go":             queryFilterTestDotGo(proj),
		"users.go":                         usersDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
		"webhooks.go":                      webhooksDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	migrations := map[string]string{
		"migrations/00001_initial.sql": baseMigrationDotSQL(proj),
	}

	for path, file := range migrations {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, false); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{
		"migrate.go":      migrateDotGo(proj),
		"migrate_test.go": migrateTestDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ, wordsmith.FromSingularPascalCase("Postgres"))
		jenFiles[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ, wordsmith.FromSingularPascalCase("Postgres"))
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

func queryBuilderForDatabase(db wordsmith.SuperPalabra) squirrel.StatementBuilderType {
	switch db.LowercaseAbbreviation() {
	case "p":
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	case "s", "m":
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	default:
		panic(fmt.Sprintf("invalid database type! %q", db.LowercaseAbbreviation()))
	}
}

func unixTimeForDatabase(db wordsmith.SuperPalabra) string {
	switch db.LowercaseAbbreviation() {
	case "m":
		return "UNIX_TIMESTAMP()"
	case "p":
		return "extract(epoch FROM NOW())"
	case "s":
		return "(strftime('%s','now'))"
	default:
		panic(fmt.Sprintf("invalid database type! %q", db.LowercaseAbbreviation()))
	}
}

//go:embed migrations/00001_initial.sql
var baseMigration string

func baseMigrationDotSQL(proj *models.Project) string {
	return models.RenderCodeFile(proj, baseMigration, nil)
}

//go:embed accounts.gotpl
var accountsTemplate string

func accountsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTemplate, nil)
}

//go:embed accounts_test.gotpl
var accountsTestTemplate string

func accountsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountsTestTemplate, nil)
}

//go:embed account_user_memberships.gotpl
var accountUserMembershipsTemplate string

func accountUserMembershipsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipsTemplate, nil)
}

//go:embed account_user_memberships_test.gotpl
var accountUserMembershipsTestTemplate string

func accountUserMembershipsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipsTestTemplate, nil)
}

//go:embed admin.gotpl
var adminTemplate string

func adminDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTemplate, nil)
}

//go:embed admin_test.gotpl
var adminTestTemplate string

func adminTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminTestTemplate, nil)
}

//go:embed api_clients.gotpl
var apiClientsTemplate string

func apiClientsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTemplate, nil)
}

//go:embed api_clients_test.gotpl
var apiClientsTestTemplate string

func apiClientsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientsTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed errors.gotpl
var errorsTemplate string

func errorsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, errorsTemplate, nil)
}

//go:embed postgres.gotpl
var postgresTemplate string

func postgresDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, postgresTemplate, nil)
}

//go:embed postgres_test.gotpl
var postgresTestTemplate string

func postgresTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, postgresTestTemplate, nil)
}

//go:embed queries.gotpl
var queriesTemplate string

func queriesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queriesTemplate, nil)
}

//go:embed queries_test.gotpl
var queriesTestTemplate string

func queriesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queriesTestTemplate, nil)
}

//go:embed query_filters.gotpl
var queryFiltersTemplate string

func queryFiltersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFiltersTemplate, nil)
}

//go:embed query_filter_test.gotpl
var queryFilterTestTemplate string

func queryFilterTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFilterTestTemplate, nil)
}

//go:embed users.gotpl
var usersTemplate string

func usersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTemplate, nil)
}

//go:embed users_test.gotpl
var usersTestTemplate string

func usersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersTestTemplate, nil)
}

//go:embed webhooks.gotpl
var webhooksTemplate string

func webhooksDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTemplate, nil)
}

//go:embed webhooks_test.gotpl
var webhooksTestTemplate string

func webhooksTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTestTemplate, nil)
}
