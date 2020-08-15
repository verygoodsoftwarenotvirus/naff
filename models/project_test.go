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

func buildExampleTodoListProject() *Project {
	return &Project{
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/example",
		DataTypes: []DataType{
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
						DefaultValue:          "''",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToUser:    true,
				RestrictedToUser: true,
				SearchEnabled:    true,
			},
		},
	}
}

func buildExampleForumsListProject() *Project {
	p := &Project{
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/naff/example_output",
		Name:       wordsmith.FromSingularPascalCase("Discussion"),
		DataTypes: []DataType{
			{
				Name: wordsmith.FromSingularPascalCase("Forum"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				IsEnumeration: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Subforum"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
			},
			{
				Name: wordsmith.FromSingularPascalCase("Thread"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Title"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct:  wordsmith.FromSingularPascalCase("Subforum"),
				BelongsToUser:    true,
				RestrictedToUser: false,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Post"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Content"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct:  wordsmith.FromSingularPascalCase("Thread"),
				BelongsToUser:    true,
				RestrictedToUser: false,
			},
			{
				Name: wordsmith.FromSingularPascalCase("ReactionIcon"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Icon"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				IsEnumeration: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("PostReaction"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("PostReactionIcon"),
						Type:                  "uint64",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct:  wordsmith.FromSingularPascalCase("Post"),
				BelongsToUser:    true,
				RestrictedToUser: false,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Notification"),
				Fields: []DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Text"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToUser:    true,
				RestrictedToUser: true,
			},
		},
	}

	p.EnableDatabase(Postgres)
	p.EnableDatabase(Sqlite)
	p.EnableDatabase(MariaDB)

	return p
}

func TestProject_Validate(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			enabledDatabases: validDatabaseMap,
			DataTypes: []DataType{
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
							DefaultValue:          "''",
							Pointer:               false,
							ValidForCreationInput: true,
							ValidForUpdateInput:   true,
						},
					},
					BelongsToUser:    true,
					RestrictedToUser: true,
					SearchEnabled:    true,
				},
			},
		}

		p.Validate()
	})

	T.Run("with no databases", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()

		defer func() {
			if r := recover(); r == nil {
				t.Error("function didn't panic when expected")
			}
		}()

		p.Validate()
	})

	T.Run("with no types", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			enabledDatabases: validDatabaseMap,
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("function didn't panic when expected")
			}
		}()

		p.Validate()
	})

	T.Run("with no fields in type", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			enabledDatabases: validDatabaseMap,
			DataTypes: []DataType{
				{
					Name:             wordsmith.FromSingularPascalCase("Item"),
					BelongsToUser:    true,
					RestrictedToUser: true,
					SearchEnabled:    true,
				},
			},
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("function didn't panic when expected")
			}
		}()

		p.Validate()
	})

	T.Run("with cyclic arrangement", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			enabledDatabases: validDatabaseMap,
			DataTypes: []DataType{
				{
					Name: wordsmith.FromSingularPascalCase("A"),
					Fields: []DataField{
						{
							Name:                  wordsmith.FromSingularPascalCase("Name"),
							Type:                  "string",
							Pointer:               false,
							ValidForCreationInput: true,
							ValidForUpdateInput:   true,
						},
					},
					BelongsToStruct: wordsmith.FromSingularPascalCase("C"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("B"),
					Fields: []DataField{
						{
							Name:                  wordsmith.FromSingularPascalCase("Name"),
							Type:                  "string",
							Pointer:               false,
							ValidForCreationInput: true,
							ValidForUpdateInput:   true,
						},
					},
					BelongsToStruct: wordsmith.FromSingularPascalCase("A"),
				},
				{
					Name: wordsmith.FromSingularPascalCase("C"),
					Fields: []DataField{
						{
							Name:                  wordsmith.FromSingularPascalCase("Name"),
							Type:                  "string",
							Pointer:               false,
							ValidForCreationInput: true,
							ValidForUpdateInput:   true,
						},
					},
					BelongsToStruct: wordsmith.FromSingularPascalCase("B"),
				},
			},
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("function didn't panic when expected")
			}
		}()

		p.Validate()
	})
}

