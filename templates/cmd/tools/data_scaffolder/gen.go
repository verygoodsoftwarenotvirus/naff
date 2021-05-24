package data_scaffolder

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "data_scaffolder"

	basePackagePath = "cmd/tools/data_scaffolder"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"exiter.go":      exiterDotGo(proj),
		"exiter_test.go": exiterTestDotGo(proj),
		"main.go":        mainDotGo(proj),
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
