package querybuilding

import (
	"bytes"
	_ "embed"
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/database/querybuilding"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"external_id_generator.go":      externalIDGeneratorDotGo(proj),
		"mock_external_id_generator.go": mockExternalIDGeneratorDotGo(proj),
		"query_builders.go":             queryBuildersDotGo(proj),
		"query_constants.go":            queryConstantsDotGo(proj),
		"query_filter_test.go":          queryFilterTestDotGo(proj),
		"query_filters.go":              queryFiltersDotGo(proj),
		"column_lists.go":               columnListsDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed external_id_generator.gotpl
var externalIDGeneratorTemplate string

func externalIDGeneratorDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, externalIDGeneratorTemplate, nil)
}

//go:embed mock_external_id_generator.gotpl
var mockExternalIDGeneratorTemplate string

func mockExternalIDGeneratorDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockExternalIDGeneratorTemplate, nil)
}

//go:embed query_builders.gotpl
var queryBuildersTemplate string

func queryBuildersDotGo(proj *models.Project) string {
	querybuildingInterfaceDeclarations := []jen.Code{}
	querybuildingInterfaces := []string{}
	for _, typ := range proj.DataTypes {
		sn := typ.Name.Singular()
		pn := typ.Name.Plural()
		uvn := typ.Name.UnexportedVarName()

		querybuildingInterfaces = append(querybuildingInterfaces, fmt.Sprintf("%sSQLQueryBuilder", sn))
		querybuildingInterfaceDeclarations = append(querybuildingInterfaceDeclarations,
			jen.Newline(),
			jen.Newline(),
			jen.Commentf("%sSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.", sn),
			jen.Newline(),
			jen.IDf("%sSQLQueryBuilder", sn).Interface(
				jen.IDf("Build%sExistsQuery", sn).Params(typ.BuildDBQuerierExistenceQueryMethodParams(proj)...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildGet%sQuery", sn).Params(typ.BuildDBQuerierRetrievalMethodParams(proj)...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildGetAll%sCountQuery", pn).Params(constants.CtxParam()).String(),
				jen.IDf("BuildGetBatchOf%sQuery", pn).Params(constants.CtxParam(), jen.List(jen.ID("beginID"), jen.ID("endID")).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildGet%sQuery", pn).Params(typ.BuildDBQuerierListRetrievalQueryBuildingMethodParams(proj)...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildGet%sWithIDsQuery", pn).Params(typ.BuildGetListOfSomethingFromIDsParams(proj)...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildCreate%sQuery", sn).Params(typ.BuildDBQuerierCreationQueryBuildingMethodParams(proj, false)...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildUpdate%sQuery", sn).Params(typ.BuildDBQuerierUpdateQueryBuildingMethodParams(proj)...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildArchive%sQuery", sn).Params(typ.BuildDBQuerierArchiveQueryMethodParams()...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
				jen.IDf("BuildGetAuditLogEntriesFor%sQuery", sn).Params(constants.CtxParam(), jen.IDf("%sID", uvn).Uint64()).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()),
			),
			jen.Newline(),
		)
	}

	var b bytes.Buffer
	if err := jen.Null().Add(querybuildingInterfaceDeclarations...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"queryBuilderInterfaceDeclarations": b.String(),
		"querybuilderInterfaces":            strings.Join(querybuildingInterfaces, "\n"),
	}

	return models.RenderCodeFile(proj, queryBuildersTemplate, generated)
}

//go:embed query_constants.gotpl
var queryConstantsTemplate string

func queryConstantsDotGo(proj *models.Project) string {
	typeDefinitions := []jen.Code{}
	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		prn := typ.Name.PluralRouteName()
		pcn := typ.Name.PluralCommonName()

		typeDefinitions = append(typeDefinitions,
			jen.Newline(),
			jen.Newline(),
			jen.Comment(""),
			jen.Newline(),
			jen.Commentf("%s Table.", pn),
			jen.Newline(),
			jen.Comment(""),
			jen.Newline(),
			jen.Newline(),
			jen.Commentf("%sTableName is what the %s table calls itself.", pn, pcn),
			jen.Newline(),
			jen.IDf("%sTableName", pn).Equals().Lit(prn),
			jen.Newline(),
		)

		for _, field := range typ.Fields {
			fsn := field.Name.Singular()

			typeDefinitions = append(typeDefinitions,
				jen.Commentf("%sTable%sColumn is what the %s table calls the %s column.", pn, fsn, pcn, field.Name.RouteName()),
				jen.Newline(),
				jen.IDf("%sTable%sColumn", pn, field.Name.Singular()).Equals().Lit(field.Name.RouteName()),
				jen.Newline(),
			)
		}

		typeDefinitions = append(typeDefinitions,
			jen.Commentf("%sTableAccountOwnershipColumn is what the %s table calls the ownership column.", pn, pcn),
			jen.Newline(),
			jen.IDf("%sTableAccountOwnershipColumn", pn).Equals().ID("accountOwnershipColumn"),
			jen.Newline(),
		)
	}

	var b bytes.Buffer
	if err := jen.Null().Add(typeDefinitions...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"queryConstants": b.String(),
	}

	return models.RenderCodeFile(proj, queryConstantsTemplate, generated)
}

//go:embed query_filter_test.gotpl
var queryFilterTestTemplate string

func queryFilterTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFilterTestTemplate, nil)
}

//go:embed query_filters.gotpl
var queryFiltersTemplate string

func queryFiltersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, queryFiltersTemplate, nil)
}

//go:embed column_lists.gotpl
var columnListsTemplate string

func columnListsDotGo(proj *models.Project) string {
	typeDefinitions := []jen.Code{}
	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		pcn := typ.Name.PluralCommonName()

		typeDefinitions = append(typeDefinitions,
			jen.Newline(),
			jen.Newline(),
			jen.Comment(""),
			jen.Newline(),
			jen.Commentf("%s Table.", pn),
			jen.Newline(),
			jen.Comment(""),
			jen.Newline(),
			jen.Newline(),
			jen.Commentf("%sTableColumns are the columns for the %s table.", pn, pcn),
			jen.Newline(),
		)

		columns := []jen.Code{
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", pn), jen.ID("IDColumn")),
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", pn), jen.ID("ExternalIDColumn")),
		}

		for _, field := range typ.Fields {
			columns = append(columns,
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", pn), jen.IDf("%sTable%sColumn", pn, field.Name.Singular())),
			)
		}

		columns = append(columns,
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", pn), jen.ID("CreatedOnColumn")),
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", pn), jen.ID("LastUpdatedOnColumn")),
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", pn), jen.ID("ArchivedOnColumn")),
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", pn), jen.IDf("%sTableAccountOwnershipColumn", pn)),
		)

		typeDefinitions = append(typeDefinitions,
			jen.IDf("%sTableColumns", pn).Equals().Index().String().Valuesln(
				columns...,
			),
			jen.Newline(),
		)
	}

	var b bytes.Buffer
	if err := jen.Null().Add(typeDefinitions...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"columnLists": b.String(),
	}

	return models.RenderCodeFile(proj, columnListsTemplate, generated)
}
