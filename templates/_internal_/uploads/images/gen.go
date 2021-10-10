package images

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/uploads/images"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"thumbnails.go":      thumbnailsDotGo(proj),
		"thumbnails_test.go": thumbnailsTestDotGo(proj),
		"wire.go":            wireDotGo(proj),
		"doc.go":             docDotGo(proj),
		"images.go":          imagesDotGo(proj),
		"images_test.go":     imagesTestDotGo(proj),
		"mock.go":            mockDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed thumbnails.gotpl
var thumbnailsTemplate string

func thumbnailsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, thumbnailsTemplate, nil)
}

//go:embed thumbnails_test.gotpl
var thumbnailsTestTemplate string

func thumbnailsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, thumbnailsTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed images.gotpl
var imagesTemplate string

func imagesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, imagesTemplate, nil)
}

//go:embed images_test.gotpl
var imagesTestTemplate string

func imagesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, imagesTestTemplate, nil)
}

//go:embed mock.gotpl
var mockTemplate string

func mockDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockTemplate, nil)
}
