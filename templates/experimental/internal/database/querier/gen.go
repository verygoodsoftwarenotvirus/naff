package querier

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "querier"

	basePackagePath = "internal/database/querier"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"admin.go":                          adminDotGo(proj),
		"doc.go":                            docDotGo(proj),
		"migrate_test_.go":                  migrateTestDotGo(proj),
		"querier_test_.go":                  querierTestDotGo(proj),
		"account_user_memberships_test_.go": accountUserMembershipsTestDotGo(proj),
		"users_test_.go":                    usersTestDotGo(proj),
		"errors.go":                         errorsDotGo(proj),
		"account_user_memberships.go":       accountUserMembershipsDotGo(proj),
		"admin_test_.go":                    adminTestDotGo(proj),
		"querier.go":                        querierDotGo(proj),
		"users.go":                          usersDotGo(proj),
		"migrate.go":                        migrateDotGo(proj),
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
