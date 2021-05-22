package fake

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "fakemodels"

	basePackagePath = "pkg/types/fakes"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":           docDotGo(proj),
		"fake.go":          fakeDotGo(proj),
		"oauth2_client.go": oauth2ClientDotGo(proj),
		"query_filter.go":  queryFilterDotGo(proj),
		"user.go":          userDotGo(proj),
		"webhook.go":       webhookDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s.go", typ.Name.RouteName())] = iterablesDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
