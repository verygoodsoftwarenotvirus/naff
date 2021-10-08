package elasticsearch

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/search/elasticsearch"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"doc.go":  docDotGo(proj),
		"elasticsearch.go":  elasticsearchDotGo(proj),
		"elasticsearch_test.go":  elasticsearchTestDotGo(proj),
		"wire.go":  wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
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


//go:embed elasticsearch.gotpl
var elasticsearch string
func elasticsearchDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, elasticsearch, nil)
}

//go:embed elasticsearch_test.gotpl
var elasticsearchTest string
func elasticsearchTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, elasticsearchTest, nil)
}

//go:embed wire.gotpl
var wire string
func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wire, nil)
}

