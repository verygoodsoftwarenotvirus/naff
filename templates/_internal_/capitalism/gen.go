package capitalism

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "capitalism"

	basePackagePath = "internal/capitalism"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"config.go":               configDotGo(proj),
		"config_test.go":          configTestDotGo(proj),
		"mock.go":                 mockDotGo(proj),
		"noop_payment_manager.go": noopPaymentManagerDotGo(proj),
		"payment_manager.go":      paymentManagerDotGo(proj),
		"wire.go":                 wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
