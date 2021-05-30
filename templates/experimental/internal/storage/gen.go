package storage

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "storage"

	basePackagePath = "internal/storage"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"bucket_filesystem_test.go": bucketFilesystemTestDotGo(proj),
		"bucket_gcs.go":             bucketGcsDotGo(proj),
		"bucket_s3.go":              bucketS3DotGo(proj),
		"bucket_s3_test.go":         bucketS3TestDotGo(proj),
		"files.go":                  filesDotGo(proj),
		"uploader.go":               uploaderDotGo(proj),
		"bucket_azure_test.go":      bucketAzureTestDotGo(proj),
		"bucket_filesystem.go":      bucketFilesystemDotGo(proj),
		"bucket_gcs_test.go":        bucketGcsTestDotGo(proj),
		"files_test.go":             filesTestDotGo(proj),
		"uploader_test.go":          uploaderTestDotGo(proj),
		"wire.go":                   wireDotGo(proj),
		"bucket_azure.go":           bucketAzureDotGo(proj),
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
