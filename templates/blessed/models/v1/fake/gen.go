package fake

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "fakemodels"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"models/v1/fake/doc.go":           docDotGo(proj),
		"models/v1/fake/fake.go":          fakeDotGo(proj),
		"models/v1/fake/oauth2_client.go": oauth2ClientDotGo(proj),
		"models/v1/fake/query_filter.go":  queryFilterDotGo(proj),
		"models/v1/fake/user.go":          userDotGo(proj),
		"models/v1/fake/webhook.go":       webhookDotGo(proj),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("models/v1/fake/%s.go", typ.Name.RouteName())] = iterablesDotGo(proj, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}
