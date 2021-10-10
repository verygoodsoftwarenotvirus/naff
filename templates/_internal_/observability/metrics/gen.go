package metrics

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/observability/metrics"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"config_test.go":  configTestDotGo(proj),
		"counter.go":      counterDotGo(proj),
		"counter_test.go": counterTestDotGo(proj),
		"doc.go":          docDotGo(proj),
		"types.go":        typesDotGo(proj),
		"types_test.go":   typesTestDotGo(proj),
		"wire.go":         wireDotGo(proj),
		"config.go":       configDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed counter.gotpl
var counterTemplate string

func counterDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, counterTemplate, nil)
}

//go:embed counter_test.gotpl
var counterTestTemplate string

func counterTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, counterTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed types.gotpl
var typesTemplate string

func typesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, typesTemplate, nil)
}

//go:embed types_test.gotpl
var typesTestTemplate string

func typesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, typesTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}
