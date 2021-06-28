package types

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "types"

	basePackagePath = "pkg/types"
)

func jsonTag(val string) map[string]string {
	if val == "" {
		val = "-"
	}
	return map[string]string{"json": val}
}

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"webhook.go":                      webhookDotGo(proj),
		"api_client_test.go":              apiClientTestDotGo(proj),
		"user_test.go":                    userTestDotGo(proj),
		"doc.go":                          docDotGo(proj),
		"main.go":                         mainDotGo(proj),
		"query_filter_test.go":            queryFilterTestDotGo(proj),
		"validators.go":                   validatorsDotGo(proj),
		"validators_test.go":              validatorsTestDotGo(proj),
		"webhook_test.go":                 webhookTestDotGo(proj),
		"account_user_membership_test.go": accountUserMembershipTestDotGo(proj),
		"audit_log_entry.go":              auditLogEntryDotGo(proj),
		"admin.go":                        adminDotGo(proj),
		"admin_test.go":                   adminTestDotGo(proj),
		"api_client.go":                   apiClientDotGo(proj),
		"audit_log_entry_test.go":         auditLogEntryTestDotGo(proj),
		"auth.go":                         authDotGo(proj),
		"auth_test.go":                    authTestDotGo(proj),
		"account.go":                      accountDotGo(proj),
		"account_test.go":                 accountTestDotGo(proj),
		"main_test.go":                    mainTestDotGo(proj),
		"query_filter.go":                 queryFilterDotGo(proj),
		"user.go":                         userDotGo(proj),
		"account_user_membership.go":      accountUserMembershipDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{}
	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s.go", typ.Name.RouteName())] = iterableDotGo(proj, typ)
		jenFiles[fmt.Sprintf("%s_test.go", typ.Name.RouteName())] = iterableTestDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed webhook.gotpl
var webhookTemplate string

func webhookDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhookTemplate, nil)
}

//go:embed api_client_test.gotpl
var apiClientTestTemplate string

func apiClientTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientTestTemplate, nil)
}

//go:embed user_test.gotpl
var userTestTemplate string

func userTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, userTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed main.gotpl
var mainTemplate string

func mainDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mainTemplate, nil)
}

//go:embed query_filter_test.gotpl
var queryFilterTestTemplate string

func queryFilterTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFilterTestTemplate, nil)
}

//go:embed validators.gotpl
var validatorsTemplate string

func validatorsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, validatorsTemplate, nil)
}

//go:embed validators_test.gotpl
var validatorsTestTemplate string

func validatorsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, validatorsTestTemplate, nil)
}

//go:embed webhook_test.gotpl
var webhookTestTemplate string

func webhookTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhookTestTemplate, nil)
}

//go:embed account_user_membership_test.gotpl
var accountUserMembershipTestTemplate string

func accountUserMembershipTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipTestTemplate, nil)
}

//go:embed audit_log_entry.gotpl
var auditLogEntryTemplate string

func auditLogEntryDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, auditLogEntryTemplate, nil)
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

//go:embed api_client.gotpl
var apiClientTemplate string

func apiClientDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientTemplate, nil)
}

//go:embed audit_log_entry_test.gotpl
var auditLogEntryTestTemplate string

func auditLogEntryTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, auditLogEntryTestTemplate, nil)
}

//go:embed auth.gotpl
var authTemplate string

func authDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTemplate, nil)
}

//go:embed auth_test.gotpl
var authTestTemplate string

func authTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authTestTemplate, nil)
}

//go:embed account.gotpl
var accountTemplate string

func accountDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountTemplate, nil)
}

//go:embed account_test.gotpl
var accountTestTemplate string

func accountTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountTestTemplate, nil)
}

//go:embed main_test.gotpl
var mainTestTemplate string

func mainTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mainTestTemplate, nil)
}

//go:embed query_filter.gotpl
var queryFilterTemplate string

func queryFilterDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFilterTemplate, nil)
}

//go:embed user.gotpl
var userTemplate string

func userDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, userTemplate, nil)
}

//go:embed account_user_membership.gotpl
var accountUserMembershipTemplate string

func accountUserMembershipDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipTemplate, nil)
}
