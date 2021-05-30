package images

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "images"

	basePackagePath = "internal/uploads/images"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"thumbnails.go":      thumbnailsDotGo(proj),
		"thumbnails_test.go": thumbnailsTestDotGo(proj),
		"wire.go":            wireDotGo(proj),
		"doc.go":             docDotGo(proj),
		"images.go":          imagesDotGo(proj),
		"images_test.go":     imagesTestDotGo(proj),
		"mock.go":            mockDotGo(proj),
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
