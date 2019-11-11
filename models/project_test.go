package models

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

func TestParseModels(T *testing.T) {

	T.Run("happy path", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name string
	Details string
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		expectedDataTypes := []DataType{
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})

		assert.Equal(t, expectedDataTypes, actualDataTypes)
		assert.Equal(t, expectedImports, actualImports)
	})

	T.Run("ignores invalid declarations", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Thing string

type Item struct{
	Name string
	Details string
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		expectedDataTypes := []DataType{
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})

		assert.Equal(t, expectedDataTypes, actualDataTypes)
		assert.Equal(t, expectedImports, actualImports)
	})

	T.Run("with pointers", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name *string
	Details string
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		expectedDataTypes := []DataType{
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               true,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})

		assert.Equal(t, expectedDataTypes, actualDataTypes)
		assert.Equal(t, expectedImports, actualImports)
	})

	T.Run("with creatable and editable bangflags", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name *string ` + "`" + `naff:"!creatable"` + "`" + `
	Details string ` + "`" + `naff:"!editable"` + "`" + `
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		expectedDataTypes := []DataType{
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               true,
						ValidForCreationInput: false,
						ValidForUpdateInput:   true,
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   false,
					},
				},
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})

		assert.Equal(t, expectedDataTypes, actualDataTypes)
		assert.Equal(t, expectedImports, actualImports)
	})
}
