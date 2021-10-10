package config

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "config"

	basePackagePath = "internal/config"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"doc.go":       docDotGo(proj),
		"meta.go":      metaDotGo(proj),
		"meta_test.go": metaTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{
		"config.go":      configDotGo(proj),
		"config_test.go": configTestDotGo(proj),
		"wire.go":        wireDotGo(proj),
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed meta.gotpl
var metaTemplate string

func metaDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, metaTemplate, nil)
}

//go:embed meta_test.gotpl
var metaTestTemplate string

func metaTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, metaTestTemplate, nil)
}
