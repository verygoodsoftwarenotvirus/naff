package authorization

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

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
	files := map[string]string{
		"service_role_test.go": serviceRoleTestDotGo(proj),
		"permissions.go":       permissionsDotGo(proj),
		"rbac.go":              rbacDotGo(proj),
		"service_role.go":      serviceRoleDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{
		"account_role.go":       accountRoleDotGo(proj),
		"account_role_test.go":  accountRoleTestDotGo(proj),
		"authorization.go":      authorizationDotGo(proj),
		"authorization_test.go": authorizationTestDotGo(proj),
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
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

//go:embed permissions.gotpl
var permissionsTemplate string

func permissionsDotGo(proj *models.Project) string {
	adminTypePermissions := []string{}
	memberTypePermissions := []string{}
	accountAdminPermissionsSetDecl := []string{}
	accountMemberPermissionsSetDecl := []string{}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		prn := typ.Name.PluralRouteName()

		memberTypePermissions = append(memberTypePermissions, fmt.Sprintf(`
// Create%sPermission is an account user permission.
Create%sPermission Permission = "create.%s"
// Read%sPermission is an account user permission.
Read%sPermission Permission = "read.%s"
// Search%sPermission is an account user permission.
Search%sPermission Permission = "search.%s"
// Update%sPermission is an account user permission.
Update%sPermission Permission = "update.%s"
// Archive%sPermission is an account user permission.
Archive%sPermission Permission = "archive.%s"
`, pn, pn, prn, pn, pn, prn, pn, pn, prn, pn, pn, prn, pn, pn, prn))
		accountMemberPermissionsSetDecl = append(accountMemberPermissionsSetDecl, fmt.Sprintf(`Create%sPermission.ID():  Create%sPermission,
Read%sPermission.ID():    Read%sPermission,
Search%sPermission.ID():  Search%sPermission,
Update%sPermission.ID():  Update%sPermission,
Archive%sPermission.ID(): Archive%sPermission,
`, pn, pn, pn, pn, pn, pn, pn, pn, pn, pn))
	}

	generated := map[string]string{
		"adminTypePermissions":            strings.Join(adminTypePermissions, "\n"),
		"memberTypePermissions":           strings.Join(memberTypePermissions, "\n"),
		"accountAdminPermissionsSetDecl":  strings.Join(accountAdminPermissionsSetDecl, "\n"),
		"accountMemberPermissionsSetDecl": strings.Join(accountMemberPermissionsSetDecl, "\n"),
	}

	return models.RenderCodeFile(proj, permissionsTemplate, generated)
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
