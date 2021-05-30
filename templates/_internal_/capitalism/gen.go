package capitalism

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/capitalism"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"config.go":               configDotGo(proj),
		"config_test.go":          configTestDotGo(proj),
		"mock.go":                 mockDotGo(proj),
		"noop_payment_manager.go": noopPaymentManagerDotGo(proj),
		"payment_manager.go":      paymentManagerDotGo(proj),
		"wire.go":                 wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate)
}

//go:embed mock.gotpl
var mockTemplate string

func mockDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockTemplate)
}

//go:embed noop_payment_manager.gotpl
var noopPaymentManagerTemplate string

func noopPaymentManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, noopPaymentManagerTemplate)
}

//go:embed payment_manager.gotpl
var paymentManagerTemplate string

func paymentManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, paymentManagerTemplate)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate)
}
