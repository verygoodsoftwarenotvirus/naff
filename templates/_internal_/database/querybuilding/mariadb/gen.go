package mariadb

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
	packageName = "mariadb"

	basePackagePath = "internal/database/querybuilding/mariadb"
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
		"mariadb.go":                       mariadbDotGo(proj),
		"mariadb_test.go":                  mariadbTestDotGo(proj),
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

//go:embed mariadb.gotpl
var mariadbTemplate string

func mariadbDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mariadbTemplate, nil)
}

//go:embed mariadb_test.gotpl
var mariadbTestTemplate string

func mariadbTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mariadbTestTemplate, nil)
}

func typeToMariaDBType(t string) string {
	typeMap := map[string]string{
		//"[]string": "LONGTEXT",
		"string":   "LONGTEXT",
		"*string":  "LONGTEXT",
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
		"uint":     "INTEGER UNSIGNED",
		"*uint":    "INTEGER UNSIGNED",
		"uint8":    "INTEGER UNSIGNED",
		"*uint8":   "INTEGER UNSIGNED",
		"uint16":   "INTEGER UNSIGNED",
		"*uint16":  "INTEGER UNSIGNED",
		"uint32":   "INTEGER UNSIGNED",
		"*uint32":  "INTEGER UNSIGNED",
		"uint64":   "INTEGER UNSIGNED",
		"*uint64":  "INTEGER UNSIGNED",
		"float32":  "DOUBLE PRECISION",
		"*float32": "DOUBLE PRECISION",
		"float64":  "DOUBLE PRECISION",
		"*float64": "DOUBLE PRECISION",
	}

	if x, ok := typeMap[strings.TrimSpace(t)]; ok {
		return x
	}

	panic(fmt.Sprintf("unknown type!: %q", t))
}

//go:embed migrations.gotpl
var migrationsTemplate string

func migrationsDotGo(proj *models.Project) string {
	typeMigrations := []jen.Code{}
	for i, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		prn := typ.Name.PluralRouteName()

		migrationLines := []jen.Code{
			jen.Litf("CREATE TABLE IF NOT EXISTS %s (", prn),
			jen.Lit("    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,"),
			jen.Lit("    `external_id` VARCHAR(36) NOT NULL,"),
		}
		for _, field := range typ.Fields {
			rn := field.Name.RouteName()

			query := fmt.Sprintf("    `%s` %s", rn, typeToMariaDBType(field.Type))
			if !field.Pointer {
				query += ` NOT NULL`
			}
			if field.DefaultValue != "" {
				query += fmt.Sprintf(` DEFAULT %s`, field.DefaultValue)
			}

			migrationLines = append(migrationLines, jen.Lit(fmt.Sprintf("%s,", query)))
		}
		migrationLines = append(migrationLines,
			jen.Lit("    `created_on` BIGINT UNSIGNED,"),
			jen.Lit("    `last_updated_on` BIGINT UNSIGNED DEFAULT NULL,"),
			jen.Lit("    `archived_on` BIGINT UNSIGNED DEFAULT NULL,"),
			jen.Lit("    `belongs_to_account` BIGINT UNSIGNED NOT NULL,"),
			jen.Lit("    PRIMARY KEY (`id`),"),
			jen.Lit("    FOREIGN KEY (`belongs_to_account`) REFERENCES accounts(`id`) ON DELETE CASCADE"),
			jen.Lit(");"),
		)

		typeMigrations = append(typeMigrations,
			jen.Valuesln(
				jen.ID("Version").MapAssign().Lit(0.15+float64(i)*.01),
				jen.ID("Description").MapAssign().Litf("create %s table", prn),
				jen.ID("Script").MapAssign().Qual("strings", "Join").Call(jen.Index().String().Valuesln(migrationLines...), jen.Lit("\n")),
			),
			jen.Valuesln(
				jen.ID("Version").MapAssign().Lit(0.16+float64(i)*.01),
				jen.ID("Description").MapAssign().Litf("create %s table creation trigger", prn),
				jen.ID("Script").MapAssign().ID("buildCreationTriggerScript").Call(jen.Qual(proj.QuerybuildersPackage(), fmt.Sprintf("%sTableName", pn))),
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
