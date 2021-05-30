package types

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "types"

	basePackagePath = "pkg/types"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"audit_log_entry_test.go":         auditLogEntryTestDotGo(proj),
		"auth_test.go":                    authTestDotGo(proj),
		"item_test.go":                    itemTestDotGo(proj),
		"query_filter_test.go":            queryFilterTestDotGo(proj),
		"validators.go":                   validatorsDotGo(proj),
		"validators_test.go":              validatorsTestDotGo(proj),
		"webhook.go":                      webhookDotGo(proj),
		"account_test.go":                 accountTestDotGo(proj),
		"account_user_membership_test.go": accountUserMembershipTestDotGo(proj),
		"admin_test.go":                   adminTestDotGo(proj),
		"item.go":                         itemDotGo(proj),
		"main_test.go":                    mainTestDotGo(proj),
		"user.go":                         userDotGo(proj),
		"webhook_test.go":                 webhookTestDotGo(proj),
		"account_user_membership.go":      accountUserMembershipDotGo(proj),
		"audit_log_entry.go":              auditLogEntryDotGo(proj),
		"auth.go":                         authDotGo(proj),
		"doc.go":                          docDotGo(proj),
		"api_client.go":                   apiClientDotGo(proj),
		"admin.go":                        adminDotGo(proj),
		"api_client_test.go":              apiClientTestDotGo(proj),
		"main.go":                         mainDotGo(proj),
		"query_filter.go":                 queryFilterDotGo(proj),
		"user_test.go":                    userTestDotGo(proj),
		"account.go":                      accountDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
