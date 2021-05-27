package querybuilding

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "querybuilding"

	basePackagePath = "internal/database/querybuilding"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"external_id_generator.go":      externalIDGeneratorDotGo(proj),
		"mock_external_id_generator.go": mockExternalIDGeneratorDotGo(proj),
		"query_builders.go":             queryBuildersDotGo(proj),
		"query_constants.go":            queryConstantsDotGo(proj),
		"query_filter_test.go":          queryFilterTestDotGo(proj),
		"query_filters.go":              queryFiltersDotGo(proj),
		"column_lists.go":               columnListsDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
