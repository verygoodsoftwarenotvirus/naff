package client

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"
)

func Test_buildParamsForMethodThatHandlesAnInstanceOfADataType(T *testing.T) {
	T.Parallel()

	T.Run("normal operation with dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("CherryName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, cherry)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, appleID, bananaID, cherryID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation with fewer dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("BananaName"),
					Type: "string",
				},
			},
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, banana)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, appleID, bananaID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})

	T.Run("normal operation without dependencies", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
			Fields: []models.DataField{
				{
					Name: wordsmith.FromSingularPascalCase("AppleName"),
					Type: "string",
				},
			},
			BelongsToNobody: true,
		}
		proj := &models.Project{
			DataTypes: []models.DataType{apple},
		}

		ret := jen.NewFile("farts")

		ret.Add(
			jen.Func().ID("doSomething").Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(proj, apple)...).Block(),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"context"
)

func doSomething(ctx context.Context, appleID uint64) {}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
