package integration

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func Test_buildRequisiteCreationCode(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
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
		proj := &models.Project{DataTypes: []models.DataType{apple, banana, cherry}}

		ret := jen.NewFile("farts")
		ret.Add(jen.Func().ID("doSomething").Params().Block(
			buildRequisiteCreationCode(proj, cherry)...,
		))

		var b bytes.Buffer
		require.NoError(t, ret.Render(&b))

		expected := `package farts

import (
	gofakeit "github.com/brianvoe/gofakeit"
	v1 "models/v1"
)

func doSomething() {
	// Create apple
	exampleApple := &v1.Apple{
		AppleName: gofakeit.Word(),
	}

	createdApple, err := todoClient.CreateApple(ctx, &v1.AppleCreationInput{
		AppleName: exampleApple.AppleName,
	})
	checkValueAndError(t, createdApple, err)

	// Create banana
	exampleBanana := &v1.Banana{
		BananaName: gofakeit.Word(),
	}

	createdBanana, err := todoClient.CreateBanana(ctx, &v1.BananaCreationInput{
		BananaName: exampleBanana.BananaName,
	}, createdApple.ID)
	checkValueAndError(t, createdBanana, err)

	// Create cherry
	exampleCherry := &v1.Cherry{
		CherryName: gofakeit.Word(),
	}

	createdCherry, err := todoClient.CreateCherry(ctx, &v1.CherryCreationInput{
		CherryName: exampleCherry.CherryName,
	}, createdApple.ID, createdBanana.ID)
	checkValueAndError(t, createdCherry, err)

}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}

//func Test_buildTestCreating(T *testing.T) {
//	T.Parallel()
//
//	T.Run("normal operation", func(t *testing.T) {
//		apple := models.DataType{
//			Name: wordsmith.FromSingularPascalCase("Apple"),
//			Fields: []models.DataField{
//				{
//					Name: wordsmith.FromSingularPascalCase("AppleName"),
//					Type: "string",
//				},
//			},
//		}
//		banana := models.DataType{
//			Name:            wordsmith.FromSingularPascalCase("Banana"),
//			BelongsToStruct: apple.Name,
//			Fields: []models.DataField{
//				{
//					Name: wordsmith.FromSingularPascalCase("BananaName"),
//					Type: "string",
//				},
//			},
//		}
//		cherry := models.DataType{
//			Name:            wordsmith.FromSingularPascalCase("Cherry"),
//			BelongsToStruct: banana.Name,
//			Fields: []models.DataField{
//				{
//					Name: wordsmith.FromSingularPascalCase("CherryName"),
//					Type: "string",
//				},
//			},
//		}
//
//		proj := &models.Project{
//			DataTypes: []models.DataType{apple, banana, cherry},
//		}
//
//		ret := jen.NewFile("farts")
//
//		ret.Add(
//			jen.Func().ID("doSomething").Params().Block(
//				buildTestCreating(proj, cherry)...,
//			),
//		)
//
//		var b bytes.Buffer
//		err := ret.Render(&b)
//		require.NoError(t, err)
//
//		expected := `package farts
//
//`
//		actual := b.String()
//
//		assert.Equal(t, expected, actual)
//	})
//}

func Test_buildCreationArguments(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
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
			jen.Func().ID("doSomething").Params().Block(
				buildCreationArguments(proj, "expected", cherry)...,
			),
		)

		var b bytes.Buffer
		err := ret.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import ()

func doSomething() {
	expectedApple.ID
	expectedBanana.ID
	expectedCherry.ID
}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
