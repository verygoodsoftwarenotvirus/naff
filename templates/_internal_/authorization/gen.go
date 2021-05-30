package authorization

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/authorization"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"service_role_test.go":  serviceRoleTestDotGo(proj),
		"account_role.go":       accountRoleDotGo(proj),
		"account_role_test.go":  accountRoleTestDotGo(proj),
		"authorization.go":      authorizationDotGo(proj),
		"authorization_test.go": authorizationTestDotGo(proj),
		"permissions.go":        permissionsDotGo(proj),
		"rbac.go":               rbacDotGo(proj),
		"service_role.go":       serviceRoleDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed service_role_test.gotpl
var serviceRoleTestTemplate string

func serviceRoleTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serviceRoleTestTemplate, nil)
}

//go:embed account_role.gotpl
var accountRoleTemplate string

func accountRoleDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountRoleTemplate, nil)
}

//go:embed account_role_test.gotpl
var accountRoleTestTemplate string

func accountRoleTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, accountRoleTestTemplate, nil)
}

//go:embed authorization.gotpl
var authorizationTemplate string

func authorizationDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authorizationTemplate, nil)
}

//go:embed authorization_test.gotpl
var authorizationTestTemplate string

func authorizationTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, authorizationTestTemplate, nil)
}

//go:embed permissions.gotpl
var permissionsTemplate string

func permissionsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, permissionsTemplate, nil)
}

//go:embed rbac.gotpl
var rbacTemplate string

func rbacDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, rbacTemplate, nil)
}

//go:embed service_role.gotpl
var serviceRoleTemplate string

func serviceRoleDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serviceRoleTemplate, nil)
}
