package viper

import (
	"bytes"
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/config/viper"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"config.go":      configDotGo(proj),
		"config_test.go": configTestDotGo(proj),
		"doc.go":         docDotGo(proj),
		"keys.go":        keysDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	serviceConfigCodes := []jen.Code{}
	for _, typ := range proj.DataTypes {
		serviceConfigCodes = append(serviceConfigCodes,
			jen.ID("cfg").Dot("Set").Call(
				jen.IDf("ConfigKey%sLogging", typ.Name.Plural()),
				jen.ID("input").Dot("Services").Dot(typ.Name.Plural()).Dot("Logging"),
			),
			jen.Newline(),
		)
		if typ.SearchEnabled {
			serviceConfigCodes = append(serviceConfigCodes,
				jen.ID("cfg").Dot("Set").Call(
					jen.IDf("ConfigKey%sSearchIndexPath", typ.Name.Plural()),
					jen.ID("input").Dot("Services").Dot(typ.Name.Plural()).Dot("SearchIndexPath"),
				),
				jen.Newline(),
			)
		}
	}

	var b bytes.Buffer
	if err := jen.Null().Add(serviceConfigCodes...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"serviceConfigs": b.String(),
	}

	return models.RenderCodeFile(proj, configTemplate, generated)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	configFields := []string{}
	imports := []string{}

	for _, typ := range proj.DataTypes {
		imports = append(imports, fmt.Sprintf("\t%q", proj.ServicePackage(typ.Name.PackageName())))
		var b bytes.Buffer
		if err := jen.Null().Add(
			jen.Newline(),
			jen.ID(typ.Name.Plural()).MapAssign().Qual(proj.ServicePackage(typ.Name.PackageName()), "Config").Valuesln(
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.ID("SearchIndexPath").MapAssign().Litf("/%s_index_path", typ.Name.PluralRouteName())
					}
					return jen.Null()
				}(),
				jen.ID("Logging").MapAssign().Qual(proj.InternalLoggingPackage(), "Config").Valuesln(
					jen.ID("Name").MapAssign().Lit(typ.Name.PluralRouteName()),
					jen.ID("Level").MapAssign().Qual(proj.InternalLoggingPackage(), "InfoLevel"),
					jen.ID("Provider").MapAssign().Qual(proj.InternalLoggingPackage(), "ProviderZerolog"),
				),
			),
			jen.Newline(),
		).RenderWithoutFormatting(&b); err != nil {
			panic(err)
		}

		configFields = append(configFields, fmt.Sprintf("%s,\n", strings.TrimSpace(b.String())))
	}

	generated := map[string]string{
		"typeImports":      strings.Join(imports, "\n"),
		"testConfigFields": strings.Join(configFields, ""),
	}

	return models.RenderCodeFile(proj, configTestTemplate, generated)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed keys.gotpl
var keysTemplate string

func keysDotGo(proj *models.Project) string {
	keyDeclarations := []jen.Code{}
	for _, typ := range proj.DataTypes {
		keyDeclarations = append(keyDeclarations,
			jen.Newline(),
			jen.Newline(),
			jen.IDf("%sKey", typ.Name.PluralUnexportedVarName()).Equals().ID("servicesKey").Plus().ID("x").Plus().Lit(typ.Name.PluralRouteName()),
			jen.Newline(),
		)

		keyDeclarations = append(keyDeclarations,
			jen.Newline(),
			jen.Commentf("ConfigKey%sLogging controls logging for the %s service.", typ.Name.Plural(), typ.Name.Plural()),
			jen.Newline(),
			jen.IDf("ConfigKey%sLogging", typ.Name.Plural()).Equals().IDf("%sKey", typ.Name.PluralUnexportedVarName()).Plus().ID("x").Plus().ID("loggingKey"),
		)

		if typ.SearchEnabled {
			keyDeclarations = append(keyDeclarations,
				jen.Newline(),
				jen.Commentf("ConfigKey%sSearchIndexPath is the key viper will use to refer to the SearchSettings.%sSearchIndexPath setting.", typ.Name.Plural(), typ.Name.Plural()),
				jen.Newline(),
				jen.IDf("ConfigKey%sSearchIndexPath", typ.Name.Plural()).Equals().IDf("%sKey", typ.Name.PluralUnexportedVarName()).Plus().ID("x").Plus().Lit("search_index_path"),
			)
		}
	}

	var b bytes.Buffer
	if err := jen.Null().Add(keyDeclarations...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"keys": b.String(),
	}

	return models.RenderCodeFile(proj, keysTemplate, generated)
}
