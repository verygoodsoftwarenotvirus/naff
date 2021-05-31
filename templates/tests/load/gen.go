package load

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "load"

	basePackagePath = "tests/load"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"actions.go":  actionsDotGo(proj),
		"init.go":     initDotGo(proj),
		"main.go":     mainDotGo(proj),
		"webhooks.go": webhooksDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{}
	for _, typ := range proj.DataTypes {
		jenFiles[fmt.Sprintf("%s.go", typ.Name.PluralRouteName())] = iterablesDotGo(proj, typ)
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed actions.gotpl
var actionsTemplate string

func actionsDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, actionsTemplate, nil)
}

//go:embed init.gotpl
var initTemplate string

func initDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, initTemplate, nil)
}

//go:embed main.gotpl
var mainTemplate string

func mainDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mainTemplate, nil)
}

//go:embed webhooks.gotpl
var webhooksTemplate string

func webhooksDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, webhooksTemplate, nil)
}
