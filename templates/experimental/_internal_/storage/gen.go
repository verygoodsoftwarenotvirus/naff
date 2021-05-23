package storage

import (
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
		"bucket_azure.go":            bucketAzureDotGo(proj),
		"bucket_azure_test_.go":      bucketAzureTestDotGo(proj),
		"bucket_s3_test_.go":         bucketS3TestDotGo(proj),
		"files.go":                   filesDotGo(proj),
		"files_test_.go":             filesTestDotGo(proj),
		"bucket_filesystem.go":       bucketFilesystemDotGo(proj),
		"bucket_filesystem_test_.go": bucketFilesystemTestDotGo(proj),
		"bucket_gcs.go":              bucketGcsDotGo(proj),
		"bucket_gcs_test_.go":        bucketGcsTestDotGo(proj),
		"bucket_s3.go":               bucketS3DotGo(proj),
		"uploader.go":                uploaderDotGo(proj),
		"uploader_test_.go":          uploaderTestDotGo(proj),
		"wire.go":                    wireDotGo(proj),
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