func TestProject_SearchEnabled(T *testing.T) {
	T.Parallel()

	T.Run("expecting true", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes[0].SearchEnabled = true

		assert.True(t, p.SearchEnabled())
	})

	T.Run("expecting false", func(t *testing.T) {
		t.Parallel()

		p := buildExampleTodoListProject()
		p.DataTypes[0].SearchEnabled = false

		assert.False(t, p.SearchEnabled())
	})
}

func TestProject_ParseModels(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			sourcePackage:    "gitlab.com/verygoodsoftwarenotvirus/naff/example_models/todo",
			enabledDatabases: validDatabaseMap,
		}

		assert.NoError(t, p.ParseModels())
	})

	T.Run("with invalid path", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			sourcePackage: "gitlab.com/verygoodsoftwarenotvirus/naff/does/not/exist/lol",
		}

		assert.Error(t, p.ParseModels())
	})

	T.Run("with invalid ownership", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			sourcePackage:    "gitlab.com/verygoodsoftwarenotvirus/naff/example_models/invalid_ownership",
			enabledDatabases: validDatabaseMap,
		}

		assert.Error(t, p.ParseModels())
	})
}

func TestProject_FindOwnerTypeChainWithoutReversing(T *testing.T) {
	T.Parallel()

	T.Run("basic usecase", func(t *testing.T) {
		t.Parallel()

		forumType := DataType{
			Name: wordsmith.FromSingularPascalCase("Forum"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Name"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			IsEnumeration: true,
		}
		subForumType := DataType{
			Name: wordsmith.FromSingularPascalCase("Subforum"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Name"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
		}
		threadType := DataType{
			Name: wordsmith.FromSingularPascalCase("Thread"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Title"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Subforum"),
			BelongsToUser:    true,
			RestrictedToUser: false,
		}
		postType := DataType{
			Name: wordsmith.FromSingularPascalCase("Post"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Content"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thread"),
			BelongsToUser:    true,
			RestrictedToUser: false,
		}

		p := &Project{
			enabledDatabases: validDatabaseMap,
			DataTypes: []DataType{
				forumType,
				subForumType,
				threadType,
				postType,
			},
		}

		expected := []DataType{
			threadType,
			subForumType,
			forumType,
		}
		actual := p.FindOwnerTypeChainWithoutReversing(postType)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_FindType(T *testing.T) {
	T.Parallel()

	T.Run("basic usecase", func(t *testing.T) {
		t.Parallel()

		forumType := DataType{
			Name: wordsmith.FromSingularPascalCase("Forum"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Name"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			IsEnumeration: true,
		}
		subForumType := DataType{
			Name: wordsmith.FromSingularPascalCase("Subforum"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Name"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
		}
		threadType := DataType{
			Name: wordsmith.FromSingularPascalCase("Thread"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Title"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Subforum"),
			BelongsToUser:    true,
			RestrictedToUser: false,
		}
		postType := DataType{
			Name: wordsmith.FromSingularPascalCase("Post"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Content"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thread"),
			BelongsToUser:    true,
			RestrictedToUser: false,
		}

		p := &Project{
			enabledDatabases: validDatabaseMap,
			DataTypes: []DataType{
				forumType,
				subForumType,
				threadType,
				postType,
			},
		}

		expected := &threadType
		actual := p.FindType("Thread")

		assert.Equal(t, expected, actual)
	})

	T.Run("missing type", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		assert.Nil(t, p.FindType("whatever"))
	})
}

func TestProject_FindDependentsOfType(T *testing.T) {
	T.Parallel()

	T.Run("basic usecase", func(t *testing.T) {
		t.Parallel()

		forumType := DataType{
			Name: wordsmith.FromSingularPascalCase("Forum"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Name"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			IsEnumeration: true,
		}
		subForumType := DataType{
			Name: wordsmith.FromSingularPascalCase("Subforum"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Name"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
		}
		threadType := DataType{
			Name: wordsmith.FromSingularPascalCase("Thread"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Title"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Subforum"),
			BelongsToUser:    true,
			RestrictedToUser: false,
		}
		postType := DataType{
			Name: wordsmith.FromSingularPascalCase("Post"),
			Fields: []DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Content"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thread"),
			BelongsToUser:    true,
			RestrictedToUser: false,
		}

		p := &Project{
			enabledDatabases: validDatabaseMap,
			DataTypes: []DataType{
				forumType,
				subForumType,
				threadType,
				postType,
			},
		}

		expected := []DataType{subForumType}
		actual := p.FindDependentsOfType(forumType)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ensureNoNilFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		assert.Nil(t, p.enabledDatabases)
		p.ensureNoNilFields()
		assert.NotNil(t, p.enabledDatabases)
	})
}

func TestProject_EnableDatabase(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		assert.False(t, p.DatabaseIsEnabled(Postgres))
		p.EnableDatabase(Postgres)
		assert.True(t, p.DatabaseIsEnabled(Postgres))
	})

	T.Run("with invalid database", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		defer func() {
			if r := recover(); r == nil {
				t.Error("function didn't panic as expected")
			}
		}()

		p.EnableDatabase("fart")
	})
}

func TestProject_DisableDatabase(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		p.EnableDatabase(Postgres)
		assert.True(t, p.DatabaseIsEnabled(Postgres))

		p.DisableDatabase(Postgres)
		assert.False(t, p.DatabaseIsEnabled(Postgres))
	})

	T.Run("with invalid database", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		defer func() {
			if r := recover(); r == nil {
				t.Error("function didn't panic as expected")
			}
		}()

		p.DisableDatabase("fart")
	})
}

func TestProject_DatabaseIsEnabled(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		assert.False(t, p.DatabaseIsEnabled(Postgres))
		p.EnableDatabase(Postgres)
		assert.True(t, p.DatabaseIsEnabled(Postgres))
	})

	T.Run("with invalid database", func(t *testing.T) {
		t.Parallel()

		p := &Project{}

		defer func() {
			if r := recover(); r == nil {
				t.Error("function didn't panic as expected")
			}
		}()

		p.DatabaseIsEnabled("fart")
	})
}

