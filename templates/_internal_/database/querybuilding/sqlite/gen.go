package sqlite

import (
	"bytes"
	_ "embed"
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "sqlite"

	basePackagePath = "internal/database/querybuilding/sqlite"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"accounts.go":                      accountsDotGo(proj),
		"accounts_test.go":                 accountsTestDotGo(proj),
		"account_user_memberships.go":      accountUserMembershipsDotGo(proj),
		"account_user_memberships_test.go": accountUserMembershipsTestDotGo(proj),
		"api_clients.go":                   apiClientsDotGo(proj),
		"api_clients_test.go":              apiClientsTestDotGo(proj),
		"audit_log_entries.go":             auditLogEntriesDotGo(proj),
		"audit_log_entries_test.go":        auditLogEntriesTestDotGo(proj),
		"doc.go":                           docDotGo(proj),
		"generic.go":                       genericDotGo(proj),
		"generic_test.go":                  genericTestDotGo(proj),
		"sqlite.go":                        sqliteDotGo(proj),
		"sqlite_test.go":                   sqliteTestDotGo(proj),
		"migrations.go":                    migrationsDotGo(proj),
		"users.go":                         usersDotGo(proj),
		"users_test.go":                    usersTestDotGo(proj),
		"webhooks.go":                      webhooksDotGo(proj),
		"webhooks_test.go":                 webhooksTestDotGo(proj),
		"wire.go":                          wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{}
	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
		jenFiles[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
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

//go:embed audit_log_entries.gotpl
var auditLogEntriesTemplate string

func auditLogEntriesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, auditLogEntriesTemplate, nil)
}

//go:embed audit_log_entries_test.gotpl
var auditLogEntriesTestTemplate string

func auditLogEntriesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, auditLogEntriesTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed generic.gotpl
var genericTemplate string

func genericDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, genericTemplate, nil)
}

//go:embed generic_test.gotpl
var genericTestTemplate string

func genericTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, genericTestTemplate, nil)
}

//go:embed sqlite.gotpl
var sqliteTemplate string

func sqliteDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, sqliteTemplate, nil)
}

//go:embed sqlite_test.gotpl
var sqliteTestTemplate string

func sqliteTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, sqliteTestTemplate, nil)
}

func typeToSqliteType(t string) string {
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
		"uint32":   "INTEGER",
		"*uint32":  "INTEGER",
		"uint64":   "INTEGER",
		"*uint64":  "INTEGER",
		"float32":  "REAL",
		"*float32": "REAL",
		"float64":  "REAL",
		"*float64": "REAL",
	}

	if x, ok := typeMap[t]; ok {
		return x
	}

	panic(fmt.Sprintf("unknown type!: %q", t))
}

//go:embed migrations.gotpl
var migrationsTemplate string

func migrationsDotGo(proj *models.Project) string {
	typeMigrations := []jen.Code{}
	for i, typ := range proj.DataTypes {
		prn := typ.Name.PluralRouteName()

		scriptParts := []string{
			fmt.Sprintf("\n			CREATE TABLE IF NOT EXISTS %s (", typ.Name.PluralRouteName()),
			`				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`,
			`				external_id TEXT NOT NULL,`,
		}

		for _, field := range typ.Fields {
			rn := field.Name.RouteName()

			query := fmt.Sprintf("				%s %s", rn, typeToSqliteType(field.Type))
			if !field.IsPointer {
				query += ` NOT NULL`
			}
			if field.DefaultValue != "" {
				query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
			}

			scriptParts = append(scriptParts, fmt.Sprintf("%s,", query))
		}

		scriptParts = append(scriptParts,
			`				created_on INTEGER NOT NULL DEFAULT (strftime('%s','now')),`,
			`				last_updated_on INTEGER DEFAULT NULL,`,
			`				archived_on INTEGER DEFAULT NULL,`,
		)

		if typ.BelongsToAccount {
			scriptParts = append(scriptParts,
				`				belongs_to_account INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE`,
			)
		}
		if typ.BelongsToStruct != nil {
			scriptParts = append(scriptParts,
				fmt.Sprintf(`				belongs_to_%s INTEGER NOT NULL,`, typ.BelongsToStruct.RouteName()),
				fmt.Sprintf(`				FOREIGN KEY(belongs_to_%s) REFERENCES %s(id)`, typ.BelongsToStruct.RouteName(), typ.BelongsToStruct.PluralRouteName()),
			)
		}

		scriptParts = append(scriptParts,
			`			);`,
		)

		typeMigrations = append(typeMigrations,
			jen.Valuesln(
				jen.ID("Version").MapAssign().Lit(0.08+float64(i)*.01),
				jen.ID("Description").MapAssign().Litf("create %s table", prn),
				jen.ID("Script").MapAssign().RawString(strings.Join(scriptParts, "\n")),
			),
		)
	}

	var b bytes.Buffer
	if err := jen.Listln(typeMigrations...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"typeMigrations": fmt.Sprintf("%s,", b.String()),
	}

	return models.RenderCodeFile(proj, migrationsTemplate, generated)
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

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}
