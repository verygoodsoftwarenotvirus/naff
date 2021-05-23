package authorization

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "authorization"

	basePackagePath = "internal/authorization"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"service_role_test_.go":  serviceRoleTestDotGo(proj),
		"account_role.go":        accountRoleDotGo(proj),
		"account_role_test_.go":  accountRoleTestDotGo(proj),
		"authorization.go":       authorizationDotGo(proj),
		"authorization_test_.go": authorizationTestDotGo(proj),
		"permissions.go":         permissionsDotGo(proj),
		"rbac.go":                rbacDotGo(proj),
		"service_role.go":        serviceRoleDotGo(proj),
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
