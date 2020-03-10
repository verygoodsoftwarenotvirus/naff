package queriers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"
)

var (
	a = models.DataType{
		Name: wordsmith.FromSingularPascalCase("Grandparent"),
		Fields: []models.DataField{
			{
				Name:                  wordsmith.FromSingularPascalCase("GrandparentName"),
				Type:                  "string",
				ValidForCreationInput: true,
			},
		},
	}
	b = models.DataType{
		Name:            wordsmith.FromSingularPascalCase("Parent"),
		BelongsToStruct: a.Name,
		Fields: []models.DataField{
			{
				Name:                  wordsmith.FromSingularPascalCase("ParentName"),
				Type:                  "string",
				ValidForCreationInput: true,
			},
		},
	}
	c = models.DataType{
		Name:            wordsmith.FromSingularPascalCase("Child"),
		BelongsToStruct: b.Name,
		Fields: []models.DataField{
			{
				Name:                  wordsmith.FromSingularPascalCase("ChildName"),
				Type:                  "string",
				ValidForCreationInput: true,
			},
		},
	}
)

func Test_buildTestDBGetSomething(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetSomething(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetSomething(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetSomething(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBBuildGetSomethingCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetSomethingCountQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetSomethingCountQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetSomethingCountQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBGetSomethingCount(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetSomethingCount(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetSomethingCount(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetSomethingCount(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBBuildGetAllSomethingCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetAllSomethingCountQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetAllSomethingCountQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetAllSomethingCountQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBGetAllSomethingCount(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetAllSomethingCount(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetAllSomethingCount(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetAllSomethingCount(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBGetListOfSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetListOfSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetListOfSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetListOfSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBGetListOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetListOfSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetListOfSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetListOfSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBGetAllSomethingForSomethingElseFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetAllSomethingForSomethingElseFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetAllSomethingForSomethingElseFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBGetAllSomethingForSomethingElseFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBCreateSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBCreateSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBCreateSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBCreateSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBCreateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBCreateSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBCreateSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBCreateSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestBuildUpdateSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestBuildUpdateSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestBuildUpdateSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestBuildUpdateSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBUpdateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBUpdateSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBUpdateSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBUpdateSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBArchiveSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBArchiveSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBArchiveSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBArchiveSomethingQueryFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBArchiveSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBArchiveSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBArchiveSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBArchiveSomethingFuncDecl(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTestDBBuildGetSomethingQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres high dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetSomethingQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), c)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres one dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetSomethingQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), b)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})

	T.Run("postgres no dependency", func(t *testing.T) {
		proj := &models.Project{
			DataTypes: []models.DataType{a, b, c},
		}

		f := jen.NewFile("farts")
		lines := buildTestDBBuildGetSomethingQuery(proj, wordsmith.FromSingularPascalCase("Postgres"), a)
		f.Add(lines...)

		var b bytes.Buffer
		assert.NoError(t, f.Render(&b))

		actual := b.String()
		expected := `package farts

`

		assert.Equal(t, expected, actual)
	})
}
