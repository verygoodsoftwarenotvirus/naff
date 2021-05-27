package types

import (
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
	files := map[string]*jen.File{
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

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s.go", typ.Name.RouteName())] = iterableDotGo(proj, typ)
		files[fmt.Sprintf("%s_test.go", typ.Name.RouteName())] = iterableTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
