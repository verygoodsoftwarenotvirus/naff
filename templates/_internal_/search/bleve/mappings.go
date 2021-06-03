package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mappingsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	for _, typ := range proj.DataTypes {
		code.Add(
			jen.Func().IDf("build%sMapping", typ.Name.Singular()).Params().Params(jen.Op("*").Qual("github.com/blevesearch/bleve/v2/mapping", "IndexMappingImpl")).Body(
				jen.ID("m").Op(":=").Qual("github.com/blevesearch/bleve/v2/mapping", "NewIndexMapping").Call(),
				jen.Newline(),
				jen.ID("englishTextFieldMapping").Op(":=").Qual(constants.SearchLibrary, "NewTextFieldMapping").Call(),
				jen.ID("englishTextFieldMapping").Dot("Analyzer").Op("=").Qual("github.com/blevesearch/bleve/v2/analysis/lang/en", "AnalyzerName"),
				jen.Newline(),
				jen.IDf("%sMapping", typ.Name.UnexportedVarName()).Op(":=").Qual(constants.SearchLibrary, "NewDocumentMapping").Call(),
				jen.IDf("%sMapping", typ.Name.UnexportedVarName()).Dot("AddFieldMappingsAt").Call(
					jen.Lit("name"),
					jen.ID("englishTextFieldMapping"),
				),
				jen.IDf("%sMapping", typ.Name.UnexportedVarName()).Dot("AddFieldMappingsAt").Call(
					jen.Lit("details"),
					jen.ID("englishTextFieldMapping"),
				),
				jen.IDf("%sMapping", typ.Name.UnexportedVarName()).Dot("AddFieldMappingsAt").Call(
					jen.Lit("belongsToAccount"),
					jen.Qual(constants.SearchLibrary, "NewNumericFieldMapping").Call(),
				),
				jen.ID("m").Dot("AddDocumentMapping").Call(
					jen.Litf(typ.Name.RouteName()),
					jen.IDf("%sMapping", typ.Name.UnexportedVarName()),
				),
				jen.Newline(),
				jen.Return().ID("m"),
			),
			jen.Newline(),
		)
	}

	return code
}
