package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"path/filepath"
)

const (
	packageName = "build"

	basePackagePrefix = "internal/build/server"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"build.go": buildDotGo(proj),
		"doc.go":   docDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePrefix, path), file); err != nil {
			return err
		}
	}

	return nil
}
