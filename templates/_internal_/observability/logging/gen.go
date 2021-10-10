package logging

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/observability/logging"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"config.go":              configDotGo(proj),
		"config_test.go":         configTestDotGo(proj),
		"logging.go":             loggingDotGo(proj),
		"logging_test.go":        loggingTestDotGo(proj),
		"noop_logger.go":         noopLoggerDotGo(proj),
		"zerolog_logger.go":      zerologLoggerDotGo(proj),
		"zerolog_logger_test.go": zerologLoggerTestDotGo(proj),
	}

	for path, file := range files {
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

//go:embed logging.gotpl
var loggingTemplate string

func loggingDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, loggingTemplate, nil)
}

//go:embed logging_test.gotpl
var loggingTestTemplate string

func loggingTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, loggingTestTemplate, nil)
}

//go:embed noop_logger.gotpl
var noopLoggerTemplate string

func noopLoggerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, noopLoggerTemplate, nil)
}

//go:embed zerolog_logger.gotpl
var zerologLoggerTemplate string

func zerologLoggerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, zerologLoggerTemplate, nil)
}

//go:embed zerolog_logger_test.gotpl
var zerologLoggerTestTemplate string

func zerologLoggerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, zerologLoggerTestTemplate, nil)
}
