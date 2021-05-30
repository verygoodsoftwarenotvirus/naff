package fakes

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "fakes"

	basePackagePath = "pkg/types/fakes"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"fake.go":                    fakeDotGo(proj),
		"misc.go":                    miscDotGo(proj),
		"query_filter.go":            queryFilterDotGo(proj),
		"account.go":                 accountDotGo(proj),
		"audit_log_entry.go":         auditLogEntryDotGo(proj),
		"auth.go":                    authDotGo(proj),
		"doc.go":                     docDotGo(proj),
		"item.go":                    itemDotGo(proj),
		"user.go":                    userDotGo(proj),
		"webhook.go":                 webhookDotGo(proj),
		"account_user_membership.go": accountUserMembershipDotGo(proj),
		"admin.go":                   adminDotGo(proj),
		"delegated_client.go":        delegatedClientDotGo(proj),
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
