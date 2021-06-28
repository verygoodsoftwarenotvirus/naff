package secrets

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/secrets"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"config.go":       configDotGo(proj),
		"config_test.go":  configTestDotGo(proj),
		"secrets.go":      secretsDotGo(proj),
		"secrets_test.go": secretsTestDotGo(proj),
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
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed secrets.gotpl
var secretsTemplate string

func secretsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, secretsTemplate, nil)
}

//go:embed secrets_test.gotpl
var secretsTestTemplate string

func secretsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, secretsTestTemplate, nil)
}
