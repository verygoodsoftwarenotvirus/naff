package storage

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/storage"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"bucket_azure.go":           bucketAzureDotGo(proj),
		"bucket_azure_test.go":      bucketAzureTestDotGo(proj),
		"bucket_s3_test.go":         bucketS3TestDotGo(proj),
		"files.go":                  filesDotGo(proj),
		"files_test.go":             filesTestDotGo(proj),
		"bucket_filesystem.go":      bucketFilesystemDotGo(proj),
		"bucket_filesystem_test.go": bucketFilesystemTestDotGo(proj),
		"bucket_gcs.go":             bucketGcsDotGo(proj),
		"bucket_gcs_test.go":        bucketGcsTestDotGo(proj),
		"bucket_s3.go":              bucketS3DotGo(proj),
		"uploader.go":               uploaderDotGo(proj),
		"uploader_test.go":          uploaderTestDotGo(proj),
		"wire.go":                   wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed bucket_azure.gotpl
var bucketAzureTemplate string

func bucketAzureDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketAzureTemplate, nil)
}

//go:embed bucket_azure_test.gotpl
var bucketAzureTestTemplate string

func bucketAzureTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketAzureTestTemplate, nil)
}

//go:embed bucket_s3_test.gotpl
var bucketS3TestTemplate string

func bucketS3TestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketS3TestTemplate, nil)
}

//go:embed files.gotpl
var filesTemplate string

func filesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, filesTemplate, nil)
}

//go:embed files_test.gotpl
var filesTestTemplate string

func filesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, filesTestTemplate, nil)
}

//go:embed bucket_filesystem.gotpl
var bucketFilesystemTemplate string

func bucketFilesystemDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketFilesystemTemplate, nil)
}

//go:embed bucket_filesystem_test.gotpl
var bucketFilesystemTestTemplate string

func bucketFilesystemTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketFilesystemTestTemplate, nil)
}

//go:embed bucket_gcs.gotpl
var bucketGcsTemplate string

func bucketGcsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketGcsTemplate, nil)
}

//go:embed bucket_gcs_test.gotpl
var bucketGcsTestTemplate string

func bucketGcsTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketGcsTestTemplate, nil)
}

//go:embed bucket_s3.gotpl
var bucketS3Template string

func bucketS3DotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, bucketS3Template, nil)
}

//go:embed uploader.gotpl
var uploaderTemplate string

func uploaderDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, uploaderTemplate, nil)
}

//go:embed uploader_test.gotpl
var uploaderTestTemplate string

func uploaderTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, uploaderTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}
