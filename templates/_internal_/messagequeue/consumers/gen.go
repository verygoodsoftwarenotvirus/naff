package consumers

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "consumers"

	basePackagePath = "internal/messagequeue/consumers"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"consumers.go":           consumersDotGo(proj),
		"redis_consumer.go":      redisConsumerDotGo(proj),
		"redis_consumer_test.go": redisConsumerTestDotGo(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed consumers.gotpl
var consumersTemplate string

func consumersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, consumersTemplate, nil)
}

//go:embed redis_consumer.gotpl
var redisConsumerTemplate string

func redisConsumerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, redisConsumerTemplate, nil)
}

//go:embed redis_consumer_test.gotpl
var redisConsumerTestTemplate string

func redisConsumerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, redisConsumerTestTemplate, nil)
}
