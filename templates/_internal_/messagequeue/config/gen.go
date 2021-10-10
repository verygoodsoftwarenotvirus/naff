package config

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "config"

	basePackagePath = "internal/messagequeue/config"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"config.go":      configDotGo(proj),
		"config_test.go": configTestDotGo(proj),
		"wire.go":        wireDotGo(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
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

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}
