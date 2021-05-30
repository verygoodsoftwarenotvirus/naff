package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mappingsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildItemMapping").Params().Params(jen.Op("*").Qual("github.com/blevesearch/bleve/v2/mapping", "IndexMappingImpl")).Body(
			jen.ID("m").Op(":=").Qual("github.com/blevesearch/bleve/v2/mapping", "NewIndexMapping").Call(),
			jen.ID("englishTextFieldMapping").Op(":=").Qual("github.com/blevesearch/bleve/v2", "NewTextFieldMapping").Call(),
			jen.ID("englishTextFieldMapping").Dot("Analyzer").Op("=").Qual("github.com/blevesearch/bleve/v2/analysis/lang/en", "AnalyzerName"),
			jen.ID("itemMapping").Op(":=").Qual("github.com/blevesearch/bleve/v2", "NewDocumentMapping").Call(),
			jen.ID("itemMapping").Dot("AddFieldMappingsAt").Call(
				jen.Lit("name"),
				jen.ID("englishTextFieldMapping"),
			),
			jen.ID("itemMapping").Dot("AddFieldMappingsAt").Call(
				jen.Lit("details"),
				jen.ID("englishTextFieldMapping"),
			),
			jen.ID("itemMapping").Dot("AddFieldMappingsAt").Call(
				jen.Lit("belongsToAccount"),
				jen.Qual("github.com/blevesearch/bleve/v2", "NewNumericFieldMapping").Call(),
			),
			jen.ID("m").Dot("AddDocumentMapping").Call(
				jen.Lit("item"),
				jen.ID("itemMapping"),
			),
			jen.Return().ID("m"),
		),
		jen.Line(),
	)

	return code
}