func TestProject_EnabledDatabases(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		p := &Project{
			enabledDatabases: validDatabaseMap,
		}

		actual := p.EnabledDatabases()

		for k := range validDatabaseMap {
			assert.Contains(t, actual, string(k))
		}
	})
}

func TestParseModels(T *testing.T) {
	T.Parallel()

	T.Run("basic usecase", func(t *testing.T) {
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

	T.Run("panics with uintptr type in wrong place", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name string
	Details uintptr
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic as expected")
			}
		}()

		_, _, _ = parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
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

	T.Run("with meta field indicating belonging to too many owners", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Owner struct {
	FirstName string
}

type Loaner struct {
	FirstName string
}

type Groaner struct {
	FirstName string
}

type Item struct{
	Name string
	Details string

	_META_ uintptr ` + "`" + `belongs_to:"Owner,Loaner,Groaner"` + "`" + `
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		_, _, _ = parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
	})

	T.Run("with meta field indicating belonging to two owners, but no user", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Owner struct {
	FirstName string
}

type Loaner struct {
	FirstName string
}

type Item struct{
	Name string
	Details string

	_META_ uintptr ` + "`" + `belongs_to:"Owner,Loaner"` + "`" + `
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		_, _, _ = parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
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

	_META_ uintptr ` + "`" + `belongs_to:"User,Owner" restricted_to_user:"true"` + "`" + `
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
				BelongsToUser: false,
				IsEnumeration: true,
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

	T.Run("with search enabled", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name string
	Details string

	_META_ uintptr ` + "`" + `search_enabled:"true"` + "`" + `
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
				SearchEnabled: true,
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

	T.Run("with search enabled, but no string field", func(t *testing.T) {
		exampleOutputPath := "things/stuff"
		exampleCode := `
package whatever

type Item struct{
	Name uint64
	Details uint64

	_META_ uintptr ` + "`" + `search_enabled:"true"` + "`" + `
}
`
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", exampleCode, parser.AllErrors)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		_, _, _ = parseModels(exampleOutputPath, map[string]*ast.File{f.Name.String(): f})
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

func TestGetTypeForTypeName(T *testing.T) {
	T.Parallel()

	T.Run("non-nil expectations", func(t *testing.T) {
		t.Parallel()

		expectedToNotReturnNil := []string{
			"bool",
			"int",
			"int8",
			"int16",
			"int32",
			"int64",
			"uint",
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"uintptr",
			"float32",
			"float64",
			"string",
		}

		for _, input := range expectedToNotReturnNil {
			assert.NotNil(t, GetTypeForTypeName(input))
		}
	})

	T.Run("with expected panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic did not occur")
			}
		}()

		GetTypeForTypeName("fart")
	})
}
