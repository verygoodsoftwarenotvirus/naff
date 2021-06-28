package server

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "server"

	basePackagePath = "internal/server"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"wire.go":             wireDotGo(proj),
		"config.go":           configDotGo(proj),
		"config_test.go":      configTestDotGo(proj),
		"doc.go":              docDotGo(proj),
		"http_server.go":      httpServerDotGo(proj),
		"http_server_test.go": httpServerTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	jenFiles := map[string]*jen.File{
		"http_routes.go": httpRoutesDotGo(proj),
	}

	for path, file := range jenFiles {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed http_server.gotpl
var httpServerTemplate string

func httpServerDotGo(proj *models.Project) string {
	typeServiceDeclarationFields := []string{}
	typeServiceParams := []string{}
	typeServiceConstructorFields := []string{}

	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()
		puvn := typ.Name.PluralUnexportedVarName()

		typeServiceDeclarationFields = append(typeServiceDeclarationFields, fmt.Sprintf("%sService      types.%sDataService", puvn, sn))
		typeServiceParams = append(typeServiceParams, fmt.Sprintf("%sService types.%sDataService,", puvn, sn))
		typeServiceConstructorFields = append(typeServiceConstructorFields, fmt.Sprintf("%sService:      %sService,", puvn, puvn))
	}

	generated := map[string]string{
		"typeServiceDeclarationFields": strings.Join(typeServiceDeclarationFields, "\n\t"),
		"typeServiceParams":            strings.Join(typeServiceParams, "\n\t"),
		"typeServiceConstructorFields": strings.Join(typeServiceConstructorFields, "\n\t"),
	}

	return models.RenderCodeFile(proj, httpServerTemplate, generated)
}

//go:embed http_server_test.gotpl
var httpServerTestTemplate string

func httpServerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpServerTestTemplate, nil)
}
