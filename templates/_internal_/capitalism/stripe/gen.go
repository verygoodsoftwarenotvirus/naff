package stripe

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/capitalism/stripe"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"mock_backend_test.go": mockBackendTestDotGo(proj),
		"stripe.go":            stripeDotGo(proj),
		"stripe_test.go":       stripeTestDotGo(proj),
		"wire.go":              wireDotGo(proj),
		"wire_test.go":         wireTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed mock_backend_test.gotpl
var mockBackendTestTemplate string

func mockBackendTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockBackendTestTemplate)
}

//go:embed stripe.gotpl
var stripeTemplate string

func stripeDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, stripeTemplate)
}

//go:embed stripe_test.gotpl
var stripeTestTemplate string

func stripeTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, stripeTestTemplate)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate)
}

//go:embed wire_test.gotpl
var wireTestTemplate string

func wireTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTestTemplate)
}
