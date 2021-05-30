package events

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/events"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"wire.go":            wireDotGo(proj),
		"config.go":          configDotGo(proj),
		"config_test.go":     configTestDotGo(proj),
		"publisher.go":       publisherDotGo(proj),
		"publisher_test.go":  publisherTestDotGo(proj),
		"subscriber.go":      subscriberDotGo(proj),
		"subscriber_test.go": subscriberTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
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

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed publisher.gotpl
var publisherTemplate string

func publisherDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, publisherTemplate, nil)
}

//go:embed publisher_test.gotpl
var publisherTestTemplate string

func publisherTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, publisherTestTemplate, nil)
}

//go:embed subscriber.gotpl
var subscriberTemplate string

func subscriberDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, subscriberTemplate, nil)
}

//go:embed subscriber_test.gotpl
var subscriberTestTemplate string

func subscriberTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, subscriberTestTemplate, nil)
}
