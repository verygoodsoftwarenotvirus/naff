package querybuilding

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/database/querybuilding"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"external_id_generator.go":      externalIDGeneratorDotGo(proj),
		"mock_external_id_generator.go": mockExternalIDGeneratorDotGo(proj),
		"query_builders.go":             queryBuildersDotGo(proj),
		"query_constants.go":            queryConstantsDotGo(proj),
		"query_filter_test.go":          queryFilterTestDotGo(proj),
		"query_filters.go":              queryFiltersDotGo(proj),
		"column_lists.go":               columnListsDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed external_id_generator.gotpl
var externalIDGeneratorTemplate string

func externalIDGeneratorDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, externalIDGeneratorTemplate, nil)
}

//go:embed mock_external_id_generator.gotpl
var mockExternalIDGeneratorTemplate string

func mockExternalIDGeneratorDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockExternalIDGeneratorTemplate, nil)
}

//go:embed query_builders.gotpl
var queryBuildersTemplate string

func queryBuildersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryBuildersTemplate, nil)
}

//go:embed query_constants.gotpl
var queryConstantsTemplate string

func queryConstantsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryConstantsTemplate, nil)
}

//go:embed query_filter_test.gotpl
var queryFilterTestTemplate string

func queryFilterTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFilterTestTemplate, nil)
}

//go:embed query_filters.gotpl
var queryFiltersTemplate string

func queryFiltersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFiltersTemplate, nil)
}

//go:embed column_lists.gotpl
var columnListsTemplate string

func columnListsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, columnListsTemplate, nil)
}
