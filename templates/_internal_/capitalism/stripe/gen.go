package stripe

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "stripe"

	basePackagePath = "internal/capitalism/stripe"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"mock_backend_test.go": mockBackendTestDotGo(proj),
		"stripe.go":            stripeDotGo(proj),
		"stripe_test.go":       stripeTestDotGo(proj),
		"wire.go":              wireDotGo(proj),
		"wire_test.go":         wireTestDotGo(proj),
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
