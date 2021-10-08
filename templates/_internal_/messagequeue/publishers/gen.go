package publishers

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "main"

	basePackagePath = "internal/messagequeue/publishers"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"publishers.go":           publishersDotGo(proj),
		"redis_publisher.go":      redisPublisherDotGo(proj),
		"redis_publisher_test.go": redisPublisherTestDotGo(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed publishers.gotpl
var publishersTemplate string

func publishersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, publishersTemplate, nil)
}

//go:embed redis_publisher.gotpl
var redisPublisherTemplate string

func redisPublisherDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, redisPublisherTemplate, nil)
}

//go:embed redis_publisher_test.gotpl
var redisPublisherTestTemplate string

func redisPublisherTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, redisPublisherTestTemplate, nil)
}
