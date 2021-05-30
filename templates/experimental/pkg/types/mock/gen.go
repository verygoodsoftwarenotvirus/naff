package mock

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "mock"

	basePackagePath = "pkg/types/mock"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"item_data_manager.go":                    itemDataManagerDotGo(proj),
		"account_user_membership_data_manager.go": accountUserMembershipDataManagerDotGo(proj),
		"admin_user_data_manager.go":              adminUserDataManagerDotGo(proj),
		"auth_audit_manager.go":                   authAuditManagerDotGo(proj),
		"auth_service.go":                         authServiceDotGo(proj),
		"user_data_manager.go":                    userDataManagerDotGo(proj),
		"users_service.go":                        usersServiceDotGo(proj),
		"webhook_data_manager.go":                 webhookDataManagerDotGo(proj),
		"account_data_manager.go":                 accountDataManagerDotGo(proj),
		"api_client_data_manager.go":              apiClientDataManagerDotGo(proj),
		"audit_log_entry_data_manager.go":         auditLogEntryDataManagerDotGo(proj),
		"doc.go":                                  docDotGo(proj),
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
