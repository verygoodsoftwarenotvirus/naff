package integration

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "integration"

	basePackagePath = "tests/integration"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"init.go":          initDotGo(proj),
		"meta_test.go":     metaTestDotGo(proj),
		"oauth2_test.go":   oauth2TestDotGo(proj),
		"users_test.go":    usersTestDotGo(proj),
		"webhooks_test.go": webhooksTestDotGo(proj),
		"auth_test.go":     authTestDotGo(proj),
		"doc.go":           docDotGo(),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
