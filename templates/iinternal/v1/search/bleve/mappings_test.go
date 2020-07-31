package bleve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mappingsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mappingsDotGo(proj)

		expected := `
package example

import (
	bleve "github.com/blevesearch/bleve"
	en "github.com/blevesearch/bleve/analysis/lang/en"
	mapping "github.com/blevesearch/bleve/mapping"
)

func buildItemMapping() *mapping.IndexMappingImpl {
	m := mapping.NewIndexMapping()

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	itemMapping := bleve.NewDocumentMapping()
	itemMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	itemMapping.AddFieldMappingsAt("details", englishTextFieldMapping)
	itemMapping.AddFieldMappingsAt("belongsToUser", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("item", itemMapping)

	return m
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingMapper(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingMapper(typ)

		expected := `
package example

import (
	bleve "github.com/blevesearch/bleve"
	en "github.com/blevesearch/bleve/analysis/lang/en"
	mapping "github.com/blevesearch/bleve/mapping"
)

func buildItemMapping() *mapping.IndexMappingImpl {
	m := mapping.NewIndexMapping()

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	itemMapping := bleve.NewDocumentMapping()
	itemMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	itemMapping.AddFieldMappingsAt("details", englishTextFieldMapping)
	itemMapping.AddFieldMappingsAt("belongsToUser", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("item", itemMapping)

	return m
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
