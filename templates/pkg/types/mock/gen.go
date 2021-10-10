package mock

import (
	_ "embed"
	"fmt"
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
	files := map[string]string{
		"account_data_manager.go":                 accountDataManagerDotGo(proj),
		"account_user_membership_data_manager.go": accountUserMembershipDataManagerDotGo(proj),
		"admin_user_data_manager.go":              adminUserDataManagerDotGo(proj),
		"api_client_data_manager.go":              apiClientDataManagerDotGo(proj),
		"users_service.go":                        usersServiceDotGo(proj),
		"webhook_data_manager.go":                 webhookDataManagerDotGo(proj),
		"auth_service.go":                         authServiceDotGo(proj),
		"user_data_manager.go":                    userDataManagerDotGo(proj),
		"doc.go":                                  docDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{}
	for _, typ := range proj.DataTypes {
		rn := typ.Name.RouteName()
		jenFiles[fmt.Sprintf("%s_data_manager.go", rn)] = mockIterableDataManagerDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed account_data_manager.gotpl
var accountDataManagerTemplate string

func accountDataManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountDataManagerTemplate, nil)
}

//go:embed account_user_membership_data_manager.gotpl
var accountUserMembershipDataManagerTemplate string

func accountUserMembershipDataManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountUserMembershipDataManagerTemplate, nil)
}

//go:embed admin_user_data_manager.gotpl
var adminUserDataManagerTemplate string

func adminUserDataManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, adminUserDataManagerTemplate, nil)
}

//go:embed api_client_data_manager.gotpl
var apiClientDataManagerTemplate string

func apiClientDataManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, apiClientDataManagerTemplate, nil)
}

//go:embed users_service.gotpl
var usersServiceTemplate string

func usersServiceDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, usersServiceTemplate, nil)
}

//go:embed webhook_data_manager.gotpl
var webhookDataManagerTemplate string

func webhookDataManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhookDataManagerTemplate, nil)
}

//go:embed auth_service.gotpl
var authServiceTemplate string

func authServiceDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authServiceTemplate, nil)
}

//go:embed user_data_manager.gotpl
var userDataManagerTemplate string

func userDataManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, userDataManagerTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}
