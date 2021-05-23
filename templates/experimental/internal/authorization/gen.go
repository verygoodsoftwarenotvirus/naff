package authorization

import (
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
		"internal/authorization/authorization_test.go": authorizationTestDotGo(),
		"internal/authorization/permissions.go":        permissionsDotGo(),
		"internal/authorization/rbac.go":               rbacDotGo(),
		"internal/authorization/service_role.go":       serviceRoleDotGo(),
		"internal/authorization/service_role_test.go":  serviceRoleTestDotGo(),
		"internal/authorization/account_role.go":       accountRoleDotGo(),
		"internal/authorization/account_role_test.go":  accountRoleTestDotGo(),
		"internal/authorization/authorization.go":      authorizationDotGo(),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
