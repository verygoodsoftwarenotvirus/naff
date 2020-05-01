package models

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"github.com/stretchr/testify/assert"
)

func TestParseModels(T *testing.T) {
	T.Parallel()

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
						Pos:                   token.Pos(39),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						Pos:                   token.Pos(52),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
				},
				BelongsToUser: true,
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

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
						Pos:                   token.Pos(58),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						Pos:                   token.Pos(71),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
				},
				BelongsToUser: true,
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

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
						Pos:                   token.Pos(39),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						Pos:                   token.Pos(53),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
				},
				BelongsToUser: true,
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

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
						Pos:                   token.Pos(39),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   false,
						Pos:                   token.Pos(73),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
				},
				BelongsToUser: true,
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

		assert.Equal(t, expectedDataTypes, actualDataTypes)
		assert.Equal(t, expectedImports, actualImports)
	})

	T.Run("with meta field indicating belonging to another object", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Owner struct {
	FirstName string
}

type Item struct{
	Name string
	Details string
	_META_ uintptr ` + "`" + `belongs_to:"Owner"` + "`" + `
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
				Name: wordsmith.FromSingularPascalCase("Owner"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("FirstName"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(41),
					},
				},
				BelongsToUser: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(80),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(93),
					},
				},
				BelongsToUser:   false,
				BelongsToStruct: wordsmith.FromSingularPascalCase("Owner"),
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/owners", exampleOutputPath),
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

		assert.Equal(t, expectedImports, actualImports)
		assert.Equal(t, len(expectedDataTypes), len(actualDataTypes))
		assert.Equal(t, expectedDataTypes, actualDataTypes)
	})

	T.Run("with meta field indicating belonging to another object and a user", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Owner struct {
	FirstName string
}

type Item struct{
	Name string
	Details string
	_META_ uintptr ` + "`" + `belongs_to:"Owner,User"` + "`" + `
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
				Name: wordsmith.FromSingularPascalCase("Owner"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("FirstName"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(41),
					},
				},
				BelongsToUser: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(80),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(93),
					},
				},
				BelongsToUser:   true,
				BelongsToStruct: wordsmith.FromSingularPascalCase("Owner"),
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/owners", exampleOutputPath),
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

		assert.Equal(t, expectedImports, actualImports)
		assert.Equal(t, len(expectedDataTypes), len(actualDataTypes))
		assert.Equal(t, expectedDataTypes, actualDataTypes)
	})

	T.Run("with meta field indicating belonging to another object and a user, restricted to that user", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Owner struct {
	FirstName string
}

type Item struct{
	Name string
	Details string
	_META_ uintptr ` + "`" + `belongs_to:"Owner,User" restricted_to_user:"true"` + "`" + `
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
				Name: wordsmith.FromSingularPascalCase("Owner"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("FirstName"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(41),
					},
				},
				BelongsToUser: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(80),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						UnderlyingType:        GetTypeForTypeName("string"),
						Pos:                   token.Pos(93),
					},
				},
				BelongsToUser:    true,
				RestrictedToUser: true,
				BelongsToStruct:  wordsmith.FromSingularPascalCase("Owner"),
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/owners", exampleOutputPath),
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

		assert.Equal(t, expectedImports, actualImports)
		assert.Equal(t, len(expectedDataTypes), len(actualDataTypes))
		assert.Equal(t, expectedDataTypes, actualDataTypes)
	})

	T.Run("with meta field indicating belonging to nothing", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name string
	Details string
	_META_ uintptr ` + "`" + `belongs_to:"__nobody__"` + "`" + `
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
						Pos:                   token.Pos(39),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
						Pos:                   token.Pos(52),
						UnderlyingType:        GetTypeForTypeName("string"),
					},
				},
				BelongsToUser:   false,
				BelongsToNobody: true,
			},
		}
		expectedImports := []string{
			fmt.Sprintf("%s/services/v1/items", exampleOutputPath),
		}

		actualDataTypes, actualImports, err := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.NoError(t, err)

		assert.Equal(t, expectedImports, actualImports)
		assert.Equal(t, len(expectedDataTypes), len(actualDataTypes))
		assert.Equal(t, expectedDataTypes, actualDataTypes)
	})

	T.Run("with invalid ownership arrangement", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name string
	Details string
	_META_ uintptr ` + "`" + `belongs_to:"Owner"` + "`" + `
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, _, actualErr := parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
		assert.Error(t, actualErr)
	})
}

func TestProject_containsCyclicOwnerships(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p := &Project{
			DataTypes: []DataType{
				{
					Name: wordsmith.FromSingularPascalCase("A"),
				},
				{
					Name:            wordsmith.FromSingularPascalCase("B"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("A"),
				},
				{
					Name:            wordsmith.FromSingularPascalCase("C"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("B"),
				},
			},
		}

		actual := p.containsCyclicOwnerships()

		assert.False(t, actual)
	})

	T.Run("with violation", func(t *testing.T) {
		p := &Project{
			DataTypes: []DataType{
				{
					Name:            wordsmith.FromSingularPascalCase("A"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("C"),
				},
				{
					Name:            wordsmith.FromSingularPascalCase("B"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("A"),
				},
				{
					Name:            wordsmith.FromSingularPascalCase("C"),
					BelongsToStruct: wordsmith.FromSingularPascalCase("B"),
				},
			},
		}

		actual := p.containsCyclicOwnerships()

		assert.True(t, actual)
	})
}

func TestProject_FindOwnerTypeChain(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		apple := DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &Project{
			DataTypes: []DataType{apple, banana, cherry},
		}

		expected := []DataType{apple, banana}
		actual := proj.FindOwnerTypeChain(cherry)

		assert.Equal(t, expected, actual)
	})
}
