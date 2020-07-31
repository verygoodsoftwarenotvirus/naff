package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"strings"
)

func mappingsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			code.Add(buildSomethingMapper(typ)...)
		}
	}

	return code
}

func buildSomethingMapper(typ models.DataType) []jen.Code {
	uvn := typ.Name.UnexportedVarName()

	block := []jen.Code{
		jen.ID("m").Assign().Qual(
			"github.com/blevesearch/bleve/mapping",
			"NewIndexMapping",
		).Call(),
		jen.Line(),
		jen.ID("englishTextFieldMapping").Assign().Qual(searchPackage, "NewTextFieldMapping").Call(),
		jen.ID("englishTextFieldMapping").Dot("Analyzer").Equals().Qual(
			"github.com/blevesearch/bleve/analysis/lang/en",
			"AnalyzerName",
		),
		jen.Line(),
		jen.IDf("%sMapping", uvn).Assign().Qual(searchPackage, "NewDocumentMapping").Call(),
	}

	for _, field := range typ.Fields {
		if strings.ToLower(field.Type) == "string" {
			block = append(block,
				jen.IDf("%sMapping", uvn).Dot("AddFieldMappingsAt").Call(
					jen.Lit(field.Name.UnexportedVarName()),
					jen.ID("englishTextFieldMapping"),
				),
			)
		}
	}

	block = append(block,
		jen.IDf("%sMapping", uvn).Dot("AddFieldMappingsAt").Call(
			jen.Lit("belongsToUser"),
			jen.Qual(searchPackage, "NewNumericFieldMapping").Call(),
		),
		jen.ID("m").Dot("AddDocumentMapping").Call(
			jen.Lit(uvn),
			jen.IDf("%sMapping", uvn),
		),
		jen.Line(),
		jen.Return(jen.ID("m")),
	)

	return []jen.Code{
		jen.Func().IDf("build%sMapping", typ.Name.Singular()).Params().Params(
			jen.PointerTo().Qual(
				"github.com/blevesearch/bleve/mapping",
				"IndexMappingImpl",
			),
		).Block(
			block...,
		),
		jen.Line(),
	}
}
