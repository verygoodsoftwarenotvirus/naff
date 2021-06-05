package database

import (
	"bytes"
	_ "embed"
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/database"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"database.go":      databaseDotGo(proj),
		"database_mock.go": databaseMockDotGo(proj),
		"doc.go":           docDotGo(proj),
		"wire.go":          wireDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed database.gotpl
var databaseTemplate string

func databaseDotGo(proj *models.Project) string {
	dataManagers := []string{}
	for _, typ := range proj.DataTypes {
		dataManagers = append(dataManagers, fmt.Sprintf("types.%sDataManager", typ.Name.Singular()))
	}

	generated := map[string]string{
		"dataManagers": strings.Join(dataManagers, "\n\t\t"),
	}

	return models.RenderCodeFile(proj, databaseTemplate, generated)
}

//go:embed database_mock.gotpl
var databaseMockTemplate string

func databaseMockDotGo(proj *models.Project) string {
	mockDatabaseFields := []string{}
	mockDataManagers := []string{}
	mockSQLBuilders := []string{}
	mockSQLBuilderEmbeds := []string{}
	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()

		mockDatabaseFields = append(mockDatabaseFields, fmt.Sprintf("%sDataManager:                  &mocktypes.%sDataManager{},", sn, sn))
		mockDataManagers = append(mockDataManagers, fmt.Sprintf("*mocktypes.%sDataManager", sn))
		mockSQLBuilders = append(mockSQLBuilders, fmt.Sprintf("%sSQLQueryBuilder: &mockquerybuilding.%sSQLQueryBuilder{},", sn, sn))
		mockSQLBuilderEmbeds = append(mockSQLBuilderEmbeds, fmt.Sprintf("*mockquerybuilding.%sSQLQueryBuilder", sn))
	}

	generated := map[string]string{
		"mockDatabaseFields":   strings.Join(mockDatabaseFields, "\n\t\t"),
		"mockDataManagers":     strings.Join(mockDataManagers, "\n\t\t"),
		"mockSQLBuilders":      strings.Join(mockSQLBuilders, "\n\t\t"),
		"mockSQLBuilderEmbeds": strings.Join(mockSQLBuilderEmbeds, "\n\t\t"),
	}

	return models.RenderCodeFile(proj, databaseMockTemplate, generated)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	dataManagerProviders := []jen.Code{}
	typeDataManagerProvidersArgs := []string{}

	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()

		typeDataManagerProvidersArgs = append(typeDataManagerProvidersArgs, fmt.Sprintf("Provide%sDataManager,", sn))
		dataManagerProviders = append(dataManagerProviders,
			jen.Newline(),
			jen.Commentf("Provide%sDataManager is an arbitrary function for dependency injection's sake.", sn),
			jen.Newline(),
			jen.Func().IDf("Provide%sDataManager", sn).Params(jen.ID("db").ID("DataManager")).Params(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn))).Body(
				jen.Return(jen.ID("db")),
			),
		)
	}

	var b bytes.Buffer
	if err := jen.Null().Add(dataManagerProviders...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"typeDataManagerProvidersArgs": strings.Join(typeDataManagerProvidersArgs, "\n\t\t"),
		"dataManagerProviders":         b.String(),
	}

	return models.RenderCodeFile(proj, wireTemplate, generated)
}
